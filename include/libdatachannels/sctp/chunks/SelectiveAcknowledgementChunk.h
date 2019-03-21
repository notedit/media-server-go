#ifndef SCTP_SELECTIVEACKNOWLEDGEMENTCHUNK_H_
#define SCTP_SELECTIVEACKNOWLEDGEMENTCHUNK_H_


#include <vector>
#include <utility>
#include "Buffer.h"
#include "sctp/Chunk.h"

namespace sctp
{

class SelectiveAcknowledgementChunk  : public Chunk
{
public:
	SelectiveAcknowledgementChunk () : Chunk(Chunk::SACK) {}
	virtual ~SelectiveAcknowledgementChunk() = default;
	
	virtual size_t Serialize(BufferWritter& buffer) const override;
	virtual size_t GetSize() const override;

	static Chunk::shared Parse(BufferReader& reader);
public:
	//        0                   1                   2                   3
	//        0 1 2 3 4 5 6 7 8 9 0 1 2 3 4 5 6 7 8 9 0 1 2 3 4 5 6 7 8 9 0 1
	//       +-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+
	//       |   Type = 3    |Chunk  Flags   |      Chunk Length             |
	//       +-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+
	//       |                      Cumulative TSN Ack                       |
	//       +-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+
	//       |          Advertised Receiver Window Credit (a_rwnd)           |
	//       +-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+
	//       | Number of Gap Ack Blocks = N  |  Number of Duplicate TSNs = X |
	//       +-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+
	//       |  Gap Ack Block #1 Start       |   Gap Ack Block #1 End        |
	//       +-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+
	//       /                                                               /
	//       \                              ...                              \
	//       /                                                               /
	//       +-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+
	//       |   Gap Ack Block #N Start      |  Gap Ack Block #N End         |
	//       +-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+
	//       |                       Duplicate TSN 1                         |
	//       +-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+
	//       /                                                               /
	//       \                              ...                              \
	//       /                                                               /
	//       +-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+
	//       |                       Duplicate TSN X                         |
	//       +-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+
	
	uint32_t cumulativeTrasnmissionSequenceNumberAck = 0;
	uint32_t adveritsedReceiverWindowCredit = 0;
	std::vector<std::pair<uint16_t,uint16_t>> gapAckBlocks;
	std::vector<uint32_t> duplicateTuplicateTrasnmissionSequenceNumbers;
};	
	
}; // namespace sctp

#endif

