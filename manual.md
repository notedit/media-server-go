# Tutorial
This document will show you how to setup and endpoint and transport manually.

## Initialization

First import both media server go and sdp module:

```go
import "github.com/notedit/media-server-go"
import "github.com/notedit/media-server-go/sdp"
```


Then you need to create an Endpoint, which will create an UDP socket for receiving connection. You need to pass the `ip` address in which this Enpoint will be accesible by the WebRTC clients. This is typically the public IP address of the server, ans will be used on the ICE candidates sent to the browser on the SDP.

You can also pass a `port` with a `ip`, see `NewEndpointWithPort(ip string, port int)`, this endpoint will listen on this port. 

```go
//Create UDP server endpoint
endpoint := mediaserver.NewEndpoint("127.0.0.1")
// Or
endpoint := mediaserver.NewEndpointWithPort("127.0.0.1", 50000) 
```

Now you are ready to connect to your server.

## Connect a client

On your browser, create an SDP offer and sent it to your server (via websockets for example). Once you have it, you will have to parse it to extract the requried information. 
With that information, you can create an ICE+DTLS transport on the `Endpoint`.

```go
//Process the sdp
offer, err := sdp.Parse(offerStr)

//Create an DTLS ICE transport in that enpoint
transport = endpoint.CreateTransport(offer, nil)

```

Now set the RTP remote properties for both audio and video:

```go	
//Set RTP remote properties
transport.SetRemoteProperties(offer.GetMedia("audio"), offer.GetMedia("video"))
```

You can start creating the answer now. First get the ICE and DTLS info from the `Transport` and the ICE candidate into from the `Endpoint`

```go
//Get local DTLS and ICE info
ice := transport.GetLocalICEInfo()
dtls := transport.GetLocalDTLSInfo()
candidates := endpoint.GetLocalCandidates()

answer := sdp.NewSDPInfo()

//Add ice and dtls info
answer.SetDTLS(dtls)
answer.SetICE(ice)
//Add candidates
answer.AddCandidates(candidates)
```

Choose your codecs and set RTP parameters to answer the offer:
 
```go
//Get remote audio m-line info 
if offer.GetMedia("audio") != nil {
	audioMedia := offer.GetMedia("audio").AnswerCapability(audioCapability)
	answer.AddMedia(audioMedia)
}

//Get remote video m-line info 
const videoOffer = offer.getMedia("video");

//If offer had video
if offer.GetMedia("video") != nil {
	videoMedia := offer.GetMedia("video").AnswerCapability(videoCapability)
	answer.AddMedia(videoMedia)
}
```

Set the our negotiated RTP properties on the transport

```go
//Set RTP local  properties
transport.SetLocalProperties(answer.GetMedia("audio"), answer.GetMedia("video"))
```

## Stream management

You need to process the stream offered by the client, so extract the stream info from the SDP offer, and create an `IncomingStream` object.

```go
//Get stream 
for _, stream := range offer.GetStreams() {

	incomingStream := transport.CreateIncomingStream(stream)
	outgoingStream := transport.CreateOutgoingStream2(stream.Clone())

	outgoingStream.AttachTo(incomingStream)

	answer.AddStream(outgoingStream.GetStreamInfo())
}
//Create the remote stream into the transport
const incomingStream = transport.createIncomingStream(offered);
```

Now, for example, create an outgoing stream, and add it to the answer so the browser is aware of it.

```go
//Create new local stream
outgoingStream := transport.CreateOutgoingStream2(stream.Clone())
//Add local stream info it to the answer
answer.AddStream(outgoingStream.GetStreamInfo())
```

You can attach an `OutgoingStream` to an `IncomingStream`, this will create a `Transponder` array that will forward the incoming data to the ougoing stream, it will allow you also to apply transoformations to it (like SVC layer selection).

In this case, as you are attaching an incoming stream to an outgoing stream from the same client, you will get audio and video loopback on the client.

```go
//Copy incoming data from the remote stream to the local one
outgoingStream.AttachTo(incomingStream)
```

You can now send answer the SDP to the client.
```go
//Get answer SDP
const str = answer.toString()
```