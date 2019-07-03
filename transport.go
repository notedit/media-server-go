package mediaserver

import (
	"fmt"
	"sync"

	"github.com/gofrs/uuid"
	native "github.com/notedit/media-server-go/wrapper"
	"github.com/notedit/sdp"
)

type senderSideEstimatorListener interface {
	native.SenderSideEstimatorListener
	deleteSenderSideEstimatorListener()
}

type goSenderSideEstimatorListener struct {
	native.SenderSideEstimatorListener
}

func (r *goSenderSideEstimatorListener) deleteSenderSideEstimatorListener() {
	native.DeleteDirectorSenderSideEstimatorListener(r.SenderSideEstimatorListener)
}

type overwrittenSenderSideEstimatorListener struct {
	p native.SenderSideEstimatorListener
}

func (p *overwrittenSenderSideEstimatorListener) OnTargetBitrateRequested(bitrate uint) {
	fmt.Println(bitrate)
}

type dtlsICETransportListener interface {
	native.DTLSICETransportListener
	deleteDTLSICETransportListener()
}

type goDTLSICETransportListener struct {
	native.DTLSICETransportListener
}

func (d *goDTLSICETransportListener) deleteDTLSICETransportListener() {
	native.DeleteDTLSICETransportListener(d.DTLSICETransportListener)
}

type overwrittenDTLSICETransportListener struct {
	p native.DTLSICETransportListener
}

func (p *overwrittenDTLSICETransportListener) OnDTLSStateChange(state uint) {
	fmt.Println("OnDTLSStateChange", state)
}

type (
	// TransportStopListener listener
	TransportStopListener func()
	// IncomingTrackListener new track listener
	IncomingTrackListener func(*IncomingStreamTrack, *IncomingStream)
	// OutgoingTrackListener new track listener
	OutgoingTrackListener func(*OutgoingStreamTrack, *OutgoingStream)
	// DTLSStateListener listener
	DTLSStateListener func(state string)
)

// ICEStats ice stats for this connection
type ICEStats struct {
	RequestsSent      int64
	RequestsReceived  int64
	ResponsesSent     int64
	ResponsesReceived int64
}

// Transport represent a connection between a local ICE candidate and a remote set of ICE candidates over a single DTLS session
type Transport struct {
	localIce         *sdp.ICEInfo
	localDtls        *sdp.DTLSInfo
	localCandidates  []*sdp.CandidateInfo
	remoteIce        *sdp.ICEInfo
	remoteDtls       *sdp.DTLSInfo
	remoteCandidates []*sdp.CandidateInfo
	bundle           native.RTPBundleTransport
	transport        native.DTLSICETransport
	connection       native.RTPBundleTransportConnection
	dtlsState        string

	username             native.StringFacade
	incomingStreams      map[string]*IncomingStream
	outgoingStreams      map[string]*OutgoingStream
	incomingStreamTracks map[string]*IncomingStreamTrack
	outgoingStreamTracks map[string]*OutgoingStreamTrack

	iceStats ICEStats

	senderSideListener       senderSideEstimatorListener
	dtlsICEListener          dtlsICETransportListener
	outDTLSStateListener     DTLSStateListener
	onTransportStopListeners []TransportStopListener
	onIncomingTrackListeners []IncomingTrackListener
	onOutgoingTrackListeners []OutgoingTrackListener
	sync.Mutex
}

// NewTransport create a new transport
func NewTransport(bundle native.RTPBundleTransport, remoteIce *sdp.ICEInfo, remoteDtls *sdp.DTLSInfo, remoteCandidates []*sdp.CandidateInfo,
	localIce *sdp.ICEInfo, localDtls *sdp.DTLSInfo, localCandidates []*sdp.CandidateInfo, disableSTUNKeepAlive bool) *Transport {

	transport := new(Transport)
	transport.remoteIce = remoteIce
	transport.remoteDtls = remoteDtls
	transport.localIce = localIce
	transport.localDtls = localDtls
	transport.bundle = bundle
	transport.dtlsState = "new"

	properties := native.NewProperties()

	properties.SetProperty("ice.localUsername", localIce.GetUfrag())
	properties.SetProperty("ice.localPassword", localIce.GetPassword())
	properties.SetProperty("ice.remoteUsername", remoteIce.GetUfrag())
	properties.SetProperty("ice.remotePassword", remoteIce.GetPassword())

	properties.SetProperty("dtls.setup", remoteDtls.GetSetup().String())
	properties.SetProperty("dtls.hash", remoteDtls.GetHash())
	properties.SetProperty("dtls.fingerprint", remoteDtls.GetFingerprint())

	stunKeepAlive := "false"
	if disableSTUNKeepAlive {
		stunKeepAlive = "true"
	}

	properties.SetProperty("disableSTUNKeepAlive", stunKeepAlive)

	transport.username = native.NewStringFacade(localIce.GetUfrag() + ":" + remoteIce.GetUfrag())
	transport.connection = bundle.AddICETransport(transport.username, properties)
	transport.transport = transport.connection.GetTransport()

	transport.iceStats = ICEStats{}

	native.DeleteProperties(properties)

	sseListener := &overwrittenSenderSideEstimatorListener{}
	p := native.NewDirectorSenderSideEstimatorListener(sseListener)
	sseListener.p = p

	transport.senderSideListener = &goSenderSideEstimatorListener{SenderSideEstimatorListener: p}
	transport.transport.SetSenderSideEstimatorListener(transport.senderSideListener)

	dtlsListener := &overwrittenDTLSICETransportListener{}
	dtlsl := native.NewDirectorDTLSICETransportListener(dtlsListener)
	dtlsListener.p = dtlsl

	transport.dtlsICEListener = &goDTLSICETransportListener{DTLSICETransportListener: dtlsl}
	transport.transport.SetListener(transport.dtlsICEListener)

	var address string
	var port int
	for _, candidate := range remoteCandidates {
		if candidate.GetType() == "relay" {
			address = candidate.GetRelAddr()
			port = candidate.GetRelPort()
		} else {
			address = candidate.GetAddress()
			port = candidate.GetPort()
		}
		bundle.AddRemoteCandidate(transport.username, address, uint16(port))
	}

	transport.localCandidates = localCandidates
	transport.remoteCandidates = remoteCandidates

	transport.incomingStreams = make(map[string]*IncomingStream)
	transport.outgoingStreams = make(map[string]*OutgoingStream)

	transport.incomingStreamTracks = make(map[string]*IncomingStreamTrack)
	transport.outgoingStreamTracks = make(map[string]*OutgoingStreamTrack)

	transport.onTransportStopListeners = make([]TransportStopListener, 0)

	transport.onIncomingTrackListeners = make([]IncomingTrackListener, 0)
	transport.onOutgoingTrackListeners = make([]OutgoingTrackListener, 0)

	return transport
}

// Dump  dump incoming and outgoint rtp and rtcp packets into a pcap file
func (t *Transport) Dump(filename string, incoming bool, outgoing bool, rtcp bool) bool {
	ret := t.transport.Dump(filename, incoming, outgoing, rtcp)
	if ret == 0 {
		return false
	}
	return true
}

// SetBandwidthProbing Enable/Disable bitrate probing
// This will send padding only RTX packets to allow bandwidth estimation algortithm to probe bitrate beyonf current sent values.
// The ammoung of probing bitrate would be limited by the sender bitrate estimation and the limit set on the setMaxProbing Bitrate.
func (t *Transport) SetBandwidthProbing(probe bool) {
	t.transport.SetBandwidthProbing(probe)
}

// SetMaxProbingBitrate Set the maximum bitrate to be used if probing is enabled.
func (t *Transport) SetMaxProbingBitrate(bitrate uint) {
	t.transport.SetMaxProbingBitrate(bitrate)
}

// GetDTLSState  get dtls state
func (t *Transport) GetDTLSState() string {
	return t.dtlsState
}

// GetICEStats  get ice stats
func (t *Transport) GetICEStats() ICEStats {

	iceStats := ICEStats{}
	iceStats.RequestsSent = t.connection.GetIceRequestsSent()
	iceStats.RequestsReceived = t.connection.GetIceRequestsReceived()
	iceStats.ResponsesSent = t.connection.GetIceResponsesSent()
	iceStats.ResponsesReceived = t.connection.GetIceResponsesReceived()

	return iceStats
}

// SetRemoteProperties  Set remote RTP properties
func (t *Transport) SetRemoteProperties(audio *sdp.MediaInfo, video *sdp.MediaInfo) {
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

	t.transport.SetRemoteProperties(properties)

	native.DeleteProperties(properties)
}

// SetLocalProperties Set local RTP properties
func (t *Transport) SetLocalProperties(audio *sdp.MediaInfo, video *sdp.MediaInfo) {

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

	t.transport.SetLocalProperties(properties)
	native.DeleteProperties(properties)
}

// GetLocalDTLSInfo Get transport local DTLS info
func (t *Transport) GetLocalDTLSInfo() *sdp.DTLSInfo {

	return t.localDtls
}

// GetLocalICEInfo Get transport local ICE info
func (t *Transport) GetLocalICEInfo() *sdp.ICEInfo {

	return t.localIce
}

// GetLocalCandidates Get local ICE candidates for this transport
func (t *Transport) GetLocalCandidates() []*sdp.CandidateInfo {

	return t.localCandidates
}

// GetRemoteCandidates Get remote ICE candidates for this transport
func (t *Transport) GetRemoteCandidates() []*sdp.CandidateInfo {
	return t.remoteCandidates
}

// AddRemoteCandidate register a remote candidate info. Only needed for ice-lite to ice-lite endpoints
func (t *Transport) AddRemoteCandidate(candidate *sdp.CandidateInfo) {

	var address string
	var port int

	if candidate.GetType() == "relay" {
		address = candidate.GetRelAddr()
		port = candidate.GetRelPort()
	} else {
		address = candidate.GetAddress()
		port = candidate.GetPort()
	}

	if t.bundle.AddRemoteCandidate(t.username, address, uint16(port)) != 0 {
		return
	}

	t.remoteCandidates = append(t.remoteCandidates, candidate)
}

// CreateOutgoingStream Create new outgoing stream in this transport using StreamInfo
func (t *Transport) CreateOutgoingStream(streamInfo *sdp.StreamInfo) *OutgoingStream {

	if _, ok := t.outgoingStreams[streamInfo.GetID()]; ok {
		return nil
	}

	info := streamInfo.Clone()
	outgoingStream := NewOutgoingStream(t.transport, info)

	outgoingStream.OnStop(func() {
		t.Lock()
		delete(t.outgoingStreams, outgoingStream.GetID())
		t.Unlock()
	})

	t.Lock()
	t.outgoingStreams[outgoingStream.GetID()] = outgoingStream
	t.Unlock()

	outgoingStream.OnTrack(func(track *OutgoingStreamTrack) {
		for _, trackFunc := range t.onOutgoingTrackListeners {
			trackFunc(track, outgoingStream)
		}
	})

	for _, track := range outgoingStream.GetTracks() {
		for _, trackFunc := range t.onOutgoingTrackListeners {
			trackFunc(track, outgoingStream)
		}
	}

	return outgoingStream
}

// CreateOutgoingStreamWithID  alias CreateOutgoingStream
func (t *Transport) CreateOutgoingStreamWithID(streamID string, audio bool, video bool) *OutgoingStream {

	streamInfo := sdp.NewStreamInfo(streamID)
	if audio {
		audioTrack := sdp.NewTrackInfo(uuid.Must(uuid.NewV4()).String(), "audio")
		ssrc := NextSSRC()
		audioTrack.AddSSRC(ssrc)
		streamInfo.AddTrack(audioTrack)
	}

	if video {
		videoTrack := sdp.NewTrackInfo(uuid.Must(uuid.NewV4()).String(), "video")
		ssrc := NextSSRC()
		videoTrack.AddSSRC(ssrc)
		streamInfo.AddTrack(videoTrack)
	}

	stream := t.CreateOutgoingStream(streamInfo)
	return stream
}

// CreateOutgoingStreamTrack Create new outgoing track in this transport
func (t *Transport) CreateOutgoingStreamTrack(media string, trackId string, ssrcs map[string]uint) *OutgoingStreamTrack {

	var mediaType native.MediaFrameType = 0
	if media == "video" {
		mediaType = 1
	}

	if trackId == "" {
		trackId = uuid.Must(uuid.NewV4()).String()
	}

	source := native.NewRTPOutgoingSourceGroup(mediaType)

	if ssrc, ok := ssrcs["media"]; ok {
		source.GetMedia().SetSsrc(ssrc)
	} else {
		source.GetMedia().SetSsrc(NextSSRC())
	}

	if ssrc, ok := ssrcs["rtx"]; ok {
		source.GetRtx().SetSsrc(ssrc)
	} else {
		source.GetRtx().SetSsrc(NextSSRC())
	}

	if ssrc, ok := ssrcs["fec"]; ok {
		source.GetFec().SetSsrc(ssrc)
	} else {
		source.GetFec().SetSsrc(NextSSRC())
	}

	// todo error handle
	t.transport.AddOutgoingSourceGroup(source)

	outgoingTrack := newOutgoingStreamTrack(media, trackId, native.TransportToSender(t.transport), source)

	outgoingTrack.OnStop(func() {
		t.transport.RemoveOutgoingSourceGroup(source)
	})

	for _, trackFunc := range t.onOutgoingTrackListeners {
		trackFunc(outgoingTrack, nil)
	}

	return outgoingTrack
}

// CreateIncomingStream Create an incoming stream object from the media stream info objet
func (t *Transport) CreateIncomingStream(streamInfo *sdp.StreamInfo) *IncomingStream {

	if _, ok := t.incomingStreams[streamInfo.GetID()]; ok {
		return nil
	}

	incomingStream := newIncomingStream(t.transport, native.TransportToReceiver(t.transport), streamInfo)

	t.Lock()
	t.incomingStreams[incomingStream.GetID()] = incomingStream
	t.Unlock()

	return incomingStream
}

// CreateIncomingStreamTrack Create new incoming stream in this transport. TODO: Simulcast is still not supported
// You can use IncomingStream's CreateTrack
func (t *Transport) CreateIncomingStreamTrack(media string, trackId string, ssrcs map[string]uint) *IncomingStreamTrack {

	var mediaType native.MediaFrameType = 0
	if media == "video" {
		mediaType = 1
	}

	if trackId == "" {
		trackId = uuid.Must(uuid.NewV4()).String()
	}

	source := native.NewRTPIncomingSourceGroup(mediaType, t.transport.GetTimeService())

	if ssrc, ok := ssrcs["media"]; ok {
		source.GetMedia().SetSsrc(ssrc)
	} else {
		source.GetMedia().SetSsrc(NextSSRC())
	}

	if ssrc, ok := ssrcs["rtx"]; ok {
		source.GetRtx().SetSsrc(ssrc)
	} else {
		source.GetRtx().SetSsrc(NextSSRC())
	}

	if ssrc, ok := ssrcs["fec"]; ok {
		source.GetFec().SetSsrc(ssrc)
	} else {
		source.GetFec().SetSsrc(NextSSRC())
	}

	t.transport.AddIncomingSourceGroup(source)

	sources := map[string]native.RTPIncomingSourceGroup{"": source}

	incomingTrack := NewIncomingStreamTrack(media, trackId, native.TransportToReceiver(t.transport), sources)

	// todo: remove callback
	incomingTrack.OnStop(func() {
		for _, item := range sources {
			t.transport.RemoveIncomingSourceGroup(item)
		}
	})

	for _, trackFunc := range t.onIncomingTrackListeners {
		trackFunc(incomingTrack, nil)
	}

	return incomingTrack
}

func (t *Transport) RemoveIncomingStream(incomingStream *IncomingStream) {

	t.Lock()
	delete(t.incomingStreams, incomingStream.GetID())
	t.Unlock()
}

// GetIncomingStreams get all incoming streams
func (t *Transport) GetIncomingStreams() []*IncomingStream {
	incomings := []*IncomingStream{}
	for _, stream := range t.incomingStreams {
		incomings = append(incomings, stream)
	}
	return incomings
}

// GetIncomingStream  get one incoming stream
func (t *Transport) GetIncomingStream(streamId string) *IncomingStream {
	t.Lock()
	defer t.Unlock()
	return t.incomingStreams[streamId]
}

// GetOutgoingStreams get all outgoing streams
func (t *Transport) GetOutgoingStreams() []*OutgoingStream {
	outgoings := []*OutgoingStream{}
	for _, stream := range t.outgoingStreams {
		outgoings = append(outgoings, stream)
	}
	return outgoings
}

// GetOutgoingStream get one outgoing stream
func (t *Transport) GetOutgoingStream(streamId string) *OutgoingStream {
	t.Lock()
	defer t.Unlock()
	return t.outgoingStreams[streamId]
}

// OnStop register a stop listener
func (t *Transport) OnStop(stop TransportStopListener) {
	t.onTransportStopListeners = append(t.onTransportStopListeners, stop)
}

// OnIncomingTrack register incoming track
func (t *Transport) OnIncomingTrack(listener IncomingTrackListener) {
	t.onIncomingTrackListeners = append(t.onIncomingTrackListeners, listener)
}

// OnOutgoingTrack register outgoing track
func (t *Transport) OnOutgoingTrack(listener OutgoingTrackListener) {
	t.onOutgoingTrackListeners = append(t.onOutgoingTrackListeners, listener)
}

// OnDTLSICEState  OnDTLSICEState
func (t *Transport) OnDTLSICEState(listener DTLSStateListener) {
	t.outDTLSStateListener = listener
}

// Stop stop this transport
func (t *Transport) Stop() {

	if t.bundle == nil {
		return
	}

	for _, incoming := range t.incomingStreams {
		incoming.Stop()
	}

	for _, outgoing := range t.outgoingStreams {
		outgoing.Stop()
	}

	if t.senderSideListener != nil {
		t.senderSideListener.deleteSenderSideEstimatorListener()
		t.senderSideListener = nil
	}

	if t.dtlsICEListener != nil {
		t.dtlsICEListener.deleteDTLSICETransportListener()
		t.dtlsICEListener = nil
	}

	t.bundle.RemoveICETransport(t.username)

	for _, stopFunc := range t.onTransportStopListeners {
		stopFunc()
	}

	native.DeleteStringFacade(t.username)

	t.incomingStreams = nil
	t.outgoingStreams = nil

	t.connection = nil
	t.transport = nil

	t.username = nil
	t.bundle = nil

}
