#!/bin/sh

set -e

echo "run db migration"
source /app/app.env
/app/soda migrate up -p /db -c /db/database.yml

echo "start the application"
exec "$@"