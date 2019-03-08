#!/bin/sh

set -o errexit

echo "STEP 1: compiling ingestor"

docker run --rm -it \
  -v "$(pwd)/go_pkg_mod/:/go/pkg/mod/" \
  -v "$(pwd):/tsv_load/" \
  go_build \
  /bin/sh -c 'cd /tsv_load/ingestor && go build .'

echo "STEP 2: compiling upserter"

docker run --rm -it \
  -v "$(pwd)/go_pkg_mod/:/go/pkg/mod/" \
  -v "$(pwd):/tsv_load/" \
  go_build \
  /bin/sh -c 'cd /tsv_load/upserter && go build .'

echo "STEP 3: purging previously built dockers"

docker-compose rm -f

echo "STEP 4: purging previously built and now untagged docker images"

docker images | grep '<none>' | sed -E 's/\s+/\t/g' | cut -f 3 | xargs -r docker rmi -f
docker system prune --volumes --force

echo "STEP 5: building dockers"

docker-compose build
