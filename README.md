# media-server-go
WebRTC media server for go

|         | x86 | x64 |
|:------- |:--- |:--- |
| Linux   | -   | ✔︎   | 
| macOS   | -   | ✔︎   | 

only support Linux x64 and macOS x64 for now


## How to use 

Yon can see the demos from here [Demos](https://github.com/notedit/media-server-go-demo)


## Install 

```
go get github.com/notedit/media-server-go
```


## Thanks 

 - [Media Server](https://github.com/medooze/media-server)
 - [Media Server for Node.js](https://github.com/medooze/media-server-node)



## How to build manually 

you should install `libtool` and `autoconf` before you build 

ubuntu
```
apt install autoconf
apt install libtool
```
macOS
```
brew install libtool
brew install autoconf
```


- clone the code

- bash build.sh

- go build 
