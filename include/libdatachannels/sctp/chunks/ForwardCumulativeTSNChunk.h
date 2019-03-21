#ifndef SCTP_FORWARDCUMULATIVETSNCHUNK_H
#define SCTP_FORWARDCUMULATIVETSNCHUNK_H

#include <map>

#include "sctp/Chunk.h"

namespace sctp
{
	

class ForwardCumulativeTSNChunk : public Chunk
{
public:
	ForwardCumulativeTSNChunk () : Chunk(Chunk::FORWARD_CUMULATIVE_TSN) {}
	virtual ~ForwardCumulativeTSNChunk() = default;
	
	virtual size_t Serialize(BufferWritter& buffer) const override;
	virtual size_t GetSize() const override;

	static Chunk::shared Parse(BufferReader& reader);
public:
	//    0                   1                   2                   3
	//    0 1 2 3 4 5 6 7 8 9 0 1 2 3 4 5 6 7 8 9 0 1 2 3 4 5 6 7 8 9 0 1
	//   +-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+
	//   |   Type = 192  |  Flags = 0x00 |        Length = Variable      |
	//   +-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+
	//   |                      New Cumulative TSN                       |
	//   +-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+
	//   |         Stream-1              |       Stream Sequence-1       |
	//   +-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+
	//   \                                                               /
	//   /                                                               \
	//   +-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+
	//   |         Stream-N              |       Stream Sequence-N       |
	//   +-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+	
	uint32_t newCumulativeTSN;
	std::map<uint16_t,uint16_t> streamsSequence;

};

}; // namespace

#endif /* SCTP_FORWARDCUMULATIVETSNCHUNK_H */

