package mediaserver

import (
	native "github.com/notedit/media-server-go/wrapper"
)

func init() {
	native.MediaServerInitialize()
}

func EnableLog(flag bool) {
	native.MediaServerEnableLog(flag)
}

func EnableDebug(flag bool) {
	native.MediaServerEnableDebug(flag)
}

func SetPortRange(minPort, maxPort int) bool {
	return native.MediaServerSetPortRange(minPort, maxPort)
}

func EnableUltraDebug(flag bool) {
	native.MediaServerEnableUltraDebug(flag)
}
