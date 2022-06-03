# syntax=docker/dockerfile:1.2
FROM golang:1.18-alpine as builder_img_

COPY services/common /go/src/common
COPY services/imdb/ /go/src/imdb

WORKDIR /go/src/imdb

RUN --mount=type=cache,target=/go/pkg \
    go mod download

# go mod verify is a bit broken. It fails with common v0.0.1: missing ziphash: open hash: no such file or directory
# it looks like it is a tooling issue

RUN --mount=type=cache,target=/go/pkg \
    go build -v -o /go/bin/imdb .

FROM alpine:latest

WORKDIR /app

COPY --from=builder_img_ /go/bin/imdb /app/

ENTRYPOINT /app/imdb