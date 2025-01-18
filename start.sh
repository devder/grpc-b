#!/bin/sh

# Exit immediately if a command exits with a non-zero status
set -e

# temporary remove it to run the db migrations in the code itself
# echo "run db migration"
# # load the environment variables
# source /usr/src/app/app.env
# /usr/src/app/migrate -path db/migrations -database "$DB_SOURCE" -verbose up

echo "start the app"
# means run the command that is passed as an argument to this script
# in this case the CMD from the Dockerfile
exec "$@" 