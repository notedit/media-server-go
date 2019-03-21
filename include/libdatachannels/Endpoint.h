#ifndef DATACHANNEL_IMPL_ENDPOINT_H_
#define DATACHANNEL_IMPL_ENDPOINT_H_
#include "Datachannels.h"
#include "sctp/Association.h"

namespace datachannels
{
namespace impl
{

class Endpoint : public datachannels::Endpoint
{
public:
	Endpoint(TimeService& timeService);
	virtual ~Endpoint();
	
	virtual bool Init(const Options& options)  override;
	virtual Datachannel::shared CreateDatachannel(const Datachannel::Options& options)  override;
	virtual bool Close()  override;
	
	// Getters
	virtual uint16_t GetLocalPort() const override;
	virtual uint16_t GetRemotePort() const override;
	virtual datachannels::Transport& GetTransport() override;
private:
	sctp::Association association;
};

}; //namespace impl
}; //namespace datachannels
#endif 
