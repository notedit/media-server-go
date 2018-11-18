package sdptransform

type OriginStruct struct {
	Username       string `json:"username"`
	SessionId      string `json:"sessionId"`
	SessionVersion int    `json:"sessionVersion"`
	NetType        string `json:"netType"`
	Address        string `json:"address"`
}

type GroupStruct struct {
	Type string `json:"type"`
	Mids string `json:"mids"`
}

type MsidSemanticStruct struct {
	Semantic string `json:"semantic"`
	Token    string `json:"token"`
}

type ConnectionStruct struct {
	Version int    `json:"version"`
	Ip      string `json:"ip"`
}

type RtpStruct struct {
	Payload  int    `json:"payload"`
	Codec    string `json:"codec"`
	Rate     int    `json:"rate"`
	Encoding int    `json:"encoding"`
}

type RtcpStruct struct {
	Port    int    `json:"port"`
	NetType string `json:"netType"`
	IpVer   int    `json:"ipVer"`
	Address string `json:"address"`
}

type FmtpStruct struct {
	Payload int    `json:"payload"`
	Config  string `json:"config"`
}

type FingerprintStruct struct {
	Type string `json:"type"`
	Hash string `json:"hash"`
}

type ExtStruct struct {
	Value int    `json:"value"`
	Uri   string `json:"uri"`
}

type RtcpFbStruct struct {
	Payload int    `json:"payload"`
	Type    string `json:"type"`
	Subtype string `json:"subtype"`
}

type SsrcGroupStruct struct {
	Semantics string `json:"semantics"`
	Ssrcs     string `json:"ssrcs"`
	SsrcArr   []string
}

type SsrcStruct struct {
	Id        uint   `json:"id"`
	Attribute string `json:"attribute"`
	Value     string `json:"value"`
}

type BandwithStruct struct {
	Type  string `json:"type"`
	Limit int    `json:"limit"`
}

type CandidateStruct struct {
	Foundation string `json:"foundation"`
	Component  int    `json:"component"`
	Transport  string `json:"transport"`
	Priority   int    `json:"priority"`
	Ip         string `json:"ip"`
	Port       int    `json:"port"`
	Type       string `json:"type"`
	Raddr      string `json:"raddr"`
	Rport      int    `json:"aport"`
}

type RidStruct struct {
	Id        string `json:"id"`
	Direction string `json:"direction"`
	Params    string `json:"params"`
}

type MediaStruct struct {
	Rtp         []*RtpStruct       `json:"rtp"`
	Fmtp        []*FmtpStruct      `json:"fmtp"`
	Type        string             `json:"type"`
	Port        int                `json:"port"`
	Protocal    string             `json:"protocal"`
	Payloads    string             `json:"payloads"`
	Connection  *ConnectionStruct  `json:"connection"`
	Rtcp        *RtcpStruct        `json:"rtcp"`
	IceUfrag    string             `json:"iceUfrag"`
	IcePwd      string             `json:"icePwd"`
	Fingerprint *FingerprintStruct `json:"fingerprint"`
	Setup       string             `json:"setup"`
	Mid         string             `json:"mid"`
	Msid        string             `json:"msid"`
	Ext         []*ExtStruct       `json:"ext"`
	Direction   string             `json:"direction"`
	RtcpRsize   string             `json:"rtcpRsize"`
	RtcpMux     string             `json:"rtcpMux"`
	RtcpFb      []*RtcpFbStruct    `json:"rtcpFb"`
	Rids        []*RidStruct       `json:"rids"`
	SsrcGroups  []*SsrcGroupStruct `json:"ssrcGroups"`
	Ssrcs       []*SsrcStruct      `json:"ssrcs"`
	Candidates  []*CandidateStruct `json:"candidates"`
	Bandwidth   []*BandwithStruct  `json:"bandwidth"`
}

type SdpStruct struct {
	Version      int                 `json:"version"`
	Origin       *OriginStruct       `json:"origin"`
	Name         string              `json:"name"`
	Timing       interface{}         `json:"timing"`
	Groups       []*GroupStruct      `json:"groups"`
	MsidSemantic *MsidSemanticStruct `json:"msidSemantic"`
	Media        []*MediaStruct      `json:"media"`
	Fingerprint  *FingerprintStruct  `json:"fingerprint"`
	Connection   *ConnectionStruct   `json:"connection"`
	Icelite      string              `json:"icelite"`
}
