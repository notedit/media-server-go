package mediaserver

import "C"
import (
	"fmt"
	"unsafe"

	native "github.com/notedit/media-server-go/wrapper"
)

// MediaStreamDuplicater we can make a copy of the incoming stream and callback the mediaframe data
type MediaFrameMultiplexer struct {
	track      *IncomingStreamTrack
	multiplexer native.MediaFrameMultiplexer
	listener   mediaframeListener // used for native wrapper, see swig's doc

	mediaframeListener func([]byte, uint) // used for outside
}


type mediaframeListener interface {
	native.MediaFrameListener
	deleteMediaFrameListener()
}

type goMediaFrameListener struct {
	native.MediaFrameListener
}

func (m *goMediaFrameListener) deleteMediaFrameListener() {
	native.DeleteDirectorMediaFrameListener(m.MediaFrameListener)
}

type overwrittenMediaFrameListener struct {
	p          native.MediaFrameListener
	multiplexer *MediaFrameMultiplexer
}

func (p *overwrittenMediaFrameListener) OnMediaFrame(frame native.MediaFrame) {

	if p.multiplexer != nil && p.multiplexer.mediaframeListener != nil {
		buffer := C.GoBytes(unsafe.Pointer(frame.GetData()), C.int(frame.GetLength()))
		if frame.GetType() == native.MediaFrameVideo {
			data, err := annexbConvert(buffer)
			if err == nil {
				p.multiplexer.mediaframeListener(data, frame.GetTimeStamp())
			} else {
				fmt.Println(err)
			}
		} else {
			p.multiplexer.mediaframeListener(buffer, frame.GetTimeStamp())
		}

	}
}

// NewMediaStreamDuplicater duplicate this IncomingStreamTrack and callback the mediaframe
func NewMediaFrameMultiplexer(track *IncomingStreamTrack) *MediaFrameMultiplexer {

	duplicater := &MediaFrameMultiplexer{}
	duplicater.track = track

	// We should make sure this source is the main source
	source := track.GetFirstEncoding().GetSource()
	duplicater.multiplexer = native.NewMediaFrameMultiplexer(source)

	listener := &overwrittenMediaFrameListener{
		multiplexer: duplicater,
	}
	p := native.NewDirectorMediaFrameListener(listener)
	listener.p = p

	duplicater.listener = &goMediaFrameListener{MediaFrameListener: p}

	duplicater.multiplexer.AddMediaListener(duplicater.listener)

	return duplicater
}

// SetMediaFrameListener set outside mediaframe listener
func (d *MediaFrameMultiplexer) SetMediaFrameListener(listener func([]byte, uint)) {
	d.mediaframeListener = listener
}

// Stop stop this
func (d *MediaFrameMultiplexer) Stop() {

	if d.track == nil {
		return
	}

	if d.listener != nil {
		d.multiplexer.RemoveMediaListener(d.listener)
		d.listener.deleteMediaFrameListener()
	}

	d.track = nil
}
