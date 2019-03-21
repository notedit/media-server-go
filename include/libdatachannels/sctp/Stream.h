#ifndef SCTP_STREAM_H
#define SCTP_STREAM_H

#include "Datachannels.h"

#include <list>
#include <memory>

#include "Buffer.h"


namespace sctp
{

class Association;
	
class Stream
{
public:
	using shared = std::shared_ptr<Stream>;
public:
	Stream(Association &association, uint16_t id);
	virtual ~Stream();
	
	bool Recv(const uint8_t ppid, const uint8_t* buffer, const size_t size);
	bool Send(const uint8_t ppid, const uint8_t* buffer, const size_t size);
	
	uint16_t GetId() const { return id; }
	
	// Event handlers
	void OnMessage(std::function<void(uint8_t, const uint8_t*,uint64_t)> callback)
	{
		//Store callback
		onMessage = callback;
	}
private:
	uint16_t id;
	Association &association;
	std::list<std::pair<uint8_t,Buffer>> outgoingMessages;
	Buffer incomingMessage;
	
	std::function<void(uint8_t, const uint8_t*,uint64_t)> onMessage;
};

}; // namespace
#endif /* SCTP_STREAM_H */

