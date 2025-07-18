#!/bin/bash

export PGPASSWORD=password

echo "=== Checking city and country columns in marketplace_listings ==="
psql -U postgres -d hostel_db -h localhost -t -c "SELECT column_name, data_type FROM information_schema.columns WHERE table_name = 'marketplace_listings' AND column_name IN ('city', 'country') ORDER BY column_name" 2>&1

echo -e "\n=== Checking migrations >= 90 ==="
psql -U postgres -d hostel_db -h localhost -t -c "SELECT id, name, CASE WHEN applied_at IS NOT NULL THEN 'applied' ELSE 'not applied' END as status FROM migrations WHERE id >= 90 ORDER BY id" 2>&1