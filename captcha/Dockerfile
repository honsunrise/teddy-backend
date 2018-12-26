# Stage 1 (to create a "build" image)
FROM golang:1.11.2-alpine AS builder
RUN go version
RUN apk update && \
    apk add protobuf && \
    apk add git && \
    apk add openssh-client
RUN TOKEN=17cae3fbe874045feba5f8c90703ebc922a89e64 && \
    git config --global url."https://${TOKEN}@github.com/".insteadOf "https://github.com/"
COPY . /go/src/github.com/zhsyourai/teddy-backend/captcha
WORKDIR /go/src/github.com/zhsyourai/teddy-backend/captcha
RUN set -x && \
    go get github.com/golang/dep/cmd/dep && \
    dep ensure -v
RUN set -x && \
    GIT_TAG="v1.2.0" && \
    go get -d -u github.com/golang/protobuf/protoc-gen-go && \
    git -C "$(go env GOPATH)"/src/github.com/golang/protobuf checkout $GIT_TAG && \
    go install github.com/golang/protobuf/protoc-gen-go
RUN GOPATH=$(go env GOPATH) && \
    set -x && \
    protoc --proto_path=${GOPATH}/src:./vendor \
        --go_out=plugins=grpc:./vendor \
        ./vendor/github.com/zhsyourai/teddy-backend/common/proto/captcha/captcha.proto
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -a -o main .

# Stage 2 (to create a downsized "container executable")
FROM alpine:3.2
RUN apk update && \
    apk add ca-certificates wget && \
    update-ca-certificates
RUN GRPC_HEALTH_PROBE_VERSION=v0.2.0 && \
    wget -qO/bin/grpc_health_probe https://github.com/grpc-ecosystem/grpc-health-probe/releases/download/${GRPC_HEALTH_PROBE_VERSION}/grpc_health_probe-linux-amd64 && \
    chmod +x /bin/grpc_health_probe
RUN mkdir /app
WORKDIR /app
COPY --from=builder /go/src/github.com/zhsyourai/teddy-backend/captcha/main .
ENTRYPOINT [ "/app/main" ]
