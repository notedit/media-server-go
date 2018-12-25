# media-server-go
WebRTC media server for go

|         | x86 | x64 |
|:------- |:--- |:--- |
| Linux   | -   | ✔︎   | 
| macOS   | -   | ✔︎   | 

only support Linux x64 and macOS x64 for now

## How to generate swig code  

swig -go -cgo -c++ -intgosize 64 mediaserver.i


## Install 

```
go get github.com/notedit/media-server-go
```


## How to build manually 

1,  clone the code

2,  bash build.sh

3,  go build 
