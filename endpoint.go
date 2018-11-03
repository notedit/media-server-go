package mediaserver

type Endpoint struct {
	ip          string
	bundle      RTPBundleTransport
	transports  map[string]*Transport
	candidate   interface{}
	fingerprint string
}

func NewEndpoint(ip string) *Endpoint {
	bundle := NewRTPBundleTransport()
	bundle.Init()
	fingerprint := MediaServerGetFingerprint().ToString()
	// candidate todo
	return &Endpoint{
		ip:          ip,
		bundle:      bundle,
		transports:  make(map[string]*Transport),
		fingerprint: fingerprint,
	}
}

func (e *Endpoint) CreateTransport() {

}

func (e *Endpoint) GetLocalCandidates() {

}

func (e *Endpoint) GetDTLSFingerprint() string {
	return e.fingerprint
}

func (e *Endpoint) CreateOffer() {

}

func (e *Endpoint) Stop() {

}
