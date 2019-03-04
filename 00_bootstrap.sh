#!/bin/bash

echo "you might pass 'force' as first param to this script to force dependency/protobuf re-fetch"
sleep 1

docker run --rm -it \
  -v "$(pwd)/go_pkg_mod:/go/pkg/mod/" \
  -v "$(pwd):/tsv_load/" \
  golang:1.12.0-alpine3.9 \
  /bin/sh '/tsv_load/00_bootstrap_docker.sh' "$1"

echo "docker may have generated some root-owned files -- reclaim ownership?"
sudo chown -R "$(whoami):$(whoami)" .
