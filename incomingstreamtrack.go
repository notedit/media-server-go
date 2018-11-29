package mediaserver

import (
	"sort"
	"time"

	"github.com/chuckpreslar/emission"
	"github.com/notedit/media-server-go/sdp"
)

type Layer struct {
	SpatialLayerId  byte
	TemporalLayerId byte
	TotalBytes      uint
	NumPackets      uint
	Bitrate         uint
	SimulcastIdx    int
}

type Encoding struct {
	id           string
	source       RTPIncomingSourceGroup
	depacketizer StreamTrackDepacketizer
}

func (e *Encoding) GetID() string {
	return e.id
}

func (e *Encoding) GetSource() RTPIncomingSourceGroup {
	return e.source
}

func (e *Encoding) GetDepacketizer() StreamTrackDepacketizer {
	return e.depacketizer
}

type IncomingStreamTrack struct {
	id        string
	media     string
	receiver  RTPReceiverFacade
	counter   int
	encodings map[string]*Encoding
	trackInfo *sdp.TrackInfo
	stats     map[string]*IncomingAllStats
	*emission.Emitter
}

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

func getStatsFromIncomingSource(source RTPIncomingSource) *IncomingStats {

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
			SpatialLayerId:  layer.GetSpatialLayerId(),
			TemporalLayerId: layer.GetTemporalLayerId(),
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

func newIncomingStreamTrack(media string, id string, receiver RTPReceiverFacade, souces map[string]RTPIncomingSourceGroup) *IncomingStreamTrack {
	track := &IncomingStreamTrack{}

	track.id = id
	track.media = media
	track.receiver = receiver
	track.counter = 0
	track.encodings = make(map[string]*Encoding)
	track.Emitter = emission.NewEmitter()

	for k, source := range souces {
		encoding := &Encoding{
			id:           k,
			source:       source,
			depacketizer: NewStreamTrackDepacketizer(source),
		}
		track.encodings[k] = encoding
	}

	return track
}

func (i *IncomingStreamTrack) GetID() string {
	return i.id
}

func (i *IncomingStreamTrack) GetMedia() string {
	return i.media
}

func (i *IncomingStreamTrack) GetTrackInfo() *sdp.TrackInfo {
	return i.trackInfo
}

func (i *IncomingStreamTrack) GetSSRCs() []map[string]RTPIncomingSource {

	ssrcs := make([]map[string]RTPIncomingSource, 0)

	for _, encoding := range i.encodings {
		ssrcs = append(ssrcs, map[string]RTPIncomingSource{
			"media": encoding.source.GetMedia(),
			"rtx":   encoding.source.GetRtx(),
			"fec":   encoding.source.GetFec(),
		})
	}
	return ssrcs
}

func (i *IncomingStreamTrack) GetStats() map[string]*IncomingAllStats {

	if i.stats == nil {
		i.stats = map[string]*IncomingAllStats{}
	}

	for id, encoding := range i.encodings {
		state := i.stats[id]
		if state == nil || (state != nil && time.Now().UnixNano()-state.timestamp > 200000000) {

			encoding.GetSource().Update()

			media := getStatsFromIncomingSource(encoding.GetSource().GetMedia())
			fec := getStatsFromIncomingSource(encoding.GetSource().GetFec())
			rtx := getStatsFromIncomingSource(encoding.GetSource().GetRtx())

			i.stats[id] = &IncomingAllStats{
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

func (i *IncomingStreamTrack) GetActiveLayers() {

}

func (i *IncomingStreamTrack) GetEncodings() map[string]*Encoding {

	return i.encodings
}

func (i *IncomingStreamTrack) GetFirstEncoding() *Encoding {

	for _, encoding := range i.encodings {
		if encoding != nil {
			return encoding
		}
	}
	return nil
}

func (i *IncomingStreamTrack) Attached() {

	i.counter = i.counter + 1

	if i.counter == 1 {
		i.EmitSync("attached")
	}
}

func (i *IncomingStreamTrack) Refresh() {

	for _, encoding := range i.encodings {
		//Request an iframe on main ssrc
		i.receiver.SendPLI(encoding.source.GetMedia().GetSsrc())
	}
}

func (i *IncomingStreamTrack) Detached() {

	i.counter = i.counter - 1

	if i.counter == 0 {
		i.EmitSync("detached")
	}
}

func (i *IncomingStreamTrack) Stop() {

	if i.receiver == nil {
		return
	}

	for _, encoding := range i.encodings {
		if encoding.depacketizer != nil {
			encoding.depacketizer.Stop()
		}
	}

	i.EmitSync("stopped")

	i.encodings = nil

	i.receiver = nil
}
