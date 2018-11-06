package mediaserver

import (
	"fmt"
	"strings"

	"./sdp"
)

type StreamerSession struct {
	local    bool
	port     int
	ip       string
	incoming *IncomingStreamTrack
	outgoing *OutgoingStreamTrack
	session  RTPSessionFacade
}

func NewStreamerSession(local bool, ip string, port int, media *sdp.MediaInfo) *StreamerSession {

	streamerSession := &StreamerSession{}
	var mediaType MediaFrameType = 0
	if strings.ToLower(media.GetType()) == "video" {
		mediaType = 1
	}
	session := NewRTPSessionFacade(mediaType)
	if local {
		session.SetLocalPort(port)
	} else {
		session.SetRemotePort(ip, port)
	}

	properties := NewProperties()
	if media != nil {
		num := 0
		for _, codec := range media.GetCodecs() {
			item := fmt.Sprintf("codecs.%d", num)
			properties.SetProperty(item+".codec", codec.GetCodec())
			properties.SetProperty(item+".pt", codec.GetType())
			if codec.HasRTX() {
				properties.SetProperty(item+".rtx", codec.GetRTX())
			}
			num = num + 1
		}
		properties.SetProperty("codecs.length", num)
	}

	session.Init(properties)
	streamerSession.session = session

	streamerSession.incoming = newIncomingStreamTrack(media.GetType(), media.GetType(), SessionToReceiver(session), []RTPIncomingSourceGroup{session.GetIncomingSourceGroup()})

	streamerSession.outgoing = newOutgoingStreamTrack(media.GetType(), media.GetType(), SessionToSender(session), session.GetOutgoingSourceGroup())

	// some callback event

	return streamerSession
}

func (s *StreamerSession) GetIncomingStreamTrack() *IncomingStreamTrack {
	return s.incoming
}

func (s *StreamerSession) GetOutgoingStreamTrack() *OutgoingStreamTrack {
	return s.outgoing
}

func (s *StreamerSession) Stop() {

	if s.session == nil {
		return
	}

	if s.incoming != nil {
		s.incoming.Stop()
	}

	if s.outgoing != nil {
		s.outgoing.Stop()
	}

	DeleteRTPSessionFacade(s.session)

	s.session = nil
}
