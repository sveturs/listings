#!/bin/bash

# Setup variables
COMPONENT_NAME=$1
COMPONENT_PATH=$2

# Check if component name and path were provided
if [ -z "$COMPONENT_NAME" ] || [ -z "$COMPONENT_PATH" ]; then
  echo "Error: Please provide both component name and path"
  echo "Usage: ./create-ts-component.sh ComponentName src/components/path"
  exit 1
fi

# Create full path
FULL_PATH="$COMPONENT_PATH/$COMPONENT_NAME.tsx"

# Check if the directory exists
if [ ! -d "$COMPONENT_PATH" ]; then
  echo "Creating directory: $COMPONENT_PATH"
  mkdir -p "$COMPONENT_PATH"
fi

# Check if the file already exists
if [ -f "$FULL_PATH" ]; then
  echo "Error: Component $FULL_PATH already exists"
  exit 1
fi

# Create the component file
cat > "$FULL_PATH" << EOF
import React from 'react';

interface ${COMPONENT_NAME}Props {
  // Define your props here
}

const $COMPONENT_NAME: React.FC<${COMPONENT_NAME}Props> = (props) => {
  return (
    <div>
      {/* Component content */}
    </div>
  );
};

export default $COMPONENT_NAME;
EOF

echo "✅ TypeScript component created at $FULL_PATH"