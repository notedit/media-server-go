#ifndef SCTP_SHUTDOWNASSOCIATIONCHUNK_H_
#define SCTP_SHUTDOWNASSOCIATIONCHUNK_H_


#include "Buffer.h"
#include "sctp/Chunk.h"

namespace sctp
{

class ShutdownAssociationChunk  : public Chunk
{
public:
	ShutdownAssociationChunk () : Chunk(Chunk::SHUTDOWN) {}
	virtual ~ShutdownAssociationChunk() = default;
	
	virtual size_t Serialize(BufferWritter& buffer) const override;
	virtual size_t GetSize() const override;

	static Chunk::shared Parse(BufferReader& reader);
public:
	//        0                   1                   2                   3
	//        0 1 2 3 4 5 6 7 8 9 0 1 2 3 4 5 6 7 8 9 0 1 2 3 4 5 6 7 8 9 0 1
	//       +-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+
	//       |   Type = 7    | Chunk  Flags  |      Length = 8               |
	//       +-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+
	//       |                      Cumulative TSN Ack                       |
	//       +-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+
	uint32_t cumulativeTrasnmissionSequenceNumberAck;
};
	
}; // namespace sctp

#endif

