const fs = require('fs');
const path = require('path');
const glob = require('glob');

function fixStorefrontSyntaxErrors(filePath) {
  let content = fs.readFileSync(filePath, 'utf8');
  let modified = false;
  const changes = [];

  // Паттерны для исправления
  const patterns = [
    // Исправление tStorefronts.products('title') -> t('storefronts.products.title')
    {
      regex: /tStorefronts\.products\('([^']+)'\)/g,
      replacement: "t('storefronts.products.$1')"
    },
    // Исправление tStorefronts.payment_methods('cash') -> t('storefronts.payment_methods.cash')
    {
      regex: /tStorefronts\.payment_methods\('([^']+)'\)/g,
      replacement: "t('storefronts.payment_methods.$1')"
    },
    // Исправление tCreate_storefront.business_hours = useTranslations -> const t = useTranslations
    {
      regex: /const tCreate_storefront\.[a-zA-Z_]+ = useTranslations\('create_storefront'\);/g,
      replacement: "const t = useTranslations('create_storefront');"
    },
    // Исправление tCreate_storefront.business_hours('title') -> t('business_hours.title')
    {
      regex: /tCreate_storefront\.business_hours\('([^']+)'\)/g,
      replacement: "t('business_hours.$1')"
    },
    // Исправление tCreate_storefront.errors('name_required') -> t('errors.name_required')
    {
      regex: /tCreate_storefront\.errors\('([^']+)'\)/g,
      replacement: "t('errors.$1')"
    },
    // Исправление tCreate_storefront.business_types('retail') -> t('business_types.retail')
    {
      regex: /tCreate_storefront\.business_types\('([^']+)'\)/g,
      replacement: "t('business_types.$1')"
    },
    // Исправление tCreate_storefront.location('title') -> t('location.title')
    {
      regex: /tCreate_storefront\.location\('([^']+)'\)/g,
      replacement: "t('location.$1')"
    },
    // Исправление tCreate_storefront.preview('title') -> t('preview.title')
    {
      regex: /tCreate_storefront\.preview\('([^']+)'\)/g,
      replacement: "t('preview.$1')"
    },
    // Исправление tCreate_storefront.basic_info('name') -> t('basic_info.name')
    {
      regex: /tCreate_storefront\.basic_info\('([^']+)'\)/g,
      replacement: "t('basic_info.$1')"
    },
    // Исправление tCreate_storefront.business_details('tax_number') -> t('business_details.tax_number')
    {
      regex: /tCreate_storefront\.business_details\('([^']+)'\)/g,
      replacement: "t('business_details.$1')"
    },
    // Исправление tStorefronts('verified') -> t('storefronts.verified')
    {
      regex: /tStorefronts\('([^']+)'\)/g,
      replacement: "t('storefronts.$1')"
    },
    // Исправление tRoles('staff') -> t('roles.staff')
    {
      regex: /tRoles\('([^']+)'\)/g,
      replacement: "t('roles.$1')"
    },
    // Исправление tPermissions('products') -> t('permissions.products')
    {
      regex: /tPermissions\('([^']+)'\)/g,
      replacement: "t('permissions.$1')"
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

  // Удаление дублирующихся const t = useTranslations
  const lines = content.split('\n');
  const seenTranslations = new Set();
  const filteredLines = [];
  
  lines.forEach(line => {
    if (line.includes('const t = useTranslations(') || 
        line.includes('const tCommon = useTranslations(') ||
        line.includes('const tStorefronts = useTranslations(') ||
        line.includes('const tRoles = useTranslations(') ||
        line.includes('const tPermissions = useTranslations(')) {
      const trimmedLine = line.replace(/^\s+/, '').replace(/\s+$/, '');
      if (!seenTranslations.has(trimmedLine)) {
        seenTranslations.add(trimmedLine);
        filteredLines.push(line);
      } else {
        modified = true;
        changes.push(`Removed duplicate: ${trimmedLine}`);
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

// Найти все файлы storefronts
const files = glob.sync('src/components/storefronts/**/*.{tsx,jsx}', { cwd: process.cwd() });

console.log(`Found ${files.length} storefront component files to check...`);

let totalFixed = 0;
const fixedFiles = [];

files.forEach(file => {
  const fullPath = path.join(process.cwd(), file);
  const { modified, changes } = fixStorefrontSyntaxErrors(fullPath);
  
  if (modified) {
    totalFixed++;
    fixedFiles.push(file);
    console.log(`\n✓ Fixed ${file}:`);
    changes.forEach(change => console.log(`  - ${change}`));
  }
});

console.log(`\n✅ Fixed ${totalFixed} files:`);
fixedFiles.forEach(file => console.log(`  - ${file}`));