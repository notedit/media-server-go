package mediaserver

import "github.com/chuckpreslar/emission"

type ActiveSpeakerDetector struct {
	tracks   map[string]*IncomingStreamTrack
	detector ActiveSpeakerDetectorFacade
	*emission.Emitter
}

func NewActiveSpeakerDetector() *ActiveSpeakerDetector {

	detector := new(ActiveSpeakerDetector)
	detector.tracks = map[string]*IncomingStreamTrack{}
	detector.detector = NewActiveSpeakerDetectorFacade()

	// todo
	return detector
}
