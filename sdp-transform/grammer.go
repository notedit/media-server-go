package sdptransform

import (
	"regexp"

	"github.com/Jeffail/gabs"
)

type Rule struct {
	Name       string
	Push       string
	Reg        *regexp.Regexp
	Names      []string
	Types      []rune
	Format     string
	FormatFunc func(obj *gabs.Container) string
}

func hasValue(obj *gabs.Container, key string) bool {

	if !obj.Exists(key) {
		return false
	}
	value := obj.Path(key)
	if str, ok := value.Data().(string); ok {
		if len(str) == 0 {
			return false
		}
		return true
	} else if _, ok := value.Data().(int); ok {
		return true
	} else {
		return false
	}
}

var rulesMap map[byte][]*Rule = map[byte][]*Rule{
	'v': []*Rule{
		&Rule{
			Name:   "version",
			Push:   "",
			Reg:    regexp.MustCompile("^(\\d*)$"),
			Names:  []string{},
			Types:  []rune{'d'},
			Format: "%d",
		},
	},
	'o': []*Rule{
		&Rule{
			Name:   "origin",
			Push:   "",
			Reg:    regexp.MustCompile("^(\\S*) (\\S*) (\\d*) (\\S*) IP(\\d) (\\S*)"),
			Names:  []string{"username", "sessionId", "sessionVersion", "netType", "ipVer", "address"},
			Types:  []rune{'s', 's', 'd', 's', 'd', 's'},
			Format: "%s %s %d %s IP%d %s",
		},
	},
	's': []*Rule{
		&Rule{
			Name:   "name",
			Push:   "",
			Reg:    regexp.MustCompile("(.*)"),
			Names:  []string{},
			Types:  []rune{'s'},
			Format: "%s",
		},
	},
	'i': []*Rule{
		&Rule{
			Name:   "description",
			Push:   "",
			Reg:    regexp.MustCompile("(.*)"),
			Names:  []string{},
			Types:  []rune{'s'},
			Format: "%s",
		},
	},
	'u': []*Rule{
		&Rule{
			Name:   "uri",
			Push:   "",
			Reg:    regexp.MustCompile("(.*)"),
			Names:  []string{},
			Types:  []rune{'s'},
			Format: "%s",
		},
	},
	'e': []*Rule{
		&Rule{
			Name:   "email",
			Push:   "",
			Reg:    regexp.MustCompile("(.*)"),
			Names:  []string{},
			Types:  []rune{'s'},
			Format: "%s",
		},
	},
	'p': []*Rule{
		&Rule{
			Name:   "phone",
			Push:   "",
			Reg:    regexp.MustCompile("(.*)"),
			Names:  []string{},
			Types:  []rune{'s'},
			Format: "%s",
		},
	},
	'z': []*Rule{
		&Rule{
			Name:   "timezones",
			Push:   "",
			Reg:    regexp.MustCompile("(.*)"),
			Names:  []string{},
			Types:  []rune{'s'},
			Format: "%s",
		},
	},
	'r': []*Rule{
		&Rule{
			Name:   "repeats",
			Push:   "",
			Reg:    regexp.MustCompile("(.*)"),
			Names:  []string{},
			Types:  []rune{'s'},
			Format: "%s",
		},
	},
	't': []*Rule{
		&Rule{
			Name:   "timing",
			Push:   "",
			Reg:    regexp.MustCompile("^(\\d*) (\\d*)"),
			Names:  []string{"start", "stop"},
			Types:  []rune{'d', 'd'},
			Format: "%d %d",
		},
	},
	'c': []*Rule{
		&Rule{
			Name:   "connection",
			Push:   "",
			Reg:    regexp.MustCompile("^IN IP(\\d) (\\S*)"),
			Names:  []string{"version", "ip"},
			Types:  []rune{'d', 's'},
			Format: "IN IP%d %s",
		},
	},
	'b': []*Rule{
		&Rule{
			Name:   "",
			Push:   "bandwidth",
			Reg:    regexp.MustCompile("^(TIAS|AS|CT|RR|RS):(\\d*)"),
			Names:  []string{"type", "limit"},
			Types:  []rune{'s', 'd'},
			Format: "%s:%d",
		},
	},
	'm': []*Rule{ // m=video 51744 RTP/AVP 126 97 98 34 31
		&Rule{
			Name:   "",
			Push:   "",
			Reg:    regexp.MustCompile("^(\\w*) (\\d*) ([\\w\\/]*)(?: (.*))?"),
			Names:  []string{"type", "port", "protocal", "payloads"},
			Types:  []rune{'s', 'd', 's', 's'},
			Format: "%s %d %s %s",
		},
	},
	'a': []*Rule{ // a=rtpmap:110 opus/48000/2
		&Rule{
			Name:   "",
			Push:   "rtp",
			Reg:    regexp.MustCompile("^rtpmap:(\\d*) ([\\w\\-\\.]*)(?:\\s*\\/(\\d*)(?:\\s*\\/(\\d*))?)?"),
			Names:  []string{"payload", "codec", "rate", "encoding"},
			Types:  []rune{'d', 's', 'd', 'd'},
			Format: "",
			FormatFunc: func(obj *gabs.Container) string {
				var ret string
				if hasValue(obj, "encoding") {
					ret = "rtpmap:%d %s/%d/%d"
				} else {
					if hasValue(obj, "rate") {
						ret = "rtpmap:%d %s/%d"
					} else {
						ret = "rtpmap:%d %s"
					}
				}
				return ret
			},
		},
		// a=fmtp:108 profile-level-id=24;object=23;bitrate=64000
		// a=fmtp:111 minptime=10; useinbandfec=1
		&Rule{
			Name:   "",
			Push:   "fmtp",
			Reg:    regexp.MustCompile("^fmtp:(\\d*) ([\\S| ]*)"),
			Names:  []string{"payload", "config"},
			Types:  []rune{'d', 's'},
			Format: "fmtp:%d %s",
		},
		// a=control:streamid=0
		&Rule{
			Name:   "control",
			Push:   "",
			Reg:    regexp.MustCompile("^control:(.*)"),
			Names:  []string{},
			Types:  []rune{'s'},
			Format: "controle:%s",
		},
		// a=rtcp:65179 IN IP4 193.84.77.194
		&Rule{
			Name:   "rtcp",
			Push:   "",
			Reg:    regexp.MustCompile("^rtcp:(\\d*)(?: (\\S*) IP(\\d) (\\S*))?"),
			Names:  []string{"port", "netType", "ipVer", "address"},
			Types:  []rune{'d', 's', 'd', 's'},
			Format: "",
			FormatFunc: func(obj *gabs.Container) string {
				if hasValue(obj, "address") {
					return "rtcp:%d %s IP%d %s"
				} else {
					return "rtcp:%d"
				}
			},
		},
		// a=rtcp-fb:98 trr-int 100
		&Rule{
			Name:   "",
			Push:   "rtcpFbTrrInt",
			Reg:    regexp.MustCompile("^rtcp-fb:(\\*|\\d*) trr-int (\\d*)"),
			Names:  []string{"payload", "value"},
			Types:  []rune{'s', 'd'},
			Format: "rtcp-fb:%s trr-int %d",
		},
		// a=rtcp-fb:98 nack rpsi
		&Rule{
			Name:   "",
			Push:   "rtcpFb",
			Reg:    regexp.MustCompile("^rtcp-fb:(\\*|\\d*) ([\\w\\-_]*)(?: ([\\w\\-_]*))?"),
			Names:  []string{"payload", "type", "subtype"},
			Types:  []rune{'d', 's', 's'},
			Format: "",
			FormatFunc: func(obj *gabs.Container) string {
				if hasValue(obj, "subtype") {
					return "rtcp-fb:%d %s %s"
				} else {
					return "rtcp-fb:%d %s"
				}
			},
		},
		// a=extmap:2 urn:ietf:params:rtp-hdrext:toffset
		// a=extmap:1/recvonly URI-gps-string
		&Rule{
			Name:   "",
			Push:   "ext",
			Reg:    regexp.MustCompile("^extmap:(\\d+)(?:\\/(\\w+))? (\\S*)(?: (\\S*))?"),
			Names:  []string{"value", "direction", "uri", "config"},
			Types:  []rune{'d', 's', 's', 's'},
			Format: "",
			FormatFunc: func(obj *gabs.Container) string {
				ret := "extmap:%d"
				if hasValue(obj, "direction") {
					ret = ret + "/%s"
				} else {
					ret = ret + "%v"
				}
				ret = ret + " %s"
				if hasValue(obj, "config") {
					ret = ret + " %s"
				}
				return ret
			},
		},
		// a=crypto:1 AES_CM_128_HMAC_SHA1_80 inline:PS1uQCVeeCFCanVmcjkpPywjNWhcYD0mXXtxaVBR|2^20|1:32
		&Rule{
			Name:   "",
			Push:   "crypto",
			Reg:    regexp.MustCompile("^crypto:(\\d*) ([\\w_]*) (\\S*)(?: (\\S*))?"),
			Names:  []string{"id", "suite", "config", "sessionConfig"},
			Types:  []rune{'d', 's', 's', 's'},
			Format: "",
			FormatFunc: func(obj *gabs.Container) string {
				if hasValue(obj, "sessionConfig") {
					return "crypto:%d %s %s %s"
				} else {
					return "crypto:%d %s %s"
				}
			},
		},
		// a=setup:actpass
		&Rule{
			Name:   "setup",
			Push:   "",
			Reg:    regexp.MustCompile("^setup:(\\w*)"),
			Names:  []string{},
			Types:  []rune{'s'},
			Format: "setup:%s",
		},
		// a=mid:1
		&Rule{
			Name:   "mid",
			Push:   "",
			Reg:    regexp.MustCompile("^mid:([^\\s]*)"),
			Names:  []string{},
			Types:  []rune{'s'},
			Format: "mid:%s",
		},
		// a=msid:0c8b064d-d807-43b4-b434-f92a889d8587 98178685-d409-46e0-8e16-7ef0db0db64a
		&Rule{
			Name:   "msid",
			Push:   "",
			Reg:    regexp.MustCompile("^msid:(.*)"),
			Names:  []string{},
			Types:  []rune{'s'},
			Format: "msid:%s",
		},
		// a=ptime:20
		&Rule{
			Name:   "ptime",
			Push:   "",
			Reg:    regexp.MustCompile("^ptime:(\\d*)"),
			Names:  []string{},
			Types:  []rune{'d'},
			Format: "ptime:%d",
		},
		// a=maxptime:60
		&Rule{
			Name:   "maxptime",
			Push:   "",
			Reg:    regexp.MustCompile("^maxptime:(\\d*)"),
			Names:  []string{},
			Types:  []rune{'d'},
			Format: "maxptime:%d",
		},
		// a=sendrecv
		&Rule{
			Name:   "direction",
			Push:   "",
			Reg:    regexp.MustCompile("^(sendrecv|recvonly|sendonly|inactive)"),
			Names:  []string{},
			Types:  []rune{'s'},
			Format: "%s",
		},
		// a=ice-lite
		&Rule{
			Name:   "icelite",
			Push:   "",
			Reg:    regexp.MustCompile("^(ice-lite)"),
			Names:  []string{},
			Types:  []rune{'s'},
			Format: "%s",
		},
		// a=ice-ufrag:F7gI
		&Rule{
			Name:   "iceUfrag",
			Push:   "",
			Reg:    regexp.MustCompile("^ice-ufrag:(\\S*)"),
			Names:  []string{},
			Types:  []rune{'s'},
			Format: "ice-ufrag:%s",
		},
		// a=ice-pwd:x9cml/YzichV2+XlhiMu8g
		&Rule{
			Name:   "icePwd",
			Push:   "",
			Reg:    regexp.MustCompile("^ice-pwd:(\\S*)"),
			Names:  []string{},
			Types:  []rune{'s'},
			Format: "ice-pwd:%s",
		},
		// a=fingerprint:SHA-1 00:11:22:33:44:55:66:77:88:99:AA:BB:CC:DD:EE:FF:00:11:22:33
		&Rule{
			Name:   "fingerprint",
			Push:   "",
			Reg:    regexp.MustCompile("^fingerprint:(\\S*) (\\S*)"),
			Names:  []string{"type", "hash"},
			Types:  []rune{'s', 's'},
			Format: "fingerprint:%s %s",
		},
		// a=candidate:0 1 UDP 2113667327 203.0.113.1 54400 typ host
		// a=candidate:1162875081 1 udp 2113937151 192.168.34.75 60017 typ host generation 0 network-id 3 network-cost 10
		// a=candidate:3289912957 2 udp 1845501695 193.84.77.194 60017 typ srflx raddr 192.168.34.75 rport 60017 generation 0 network-id 3 network-cost 10
		// a=candidate:229815620 1 tcp 1518280447 192.168.150.19 60017 typ host tcptype active generation 0 network-id 3 network-cost 10
		// a=candidate:3289912957 2 tcp 1845501695 193.84.77.194 60017 typ srflx raddr 192.168.34.75 rport 60017 tcptype passive generation 0 network-id 3 network-cost 10
		&Rule{
			Name:   "",
			Push:   "candidates",
			Reg:    regexp.MustCompile("^candidate:(\\S*) (\\d*) (\\S*) (\\d*) (\\S*) (\\d*) typ (\\S*)(?: raddr (\\S*) rport (\\d*))?(?: tcptype (\\S*))?(?: generation (\\d*))?(?: network-id (\\d*))?(?: network-cost (\\d*))?"),
			Names:  []string{"foundation", "component", "transport", "priority", "ip", "port", "type", "raddr", "rport", "tcptype", "generation", "network-id", "network-cost"},
			Types:  []rune{'s', 'd', 's', 'd', 's', 'd', 's', 's', 'd', 's', 'd', 'd', 'd', 'd'},
			Format: "",
			FormatFunc: func(obj *gabs.Container) string {
				ret := "candidate:%s %d %s %d %s %d typ %s"
				if hasValue(obj, "raddr") {
					ret = ret + " raddr %s rport %d"
				} else {
					ret = ret + "%v%v"
				}
				if hasValue(obj, "tcptype") {
					ret = ret + " tcptype %s"
				} else {
					ret = ret + "%v"
				}
				if hasValue(obj, "generation") {
					ret = ret + " generation %d"
				}
				if hasValue(obj, "network-id") {
					ret = ret + " network-id %d"
				} else {
					ret = ret + "%v"
				}
				if hasValue(obj, "network-cost") {
					ret = ret + " network-cost %d"
				} else {
					ret = ret + "%v"
				}
				return ret
			},
		},
		// a=end-of-candidates
		&Rule{
			Name:   "endOfCandidates",
			Push:   "",
			Reg:    regexp.MustCompile("^(end-of-candidates)"),
			Names:  []string{},
			Types:  []rune{'s'},
			Format: "%s",
		},
		// a=remote-candidates:1 203.0.113.1 54400 2 203.0.113.1 54401
		&Rule{
			Name:   "remoteCandidates",
			Push:   "",
			Reg:    regexp.MustCompile("^remote-candidates:(.*)"),
			Names:  []string{},
			Types:  []rune{'s'},
			Format: "remote-candidates:%s",
		},
		// a=ice-options:google-ice
		&Rule{
			Name:   "iceOptions",
			Push:   "",
			Reg:    regexp.MustCompile("^ice-options:(\\S*)"),
			Names:  []string{},
			Types:  []rune{'s'},
			Format: "ice-options:%s",
		},
		// a=ssrc:2566107569 cname:t9YU8M1UxTF8Y1A1
		&Rule{
			Name:   "",
			Push:   "ssrcs",
			Reg:    regexp.MustCompile("^ssrc:(\\d*) ([\\w_-]*)(?::(.*))?"),
			Names:  []string{"id", "attribute", "value"},
			Types:  []rune{'d', 's', 's'},
			Format: "",
			FormatFunc: func(obj *gabs.Container) string {
				ret := "ssrc:%d"
				if hasValue(obj, "attribute") {
					ret = ret + " %s"
					if hasValue(obj, "value") {
						ret = ret + ":%s"
					}
				}
				return ret
			},
		},
		// a=ssrc-group:FEC 1 2
		// a=ssrc-group:FEC-FR 3004364195 1080772241
		&Rule{
			Name:   "",
			Push:   "ssrcGroups",
			Reg:    regexp.MustCompile("^ssrc-group:([\x21\x23\x24\x25\x26\x27\x2A\x2B\x2D\x2E\\w]*) (.*)"),
			Names:  []string{"semantics", "ssrcs"},
			Types:  []rune{'s', 's'},
			Format: "ssrc-group:%s %s",
		},
		// a=msid-semantic: WMS Jvlam5X3SX1OP6pn20zWogvaKJz5Hjf9OnlV
		&Rule{
			Name:   "msidSemantic",
			Push:   "",
			Reg:    regexp.MustCompile("^msid-semantic:\\s?(\\w*) (\\S*)"),
			Names:  []string{"semantic", "token"},
			Types:  []rune{'s', 's'},
			Format: "msid-semantic: %s %s",
		},
		// a=group:BUNDLE audio video
		&Rule{
			Name:   "",
			Push:   "groups",
			Reg:    regexp.MustCompile("^group:(\\w*) (.*)"),
			Names:  []string{"type", "mids"},
			Types:  []rune{'s', 's'},
			Format: "group:%s %s",
		},
		// a=rtcp-mux
		&Rule{
			Name:   "rtcpMux",
			Push:   "",
			Reg:    regexp.MustCompile("^(rtcp-mux)"),
			Names:  []string{},
			Types:  []rune{'s'},
			Format: "%s",
		},
		// a=rtcp-rsize
		&Rule{
			Name:   "rtcpRsize",
			Push:   "",
			Reg:    regexp.MustCompile("^(rtcp-rsize)"),
			Names:  []string{},
			Types:  []rune{'s'},
			Format: "%s",
		},
		// a=sctpmap:5000 webrtc-datachannel 1024
		&Rule{
			Name:   "sctpmap",
			Push:   "",
			Reg:    regexp.MustCompile("^sctpmap:(\\d+) (\\S*)(?: (\\d*))?"),
			Names:  []string{"sctpmapNumber", "app", "maxMessageSize"},
			Types:  []rune{'d', 's', 'd'},
			Format: "",
			FormatFunc: func(obj *gabs.Container) string {
				if hasValue(obj, "maxMessageSize") {
					return "sctpmap:%s %s %s"
				} else {
					return "sctpmap:%s %s"
				}
			},
		},
		// a=x-google-flag:conference
		&Rule{
			Name:   "xGoogleFlag",
			Push:   "",
			Reg:    regexp.MustCompile("x-google-flag:([^\\s]*)"),
			Names:  []string{},
			Types:  []rune{'s'},
			Format: "x-google-flag:%s",
		},
		// a=rid:1 send max-width=1280;max-height=720;max-fps=30;depend=0
		&Rule{
			Name:   "",
			Push:   "rids",
			Reg:    regexp.MustCompile("^rid:([\\d\\w]+) (\\w+)(?: ([\\S| ]*))?"),
			Names:  []string{"id", "direction", "params"},
			Types:  []rune{'s', 's', 's'},
			Format: "",
			FormatFunc: func(obj *gabs.Container) string {
				if hasValue(obj, "params") {
					return "rid:%s %s %s"
				} else {
					return "rid:%s %s"
				}
			},
		},
		// a=imageattr:97 send [x=800,y=640,sar=1.1,q=0.6] [x=480,y=320] recv [x=330,y=250]
		// a=imageattr:* send [x=800,y=640] recv *
		// a=imageattr:100 recv [x=320,y=240]
		&Rule{
			Name: "",
			Push: "imageattrs",
			Reg: regexp.MustCompile(
				"^imageattr:(\\d+|\\*)" +
					"[\\s\\t]+(send|recv)[\\s\\t]+(\\*|\\[\\S+\\](?:[\\s\\t]+\\[\\S+\\])*)" +
					"(?:[\\s\\t]+(recv|send)[\\s\\t]+(\\*|\\[\\S+\\](?:[\\s\\t]+\\[\\S+\\])*))?"),
			Names: []string{
				"pt",
				"dir1",
				"attrs1",
				"dir2",
				"attrs2",
			},
			Types:  []rune{'s', 's', 's', 's', 's'},
			Format: "",
			FormatFunc: func(obj *gabs.Container) string {
				ret := "imageattr:%s %s %s"
				if hasValue(obj, "dir2") {
					ret = ret + " %s %s"
				}
				return ret
			},
		},
		// a=simulcast:send 1,2,3;~4,~5 recv 6;~7,~8
		// a=simulcast:recv 1;4,5 send 6;7
		&Rule{
			Name: "simulcast",
			Push: "",
			Reg: regexp.MustCompile(
				"^simulcast:" +
					"(send|recv) ([a-zA-Z0-9\\-_~;,]+)" +
					"(?:\\s?(send|recv) ([a-zA-Z0-9\\-_~;,]+))?" +
					"$"),
			Names:  []string{"dir1", "list1", "dir2", "list2"},
			Types:  []rune{'s', 's', 's', 's'},
			Format: "",
			FormatFunc: func(obj *gabs.Container) string {
				ret := "simulcast:%s %s"
				if hasValue(obj, "dir2") {
					ret = ret + " %s %s"
				}
				return ret
			},
		},
		// Old simulcast draft 03 (implemented by Firefox).
		// https://tools.ietf.org/html/draft-ietf-mmusic-sdp-simulcast-03
		// a=simulcast: recv pt=97;98 send pt=97
		// a=simulcast: send rid=5;6;7 paused=6,7
		&Rule{
			Name:   "simulcast_03",
			Push:   "",
			Reg:    regexp.MustCompile("^simulcast:[\\s\\t]+([\\S+\\s\t]+)$"),
			Names:  []string{"value"},
			Types:  []rune{'s'},
			Format: "simulcast: %s",
		},
		// a=framerate:25
		// a=framerate:29.97
		&Rule{
			Name:   "framerate",
			Push:   "",
			Reg:    regexp.MustCompile("^framerate:(\\d+(?:$|\\.\\d+))"),
			Names:  []string{},
			Types:  []rune{'f'},
			Format: "framerate:%s",
		},
		// a=source-filter: incl IN IP4 239.5.2.31 10.1.15.5
		&Rule{
			Name:   "sourceFilter",
			Push:   "",
			Reg:    regexp.MustCompile("^source-filter:[\\s\\t]+(excl|incl) (\\S*) (IP4|IP6|\\*) (\\S*) (.*)"),
			Names:  []string{"filterMode", "netType", "addressTypes", "destAddress", "srcList"},
			Types:  []rune{'s', 's', 's', 's', 's'},
			Format: "source-filter: %s %s %s %s %s",
		},
		// Any a= that we don't understand is kepts verbatim on media.invalid.
		&Rule{
			Name:   "",
			Push:   "invalid",
			Reg:    regexp.MustCompile("(.*)"),
			Names:  []string{"value"},
			Types:  []rune{'s'},
			Format: "%s",
		},
	},
}
