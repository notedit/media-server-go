#include "sctp/Association.h"
#include "sctp/Chunk.h"

#include <chrono>
#include <random>
#include <crc32c/crc32c.h>
#include <condition_variable>

namespace sctp
{
	
//Random stuff
std::random_device rd;
std::mt19937 gen{rd()};
std::uniform_int_distribution<unsigned long> dis{1, 4294967295};

Association::Association(datachannels::TimeService& timeService) :
	timeService(timeService)
{
}

void Association::SetState(State state)
{
	this->state = state;
}

bool Association::Associate()
{
	//Check state
	if (state!=State::Closed)
		//Error
		return false;
	
	//	"A" first sends an INIT chunk to "Z".  In the INIT, "A" must
	//	provide its Verification Tag (Tag_A) in the Initiate Tag field.
	//	Tag_A SHOULD be a random number in the range of 1 to 4294967295
	//	(see Section 5.3.1 for Tag value selection).  After sending the
	//	INIT, "A" starts the T1-init timer and enters the COOKIE-WAIT
	//	state.
	
	//Create new verification tag
	localVerificationTag = dis(gen);
	
	//Reset init retransmissions
	initRetransmissions = 0;
	
	//Enqueue new INIT chunk
	auto init = std::make_shared<InitiationChunk>();
	
	//Set params
	init->initiateTag			= localVerificationTag;
	init->advertisedReceiverWindowCredit	= 0;
	init->numberOfOutboundStreams		= 0xFFFF;
	init->numberOfInboundStreams		= 0xFFFF;
	init->initialTransmissionSequenceNumber = 0;
	
	// draft-ietf-rtcweb-data-channel-13
	//	The INIT and INIT-ACK chunk MUST NOT contain any IPv4 Address or
	//	IPv6 Address parameters.  The INIT chunk MUST NOT contain the
 	//	Supported Address Types parameter.
	
	// draft-ietf-rtcweb-data-channel-13#page-7
	//	The dynamic address reconfiguration extension defined in [RFC5061]
	//	MUST be used to signal the support of the stream reset extension
	//	defined in [RFC6525].  Other features of [RFC5061] are OPTIONAL.
	init->supportedExtensions.push_back(Chunk::Type::RE_CONFIG);
		
	//Set timer
	initTimer = timeService.CreateTimer(InitRetransmitTimeout,[=](...){
		//Retransmit init chunk
		if (initRetransmissions++<MaxInitRetransmits)
		{
			//Enquee
			Enqueue(std::static_pointer_cast<Chunk>(init));
			//Retry again
			initTimer->Again(InitRetransmitTimeout);
		} else {
			//Close
			SetState(State::Closed);
		}
	});
	
	//Change state
	SetState(State::CookieWait);
	
	//Enquee
	Enqueue(std::static_pointer_cast<Chunk>(init));
	
	//Done
	return true;
}

bool Association::Shutdown()
{
	return true;
}

bool Association::Abort()
{
	return true;
}

size_t Association::WritePacket(uint8_t *data, uint32_t size)
{
	//Create reader
	BufferReader reader(data,size);
	
	//TODO: Check crc 
	
	
	
	//Parse packet header
	auto header = PacketHeader::Parse(reader);

	//Ensure it was correctly parsed
	if (!header)
		//Error
		return false;

	//Check correct local and remote port
	if (header->sourcePortNumber!=remotePort || header->destinationPortNumber!=localPort || header->verificationTag!=localVerificationTag)
		//Error
		return false;
	
	//Read chunks
	while (reader.GetLeft()>4)
	{
		//Parse chunk
		auto chunk = Chunk::Parse(reader);
		//Check 
		if (!chunk)
			//Skip
			continue;
		//Process it
		Process(chunk);
	}
	
	//If we need to acknowledge
	if (pendingAcknowledge)
	{
		//If we have to do it now
		if (pendingAcknowledgeTimeout == 0ms)
			//Acknoledge now
			Acknowledge();
		//rfc4960#page-78
		//	Specifically, an
		//	acknowledgement SHOULD be generated for at least every second packet
		//	(not every second DATA chunk) received, and SHOULD be generated
		//	within 200 ms of the arrival of any unacknowledged DATA chunk.

		//If there was already a timeout
		else if (sackTimer.get())
			//We should do sack now
			Acknowledge();
		else
			//Schedule timer
			sackTimer = timeService.CreateTimer(pendingAcknowledgeTimeout,[this](...){
				//In the future
				Acknowledge();
			});
	}
		
	//Done
	return true;
}

size_t Association::ReadPacket(uint8_t *data, uint32_t size)
{
	//Check there is pending data
	if (!pendingData)
		//Nothing to do
		return 0;
	
	//Create buffer writter
	BufferWritter writter(data,size);
	
	//Create new packet header
	PacketHeader header(localPort,remotePort,remoteVerificationTag);
	
	//Serialize it
	if (!header.Serialize(writter))
		//Error
		return 0;

	size_t num = 0;
	
	//Fill chunks from control queue first
	for (auto it=queue.begin();it!=queue.end();)
	{
		//Get chunk
		auto chunk = *it;

		//Ensure we have enought space for chunk
		if (writter.GetLeft()<chunk->GetSize())
			//We cant send more on this packet
			break;
		
		//Check if it must be sent alone
		if (chunk->type==Chunk::Type::INIT || chunk->type==Chunk::Type::INIT_ACK || chunk->type==Chunk::Type::COOKIE_ECHO)
		{
			//IF it is not firest
			if (num)
				//Flush all before this
				break;
		}
		
		//Remove from queue and move to next chunk
		it = queue.erase(it);
		
		//Serialize chunk
		chunk->Serialize(writter);
		
		//Check if it must be sent alone
		if (chunk->type==Chunk::Type::INIT || chunk->type==Chunk::Type::INIT_ACK || chunk->type==Chunk::Type::COOKIE_ECHO)
			//Send alone
			break;
	}

	//TODO:Check in which stata data can be sent
		//TODO:Now fill data chunks from streams

	//Get length
	size_t length = writter.GetLength();
	//Calculate crc
	header.checksum  = crc32c::Crc32c(data,length);
	//Go to the begining
	writter.GoTo(0);
	
	//Serialize it now with checksum
	header.Serialize(writter);
	
	//Check if there is more data to send
	if (!queue.size())
		//No
		pendingData = false;
	//Done
	return length;
}

void Association::Process(const Chunk::shared& chunk)
{
	//Depending onthe state
	switch (state)
	{
		case State::Closed:
		{
			switch(chunk->type)
			{
				case Chunk::Type::INIT:
				{
					//Get chunk of correct type
					auto init = std::static_pointer_cast<InitiationChunk>(chunk);
					
					//rfc4960#page-55
					//	"Z" shall respond immediately with an INIT ACK chunk.  The
					//	destination IP address of the INIT ACK MUST be set to the source
					//	IP address of the INIT to which this INIT ACK is responding.  In
					//	the response, besides filling in other parameters, "Z" must set
					//	the Verification Tag field to Tag_A, and also provide its own
					//	Verification Tag (Tag_Z) in the Initiate Tag field.
					//
					
					//Get remote verification tag
					remoteVerificationTag = init->initiateTag;
						
					//Create new verification tag
					localVerificationTag = dis(gen);

					//Reset init retransmissions
					initRetransmissions = 0;

					//Enqueue new INIT chunk
					auto initAck = std::make_shared<InitiationAcknowledgementChunk>();

					//Set params
					initAck->initiateTag			= localVerificationTag;
					initAck->advertisedReceiverWindowCredit	= localAdvertisedReceiverWindowCredit;
					initAck->numberOfOutboundStreams	= 0xFFFF;
					initAck->numberOfInboundStreams		= 0xFFFF;
					initAck->initialTransmissionSequenceNumber = 0;

					// draft-ietf-rtcweb-data-channel-13
					//	The INIT and INIT-ACK chunk MUST NOT contain any IPv4 Address or
					//	IPv6 Address parameters.  The INIT chunk MUST NOT contain the
					//	Supported Address Types parameter.

					// draft-ietf-rtcweb-data-channel-13#page-7
					//	The dynamic address reconfiguration extension defined in [RFC5061]
					//	MUST be used to signal the support of the stream reset extension
					//	defined in [RFC6525].  Other features of [RFC5061] are OPTIONAL.
					initAck->supportedExtensions.push_back(Chunk::Type::RE_CONFIG);
					
					//rfc4960#page-55
					//	Moreover, "Z" MUST generate and send along with the INIT ACK a
					//	State Cookie.  See Section 5.1.3 for State Cookie generation.
					
					// AAA is ensured by the DTLS & ICE layer, so createing a complex cookie is unnecessary IMHO
					initAck->stateCookie.SetData((uint8_t*)"dtls",strlen("dtls"));

					//Send back unkown parameters
					for (const auto& unknown : init->unknownParameters)
						//Copy as unrecognized
						initAck->unrecognizedParameters.push_back(unknown.second.Clone());
					
					///Enquee
					Enqueue(std::static_pointer_cast<Chunk>(initAck));
					
					break;
				}
				case Chunk::Type::COOKIE_ECHO:
				{
					//rfc4960#page-55
					//	D) Upon reception of the COOKIE ECHO chunk, endpoint "Z" will reply
					//	with a COOKIE ACK chunk after building a TCB and moving to the
					//	ESTABLISHED state.  A COOKIE ACK chunk may be bundled with any
					//	pending DATA chunks (and/or SACK chunks), but the COOKIE ACK chunk
					//	MUST be the first chunk in the packet.
					
					//Enqueue new INIT chunk
					auto cookieAck = std::make_shared<CookieAckChunk>();
					
					//We don't check the cookie for seame reasons as we don't create one
					
					///Enquee
					Enqueue(std::static_pointer_cast<Chunk>(cookieAck));
					
					//Change state
					SetState(State::Established);
				}
				default:
					//Error
					break;
			}
			break;
		}
		case State::CookieWait:
		{
			switch(chunk->type)
			{
				case Chunk::Type::INIT_ACK:
				{
					//Get chunk of correct type
					auto initAck = std::static_pointer_cast<InitiationAcknowledgementChunk>(chunk);
					
					//	C) Upon reception of the INIT ACK from "Z", "A" shall stop the T1-
					//	init timer and leave the COOKIE-WAIT state.  "A" shall then send
					//	the State Cookie received in the INIT ACK chunk in a COOKIE ECHO
					//	chunk, start the T1-cookie timer, and enter the COOKIE-ECHOED
					//	state.
					//
					
					// Stop timer
					initTimer->Cancel();
					
					//Enqueue new INIT chunk
					auto cookieEcho = std::make_shared<CookieEchoChunk>();
					
					//Copy cookie
					cookieEcho->cookie.SetData(initAck->stateCookie);
					
					//Reset cookie retransmissions
					initRetransmissions = 0;
					
					//Set timer
					cookieEchoTimer = timeService.CreateTimer(100ms,[=](...){
						//3)  If the T1-cookie timer expires, the endpoint MUST retransmit
						//    COOKIE ECHO and restart the T1-cookie timer without changing
						//    state.  This MUST be repeated up to 'Max.Init.Retransmits' times.
						//    After that, the endpoint MUST abort the initialization process
						//    and report the error to the SCTP user.
						if (initRetransmissions++<MaxInitRetransmits)
						{
							//Retransmit
							Enqueue(std::static_pointer_cast<Chunk>(cookieEcho));
							//Retry again
							initTimer->Again(100ms);
						} else {
							//Close
							SetState(State::Closed);
						}
					});
					
					///Enquee
					Enqueue(std::static_pointer_cast<Chunk>(cookieEcho));
					
					//Set new state
					SetState(State::CookieEchoed);
				}
				default:
					//Error
					break;
			}
			break;
		}
		case State::CookieEchoed:
		{
			switch(chunk->type)
			{
				case Chunk::Type::COOKIE_ACK:
				{
					//	E) Upon reception of the COOKIE ACK, endpoint "A" will move from the
					//	COOKIE-ECHOED state to the ESTABLISHED state, stopping the T1-
					//	cookie timer.
					
					// Stop timer
					cookieEchoTimer->Cancel();
					
					//Change state
					SetState(State::Established);
				}
				default:
					//Error
					break;
			}
			break;
		}
		case State::Established:
		{
			switch(chunk->type)
			{
				case Chunk::Type::PDATA:
				{
					//Get chunk of correct type
					auto pdata = std::static_pointer_cast<PayloadDataChunk>(chunk);
					
					//	After the reception of the first DATA chunk in an association the
					//	endpoint MUST immediately respond with a SACK to acknowledge the DATA
					//	chunk.  Subsequent acknowledgements should be done as described in
					bool first = lastReceivedTransmissionSequenceNumber == MaxTransmissionSequenceNumber;
					
					//Get tsn
					auto tsn = receivedTransmissionSequenceNumberWrapper.Wrap(pdata->transmissionSequenceNumber);
					
					//Storea tsn, if the container has elements with equivalent key, inserts at the upper bound of that range 
					auto it = receivedTransmissionSequenceNumbers.insert(tsn);
					
					//	When a packet arrives with duplicate DATA chunk(s) and with no new
					//	DATA chunk(s), the endpoint MUST immediately send a SACK with no
					//	delay.  If a packet arrives with duplicate DATA chunk(s) bundled with
					//	new DATA chunks, the endpoint MAY immediately send a SACK.
					bool duplicated = it != receivedTransmissionSequenceNumbers.begin() && *(--it)!=tsn;
					
					//rfc4960#page-89
					//	Upon the reception of a new DATA chunk, an endpoint shall examine the
					//	continuity of the TSNs received.  If the endpoint detects a gap in
					//	the received DATA chunk sequence, it SHOULD send a SACK with Gap Ack
					//	Blocks immediately.  The data receiver continues sending a SACK after
					//	receipt of each SCTP packet that doesn't fill the gap.
					bool hasGaps = false;
					
					//Iterate the received transmission numbers
					uint64_t prev = MaxTransmissionSequenceNumber;
					for (auto curr : receivedTransmissionSequenceNumbers)
					{
						//Check if not first or if ther is a jump bigger than 1 seq num
						if (prev!=MaxTransmissionSequenceNumber && curr>prev+1)
						{
							//It has a gap at least
							hasGaps = true;
							break;
						}
						//Next
						prev = curr;
					}

					//rfc4960#page-78
					//	When the receiver's advertised window is 0, the receiver MUST drop
					//	any new incoming DATA chunk with a TSN larger than the largest TSN
					//	received so far.  If the new incoming DATA chunk holds a TSN value
					//	less than the largest TSN received so far, then the receiver SHOULD
					//	drop the largest TSN held for reordering and accept the new incoming
					//	DATA chunk.  In either case, if such a DATA chunk is dropped, the
					//	receiver MUST immediately send back a SACK with the current receive
					//	window showing only DATA chunks received and accepted so far.  The
					//	dropped DATA chunk(s) MUST NOT be included in the SACK, as they were
					//	not accepted. 
					
					//We need to acknoledfe
					pendingAcknowledge = true;
					
					//If we need to send it now
					if (first || hasGaps || duplicated || sackTimer)
						//Acknoledge now
						pendingAcknowledgeTimeout = 0ms; 
					else 
						//Create timer
						pendingAcknowledgeTimeout = SackTimeout; 
					break;
				}
				case Chunk::Type::SACK:
				{
					break;
				}
			}
			break;
		}
		case State::ShutdownPending:
		{
			break;
		}
		case State::ShutDownSent:
		{
			break;
		}
		case State::ShutDownReceived:
		{
			break;
		}
		case State::ShutDown:
		{
			break;
		}
		case State::ShutDownAckSent:
		{
			break;
		}
	}
}


void Association::Acknowledge()
{
	//New sack message
	auto sack = std::make_shared<SelectiveAcknowledgementChunk>();
	
	//rfc4960#page-34
	//	By definition, the value of the Cumulative TSN Ack parameter is the
	//	last TSN received before a break in the sequence of received TSNs
	//	occurs; the next TSN value following this one has not yet been
	//	received at the endpoint sending the SACK.  This parameter therefore
	//	acknowledges receipt of all TSNs less than or equal to its value.

	//rfc4960#page-34
	//	The SACK also contains zero or more Gap Ack Blocks.  Each Gap Ack
	//	Block acknowledges a subsequence of TSNs received following a break
	//	in the sequence of received TSNs.  By definition, all TSNs
	//	acknowledged by Gap Ack Blocks are greater than the value of the
	//	Cumulative TSN Ack.
	
	//rfc4960#page-78
	//	Any received DATA chunks with TSN
	//	greater than the value in the Cumulative TSN Ack field are reported
	//	in the Gap Ack Block fields.  The SCTP endpoint MUST report as many
	//	Gap Ack Blocks as can fit in a single SACK chunk limited by the
	//	current path MTU.
	
	//Iterate the received transmission numbers
	uint64_t prev = MaxTransmissionSequenceNumber;
	uint64_t gap  = MaxTransmissionSequenceNumber;
	for (auto it = receivedTransmissionSequenceNumbers.begin();it!=receivedTransmissionSequenceNumbers.end();)
	{
		//Get current
		uint64_t current = *it;
		//If we have a continous sequence number or it is the first one
		if (lastReceivedTransmissionSequenceNumber==MaxTransmissionSequenceNumber || current==lastReceivedTransmissionSequenceNumber+1)
		{
			//Update last received tsn
			lastReceivedTransmissionSequenceNumber = current;
		}
		//Check if is duplicated
		else if (prev==current)
		{
			//It is duplicated
			sack->duplicateTuplicateTrasnmissionSequenceNumbers.push_back(current);
		}
		//If it is a gap
		else if (prev!=MaxTransmissionSequenceNumber && current>prev+1)
		{
			//If we had a gap start
			if (gap)
				//Add block ending at previous one
				sack->gapAckBlocks.push_back({
					static_cast<uint16_t>(gap-lastReceivedTransmissionSequenceNumber),
					static_cast<uint16_t>(prev-lastReceivedTransmissionSequenceNumber)
				});
			//Start new gap
			gap = current;
		}
		
		//Move next
		prev = current;
		
		//Remove everything up to the cumulative sequence number
		if (lastReceivedTransmissionSequenceNumber==current)
			//Erase and move
			it = receivedTransmissionSequenceNumbers.erase(it);
		else
			//Next
			++it;
	}
	
	//If we had a gap start
	if (gap!=MaxTransmissionSequenceNumber)
		//Add block ending at last one
		sack->gapAckBlocks.push_back({
			static_cast<uint16_t>(gap-lastReceivedTransmissionSequenceNumber),
			static_cast<uint16_t>(prev-lastReceivedTransmissionSequenceNumber)
		});
		
	//Set last consecutive recevied number
	sack->cumulativeTrasnmissionSequenceNumberAck = receivedTransmissionSequenceNumberWrapper.UnWrap(lastReceivedTransmissionSequenceNumber);
	
	//Set window
	sack->adveritsedReceiverWindowCredit = localAdvertisedReceiverWindowCredit;
	
	//Send it
	Enqueue(sack);
	
	//No need to acknoledge
	pendingAcknowledge = false;
	
	//Reset any pending sack timer
	if (sackTimer)
	{
		//Cancel it
		sackTimer->Cancel();
		//Reset it
		sackTimer = nullptr;
	}
}

void Association::Enqueue(const Chunk::shared& chunk)
{
	bool wasPending = pendingData;
	//Push back
	queue.push_back(chunk);
	//Reset flag
	pendingData = true;
	//If it is first
	if (!wasPending && onPendingData)
		//Call callback
		onPendingData();
}

}; //namespace sctp
