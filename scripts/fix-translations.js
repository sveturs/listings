#!/usr/bin/env node

const fs = require('fs');
const path = require('path');
const glob = require('glob');

// Функция для исправления вложенных namespace
function fixNestedNamespaces(content) {
  // Исправляем useTranslations с вложенными namespace
  content = content.replace(/useTranslations\(['"]([^'"]+)\.([^'"]+)['"]\)/g, (match, module) => {
    const baseModule = module.split('.')[0];
    console.log(`  Fixing: ${match} -> useTranslations('${baseModule}')`);
    return `useTranslations('${baseModule}')`;
  });
  
  return content;
}

// Функция для исправления пустых useTranslations()
function fixEmptyTranslations(content, filePath) {
  if (content.includes('useTranslations()')) {
    // Определяем модуль на основе пути файла
    let module = 'common';
    
    if (filePath.includes('/marketplace/')) module = 'marketplace';
    else if (filePath.includes('/auth/')) module = 'auth';
    else if (filePath.includes('/profile/')) module = 'profile';
    else if (filePath.includes('/admin/')) module = 'admin';
    else if (filePath.includes('/orders/')) module = 'orders';
    else if (filePath.includes('/cars/')) module = 'cars';
    else if (filePath.includes('/storefront/')) module = 'storefront';
    else if (filePath.includes('/map/')) module = 'map';
    else if (filePath.includes('/examples/')) module = 'examples';
    
    content = content.replace(/useTranslations\(\)/g, `useTranslations('${module}')`);
    console.log(`  Fixed empty useTranslations() with module: ${module}`);
  }
  
  return content;
}

// Главная функция
async function main() {
  const srcPath = '/data/hostel-booking-system/frontend/svetu/src';
  const patterns = [
    '**/*.tsx',
    '**/*.ts'
  ];
  
  console.log('🔍 Searching for files with translation issues...\n');
  
  let totalFixed = 0;
  
  for (const pattern of patterns) {
    const files = glob.sync(path.join(srcPath, pattern), {
      ignore: ['**/node_modules/**', '**/dist/**', '**/build/**']
    });
    
    for (const file of files) {
      const content = fs.readFileSync(file, 'utf8');
      let modified = content;
      
      // Проверяем, нужно ли исправление
      const hasNestedNamespace = /useTranslations\(['"][^'"]+\.[^'"]+['"]\)/.test(content);
      const hasEmptyTranslations = /useTranslations\(\)/.test(content);
      
      if (hasNestedNamespace || hasEmptyTranslations) {
        console.log(`📝 Processing: ${file.replace(srcPath, '.')}`);
        
        // Применяем исправления
        modified = fixNestedNamespaces(modified);
        modified = fixEmptyTranslations(modified, file);
        
        // Сохраняем файл
        fs.writeFileSync(file, modified, 'utf8');
        totalFixed++;
        console.log(`  ✅ Fixed!\n`);
      }
    }
  }
  
  console.log(`\n✨ Done! Fixed ${totalFixed} files.`);
  
  // Дополнительная проверка
  console.log('\n🔍 Checking for remaining issues...\n');
  
  const remainingIssues = [];
  
  for (const pattern of patterns) {
    const files = glob.sync(path.join(srcPath, pattern), {
      ignore: ['**/node_modules/**', '**/dist/**', '**/build/**']
    });
    
    for (const file of files) {
      const content = fs.readFileSync(file, 'utf8');
      
      if (/useTranslations\(['"][^'"]+\.[^'"]+['"]\)/.test(content)) {
        remainingIssues.push({ file: file.replace(srcPath, '.'), issue: 'nested namespace' });
      }
      
      if (/useTranslations\(\)/.test(content)) {
        remainingIssues.push({ file: file.replace(srcPath, '.'), issue: 'empty call' });
      }
    }
  }
  
  if (remainingIssues.length > 0) {
    console.log('⚠️ Remaining issues found:');
    remainingIssues.forEach(({ file, issue }) => {
      console.log(`  - ${file}: ${issue}`);
    });
  } else {
    console.log('✅ No remaining issues found!');
  }
}

// Запуск
main().catch(console.error);