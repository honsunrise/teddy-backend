
GOPATH:=$(shell go env GOPATH)

.PHONY: proto test docker


proto:
	protoc --proto_path=${GOPATH}/src:. --go_out=plugins=grpc:. proto/captcha.proto

build: proto
	CGO_ENABLED=1 GOOS=linux GOARCH=amd64 go build --ldflags "-extldflags -static" -o captcha-srv main.go config.go

test:
	go test -v ./... -cover

docker: build
	docker build . -t captcha-srv:latest

docker-tag: docker
	docker tag captcha-srv:latest ${DOCKER_ID_USER}/captcha-srv:latest

docker-push: docker-tag
	docker push ${DOCKER_ID_USER}/captcha-srv:latest
