package mediaserver

import "github.com/notedit/media-server-go/sdp"

// RawRTPStreamer streamer that can send raw rtp data
type RawRTPStreamer struct {
	sessions map[string]*RawRTPStreamerSession
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

	s.sessions[session.id] = session

	return session
}

// Stop stop this streamer
func (s *RawRTPStreamer) Stop() {

	for _, session := range s.sessions {
		session.Stop()
	}

}
