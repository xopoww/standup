#!/bin/bash

set -euo pipefail

PROTOC_VERSION=24.4
PROTOC_DIR="$(dirname "$0")/../dev/tools/protoc"

PROTOC_BIN="$PROTOC_DIR/bin/protoc"
if [ -e "$PROTOC_BIN" ] && [ "$("$PROTOC_BIN" --version | cut -d' ' -f2)" = "$PROTOC_VERSION" ]; then
    echo "protoc v$PROTOC_VERSION already installed at $(realpath "$PROTOC_BIN")"
    exit 0
fi

rm -rf "PROTOC_DIR" || true
mkdir -p "$PROTOC_DIR"

PLATFORM="$(uname | tr '[:upper:]' '[:lower:]')"
PROTOC_PLATFORM=
case "$PLATFORM" in
linux)
    PROTOC_PLATFORM=linux;;
darwin)
    PROTOC_PLATFORM=osx;;
*)
    echo "Unsupported platform $PLATFORM" >&2
    exit 1
esac

ARCHIVE_NAME="protoc-$PROTOC_VERSION-$PROTOC_PLATFORM-x86_64.zip"

DOWNLOAD_PATH="$PROTOC_DIR/$ARCHIVE_NAME"
wget -O "$DOWNLOAD_PATH" "https://github.com/protocolbuffers/protobuf/releases/download/v$PROTOC_VERSION/$ARCHIVE_NAME"
unzip "$DOWNLOAD_PATH" -d "$PROTOC_DIR"