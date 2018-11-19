package sdp

type Direction uint8

const (
	SENDRECV Direction = iota + 1
	SENDONLY
	RECVONLY
	INACTIVE
)

func DirectionbyValue(direction string) Direction {

	switch direction {
	case "sendrecv":
		return SENDRECV
	case "sendonly":
		return SENDONLY
	case "recvonly":
		return RECVONLY
	case "inactive":
		return INACTIVE
	default:
		return 0
	}

}

func (d Direction) Reverse() Direction {
	switch d {
	case SENDRECV:
		return SENDRECV
	case SENDONLY:
		return RECVONLY
	case RECVONLY:
		return SENDONLY
	case INACTIVE:
		return INACTIVE
	default:
		return 0
	}
}

func (d Direction) String() string {
	switch d {
	case SENDRECV:
		return "sendrecv"
	case SENDONLY:
		return "sendonly"
	case RECVONLY:
		return "recvonly"
	case INACTIVE:
		return "inactive"
	default:
		return ""
	}
}
