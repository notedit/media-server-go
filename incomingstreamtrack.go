package mediaserver

import "github.com/chuckpreslar/emission"

type Encoding struct {
	id           string
	source       RTPIncomingSourceGroup
	depacketizer StreamTrackDepacketizer
}

func (e *Encoding) GetID() string {
	return e.id
}

func (e *Encoding) GetSource() RTPIncomingSourceGroup {
	return e.source
}

func (e *Encoding) GetDepacketizer() StreamTrackDepacketizer {
	return e.depacketizer
}

type IncomingStreamTrack struct {
	id        string
	media     string
	receiver  RTPReceiverFacade
	counter   int
	encodings map[string]*Encoding
	stats     *IncomingStats // buffer the last stats
	*emission.Emitter
}

type IncomingStats struct {
	LostPackets    int
	DropPackets    int
	NumPackets     int
	NumRTCPPackets int
	TotalBytes     int
	TotalRTCPBytes int
	TotalPLIs      int
	TotalNACKs     int
	Bitrate        int
}

func newIncomingStreamTrack(media string, id string, receiver RTPReceiverFacade, souces map[string]RTPIncomingSourceGroup) *IncomingStreamTrack {
	track := &IncomingStreamTrack{}

	track.id = id
	track.media = media
	track.receiver = receiver
	track.counter = 0
	track.encodings = make([]*Encoding)
	track.Emitter = emission.NewEmitter()

	for k, source := range souces {
		encoding := &Encoding{
			id:           k,
			source:       source,
			depacketizer: NewStreamTrackDepacketizer(source),
		}
		track.encodings = append(track.encodings, encoding)
	}

	return track
}

func (i *IncomingStreamTrack) GetID() string {
	return i.id
}

func (i *IncomingStreamTrack) GetMedia() string {
	return i.media
}

func (i *IncomingStreamTrack) GetSSRCs() []map[string]RTPIncomingSource {

	ssrcs := make([]map[string]RTPIncomingSource, 0)

	for _, encoding := range i.encodings {
		ssrcs = append(ssrcs, map[string]RTPIncomingSource{
			"media": encoding.source.GetMedia(),
			"rtx":   encoding.source.GetRtx(),
			"fec":   encoding.source.GetFec(),
		})
	}
	return ssrcs
}

func (i *IncomingStreamTrack) GetStats() *IncomingStats {

	// todo
	return nil
}

func (i *IncomingStreamTrack) GetActiveLayers() {

	// todo
}

func (i *IncomingStreamTrack) GetEncodings() []*Encoding {

	return i.encodings
}

func (i *IncomingStreamTrack) Attached() {

	i.counter = i.counter + 1

	if i.counter == 1 {
		i.EmitSync("attached")
	}
}

func (i *IncomingStreamTrack) Refresh() {

	for _, encoding := range i.encodings {
		//Request an iframe on main ssrc
		i.receiver.SendPLI(encoding.source.GetMedia().GetSsrc())
	}
}

func (i *IncomingStreamTrack) Detached() {

	i.counter = i.counter - 1

	if i.counter == 0 {
		i.EmitSync("detached")
	}
}

func (i *IncomingStreamTrack) Stop() {

	if i.receiver == nil {
		return
	}

	for _, encoding := range i.encodings {
		if encoding.depacketizer != nil {
			encoding.depacketizer.Stop()
		}
	}

	i.EmitSync("stopped")

	i.encodings = nil

	i.receiver = nil
}
