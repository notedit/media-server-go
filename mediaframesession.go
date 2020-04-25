package mediaserver

import (
	"fmt"
	native "github.com/notedit/media-server-go/wrapper"
	"github.com/notedit/sdp"
	"strings"
)

type MediaFrameSession struct {
	incoming *IncomingStreamTrack
	session  native.MediaFrameSessionFacade
}

// NewMediaFrameSession create media frame session
func NewMediaFrameSession(media *sdp.MediaInfo) *MediaFrameSession {

	mediaSession := &MediaFrameSession{}
	var mediaType native.MediaFrameType = 0
	if strings.ToLower(media.GetType()) == "video" {
		mediaType = 1
	}

	session := native.NewMediaFrameSessionFacade(mediaType)

	properties := native.NewPropertiesFacade()
	if media != nil {
		num := 0
		for _, codec := range media.GetCodecs() {
			item := fmt.Sprintf("codecs.%d", num)
			properties.SetPropertyStr(item+".codec", codec.GetCodec())
			properties.SetPropertyInt(item+".pt", codec.GetType())
			num = num + 1
		}
		properties.SetPropertyInt("codecs.length", num)
	}

	session.Init(properties)
	native.DeletePropertiesFacade(properties)
	mediaSession.session = session
	mediaSession.incoming = NewIncomingStreamTrack(media.GetType(), media.GetType(), native.RTPSessionToReceiver(session), map[string]native.RTPIncomingSourceGroup{"": session.GetIncomingSourceGroup()})

	return mediaSession
}

// GetIncomingStreamTrack get incoming stream track
func (s *MediaFrameSession) GetIncomingStreamTrack() *IncomingStreamTrack {
	return s.incoming
}

// Push push raw media frame
func (s *MediaFrameSession) Push(rtp []byte) {
	if rtp == nil || len(rtp) == 0 {
		return
	}
	fmt.Println("push buffer 22222", len(rtp))
	s.session.OnRTPPacket(&rtp[0], len(rtp))
}

// Stop stop this
func (s *MediaFrameSession) Stop() {

	if s.session == nil {
		return
	}

	if s.incoming != nil {
		s.incoming.Stop()
	}

	s.session.End()

	native.DeleteMediaFrameSessionFacade(s.session)

	s.session = nil
}
