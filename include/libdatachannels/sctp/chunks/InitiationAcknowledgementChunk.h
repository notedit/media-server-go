#ifndef SCTP_INITIATIONACKNOWLEDGEMENTCHUNK_H_
#define SCTP_INITIATIONACKNOWLEDGEMENTCHUNK_H_

#include "Buffer.h"
#include "sctp/Chunk.h"

#include <optional>

namespace sctp
{

class InitiationAcknowledgementChunk  : public Chunk
{
public:
	InitiationAcknowledgementChunk () : Chunk(Chunk::INIT_ACK) {}
	virtual ~InitiationAcknowledgementChunk() = default;
	
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
	Buffer stateCookie;
	
	std::vector<std::array<uint8_t, 8>>  ipV4Addresses;		// IPv4 Address Parameter (5)
	std::vector<std::array<uint8_t, 20>> ipV6Addresses;		// IPv6 Address Parameter (6)
	std::optional<std::string> hostName;				// Host Name Address (11)
	std::vector<Buffer> unrecognizedParameters;			// Unrecognized Parameter (8)
	std::vector<std::pair<uint8_t,Buffer>> unknownParameters;
	std::vector<uint8_t> supportedExtensions;			// Supported Extensions  (0x8008) rfc5061#page-13
	
	bool forwardTSNSupported = true;				//  Forward-TSN-Supported 49152 (0xC000) rfc3758
};	
	
}; // namespace sctp

#endif

