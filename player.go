package mediaserver

import "errors"

type Player struct {
	player PlayerFacade
	tracks map[string]*IncomingStreamTrack
}

func NewPlayer(filename string) (*Player, error) {
	player := &Player{}
	player.player = NewPlayerFacade(player)
	player.tracks = make(map[string]*IncomingStreamTrack)

	if player.player.Open(filename) == 0 {
		return nil, errors.New("player can not open filanme:" + filename)
	}

	if player.player.HasVideoTrack() {

		trackID := "video"
		source := player.player.GetVideoSource()

		// todo  fix source
		incoming := newIncomingStreamTrack("video", trackID, nil, map[string]RTPIncomingSourceGroup{"": source})

		// todo event
		player.tracks[trackID] = incoming
	}

	if player.player.HasAudioTrack() {

		trackID := "audio"
		source := player.player.GetAudioSource()

		incoming := newIncomingStreamTrack("audio", trackID, nil, map[string]RTPIncomingSourceGroup{"": source})

		// todo
		player.tracks[trackID] = incoming

	}

	// todo OnEnd callback

	return player, nil
}

func (p *Player) GetTracks() []*IncomingStreamTrack {
	tracks := []*IncomingStreamTrack{}
	for _, track := range p.tracks {
		tracks = append(tracks, track)
	}
	return tracks
}

func (p *Player) GetAudioTracks() []*IncomingStreamTrack {
	tracks := []*IncomingStreamTrack{}
	for _, track := range p.tracks {
		if track.GetMedia() == "audio" {
			tracks = append(tracks, track)
		}
	}
	return tracks
}

func (p *Player) GetVideoTracks() []*IncomingStreamTrack {
	tracks := []*IncomingStreamTrack{}
	for _, track := range p.tracks {
		if track.GetMedia() == "video" {
			tracks = append(tracks, track)
		}
	}
	return tracks
}

func (p *Player) Play() {

	if p.player != nil {
		p.player.Play()
	}
}

func (p *Player) Resume() {

	if p.player != nil {
		p.player.Play()
	}
}

func (p *Player) Pause() {

	if p.player != nil {
		p.player.Stop()
	}
}

func (p *Player) Seek(time uint64) {

	if p.player != nil {
		p.player.Seek(time)
	}
}

func (p *Player) Stop() {

	if p.player == nil {
		return
	}

	for _, track := range p.tracks {
		track.Stop()
	}

	p.tracks = nil

	p.player.Close()

	p.player = nil
}

// add fake interface, make Player can build,  todo
func (p *Player) Swigcptr() uintptr {
	return 0
}
func (p *Player) SwigIsPlayerListener() {

}
func (p *Player) OnEnd() {

}

func (p *Player) DirectorInterface() interface{} {

	return nil
}
