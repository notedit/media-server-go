#include "Datachannel.h"

namespace datachannels
{
namespace impl
{

Datachannel::Datachannel(const sctp::Stream::shared& stream)
{
	this->stream = stream;
	this->stream->OnMessage([&](const uint8_t ppid, const uint8_t* buffer, const size_t size){
		
	});
}
	
bool Datachannel::Send(MessageType type, const uint8_t* data, const uint64_t size)
{
	uint8_t empty = 0;
	
	if (!data || !size)
		//   SCTP does not support the sending of empty user messages.  Therefore,
		//   if an empty message has to be sent, the appropriate PPID (WebRTC
		//   String Empty or WebRTC Binary Empty) is used and the SCTP user
		//   message of one zero byte is sent.  When receiving an SCTP user
		//   message with one of these PPIDs, the receiver MUST ignore the SCTP
		//   user message and process it as an empty message.		
		return stream->Send(type==UTF8 ? WebRTCStringEmpty : WebRTCBinaryEmpty, &empty, 1);
	else 
		return stream->Send(type==UTF8 ? WebRTCString : WebRTCBinary, data, size);
}

bool Datachannel::Close()
{
	//   Closing of a data channel MUST be signaled by resetting the
	//   corresponding outgoing streams [RFC6525].  This means that if one
	//   side decides to close the data channel, it resets the corresponding
	//   outgoing stream.  When the peer sees that an incoming stream was
	//   reset, it also resets its corresponding outgoing stream.  Once this
	//   is completed, the data channel is closed.  Resetting a stream sets
	//   the Stream Sequence Numbers (SSNs) of the stream back to 'zero' with
	//   a corresponding notification to the application layer that the reset
	//   has been performed.  Streams are available for reuse after a reset
	//   has been performed.
	//
	//   [RFC6525] also guarantees that all the messages are delivered (or
	//   abandoned) before the stream is reset.	
	return true;
}

}; //namespace impl
}; //namespace datachannel
