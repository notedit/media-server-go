package sdptransform

type OriginStruct struct {
	Username       string `json:"username"`
	SessionId      string `json:"sessionId"`
	SessionVersion int    `json:"sessionVersion"`
	NetType        string `json:"netType"`
	IpVer          int    `json:"ipVer"`
	Address        string `json:"address"`
}

type GroupStruct struct {
	Type string `json:"type,omitempty"`
	Mids string `json:"mids,omitempty"`
}

type TimingStruct struct {
	Start int `json:"start"`
	Stop  int `json:"stop"`
}

type MsidSemanticStruct struct {
	Semantic string `json:"semantic,omitempty"`
	Token    string `json:"token,omitempty"`
}

type ConnectionStruct struct {
	Version int    `json:"version,omitempty"`
	Ip      string `json:"ip,omitempty"`
}

type RtpStruct struct {
	Payload  int    `json:"payload"`
	Codec    string `json:"codec"`
	Rate     int    `json:"rate,omitempty"`
	Encoding int    `json:"encoding,omitempty"`
}

type RtcpStruct struct {
	Port    int    `json:"port,omitempty"`
	NetType string `json:"netType,omitempty"`
	IpVer   int    `json:"ipVer,omitempty"`
	Address string `json:"address,omitempty"`
}

type FmtpStruct struct {
	Payload int    `json:"payload,omitempty"`
	Config  string `json:"config,omitempty"`
}

type FingerprintStruct struct {
	Type string `json:"type,omitempty"`
	Hash string `json:"hash,omitempty"`
}

type ExtStruct struct {
	Value int    `json:"value,omitempty"`
	Uri   string `json:"uri,omitempty"`
}

type RtcpFbStruct struct {
	Payload int    `json:"payload,omitempty"`
	Type    string `json:"type,omitempty"`
	Subtype string `json:"subtype,omitempty"`
}

type SsrcGroupStruct struct {
	Semantics string `json:"semantics,omitempty"`
	Ssrcs     string `json:"ssrcs,omitempty"`
}

type SsrcStruct struct {
	Id        uint   `json:"id,omitempty"`
	Attribute string `json:"attribute,omitempty"`
	Value     string `json:"value,omitempty"`
}

type BandwithStruct struct {
	Type  string `json:"type,omitempty"`
	Limit int    `json:"limit,omitempty"`
}

type CandidateStruct struct {
	Foundation string `json:"foundation,omitempty"`
	Component  int    `json:"component,omitempty"`
	Transport  string `json:"transport,omitempty"`
	Priority   int    `json:"priority,omitempty"`
	Ip         string `json:"ip,omitempty"`
	Port       int    `json:"port,omitempty"`
	Type       string `json:"type,omitempty"`
	Raddr      string `json:"raddr,omitempty"`
	Rport      int    `json:"aport,omitempty"`
}

type RidStruct struct {
	Id        string `json:"id,omitempty"`
	Direction string `json:"direction,omitempty"`
	Params    string `json:"params,omitempty"`
}

type SimulcastStruct struct {
	Dir1  string `json:"dir1"`
	List1 string `json:"list1"`
	Dir2  string `json:"dir2"`
	List2 string `json:"list2"`
}

// todo
type Simulcast03Struct struct {
}

type MediaStruct struct {
	Rtp         []*RtpStruct       `json:"rtp,omitempty"`
	Fmtp        []*FmtpStruct      `json:"fmtp,omitempty"`
	Type        string             `json:"type,omitempty"`
	Port        int                `json:"port,omitempty"`
	Protocal    string             `json:"protocal,omitempty"`
	Payloads    string             `json:"payloads,omitempty"`
	Connection  *ConnectionStruct  `json:"connection,omitempty"`
	Rtcp        *RtcpStruct        `json:"rtcp,omitempty"`
	IceUfrag    string             `json:"iceUfrag,omitempty"`
	IcePwd      string             `json:"icePwd,omitempty"`
	Fingerprint *FingerprintStruct `json:"fingerprint,omitempty"`
	Setup       string             `json:"setup,omitempty"`
	Mid         string             `json:"mid,omitempty"`
	Msid        string             `json:"msid,omitempty"`
	Ext         []*ExtStruct       `json:"ext,omitempty"`
	Direction   string             `json:"direction,omitempty"`
	RtcpRsize   string             `json:"rtcpRsize,omitempty"`
	RtcpMux     string             `json:"rtcpMux,omitempty"`
	RtcpFb      []*RtcpFbStruct    `json:"rtcpFb,omitempty"`
	Rids        []*RidStruct       `json:"rids,omitempty"`
	SsrcGroups  []*SsrcGroupStruct `json:"ssrcGroups,omitempty"`
	Ssrcs       []*SsrcStruct      `json:"ssrcs,omitempty"`
	Candidates  []*CandidateStruct `json:"candidates,omitempty"`
	Bandwidth   []*BandwithStruct  `json:"bandwidth,omitempty"`
	Simulcast   *SimulcastStruct   `json:"simulcast,omitempty"`
}

type SdpStruct struct {
	Version      int                 `json:"version"`
	Origin       *OriginStruct       `json:"origin"`
	Name         string              `json:"name"`
	Timing       *TimingStruct       `json:"timing,omitempty"`
	Groups       []*GroupStruct      `json:"groups,omitempty"`
	MsidSemantic *MsidSemanticStruct `json:"msidSemantic,omitempty"`
	Media        []*MediaStruct      `json:"media,omitempty"`
	Fingerprint  *FingerprintStruct  `json:"fingerprint,omitempty"`
	Connection   *ConnectionStruct   `json:"connection,omitempty"`
	Icelite      string              `json:"icelite,omitempty"`
}
