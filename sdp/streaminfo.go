package sdp

import (
	"strings"
)

type StreamInfo struct {
	id     string
	tracks map[string]*TrackInfo
}

func NewStreamInfo(streamID string) *StreamInfo {

	return &StreamInfo{
		id:     streamID,
		tracks: map[string]*TrackInfo{},
	}
}

func (s *StreamInfo) Clone() *StreamInfo {
	stream := &StreamInfo{
		id:     s.id,
		tracks: make(map[string]*TrackInfo),
	}

	for k, v := range s.tracks {
		stream.tracks[k] = v.Clone()
	}
	return stream
}

func (s *StreamInfo) GetID() string {

	return s.id
}

func (s *StreamInfo) AddTrack(track *TrackInfo) {

	s.tracks[track.GetID()] = track
}

func (s *StreamInfo) RemoveTrack(track *TrackInfo) {

	delete(s.tracks, track.GetID())
}

func (s *StreamInfo) RemoveTrackById(trackId string) {
	delete(s.tracks, trackId)
}

func (s *StreamInfo) GetFirstTrack(media string) *TrackInfo {

	for _, trak := range s.tracks {

		if strings.ToLower(trak.GetMedia()) == strings.ToLower(media) {

			return trak
		}
	}
	return nil
}

func (s *StreamInfo) GetTracks() map[string]*TrackInfo {

	return s.tracks
}

func (s *StreamInfo) RemoveAllTracks() {

	for trackId, _ := range s.tracks {
		delete(s.tracks, trackId)
	}
}

func (s *StreamInfo) GetTrack(trackID string) *TrackInfo {

	return s.tracks[trackID]
}
