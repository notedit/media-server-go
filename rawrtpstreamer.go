package mediaserver

import "github.com/notedit/media-server-go/sdp"

type RawRTPStreamer struct {
	sessions map[string]*RawRTPStreamerSession
}

func NewRawRTPStreamer() *RawRTPStreamer {
	streamer := &RawRTPStreamer{}
	streamer.sessions = make(map[string]*RawRTPStreamerSession)
	return streamer
}

func (s *RawRTPStreamer) CreateSession(media *sdp.MediaInfo) *RawRTPStreamerSession {

	session := NewRawRTPStreamerSession(media)

	session.Once("stopped", func() {
		delete(s.sessions, session.GetID())
	})

	s.sessions[session.id] = session
	return session
}

func (s *RawRTPStreamer) Stop() {

	for _, session := range s.sessions {
		session.Stop()
	}

}
