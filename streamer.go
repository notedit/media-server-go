package mediaserver

import sdp "./sdp"

type Streamer struct {
	sessions []*StreamerSession
}

func NewStreamer() *Streamer {
	streamer := &Streamer{}
	streamer.sessions = []*StreamerSession{}
	return streamer
}

func (s *Streamer) CreateSession(local bool, ip string, port int, media *sdp.MediaInfo) *StreamerSession {

	session := NewStreamerSession(local, ip, port, media)
	// todo stopped event
	s.sessions = append(s.sessions, session)
	return session
}

func (s *Streamer) Stop() {

	for _, session := range s.sessions {
		session.Stop()
	}
	s.sessions = nil
}
