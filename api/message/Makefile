
GOPATH:=$(shell go env GOPATH)

.PHONY: test docker

test:
	go test -v ./... -cover

docker:
	docker build . -t api-message:latest

docker-tag: docker
	docker tag api-message:latest ${DOCKER_ID_USER}/api-message:latest

docker-push: docker-tag
	docker push ${DOCKER_ID_USER}/api-message:latest
