#!/usr/bin/env node

const fs = require('fs');
const path = require('path');
const glob = require('glob');

// Цвета для вывода
const colors = {
  red: '\x1b[31m',
  green: '\x1b[32m',
  yellow: '\x1b[33m',
  blue: '\x1b[34m',
  reset: '\x1b[0m',
};

let totalFixed = 0;
const fixedFiles = [];

// Паттерн для поиска неправильного использования
const nestedNamespacePattern =
  /useTranslations\s*\(\s*['"`]([^'"`]+\.[^'"`]+)['"`]\s*\)/g;

// Найти все файлы TypeScript/JavaScript
const files = glob.sync(
  path.join(__dirname, '..', 'src', '**', '*.{ts,tsx,js,jsx}'),
  {
    ignore: [
      '**/node_modules/**',
      '**/scripts/**',
      '**/*.test.*',
      '**/*.spec.*',
    ],
  }
);

console.log(
  `${colors.blue}Поиск файлов с неправильным использованием namespace...${colors.reset}\n`
);

files.forEach((file) => {
  const content = fs.readFileSync(file, 'utf8');
  let modified = false;
  let newContent = content;

  // Проверяем и исправляем
  const matches = content.match(nestedNamespacePattern);
  if (matches) {
    console.log(
      `${colors.yellow}Найдено в ${path.relative(process.cwd(), file)}:${colors.reset}`
    );

    matches.forEach((match) => {
      const namespace = match.match(/['"`]([^'"`]+)['"`]/)[1];
      const [module, ...rest] = namespace.split('.');

      console.log(`  ${colors.red}- ${match}${colors.reset}`);
      console.log(
        `  ${colors.green}+ useTranslations('${module}')${colors.reset}`
      );

      // Заменяем в контенте
      newContent = newContent.replace(match, `useTranslations('${module}')`);
      modified = true;
    });

    if (modified) {
      // Сохраняем файл
      fs.writeFileSync(file, newContent, 'utf8');
      fixedFiles.push(file);
      totalFixed++;
      console.log(`  ${colors.green}✓ Исправлено${colors.reset}\n`);
    }
  }
});

// Также нужно проверить использование переводов
console.log(
  `\n${colors.blue}Проверка использования переводов после namespace...${colors.reset}\n`
);

fixedFiles.forEach((file) => {
  const content = fs.readFileSync(file, 'utf8');

  // Находим переменные от useTranslations
  const translationVars = [];
  const varPattern =
    /const\s+(\w+)\s*=\s*useTranslations\s*\(\s*['"`](\w+)['"`]\s*\)/g;
  let match;

  while ((match = varPattern.exec(content)) !== null) {
    translationVars.push({
      varName: match[1],
      module: match[2],
    });
  }

  // Проверяем использование
  translationVars.forEach(({ varName, module }) => {
    const usagePattern = new RegExp(
      `${varName}\\s*\\(\\s*['"\`]([^'"\`]+\\.[^'"\`]+)['"\`]`,
      'g'
    );
    const usages = content.match(usagePattern);

    if (usages) {
      console.log(
        `${colors.yellow}⚠️  В файле ${path.relative(process.cwd(), file)} нужно проверить использование:${colors.reset}`
      );
      console.log(`   Модуль '${module}', переменная '${varName}'`);
      console.log(
        `   Найдены вложенные обращения, которые могут требовать ручной проверки`
      );
    }
  });
});

// Итоговый отчет
console.log(
  `\n${colors.green}========================================${colors.reset}`
);
console.log(`${colors.green}Исправлено файлов: ${totalFixed}${colors.reset}`);
console.log(
  `${colors.green}========================================${colors.reset}\n`
);

if (totalFixed > 0) {
  console.log('Исправленные файлы:');
  fixedFiles.forEach((file) => {
    console.log(`  - ${path.relative(process.cwd(), file)}`);
  });

  console.log(
    `\n${colors.yellow}ВАЖНО: Проверьте использование переводов в исправленных файлах!${colors.reset}`
  );
  console.log(
    `${colors.yellow}Если использовались вложенные ключи типа t('products.title'),${colors.reset}`
  );
  console.log(
    `${colors.yellow}их нужно изменить на t('products.title') без изменений.${colors.reset}`
  );
} else {
  console.log(
    `${colors.green}Файлов с неправильным использованием namespace не найдено!${colors.reset}`
  );
}
