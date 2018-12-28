package mediaserver

import (
	"fmt"
	"time"

	"github.com/chuckpreslar/emission"
	"github.com/notedit/media-server-go/sdp"
	native "github.com/notedit/media-server-go/wrapper"
)

type OutgoingStreamTrack struct {
	id            string
	media         string
	muted         bool
	sender        native.RTPSenderFacade
	source        native.RTPOutgoingSourceGroup
	transpoder    *Transponder
	trackInfo     *sdp.TrackInfo
	interCallback rembBitrateListener
	statss        *OutgoingStatss
	// todo outercallback
	*emission.Emitter
}

type OutgoingStats struct {
	NumPackets     uint
	NumRTCPPackets uint
	TotalBytes     uint
	TotalRTCPBytes uint
	Bitrate        uint
}

type OutgoingStatss struct {
	Media     *OutgoingStats
	Rtx       *OutgoingStats
	Fec       *OutgoingStats
	timestamp int64
}

type rembBitrateListener interface {
	native.REMBBitrateListener
	deleteREMBBitrateListener()
}

type goREMBBitrateListener struct {
	native.REMBBitrateListener
}

func (r *goREMBBitrateListener) deleteREMBBitrateListener() {
	native.DeleteDirectorREMBBitrateListener(r.REMBBitrateListener)
}

type overwrittenREMBBitrateListener struct {
	p     native.REMBBitrateListener
	track *OutgoingStreamTrack
}

func (p *overwrittenREMBBitrateListener) OnREMB() {

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

func newOutgoingStreamTrack(media string, id string, sender native.RTPSenderFacade, source native.RTPOutgoingSourceGroup) *OutgoingStreamTrack {

	track := &OutgoingStreamTrack{}
	track.id = id
	track.media = media
	track.sender = sender
	track.muted = false
	track.source = source
	track.Emitter = emission.NewEmitter()
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

	// callback
	callback := &overwrittenREMBBitrateListener{
		track: track,
	}
	p := native.NewDirectorREMBBitrateListener(callback)
	callback.p = p

	track.interCallback = &goREMBBitrateListener{REMBBitrateListener: p}

	return track
}

func (o *OutgoingStreamTrack) GetID() string {
	return o.id
}

func (o *OutgoingStreamTrack) GetMedia() string {
	return o.media
}

func (o *OutgoingStreamTrack) GetTrackInfo() *sdp.TrackInfo {
	return o.trackInfo
}

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

func (o *OutgoingStreamTrack) GetSSRCs() map[string]native.RTPOutgoingSource {

	return map[string]native.RTPOutgoingSource{
		"media": o.source.GetMedia(),
		"rtx":   o.source.GetRtx(),
		"fec":   o.source.GetFec(),
	}
}

func (o *OutgoingStreamTrack) IsMuted() bool {
	return o.muted
}

func (o *OutgoingStreamTrack) Mute(muting bool) {

	if o.transpoder != nil {
		o.transpoder.Mute(muting)
	}

	if o.muted != muting {
		o.muted = muting
		o.EmitSync("muted", o.muted)
	}
}

func (o *OutgoingStreamTrack) AttachTo(incomingTrack *IncomingStreamTrack) *Transponder {

	// detach first
	o.Detach()

	// todo add remblistener
	transponder := native.NewRTPStreamTransponderFacade(o.source, o.sender, o.interCallback)

	o.transpoder = NewTransponder(transponder)

	if o.muted {
		o.transpoder.Mute(o.muted)
	}

	fmt.Println(" incomingTrack", incomingTrack)

	o.transpoder.SetIncomingTrack(incomingTrack)

	o.transpoder.Once("stopped", o.onTransponderStopped)

	return o.transpoder
}

func (o *OutgoingStreamTrack) Detach() {

	if o.transpoder == nil {
		return
	}

	o.transpoder.Off("stopped", o.onTransponderStopped)

	o.transpoder.Stop()

	o.transpoder = nil
}

func (o *OutgoingStreamTrack) GetTransponder() *Transponder {

	return o.transpoder
}

func (o *OutgoingStreamTrack) Stop() {

	if o.sender == nil {
		return
	}

	// swig memory clean
	o.interCallback.deleteREMBBitrateListener()

	o.Detach()

	o.EmitSync("stopped")

	if o.source != nil {
		native.DeleteRTPOutgoingSourceGroup(o.source)
	}

	o.transpoder.Stop()

	o.source = nil

	o.sender = nil
}

func (o *OutgoingStreamTrack) onTransponderStopped() {

	o.transpoder = nil
}
