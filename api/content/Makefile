
GOPATH:=$(shell go env GOPATH)

.PHONY: test docker

build:
	CGO_ENABLED=1 GOOS=linux GOARCH=amd64 go build --ldflags "-extldflags -static" -o content-api main.go config.go

test:
	go test -v ./... -cover

docker: build
	docker build . -t content-api:latest

docker-tag: docker
	docker tag content-api:latest ${DOCKER_ID_USER}/content-api:latest

docker-push: docker-tag
	docker push ${DOCKER_ID_USER}/content-api:latest
