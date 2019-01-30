package mediaserver

import (
	"sync"
	"time"
)

type Refresher struct {
	period int
	tracks map[string]*IncomingStreamTrack
	ticker *time.Ticker
	sync.Mutex
}

func NewRefresher(period int) *Refresher {
	refresher := &Refresher{}
	refresher.tracks = map[string]*IncomingStreamTrack{}
	refresher.period = period

	return refresher
}

func (r *Refresher) Add(incom *IncomingStreamTrack) {

	if incom.GetMedia() == "video" {

		r.Lock()
		r.tracks[incom.GetID()] = incom
		r.Unlock()

		incom.OnStop(func() {
			r.Lock()
			delete(r.tracks, incom.GetID())
			r.Unlock()
		})
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
	r.tracks = nil
}

