#!/bin/bash

echo "Stopping containers gracefully..."
docker-compose down -t 1 || true

echo "Force killing containers..."
running_containers=$(docker ps -q)
if [ ! -z "$running_containers" ]; then
    docker kill $running_containers || true
fi

echo "Removing stopped containers..."
stopped_containers=$(docker ps -aq)
if [ ! -z "$stopped_containers" ]; then
    docker rm $stopped_containers || true
fi

echo "Cleaning up networks..."
docker network prune -f || true


echo "Cleanup complete"