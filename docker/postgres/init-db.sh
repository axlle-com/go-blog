#!/bin/bash
set -e

# Создание таблицы user
psql -v ON_ERROR_STOP=1 --username "$POSTGRES_USER" --dbname "$POSTGRES_DB" <<-EOSQL
    CREATE TABLE user (
        id SERIAL PRIMARY KEY,
        name VARCHAR(100) NOT NULL
    );
EOSQL
