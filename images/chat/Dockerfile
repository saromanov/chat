# Image page: <https://hub.docker.com/_/golang>

FROM golang:1.14-alpine as builder
COPY . /src
WORKDIR /src
COPY ./images/chat/entrypoint.sh /tmp/entrypoint.sh
RUN go build -ldflags="-s -w" -o /tmp/app ./cmd/.

FROM alpine:latest
RUN goose
COPY --from=builder /tmp/app /bin/app
COPY --from=builder /tmp/entrypoint.sh /bin/entrypoint.sh
RUN chmod +x /bin/entrypoint.sh
CMD ["/bin/entrypoint.sh"]