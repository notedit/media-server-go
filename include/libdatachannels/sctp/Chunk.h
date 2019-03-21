#ifndef SCTP_CHUNK_H_
#define SCTP_CHUNK_H_
#include <stdint.h>
#include <memory>

#include "Buffer.h"
#include "BufferReader.h"
#include "BufferWritter.h"

namespace sctp
{

inline size_t SizePad(size_t size, size_t num)
{
	return ((size + num -1) / num ) * num;
}

class Chunk
{
public:
	using shared = std::shared_ptr<Chunk>;
	
	enum Type
	{
		PDATA			= 0,  //  Payload Data (DATA)
		INIT			= 1,  //  Initiation (INIT)
		INIT_ACK		= 2,  //  Initiation Acknowledgement (INIT ACK)
		SACK			= 3,  //  Selective Acknowledgement (SACK)
		HEARTBEAT		= 4,  //  Heartbeat Request ()
		HEARTBEAT_ACK		= 5,  //  Heartbeat Acknowledgement (HEARTBEAT ACK)
		ABORT			= 6,  //  Abort (ABORT)
		SHUTDOWN		= 7,  //  Shutdown (SHUTDOWN)
		SHUTDOWN_ACK		= 8,  //  Shutdown Acknowledgement (SHUTDOWN ACK)
		ERROR			= 9,  //  Operation Error (ERROR)
		COOKIE_ECHO		= 10, //  State Cookie (COOKIE ECHO)
		COOKIE_ACK		= 11, //  Cookie Acknowledgement (COOKIE ACK),
		ECNE			= 12, //  Reserved for Explicit Congestion Notification Echo
		CWR			= 13, //  Reserved for Congestion Window Reduced (CWR)
		SHUTDOWN_COMPLETE	= 14, //  Shutdown Complete
		// 15 to 62   - available
		// 63         - reserved for IETF-defined Chunk Extensions
		// 64 to 126  - available
		PAD			= 84,
		// 127        - reserved for IETF-defined Chunk Extensions
		// 128 to 190 - available
		RE_CONFIG		= 130, // Re-configuration Chunk (RE-CONFIG) rfc6525
		// 191        - reserved for IETF-defined Chunk Extensions
		// 192 to 254 - available
		FORWARD_CUMULATIVE_TSN	= 192,
		// 255        - reserved for IETF-defined Chunk Extensions
	};
	
	enum Parameter
	{
		
		HeartbeatInfo				= 1,
		IPv4Address				= 5,
		IPv6Address				= 6,
		StateCookie				= 7,
		UnrecognizedParameter			= 8,
		CookiePreservative			= 9,
		ReservedforECNCapable			= 32768,
		HostNameAddress				= 11,
		SupportedAddressTypes			= 12,
		OutgoingSSNResetRequestParameter	= 13,
		IncomingSSNResetRequestParameter	= 14,
		SSNTSNResetRequestParameter		= 15,
		ReCconfigurationResponseParameter	= 16,
		AddOutgoingStreamsRequestParameter	= 17,
		AddIncomingStreamsRequestParameter	= 18,
		SupportedExtensions			= 0x8008, //rfc5061
		ForwardTSNSupported			= 49152,  //rfc3758 (0xC000)
		Padding					= 0x8005, //rfc480
	};
	
	Chunk(uint8_t type)
	{
		this->type = type;
	}
	virtual ~Chunk() = default;
	
	static Chunk::shared Parse(BufferReader& buffer);
	virtual size_t GetSize() const = 0;
	virtual size_t Serialize(BufferWritter& buffer) const = 0;
	
public:
	//	Chunk Types are encoded such that the highest-order 2 bits specify
	//	the action that must be taken if the processing endpoint does not
	//	recognize the Chunk Type.
	//
	//	00 -  Stop processing this SCTP packet and discard it, do not
	//	    process any further chunks within it.
	//
	//	01 -  Stop processing this SCTP packet and discard it, do not
	//	    process any further chunks within it, and report the
	//	    unrecognized chunk in an 'Unrecognized Chunk Type'.
	//
	//	10 -  Skip this chunk and continue processing.
	//
	//	11 -  Skip this chunk and continue processing, but report in an
	//	    ERROR chunk using the 'Unrecognized Chunk Type' cause of
	static bool SkipOnUnknown(uint8_t type)		{ return type & 0x80;	}
	static bool ReportOnUnknown(uint8_t type)	{ return type & 0x40;	}	
	
public:	
	uint8_t type;
	uint8_t flag = 0;
};

}; // namespace sctp

#include "sctp/chunks/AbortAssociationChunk.h"
#include "sctp/chunks/HeartbeatRequestChunk.h"
#include "sctp/chunks/HeartbeatAckChunk.h"
#include "sctp/chunks/OperationErrorChunk.h"
#include "sctp/chunks/SelectiveAcknowledgementChunk.h"
#include "sctp/chunks/ShutdownCompleteChunk.h"
#include "sctp/chunks/CookieAckChunk.h"
#include "sctp/chunks/InitiationAcknowledgementChunk.h"
#include "sctp/chunks/PayloadDataChunk.h"
#include "sctp/chunks/ShutdownAcknowledgementChunk.h"
#include "sctp/chunks/CookieEchoChunk.h"
#include "sctp/chunks/InitiationChunk.h"
#include "sctp/chunks/ReConfigChunk.h"
#include "sctp/chunks/ShutdownAssociationChunk.h"
#include "sctp/chunks/ForwardCumulativeTSNChunk.h"
#include "sctp/chunks/UnknownChunk.h"
#include "sctp/chunks/PaddingChunk.h"



#endif
