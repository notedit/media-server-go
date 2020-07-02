package mediaserver

import (
	"time"

	native "github.com/notedit/media-server-go/wrapper"
	"github.com/notedit/sdp"
)

// OutgoingStreamTrack Audio or Video track of a media stream sent to a remote peer
type OutgoingStreamTrack struct {
	id              string
	media           string
	muted           bool
	sender          native.RTPSenderFacade
	source          native.RTPOutgoingSourceGroup
	transpoder      *Transponder
	trackInfo       *sdp.TrackInfo
	statss          *OutgoingStatss
	onMuteListeners []func(bool)
	onStopListeners []func()
	// todo outercallback
}

// OutgoingStats stats info
type OutgoingStats struct {
	NumPackets     uint
	NumRTCPPackets uint
	TotalBytes     uint
	TotalRTCPBytes uint
	Bitrate        uint
}

// OutgoingStatss stats info
type OutgoingStatss struct {
	Media     *OutgoingStats
	Rtx       *OutgoingStats
	Fec       *OutgoingStats
	timestamp int64
}

func getStatsFromOutgoingSource(source native.RTPOutgoingSource) *OutgoingStats {

	stats := &OutgoingStats{
		NumPackets:     source.GetNumPackets(),
		NumRTCPPackets: source.GetNumRTCPPackets(),
		TotalBytes:     source.GetTotalBytes(),
		TotalRTCPBytes: source.GetTotalRTCPBytes(),
		Bitrate:        source.GetBitrate(),
	}

	return stats
}

// NewOutgoingStreamTrack create outgoing stream track
func newOutgoingStreamTrack(media string, id string, sender native.RTPSenderFacade, source native.RTPOutgoingSourceGroup) *OutgoingStreamTrack {

	track := &OutgoingStreamTrack{}
	track.id = id
	track.media = media
	track.sender = sender
	track.muted = false
	track.source = source
	track.trackInfo = sdp.NewTrackInfo(id, media)

	track.trackInfo.AddSSRC(source.GetMedia().GetSsrc())

	if source.GetRtx().GetSsrc() > 0 {
		track.trackInfo.AddSSRC(source.GetRtx().GetSsrc())
	}

	if source.GetFec().GetSsrc() > 0 {
		track.trackInfo.AddSSRC(source.GetFec().GetSsrc())
	}

	if source.GetRtx().GetSsrc() > 0 {
		sourceGroup := sdp.NewSourceGroupInfo("FID", []uint{source.GetMedia().GetSsrc(), source.GetRtx().GetSsrc()})
		track.trackInfo.AddSourceGroup(sourceGroup)
	}

	if source.GetFec().GetSsrc() > 0 {
		sourceGroup := sdp.NewSourceGroupInfo("FEC-FR", []uint{source.GetMedia().GetSsrc(), source.GetFec().GetSsrc()})
		track.trackInfo.AddSourceGroup(sourceGroup)
	}


	track.onMuteListeners = make([]func(bool), 0)
	track.onStopListeners = make([]func(), 0)

	return track
}

// GetID  get track id
func (o *OutgoingStreamTrack) GetID() string {
	return o.id
}

// GetMedia get media type
func (o *OutgoingStreamTrack) GetMedia() string {
	return o.media
}

// GetTrackInfo get track info
func (o *OutgoingStreamTrack) GetTrackInfo() *sdp.TrackInfo {
	return o.trackInfo
}

// GetStats get stats info
func (o *OutgoingStreamTrack) GetStats() *OutgoingStatss {

	if o.statss == nil {
		o.statss = &OutgoingStatss{}
	}

	if time.Now().UnixNano()-o.statss.timestamp > 200000000 {
		o.statss.Media = getStatsFromOutgoingSource(o.source.GetMedia())
		o.statss.Rtx = getStatsFromOutgoingSource(o.source.GetRtx())
		o.statss.Fec = getStatsFromOutgoingSource(o.source.GetFec())
		o.statss.timestamp = time.Now().UnixNano()
	}

	return o.statss
}

// GetSSRCs get ssrcs map
func (o *OutgoingStreamTrack) GetSSRCs() map[string]native.RTPOutgoingSource {

	return map[string]native.RTPOutgoingSource{
		"media": o.source.GetMedia(),
		"rtx":   o.source.GetRtx(),
		"fec":   o.source.GetFec(),
	}
}

// IsMuted Check if the track is muted or not
func (o *OutgoingStreamTrack) IsMuted() bool {
	return o.muted
}

// Mute Mute/Unmute the track
func (o *OutgoingStreamTrack) Mute(muting bool) {

	if o.transpoder != nil {
		o.transpoder.Mute(muting)
	}

	if o.muted != muting {
		o.muted = muting

		for _, mutefunc := range o.onMuteListeners {
			mutefunc(muting)
		}
	}
}

// AttachTo Listen media from the incoming stream track and send it to the remote peer of the associated transport
func (o *OutgoingStreamTrack) AttachTo(incomingTrack *IncomingStreamTrack) *Transponder {

	// detach first
	o.Detach()

	transponder := native.NewRTPStreamTransponderFacade(o.source, o.sender)

	o.transpoder = NewTransponder(transponder)

	if o.muted {
		o.transpoder.Mute(o.muted)
	}

	o.transpoder.SetIncomingTrack(incomingTrack)

	return o.transpoder
}

// Detach Stop forwarding any previous attached track
func (o *OutgoingStreamTrack) Detach() {

	if o.transpoder == nil {
		return
	}

	o.transpoder.Stop()

	o.transpoder = nil
}

// GetTransponder Get attached transpoder for this track
func (o *OutgoingStreamTrack) GetTransponder() *Transponder {
	return o.transpoder
}

func (o *OutgoingStreamTrack) OnMute(mute func(bool)) {
	o.onMuteListeners = append(o.onMuteListeners, mute)
}

// Stop Removes the track from the outgoing stream and also detaches from any attached incoming track
func (o *OutgoingStreamTrack) Stop() {

	if o.sender == nil {
		return
	}

	if o.transpoder != nil { // maybe = nil at onTransponderStopped
		o.transpoder.Stop()
		o.transpoder = nil
	}

	native.DeleteRTPSenderFacade(o.sender)
	o.sender = nil
}

func (o *OutgoingStreamTrack) GetRTPOutgoingSourceGroup() native.RTPOutgoingSourceGroup {
	return o.source
}

func (o *OutgoingStreamTrack) DeleteRTPOutgoingSourceGroup(bundle native.RTPBundleTransport) {

	if o.source != nil {
		bundle.DeleteOutGoingSourceGroup(o.source)
	}

	o.source = nil
}
