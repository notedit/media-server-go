package native

/*
#cgo CXXFLAGS: -std=c++1z
#cgo CPPFLAGS: -I/usr/local/include
#cgo CPPFLAGS: -I${SRCDIR}/../include/crc32c/include/
#cgo CPPFLAGS: -I${SRCDIR}/../include/libdatachannels/
#cgo CPPFLAGS: -I${SRCDIR}/../include/libdatachannels/internal/
#cgo CPPFLAGS: -I${SRCDIR}/../include/media-server/include/
#cgo CPPFLAGS: -I${SRCDIR}/../include/media-server/src/
#cgo LDFLAGS: -L/usr/local/lib -lmediaserver -lssl -lcrypto -lsrtp2
#cgo LDFLAGS: /usr/local/lib/libmp4v2.a
#cgo LDFLAGS: -ldl
*/
import "C"
