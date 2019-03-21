#ifndef SCTP_HEARTBEATREQUESTCHUNK_H_
#define SCTP_HEARTBEATREQUESTCHUNK_H_

#include "Buffer.h"
#include "sctp/Chunk.h"

namespace sctp
{
	
class HeartbeatRequestChunk  : public Chunk
{
public:
	HeartbeatRequestChunk () : Chunk(Chunk::HEARTBEAT) {}
	virtual ~HeartbeatRequestChunk() = default;
	
	virtual size_t Serialize(BufferWritter& buffer) const override;
	virtual size_t GetSize() const override;

	static Chunk::shared Parse(BufferReader& reader);
public:
	//        0                   1                   2                   3
	//        0 1 2 3 4 5 6 7 8 9 0 1 2 3 4 5 6 7 8 9 0 1 2 3 4 5 6 7 8 9 0 1
	//       +-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+
	//       |   Type = 4    | Chunk  Flags  |      Heartbeat Length         |
	//       +-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+
	//       |    Heartbeat Info Type=1      |         HB Info Length        |
	//       +-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+
	//       /                  Sender-Specific Heartbeat Info               /
	//       \                                                               \
	//       +-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+
	
	//        0                   1                   2                   3
	//        0 1 2 3 4 5 6 7 8 9 0 1 2 3 4 5 6 7 8 9 0 1 2 3 4 5 6 7 8 9 0 1
	//       +-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+
	//       |    Heartbeat Info Type=1      |         HB Info Length        |
	//       +-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+
	//       /                  Sender-Specific Heartbeat Info               /
	//       \                                                               \
	//       +-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+	
	
	Buffer senderSpecificHearbeatInfo;
};

}; // namespace sctp

#endif

