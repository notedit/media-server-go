package mediaserver

import (
	"github.com/notedit/sdp"
)

// SDPManagerPlanb  planb manager
type SDPManagerPlanb struct {
	state           string
	endpoint        *Endpoint
	transport       *Transport
	capabilities    map[string]*sdp.Capability
	remoteInfo      *sdp.SDPInfo
	localInfo       *sdp.SDPInfo
	OnRenegotiation RenegotiationCallback
}

func NewSDPManagerPlanb(endpoint *Endpoint, capabilities map[string]*sdp.Capability) *SDPManagerPlanb {
	sdpManager := &SDPManagerPlanb{}
	sdpManager.endpoint = endpoint
	sdpManager.capabilities = capabilities
	sdpManager.state = "initial"

	return sdpManager
}

func (s *SDPManagerPlanb) GetState() string {
	return s.state
}

func (s *SDPManagerPlanb) GetTransport() *Transport {
	return s.transport
}

func (s *SDPManagerPlanb) CreateLocalDescription() (*sdp.SDPInfo, error) {
	if s.localInfo == nil {
		ice := sdp.ICEInfoGenerate(true)
		dtls := sdp.NewDTLSInfo(sdp.SETUPACTPASS, "sha-256", s.endpoint.GetDTLSFingerprint())
		candidates := s.endpoint.GetLocalCandidates()
		s.localInfo = sdp.Create(ice, dtls, candidates, s.capabilities)
	}
	s.localInfo.RemoveAllStreams()
	if s.transport != nil {
		for _, stream := range s.transport.GetOutgoingStreams() {
			s.localInfo.AddStream(stream.GetStreamInfo())
		}
	}
	switch s.state {
	case "initial", "stable":
		s.state = "local-offer"
		break
	case "remote-offer":
		s.state = "stable"
		break
	}
	return s.localInfo, nil
}

func (s *SDPManagerPlanb) ProcessRemoteDescription(sdpStr string) (*sdp.SDPInfo, error) {
	info, err := sdp.Parse(sdpStr)
	if err != nil {
		return nil, err
	}
	s.remoteInfo = info

	if s.transport == nil {
		s.transport = s.endpoint.CreateTransport(s.remoteInfo, s.localInfo, false)
		if s.localInfo == nil {
			s.localInfo = s.remoteInfo.Answer(s.transport.GetLocalICEInfo(), s.transport.GetLocalDTLSInfo(),
				s.endpoint.GetLocalCandidates(), s.capabilities)

		}
		s.transport.SetLocalProperties(s.localInfo.GetAudioMedia(), s.localInfo.GetVideoMedia())
		s.transport.SetRemoteProperties(s.remoteInfo.GetAudioMedia(), s.remoteInfo.GetVideoMedia())

		s.transport.OnOutgoingTrack(func(track *OutgoingStreamTrack, stream *OutgoingStream) {
			track.OnStop(func() {
				s.renegotiate()
			})
			s.renegotiate()
		})

	}

	if s.state != "local-offer" {
		s.localInfo = s.remoteInfo.Answer(s.transport.GetLocalICEInfo(),
			s.transport.GetLocalDTLSInfo(),
			s.endpoint.GetLocalCandidates(),
			s.capabilities)
	}

	for _, stream := range s.transport.GetIncomingStreams() {
		streamInfo := s.remoteInfo.GetStream(stream.GetID())
		if streamInfo == nil {
			stream.Stop()
			continue
		}
		for _, track := range stream.GetTracks() {
			trackInfo := streamInfo.GetTrack(track.GetID())
			if trackInfo == nil {
				track.Stop()
			}
		}
	}

	for _, streamInfo := range s.remoteInfo.GetStreams() {
		stream := s.transport.GetIncomingStream(streamInfo.GetID())
		if stream == nil {
			s.transport.CreateIncomingStream(streamInfo)
			continue
		}
		for _, trackInfo := range streamInfo.GetTracks() {
			track := stream.GetTrack(trackInfo.GetID())
			if track == nil {
				stream.CreateTrack(trackInfo)
			}
		}
	}

	switch s.state {
	case "initial", "stable":
		s.state = "remote-offer"
		break
	case "local-offer":
		s.state = "stable"
		break
	}
	return s.remoteInfo, nil
}

func (s *SDPManagerPlanb) renegotiate() {
	if s.state == "initial" || s.state == "stable" {
		if s.OnRenegotiation != nil {
			s.OnRenegotiation(s.transport)
		}
	}
}
