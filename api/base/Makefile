
GOPATH:=$(shell go env GOPATH)

.PHONY: test docker

build:
	go build -o base-api main.go plugin.go

test:
	go test -v ./... -cover

docker:
	docker build . -t base-api:latest
