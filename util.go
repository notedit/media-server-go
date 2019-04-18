package mediaserver

import (
	"errors"
)

const ssrcMin uint = 1000000000
const ssrcMax uint = 4294967295

var ssrcValue = ssrcMin

func Min(x, y int) int {
	if x < y {
		return x
	}
	return y
}

func Max(x, y int) int {
	if x < y {
		return y
	}
	return x
}

func NextSSRC() uint {

	if ssrcValue == ssrcMax {
		ssrcValue = ssrcMin
	}
	ssrcValue = ssrcValue + 1
	return ssrcValue
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

var nalu_prefix = []byte{0, 0, 0, 1}

func annexbConvert(avc []byte) ([]byte, error) {
	if len(avc) < 4 {
		return nil, errors.New("too short")
	}
	val4 := u32be(avc)
	_val4 := val4
	_b := avc[4:]
	annexb := []byte{}

	for {
		annexb = append(annexb, nalu_prefix...)
		annexb = append(annexb, _b[:_val4]...)
		_b = _b[_val4:]
		if len(_b) < 4 {
			break
		}
		_val4 = u32be(_b)
		_b = _b[4:]
		if _val4 > uint32(len(_b)) {
			break
		}
	}

	return annexb, nil
}
