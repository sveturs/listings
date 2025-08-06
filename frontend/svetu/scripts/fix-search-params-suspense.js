#!/usr/bin/env node

const fs = require('fs');
const path = require('path');
const glob = require('glob');

// Функция для чтения файла
function readFile(filePath) {
  return fs.readFileSync(filePath, 'utf8');
}

// Функция для записи файла
function writeFile(filePath, content) {
  fs.writeFileSync(filePath, content);
}

// Проверка, содержит ли файл useSearchParams
function hasUseSearchParams(content) {
  return content.includes('useSearchParams');
}

// Проверка, есть ли уже Suspense
function hasSuspense(content) {
  return content.includes('Suspense') && content.includes('<Suspense');
}

// Создание клиентского компонента
function createClientComponent(pageFile, componentName) {
  const dir = path.dirname(pageFile);
  const clientFile = path.join(dir, `${componentName}Client.tsx`);

  // Читаем оригинальный файл
  let content = readFile(pageFile);

  // Меняем имя компонента
  content = content.replace(
    /export default function \w+Page/,
    `export default function ${componentName}Client`
  );

  // Сохраняем клиентский компонент
  writeFile(clientFile, content);

  return clientFile;
}

// Обновление page.tsx для использования Suspense
function updatePageWithSuspense(pageFile, componentName) {
  const content = `import { Suspense } from 'react';
import ${componentName}Client from './${componentName}Client';

export default function ${componentName}Page() {
  return (
    <Suspense fallback={<div className="min-h-screen flex items-center justify-center">Loading...</div>}>
      <${componentName}Client />
    </Suspense>
  );
}`;

  writeFile(pageFile, content);
}

// Получаем имя компонента из пути
function getComponentName(filePath) {
  const parts = filePath.split('/');
  const parentDir = parts[parts.length - 2];

  // Конвертируем kebab-case в PascalCase
  return parentDir
    .split('-')
    .map((word) => word.charAt(0).toUpperCase() + word.slice(1))
    .join('');
}

// Находим все файлы page.tsx
const files = glob.sync('src/app/**/page.tsx', {
  ignore: ['**/node_modules/**'],
});

let fixedCount = 0;
const problematicFiles = [];

files.forEach((file) => {
  const content = readFile(file);

  // Проверяем, использует ли файл useSearchParams
  if (hasUseSearchParams(content) && !hasSuspense(content)) {
    console.log(`Found useSearchParams without Suspense in: ${file}`);

    // Проверяем, является ли это клиентским компонентом
    if (content.includes("'use client'")) {
      const componentName = getComponentName(file);

      try {
        // Создаем клиентский компонент
        const clientFile = createClientComponent(file, componentName);
        console.log(`  Created client component: ${clientFile}`);

        // Обновляем page.tsx
        updatePageWithSuspense(file, componentName);
        console.log(`  Updated page with Suspense: ${file}`);

        fixedCount++;
      } catch (error) {
        console.error(`  Error processing ${file}: ${error.message}`);
        problematicFiles.push(file);
      }
    } else {
      console.log(`  Skipping server component: ${file}`);
    }
  }
});

console.log(`\n✅ Fixed ${fixedCount} files`);

if (problematicFiles.length > 0) {
  console.log(`\n⚠️  Problematic files that need manual review:`);
  problematicFiles.forEach((file) => console.log(`  - ${file}`));
}
