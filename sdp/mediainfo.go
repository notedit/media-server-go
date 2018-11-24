package sdp

import (
	"strings"
)

type MediaInfo struct {
	id            string
	mtype         string // "audio" | "video"
	direction     Direction
	extensions    map[int]string        // Add rtp header extension support
	codecs        map[string]*CodecInfo // key: pt   value:  codec info
	rids          map[string]*RIDInfo
	simulcast     bool
	simulcastInfo *SimulcastInfo
	bitrate       int
}

// TODO add simulcast or rtx
func NewMediaInfo(id string, mtype string) *MediaInfo {

	media := &MediaInfo{
		id:         id,
		mtype:      mtype,
		direction:  SENDRECV,
		extensions: map[int]string{},
		codecs:     map[string]*CodecInfo{},
		rids:       map[string]*RIDInfo{},
		bitrate:    0,
	}

	return media
}

func (m *MediaInfo) Clone() *MediaInfo {

	cloned := NewMediaInfo(m.id, m.mtype)
	cloned.SetDirection(m.direction)
	cloned.SetBitrate(m.bitrate)

	for _, codec := range m.codecs {
		cloned.AddCodec(codec.Clone())
	}
	for id, name := range m.extensions {
		cloned.AddExtension(id, name)
	}
	for _, rid := range m.rids {
		cloned.AddRID(rid.Clone())
	}
	if m.simulcastInfo != nil {
		cloned.SetSimulcastInfo(m.simulcastInfo.Clone())
	}
	if m.simulcast {
		cloned.SetSimulcast(m.simulcast)
	}
	return cloned
}

func (m *MediaInfo) GetType() string {

	return m.mtype
}

func (m *MediaInfo) GetID() string {

	return m.id
}

func (m *MediaInfo) SetID(id string) {

	m.id = id
}

func (m *MediaInfo) AddExtension(id int, name string) {

	m.extensions[id] = name
}

func (m *MediaInfo) AddRID(ridInfo *RIDInfo) {

	m.rids[ridInfo.GetID()] = ridInfo
}

func (m *MediaInfo) AddCodec(codecInfo *CodecInfo) {

	m.codecs[codecInfo.GetCodec()] = codecInfo
}

func (m *MediaInfo) SetCodecs(codecs map[string]*CodecInfo) {

	m.codecs = codecs
}

func (m *MediaInfo) GetCodec(codec string) *CodecInfo {

	for _, codecInfo := range m.codecs {
		if strings.ToLower(codecInfo.GetCodec()) == strings.ToLower(codec) {
			return codecInfo
		}
	}
	return nil
}

func (m *MediaInfo) GetCodecForType(pt int) *CodecInfo {

	for _, codecInfo := range m.codecs {
		if codecInfo.GetType() == pt {
			return codecInfo
		}
	}
	return nil
}

func (m *MediaInfo) GetCodecs() map[string]*CodecInfo {

	return m.codecs
}

func (m *MediaInfo) HasRTX() bool {

	for _, codecInfo := range m.codecs {
		if codecInfo.HasRTX() {
			return true
		}
	}

	return false
}

func (m *MediaInfo) GetExtensions() map[int]string {

	return m.extensions
}

func (m *MediaInfo) HasExtension(uri string) bool {

	for _, extension := range m.extensions {
		if extension == uri {
			return true
		}
	}
	return false
}

func (m *MediaInfo) GetRIDS() map[string]*RIDInfo {

	return m.rids
}

func (m *MediaInfo) GetRID(id string) *RIDInfo {

	return m.rids[id]
}

func (m *MediaInfo) GetBitrate() int {

	return m.bitrate
}

func (m *MediaInfo) SetBitrate(bitrate int) {

	m.bitrate = bitrate
}

func (m *MediaInfo) GetDirection() Direction {
	return m.direction
}

func (m *MediaInfo) SetDirection(direction Direction) {

	m.direction = direction
}

func (m *MediaInfo) GetSimulcast() bool {

	return m.simulcast
}

func (m *MediaInfo) SetSimulcast(simulcast bool) {

	m.simulcast = simulcast
}

func (m *MediaInfo) GetSimulcastInfo() *SimulcastInfo {

	return m.simulcastInfo
}

func (m *MediaInfo) SetSimulcastInfo(info *SimulcastInfo) {

	m.simulcastInfo = info
}

func (m *MediaInfo) Answer(supportedMedia *MediaInfo) *MediaInfo {

	answer := NewMediaInfo(m.id, m.mtype)
	answer.SetDirection(m.direction.Reverse())

	for _, codec := range m.codecs {
		// If we support this codec
		if supportedMedia.GetCodec(strings.ToLower(codec.GetCodec())) != nil {
			supported := supportedMedia.GetCodec(strings.ToLower(codec.GetCodec()))
			if supported.GetCodec() == "h264" && supported.HasParam("packetization-mode") && supported.GetParam("packetization-mode") != codec.GetParam("packetization-mode") {
				continue
			}
			if supported.GetCodec() == "h264" && supported.HasParam("profile-level-id") && supported.GetParam("profile-level-id") != codec.GetParam("profile-level-id") {
				continue
			}
			cloned := supported.Clone()
			cloned.SetType(codec.GetType())
			if cloned.HasRTX() {
				cloned.SetRTX(codec.GetRTX())
			}
			cloned.AddParams(codec.GetParams())
			answer.AddCodec(cloned)
		}
	}

	//extentions
	for i, uri := range m.extensions {
		if supportedMedia.HasExtension(uri) {
			answer.AddExtension(i, uri)
		}
	}

	// todo add simulcast support
	if supportedMedia.simulcast && m.simulcast && m.simulcastInfo != nil {
		// simulcast := NewSimulcastInfo()

		// send := m.simulcastInfo.GetSimulcastStreams(SEND)
		// if send != nil {

		// }
	}

	return answer
}

func (m *MediaInfo) AnswerCapability(cap *Capability) *MediaInfo {

	answer := NewMediaInfo(m.id, m.mtype)
	answer.SetDirection(m.direction.Reverse())

	rtcpfbs := []*RTCPFeedbackInfo{}
	for _, rtcpfb := range cap.Rtcpfbs {
		rtcpfbs = append(rtcpfbs, NewRTCPFeedbackInfo(rtcpfb.ID, rtcpfb.Params))
	}
	codecs := CodecMapFromNames(cap.Codecs, cap.Rtx, rtcpfbs)

	for codecName, codec := range m.codecs {
		// If we support this codec
		if codecs[codecName] != nil {
			supported := codecs[codecName]
			if supported.GetCodec() == "h264" && supported.HasParam("packetization-mode") && supported.GetParam("packetization-mode") != codec.GetParam("packetization-mode") {
				continue
			}
			if supported.GetCodec() == "h264" && supported.HasParam("profile-level-id") && supported.GetParam("profile-level-id") != codec.GetParam("profile-level-id") {
				continue
			}
			cloned := supported.Clone()
			cloned.SetType(codec.GetType())
			if cloned.HasRTX() {
				cloned.SetRTX(codec.GetRTX())
			}
			cloned.AddParams(codec.GetParams())
			answer.AddCodec(cloned)
		}
	}

	//extentions
	for i, uri := range m.extensions {
		if contains(cap.Extensions, uri) {
			answer.AddExtension(i, uri)
		}
	}

	if cap.Simulcast && m.simulcast && m.simulcastInfo != nil {
		// simulcast := NewSimulcastInfo()

		// send := m.simulcastInfo.GetSimulcastStreams(SEND)
		// if send != nil {

		// }
	}

	return answer
}

func MediaInfoCreate(mType string, capability *Capability) *MediaInfo {

	mediaInfo := NewMediaInfo(mType, mType)

	if capability != nil {
		if capability.Codecs != nil {
			rtcpfbs := []*RTCPFeedbackInfo{}
			for _, rtcpfb := range capability.Rtcpfbs {
				rtcpfbs = append(rtcpfbs, NewRTCPFeedbackInfo(rtcpfb.ID, rtcpfb.Params))
			}
			codecs := CodecMapFromNames(capability.Codecs, capability.Rtx, rtcpfbs)
			mediaInfo.SetCodecs(codecs)
		}
		for i, extension := range capability.Extensions {
			mediaInfo.AddExtension(i, extension)
		}
	} else {
		mediaInfo.SetDirection(INACTIVE)
	}

	return mediaInfo
}
