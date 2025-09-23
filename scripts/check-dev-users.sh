#!/bin/bash

# Script to monitor users on dev.svetu.rs
# Shows user activity, listings, and other useful information

echo "=================================================="
echo "     DEV.SVETU.RS USER MONITORING DASHBOARD      "
echo "=================================================="
echo ""

# Get current date/time on server
echo "📅 Report generated: $(date '+%Y-%m-%d %H:%M:%S')"
echo ""

# Total users count
echo "📊 USER STATISTICS"
echo "──────────────────"
ssh svetu@svetu.rs "docker exec svetu-dev_db_1 psql -U svetu_dev_user -d svetu_dev_db -t -c \"
SELECT
    'Total users: ' || COUNT(*)
FROM users;
\"" | sed 's/^[ \t]*//'

ssh svetu@svetu.rs "docker exec svetu-dev_db_1 psql -U svetu_dev_user -d svetu_dev_db -t -c \"
SELECT
    'Active in last 7 days: ' || COUNT(*)
FROM users
WHERE updated_at > NOW() - INTERVAL '7 days';
\"" | sed 's/^[ \t]*//'

ssh svetu@svetu.rs "docker exec svetu-dev_db_1 psql -U svetu_dev_user -d svetu_dev_db -t -c \"
SELECT
    'New users (last 30 days): ' || COUNT(*)
FROM users
WHERE created_at > NOW() - INTERVAL '30 days';
\"" | sed 's/^[ \t]*//'

echo ""
echo "👥 ALL USERS (sorted by last activity)"
echo "──────────────────────────────────────"
echo ""

ssh svetu@svetu.rs "docker exec svetu-dev_db_1 psql -U svetu_dev_user -d svetu_dev_db -c \"
SELECT
    u.id AS \\\"ID\\\",
    CASE
        WHEN LENGTH(u.email) > 30 THEN SUBSTRING(u.email, 1, 27) || '...'
        ELSE RPAD(u.email, 30, ' ')
    END AS \\\"Email\\\",
    CASE
        WHEN LENGTH(u.name) > 20 THEN SUBSTRING(u.name, 1, 17) || '...'
        ELSE RPAD(COALESCE(u.name, '-'), 20, ' ')
    END AS \\\"Name\\\",
    TO_CHAR(u.created_at, 'DD.MM.YY HH24:MI') AS \\\"Registered\\\",
    TO_CHAR(u.updated_at, 'DD.MM.YY HH24:MI') AS \\\"Last Activity\\\",
    CASE
        WHEN u.updated_at > NOW() - INTERVAL '1 day' THEN '🟢'
        WHEN u.updated_at > NOW() - INTERVAL '7 days' THEN '🟡'
        WHEN u.updated_at > NOW() - INTERVAL '30 days' THEN '🟠'
        ELSE '⚫'
    END AS \\\"St\\\",
    COALESCE((SELECT COUNT(*) FROM marketplace_listings WHERE user_id = u.id)::TEXT, '0') AS \\\"Ads\\\",
    COALESCE((SELECT COUNT(*) FROM storefronts WHERE user_id = u.id)::TEXT, '0') AS \\\"Stores\\\",
    CASE
        WHEN u.account_status = 'active' THEN '✅'
        WHEN u.account_status = 'suspended' THEN '🚫'
        WHEN u.account_status = 'deleted' THEN '❌'
        ELSE ''
    END AS \\\"Status\\\"
FROM users u
ORDER BY u.updated_at DESC;
\""

echo ""
echo "📈 USER ACTIVITY (last 24 hours)"
echo "─────────────────────────────────"
ssh svetu@svetu.rs "docker exec svetu-dev_db_1 psql -U svetu_dev_user -d svetu_dev_db -t -c \"
SELECT
    CASE
        WHEN COUNT(*) = 0 THEN '  No activity in the last 24 hours'
        ELSE ''
    END
FROM users
WHERE updated_at > NOW() - INTERVAL '24 hours';
\"" | sed 's/^[ \t]*//'

ssh svetu@svetu.rs "docker exec svetu-dev_db_1 psql -U svetu_dev_user -d svetu_dev_db -c \"
SELECT
    u.email AS \\\"User\\\",
    TO_CHAR(u.updated_at, 'HH24:MI:SS') AS \\\"Time\\\",
    CASE
        WHEN ml.created_at IS NOT NULL THEN 'Created listing: ' || ml.title
        WHEN sf.created_at IS NOT NULL THEN 'Created storefront: ' || sf.name
        ELSE 'Profile update'
    END AS \\\"Action\\\"
FROM users u
LEFT JOIN LATERAL (
    SELECT title, created_at
    FROM marketplace_listings
    WHERE user_id = u.id
    AND created_at > NOW() - INTERVAL '24 hours'
    ORDER BY created_at DESC
    LIMIT 1
) ml ON true
LEFT JOIN LATERAL (
    SELECT name, created_at
    FROM storefronts
    WHERE user_id = u.id
    AND created_at > NOW() - INTERVAL '24 hours'
    ORDER BY created_at DESC
    LIMIT 1
) sf ON true
WHERE u.updated_at > NOW() - INTERVAL '24 hours'
ORDER BY u.updated_at DESC
LIMIT 10;
\"" 2>/dev/null || echo "  No recent activity"

echo ""
echo "🆕 NEWEST USERS (last 5 registered)"
echo "────────────────────────────────────"
ssh svetu@svetu.rs "docker exec svetu-dev_db_1 psql -U svetu_dev_user -d svetu_dev_db -c \"
SELECT
    RPAD(u.email, 35, ' ') AS \\\"Email\\\",
    RPAD(COALESCE(u.name, '-'), 25, ' ') AS \\\"Name\\\",
    TO_CHAR(u.created_at, 'DD.MM.YYYY HH24:MI') AS \\\"Registered\\\"
FROM users u
ORDER BY u.created_at DESC
LIMIT 5;
\""

echo ""
echo "📦 TOP CONTENT CREATORS"
echo "────────────────────────"
ssh svetu@svetu.rs "docker exec svetu-dev_db_1 psql -U svetu_dev_user -d svetu_dev_db -c \"
SELECT
    RPAD(u.email, 30, ' ') AS \\\"User\\\",
    COUNT(ml.id) AS \\\"Listings\\\",
    COUNT(DISTINCT sf.id) AS \\\"Stores\\\",
    TO_CHAR(MAX(ml.created_at), 'DD.MM.YY') AS \\\"Last Listing\\\"
FROM users u
LEFT JOIN marketplace_listings ml ON u.id = ml.user_id
LEFT JOIN storefronts sf ON u.id = sf.user_id
GROUP BY u.id, u.email
HAVING COUNT(ml.id) > 0 OR COUNT(DISTINCT sf.id) > 0
ORDER BY COUNT(ml.id) DESC, COUNT(DISTINCT sf.id) DESC
LIMIT 5;
\""

echo ""
echo "📋 LEGEND"
echo "─────────"
echo "Status indicators:"
echo "  🟢 Active today"
echo "  🟡 Active this week"
echo "  🟠 Active this month"
echo "  ⚫ Inactive > 30 days"
echo "  ✅ Active account"
echo "  🚫 Suspended account"
echo "  ❌ Deleted account"
echo ""
echo "=================================================="