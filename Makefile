export ROOT_DIR=${PWD}

include config.mk

all:OPENSSL SRTP MP4V2 MEDIASERVER_STATIC
	echo $(ROOT_DIR)

OPENSSL:
	cd ${OPENSSL_DIR} && make clean && ./Configure darwin64-x86_64-cc && make 
	cd $(ROOT_DIR)

SRTP:
	cd ${LIBSRTP_DIR} && ./configure && make  
	cd $(ROOT_DIR) 

MP4V2:
	cd ${LIBMP4_DIR} && autoreconf -i && ./configure && make 
	cd $(ROOT_DIR)

MEDIASERVER_STATIC:
	cp config.mk  ./mediaserver/ && make -C mediaserver libmediaserver.a 
	echo ${ROOT_DIR}

ECHO:
	echo $(ROOT_DIR)
	echo $(OPENSSL_DIR)
	