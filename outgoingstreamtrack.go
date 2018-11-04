package mediaserver

type OutgoingStreamTrack struct {
	id     string
	media  string
	sender RTPSenderFacade
	source interface{}
	muted  bool
}

type OutgoingStats struct {
	NumPackets     int
	NumRTCPPackets int
	TotalBytes     int
	TotalRTCPBytes int
	Bitrate        int
}

func NewOutgoingStreamTrack(id string, media string, sender RTPSenderFacade, source interface{}) *OutgoingStreamTrack {

	track := &OutgoingStreamTrack{}
	track.id = id
	track.media = media
	track.sender = sender
	track.muted = false
	track.source = source

	// todo onremb callback
	return track
}

func (o *OutgoingStreamTrack) GetID() string {
	return o.id
}

func (o *OutgoingStreamTrack) GetMedia() string {
	return o.media
}

func (o *OutgoingStreamTrack) GetStats() {

}

func (o *OutgoingStreamTrack) GetSSRCs() {

}

func (o *OutgoingStreamTrack) IsMuted() bool {
	return o.muted
}

func (o *OutgoingStreamTrack) Mute(muting bool) {

}

func (o *OutgoingStreamTrack) AttachTo(incomingTrack *IncomingStreamTrack) {

}

// todo

func (o *OutgoingStreamTrack) Stop() {

}
