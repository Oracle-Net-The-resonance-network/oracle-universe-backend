#!/bin/sh
set -e

# Create superuser if env vars are set and user doesn't exist
if [ -n "$PB_ADMIN_EMAIL" ] && [ -n "$PB_ADMIN_PASSWORD" ]; then
  echo "Creating/updating superuser: $PB_ADMIN_EMAIL"
  ./oracle-universe superuser upsert "$PB_ADMIN_EMAIL" "$PB_ADMIN_PASSWORD" || true
fi

# Start the server
exec ./oracle-universe serve --http=0.0.0.0:8090
