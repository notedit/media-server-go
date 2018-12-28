package mediaserver

import (
	"github.com/chuckpreslar/emission"
	native "github.com/notedit/media-server-go/wrapper"
)

type ActiveSpeakerDetector struct {
	tracks   map[string]*IncomingStreamTrack
	detector native.ActiveSpeakerDetectorFacade
	*emission.Emitter
}

func NewActiveSpeakerDetector() *ActiveSpeakerDetector {

	detector := &ActiveSpeakerDetector{}
	detector.tracks = map[string]*IncomingStreamTrack{}
	detector.detector = native.NewActiveSpeakerDetectorFacade()
	detector.Emitter = emission.NewEmitter()

	// todo
	return detector
}

func (a *ActiveSpeakerDetector) SetMinChangePeriod(minChangePeriod uint) {
	a.detector.SetMinChangePeriod(minChangePeriod)
}

func (a *ActiveSpeakerDetector) AddSpeaker(track *IncomingStreamTrack) {

}

func (a *ActiveSpeakerDetector) RemoveSpeaker(track *IncomingStreamTrack) {

}

func (a *ActiveSpeakerDetector) Stop() {

}
