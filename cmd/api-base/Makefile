
GOPATH:=$(shell go env GOPATH)

.PHONY: test docker

docker:
	docker build . -t api-base:latest

docker-tag: docker
	docker tag base-api:latest ${DOCKER_ID_USER}/api-base:latest

docker-push: docker-tag
	docker push ${DOCKER_ID_USER}/api-base:latest
