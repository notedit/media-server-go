package mediaserver

import (
	"sync"

	native "github.com/notedit/media-server-go/wrapper"
	"github.com/notedit/sdp"
)

// Endpoint is an endpoint represent an UDP server socket.
// The endpoint will process STUN requests in order to be able to associate the remote ip:port with the registered transport and forward any further data comming from that transport.
// Being a server it is ICE-lite.
type Endpoint struct {
	ip              string
	bundle          native.RTPBundleTransport
	transports      map[string]*Transport
	candidate       *sdp.CandidateInfo
	mirroredStreams map[string]*IncomingStream
	mirroredTracks  map[string]*IncomingStreamTrack
	fingerprint     string
	sync.Mutex
}

// NewEndpoint create a new endpoint with given ip
func NewEndpoint(ip string) *Endpoint {
	endpoint := &Endpoint{}
	endpoint.bundle = native.NewRTPBundleTransport()
	endpoint.bundle.Init()
	endpoint.transports = make(map[string]*Transport)
	endpoint.fingerprint = native.MediaServerGetFingerprint().ToString()
	endpoint.mirroredStreams = make(map[string]*IncomingStream)
	endpoint.mirroredTracks = make(map[string]*IncomingStreamTrack)
	endpoint.candidate = sdp.NewCandidateInfo("1", 1, "UDP", 33554431, ip, endpoint.bundle.GetLocalPort(), "host", "", 0)
	return endpoint
}

// NewEndpointWithPort create a new endpint with given ip and port
func NewEndpointWithPort(ip string, port int) *Endpoint {
	endpoint := &Endpoint{}
	endpoint.bundle = native.NewRTPBundleTransport()
	endpoint.bundle.Init(port)
	endpoint.transports = make(map[string]*Transport)
	endpoint.fingerprint = native.MediaServerGetFingerprint().ToString()
	endpoint.candidate = sdp.NewCandidateInfo("1", 1, "UDP", 33554431, ip, endpoint.bundle.GetLocalPort(), "host", "", 0)
	return endpoint
}

//SetAffinity Set cpu affinity
func (e *Endpoint) SetAffinity(cpu int) {
	e.bundle.SetAffinity(cpu)
}

// CreateTransport create a new transport object and register it with the remote ICE username and password
// disableSTUNKeepAlive - Disable ICE/STUN keep alives, required for server to server transports, set this to false if you do not how to use it
func (e *Endpoint) CreateTransport(remoteSdp *sdp.SDPInfo, localSdp *sdp.SDPInfo, options ...bool) *Transport {

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

	disableSTUNKeepAlive := false

	if len(options) > 0 {
		disableSTUNKeepAlive = options[0]
	}

	transport := NewTransport(e.bundle, remoteIce, remoteDtls, remoteCandidates,
		localIce, localDtls, localCandidates, disableSTUNKeepAlive)

	e.Lock()
	e.transports[transport.username.ToString()] = transport
	e.Unlock()

	transport.OnStop(func() {
		e.Lock()
		delete(e.transports, transport.username.ToString())
		e.Unlock()
	})

	return transport
}

// GetLocalCandidates Get local ICE candidates for this endpoint. It will be shared by all the transport associated to this endpoint.
func (e *Endpoint) GetLocalCandidates() []*sdp.CandidateInfo {
	return []*sdp.CandidateInfo{e.candidate}
}

// GetDTLSFingerprint Get local DTLS fingerprint for this endpoint. It will be shared by all the transport associated to this endpoint
func (e *Endpoint) GetDTLSFingerprint() string {
	return e.fingerprint
}

// CreateOffer  create offer based on audio and video capability
// It generates a random ICE username and password and gets endpoint fingerprint
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



// CreateSDPManager Create new SDP manager, this object will manage the SDP O/A for you and produce a suitable trasnport.
func (e *Endpoint) CreateSDPManager(sdpSemantics string, capabilities map[string]*sdp.Capability) SDPManager {

	if sdpSemantics == "plan-b" {
		return NewSDPManagerPlanb(e, capabilities)
	} else if sdpSemantics == "unified-plan" {
		return NewSDPManagerUnified(e, capabilities)
	}
	return nil
}

// Stop stop the endpoint UDP server and terminate any associated transport
func (e *Endpoint) Stop() {

	if e.bundle == nil {
		return
	}

	for _, transport := range e.transports {
		transport.Stop()
	}

	e.transports = nil

	e.bundle.End()

	native.DeleteRTPBundleTransport(e.bundle)

}
