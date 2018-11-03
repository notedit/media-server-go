package sdp

type RTCPFeedbackInfo struct {
	id     string
	params []string
}

func NewRTCPFeedbackInfo(id string, params []string) *RTCPFeedbackInfo {
	return &RTCPFeedbackInfo{id: id, params: params}
}

func (r *RTCPFeedbackInfo) Clone() *RTCPFeedbackInfo {
	rtcpfeedback := &RTCPFeedbackInfo{id: r.id}
	rtcpfeedback.params = make([]string, len(r.params))
	copy(rtcpfeedback.params, r.params)
	return rtcpfeedback
}

func (r *RTCPFeedbackInfo) GetID() string {
	return r.id
}

func (r *RTCPFeedbackInfo) GetParams() []string {
	return r.params
}
