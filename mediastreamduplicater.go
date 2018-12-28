package mediaserver

import "C"
import (
	"unsafe"

	native "github.com/notedit/media-server-go/wrapper"
)

type MediaFrameCallback func(frame []byte, duration uint, timestamp uint)

type MediaStreamDuplicater struct {
	track      *IncomingStreamTrack
	duplicater native.MediaStreamDuplicaterFacade
	callback   MediaFrameCallback
	listener   mediaframeListener
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
				p.duplicater.callback(data, frame.GetDuration(), frame.GetTimeStamp())
			}
		} else {
			p.duplicater.callback(buffer, frame.GetDuration(), frame.GetTimeStamp())
		}

	}
}

func NewMediaStreamDuplicater(track *IncomingStreamTrack, callback MediaFrameCallback) *MediaStreamDuplicater {

	duplicater := &MediaStreamDuplicater{}
	duplicater.track = track

	// We should make sure this source is the main source
	source := track.GetFirstEncoding().GetSource()
	duplicater.duplicater = native.NewMediaStreamDuplicaterFacade(source)

	track.On("stopped", func() {
		duplicater.Stop()
	})

	duplicater.callback = callback

	listener := &overwrittenMediaFrameListener{
		duplicater: duplicater,
	}
	p := native.NewDirectorMediaFrameListener(listener)
	listener.p = p

	duplicater.listener = &goMediaFrameListener{MediaFrameListener: p}

	duplicater.duplicater.AddMediaListener(duplicater.listener)

	return duplicater
}

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
