
GOPATH:=$(shell go env GOPATH)

.PHONY: test docker

test:
	go test -v ./... -cover

docker:
	docker build . -t api-uaa:latest

docker-tag: docker
	docker tag api-uaa:latest ${DOCKER_ID_USER}/api-uaa:latest

docker-push: docker-tag
	docker push ${DOCKER_ID_USER}/api-uaa:latest
