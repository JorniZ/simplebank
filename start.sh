#!/bin/sh

set -e

echo "run db migration"
source /simple-bank-api/app.env
/simple-bank-api/migrate -path /simple-bank-api/migration -database "$DB_SOURCE" -verbose up

echo "start the app"
exec "$@"