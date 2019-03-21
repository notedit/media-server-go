#ifndef SCTP_ERROR_CAUSE_H_
#define SCTP_ERROR_CAUSE_H_
#include <string>

#include "Buffer.h"

namespace sctp
{
	
class ErrorCause
{
public:
	enum Cause
	{
		InvalidStreamIdentifier			= 1,
		MissingMandatoryParameter		= 2,
		StaleCookieError			= 3,
		OutofResource				= 4,
		UnresolvableAddress			= 5,
		UnrecognizedChunkType			= 6,
		InvalidMandatoryParameter		= 7,
		UnrecognizedParameters			= 8,
		NoUserData				= 9,
		CookieReceivedWhileShuttingDown		= 10,
		RestartofanAssociationwithNewAddresses	= 11,
		UserInitiatedAbort			= 12,
		ProtocolViolation			= 13,
	};
public:
	uint16_t code;
	Buffer info;
};

}; // namespace sctp

#endif
