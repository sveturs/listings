const fs = require('fs');
const path = require('path');
const glob = require('glob');

function fixCreateListingSyntaxErrors(filePath) {
  let content = fs.readFileSync(filePath, 'utf8');
  let modified = false;
  const changes = [];

  // Паттерны для исправления
  const patterns = [
    // Исправление const tCreate_listing.attributes = useTranslations -> const t = useTranslations
    {
      regex:
        /const tCreate_listing\.[a-zA-Z_\.]+ = useTranslations\('create_listing'\);/g,
      replacement: "const t = useTranslations('create_listing');",
    },
    // Исправление const tStorefronts.products.attributeGroups = useTranslations
    {
      regex:
        /const tStorefronts\.[a-zA-Z_\.]+ = useTranslations\('storefronts'\);/g,
      replacement: "const t = useTranslations('storefronts');",
    },
    // Исправление tCreate_listing.basic_info('title') -> t('basic_info.title')
    {
      regex: /tCreate_listing\.basic_info\('([^']+)'\)/g,
      replacement: "t('basic_info.$1')",
    },
    // Исправление tCreate_listing.regional_tips('serbia') -> t('regional_tips.serbia')
    {
      regex: /tCreate_listing\.regional_tips\('([^']+)'\)/g,
      replacement: "t('regional_tips.$1')",
    },
    // Исправление tCreate_listing.errors('title_required') -> t('errors.title_required')
    {
      regex: /tCreate_listing\.errors\('([^']+)'\)/g,
      replacement: "t('errors.$1')",
    },
    // Исправление tCreate_listing.preview.rules('rule1') -> t('preview.rules.rule1')
    {
      regex: /tCreate_listing\.preview\.rules\('([^']+)'\)/g,
      replacement: "t('preview.rules.$1')",
    },
    // Исправление tCreate_listing.photos('add_photos') -> t('photos.add_photos')
    {
      regex: /tCreate_listing\.photos\('([^']+)'\)/g,
      replacement: "t('photos.$1')",
    },
    // Исправление tCreate_listing.attributes('title') -> t('attributes.title')
    {
      regex: /tCreate_listing\.attributes\('([^']+)'\)/g,
      replacement: "t('attributes.$1')",
    },
    // Исправление tStorefronts.products.attributeGroups('key') -> t('products.attributeGroups.key')
    {
      regex: /tStorefronts\.products\.attributeGroups\('([^']+)'\)/g,
      replacement: "t('products.attributeGroups.$1')",
    },
    // Исправление tStorefronts.products('title') -> t('products.title')
    {
      regex: /tStorefronts\.products\('([^']+)'\)/g,
      replacement: "t('products.$1')",
    },
  ];

  patterns.forEach((pattern) => {
    const matches = content.match(pattern.regex);
    if (matches) {
      content = content.replace(pattern.regex, pattern.replacement);
      modified = true;
      matches.forEach((match) => {
        changes.push(`Fixed: ${match}`);
      });
    }
  });

  // Удаление дублирующихся const t = useTranslations
  const lines = content.split('\n');
  const seenTranslations = new Set();
  const filteredLines = [];

  lines.forEach((line) => {
    if (
      line.includes('const t = useTranslations(') ||
      line.includes('const tCommon = useTranslations(') ||
      line.includes('const tCreate_listing = useTranslations(') ||
      line.includes('const tStorefronts = useTranslations(')
    ) {
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

// Найти все файлы create-listing и products
const patterns = [
  'src/components/create-listing/**/*.{tsx,jsx}',
  'src/components/products/**/*.{tsx,jsx}',
];

let files = [];
patterns.forEach((pattern) => {
  files = files.concat(glob.sync(pattern, { cwd: process.cwd() }));
});

console.log(`Found ${files.length} files to check...`);

let totalFixed = 0;
const fixedFiles = [];

files.forEach((file) => {
  const fullPath = path.join(process.cwd(), file);
  const { modified, changes } = fixCreateListingSyntaxErrors(fullPath);

  if (modified) {
    totalFixed++;
    fixedFiles.push(file);
    console.log(`\n✓ Fixed ${file}:`);
    changes.forEach((change) => console.log(`  - ${change}`));
  }
});

console.log(`\n✅ Fixed ${totalFixed} files:`);
fixedFiles.forEach((file) => console.log(`  - ${file}`));
