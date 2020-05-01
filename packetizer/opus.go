package packetizer

type OpusPacketier struct{}

func (p *OpusPacketier) Packetize(payload []byte, mtu int) [][]byte {

	if payload == nil {
		return [][]byte{}
	}

	out := make([]byte, len(payload))
	copy(out, payload)
	return [][]byte{out}
}
