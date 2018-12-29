package mediaserver

import (
	native "github.com/notedit/media-server-go/wrapper"
)

type ActiveSpeakerDetector struct {
	tracks   map[uint]*IncomingStreamTrack
	detector native.ActiveSpeakerDetectorFacade
}

// TODO, add active callback

func NewActiveSpeakerDetector() *ActiveSpeakerDetector {

	detector := &ActiveSpeakerDetector{}
	detector.tracks = map[uint]*IncomingStreamTrack{}
	detector.detector = native.NewActiveSpeakerDetectorFacade()

	return detector
}

func (a *ActiveSpeakerDetector) SetMinChangePeriod(minChangePeriod uint) {
	a.detector.SetMinChangePeriod(minChangePeriod)
}

func (a *ActiveSpeakerDetector) AddSpeaker(track *IncomingStreamTrack) {

	// We should make sure this source is the main source
	source := track.GetFirstEncoding().GetSource()
	if source == nil {
		return
	}
	ssrc := source.GetMedia().GetSsrc()
	a.tracks[ssrc] = track
	a.detector.AddIncomingSourceGroup(source)
}

func (a *ActiveSpeakerDetector) RemoveSpeaker(track *IncomingStreamTrack) {
	source := track.GetFirstEncoding().GetSource()
	if source == nil {
		return
	}
	delete(a.tracks, source.GetMedia().GetSsrc())
	a.detector.RemoveIncomingSourceGroup(source)
}

func (a *ActiveSpeakerDetector) Stop() {

	for _, track := range a.tracks {
		source := track.GetFirstEncoding().GetSource()
		a.detector.RemoveIncomingSourceGroup(source)
	}

	a.tracks = map[uint]*IncomingStreamTrack{}

	native.DeleteActiveSpeakerDetectorFacade(a.detector)
}
