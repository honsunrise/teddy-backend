FROM alpine:3.2
RUN mkdir /app
WORKDIR /app
COPY message-api message-api
RUN chmod +x message-api
ENTRYPOINT [ "/app/message-api" ]
