package mediaserver

type RenegotiationCallback func(transport *Transport)

type SDPManager interface {
	GetState() string
	GetTransport() *Transport
	CreateLocalDescription() string
	ProcessRemoteDescription() string
}
