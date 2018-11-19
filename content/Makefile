
GOPATH:=$(shell go env GOPATH)

.PHONY: proto test docker


proto:
	protoc --proto_path=${GOPATH}/src:. --go_out=plugins=grpc:. proto/content.proto

build: proto
	go build -o content-srv main.go

test:
	go test -v ./... -cover

docker: build
	docker build . -t content-srv:latest

docker-tag: docker
	docker tag content-srv:latest ${DOCKER_ID_USER}/content-srv:latest

docker-push: docker-tag
	docker push ${DOCKER_ID_USER}/content-srv:latest