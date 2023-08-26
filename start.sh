#!/bin/sh

set -e

echo "run db migration"
. /app/app.env
/app/soda migrate up -p /db -c /db/database.yml -e "production"

echo "start the application"
exec "$@"