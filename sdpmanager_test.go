package mediaserver

import (
	"testing"

	"github.com/notedit/sdp"
)

var Capabilities = map[string]*sdp.Capability{
	"audio": &sdp.Capability{
		Codecs: []string{"opus"},
	},
	"video": &sdp.Capability{
		Codecs: []string{"h264"},
		Rtx:    true,
		Rtcpfbs: []*sdp.RtcpFeedback{
			&sdp.RtcpFeedback{
				ID: "goog-remb",
			},
			&sdp.RtcpFeedback{
				ID: "transport-cc",
			},
			&sdp.RtcpFeedback{
				ID:     "ccm",
				Params: []string{"fir"},
			},
			&sdp.RtcpFeedback{
				ID:     "nack",
				Params: []string{"pli"},
			},
		},
		Extensions: []string{
			"urn:3gpp:video-orientation",
			"http://www.ietf.org/id/draft-holmer-rmcat-transport-wide-cc-extensions-01",
			"http://www.webrtc.org/experiments/rtp-hdrext/abs-send-time",
			"urn:ietf:params:rtp-hdrext:toffse",
			"urn:ietf:params:rtp-hdrext:sdes:rtp-stream-id",
			"urn:ietf:params:rtp-hdrext:sdes:mid",
		},
	},
}

func TestSDPManagerCreate(t *testing.T) {

	endpoint := NewEndpoint("127.0.0.1")
	sdpManager := endpoint.CreateSDPManager("unified-plan", Capabilities)
	if sdpManager.GetState() != "initial" {
		t.Error("sdpmanager create error")
	}
}

func TestSDPManagerOfferAnswer(t *testing.T) {

	endpoint1 := NewEndpoint("127.0.0.1")
	endpoint2 := NewEndpoint("127.0.0.1")

	sdpmanager1 := endpoint1.CreateSDPManager("unified-plan", Capabilities)
	sdpmanager2 := endpoint2.CreateSDPManager("unified-plan", Capabilities)

	offer, _ := sdpmanager1.CreateLocalDescription()

	if sdpmanager1.GetState() != "local-offer" {
		t.Error("create local sdp error")
	}

	sdpmanager2.ProcessRemoteDescription(offer.String())

	if sdpmanager2.GetState() != "remote-offer" {
		t.Error("process remote sdp error")
	}

	answer, _ := sdpmanager2.CreateLocalDescription()

	if sdpmanager2.GetState() != "stable" {
		t.Error("state error ")
	}

	sdpmanager1.ProcessRemoteDescription(answer.String())

	if sdpmanager1.GetState() != "stable" {
		t.Error("process remote sdp error")
	}
}
