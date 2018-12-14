#!/usr/bin/env sh


pushd uaa
make build
popd

pushd message
make build
popd

pushd content
make build
popd

pushd captcha
make build
popd

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
