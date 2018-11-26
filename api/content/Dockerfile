FROM alpine:3.2
RUN mkdir /app
WORKDIR /app
COPY content-api content-api
RUN chmod +x content-api
ENTRYPOINT [ "/app/content-api" ]
