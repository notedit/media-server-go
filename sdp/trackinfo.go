package sdp

import (
	"strings"
)

type TrackInfo struct {
	id        string
	mediaID   string
	media     string // "audio" | "video"
	ssrcs     []uint
	groups    []*SourceGroupInfo
	encodings [][]*TrackEncodingInfo
}

func NewTrackInfo(id, media string) *TrackInfo {

	info := &TrackInfo{
		id:        id,
		media:     media,
		ssrcs:     []uint{},
		groups:    []*SourceGroupInfo{},
		encodings: [][]*TrackEncodingInfo{},
	}

	return info
}

func (t *TrackInfo) Clone() *TrackInfo {

	cloned := &TrackInfo{
		id:        t.id,
		media:     t.media,
		ssrcs:     make([]uint, len(t.ssrcs)),
		groups:    make([]*SourceGroupInfo, 0),
		encodings: make([][]*TrackEncodingInfo, len(t.encodings)),
	}
	copy(cloned.ssrcs, t.ssrcs)
	for _, group := range t.groups {
		cloned.groups = append(cloned.groups, group.Clone())
	}
	for i := range t.encodings {
		for _, v := range t.encodings[i] {
			cloned.encodings[i] = append(cloned.encodings[i], v.Clone())
		}
	}
	return cloned
}

func (t *TrackInfo) GetMedia() string {

	return t.media
}

func (t *TrackInfo) SetMediaID(mediaID string) {

	t.mediaID = mediaID
}

func (t *TrackInfo) GetMediaID() string {

	return t.mediaID
}

func (t *TrackInfo) GetID() string {

	return t.id
}

func (t *TrackInfo) AddSSRC(ssrc uint) {

	t.ssrcs = append(t.ssrcs, ssrc)
}

func (t *TrackInfo) GetSSRCS() []uint {

	return t.ssrcs
}

func (t *TrackInfo) AddSourceGroup(group *SourceGroupInfo) {

	t.groups = append(t.groups, group)
}

func (t *TrackInfo) GetSourceGroup(schematics string) *SourceGroupInfo {

	for _, group := range t.groups {
		if strings.ToLower(group.GetSemantics()) == strings.ToLower(schematics) {
			return group
		}
	}
	return nil
}

func (t *TrackInfo) GetSourceGroupS() []*SourceGroupInfo {

	return t.groups
}

func (t *TrackInfo) GetEncodings() [][]*TrackEncodingInfo {

	return t.encodings
}

func (t *TrackInfo) AddEncoding(encoding *TrackEncodingInfo) {

	t.encodings = append(t.encodings, []*TrackEncodingInfo{encoding})
}

func (t *TrackInfo) AddAlternativeEncodings(alternatives []*TrackEncodingInfo) {

	t.encodings = append(t.encodings, alternatives)
}

func (t *TrackInfo) SetEncodings(encodings [][]*TrackEncodingInfo) {

	t.encodings = encodings
}
