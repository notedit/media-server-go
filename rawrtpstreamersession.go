package mediaserver

import (
	"github.com/chuckpreslar/emission"
	"github.com/notedit/media-server-go/sdp"
)

type RawRTPStreamerSession struct {
	id       string
	incoming *IncomingStreamTrack
	*emission.Emitter
}

func NewRawRTPStreamerSession(media *sdp.MediaInfo) *RawRTPStreamerSession {

	return nil
}

func (s *RawRTPStreamerSession) GetID() string {
	return s.id
}

func (s *RawRTPStreamerSession) GetIncomingStreamTrack() *IncomingStreamTrack {
	return s.incoming
}

func (s *RawRTPStreamerSession) Stop() {

}
