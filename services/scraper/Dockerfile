FROM golang:1.18-alpine as builder_img_
RUN apk add --no-cache --update gcc g++
WORKDIR /go/src
COPY go.mod go.sum ./
RUN go mod download && go mod verify
COPY . .
RUN go build -v -o /go/bin/scraper .

FROM alpine:latest
COPY --from=builder_img_ /go/bin/scraper .
CMD ["./scraper"]