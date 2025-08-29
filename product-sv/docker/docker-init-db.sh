#!/bin/bash

# WARNING! This script is for local development only.

set -e

echo "Step 1. Running initial script to create database and user (ROLE) with granted access to it"

psql -v ON_ERROR_STOP=1 --username "$POSTGRES_USER" --dbname "$POSTGRES_DB" <<-EOSQL
    -- Create role if it doesn't exist
    DO \$\$
    BEGIN
        IF NOT EXISTS (SELECT FROM pg_roles WHERE rolname = '${POSTGRES_USER}') THEN
            CREATE ROLE "${POSTGRES_USER}" WITH LOGIN PASSWORD '${POSTGRES_PASSWORD}';
        END IF;
    END
    \$\$;

    -- Create database and grant access
    DO \$\$
    BEGIN
        IF NOT EXISTS (SELECT FROM pg_database WHERE datname = '${POSTGRES_DB}') THEN
            EXECUTE format('CREATE DATABASE "%I";', '${POSTGRES_DB}');
        END IF;
    END
    \$\$;

    GRANT ALL PRIVILEGES ON DATABASE "${POSTGRES_DB}" TO "${POSTGRES_USER}";
EOSQL

echo "Step 2. Running initial schema migration. Migrations SQL files located in docker-entrypoint-initdb.d/migrations ..."

for f in /docker-entrypoint-initdb.d/migrations/*.sql; do
  echo "Applying $f..."
  psql -v ON_ERROR_STOP=1 --username "$POSTGRES_USER" --dbname "$POSTGRES_DB" -f "$f"
done
