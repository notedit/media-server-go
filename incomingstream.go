package mediaserver

import (
	"strings"

	"github.com/chuckpreslar/emission"
	"github.com/notedit/media-server-go/sdp"
)

type IncomingStream struct {
	id        string
	info      *sdp.StreamInfo
	transport DTLSICETransport
	tracks    map[string]*IncomingStreamTrack
	*emission.Emitter
}

func newIncomingStream(transport DTLSICETransport, receiver RTPReceiverFacade, info *sdp.StreamInfo) *IncomingStream {
	stream := &IncomingStream{}

	stream.id = info.GetID()
	stream.transport = transport
	stream.tracks = make(map[string]*IncomingStreamTrack)
	stream.Emitter = emission.NewEmitter()

	for _, track := range info.GetTracks() {

		var mediaType MediaFrameType = 0
		if track.GetMedia() == "video" {
			mediaType = 1
		}

		sources := map[string]RTPIncomingSourceGroup{}

		encodings := track.GetEncodings()

		if len(encodings) > 0 {
			// simulcast encoding
			// we handle it later
			// todo
		} else if track.GetSourceGroup("SIM") != nil {
			// chrome like simulcast
			// todo
		} else {
			source := NewRTPIncomingSourceGroup(mediaType)

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

			stream.transport.AddIncomingSourceGroup(source)

			// Append to soruces with empty rid
			sources[""] = source
		}

		incomingTrack := newIncomingStreamTrack(track.GetMedia(), track.GetID(), receiver, sources)

		incomingTrack.Once("stopped", func() {

			delete(stream.tracks, incomingTrack.GetID())

			for _, source := range sources {
				stream.transport.RemoveIncomingSourceGroup(source)
			}
		})

		stream.tracks[incomingTrack.GetID()] = incomingTrack
	}

	return stream
}

func (i *IncomingStream) GetID() string {
	return i.id
}

func (i *IncomingStream) GetStreamInfo() *sdp.StreamInfo {
	return i.info
}

func (i *IncomingStream) GetStats() map[string]*IncomingStats {

	stats := map[string]*IncomingStats{}

	for _, track := range i.tracks {
		stats[track.GetID()] = track.GetStats()
	}

	return stats
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

	i.Emit("stopped")

	i.transport = nil
}
