#!/bin/bash
# Build ingestwebhook command-line tool for 64-bit MacOS and Linux

VERSION=$(git describe --tags)

[[ -d build ]] || mkdir build
GOOS=darwin GOARCH=amd64 go build -o build/ingestwebhook.${VERSION}.darwin-amd64 . || exit 1
GOOS=linux GOARCH=amd64 go build -o build/ingestwebhook.${VERSION}.linux-amd64 . || exit 1
openssl dgst -sha256 build/*.${VERSION}.* | sed -e 's|build/||g'
