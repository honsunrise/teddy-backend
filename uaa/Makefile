
GOPATH:=$(shell go env GOPATH)

.PHONY: proto test docker


proto:
	protoc --proto_path=${GOPATH}/src:. --go_out=plugins=grpc:. proto/uaa.proto

build: proto
	CGO_ENABLED=1 GOOS=linux GOARCH=amd64 go build --ldflags "-extldflags -static" -o uaa-srv main.go config.go

test:
	go test -v ./... -cover

docker: build
	docker build . -t uaa-srv:latest

docker-tag: docker
	docker tag uaa-srv:latest ${DOCKER_ID_USER}/uaa-srv:latest

docker-push: docker-tag
	docker push ${DOCKER_ID_USER}/uaa-srv:latest