
GOPATH:=$(shell go env GOPATH)

.PHONY: proto test docker


proto:
	protoc --proto_path=${GOPATH}/src:. --go_out=plugins=grpc:. proto/message.proto

build: proto
	CGO_ENABLED=1 GOOS=linux GOARCH=amd64 go build --ldflags "-extldflags -static" -o message-srv main.go config.go

test:
	go test -v ./... -cover

docker: build
	docker build . -t message-srv:latest

docker-tag: docker
	docker tag message-srv:latest ${DOCKER_ID_USER}/message-srv:latest

docker-push: docker-tag
	docker push ${DOCKER_ID_USER}/message-srv:latest