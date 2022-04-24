# syntax=docker/dockerfile:1.2
FROM golang:1.18-alpine as builder_img_

RUN apk add --no-cache --update gcc g++

WORKDIR /go/src

COPY go.mod go.sum ./

RUN --mount=type=cache,target=/go/pkg \
      go mod download && go mod verify

COPY . .

RUN --mount=type=cache,target=/go/pkg \
    go build -v -o /go/bin/scraper .

FROM alpine

WORKDIR /app

COPY --from=builder_img_ /go/bin/scraper /app

ENTRYPOINT "./scraper"