package mediaserver

type RecorderTrackStopListener func()

// RecorderTrack  a track to record
type RecorderTrack struct {
	id              string
	track           *IncomingStreamTrack
	encoding        *Encoding
	onStopListeners []RecorderTrackStopListener
}

// NewRecorderTrack create a new recorder track
func NewRecorderTrack(id string, track *IncomingStreamTrack, encoding *Encoding) *RecorderTrack {

	recorderTrack := &RecorderTrack{}
	recorderTrack.id = id
	recorderTrack.track = track
	recorderTrack.encoding = encoding

	track.OnStop(func() {
		recorderTrack.Stop()
	})

	recorderTrack.onStopListeners = make([]RecorderTrackStopListener, 0)

	return recorderTrack
}

// GetID  get recorder track id
func (r *RecorderTrack) GetID() string {
	return r.id
}

// GetTrack get internal IncomingStreamTrack
func (r *RecorderTrack) GetTrack() *IncomingStreamTrack {
	return r.track
}

// GetEncoding get encoding info
func (r *RecorderTrack) GetEncoding() *Encoding {
	return r.encoding
}

// OnStop register a stop listener
func (r *RecorderTrack) OnStop(stop RecorderTrackStopListener) {
	r.onStopListeners = append(r.onStopListeners, stop)
}

// Stop stop the recorder track
func (r *RecorderTrack) Stop() {

	if r.track == nil {
		return
	}

	for _, stopFunc := range r.onStopListeners {
		stopFunc()
	}

	r.onStopListeners = nil
	r.track = nil
	r.encoding = nil
}
