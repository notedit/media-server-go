package mediaserver

type encoding struct {
	id           string
	souce        interface{}
	depacketizer StreamTrackDepacketizer
}

type IncomingStreamTrack struct {
	id        string
	media     string
	receiver  RTPReceiverFacade
	counter   int
	encodings map[string]*encoding
}

type Stats struct {
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

func newIncomingStreamTrack(media string, id string, receiver RTPReceiverFacade, souces map[string]interface{}) *IncomingStreamTrack {
	track := &IncomingStreamTrack{}

	track.id = id
	track.media = media
	track.receiver = receiver
	track.counter = 0
	track.encodings = make(map[string]*encoding)

	for k, souce := range souces {
		track.encodings[k] = &encoding{
			id:           k,
			souce:        souce,
			depacketizer: nil, //NewStreamTrackDepacketizer()
		}
	}

	return track
}

func (i *IncomingStreamTrack) GetMedia() string {
	return i.media
}

func (i *IncomingStreamTrack) GetStats() *Stats {
	return nil
}

func (i *IncomingStreamTrack) GetActiveLayers() {

}

func (i *IncomingStreamTrack) Stop() {
}
