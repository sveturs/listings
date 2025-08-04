const fs = require('fs');
const path = require('path');
const glob = require('glob');

function fixAllSyntaxErrors(filePath) {
  let content = fs.readFileSync(filePath, 'utf8');
  let modified = false;
  const changes = [];

  // Паттерны для исправления
  const patterns = [
    // Исправление const tAdmin.balance.error = useTranslations -> const t = useTranslations
    {
      regex: /const tAdmin\.[a-zA-Z_\.]+ = useTranslations\('admin'\);/g,
      replacement: "const tAdmin = useTranslations('admin');"
    },
    // Исправление const tStorefronts.products = useTranslations -> const t = useTranslations
    {
      regex: /const tStorefronts\.[a-zA-Z_\.]+ = useTranslations\('storefronts'\);/g,
      replacement: "const tStorefronts = useTranslations('storefronts');"
    },
    // Исправление const tCreate_listing.* = useTranslations -> const t = useTranslations
    {
      regex: /const tCreate_listing\.[a-zA-Z_\.]+ = useTranslations\('create_listing'\);/g,
      replacement: "const tCreate_listing = useTranslations('create_listing');"
    },
    // Исправление const tCreate_storefront.* = useTranslations -> const t = useTranslations
    {
      regex: /const tCreate_storefront\.[a-zA-Z_\.]+ = useTranslations\('create_storefront'\);/g,
      replacement: "const tCreate_storefront = useTranslations('create_storefront');"
    },
    // Исправление tAdmin.balance.error('key') -> tAdmin('balance.error.key')
    {
      regex: /tAdmin\.balance\.error\('([^']+)'\)/g,
      replacement: "tAdmin('balance.error.$1')"
    },
    // Исправление tAdmin.common('key') -> tAdmin('common.key')
    {
      regex: /tAdmin\.common\('([^']+)'\)/g,
      replacement: "tAdmin('common.$1')"
    },
    // Исправление tStorefronts.products('key') -> tStorefronts('products.key')
    {
      regex: /tStorefronts\.products\('([^']+)'\)/g,
      replacement: "tStorefronts('products.$1')"
    },
    // Исправление tStorefronts.delivery('key') -> tStorefronts('delivery.key')
    {
      regex: /tStorefronts\.delivery\('([^']+)'\)/g,
      replacement: "tStorefronts('delivery.$1')"
    }
  ];

  patterns.forEach(pattern => {
    const matches = content.match(pattern.regex);
    if (matches) {
      content = content.replace(pattern.regex, pattern.replacement);
      modified = true;
      matches.forEach(match => {
        changes.push(`Fixed: ${match}`);
      });
    }
  });

  // Удаление дублирующихся объявлений переменных
  const lines = content.split('\n');
  const seenTranslations = new Map();
  const filteredLines = [];
  
  lines.forEach((line, index) => {
    const match = line.match(/^\s*const\s+(t|tCommon|tAdmin|tStorefronts|tCreate_listing|tCreate_storefront|tProfile|tRoles|tPermissions)\s*=\s*useTranslations\('([^']+)'\);?\s*$/);
    if (match) {
      const varName = match[1];
      const module = match[2];
      const key = `${varName}-${module}`;
      
      if (!seenTranslations.has(key)) {
        seenTranslations.set(key, index);
        filteredLines.push(line);
      } else {
        modified = true;
        changes.push(`Removed duplicate: ${line.trim()}`);
      }
    } else {
      filteredLines.push(line);
    }
  });

  if (modified) {
    content = filteredLines.join('\n');
    fs.writeFileSync(filePath, content);
  }

  return { modified, changes };
}

// Найти все файлы tsx/jsx
const files = glob.sync('src/**/*.{tsx,jsx}', { cwd: process.cwd() });

console.log(`Found ${files.length} files to check...`);

let totalFixed = 0;
const fixedFiles = [];

files.forEach(file => {
  const fullPath = path.join(process.cwd(), file);
  const { modified, changes } = fixAllSyntaxErrors(fullPath);
  
  if (modified) {
    totalFixed++;
    fixedFiles.push(file);
    console.log(`\n✓ Fixed ${file}:`);
    changes.forEach(change => console.log(`  - ${change}`));
  }
});

console.log(`\n✅ Fixed ${totalFixed} files:`);
fixedFiles.forEach(file => console.log(`  - ${file}`));