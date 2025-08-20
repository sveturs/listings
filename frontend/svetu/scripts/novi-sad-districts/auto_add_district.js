#!/usr/bin/env node

/**
 * Скрипт для автоматического добавления района на карту Нови-Сада
 * Использование: node auto_add_district.js <название_района> <город> <страна>
 * Пример: node auto_add_district.js Salajka "Novi Sad" Serbia
 */

const fs = require('fs');
const { exec } = require('child_process');
const util = require('util');
const execPromise = util.promisify(exec);

// Цвета для новых районов (будем брать по кругу)
const COLORS = [
  '#06B6D4', // cyan-500
  '#A855F7', // purple-500
  '#EC4899', // pink-500
  '#14B8A6', // teal-500
  '#F97316', // orange-500
  '#EAB308', // yellow-500
  '#22C55E', // green-500
  '#3B82F6', // blue-500
  '#EF4444', // red-500
];

// Путь к файлу с районами
const DISTRICTS_FILE =
  '/data/hostel-booking-system/frontend/svetu/src/app/[locale]/examples/novi-sad-districts/page.tsx';

async function searchDistrict(districtName, city, country) {
  console.log(`🔍 Поиск района "${districtName}" в ${city}, ${country}...`);

  const query = encodeURIComponent(`${districtName},${city},${country}`);
  const url = `https://nominatim.openstreetmap.org/search?q=${query}&format=json&addressdetails=1&limit=10`;

  try {
    const { stdout } = await execPromise(`curl -s "${url}"`);
    const results = JSON.parse(stdout);

    if (!results || results.length === 0) {
      throw new Error(`Район "${districtName}" не найден`);
    }

    // Ищем administrative relation
    const relation = results.find(
      (r) =>
        r.osm_type === 'relation' &&
        (r.type === 'administrative' ||
          r.type === 'suburb' ||
          r.type === 'neighbourhood')
    );

    if (!relation) {
      console.log(
        '⚠️  Administrative relation не найден, пробуем первый relation...'
      );
      const anyRelation = results.find((r) => r.osm_type === 'relation');
      if (!anyRelation) {
        throw new Error(`Не найден OSM relation для "${districtName}"`);
      }
      return anyRelation;
    }

    return relation;
  } catch (error) {
    throw new Error(`Ошибка поиска: ${error.message}`);
  }
}

async function getPolygon(osmId) {
  console.log(`📍 Получение полигона для OSM relation ${osmId}...`);

  const url = `https://nominatim.openstreetmap.org/lookup?osm_ids=R${osmId}&format=geojson&polygon_geojson=1`;

  try {
    const { stdout } = await execPromise(`curl -s "${url}"`);
    const geojson = JSON.parse(stdout);

    if (!geojson.features || geojson.features.length === 0) {
      throw new Error('Не удалось получить полигон');
    }

    const coordinates = geojson.features[0].geometry.coordinates[0];
    console.log(`✅ Получен полигон с ${coordinates.length} точками`);

    return coordinates;
  } catch (error) {
    throw new Error(`Ошибка получения полигона: ${error.message}`);
  }
}

function generateDistrictCode(name, nameEn, coordinates, osmId, color) {
  // Генерируем ID из английского названия
  const id = nameEn.toLowerCase().replace(/\s+/g, '-');

  // Форматируем координаты для JavaScript (только первые и последние 10 точек для краткости)
  let coordsStr;
  if (coordinates.length <= 20) {
    coordsStr = coordinates
      .map((c) => `        [${c[0]}, ${c[1]}]`)
      .join(',\n');
  } else {
    const first10 = coordinates
      .slice(0, 10)
      .map((c) => `        [${c[0]}, ${c[1]}]`)
      .join(',\n');
    const last10 = coordinates
      .slice(-10)
      .map((c) => `        [${c[0]}, ${c[1]}]`)
      .join(',\n');
    coordsStr = `${first10},\n        // ... еще ${coordinates.length - 20} точек опущено для краткости ...\n${last10}`;
  }

  return `    {
      id: '${id}',
      name: '${name}',
      nameEn: '${nameEn}',
      population: 'Неизвестно',
      area: 'Неизвестно',
      color: '${color}',
      description: 'Автоматически добавлен из OSM relation ${osmId} (${coordinates.length} точек)',
      coordinates: [
${coordsStr}
      ],
      stats: {
        properties: 0,
        avgPrice: '€0',
        growth: '+0%',
      },
    },`;
}

async function addDistrictToFile(districtCode) {
  console.log('📝 Добавление района в файл...');

  try {
    let content = fs.readFileSync(DISTRICTS_FILE, 'utf8');

    // Находим массив districts
    const districtsStart = content.indexOf('const districts = [');
    if (districtsStart === -1) {
      throw new Error('Не найден массив districts в файле');
    }

    // Находим последний район перед закрытием массива
    const lastDistrictPattern = /\},\s*\{[^}]*id:\s*'salajka'/s;
    const match = content.match(lastDistrictPattern);

    if (!match) {
      throw new Error('Не найден последний район (salajka) в массиве');
    }

    const insertPosition = content.indexOf(match[0]) + match[0].length;

    // Находим конец объекта salajka
    let bracketCount = 1;
    let i = insertPosition;
    while (i < content.length && bracketCount > 0) {
      if (content[i] === '{') bracketCount++;
      if (content[i] === '}') bracketCount--;
      i++;
    }

    // Вставляем новый район после запятой
    const beforeInsert = content.substring(0, i);
    const afterInsert = content.substring(i);

    content = beforeInsert + ',\n' + districtCode + afterInsert;

    fs.writeFileSync(DISTRICTS_FILE, content);
    console.log('✅ Район успешно добавлен в файл!');
  } catch (error) {
    throw new Error(`Ошибка записи в файл: ${error.message}`);
  }
}

async function getExistingDistrictsCount() {
  try {
    const content = fs.readFileSync(DISTRICTS_FILE, 'utf8');
    const matches = content.match(/id:\s*'/g);
    return matches ? matches.length : 0;
  } catch {
    return 0;
  }
}

async function main() {
  const args = process.argv.slice(2);

  if (args.length < 3) {
    console.log(
      'Использование: node auto_add_district.js <название_района> <город> <страна>'
    );
    console.log('Пример: node auto_add_district.js Salajka "Novi Sad" Serbia');
    console.log(
      'Пример с пробелом: node auto_add_district.js "Slana bara" "Novi Sad" Serbia'
    );
    process.exit(1);
  }

  // Обрабатываем возможность названий с пробелами
  let districtName, city, country;

  // Если больше 3 аргументов и нет кавычек, возможно название из нескольких слов
  if (args.length > 3) {
    // Ищем где начинается город (обычно это Novi)
    const noviIndex = args.findIndex(
      (arg) => arg === 'Novi' || arg === '"Novi'
    );
    if (noviIndex > 0) {
      districtName = args.slice(0, noviIndex).join(' ');
      city = args
        .slice(noviIndex, noviIndex + 2)
        .join(' ')
        .replace(/"/g, '');
      country = args[noviIndex + 2] || args[args.length - 1];
    } else {
      // Если не нашли, берем последние 2 как город и страну
      country = args[args.length - 1];
      city = args[args.length - 2];
      districtName = args.slice(0, -2).join(' ');
    }
  } else {
    [districtName, city, country] = args;
  }

  // Убираем кавычки если есть
  districtName = districtName.replace(/"/g, '');
  city = city.replace(/"/g, '');
  country = country.replace(/"/g, '');

  try {
    // 1. Поиск района
    const searchResult = await searchDistrict(districtName, city, country);
    console.log(`✅ Найден: ${searchResult.display_name}`);
    console.log(`   OSM ID: ${searchResult.osm_id}`);
    console.log(`   Тип: ${searchResult.type}`);

    // 2. Получение полигона
    const coordinates = await getPolygon(searchResult.osm_id);

    // 3. Определение цвета (циклически из массива)
    const districtsCount = await getExistingDistrictsCount();
    const color = COLORS[districtsCount % COLORS.length];

    // 4. Генерация кода
    const nameEn = districtName; // Можно улучшить транслитерацию
    const districtCode = generateDistrictCode(
      districtName,
      nameEn,
      coordinates,
      searchResult.osm_id,
      color
    );

    // 5. Добавление в файл
    await addDistrictToFile(districtCode);

    console.log('\n🎉 Готово! Район добавлен на карту.');
    console.log('   Перезагрузите страницу чтобы увидеть изменения.');
  } catch (error) {
    console.error(`\n❌ Ошибка: ${error.message}`);
    process.exit(1);
  }
}

// Запуск
main();
