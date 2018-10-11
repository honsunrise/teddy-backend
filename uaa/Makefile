
GOPATH:=$(shell go env GOPATH)

.PHONY: proto test docker


proto:
	protoc --proto_path=${GOPATH}/src:. --go_out=plugins=grpc:. proto/uaa.proto

build: proto
	go build -o uaa-srv main.go

test:
	go test -v ./... -cover

docker:
	docker build . -t uaa-srv:latest
