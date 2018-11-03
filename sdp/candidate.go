package sdp

type CandidateInfo struct {
	foundation  string
	componentID int
	transport   string
	priority    int
	address     string
	port        int
	ctype       string
	relAddr     string
	relPort     int
}

func NewCandidateInfo(foundation string, componentID int, transport string,
	priority int, address string, port int, ctype string, relAddr string, relPort int) *CandidateInfo {

	candidate := &CandidateInfo{
		foundation:  foundation,
		componentID: componentID,
		transport:   transport,
		priority:    priority,
		address:     address,
		port:        port,
		ctype:       ctype,
		relAddr:     relAddr,
		relPort:     relPort,
	}
	return candidate

}

func (c *CandidateInfo) Clone() *CandidateInfo {
	candidate := &CandidateInfo{
		foundation:  c.foundation,
		componentID: c.componentID,
		transport:   c.transport,
		priority:    c.priority,
		address:     c.address,
		port:        c.port,
		ctype:       c.ctype,
		relAddr:     c.relAddr,
		relPort:     c.relPort,
	}
	return candidate
}

func (c *CandidateInfo) GetFoundation() string {
	return c.foundation
}

func (c *CandidateInfo) GetComponentID() int {
	return c.componentID
}

func (c *CandidateInfo) GetTransport() string {
	return c.transport
}

func (c *CandidateInfo) GetPriority() int {
	return c.priority
}

func (c *CandidateInfo) GetAddress() string {
	return c.address
}

func (c *CandidateInfo) GetPort() int {
	return c.port
}

func (c *CandidateInfo) GetType() string {
	return c.ctype
}

func (c *CandidateInfo) GetRelAddr() string {
	return c.relAddr
}

func (c *CandidateInfo) GetRelPort() int {
	return c.relPort
}
