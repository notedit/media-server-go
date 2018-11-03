package sdp

type Setup string

const (
	SETUPACTIVE   Setup = "active"
	SETUPPASSIVE  Setup = "passive"
	SETUPACTPASS  Setup = "actpass"
	SETUPINACTIVE Setup = "inactive"
)

func SetupByValue(s string) Setup {

	switch s {
	case "active":
		return SETUPACTIVE
	case "passive":
		return SETUPPASSIVE
	case "actpass":
		return SETUPACTPASS
	case "inactive":
		return SETUPINACTIVE
	default:
		return ""
	}
}

func (s Setup) Reverse() Setup {

	switch s {
	case SETUPACTIVE:
		return SETUPPASSIVE
	case SETUPPASSIVE:
		return SETUPACTIVE
	case SETUPACTPASS:
		return SETUPPASSIVE
	case SETUPINACTIVE:
		return SETUPINACTIVE
	default:
		return ""
	}
}

func (s Setup) String() string {

	switch s {
	case SETUPACTIVE:
		return "active"
	case SETUPPASSIVE:
		return "passive"
	case SETUPACTPASS:
		return "actpass"
	case SETUPINACTIVE:
		return "inactive"
	default:
		return ""
	}
}
