package mediaserver

/*
#include <stdlib.h>
#include <stdlib.h>
*/
import "C"
import (
	"fmt"
	"strings"
	"unsafe"

	"github.com/chuckpreslar/emission"
	"github.com/notedit/media-server-go/sdp"
	"github.com/satori/go.uuid"
)

type RawRTPStreamerSession struct {
	id       string
	incoming *IncomingStreamTrack
	session  RawRTPSessionFacade
	*emission.Emitter
}

func NewRawRTPStreamerSession(media *sdp.MediaInfo) *RawRTPStreamerSession {

	streamerSession := &RawRTPStreamerSession{}
	var mediaType MediaFrameType = 0
	if strings.ToLower(media.GetType()) == "video" {
		mediaType = 1
	}
	session := NewRawRTPSessionFacade(mediaType)
	streamerSession.id = uuid.Must(uuid.NewV4()).String()

	properties := NewProperties()
	if media != nil {
		num := 0
		for _, codec := range media.GetCodecs() {
			item := fmt.Sprintf("codecs.%d", num)
			properties.SetProperty(item+".codec", codec.GetCodec())
			properties.SetProperty(item+".pt", codec.GetType())
			if codec.HasRTX() {
				properties.SetProperty(item+".rtx", codec.GetRTX())
			}
			num = num + 1
		}
		properties.SetProperty("codecs.length", num)
	}

	session.Init(properties)
	DeleteProperties(properties)
	streamerSession.session = session
	streamerSession.incoming = newIncomingStreamTrack(media.GetType(), media.GetType(), nil, map[string]RTPIncomingSourceGroup{"": session.GetIncomingSourceGroup()})

	streamerSession.incoming.Once("stopped", func() {
		streamerSession.incoming = nil
	})

	return streamerSession
}

func (s *RawRTPStreamerSession) GetID() string {
	return s.id
}

func (s *RawRTPStreamerSession) GetIncomingStreamTrack() *IncomingStreamTrack {
	return s.incoming
}

func (s *RawRTPStreamerSession) Push(rtp []byte) {
	b := C.CBytes(rtp)
	defer C.free(unsafe.Pointer(b))
	s.session.OnRTPPacket((*byte)(b), len(rtp))
}

func (s *RawRTPStreamerSession) Stop() {

	if s.session == nil {
		return
	}

	if s.incoming != nil {
		s.incoming.Stop()
	}

	s.session.End()

	DeleteRawRTPSessionFacade(s.session)

	s.EmitSync("stopped")

}
