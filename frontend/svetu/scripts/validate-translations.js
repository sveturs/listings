const fs = require('fs');
const path = require('path');

const locales = ['en', 'ru', 'sr'];
const requiredKeys = {
  common: [
    'navigation.home',
    'navigation.search',
    'navigation.create',
    'navigation.chats',
    'navigation.profile',
    'navigation.close',
  ],
  marketplace: ['home', 'homeContent'],
  admin: ['title'],
};

function checkKeys(obj, keys) {
  const missing = [];
  for (const key of keys) {
    const parts = key.split('.');
    let current = obj;
    for (const part of parts) {
      if (!current || !current[part]) {
        missing.push(key);
        break;
      }
      current = current[part];
    }
  }
  return missing;
}

console.log('Validating translations...\n');

let hasErrors = false;

for (const locale of locales) {
  console.log(`Checking ${locale}:`);

  for (const [module, keys] of Object.entries(requiredKeys)) {
    const filePath = path.join(
      __dirname,
      `../src/messages/${locale}/${module}.json`
    );

    if (!fs.existsSync(filePath)) {
      console.log(`  ❌ Missing file: ${module}.json`);
      hasErrors = true;
      continue;
    }

    try {
      const content = JSON.parse(fs.readFileSync(filePath, 'utf8'));
      const missing = checkKeys(content, keys);

      if (missing.length > 0) {
        console.log(`  ❌ ${module}: Missing keys: ${missing.join(', ')}`);
        hasErrors = true;
      } else {
        console.log(`  ✅ ${module}: All keys present`);
      }
    } catch (e) {
      console.log(`  ❌ ${module}: Parse error - ${e.message}`);
      hasErrors = true;
    }
  }

  console.log();
}

if (hasErrors) {
  console.log('❌ Validation failed!');
  process.exit(1);
} else {
  console.log('✅ All translations are valid!');
}
