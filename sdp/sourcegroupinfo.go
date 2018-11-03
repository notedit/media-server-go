package sdp

type SourceGroupInfo struct {
	semantics string
	ssrcs     []int
}

func NewSourceGroupInfo(semantics string, ssrcs []int) *SourceGroupInfo {

	info := &SourceGroupInfo{
		semantics: semantics,
		ssrcs:     make([]int, len(ssrcs)),
	}

	copy(info.ssrcs, ssrcs)
	return info
}

func (s *SourceGroupInfo) Clone() *SourceGroupInfo {

	cloned := &SourceGroupInfo{
		semantics: s.semantics,
		ssrcs:     make([]int, len(s.ssrcs)),
	}

	copy(cloned.ssrcs, s.ssrcs)
	return cloned
}

func (s *SourceGroupInfo) GetSemantics() string {

	return s.semantics
}

func (s *SourceGroupInfo) GetSSRCs() []int {

	return s.ssrcs
}
