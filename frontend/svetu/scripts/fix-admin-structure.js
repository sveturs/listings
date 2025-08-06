const fs = require('fs');
const path = require('path');

// Функция для исправления структуры admin.json
function fixAdminStructure(locale) {
  const filePath = path.join(__dirname, `../src/messages/${locale}/admin.json`);

  console.log(`Обработка файла: ${filePath}`);

  const content = fs.readFileSync(filePath, 'utf8');
  const data = JSON.parse(content);

  // Проверяем, есть ли обертка "admin"
  if (data.admin && typeof data.admin === 'object') {
    console.log(`Найдена обертка "admin" в ${locale}/admin.json, удаляем...`);

    // Извлекаем содержимое из обертки
    const adminContent = data.admin;

    // Перемещаем все остальные ключи (если есть) в корень
    const otherKeys = {};
    for (const key in data) {
      if (key !== 'admin') {
        otherKeys[key] = data[key];
      }
    }

    // Создаем новую структуру без обертки
    const newData = {
      ...adminContent,
      ...otherKeys,
    };

    // Записываем обратно
    fs.writeFileSync(filePath, JSON.stringify(newData, null, 2) + '\n', 'utf8');
    console.log(`Файл ${locale}/admin.json исправлен успешно`);
  } else {
    console.log(`Файл ${locale}/admin.json уже имеет правильную структуру`);
  }
}

// Обрабатываем все языки
['ru', 'en', 'sr'].forEach((locale) => {
  try {
    fixAdminStructure(locale);
  } catch (error) {
    console.error(`Ошибка при обработке ${locale}:`, error);
  }
});

console.log('\nИсправление структуры admin.json завершено!');
