#!/bin/sh
set -e

# Restore from R2 if no local database exists
if [ ! -f /app/pb_data/data.db ]; then
  echo "No local database found. Restoring from R2..."
  litestream restore -if-replica-exists -config /etc/litestream.yml /app/pb_data/data.db
fi

# Create superuser if env vars are set
if [ -n "$PB_ADMIN_EMAIL" ] && [ -n "$PB_ADMIN_PASSWORD" ]; then
  echo "Creating/updating superuser: $PB_ADMIN_EMAIL"
  ./oracle-universe superuser upsert "$PB_ADMIN_EMAIL" "$PB_ADMIN_PASSWORD" || true
fi

# Litestream as parent process, PocketBase as child
# Litestream continuously replicates WAL changes to R2
exec litestream replicate -config /etc/litestream.yml \
  -exec "./oracle-universe serve --http=0.0.0.0:8090"
