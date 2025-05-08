#!/bin/bash
export NODE_OPTIONS=--openssl-legacy-provider
npx tsc src/components/icons/SveTuLogo.tsx --noEmit --jsx react --esModuleInterop --skipLibCheck