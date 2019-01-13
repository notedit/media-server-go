package mediaserver

import (
	"testing"

	"github.com/notedit/media-server-go/sdp"
)

func Test_TransportCreate(t *testing.T) {

	EnableLog(true)
	endpoint := NewEndpoint("127.0.0.1")

	iceInfo := sdp.GenerateICEInfo(true)
	dtlsInfo := sdp.NewDTLSInfo(sdp.SETUPACTPASS, "sha-256", "F2:AA:0E:C3:22:59:5E:14:95:69:92:3D:13:B4:84:24:2C:C2:A2:C0:3E:FD:34:8E:5E:EA:6F:AF:52:CE:E6:0F")
	sdpInfo := sdp.NewSDPInfo()
	sdpInfo.SetICE(iceInfo)
	sdpInfo.SetDTLS(dtlsInfo)
	transport := endpoint.CreateTransport(sdpInfo, nil)

	if transport == nil {
		t.Error("can not create transport")
	}
	t.Log("yes")
}

func Test_CreateIncomingTrack(t *testing.T) {

	endpoint := NewEndpoint("127.0.0.1")
	iceInfo := sdp.ICEInfoGenerate(true)
	dtlsInfo := sdp.NewDTLSInfo(sdp.SETUPACTPASS, "sha-256", "F2:AA:0E:C3:22:59:5E:14:95:69:92:3D:13:B4:84:24:2C:C2:A2:C0:3E:FD:34:8E:5E:EA:6F:AF:52:CE:E6:0F")
	sdpInfo := sdp.NewSDPInfo()
	sdpInfo.SetICE(iceInfo)
	sdpInfo.SetDTLS(dtlsInfo)

	transport := endpoint.CreateTransport(sdpInfo, nil)

	incomingTrack := transport.CreateIncomingStreamTrack("audio", "audiotrack", map[string]uint{})

	if incomingTrack.GetID() != "audiotrack" {
		t.Error("create incoming track error")
	}
	t.Log("yes")
}

func Test_CreateOutgoingTrack(t *testing.T) {

	endpoint := NewEndpoint("127.0.0.1")
	iceInfo := sdp.ICEInfoGenerate(true)
	dtlsInfo := sdp.NewDTLSInfo(sdp.SETUPACTPASS, "sha-256", "F2:AA:0E:C3:22:59:5E:14:95:69:92:3D:13:B4:84:24:2C:C2:A2:C0:3E:FD:34:8E:5E:EA:6F:AF:52:CE:E6:0F")
	sdpInfo := sdp.NewSDPInfo()
	sdpInfo.SetICE(iceInfo)
	sdpInfo.SetDTLS(dtlsInfo)

	transport := endpoint.CreateTransport(sdpInfo, nil)
	outgoingTrack := transport.CreateOutgoingStreamTrack("video", "videotrack", map[string]uint{})

	if outgoingTrack.GetID() != "videotrack" {
		t.Error("create outgoing track error")
	}
}

func Test_TransportStop(t *testing.T) {

	EnableLog(true)
	endpoint := NewEndpoint("127.0.0.1")

	iceInfo := sdp.ICEInfoGenerate(true)
	dtlsInfo := sdp.NewDTLSInfo(sdp.SETUPACTPASS, "sha-256", "F2:AA:0E:C3:22:59:5E:14:95:69:92:3D:13:B4:84:24:2C:C2:A2:C0:3E:FD:34:8E:5E:EA:6F:AF:52:CE:E6:0F")
	sdpInfo := sdp.NewSDPInfo()
	sdpInfo.SetICE(iceInfo)
	sdpInfo.SetDTLS(dtlsInfo)

	transport := endpoint.CreateTransport(sdpInfo, nil)

	transport.OnStop(func() {
		t.Log("transport stopped")
	})
}
