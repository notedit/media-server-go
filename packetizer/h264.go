package packetizer


type H264Packetier struct{}

func (p *H264Packetier) Packetize(payload []byte, mtu int) (payloads [][]byte) {

	if payload == nil {
		return
	}

	nalus := splitnalus(payload)
	if nalus == nil {
		return
	}

	for _, nalu := range nalus {
		naluType := nalu[0] & 0x1F
		naluRefIdc := nalu[0] & 0x60

		if !(naluType < 9 && naluType > 4) {
			continue
		}

		if len(nalu) <= mtu {
			out := make([]byte, len(nalu))
			copy(out, nalu)
			payloads = append(payloads, out)
			continue
		}

		// remove the fua header
		maxFragmentSize := mtu - 2

		naluData := nalu

		naluDataIndex := 1
		naluDataLength := len(nalu) - naluDataIndex
		naluDataRemaining := naluDataLength

		if min(maxFragmentSize, naluDataRemaining) <= 0 {
			continue
		}

		for naluDataRemaining > 0 {
			currentFragmentSize := min(maxFragmentSize, naluDataRemaining)
			out := make([]byte, 2+currentFragmentSize)

			// +---------------+
			// |0|1|2|3|4|5|6|7|
			// +-+-+-+-+-+-+-+-+
			// |F|NRI|  Type   |
			// +---------------+
			out[0] = 28
			out[0] |= naluRefIdc

			// +---------------+
			//|0|1|2|3|4|5|6|7|
			//+-+-+-+-+-+-+-+-+
			//|S|E|R|  Type   |
			//+---------------+

			out[1] = naluType
			if naluDataRemaining == naluDataLength {
				// Set start bit
				out[1] |= 1 << 7
			} else if naluDataRemaining-currentFragmentSize == 0 {
				// Set end bit
				out[1] |= 1 << 6
			}

			copy(out[2:], naluData[naluDataIndex:naluDataIndex+currentFragmentSize])
			payloads = append(payloads, out)

			naluDataRemaining -= currentFragmentSize
			naluDataIndex += currentFragmentSize
		}

	}

	return payloads
}

func splitnalus(b []byte) (nalus [][]byte) {

	if len(b) < 4 {
		return
	}

	val3 := u24be(b)
	val4 := u32be(b)

	if val3 == 1 || val4 == 1 {
		_val3 := val3
		_val4 := val4
		start := 0
		pos := 0
		for {
			if start != pos {
				nalus = append(nalus, b[start:pos])
			}
			if _val3 == 1 {
				pos += 3
			} else if _val4 == 1 {
				pos += 4
			}
			start = pos
			if start == len(b) {
				break
			}
			_val3 = 0
			_val4 = 0

			for pos < len(b) {
				if pos+2 < len(b) && b[pos] == 0 {
					_val3 = u24be(b[pos:])
					if _val3 == 0 {
						if pos+3 < len(b) {
							_val4 = uint32(b[pos+3])
							if _val4 == 1 {
								break
							}
						}
					} else if _val3 == 1 {
						break
					}
					pos++
				} else {
					pos++
				}
			}

		}
	}
	return
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

func u32be(b []byte) (i uint32) {
	i = uint32(b[0])
	i <<= 8
	i |= uint32(b[1])
	i <<= 8
	i |= uint32(b[2])
	i <<= 8
	i |= uint32(b[3])
	return
}

func u24be(b []byte) (i uint32) {
	i = uint32(b[0])
	i <<= 8
	i |= uint32(b[1])
	i <<= 8
	i |= uint32(b[2])
	return
}