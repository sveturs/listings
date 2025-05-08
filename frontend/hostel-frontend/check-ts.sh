#!/bin/bash

# Setup variables
COMPONENT_PATH=$1
NODE_OPTIONS="--openssl-legacy-provider"

# Check if a component path was provided
if [ -z "$COMPONENT_PATH" ]; then
  echo "Error: Please provide a component path to check"
  echo "Usage: ./check-ts.sh src/components/path/to/Component.tsx"
  exit 1
fi

# Check if the file exists
if [ ! -f "$COMPONENT_PATH" ]; then
  echo "Error: File $COMPONENT_PATH does not exist"
  exit 1
fi

# Run TypeScript compiler to check types
echo "Checking TypeScript types for $COMPONENT_PATH..."
export NODE_OPTIONS
npx tsc "$COMPONENT_PATH" --noEmit --jsx react --esModuleInterop --skipLibCheck

# Check if tsc succeeded
if [ $? -eq 0 ]; then
  echo "✅ TypeScript check passed for $COMPONENT_PATH"
else
  echo "❌ TypeScript check failed for $COMPONENT_PATH"
  exit 1
fi