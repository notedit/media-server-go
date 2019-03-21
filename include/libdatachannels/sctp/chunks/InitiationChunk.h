#ifndef SCTP_INITIATIONCHUNK_H_
#define SCTP_INITIATIONCHUNK_H_


#include <array>
#include <vector>
#include <optional>
#include "Buffer.h"
#include "sctp/Chunk.h"

namespace sctp
{

class InitiationChunk : public Chunk
{
public:
	InitiationChunk () : Chunk(Chunk::INIT) {}
	virtual ~InitiationChunk() = default;
	
	virtual size_t Serialize(BufferWritter& buffer) const override;
	virtual size_t GetSize() const override;

	static Chunk::shared Parse(BufferReader& reader);
public:	
	//	0                   1                   2                   3
	//	0 1 2 3 4 5 6 7 8 9 0 1 2 3 4 5 6 7 8 9 0 1 2 3 4 5 6 7 8 9 0 1
	//	+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+
	//	|   Type = 1    |  Chunk Flags  |      Chunk Length             |
	//	+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+
	//	|                         Initiate Tag                          |
	//	+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+
	//	|           Advertised Receiver Window Credit (a_rwnd)          |
	//	+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+
	//	|  Number of Outbound Streams   |  Number of Inbound Streams    |
	//	+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+
	//	|                          Initial TSN                          |
	//	+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+
	//	\                                                               \
	//	/              Optional/Variable-Length Parameters              /
	//	\                                                               \
	//	+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+
	uint32_t initiateTag = 0;
	uint32_t advertisedReceiverWindowCredit = 0;
	uint16_t numberOfOutboundStreams = 0;
	uint16_t numberOfInboundStreams = 0;
	uint32_t initialTransmissionSequenceNumber = 0;
	//Optional parameters
	std::vector<std::array<uint8_t, 8>>   ipV4Addresses;		// IPv4 Address Parameter (5)
	std::vector<std::array<uint8_t, 20>>  ipV6Addresses;		// IPv6 Address Parameter (6)
	std::optional<uint64_t>	suggestedCookieLifeSpanIncrement;	// Cookie Preservative (9)
	std::optional<std::string> hostName;				// Host Name Address (11)
	std::vector<uint16_t> supportedAddressTypes;			// Supported Address Types (12)
	std::vector<uint8_t> supportedExtensions;			// Supported Extensions  (0x8008) rfc5061#page-13
	std::vector<std::pair<uint8_t,Buffer>> unknownParameters;
	bool forwardTSNSupported = true;				//  Forward-TSN-Supported 49152 (0xC000) rfc3758
};
	
}; // namespace sctp

#endif

