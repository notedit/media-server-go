package mediaserver

import "C"
import (
	"unsafe"

	native "github.com/notedit/media-server-go/wrapper"
)

// MediaStreamDuplicater we can make a copy of the incoming stream and callback the mediaframe data
type MediaStreamDuplicater struct {
	MediaFrames chan []byte
	track       *IncomingStreamTrack
	duplicater  native.MediaStreamDuplicaterFacade
	listener    mediaframeListener
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
	duplicater *MediaStreamDuplicater
}

func (p *overwrittenMediaFrameListener) OnMediaFrame(frame native.MediaFrame) {

	if p.duplicater != nil {
		buffer := C.GoBytes(unsafe.Pointer(frame.GetData()), C.int(frame.GetLength()))
		if frame.GetType() == native.MediaFrameVideo {
			data, err := annexbConvert(buffer)
			if err == nil {
				p.duplicater.MediaFrames <- data
			}
		} else {
			p.duplicater.MediaFrames <- buffer
		}

	}
}

// NewMediaStreamDuplicater duplicate this IncomingStreamTrack and callback the mediaframe
func NewMediaStreamDuplicater(track *IncomingStreamTrack) *MediaStreamDuplicater {

	duplicater := &MediaStreamDuplicater{}
	duplicater.track = track

	// We should make sure this source is the main source
	source := track.GetFirstEncoding().GetSource()
	duplicater.duplicater = native.NewMediaStreamDuplicaterFacade(source)

	track.OnStop(func() {
		duplicater.Stop()
	})

	listener := &overwrittenMediaFrameListener{
		duplicater: duplicater,
	}
	p := native.NewDirectorMediaFrameListener(listener)
	listener.p = p

	duplicater.listener = &goMediaFrameListener{MediaFrameListener: p}

	duplicater.duplicater.AddMediaListener(duplicater.listener)

	duplicater.MediaFrames = make(chan []byte, 5)
	return duplicater
}

// Stop stop this
func (d *MediaStreamDuplicater) Stop() {

	if d.track == nil {
		return
	}

	if d.listener != nil {
		d.duplicater.RemoveMediaListener(d.listener)
		d.listener.deleteMediaFrameListener()
	}

	d.track = nil
}
