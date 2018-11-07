package mediaserver

import (
	"github.com/notedit/media-server-go/sdp"
)

type Streamer struct {
	sessions map[string]*StreamerSession
}

func NewStreamer() *Streamer {
	streamer := &Streamer{}
	streamer.sessions = make(map[string]*StreamerSession)
	return streamer
}

func (s *Streamer) CreateSession(local bool, ip string, port int, media *sdp.MediaInfo) *StreamerSession {

	session := NewStreamerSession(local, ip, port, media)

	session.Once("stopped", func() {
		delete(s.sessions, session.GetID())
	})

	s.sessions[session.id] = session
	return session
}

func (s *Streamer) Stop() {

	for _, session := range s.sessions {
		session.Stop()
	}
	s.sessions = nil
}
