
GOPATH:=$(shell go env GOPATH)

.PHONY: test docker

build:
	go build -o message-api main.go

test:
	go test -v ./... -cover

docker:
	docker build . -t message-api:latest
