package mediaserver

import (
	"errors"

	"github.com/chuckpreslar/emission"
)

type BitrateTraversal string

const (
	TraversalDefault               BitrateTraversal = "default"
	TraversalSpatialTemporal       BitrateTraversal = "spatial-temporal"
	TraversalZigZagSpatialTemporal BitrateTraversal = "zig-zag-spatial-temporal"
	TraversalZigZagTemporalSpatial BitrateTraversal = "zig-zag-temporal-spatial"
)

type Transponder struct {
	muted              bool
	track              *IncomingStreamTrack
	transponder        RTPStreamTransponderFacade
	encodingId         string
	spatialLayerId     int
	temporalLayerId    int
	maxSpatialLayerId  int
	maxTemporalLayerId int
	*emission.Emitter
}

func NewTransponder(transponderFacade RTPStreamTransponderFacade) *Transponder {
	transponder := new(Transponder)
	transponder.muted = false

	transponder.transponder = transponderFacade
	transponder.spatialLayerId = MaxLayerId
	transponder.temporalLayerId = MaxLayerId
	transponder.maxSpatialLayerId = MaxLayerId
	transponder.maxTemporalLayerId = MaxLayerId
	transponder.Emitter = emission.NewEmitter()

	return transponder
}

func (t *Transponder) SetIncomingTrack(incomingTrack *IncomingStreamTrack) error {

	if t.transponder == nil {
		return errors.New("Transponder is already closed")
	}

	if incomingTrack == nil {
		return errors.New("Track can not be nil")
	}

	if t.track != nil {
		t.track.Off("stopped", t.onAttachedTrackStopped)
		t.track.Detached()
	}

	t.track = incomingTrack

	// we need make sure first encoding id not nil
	// todo check
	// get first encoding
	encoding := t.track.GetFirstEncoding()

	t.transponder.SetIncoming(encoding.GetSource(), incomingTrack.receiver)

	t.encodingId = encoding.GetID()

	t.spatialLayerId = MaxLayerId
	t.temporalLayerId = MaxLayerId
	t.maxSpatialLayerId = MaxLayerId
	t.maxTemporalLayerId = MaxLayerId

	t.track.Once("stopped", t.onAttachedTrackStopped)

	t.track.Attached()

	return nil
}

func (t *Transponder) GetIncomingTrack() *IncomingStreamTrack {
	return t.track
}

// GetAvailableLayers   Get available encodings and layers
func (t *Transponder) GetAvailableLayers() {
	// todo
}

func (t *Transponder) IsMuted() bool {

	return t.muted
}

func (t *Transponder) Mute(muting bool) {

	if t.muted != muting {
		t.muted = muting
		if t.transponder != nil {
			t.transponder.Mute(muting)
		}
		t.EmitSync("muted")
	}
}

func (t *Transponder) SetTargetBitrate(bitrate int, traversal BitrateTraversal, strict bool) int {

	if t.track == nil {
		return 0
	}

	// current := -1
	// encodingId := 0
	// spatialLayerId := MaxLayerId
	// temporalLayerId := MaxLayerId

	// min := math.MaxInt32
	// encodingIdMin := 0

	// spatialLayerIdMin := MaxLayerId
	// temporalLayerIdMin := MaxLayerId

	// ordering := false

	// todo

	return 0
}

func (t *Transponder) SelectEncoding() {

	// todo
}

func (t *Transponder) GetSelectedEncoding() string {

	return t.encodingId
}

func (t *Transponder) GetSelectedSpatialLayerId() int {

	return t.spatialLayerId
}

func (t *Transponder) GetSelectedTemporalLayerId() int {

	return t.temporalLayerId
}

// SelectLayer Select SVC temporatl and spatial layers. Only available for VP9 media.
func (t *Transponder) SelectLayer(spatialLayerId, temporalLayerId int) {

	spatialLayerId = Min(spatialLayerId, t.maxSpatialLayerId)
	temporalLayerId = Min(temporalLayerId, t.maxTemporalLayerId)

	if t.spatialLayerId == spatialLayerId && t.temporalLayerId == temporalLayerId {
		return
	}

	t.transponder.SelectLayer(spatialLayerId, temporalLayerId)

	t.spatialLayerId = spatialLayerId
	t.temporalLayerId = temporalLayerId
}

func (t *Transponder) SetMaximumLayers(maxSpatialLayerId, maxTemporalLayerId int) {

	if maxSpatialLayerId < 0 || maxTemporalLayerId < 0 {
		return
	}

	t.maxSpatialLayerId = maxSpatialLayerId
	t.maxTemporalLayerId = maxTemporalLayerId

}

func (t *Transponder) Stop() {

	if t.transponder == nil {
		return
	}

	if t.track != nil {
		t.track.Off("stopped", t.onAttachedTrackStopped)
		t.track.Detached()
	}

	t.transponder.Close()

	DeleteRTPStreamTransponderFacade(t.transponder)

	t.EmitSync("stopped")

	t.transponder = nil

	t.track = nil
}

func (t *Transponder) onAttachedTrackStopped() {
	t.Stop()
}
