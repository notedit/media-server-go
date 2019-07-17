package mediaserver

import (
	"sort"
	"strconv"
	"time"

	native "github.com/notedit/media-server-go/wrapper"
	"github.com/notedit/sdp"
)

// Layer info
type Layer struct {
	// EncodingId str
	EncodingId string
	// SpatialLayerId int
	SpatialLayerId int
	// TemporalLayerId  int
	TemporalLayerId int
	// SimulcastIdx int
	SimulcastIdx int
	// TotalBytes uint
	TotalBytes uint
	// NumPackets uint
	NumPackets uint
	// Bitrate  uint
	Bitrate uint
}

// Encoding info
type Encoding struct {
	id           string
	source       native.RTPIncomingSourceGroup
	depacketizer native.StreamTrackDepacketizer
}

// GetID encoding id
func (e *Encoding) GetID() string {
	return e.id
}

// GetSource  get native RTPIncomingSourceGroup
func (e *Encoding) GetSource() native.RTPIncomingSourceGroup {
	return e.source
}

// GetDepacketizer  get native StreamTrackDepacketizer
func (e *Encoding) GetDepacketizer() native.StreamTrackDepacketizer {
	return e.depacketizer
}

// IncomingTrackStopListener stop listener
type IncomingTrackStopListener func()

// IncomingStreamTrack Audio or Video track of a remote media stream
type IncomingStreamTrack struct {
	id                    string
	media                 string
	receiver              native.RTPReceiverFacade
	counter               int
	encodings             []*Encoding
	trackInfo             *sdp.TrackInfo
	stats                 map[string]*IncomingAllStats
	mediaframeMultiplexer *MediaFrameMultiplexer
	onStopListeners       []func()
	onAttachedListeners   []func()
	onDetachedListeners   []func()
}

// IncomingStats info
type IncomingStats struct {
	LostPackets    uint
	DropPackets    uint
	NumPackets     uint
	NumRTCPPackets uint
	TotalBytes     uint
	TotalRTCPBytes uint
	TotalPLIs      uint
	TotalNACKs     uint
	Bitrate        uint
	Layers         []*Layer
}

// IncomingAllStats info
type IncomingAllStats struct {
	Rtt          uint
	MinWaitTime  uint
	MaxWaitTime  uint
	AvgWaitTime  float64
	Media        *IncomingStats
	Rtx          *IncomingStats
	Fec          *IncomingStats
	Bitrate      uint
	Total        uint
	Remb         uint
	SimulcastIdx int
	timestamp    int64
}

// ActiveEncoding info
type ActiveEncoding struct {
	EncodingId   string
	SimulcastIdx int
	Bitrate      uint
	Layers       []*Layer
}

// ActiveLayersInfo info
type ActiveLayersInfo struct {
	Active   []*ActiveEncoding
	Inactive []*ActiveEncoding
	Layers   []*Layer
}

func getStatsFromIncomingSource(source native.RTPIncomingSource) *IncomingStats {

	stats := &IncomingStats{
		LostPackets:    source.GetLostPackets(),
		DropPackets:    source.GetDropPackets(),
		NumPackets:     source.GetNumPackets(),
		NumRTCPPackets: source.GetNumRTCPPackets(),
		TotalBytes:     source.GetTotalBytes(),
		TotalRTCPBytes: source.GetTotalRTCPBytes(),
		TotalPLIs:      source.GetTotalPLIs(),
		TotalNACKs:     source.GetTotalNACKs(),
		Bitrate:        source.GetBitrate(),
		Layers:         []*Layer{},
	}

	layers := source.Layers()

	individual := []*Layer{}

	var i int64
	for i = 0; i < layers.Size(); i++ {
		layer := layers.Get(int64(i))

		layerInfo := &Layer{
			SpatialLayerId:  int(layer.GetSpatialLayerId()),
			TemporalLayerId: int(layer.GetTemporalLayerId()),
			TotalBytes:      layer.GetTotalBytes(),
			NumPackets:      layer.GetNumPackets(),
			Bitrate:         layer.GetBitrate(),
		}

		individual = append(individual, layerInfo)
	}

	for _, layer := range individual {

		aggregated := &Layer{
			SpatialLayerId:  layer.SpatialLayerId,
			TemporalLayerId: layer.TemporalLayerId,
			TotalBytes:      0,
			NumPackets:      0,
			Bitrate:         0,
		}

		for _, layer2 := range individual {

			if layer2.SpatialLayerId <= aggregated.SpatialLayerId && layer2.TemporalLayerId <= aggregated.TemporalLayerId {

				aggregated.TotalBytes += layer2.TotalBytes
				aggregated.NumPackets += layer2.NumPackets
				aggregated.Bitrate += layer2.Bitrate
			}
		}

		stats.Layers = append(stats.Layers, aggregated)
	}

	return stats
}

// NewIncomingStreamTrack Create incoming audio/video track
func NewIncomingStreamTrack(media string, id string, receiver native.RTPReceiverFacade, sources map[string]native.RTPIncomingSourceGroup) *IncomingStreamTrack {
	track := &IncomingStreamTrack{}

	track.id = id
	track.media = media
	track.receiver = receiver
	track.counter = 0
	track.encodings = make([]*Encoding, 0)

	track.trackInfo = sdp.NewTrackInfo(id, media)

	for k, source := range sources {
		encoding := &Encoding{
			id:           k,
			source:       source,
			depacketizer: native.NewStreamTrackDepacketizer(source),
		}

		track.encodings = append(track.encodings, encoding)

		//Add ssrcs to track info
		if source.GetMedia().GetSsrc() > 0 {
			track.trackInfo.AddSSRC(source.GetMedia().GetSsrc())
		}

		if source.GetRtx().GetSsrc() > 0 {
			track.trackInfo.AddSSRC(source.GetRtx().GetSsrc())
		}

		if source.GetFec().GetSsrc() > 0 {
			track.trackInfo.AddSSRC(source.GetFec().GetSsrc())
		}

		//Add RTX and FEC groups
		if source.GetRtx().GetSsrc() > 0 {
			sourceGroup := sdp.NewSourceGroupInfo("FID", []uint{source.GetMedia().GetSsrc(), source.GetRtx().GetSsrc()})
			track.trackInfo.AddSourceGroup(sourceGroup)
		}

		if source.GetFec().GetSsrc() > 0 {
			sourceGroup := sdp.NewSourceGroupInfo("FEC-FR", []uint{source.GetMedia().GetSsrc(), source.GetFec().GetSsrc()})
			track.trackInfo.AddSourceGroup(sourceGroup)
		}

		// if simulcast
		if len(k) > 0 {
			// make soure the pasused
			encodingInfo := sdp.NewTrackEncodingInfo(k, false)
			if source.GetMedia().GetSsrc() > 0 {
				ssrc := strconv.FormatUint(uint64(source.GetMedia().GetSsrc()), 10)
				encodingInfo.AddParam("ssrc", ssrc)
			}
			track.trackInfo.AddEncoding(encodingInfo)
		}
	}

	track.onAttachedListeners = make([]func(), 0)
	track.onDetachedListeners = make([]func(), 0)
	track.onStopListeners = make([]func(), 0)

	sort.SliceStable(track.encodings, func(i, j int) bool {
		return track.encodings[i].id < track.encodings[j].id
	})

	return track
}

// GetID get track id
func (i *IncomingStreamTrack) GetID() string {
	return i.id
}

// GetMedia get track media type "video" or "audio"
func (i *IncomingStreamTrack) GetMedia() string {
	return i.media
}

// GetTrackInfo get track info
func (i *IncomingStreamTrack) GetTrackInfo() *sdp.TrackInfo {
	return i.trackInfo
}

// GetSSRCs get all RTPIncomingSource include "media" "rtx" "fec"
func (i *IncomingStreamTrack) GetSSRCs() []map[string]native.RTPIncomingSource {

	ssrcs := make([]map[string]native.RTPIncomingSource, 0)

	for _, encoding := range i.encodings {
		ssrcs = append(ssrcs, map[string]native.RTPIncomingSource{
			"media": encoding.source.GetMedia(),
			"rtx":   encoding.source.GetRtx(),
			"fec":   encoding.source.GetFec(),
		})
	}
	return ssrcs
}

// GetStats Get stats for all encodings
func (i *IncomingStreamTrack) GetStats() map[string]*IncomingAllStats {

	if i.stats == nil {
		i.stats = map[string]*IncomingAllStats{}
	}

	for _, encoding := range i.encodings {
		state := i.stats[encoding.id]
		if state == nil || (state != nil && time.Now().UnixNano()-state.timestamp > 200000000) {

			encoding.GetSource().Update()

			media := getStatsFromIncomingSource(encoding.GetSource().GetMedia())
			fec := getStatsFromIncomingSource(encoding.GetSource().GetFec())
			rtx := getStatsFromIncomingSource(encoding.GetSource().GetRtx())

			i.stats[encoding.id] = &IncomingAllStats{
				Rtt:         encoding.GetSource().GetRtt(),
				MinWaitTime: encoding.GetSource().GetMinWaitedTime(),
				MaxWaitTime: encoding.GetSource().GetMaxWaitedTime(),
				AvgWaitTime: encoding.GetSource().GetAvgWaitedTime(),
				Media:       media,
				Rtx:         rtx,
				Fec:         fec,
				Bitrate:     media.Bitrate,
				Total:       media.Bitrate + fec.Bitrate + rtx.Bitrate,
				timestamp:   time.Now().UnixNano(),
			}
		}
	}

	simulcastIdx := 0

	stats := []*IncomingAllStats{}

	for _, state := range i.stats {
		stats = append(stats, state)
	}

	sort.Slice(stats, func(i, j int) bool { return stats[i].Bitrate > stats[j].Bitrate })

	for _, state := range stats {
		if state.Bitrate > 0 {
			simulcastIdx += 1
			state.SimulcastIdx = simulcastIdx
		} else {
			state.SimulcastIdx = -1
		}

		for _, layer := range state.Media.Layers {
			layer.SimulcastIdx = state.SimulcastIdx
		}
	}

	return i.stats
}

// GetActiveLayers Get active encodings and layers ordered by bitrate
func (i *IncomingStreamTrack) GetActiveLayers() *ActiveLayersInfo {

	active := []*ActiveEncoding{}
	inactive := []*ActiveEncoding{}
	all := []*Layer{}

	stats := i.GetStats()

	for id, state := range stats {

		if state.Bitrate == 0 {
			inactive = append(inactive, &ActiveEncoding{
				EncodingId: id,
			})
			continue
		}

		encoding := &ActiveEncoding{
			EncodingId:   id,
			SimulcastIdx: state.SimulcastIdx,
			Bitrate:      state.Bitrate,
			Layers:       []*Layer{},
		}

		layers := state.Media.Layers

		for _, layer := range layers {
			encoding.Layers = append(encoding.Layers, &Layer{
				SimulcastIdx:    layer.SimulcastIdx,
				SpatialLayerId:  layer.SpatialLayerId,
				TemporalLayerId: layer.TemporalLayerId,
				Bitrate:         layer.Bitrate,
			})

			all = append(all, &Layer{
				EncodingId:      id,
				SimulcastIdx:    layer.SimulcastIdx,
				SpatialLayerId:  layer.SpatialLayerId,
				TemporalLayerId: layer.TemporalLayerId,
				Bitrate:         layer.Bitrate,
			})

		}

		if len(encoding.Layers) > 0 {
			sort.SliceStable(encoding.Layers, func(i, j int) bool { return encoding.Layers[i].Bitrate < encoding.Layers[j].Bitrate })
		} else {

			all = append(all, &Layer{
				EncodingId:      encoding.EncodingId,
				SimulcastIdx:    encoding.SimulcastIdx,
				SpatialLayerId:  MaxLayerId,
				TemporalLayerId: MaxLayerId,
				Bitrate:         encoding.Bitrate,
			})
		}
		active = append(active, encoding)
	}

	if len(active) > 0 {
		sort.Slice(active, func(i, j int) bool { return active[i].Bitrate < active[j].Bitrate })
	}

	if len(all) > 0 {
		sort.Slice(all, func(i, j int) bool { return all[i].Bitrate < all[j].Bitrate })
	}

	return &ActiveLayersInfo{
		Active:   active,
		Inactive: inactive,
		Layers:   all,
	}

}

// GetEncodings  get all encodings
func (i *IncomingStreamTrack) GetEncodings() []*Encoding {

	return i.encodings
}

// GetFirstEncoding get the first Encoding
func (i *IncomingStreamTrack) GetFirstEncoding() *Encoding {

	for _, encoding := range i.encodings {
		if encoding != nil {
			return encoding
		}
	}
	return nil
}

// GetEncoding get Encoding by id
func (i *IncomingStreamTrack) GetEncoding(encodingID string) *Encoding {

	for _, encoding := range i.encodings {
		if encoding.id == encodingID {
			return encoding
		}
	}
	return nil
}

// Attached Signal that this track has been attached.
func (i *IncomingStreamTrack) Attached() {

	i.counter = i.counter + 1

	if i.counter == 1 {
		for _, attach := range i.onAttachedListeners {
			attach()
		}
	}
}

// Refresh Request an intra refres
func (i *IncomingStreamTrack) Refresh() {

	for _, encoding := range i.encodings {
		//Request an iframe on main ssrc
		i.receiver.SendPLI(encoding.source.GetMedia().GetSsrc())
	}
}

// Detached Signal that this track has been detached.
func (i *IncomingStreamTrack) Detached() {

	i.counter = i.counter - 1

	if i.counter == 0 {
		for _, detach := range i.onDetachedListeners {
			detach()
		}
	}
}

// OnDetach
func (i *IncomingStreamTrack) OnDetach(detach func()) {
	i.onDetachedListeners = append(i.onDetachedListeners, detach)
}

// OnAttach  run this func when attached
func (i *IncomingStreamTrack) OnAttach(attach func()) {
	i.onAttachedListeners = append(i.onAttachedListeners, attach)
}

// OnStop register stop callback
func (i *IncomingStreamTrack) OnStop(stop func()) {
	i.onStopListeners = append(i.onStopListeners, stop)
}

// OnMediaFrame callback
func (i *IncomingStreamTrack) OnMediaFrame(listener func([]byte, uint)) {

	if i.mediaframeMultiplexer == nil {
		i.mediaframeMultiplexer = NewMediaFrameMultiplexer(i)
	}

	i.mediaframeMultiplexer.SetMediaFrameListener(listener)
}

// Stop Removes the track from the incoming stream and also detaches any attached outgoing track or recorder
func (i *IncomingStreamTrack) Stop() {

	if i.receiver == nil {
		return
	}
	
	if i.mediaframeMultiplexer != nil {
		i.mediaframeMultiplexer.Stop()
		i.mediaframeMultiplexer = nil
	}

	for _, stopFunc := range i.onStopListeners {
		stopFunc()
	}

	for _, encoding := range i.encodings {
		if encoding.depacketizer != nil {
			encoding.depacketizer.Stop()
			native.DeleteStreamTrackDepacketizer(encoding.depacketizer)
		}
		if encoding.source != nil {
			native.DeleteRTPIncomingSourceGroup(encoding.source)
		}
	}

	i.encodings = nil

	i.receiver = nil
}
