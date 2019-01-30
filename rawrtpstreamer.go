package mediaserver

import (
	"github.com/notedit/media-server-go/sdp"
	"sync"
)

// RawRTPStreamer streamer that can send raw rtp data
type RawRTPStreamer struct {
	sessions map[string]*RawRTPStreamerSession
	sync.Mutex
}

// NewRawRTPStreamer create new raw rtp streamer
func NewRawRTPStreamer() *RawRTPStreamer {
	streamer := &RawRTPStreamer{}
	streamer.sessions = make(map[string]*RawRTPStreamerSession)
	return streamer
}

// CreateSession create a audio/media session
func (s *RawRTPStreamer) CreateSession(media *sdp.MediaInfo) *RawRTPStreamerSession {

	session := NewRawRTPStreamerSession(media)

	s.Lock()
	s.sessions[session.id] = session
	s.Unlock()

	return session
}

func (s *RawRTPStreamer) RemoveSession(session *RawRTPStreamerSession) {

	s.Lock()
	delete(s.sessions, session.id)
	s.Unlock()

}

// Stop stop this streamer
func (s *RawRTPStreamer) Stop() {

	for _, session := range s.sessions {
		session.Stop()
	}

	s.sessions = nil

}
