#include "sctp/chunks/UnknownChunk.h"

namespace sctp
{
	
size_t UnknownChunk::GetSize() const
{
	//Header + buffer
	return SizePad(4, 4+buffer.GetSize());
}

size_t UnknownChunk::Serialize(BufferWritter& writter) const
{
	//Get init pos
	size_t ini = writter.Mark();
	
	//Check header length
	if (!writter.Assert(4))
		return 0;
	
	//Write header
	writter.Set1(type);
	writter.Set1(flag);
	
	//Skip length position
	size_t mark = writter.Skip(2);
	
	//Check buffer size
	if (!writter.Assert(buffer.GetSize()))
		return 0;
	
	//Write buffer
	writter.Set(buffer);
	
	//Get length
	size_t length = writter.GetOffset(ini);
	//Set it
	writter.Set2(mark,length);
	
	//Pad
	return writter.PadTo(4);
}
	
Chunk::shared UnknownChunk::Parse(BufferReader& reader)
{
	//Check size
	if (!reader.Assert(4)) 
		//Error
		return nullptr;
	
	//Get header
	uint8_t type	= reader.Get1();
	uint8_t flag	= reader.Get1(); //Ignored, should be 0
	uint16_t length	= reader.Get2();
	
	//Check size
	if (length<4 || !reader.Assert(length-4)) 
		//Error
		return nullptr;
	
	//Create chunk
	auto unknown = std::make_shared<UnknownChunk>(type);
	
	//Set attributes
	unknown->flag = flag;
	unknown->buffer = reader.GetBuffer(length-4);
		
	//Done
	return std::static_pointer_cast<Chunk>(unknown);
}
	
};
