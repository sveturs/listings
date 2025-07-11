# Контекст для продолжения Phase 1.4 - Зона пешей доступности

## Текущий статус

### ✅ Что уже сделано:
1. **UI компонент создан и интегрирован**:
   - `WalkingAccessibilityControl` полностью функционален
   - Интегрирован в `/app/[locale]/map/page.tsx` 
   - Добавлен в мобильные фильтры
   - Переводы добавлены (ru/en)

2. **Функции изохрон готовы**:
   - `/frontend/svetu/src/components/GIS/utils/isochrone.ts`
   - `generateIsochrone()` - базовая функция
   - `generateStylizedIsochrone()` - с волнообразными краями
   - Конвертеры времени/расстояния

3. **Хук для расчета расстояний**:
   - `/frontend/svetu/src/components/GIS/hooks/useDistanceCalculation.ts`
   - Использует Turf.js для точных расчетов

## ❌ Что НЕ работает:
Переключатель режима не влияет на отображение зоны поиска на карте. Всегда показывается круг радиуса.

## 🔧 Что нужно доделать:

### 1. В `/app/[locale]/map/page.tsx`:
```typescript
// Сейчас состояния есть:
const [walkingMode, setWalkingMode] = useState<'radius' | 'walking'>('radius');
const [walkingTime, setWalkingTime] = useState(15);

// Но они НЕ передаются в StorefrontMap!
// Нужно передать эти props в компонент StorefrontMap
```

### 2. В `StorefrontMap` компоненте:
- Принять новые props: `walkingMode` и `walkingTime`
- Передать их дальше в `InteractiveMap`

### 3. В `InteractiveMap` компоненте:
- Принять props для режима отображения
- Модифицировать `radiusCircleGeoJSON` чтобы:
  - При `mode === 'radius'` - использовать текущий circle из Turf
  - При `mode === 'walking'` - использовать `generateStylizedIsochrone()`

### 4. Визуальное отличие:
- Круг радиуса: синий с прозрачной заливкой
- Изохрона: зеленоватый оттенок с волнообразными краями

## 📁 Ключевые файлы для доработки:

1. `/data/hostel-booking-system/frontend/svetu/src/app/[locale]/map/page.tsx`
2. `/data/hostel-booking-system/frontend/svetu/src/components/GIS/Map/StorefrontMap.tsx`
3. `/data/hostel-booking-system/frontend/svetu/src/components/GIS/Map/InteractiveMap.tsx`

## 💡 Подсказка для следующей сессии:

В `InteractiveMap.tsx` нужно изменить блок (строки ~285-302):
```typescript
const radiusCircleGeoJSON = useMemo(() => {
  if (!showBuyerMarker) return null;

  const center = [
    internalBuyerLocation.longitude,
    internalBuyerLocation.latitude,
  ];
  
  // Здесь добавить условие:
  if (walkingMode === 'walking') {
    return generateStylizedIsochrone(center, walkingTime);
  } else {
    const radiusInKm = searchRadius / 1000;
    return circle(center, radiusInKm, {
      steps: 64,
      units: 'kilometers',
    });
  }
}, [showBuyerMarker, internalBuyerLocation, searchRadius, walkingMode, walkingTime]);
```

## 🎯 Ожидаемый результат:
- При переключении на "Пешком" - вместо идеального круга показывается изохрона
- Изохрона учитывает реальное время ходьбы (5 км/ч)
- Визуально отличается от круга радиуса
- Фильтрация объявлений работает для обоих режимов