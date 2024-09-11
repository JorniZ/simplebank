set -e

echo "run db migration"
/simple-bank-api/migrate -path /simple-bank-api/migration -datebase "$DB_SOURCE" -verbose up

echo "start the app"
exec "$@"