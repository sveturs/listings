#!/bin/bash
# Script to run TypeScript compiler with skipLibCheck flag
# Usage: ./run-tsc.sh [files...]

# Get the directory of this script
DIR="$( cd "$( dirname "${BASH_SOURCE[0]}" )" && pwd )"
cd "$DIR"

# Check if specific files are provided
if [ $# -gt 0 ]; then
  # Run TypeScript compiler on specific files
  echo "Checking TypeScript files: $@"
  npx tsc --noEmit --skipLibCheck --jsx react --jsxFactory React.createElement --moduleResolution node --esModuleInterop --allowSyntheticDefaultImports --ignoreDeprecations 5.0 --skipDefaultLibCheck $@
else
  # Run TypeScript compiler on all files
  echo "Checking all TypeScript files"
  npx tsc --noEmit --skipLibCheck --jsx react --jsxFactory React.createElement --moduleResolution node --esModuleInterop --allowSyntheticDefaultImports --ignoreDeprecations 5.0 --skipDefaultLibCheck
fi