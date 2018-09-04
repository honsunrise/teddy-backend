
GOPATH:=$(shell go env GOPATH)

.PHONY: test docker

build:
	go build -o uaa-api main.go plugin.go

test:
	go test -v ./... -cover

docker:
	docker build . -t uaa-api:latest
