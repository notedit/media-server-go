package mediaserver

import (
	"errors"
	"math"
	"sort"

	native "github.com/notedit/media-server-go/wrapper"
)

type BitrateTraversal string

const (
	TraversalDefault               BitrateTraversal = "default"
	TraversalSpatialTemporal       BitrateTraversal = "spatial-temporal"
	TraversalTemporalSpatial       BitrateTraversal = "temporal-spatial"
	TraversalZigZagSpatialTemporal BitrateTraversal = "zig-zag-spatial-temporal"
	TraversalZigZagTemporalSpatial BitrateTraversal = "zig-zag-temporal-spatial"
)

// Transponder
type Transponder struct {
	muted              bool
	track              *IncomingStreamTrack
	transponder        native.RTPStreamTransponderFacade
	encodingId         string
	spatialLayerId     int
	temporalLayerId    int
	maxSpatialLayerId  int
	maxTemporalLayerId int
	onMuteListeners    []func(bool)
	onStopListeners    []func()
}

func NewTransponder(transponderFacade native.RTPStreamTransponderFacade) *Transponder {
	transponder := new(Transponder)
	transponder.muted = false

	transponder.transponder = transponderFacade
	transponder.spatialLayerId = MaxLayerId
	transponder.temporalLayerId = MaxLayerId
	transponder.maxSpatialLayerId = MaxLayerId
	transponder.maxTemporalLayerId = MaxLayerId

	transponder.onMuteListeners = make([]func(bool), 0)
	transponder.onStopListeners = make([]func(), 0)

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
		t.track.Detached()
	}

	t.track = incomingTrack

	// we need make sure first encoding id not nil
	// todo check
	// get first encoding
	encoding := t.track.GetFirstEncoding()
	if encoding == nil {
		panic("encoding is nil")
	}

	t.transponder.SetIncoming(encoding.GetSource(), incomingTrack.receiver)

	t.encodingId = encoding.GetID()

	t.spatialLayerId = MaxLayerId
	t.temporalLayerId = MaxLayerId
	t.maxSpatialLayerId = MaxLayerId
	t.maxTemporalLayerId = MaxLayerId

	t.track.OnStop(t.onAttachedTrackStopped)

	t.track.Attached()

	return nil
}

func (t *Transponder) GetIncomingTrack() *IncomingStreamTrack {
	return t.track
}

// GetAvailableLayers   Get available encodings and layers
func (t *Transponder) GetAvailableLayers() *ActiveLayersInfo {
	if t.track != nil {
		return t.track.GetActiveLayers()
	}
	return nil
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

		for _, mutefunc := range t.onMuteListeners {
			mutefunc(muting)
		}
	}
}

func getSpatialLayerId(layer *Layer) int {
	if layer.SpatialLayerId != MaxLayerId {
		return layer.SpatialLayerId
	}
	return layer.SimulcastIdx
}

func (t *Transponder) SetTargetBitrate(bitrate uint, traversal BitrateTraversal, strict bool) uint {

	if t.track == nil {
		return 0
	}

	var current uint
	var encodingId string
	var encodingIdMin string

	spatialLayerId := MaxLayerId
	temporalLayerId := MaxLayerId

	min := uint(math.MaxInt32)

	spatialLayerIdMin := MaxLayerId
	temporalLayerIdMin := MaxLayerId

	var orderfunc func(i int, j int) bool
	layers := t.track.GetActiveLayers().Layers

	if len(layers) == 0 {
		t.Mute(false)
		return 0
	}

	switch traversal {
	case TraversalSpatialTemporal:
		orderfunc = func(i, j int) bool {
			a := layers[i]
			b := layers[j]

			ret1 := getSpatialLayerId(b)*MaxLayerId + b.TemporalLayerId
			ret2 := getSpatialLayerId(a)*MaxLayerId + a.TemporalLayerId
			if ret1-ret2 < 0 {
				return true
			}
			return false
		}
		break
	case TraversalZigZagSpatialTemporal:
		orderfunc = func(i, j int) bool {
			a := layers[i]
			b := layers[j]

			ret1 := (getSpatialLayerId(b)+b.TemporalLayerId+1)*MaxLayerId - b.TemporalLayerId
			ret2 := (getSpatialLayerId(a)+a.TemporalLayerId+1)*MaxLayerId - a.TemporalLayerId
			if ret1-ret2 < 0 {
				return true
			}
			return false
		}
		break
	case TraversalTemporalSpatial:
		orderfunc = func(i, j int) bool {
			a := layers[i]
			b := layers[j]

			ret1 := b.TemporalLayerId*MaxLayerId + getSpatialLayerId(b)
			ret2 := a.TemporalLayerId*MaxLayerId + getSpatialLayerId(a)
			if ret1-ret2 < 0 {
				return true
			}
			return false
		}
		break
	case TraversalZigZagTemporalSpatial:
		orderfunc = func(i, j int) bool {
			a := layers[i]
			b := layers[j]

			ret1 := (getSpatialLayerId(b)+b.TemporalLayerId+1)*MaxLayerId - getSpatialLayerId(b)
			ret2 := (getSpatialLayerId(a)+a.TemporalLayerId+1)*MaxLayerId - getSpatialLayerId(a)
			if ret1-ret2 < 0 {
				return true
			}
			return false
		}
		break
	default:
		break
	}

	if orderfunc != nil {
		sort.SliceStable(layers, orderfunc)
	}

	for _, layer := range layers {

		if layer.Bitrate <= bitrate && layer.Bitrate > current && t.maxSpatialLayerId >= layer.SpatialLayerId && t.maxSpatialLayerId >= layer.TemporalLayerId {
			encodingId = layer.EncodingId
			spatialLayerId = layer.SpatialLayerId
			temporalLayerId = layer.TemporalLayerId
			current = layer.Bitrate
			break
		}

		if layer.Bitrate > 0 && layer.Bitrate < min && t.maxSpatialLayerId >= layer.SpatialLayerId {
			encodingIdMin = layer.EncodingId
			spatialLayerIdMin = layer.SpatialLayerId
			temporalLayerIdMin = layer.TemporalLayerId
			min = layer.Bitrate
		}
	}

	// Check if we have been able to find a layer that matched the target bitrate
	if current <= 0 {

		if strict == false {
			t.Mute(false)
			t.SelectEncoding(encodingIdMin)
			t.SelectLayer(spatialLayerIdMin, temporalLayerIdMin)
			return min
		} else {
			t.Mute(true)
			return 0
		}
	}
	t.Mute(false)
	t.SelectEncoding(encodingId)
	t.SelectLayer(spatialLayerId, temporalLayerId)
	return current
}

// SelectEncoding by id
func (t *Transponder) SelectEncoding(encodingId string) {

	if t.encodingId == encodingId {
		return
	}
	encoding := t.track.GetEncoding(encodingId)
	if encoding == nil {
		return
	}
	t.transponder.SetIncoming(encoding.GetSource(), t.track.receiver)
	t.encodingId = encodingId
}

// GetSelectedEncoding get selected encoding id
func (t *Transponder) GetSelectedEncoding() string {
	return t.encodingId
}

// GetSelectedSpatialLayerId  return int
func (t *Transponder) GetSelectedSpatialLayerId() int {
	return t.spatialLayerId
}

// GetSelectedTemporalLayerId  return int
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

// OnMute register mute listener
func (t *Transponder) OnMute(listener func(bool)) {
	t.onMuteListeners = append(t.onMuteListeners, listener)
}

// OnStop register stop listener
func (t *Transponder) OnStop(stop func()) {
	t.onStopListeners = append(t.onStopListeners, stop)
}

// Stop stop this transponder
func (t *Transponder) Stop() {

	if t.transponder == nil {
		return
	}

	if t.track != nil {
		t.track.Detached()
	}

	t.transponder.Close()

	native.DeleteRTPStreamTransponderFacade(t.transponder)

	for _, stopFunc := range t.onStopListeners {
		stopFunc()
	}

	t.transponder = nil

	t.track = nil
}

func (t *Transponder) onAttachedTrackStopped() {
	t.Stop()
}
