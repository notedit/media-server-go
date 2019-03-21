#ifndef SCTP_SHUTDOWNACKNOWLEDGEMENTCHUNK_H_
#define SCTP_SHUTDOWNACKNOWLEDGEMENTCHUNK_H_


#include "Buffer.h"
#include "sctp/Chunk.h"

namespace sctp
{

class ShutdownAcknowledgementChunk  : public Chunk
{
public:
	ShutdownAcknowledgementChunk () : Chunk(Chunk::SHUTDOWN_ACK) {}
	virtual ~ShutdownAcknowledgementChunk() = default;
	
	virtual size_t Serialize(BufferWritter& buffer) const override;
	virtual size_t GetSize() const override;

	static Chunk::shared Parse(BufferReader& reader);
public:
	//        0                   1                   2                   3
	//        0 1 2 3 4 5 6 7 8 9 0 1 2 3 4 5 6 7 8 9 0 1 2 3 4 5 6 7 8 9 0 1
	//       +-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+
	//       |   Type = 8    | Chunk  Flags  |      Length = 8               |
	//       +-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+
};

}; // namespace sctp

#endif

