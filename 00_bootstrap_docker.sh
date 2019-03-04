#!/bin/sh

set -o errexit

if "$1" != "force" && test -d /tsv_load/vendor ; then
  echo libs already vendored, skipping
else
  echo "STEP 1: installing git and protobuf (in golang docker)"
  apk update
  apk upgrade
  apk add git protobuf

  echo "STEP 2: installing protoc-gen-go"
  go get -u github.com/golang/protobuf/protoc-gen-go

  cd /tsv_load/

  echo "STEP 3: generating protobuf golang API"
  protoc \
    --plugin=/go/bin/protoc-gen-go \
    --go_out=plugins=grpc:proto \
    -I proto \
    proto/upserter.proto

  echo "STEP 4: vendoring library dependencies"
  go mod vendor
  go get ./parser ./sender ./upserter

fi
