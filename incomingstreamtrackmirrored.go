package mediaserver

import (
	"github.com/notedit/media-server-go/wrapper"
)

type IncomingStreamTrackMirrored struct {
	track     *IncomingStreamTrack
	receiver  native.RTPReceiverFacade
	counter   int
	encodings []*Encoding
}

func NewMirrorIncomingTrack(track *IncomingStreamTrack, timeService native.TimeService) *IncomingStreamTrackMirrored {

	mirror := &IncomingStreamTrackMirrored{}

	mirror.track = track
	mirror.receiver = track.receiver
	mirror.counter = 0
	mirror.encodings = []*Encoding{}

	return mirror
}
