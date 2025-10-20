# 📋 План улучшения покрытия тестами Frontend

**Дата создания:** 2025-10-20
**Текущее покрытие:** 64.89% (Statements), 57.8% (Branches)
**Целевое покрытие:** 75%+ (Statements), 70%+ (Branches)

---

## 🎯 Приоритет 1: Критические компоненты (3-5%)

### 1. AutocompleteAttributeField.tsx
**Текущее покрытие:** 3.03% ⚠️
**Целевое покрытие:** 80%+
**Файл:** `src/components/shared/AutocompleteAttributeField.tsx`

#### Анализ кода:
- **330 строк** сложного компонента с автокомплитом
- Использует `useAttributeAutocomplete` хук
- Управляет состоянием: `inputValue`, `showSuggestions`, `selectedIndex`, `suggestions`
- Генерирует умные предложения на основе типа атрибута
- Обрабатывает клавиатурную навигацию (Arrow Up/Down, Enter, Escape)

#### Что тестировать:

**Базовая функциональность:**
```typescript
describe('AutocompleteAttributeField', () => {
  const mockAttribute = {
    id: 1,
    name: 'brand',
    display_name: 'Бренд',
    is_required: false,
    options: ['Apple', 'Samsung', 'Xiaomi']
  };

  const mockOnChange = jest.fn();

  test('рендерит поле ввода с правильным placeholder', () => {
    render(
      <AutocompleteAttributeField
        attribute={mockAttribute}
        onChange={mockOnChange}
      />
    );
    expect(screen.getByPlaceholderText('Бренд')).toBeInTheDocument();
  });

  test('показывает required индикатор если is_required=true', () => {
    render(
      <AutocompleteAttributeField
        attribute={{...mockAttribute, is_required: true}}
        onChange={mockOnChange}
      />
    );
    expect(screen.getByText('*')).toBeInTheDocument();
  });

  test('вызывает onChange при вводе текста', () => {
    render(
      <AutocompleteAttributeField
        attribute={mockAttribute}
        onChange={mockOnChange}
      />
    );

    const input = screen.getByPlaceholderText('Бренд');
    fireEvent.change(input, { target: { value: 'Apple' } });

    expect(mockOnChange).toHaveBeenCalledWith({
      attribute_id: 1,
      text_value: 'Apple'
    });
  });
});
```

**Автокомплит предложений:**
```typescript
test('показывает предложения при фокусе', async () => {
  const { useAttributeAutocomplete } = require('@/hooks/useAttributeAutocomplete');

  useAttributeAutocomplete.mockReturnValue({
    getFilteredSuggestions: () => [
      { value: 'Apple', type: 'popular' },
      { value: 'Samsung', type: 'recent' }
    ],
    saveValue: jest.fn()
  });

  render(<AutocompleteAttributeField attribute={mockAttribute} onChange={mockOnChange} />);

  const input = screen.getByPlaceholderText('Бренд');
  fireEvent.focus(input);

  await waitFor(() => {
    expect(screen.getByText('Apple')).toBeInTheDocument();
    expect(screen.getByText('Samsung')).toBeInTheDocument();
  });
});

test('скрывает предложения при выборе', async () => {
  // Mock hook
  render(<AutocompleteAttributeField attribute={mockAttribute} onChange={mockOnChange} />);

  const input = screen.getByPlaceholderText('Бренд');
  fireEvent.focus(input);

  await waitFor(() => screen.getByText('Apple'));

  fireEvent.click(screen.getByText('Apple'));

  await waitFor(() => {
    expect(screen.queryByText('Samsung')).not.toBeInTheDocument();
  });
  expect(input).toHaveValue('Apple');
});
```

**Клавиатурная навигация:**
```typescript
test('навигация по предложениям стрелками', async () => {
  render(<AutocompleteAttributeField attribute={mockAttribute} onChange={mockOnChange} />);

  const input = screen.getByPlaceholderText('Бренд');
  fireEvent.focus(input);

  await waitFor(() => screen.getByText('Apple'));

  // Arrow Down
  fireEvent.keyDown(input, { key: 'ArrowDown' });
  expect(screen.getByText('Apple')).toHaveClass('bg-primary');

  // Arrow Down again
  fireEvent.keyDown(input, { key: 'ArrowDown' });
  expect(screen.getByText('Samsung')).toHaveClass('bg-primary');

  // Enter to select
  fireEvent.keyDown(input, { key: 'Enter' });
  expect(input).toHaveValue('Samsung');
});

test('Escape закрывает предложения', async () => {
  render(<AutocompleteAttributeField attribute={mockAttribute} onChange={mockOnChange} />);

  const input = screen.getByPlaceholderText('Бренд');
  fireEvent.focus(input);

  await waitFor(() => screen.getByText('Apple'));

  fireEvent.keyDown(input, { key: 'Escape' });

  await waitFor(() => {
    expect(screen.queryByText('Apple')).not.toBeInTheDocument();
  });
});
```

**Умные предложения (generateSmartSuggestions):**
```typescript
test('генерирует умные предложения для цен', () => {
  const priceAttribute = { ...mockAttribute, name: 'price', display_name: 'Цена' };

  render(<AutocompleteAttributeField attribute={priceAttribute} onChange={mockOnChange} />);

  const input = screen.getByPlaceholderText('Цена');
  fireEvent.focus(input);

  // Должны быть предложения: 50000, 100000, 150000...
  await waitFor(() => {
    expect(screen.getByText(/50000/)).toBeInTheDocument();
  });
});

test('генерирует умные предложения для годов', () => {
  const yearAttribute = { ...mockAttribute, name: 'year', display_name: 'Год' };

  render(<AutocompleteAttributeField attribute={yearAttribute} onChange={mockOnChange} />);

  const input = screen.getByPlaceholderText('Год');
  fireEvent.focus(input);

  // Должны быть предложения: 2024, 2023, 2022...
  await waitFor(() => {
    expect(screen.getByText('2024')).toBeInTheDocument();
  });
});
```

**Иконки предложений:**
```typescript
test('показывает правильные иконки для типов предложений', async () => {
  // Mock hook to return different types
  useAttributeAutocomplete.mockReturnValue({
    getFilteredSuggestions: () => [
      { value: 'Apple', type: 'exact' },
      { value: 'Samsung', type: 'popular' },
      { value: 'Xiaomi', type: 'recent' }
    ],
    saveValue: jest.fn()
  });

  render(<AutocompleteAttributeField attribute={mockAttribute} onChange={mockOnChange} />);

  const input = screen.getByPlaceholderText('Бренд');
  fireEvent.focus(input);

  await waitFor(() => {
    expect(screen.getByText('🎯')).toBeInTheDocument(); // exact
    expect(screen.getByText('⭐')).toBeInTheDocument(); // popular
    expect(screen.getByText('🕒')).toBeInTheDocument(); // recent
  });
});
```

**Приоритет:** 🔴 **Критический**
**Оценка времени:** 4-6 часов

---

### 2. useAttributeAutocomplete.ts
**Текущее покрытие:** 4.27% ⚠️
**Целевое покрытие:** 80%+
**Файл:** `src/hooks/useAttributeAutocomplete.ts`

#### Анализ кода:
- **295 строк** кастомного хука
- Управляет localStorage для популярных и недавних значений
- Debouncing для оптимизации записи
- Очистка старых данных при превышении квоты
- Фильтрация и ранжирование предложений

#### Что тестировать:

**Базовая функциональность:**
```typescript
describe('useAttributeAutocomplete', () => {
  beforeEach(() => {
    localStorage.clear();
  });

  test('инициализируется с пустыми значениями', () => {
    const { result } = renderHook(() =>
      useAttributeAutocomplete({ attributeId: 1, attributeName: 'brand' })
    );

    expect(result.current.popularValues).toEqual([]);
    expect(result.current.recentValues).toEqual([]);
  });

  test('загружает данные из localStorage', () => {
    localStorage.setItem('recent_v1_1', JSON.stringify(['Apple', 'Samsung']));
    localStorage.setItem('popular_v1_brand', JSON.stringify(['Xiaomi', 'Huawei']));

    const { result } = renderHook(() =>
      useAttributeAutocomplete({ attributeId: 1, attributeName: 'brand' })
    );

    expect(result.current.recentValues).toEqual(['Apple', 'Samsung']);
    expect(result.current.popularValues).toEqual(['Xiaomi', 'Huawei']);
  });
});
```

**saveValue и addRecentValue:**
```typescript
test('добавляет значение в недавние', async () => {
  const { result } = renderHook(() =>
    useAttributeAutocomplete({ attributeId: 1, attributeName: 'brand' })
  );

  act(() => {
    result.current.addRecentValue('Apple');
  });

  expect(result.current.recentValues).toContain('Apple');
});

test('ограничивает недавние значения до MAX_RECENT_VALUES (5)', async () => {
  const { result } = renderHook(() =>
    useAttributeAutocomplete({ attributeId: 1, attributeName: 'brand' })
  );

  act(() => {
    result.current.addRecentValue('Apple');
    result.current.addRecentValue('Samsung');
    result.current.addRecentValue('Xiaomi');
    result.current.addRecentValue('Huawei');
    result.current.addRecentValue('Sony');
    result.current.addRecentValue('LG'); // 6-е значение
  });

  expect(result.current.recentValues).toHaveLength(5);
  expect(result.current.recentValues[0]).toBe('LG'); // Последнее добавленное
});

test('перемещает существующее значение в начало', () => {
  const { result } = renderHook(() =>
    useAttributeAutocomplete({ attributeId: 1, attributeName: 'brand' })
  );

  act(() => {
    result.current.addRecentValue('Apple');
    result.current.addRecentValue('Samsung');
    result.current.addRecentValue('Apple'); // Повторное добавление
  });

  expect(result.current.recentValues[0]).toBe('Apple');
  expect(result.current.recentValues).toHaveLength(2); // Без дубликатов
});
```

**incrementPopularity:**
```typescript
test('увеличивает популярность значения', () => {
  const { result } = renderHook(() =>
    useAttributeAutocomplete({ attributeId: 1, attributeName: 'brand' })
  );

  act(() => {
    result.current.incrementPopularity('Apple');
    result.current.incrementPopularity('Apple');
    result.current.incrementPopularity('Samsung');
  });

  // Apple должен быть первым (2 раза vs 1 раз)
  expect(result.current.popularValues[0]).toBe('Apple');
});
```

**getFilteredSuggestions:**
```typescript
test('возвращает все предложения для пустого запроса', () => {
  const { result } = renderHook(() =>
    useAttributeAutocomplete({ attributeId: 1, attributeName: 'brand' })
  );

  act(() => {
    result.current.saveValue('Apple');
    result.current.saveValue('Samsung');
  });

  const suggestions = result.current.getFilteredSuggestions('');
  expect(suggestions.length).toBeGreaterThan(0);
});

test('фильтрует по запросу (startsWith)', () => {
  const { result } = renderHook(() =>
    useAttributeAutocomplete({ attributeId: 1, attributeName: 'brand' })
  );

  act(() => {
    result.current.saveValue('Apple');
    result.current.saveValue('Samsung');
    result.current.saveValue('Xiaomi');
  });

  const suggestions = result.current.getFilteredSuggestions('Sam');
  expect(suggestions).toHaveLength(1);
  expect(suggestions[0].value).toBe('Samsung');
});

test('фильтрует по запросу (contains)', () => {
  const { result } = renderHook(() =>
    useAttributeAutocomplete({ attributeId: 1, attributeName: 'brand' })
  );

  act(() => {
    result.current.saveValue('iPhone 15');
    result.current.saveValue('Samsung Galaxy');
  });

  const suggestions = result.current.getFilteredSuggestions('phone');
  expect(suggestions[0].value).toBe('iPhone 15');
});

test('ранжирует точное совпадение выше', () => {
  const { result } = renderHook(() =>
    useAttributeAutocomplete({ attributeId: 1, attributeName: 'brand' })
  );

  act(() => {
    result.current.saveValue('Apple iPhone');
    result.current.saveValue('Apple');
  });

  const suggestions = result.current.getFilteredSuggestions('Apple');
  expect(suggestions[0].value).toBe('Apple'); // Точное совпадение первым
});
```

**Debouncing и localStorage:**
```typescript
test('сохраняет в localStorage с debouncing', async () => {
  jest.useFakeTimers();

  const { result } = renderHook(() =>
    useAttributeAutocomplete({ attributeId: 1, attributeName: 'brand' })
  );

  act(() => {
    result.current.addRecentValue('Apple');
  });

  // До истечения debounce - не должно быть в localStorage
  expect(localStorage.getItem('recent_v1_1')).toBeNull();

  // Ждем debounce (100ms)
  act(() => {
    jest.advanceTimersByTime(100);
  });

  await waitFor(() => {
    expect(localStorage.getItem('recent_v1_1')).toBeTruthy();
  });

  jest.useRealTimers();
});
```

**clearOldStorageData:**
```typescript
test('очищает старые ключи без версии', () => {
  localStorage.setItem('recent_1', '["old"]');
  localStorage.setItem('recent_v1_1', '["new"]');

  const { result } = renderHook(() =>
    useAttributeAutocomplete({ attributeId: 1, attributeName: 'brand' })
  );

  act(() => {
    result.current.clearData();
  });

  // Только новый ключ с версией должен остаться
  expect(localStorage.getItem('recent_1')).toBeNull();
});
```

**Приоритет:** 🔴 **Критический**
**Оценка времени:** 4-5 часов

---

### 3. cars.ts Service
**Текущее покрытие:** 5.71% ⚠️
**Целевое покрытие:** 80%+
**Файл:** `src/services/cars.ts`

#### Анализ кода:
- **145 строк** API сервиса
- 4 метода: `getMakes()`, `getModelsByMake()`, `getGenerationsByModel()`, `searchMakes()`
- Использует fetch API
- Обрабатывает ошибки

#### Что тестировать:

**getMakes:**
```typescript
describe('CarsService', () => {
  beforeEach(() => {
    global.fetch = jest.fn();
  });

  afterEach(() => {
    jest.resetAllMocks();
  });

  describe('getMakes', () => {
    test('возвращает список марок при успешном запросе', async () => {
      const mockMakes = [
        { id: 1, name: 'BMW', slug: 'bmw' },
        { id: 2, name: 'Mercedes', slug: 'mercedes' }
      ];

      (global.fetch as jest.Mock).mockResolvedValueOnce({
        ok: true,
        json: async () => ({ data: mockMakes })
      });

      const result = await CarsService.getMakes();

      expect(result.success).toBe(true);
      expect(result.data).toEqual(mockMakes);
      expect(global.fetch).toHaveBeenCalledWith(
        expect.stringContaining('/api/v1/cars/makes'),
        expect.objectContaining({ method: 'GET' })
      );
    });

    test('обрабатывает ошибку HTTP', async () => {
      (global.fetch as jest.Mock).mockResolvedValueOnce({
        ok: false,
        status: 404
      });

      const result = await CarsService.getMakes();

      expect(result.success).toBe(false);
      expect(result.error).toContain('404');
    });

    test('обрабатывает network ошибку', async () => {
      (global.fetch as jest.Mock).mockRejectedValueOnce(
        new Error('Network error')
      );

      const result = await CarsService.getMakes();

      expect(result.success).toBe(false);
      expect(result.error).toBe('Network error');
    });

    test('обрабатывает данные без обертки .data', async () => {
      const mockMakes = [{ id: 1, name: 'BMW' }];

      (global.fetch as jest.Mock).mockResolvedValueOnce({
        ok: true,
        json: async () => mockMakes // Без обертки
      });

      const result = await CarsService.getMakes();

      expect(result.success).toBe(true);
      expect(result.data).toEqual(mockMakes);
    });
  });
});
```

**getModelsByMake:**
```typescript
describe('getModelsByMake', () => {
  test('возвращает модели для указанной марки', async () => {
    const mockModels = [
      { id: 1, name: 'X5', make_id: 1 },
      { id: 2, name: 'X7', make_id: 1 }
    ];

    (global.fetch as jest.Mock).mockResolvedValueOnce({
      ok: true,
      json: async () => ({ data: mockModels })
    });

    const result = await CarsService.getModelsByMake('bmw');

    expect(result.success).toBe(true);
    expect(result.data).toEqual(mockModels);
    expect(global.fetch).toHaveBeenCalledWith(
      expect.stringContaining('/makes/bmw/models'),
      expect.any(Object)
    );
  });

  test('правильно кодирует slug с пробелами', async () => {
    (global.fetch as jest.Mock).mockResolvedValueOnce({
      ok: true,
      json: async () => ({ data: [] })
    });

    await CarsService.getModelsByMake('aston-martin');

    expect(global.fetch).toHaveBeenCalledWith(
      expect.stringContaining('/makes/aston-martin/models'),
      expect.any(Object)
    );
  });
});
```

**getGenerationsByModel:**
```typescript
describe('getGenerationsByModel', () => {
  test('возвращает поколения для модели', async () => {
    const mockGenerations = [
      { id: 1, name: 'F15 (2013-2018)', model_id: 10 },
      { id: 2, name: 'G05 (2018-present)', model_id: 10 }
    ];

    (global.fetch as jest.Mock).mockResolvedValueOnce({
      ok: true,
      json: async () => ({ data: mockGenerations })
    });

    const result = await CarsService.getGenerationsByModel(10);

    expect(result.success).toBe(true);
    expect(result.data).toEqual(mockGenerations);
    expect(global.fetch).toHaveBeenCalledWith(
      expect.stringContaining('/models/10/generations'),
      expect.any(Object)
    );
  });
});
```

**searchMakes:**
```typescript
describe('searchMakes', () => {
  test('ищет марки по запросу', async () => {
    const mockResults = [
      { id: 1, name: 'BMW', slug: 'bmw' }
    ];

    (global.fetch as jest.Mock).mockResolvedValueOnce({
      ok: true,
      json: async () => ({ data: mockResults })
    });

    const result = await CarsService.searchMakes('BM');

    expect(result.success).toBe(true);
    expect(result.data).toEqual(mockResults);
    expect(global.fetch).toHaveBeenCalledWith(
      expect.stringContaining('/makes/search?q=BM'),
      expect.any(Object)
    );
  });

  test('правильно кодирует спецсимволы в запросе', async () => {
    (global.fetch as jest.Mock).mockResolvedValueOnce({
      ok: true,
      json: async () => ({ data: [] })
    });

    await CarsService.searchMakes('BMW & Mercedes');

    expect(global.fetch).toHaveBeenCalledWith(
      expect.stringContaining(encodeURIComponent('BMW & Mercedes')),
      expect.any(Object)
    );
  });
});
```

**Приоритет:** 🔴 **Критический**
**Оценка времени:** 2-3 часа

---

## 🎯 Приоритет 2: Важные утилиты (20-40%)

### 4. iconMapper.tsx
**Текущее покрытие:** 20% ⚠️
**Целевое покрытие:** 80%+
**Файл:** `src/utils/iconMapper.tsx`

#### Анализ кода:
- **128 строк** маппинга иконок
- 2 экспортируемые функции: `getCategoryIcon()`, `renderCategoryIcon()`
- Поддержка эмодзи
- Fallback на Package иконку

#### Что тестировать:

```typescript
describe('iconMapper', () => {
  describe('getCategoryIcon', () => {
    test('возвращает правильную иконку для известного имени', () => {
      const IconComponent = getCategoryIcon('car');
      expect(IconComponent).toBeDefined();
      expect(IconComponent).not.toBe(Package); // Не fallback
    });

    test('возвращает Package для неизвестного имени', () => {
      const IconComponent = getCategoryIcon('unknown-icon-name');
      expect(IconComponent).toBe(Package);
    });

    test('возвращает null для пустого имени', () => {
      expect(getCategoryIcon('')).toBeNull();
      expect(getCategoryIcon(undefined)).toBeNull();
    });

    test('не чувствителен к регистру', () => {
      expect(getCategoryIcon('CAR')).toBe(getCategoryIcon('car'));
      expect(getCategoryIcon('Truck')).toBe(getCategoryIcon('truck'));
    });

    test('обрабатывает все транспортные иконки', () => {
      const transportIcons = ['car', 'truck', 'motorcycle', 'ship', 'sailboat'];

      transportIcons.forEach(iconName => {
        const IconComponent = getCategoryIcon(iconName);
        expect(IconComponent).toBeDefined();
        expect(IconComponent).not.toBe(Package);
      });
    });

    test('обрабатывает индустриальные иконки', () => {
      expect(getCategoryIcon('factory')).toBe(Factory);
      expect(getCategoryIcon('tractor')).toBe(Tractor);
      expect(getCategoryIcon('wheat')).toBe(Wheat);
    });
  });

  describe('renderCategoryIcon', () => {
    test('рендерит иконку компонент', () => {
      const { container } = render(
        <>{renderCategoryIcon('car', 'w-6 h-6')}</>
      );

      const svg = container.querySelector('svg');
      expect(svg).toBeInTheDocument();
      expect(svg).toHaveClass('w-6', 'h-6');
    });

    test('рендерит эмодзи как текст', () => {
      const { container } = render(
        <>{renderCategoryIcon('🚗', 'text-2xl')}</>
      );

      const span = container.querySelector('span');
      expect(span).toBeInTheDocument();
      expect(span).toHaveTextContent('🚗');
      expect(span).toHaveClass('text-2xl');
    });

    test('возвращает null для пустого имени', () => {
      expect(renderCategoryIcon('')).toBeNull();
      expect(renderCategoryIcon(undefined)).toBeNull();
    });

    test('применяет custom className', () => {
      const { container } = render(
        <>{renderCategoryIcon('car', 'custom-class')}</>
      );

      const svg = container.querySelector('svg');
      expect(svg).toHaveClass('custom-class');
    });

    test('обрабатывает многобайтные эмодзи', () => {
      const emojis = ['🚗', '🏠', '📱', '⚽'];

      emojis.forEach(emoji => {
        const { container } = render(
          <>{renderCategoryIcon(emoji)}</>
        );

        expect(container.querySelector('span')).toHaveTextContent(emoji);
      });
    });
  });
});
```

**Приоритет:** 🟡 **Средний**
**Оценка времени:** 1-2 часа

---

### 5. env.ts
**Текущее покрытие:** 41.66%
**Целевое покрытие:** 80%+
**Файл:** `src/utils/env.ts`

#### Анализ кода:
- **37 строк** утилит для переменных окружения
- Server-side vs Client-side логика
- Типизированный доступ через `publicEnv`

#### Что тестировать:

```typescript
describe('env', () => {
  describe('getEnv', () => {
    test('возвращает значение из process.env на сервере', () => {
      // Mock window as undefined (server-side)
      const originalWindow = global.window;
      // @ts-ignore
      delete global.window;

      process.env.TEST_VAR = 'server-value';

      expect(getEnv('TEST_VAR')).toBe('server-value');

      // Restore
      global.window = originalWindow;
    });

    test('возвращает defaultValue если переменная не найдена', () => {
      delete process.env.NON_EXISTENT_VAR;

      expect(getEnv('NON_EXISTENT_VAR', 'default')).toBe('default');
    });

    test('использует runtime env на клиенте', () => {
      // Mock window (client-side)
      global.window = {} as any;

      const { env } = require('next-runtime-env');
      env.mockReturnValue('client-value');

      expect(getEnv('TEST_VAR')).toBe('client-value');
    });
  });

  describe('publicEnv', () => {
    test('возвращает правильный API_URL', () => {
      process.env.NEXT_PUBLIC_API_URL = 'http://api.example.com';

      expect(publicEnv.API_URL).toBe('http://api.example.com');
    });

    test('использует дефолтный API_URL если не задан', () => {
      delete process.env.NEXT_PUBLIC_API_URL;

      expect(publicEnv.API_URL).toBe('http://localhost:3000');
    });

    test('возвращает правильный MINIO_URL', () => {
      process.env.NEXT_PUBLIC_MINIO_URL = 'http://minio.example.com';

      expect(publicEnv.MINIO_URL).toBe('http://minio.example.com');
    });

    test('парсит ENABLE_PAYMENTS как boolean', () => {
      process.env.NEXT_PUBLIC_ENABLE_PAYMENTS = 'true';
      expect(publicEnv.ENABLE_PAYMENTS).toBe(true);

      process.env.NEXT_PUBLIC_ENABLE_PAYMENTS = 'false';
      expect(publicEnv.ENABLE_PAYMENTS).toBe(false);

      process.env.NEXT_PUBLIC_ENABLE_PAYMENTS = '';
      expect(publicEnv.ENABLE_PAYMENTS).toBe(false);
    });

    test('возвращает undefined для необязательных переменных', () => {
      delete process.env.NEXT_PUBLIC_WEBSOCKET_URL;

      expect(publicEnv.WEBSOCKET_URL).toBeUndefined();
    });
  });
});
```

**Приоритет:** 🟡 **Средний**
**Оценка времени:** 1-2 часа

---

### 6. config/index.ts
**Текущее покрытие:** 36.14%
**Целевое покрытие:** 70%+
**Файл:** `src/config/index.ts`

#### Анализ кода:
- Сложный класс ConfigManager
- Валидация через Zod схемы
- Lazy initialization
- Server-side vs Client-side логика

#### Что тестировать:

```typescript
describe('ConfigManager', () => {
  beforeEach(() => {
    // Reset singleton
    jest.resetModules();
  });

  test('инициализируется с дефолтными значениями', () => {
    const configManager = require('@/config').default;

    expect(configManager.getApiUrl()).toBe('http://localhost:3000');
  });

  test('валидирует публичные переменные', () => {
    process.env.NEXT_PUBLIC_API_URL = 'invalid-url'; // Невалидный URL

    const configManager = require('@/config').default;

    // Должен использовать дефолт при ошибке валидации
    expect(configManager.getApiUrl()).toBeTruthy();
  });

  test('возвращает IMAGE_HOSTS как массив', () => {
    process.env.NEXT_PUBLIC_IMAGE_HOSTS = 's3.example.com,cdn.example.com';

    const configManager = require('@/config').default;
    const hosts = configManager.getImageHosts();

    expect(Array.isArray(hosts)).toBe(true);
    expect(hosts).toContain('s3.example.com');
    expect(hosts).toContain('cdn.example.com');
  });

  test('обрабатывает пустой IMAGE_HOSTS', () => {
    delete process.env.NEXT_PUBLIC_IMAGE_HOSTS;

    const configManager = require('@/config').default;
    const hosts = configManager.getImageHosts();

    expect(Array.isArray(hosts)).toBe(true);
    expect(hosts).toHaveLength(0);
  });

  test('isPaymentsEnabled возвращает boolean', () => {
    process.env.NEXT_PUBLIC_ENABLE_PAYMENTS = 'true';

    const configManager = require('@/config').default;

    expect(configManager.isPaymentsEnabled()).toBe(true);
  });
});
```

**Приоритет:** 🟢 **Низкий** (сложная интеграция, можно отложить)
**Оценка времени:** 3-4 часа

---

## 📊 Итоговая оценка

| Приоритет | Компонент | Текущее | Цель | Время | Статус |
|-----------|-----------|---------|------|-------|--------|
| 🔴 P1 | AutocompleteAttributeField | 3.03% | 80%+ | 4-6ч | Pending |
| 🔴 P1 | useAttributeAutocomplete | 4.27% | 80%+ | 4-5ч | Pending |
| 🔴 P1 | cars.ts | 5.71% | 80%+ | 2-3ч | Pending |
| 🟡 P2 | iconMapper.tsx | 20% | 80%+ | 1-2ч | Pending |
| 🟡 P2 | env.ts | 41.66% | 80%+ | 1-2ч | Pending |
| 🟢 P3 | config/index.ts | 36.14% | 70%+ | 3-4ч | Pending |

**Общее время:** ~16-22 часа
**Ожидаемый результат:** Покрытие увеличится с 64.89% до ~73-75%

---

## 🚀 Рекомендации по выполнению

### Порядок выполнения:
1. **День 1-2**: AutocompleteAttributeField + useAttributeAutocomplete (критично, связаны)
2. **День 3**: cars.ts (простой, быстрый)
3. **День 4**: iconMapper.tsx + env.ts (утилиты, средний приоритет)
4. **День 5** (опционально): config/index.ts (можно отложить)

### Полезные паттерны:

**1. Mock localStorage:**
```typescript
beforeEach(() => {
  const localStorageMock = {
    getItem: jest.fn(),
    setItem: jest.fn(),
    removeItem: jest.fn(),
    clear: jest.fn(),
  };
  global.localStorage = localStorageMock as any;
});
```

**2. Mock fetch:**
```typescript
global.fetch = jest.fn() as jest.Mock;

(global.fetch as jest.Mock).mockResolvedValueOnce({
  ok: true,
  json: async () => ({ data: mockData })
});
```

**3. Mock hooks:**
```typescript
jest.mock('@/hooks/useAttributeAutocomplete', () => ({
  useAttributeAutocomplete: jest.fn(() => ({
    getFilteredSuggestions: jest.fn(),
    saveValue: jest.fn()
  }))
}));
```

**4. Fake timers для debounce:**
```typescript
jest.useFakeTimers();
act(() => {
  jest.advanceTimersByTime(100);
});
jest.useRealTimers();
```

---

## ✅ Критерии успеха

- [ ] Все новые тесты проходят успешно
- [ ] Покрытие statements > 73%
- [ ] Покрытие branches > 65%
- [ ] Нет console warnings в тестах
- [ ] CI/CD pipeline успешен
- [ ] Документация обновлена

---

**Автор:** Claude Code
**Дата создания:** 2025-10-20
**Версия:** 1.0
