target_sources(libdatachannels PUBLIC 
	${CMAKE_CURRENT_SOURCE_DIR}/AbortAssociationChunk.cpp
	${CMAKE_CURRENT_SOURCE_DIR}/CookieAckChunk.cpp
	${CMAKE_CURRENT_SOURCE_DIR}/CookieEchoChunk.cpp
	${CMAKE_CURRENT_SOURCE_DIR}/ForwardCumulativeTSNChunk.cpp
	${CMAKE_CURRENT_SOURCE_DIR}/HeartbeatAckChunk.cpp
	${CMAKE_CURRENT_SOURCE_DIR}/HeartbeatRequestChunk.cpp
	${CMAKE_CURRENT_SOURCE_DIR}/InitiationAcknowledgementChunk.cpp
	${CMAKE_CURRENT_SOURCE_DIR}/InitiationChunk.cpp
	${CMAKE_CURRENT_SOURCE_DIR}/OperationErrorChunk.cpp
	${CMAKE_CURRENT_SOURCE_DIR}/PayloadDataChunk.cpp
	${CMAKE_CURRENT_SOURCE_DIR}/ReConfigChunk.cpp
	${CMAKE_CURRENT_SOURCE_DIR}/SelectiveAcknowledgementChunk.cpp
	${CMAKE_CURRENT_SOURCE_DIR}/ShutdownAcknowledgementChunk.cpp
	${CMAKE_CURRENT_SOURCE_DIR}/ShutdownAssociationChunk.cpp
	${CMAKE_CURRENT_SOURCE_DIR}/ShutdownCompleteChunk.cpp
	${CMAKE_CURRENT_SOURCE_DIR}/UnknownChunk.cpp
	${CMAKE_CURRENT_SOURCE_DIR}/PaddingChunk.cpp
)


