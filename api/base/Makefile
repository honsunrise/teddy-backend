
GOPATH:=$(shell go env GOPATH)

.PHONY: test docker

build:
	go build -o base-api main.go

test:
	go test -v ./... -cover

docker: build
	docker build . -t base-api:latest

docker-tag: docker
	docker tag base-api:latest ${DOCKER_ID_USER}/base-api:latest

docker-push: docker-tag
	docker push ${DOCKER_ID_USER}/base-api:latest
