package native

/*
#cgo CXXFLAGS: -std=c++1z
#cgo CPPFLAGS: -I/usr/local/include:/usr/local/ssl/include/
#cgo CPPFLAGS: -I${SRCDIR}/../include/crc32c/include/
#cgo CPPFLAGS: -I${SRCDIR}/../include/libdatachannels/
#cgo CPPFLAGS: -I${SRCDIR}/../include/libdatachannels/internal/
#cgo CPPFLAGS: -I${SRCDIR}/../include/media-server/include/

#cgo LDFLAGS: -L/usr/local/ssl/lib/ -lssl  -lcrypto
#cgo LDFLAGS: -L/usr/local/lib  -lmp4v2 -lsrtp2 -lmediaserver
#cgo LDFLAGS: -ldl
*/
import "C"
