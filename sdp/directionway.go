package sdp

type DirectionWay uint8

const (
	SEND DirectionWay = iota + 1
	RECV
)

func DirectionWaybyValue(d string) DirectionWay {

	switch d {
	case "recv":
		return RECV
	case "send":
		return SEND
	default:
		return 0
	}

}

func (d DirectionWay) Reverse() DirectionWay {

	switch d {
	case SEND:
		return RECV
	case RECV:
		return SEND
	default:
		return 0
	}
}

func (d DirectionWay) String() string {

	switch d {
	case SEND:
		return "send"
	case RECV:
		return "recv"
	default:
		return ""
	}
}
