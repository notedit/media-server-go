package mediaserver

import (
	"github.com/notedit/media-server-go/wrapper"
	"github.com/notedit/sdp"
)



type mirrorEncoding struct {
	id           string
	source       native.RTPIncomingMediaStreamMultiplexer
	depacketizer native.StreamTrackDepacketizer
}


type IncomingStreamTrackMirrored struct {
	track     *IncomingStreamTrack
	receiver  native.RTPReceiverFacade
	counter   int
	encodings []*mirrorEncoding
}

func NewMirrorIncomingTrack(track *IncomingStreamTrack, timeService native.TimeService) *IncomingStreamTrackMirrored {

	mirror := &IncomingStreamTrackMirrored{}

	mirror.track = track
	mirror.receiver = track.receiver
	mirror.counter = 0
	mirror.encodings = []*mirrorEncoding{}

	for _,encoding :=  range track.GetEncodings() {
		source := native.NewRTPIncomingMediaStreamMultiplexer(encoding.source.GetMedia().GetSsrc(), timeService)
		encoding.source.AddListener(source)

		newEncoding := &mirrorEncoding{
			id:           encoding.id,
			source:       source,
			depacketizer: native.NewStreamTrackDepacketizer(source.SwigGetRTPIncomingMediaStream()),
		}

		mirror.encodings = append(mirror.encodings, newEncoding)

	}

	mirror.track.Attached()

	return mirror
}


func (t *IncomingStreamTrackMirrored) GetStats() map[string]*IncomingAllStats {
	return t.track.GetStats()
}


func (t *IncomingStreamTrackMirrored) GetActiveLayers() *ActiveLayersInfo {
	return t.track.GetActiveLayers()
}


func (t *IncomingStreamTrackMirrored) GetID() string {
	return t.track.GetID()
}

func (t *IncomingStreamTrackMirrored) GetTrackInfo() *sdp.TrackInfo {
	return t.GetTrackInfo()
}

func (t *IncomingStreamTrackMirrored) GetSSRCs() {

}

func (t *IncomingStreamTrackMirrored) GetMedia() string {
	return t.track.GetMedia()
}


func (t *IncomingStreamTrackMirrored) Attached() bool {

	t.counter += 1

	if t.counter == 1 {
		return true
	}

	return false
}


func (t *IncomingStreamTrackMirrored) Refresh()  {

	t.track.Refresh()
}


func (t *IncomingStreamTrackMirrored) Detached() bool {

	if t.counter == 0 {
		return true
	}

	t.counter -= 1

	if t.counter == 0 {
		return true
	}

	return false
}


func (t *IncomingStreamTrackMirrored) Stop() {

	if t.track == nil {
		return
	}

	for i,encoding := range t.track.GetEncodings()  {
		mencoding := t.encodings[i]
		encoding.GetSource().RemoveListener(mencoding.source)
		native.DeleteRTPIncomingMediaStreamMultiplexer(t.encodings[i].source)
		mencoding.depacketizer.Stop()
		native.DeleteStreamTrackDepacketizer(mencoding.depacketizer)
	}


	t.track.Detached()

	t.encodings = nil

	t.track = nil

	t.receiver = nil

}













