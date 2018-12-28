package native

/*
#cgo CXXFLAGS: -std=c++0x
#cgo CPPFLAGS: -I${SRCDIR}/../include/srtp/
#cgo CPPFLAGS: -I${SRCDIR}/../include/mp4v2/
#cgo CPPFLAGS: -I${SRCDIR}/../media-server/include
#cgo darwin,amd64 CPPFLAGS: -I${SRCDIR}/../include/openssl/darwin-amd64
#cgo darwin,amd64 LDFLAGS: -L${SRCDIR}/../lib -lmediaserver-darwin-amd64  -lssl-darwin-amd64  -lsrtp2-darwin-amd64  -lcrypto-darwin-amd64 -lmp4v2-darwin-amd64
#cgo linux,amd64 CPPFLAGS: -I${SRCDIR}/../include/openssl/linux-amd64
#cgo linux,amd64 LDFLAGS: -L${SRCDIR}/../lib -lmediaserver-linux-amd64 -lssl-linux-amd64 -lcrypto-linux-amd64 -lsrtp2-linux-amd64 -lmp4v2-linux-amd64
#cgo LDFLAGS: -ldl
*/
import "C"
