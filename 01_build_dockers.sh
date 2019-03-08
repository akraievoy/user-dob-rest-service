#!/bin/sh

set -o errexit

echo "STEP 1: go vet all local sources"

docker run --rm -it \
  -v "$(pwd)/go_pkg_mod/:/go/pkg/mod/" \
  -v "$(pwd):/tsv_load/" \
  go_build \
  /bin/sh -c 'cd /tsv_load/ && go vet ./ingestor ./parser ./service_utils ./upserter ./verifier'

echo "STEP 2: compiling ingestor"

docker run --rm -it \
  -v "$(pwd)/go_pkg_mod/:/go/pkg/mod/" \
  -v "$(pwd):/tsv_load/" \
  go_build \
  /bin/sh -c 'cd /tsv_load/ingestor && go test . && go build .'

echo "STEP 3: compiling upserter"

docker run --rm -it \
  -v "$(pwd)/go_pkg_mod/:/go/pkg/mod/" \
  -v "$(pwd):/tsv_load/" \
  go_build \
  /bin/sh -c 'cd /tsv_load/upserter && go test . && go build .'

echo "STEP 4: purging previously built and now untagged docker images"

docker-compose rm -f
docker images | grep '<none>' | sed -E 's/\s+/\t/g' | cut -f 3 | xargs -r docker rmi -f
docker system prune --volumes --force

echo "STEP 5: building dockers"

docker-compose build
