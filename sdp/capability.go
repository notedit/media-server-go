package sdp

import (
	"encoding/json"
)

type RtcpFeedback struct {
	ID     string   `json:"id,omitempty"`
	Params []string `json:"params,omitempty"`
}

type Capability struct {
	Codecs     []string        `json:"codecs"`
	Rtx        bool            `json:"rtx,omitempty"`
	Rtcpfbs    []*RtcpFeedback `json:"rtcpfbs,omitempty"`
	Extensions []string        `json:"extensions,omitempty"`
	Simulcast  bool            `json:"simulcast,omitempty"`
}

func CapabilityFromJSON(cap []byte) (*Capability, error) {
	var capability Capability
	err := json.Unmarshal(cap, &capability)
	if err != nil {
		return nil, err
	}
	return &capability, nil
}
