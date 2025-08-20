#!/usr/bin/env node

const fs = require('fs');

const filePath =
  '/data/hostel-booking-system/frontend/svetu/src/app/[locale]/examples/novi-sad-districts/page.tsx';

console.log('🔧 Исправляем синтаксические ошибки после удаления районов...\n');

// Читаем файл
let content = fs.readFileSync(filePath, 'utf8');

// Удаляем пустые строки с одинокими закрывающими скобками между объектами районов
// Паттерн: }, пустые строки, }, {
content = content.replace(/\},\s*\n\s*\n\s*\},\s*\n\s*\{/g, '},\n    {');

// Удаляем множественные пустые }, между объектами
content = content.replace(/\},\s*\n\s*\},\s*\n\s*\{/g, '},\n    {');

// Удаляем лишние }, перед концом массива
content = content.replace(/\},\s*\n\s*\}\s*\n\s*\]/g, '}\n  ]');

// Удаляем множественные пустые строки
content = content.replace(/\n\s*\n\s*\n/g, '\n\n');

// Записываем обратно
fs.writeFileSync(filePath, content, 'utf8');

console.log('✅ Синтаксические ошибки исправлены!');

// Проверяем количество оставшихся районов
const districts = content.match(/id:\s*'/g);
const count = districts ? districts.length : 0;
console.log(`📍 Осталось районов: ${count}`);
