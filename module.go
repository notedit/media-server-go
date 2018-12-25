package mediaserver

/*
#cgo CPPFLAGS: -I${SRCDIR}/external/libsrtp/include
#cgo CPPFLAGS: -I${SRCDIR}/external/mp4v2/include
#cgo CPPFLAGS: -I${SRCDIR}/media-server/include
#cgo darwin,amd64 CPPFLAGS: -I${SRCDIR}/external/opensslinclude/darwin-amd64
#cgo darwin,amd64 LDFLAGS: -L${SRCDIR}/lib -lmediaserver-darwin-amd64  -lsrtp2-darwin-amd64  -lcrypto-darwin-amd64 -lssl-darwin-amd64 -lmp4v2-darwin-amd64
#cgo linux,amd64 CPPFLAGS: -I${SRCDIR}/external/opensslinclude/linux-amd64
#cgo linux,amd64 LDFLAGS: -L${SRCDIR}/lib -lmediaserver-linux-amd64  -lsrtp2-linux-amd64  -lcrypto-linux-amd64 -lssl-linux-amd64 -lmp4v2-linux-amd64
*/
import "C"

func init() {
	MediaServerInitialize()
}

func EnableLog(flag bool) {
	MediaServerEnableLog(flag)
}

func EnableDebug(flag bool) {
	MediaServerEnableDebug(flag)
}

func SetPortRange(minPort, maxPort int) bool {
	return MediaServerSetPortRange(minPort, maxPort)
}

func EnableUltraDebug(flag bool) {
	MediaServerEnableUltraDebug(flag)
}
