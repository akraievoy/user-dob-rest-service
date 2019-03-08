#!/bin/sh

set -o errexit

echo "STEP 1: clean"

docker-compose rm -f
docker images | grep '<none>' | sed -E 's/\s+/\t/g' | cut -f 3 | xargs -r docker rmi -f
docker system prune --volumes --force

echo "STEP 2: creating/starting dockers"

if docker-compose up --abort-on-container-exit --exit-code-from verifier ; then
  echo "END TO END PASSED"
else
  echo "END TO END FAILED"
fi
