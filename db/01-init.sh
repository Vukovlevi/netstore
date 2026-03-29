#!/bin/bash
set -e

echo "DB_TYPE: $DB_TYPE"

# Default fallback
DB_TYPE=${DB_TYPE:-szerkezet}

if [ "$DB_TYPE" = "adatok" ]; then
  echo "Using data.sql"
  rm -f /docker-entrypoint-initdb.d/netstore_szerkezet.sql
else
  echo "Using structure.sql"
  rm -f /docker-entrypoint-initdb.d/netstore_adatok.sql
fi
