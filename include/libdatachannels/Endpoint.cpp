#include "Endpoint.h"
#include <memory>

namespace datachannels
{
	
Endpoint::shared Endpoint::Create(TimeService& timeService) 
{
	//Create endpoint
	auto endpoint = std::make_shared<datachannels::impl::Endpoint>(timeService);
	//Cast and return
	return std::static_pointer_cast<Endpoint>(endpoint);
}

namespace impl
{

Endpoint::Endpoint(datachannels::TimeService& timeService) :
	association(timeService)
{
	
}

Endpoint::~Endpoint()
{
	//Terminate association now!
	association.Abort();
}

bool Endpoint::Init(const Options& options)
{
	//Set ports on sctp
	association.SetLocalPort(options.localPort);
	association.SetRemotePort(options.remotePort);
	
	//If we are clients
	if (options.setup==Setup::Client)
		//Start association
		return association.Associate();
	
	//OK, wait for client to associate
	return true;
}

Datachannel::shared Endpoint::CreateDatachannel(const Datachannel::Options& options)
{
	return nullptr;
}

bool Endpoint::Close()
{
	//Gracefuly stop association
	return association.Shutdown();
}

// Getters
uint16_t Endpoint::GetLocalPort() const
{
	return association.GetLocalPort();
}

uint16_t Endpoint::GetRemotePort() const
{
	return association.GetRemotePort();
}

datachannels::Transport& Endpoint::GetTransport()
{
	return association;
}
	
}; // namespace impl
}; // namespace datachannel
