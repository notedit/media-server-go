#ifndef DATACHANNEL_IMPL_DATACHANNEL_H_
#define DATACHANNEL_IMPL_DATACHANNEL_H_
#include "Datachannels.h"

#include "sctp/Stream.h"

namespace datachannels
{
namespace impl
{
	
class Datachannel : public datachannels::Datachannel
{
public:
	enum Payload 
	{
		WebRTCString	  = 51,
		WebRTCBinary	  = 53,
		WebRTCStringEmpty = 56,
		WebRTCBinaryEmpty = 57,
		
	};
	
public:
	Datachannel(const sctp::Stream::shared& stream);
	virtual ~Datachannel() = default;
	virtual bool Send(MessageType type, const uint8_t* data = nullptr, const uint64_t size = 0) override;
	virtual bool Close() override;
	
	// Event handlers
	virtual void OnMessage(const std::function<void(MessageType, const uint8_t*,uint64_t)>& callback) override
	{
		//Store callback
		onMessage = callback;
	}	
private:
	sctp::Stream::shared stream;
	std::function<void(MessageType, const uint8_t*,uint64_t)> onMessage;
};

}; //namespace impl
}; //namespace datachannel
#endif 
