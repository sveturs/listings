#!/bin/bash
set -e

echo "Running reindexing script..."

# Try reindexing all listings with translations
echo "Calling reindexing endpoint..."
curl -X POST http://localhost:4000/api/v1/admin/reindex-listings-with-translations \
  -H "Content-Type: application/json" \
  -H "Cookie: session=YOUR_SESSION_COOKIE" || true

echo "Done"