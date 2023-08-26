#!/bin/sh

set -e

echo "run db migration"
. /app/app.env
echo $GO_ENV
echo $DB_SOURCE
/app/soda migrate up -p /db -c /db/database.yml

echo "start the application"
exec "$@"