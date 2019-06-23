package mediaserver

import (
	"strconv"
	"time"

	native "github.com/notedit/media-server-go/wrapper"
)

// Recorder represent a file recorder
type Recorder struct {
	tracks     map[string]*RecorderTrack
	recorder   native.MP4Recorder
	ticker     *time.Ticker
	refresher  *Refresher
	maxTrackId int
}

// NewRecorder create a new recorder
func NewRecorder(filename string, waitForIntra bool, refresh int) *Recorder {
	recorder := &Recorder{}
	recorder.recorder = native.NewMP4Recorder()
	recorder.recorder.Create(filename)
	recorder.recorder.Record(waitForIntra)
	recorder.tracks = map[string]*RecorderTrack{}
	recorder.maxTrackId = 1

	if refresh > 0 {
		recorder.refresher = NewRefresher(refresh)
	}

	return recorder
}

// Record start record an incoming track
func (r *Recorder) Record(incoming *IncomingStreamTrack) {

	for _, encoding := range incoming.GetEncodings() {
		encoding.GetDepacketizer().AddMediaListener(r.recorder)

		r.maxTrackId += 1
		recorderTrack := NewRecorderTrack(strconv.Itoa(r.maxTrackId), incoming, encoding)

		recorderTrack.OnStop(func() {
			recorderTrack.encoding.depacketizer.RemoveMediaListener(r.recorder)
			delete(r.tracks, recorderTrack.GetID())
		})

		r.tracks[recorderTrack.GetID()] = recorderTrack
	}

	if r.refresher != nil {
		r.refresher.Add(incoming)
	}
}

// RecordStream start record an incoming stream
func (r *Recorder) RecordStream(incoming *IncomingStream) {

	for _, track := range incoming.GetTracks() {
		r.Record(track)
	}
}

// Stop  stop the recorder
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

	native.DeleteMP4Recorder(r.recorder)

	r.refresher = nil
	r.recorder = nil
}
