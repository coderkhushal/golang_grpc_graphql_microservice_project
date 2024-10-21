
#!/bin/bash

# Start PostgreSQL in the background
/usr/local/bin/docker-entrypoint.sh postgres &

# Wait until PostgreSQL is ready
until pg_isready -h localhost -U "$POSTGRES_USER"; do
  echo 'Waiting for PostgreSQL to be ready...'
  sleep 2
done

# Run the migration script
echo 'Running migrations...'
PGPASSWORD="$POSTGRES_PASSWORD" psql -U "$POSTGRES_USER" -d "$POSTGRES_DB" -f /docker-entrypoint-initdb.d/1.sql

# Keep the PostgreSQL server running
wait
