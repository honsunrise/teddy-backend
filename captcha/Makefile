
GOPATH:=$(shell go env GOPATH)

.PHONY: proto test docker


proto:
	protoc --proto_path=${GOPATH}/src:. --go_out=plugins=grpc:. proto/captcha.proto

build: proto
	go build -o captcha-srv main.go

test:
	go test -v ./... -cover

docker:
	docker build . -t captcha-srv:latest
