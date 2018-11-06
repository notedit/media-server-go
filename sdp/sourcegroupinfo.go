package sdp

type SourceGroupInfo struct {
	semantics string
	ssrcs     []uint
}

func NewSourceGroupInfo(semantics string, ssrcs []uint) *SourceGroupInfo {

	info := &SourceGroupInfo{
		semantics: semantics,
		ssrcs:     make([]uint, len(ssrcs)),
	}

	copy(info.ssrcs, ssrcs)
	return info
}

func (s *SourceGroupInfo) Clone() *SourceGroupInfo {

	cloned := &SourceGroupInfo{
		semantics: s.semantics,
		ssrcs:     make([]uint, len(s.ssrcs)),
	}

	copy(cloned.ssrcs, s.ssrcs)
	return cloned
}

func (s *SourceGroupInfo) GetSemantics() string {

	return s.semantics
}

func (s *SourceGroupInfo) GetSSRCs() []uint {

	return s.ssrcs
}
