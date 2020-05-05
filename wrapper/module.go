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
#cgo LDFLAGS: ${SRCDIR}/../media-server/bin/release/libmediaserver.a
#cgo LDFLAGS: ${SRCDIR}/../thirdparty/openssl/build/libssl.a
#cgo LDFLAGS: ${SRCDIR}/../thirdparty/openssl/build/libcrypto.a
#cgo LDFLAGS: ${SRCDIR}/../thirdparty/libsrtp/build/libsrtp2.a
#cgo LDFLAGS: ${SRCDIR}/../thirdparty/mp4v2/build/libmp4v2.a
#cgo LDFLAGS: -ldl
*/
import "C"
