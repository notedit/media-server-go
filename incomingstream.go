package mediaserver

import (
	"strings"

	"github.com/chuckpreslar/emission"

	"./sdp"
)

type IncomingStream struct {
	id        string
	info      *sdp.StreamInfo
	transport *Transport
	tracks    map[string]*IncomingStreamTrack
	*emission.Emitter
}

func newIncomingStream(transport *Transport, receiver RTPReceiverFacade, info *sdp.StreamInfo) *IncomingStream {
	stream := &IncomingStream{}

	stream.id = info.GetID()
	stream.transport = transport
	stream.tracks = make(map[string]*IncomingStreamTrack)
	stream.Emitter = emission.NewEmitter()
	// todo init track

	return stream
}

func (i *IncomingStream) GetID() string {
	return i.id
}

func (i *IncomingStream) GetStreamInfo() *sdp.StreamInfo {
	return i.info
}

func (i *IncomingStream) GetStats() map[string]interface{} {

	return nil
}

func (i *IncomingStream) GetTrack(trackID string) *IncomingStreamTrack {
	return i.tracks[trackID]
}

func (i *IncomingStream) GetTracks() []*IncomingStreamTrack {
	tracks := []*IncomingStreamTrack{}
	for _, track := range i.tracks {
		tracks = append(tracks, track)
	}
	return tracks
}

func (i *IncomingStream) GetAudioTracks() []*IncomingStreamTrack {
	audioTracks := []*IncomingStreamTrack{}
	for _, track := range i.tracks {
		if strings.ToLower(track.GetMedia()) == "audio" {
			audioTracks = append(audioTracks, track)
		}
	}
	return audioTracks
}

func (i *IncomingStream) GetVideoTracks() []*IncomingStreamTrack {
	videoTracks := []*IncomingStreamTrack{}
	for _, track := range i.tracks {
		if strings.ToLower(track.GetMedia()) == "video" {
			videoTracks = append(videoTracks, track)
		}
	}
	return videoTracks

}

func (i *IncomingStream) Stop() {

	if i.transport == nil {
		return
	}

	for k, track := range i.tracks {
		track.Stop()
		delete(i.tracks, k)
	}

	i.transport = nil
}
