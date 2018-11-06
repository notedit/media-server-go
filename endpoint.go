package mediaserver

import (
	"./sdp"

	"github.com/chuckpreslar/emission"
)

type Endpoint struct {
	ip          string
	bundle      RTPBundleTransport
	transports  map[string]*Transport
	candidate   *sdp.CandidateInfo
	fingerprint string
	*emission.Emitter
}

// NewEndpoint create a endpoint
func NewEndpoint(ip string) *Endpoint {
	endpoint := new(Endpoint)
	endpoint.bundle = NewRTPBundleTransport()
	endpoint.bundle.Init()
	endpoint.transports = make(map[string]*Transport)
	endpoint.fingerprint = MediaServerGetFingerprint().ToString()
	endpoint.candidate = sdp.NewCandidateInfo("1", 1, "UDP", 33554431, ip, endpoint.bundle.GetLocalPort(), "host", "", 0)
	endpoint.Emitter = emission.NewEmitter()
	return endpoint
}

func (e *Endpoint) CreateTransport(remoteIce *sdp.ICEInfo, remoteDtls *sdp.DTLSInfo, remoteCandidates []*sdp.CandidateInfo,
	localIce *sdp.ICEInfo, localDtls *sdp.DTLSInfo, localCandidates []*sdp.CandidateInfo, disableSTUNKeepAlive bool) *Transport {

	if localIce == nil {
		localIce = sdp.GenerateIce(true)
	}

	if localDtls == nil {
		localDtls = sdp.NewDTLSInfo(remoteDtls.GetSetup().Reverse(), "sha-256", e.fingerprint)
	}

	if localCandidates == nil {
		localCandidates = []*sdp.CandidateInfo{e.candidate}
	}

	remoteCandidatesClone := []*sdp.CandidateInfo{}
	for _, candidate := range remoteCandidates {
		remoteCandidatesClone = append(remoteCandidatesClone, candidate.Clone())
	}

	localCandidatesClone := []*sdp.CandidateInfo{}
	for _, candidate := range localCandidates {
		localCandidatesClone = append(localCandidatesClone, candidate.Clone())
	}

	transport := NewTransport(e.bundle, remoteIce.Clone(), remoteDtls.Clone(), remoteCandidatesClone,
		localIce.Clone(), localDtls.Clone(), localCandidatesClone, disableSTUNKeepAlive)

	e.transports[transport.username.ToString()] = transport

	transport.Once("stopped", func() {
		delete(e.transports, transport.username.ToString())
	})

	return transport
}

func (e *Endpoint) GetLocalCandidates() []*sdp.CandidateInfo {
	return []*sdp.CandidateInfo{e.candidate}
}

func (e *Endpoint) GetDTLSFingerprint() string {
	return e.fingerprint
}

// CreateOffer  create offer based on audio and video capability
func (e *Endpoint) CreateOffer(video *sdp.Capability, audio *sdp.Capability) *sdp.SDPInfo {

	dtls := sdp.NewDTLSInfo(sdp.SETUPACTPASS, "sha-256", e.fingerprint)

	ice := sdp.GenerateICEInfo(true)

	candidates := e.GetLocalCandidates()

	capabilities := make(map[string]*sdp.Capability)

	if video != nil {
		capabilities["video"] = video
	}

	if audio != nil {
		capabilities["audio"] = audio
	}

	return sdp.Create(ice, dtls, candidates, capabilities)
}

// Stop  stop this endpoint
func (e *Endpoint) Stop() {

	if e.bundle == nil {
		return
	}

	for _, transport := range e.transports {
		transport.Stop()
	}

	e.transports = nil

	e.bundle.End()

	e.bundle = nil
}
