target_sources(libdatachannels PUBLIC 
	${CMAKE_CURRENT_SOURCE_DIR}/Association.cpp
	${CMAKE_CURRENT_SOURCE_DIR}/PacketHeader.cpp
	${CMAKE_CURRENT_SOURCE_DIR}/Stream.cpp
	${CMAKE_CURRENT_SOURCE_DIR}/Chunk.cpp
)

add_subdirectory (chunks)
