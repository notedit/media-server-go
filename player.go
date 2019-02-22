package mediaserver

import (
	"errors"

	native "github.com/notedit/media-server-go/wrapper"
)

// Player can player the file as webrtc source
type Player struct {
	player      native.PlayerFacade
	tracks      map[string]*IncomingStreamTrack
	endCallback playerEndCallback
	endListener PlayerEndListener
}

// PlayerEndListener end listener
type PlayerEndListener func()

type playerEndCallback interface {
	native.PlayerEndListener
	deletePlayerListener()
}

type goplayerEndCallback struct {
	native.PlayerEndListener
}

func (r *goplayerEndCallback) deletePlayerListener() {
	native.DeleteDirectorPlayerEndListener(r.PlayerEndListener)
}

type overwrittenEndCallback struct {
	p      native.PlayerEndListener
	player *Player
}

func (p *overwrittenEndCallback) OnEnd() {
	if p.player != nil && p.player.endListener != nil {
		p.player.endListener()
	}
}

// NewPlayer create new file player
func NewPlayer(filename string, listener PlayerEndListener) (*Player, error) {
	player := &Player{}
	player.player = native.NewPlayerFacade()
	player.tracks = make(map[string]*IncomingStreamTrack)

	if player.player.Open(filename) == 0 {
		native.DeletePlayerFacade(player.player)
		return nil, errors.New("player can not open filanme:" + filename)
	}

	if player.player.HasVideoTrack() {

		trackID := "video"
		source := player.player.GetVideoSource()

		incoming := NewIncomingStreamTrack("video", trackID, nil, map[string]native.RTPIncomingSourceGroup{"": source})

		player.tracks[trackID] = incoming
	}

	if player.player.HasAudioTrack() {

		trackID := "audio"
		source := player.player.GetAudioSource()

		incoming := NewIncomingStreamTrack("audio", trackID, nil, map[string]native.RTPIncomingSourceGroup{"": source})

		player.tracks[trackID] = incoming

	}

	callback := &overwrittenEndCallback{
		player: player,
	}
	p := native.NewDirectorPlayerEndListener(callback)
	callback.p = p

	player.endCallback = &goplayerEndCallback{PlayerEndListener: p}

	player.endListener = listener

	player.player.SetPlayEndListener(player.endCallback)

	return player, nil
}

// GetTracks tracks this file contains
func (p *Player) GetTracks() []*IncomingStreamTrack {
	tracks := []*IncomingStreamTrack{}
	for _, track := range p.tracks {
		tracks = append(tracks, track)
	}
	return tracks
}

// GetAudioTracks audio tracks this file contains
func (p *Player) GetAudioTracks() []*IncomingStreamTrack {
	tracks := []*IncomingStreamTrack{}
	for _, track := range p.tracks {
		if track.GetMedia() == "audio" {
			tracks = append(tracks, track)
		}
	}
	return tracks
}

// GetVideoTracks video tracks this file contains
func (p *Player) GetVideoTracks() []*IncomingStreamTrack {
	tracks := []*IncomingStreamTrack{}
	for _, track := range p.tracks {
		if track.GetMedia() == "video" {
			tracks = append(tracks, track)
		}
	}
	return tracks
}

// Play start
func (p *Player) Play() {
	if p.player != nil {
		p.player.Play()
	}
}

// Resume play
func (p *Player) Resume() {

	if p.player != nil {
		p.player.Play()
	}
}

// Pause  play
func (p *Player) Pause() {

	if p.player != nil {
		p.player.Stop()
	}
}

// Seek seek to
func (p *Player) Seek(time uint64) {

	if p.player != nil {
		p.player.Seek(time)
	}
}

// Stop  stop play
func (p *Player) Stop() {

	if p.player == nil {
		return
	}

	if p.endCallback != nil {
		p.endCallback.deletePlayerListener()
	}

	for _, track := range p.tracks {
		track.Stop()
	}

	p.tracks = nil

	p.player.Close()

	p.player = nil
}
