# syntax=docker/dockerfile:1.2
FROM golang:1.18-alpine as builder_img_

RUN apk add --no-cache --update gcc g++

COPY services/common /go/src/common
COPY services/api/ /go/src/api

WORKDIR /go/src/api

RUN --mount=type=cache,target=/go/pkg \
    go mod download

RUN --mount=type=cache,target=/go/pkg \
    CGO_ENABLED=1 go build -v -o /go/bin/api .

FROM alpine:latest

WORKDIR /app

RUN apk add --no-cache --update tzdata

COPY --from=builder_img_ /go/bin /app

ENTRYPOINT /app/api