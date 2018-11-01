# media-server-go
WebRTC media server for go



## How to build

1,  build external mp4v2  openssl  libsrtp

2,  cp config.mk  mediaserver/  then build  make libmediaserver.a

3,  swig -go -cgo -c++ -intgosize 64 mediaserver.i

4,  go build 
