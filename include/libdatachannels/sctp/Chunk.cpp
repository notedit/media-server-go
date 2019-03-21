#include "sctp/Chunk.h"

namespace sctp
{
	
Chunk::shared Chunk::Parse(BufferReader& reader)
{
	//Ensure we have at laast the header
	if (!reader.Assert(4))
		//Error
		return nullptr;
	
	// Peek type
	switch((Type)reader.Peek1())
	{
		case Type::PDATA:
			return PayloadDataChunk::Parse(reader);
		case Type::INIT:
			return InitiationChunk::Parse(reader);
		case Type::INIT_ACK:
			return InitiationAcknowledgementChunk::Parse(reader);
		case Type::SACK:
			return SelectiveAcknowledgementChunk::Parse(reader);
		case Type::HEARTBEAT:
			return HeartbeatRequestChunk::Parse(reader);
		case Type::HEARTBEAT_ACK:
			return HeartbeatAckChunk::Parse(reader);
		case Type::ABORT:
			return AbortAssociationChunk::Parse(reader);
		case Type::SHUTDOWN:
			return ShutdownAssociationChunk::Parse(reader);
		case Type::SHUTDOWN_ACK:
			return ShutdownAcknowledgementChunk::Parse(reader);
		case Type::ERROR:
			return OperationErrorChunk::Parse(reader);
		case Type::COOKIE_ECHO:
			return CookieEchoChunk::Parse(reader);
		case Type::COOKIE_ACK:
			return CookieAckChunk::Parse(reader);
		case Type::ECNE:
			return UnknownChunk::Parse(reader);
		case Type::CWR:
			return UnknownChunk::Parse(reader);
		case Type::SHUTDOWN_COMPLETE:
			return ShutdownCompleteChunk::Parse(reader);
		case Type::RE_CONFIG:
			return ReConfigChunk::Parse(reader);
		case Type::FORWARD_CUMULATIVE_TSN:
			return ForwardCumulativeTSNChunk::Parse(reader);
	}
	
	return UnknownChunk::Parse(reader);
}

};
