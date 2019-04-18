package mediaserver

import (
	"github.com/notedit/sdp"
)

type RenegotiationCallback func(transport *Transport)

// SDPManager interface
type SDPManager interface {
	GetState() string
	GetTransport() *Transport
	CreateLocalDescription() (*sdp.SDPInfo, error)
	ProcessRemoteDescription(sdp string) (*sdp.SDPInfo, error)
}
