package sdp

type DirectionWay string

const (
	DirectionWaySEND DirectionWay = "send"
	DirectionWayRECV DirectionWay = "recv"
)

func DirectionWaybyValue(d string) DirectionWay {

	switch d {
	case "recv":
		return DirectionWayRECV
	case "send":
		return DirectionWaySEND
	default:
		return ""
	}

}

func (d DirectionWay) Reverse() DirectionWay {

	switch d {
	case DirectionWaySEND:
		return DirectionWayRECV
	case DirectionWayRECV:
		return DirectionWaySEND
	default:
		return ""
	}
}

func (d DirectionWay) String() string {

	switch d {
	case DirectionWaySEND:
		return "send"
	case DirectionWayRECV:
		return "recv"
	default:
		return ""
	}
}
