#!/usr/bin/env sh

DB_FILE="sqlite-cli.db"
DB_MIGRATION="migration.sql"

touch "$DB_FILE"
sqlite3 "$DB_FILE" < "$DB_MIGRATION"