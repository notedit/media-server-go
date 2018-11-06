package mediaserver

import "github.com/chuckpreslar/emission"

type encoding struct {
	id           int
	souce        RTPIncomingSourceGroup
	depacketizer StreamTrackDepacketizer
}

type IncomingStreamTrack struct {
	id        string
	media     string
	receiver  RTPReceiverFacade
	counter   int
	encodings map[int]*encoding
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

func newIncomingStreamTrack(media string, id string, receiver RTPReceiverFacade, souces []RTPIncomingSourceGroup) *IncomingStreamTrack {
	track := &IncomingStreamTrack{}

	track.id = id
	track.media = media
	track.receiver = receiver
	track.counter = 0
	track.encodings = make(map[int]*encoding)
	track.Emitter = emission.NewEmitter()

	for k, souce := range souces {
		track.encodings[k] = &encoding{
			id:           k,
			souce:        souce,
			depacketizer: nil, //NewStreamTrackDepacketizer()
		}
	}

	return track
}

func (i *IncomingStreamTrack) GetID() string {
	return i.id
}

func (i *IncomingStreamTrack) GetMedia() string {
	return i.media
}

func (i *IncomingStreamTrack) GetSSRCs() {

}

func (i *IncomingStreamTrack) GetStats() *IncomingStats {
	return nil
}

func (i *IncomingStreamTrack) GetActiveLayers() {

}

func (i *IncomingStreamTrack) Refresh() {

}

func (i *IncomingStreamTrack) Stop() {
}
