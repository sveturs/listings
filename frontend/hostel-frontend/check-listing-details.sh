#!/bin/bash

cd /data/hostel-booking-system/frontend/hostel-frontend/

# Используем опцию --isolatedModules для игнорирования проверки зависимостей
NODE_OPTIONS="--max-old-space-size=4096" npx tsc src/pages/marketplace/ListingDetailsPage.tsx \
  --isolatedModules \
  --skipLibCheck \
  --esModuleInterop \
  --target es5 \
  --jsx react \
  --moduleResolution node \
  --noEmit \
  --allowJs