package sdp

import (
	"fmt"
	"strings"
)

type MediaInfo struct {
	id         string
	mtype      string // "audio" | "video"
	direction  Direction
	extensions map[int]string        // Add rtp header extension support
	codecs     map[string]*CodecInfo // key: pt   value:  codec info
	rids       map[string]*RIDInfo
	simulcast  *SimulcastInfo
	bitrate    int
}

func NewMediaInfo(id string, mtype string) *MediaInfo {

	media := &MediaInfo{
		id:         id,
		mtype:      mtype,
		direction:  DirectionSENDRECV,
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
	if m.simulcast != nil {
		cloned.SetSimulcast(m.simulcast.Clone())
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
		fmt.Println("GetCodec ", codecInfo.GetCodec(), codecInfo.GetType())
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

func (m *MediaInfo) GetSimulcast() *SimulcastInfo {

	return m.simulcast
}

func (m *MediaInfo) SetSimulcast(simulcast *SimulcastInfo) {

	m.simulcast = simulcast
}

func (m *MediaInfo) Answer(capability *Capability) *MediaInfo {

	return nil
}

func MediaInfoCreate(mType string, capability *Capability) *MediaInfo {

	mediaInfo := NewMediaInfo(mType, mType)

	if capability != nil {
		if capability.Codecs != nil {
			codecs := CodecMapFromNames(capability.Codecs, capability.Rtx, capability.Rtcpfbs)
			mediaInfo.SetCodecs(codecs)
		}
		for i, extension := range capability.Extensions {
			mediaInfo.AddExtension(i, extension)
		}
	} else {
		mediaInfo.SetDirection(DirectionINACTIVE)
	}

	return mediaInfo
}
