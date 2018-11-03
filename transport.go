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
}

func NewTransport(bundle RTPBundleTransport, remoteIce *sdp.ICEInfo, remoteDtls *sdp.DTLSInfo,
	localIce *sdp.ICEInfo, localDtls *sdp.DTLSInfo, disableSTUNKeepAlive bool) *Transport {

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

	properties.SetProperty("disableSTUNKeepAlive", disableSTUNKeepAlive)

	// todo set srtpProtectionProfiles, when will we use this?

	username := NewStringFacade(localIce.GetUfrag() + ":" + remoteIce.GetUfrag())
	transport.dtlsTransport = bundle.AddICETransport(username, properties)

	return transport
}
