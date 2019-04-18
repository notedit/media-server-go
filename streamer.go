package mediaserver

import (
	"sync"

	"github.com/notedit/sdp"
)

// Streamer
type Streamer struct {
	sessions map[string]*StreamerSession
	sync.Mutex
}

// NewStreamer create a streamer
func NewStreamer() *Streamer {
	streamer := &Streamer{}
	streamer.sessions = make(map[string]*StreamerSession)
	return streamer
}

// CreateSession
func (s *Streamer) CreateSession(local bool, ip string, port int, media *sdp.MediaInfo) *StreamerSession {

	session := NewStreamerSession(local, ip, port, media)

	s.Lock()
	s.sessions[session.id] = session
	s.Unlock()

	return session
}

// RemoveSession remove a session
func (s *Streamer) RemoveSession(session *StreamerSession) {
	s.Lock()
	delete(s.sessions, session.id)
	s.Unlock()
}

// Stop stop all sessions
func (s *Streamer) Stop() {

	for _, session := range s.sessions {
		session.Stop()
	}
	s.sessions = nil
}
