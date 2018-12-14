package mediaserver

import (
	"time"

	"github.com/chuckpreslar/emission"
)

type Refresher struct {
	period int
	tracks map[string]*IncomingStreamTrack
	ticker *time.Ticker
	*emission.Emitter
}

func NewRefresher(period int) *Refresher {
	refresher := &Refresher{}
	refresher.tracks = map[string]*IncomingStreamTrack{}
	refresher.period = period
	refresher.Emitter = emission.NewEmitter()

	return refresher
}

func (r *Refresher) Add(incom *IncomingStreamTrack) {

	if incom.GetMedia() == "video" {
		r.tracks[incom.GetID()] = incom

		incom.Once("stopped", r.refresherOntrackStopped)
	}

	if r.ticker == nil {
		r.ticker = time.NewTicker(time.Duration(r.period) * time.Millisecond)
		go func() {
			for _ = range r.ticker.C {
				for _, track := range r.tracks {
					track.Refresh()
				}
			}
		}()
	}
}

func (r *Refresher) AddStream(incoming *IncomingStream) {

	for _, track := range incoming.GetTracks() {
		if track.GetMedia() == "video" {
			r.Add(track)
		}
	}
}

func (r *Refresher) Stop() {

	if r.ticker != nil {
		r.ticker.Stop()
		r.ticker = nil
	}

	for _, track := range r.tracks {
		track.Off("stopped", r.refresherOntrackStopped)
	}

	r.tracks = nil
}

func (r *Refresher) refresherOntrackStopped(track *IncomingStreamTrack) {

	delete(r.tracks, track.GetID())
}
