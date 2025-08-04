const fs = require('fs');

const filePath = './src/components/Storefront/ProductVariants/VariantStockTable.tsx';
let content = fs.readFileSync(filePath, 'utf8');

// Fix all sort field names
content = content.replace(/handleSort\('variants\.name'\)/g, "handleSort('name')");
content = content.replace(/handleSort\('variants\.sku'\)/g, "handleSort('sku')");
content = content.replace(/handleSort\('variants\.price'\)/g, "handleSort('price')");
content = content.replace(/handleSort\('variants\.stock'\)/g, "handleSort('stock')");
content = content.replace(/handleSort\('variants\.updated_at'\)/g, "handleSort('updated_at')");

// Fix sort.field comparisons if any
content = content.replace(/sort\.field === 'variants\.name'/g, "sort.field === 'name'");
content = content.replace(/sort\.field === 'variants\.sku'/g, "sort.field === 'sku'");
content = content.replace(/sort\.field === 'variants\.price'/g, "sort.field === 'price'");
content = content.replace(/sort\.field === 'variants\.stock'/g, "sort.field === 'stock'");
content = content.replace(/sort\.field === 'variants\.updated_at'/g, "sort.field === 'updated_at'");

fs.writeFileSync(filePath, content);
console.log('Fixed variant table sort fields successfully!');