#ifndef SCTP_UNKNOWNCHUNK_H_
#define SCTP_UNKNOWNCHUNK_H_

#include "Buffer.h"
#include "sctp/Chunk.h"

namespace sctp
{

class UnknownChunk  : public Chunk
{
public:
	UnknownChunk(uint8_t type) : Chunk(type) {}
	virtual ~UnknownChunk() = default;
	
	virtual size_t Serialize(BufferWritter& buffer) const override;
	virtual size_t GetSize() const override;

	static Chunk::shared Parse(BufferReader& reader);
public:
	//        0                   1                   2                   3
	//        0 1 2 3 4 5 6 7 8 9 0 1 2 3 4 5 6 7 8 9 0 1 2 3 4 5 6 7 8 9 0 1
	//       +-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+
	//       |   Type = ??   |Chunk  Flags   |         Length                |
	//       +-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+
	//       /                     Buffer                                    /
	//       \                                                               \
	//       +-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+
	Buffer buffer;
};
	
}; // namespace sctp

#endif

