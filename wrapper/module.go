package native


/*
#cgo CXXFLAGS: -std=c++1z
#cgo CPPFLAGS: -I${SRCDIR}/../thirdparty/openssl/build/include/
#cgo CPPFLAGS: -I${SRCDIR}/../thirdparty/libsrtp/build/include/
#cgo CPPFLAGS: -I${SRCDIR}/../thirdparty/mp4v2/build/include/
#cgo CPPFLAGS: -I${SRCDIR}/../media-server/ext/crc32c/include/
#cgo CPPFLAGS: -I${SRCDIR}/../media-server/ext/libdatachannels/src/
#cgo CPPFLAGS: -I${SRCDIR}/../media-server/ext/libdatachannels/src/internal/
#cgo CPPFLAGS: -I${SRCDIR}/../media-server/include/
#cgo CPPFLAGS: -I${SRCDIR}/../media-server/src/
#cgo LDFLAGS: -L${SRCDIR}/../media-server/bin/release  -lmediaserver
#cgo LDFLAGS: -L${SRCDIR}/../thirdparty/openssl/build  -lssl -lcrypto
#cgo LDFLAGS: -L${SRCDIR}/../thirdparty/libsrtp/build  -lsrtp2
#cgo LDFLAGS: -L${SRCDIR}/../thirdparty/mp4v2/build  -lmp4v2
#cgo LDFLAGS: -ldl
*/
import "C"
