#!/usr/bin/env node

const fs = require('fs');
const path = require('path');
const glob = require('glob');

console.log('🔍 Поиск файлов с useTranslations() без модуля...\n');

// Паттерны для поиска
const patterns = [
  'src/**/*.tsx',
  'src/**/*.ts',
];

let totalFiles = 0;
let fixedFiles = 0;
const problematicFiles = [];

patterns.forEach(pattern => {
  const files = glob.sync(pattern, {
    cwd: path.join(__dirname, '..'),
    ignore: ['**/node_modules/**', '**/*.test.*', '**/*.spec.*']
  });

  files.forEach(file => {
    totalFiles++;
    const filePath = path.join(__dirname, '..', file);
    let content = fs.readFileSync(filePath, 'utf8');
    const originalContent = content;

    // Проверяем, есть ли useTranslations() без аргумента
    if (content.includes('useTranslations()')) {
      console.log(`\n📄 ${file}`);
      
      // Анализируем, какие модули используются через t('')
      const moduleUsage = new Map();
      
      // Ищем все использования t('module.key')
      const tUsageRegex = /t\(['"]([^'"]+)\./g;
      let match;
      while ((match = tUsageRegex.exec(content)) !== null) {
        const module = match[1];
        moduleUsage.set(module, (moduleUsage.get(module) || 0) + 1);
      }

      if (moduleUsage.size > 0) {
        // Определяем основной модуль (наиболее используемый)
        let mainModule = '';
        let maxUsage = 0;
        moduleUsage.forEach((count, module) => {
          if (count > maxUsage) {
            maxUsage = count;
            mainModule = module;
          }
        });

        console.log(`  ⚡ Обнаружены модули: ${Array.from(moduleUsage.keys()).join(', ')}`);
        console.log(`  📦 Основной модуль: ${mainModule} (${maxUsage} использований)`);

        // Заменяем useTranslations()
        content = content.replace(
          /const\s+(\w+)\s*=\s*useTranslations\(\)/g,
          `const $1 = useTranslations('${mainModule}')`
        );

        // Если есть другие модули, добавляем их
        const otherModules = Array.from(moduleUsage.keys()).filter(m => m !== mainModule);
        if (otherModules.length > 0) {
          // Находим место после первого useTranslations
          const insertIndex = content.indexOf(`useTranslations('${mainModule}')`) + `useTranslations('${mainModule}')`.length;
          const lineEnd = content.indexOf('\n', insertIndex);
          
          const additionalTranslations = otherModules.map(module => {
            const varName = module === 'common' ? 'tCommon' : 
                           module === 'admin' ? 'tAdmin' :
                           module === 'auth' ? 'tAuth' :
                           `t${module.charAt(0).toUpperCase() + module.slice(1)}`;
            return `\n  const ${varName} = useTranslations('${module}');`;
          }).join('');

          content = content.slice(0, lineEnd) + additionalTranslations + content.slice(lineEnd);

          console.log(`  ➕ Добавлены дополнительные модули: ${otherModules.join(', ')}`);
        }

        // Теперь обновляем использование t('module.key')
        otherModules.forEach(module => {
          const varName = module === 'common' ? 'tCommon' : 
                         module === 'admin' ? 'tAdmin' :
                         module === 'auth' ? 'tAuth' :
                         `t${module.charAt(0).toUpperCase() + module.slice(1)}`;
          
          // Заменяем t('module.key') на tModule('key')
          const regex = new RegExp(`t\\(['"]${module}\\.([^'"]+)['"]\\)`, 'g');
          content = content.replace(regex, `${varName}('$1')`);
        });

        // Обновляем использование основного модуля
        const mainModuleRegex = new RegExp(`t\\(['"]${mainModule}\\.([^'"]+)['"]\\)`, 'g');
        content = content.replace(mainModuleRegex, "t('$1')");

        if (content !== originalContent) {
          fs.writeFileSync(filePath, content);
          fixedFiles++;
          console.log('  ✅ Исправлено!');
        }
      } else {
        // Не можем определить модуль автоматически
        problematicFiles.push(file);
        console.log('  ⚠️  Не удалось определить модуль автоматически');
      }
    }
  });
});

console.log('\n📊 Результаты:');
console.log(`   Проверено файлов: ${totalFiles}`);
console.log(`   Исправлено файлов: ${fixedFiles}`);
console.log(`   Требуют ручного исправления: ${problematicFiles.length}`);

if (problematicFiles.length > 0) {
  console.log('\n⚠️  Файлы, требующие ручного исправления:');
  problematicFiles.forEach(file => {
    console.log(`   - ${file}`);
  });
}

console.log('\n✨ Готово!');