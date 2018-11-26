FROM alpine:3.2
RUN mkdir /app
WORKDIR /app
COPY base-api base-api
RUN chmod +x base-api
ENTRYPOINT [ "/app/base-api" ]
