FROM golang:1.12.0-alpine3.9

RUN apk update && \
    apk upgrade && \
    apk add build-base git protobuf && \
    go get -u github.com/golang/protobuf/protoc-gen-go

