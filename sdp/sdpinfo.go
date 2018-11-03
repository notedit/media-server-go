package sdp

import (
	"strconv"
	"strings"
	"time"

	sdptransform "github.com/notedit/go-sdp-transform"
)

type SDPInfo struct {
	version    int
	streams    map[string]*StreamInfo
	medias     []*MediaInfo     // as we need to keep order
	candidates []*CandidateInfo // keep order
	ice        *ICEInfo
	dtls       *DTLSInfo
}

func NewSDPInfo() *SDPInfo {

	sdp := &SDPInfo{
		version:    1,
		streams:    map[string]*StreamInfo{},
		medias:     []*MediaInfo{},
		candidates: []*CandidateInfo{},
	}

	return sdp
}

func (s *SDPInfo) SetVersion(version int) {

	s.version = version
}

func (s *SDPInfo) AddMedia(media *MediaInfo) {

	s.medias = append(s.medias, media)
}

func (s *SDPInfo) GetMedia(mtype string) *MediaInfo {

	for _, media := range s.medias {
		if strings.ToLower(media.GetType()) == strings.ToLower(mtype) {
			return media
		}
	}
	return nil
}

func (s *SDPInfo) GetMediasByType(mtype string) []*MediaInfo {

	medias := []*MediaInfo{}
	for _, media := range s.medias {
		if strings.ToLower(media.GetType()) == strings.ToLower(mtype) {
			medias = append(medias, media)
		}
	}
	return medias
}

func (s *SDPInfo) GetMediaByID(mid string) *MediaInfo {

	for _, media := range s.medias {
		if strings.ToLower(media.GetID()) == strings.ToLower(mid) {
			return media
		}
	}

	return nil
}

func (s *SDPInfo) GetMedias() []*MediaInfo {

	return s.medias
}

func (s *SDPInfo) GetVersion() int {

	return s.version
}

func (s *SDPInfo) GetDTLS() *DTLSInfo {

	return s.dtls
}

func (s *SDPInfo) SetDTLS(dtls *DTLSInfo) {

	s.dtls = dtls
}

func (s *SDPInfo) GetICE() *ICEInfo {

	return s.ice
}

func (s *SDPInfo) SetICE(ice *ICEInfo) {

	s.ice = ice
}

func (s *SDPInfo) AddCandidate(candidate *CandidateInfo) {

	// todo check
	s.candidates = append(s.candidates, candidate)
}

func (s *SDPInfo) AddCandidates(candidates []*CandidateInfo) {

	for _, candidate := range candidates {
		s.AddCandidate(candidate)
	}
}

func (s *SDPInfo) GetCandidates() []*CandidateInfo {

	return s.candidates
}

func (s *SDPInfo) GetStream(id string) *StreamInfo {

	return s.streams[id]
}

func (s *SDPInfo) GetStreams() map[string]*StreamInfo {

	return s.streams
}

func (s *SDPInfo) GetFirstStream() *StreamInfo {

	for _, stream := range s.streams {
		return stream
	}
	return nil
}

func (s *SDPInfo) AddStream(stream *StreamInfo) {

	s.streams[stream.GetID()] = stream
}

func (s *SDPInfo) RemoveStream(stream *StreamInfo) {

	delete(s.streams, stream.GetID())
}

func (s *SDPInfo) Answer(ice *ICEInfo, dtls *DTLSInfo, candidates []*CandidateInfo, capabilities map[string]*Capability) *SDPInfo {

	sdpInfo := NewSDPInfo()

	if ice != nil {
		sdpInfo.SetICE(ice.Clone())
	}

	if dtls != nil {
		sdpInfo.SetDTLS(dtls)
	}

	for _, candidate := range candidates {
		sdpInfo.AddCandidate(candidate)
	}

	for _, media := range s.medias {

		capability := capabilities[media.GetType()]
		if capability != nil {
			media.Answer(capability)
			sdpInfo.AddMedia(media)
		}
	}

	return sdpInfo
}

func (s *SDPInfo) String() string {

	sdpMap := &sdptransform.SdpStruct{
		Version: 0,
		Media:   []*sdptransform.MediaStruct{},
		Groups:  []*sdptransform.GroupStruct{},
	}

	sdpMap.Origin = &sdptransform.OriginStruct{
		Username:       "-",
		SessionId:      strconv.FormatInt(time.Now().UnixNano(), 10),
		SessionVersion: s.version,
		NetType:        "IN",
		Address:        "127.0.0.1",
	}

	sdpMap.Connection = &sdptransform.ConnectionStruct{
		Version: 4,
		Ip:      "0.0.0.0",
	}

	sdpMap.Name = "media"

	if s.GetICE().IsLite() {
		sdpMap.Icelite = "ice-lite"
	}

	sdpMap.MsidSemantic = &sdptransform.MsidSemanticStruct{
		Semantic: "WMS",
		Token:    "*",
	}

	bundleType := "BUNDLE"
	bundleMids := []string{}

	for _, media := range s.medias {

		mediaMap := &sdptransform.MediaStruct{
			Type:       media.GetType(),
			Port:       9,
			Protocal:   "UDP/TLS/RTP/SAVP",
			Fmtp:       []*sdptransform.FmtpStruct{},
			Rtp:        []*sdptransform.RtpStruct{},
			RtcpFb:     []*sdptransform.RtcpFbStruct{},
			Ext:        []*sdptransform.ExtStruct{},
			Bandwidth:  []*sdptransform.BandwithStruct{},
			Candidates: []*sdptransform.CandidateStruct{},
			SsrcGroups: []*sdptransform.SsrcGroupStruct{},
			Ssrcs:      []*sdptransform.SsrcStruct{},
			Rids:       []*sdptransform.RidStruct{},
		}

		mediaMap.Direction = media.GetDirection().String()

		mediaMap.RtcpMux = "rtcp-mux"

		mediaMap.RtcpRsize = "rtcp-rsize"

		mediaMap.Mid = media.GetID()

		bundleMids = append(bundleMids, media.GetID())

		if media.GetBitrate() > 0 {
			mediaMap.Bandwidth = append(mediaMap.Bandwidth, &sdptransform.BandwithStruct{
				Type:  "AS",
				Limit: media.GetBitrate(),
			})
		}

		for _, candidate := range s.GetCandidates() {

			mediaMap.Candidates = append(mediaMap.Candidates, &sdptransform.CandidateStruct{
				Foundation: candidate.GetFoundation(),
				Component:  candidate.GetComponentID(),
				Transport:  candidate.GetTransport(),
				Priority:   candidate.GetPriority(),
				Ip:         candidate.GetAddress(),
				Port:       candidate.GetPort(),
				Type:       candidate.GetType(),
				Raddr:      candidate.GetRelAddr(),
				Rport:      candidate.GetRelPort(),
			})
		}

		mediaMap.IceUfrag = s.GetICE().GetUfrag()
		mediaMap.IcePwd = s.GetICE().GetPassword()

		mediaMap.Fingerprint = &sdptransform.FingerprintStruct{
			Type: s.GetDTLS().GetHash(),
			Hash: s.GetDTLS().GetFingerprint(),
		}

		mediaMap.Setup = s.GetDTLS().GetSetup().String()

		for _, codec := range media.GetCodecs() {

			if "video" == strings.ToLower(media.GetType()) {

				mediaMap.Rtp = append(mediaMap.Rtp, &sdptransform.RtpStruct{
					Payload: codec.GetType(),
					Codec:   strings.ToUpper(codec.GetCodec()),
					Rate:    90000,
				})
			} else {

				if "opus" == strings.ToLower(codec.GetCodec()) {

					mediaMap.Rtp = append(mediaMap.Rtp, &sdptransform.RtpStruct{
						Payload:  codec.GetType(),
						Codec:    codec.GetCodec(),
						Rate:     48000,
						Encoding: 2,
					})
				} else {

					mediaMap.Rtp = append(mediaMap.Rtp, &sdptransform.RtpStruct{
						Payload: codec.GetType(),
						Codec:   codec.GetCodec(),
						Rate:    8000,
					})
				}
			}

			for _, rtcpfb := range codec.GetRTCPFeedbacks() {
				mediaMap.RtcpFb = append(mediaMap.RtcpFb, &sdptransform.RtcpFbStruct{
					Payload: codec.GetType(),
					Type:    rtcpfb.GetID(),
					Subtype: strings.Join(rtcpfb.GetParams(), " "),
				})
			}

			if codec.HasRTX() {
				mediaMap.Rtp = append(mediaMap.Rtp, &sdptransform.RtpStruct{
					Payload: codec.GetRTX(),
					Codec:   "rtx",
					Rate:    90000,
				})
				mediaMap.Fmtp = append(mediaMap.Fmtp, &sdptransform.FmtpStruct{
					Payload: codec.GetRTX(),
					Config:  "apt=" + strconv.Itoa(codec.GetType()),
				})
			}

			params := codec.GetParams()

			if params != nil && len(params) > 0 {

				fmtp := &sdptransform.FmtpStruct{
					Payload: codec.GetType(),
					Config:  "",
				}

				for k, v := range params {

					if fmtp.Config != "" {
						fmtp.Config = fmtp.Config + ";"
					}

					// k and value
					if v != "" {
						fmtp.Config = fmtp.Config + k + "=" + v
					} else {
						fmtp.Config = fmtp.Config + k
					}
				}

				mediaMap.Fmtp = append(mediaMap.Fmtp, fmtp)
			}
		}

		payloads := []int{}

		for _, rtp := range mediaMap.Rtp {
			payloads = append(payloads, rtp.Payload)
		}

		mediaMap.Payloads = arrayToString(payloads, " ")

		for id, uri := range media.GetExtensions() {

			mediaMap.Ext = append(mediaMap.Ext, &sdptransform.ExtStruct{
				Value: id,
				Uri:   uri,
			})
		}

		for _, ridInfo := range media.GetRIDS() {

			rid := &sdptransform.RidStruct{
				Id:        ridInfo.GetID(),
				Direction: ridInfo.GetDirection().String(),
				Params:    "",
			}

			if len(ridInfo.GetFormats()) > 0 {
				rid.Params = "pt=" + strings.Join(ridInfo.GetFormats(), ",")
			}

			for key, val := range ridInfo.GetParams() {
				if rid.Params == "" {
					rid.Params = key + "=" + val
				} else {
					rid.Params = rid.Params + ";" + key + "=" + val
				}
			}

			mediaMap.Rids = append(mediaMap.Rids, rid)
		}

		// Todo simulcast
		sdpMap.Media = append(sdpMap.Media, mediaMap)
	}

	// streams
	for _, stream := range s.GetStreams() {
		for _, track := range stream.GetTracks() {
			for _, md := range sdpMap.Media {
				// check if it is unified or plan b
				if track.GetMediaID() != "" {
					if track.GetMediaID() == md.Mid {
						groups := track.GetSourceGroupS()
						for _, group := range groups {
							md.SsrcGroups = append(md.SsrcGroups, &sdptransform.SsrcGroupStruct{
								Semantics: group.GetSemantics(),
								Ssrcs:     arrayToString(group.GetSSRCs(), " "),
							})
						}
						ssrcs := track.GetSSRCS()
						for _, ssrc := range ssrcs {
							md.Ssrcs = append(md.Ssrcs, &sdptransform.SsrcStruct{
								Id:        ssrc,
								Attribute: "cname",
								Value:     stream.GetID(),
							})
						}
						md.Msid = stream.GetID() + " " + track.GetID()
						break
					}
				} else if strings.ToLower(md.Type) == strings.ToLower(track.GetMedia()) {

					groups := track.GetSourceGroupS()

					for _, group := range groups {
						md.SsrcGroups = append(md.SsrcGroups, &sdptransform.SsrcGroupStruct{
							Semantics: group.GetSemantics(),
							Ssrcs:     arrayToString(group.GetSSRCs(), " "),
						})
					}
					ssrcs := track.GetSSRCS()
					for _, ssrc := range ssrcs {
						md.Ssrcs = append(md.Ssrcs, &sdptransform.SsrcStruct{
							Id:        ssrc,
							Attribute: "cname",
							Value:     stream.GetID(),
						})
						md.Ssrcs = append(md.Ssrcs, &sdptransform.SsrcStruct{
							Id:        ssrc,
							Attribute: "msid",
							Value:     stream.GetID() + " " + track.GetID(),
						})
					}
					break
				}
			}
		}
	}
	sdpMap.Groups = append(sdpMap.Groups, &sdptransform.GroupStruct{
		Mids: strings.Join(bundleMids, " "),
		Type: bundleType,
	})

	sdpStr, err := sdptransform.Write(sdpMap)
	if err != nil {
		println(err)
	}

	return sdpStr
}

func (s *SDPInfo) Clone() *SDPInfo {

	cloned := NewSDPInfo()
	cloned.SetVersion(s.GetVersion())
	for _, media := range s.GetMedias() {
		cloned.AddMedia(media.Clone())
	}
	for _, stream := range s.GetStreams() {
		cloned.AddStream(stream.Clone())
	}
	for _, candidate := range s.GetCandidates() {
		cloned.AddCandidate(candidate)
	}
	cloned.SetICE(s.GetICE().Clone())
	cloned.SetDTLS(s.GetDTLS().Clone())
	return cloned
}

// Unify return an unified plan version of the SDP info
func (s *SDPInfo) Unify() *SDPInfo {
	cloned := NewSDPInfo()

	cloned.version = s.version

	for _, media := range s.medias {
		cloned.AddMedia(media.Clone())
	}

	medias := map[string][]*MediaInfo{
		"audio": cloned.GetMediasByType("audio"),
		"video": cloned.GetMediasByType("video"),
	}

	for _, stream := range s.streams {
		clonedStream := stream.Clone()
		for _, clonedTrack := range clonedStream.GetTracks() {
			var clonedMedia *MediaInfo
			if len(medias[clonedTrack.GetMedia()]) == 0 {
				media := s.GetMedia(clonedTrack.GetMedia())
				clonedMedia = media.Clone()
				clonedMedia.SetID(clonedTrack.GetID())
				cloned.AddMedia(clonedMedia)
			} else {
				mediaList := medias[clonedTrack.GetMedia()]
				clonedMedia = mediaList[len(mediaList)-1]
				mediaList = mediaList[:len(mediaList)-1]
				medias[clonedTrack.GetMedia()] = mediaList
			}
			clonedTrack.SetMediaID(clonedMedia.GetID())
		}
		cloned.AddStream(clonedStream)
	}

	for _, candidate := range s.GetCandidates() {
		cloned.AddCandidate(candidate.Clone())
	}

	cloned.SetICE(s.GetICE().Clone())
	cloned.SetDTLS(s.GetDTLS().Clone())

	return cloned
}

type Capability struct {
	Codecs     []string
	Rtx        bool
	Rtcpfbs    []*RTCPFeedbackInfo
	Extensions []string
}

func Create(ice *ICEInfo, dtls *DTLSInfo, candidates []*CandidateInfo, capabilities map[string]*Capability) *SDPInfo {

	sdpInfo := NewSDPInfo()

	if ice != nil {
		sdpInfo.SetICE(ice.Clone())
	}

	if dtls != nil {
		sdpInfo.SetDTLS(dtls)
	}

	for _, candidate := range candidates {
		sdpInfo.AddCandidate(candidate)
	}

	dyn := 96

	for mType, capability := range capabilities {

		media := MediaInfoCreate(mType, capability)

		for _, codec := range media.GetCodecs() {

			if codec.GetType() >= 96 {
				dyn++
				codec.SetType(dyn)
			}

			if codec.GetRTX() > 0 {
				dyn++
				codec.SetRTX(dyn)
			}
		}

		sdpInfo.AddMedia(media)
	}

	return sdpInfo
}

func Parse(sdp string) (*SDPInfo, error) {

	sdpMap, err := sdptransform.Parse(sdp)

	if err != nil {
		return nil, err
	}

	sdpInfo := NewSDPInfo()

	sdpInfo.SetVersion(sdpMap.Version)

	for _, md := range sdpMap.Media {

		media := md.Type
		mid := md.Mid

		mediaInfo := NewMediaInfo(mid, media)

		ufrag := md.IceUfrag
		pwd := md.IcePwd

		sdpInfo.SetICE(NewICEInfo(ufrag, pwd))

		for _, candiate := range md.Candidates {

			candidateInfo := NewCandidateInfo(
				candiate.Foundation,
				candiate.Component,
				candiate.Transport,
				candiate.Priority,
				candiate.Ip,
				candiate.Port,
				candiate.Type,
				candiate.Raddr,
				candiate.Rport)

			sdpInfo.AddCandidate(candidateInfo)
		}

		var fingerpirnt *sdptransform.FingerprintStruct

		if sdpMap.Fingerprint != nil {
			fingerpirnt = sdpMap.Fingerprint
		}

		if md.Fingerprint != nil {
			fingerpirnt = md.Fingerprint
		}

		remoteHash := fingerpirnt.Type
		remoteFingerprint := fingerpirnt.Hash

		setup := SETUPACTPASS

		if md.Setup != "" {
			setup = SetupByValue(md.Setup)
		}

		sdpInfo.SetDTLS(NewDTLSInfo(setup, remoteHash, remoteFingerprint))

		direction := DirectionSENDRECV

		if md.Direction != "" {
			direction = DirectionbyValue(md.Direction)
		}

		mediaInfo.SetDirection(direction)

		apts := map[int]int{}

		for _, fmt := range md.Rtp {

			payload := fmt.Payload
			codec := fmt.Codec

			if "RED" == strings.ToUpper(codec) || "ULPFEC" == strings.ToUpper(codec) {
				continue
			}

			params := map[string]string{}

			for _, fmtp := range md.Fmtp {

				if fmtp.Payload == payload {
					list := strings.Split(fmtp.Config, ";")

					for _, kv := range list {
						param := strings.Split(kv, "=")
						if len(param) < 2 {
							println("param too shot", param)
							continue
						}
						params[param[0]] = param[1]
					}
				}
			}

			if "RTX" == strings.ToUpper(codec) {
				if apt, ok := params["apt"]; ok {
					aptint, _ := strconv.Atoi(apt)
					apts[aptint] = payload
				}
			} else {
				mediaInfo.AddCodec(NewCodecInfo(codec, payload, params))
			}
		}

		// rtx
		for pt1, pt2 := range apts {
			codecInfo := mediaInfo.GetCodecForType(pt1)
			if codecInfo != nil {
				codecInfo.SetRTX(pt2)
			}
		}

		// rtcpFb
		if md.RtcpFb != nil {
			for _, rtcfb := range md.RtcpFb {
				codecInfo := mediaInfo.GetCodecForType(rtcfb.Payload)
				if codecInfo != nil {
					id := rtcfb.Type
					params := []string{}
					if rtcfb.Subtype != "" {
						params = strings.Split(rtcfb.Subtype, " ")
					}
					codecInfo.AddRTCPFeedback(NewRTCPFeedbackInfo(id, params))
				}
			}
		}

		// extmap
		for _, extmap := range md.Ext {
			mediaInfo.AddExtension(extmap.Value, extmap.Uri)
		}

		for _, rid := range md.Rids {
			direction := DirectionWaybyValue(rid.Direction)
			ridInfo := NewRIDInfo(rid.Id, direction)

			formats := []string{}
			params := map[string]string{}

			if rid.Params != "" {
				list := sdptransform.ParseParams(rid.Params)
				for k, v := range list {
					if k == "pt" {
						formats = strings.Split(v, ",")
					} else {
						params[k] = v
					}
				}
				ridInfo.SetFormats(formats)
				ridInfo.SetParams(params)
			}

			mediaInfo.AddRID(ridInfo)
		}

		// todo simulcast

		encodings := [][]*TrackEncodingInfo{}

		sources := map[int]*SourceInfo{}

		if md.Ssrcs != nil {

			for _, ssrcAttr := range md.Ssrcs {

				ssrc := ssrcAttr.Id
				key := ssrcAttr.Attribute
				value := ssrcAttr.Value

				source, ok := sources[ssrc]

				if !ok {
					source = NewSourceInfo(ssrc)
					sources[ssrc] = source
				}

				if strings.ToLower(key) == "cname" {
					source.SetCName(value)
				} else if strings.ToLower(key) == "msid" {
					ids := strings.Split(value, " ")
					// get stream id and track id
					streamId := ids[0]
					trackId := ids[1]

					source.SetStreamID(streamId)
					source.SetTrackID(trackId)

					stream := sdpInfo.GetStream(streamId)

					if stream == nil {
						stream = NewStreamInfo(streamId)
						sdpInfo.AddStream(stream)
					}

					track := stream.GetTrack(trackId)

					if track == nil {
						track = NewTrackInfo(trackId, media)
						// Set simulcast encodings (if any)
						// todo
						track.SetEncodings(encodings)
						stream.AddTrack(track)
					}
					// Add ssrc
					track.AddSSRC(ssrc)
				}

			}
		}

		// Check if ther is a global msid
		// Why this?
		if md.Msid != "" {
			ids := strings.Split(md.Msid, " ")
			streamId := ids[0]
			trackId := ids[1]

			stream := sdpInfo.GetStream(streamId)

			if stream == nil {
				stream = NewStreamInfo(streamId)
				sdpInfo.AddStream(stream)
			}

			track := stream.GetTrack(trackId)

			if track == nil {
				track = NewTrackInfo(trackId, media)
				track.SetMediaID(mid)
				track.SetEncodings(encodings)
				stream.AddTrack(track)
			}

			for ssrc, source := range sources {

				if source.GetStreamID() == "" {
					source.SetStreamID(streamId)
					source.SetTrackID(trackId)
					track.AddSSRC(ssrc)
				}
			}
		}

		for ssrc, source := range sources {

			if source.GetStreamID() == "" {

				streamId := source.GetCName()
				trackId := mid

				source.SetStreamID(streamId)
				source.SetTrackID(trackId)

				stream := sdpInfo.GetStream(streamId)

				if stream == nil {
					stream = NewStreamInfo(streamId)
					sdpInfo.AddStream(stream)
				}

				track := stream.GetTrack(trackId)

				if track == nil {
					track = NewTrackInfo(trackId, media)
					track.SetEncodings(encodings)
					stream.AddTrack(track)
				}

				track.AddSSRC(ssrc)
			}
		}

		if md.SsrcGroups != nil {
			for _, ssrcGroupAttr := range md.SsrcGroups {
				ssrcs := strings.Split(ssrcGroupAttr.Ssrcs, " ")
				ssrcsint := []int{}
				for _, ssrcstr := range ssrcs {
					ssrcint, _ := strconv.Atoi(ssrcstr)
					ssrcsint = append(ssrcsint, ssrcint)
				}
				group := NewSourceGroupInfo(ssrcGroupAttr.Semantics, ssrcsint)
				ssrc := ssrcsint[0]
				source := sources[ssrc]

				streamInfo := sdpInfo.GetStream(source.GetStreamID())
				if streamInfo != nil && streamInfo.GetTrack(source.GetTrackID()) != nil {
					streamInfo.GetTrack(source.GetTrackID()).AddSourceGroup(group)
				}
			}
		}

		sdpInfo.AddMedia(mediaInfo)

	}

	return sdpInfo, nil
}
