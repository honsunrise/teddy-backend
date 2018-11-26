FROM alpine:3.2
RUN mkdir /app
WORKDIR /app
COPY uaa-api uaa-api
RUN chmod +x uaa-api
ENTRYPOINT [ "/app/uaa-api" ]
