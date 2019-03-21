#ifndef SCTP_PAYLOADDATACHUNK_H_
#define SCTP_PAYLOADDATACHUNK_H_


#include "Buffer.h"
#include "sctp/Chunk.h"

namespace sctp
{

	
class PayloadDataChunk :  public Chunk
{
public:
	PayloadDataChunk () : Chunk(Chunk::PDATA) {}
	virtual ~PayloadDataChunk() = default;
	
	virtual size_t Serialize(BufferWritter& buffer) const override;
	virtual size_t GetSize() const override;

	static Chunk::shared Parse(BufferReader& reader);
public:
	//        0                   1                   2                   3
	//        0 1 2 3 4 5 6 7 8 9 0 1 2 3 4 5 6 7 8 9 0 1 2 3 4 5 6 7 8 9 0 1
	//       +-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+
	//       |   Type = 0    | Reserved|U|B|E|    Length                     |
	//       +-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+
	//       |                              TSN                              |
	//       +-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+
	//       |      Stream Identifier S      |   Stream Sequence Number n    |
	//       +-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+
	//       |                  Payload Protocol Identifier                  |
	//       +-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+
	//       \                                                               \
	//       /                 User Data (seq n of Stream S)                 /
	//       \                                                               \
	//       +-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+	
	bool unordered				= false;
	bool beginingFragment			= false;
	bool endingFragment			= false;
	uint32_t transmissionSequenceNumber     = 0;
	uint16_t streamIdentifier		= 0;
	uint16_t streamSequenceNumber		= 0;
	uint32_t payloadProtocolIdentifier	= 0;
	Buffer userData				= 0;
};


}; // namespace sctp

#endif

