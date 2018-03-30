#!/bin/sh

set -e
set -o
set -x

export PACKAGE=$1

go get github.com/golang/lint/golint

cd /go/src/${PACKAGE}

go get -v -t ./...

golint -set_exit_status

go vet .

go test .

go build -o main ./*.go
