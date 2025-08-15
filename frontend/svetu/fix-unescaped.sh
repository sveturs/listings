#!/bin/bash

# Fix unescaped quotes in React components
find src -name "*.tsx" -type f -exec sed -i \
  -e 's/"\([^"]*\)"/\&quot;\1\&quot;/g' \
  -e "s/'/\&apos;/g" \
  {} \;

echo "Fixed unescaped entities"