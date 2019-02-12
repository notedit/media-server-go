package mediaserver

/*
#include <stdlib.h>
*/
import "C"
import (
	"fmt"
	"strings"

	"github.com/gofrs/uuid"
	"github.com/notedit/media-server-go/sdp"
	native "github.com/notedit/media-server-go/wrapper"
)

type RawStreamerSession struct {
	id       string
	incoming *IncomingStreamTrack
	session  native.RawRTPSessionFacade
}

func NewRawStreamerSession(media *sdp.MediaInfo) *RawStreamerSession {

	streamerSession := &RawStreamerSession{}
	var mediaType native.MediaFrameType = 0
	if strings.ToLower(media.GetType()) == "video" {
		mediaType = 1
	}
	session := native.NewRawRTPSessionFacade(mediaType)
	streamerSession.id = uuid.Must(uuid.NewV4()).String()

	properties := native.NewProperties()
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
	native.DeleteProperties(properties)
	streamerSession.session = session
	streamerSession.incoming = newIncomingStreamTrack(media.GetType(), media.GetType(), native.RTPSessionToReceiver(session), map[string]native.RTPIncomingSourceGroup{"": session.GetIncomingSourceGroup()})

	streamerSession.incoming.OnStop(func() {
		streamerSession.incoming = nil
	})

	return streamerSession
}

func (s *RawStreamerSession) GetID() string {
	return s.id
}

func (s *RawStreamerSession) GetIncomingStreamTrack() *IncomingStreamTrack {
	return s.incoming
}

func (s *RawStreamerSession) Push(rtp []byte) {
	if rtp == nil || len(rtp) == 0 {
		return
	}
	s.session.OnRTPPacket(&rtp[0], len(rtp))
}

func (s *RawStreamerSession) Stop() {

	if s.session == nil {
		return
	}

	if s.incoming != nil {
		s.incoming.Stop()
	}

	s.session.End()

	native.DeleteRawRTPSessionFacade(s.session)

	s.session = nil
}
