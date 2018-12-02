package mediaserver

import "C"

import (
	"unsafe"
)

type RtpCallback func(data []byte)

type MediatrackDuplicater struct {
	track      *IncomingStreamTrack
	duplicater RTPPacketDuplicaterFacade
	callback   RtpCallback
	listener   rtpPackerListener
}

type rtpPackerListener interface {
	RTPPacketListener
	deleteRTPPacketListener()
}

type gortpPacketListener struct {
	RTPPacketListener
}

func (r *gortpPacketListener) deleteRTPPacketListener() {
	DeleteDirectorRTPPacketListener(r.RTPPacketListener)
}

type overwrittenRtpPacketListener struct {
	p          RTPPacketListener
	duplicater *MediatrackDuplicater
}

func (o *overwrittenRtpPacketListener) OnRTP(data *byte, length uint) {

	if o.duplicater != nil {
		rtpdata := C.GoBytes(unsafe.Pointer(data), C.int(length))
		o.duplicater.callback(rtpdata)
	}
}

func NewMediatrackDuplicater(track *IncomingStreamTrack, callback RtpCallback) *MediatrackDuplicater {

	duplicater := &MediatrackDuplicater{}
	duplicater.track = track
	source := track.GetFirstEncoding().GetSource()

	duplicater.duplicater = NewRTPPacketDuplicaterFacade(source)

	track.On("stopped", func() {
		duplicater.Stop()
	})

	duplicater.callback = callback

	listener := &overwrittenRtpPacketListener{
		duplicater: duplicater,
	}
	p := NewDirectorRTPPacketListener(listener)
	listener.p = p

	duplicater.listener = &gortpPacketListener{RTPPacketListener: p}

	duplicater.duplicater.SetRTPListener(duplicater.listener)

	return duplicater
}

func (m *MediatrackDuplicater) Stop() {
	if m.track == nil {
		return
	}

	if m.listener != nil {
		m.listener.deleteRTPPacketListener()
	}

	m.track = nil
}
