# media-server-go

[![Build Status](https://travis-ci.com/notedit/media-server-go.svg?branch=master)](https://travis-ci.com/notedit/media-server-go)

WebRTC media server for go





## How to use 

[Read the Tutorial](https://github.com/notedit/media-server-go/blob/master/manual.md)


Yon can see the demos from here [Demos](https://github.com/notedit/media-server-go-demo)



## Examples

- [WebRTC-Broadcast](https://github.com/notedit/media-server-go-demo/tree/master/broadcast): WebRTC publish and play 
- [Raw-RTP-Input](https://github.com/notedit/media-server-go-demo/tree/master/raw-rtp-input): Send raw rtp data into webrtc
- [WebRTC-Record](https://github.com/notedit/media-server-go-demo/tree/master/recording): WebRTC record
- [RTMP-To-WebRTC](https://github.com/notedit/media-server-go-demo/tree/master/rtmp-to-webrtc): Rtmp to webrtc
- [Server-To-Server](https://github.com/notedit/media-server-go-demo/tree/master/server-to-server): WebRTC server relay
- [WebRTC-To-RTMP](https://github.com/notedit/media-server-go-demo/tree/master/webrtc-to-rtmp): WebRTC to rtmp
- [WebRTC-To-HLS](https://github.com/notedit/media-server-go-demo/tree/master/webrtc-to-hls): WebRTC to hls



## Install 


`media-server-go` is not go getable, so you should clone it and build it yourself.

You should install `libtool` and `autoconf` `automake` before you build 


On ubuntu
```sh
apt install autoconf
apt install libtool
apt install automake
```


On macOS

```sh
brew install libtool
brew install autoconf
brew install automake
```


Your compiler should support `c++17`, for linux, you should update your `gcc/g++` to `7.0+`

for macos, clang should support `c++17`.


```sh
sudo add-apt-repository -y ppa:ubuntu-toolchain-r/test
sudo apt-get update -qq
sudo apt-get install g++-7
sudo update-alternatives --install /usr/bin/g++ g++ /usr/bin/g++-7 90
```


```sh
git clone --recurse-submodules https://github.com/notedit/media-server-go.git  

cd media-server-go

make

go install 

```


then you can use media-server-go in your project.




## Thanks 

 - [Media Server](https://github.com/medooze/media-server)
 - [Media Server for Node.js](https://github.com/medooze/media-server-node)






