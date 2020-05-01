export ROOT_DIR=${PWD}

include config.mk


CPPFLAGS = -I${ROOT_DIR}/media-server/ext/crc32/include/  -I${ROOT_DIR}/media-server/ext/libdatachannels/  -I${ROOT_DIR}/media-server/ext/libdatachannels/src/

all:OPENSSL SRTP MP4V2 MEDIASERVER_STATIC
	echo $(ROOT_DIR)

OPENSSL:
	cd ${OPENSSL_SRC} &&  export KERNEL_BITS=64 && ./config --prefix=${OPENSSL_DIR} && make && make install &&  cp -rf ${OPENSSL_DIR}/lib/*.a  ${OPENSSL_DIR}/


SRTP:
	cd ${LIBSRTP_SRC} && ./configure --prefix=${LIBSRTP_DIR} && make && make install && cp -rf ${LIBSRTP_DIR}/lib/*.a ${LIBSRTP_DIR}/


MP4V2:
	cd ${LIBMP4_SRC} && autoreconf -i && ./configure --prefix=${LIBMP4_DIR} && make && make install && cp -rf ${LIBMP4_DIR}/lib/* ${LIBMP4_DIR}/


MEDIASERVER_STATIC:
	cp config.mk  ./media-server/ && make -C media-server libmediaserver.a CPPFLAGS=${CPPFLAGS}
	echo ${ROOT_DIR}

ECHO:
	echo $(ROOT_DIR)
	echo $(OPENSSL_DIR)
	
