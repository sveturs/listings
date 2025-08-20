#!/usr/bin/env node

const https = require('https');
const fs = require('fs');
const path = require('path');

// Функция для выполнения запроса к Overpass API
async function queryOverpass(query) {
  return new Promise((resolve, reject) => {
    const postData = `data=${encodeURIComponent(query)}`;

    const options = {
      hostname: 'overpass-api.de',
      port: 443,
      path: '/api/interpreter',
      method: 'POST',
      headers: {
        'Content-Type': 'application/x-www-form-urlencoded',
        'Content-Length': Buffer.byteLength(postData),
      },
    };

    const req = https.request(options, (res) => {
      let data = '';

      res.on('data', (chunk) => {
        data += chunk;
      });

      res.on('end', () => {
        try {
          const result = JSON.parse(data);
          resolve(result);
        } catch (e) {
          reject(new Error(`Failed to parse Overpass response: ${e.message}`));
        }
      });
    });

    req.on('error', (e) => {
      reject(e);
    });

    req.write(postData);
    req.end();
  });
}

// Функция для соединения way-элементов multipolygon в правильном порядке
function assembleMultipolygonCoordinates(members) {
  const outerWays = members.filter(
    (member) =>
      member.role === 'outer' && member.geometry && member.geometry.length > 0
  );

  if (outerWays.length === 0) {
    return [];
  }

  if (outerWays.length === 1) {
    // Простой случай - один way
    return outerWays[0].geometry.map((coord) => [coord.lon, coord.lat]);
  }

  // Сложный случай - множественные way нужно соединить
  console.log(`🔗 Соединение ${outerWays.length} way-элементов в multipolygon`);

  const allCoordinates = [];
  const usedWays = new Set();

  // Начинаем с первого way
  let currentWay = outerWays[0];
  usedWays.add(0);
  allCoordinates.push(
    ...currentWay.geometry.map((coord) => [coord.lon, coord.lat])
  );

  // Получаем последнюю точку текущего way
  let lastPoint = currentWay.geometry[currentWay.geometry.length - 1];

  // Ищем следующий way, который начинается или заканчивается в lastPoint
  while (usedWays.size < outerWays.length) {
    let foundNextWay = false;

    for (let i = 0; i < outerWays.length; i++) {
      if (usedWays.has(i)) continue;

      const way = outerWays[i];
      const firstPoint = way.geometry[0];
      const wayLastPoint = way.geometry[way.geometry.length - 1];

      // Проверяем, стыкуется ли way с текущей позицией
      const tolerance = 0.000001; // Допуск для сравнения координат

      if (
        Math.abs(firstPoint.lat - lastPoint.lat) < tolerance &&
        Math.abs(firstPoint.lon - lastPoint.lon) < tolerance
      ) {
        // Way начинается там, где закончился предыдущий - добавляем как есть (пропуская первую точку)
        allCoordinates.push(
          ...way.geometry.slice(1).map((coord) => [coord.lon, coord.lat])
        );
        lastPoint = wayLastPoint;
        usedWays.add(i);
        foundNextWay = true;
        break;
      } else if (
        Math.abs(wayLastPoint.lat - lastPoint.lat) < tolerance &&
        Math.abs(wayLastPoint.lon - lastPoint.lon) < tolerance
      ) {
        // Way заканчивается там, где закончился предыдущий - добавляем в обратном порядке (пропуская последнюю точку)
        allCoordinates.push(
          ...way.geometry
            .slice(0, -1)
            .reverse()
            .map((coord) => [coord.lon, coord.lat])
        );
        lastPoint = firstPoint;
        usedWays.add(i);
        foundNextWay = true;
        break;
      }
    }

    if (!foundNextWay) {
      console.log(
        `⚠️ Не удалось найти следующий way для соединения. Использовано ${usedWays.size} из ${outerWays.length}`
      );
      break;
    }
  }

  console.log(
    `✅ Собрано ${allCoordinates.length} координат из ${usedWays.size} way-элементов`
  );
  return allCoordinates;
}

// Функция для получения координат района через Overpass API
async function getDistrictFromOverpass(
  districtName,
  osmId = null,
  osmType = null
) {
  console.log(`🔍 Поиск района "${districtName}" через Overpass API...`);

  let query = '';

  if (osmId && osmType) {
    // Если есть OSM ID и тип, используем их напрямую
    if (osmType === 'node') {
      // Для node нужен специальный запрос
      query = `[out:json][timeout:25];
        node(${osmId});
        out;`;
    } else {
      query = `[out:json][timeout:25];
        ${osmType}(${osmId});
        out geom;`;
    }
  } else {
    // Иначе ищем по имени
    query = `[out:json][timeout:25];
      area[name="Novi Sad"]->.novisad;
      (
        relation(area.novisad)["name"~"${districtName}",i]["place"];
        relation(area.novisad)["name:en"~"${districtName}",i]["place"];
        relation(area.novisad)["name:sr"~"${districtName}",i]["place"];
      );
      out geom;`;
  }

  try {
    const result = await queryOverpass(query);

    if (!result.elements || result.elements.length === 0) {
      console.log('❌ Район не найден в Overpass API');
      return null;
    }

    const element = result.elements[0];
    console.log(`✅ Найден район: ${element.tags?.name || districtName}`);
    console.log(`   OSM ID: ${element.id}, Тип: ${element.type}`);
    console.log(`   Place: ${element.tags?.place || 'unknown'}`);

    // Извлекаем координаты в зависимости от типа объекта
    let coordinates = [];

    if (element.type === 'relation' && element.members) {
      // Для relation используем новую функцию для правильного соединения way-элементов
      coordinates = assembleMultipolygonCoordinates(element.members);
    } else if (element.type === 'way' && element.geometry) {
      // Для way используем geometry напрямую
      coordinates = element.geometry.map((coord) => [coord.lon, coord.lat]);
    } else if (element.type === 'node' && element.lat && element.lon) {
      // Node типы не подходят для районов - они представляют точки, а не области
      console.log(
        '⚠️ Node тип не подходит для района - нужны way или relation с полигонами'
      );
      return null;
    }

    if (coordinates.length === 0) {
      console.log('⚠️ Не удалось извлечь координаты из результата');
      console.log('Элемент:', JSON.stringify(element, null, 2));
      return null;
    }

    console.log(`📍 Найдено ${coordinates.length} координат`);

    return {
      id: element.id,
      type: element.type,
      name: element.tags?.name || districtName,
      nameEn: element.tags?.['name:en'] || districtName,
      nameSr: element.tags?.['name:sr'] || element.tags?.name || districtName,
      place: element.tags?.place || 'quarter',
      coordinates: coordinates,
      tags: element.tags,
    };
  } catch (error) {
    console.error('❌ Ошибка при запросе к Overpass API:', error.message);
    return null;
  }
}

// Функция для добавления района в файл
function addDistrictToFile(district) {
  const filePath = path.join(
    __dirname,
    '../../src/app/[locale]/examples/novi-sad-districts/page.tsx'
  );

  // Создаем резервную копию
  const backupPath = filePath + `.backup.${Date.now()}`;
  fs.copyFileSync(filePath, backupPath);
  console.log(`📁 Создана резервная копия: ${backupPath}`);

  // Читаем файл
  let content = fs.readFileSync(filePath, 'utf8');

  // Генерируем ID для района
  let districtId =
    (district.nameEn || district.name)
      .toLowerCase()
      .replace(/\s+/g, '-')
      .replace(/[^\w\-]/g, '')
      .replace(/^-+|-+$/g, '') || // удаляем дефисы в начале и конце
    `district-${district.id}`; // fallback на OSM ID если не получилось

  // Проверяем, существует ли уже район с таким ID
  if (content.includes(`id: '${districtId}'`)) {
    console.log(
      `⚠️ Район с ID '${districtId}' уже существует, добавляю суффикс`
    );
    districtId = `${districtId}-${district.id}`; // добавляем OSM ID как суффикс для уникальности
  }

  // Генерируем случайный цвет
  const colors = [
    '#10B981',
    '#F59E0B',
    '#EF4444',
    '#8B5CF6',
    '#F97316',
    '#22D3EE',
    '#3B82F6',
    '#22C55E',
    '#EAB308',
    '#14B8A6',
    '#A855F7',
    '#06B6D4',
    '#EC4899',
  ];
  const color = colors[Math.floor(Math.random() * colors.length)];

  // Создаем объект района
  const newDistrictCode = `
    {
      id: '${districtId}',
      name: '${district.nameSr || district.name}',
      nameEn: '${district.nameEn}',
      population: 'Неизвестно',
      area: 'Неизвестно',
      color: '${color}',
      description: 'Автоматически добавлен из OSM ${district.type} ${district.id} через Overpass API (${district.coordinates.length} точек)',
      coordinates: [
${district.coordinates
  .slice(0, -1)
  .map((coord) => `        [${coord[0]}, ${coord[1]}]`)
  .join(',\n')}
      ],
      stats: {
        properties: 0,
        avgPrice: '€0',
        growth: '+0%',
      }
    },`;

  // Находим место для вставки (перед закрывающей скобкой массива districts)
  const insertPosition = content.lastIndexOf('];');

  if (insertPosition === -1) {
    console.error('❌ Не удалось найти место для вставки в файле');
    return false;
  }

  // Проверяем, есть ли уже районы в массиве (нужно ли добавить запятую после предыдущего)
  const beforeInsert = content.slice(0, insertPosition).trimEnd();
  let finalContent;

  // Если перед позицией вставки есть закрывающая фигурная скобка без запятой - добавляем запятую
  if (beforeInsert.endsWith('}') && !beforeInsert.endsWith('},')) {
    // Добавляем запятую после последнего существующего района
    const lastBracePos = beforeInsert.lastIndexOf('}');
    finalContent =
      beforeInsert.slice(0, lastBracePos + 1) +
      ',' +
      beforeInsert.slice(lastBracePos + 1);
    content =
      finalContent + newDistrictCode + '\n  ' + content.slice(insertPosition);
  } else {
    // Вставляем новый район как есть
    content =
      content.slice(0, insertPosition) +
      newDistrictCode +
      '\n  ' +
      content.slice(insertPosition);
  }

  // Записываем обновленный файл
  fs.writeFileSync(filePath, content);
  console.log(`✅ Район "${district.name}" успешно добавлен в файл!`);
  console.log(`   ID: ${districtId}`);
  console.log(`   Координат: ${district.coordinates.length}`);

  return true;
}

// Главная функция
async function main() {
  const args = process.argv.slice(2);

  if (args.length < 1) {
    console.log('Использование:');
    console.log(
      '  node add_district_overpass.js <имя_района> [osm_id] [osm_type]'
    );
    console.log('');
    console.log('Примеры:');
    console.log('  node add_district_overpass.js "Liman 1"');
    console.log('  node add_district_overpass.js "Лиман 2"');
    console.log('  node add_district_overpass.js "Veternik" 9274293 relation');
    process.exit(1);
  }

  const districtName = args[0];
  const osmId = args[1] || null;
  const osmType = args[2] || null;

  console.log(`🚀 Добавление района "${districtName}" через Overpass API`);

  // Получаем данные района
  const district = await getDistrictFromOverpass(districtName, osmId, osmType);

  if (!district) {
    console.error('❌ Не удалось получить данные района');
    process.exit(1);
  }

  // Добавляем район в файл
  const success = addDistrictToFile(district);

  if (success) {
    console.log('🎉 Готово! Район успешно добавлен.');
    console.log(
      '📝 Не забудьте перезапустить сервер разработки для отображения изменений.'
    );
  } else {
    console.error('❌ Не удалось добавить район в файл');
    process.exit(1);
  }
}

// Запуск
main().catch((error) => {
  console.error('❌ Критическая ошибка:', error);
  process.exit(1);
});
