#!/usr/bin/env node

const fs = require('fs');
const path = require('path');

// Required environment variables
const requiredVars = ['NEXT_PUBLIC_API_URL', 'NEXT_PUBLIC_MINIO_URL'];

// Optional but recommended
const recommendedVars = [
  'NEXT_PUBLIC_WEBSOCKET_URL',
  'NEXT_PUBLIC_GOOGLE_CLIENT_ID',
];

console.log('🔍 Checking environment variables...\n');

let hasErrors = false;
let hasWarnings = false;

// Check required variables
console.log('Required variables:');
requiredVars.forEach((varName) => {
  if (process.env[varName]) {
    console.log(`✅ ${varName}: ${process.env[varName]}`);
  } else {
    console.log(`❌ ${varName}: NOT SET`);
    hasErrors = true;
  }
});

console.log('\nRecommended variables:');
recommendedVars.forEach((varName) => {
  if (process.env[varName]) {
    console.log(`✅ ${varName}: ${process.env[varName]}`);
  } else {
    console.log(`⚠️  ${varName}: not set (optional)`);
    hasWarnings = true;
  }
});

// Check for .env.local file
const envLocalPath = path.join(process.cwd(), '.env.local');
if (!fs.existsSync(envLocalPath)) {
  console.log(
    '\n⚠️  No .env.local file found. Using defaults or system environment.'
  );
  hasWarnings = true;
}

// Summary
console.log('\n' + '='.repeat(50));
if (hasErrors) {
  console.log('❌ Environment check failed! Missing required variables.');
  process.exit(1);
} else if (hasWarnings) {
  console.log('⚠️  Environment check passed with warnings.');
} else {
  console.log('✅ Environment check passed!');
}
