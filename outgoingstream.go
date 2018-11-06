package mediaserver

import (
	"strings"

	"./sdp"
	"github.com/chuckpreslar/emission"
)

type OutgoingStream struct {
	id        string
	transport *Transport
	info      *sdp.StreamInfo
	muted     bool
	tracks    map[string]*OutgoingStreamTrack
	*emission.Emitter
}

func NewOutgoingStream(transport *Transport, info *sdp.StreamInfo) *OutgoingStream {
	stream := new(OutgoingStream)

	stream.id = info.GetID()
	stream.transport = transport
	stream.info = info
	stream.tracks = make(map[string]*OutgoingStreamTrack)
	stream.Emitter = emission.NewEmitter()

	return stream
}

func (o *OutgoingStream) GetID() string {
	return o.id
}

func (o *OutgoingStream) GetStats() {

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
