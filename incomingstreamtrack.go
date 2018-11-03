package mediaserver

type IncomingStreamTrack struct {
	id    string
	media string
}

func (i *IncomingStreamTrack) GetMedia() string {
	return i.media
}

func (i *IncomingStreamTrack) Stop() {

}
