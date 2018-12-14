
GOPATH:=$(shell go env GOPATH)

.PHONY: proto test docker

proto:
	protoc --proto_path=${GOPATH}/src \
	    --go_out=plugins=grpc:${GOPATH}/src \
	    ${GOPATH}/src/github.com/zhsyourai/teddy-backend/common/proto/content.proto

test:
	go test -v ./... -cover

docker:
	docker build . -t srv-content:latest

docker-tag: docker
	docker tag srv-content:latest ${DOCKER_ID_USER}/srv-content:latest

docker-push: docker-tag
	docker push ${DOCKER_ID_USER}/srv-content:latest