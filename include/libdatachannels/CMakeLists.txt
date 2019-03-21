target_include_directories (libdatachannels PUBLIC ${CMAKE_CURRENT_SOURCE_DIR}/internal)
target_include_directories (libdatachannels PUBLIC ${CMAKE_CURRENT_SOURCE_DIR})

target_sources(libdatachannels PUBLIC 
	${CMAKE_CURRENT_SOURCE_DIR}/Endpoint.cpp
	${CMAKE_CURRENT_SOURCE_DIR}/Datachannel.cpp
)

add_subdirectory (sctp)
