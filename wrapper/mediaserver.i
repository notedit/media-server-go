%module(directors="1") native
%{

#include <string>
#include <list>
#include <functional>
#include "../media-server/include/config.h"
#include "../media-server/include/dtls.h"
#include "../media-server/include/OpenSSL.h"
#include "../media-server/include/media.h"
#include "../media-server/include/video.h"
#include "../media-server/include/audio.h"
#include "../media-server/include/rtp.h"
#include "../media-server/include/tools.h"
#include "../media-server/include/rtpsession.h"
#include "../media-server/include/DTLSICETransport.h"
#include "../media-server/include/RTPBundleTransport.h"
#include "../media-server/include/mp4recorder.h"
#include "../media-server/include/mp4streamer.h"
#include "../media-server/include/rtp/RTPStreamTransponder.h"
#include "../media-server/include/ActiveSpeakerDetector.h"
#include "../media-server/include/EventLoop.h"


using RTPBundleTransportConnection = RTPBundleTransport::Connection;
using MediaFrameListener = MediaFrame::Listener;


class PropertiesFacade : private Properties
{
public:
	void SetPropertyInt(const char* key,int intval)
	{
		Properties::SetProperty(key,intval);
	}
	void SetPropertyStr(const char* key,const char* val)
	{
		Properties::SetProperty(key,val);
	}
	void SetPropertyBool(const char* key,bool boolval)
	{
		Properties::SetProperty(key,boolval);
	}
};


class MediaServer
{

public:

	static void Initialize()
	{
		//Initialize ssl
		OpenSSL::ClassInit();
		
		//Start DTLS
		DTLSConnection::Initialize();
	}
	
	static void EnableLog(bool flag)
	{
		//Enable log
		Log("-EnableLog [%d]\n",flag);
		Logger::EnableLog(flag);
	}
	
	static void EnableDebug(bool flag)
	{
		//Enable debug
		Logger::EnableDebug(flag);
	}
	
	static void EnableUltraDebug(bool flag)
	{
		//Enable debug
		Log("-EnableUltraDebug [%d]\n",flag);
		Logger::EnableUltraDebug(flag);
	}
	
	static bool SetPortRange(int minPort, int maxPort)
	{
		return RTPTransport::SetPortRange(minPort,maxPort);
	}
	
	static std::string GetFingerprint()
	{
		return DTLSConnection::GetCertificateFingerPrint(DTLSConnection::Hash::SHA256);
	}
};


class RTPSessionFacade : 	
	public RTPSender,
	public RTPReceiver,
	public RTPSession
{
public:
	RTPSessionFacade(MediaFrame::Type media) : RTPSession(media,NULL)
	{
		//Delegate to group
		delegate = true;
		//Start group dispatch
		GetIncomingSourceGroup()->Start();
	}
	virtual ~RTPSessionFacade() = default;

	virtual int Enqueue(const RTPPacket::shared& packet)	 { return SendPacket(packet); }
	virtual int Enqueue(const RTPPacket::shared& packet,std::function<RTPPacket::shared(const RTPPacket::shared&)> modifier) { return SendPacket(modifier(packet)); }
	virtual int SendPLI(DWORD ssrc)				 { return RequestFPU();}
	
	int Init(const Properties &properties)
	{
		RTPMap rtp;
		RTPMap apt;
		
		//Get codecs
		std::vector<Properties> codecs;
		properties.GetChildrenArray("codecs",codecs);

		//For each codec
		for (auto it = codecs.begin(); it!=codecs.end(); ++it)
		{
			
			BYTE codec;
			//Depending on the type
			switch (GetMediaType())
			{
				case MediaFrame::Audio:
					codec = (BYTE)AudioCodec::GetCodecForName(it->GetProperty("codec"));
					break;
				case MediaFrame::Video:
					codec = (BYTE)VideoCodec::GetCodecForName(it->GetProperty("codec"));
					break;
				default:
					//skip 
					continue;
					
			}

			if (codec == (BYTE)-1) {
				continue;
			}

			//Get codec type
			BYTE type = it->GetProperty("pt",0);
			//ADD it
			rtp[type] = codec;
		}
	
		//Set local 
		RTPSession::SetSendingRTPMap(rtp,apt);
		RTPSession::SetReceivingRTPMap(rtp,apt);
		
		//Call parent
		return RTPSession::Init();
	}
};


class MP4RecorderFacade :
    public MP4Recorder,
    public MP4Recorder::Listener
{
public:
    MP4RecorderFacade() :
        MP4Recorder(this)
    {

    }

    void onFirstFrame(QWORD time) override
    {
        // todo
    }

    void onClosed() override
    {
        // todo
    }
};


class ActiveTrackListener {
public:
	ActiveTrackListener()
	{

	}
	virtual ~ActiveTrackListener() {

	}
	virtual void onActiveTrackchanged(uint32_t id){

	}
};




class MediaFrameSessionFacade :
	public RTPReceiver
{
public:
	MediaFrameSessionFacade(MediaFrame::Type media):
	source(media,loop)
	{

		loop.Start(-1);
		source.Start();
		mediatype = media;
	}
	int Init(const Properties &properties)
	{
		//Get codecs
		std::vector<Properties> codecs;
		properties.GetChildrenArray("codecs",codecs);

		//For each codec
		for (auto it = codecs.begin(); it!=codecs.end(); ++it)
		{

			BYTE codec;
			//Depending on the type
			switch (mediatype)
			{
				case MediaFrame::Audio:
					codec = (BYTE)AudioCodec::GetCodecForName(it->GetProperty("codec"));
					audioCodec = codec;
					break;
				case MediaFrame::Video:
					codec = (BYTE)VideoCodec::GetCodecForName(it->GetProperty("codec"));
					videoCodec = codec;
					break;
				default:
					///Ignore
					codec = (BYTE)-1;
					break;
			}

			//Get codec type
			BYTE type = it->GetProperty("pt",0);
			//ADD it
			rtp[type] = codec;
		}

		return 1;
	}

	void onRTPPacket(uint8_t* data, int size)
	{


		// Run on thread
		loop.Async([=](...)  {

			RTPHeader header;
			RTPHeaderExtension extension;

			int lsize = size;
			int len = header.Parse(data,lsize);

			if (!len)
			{
				//Debug
				Debug("-MediaFrameSessionFacade::onRTPPacket() | Could not parse RTP header\n");
				return;
			}

			if (header.extension)
			{

				//Parse extension
				int l = extension.Parse(extMap,data+len,lsize-len);
				//If not parsed
				if (!l)
				{
					///Debug
					Debug("-MediaFrameSessionFacade::onRTPPacket() | Could not parse RTP header extension\n");
					//Exit
					return;
				}
				//Inc ini
				len += l;
			}

			if (header.padding)
			{
				//Get last 2 bytes
				WORD padding = get1(data,lsize-1);
				//Ensure we have enought size
				if (size-len<padding)
				{
					///Debug
					Debug("-PCAPTransportEmulator::Run() | RTP padding is bigger than size [padding:%u,size%u]\n",padding,size);
					//Ignore this try again
					return;
				}
				//Remove from size
				lsize -= padding;
			}


			DWORD ssrc = header.ssrc;
			BYTE type  = header.payloadType;
			//Get initial codec
			BYTE codec = rtp.GetCodecForType(header.payloadType);;

			//Check codec
			if (codec==RTPMap::NotFound)
			{
				//Exit
				Error("-MediaFrameSessionFacade::onRTPPacket(%s) | RTP packet type unknown [%d]\n",MediaFrame::TypeToString(mediatype),type);
				//Exit
				return;
			}

			auto packet = std::make_shared<RTPPacket>(mediatype,codec,header,extension);

			//Set the payload
			packet->SetPayload(data+len,lsize-len);

			WORD seq = packet->GetSeqNum();

			source.media.SetSeqNum(seq);

			if (source.media.ssrc != ssrc) {
				source.media.Reset();
				source.media.ssrc = ssrc;
			}

			source.media.Update(getTimeMS(),packet->GetSeqNum(),packet->GetRTPHeader().GetSize()+packet->GetMediaLength());

			WORD cycles = source.media.SetSeqNum(packet->GetSeqNum());
            //Set cycles back
            packet->SetSeqCycles(cycles);

			source.AddPacket(packet,0);

		});
	}

	void onRTPData(uint8_t* data, int size, uint32_t timestamp)
	{

	}

	RTPIncomingSourceGroup* GetIncomingSourceGroup()
	{
		return &source;
	}
	int End()
	{
		Log("MediaFrameSessionFacade End\n");
		return 1;
	}
	virtual int SendPLI(DWORD ssrc) {
		return 1;
	}

private:
	RTPMap extMap;
	RTPMap rtp;
	RTPMap apt;
	BYTE audioCodec;
	BYTE videoCodec;

	DWORD ssrc = 0;
	DWORD extSeqNum = 0;

	MediaFrame::Type mediatype;
	EventLoop loop;
	RTPIncomingSourceGroup source;
};




class RTPSenderFacade
{
public:	
	RTPSenderFacade(DTLSICETransport* transport)
	{
		sender = transport;
	}

	RTPSenderFacade(RTPSessionFacade* session)
	{
		sender = session;
	}
	
	RTPSender* get() { return sender;}
private:
	RTPSender* sender;
};

class RTPReceiverFacade
{
public:	
	RTPReceiverFacade(DTLSICETransport* transport)
	{
		receiver = transport;
	}

	RTPReceiverFacade(RTPSessionFacade* session)
	{
		receiver = session;
	}

    RTPReceiverFacade(MediaFrameSessionFacade* session)
    {
        receiver = session;
    }


	int SendPLI(DWORD ssrc)
	{
		return receiver ? receiver->SendPLI(ssrc) : 0;
	}
	
	RTPReceiver* get() { return receiver;}
private:
	RTPReceiver* receiver;
};


RTPSenderFacade* TransportToSender(DTLSICETransport* transport)
{
	return new RTPSenderFacade(transport);
}
RTPReceiverFacade* TransportToReceiver(DTLSICETransport* transport)
{
	return new RTPReceiverFacade(transport);
}

RTPSenderFacade* SessionToSender(RTPSessionFacade* session)
{
	return new RTPSenderFacade(session);	
}

RTPReceiverFacade* SessionToReceiver(RTPSessionFacade* session)
{
	return new RTPReceiverFacade(session);
}


RTPReceiverFacade* RTPSessionToReceiver(MediaFrameSessionFacade* session)
{
	return new RTPReceiverFacade(session);
}




class RTPStreamTransponderFacade : 
	public RTPStreamTransponder
{
public:
	RTPStreamTransponderFacade(RTPOutgoingSourceGroup* outgoing,RTPSenderFacade* sender) :
		RTPStreamTransponder(outgoing, sender ? sender->get() : NULL)
	{}

	bool SetIncoming(RTPIncomingMediaStream* incoming, RTPReceiverFacade* receiver)
	{
		return RTPStreamTransponder::SetIncoming(incoming, receiver ? receiver->get() : NULL);
	}

	bool SetIncoming(RTPIncomingMediaStream* incoming, RTPReceiver* receiver)
	{
		return RTPStreamTransponder::SetIncoming(incoming, receiver);
	}
	
	virtual void onREMB(RTPOutgoingSourceGroup* group,DWORD ssrc, DWORD bitrate) override
	{
		// todo  make callback
		// Log("onREMB\n");
	}
	void SetMinPeriod(DWORD period) { this->period = period; }

private:
	DWORD period = 1000;
	QWORD last = 0;
};



class StreamTrackDepacketizer :
	public RTPIncomingMediaStream::Listener
{
public:
	StreamTrackDepacketizer(RTPIncomingMediaStream* incomingSource)
	{
		//Store
		this->incomingSource = incomingSource;
		//Add us as RTP listeners
		this->incomingSource->AddListener(this);
		//No depkacketixer yet
		depacketizer = NULL;
	}

	virtual ~StreamTrackDepacketizer()
	{
		//JIC
		Stop();
		//Check 
		if (depacketizer)
			//Delete depacketier
			delete(depacketizer);
	}

	virtual void onRTP(RTPIncomingMediaStream* group,const RTPPacket::shared& packet)
	{

	    if (listeners.empty())
	           return;


		//If depacketizer is not the same codec 
		if (depacketizer && depacketizer->GetCodec()!=packet->GetCodec())
		{
			//Delete it
			delete(depacketizer);
			//Create it next
			depacketizer = NULL;
		}
		//If we don't have a depacketized
		if (!depacketizer)
			//Create one
			depacketizer = RTPDepacketizer::Create(packet->GetMedia(),packet->GetCodec());
		//Ensure we have it
		if (!depacketizer)
			//Do nothing
			return;
		//Pass the pakcet to the depacketizer
		 MediaFrame* frame = depacketizer->AddPacket(packet);
		 
		 //If we have a new frame
		 if (frame)
		 {
			 //Call all listeners
			 for (const auto& listener : listeners)
				 //Call listener
				 listener->onMediaFrame(packet->GetSSRC(),*frame);
			 //Next
			 depacketizer->ResetFrame();
		 }	
	}

	virtual void onBye(RTPIncomingMediaStream* group) 
	{
		if(depacketizer)
			//Skip current
			depacketizer->ResetFrame();
	}
	
	virtual void onEnded(RTPIncomingMediaStream* group) 
	{
		if (incomingSource==group)
			incomingSource = nullptr;
	}
	
	void AddMediaListener(MediaFrame::Listener *listener)
	{
		//Add to set
		if (!incomingSource || !listener)
			//Done
			return;
		//Add listener async
		incomingSource->GetTimeService().Async([=](...){
			//Add to set
			listeners.insert(listener);
		});
	}
	
	void RemoveMediaListener(MediaFrame::Listener *listener)
	{
		//Remove from set
		if (!incomingSource)
			//Done
			return;

		//Add listener sync so it can be deleted after this call
		incomingSource->GetTimeService().Sync([=](...){
			//Remove from set
			listeners.erase(listener);
		});
	}
	
	void Stop()
	{
		//If already stopped
		if (!incomingSource)
			//Done
			return;
		
		//Stop listeneing
		incomingSource->RemoveListener(this);
		//Clean it
		incomingSource = NULL;
	}
	
private:
    std::set<MediaFrame::Listener*> listeners;
	RTPDepacketizer* depacketizer;
	RTPIncomingMediaStream* incomingSource;
};



class DTLSICETransportListener :
	public DTLSICETransport::Listener
{
public:
	DTLSICETransportListener()
	{

 	}

 	virtual ~DTLSICETransportListener() = default;

 	virtual void onRemoteICECandidateActivated(const std::string& ip, uint16_t port, uint32_t priority) override
 	{

 	    // todo
 	}

 	virtual void onDTLSStateChanged(const DTLSICETransport::DTLSState state) override 
	{

		switch(state)
		{
			case DTLSICETransport::DTLSState::New:
				onDTLSStateChange(0);
				break;
			case DTLSICETransport::DTLSState::Connecting:
				onDTLSStateChange(1);
				break;
			case DTLSICETransport::DTLSState::Connected:
				onDTLSStateChange(2);
				break;
			case DTLSICETransport::DTLSState::Closed:
				onDTLSStateChange(3);
				break;
			case DTLSICETransport::DTLSState::Failed:
				onDTLSStateChange(4);
				break;
		}
	}

	virtual void onDTLSStateChange(uint32_t state)
	{

	}
};




class SenderSideEstimatorListener : 
	public RemoteRateEstimator::Listener
{
public:
	SenderSideEstimatorListener()
	{
		
	}
	
	virtual void onTargetBitrateRequested(DWORD bitrate) override 
	{
        // todo make callback
	}

private:
	DWORD period  = 500;
	QWORD last = 0;
};



EvenSource::EvenSource()
{
}

EvenSource::EvenSource(const char* str)
{
}

EvenSource::EvenSource(const std::wstring &str)
{
}

EvenSource::~EvenSource()
{
}

void EvenSource::SendEvent(const char* type,const char* msg,...)
{
}


class LayerSources : public std::vector<LayerSource*>
{
public:
	size_t size() const  { return std::vector<LayerSource*>::size(); }
	LayerSource* get(size_t i)	{ return  std::vector<LayerSource*>::at(i); }
};



class ActiveSpeakerDetectorFacade :
	public ActiveSpeakerDetector,
	public ActiveSpeakerDetector::Listener,
	public RTPIncomingMediaStream::Listener
{
public:	
	ActiveSpeakerDetectorFacade(ActiveTrackListener* listener) :
		ActiveSpeakerDetector(this),
		listener(listener)
	{};
		
	virtual void onActiveSpeakerChanded(uint32_t id) override
	{
        // todo make callback

		if (listener) 
		{
			listener->onActiveTrackchanged(id);
		}
	}
	
	void AddIncomingSourceGroup(RTPIncomingMediaStream* incoming, uint32_t id)
	{
			if (incoming)
    		{
    			ScopedLock lock(mutex);
    			//Insert new
    			auto [it,inserted] = sources.try_emplace(incoming,id);
    			//If already present
    			if (!inserted)
    				//do nothing
    				return;
    			//Add us as rtp listeners
    			incoming->AddListener(this);
    			//initialize to silence
    			ActiveSpeakerDetector::Accumulate(id, false, 127, getTimeMS());
    		}
	}
	
	void RemoveIncomingSourceGroup(RTPIncomingMediaStream* incoming)
	{
		if (incoming)
		{	
			ScopedLock lock(mutex);
			//Get map
			auto it = sources.find(incoming);
			//check it was present
			if (it==sources.end())
				//Do nothing
				return;
			//Remove listener
			incoming->RemoveListener(this);
			//RElease id
			ActiveSpeakerDetector::Release(it->second);
			//Erase
			sources.erase(it);
		}
	}
	
	virtual void onRTP(RTPIncomingMediaStream* incoming,const RTPPacket::shared& packet) override
	{
        if (packet->HasAudioLevel())
        {
            ScopedLock lock(mutex);
            //Get map
            auto it = sources.find(incoming);
            //check it was present
            if (it==sources.end())
                //Do nothing
                return;
            //Accumulate on id
            ActiveSpeakerDetector::Accumulate(it->second, packet->GetVAD(),packet->GetLevel(), getTimeMS());
        }
	}		

	virtual void onBye(RTPIncomingMediaStream* incoming) override
	{

	}
	
	virtual void onEnded(RTPIncomingMediaStream* incoming) override
	{
		if (incoming)
		{
			ScopedLock lock(mutex);
			//Get map
			auto it = sources.find(incoming);
			//check it was present
			if (it==sources.end())
				//Do nothing
			//Release id
			ActiveSpeakerDetector::Release(it->second);
			//Erase
			sources.erase(it);
		}
	}
private:
	Mutex mutex;
	std::map<RTPIncomingMediaStream*,uint32_t> sources;
	ActiveTrackListener* listener;
};


class  MediaFrameListenerFacade:
	public MediaFrameListener
{
public:
	MediaFrameListenerFacade()
	{

	}

	virtual void onMediaFrame(const MediaFrame &frame)  {

	}
	virtual void onMediaFrame(DWORD ssrc, const MediaFrame &frame) {

		onMediaFrame(frame);
	}
	
};


class MediaFrameMultiplexer :
	public RTPIncomingMediaStream::Listener
{
public:
	MediaFrameMultiplexer(RTPIncomingMediaStream* incomingSource)
	{
		//Store
		this->incomingSource = incomingSource;
		//Add us as RTP listeners
		this->incomingSource->AddListener(this);
		//No depkacketixer yet
		depacketizer = NULL;
	}

	virtual ~MediaFrameMultiplexer()
	{
		//JIC
		Stop();
		//Check 
		if (depacketizer)
			//Delete depacketier
			delete(depacketizer);
	}

	virtual void onRTP(RTPIncomingMediaStream* group,const RTPPacket::shared& packet)
	{

		if (listeners.empty()) 
			return;


		if (depacketizer && depacketizer->GetCodec()!=packet->GetCodec())
		{
			//Delete it
			delete(depacketizer);
			//Create it next
			depacketizer = NULL;
		}
		//If we don't have a depacketized
		if (!depacketizer)
			//Create one
			depacketizer = RTPDepacketizer::Create(packet->GetMedia(),packet->GetCodec());
		//Ensure we have it
		if (!depacketizer)
			//Do nothing
			return;
		//Pass the pakcet to the depacketizer
		 MediaFrame* frame = depacketizer->AddPacket(packet);
		 
		 //If we have a new frame
		 if (frame)
		 {
			 //Call all listeners
			 for (Listeners::const_iterator it = listeners.begin();it!=listeners.end();++it)
				 //Call listener
				 (*it)->onMediaFrame(*frame);
			 //Next
			 depacketizer->ResetFrame();
		 }

	}

	virtual void onBye(RTPIncomingMediaStream* group) 
	{
		if(depacketizer)
			//Skip current
			depacketizer->ResetFrame();
	}
	
	virtual void onEnded(RTPIncomingMediaStream* group) 
	{
		if (incomingSource==group)
			incomingSource = nullptr;
	}
	
	void AddMediaListener(MediaFrameListener *listener)
	{
		//Add to set
		listeners.insert(listener);
	}
	
	void RemoveMediaListener(MediaFrameListener *listener)
	{
		//Remove from set
		listeners.erase(listener);
	}
	
	void Stop()
	{
		//If already stopped
		if (!incomingSource)
			//Done
			return;
		
		//Stop listeneing
		incomingSource->RemoveListener(this);
		//Clean it
		incomingSource = NULL;
	}
	
private:
	typedef std::set<MediaFrameListener*> Listeners;
private:
	Listeners listeners;
	RTPDepacketizer* depacketizer;
	RTPIncomingMediaStream* incomingSource;
};

%}


%feature("director") SenderSideEstimatorListener;
%feature("director") MediaFrameListenerFacade;
%feature("director") ActiveTrackListener;
%feature("director") DTLSICETransportListener;



%include <typemaps.i>
%include "stdint.i"
%include "std_string.i"
%include "std_vector.i"


#define QWORD		uint64_t
#define DWORD		uint32_t
#define WORD		uint16_t
#define SWORD		int16_t
#define BYTE		uint8_t
#define SBYTE		char


%include "../media-server/include/acumulator.h"


%{
using MediaFrameType = MediaFrame::Type;
%}
enum MediaFrameType;




struct LayerInfo
{
	static BYTE MaxLayerId; 
	BYTE temporalLayerId = MaxLayerId;
	BYTE spatialLayerId  = MaxLayerId;
};

struct LayerSource : public LayerInfo
{
	DWORD		numPackets;
	DWORD		totalBytes;
	DWORD		bitrate;
};

class LayerSources
{
public:
	size_t size() const;
	LayerSource* get(size_t i);
};

struct RTPSource 
{
	DWORD ssrc;
	DWORD extSeqNum;
	DWORD cycles;
	DWORD jitter;
	DWORD numPackets;
	DWORD numRTCPPackets;
	DWORD totalBytes;
	DWORD totalRTCPBytes;
	DWORD bitrate;
};


struct RTPIncomingSource : public RTPSource
{
	DWORD lostPackets;
	DWORD dropPackets;
	DWORD totalPacketsSinceLastSR;
	DWORD totalBytesSinceLastSR;
	DWORD minExtSeqNumSinceLastSR ;
	DWORD lostPacketsSinceLastSR;
	QWORD lastReceivedSenderNTPTimestamp;
	QWORD lastReceivedSenderReport;
	QWORD lastReport;
	QWORD lastPLI;
	DWORD totalPLIs;
	DWORD totalNACKs;
	QWORD lastNACKed;
	
	%extend 
	{
		LayerSources layers() 
		{
			LayerSources layers;
			for(auto it = $self->layers.begin(); it != $self->layers.end(); ++it )
				layers.push_back(&(it->second));
			return layers;
		}
	}
};
	
struct RTPOutgoingSource : public RTPSource
{
	
	DWORD time;
	DWORD lastTime;
	DWORD numPackets;
	DWORD numRTCPPackets;
	DWORD totalBytes;
	DWORD totalRTCPBytes;
	QWORD lastSenderReport;
	QWORD lastSenderReportNTP;
};


%nodefaultctor TimeService;
struct TimeService
{
	
};


struct RTPOutgoingSourceGroup
{
	RTPOutgoingSourceGroup(MediaFrameType type);
	RTPOutgoingSourceGroup(std::string &streamId,MediaFrameType type);
	
	MediaFrameType  type;
	RTPOutgoingSource media;
	RTPOutgoingSource fec;
	RTPOutgoingSource rtx;

	void Update();
};


%nodefaultctor RTPSender;
%nodefaultdtor RTPSender; 
struct RTPSender {};

%nodefaultctor RTPReceiver;
%nodefaultdtor RTPReceiver; 
struct RTPReceiver {};

%{
using RTPIncomingMediaStreamListener = RTPIncomingMediaStream::Listener;
%}
%nodefaultctor RTPIncomingMediaStreamListener;
struct RTPIncomingMediaStreamListener
{

};


%nodefaultctor RTPIncomingMediaStream;
%nodefaultdtor RTPIncomingMediaStream; 
struct RTPIncomingMediaStream {

};



struct RTPIncomingSourceGroup : public RTPIncomingMediaStream
{
	RTPIncomingSourceGroup(MediaFrameType type, TimeService& TimeService);
	std::string rid;
	std::string mid;
	DWORD rtt;
	MediaFrameType  type;
	RTPIncomingSource media;
	RTPIncomingSource fec;
	RTPIncomingSource rtx;

	DWORD lost;
	DWORD minWaitedTime;
	DWORD maxWaitedTime;
	double avgWaitedTime;

	void AddListener(RTPIncomingMediaStreamListener* listener);
	void RemoveListener(RTPIncomingMediaStreamListener* listener);

	void Update();
};


struct RTPIncomingMediaStreamMultiplexer :  public RTPIncomingMediaStreamListener,public RTPIncomingMediaStream
{
	RTPIncomingMediaStreamMultiplexer(DWORD ssrc, TimeService& TimeService);
	void Stop();
};



class PropertiesFacade : private Properties
{
public:
	void SetPropertyInt(const char* key,int intval);
	void SetPropertyStr(const char* key,const char* val);
	void SetPropertyBool(const char* key,bool boolval);
};

class MediaServer
{
public:
	static void Initialize();
	static void EnableLog(bool flag);
	static void EnableDebug(bool flag);
	static void EnableUltraDebug(bool flag);
	static std::string GetFingerprint();
	static bool SetPortRange(int minPort, int maxPort);
};



%nodefaultctor RTPBundleTransportConnection;
%nodefaultdtor RTPBundleTransportConnection;
struct RTPBundleTransportConnection
{
	DTLSICETransport* transport;
	bool disableSTUNKeepAlive;
	size_t iceRequestsSent		= 0;
	size_t iceRequestsReceived	= 0;
	size_t iceResponsesSent		= 0;
	size_t iceResponsesReceived	= 0;
};



class RTPBundleTransport
{
public:
	RTPBundleTransport();
	int Init();
	int Init(int port);
	RTPBundleTransportConnection* AddICETransport(const std::string &username,const Properties& properties);
	int RemoveICETransport(const std::string &username);
	int End();
	int GetLocalPort() const { return port; }
	int AddRemoteCandidate(const std::string& username,const char* ip, WORD port);
	bool SetAffinity(int cpu);
	void SetIceTimeout(uint32_t timeout);
	TimeService& GetTimeService();
	void DeleteOutGoingSourceGroup(RTPOutgoingSourceGroup *outgoing);
};



class DTLSICETransportListener
{
public:
	DTLSICETransportListener();
	virtual ~DTLSICETransportListener() {};
	// swig does not support inter class
	virtual void onDTLSStateChange(uint32_t state);
};


%{
using RemoteRateEstimatorListener = RemoteRateEstimator::Listener;
%}
%nodefaultctor RemoteRateEstimatorListener;
struct RemoteRateEstimatorListener
{
};


%nodefaultctor DTLSICETransport; 
class DTLSICETransport
{
public:

	void SetListener(DTLSICETransportListener* listener);

	void Start();
	void Stop();
	
	void SetSRTPProtectionProfiles(const std::string& profiles);
	void SetRemoteProperties(const Properties& properties);
	void SetLocalProperties(const Properties& properties);
	virtual int SendPLI(DWORD ssrc) override;
	virtual int Enqueue(const RTPPacket::shared& packet) override;
	int Dump(const char* filename, bool inbound = true, bool outbound = true, bool rtcp = true, bool rtpHeadersOnly = false);
	int Dump(UDPDumper* dumper, bool inbound = true, bool outbound = true, bool rtcp = true, bool rtpHeadersOnly = false);
	int DumpBWEStats(const char* filename);
	void Reset();
	
	void ActivateRemoteCandidate(ICERemoteCandidate* candidate,bool useCandidate, DWORD priority);
	int SetRemoteCryptoDTLS(const char *setup,const char *hash,const char *fingerprint);
	int SetLocalSTUNCredentials(const char* username, const char* pwd);
	int SetRemoteSTUNCredentials(const char* username, const char* pwd);
	bool AddOutgoingSourceGroup(RTPOutgoingSourceGroup *group);
	bool RemoveOutgoingSourceGroup(RTPOutgoingSourceGroup *group);
	bool AddIncomingSourceGroup(RTPIncomingSourceGroup *group);
	bool RemoveIncomingSourceGroup(RTPIncomingSourceGroup *group);
	
	void SetBandwidthProbing(bool probe);
	void SetMaxProbingBitrate(DWORD bitrate);
	void SetProbingBitrateLimit(DWORD bitrate);
	void SetSenderSideEstimatorListener(RemoteRateEstimatorListener* listener);
	
	const char* GetRemoteUsername() const;
	const char* GetRemotePwd()	const;
	const char* GetLocalUsername()	const;
	const char* GetLocalPwd()	const;
	
	DWORD GetRTT() const { return rtt; }

	QWORD GetLastActiveTime() { return lastActiveTime; }
	
	TimeService& GetTimeService();
};



class RTPSessionFacade :
	public RTPSender,
	public RTPReceiver
{
public:
	RTPSessionFacade(MediaFrameType media);
	int Init(const Properties &properties);
	int SetLocalPort(int recvPort);
	int GetLocalPort();
	int SetRemotePort(char *ip,int sendPort);
	RTPOutgoingSourceGroup* GetOutgoingSourceGroup();
	RTPIncomingSourceGroup* GetIncomingSourceGroup();
	int End();
	virtual int Enqueue(const RTPPacket::shared& packet);
	virtual int SendPLI(DWORD ssrc);
};


class RTPSenderFacade
{
public:	
	RTPSenderFacade(DTLSICETransport* transport);
	RTPSenderFacade(RTPSessionFacade* session);
	RTPSender* get();

};

class RTPReceiverFacade
{
public:	
	RTPReceiverFacade(DTLSICETransport* transport);
	RTPReceiverFacade(RTPSessionFacade* session);
	RTPReceiver* get();
	int SendPLI(DWORD ssrc);
};


RTPSenderFacade*	TransportToSender(DTLSICETransport* transport);
RTPReceiverFacade*	TransportToReceiver(DTLSICETransport* transport);
RTPSenderFacade*	SessionToSender(RTPSessionFacade* session);
RTPReceiverFacade*	SessionToReceiver(RTPSessionFacade* session);
RTPReceiverFacade*  RTPSessionToReceiver(MediaFrameSessionFacade* session);


class RTPStreamTransponderFacade 
{
public:
	RTPStreamTransponderFacade(RTPOutgoingSourceGroup* outgoing,RTPSenderFacade* sender);
	bool SetIncoming(RTPIncomingMediaStream* incoming, RTPReceiverFacade* receiver);
	bool SetIncoming(RTPIncomingMediaStream* incoming, RTPReceiver* receiver);
	void SelectLayer(int spatialLayerId,int temporalLayerId);
	void Mute(bool muting);
	void Close();
};


%nodefaultctor MediaFrameListener;
%nodefaultdtor MediaFrameListener;
struct MediaFrameListener
{
};


class StreamTrackDepacketizer 
{
public:
	StreamTrackDepacketizer(RTPIncomingMediaStream* incomingSource);
	void AddMediaListener(MediaFrameListener* listener);
	void RemoveMediaListener(MediaFrameListener* listener);
	void Stop();
};




class MP4RecorderFacade :
	public MediaFrameListener
{
public:
	MP4RecorderFacade();

	//Recorder interface
	virtual bool Create(const char *filename);
	virtual bool Record();
	virtual bool Record(bool waitVideo);
	virtual bool Stop();
	virtual bool Close();
	void SetTimeShiftDuration(DWORD duration);
	bool Close(bool async);
};


class MediaFrameSessionFacade :
	public RTPReceiver
{
public:
	MediaFrameSessionFacade(MediaFrameType media);
	int Init(const Properties &properties);
	void onRTPPacket(uint8_t* buffer, int len);
	void onRTPData(uint8_t* buffer, int size, uint8_t payloadType);
	RTPIncomingSourceGroup* GetIncomingSourceGroup();
	int End();
	virtual int SendPLI(DWORD ssrc);
};


class SenderSideEstimatorListener :
	public RemoteRateEstimatorListener
{
public:
	SenderSideEstimatorListener();
	virtual ~SenderSideEstimatorListener() {}
	void onTargetBitrateRequested(DWORD bitrate);
};


class ActiveSpeakerDetectorFacade
{
public:	
	ActiveSpeakerDetectorFacade(ActiveTrackListener* listener);
	void SetMinChangePeriod(uint32_t minChangePeriod);
	void SetMaxAccumulatedScore(uint64_t maxAcummulatedScore);
	void SetNoiseGatingThreshold(uint8_t noiseGatingThreshold);
	void SetMinActivationScore(uint32_t minActivationScore);
	void AddIncomingSourceGroup(RTPIncomingMediaStream* incoming, uint32_t id);
	void RemoveIncomingSourceGroup(RTPIncomingMediaStream* incoming);
};


class MediaFrameListenerFacade
{
public:
	MediaFrameListenerFacade();
	virtual ~MediaFrameListenerFacade() {}
	virtual void onMediaFrame(const MediaFrame &frame);
};


class MediaFrameMultiplexer
{
public:
	MediaFrameMultiplexer(RTPIncomingMediaStream* incomingSource);
	void AddMediaListener(MediaFrameListenerFacade* listener);
	void RemoveMediaListener(MediaFrameListenerFacade* listener);
	void Stop();
};


class ActiveTrackListener {
public:
	ActiveTrackListener();
	virtual ~ActiveTrackListener() {}
	virtual void onActiveTrackchanged(uint32_t id);
};




