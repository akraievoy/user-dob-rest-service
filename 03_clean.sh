#!/bin/bash

sudo chown -R "${whoami}:${whoami}" .
chmod -R 770 .idea go_pkg_mod vendor
rm -rf .idea go_pkg_mod vendor ingestor/ingestor upserter/upserter
git gc
docker-compose rm -f
docker images | grep '<none>' | sed -E 's/\s+/\t/g' | cut -f 3 | xargs -r docker rmi -f
docker system prune --volumes --force
