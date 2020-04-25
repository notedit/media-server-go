package packetizer



const (
	vp8HeaderSize = 1
)

type VP8Packetier struct{}

func (p *VP8Packetier) Packetize(payload []byte, mtu int) (payloads [][]byte) {

	maxFragmentSize := mtu - vp8HeaderSize

	payloadData := payload
	payloadDataRemaining := len(payload)

	payloadDataIndex := 0

	if min(maxFragmentSize, payloadDataRemaining) <= 0 {
		return payloads
	}

	for payloadDataRemaining > 0 {
		currentFragmentSize := min(maxFragmentSize, payloadDataRemaining)
		out := make([]byte, vp8HeaderSize+currentFragmentSize)
		if payloadDataRemaining == len(payload) {
			out[0] = 0x10
		}

		copy(out[vp8HeaderSize:], payloadData[payloadDataIndex:payloadDataIndex+currentFragmentSize])
		payloads = append(payloads, out)

		payloadDataRemaining -= currentFragmentSize
		payloadDataIndex += currentFragmentSize
	}

	return
}