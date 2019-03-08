#!/bin/bash

echo "you might pass 'force' as first param to this script to force dependency/protobuf re-fetch"
sleep 1

echo "STEP 1 building go_build docker"

docker build -t go_build go_build

echo "STEP 2 starting go_build"

docker run --rm -it \
  -v "$(pwd)/go_pkg_mod:/go/pkg/mod/" \
  -v "$(pwd):/tsv_load/" \
  go_build \
  /bin/sh '/tsv_load/00_bootstrap_docker.sh' "$1"

echo "docker may have generated some root-owned files -- reclaim ownership?"
sudo chown -R "$(whoami):$(whoami)" .
