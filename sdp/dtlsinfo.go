package sdp

type DTLSInfo struct {
	setup       Setup
	hash        string
	fingerprint string
}

func NewDTLSInfo(setup Setup, hash string, fingerprint string) *DTLSInfo {

	return &DTLSInfo{
		setup:       setup,
		hash:        hash,
		fingerprint: fingerprint,
	}
}

func (d *DTLSInfo) Clone() *DTLSInfo {
	dtls := &DTLSInfo{
		setup:       d.setup,
		hash:        d.hash,
		fingerprint: d.fingerprint,
	}
	return dtls
}

func (d *DTLSInfo) GetFingerprint() string {
	return d.fingerprint
}

func (d *DTLSInfo) GetHash() string {
	return d.hash
}

func (d *DTLSInfo) GetSetup() Setup {
	return d.setup
}

func (d *DTLSInfo) SetSetup(setup Setup) {
	d.setup = setup
}
