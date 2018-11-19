package sdp

type SimulcastInfo struct {
	send [][]*SimulcastStreamInfo
	recv [][]*SimulcastStreamInfo
}

func NewSimulcastInfo() *SimulcastInfo {
	info := &SimulcastInfo{}
	info.send = make([][]*SimulcastStreamInfo, 0)
	info.recv = make([][]*SimulcastStreamInfo, 0)
	return info
}

func (s *SimulcastInfo) Clone() *SimulcastInfo {
	cloned := new(SimulcastInfo)
	cloned.send = make([][]*SimulcastStreamInfo, len(s.send))
	cloned.recv = make([][]*SimulcastStreamInfo, len(s.recv))

	for i := range s.send {
		cloned.send[i] = make([]*SimulcastStreamInfo, len(s.send[i]))
		copy(cloned.send[i], s.send[i])
	}

	for i := range s.recv {
		cloned.recv[i] = make([]*SimulcastStreamInfo, len(s.recv[i]))
		copy(cloned.recv[i], s.recv[i])
	}
	return cloned
}

func (s *SimulcastInfo) AddSimulcastAlternativeStreams(direction DirectionWay, streams []*SimulcastStreamInfo) {

	if direction == SEND {
		s.send = append(s.send, streams)
	} else {
		s.recv = append(s.recv, streams)
	}
}

func (s *SimulcastInfo) AddSimulcastStream(direction DirectionWay, stream *SimulcastStreamInfo) {

	if direction == SEND {
		s.send = append(s.send, []*SimulcastStreamInfo{stream})
	} else {
		s.recv = append(s.recv, []*SimulcastStreamInfo{stream})
	}
}

func (s *SimulcastInfo) GetSimulcastStreams(direction DirectionWay) [][]*SimulcastStreamInfo {

	if direction == SEND {
		return s.send
	}

	return s.recv
}
