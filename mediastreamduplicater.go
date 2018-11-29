package mediaserver

type MediaFrameCallback func(frame MediaFrame)

type MediaStreamDuplicater struct {
	track      *IncomingStreamTrack
	duplicater MediaStreamDuplicaterFacade
	callback   MediaFrameCallback
	listener   mediaframeListener
}

type mediaframeListener interface {
	MediaFrameListener
	deleteMediaFrameListener()
}

type goMediaFrameListener struct {
	MediaFrameListener
}

func (m *goMediaFrameListener) deleteMediaFrameListener() {
	DeleteDirectorMediaFrameListener(m.MediaFrameListener)
}

type overwrittenMediaFrameListener struct {
	p          MediaFrameListener
	duplicater *MediaStreamDuplicater
}

func (p *overwrittenMediaFrameListener) OnMediaFrame(frame MediaFrame) {

	if p.duplicater != nil {
		p.duplicater.callback(frame)
	}
}

func NewMediaStreamDuplicater(track *IncomingStreamTrack, callback MediaFrameCallback) *MediaStreamDuplicater {

	duplicater := &MediaStreamDuplicater{}
	duplicater.track = track

	// We should make sure this source is the main source
	source := track.GetFirstEncoding().GetSource()
	duplicater.duplicater = NewMediaStreamDuplicaterFacade(source)

	track.On("stopped", func() {
		duplicater.Stop()
	})

	duplicater.callback = callback

	listener := &overwrittenMediaFrameListener{
		duplicater: duplicater,
	}
	p := NewDirectorMediaFrameListener(listener)
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
