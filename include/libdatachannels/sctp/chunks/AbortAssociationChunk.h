#ifndef SCTP_ABORTASSOCIATIONCHUNK_H_
#define SCTP_ABORTASSOCIATIONCHUNK_H_

#include <stdint.h>
#include <vector>
#include "Buffer.h"
#include "sctp/Chunk.h"
#include "sctp/ErrorCause.h"

namespace sctp
{

class AbortAssociationChunk  : public Chunk
{
public:
	AbortAssociationChunk () : Chunk(Chunk::ABORT ) {}
	virtual ~AbortAssociationChunk() = default;
	
	virtual size_t Serialize(BufferWritter& buffer) const override;
	virtual size_t GetSize() const override;

	static Chunk::shared Parse(BufferReader& reader);
public:
	//        0                   1                   2                   3
	//        0 1 2 3 4 5 6 7 8 9 0 1 2 3 4 5 6 7 8 9 0 1 2 3 4 5 6 7 8 9 0 1
	//       +-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+
	//       |   Type = 6    |Reserved     |T|           Length              |
	//       +-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+
	//       \                                                               \
	//       /                   zero or more Error Causes                   /
	//       \                                                               \
	//       +-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+
	//	
	bool verificationTag;
	std::vector<ErrorCause> errorCauses;
};
	
}; // namespace sctp

#endif

