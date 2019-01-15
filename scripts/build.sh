#!/usr/bin/env bash

SOURCE="${BASH_SOURCE[0]}"
while [[ -h "$SOURCE" ]] ; do SOURCE="$(readlink "$SOURCE")"; done
DIR="$( cd -P "$( dirname "$SOURCE" )/.." && pwd )"

pushd "$DIR/uaa"
make build
popd > /dev/null

pushd "$DIR/message"
make build
popd > /dev/null

pushd "$DIR/content"
make build
popd > /dev/null

pushd "$DIR/captcha"
make build
popd > /dev/null

pushd api/uaa
make build
popd

pushd api/base
make build
popd

pushd api/content
make build
popd

pushd api/message
make build
popd
