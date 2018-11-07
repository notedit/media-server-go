package mediaserver

import (
	"strconv"
	"time"
)

type Recorder struct {
	tracks     map[string]*RecorderTrack
	recorder   MP4Recorder
	ticker     *time.Ticker
	refresher  *Refresher
	maxTrackId int
}

func NewRecorder(filename string, waitForIntra bool, refresh int) *Recorder {
	recorder := &Recorder{}
	recorder.recorder = NewMP4Recorder()
	recorder.recorder.Create(filename)
	recorder.recorder.Record(waitForIntra)
	recorder.tracks = map[string]*RecorderTrack{}
	recorder.maxTrackId = 1

	if refresh > 0 {
		recorder.refresher = NewRefresher(refresh)
	}

	return recorder
}

func (r *Recorder) Record(incoming *IncomingStreamTrack) {

	for _, encoding := range incoming.GetEncodings() {
		encoding.GetDepacketizer().AddMediaListener(r.recorder)

		r.maxTrackId += 1
		recorderTrack := NewRecorderTrack(strconv.Itoa(r.maxTrackId), incoming, encoding)

		recorderTrack.Once("stopped", func() {

			recorderTrack.encoding.depacketizer.RemoveMediaListener(r.recorder)

			delete(r.tracks, recorderTrack.GetID())
		})
		r.tracks[recorderTrack.GetID()] = recorderTrack
	}

	if r.refresher != nil {
		r.refresher.Add(incoming)
	}
}

func (r *Recorder) RecordStream(incoming *IncomingStream) {

	for _, track := range incoming.GetTracks() {
		r.Record(track)
	}
}

func (r *Recorder) Stop() {

	if r.recorder == nil {
		return
	}

	for _, track := range r.tracks {
		track.Stop()
	}

	if r.refresher != nil {
		r.refresher.Stop()
	}

	r.recorder.Close()

	r.refresher = nil
	r.recorder = nil
}
