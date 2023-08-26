#!/bin/sh

set -e

echo "run db migration"
cp /app/app.env .env
/app/soda migrate up -p /db -c /db/database.yml

echo "start the application"
exec "$@"