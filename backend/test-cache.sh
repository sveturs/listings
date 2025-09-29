#!/bin/bash

echo "=== Redis Cache Testing Script ==="
echo "Testing marketplace cache implementation..."

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

# Check Redis is running
echo -e "\n${YELLOW}1. Checking Redis connection...${NC}"
if redis-cli ping > /dev/null 2>&1; then
    echo -e "${GREEN}✓ Redis is running${NC}"
else
    echo -e "${RED}✗ Redis is not running${NC}"
    exit 1
fi

# Clear cache for clean test
echo -e "\n${YELLOW}2. Clearing Redis cache (DB 1)...${NC}"
redis-cli -n 1 FLUSHDB > /dev/null 2>&1
echo -e "${GREEN}✓ Cache cleared${NC}"

# Test search API caching
echo -e "\n${YELLOW}3. Testing Search API caching...${NC}"

echo "   First request (no cache):"
START=$(date +%s%N)
RESULT1=$(curl -s "http://localhost:3000/api/v1/marketplace/search?category_id=1301&limit=5" 2>/dev/null)
END=$(date +%s%N)
TIME1=$((($END - $START)/1000000))
echo -e "   Time: ${TIME1}ms"

echo "   Second request (should be cached):"
START=$(date +%s%N)
RESULT2=$(curl -s "http://localhost:3000/api/v1/marketplace/search?category_id=1301&limit=5" 2>/dev/null)
END=$(date +%s%N)
TIME2=$((($END - $START)/1000000))
echo -e "   Time: ${TIME2}ms"

if [ "$TIME2" -lt "$TIME1" ]; then
    echo -e "   ${GREEN}✓ Cache is working! Second request was faster (${TIME2}ms < ${TIME1}ms)${NC}"
else
    echo -e "   ${YELLOW}⚠ Cache might not be working properly${NC}"
fi

# Test listing details caching
echo -e "\n${YELLOW}4. Testing Listing Details API caching...${NC}"

echo "   First request (no cache):"
START=$(date +%s%N)
RESULT3=$(curl -s "http://localhost:3000/api/v1/marketplace/listings/329" 2>/dev/null)
END=$(date +%s%N)
TIME3=$((($END - $START)/1000000))
echo -e "   Time: ${TIME3}ms"

echo "   Second request (should be cached):"
START=$(date +%s%N)
RESULT4=$(curl -s "http://localhost:3000/api/v1/marketplace/listings/329" 2>/dev/null)
END=$(date +%s%N)
TIME4=$((($END - $START)/1000000))
echo -e "   Time: ${TIME4}ms"

if [ "$TIME4" -lt "$TIME3" ]; then
    echo -e "   ${GREEN}✓ Cache is working! Second request was faster (${TIME4}ms < ${TIME3}ms)${NC}"
else
    echo -e "   ${YELLOW}⚠ Cache might not be working properly${NC}"
fi

# Check Redis keys
echo -e "\n${YELLOW}5. Checking Redis keys created...${NC}"
SEARCH_KEYS=$(redis-cli -n 1 keys "search:*" | wc -l)
LISTING_KEYS=$(redis-cli -n 1 keys "listing:*" | wc -l)
STATS_KEYS=$(redis-cli -n 1 keys "cache:stats:*" | wc -l)

echo "   Search keys: $SEARCH_KEYS"
echo "   Listing keys: $LISTING_KEYS"
echo "   Stats keys: $STATS_KEYS"

if [ "$SEARCH_KEYS" -gt 0 ] || [ "$LISTING_KEYS" -gt 0 ]; then
    echo -e "   ${GREEN}✓ Cache keys are being created${NC}"
else
    echo -e "   ${RED}✗ No cache keys found - cache might not be working${NC}"
fi

# Show cache statistics
echo -e "\n${YELLOW}6. Cache Statistics:${NC}"
HITS=$(redis-cli -n 1 get "cache:stats:hits" 2>/dev/null || echo "0")
MISSES=$(redis-cli -n 1 get "cache:stats:misses" 2>/dev/null || echo "0")
echo "   Cache hits: $HITS"
echo "   Cache misses: $MISSES"

if [ "$HITS" != "0" ] || [ "$MISSES" != "0" ]; then
    HIT_RATE=$(echo "scale=2; $HITS * 100 / ($HITS + $MISSES)" | bc 2>/dev/null || echo "0")
    echo -e "   Hit rate: ${HIT_RATE}%"
fi

echo -e "\n${GREEN}=== Cache testing completed ===${NC}"