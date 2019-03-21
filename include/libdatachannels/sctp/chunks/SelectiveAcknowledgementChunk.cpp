#include "sctp/chunks/SelectiveAcknowledgementChunk.h"

namespace sctp
{
	
size_t SelectiveAcknowledgementChunk::GetSize() const
{
	//Header + attributes
	size_t size = 16 + gapAckBlocks.size()*4 + duplicateTuplicateTrasnmissionSequenceNumbers.size()*4;
	
	//Done
	return size;
}

size_t SelectiveAcknowledgementChunk::Serialize(BufferWritter& writter) const
{
	//Check header length
	if (!writter.Assert(16))
		return 0;
	
	//Get init pos
	size_t ini = writter.Mark();
	
	//Write header
	writter.Set1(type);
	writter.Set1(0);
	//Skip length position
	size_t mark = writter.Skip(2);
	
	//Set attributes
	writter.Set4(cumulativeTrasnmissionSequenceNumberAck);
	writter.Set4(adveritsedReceiverWindowCredit);
	writter.Set2(gapAckBlocks.size());
	writter.Set2(duplicateTuplicateTrasnmissionSequenceNumbers.size());
	
	//For each gap
	for (const auto& gap : gapAckBlocks)
	{
		//Check header length
		if (!writter.Assert(4))
			return 0;
		///Write gap
		writter.Set2(gap.first);
		writter.Set2(gap.second);
	}
	
	//For each duplicated tsn
	for (const auto& duplicated : duplicateTuplicateTrasnmissionSequenceNumbers)
	{
		//Check header length
		if (!writter.Assert(4))
			return 0;
		///Write gap
		writter.Set4(duplicated);
	}
	
	//Get length
	size_t length = writter.GetOffset(ini);
	//Set it
	writter.Set2(mark,length);
	
	//Done
	return length;
}
	
Chunk::shared SelectiveAcknowledgementChunk::Parse(BufferReader& reader)
{
	//Check size
	if (!reader.Assert(16)) 
		//Error
		return nullptr;
	
	//Get header
	size_t mark	= reader.Mark();
	uint8_t type	= reader.Get1();
	uint8_t flag	= reader.Get1(); //Ignored, should be 0
	uint16_t length	= reader.Get2();
	
	//Check type
	if (type!=Type::SACK)
		//Error
		return nullptr;
		
	//Create chunk
	auto ack = std::make_shared<SelectiveAcknowledgementChunk>();
	
	//Read params
	ack->cumulativeTrasnmissionSequenceNumberAck	= reader.Get4();
	ack->adveritsedReceiverWindowCredit		= reader.Get4();
	const auto numGapAckBlocks			= reader.Get2();
	const auto numDuplicatedTSNs			= reader.Get2();
	
	//For each gap
	for (size_t i=0;i<numGapAckBlocks;++i)
	{
		//Check size
		if (!reader.Assert(4)) 
			//Error
			return nullptr;
		//Read gap
		ack->gapAckBlocks.push_back({
			reader.Get2(),
			reader.Get2()
		});
	}
	
	//For each duplicated tsn
	for (size_t i=0;i<numDuplicatedTSNs;++i)
	{
		//Check size
		if (!reader.Assert(4)) 
			//Error
			return nullptr;
		//Read gap
		ack->duplicateTuplicateTrasnmissionSequenceNumbers.push_back(reader.Get4());
	}
	
	//Check size
	if (!reader.Assert(length-16)) 
		//Error
		return nullptr;
		
	//Done
	return std::static_pointer_cast<Chunk>(ack);
}
	
};
