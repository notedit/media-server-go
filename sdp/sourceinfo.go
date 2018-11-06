package sdp

type SourceInfo struct {
	ssrc     uint
	cname    string
	streamID string
	trackID  string
}

func NewSourceInfo(ssrc uint) *SourceInfo {

	return &SourceInfo{ssrc: ssrc}
}

func (s *SourceInfo) Clone() *SourceInfo {
	cloned := &SourceInfo{}
	cloned.ssrc = s.ssrc
	cloned.cname = s.cname
	cloned.streamID = s.streamID
	cloned.trackID = s.trackID
	return cloned
}

func (s *SourceInfo) GetCName() string {
	return s.cname
}

func (s *SourceInfo) SetCName(cname string) {
	s.cname = cname
}

func (s *SourceInfo) GetStreamID() string {
	return s.streamID
}

func (s *SourceInfo) SetStreamID(streamID string) {
	s.streamID = streamID
}

func (s *SourceInfo) GetTrackID() string {
	return s.trackID
}

func (s *SourceInfo) SetTrackID(trackID string) {
	s.trackID = trackID
}

func (s *SourceInfo) GetSSRC() uint {
	return s.ssrc
}
