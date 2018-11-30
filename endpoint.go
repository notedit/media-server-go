package mediaserver

import (
	"runtime"

	"github.com/chuckpreslar/emission"
	"github.com/notedit/media-server-go/sdp"
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
	endpoint := &Endpoint{}
	endpoint.bundle = NewRTPBundleTransport()
	endpoint.bundle.Init()
	endpoint.transports = make(map[string]*Transport)
	endpoint.fingerprint = MediaServerGetFingerprint().ToString()
	endpoint.candidate = sdp.NewCandidateInfo("1", 1, "UDP", 33554431, ip, endpoint.bundle.GetLocalPort(), "host", "", 0)
	endpoint.Emitter = emission.NewEmitter()
	runtime.SetFinalizer(endpoint, endpoint.deleteRTPBundleTransport)
	return endpoint
}

func NewEndpointWithPort(ip string, port int) *Endpoint {
	endpoint := &Endpoint{}
	endpoint.bundle = NewRTPBundleTransport()
	endpoint.bundle.Init(port)
	endpoint.transports = make(map[string]*Transport)
	endpoint.fingerprint = MediaServerGetFingerprint().ToString()
	endpoint.candidate = sdp.NewCandidateInfo("1", 1, "UDP", 33554431, ip, endpoint.bundle.GetLocalPort(), "host", "", 0)
	endpoint.Emitter = emission.NewEmitter()
	runtime.SetFinalizer(endpoint, endpoint.deleteRTPBundleTransport)
	return endpoint
}

func (e *Endpoint) CreateTransport(remoteSdp *sdp.SDPInfo, localSdp *sdp.SDPInfo, disableSTUNKeepAlive bool) *Transport {

	var localIce *sdp.ICEInfo
	var localDtls *sdp.DTLSInfo
	var localCandidates []*sdp.CandidateInfo

	if localSdp == nil {
		localIce = sdp.ICEInfoGenerate(true)
		localDtls = sdp.NewDTLSInfo(remoteSdp.GetDTLS().GetSetup().Reverse(), "sha-256", e.fingerprint)
		localCandidates = []*sdp.CandidateInfo{e.candidate}
	} else {
		localIce = localSdp.GetICE().Clone()
		localDtls = localSdp.GetDTLS().Clone()
		localCandidates = localSdp.GetCandidates()
	}

	remoteIce := remoteSdp.GetICE().Clone()
	remoteDtls := remoteSdp.GetDTLS().Clone()
	remoteCandidates := remoteSdp.GetCandidates()

	localIce.SetLite(true)
	localIce.SetEndOfCandidate(true)

	transport := NewTransport(e.bundle, remoteIce, remoteDtls, remoteCandidates,
		localIce, localDtls, localCandidates, disableSTUNKeepAlive)

	e.transports[transport.username.ToString()] = transport

	transport.Once("stopped", func() {
		delete(e.transports, transport.username.ToString())
	})

	return transport

}

func (e *Endpoint) CreateTransportWithRemote(sdpInfo *sdp.SDPInfo, disableSTUNKeepAlive bool) *Transport {

	localIce := sdp.ICEInfoGenerate(true)
	localDtls := sdp.NewDTLSInfo(sdpInfo.GetDTLS().GetSetup().Reverse(), "sha-256", e.fingerprint)
	localCandidates := []*sdp.CandidateInfo{e.candidate.Clone()}

	remoteCandidatesClone := []*sdp.CandidateInfo{}
	for _, candidate := range sdpInfo.GetCandidates() {
		remoteCandidatesClone = append(remoteCandidatesClone, candidate.Clone())
	}

	remoteIceClone := sdpInfo.GetICE().Clone()
	remoteDtlsClone := sdpInfo.GetDTLS().Clone()

	transport := NewTransport(e.bundle, remoteIceClone, remoteDtlsClone, remoteCandidatesClone, localIce, localDtls, localCandidates, disableSTUNKeepAlive)

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

	e.EmitSync("stopped")

	e.bundle.End()

	runtime.SetFinalizer(e, nil)
	DeleteRTPBundleTransport(e.bundle)

	e.bundle = nil
}

func (e *Endpoint) deleteRTPBundleTransport() {

	DeleteRTPBundleTransport(e.bundle)
}
