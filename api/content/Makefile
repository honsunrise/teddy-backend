
GOPATH:=$(shell go env GOPATH)

.PHONY: test docker

test:
	go test -v ./... -cover

docker:
	docker build . -t api-content:latest

docker-tag: docker
	docker tag api-content:latest ${DOCKER_ID_USER}/api-content:latest

docker-push: docker-tag
	docker push ${DOCKER_ID_USER}/api-content:latest
