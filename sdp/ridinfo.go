package sdp

import (
	"strconv"
)

type RIDInfo struct {
	id        string
	direction DirectionWay // TODO
	formats   []int
	params    map[string]string
}

func NewRIDInfo(id string, direction DirectionWay) *RIDInfo {

	ridInfo := &RIDInfo{
		id:        id,
		direction: direction,
		formats:   []int{},
		params:    map[string]string{},
	}
	return ridInfo
}

func (r *RIDInfo) Clone() *RIDInfo {

	ridInfo := &RIDInfo{}
	ridInfo.id = r.id
	ridInfo.direction = r.direction
	ridInfo.formats = make([]int, len(r.formats))
	ridInfo.params = make(map[string]string)
	copy(ridInfo.formats, r.formats)
	for k, v := range r.params {
		ridInfo.params[k] = v
	}
	return ridInfo
}

func (r *RIDInfo) GetID() string {
	return r.id
}

func (r *RIDInfo) GetDirection() DirectionWay {
	return r.direction
}

func (r *RIDInfo) SetDirection(direction DirectionWay) {
	r.direction = direction
}

func (r *RIDInfo) GetFormats() []int {
	return r.formats
}

func (r *RIDInfo) SetFormats(formats []string) {
	r.formats = []int{}
	for _, format := range formats {
		formatInt, err := strconv.Atoi(format)
		if err != nil {
			r.formats = append(r.formats, formatInt)
		}
	}
}

func (r *RIDInfo) GetParams() map[string]string {
	return r.params
}

func (r *RIDInfo) SetParams(params map[string]string) {
	r.params = params
}

func (r *RIDInfo) AddParam(id, param string) {
	r.params[id] = param
}
