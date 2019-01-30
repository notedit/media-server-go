package mediaserver

import (
	"github.com/notedit/media-server-go/sdp"
	"sync"
)

type Streamer struct {
	sessions map[string]*StreamerSession
	sync.Mutex
}

func NewStreamer() *Streamer {
	streamer := &Streamer{}
	streamer.sessions = make(map[string]*StreamerSession)
	return streamer
}

func (s *Streamer) CreateSession(local bool, ip string, port int, media *sdp.MediaInfo) *StreamerSession {

	session := NewStreamerSession(local, ip, port, media)

	session.OnStop(func() {
		s.Lock()
		delete(s.sessions, session.GetID())
		s.Unlock()
	})

	s.Lock()
	s.sessions[session.id] = session
	s.Unlock()

	return session
}

func (s *Streamer) Stop() {

	for _, session := range s.sessions {
		session.Stop()
	}
	s.sessions = nil
}
