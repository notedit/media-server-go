package sdp

type Direction string

const (
	DirectionSENDRECV Direction = "sendrecv"
	DirectionSENDONLY Direction = "sendonly"
	DirectionRECVONLY Direction = "recvonly"
	DirectionINACTIVE Direction = "inactive"
)

func DirectionbyValue(direction string) Direction {

	switch direction {
	case "sendrecv":
		return DirectionSENDRECV
	case "sendonly":
		return DirectionSENDONLY
	case "recvonly":
		return DirectionRECVONLY
	case "inactive":
		return DirectionINACTIVE
	default:
		return ""
	}

}

func (d Direction) Reverse() Direction {
	switch d {
	case DirectionSENDRECV:
		return DirectionSENDRECV
	case DirectionSENDONLY:
		return DirectionRECVONLY
	case DirectionRECVONLY:
		return DirectionSENDONLY
	case DirectionINACTIVE:
		return DirectionINACTIVE
	default:
		return ""
	}
}

func (d Direction) String() string {
	switch d {
	case DirectionSENDRECV:
		return "sendrecv"
	case DirectionSENDONLY:
		return "sendonly"
	case DirectionRECVONLY:
		return "recvonly"
	case DirectionINACTIVE:
		return "inactive"
	default:
		return ""
	}
}
