#!/usr/bin/env node

const fs = require('fs');
const path = require('path');

/**
 * Безопасный скрипт для удаления районов из карты Нови-Сада
 * Корректно обрабатывает структуру объектов и синтаксис массива
 */

const filePath =
  '/data/hostel-booking-system/frontend/svetu/src/app/[locale]/examples/novi-sad-districts/page.tsx';

// Получаем список районов для удаления из аргументов командной строки
const args = process.argv.slice(2);

if (args.length === 0) {
  console.log(
    '🔧 Использование: node safe_remove_districts.js <id1> <id2> ...'
  );
  console.log(
    'Пример: node safe_remove_districts.js petrovaradin veternik futog'
  );
  process.exit(1);
}

const districtsToRemove = args;

console.log('🗑️  Безопасное удаление районов из файла...\n');
console.log('Районы для удаления:', districtsToRemove.join(', '));
console.log('=' + '='.repeat(60));

// Читаем файл
let content = fs.readFileSync(filePath, 'utf8');
const originalLength = content.length;

// Сохраняем бэкап
const backupPath = filePath + '.backup_' + Date.now();
fs.writeFileSync(backupPath, content);
console.log(`📁 Создан бэкап: ${backupPath}\n`);

/**
 * Функция для безопасного удаления района
 * Находит полный объект района от { до закрывающей } включая все вложенные структуры
 */
function removeDistrictSafely(content, districtId) {
  // Сначала находим начало объекта с нужным id
  const idPattern = new RegExp(
    `\\{\\s*\\n\\s*id:\\s*['"\`]${districtId}['"\`]`,
    'g'
  );

  let match;
  let removed = false;

  while ((match = idPattern.exec(content)) !== null) {
    const startPos = match.index;

    // Теперь нужно найти закрывающую скобку для этого объекта
    // Считаем баланс скобок
    let bracketCount = 1; // Начинаем с 1, так как уже нашли открывающую
    let currentPos = match.index + match[0].length;
    let insideString = false;
    let stringChar = null;

    while (currentPos < content.length && bracketCount > 0) {
      const char = content[currentPos];
      const prevChar = currentPos > 0 ? content[currentPos - 1] : '';

      // Обработка строк (чтобы не считать скобки внутри строк)
      if (
        !insideString &&
        (char === '"' || char === "'" || char === '`') &&
        prevChar !== '\\'
      ) {
        insideString = true;
        stringChar = char;
      } else if (insideString && char === stringChar && prevChar !== '\\') {
        insideString = false;
        stringChar = null;
      }

      // Считаем скобки только вне строк
      if (!insideString) {
        if (char === '{') {
          bracketCount++;
        } else if (char === '}') {
          bracketCount--;
        }
      }

      currentPos++;
    }

    if (bracketCount === 0) {
      // Нашли закрывающую скобку
      // Теперь проверяем, есть ли запятая после объекта
      let endPos = currentPos;

      // Пропускаем пробелы после }
      while (endPos < content.length && /\s/.test(content[endPos])) {
        endPos++;
      }

      // Если есть запятая, включаем её в удаление
      if (content[endPos] === ',') {
        endPos++;
      }

      // Удаляем объект
      const before = content.substring(0, startPos);
      const after = content.substring(endPos);

      // Проверяем, не создаём ли мы двойную запятую
      // Если перед удаляемым объектом есть запятая, а после тоже объект, нужно быть осторожным
      const beforeTrimmed = before.trimEnd();
      const afterTrimmed = after.trimStart();

      if (beforeTrimmed.endsWith(',') && afterTrimmed.startsWith('{')) {
        // Всё в порядке, просто соединяем
        content = before + after;
      } else if (beforeTrimmed.endsWith(',') && afterTrimmed.startsWith(']')) {
        // Удаляем последний элемент массива, нужно убрать запятую перед
        content = before.trimEnd().slice(0, -1) + '\n  ' + after;
      } else {
        content = before + after;
      }

      removed = true;
      console.log(`✅ Удален район: ${districtId}`);
      break;
    }
  }

  if (!removed) {
    console.log(`⚠️  Район не найден: ${districtId}`);
  }

  return { content, removed };
}

let removedCount = 0;
let notFoundCount = 0;

// Удаляем каждый район
for (const districtId of districtsToRemove) {
  const result = removeDistrictSafely(content, districtId);
  content = result.content;

  if (result.removed) {
    removedCount++;
  } else {
    notFoundCount++;
  }
}

// Финальная очистка синтаксиса
function cleanupSyntax(content) {
  // Удаляем множественные запятые
  content = content.replace(/,\s*,+/g, ',');

  // Удаляем запятую перед закрывающей скобкой массива
  content = content.replace(/,\s*\n\s*\]/g, '\n  ]');

  // Удаляем лишние пустые строки (больше 2 подряд)
  content = content.replace(/\n\s*\n\s*\n+/g, '\n\n');

  // Исправляем отступы для закрывающей скобки массива
  content = content.replace(/\n\s*\];\s*$/m, '\n  ];');

  return content;
}

content = cleanupSyntax(content);

// Проверяем корректность синтаксиса
function validateSyntax(content) {
  let issues = [];

  // Проверка 1: Баланс скобок
  let brackets = { '{': 0, '}': 0, '[': 0, ']': 0 };
  let insideString = false;
  let stringChar = null;

  for (let i = 0; i < content.length; i++) {
    const char = content[i];
    const prevChar = i > 0 ? content[i - 1] : '';

    if (
      !insideString &&
      (char === '"' || char === "'" || char === '`') &&
      prevChar !== '\\'
    ) {
      insideString = true;
      stringChar = char;
    } else if (insideString && char === stringChar && prevChar !== '\\') {
      insideString = false;
      stringChar = null;
    }

    if (!insideString) {
      if (brackets.hasOwnProperty(char)) {
        brackets[char]++;
      }
    }
  }

  if (brackets['{'] !== brackets['}']) {
    issues.push(
      `❌ Несбалансированные фигурные скобки: { = ${brackets['{']}, } = ${brackets['}']}`
    );
  }
  if (brackets['['] !== brackets[']']) {
    issues.push(
      `❌ Несбалансированные квадратные скобки: [ = ${brackets['[']}, ] = ${brackets[']']}`
    );
  }

  // Проверка 2: Двойные запятые
  if (/,\s*,/.test(content)) {
    issues.push('❌ Обнаружены двойные запятые');
  }

  // Проверка 3: Запятая перед ]
  if (/,\s*\]/.test(content)) {
    issues.push('❌ Обнаружена запятая перед закрывающей скобкой массива');
  }

  return issues;
}

const syntaxIssues = validateSyntax(content);

// Записываем результат
fs.writeFileSync(filePath, content, 'utf8');

const bytesRemoved = originalLength - content.length;

// Подсчитываем оставшиеся районы
const remainingDistricts = content.match(/id:\s*['"`]/g);
const remainingCount = remainingDistricts ? remainingDistricts.length : 0;

console.log('\n' + '='.repeat(60));
console.log('📊 ИТОГ:');
console.log(`✅ Успешно удалено: ${removedCount} районов`);
console.log(`⚠️  Не найдено: ${notFoundCount} районов`);
console.log(`📉 Удалено байт: ${bytesRemoved}`);
console.log(`📍 Осталось районов: ${remainingCount}`);
console.log('='.repeat(60));

if (syntaxIssues.length > 0) {
  console.log('\n⚠️  ПРЕДУПРЕЖДЕНИЯ О СИНТАКСИСЕ:');
  syntaxIssues.forEach((issue) => console.log(issue));
  console.log('\n💡 Возможно, требуется ручная проверка файла.');
} else {
  console.log('\n✅ Синтаксис файла корректен!');
}

if (removedCount > 0) {
  console.log('\n🎉 Районы успешно удалены!');
  console.log(
    '📝 Теперь можно добавлять их заново с правильными координатами:'
  );
  console.log(
    '   node /tmp/auto_add_district.js "Название" "Novi Sad" Serbia\n'
  );
}

// Информация о восстановлении
console.log('💾 Для восстановления из бэкапа используйте:');
console.log(`   cp ${backupPath} ${filePath}\n`);
