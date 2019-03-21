#ifndef SCTP_OPERATIONERRORCHUNK_H_
#define SCTP_OPERATIONERRORCHUNK_H_


#include <vector>
#include "Buffer.h"
#include "sctp/Chunk.h"
#include "sctp/ErrorCause.h"

namespace sctp
{

class OperationErrorChunk  : public Chunk
{
public:
	OperationErrorChunk () : Chunk(Chunk::ERROR ) {}
	virtual ~OperationErrorChunk() = default;
	
	virtual size_t Serialize(BufferWritter& buffer) const override;
	virtual size_t GetSize() const override;

	static Chunk::shared Parse(BufferReader& reader);
public:
	//        0                   1                   2                   3
	//        0 1 2 3 4 5 6 7 8 9 0 1 2 3 4 5 6 7 8 9 0 1 2 3 4 5 6 7 8 9 0 1
	//       +-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+
	//       |   Type = 9    | Chunk  Flags  |           Length              |
	//       +-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+
	//       \                                                               \
	//       /                    one or more Error Causes                   /
	//       \                                                               \
	//       +-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+
	std::vector<ErrorCause> errorCauses;
};

}; // namespace sctp

#endif

