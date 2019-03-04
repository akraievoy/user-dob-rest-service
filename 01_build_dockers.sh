#!/bin/sh

set -o errexit

echo "STEP 1: compiling sender"

docker run --rm -it \
  -v "$(pwd)/go_pkg_mod/:/go/pkg/mod/" \
  -v "$(pwd):/tsv_load/" \
  golang:1.12.0-alpine3.9 \
  /bin/sh -c 'cd /tsv_load/sender && go build .'

echo "STEP 2: building sender docker"

docker build -t sender -f sender.Dockerfile .


echo "STEP 3: compiling upserter"

docker run --rm -it \
  -v "$(pwd)/go_pkg_mod/:/go/pkg/mod/" \
  -v "$(pwd):/tsv_load/" \
  golang:1.12.0-alpine3.9 \
  /bin/sh -c 'cd /tsv_load/upserter && go build .'

echo "STEP 4: building upserter docker"

docker build -t upserter -f upserter.Dockerfile .
