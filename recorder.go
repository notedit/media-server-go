package mediaserver

import "time"

type Recorder struct {
	recorder  MP4Recorder
	tracks    map[string]*IncomingStreamTrack
	ticker    *time.Ticker
	refresher *Refresher
}

func NewRecorder(filename string, waitForIntra bool, refresh int) *Recorder {
	recorder := &Recorder{}
	recorder.recorder = NewMP4Recorder()
	recorder.recorder.Create(filename)
	recorder.recorder.Record(waitForIntra)

	if refresh > 0 {
		recorder.refresher = NewRefresher(refresh)
	}

	// todo track stoped event
	return recorder
}

func (r *Recorder) Record(incom *IncomingStreamTrack) {

	// todo
}

func (r *Recorder) Stop() {
	// todo
}
