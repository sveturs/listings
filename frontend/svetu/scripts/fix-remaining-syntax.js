const fs = require('fs');
const path = require('path');
const glob = require('glob');

function fixRemainingSyntaxErrors(filePath) {
  let content = fs.readFileSync(filePath, 'utf8');
  let modified = false;
  const changes = [];

  // Паттерны для исправления
  const patterns = [
    // Исправление tStorefronts.products.steps('basic') -> t('products.steps.basic')
    {
      regex: /tStorefronts\.products\.steps\('([^']+)'\)/g,
      replacement: "t('products.steps.$1')"
    },
    // Исправление tStorefronts.products.errors('key') -> t('products.errors.key')
    {
      regex: /tStorefronts\.products\.errors\('([^']+)'\)/g,
      replacement: "t('products.errors.$1')"
    },
    // Исправление tStorefronts.products('key') -> t('products.key')
    {
      regex: /tStorefronts\.products\('([^']+)'\)/g,
      replacement: "t('products.$1')"
    },
    // Исправление любых других tStorefronts.X.Y('key') -> t('X.Y.key')
    {
      regex: /tStorefronts\.([a-zA-Z_]+)\.([a-zA-Z_]+)\('([^']+)'\)/g,
      replacement: "t('$1.$2.$3')"
    },
    // Исправление tCreate_listing.X.Y('key') -> t('X.Y.key')
    {
      regex: /tCreate_listing\.([a-zA-Z_]+)\.([a-zA-Z_]+)\('([^']+)'\)/g,
      replacement: "t('$1.$2.$3')"
    },
    // Исправление tCreate_listing.X('key') -> t('X.key')
    {
      regex: /tCreate_listing\.([a-zA-Z_]+)\('([^']+)'\)/g,
      replacement: "t('$1.$2')"
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

  if (modified) {
    fs.writeFileSync(filePath, content);
  }

  return { modified, changes };
}

// Найти все файлы в products и create-listing
const patterns = [
  'src/components/products/**/*.{tsx,jsx}',
  'src/components/create-listing/**/*.{tsx,jsx}'
];

let files = [];
patterns.forEach(pattern => {
  files = files.concat(glob.sync(pattern, { cwd: process.cwd() }));
});

console.log(`Found ${files.length} files to check...`);

let totalFixed = 0;
const fixedFiles = [];

files.forEach(file => {
  const fullPath = path.join(process.cwd(), file);
  const { modified, changes } = fixRemainingSyntaxErrors(fullPath);
  
  if (modified) {
    totalFixed++;
    fixedFiles.push(file);
    console.log(`\n✓ Fixed ${file}:`);
    changes.forEach(change => console.log(`  - ${change}`));
  }
});

console.log(`\n✅ Fixed ${totalFixed} files:`);
fixedFiles.forEach(file => console.log(`  - ${file}`));