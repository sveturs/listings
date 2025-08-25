#!/usr/bin/env node

/**
 * Улучшенный скрипт для получения точных границ районов крупных городов Сербии
 * Использует многоуровневую fallback систему для максимальной точности данных
 *
 * Использование:
 * node enhanced_district_fetcher.js "Veternik" "Novi Sad" --cache --validate
 * node enhanced_district_fetcher.js --city "Beograd" --fetch-all
 * node enhanced_district_fetcher.js --import-from-hdx serbia-admin-boundaries.geojson
 */

const fs = require('fs');
const https = require('https');
const http = require('http');
const path = require('path');
const crypto = require('crypto');

// Конфигурация
const CONFIG = {
  // Пути к файлам
  DISTRICTS_FILE:
    '/data/hostel-booking-system/frontend/svetu/src/app/[locale]/examples/novi-sad-districts/page.tsx',
  CACHE_DIR: '/tmp/districts_cache',

  // API конфигурация
  NOMINATIM_BASE: 'https://nominatim.openstreetmap.org',
  OVERPASS_BASE: 'https://overpass-api.de/api/interpreter',
  MAPBOX_BASE: 'https://api.mapbox.com',

  // Лимиты и таймауты
  REQUEST_TIMEOUT: 30000,
  RATE_LIMIT_DELAY: 1000,
  MAX_RETRIES: 3,

  // Валидация
  MIN_COORDINATE_COUNT: 4,
  MAX_COORDINATE_COUNT: 10000,
  MIN_AREA_SIZE: 0.0001, // квадратные градусы
  MAX_AREA_SIZE: 10.0,

  // Цвета для новых районов
  COLORS: [
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
  ],
};

class DistrictFetcher {
  constructor(options = {}) {
    this.options = { cache: true, validate: true, ...options };
    this.initCache();
  }

  // Инициализация кэша
  initCache() {
    if (!fs.existsSync(CONFIG.CACHE_DIR)) {
      fs.mkdirSync(CONFIG.CACHE_DIR, { recursive: true });
    }
  }

  // Генерация ключа кэша
  getCacheKey(query) {
    return crypto.createHash('md5').update(JSON.stringify(query)).digest('hex');
  }

  // Получение данных из кэша
  getFromCache(key) {
    if (!this.options.cache) return null;

    const cacheFile = path.join(CONFIG.CACHE_DIR, `${key}.json`);
    try {
      if (fs.existsSync(cacheFile)) {
        const cached = JSON.parse(fs.readFileSync(cacheFile, 'utf8'));
        const age = Date.now() - cached.timestamp;
        // Кэш действителен 7 дней
        if (age < 7 * 24 * 60 * 60 * 1000) {
          console.log(`📦 Загружено из кэша: ${key}`);
          return cached.data;
        }
      }
    } catch (error) {
      console.warn(`⚠️ Ошибка чтения кэша: ${error.message}`);
    }
    return null;
  }

  // Сохранение в кэш
  saveToCache(key, data) {
    if (!this.options.cache) return;

    const cacheFile = path.join(CONFIG.CACHE_DIR, `${key}.json`);
    try {
      fs.writeFileSync(
        cacheFile,
        JSON.stringify(
          {
            timestamp: Date.now(),
            data: data,
          },
          null,
          2
        )
      );
    } catch (error) {
      console.warn(`⚠️ Ошибка записи в кэш: ${error.message}`);
    }
  }

  // HTTP запрос с retry логикой
  async makeRequest(url, options = {}, retries = CONFIG.MAX_RETRIES) {
    return new Promise((resolve, reject) => {
      const protocol = url.startsWith('https') ? https : http;
      const isPost = options.method === 'POST';

      const req = protocol.request(
        url,
        {
          method: options.method || 'GET',
          headers: {
            'User-Agent': 'SvetU-District-Fetcher/1.0',
            ...options.headers,
          },
          timeout: CONFIG.REQUEST_TIMEOUT,
        },
        (res) => {
          let data = '';

          res.on('data', (chunk) => (data += chunk));
          res.on('end', () => {
            try {
              if (res.statusCode >= 200 && res.statusCode < 300) {
                resolve(data.trim() ? JSON.parse(data) : {});
              } else if (retries > 0 && res.statusCode >= 500) {
                // Повторяем запрос для серверных ошибок
                console.log(
                  `🔄 Повтор запроса через ${CONFIG.RATE_LIMIT_DELAY}мс, осталось попыток: ${retries}`
                );
                setTimeout(() => {
                  this.makeRequest(url, options, retries - 1)
                    .then(resolve)
                    .catch(reject);
                }, CONFIG.RATE_LIMIT_DELAY);
              } else {
                reject(new Error(`HTTP ${res.statusCode}: ${data}`));
              }
            } catch (parseError) {
              if (retries > 0) {
                setTimeout(() => {
                  this.makeRequest(url, options, retries - 1)
                    .then(resolve)
                    .catch(reject);
                }, CONFIG.RATE_LIMIT_DELAY);
              } else {
                reject(new Error(`JSON Parse Error: ${parseError.message}`));
              }
            }
          });
        }
      );

      req.on('error', (error) => {
        if (retries > 0) {
          setTimeout(() => {
            this.makeRequest(url, options, retries - 1)
              .then(resolve)
              .catch(reject);
          }, CONFIG.RATE_LIMIT_DELAY);
        } else {
          reject(error);
        }
      });

      if (isPost && options.data) {
        req.write(options.data);
      }

      req.end();
    });
  }

  // 1. Mapbox Geocoding API (с поддержкой полигонов)
  async fetchFromMapbox(districtName, city, mapboxToken) {
    if (!mapboxToken) {
      console.log('⚠️ MAPBOX_TOKEN не предоставлен, пропускаем Mapbox API');
      return null;
    }

    const cacheKey = this.getCacheKey({
      source: 'mapbox',
      district: districtName,
      city,
    });
    const cached = this.getFromCache(cacheKey);
    if (cached) return cached;

    console.log(`🗺️ Запрос к Mapbox Geocoding API: ${districtName}, ${city}`);

    try {
      const query = encodeURIComponent(`${districtName}, ${city}, Serbia`);
      const url =
        `${CONFIG.MAPBOX_BASE}/geocoding/v5/mapbox.places/${query}.json?` +
        `types=district,neighborhood,locality&` +
        `bbox=true&polygon=true&limit=5&` +
        `access_token=${mapboxToken}`;

      const response = await this.makeRequest(url);

      if (response.features && response.features.length > 0) {
        for (const feature of response.features) {
          if (feature.geometry && feature.geometry.type === 'Polygon') {
            const result = {
              source: 'mapbox',
              name: feature.text || districtName,
              nameEn:
                feature.properties?.name_en || feature.text || districtName,
              coordinates: feature.geometry.coordinates[0], // Берем внешнее кольцо полигона
              bbox: feature.bbox,
              properties: feature.properties,
            };

            this.saveToCache(cacheKey, result);
            return result;
          }
        }
      }
    } catch (error) {
      console.warn(`⚠️ Mapbox API ошибка: ${error.message}`);
    }

    return null;
  }

  // 2. Улучшенный Overpass API запрос
  async fetchFromOverpass(districtName, city, osmId = null) {
    const cacheKey = this.getCacheKey({
      source: 'overpass',
      district: districtName,
      city,
      osmId,
    });
    const cached = this.getFromCache(cacheKey);
    if (cached) return cached;

    console.log(`🌐 Запрос к Overpass API: ${districtName}, ${city}`);

    let query;
    if (osmId) {
      query = `[out:json][timeout:30];
        relation(${osmId});
        out geom;`;
    } else {
      // Расширенный поиск с multiple criteria
      query = `[out:json][timeout:30];
        area[name="${city}"][place=city]->.city;
        (
          // Поиск по точному названию
          relation(area.city)["name"="${districtName}"]["place"~"suburb|neighbourhood|quarter"];
          relation(area.city)["name:en"="${districtName}"]["place"~"suburb|neighbourhood|quarter"];
          relation(area.city)["name:sr"="${districtName}"]["place"~"suburb|neighbourhood|quarter"];
          
          // Поиск по частичному совпадению (case insensitive)
          relation(area.city)["name"~"${districtName}",i]["place"~"suburb|neighbourhood|quarter"];
          relation(area.city)["name:en"~"${districtName}",i]["place"~"suburb|neighbourhood|quarter"];
          
          // Альтернативный поиск как way
          way(area.city)["name"="${districtName}"]["place"~"suburb|neighbourhood|quarter"];
          way(area.city)["name"~"${districtName}",i]["place"~"suburb|neighbourhood|quarter"];
        );
        out geom;`;
    }

    try {
      const postData = `data=${encodeURIComponent(query)}`;
      const response = await this.makeRequest(CONFIG.OVERPASS_BASE, {
        method: 'POST',
        headers: {
          'Content-Type': 'application/x-www-form-urlencoded',
          'Content-Length': Buffer.byteLength(postData),
        },
        data: postData,
      });

      if (response.elements && response.elements.length > 0) {
        // Сортируем по приоритету: relation > way, точное совпадение > частичное
        const element = this.selectBestElement(response.elements, districtName);

        if (element) {
          const coordinates = this.extractCoordinates(element);
          if (coordinates && coordinates.length > 0) {
            const result = {
              source: 'overpass',
              osmId: element.id,
              osmType: element.type,
              name: element.tags?.name || districtName,
              nameEn:
                element.tags?.['name:en'] || element.tags?.name || districtName,
              nameSr:
                element.tags?.['name:sr'] || element.tags?.name || districtName,
              coordinates: coordinates,
              place: element.tags?.place,
              tags: element.tags,
            };

            this.saveToCache(cacheKey, result);
            return result;
          }
        }
      }
    } catch (error) {
      console.warn(`⚠️ Overpass API ошибка: ${error.message}`);
    }

    return null;
  }

  // Выбор лучшего элемента из результатов Overpass
  selectBestElement(elements, targetName) {
    // Приоритет: relation > way > node
    const relations = elements.filter((e) => e.type === 'relation');
    const ways = elements.filter((e) => e.type === 'way');

    const candidates = relations.length > 0 ? relations : ways;
    if (candidates.length === 0) return null;

    // Ищем точное совпадение имени
    const exactMatch = candidates.find(
      (e) =>
        e.tags?.name === targetName ||
        e.tags?.['name:en'] === targetName ||
        e.tags?.['name:sr'] === targetName
    );

    if (exactMatch) return exactMatch;

    // Если точного совпадения нет, берем первый подходящий
    return candidates[0];
  }

  // Извлечение координат из элемента Overpass
  extractCoordinates(element) {
    if (element.type === 'relation' && element.members) {
      return this.assembleRelationCoordinates(element.members);
    } else if (element.type === 'way' && element.geometry) {
      return element.geometry.map((coord) => [coord.lon, coord.lat]);
    }
    return null;
  }

  // Улучшенная сборка координат из relation members
  assembleRelationCoordinates(members) {
    const outerWays = members.filter(
      (member) =>
        member.role === 'outer' && member.geometry && member.geometry.length > 0
    );

    if (outerWays.length === 0) return null;
    if (outerWays.length === 1) {
      return outerWays[0].geometry.map((coord) => [coord.lon, coord.lat]);
    }

    console.log(`🔗 Сборка полигона из ${outerWays.length} way-элементов`);

    // Для complex multipolygon используем алгоритм поиска цепочки
    return this.buildPolygonChain(outerWays);
  }

  // Алгоритм построения цепочки полигона
  buildPolygonChain(ways) {
    if (ways.length === 0) return null;

    const coords = [];
    const used = new Set();

    // Начинаем с первого way
    let current = ways[0];
    used.add(0);
    coords.push(...current.geometry.map((coord) => [coord.lon, coord.lat]));

    let lastPoint = current.geometry[current.geometry.length - 1];

    // Строим цепочку
    while (used.size < ways.length) {
      let nextIndex = -1;
      let reversed = false;

      // Ищем следующий way
      for (let i = 0; i < ways.length; i++) {
        if (used.has(i)) continue;

        const way = ways[i];
        const first = way.geometry[0];
        const last = way.geometry[way.geometry.length - 1];

        if (this.pointsEqual(lastPoint, first)) {
          nextIndex = i;
          reversed = false;
          break;
        } else if (this.pointsEqual(lastPoint, last)) {
          nextIndex = i;
          reversed = true;
          break;
        }
      }

      if (nextIndex === -1) {
        console.log(
          `⚠️ Не удалось найти продолжение цепочки. Собрано ${coords.length} точек`
        );
        break;
      }

      const nextWay = ways[nextIndex];
      used.add(nextIndex);

      if (reversed) {
        const reversedCoords = nextWay.geometry
          .slice(0, -1)
          .reverse()
          .map((coord) => [coord.lon, coord.lat]);
        coords.push(...reversedCoords);
        lastPoint = nextWay.geometry[0];
      } else {
        const forwardCoords = nextWay.geometry
          .slice(1)
          .map((coord) => [coord.lon, coord.lat]);
        coords.push(...forwardCoords);
        lastPoint = nextWay.geometry[nextWay.geometry.length - 1];
      }
    }

    console.log(`✅ Собран полигон из ${coords.length} точек`);
    return coords;
  }

  // Проверка равенства точек с допуском
  pointsEqual(p1, p2, tolerance = 0.000001) {
    return (
      Math.abs(p1.lat - p2.lat) < tolerance &&
      Math.abs(p1.lon - p2.lon) < tolerance
    );
  }

  // 3. Fallback на Nominatim API (улучшенный)
  async fetchFromNominatim(districtName, city) {
    const cacheKey = this.getCacheKey({
      source: 'nominatim',
      district: districtName,
      city,
    });
    const cached = this.getFromCache(cacheKey);
    if (cached) return cached;

    console.log(`📍 Запрос к Nominatim API: ${districtName}, ${city}`);

    try {
      const query = encodeURIComponent(`${districtName}, ${city}, Serbia`);
      const url =
        `${CONFIG.NOMINATIM_BASE}/search?` +
        `q=${query}&format=json&addressdetails=1&limit=10&` +
        `polygon_geojson=1&extratags=1`;

      const response = await this.makeRequest(url);

      if (response && response.length > 0) {
        // Ищем relation с полигоном
        const relationWithPolygon = response.find(
          (r) =>
            r.osm_type === 'relation' &&
            r.geojson &&
            r.geojson.type === 'Polygon'
        );

        if (relationWithPolygon) {
          const result = {
            source: 'nominatim',
            osmId: relationWithPolygon.osm_id,
            osmType: relationWithPolygon.osm_type,
            name: relationWithPolygon.display_name.split(',')[0],
            nameEn: districtName,
            coordinates: relationWithPolygon.geojson.coordinates[0],
            properties: {
              address: relationWithPolygon.address,
              type: relationWithPolygon.type,
              class: relationWithPolygon.class,
            },
          };

          this.saveToCache(cacheKey, result);
          return result;
        }
      }
    } catch (error) {
      console.warn(`⚠️ Nominatim API ошибка: ${error.message}`);
    }

    return null;
  }

  // Валидация полученных данных
  validateDistrict(district) {
    if (!this.options.validate) return { valid: true };

    const issues = [];

    // 1. Проверка количества координат
    if (
      !district.coordinates ||
      district.coordinates.length < CONFIG.MIN_COORDINATE_COUNT
    ) {
      issues.push(
        `Слишком мало координат: ${district.coordinates?.length || 0}`
      );
    }
    if (
      district.coordinates &&
      district.coordinates.length > CONFIG.MAX_COORDINATE_COUNT
    ) {
      issues.push(`Слишком много координат: ${district.coordinates.length}`);
    }

    // 2. Проверка площади
    if (district.coordinates) {
      const area = this.calculatePolygonArea(district.coordinates);
      if (area < CONFIG.MIN_AREA_SIZE) {
        issues.push(`Слишком маленькая площадь: ${area.toFixed(6)}`);
      }
      if (area > CONFIG.MAX_AREA_SIZE) {
        issues.push(`Слишком большая площадь: ${area.toFixed(6)}`);
      }
    }

    // 3. Проверка замкнутости полигона
    if (district.coordinates && district.coordinates.length > 0) {
      const first = district.coordinates[0];
      const last = district.coordinates[district.coordinates.length - 1];
      if (
        !this.pointsEqual(
          { lat: first[1], lon: first[0] },
          { lat: last[1], lon: last[0] }
        )
      ) {
        issues.push('Полигон не замкнут');
      }
    }

    return {
      valid: issues.length === 0,
      issues: issues,
      area: district.coordinates
        ? this.calculatePolygonArea(district.coordinates)
        : 0,
      coordinateCount: district.coordinates?.length || 0,
    };
  }

  // Вычисление площади полигона (приблизительно)
  calculatePolygonArea(coordinates) {
    if (coordinates.length < 3) return 0;

    let area = 0;
    const n = coordinates.length;

    for (let i = 0; i < n; i++) {
      const j = (i + 1) % n;
      area += coordinates[i][0] * coordinates[j][1];
      area -= coordinates[j][0] * coordinates[i][1];
    }

    return Math.abs(area) / 2;
  }

  // Основной метод для получения района
  async fetchDistrict(districtName, city, mapboxToken = null) {
    console.log(`\n🔍 Поиск района "${districtName}" в ${city}`);
    console.log('=' + '='.repeat(60));

    const strategies = [
      {
        name: 'Mapbox Geocoding',
        method: () => this.fetchFromMapbox(districtName, city, mapboxToken),
      },
      {
        name: 'Overpass API',
        method: () => this.fetchFromOverpass(districtName, city),
      },
      {
        name: 'Nominatim API',
        method: () => this.fetchFromNominatim(districtName, city),
      },
    ];

    for (const strategy of strategies) {
      try {
        console.log(`\n🎯 Пробуем: ${strategy.name}`);
        const result = await strategy.method();

        if (result) {
          const validation = this.validateDistrict(result);

          console.log(`✅ Найдено через ${strategy.name}:`);
          console.log(`   Название: ${result.name} (${result.nameEn})`);
          console.log(`   Источник: ${result.source}`);
          console.log(`   Координат: ${result.coordinates?.length || 0}`);

          if (this.options.validate) {
            if (validation.valid) {
              console.log(
                `   ✅ Валидация пройдена (площадь: ${validation.area.toFixed(6)})`
              );
            } else {
              console.log(
                `   ⚠️ Проблемы валидации: ${validation.issues.join(', ')}`
              );
              // Продолжаем поиск если данные не прошли валидацию
              continue;
            }
          }

          return result;
        } else {
          console.log(`   ❌ Район не найден`);
        }
      } catch (error) {
        console.log(`   ❌ Ошибка: ${error.message}`);
      }

      // Пауза между запросами
      if (strategies.indexOf(strategy) < strategies.length - 1) {
        await new Promise((resolve) =>
          setTimeout(resolve, CONFIG.RATE_LIMIT_DELAY)
        );
      }
    }

    return null;
  }

  // Добавление района в файл
  async addDistrictToFile(district) {
    if (!district) {
      throw new Error('Нет данных района для добавления');
    }

    console.log('\n📝 Добавление района в файл...');

    // Создаем бэкап
    const backupPath = CONFIG.DISTRICTS_FILE + `.backup.${Date.now()}`;
    fs.copyFileSync(CONFIG.DISTRICTS_FILE, backupPath);
    console.log(`📁 Создан бэкап: ${path.basename(backupPath)}`);

    // Читаем файл
    let content = fs.readFileSync(CONFIG.DISTRICTS_FILE, 'utf8');

    // Генерируем уникальный ID
    const baseId = (district.nameEn || district.name)
      .toLowerCase()
      .replace(/\s+/g, '-')
      .replace(/[^\w\-]/g, '')
      .replace(/^-+|-+$/g, '');

    let districtId = baseId;
    let counter = 1;
    while (content.includes(`id: '${districtId}'`)) {
      districtId = `${baseId}-${counter++}`;
    }

    // Выбираем цвет
    const existingCount = (content.match(/id:\s*'/g) || []).length;
    const color = CONFIG.COLORS[existingCount % CONFIG.COLORS.length];

    // Оптимизируем координаты если их слишком много
    let coordinates = district.coordinates;
    if (coordinates.length > 1000) {
      console.log(
        `📐 Упрощение полигона: ${coordinates.length} → `,
        (end = '')
      );
      coordinates = this.simplifyPolygon(coordinates, 0.0001);
      console.log(`${coordinates.length} точек`);
    }

    // Генерируем код объекта
    const districtCode = this.generateDistrictCode(
      districtId,
      district,
      coordinates,
      color
    );

    // Вставляем в файл
    const insertPosition = content.lastIndexOf('  ];');
    if (insertPosition === -1) {
      throw new Error('Не найден массив districts в файле');
    }

    // Проверяем нужна ли запятая
    const beforeInsert = content.substring(0, insertPosition).trimEnd();
    if (beforeInsert.endsWith('}') && !beforeInsert.endsWith('},')) {
      const lastBraceIndex = beforeInsert.lastIndexOf('}');
      content =
        beforeInsert.substring(0, lastBraceIndex + 1) +
        ',' +
        beforeInsert.substring(lastBraceIndex + 1) +
        districtCode +
        '\n' +
        content.substring(insertPosition);
    } else {
      content =
        content.substring(0, insertPosition) +
        districtCode +
        '\n' +
        content.substring(insertPosition);
    }

    // Записываем файл
    fs.writeFileSync(CONFIG.DISTRICTS_FILE, content);

    console.log(`✅ Район "${district.name}" добавлен!`);
    console.log(`   ID: ${districtId}`);
    console.log(`   Цвет: ${color}`);
    console.log(`   Координат: ${coordinates.length}`);
    console.log(`   Источник: ${district.source}`);

    return districtId;
  }

  // Упрощение полигона (Douglas-Peucker algorithm)
  simplifyPolygon(coordinates, tolerance) {
    if (coordinates.length <= 2) return coordinates;

    // Простая реализация упрощения - берем каждую N-ю точку
    const step = Math.max(1, Math.floor(coordinates.length / 500));
    const simplified = [];

    for (let i = 0; i < coordinates.length; i += step) {
      simplified.push(coordinates[i]);
    }

    // Убеждаемся что первая и последняя точка одинаковые
    if (simplified.length > 0) {
      const first = simplified[0];
      const last = simplified[simplified.length - 1];
      if (first[0] !== last[0] || first[1] !== last[1]) {
        simplified.push(first);
      }
    }

    return simplified;
  }

  // Генерация кода района для вставки в файл
  generateDistrictCode(id, district, coordinates, color) {
    const coordsStr = coordinates
      .slice(0, -1) // Убираем дублирующую последнюю точку
      .map(
        (coord) => `        [${coord[0].toFixed(6)}, ${coord[1].toFixed(6)}]`
      )
      .join(',\n');

    return `    {
      id: '${id}',
      name: '${district.name}',
      nameEn: '${district.nameEn}',
      population: 'Неизвестно',
      area: 'Неизвестно',
      color: '${color}',
      description: 'Получено из ${district.source}${district.osmId ? ` (OSM ${district.osmType} ${district.osmId})` : ''} - ${coordinates.length} точек',
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

  // Получение всех районов города
  async fetchAllDistricts(cityName, mapboxToken = null) {
    console.log(`\n🏙️ Получение всех районов города ${cityName}`);

    // Список основных районов для тестирования
    const districts = [
      'Veternik',
      'Futog',
      'Petrovaradin',
      'Sremska Kamenica',
      'Liman',
      'Detelinara',
      'Grbavica',
      'Podbara',
      'Salajka',
    ];

    const results = [];

    for (const districtName of districts) {
      try {
        const district = await this.fetchDistrict(
          districtName,
          cityName,
          mapboxToken
        );
        if (district) {
          results.push(district);
          console.log(`\n✅ ${districtName}: успешно получен`);
        } else {
          console.log(`\n❌ ${districtName}: не найден`);
        }

        // Пауза между запросами
        await new Promise((resolve) =>
          setTimeout(resolve, CONFIG.RATE_LIMIT_DELAY * 2)
        );
      } catch (error) {
        console.error(`\n❌ ${districtName}: ошибка - ${error.message}`);
      }
    }

    console.log(
      `\n📊 Итого найдено: ${results.length}/${districts.length} районов`
    );
    return results;
  }
}

// CLI интерфейс
async function main() {
  const args = process.argv.slice(2);

  if (args.length === 0) {
    console.log(`
🏘️ Улучшенный инструмент получения границ районов

Использование:
  node enhanced_district_fetcher.js "Район" "Город" [опции]
  node enhanced_district_fetcher.js --city "Город" --fetch-all [опции]

Опции:
  --mapbox-token TOKEN     Токен для Mapbox API (лучшее качество)
  --no-cache              Отключить кэширование
  --no-validate           Отключить валидацию данных
  --dry-run               Только получить данные, не добавлять в файл

Примеры:
  node enhanced_district_fetcher.js "Veternik" "Novi Sad"
  node enhanced_district_fetcher.js "Лиман" "Нови Сад" --mapbox-token pk.123...
  node enhanced_district_fetcher.js --city "Novi Sad" --fetch-all
`);
    process.exit(1);
  }

  // Парсинг аргументов
  const options = {
    cache: !args.includes('--no-cache'),
    validate: !args.includes('--no-validate'),
    dryRun: args.includes('--dry-run'),
    fetchAll: args.includes('--fetch-all'),
  };

  const mapboxTokenIndex = args.indexOf('--mapbox-token');
  const mapboxToken =
    mapboxTokenIndex !== -1
      ? args[mapboxTokenIndex + 1]
      : process.env.MAPBOX_TOKEN;

  const cityIndex = args.indexOf('--city');
  const city = cityIndex !== -1 ? args[cityIndex + 1] : args[1];

  const fetcher = new DistrictFetcher(options);

  try {
    if (options.fetchAll) {
      // Получение всех районов города
      const districts = await fetcher.fetchAllDistricts(city, mapboxToken);

      if (!options.dryRun) {
        for (const district of districts) {
          await fetcher.addDistrictToFile(district);
        }
      }

      console.log(`\n🎉 Обработано ${districts.length} районов`);
    } else {
      // Получение одного района
      const districtName = args[0];

      if (!districtName || !city) {
        throw new Error('Необходимо указать название района и город');
      }

      const district = await fetcher.fetchDistrict(
        districtName,
        city,
        mapboxToken
      );

      if (district) {
        if (options.dryRun) {
          console.log('\n📋 Данные района (dry-run режим):');
          console.log(JSON.stringify(district, null, 2));
        } else {
          await fetcher.addDistrictToFile(district);
        }

        console.log('\n🎉 Готово!');
      } else {
        console.error('\n❌ Район не найден ни в одном источнике');
        process.exit(1);
      }
    }
  } catch (error) {
    console.error(`\n❌ Критическая ошибка: ${error.message}`);
    process.exit(1);
  }
}

// Запуск при прямом вызове
if (require.main === module) {
  main().catch((error) => {
    console.error('Неожиданная ошибка:', error);
    process.exit(1);
  });
}

module.exports = DistrictFetcher;
