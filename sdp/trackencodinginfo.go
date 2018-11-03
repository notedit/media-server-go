package sdp

type TrackEncodingInfo struct {
	id     string
	paused bool
	codecs map[string]*CodecInfo
	params map[string]string
}

func NewTrackEncodingInfo(id string, paused bool) *TrackEncodingInfo {

	info := &TrackEncodingInfo{
		id:     id,
		paused: paused,
		codecs: map[string]*CodecInfo{},
		params: map[string]string{},
	}
	return info
}

func (t *TrackEncodingInfo) Clone() *TrackEncodingInfo {

	cloned := NewTrackEncodingInfo(t.id, t.paused)

	for k, v := range t.codecs {
		cloned.codecs[k] = v.Clone()
	}

	for k, v := range t.params {
		cloned.params[k] = v
	}
	return cloned
}

func (t *TrackEncodingInfo) GetID() string {

	return t.id
}

func (t *TrackEncodingInfo) GetCodecs() map[string]*CodecInfo {

	return t.codecs
}

func (t *TrackEncodingInfo) AddCodec(codec *CodecInfo) {
	t.codecs[codec.GetCodec()] = codec
}

func (t *TrackEncodingInfo) GetParams() map[string]string {
	return t.params
}

func (t *TrackEncodingInfo) SetParams(params map[string]string) {
	t.params = params
}

func (t *TrackEncodingInfo) AddParam(id, param string) {
	t.params[id] = param
}

func (t *TrackEncodingInfo) IsPaused() bool {

	return t.paused
}
