#!/bin/bash

REQUIRED_VARS=(POSTGRES_USER POSTGRES_PASSWORD POSTGRES_DB POSTGRES_PORT POSTGRES_DB POSTGRES_HOST)
for var in "${REQUIRED_VARS[@]}"; do
  if [ -z "${!var}" ]; then
    echo "Error: Environment variable $var is not set."
    exit 1
  fi
done

SSL_MODE="${SSL_MODE:-disable}"
DB_URL="postgresql://${POSTGRES_USER}:${POSTGRES_PASSWORD}@${POSTGRES_HOST}:${POSTGRES_PORT}/${POSTGRES_DB}?sslmode=${SSL_MODE}"

MIGRATION_DIR="./migrations"

echo "Running migrations..."
goose -dir "${MIGRATION_DIR}" postgres "${DB_URL}" up
if [ $? -eq 0 ]; then
  echo "Migrations applied successfully."
else
  echo "Failed to apply migrations."
  exit 1
fi
