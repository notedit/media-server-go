package mediaserver

type ActiveSpeakerDetector struct {
	tracks   map[string]*IncomingStreamTrack
	detector ActiveSpeakerDetectorFacade
}

func NewActiveSpeakerDetector() *ActiveSpeakerDetector {

	return nil
}
