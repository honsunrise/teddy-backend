
GOPATH:=$(shell go env GOPATH)

.PHONY: proto test docker


proto:
	protoc --proto_path=${GOPATH}/src:. --go_out=plugins=grpc:. proto/message.proto

build: proto
	go build -o message-srv main.go

test:
	go test -v ./... -cover

docker: build
	docker build . -t message-srv:latest

docker-tag: docker
	docker tag message-srv:latest ${DOCKER_ID_USER}/message-srv:latest

docker-push: docker-tag
	docker push ${DOCKER_ID_USER}/message-srv:latest