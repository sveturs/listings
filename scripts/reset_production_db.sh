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

echo "🔄 Updating code on production server..."
ssh $SSH_HOST "cd $REMOTE_DIR && git pull"

echo "🗑️  Dropping and recreating database..."
ssh $SSH_HOST "cd $REMOTE_DIR && docker exec hostel_db psql -U postgres -c 'DROP DATABASE IF EXISTS hostel_db;'"
ssh $SSH_HOST "cd $REMOTE_DIR && docker exec hostel_db psql -U postgres -c 'CREATE DATABASE hostel_db;'"

echo "📝 Running migrations..."
ssh $SSH_HOST "cd $REMOTE_DIR && docker compose run --rm migrate"

echo "🔄 Restarting services..."
ssh $SSH_HOST "cd $REMOTE_DIR && docker compose restart backend"

echo "✅ Production database reset complete!"
echo "Don't forget to test all functionality"