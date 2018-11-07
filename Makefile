export ROOT_DIR=${PWD}


LOG				= yes
DEBUG			= no
SANITIZE		= no
STATIC			= yes
STATIC_OPENSSL	= yes
STATIC_LIBSRTP	= yes
STATIC_LIBMP4	= yes
OPENSSL_DIR		= ${ROOT_DIR}/external/openssl
LIBSRTP_DIR		= ${ROOT_DIR}/external/libsrtp
LIBMP4_DIR		= ${ROOT_DIR}/external/mp4v2
VADWEBRTC		= yes
SRCDIR			= ${ROOT_DIR}/mediaserver
IMAGEMAGICK		= no


all:OPENSSL SRTP MP4V2 MEDIASERVER
	echo $(ROOT_DIR)

OPENSSL:
	cd ${OPENSSL_DIR} && ${OPENSSL_DIR}/config no-shared && make 
	cd $(ROOT_DIR)

SRTP:
	cd ${LIBSRTP_DIR} && configure && make  
	cd $(ROOT_DIR) 

MP4V2:
	cd ${LIBMP4_DIR} && autoreconf -i && configure && make 
	cd $(ROOT_DIR)

MEDIASERVER:
	cp config.mk  ./mediaserver/ && cd ${SRCDIR} && make libmediaserver.a 
	cd ${ROOT_DIR}

ECHO:
	echo $(ROOT_DIR)
	