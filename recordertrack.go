package mediaserver

import "github.com/chuckpreslar/emission"

type RecorderTrack struct {
	id       string
	track    *IncomingStreamTrack
	encoding *Encoding
	*emission.Emitter
}

func NewRecorderTrack(id string, track *IncomingStreamTrack, encoding *Encoding) *RecorderTrack {

	recorderTrack := &RecorderTrack{}
	recorderTrack.id = id
	recorderTrack.track = track
	recorderTrack.encoding = encoding

	recorderTrack.Emitter = emission.NewEmitter()

	track.Once("stopped", recorderTrack.onTrackStopped)

	return recorderTrack
}

func (r *RecorderTrack) GetID() string {
	return r.id
}

func (r *RecorderTrack) GetTrack() *IncomingStreamTrack {
	return r.track
}

func (r *RecorderTrack) GetEncoding() *Encoding {
	return r.encoding
}

func (r *RecorderTrack) Stop() {

	if r.track == nil {
		return
	}

	r.track.Off("stopped", r.onTrackStopped)

	r.EmitSync("stopped")

	r.track = nil
	r.encoding = nil
}

func (r *RecorderTrack) onTrackStopped() {
	r.Stop()
}
