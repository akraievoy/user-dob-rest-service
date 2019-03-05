#!/bin/sh

set -o errexit

echo "STEP 1: creating/starting dockers"

if docker-compose up --abort-on-container-exit --exit-code-from verifier ; then
  echo "END TO END PASSED"
else
  echo "END TO END FAILED"
fi
