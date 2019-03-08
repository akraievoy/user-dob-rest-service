#!/bin/sh

set -o errexit

if [ "z$1" != "zforce" ]  && test -d /tsv_load/vendor/ ; then
  echo libs already vendored, skipping
else

  echo "STEP 3: generating protobuf golang API (from go_build docker)"
  cd /tsv_load/
  protoc \
    --plugin=/go/bin/protoc-gen-go \
    --go_out=plugins=grpc:proto \
    -I proto \
    proto/upserter.proto

  echo "STEP 4: vendoring library dependencies (from go_build docker)"
  go mod vendor
  go get ./parser ./ingestor ./upserter

fi
