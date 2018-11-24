package sdp

import (
	"encoding/hex"
	"math/rand"
)

type ICEInfo struct {
	ufrag           string
	password        string
	lite            bool
	endOfCandidates bool
}

func NewICEInfo(ufrag, password string) *ICEInfo {
	return &ICEInfo{
		ufrag:           ufrag,
		password:        password,
		lite:            false,
		endOfCandidates: false,
	}
}

func GenerateICEInfo(lite bool) *ICEInfo {

	ufrag := make([]byte, 8)
	password := make([]byte, 24)
	rand.Read(ufrag)
	rand.Read(password)

	return &ICEInfo{
		ufrag:           hex.EncodeToString(ufrag),
		password:        hex.EncodeToString(password),
		lite:            lite,
		endOfCandidates: false,
	}
}

func ICEInfoGenerate(lite bool) *ICEInfo {

	ufrag := make([]byte, 8)
	password := make([]byte, 24)
	rand.Read(ufrag)
	rand.Read(password)

	return &ICEInfo{
		ufrag:           hex.EncodeToString(ufrag),
		password:        hex.EncodeToString(password),
		lite:            lite,
		endOfCandidates: false,
	}
}

func (c *ICEInfo) Clone() *ICEInfo {

	return &ICEInfo{
		ufrag:           c.ufrag,
		password:        c.password,
		lite:            c.lite,
		endOfCandidates: c.endOfCandidates,
	}
}

func (c *ICEInfo) GetUfrag() string {
	return c.ufrag
}

func (c *ICEInfo) GetPassword() string {
	return c.password
}

func (c *ICEInfo) IsLite() bool {
	return c.lite
}

func (c *ICEInfo) SetLite(lite bool) {
	c.lite = lite
}

func (c *ICEInfo) IsEndOfCandidate() bool {
	return c.endOfCandidates
}

func (c *ICEInfo) SetEndOfCandidate(endOfCandidate bool) {
	c.endOfCandidates = endOfCandidate
}
