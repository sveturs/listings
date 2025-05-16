#!/bin/bash

# Script to reset production database and apply all migrations
# WARNING: This will delete all data in the database!

set -e

echo "⚠️  WARNING: This script will reset the production database!"
echo "All data will be lost. Are you sure you want to continue?"
read -p "Type 'yes' to continue: " confirmation

if [ "$confirmation" != "yes" ]; then
    echo "Aborted."
    exit 1
fi

# SSH connection details
SSH_HOST="root@svetu.rs"
REMOTE_DIR="/opt/hostel-booking-system"

echo "⏹️  Stopping backend to close database connections..."
ssh $SSH_HOST "cd $REMOTE_DIR && docker compose stop backend"

echo "🔌 Terminating remaining database connections..."
ssh $SSH_HOST "cd $REMOTE_DIR && docker exec hostel_db psql -U postgres -c \"SELECT pg_terminate_backend(pid) FROM pg_stat_activity WHERE datname = 'hostel_db' AND pid <> pg_backend_pid();\""

echo "🗑️  Dropping and recreating database..."
ssh $SSH_HOST "cd $REMOTE_DIR && docker exec hostel_db psql -U postgres -c 'DROP DATABASE IF EXISTS hostel_db;'"
ssh $SSH_HOST "cd $REMOTE_DIR && docker exec hostel_db psql -U postgres -c 'CREATE DATABASE hostel_db;'"

echo "📝 Running migrations..."
ssh $SSH_HOST "cd $REMOTE_DIR && docker compose run --rm migrate"

echo "🔄 Starting backend..."
ssh $SSH_HOST "cd $REMOTE_DIR && docker compose start backend"

echo "✅ Production database reset complete!"
echo "Don't forget to test all functionality"