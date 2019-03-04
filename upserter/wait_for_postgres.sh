#!/bin/sh

set -e

# PGPASSWORD=${DB_PASSWORD}
until psql "postgres://${DB_USERNAME}:${DB_PASSWORD}@${DB_HOST}:${DB_PORT}/${DB_NAME}?sslmode=disable" -c '\q'; do
  >&2 echo "Postgres is unavailable - sleeping"
  sleep 5
done

COMMAND="$@"
>&2 echo "Postgres is up - executing command ${COMMAND}"
/bin/sh -c "$COMMAND"