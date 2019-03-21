#ifndef SCTP_RECONFIGCHUNK_H_
#define SCTP_RECONFIGCHUNK_H_

#include "Buffer.h"
#include "sctp/Chunk.h"

namespace sctp
{

class ReConfigChunk  : public Chunk
{
public:
	ReConfigChunk () : Chunk(Chunk::RE_CONFIG) {}
	virtual ~ReConfigChunk() = default;
	
	virtual size_t Serialize(BufferWritter& buffer) const override;
	virtual size_t GetSize() const override;

	static Chunk::shared Parse(BufferReader& reader);
public:

	//	0                   1                   2                   3
	//	0 1 2 3 4 5 6 7 8 9 0 1 2 3 4 5 6 7 8 9 0 1 2 3 4 5 6 7 8 9 0 1
	//	+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+
	//	| Type = 130    |  Chunk Flags  |      Chunk Length             |
	//	+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+
	//	\                                                               \
	//	/                  Re-configuration Parameter                   /
	//	\                                                               \
	//	+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+
	//	\                                                               \
	//	/             Re-configuration Parameter (optional)             /
	//	\                                                               \
	//	+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+

	Buffer cookie;
};
	
}; // namespace sctp

#endif


