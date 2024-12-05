#!/bin/bash

echo "Checking imports..."
echo "==================="

for file in $(find . -name "*.go"); do
    echo "File: $file"
    echo "Imports:"
    grep "^import" -A 10 "$file" | grep -v "^$"
    echo "-------------------"
done