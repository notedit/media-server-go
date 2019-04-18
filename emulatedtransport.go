package mediaserver

import (
	"fmt"
	"sync"

	native "github.com/notedit/media-server-go/wrapper"
	"github.com/notedit/sdp"
)

// EmulatedTransport pcap file as a transport
type EmulatedTransport struct {
	transport                native.PCAPTransportEmulator
	streams                  map[string]*IncomingStream
	onIncomingTrackListeners []IncomingTrackListener
	sync.Mutex
}

// NewEmulatedTransport create a transport by pcap file
func NewEmulatedTransport(pcap string) *EmulatedTransport {
	transport := &EmulatedTransport{}
	transport.transport = native.NewPCAPTransportEmulator()
	transport.streams = map[string]*IncomingStream{}
	transport.transport.Open(pcap)
	transport.onIncomingTrackListeners = make([]IncomingTrackListener, 0)
	return transport
}

func (e *EmulatedTransport) SetRemoteProperties(audio *sdp.MediaInfo, video *sdp.MediaInfo) {
	properties := native.NewProperties()
	if audio != nil {
		num := 0
		for _, codec := range audio.GetCodecs() {
			item := fmt.Sprintf("audio.codecs.%d", num)
			properties.SetProperty(item+".codec", codec.GetCodec())
			properties.SetProperty(item+".pt", codec.GetType())
			if codec.HasRTX() {
				properties.SetProperty(item+".rtx", codec.GetRTX())
			}
			num = num + 1
		}
		properties.SetProperty("audio.codecs.length", num)

		num = 0
		for id, uri := range audio.GetExtensions() {
			item := fmt.Sprintf("audio.ext.%d", num)
			properties.SetProperty(item+".id", id)
			properties.SetProperty(item+".uri", uri)
			num = num + 1
		}
		properties.SetProperty("audio.ext.length", num)
	}

	if video != nil {
		num := 0
		for _, codec := range video.GetCodecs() {
			item := fmt.Sprintf("video.codecs.%d", num)
			properties.SetProperty(item+".codec", codec.GetCodec())
			properties.SetProperty(item+".pt", codec.GetType())
			if codec.HasRTX() {
				properties.SetProperty(item+".rtx", codec.GetRTX())
			}
			num = num + 1
		}
		properties.SetProperty("video.codecs.length", num)

		num = 0
		for id, uri := range video.GetExtensions() {
			item := fmt.Sprintf("video.ext.%d", num)
			properties.SetProperty(item+".id", id)
			properties.SetProperty(item+".uri", uri)
			num = num + 1
		}
		properties.SetProperty("video.ext.length", num)
	}

	e.transport.SetRemoteProperties(properties)

	native.DeleteProperties(properties)
}

// CreateIncomingStream create incoming stream base on streaminfo
func (e *EmulatedTransport) CreateIncomingStream(streamInfo *sdp.StreamInfo) *IncomingStream {

	incomingStream := NewIncomingStreamWithEmulatedTransport(e.transport, native.PCAPTransportEmulatorToReceiver(e.transport), streamInfo)

	e.Lock()
	e.streams[incomingStream.GetID()] = incomingStream
	e.Unlock()

	incomingStream.OnStop(func() {
		e.Lock()
		delete(e.streams, incomingStream.GetID())
		e.Unlock()
	})

	incomingStream.OnTrack(func(track *IncomingStreamTrack) {
		for _, addTrackFunc := range e.onIncomingTrackListeners {
			addTrackFunc(track, incomingStream)
		}
	})

	for _, track := range incomingStream.GetTracks() {
		for _, addTrackFunc := range e.onIncomingTrackListeners {
			addTrackFunc(track, incomingStream)
		}
	}
	return incomingStream
}

// OnIncomingTrack register incoming track
func (e *EmulatedTransport) OnIncomingTrack(listener IncomingTrackListener) {
	e.onIncomingTrackListeners = append(e.onIncomingTrackListeners, listener)
}

// Play play at start time
func (e *EmulatedTransport) Play(time uint64) bool {
	e.transport.Seek(time)
	return e.transport.Play()
}

// Resume resume to play
func (e *EmulatedTransport) Resume() bool {
	return e.transport.Play()
}

// Pause  pause
func (e *EmulatedTransport) Pause() bool {
	return e.transport.Stop()
}

// Seek to some time
func (e *EmulatedTransport) Seek(time uint64) bool {
	e.transport.Seek(time)
	return e.transport.Play()
}

// Stop stop this transport
func (e *EmulatedTransport) Stop() {

	if e.transport == nil {
		return
	}

	for _, stream := range e.streams {
		stream.Stop()
	}

	e.transport.Stop()
	native.DeletePCAPTransportEmulator(e.transport)

	e.streams = nil
	e.transport = nil

}
