#include "sctp/chunks/CookieEchoChunk.h"

namespace sctp
{
	
size_t CookieEchoChunk::GetSize() const
{
	//Header + cookie
	return SizePad(4,4+cookie.GetSize());
}

size_t CookieEchoChunk::Serialize(BufferWritter& writter) const
{
	//Check header length
	if (!writter.Assert(4))
		return 0;
	
	//Get init pos
	size_t ini = writter.Mark();
	
	//Write header
	writter.Set1(type);
	writter.Set1(0);	// Always 0
	
	//Skip length position
	size_t mark = writter.Skip(2);
	
	//Check cookie size
	if (!writter.Assert(cookie.GetSize()))
		return 0;
	
	//Write cooke
	writter.Set(cookie);
	
	//Get length
	size_t length = writter.GetOffset(ini);
	//Set it
	writter.Set2(mark,length);
	
	//Pad
	return writter.PadTo(4);
}
	
Chunk::shared CookieEchoChunk::Parse(BufferReader& reader)
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
	if (type!=Type::COOKIE_ECHO)
		//Error
		return nullptr;
		
	//Create chunk
	auto cookieEcho = std::make_shared<CookieEchoChunk>();
	
	//Check size
	if (!reader.Assert(length-4)) 
		//Error
		return nullptr;
	
	//Get cookie
	cookieEcho->cookie = reader.GetBuffer(length-4);
	
	//Pad input
	if (!reader.PadTo(4))
		return nullptr;
	
	//Done
	return std::static_pointer_cast<Chunk>(cookieEcho);
}
	
};
