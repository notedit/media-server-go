#ifndef SCTP_COOKIEACKCHUNK_H_
#define SCTP_COOKIEACKCHUNK_H_

#include "Buffer.h"
#include "sctp/Chunk.h"

namespace sctp
{
	
class CookieAckChunk  : public Chunk
{
public:
	CookieAckChunk () : Chunk(Chunk::COOKIE_ACK) {}
	virtual ~CookieAckChunk() = default;
	
	virtual size_t Serialize(BufferWritter& buffer) const override;
	virtual size_t GetSize() const override;

	static Chunk::shared Parse(BufferReader& reader);
public:
	//        0                   1                   2                   3
	//        0 1 2 3 4 5 6 7 8 9 0 1 2 3 4 5 6 7 8 9 0 1 2 3 4 5 6 7 8 9 0 1
	//       +-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+
	//       |   Type = 11   |Chunk  Flags   |     Length = 4                |
	//       +-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+
	//
};

}; // namespace sctp

#endif

