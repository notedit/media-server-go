package mediaserver

import (
	"strings"

	"github.com/chuckpreslar/emission"
	"github.com/notedit/media-server-go/sdp"
	native "github.com/notedit/media-server-go/wrapper"
)

// OutgoingStreamStopListener stop listener
type OutgoingStreamStopListener func()

// OutgoingStream  represent the media stream sent to a remote peer
type OutgoingStream struct {
	id              string
	transport       native.DTLSICETransport
	info            *sdp.StreamInfo
	muted           bool
	tracks          map[string]*OutgoingStreamTrack
	onStopListeners []OutgoingStreamStopListener
	*emission.Emitter
}

// NewOutgoingStream create outgoing stream
func NewOutgoingStream(transport native.DTLSICETransport, info *sdp.StreamInfo) *OutgoingStream {
	stream := new(OutgoingStream)

	stream.id = info.GetID()
	stream.transport = transport
	stream.info = info
	stream.tracks = make(map[string]*OutgoingStreamTrack)
	stream.Emitter = emission.NewEmitter()

	for _, track := range info.GetTracks() {
		stream.CreateTrack(track)
	}

	stream.onStopListeners = make([]OutgoingStreamStopListener, 0)

	return stream
}

// GetID get id
func (o *OutgoingStream) GetID() string {
	return o.id
}

// GetStats Get statistics for all tracks in the stream
func (o *OutgoingStream) GetStats() map[string]*OutgoingStatss {

	stats := map[string]*OutgoingStatss{}
	for _, track := range o.tracks {
		stats[track.GetID()] = track.GetStats()
	}
	return stats
}

// IsMuted Check if the stream is muted or not
func (o *OutgoingStream) IsMuted() bool {
	return o.muted
}

// Mute Mute/Unmute this stream and all the tracks in it
func (o *OutgoingStream) Mute(muting bool) {

	for _, track := range o.tracks {
		track.Mute(muting)
	}

	if o.muted != muting {
		o.muted = muting
		o.EmitSync("muted", o.muted)
	}
}

// AttachTo Listen media from the incoming stream and send it to the remote peer of the associated transport
func (o *OutgoingStream) AttachTo(incomingStream *IncomingStream) []*Transponder {

	o.Detach()
	transponders := []*Transponder{}
	audios := o.GetAudioTracks()
	if len(audios) > 0 {
		index := len(audios)
		tracks := incomingStream.GetAudioTracks()
		if index < len(tracks) {
			index = len(tracks)
		}

		for i, track := range tracks {
			if i < index {
				transponders = append(transponders, audios[i].AttachTo(track))
			}
		}
	}

	videos := o.GetVideoTracks()
	if len(videos) > 0 {
		index := len(videos)
		tracks := incomingStream.GetVideoTracks()
		if index < len(tracks) {
			index = len(tracks)
		}
		for i, track := range tracks {
			if i < index {
				transponders = append(transponders, videos[i].AttachTo(track))
			}
		}
	}

	return transponders
}

// Detach Stop listening for media
func (o *OutgoingStream) Detach() {

	for _, track := range o.tracks {
		track.Detach()
	}
}

// GetStreamInfo get the stream info
func (o *OutgoingStream) GetStreamInfo() *sdp.StreamInfo {

	return o.info
}

// GetTrack get one track
func (o *OutgoingStream) GetTrack(trackID string) *OutgoingStreamTrack {
	return o.tracks[trackID]
}

// GetTracks get all the tracks
func (o *OutgoingStream) GetTracks() []*OutgoingStreamTrack {
	tracks := []*OutgoingStreamTrack{}
	for _, track := range o.tracks {
		tracks = append(tracks, track)
	}
	return tracks
}

// GetAudioTracks Get an array of the media stream audio tracks
func (o *OutgoingStream) GetAudioTracks() []*OutgoingStreamTrack {
	audioTracks := []*OutgoingStreamTrack{}
	for _, track := range o.tracks {
		if strings.ToLower(track.GetMedia()) == "audio" {
			audioTracks = append(audioTracks, track)
		}
	}
	return audioTracks
}

// GetVideoTracks Get an array of the media stream video tracks
func (o *OutgoingStream) GetVideoTracks() []*OutgoingStreamTrack {
	videoTracks := []*OutgoingStreamTrack{}
	for _, track := range o.tracks {
		if strings.ToLower(track.GetMedia()) == "video" {
			videoTracks = append(videoTracks, track)
		}
	}
	return videoTracks
}

// AddTrack add one outgoing track
func (o *OutgoingStream) AddTrack(track *OutgoingStreamTrack) {

	if _, ok := o.tracks[track.GetID()]; ok {
		return
	}

	track.Once("stopped", func() {
		delete(o.tracks, track.GetID())
	})

	o.tracks[track.GetID()] = track
}

// CreateTrack Create new track from a TrackInfo object and add it to this stream
func (o *OutgoingStream) CreateTrack(track *sdp.TrackInfo) *OutgoingStreamTrack {

	var mediaType native.MediaFrameType = 0
	if track.GetMedia() == "video" {
		mediaType = 1
	}

	source := native.NewRTPOutgoingSourceGroup(mediaType)

	source.GetMedia().SetSsrc(track.GetSSRCS()[0])

	fid := track.GetSourceGroup("FID")
	fec_fr := track.GetSourceGroup("FEC-FR")

	if fid != nil {
		source.GetRtx().SetSsrc(fid.GetSSRCs()[1])
	} else {
		source.GetRtx().SetSsrc(0)
	}

	if fec_fr != nil {
		source.GetFec().SetSsrc(fec_fr.GetSSRCs()[1])
	} else {
		source.GetFec().SetSsrc(0)
	}

	o.transport.AddOutgoingSourceGroup(source)

	outgoingTrack := newOutgoingStreamTrack(track.GetMedia(), track.GetID(), native.TransportToSender(o.transport), source)

	outgoingTrack.Once("stopped", func() {
		delete(o.tracks, outgoingTrack.GetID())
		o.transport.RemoveOutgoingSourceGroup(source)
	})

	o.tracks[outgoingTrack.GetID()] = outgoingTrack

	o.EmitSync("track", outgoingTrack)

	return outgoingTrack
}

func (o *OutgoingStream) OnStop(stop OutgoingStreamStopListener) {
	o.onStopListeners = append(o.onStopListeners, stop)
}

// Stop stop the remote stream
func (o *OutgoingStream) Stop() {

	if o.transport == nil {
		return
	}

	for _, track := range o.tracks {
		track.Stop()
	}

	for _, stopFunc := range o.onStopListeners {
		stopFunc()
	}

	o.tracks = make(map[string]*OutgoingStreamTrack, 0)

	o.EmitSync("stopped")

	o.transport = nil
}
