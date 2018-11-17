package mediaserver

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
