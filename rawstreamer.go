package mediaserver

import (
	"sync"

	"github.com/notedit/sdp"
)

// RawStreamer streamer that can send raw rtp data
type RawStreamer struct {
	sessions map[string]*RawStreamerSession
	sync.Mutex
}

// NewRawStreamer create new raw rtp streamer
func NewRawStreamer() *RawStreamer {
	streamer := &RawStreamer{}
	streamer.sessions = make(map[string]*RawStreamerSession)
	return streamer
}

// CreateSession create a audio/media session
func (s *RawStreamer) CreateSession(media *sdp.MediaInfo) *RawStreamerSession {

	session := NewRawStreamerSession(media)

	s.Lock()
	s.sessions[session.id] = session
	s.Unlock()

	return session
}

// RemoveSession remove a session
func (s *RawStreamer) RemoveSession(session *RawStreamerSession) {
	s.Lock()
	delete(s.sessions, session.id)
	s.Unlock()
}

// Stop stop this streamer
func (s *RawStreamer) Stop() {

	for _, session := range s.sessions {
		session.Stop()
	}
	s.sessions = nil
}
