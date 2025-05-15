#!/bin/bash
export NODE_OPTIONS="--openssl-legacy-provider"

# List of files to check
FILES=(
  "src/components/maps/MiniMap.tsx"
  "src/pages/gis/GISMapPage.tsx"
  "src/pages/marketplace/ChatPage.tsx" 
  "src/pages/marketplace/CreateListingPage.tsx"
  "src/pages/marketplace/EditListingPage.tsx"
  "src/pages/marketplace/FavoriteListingsPage.tsx"
  "src/pages/marketplace/ListingDetailsPage.tsx"
)

# Create a temporary TypeScript configuration file that ignores i18next
cat > tsconfig.temp.json << EOF
{
  "compilerOptions": {
    "target": "es5",
    "lib": ["dom", "dom.iterable", "esnext"],
    "allowJs": true,
    "skipLibCheck": true,
    "esModuleInterop": true,
    "allowSyntheticDefaultImports": true,
    "strict": false,
    "forceConsistentCasingInFileNames": true,
    "module": "esnext",
    "moduleResolution": "node",
    "resolveJsonModule": true,
    "isolatedModules": true,
    "noEmit": true,
    "jsx": "react",
    "typeRoots": ["./src/types", "./node_modules/@types"]
  },
  "include": ["src/types/global.d.ts", "src/types/i18next.d.ts", "src/types/i18next-ignore.d.ts"],
  "files": []
}
EOF

# Check each file individually
for file in "${FILES[@]}"; do
  echo "Checking $file..."
  
  # Create a temporary version of tsconfig.temp.json with just this file
  cp tsconfig.temp.json tsconfig.check.json
  sed -i "s/\"files\": \[\]/\"files\": \[\"$file\"\]/" tsconfig.check.json
  
  # Run TypeScript check with our temporary config
  npx tsc --project tsconfig.check.json
  
  if [ $? -eq 0 ]; then
    echo "✅ $file passed"
  else
    echo "❌ $file failed"
  fi
  
  echo ""
done

# Clean up temporary files
rm tsconfig.temp.json tsconfig.check.json