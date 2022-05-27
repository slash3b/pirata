# syntax=docker/dockerfile:1.2
FROM golang:1.18-alpine as builder_img_

RUN apk add --no-cache --update gcc g++

WORKDIR /go/src

COPY go.mod go.sum ./

RUN --mount=type=cache,target=/go/pkg \
      go mod download && go mod verify

COPY . .

RUN --mount=type=cache,target=/go/pkg \
    CGO_ENABLED=1 go build -v -o /go/bin/api .

FROM alpine:latest

WORKDIR /app

RUN apk add --no-cache --update tzdata

COPY --from=builder_img_ /go/bin /app

ENTRYPOINT /app/api