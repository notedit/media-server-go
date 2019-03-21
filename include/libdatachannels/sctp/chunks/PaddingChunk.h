#ifndef SCTP_PADDINGCHUNK_H_
#define SCTP_PADDINGCHUNK_H_

#include "Buffer.h"
#include "sctp/Chunk.h"

namespace sctp
{

class PaddingChunk  : public Chunk
{
public:
	PaddingChunk() : Chunk(PAD) {}
	virtual ~PaddingChunk() = default;
	
	virtual size_t Serialize(BufferWritter& buffer) const override;
	virtual size_t GetSize() const override;
	
	static Chunk::shared Parse(BufferReader& reader);
public:
	//    0                   1                   2                   3
	//    0 1 2 3 4 5 6 7 8 9 0 1 2 3 4 5 6 7 8 9 0 1 2 3 4 5 6 7 8 9 0 1
	//   +-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+
	//   | Type = 0x84   |   Flags=0     |             Length            |
	//   +-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+
	//   |                                                               |
	//   \                         Padding Data                          /
	//   /                                                               \
	//   +-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+
	Buffer buffer;
};
	
}; // namespace sctp

#endif

