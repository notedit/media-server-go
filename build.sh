#!/usr/bin/env bash

set -e

OS=$(go env GOOS)
ARCH=$(go env GOARCH)
ROOT_DIR=$(pwd)

make 

cp -rf external/openssl/build/include/openssl  include/openssl/$OS-$ARCH/
cp -rf external/libsrtp/build/include/srtp2   include/srtp/
cp -rf external/mp4v2/build/include/mp4v2  include/mp4v2/


cp media-server/bin/release/libmediaserver.a lib/libmediaserver-$OS-$ARCH.a
cp external/libsrtp/build/lib/libsrtp2.a lib/libsrtp2-$OS-$ARCH.a
cp external/openssl/build/lib/libcrypto.a lib/libcrypto-$OS-$ARCH.a
cp external/openssl/build/lib/libssl.a lib/libssl-$OS-$ARCH.a
cp external/mp4v2/build/lib/libmp4v2.a lib/libmp4v2-$OS-$ARCH.a





