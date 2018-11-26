package mediaserver

import (
	"errors"
	"strings"

	"github.com/chuckpreslar/emission"
	"github.com/notedit/media-server-go/sdp"
)

type OutgoingStream struct {
	id        string
	transport DTLSICETransport
	info      *sdp.StreamInfo
	muted     bool
	tracks    map[string]*OutgoingStreamTrack
	*emission.Emitter
}

func NewOutgoingStream(transport DTLSICETransport, info *sdp.StreamInfo) *OutgoingStream {
	stream := new(OutgoingStream)

	stream.id = info.GetID()
	stream.transport = transport
	stream.info = info
	stream.tracks = make(map[string]*OutgoingStreamTrack)
	stream.Emitter = emission.NewEmitter()

	for _, track := range info.GetTracks() {
		stream.CreateTrack(track)
	}

	return stream
}

func (o *OutgoingStream) GetID() string {
	return o.id
}

func (o *OutgoingStream) GetStats() {
	// todo
}

func (o *OutgoingStream) IsMuted() bool {
	return o.muted
}

func (o *OutgoingStream) Mute(muting bool) {

	for _, track := range o.tracks {
		track.Mute(muting)
	}

	if o.muted != muting {
		o.muted = muting
		o.EmitSync("muted", o.muted)
	}
}

func (o *OutgoingStream) AttachTo(incomingStream *IncomingStream) []*Transponder {

	// detach first
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

func (o *OutgoingStream) Detach() {

	for _, track := range o.tracks {
		track.Detach()
	}
}

func (o *OutgoingStream) GetStreamInfo() *sdp.StreamInfo {

	return o.info
}

func (o *OutgoingStream) GetTrack(trackID string) *OutgoingStreamTrack {
	return o.tracks[trackID]
}

func (o *OutgoingStream) GetTracks() []*OutgoingStreamTrack {
	tracks := []*OutgoingStreamTrack{}
	for _, track := range o.tracks {
		tracks = append(tracks, track)
	}
	return tracks
}

func (o *OutgoingStream) GetAudioTracks() []*OutgoingStreamTrack {
	audioTracks := []*OutgoingStreamTrack{}
	for _, track := range o.tracks {
		if strings.ToLower(track.GetMedia()) == "audio" {
			audioTracks = append(audioTracks, track)
		}
	}
	return audioTracks
}

func (o *OutgoingStream) GetVideoTracks() []*OutgoingStreamTrack {
	videoTracks := []*OutgoingStreamTrack{}
	for _, track := range o.tracks {
		if strings.ToLower(track.GetMedia()) == "video" {
			videoTracks = append(videoTracks, track)
		}
	}
	return videoTracks
}

func (o *OutgoingStream) AddTrack(track *OutgoingStreamTrack) error {

	if _, ok := o.tracks[track.GetID()]; ok {
		return errors.New("Track id already present in stream")
	}

	track.Once("stopped", func() {
		delete(o.tracks, track.GetID())
	})

	o.tracks[track.GetID()] = track

	return nil
}

func (o *OutgoingStream) CreateTrack(track *sdp.TrackInfo) *OutgoingStreamTrack {

	var mediaType MediaFrameType = 0
	if track.GetMedia() == "video" {
		mediaType = 1
	}

	source := NewRTPOutgoingSourceGroup(mediaType)

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

	outgoingTrack := newOutgoingStreamTrack(track.GetMedia(), track.GetID(), TransportToSender(o.transport), source)

	outgoingTrack.Once("stopped", func() {
		delete(o.tracks, outgoingTrack.GetID())
		o.transport.RemoveOutgoingSourceGroup(source)
	})

	o.tracks[outgoingTrack.GetID()] = outgoingTrack

	o.EmitSync("track", outgoingTrack)

	return outgoingTrack
}

func (o *OutgoingStream) Stop() {

	if o.transport == nil {
		return
	}

	for _, track := range o.tracks {
		track.Stop()
	}

	o.tracks = nil

	o.EmitSync("stopped")

	o.transport = nil
}
