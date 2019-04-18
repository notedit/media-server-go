package mediaserver

import (
	"strconv"

	"github.com/notedit/sdp"
)

type outmediatrack struct {
	track  *OutgoingStreamTrack
	stream *OutgoingStream
}

type inmediatrack struct {
	track  *IncomingStreamTrack
	stream *IncomingStream
}

type rtctransceiver struct {
	mid    string
	media  string
	remote *inmediatrack
	local  *outmediatrack
}

// SDPManagerUnified unified plan manager
type SDPManagerUnified struct {
	state               string
	endpoint            *Endpoint
	transport           *Transport
	capabilities        map[string]*sdp.Capability
	pending             []*outmediatrack
	removed             []*OutgoingStreamTrack
	transceivers        []*rtctransceiver
	remoteInfo          *sdp.SDPInfo
	localInfo           *sdp.SDPInfo
	renegotiationNeeded bool
	OnRenegotiation     RenegotiationCallback
}

// NewSDPManagerUnified create SDPUnified manager
func NewSDPManagerUnified(endpoint *Endpoint, capabilities map[string]*sdp.Capability) *SDPManagerUnified {
	sdpManager := &SDPManagerUnified{}
	sdpManager.endpoint = endpoint
	sdpManager.capabilities = capabilities
	sdpManager.state = "initial"

	sdpManager.pending = []*outmediatrack{}
	sdpManager.removed = []*OutgoingStreamTrack{}
	sdpManager.transceivers = []*rtctransceiver{}

	return sdpManager
}

func (s *SDPManagerUnified) GetState() string {
	return s.state
}

func (s *SDPManagerUnified) GetTransport() *Transport {
	return s.transport
}

func (s *SDPManagerUnified) CreateLocalDescription() (*sdp.SDPInfo, error) {

	if s.localInfo == nil {
		ice := sdp.ICEInfoGenerate(true)
		dtls := sdp.NewDTLSInfo(sdp.SETUPACTPASS, "sha-256", s.endpoint.GetDTLSFingerprint())
		candidates := s.endpoint.GetLocalCandidates()

		s.localInfo = sdp.Create(ice, dtls, candidates, make(map[string]*sdp.Capability))

		for media, capability := range s.capabilities {
			mid := strconv.Itoa(len(s.transceivers))
			s.transceivers = append(s.transceivers, &rtctransceiver{
				mid:    mid,
				media:  media,
				remote: &inmediatrack{},
				local:  &outmediatrack{},
			})
			mediaInfo := sdp.MediaInfoCreate(media, capability)
			mediaInfo.SetID(mid)
			s.localInfo.AddMedia(mediaInfo)
		}
	}

	remaintracks := []*OutgoingStreamTrack{}
	for _, mediatrack := range s.removed {
		needdelete := false
		for _, transceiver := range s.transceivers {
			if transceiver.local.track == mediatrack {
				transceiver.local.track = nil
				transceiver.local.stream = nil
				needdelete = true
			}
		}
		if !needdelete {
			remaintracks = append(remaintracks, mediatrack)
		}
	}
	s.removed = remaintracks

	if s.state == "initial" || s.state == "stable" {
		for _, pending := range s.pending {
			mid := strconv.Itoa(len(s.transceivers))
			media := pending.track.GetMedia()
			s.transceivers = append(s.transceivers, &rtctransceiver{
				mid:    mid,
				media:  media,
				remote: &inmediatrack{},
				local: &outmediatrack{
					track:  pending.track,
					stream: pending.stream,
				},
			})
		}
		s.pending = make([]*outmediatrack, 0)
	}

	s.localInfo.RemoveAllStreams()

	for _, transceiver := range s.transceivers {
		if transceiver.local.track != nil {
			mediaInfo := s.localInfo.GetMediaByID(transceiver.mid)
			if mediaInfo == nil {
				mediaInfo = s.localInfo.GetMedia(transceiver.media).Clone()
				mediaInfo.SetDirection(sdp.SENDRECV)
				mediaInfo.SetID(transceiver.mid)
				s.localInfo.AddMedia(mediaInfo)
			} else if mediaInfo.GetDirection() != sdp.SENDRECV {
				mediaInfo.SetDirection(sdp.SENDONLY)
			}

			streamInfo := s.localInfo.GetStream(transceiver.local.stream.GetID())
			if streamInfo == nil {
				streamInfo = sdp.NewStreamInfo(transceiver.local.stream.GetID())
			}
			s.localInfo.AddStream(streamInfo)
			trackInfo := transceiver.local.track.GetTrackInfo()
			trackInfo.SetMediaID(transceiver.mid)
			streamInfo.AddTrack(trackInfo)
		} else {
			mediaInfo := s.localInfo.GetMediaByID(transceiver.mid)
			if mediaInfo.GetDirection() == sdp.SENDRECV {
				mediaInfo.SetDirection(sdp.RECVONLY)
			} else if mediaInfo.GetDirection() == sdp.SENDONLY {
				mediaInfo.SetDirection(sdp.INACTIVE)
			}
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

	if len(s.pending) > 0 && s.OnRenegotiation != nil {
		s.OnRenegotiation(s.transport)
	}

	return s.localInfo, nil
}

func (s *SDPManagerUnified) ProcessRemoteDescription(sdpStr string) (*sdp.SDPInfo, error) {

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
			s.pending = append(s.pending, &outmediatrack{
				track:  track,
				stream: stream,
			})

			track.OnStop(func() {
				s.removed = append(s.removed, track)
				s.renegotiate()
			})

			s.renegotiate()
		})

	}

	if s.state == "local-offer" {
		s.localInfo = s.remoteInfo.Answer(s.transport.GetLocalICEInfo(),
			s.transport.GetLocalDTLSInfo(),
			s.endpoint.GetLocalCandidates(),
			s.capabilities)

	}

	i := 0
	for _, mediaInfo := range s.remoteInfo.GetMedias() {
		mid := mediaInfo.GetID()

		streamInfo := s.remoteInfo.GetStreamByMediaID(mid)
		trackInfo := s.remoteInfo.GetTrackByMediaID(mid)

		var stream *IncomingStream
		var track *IncomingStreamTrack
		if streamInfo != nil {
			stream = s.transport.GetIncomingStream(streamInfo.GetID())
		}

		if stream != nil && trackInfo != nil {
			track = stream.GetTrack(trackInfo.GetID())
		}

		var transceiver *rtctransceiver
		if i+1 <= len(s.transceivers) {
			transceiver = s.transceivers[i]
		}
		if transceiver == nil {
			transceiver = &rtctransceiver{
				mid:    mid,
				remote: &inmediatrack{},
				local:  &outmediatrack{},
			}
			s.transceivers = append(s.transceivers, transceiver)
		}
		i += 1
		switch mediaInfo.GetDirection() {
		case sdp.SENDRECV, sdp.SENDONLY:
			if transceiver.remote.track != nil && transceiver.remote.track != track {
				transceiver.remote.track.Stop()
			}
			if stream == nil {
				stream = s.transport.CreateIncomingStream(streamInfo)
			} else if track == nil && trackInfo != nil {
				track = stream.CreateTrack(trackInfo)
			}
			transceiver.remote.track = track
			break
		case sdp.RECVONLY, sdp.INACTIVE:
			if track != nil {
				track.Stop()
			}
			if transceiver.remote.track != nil {
				transceiver.remote.track = nil
			}

			break
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

func (s *SDPManagerUnified) renegotiate() {
	if s.state == "initial" || s.state == "stable" {
		if s.OnRenegotiation != nil {
			s.OnRenegotiation(s.transport)
		}
	}
}
