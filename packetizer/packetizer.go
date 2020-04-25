package packetizer

type Packetizer interface {
	Packetize(payload []byte, mtu int) [][]byte
}
