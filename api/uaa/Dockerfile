# Stage 1 (to create a "build" image)
FROM golang:1.11.2-alpine AS builder
RUN go version
RUN apk update && \
    apk add protobuf && \
    apk add git && \
    apk add openssh-client
RUN TOKEN=17cae3fbe874045feba5f8c90703ebc922a89e64 && \
    git config --global url."https://${TOKEN}@github.com/".insteadOf "https://github.com/"
COPY . /go/src/github.com/zhsyourai/teddy-backend/api/uaa
WORKDIR /go/src/github.com/zhsyourai/teddy-backend/api/uaa
RUN set -x && \
    GIT_TAG="v1.2.0" && \
    go get -d -u github.com/golang/protobuf/protoc-gen-go && \
    git -C "$(go env GOPATH)"/src/github.com/golang/protobuf checkout $GIT_TAG && \
    go install github.com/golang/protobuf/protoc-gen-go
RUN set -x && \
    go get github.com/golang/dep/cmd/dep && \
    dep ensure -v
RUN GOPATH=$(go env GOPATH) && \
    set -x && \
    protoc --proto_path=${GOPATH}/src:./vendor \
        --go_out=plugins=grpc:./vendor \
        ./vendor/github.com/zhsyourai/teddy-backend/common/proto/captcha.proto && \
    protoc --proto_path=${GOPATH}/src:./vendor \
        --go_out=plugins=grpc:./vendor \
        ./vendor/github.com/zhsyourai/teddy-backend/common/proto/content.proto && \
    protoc --proto_path=${GOPATH}/src:./vendor \
        --go_out=plugins=grpc:./vendor \
        ./vendor/github.com/zhsyourai/teddy-backend/common/proto/message.proto && \
    protoc --proto_path=${GOPATH}/src:./vendor \
        --go_out=plugins=grpc:./vendor \
        ./vendor/github.com/zhsyourai/teddy-backend/common/proto/uaa.proto
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -a -o main .

# Stage 2 (to create a downsized "container executable")
FROM alpine:3.2
RUN mkdir /app
WORKDIR /app
COPY --from=builder /go/src/github.com/zhsyourai/teddy-backend/api/uaa/main .
ENTRYPOINT [ "/app/main" ]
