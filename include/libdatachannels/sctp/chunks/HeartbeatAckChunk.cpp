#include "sctp/chunks/HeartbeatAckChunk.h"

namespace sctp
{
	
size_t HeartbeatAckChunk::GetSize() const
{
	//Header + attributes
	size_t size = 20;
	
	//Done
	return size;
}

size_t HeartbeatAckChunk::Serialize(BufferWritter& writter) const
{
	//Check header length
	if (!writter.Assert(4))
		return 0;
	
	//Get init pos
	size_t ini = writter.Mark();
	
	//Write header
	writter.Set1(type);
	writter.Set1(flag);
	
	//Skip length position
	size_t mark = writter.Skip(2);
	
	//Heartbeat param
	{
		//Check parameter length
		size_t len = 4+senderSpecificHearbeatInfo.GetSize();
		if (!writter.Assert(len))
			return 0;
		//Write it
		writter.Set2(Parameter::HeartbeatInfo);
		writter.Set2(len);
		writter.Set(senderSpecificHearbeatInfo);
		writter.PadTo(4);
	}
	
	//Get length
	size_t length = writter.GetOffset(ini);
	//Set it
	writter.Set2(mark,length);
	
	//Done
	return length;
}
	
Chunk::shared HeartbeatAckChunk::Parse(BufferReader& reader)
{
	//Check size
	if (!reader.Assert(4)) 
		//Error
		return nullptr;
	
	//Get header
	size_t mark	= reader.Mark();
	uint8_t type	= reader.Get1();
	uint8_t flag	= reader.Get1(); //Ignored, should be 0
	uint16_t length	= reader.Get2();
	
	//Check type
	if (type!=Type::HEARTBEAT_ACK)
		//Error
		return nullptr;
		
	//Create chunk
	auto heartbeatRequest = std::make_shared<HeartbeatAckChunk>();
	
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
			case Parameter::HeartbeatInfo:
				heartbeatRequest->senderSpecificHearbeatInfo = paramReader.GetBuffer(paramReader.GetLeft());
				break;
			default:
				//Unkonwn
				return nullptr;
		}
		//Ensure all input has been consumed
		if (paramReader.GetLeft())
			return nullptr;
		//Do padding
		reader.PadTo(4);
	}
	
	//Done
	return std::static_pointer_cast<Chunk>(heartbeatRequest);
}
	
};
