package sdp

type SimulcastStreamInfo struct {
	id     string
	paused bool
}

func NewSimulcastStreamInfo(id string, paused bool) *SimulcastStreamInfo {

	return &SimulcastStreamInfo{
		id:     id,
		paused: paused,
	}
}

func (s *SimulcastStreamInfo) Clone() *SimulcastStreamInfo {

	return &SimulcastStreamInfo{
		id:     s.id,
		paused: s.paused,
	}
}

func (s *SimulcastStreamInfo) IsPaused() bool {

	return s.paused
}

func (s *SimulcastStreamInfo) GetID() string {

	return s.id
}
