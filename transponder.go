package mediaserver

import "github.com/chuckpreslar/emission"

type Transponder struct {
	muted       bool
	track       *IncomingStreamTrack
	transponder RTPStreamTransponderFacade
	*emission.Emitter
}

func NewTransponder(transponder RTPStreamTransponderFacade) *Transponder {

	return nil
}

func (t *Transponder) Mute(muting bool) {

}

func (t *Transponder) SetIncomingTrack(incomingTrack *IncomingStreamTrack) {

}

func (t *Transponder) Stop() {

}
