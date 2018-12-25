#!/usr/bin/env bash

set -e

OS=$(go env GOOS)
ARCH=$(go env GOARCH)
ROOT_DIR=$(pwd)

make 

cp -rf external/openssl/include/openssl  external/opensslinclude/$OS-$ARCH/

cp media-server/bin/release/libmediaserver.a lib/libmediaserver-$OS-$ARCH.a
cp external/libsrtp/libsrtp2.a lib/libsrtp2-$OS-$ARCH.a
cp external/openssl/libcrypto.a lib/libcrypto-$OS-$ARCH.a
cp external/openssl/libssl.a lib/libssl-$OS-$ARCH.a
cp external/mp4v2/.libs/libmp4v2.a lib/libmp4v2-$OS-$ARCH.a





