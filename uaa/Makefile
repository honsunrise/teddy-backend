
GOPATH:=$(shell go env GOPATH)

.PHONY: proto test docker


proto:
	protoc --proto_path=${GOPATH}/src:. --go_out=plugins=grpc:. proto/uaa.proto

build: proto
	go build -o uaa-srv main.go

test:
	go test -v ./... -cover

docker: build
	docker build . -t uaa-srv:latest

docker-tag: docker
	docker tag uaa-srv:latest ${DOCKER_ID_USER}/uaa-srv:latest

docker-push: docker-tag
	docker push ${DOCKER_ID_USER}/uaa-srv:latest