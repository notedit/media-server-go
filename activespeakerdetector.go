package mediaserver

import (
	"sync"

	native "github.com/notedit/media-server-go/wrapper"
)

type activeTrackListener interface {
	native.ActiveTrackListener
	deleteActiveTrackListener()
}

type goActiveTrackListener struct {
	native.ActiveTrackListener
}

func (a *goActiveTrackListener) deleteActiveTrackListener() {
	native.DeleteDirectorActiveTrackListener(a.ActiveTrackListener)
}

type overwrittenActiveTrackListener struct {
	p        native.ActiveTrackListener
	detector *ActiveSpeakerDetector
}

func (p *overwrittenActiveTrackListener) OnActiveTrackchanged(ssrc uint) {
	if p.detector != nil && p.detector.listener != nil {
		//p.detector.listener(ssrc)
		if p.detector.tracks[ssrc] != nil {
			p.detector.listener(p.detector.tracks[ssrc])
		}
	}
}

// ActiveSpeakerDetector detector the spkeaking track
type ActiveSpeakerDetector struct {
	tracks              map[uint]*IncomingStreamTrack
	detector            native.ActiveSpeakerDetectorFacade
	listener            ActiveDetectorListener
	activeTrackListener *goActiveTrackListener
	sync.Mutex
}

// ActiveDetectorListener listener
type ActiveDetectorListener func(*IncomingStreamTrack)

// NewActiveSpeakerDetector  create new  active speaker detector
func NewActiveSpeakerDetector(listener ActiveDetectorListener) *ActiveSpeakerDetector {

	detector := &ActiveSpeakerDetector{}
	detector.tracks = map[uint]*IncomingStreamTrack{}

	activeTrackListener := &overwrittenActiveTrackListener{
		detector: detector,
	}
	p := native.NewDirectorActiveTrackListener(activeTrackListener)
	activeTrackListener.p = p

	detector.activeTrackListener = &goActiveTrackListener{ActiveTrackListener: p}
	detector.detector = native.NewActiveSpeakerDetectorFacade(detector.activeTrackListener)

	detector.listener = listener

	return detector
}

// SetMinChangePeriod  set min  period change callback  in ms
func (a *ActiveSpeakerDetector) SetMinChangePeriod(minChangePeriod uint) {
	a.detector.SetMinChangePeriod(minChangePeriod)
}

// SetMaxAccumulatedScore  maximux activity score accumulated by an speaker
func (a *ActiveSpeakerDetector) SetMaxAccumulatedScore(maxAcummulatedScore uint64) {
	a.detector.SetMaxAccumulatedScore(maxAcummulatedScore)
}

// SetNoiseGatingThreshold Minimum db level to not be considered as muted
func (a *ActiveSpeakerDetector) SetNoiseGatingThreshold(noiseGatingThreshold byte) {
	a.detector.SetNoiseGatingThreshold(noiseGatingThreshold)
}

//SetMinActivationScore  Set minimum activation score to be electible as active speaker
func (a *ActiveSpeakerDetector) SetMinActivationScore(minActivationScore uint) {
	a.detector.SetMinActivationScore(minActivationScore)
}

// AddTrack  add incoming track into detector
func (a *ActiveSpeakerDetector) AddTrack(track *IncomingStreamTrack) {

	// We should make sure this source is the main source
	source := track.GetFirstEncoding().GetSource()
	if source == nil {
		return
	}
	ssrc := source.GetMedia().GetSsrc()

	a.Lock()
	a.tracks[ssrc] = track
	a.Unlock()

	a.detector.AddIncomingSourceGroup(source)
}

// RemoveTrack  remove incoming track from detector
func (a *ActiveSpeakerDetector) RemoveTrack(track *IncomingStreamTrack) {
	source := track.GetFirstEncoding().GetSource()
	if source == nil {
		return
	}
	a.Lock()
	delete(a.tracks, source.GetMedia().GetSsrc())
	a.Unlock()
	a.detector.RemoveIncomingSourceGroup(source)
}

// Stop stop the detector
func (a *ActiveSpeakerDetector) Stop() {

	for _, track := range a.tracks {
		encoding := track.GetFirstEncoding()
		if encoding != nil {
			source := encoding.GetSource()
			if source != nil {
				a.detector.RemoveIncomingSourceGroup(source)
			}

		}
	}

	if a.activeTrackListener != nil {
		a.activeTrackListener.deleteActiveTrackListener()
	}

	a.tracks = nil

	if a.detector != nil {
		native.DeleteActiveSpeakerDetectorFacade(a.detector)
		a.detector = nil
	}

}
