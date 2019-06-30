package mediaserver

import (
	"fmt"
	"strings"

	"github.com/gofrs/uuid"
	"github.com/notedit/sdp"

	native "github.com/notedit/media-server-go/wrapper"
)

// StreamerSession represent a rtp session
type StreamerSession struct {
	id              string
	local           bool
	port            int
	ip              string
	incoming        *IncomingStreamTrack
	outgoing        *OutgoingStreamTrack
	session         native.RTPSessionFacade
	onStopListeners []func()
}

// NewStreamerSession new StreamerSession with auto selectd port
func NewStreamerSession(media *sdp.MediaInfo) *StreamerSession {

	streamerSession := &StreamerSession{}
	var mediaType native.MediaFrameType = 0
	if strings.ToLower(media.GetType()) == "video" {
		mediaType = 1
	}
	session := native.NewRTPSessionFacade(mediaType)

	streamerSession.id = uuid.Must(uuid.NewV4()).String()

	properties := native.NewProperties()

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

	native.DeleteProperties(properties)

	streamerSession.session = session

	streamerSession.incoming = NewIncomingStreamTrack(media.GetType(), media.GetType(), native.SessionToReceiver(session), map[string]native.RTPIncomingSourceGroup{"": session.GetIncomingSourceGroup()})

	streamerSession.outgoing = newOutgoingStreamTrack(media.GetType(), media.GetType(), native.SessionToSender(session), session.GetOutgoingSourceGroup())

	streamerSession.onStopListeners = make([]func(), 0)

	return streamerSession
}

// NewStreamerSessionWithLocalPort  create streamer session with pre selected port
func NewStreamerSessionWithLocalPort(port int, media *sdp.MediaInfo) *StreamerSession {

	streamerSession := &StreamerSession{}
	var mediaType native.MediaFrameType = 0
	if strings.ToLower(media.GetType()) == "video" {
		mediaType = 1
	}
	session := native.NewRTPSessionFacade(mediaType)

	streamerSession.id = uuid.Must(uuid.NewV4()).String()

	properties := native.NewProperties()

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

	session.SetLocalPort(port)

	session.Init(properties)

	native.DeleteProperties(properties)

	streamerSession.session = session

	streamerSession.incoming = NewIncomingStreamTrack(media.GetType(), media.GetType(), native.SessionToReceiver(session), map[string]native.RTPIncomingSourceGroup{"": session.GetIncomingSourceGroup()})

	streamerSession.outgoing = newOutgoingStreamTrack(media.GetType(), media.GetType(), native.SessionToSender(session), session.GetOutgoingSourceGroup())

	streamerSession.onStopListeners = make([]func(), 0)

	return streamerSession
}

// GetID get id
func (s *StreamerSession) GetID() string {
	return s.id
}

func (s *StreamerSession) GetLocalPort() int {
	return s.session.GetLocalPort()
}

func (s *StreamerSession) SetRemotePort(ip string, port int) {
	s.session.SetRemotePort(ip, port)
}

// GetIncomingStreamTrack get asso incoming track,
func (s *StreamerSession) GetIncomingStreamTrack() *IncomingStreamTrack {
	return s.incoming
}

// GetOutgoingStreamTrack get asso outgoing track,
func (s *StreamerSession) GetOutgoingStreamTrack() *OutgoingStreamTrack {
	return s.outgoing
}

// OnStop register stop listener
func (s *StreamerSession) OnStop(stop func()) {
	s.onStopListeners = append(s.onStopListeners, stop)
}

// Stop it
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

	s.session.End()

	native.DeleteRTPSessionFacade(s.session)

	for _, stopFunc := range s.onStopListeners {
		stopFunc()
	}

	s.session = nil
}
