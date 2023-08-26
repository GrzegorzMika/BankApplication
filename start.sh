#!/bin/sh

set -e

echo "run db migration"
. /app/app.env
/app/soda migrate up -p /db -c /db/database.yml

echo "start the application"
exec "$@"