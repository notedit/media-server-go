#include "sctp/chunks/InitiationChunk.h"

namespace sctp
{
	
size_t InitiationChunk::GetSize() const
{
	//Header + attributes
	size_t size = 20;
	
	//Set parameters
	size += ipV4Addresses.size() * 8;
	size += ipV6Addresses.size() * 20;
	if (suggestedCookieLifeSpanIncrement)
		size += 8;
	if (hostName)
		size += SizePad(hostName->length(), 4) ;
	size += SizePad(supportedAddressTypes.size() * 2, 4);
	size += SizePad(supportedExtensions.size(), 4);
	for (const auto& unknownParameter : unknownParameters)
		size += SizePad(unknownParameter.second.GetSize(), 4);
	
	//Done
	return size;
}

size_t InitiationChunk::Serialize(BufferWritter& writter) const
{
	//Check header length
	if (!writter.Assert(20))
		return 0;
	
	//Get init pos
	size_t ini = writter.Mark();
	
	//Write header
	writter.Set1(type);
	writter.Set1(flag);
	//Skip length position
	size_t mark = writter.Skip(2);
	
	//Set attributes
	writter.Set4(initiateTag);
	writter.Set4(advertisedReceiverWindowCredit);
	writter.Set2(numberOfOutboundStreams);
	writter.Set2(numberOfInboundStreams);
	writter.Set4(initialTransmissionSequenceNumber);
	
	//IPV4 addresses
	for (const auto& ipV4Address : ipV4Addresses)
	{
		//Check parameter length
		size_t len = 12;
		if (!writter.Assert(len))
			return 0;
		//Write it
		writter.Set2(Parameter::IPv4Address);
		writter.Set2(len);
		writter.Set<8>(ipV4Address);
		//Pad input
		if (!writter.PadTo(4))
			return 0;
	}
	
	//IPV6 addresses
	for (const auto& ipV6Address : ipV6Addresses)
	{
		//Check parameter length
		size_t len = 24;
		if (!writter.Assert(len))
			return 0;
		//Write it
		writter.Set2(Parameter::IPv6Address);
		writter.Set2(len);
		writter.Set<20>(ipV6Address);
		//Pad input
		if (!writter.PadTo(4))
			return 0;
	}
	
	//Cookie life span
	if (suggestedCookieLifeSpanIncrement)
	{
		//Check parameter length
		size_t len = 12;
		if (!writter.Assert(len))
			return 0;
		//Write it
		writter.Set2(Parameter::CookiePreservative);
		writter.Set2(len);
		writter.Set8(*suggestedCookieLifeSpanIncrement);
		//Pad input
		if (!writter.PadTo(4))
			return 0;
	}
	
	//Optional Host name
	if (hostName)
	{
		//Check parameter length
		size_t len = 4+hostName->length();
		if (!writter.Assert(len))
			return 0;
		//Write it
		writter.Set2(Parameter::HostNameAddress);
		writter.Set2(len);
		writter.Set(*hostName);
		//Pad input
		if (!writter.PadTo(4))
			return 0;
	}
	
	//Supported extensions parameter
	if (supportedAddressTypes.size())
	{
		//Check parameter length
		size_t len = 4+supportedAddressTypes.size()*2;
		if (!writter.Assert(len))
			return 0;
		//Write it
		writter.Set2(Parameter::SupportedAddressTypes);
		writter.Set2(len);
		for (const auto& supportedAddressType : supportedAddressTypes)
			writter.Set2(supportedAddressType);
		//Pad input
		if (!writter.PadTo(4))
			return 0;
	}
	
	//Supported extensions parameter
	if (supportedExtensions.size())
	{
		//Check parameter length
		size_t len = 4+supportedExtensions.size();
		if (!writter.Assert(len))
			return 0;
		//Write it
		writter.Set2(Parameter::SupportedExtensions);
		writter.Set2(len);
		for (const auto& supportedExtension : supportedExtensions)
			writter.Set1(supportedExtension);
		//Pad input
		if (!writter.PadTo(4))
			return 0;
	}
	
	//Unknown parameters
	for (const auto& unknownParameter : unknownParameters)
	{
		//Check parameter length
		size_t len = 4+unknownParameter.second.GetSize();
		if (!writter.Assert(len))
			return 0;
		//Write it
		writter.Set2(unknownParameter.first);
		writter.Set2(len);
		writter.Set(unknownParameter.second);
		//Pad input
		if (!writter.PadTo(4))
			return 0;
	}
	
	//Support for ForwardTSN
	if (forwardTSNSupported)
	{
		//Check parameter length
		size_t len = 4;;
		if (!writter.Assert(len))
			return 0;
		//Write it
		writter.Set2(Parameter::ForwardTSNSupported);
		writter.Set2(len);
		//Pad input
		if (!writter.PadTo(4))
			return 0;
	}
	
	//Get length
	size_t length = writter.GetOffset(ini);
	//Set it
	writter.Set2(mark,length);
	
	//Done
	return length;
}
	
Chunk::shared InitiationChunk::Parse(BufferReader& reader)
{
	//Check size
	if (!reader.Assert(20)) 
		//Error
		return nullptr;
	
	//Get header
	size_t mark	= reader.Mark();
	uint8_t type	= reader.Get1();
	uint8_t flag	= reader.Get1(); //Ignored, should be 0
	uint16_t length	= reader.Get2();
	
	//Check type
	if (type!=Type::INIT)
		//Error
		return nullptr;
		
	//Create chunk
	auto init = std::make_shared<InitiationChunk>();
		
	//Set attributes
	init->initiateTag			= reader.Get4();
	init->advertisedReceiverWindowCredit	= reader.Get4();
	init->numberOfOutboundStreams		= reader.Get2();
	init->numberOfInboundStreams		= reader.Get2();
	init->initialTransmissionSequenceNumber = reader.Get4();
	init->forwardTSNSupported		= false;
	
	//Read parameters
	while (reader.GetLeft()>=4)
	{
		//Get parameter type
		uint16_t paramType = reader.Get2();
		uint16_t paramLength = reader.Get2();
		//Ensure lenghth is correct as it has to contain the type and length itself
		if (paramLength<4)
			return nullptr;
		//Remove header
		paramLength-=4;
		//Ensure we have enought length
		if (!reader.Assert(paramLength)) return nullptr;
		//Get reader for the param length
		BufferReader paramReader = reader.GetReader(paramLength);
		//Depending on the parameter type
		switch(paramType)
		{
			case Parameter::IPv4Address:
				if (!paramReader.Assert(8)) return nullptr;
				init->ipV4Addresses.push_back(paramReader.Get<8>());
				break;
			case Parameter::IPv6Address:
				if (!paramReader.Assert(20)) return nullptr;
				init->ipV6Addresses.push_back(paramReader.Get<20>());
				break;
			case Parameter::HostNameAddress:
				init->hostName = paramReader.GetString(paramReader.GetLeft());
				break;
			case Parameter::SupportedAddressTypes:
				while (paramReader.GetLeft())
					init->supportedAddressTypes.push_back(paramReader.Get2());
				break;
			case Parameter::SupportedExtensions:
				while (paramReader.GetLeft())
					init->supportedExtensions.push_back(paramReader.Get1());
				break;
			case Parameter::CookiePreservative:
				if (!paramReader.Assert(8)) return nullptr;
				init->suggestedCookieLifeSpanIncrement = paramReader.Get8();
				break;
			case Parameter::ForwardTSNSupported:
				init->forwardTSNSupported = true;
				break;
			case Parameter::Padding:
				//The PAD parameter MAY be included only in the INIT chunk.  It MUST
				//NOT be included in any other chunk.  The receiver of the PAD
				//parameter MUST silently discard this parameter and continue
				//processing the rest of the INIT chunk.
				reader.Skip(reader.GetLeft());
			default:
				//Unkonwn
				init->unknownParameters.push_back(std::make_pair<uint8_t,Buffer>(paramType,paramReader.GetBuffer(paramReader.GetLeft())));
		}
		//Ensure all input has been consumed
		if (paramReader.GetLeft())
			//Error
			return nullptr;
		//Do padding
		reader.PadTo(4);
	}
	
	//Done
	return std::static_pointer_cast<Chunk>(init);
}
	
};
