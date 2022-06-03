# syntax=docker/dockerfile:1.2
FROM golang:1.18-alpine as builder_img_

COPY services/common /go/src/common
COPY services/scraper/ /go/src/scraper

WORKDIR /go/src/scraper

RUN apk add --no-cache --update gcc g++


RUN --mount=type=cache,target=/go/pkg \
      go mod download

COPY . .

RUN --mount=type=cache,target=/go/pkg \
    CGO_ENABLED=1 go build -v -o /go/bin/scraper .

FROM alpine:latest

WORKDIR /app

RUN apk add --no-cache --update tzdata

COPY --from=builder_img_ /go/bin /app

ENTRYPOINT /app/scraper