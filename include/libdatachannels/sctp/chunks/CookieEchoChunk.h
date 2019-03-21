#ifndef SCTP_COOKIEECHOCHUNK_H_
#define SCTP_COOKIEECHOCHUNK_H_

#include "Buffer.h"
#include "sctp/Chunk.h"

namespace sctp
{

class CookieEchoChunk  : public Chunk
{
public:
	CookieEchoChunk () : Chunk(Chunk::COOKIE_ECHO) {}
	virtual ~CookieEchoChunk() = default;
	
	virtual size_t Serialize(BufferWritter& buffer) const override;
	virtual size_t GetSize() const override;

	static Chunk::shared Parse(BufferReader& reader);
public:
	//        0                   1                   2                   3
	//        0 1 2 3 4 5 6 7 8 9 0 1 2 3 4 5 6 7 8 9 0 1 2 3 4 5 6 7 8 9 0 1
	//       +-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+
	//       |   Type = 10   |Chunk  Flags   |         Length                |
	//       +-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+
	//       /                     Cookie                                    /
	//       \                                                               \
	//       +-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+

	Buffer cookie;
};
	
}; // namespace sctp

#endif

