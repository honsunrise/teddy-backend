
GOPATH:=$(shell go env GOPATH)

.PHONY: test docker

build:
	go build -o message-api main.go

test:
	go test -v ./... -cover

docker: build
	docker build . -t message-api:latest

docker-tag: docker
	docker tag message-api:latest ${DOCKER_ID_USER}/message-api:latest

docker-push: docker-tag
	docker push ${DOCKER_ID_USER}/message-api:latest
