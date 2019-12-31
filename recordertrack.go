package mediaserver

type RecorderTrackStopListener func()

// RecorderTrack  a track to record
type RecorderTrack struct {
	id       string
	track    *IncomingStreamTrack
	encoding *Encoding
}

// NewRecorderTrack create a new recorder track
func NewRecorderTrack(id string, track *IncomingStreamTrack, encoding *Encoding) *RecorderTrack {

	recorderTrack := &RecorderTrack{}
	recorderTrack.id = id
	recorderTrack.track = track
	recorderTrack.encoding = encoding

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

// Stop stop the recorder track
func (r *RecorderTrack) Stop() {

	if r.track == nil {
		return
	}

	r.track = nil
	r.encoding = nil
}
