#include "sctp/chunks/PayloadDataChunk.h"

namespace sctp
{
	
size_t PayloadDataChunk::GetSize() const
{
	//Header + attributes + user data
	return SizePad(4,16+userData.GetSize());
}

size_t PayloadDataChunk::Serialize(BufferWritter& writter) const
{
	//Check header length
	if (!writter.Assert(16))
		return 0;
	
	//Get init pos
	size_t ini = writter.Mark();
	
	//Creage flag
	uint8_t flag = (unordered & 0x04) | (beginingFragment & 0x02) | (endingFragment & 0x01);
	
	//Write header
	writter.Set1(type);
	writter.Set1(flag);
	//Skip length position
	size_t mark = writter.Skip(2);
	
	//Set attributes
	writter.Set4(transmissionSequenceNumber);
	writter.Set2(streamIdentifier);
	writter.Set2(streamSequenceNumber);
	writter.Set4(payloadProtocolIdentifier);

	//Check user data size
	if (!writter.Assert(userData.GetSize()))
		return 0;
	
	//Write cooke
	writter.Set(userData);
	
	///Get length
	size_t length = writter.GetOffset(ini);
	//Set it
	writter.Set2(mark,length);
	
	//Pad
	return writter.PadTo(4);
}
	
Chunk::shared PayloadDataChunk::Parse(BufferReader& reader)
{
	//Check size
	if (!reader.Assert(16)) 
		//Error
		return nullptr;
	
	//Get header
	size_t mark	= reader.Mark();
	uint8_t type	= reader.Get1();
	uint8_t flag	= reader.Get1(); 
	uint16_t length	= reader.Get2();
	
	//Check type
	if (type!=Type::PDATA)
		//Error
		return nullptr;
		
	//Create chunk
	auto data = std::make_shared<PayloadDataChunk>();
	
	//Set flag bits
	data->unordered		= flag & 0x04;
	data->beginingFragment	= flag & 0x02;
	data->beginingFragment	= flag & 0x01;
	
	//Read params
	data->transmissionSequenceNumber = reader.Get4();
	data->streamIdentifier		 = reader.Get2();
	data->streamSequenceNumber	 = reader.Get2();
	data->payloadProtocolIdentifier	 = reader.Get4();
	
	//Check size
	if (!reader.Assert(length-16)) 
		//Error
		return nullptr;
	
	//Get user data
	data->userData = reader.GetBuffer(length-16);
	
	//Pad input
	if (!reader.PadTo(4))
		return nullptr;
	
	//Done
	return std::static_pointer_cast<Chunk>(data);
}
	
};
