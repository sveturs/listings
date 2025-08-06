const fs = require('fs');

const filePath =
  './src/app/[locale]/storefronts/[slug]/products/[id]/edit/page.tsx';
let content = fs.readFileSync(filePath, 'utf8');

// Replace all occurrences of t( with tStorefronts(
content = content.replace(/\bt\(/g, 'tStorefronts(');

// Update the translation keys
content = content.replace(
  /tStorefronts\('backToProducts'\)/g,
  "tStorefronts('products.backToProducts')"
);
content = content.replace(
  /tStorefronts\('editProduct'\)/g,
  "tStorefronts('products.editProduct')"
);
content = content.replace(
  /tStorefronts\('editingProduct'\)/g,
  "tStorefronts('products.editingProduct')"
);

// Remove duplicate const declaration
content = content.replace(
  "const t = useTranslations('storefronts');\n  const tStorefronts = useTranslations('storefronts');",
  "const tStorefronts = useTranslations('storefronts');"
);

fs.writeFileSync(filePath, content);
console.log('Fixed edit page successfully!');
