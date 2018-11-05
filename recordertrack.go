package mediaserver

type RecorderTrack struct {
	id       string
	track    *IncomingStreamTrack
	encoding interface{}
}

func NewRecorderTrack(id string, track *IncomingStreamTrack, encoding interface{}) *RecorderTrack {

	recorderTrack := &RecorderTrack{}
	recorderTrack.id = id
	recorderTrack.track = track
	recorderTrack.encoding = encoding

	// todo event callback
	return recorderTrack
}

func (r *RecorderTrack) GetID() string {
	return r.id
}

func (r *RecorderTrack) GetTrack() *IncomingStreamTrack {
	return r.track
}

func (r *RecorderTrack) GetEncoding() interface{} {
	return r.encoding
}

func (r *RecorderTrack) Stop() {

	if r.track == nil {
		return
	}

	r.track = nil
	r.encoding = nil
}
