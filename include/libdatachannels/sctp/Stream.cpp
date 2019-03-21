#include "sctp/Stream.h"

namespace sctp
{

Stream::Stream(Association &association, uint16_t id) :
	association(association)
{
	this->id = id;
}

Stream::~Stream()
{
}

bool Stream::Recv(const uint8_t ppid, const uint8_t* buffer, const size_t size) //first,last?
{
	//onMessage(ppid,buff)
	return true;
}

bool Stream::Send(const uint8_t ppid, const uint8_t* buffer, const size_t size)
{
	//TODO: check max queue size?
	
	//Add new message to ougogin queue
	outgoingMessages.push_back(std::make_pair<>(ppid,Buffer(buffer,size)));
	
	//TODO: signal pending data
	
	//done
	return true;
}

}; // namespace sctp
