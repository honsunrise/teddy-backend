
GOPATH:=$(shell go env GOPATH)

.PHONY: proto test docker

proto:
	protoc --proto_path=${GOPATH}/src \
	    --go_out=plugins=grpc:${GOPATH}/src \
	    ${GOPATH}/src/github.com/zhsyourai/teddy-backend/common/proto/message/message.proto

test:
	go test -v ./... -cover

docker:
	docker build . -t srv-message:latest

docker-tag: docker
	docker tag srv-message:latest ${DOCKER_ID_USER}/srv-message:latest

docker-push: docker-tag
	docker push ${DOCKER_ID_USER}/srv-message:latest