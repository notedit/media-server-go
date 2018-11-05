package mediaserver

import (
	sdp "./sdp"
)

type Transport struct {
	localIce      *sdp.ICEInfo
	localDtls     *sdp.DTLSInfo
	remoteIce     *sdp.ICEInfo
	remoteDtls    *sdp.DTLSInfo
	bundle        RTPBundleTransport
	dtlsTransport DTLSICETransport
	username      StringFacade
}

func NewTransport(bundle RTPBundleTransport, remoteIce *sdp.ICEInfo, remoteDtls *sdp.DTLSInfo, remoteCandidates []*sdp.CandidateInfo,
	localIce *sdp.ICEInfo, localDtls *sdp.DTLSInfo, localCandidates []*sdp.CandidateInfo, disableSTUNKeepAlive bool) *Transport {

	transport := &Transport{}
	transport.remoteIce = remoteIce
	transport.remoteDtls = remoteDtls
	transport.localIce = localIce
	transport.localDtls = localDtls
	transport.bundle = bundle

	properties := NewProperties()

	properties.SetProperty("ice.localUsername", localIce.GetUfrag())
	properties.SetProperty("ice.localPassword", localIce.GetPassword())
	properties.SetProperty("ice.remoteUsername", remoteIce.GetUfrag())
	properties.SetProperty("ice.remotePassword", remoteIce.GetPassword())

	properties.SetProperty("dtls.setup", remoteDtls.GetSetup().String())
	properties.SetProperty("dtls.hash", remoteDtls.GetHash())
	properties.SetProperty("dtls.fingerprint", remoteDtls.GetFingerprint())

	properties.SetProperty("disableSTUNKeepAlive", disableSTUNKeepAlive)

	// todo set srtpProtectionProfiles, when will we use this?

	transport.username = NewStringFacade(localIce.GetUfrag() + ":" + remoteIce.GetUfrag())
	transport.dtlsTransport = bundle.AddICETransport(transport.username, properties)

	// todo ontargetbitrate callback
	// SenderSideEstimatorListener
	var address string
	var port int
	for _, candidate := range candidates {
		if candidate.GetType() == "relay" {
			address = candidate.GetRelAddr()
			port = candidate.GetRelPort()
		} else {
			address = candidate.GetAddress()
			port = candidate.GetPort()
		}
		bundle.AddRemoteCandidate(transport.username, address, uint16(port))
	}

	// todo

	return transport
}
