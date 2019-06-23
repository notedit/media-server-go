#!/usr/bin/env bash

set -e

OS=$(go env GOOS)
ARCH=$(go env GOARCH)
ROOT_DIR=$(pwd)

MEDIASERVER_INCLUDE=$ROOT_DIR/include/media-server/include
CRC32_INCLUDE=$ROOT_DIR/include/crc32c/include
DATACHANNEL_INCUDE=$ROOT_DIR/include/libdatachannels



mkdir -p $MEDIASERVER_INCLUDE
mkdir -p $CRC32_INCLUDE
mkdir -p $DATACHANNEL_INCUDE

cp -rf  ../media-server-go-native/media-server/include/*  $MEDIASERVER_INCLUDE
cp -rf  ../media-server-go-native/media-server/ext/crc32c/include/* $CRC32_INCLUDE
cp -rf  ../media-server-go-native/media-server/ext/libdatachannels/src/* $DATACHANNEL_INCUDE





