#ifndef SCTP_PACKET_HEADER_H_
#define SCTP_PACKET_HEADER_H_
#include <stdint.h>
#include <memory>
#include <list>
#include "BufferWritter.h"
#include "BufferReader.h"
#include "sctp/Chunk.h"

namespace sctp
{

class PacketHeader
{
public:
	using shared = std::shared_ptr<PacketHeader>;
public:	
	PacketHeader(uint16_t sourcePortNumber,uint16_t destinationPortNumber,uint32_t verificationTag, uint32_t checksum = 0);
	~PacketHeader() = default;
	
	static PacketHeader::shared Parse(BufferReader& buffer) ;
	size_t Serialize(BufferWritter& buffer) const;
	size_t GetSize() const;
public:
	uint16_t sourcePortNumber	 = 0;
	uint16_t destinationPortNumber	 = 0;
	uint32_t verificationTag	 = 0;
	uint32_t checksum		 = 0;
};

}; // namespace sctp

#endif


