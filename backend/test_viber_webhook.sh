#!/bin/bash

echo "🤖 Testing Viber Bot Webhook for Tracking"
echo "========================================="
echo ""

# Send the webhook request
echo "📤 Sending track command: track_TEST-TOKEN-123"
echo ""

response=$(curl -s -X POST http://localhost:3000/api/viber/webhook \
  -H "Content-Type: application/json" \
  -H "X-Viber-Content-Signature: dev" \
  -d '{
    "event": "message",
    "timestamp": 1234567890,
    "message_token": "123456789",
    "sender": {
      "id": "test-viber-id",
      "name": "Test User"
    },
    "message": {
      "type": "text",
      "text": "track_TEST-TOKEN-123"
    }
  }' 2>&1)

echo "📥 Webhook Response:"
echo "$response"
echo ""

# Check the backend logs
echo "📋 Backend Processing Log:"
tail -20 /tmp/backend.log | grep -E "Viber|track|TEST-TOKEN" | tail -10

echo ""
echo "========================================="
echo "✅ What should happen in production:"
echo ""
echo "1. Bot receives 'track_TEST-TOKEN-123' command"
echo "2. Bot queries database for shipment info"
echo "3. Bot generates static map URL with Mapbox"
echo "4. Bot sends Rich Media message with:"
echo "   - Static map showing delivery route"
echo "   - Current status and ETA"
echo "   - Button to open interactive map"
echo "   - Button to refresh status"
echo ""
echo "🗺️ Interactive map URL:"
echo "http://localhost:3001/en/track/TEST-TOKEN-123?viber=true&embedded=true"
echo ""
echo "📱 This URL opens in Viber's embedded browser:"
echo "   - Full-screen map without header/footer"
echo "   - Real-time courier location via WebSocket"
echo "   - Compact info panel at bottom"
echo "   - Optimized for mobile viewing"