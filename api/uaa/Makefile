
GOPATH:=$(shell go env GOPATH)

.PHONY: test docker

build:
	CGO_ENABLED=1 GOOS=linux GOARCH=amd64 go build --ldflags "-extldflags -static" -o uaa-api main.go config.go

test:
	go test -v ./... -cover

docker: build
	docker build . -t uaa-api:latest

docker-tag: docker
	docker tag uaa-api:latest ${DOCKER_ID_USER}/uaa-api:latest

docker-push: docker-tag
	docker push ${DOCKER_ID_USER}/uaa-api:latest
