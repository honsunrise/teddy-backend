FROM alpine:3.2
RUN apk update && \
    apk add ca-certificates wget && \
    update-ca-certificates
RUN GRPC_HEALTH_PROBE_VERSION=v0.2.0 && \
    wget -qO/bin/grpc_health_probe https://github.com/grpc-ecosystem/grpc-health-probe/releases/download/${GRPC_HEALTH_PROBE_VERSION}/grpc_health_probe-linux-amd64 && \
    chmod +x /bin/grpc_health_probe
RUN mkdir /app
WORKDIR /app
COPY content-srv content-srv
RUN chmod +x content-srv
ENTRYPOINT [ "/app/content-srv" ]
