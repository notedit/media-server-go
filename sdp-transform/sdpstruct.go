package sdptransform

type OriginStruct struct {
	Username       string
	SessionId      string
	SessionVersion int
	NetType        string
	Address        string
}

type GroupStruct struct {
	Type string
	Mids string
}

type MsidSemanticStruct struct {
	Semantic string
	Token    string
}

type ConnectionStruct struct {
	Version int
	Ip      string
}

type RtpStruct struct {
	Payload  int
	Codec    string
	Rate     int
	Encoding int
}

type RtcpStruct struct {
	Port    int
	NetType string
	IpVer   int
	Address string
}

type FmtpStruct struct {
	Payload int
	Config  string
}

type FingerprintStruct struct {
	Type string
	Hash string
}

type ExtStruct struct {
	Value int
	Uri   string
}

type RtcpFbStruct struct {
	Payload int
	Type    string
	Subtype string
}

type SsrcGroupStruct struct {
	Semantics string
	Ssrcs     string
	SsrcArr   []string
}

type SsrcStruct struct {
	Id        uint
	Attribute string
	Value     string
}

type BandwithStruct struct {
	Type  string
	Limit int
}

type CandidateStruct struct {
	Foundation string
	Component  int
	Transport  string
	Priority   int
	Ip         string
	Port       int
	Type       string
	Raddr      string
	Rport      int
}

type RidStruct struct {
	Id        string
	Direction string
	Params    string
}

type MediaStruct struct {
	Rtp         []*RtpStruct
	Fmtp        []*FmtpStruct
	Type        string
	Port        int
	Protocal    string
	Payloads    string
	Connection  *ConnectionStruct
	Rtcp        *RtcpStruct
	IceUfrag    string
	IcePwd      string
	Fingerprint *FingerprintStruct
	Setup       string
	Mid         string
	Msid        string
	Ext         []*ExtStruct
	Direction   string
	RtcpRsize   string
	RtcpMux     string
	RtcpFb      []*RtcpFbStruct
	Rids        []*RidStruct
	SsrcGroups  []*SsrcGroupStruct
	Ssrcs       []*SsrcStruct
	Candidates  []*CandidateStruct
	Bandwidth   []*BandwithStruct
}

type SdpStruct struct {
	Version      int
	Origin       *OriginStruct
	Name         string
	Timing       interface{}
	Groups       []*GroupStruct
	MsidSemantic *MsidSemanticStruct
	Media        []*MediaStruct
	Fingerprint  *FingerprintStruct
	Connection   *ConnectionStruct
	Icelite      string
}
