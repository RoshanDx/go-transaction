#!/bin/bash

set -e

current_time=$(date +%s)

POD_NAME="test-pod"
POD_NETWORK="test-network"
POD_DB_DATA="testdb-data"

# Check if .env file exists
if [ -f .env ]; then
    source .env
else
    echo ".env file not found!"
    exit 1
fi

if ! podman pod exists $POD_NAME; then
  echo "pod '$POD_NAME' does not exist, creating pod..."
  podman pod create --name $POD_NAME
else
  echo "pod '$POD_NAME' already exists."
fi

if ! podman network exists $POD_NETWORK; then
  echo "pod network '$POD_NETWORK' does not exist, creating network..."
  podman network create $POD_NETWORK
else
  echo "pod network '$POD_NETWORK' already exists."
fi

if ! podman volume exists $POD_DB_DATA; then
  echo "pod db data '$POD_DB_DATA' does not exist, creating data..."
  podman volume create $POD_DB_DATA
else
  echo "pod db data '$POD_DB_DATA' already exists."
fi

echo "running postgres container..."
podman run -d \
  --name testdb \
  --pod $POD_NAME \
  --restart always \
  -e POSTGRES_USER="$DATABASE_USERNAME" \
  -e POSTGRES_PASSWORD="$DATABASE_PASSWORD" \
  -e POSTGRES_DB="$DATABASE_NAME" \
  --network $POD_NETWORK \
  -v $POD_DB_DATA:/var/lib/postgresql/data \
  -p 13308:"$DATABASE_PORT" \
  postgres:16

until podman exec testdb pg_isready -U "$DATABASE_USERNAME" -d "$DATABASE_DATABASE"; do
  echo "waiting for database to be ready..."
  sleep 5
done

end_time=$(date +%s)
time_difference=$((end_time - current_time))
echo "script complete in $time_difference seconds"