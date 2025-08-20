'use client';

import React, { useState, useEffect } from 'react';
import Link from 'next/link';

const DistrictManager = () => {
  const [districts, setDistricts] = useState<any[]>([]);
  const [selectedDistricts, setSelectedDistricts] = useState<Set<string>>(
    new Set()
  );
  const [newDistrictName, setNewDistrictName] = useState('');
  const [loading, setLoading] = useState(false);
  const [message, setMessage] = useState<{
    type: 'success' | 'error' | 'info';
    text: string;
  } | null>(null);
  const [searchResults, setSearchResults] = useState<any[]>([]);
  const [searching, setSearching] = useState(false);
  const [dataSource, setDataSource] = useState<'nominatim' | 'overpass'>(
    'nominatim'
  );
  const [selectedCity, setSelectedCity] = useState('');
  const [cityDistricts, setCityDistricts] = useState<any[]>([]);
  const [loadingCityDistricts, setLoadingCityDistricts] = useState(false);
  const [searchArea, setSearchArea] = useState<string>(''); // Область поиска для Overpass

  // Загружаем список текущих районов
  useEffect(() => {
    loadDistricts();
  }, []);

  const loadDistricts = async () => {
    try {
      // Пробуем загрузить через API
      const response = await fetch('/api/districts');
      if (response.ok) {
        const data = await response.json();
        if (data.districts && data.districts.length > 0) {
          setDistricts(data.districts);
          return;
        }
      }
    } catch {
      console.log('Using fallback district data');
    }

    // Fallback: используем статические данные если API недоступен
    const currentDistricts = [
      {
        id: 'detelinara',
        name: 'Детелинара',
        nameEn: 'Detelinara',
        points: 22,
      },
      { id: 'adice', name: 'Адице', nameEn: 'Adice', points: 34 },
      { id: 'grbavica', name: 'Грбавица', nameEn: 'Grbavica', points: 15 },
      {
        id: 'liman-i-iv',
        name: 'Лиман I-IV',
        nameEn: 'Liman I-IV',
        points: 26,
      },
      { id: 'telep', name: 'Телеп', nameEn: 'Telep', points: 23 },
      { id: 'bistrica', name: 'Бистрица', nameEn: 'Bistrica', points: 17 },
      {
        id: 'stari-ledinci',
        name: 'Stari Ledinci',
        nameEn: 'Stari Ledinci',
        points: 631,
      },
      { id: 'ledinci', name: 'Ledinci', nameEn: 'Ledinci', points: 401 },
      { id: 'bukovac', name: 'Bukovac', nameEn: 'Bukovac', points: 739 },
      { id: 'kisač', name: 'Kisač', nameEn: 'Kisač', points: 712 },
      { id: 'budisava', name: 'Budisava', nameEn: 'Budisava', points: 65 },
      { id: 'kać', name: 'Kać', nameEn: 'Kać', points: 354 },
      { id: 'kovilj', name: 'Kovilj', nameEn: 'Kovilj', points: 609 },
      {
        id: 'stari-grad',
        name: 'Stari Grad',
        nameEn: 'Stari Grad',
        points: 38,
      },
      {
        id: 'jugovićevo',
        name: 'Jugovićevo',
        nameEn: 'Jugovićevo',
        points: 14,
      },
      { id: 'salajka', name: 'Salajka', nameEn: 'Salajka', points: 20 },
      { id: 'slana-bara', name: 'Слана бара', nameEn: 'Slana bara', points: 8 },
      {
        id: 'novo-naselje',
        name: 'Novo Naselje',
        nameEn: 'Novo Naselje',
        points: 50,
      },
      {
        id: 'rotkvarija',
        name: 'Rotkvarija',
        nameEn: 'Rotkvarija',
        points: 20,
      },
      { id: 'banatić', name: 'Banatić', nameEn: 'Banatić', points: 25 },
      {
        id: 'sremska-kamenica',
        name: 'Sremska Kamenica',
        nameEn: 'Sremska Kamenica',
        points: 1046,
      },
      { id: 'podbara', name: 'Podbara', nameEn: 'Podbara', points: 38 },
      { id: 'klisa', name: 'Klisa', nameEn: 'Klisa', points: 27 },
      {
        id: 'vidovdansko-naselje',
        name: 'Vidovdansko naselje',
        nameEn: 'Vidovdansko naselje',
        points: 14,
      },
      {
        id: 'avijatičarsko-naselje',
        name: 'Avijatičarsko naselje',
        nameEn: 'Avijatičarsko naselje',
        points: 20,
      },
      {
        id: 'radna-zona-sever',
        name: 'Radna zona Sever',
        nameEn: 'Radna zona Sever',
        points: 99,
      },
      {
        id: 'radna-zona-jug',
        name: 'Radna zona Jug',
        nameEn: 'Radna zona Jug',
        points: 99,
      },
      {
        id: 'industrijska-zona',
        name: 'Industrijska zona',
        nameEn: 'Industrijska zona',
        points: 99,
      },
    ];

    setDistricts(currentDistricts);
  };

  const handleToggleDistrict = (districtId: string) => {
    const newSelected = new Set(selectedDistricts);
    if (newSelected.has(districtId)) {
      newSelected.delete(districtId);
    } else {
      newSelected.add(districtId);
    }
    setSelectedDistricts(newSelected);
  };

  const handleDeleteSelected = async () => {
    if (selectedDistricts.size === 0) {
      setMessage({ type: 'error', text: 'Выберите районы для удаления' });
      return;
    }

    if (
      !confirm(
        `Вы уверены, что хотите удалить ${selectedDistricts.size} район(ов)?`
      )
    ) {
      return;
    }

    setLoading(true);
    try {
      const districtsToDelete = Array.from(selectedDistricts);

      // Реальный вызов API для удаления
      const response = await fetch('/api/districts', {
        method: 'DELETE',
        headers: {
          'Content-Type': 'application/json',
        },
        body: JSON.stringify({ districtIds: districtsToDelete }),
      });

      const result = await response.json();

      if (result.success) {
        // Обновляем список районов
        setDistricts((prev) =>
          prev.filter((d) => !selectedDistricts.has(d.id))
        );
        setSelectedDistricts(new Set());
        setMessage({
          type: 'success',
          text: `✅ ${result.message}. Осталось районов: ${result.remaining}`,
        });

        // Обновляем список через секунду
        setTimeout(() => loadDistricts(), 1000);
      } else {
        throw new Error(result.error || 'Неизвестная ошибка');
      }
    } catch (error) {
      setMessage({
        type: 'error',
        text: `Ошибка при удалении: ${error instanceof Error ? error.message : String(error)}`,
      });
      console.error('Delete error:', error);
    } finally {
      setLoading(false);
    }
  };

  // Координаты городов Сербии (bounding box) - увеличенный радиус
  const cityBounds: { [key: string]: [number, number, number, number] } = {
    'Novi Sad': [45.2, 19.7, 45.35, 19.95], // Увеличен радиус
    Belgrade: [44.65, 20.2, 45.0, 20.7], // Максимально увеличен радиус для Белграда
    Niš: [43.28, 21.85, 43.38, 21.97],
    Kragujevac: [43.97, 20.85, 44.08, 20.97],
    Subotica: [46.05, 19.6, 46.16, 19.72],
    Zrenjanin: [45.33, 20.32, 45.44, 20.44],
    Pančevo: [44.82, 20.59, 44.93, 20.71],
  };

  const loadAllCityDistricts = async (cityName: string) => {
    setLoadingCityDistricts(true);
    setCityDistricts([]);
    setMessage({
      type: 'info',
      text: `🔍 Загрузка всех районов города ${cityName}...`,
    });

    try {
      let query = '';

      if (dataSource === 'overpass') {
        // Получаем координаты города
        const bounds = cityBounds[cityName];
        if (!bounds) {
          setMessage({
            type: 'error',
            text: `Координаты города ${cityName} не найдены`,
          });
          setLoadingCityDistricts(false);
          return;
        }

        // Overpass API запрос для всех районов города по координатам
        // Используем relation и way, включая boundary/administrative, исключаем транспорт
        query = `[out:json][timeout:25];
          (
            relation(${bounds[0]},${bounds[1]},${bounds[2]},${bounds[3]})["place"~"suburb|quarter|neighbourhood|district|borough|village|hamlet"][!"route"][!"public_transport"];
            way(${bounds[0]},${bounds[1]},${bounds[2]},${bounds[3]})["place"~"suburb|quarter|neighbourhood|district|borough|village|hamlet"][!"route"][!"public_transport"];
            relation(${bounds[0]},${bounds[1]},${bounds[2]},${bounds[3]})["boundary"="administrative"]["admin_level"~"9|10"][!"route"];
            way(${bounds[0]},${bounds[1]},${bounds[2]},${bounds[3]})["boundary"="administrative"]["admin_level"~"9|10"][!"route"];
          );
          out geom;`;

        const response = await fetch(
          'https://overpass-api.de/api/interpreter',
          {
            method: 'POST',
            body: `data=${encodeURIComponent(query)}`,
            headers: { 'Content-Type': 'application/x-www-form-urlencoded' },
          }
        );

        if (response.ok) {
          const data = await response.json();
          const districts = (data.elements || [])
            .filter((elem: any) => {
              // Исключаем автобусные маршруты и транспортные объекты
              if (elem.tags) {
                // Исключаем если в названии есть "бус", "bus", "маршрут", "route"
                const name = (elem.tags.name || '').toLowerCase();
                if (
                  name.includes('бус') ||
                  name.includes('bus') ||
                  name.includes('маршрут') ||
                  name.includes('route') ||
                  name.includes('=>') ||
                  name.includes('->')
                ) {
                  return false;
                }

                // Исключаем объекты с тегами транспорта
                if (
                  elem.tags.route ||
                  elem.tags.public_transport ||
                  elem.tags.type === 'route' ||
                  elem.tags.type === 'route_master'
                ) {
                  return false;
                }
              }

              // Принимаем только way и relation с геометрией, исключаем node
              if (elem.type === 'way') {
                // Для way проверяем наличие geometry
                return elem.geometry && elem.geometry.length > 0;
              } else if (elem.type === 'relation') {
                // Для relation проверяем наличие members с геометрией
                return (
                  elem.members &&
                  elem.members.some(
                    (m: any) => m.geometry && m.geometry.length > 0
                  )
                );
              }
              return false; // исключаем node
            })
            .map((elem: any) => ({
              osm_id: elem.id,
              osm_type: elem.type,
              type: elem.tags?.place || 'district',
              name: elem.tags?.name || elem.tags?.['name:en'] || 'Unknown',
              nameEn: elem.tags?.['name:en'] || elem.tags?.name,
              nameSr: elem.tags?.['name:sr'] || elem.tags?.name,
              display_name: `${elem.tags?.name || 'Unknown'}, ${cityName}`,
              source: 'overpass',
              tags: elem.tags,
              // Сохраняем геометрию для дальнейшего использования
              geometry: elem.geometry || (elem.members ? elem.members : null),
              lat: elem.lat,
              lon: elem.lon,
            }));

          // Сортируем по имени
          districts.sort((a: any, b: any) => a.name.localeCompare(b.name));

          setCityDistricts(districts);
          setMessage({
            type: 'success',
            text: `✅ Найдено ${districts.length} районов в городе ${cityName} через Overpass API`,
          });
        }
      } else {
        // Nominatim API - используем поиск районов в городе
        // Поиск всех районов города через несколько запросов
        const searchQueries = [
          `suburb in ${cityName}, Serbia`,
          `neighbourhood in ${cityName}, Serbia`,
          `quarter in ${cityName}, Serbia`,
          `district in ${cityName}, Serbia`,
          `administrative in ${cityName}, Serbia`,
          `village in ${cityName}, Serbia`,
          `town in ${cityName}, Serbia`,
          `hamlet in ${cityName}, Serbia`,
          `${cityName} suburb`,
          `${cityName} neighbourhood`,
          `${cityName} district`,
          `${cityName} village`,
          `${cityName}, Serbia suburb`,
          `${cityName}, Serbia neighbourhood`,
          `${cityName}, Serbia district`,
          `${cityName}, Serbia village`,
        ];

        const allResults = await Promise.all(
          searchQueries.map(async (query) => {
            try {
              const response = await fetch(
                `https://nominatim.openstreetmap.org/search?` +
                  `q=${encodeURIComponent(query)}&` +
                  `format=json&addressdetails=1&limit=50&` +
                  `extratags=1&namedetails=1&polygon_geojson=0`
              );
              if (response.ok) {
                const data = await response.json();
                return data;
              }
              return [];
            } catch (error) {
              console.error('Nominatim query error:', error);
              return [];
            }
          })
        );

        // Объединяем и фильтруем результаты
        const combinedResults = allResults.flat();
        const uniqueDistricts = new Map();

        combinedResults.forEach((item: any) => {
          // Принимаем только relation и way с полигонами, исключаем node
          if (
            (item.type === 'suburb' ||
              item.type === 'neighbourhood' ||
              item.type === 'quarter' ||
              item.type === 'district' ||
              item.type === 'administrative' ||
              item.type === 'village' ||
              item.type === 'hamlet' ||
              item.type === 'borough' ||
              item.type === 'city_district') &&
            (item.osm_type === 'relation' || item.osm_type === 'way') && // убираем node
            !uniqueDistricts.has(item.osm_id) &&
            item.display_name &&
            (item.display_name.toLowerCase().includes(cityName.toLowerCase()) ||
              item.display_name.toLowerCase().includes('novi sad') ||
              item.display_name.toLowerCase().includes('нови сад'))
          ) {
            const name = item.display_name.split(',')[0].trim();
            uniqueDistricts.set(item.osm_id, {
              osm_id: item.osm_id,
              osm_type: item.osm_type,
              type: item.type,
              name: name,
              nameEn: item.namedetails?.['name:en'] || name,
              nameSr: item.namedetails?.['name:sr'] || name,
              display_name: item.display_name,
              source: 'nominatim',
              tags: item.extratags,
            });
          }
        });

        const districts = Array.from(uniqueDistricts.values());
        districts.sort((a: any, b: any) => a.name.localeCompare(b.name));

        setCityDistricts(districts);
        setMessage({
          type: 'success',
          text: `✅ Найдено ${districts.length} районов в городе ${cityName} через Nominatim API`,
        });
      }
    } catch (error) {
      setMessage({
        type: 'error',
        text: `Ошибка при загрузке районов: ${error instanceof Error ? error.message : String(error)}`,
      });
    } finally {
      setLoadingCityDistricts(false);
    }
  };

  const searchDistrict = async () => {
    if (!newDistrictName.trim()) {
      setMessage({ type: 'error', text: 'Введите название района' });
      return;
    }

    setSearching(true);
    setSearchResults([]);

    try {
      const searchName = newDistrictName.trim().toLowerCase();

      // Если выбран Overpass API
      if (dataSource === 'overpass') {
        setMessage({ type: 'info', text: `🔍 Поиск через Overpass API...` });

        // Построение запроса для Overpass
        let overpassQuery = '';

        if (searchName === 'liman' || searchName === 'лиман') {
          // Специальный запрос для районов Liman
          overpassQuery = `[out:json][timeout:25];
            (
              relation["name"~"Лиман [1-4]?"]["place"~"quarter|suburb"];
              relation["name:en"~"Liman [1-4]?"]["place"~"quarter|suburb"];
            );
            out body;`;
        } else {
          // Общий запрос для других районов с выбором области
          let searchQuery = '';

          if (searchArea === 'Serbia') {
            // Поиск по всей Сербии с фильтрацией по типам
            searchQuery = `
              area[name="Serbia"]->.searchArea;
              (
                relation(area.searchArea)["name"~"${newDistrictName}",i]["type"~"multipolygon|boundary"]["place"~"suburb|neighbourhood|quarter|district|village|hamlet|borough"];
                way(area.searchArea)["name"~"${newDistrictName}",i]["place"~"suburb|neighbourhood|quarter|district|village|hamlet|borough"];
                relation(area.searchArea)["name:en"~"${newDistrictName}",i]["type"~"multipolygon|boundary"]["place"~"suburb|neighbourhood|quarter|district|village|hamlet|borough"];
                way(area.searchArea)["name:en"~"${newDistrictName}",i]["place"~"suburb|neighbourhood|quarter|district|village|hamlet|borough"];
                relation(area.searchArea)["name:sr"~"${newDistrictName}",i]["type"~"multipolygon|boundary"]["place"~"suburb|neighbourhood|quarter|district|village|hamlet|borough"];
                way(area.searchArea)["name:sr"~"${newDistrictName}",i]["place"~"suburb|neighbourhood|quarter|district|village|hamlet|borough"];
              );`;
          } else if (searchArea === 'Novi Sad') {
            // Поиск в границах Нови-Сада с фильтрацией по типам
            searchQuery = `
              area[name="Novi Sad"]->.novisad;
              (
                relation(area.novisad)["name"~"${newDistrictName}",i]["type"~"multipolygon|boundary"]["place"~"suburb|neighbourhood|quarter|district|village|hamlet|borough"];
                way(area.novisad)["name"~"${newDistrictName}",i]["place"~"suburb|neighbourhood|quarter|district|village|hamlet|borough"];
                relation(area.novisad)["name:en"~"${newDistrictName}",i]["type"~"multipolygon|boundary"]["place"~"suburb|neighbourhood|quarter|district|village|hamlet|borough"];
                way(area.novisad)["name:en"~"${newDistrictName}",i]["place"~"suburb|neighbourhood|quarter|district|village|hamlet|borough"];
                relation(area.novisad)["name:sr"~"${newDistrictName}",i]["type"~"multipolygon|boundary"]["place"~"suburb|neighbourhood|quarter|district|village|hamlet|borough"];
                way(area.novisad)["name:sr"~"${newDistrictName}",i]["place"~"suburb|neighbourhood|quarter|district|village|hamlet|borough"];
              );`;
          } else {
            // Широкий поиск по координатам с фильтрацией (включает пригороды)
            searchQuery = `
              (
                relation(45.15,19.65,45.4,20.0)["name"~"${newDistrictName}",i]["type"~"multipolygon|boundary"]["place"~"suburb|neighbourhood|quarter|district|village|hamlet|borough"];
                way(45.15,19.65,45.4,20.0)["name"~"${newDistrictName}",i]["place"~"suburb|neighbourhood|quarter|district|village|hamlet|borough"];
                relation(45.15,19.65,45.4,20.0)["name:en"~"${newDistrictName}",i]["type"~"multipolygon|boundary"]["place"~"suburb|neighbourhood|quarter|district|village|hamlet|borough"];
                way(45.15,19.65,45.4,20.0)["name:en"~"${newDistrictName}",i]["place"~"suburb|neighbourhood|quarter|district|village|hamlet|borough"];
                relation(45.15,19.65,45.4,20.0)["name:sr"~"${newDistrictName}",i]["type"~"multipolygon|boundary"]["place"~"suburb|neighbourhood|quarter|district|village|hamlet|borough"];
                way(45.15,19.65,45.4,20.0)["name:sr"~"${newDistrictName}",i]["place"~"suburb|neighbourhood|quarter|district|village|hamlet|borough"];
              );`;
          }

          overpassQuery = `[out:json][timeout:25];${searchQuery}
            out body;`;
        }

        const overpassResponse = await fetch(
          'https://overpass-api.de/api/interpreter',
          {
            method: 'POST',
            body: `data=${encodeURIComponent(overpassQuery)}`,
            headers: { 'Content-Type': 'application/x-www-form-urlencoded' },
          }
        );

        if (overpassResponse.ok) {
          const overpassData = await overpassResponse.json();
          const elements = overpassData.elements || [];

          // Фильтруем и преобразуем в формат, совместимый с системой
          const overpassResults = elements
            .filter((elem: any) => {
              // Исключаем автобусные маршруты и транспортные объекты
              if (elem.tags) {
                // Исключаем если в названии есть "бус", "bus", "маршрут", "route"
                const name = (elem.tags.name || '').toLowerCase();
                if (
                  name.includes('бус') ||
                  name.includes('bus') ||
                  name.includes('маршрут') ||
                  name.includes('route') ||
                  name.includes('=>') ||
                  name.includes('->')
                ) {
                  return false;
                }

                // Исключаем объекты с тегами транспорта
                if (
                  elem.tags.route ||
                  elem.tags.public_transport ||
                  elem.tags.type === 'route' ||
                  elem.tags.type === 'route_master'
                ) {
                  return false;
                }
              }

              // Принимаем только объекты с геометрией
              if (
                elem.type === 'way' &&
                (!elem.geometry || elem.geometry.length === 0)
              ) {
                return false;
              }
              if (
                elem.type === 'relation' &&
                (!elem.members ||
                  !elem.members.some(
                    (m: any) => m.geometry && m.geometry.length > 0
                  ))
              ) {
                return false;
              }

              return true;
            })
            .map((elem: any) => ({
              osm_id: elem.id,
              osm_type: elem.type,
              type: elem.tags?.place || 'quarter',
              display_name: `${elem.tags?.name || elem.tags?.['name:en'] || 'Unknown'}, Novi Sad, Serbia`,
              class: 'place',
              tags: elem.tags,
              source: 'overpass',
            }));

          setSearchResults(overpassResults);

          if (overpassResults.length === 0) {
            setMessage({
              type: 'info',
              text: 'Районы не найдены через Overpass API. Попробуйте другое название или переключитесь на Nominatim.',
            });
          } else {
            setMessage({
              type: 'success',
              text: `✅ Найдено ${overpassResults.length} результатов через Overpass API!`,
            });
          }
        } else {
          throw new Error('Overpass API request failed');
        }
      } else {
        // Nominatim API (существующий код)
        setMessage({ type: 'info', text: `🔍 Поиск через Nominatim API...` });

        const searchQueries = [];

        if (searchName === 'liman' || searchName === 'лиман') {
          searchQueries.push(
            `Liman I, Novi Sad, Serbia`,
            `Liman II, Novi Sad, Serbia`,
            `Liman III, Novi Sad, Serbia`,
            `Liman IV, Novi Sad, Serbia`,
            `Liman, Novi Sad, Serbia`
          );
        } else if (searchName === 'rotkvarija' || searchName === 'ротквария') {
          searchQueries.push(
            `Rotkvarija, Novi Sad, Serbia`,
            `Rotkvaria, Novi Sad, Serbia`
          );
        } else if (searchName === 'detelinara' || searchName === 'детелинара') {
          searchQueries.push(
            `Detelinara, Novi Sad, Serbia`,
            `Detelinara I, Novi Sad, Serbia`,
            `Detelinara II, Novi Sad, Serbia`
          );
        } else {
          searchQueries.push(
            `${newDistrictName}, Novi Sad, Serbia`,
            `${newDistrictName}, Novi Sad`,
            `${newDistrictName} Novi Sad Serbia`
          );
        }

        // Выполняем все поиски параллельно
        const allResults = await Promise.all(
          searchQueries.map(async (query) => {
            const encodedQuery = encodeURIComponent(query);
            const response = await fetch(
              `https://nominatim.openstreetmap.org/search?q=${encodedQuery}&format=json&addressdetails=1&limit=5`
            );
            return response.json();
          })
        );

        // Объединяем и фильтруем результаты
        const combinedResults = allResults.flat();
        const filteredResults = combinedResults.filter(
          (r: any) =>
            r &&
            r.osm_type === 'relation' &&
            (r.type === 'administrative' ||
              r.type === 'suburb' ||
              r.type === 'neighbourhood' ||
              r.type === 'quarter')
        );

        // Убираем дубликаты по osm_id
        const uniqueResults = filteredResults
          .filter(
            (item: any, index: number, self: any[]) =>
              index === self.findIndex((t: any) => t.osm_id === item.osm_id)
          )
          .map((r: any) => ({ ...r, source: 'nominatim' }));

        setSearchResults(uniqueResults);

        if (uniqueResults.length === 0) {
          setMessage({
            type: 'info',
            text: 'Районы не найдены через Nominatim API. Попробуйте переключиться на Overpass API.',
          });
        } else {
          setMessage({
            type: 'success',
            text: `✅ Найдено ${uniqueResults.length} результатов через Nominatim API!`,
          });
        }
      }
    } catch (error) {
      setMessage({
        type: 'error',
        text: `Ошибка при поиске района: ${error instanceof Error ? error.message : String(error)}`,
      });
    } finally {
      setSearching(false);
    }
  };

  const handleAddDistrict = async (osmResult: any) => {
    setLoading(true);
    try {
      const name = osmResult.display_name.split(',')[0];

      // Реальный вызов API для добавления района
      const response = await fetch('/api/districts', {
        method: 'POST',
        headers: {
          'Content-Type': 'application/json',
        },
        body: JSON.stringify({
          name: name,
          city: 'Novi Sad',
          country: 'Serbia',
          source: osmResult.source || 'nominatim',
          osmId: osmResult.osm_id,
          osmType: osmResult.osm_type,
        }),
      });

      const result = await response.json();

      if (result.success) {
        setMessage({
          type: 'success',
          text: `✅ ${result.message}. Перезагрузите страницу, чтобы увидеть изменения.`,
        });

        // Очищаем форму
        setNewDistrictName('');
        setSearchResults([]);

        // Обновляем список через 2 секунды
        setTimeout(() => {
          loadDistricts();
          // Также можно автоматически перезагрузить страницу карты
          // window.location.href = '/ru/examples/novi-sad-districts';
        }, 2000);
      } else {
        throw new Error(result.error || 'Не удалось добавить район');
      }
    } catch (error) {
      setMessage({
        type: 'error',
        text: `Ошибка при добавлении: ${error instanceof Error ? error.message : String(error)}`,
      });
      console.error('Add error:', error);
    } finally {
      setLoading(false);
    }
  };

  const suggestedDistricts = [
    'suburb',
    'neighbourhood',
    'quarter',
    'district',
    'administrative',
    'village',
    'hamlet',
    'borough',
  ];

  return (
    <div className="min-h-screen bg-gradient-to-br from-base-100 to-base-200 p-8">
      <div className="max-w-7xl mx-auto">
        {/* Header */}
        <div className="flex justify-between items-center mb-8">
          <div>
            <h1 className="text-3xl font-bold mb-2">
              🛠️ Управление районами Нови-Сада
            </h1>
            <p className="text-base-content/70">
              Добавляйте и удаляйте районы на интерактивной карте
            </p>
          </div>
          <Link
            href="/ru/examples/novi-sad-districts"
            className="btn btn-primary"
          >
            ← Вернуться к карте
          </Link>
        </div>

        {/* Messages */}
        {message && (
          <div className={`alert alert-${message.type} mb-6`}>
            <span>{message.text}</span>
            <button
              onClick={() => setMessage(null)}
              className="btn btn-sm btn-ghost"
            >
              ✕
            </button>
          </div>
        )}

        <div className="grid grid-cols-1 lg:grid-cols-2 gap-8">
          {/* Add District Section */}
          <div className="card bg-base-100 shadow-xl">
            <div className="card-body">
              <h2 className="card-title text-2xl mb-4">➕ Добавить район</h2>

              {/* City Selector */}
              <div className="form-control mb-4">
                <label className="label">
                  <span className="label-text font-semibold">
                    🏙️ Выбрать город для загрузки всех районов
                  </span>
                </label>
                <div className="flex gap-2">
                  <select
                    className="select select-bordered flex-1"
                    value={selectedCity}
                    onChange={(e) => setSelectedCity(e.target.value)}
                  >
                    <option value="">Выберите город...</option>
                    <option value="Novi Sad">Нови-Сад</option>
                    <option value="Belgrade">Белград</option>
                    <option value="Niš">Ниш</option>
                    <option value="Kragujevac">Крагуевац</option>
                    <option value="Subotica">Суботица</option>
                    <option value="Zrenjanin">Зренянин</option>
                    <option value="Pančevo">Панчево</option>
                  </select>
                  <button
                    className={`btn btn-secondary ${loadingCityDistricts ? 'loading' : ''}`}
                    onClick={() => loadAllCityDistricts(selectedCity)}
                    disabled={!selectedCity || loadingCityDistricts}
                  >
                    {loadingCityDistricts
                      ? 'Загрузка...'
                      : '📍 Загрузить районы'}
                  </button>
                </div>
              </div>

              {/* City Districts List */}
              {cityDistricts.length > 0 && (
                <div className="mb-6">
                  <h3 className="font-semibold mb-3">
                    Районы города {selectedCity} ({cityDistricts.length}):
                  </h3>
                  <div className="max-h-64 overflow-y-auto space-y-2 border rounded-lg p-2 bg-base-200">
                    {cityDistricts.map((district: any) => (
                      <div
                        key={`${district.osm_type}-${district.osm_id}`}
                        className="flex justify-between items-center p-2 bg-base-100 rounded hover:bg-primary/10 transition-colors"
                      >
                        <div className="flex-1">
                          <div className="font-medium">
                            {district.name}
                            {district.nameEn &&
                              district.nameEn !== district.name && (
                                <span className="ml-2 text-sm text-base-content/60">
                                  ({district.nameEn})
                                </span>
                              )}
                          </div>
                          <div className="text-xs text-base-content/40">
                            {district.type} • OSM: {district.osm_type}/
                            {district.osm_id}
                          </div>
                        </div>
                        <button
                          className="btn btn-xs btn-primary"
                          onClick={() =>
                            handleAddDistrict({
                              ...district,
                              display_name: `${district.name}, ${selectedCity}, Serbia`,
                            })
                          }
                          disabled={loading}
                        >
                          Добавить
                        </button>
                      </div>
                    ))}
                  </div>
                  <div className="mt-2 flex gap-2">
                    <button
                      className="btn btn-sm btn-success"
                      onClick={async () => {
                        if (
                          confirm(
                            `Добавить все ${cityDistricts.length} районов?`
                          )
                        ) {
                          for (const district of cityDistricts) {
                            await handleAddDistrict({
                              ...district,
                              display_name: `${district.name}, ${selectedCity}, Serbia`,
                            });
                          }
                        }
                      }}
                      disabled={loading}
                    >
                      Добавить все районы
                    </button>
                    <button
                      className="btn btn-sm btn-ghost"
                      onClick={() => {
                        setCityDistricts([]);
                        setSelectedCity('');
                      }}
                    >
                      Очистить
                    </button>
                  </div>
                </div>
              )}

              <div className="divider">ИЛИ</div>

              <div className="form-control">
                <label className="label">
                  <span className="label-text">Название района</span>
                </label>

                {/* Переключатель источника данных */}
                <div className="flex gap-2 mb-2">
                  <div className="form-control">
                    <label className="label cursor-pointer">
                      <input
                        type="radio"
                        name="dataSource"
                        className="radio radio-primary radio-sm"
                        checked={dataSource === 'nominatim'}
                        onChange={() => setDataSource('nominatim')}
                      />
                      <span className="label-text ml-2">Nominatim API</span>
                    </label>
                  </div>
                  <div className="form-control">
                    <label className="label cursor-pointer">
                      <input
                        type="radio"
                        name="dataSource"
                        className="radio radio-primary radio-sm"
                        checked={dataSource === 'overpass'}
                        onChange={() => setDataSource('overpass')}
                      />
                      <span className="label-text ml-2">Overpass API</span>
                    </label>
                  </div>
                </div>

                {/* Выбор области поиска для Overpass API */}
                {dataSource === 'overpass' && (
                  <div className="form-control mb-2">
                    <label className="label">
                      <span className="label-text text-sm">Область поиска</span>
                    </label>
                    <select
                      className="select select-bordered select-sm"
                      value={searchArea}
                      onChange={(e) => setSearchArea(e.target.value)}
                    >
                      <option value="">
                        Широкий поиск (вся область Нови-Сада)
                      </option>
                      <option value="Novi Sad">
                        Только в границах Нови-Сада
                      </option>
                      <option value="Serbia">Вся Сербия</option>
                    </select>
                    <div className="text-xs text-base-content/60 mt-1">
                      💡 Для пригородов выберите &quot;Широкий поиск&quot;
                    </div>
                  </div>
                )}

                <div className="flex gap-2">
                  <input
                    type="text"
                    placeholder={
                      dataSource === 'overpass'
                        ? 'Например: Veternik, Liman, Futog'
                        : 'Например: Veternik, Novi Sad, Serbia'
                    }
                    className="input input-bordered flex-1"
                    value={newDistrictName}
                    onChange={(e) => setNewDistrictName(e.target.value)}
                    onKeyPress={(e) => e.key === 'Enter' && searchDistrict()}
                  />
                  <button
                    className={`btn btn-primary ${searching ? 'loading' : ''}`}
                    onClick={searchDistrict}
                    disabled={searching || !newDistrictName.trim()}
                  >
                    {searching ? 'Поиск...' : '🔍 Найти'}
                  </button>
                </div>
                <div className="text-xs text-base-content/60 mt-1">
                  {dataSource === 'overpass'
                    ? '💡 Overpass API ищет и по полигонам, включает пригороды как Veternik'
                    : '💡 Nominatim API - основной поиск с ограничением по городу'}
                </div>
              </div>

              {/* Quick suggestions */}
              <div className="mt-4">
                <p className="text-sm text-base-content/70 mb-2">
                  Быстрый выбор:
                </p>
                <div className="flex flex-wrap gap-2">
                  {suggestedDistricts.map((district) => (
                    <button
                      key={district}
                      className="btn btn-sm btn-outline"
                      onClick={() => setNewDistrictName(district)}
                    >
                      {district}
                    </button>
                  ))}
                </div>

                {/* Специальная секция для Liman I-IV */}
                <div className="mt-4 p-3 bg-primary/5 rounded-lg border border-primary/20">
                  <p className="text-sm font-semibold text-primary mb-2">
                    🏘️ Разделить район Liman на части:
                  </p>
                  <div className="flex flex-wrap gap-2">
                    <button
                      className="btn btn-xs btn-primary btn-outline"
                      onClick={() => {
                        setNewDistrictName('Liman I');
                        setMessage({
                          type: 'info',
                          text: 'ℹ️ Для добавления Liman I потребуется вручную разделить координаты большого района Liman',
                        });
                      }}
                    >
                      Liman I
                    </button>
                    <button
                      className="btn btn-xs btn-primary btn-outline"
                      onClick={() => {
                        setNewDistrictName('Liman II');
                        setMessage({
                          type: 'info',
                          text: 'ℹ️ Для добавления Liman II потребуется вручную разделить координаты большого района Liman',
                        });
                      }}
                    >
                      Liman II
                    </button>
                    <button
                      className="btn btn-xs btn-primary btn-outline"
                      onClick={() => {
                        setNewDistrictName('Liman III');
                        setMessage({
                          type: 'info',
                          text: 'ℹ️ Для добавления Liman III потребуется вручную разделить координаты большого района Liman',
                        });
                      }}
                    >
                      Liman III
                    </button>
                    <button
                      className="btn btn-xs btn-primary btn-outline"
                      onClick={() => {
                        setNewDistrictName('Liman IV');
                        setMessage({
                          type: 'info',
                          text: 'ℹ️ Для добавления Liman IV потребуется вручную разделить координаты большого района Liman',
                        });
                      }}
                    >
                      Liman IV
                    </button>
                  </div>
                  <p className="text-xs text-base-content/60 mt-2">
                    Эти районы требуют ручного разделения координат, так как в
                    OSM они объединены
                  </p>
                </div>
              </div>

              {/* Search Results */}
              {searchResults.length > 0 && (
                <div className="mt-6">
                  <h3 className="font-semibold mb-3">Найденные результаты:</h3>
                  <div className="space-y-2 max-h-96 overflow-y-auto">
                    {searchResults.map((result: any) => {
                      // Извлекаем название района более умно
                      const fullName = result.display_name || '';
                      const nameParts = fullName.split(',');
                      const districtName = nameParts[0].trim();

                      // Проверяем является ли это вариантом Liman
                      const isLimanVariant = districtName
                        .toLowerCase()
                        .includes('liman');

                      return (
                        <div
                          key={result.osm_id}
                          className={`p-3 rounded-lg flex justify-between items-center ${
                            isLimanVariant
                              ? 'bg-primary/10 border-2 border-primary'
                              : 'bg-base-200'
                          }`}
                        >
                          <div className="flex-1">
                            <div className="font-medium text-lg">
                              {districtName}
                              {isLimanVariant && (
                                <span className="ml-2 badge badge-primary badge-sm">
                                  Район Liman
                                </span>
                              )}
                            </div>
                            <div className="text-sm text-base-content/60">
                              {nameParts.slice(1, 3).join(', ')}
                            </div>
                            <div className="text-xs text-base-content/40 mt-1">
                              OSM ID: {result.osm_id} | Тип: {result.type} |{' '}
                              {result.class}
                              {result.source && (
                                <span
                                  className={`ml-2 ${result.source === 'overpass' ? 'text-primary' : 'text-secondary'}`}
                                >
                                  (через{' '}
                                  {result.source === 'overpass'
                                    ? 'Overpass'
                                    : 'Nominatim'}{' '}
                                  API)
                                </span>
                              )}
                            </div>
                          </div>
                          <button
                            className="btn btn-sm btn-success ml-4"
                            onClick={() => handleAddDistrict(result)}
                            disabled={loading}
                          >
                            Добавить
                          </button>
                        </div>
                      );
                    })}
                  </div>
                </div>
              )}

              {/* Instructions */}
              <div className="mt-6 p-4 bg-info/10 rounded-lg">
                <h4 className="font-semibold text-info mb-2">
                  💡 Как добавить район:
                </h4>
                <ol className="text-sm space-y-1 ml-4">
                  <li>1. Введите название района на английском или сербском</li>
                  <li>
                    2. Нажмите &quot;Найти&quot; для поиска в OpenStreetMap
                  </li>
                  <li>3. Выберите правильный результат из списка</li>
                  <li>4. Нажмите &quot;Добавить&quot; для добавления района</li>
                </ol>
                <div className="mt-3 p-2 bg-success/10 rounded text-sm">
                  <strong>✅ Исправлено:</strong> Overpass API теперь ищет и по
                  полигонам (way), включает административные границы и пригороды
                  как Ветерник! Выберите область поиска для точности.
                </div>
              </div>
            </div>
          </div>

          {/* Delete Districts Section */}
          <div className="card bg-base-100 shadow-xl">
            <div className="card-body">
              <div className="flex justify-between items-center mb-4">
                <h2 className="card-title text-2xl">🗑️ Удалить районы</h2>
                <div className="text-sm text-base-content/60">
                  Выбрано: {selectedDistricts.size} из {districts.length}
                </div>
              </div>

              {/* District List */}
              <div className="max-h-96 overflow-y-auto space-y-2">
                {districts.map((district) => (
                  <label
                    key={district.id}
                    className={`flex items-center p-3 rounded-lg cursor-pointer transition-all ${
                      selectedDistricts.has(district.id)
                        ? 'bg-error/20 border-2 border-error'
                        : 'bg-base-200 hover:bg-base-300'
                    }`}
                  >
                    <input
                      type="checkbox"
                      className="checkbox checkbox-error"
                      checked={selectedDistricts.has(district.id)}
                      onChange={() => handleToggleDistrict(district.id)}
                    />
                    <div className="ml-3 flex-1">
                      <div className="font-medium">{district.name}</div>
                      <div className="text-sm text-base-content/60">
                        {district.nameEn} • {district.points} точек
                      </div>
                    </div>
                    <div className="text-xs font-mono text-base-content/50">
                      {district.id}
                    </div>
                  </label>
                ))}
              </div>

              {/* Actions */}
              <div className="mt-6 flex gap-2">
                <button
                  className="btn btn-outline flex-1"
                  onClick={() => {
                    if (selectedDistricts.size === districts.length) {
                      setSelectedDistricts(new Set());
                    } else {
                      setSelectedDistricts(new Set(districts.map((d) => d.id)));
                    }
                  }}
                >
                  {selectedDistricts.size === districts.length
                    ? 'Снять все'
                    : 'Выбрать все'}
                </button>
                <button
                  className={`btn btn-error flex-1 ${loading ? 'loading' : ''}`}
                  onClick={handleDeleteSelected}
                  disabled={selectedDistricts.size === 0 || loading}
                >
                  Удалить выбранные
                </button>
              </div>

              {/* Instructions */}
              <div className="mt-6 p-4 bg-warning/10 rounded-lg">
                <h4 className="font-semibold text-warning mb-2">⚠️ Важно:</h4>
                <ul className="text-sm space-y-1 ml-4">
                  <li>• Удаление создает резервную копию</li>
                  <li>• Изменения применяются после выполнения команды</li>
                  <li>• Для восстановления используйте бэкап файл</li>
                </ul>
              </div>
            </div>
          </div>
        </div>

        {/* Command Output Section */}
        <div className="card bg-base-100 shadow-xl mt-8">
          <div className="card-body">
            <h2 className="card-title text-xl mb-4">
              📋 Команды для выполнения
            </h2>

            <div className="grid grid-cols-1 md:grid-cols-2 gap-4">
              <div className="p-4 bg-base-200 rounded-lg">
                <h3 className="font-semibold mb-2">Добавление района:</h3>
                <code className="text-sm font-mono block p-2 bg-base-300 rounded">
                  node scripts/novi-sad-districts/auto_add_district.js
                  &quot;District Name&quot; &quot;Novi Sad&quot; Serbia
                </code>
              </div>

              <div className="p-4 bg-base-200 rounded-lg">
                <h3 className="font-semibold mb-2">Удаление районов:</h3>
                <code className="text-sm font-mono block p-2 bg-base-300 rounded">
                  node scripts/novi-sad-districts/safe_remove_districts.js
                  district-id-1 district-id-2
                </code>
              </div>
            </div>

            <div className="mt-4 p-4 bg-success/10 rounded-lg">
              <h3 className="font-semibold text-success mb-2">
                ✅ Доступные скрипты:
              </h3>
              <ul className="text-sm space-y-1 ml-4">
                <li>
                  •{' '}
                  <code className="font-mono bg-base-300 px-1 rounded">
                    scripts/novi-sad-districts/auto_add_district.js
                  </code>{' '}
                  - автоматическое добавление района из OSM
                </li>
                <li>
                  •{' '}
                  <code className="font-mono bg-base-300 px-1 rounded">
                    scripts/novi-sad-districts/safe_remove_districts.js
                  </code>{' '}
                  - безопасное удаление с бэкапом
                </li>
                <li>
                  •{' '}
                  <code className="font-mono bg-base-300 px-1 rounded">
                    scripts/novi-sad-districts/fix_syntax.js
                  </code>{' '}
                  - исправление синтаксических ошибок
                </li>
              </ul>
            </div>
          </div>
        </div>
      </div>
    </div>
  );
};

export default DistrictManager;
