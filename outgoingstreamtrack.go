package mediaserver

import (
	"github.com/chuckpreslar/emission"
	"github.com/notedit/media-server-go/sdp"
)

type OutgoingStreamTrack struct {
	id         string
	media      string
	muted      bool
	sender     RTPSenderFacade
	source     RTPOutgoingSourceGroup
	transpoder *Transponder
	trackInfo  *sdp.TrackInfo
	*emission.Emitter
}

type OutgoingStats struct {
	NumPackets     int
	NumRTCPPackets int
	TotalBytes     int
	TotalRTCPBytes int
	Bitrate        int
}

func newOutgoingStreamTrack(media string, id string, sender RTPSenderFacade, source RTPOutgoingSourceGroup) *OutgoingStreamTrack {

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

	// todo onremb callback

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

func (o *OutgoingStreamTrack) GetStats() {

}

func (o *OutgoingStreamTrack) GetSSRCs() map[string]RTPOutgoingSource {

	return map[string]RTPOutgoingSource{
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
	transponder := NewRTPStreamTransponderFacade(o.source, o.sender, o)

	o.transpoder = NewTransponder(transponder)

	if o.muted {
		o.transpoder.Mute(o.muted)
	}

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

	o.Detach()

	o.EmitSync("stopped")

	o.source = nil

	o.sender = nil
}

func (o *OutgoingStreamTrack) onTransponderStopped() {

	o.transpoder = nil
}

// add fake OnREMB listender
func (o *OutgoingStreamTrack) Swigcptr() uintptr {
	return 0
}

func (o *OutgoingStreamTrack) SwigIsREMBListener() {

}

func (o *OutgoingStreamTrack) OnREMB() {

}
