#!/usr/bin/env bash

GOPATH=$(go env GOPATH)

protoc --proto_path=${GOPATH}/src:. --go_out=plugins=grpc:. policy.proto