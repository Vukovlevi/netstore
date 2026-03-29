#!/bin/bash
set -e

DB_TYPE=${DB_TYPE:-szerkezet}

echo "Initializing database with: $DB_TYPE"

if [ "$DB_TYPE" = "adatok" ]; then
  for f in /db/adatok/*.sql; do
    echo "Running $f"
    mysql -u root -p"$MYSQL_ROOT_PASSWORD" < "$f"
  done
else
  for f in /db/szerkezet/*.sql; do
    echo "Running $f"
    mysql -u root -p"$MYSQL_ROOT_PASSWORD" < "$f"
  done
fi
