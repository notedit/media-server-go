#ifndef SCTP_SHUTDOWNCOMPLETECHUNK_H_
#define SCTP_SHUTDOWNCOMPLETECHUNK_H_


#include "Buffer.h"
#include "sctp/Chunk.h"

namespace sctp
{

class ShutdownCompleteChunk  : public Chunk
{
public:
	ShutdownCompleteChunk () : Chunk(Chunk::SHUTDOWN_COMPLETE) {}
	virtual ~ShutdownCompleteChunk() = default;
	
	virtual size_t Serialize(BufferWritter& buffer) const override;
	virtual size_t GetSize() const override;

	static Chunk::shared Parse(BufferReader& reader);
public:
	//        0                   1                   2                   3
	//        0 1 2 3 4 5 6 7 8 9 0 1 2 3 4 5 6 7 8 9 0 1 2 3 4 5 6 7 8 9 0 1
	//       +-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+
	//       |   Type = 14   |Reserved     |T|      Length = 4               |
	//       +-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+
	bool verificationTag;
};

}; // namespace sctp

#endif

