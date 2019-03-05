#!/bin/bash

sudo chown -R "${whoami}:${whoami}" .
chmod -R 770 .idea go_pkg_mod vendor
rm -rf .idea go_pkg_mod vendor ingestor/ingestor upserter/upserter
git -c gc.reflogExpire=0 -c gc.reflogExpireUnreachable=0 -c gc.rerereresolved=0 -c gc.rerereunresolved=0 -c gc.pruneExpire=now gc
docker-compose rm -f
docker images | grep '<none>' | sed -E 's/\s+/\t/g' | cut -f 3 | xargs -r docker rmi -f
docker system prune --volumes --force

