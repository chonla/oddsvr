FROM golang:1.11.2-alpine

WORKDIR /go/src/github.com/chonla/oddsvr

COPY . .

RUN apk add --no-cache git \
    && go get ./...

