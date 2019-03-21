#ifndef SEQUENCENUMBERWRAPPER_H
#define SEQUENCENUMBERWRAPPER_H

#include <cassert>

namespace sctp
{
	
template <typename T,uint8_t N = sizeof(T)*8, typename E = uint64_t>
class SequenceNumberWrapper
{
public:
	static constexpr const E MaxSequenceNumber	= ~static_cast<E>(0);
	static constexpr const E Mask			= ~static_cast<E>(0) >> static_cast<uint8_t>(sizeof(E)*8-N);
	static constexpr const T OutOfOrderWindow	= ~static_cast<T>(0) >> (N/2);
	
	E Wrap(T seqNum)
	{
		//Input should be withing given range
		assert((seqNum & Mask) == seqNum);
		
		//Current war cycle
		uint64_t seqCycles = cycles;
		
		//If not the first
		if (maxExtSeqNum!=MaxSequenceNumber)
		{
			//Check if we have a sequence wrap 
			if (seqNum<maxSeqNum && maxSeqNum-seqNum>OutOfOrderWindow)
				//Increase warp cycles
				seqCycles = ++cycles;
			//Check if we have a packet from previous cycle
			else if (seqNum>maxSeqNum && seqNum-maxSeqNum>OutOfOrderWindow)
				//It is from the previous one
				--seqCycles;
		}
		
		//Generate extended sequence number
		E extSeqNum = (seqCycles << N) | seqNum;

		//Update maximum seen value
		if (extSeqNum>maxSeqNum || maxExtSeqNum==MaxSequenceNumber)
		{
			//Update max
			maxSeqNum	= seqNum;
			maxExtSeqNum	= extSeqNum;
		}

		//Done
		return extSeqNum;
	}
	
	T UnWrap(E extSeqNum)
	{
		return static_cast<T>(extSeqNum & Mask);
	}
private:
	
	E cycles	= 0;
	T maxSeqNum	= 0;
	E maxExtSeqNum	= MaxSequenceNumber;
	
};

} // namespace sctp

#endif /* SEQUENCENUMBERWRAPPER_H */

