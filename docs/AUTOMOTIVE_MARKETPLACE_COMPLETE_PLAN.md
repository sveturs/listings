# Полный план улучшения системы автомобильных объявлений

## Резюме

Этот документ представляет комплексный план трансформации автомобильного раздела маркетплейса в современную, удобную и высокопроизводительную систему. План охватывает все аспекты: от интеграции внешних API до создания специализированных UI компонентов и оптимизации поисковой системы.

## Цели и метрики успеха

### Основные цели:
1. **Создание лучшего автомобильного маркетплейса в регионе**
2. **Упрощение процесса создания объявлений до 2-3 минут**
3. **Повышение точности поиска до 95%+**
4. **Автоматизация заполнения данных на 80%**

### KPI:
- Время создания объявления: < 3 минут
- Конверсия начатых объявлений в опубликованные: > 85%
- Точность поиска по параметрам: > 95%
- Удовлетворенность пользователей (NPS): > 70

## Архитектура решения

### Компонентная структура:

```
┌─────────────────────────────────────────────────────────────┐
│                     Frontend (Next.js)                       │
├─────────────────────────────────────────────────────────────┤
│  ┌─────────────┐  ┌──────────────┐  ┌─────────────────┐   │
│  │ Car Wizard  │  │ Car Search   │  │ Car Comparison  │   │
│  │ Component   │  │ Component    │  │ Component       │   │
│  └─────────────┘  └──────────────┘  └─────────────────┘   │
├─────────────────────────────────────────────────────────────┤
│                    API Gateway (Go Fiber)                    │
├─────────────────────────────────────────────────────────────┤
│  ┌─────────────┐  ┌──────────────┐  ┌─────────────────┐   │
│  │ Car Service │  │ Search       │  │ Attribute       │   │
│  │             │  │ Service      │  │ Service         │   │
│  └─────────────┘  └──────────────┘  └─────────────────┘   │
├─────────────────────────────────────────────────────────────┤
│  ┌─────────────┐  ┌──────────────┐  ┌─────────────────┐   │
│  │ PostgreSQL  │  │ OpenSearch   │  │ Redis Cache     │   │
│  │ (Master DB) │  │ (Search)     │  │ (Performance)   │   │
│  └─────────────┘  └──────────────┘  └─────────────────┘   │
├─────────────────────────────────────────────────────────────┤
│  ┌─────────────┐  ┌──────────────┐  ┌─────────────────┐   │
│  │ CarAPI.app  │  │ NHTSA API    │  │ Google Vision   │   │
│  │ (External)  │  │ (VIN Decode) │  │ (Image AI)      │   │
│  └─────────────┘  └──────────────┘  └─────────────────┘   │
└─────────────────────────────────────────────────────────────┘
```

## Фаза 1: Базовая интеграция (2 недели)

### 1.1 Интеграция CarAPI.app

#### Backend реализация:

```go
// backend/internal/integrations/carapi/client.go
package carapi

type Client struct {
    token      string
    baseURL    string
    httpClient *http.Client
    cache      cache.Interface
    limiter    *rate.Limiter
}

type CarAPIService interface {
    GetMakes(ctx context.Context) ([]Make, error)
    GetModels(ctx context.Context, makeID int) ([]Model, error)
    GetTrims(ctx context.Context, modelID int, year int) ([]Trim, error)
    DecodeVIN(ctx context.Context, vin string) (*VehicleInfo, error)
}

// Структуры данных
type Make struct {
    ID           int    `json:"id"`
    Name         string `json:"name"`
    Slug         string `json:"slug"`
    LogoURL      string `json:"logo_url,omitempty"`
    CountryCode  string `json:"country_code"`
}

type Model struct {
    ID       int    `json:"id"`
    MakeID   int    `json:"make_id"`
    Name     string `json:"name"`
    Slug     string `json:"slug"`
    BodyType string `json:"body_type,omitempty"`
}

type Trim struct {
    ID           int     `json:"id"`
    ModelID      int     `json:"model_id"`
    Year         int     `json:"year"`
    Name         string  `json:"name"`
    Engine       string  `json:"engine"`
    Transmission string  `json:"transmission"`
    DriveType    string  `json:"drive_type"`
    FuelType     string  `json:"fuel_type"`
    PowerHP      int     `json:"power_hp"`
    TorqueNM     int     `json:"torque_nm"`
    FuelEconomy  float64 `json:"fuel_economy_combined"`
}
```

#### Кеширование стратегия:

```go
// backend/internal/integrations/carapi/cache.go
func (c *Client) GetMakesWithCache(ctx context.Context) ([]Make, error) {
    cacheKey := "carapi:makes:all"
    
    // Проверяем Redis кеш
    var makes []Make
    err := c.cache.Get(ctx, cacheKey, &makes)
    if err == nil {
        return makes, nil
    }
    
    // Запрос к API
    makes, err = c.fetchMakes(ctx)
    if err != nil {
        return nil, fmt.Errorf("fetch makes: %w", err)
    }
    
    // Сохраняем в кеш на 24 часа
    c.cache.Set(ctx, cacheKey, makes, 24*time.Hour)
    
    return makes, nil
}
```

### 1.2 Миграции базы данных

```sql
-- 001_extend_car_tables.up.sql

-- Расширение таблицы марок
ALTER TABLE car_makes 
ADD COLUMN IF NOT EXISTS external_id VARCHAR(100),
ADD COLUMN IF NOT EXISTS logo_url VARCHAR(500),
ADD COLUMN IF NOT EXISTS manufacturer_country VARCHAR(50),
ADD COLUMN IF NOT EXISTS last_sync_at TIMESTAMP,
ADD COLUMN IF NOT EXISTS metadata JSONB DEFAULT '{}';

-- Расширение таблицы моделей
ALTER TABLE car_models
ADD COLUMN IF NOT EXISTS external_id VARCHAR(100),
ADD COLUMN IF NOT EXISTS body_type VARCHAR(50),
ADD COLUMN IF NOT EXISTS segment VARCHAR(20),
ADD COLUMN IF NOT EXISTS years_range INT4RANGE,
ADD COLUMN IF NOT EXISTS metadata JSONB DEFAULT '{}';

-- Новая таблица для комплектаций
CREATE TABLE IF NOT EXISTS car_trims (
    id SERIAL PRIMARY KEY,
    model_id INTEGER NOT NULL REFERENCES car_models(id) ON DELETE CASCADE,
    generation_id INTEGER REFERENCES car_generations(id) ON DELETE SET NULL,
    year INTEGER NOT NULL,
    name VARCHAR(255) NOT NULL,
    slug VARCHAR(255) NOT NULL,
    
    -- Технические характеристики
    engine_type VARCHAR(50),
    engine_displacement DECIMAL(3,1),
    engine_cylinders INTEGER,
    power_hp INTEGER,
    power_kw INTEGER,
    torque_nm INTEGER,
    
    -- Трансмиссия
    transmission_type VARCHAR(50),
    transmission_gears INTEGER,
    drive_type VARCHAR(20),
    
    -- Топливо и экономичность
    fuel_type VARCHAR(30),
    fuel_economy_city DECIMAL(4,2),
    fuel_economy_highway DECIMAL(4,2),
    fuel_economy_combined DECIMAL(4,2),
    fuel_tank_capacity INTEGER,
    
    -- Размеры
    length_mm INTEGER,
    width_mm INTEGER,
    height_mm INTEGER,
    wheelbase_mm INTEGER,
    cargo_volume_l INTEGER,
    
    -- Производительность
    acceleration_0_100 DECIMAL(3,1),
    top_speed_kmh INTEGER,
    
    -- Метаданные
    external_id VARCHAR(100),
    is_active BOOLEAN DEFAULT true,
    metadata JSONB DEFAULT '{}',
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Индексы для производительности
CREATE INDEX idx_car_trims_model_year ON car_trims(model_id, year);
CREATE INDEX idx_car_trims_slug ON car_trims(slug);
CREATE INDEX idx_car_makes_external_id ON car_makes(external_id);
CREATE INDEX idx_car_models_external_id ON car_models(external_id);

-- Таблица для VIN декодирования
CREATE TABLE IF NOT EXISTS vin_decode_cache (
    id SERIAL PRIMARY KEY,
    vin VARCHAR(17) UNIQUE NOT NULL,
    make_id INTEGER REFERENCES car_makes(id),
    model_id INTEGER REFERENCES car_models(id),
    trim_id INTEGER REFERENCES car_trims(id),
    year INTEGER,
    decoded_data JSONB NOT NULL,
    source VARCHAR(50) NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    expires_at TIMESTAMP NOT NULL
);

CREATE INDEX idx_vin_decode_cache_vin ON vin_decode_cache(vin);
CREATE INDEX idx_vin_decode_cache_expires ON vin_decode_cache(expires_at);
```

### 1.3 Синхронизация данных

```go
// backend/internal/proj/marketplace/service/car_sync.go
package service

type CarSyncService struct {
    carAPI      carapi.Client
    storage     storage.CarStorage
    logger      *zap.Logger
    metrics     metrics.Interface
}

func (s *CarSyncService) SyncAllMakes(ctx context.Context) error {
    s.logger.Info("Starting car makes synchronization")
    
    // Получаем данные из API
    apiMakes, err := s.carAPI.GetMakes(ctx)
    if err != nil {
        return fmt.Errorf("get makes from API: %w", err)
    }
    
    // Получаем существующие марки
    dbMakes, err := s.storage.GetAllMakes(ctx)
    if err != nil {
        return fmt.Errorf("get makes from DB: %w", err)
    }
    
    // Создаем map для быстрого поиска
    dbMakesMap := make(map[string]*models.CarMake)
    for _, make := range dbMakes {
        dbMakesMap[make.ExternalID] = &make
    }
    
    // Синхронизируем
    var created, updated int
    for _, apiMake := range apiMakes {
        if existing, ok := dbMakesMap[apiMake.ID]; ok {
            // Обновляем существующую
            if s.needsUpdate(existing, apiMake) {
                err = s.storage.UpdateMake(ctx, s.mapAPIToDB(apiMake))
                if err != nil {
                    s.logger.Error("Failed to update make", 
                        zap.String("make", apiMake.Name),
                        zap.Error(err))
                    continue
                }
                updated++
            }
        } else {
            // Создаем новую
            err = s.storage.CreateMake(ctx, s.mapAPIToDB(apiMake))
            if err != nil {
                s.logger.Error("Failed to create make",
                    zap.String("make", apiMake.Name),
                    zap.Error(err))
                continue
            }
            created++
        }
    }
    
    s.logger.Info("Car makes synchronization completed",
        zap.Int("created", created),
        zap.Int("updated", updated))
    
    // Метрики
    s.metrics.RecordSync("car_makes", created, updated)
    
    return nil
}
```

## Фаза 2: UI/UX компоненты (3 недели)

### 2.1 Умный конструктор объявлений (Car Listing Wizard)

```typescript
// frontend/svetu/src/components/cars/CarListingWizard.tsx
import React, { useState } from 'react';
import { useTranslations } from 'next-intl';
import { motion, AnimatePresence } from 'framer-motion';

interface CarListingWizardProps {
  onComplete: (data: CarListingData) => void;
  initialData?: Partial<CarListingData>;
}

export function CarListingWizard({ onComplete, initialData }: CarListingWizardProps) {
  const t = useTranslations('cars.wizard');
  const [currentStep, setCurrentStep] = useState(0);
  const [formData, setFormData] = useState<CarListingData>(initialData || {});

  const steps = [
    { id: 'vin', component: VINScannerStep },
    { id: 'basic', component: BasicInfoStep },
    { id: 'specs', component: SpecificationsStep },
    { id: 'condition', component: ConditionStep },
    { id: 'photos', component: PhotoUploadStep },
    { id: 'pricing', component: PricingStep },
    { id: 'review', component: ReviewStep }
  ];

  const CurrentStepComponent = steps[currentStep].component;

  return (
    <div className="max-w-4xl mx-auto">
      {/* Прогресс бар */}
      <div className="mb-8">
        <div className="flex justify-between items-center mb-2">
          {steps.map((step, index) => (
            <div
              key={step.id}
              className={cn(
                "flex items-center",
                index <= currentStep ? "text-primary" : "text-base-300"
              )}
            >
              <div className={cn(
                "w-10 h-10 rounded-full flex items-center justify-center",
                index < currentStep ? "bg-primary text-primary-content" :
                index === currentStep ? "bg-primary text-primary-content animate-pulse" :
                "bg-base-200"
              )}>
                {index < currentStep ? <CheckIcon /> : index + 1}
              </div>
              {index < steps.length - 1 && (
                <div className={cn(
                  "w-20 h-1 mx-2",
                  index < currentStep ? "bg-primary" : "bg-base-200"
                )} />
              )}
            </div>
          ))}
        </div>
      </div>

      {/* Контент шага */}
      <AnimatePresence mode="wait">
        <motion.div
          key={currentStep}
          initial={{ opacity: 0, x: 20 }}
          animate={{ opacity: 1, x: 0 }}
          exit={{ opacity: 0, x: -20 }}
          transition={{ duration: 0.3 }}
        >
          <CurrentStepComponent
            data={formData}
            onUpdate={(updates) => setFormData({ ...formData, ...updates })}
            onNext={() => {
              if (currentStep < steps.length - 1) {
                setCurrentStep(currentStep + 1);
              } else {
                onComplete(formData);
              }
            }}
            onBack={() => setCurrentStep(Math.max(0, currentStep - 1))}
          />
        </motion.div>
      </AnimatePresence>
    </div>
  );
}
```

### 2.2 VIN сканер с AI

```typescript
// frontend/svetu/src/components/cars/VINScannerStep.tsx
import React, { useState, useRef } from 'react';
import { useTranslations } from 'next-intl';
import { Camera, Upload, Keyboard } from 'lucide-react';
import { decodeVIN } from '@/services/cars';

export function VINScannerStep({ data, onUpdate, onNext }: StepProps) {
  const t = useTranslations('cars.vin');
  const [mode, setMode] = useState<'scan' | 'upload' | 'manual'>('scan');
  const [isProcessing, setIsProcessing] = useState(false);
  const videoRef = useRef<HTMLVideoElement>(null);
  const canvasRef = useRef<HTMLCanvasElement>(null);

  const handleVINDecode = async (vin: string) => {
    setIsProcessing(true);
    try {
      const decodedData = await decodeVIN(vin);
      
      // Автозаполнение данных
      onUpdate({
        vin,
        make: decodedData.make,
        model: decodedData.model,
        year: decodedData.year,
        trim: decodedData.trim,
        engine: decodedData.engine,
        transmission: decodedData.transmission,
        driveType: decodedData.driveType,
        fuelType: decodedData.fuelType,
        // ... остальные данные
      });
      
      // Показываем результат
      toast.success(t('decoded_success'));
      onNext();
    } catch (error) {
      toast.error(t('decode_error'));
    } finally {
      setIsProcessing(false);
    }
  };

  const startCamera = async () => {
    try {
      const stream = await navigator.mediaDevices.getUserMedia({
        video: { facingMode: 'environment' }
      });
      if (videoRef.current) {
        videoRef.current.srcObject = stream;
      }
      
      // Запускаем распознавание
      startVINRecognition();
    } catch (error) {
      console.error('Camera error:', error);
      setMode('manual');
    }
  };

  const startVINRecognition = () => {
    // Используем TensorFlow.js или отправляем на сервер
    // для распознавания VIN из видеопотока
  };

  return (
    <div className="space-y-6">
      <div className="text-center">
        <h2 className="text-2xl font-bold mb-2">{t('title')}</h2>
        <p className="text-base-content/70">{t('description')}</p>
      </div>

      {/* Выбор способа ввода */}
      <div className="flex gap-4 justify-center">
        <button
          onClick={() => setMode('scan')}
          className={cn(
            "btn",
            mode === 'scan' ? "btn-primary" : "btn-outline"
          )}
        >
          <Camera className="w-5 h-5 mr-2" />
          {t('scan_vin')}
        </button>
        <button
          onClick={() => setMode('upload')}
          className={cn(
            "btn",
            mode === 'upload' ? "btn-primary" : "btn-outline"
          )}
        >
          <Upload className="w-5 h-5 mr-2" />
          {t('upload_photo')}
        </button>
        <button
          onClick={() => setMode('manual')}
          className={cn(
            "btn",
            mode === 'manual' ? "btn-primary" : "btn-outline"
          )}
        >
          <Keyboard className="w-5 h-5 mr-2" />
          {t('enter_manually')}
        </button>
      </div>

      {/* Контент в зависимости от режима */}
      {mode === 'scan' && (
        <div className="relative">
          <video
            ref={videoRef}
            autoPlay
            playsInline
            className="w-full rounded-lg"
          />
          <canvas
            ref={canvasRef}
            className="hidden"
          />
          <div className="absolute inset-0 flex items-center justify-center">
            <div className="border-2 border-primary border-dashed w-80 h-20 rounded-lg">
              <p className="text-center mt-24 text-primary font-semibold">
                {t('position_vin_here')}
              </p>
            </div>
          </div>
        </div>
      )}

      {mode === 'manual' && (
        <div className="form-control">
          <label className="label">
            <span className="label-text">{t('vin_number')}</span>
          </label>
          <input
            type="text"
            value={data.vin || ''}
            onChange={(e) => {
              const vin = e.target.value.toUpperCase();
              onUpdate({ vin });
              if (vin.length === 17) {
                handleVINDecode(vin);
              }
            }}
            placeholder="WVWZZZ1JZ3W386752"
            className="input input-bordered"
            maxLength={17}
          />
          <label className="label">
            <span className="label-text-alt">{t('vin_help')}</span>
          </label>
        </div>
      )}

      {/* Кнопки навигации */}
      <div className="flex justify-between">
        <button
          onClick={() => onNext()}
          className="btn btn-ghost"
        >
          {t('skip')}
        </button>
        <button
          onClick={onNext}
          disabled={!data.vin || data.vin.length !== 17}
          className="btn btn-primary"
        >
          {t('continue')}
        </button>
      </div>
    </div>
  );
}
```

### 2.3 Расширенный поиск автомобилей

```typescript
// frontend/svetu/src/components/cars/AdvancedCarSearch.tsx
import React, { useState, useEffect } from 'react';
import { useTranslations } from 'next-intl';
import { Range } from 'react-range';
import { searchCars } from '@/services/cars';

interface SearchFilters {
  make?: string[];
  model?: string[];
  yearFrom?: number;
  yearTo?: number;
  priceFrom?: number;
  priceTo?: number;
  mileageMax?: number;
  fuelType?: string[];
  transmission?: string[];
  bodyType?: string[];
  engineSizeFrom?: number;
  engineSizeTo?: number;
  powerFrom?: number;
  condition?: string[];
  features?: string[];
  location?: {
    lat: number;
    lng: number;
    radius: number;
  };
}

export function AdvancedCarSearch() {
  const t = useTranslations('cars.search');
  const [filters, setFilters] = useState<SearchFilters>({});
  const [savedSearches, setSavedSearches] = useState([]);
  const [isLoading, setIsLoading] = useState(false);

  // Популярные фильтры для быстрого доступа
  const quickFilters = [
    { id: 'new', label: t('new_cars'), filters: { condition: ['new'] } },
    { id: 'suv', label: t('suv'), filters: { bodyType: ['suv'] } },
    { id: 'electric', label: t('electric'), filters: { fuelType: ['electric'] } },
    { id: 'automatic', label: t('automatic'), filters: { transmission: ['automatic'] } },
    { id: 'diesel', label: t('diesel'), filters: { fuelType: ['diesel'] } },
  ];

  const handleSearch = async () => {
    setIsLoading(true);
    try {
      const results = await searchCars(filters);
      // Обработка результатов
    } finally {
      setIsLoading(false);
    }
  };

  return (
    <div className="drawer drawer-mobile">
      <input id="search-drawer" type="checkbox" className="drawer-toggle" />
      
      <div className="drawer-content">
        {/* Результаты поиска */}
        <div className="p-4">
          <div className="flex justify-between items-center mb-4">
            <h2 className="text-2xl font-bold">
              {t('results_count', { count: 156 })}
            </h2>
            <div className="flex gap-2">
              <select className="select select-bordered">
                <option>{t('sort.relevance')}</option>
                <option>{t('sort.price_low')}</option>
                <option>{t('sort.price_high')}</option>
                <option>{t('sort.year_new')}</option>
                <option>{t('sort.mileage_low')}</option>
              </select>
              <button className="btn btn-outline">
                {t('save_search')}
              </button>
            </div>
          </div>
          
          {/* Быстрые фильтры */}
          <div className="flex gap-2 mb-4 flex-wrap">
            {quickFilters.map(filter => (
              <button
                key={filter.id}
                onClick={() => setFilters({ ...filters, ...filter.filters })}
                className="btn btn-sm btn-outline"
              >
                {filter.label}
              </button>
            ))}
          </div>
          
          {/* Сетка результатов */}
          <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-4">
            {/* Карточки автомобилей */}
          </div>
        </div>
      </div>
      
      <div className="drawer-side">
        <label htmlFor="search-drawer" className="drawer-overlay"></label>
        <aside className="w-80 min-h-full bg-base-200">
          <div className="p-4">
            <h3 className="text-lg font-bold mb-4">{t('filters')}</h3>
            
            {/* Марка и модель */}
            <div className="form-control mb-4">
              <label className="label">
                <span className="label-text">{t('make')}</span>
              </label>
              <CarMakeMultiSelect
                value={filters.make || []}
                onChange={(make) => setFilters({ ...filters, make })}
              />
            </div>
            
            {filters.make?.length > 0 && (
              <div className="form-control mb-4">
                <label className="label">
                  <span className="label-text">{t('model')}</span>
                </label>
                <CarModelMultiSelect
                  makes={filters.make}
                  value={filters.model || []}
                  onChange={(model) => setFilters({ ...filters, model })}
                />
              </div>
            )}
            
            {/* Год выпуска */}
            <div className="form-control mb-4">
              <label className="label">
                <span className="label-text">{t('year')}</span>
              </label>
              <YearRangeSlider
                min={1990}
                max={new Date().getFullYear() + 1}
                value={[filters.yearFrom || 2000, filters.yearTo || 2024]}
                onChange={([yearFrom, yearTo]) => 
                  setFilters({ ...filters, yearFrom, yearTo })
                }
              />
            </div>
            
            {/* Цена */}
            <div className="form-control mb-4">
              <label className="label">
                <span className="label-text">{t('price')}</span>
              </label>
              <PriceRangeInput
                value={[filters.priceFrom || 0, filters.priceTo || 100000]}
                onChange={([priceFrom, priceTo]) =>
                  setFilters({ ...filters, priceFrom, priceTo })
                }
              />
            </div>
            
            {/* Пробег */}
            <div className="form-control mb-4">
              <label className="label">
                <span className="label-text">{t('mileage')}</span>
              </label>
              <input
                type="number"
                value={filters.mileageMax || ''}
                onChange={(e) => setFilters({
                  ...filters,
                  mileageMax: parseInt(e.target.value)
                })}
                placeholder={t('max_mileage')}
                className="input input-bordered"
              />
            </div>
            
            {/* Расширенные фильтры */}
            <div className="collapse collapse-arrow mb-4">
              <input type="checkbox" />
              <div className="collapse-title font-medium">
                {t('advanced_filters')}
              </div>
              <div className="collapse-content">
                {/* Тип топлива */}
                <div className="form-control mb-4">
                  <label className="label">
                    <span className="label-text">{t('fuel_type')}</span>
                  </label>
                  {['petrol', 'diesel', 'electric', 'hybrid', 'lpg'].map(fuel => (
                    <label key={fuel} className="label cursor-pointer">
                      <span className="label-text">{t(`fuel.${fuel}`)}</span>
                      <input
                        type="checkbox"
                        checked={filters.fuelType?.includes(fuel)}
                        onChange={(e) => {
                          const newFuelTypes = e.target.checked
                            ? [...(filters.fuelType || []), fuel]
                            : filters.fuelType?.filter(f => f !== fuel) || [];
                          setFilters({ ...filters, fuelType: newFuelTypes });
                        }}
                        className="checkbox"
                      />
                    </label>
                  ))}
                </div>
                
                {/* Коробка передач */}
                <div className="form-control mb-4">
                  <label className="label">
                    <span className="label-text">{t('transmission')}</span>
                  </label>
                  {['manual', 'automatic', 'semi-automatic'].map(trans => (
                    <label key={trans} className="label cursor-pointer">
                      <span className="label-text">{t(`transmission.${trans}`)}</span>
                      <input
                        type="checkbox"
                        checked={filters.transmission?.includes(trans)}
                        onChange={(e) => {
                          const newTrans = e.target.checked
                            ? [...(filters.transmission || []), trans]
                            : filters.transmission?.filter(t => t !== trans) || [];
                          setFilters({ ...filters, transmission: newTrans });
                        }}
                        className="checkbox"
                      />
                    </label>
                  ))}
                </div>
              </div>
            </div>
            
            {/* Кнопки действий */}
            <div className="flex gap-2">
              <button
                onClick={() => setFilters({})}
                className="btn btn-ghost flex-1"
              >
                {t('reset')}
              </button>
              <button
                onClick={handleSearch}
                className="btn btn-primary flex-1"
                disabled={isLoading}
              >
                {isLoading && <span className="loading loading-spinner"></span>}
                {t('search')}
              </button>
            </div>
          </div>
        </aside>
      </div>
    </div>
  );
}
```

## Фаза 3: Интеграция с системой атрибутов (2 недели)

### 3.1 Расширение атрибутов для автомобилей

```sql
-- 002_car_category_attributes.up.sql

-- Создаем атрибуты специфичные для автомобилей
INSERT INTO category_attributes (name, display_name, attribute_type, options, is_searchable, is_filterable, is_variant_compatible) VALUES
-- Основные
('car_make', 'Марка', 'select', '{"source": "car_makes"}', true, true, false),
('car_model', 'Модель', 'select', '{"source": "car_models", "depends_on": "car_make"}', true, true, false),
('car_generation', 'Поколение', 'select', '{"source": "car_generations", "depends_on": "car_model"}', false, true, false),
('car_trim', 'Комплектация', 'select', '{"source": "car_trims", "depends_on": "car_model"}', false, true, true),
('car_year', 'Год выпуска', 'number', '{"min": 1900, "max": 2025}', true, true, false),
('car_mileage', 'Пробег (км)', 'number', '{"min": 0, "max": 1000000, "step": 1000}', true, true, false),
('car_vin', 'VIN номер', 'text', '{"pattern": "^[A-HJ-NPR-Z0-9]{17}$"}', false, false, false),

-- Двигатель
('car_engine_type', 'Тип двигателя', 'select', '{"options": ["petrol", "diesel", "electric", "hybrid", "lpg"]}', true, true, false),
('car_engine_size', 'Объем двигателя (л)', 'decimal', '{"min": 0.5, "max": 8.0, "step": 0.1}', false, true, false),
('car_engine_power_hp', 'Мощность (л.с.)', 'number', '{"min": 50, "max": 1500}', false, true, false),
('car_engine_power_kw', 'Мощность (кВт)', 'number', '{"min": 35, "max": 1100}', false, true, false),

-- Трансмиссия
('car_transmission', 'Коробка передач', 'select', '{"options": ["manual", "automatic", "semi-automatic", "cvt"]}', true, true, false),
('car_drive_type', 'Привод', 'select', '{"options": ["fwd", "rwd", "awd", "4wd"]}', false, true, false),

-- Кузов
('car_body_type', 'Тип кузова', 'select', '{"options": ["sedan", "hatchback", "suv", "wagon", "coupe", "convertible", "minivan", "pickup"]}', true, true, false),
('car_doors', 'Количество дверей', 'select', '{"options": ["2", "3", "4", "5"]}', false, true, false),
('car_seats', 'Количество мест', 'number', '{"min": 2, "max": 9}', false, true, false),
('car_color', 'Цвет', 'select', '{"options": ["white", "black", "gray", "silver", "red", "blue", "green", "yellow", "brown", "other"]}', false, true, true),

-- Состояние
('car_condition', 'Состояние', 'select', '{"options": ["new", "used", "damaged", "for_parts"]}', true, true, false),
('car_owners_count', 'Количество владельцев', 'number', '{"min": 0, "max": 10}', false, true, false),
('car_service_book', 'Сервисная книжка', 'boolean', '{}', false, true, false),

-- Особенности
('car_features', 'Особенности', 'multiselect', '{"options": ["abs", "esp", "airbags", "cruise_control", "parking_sensors", "backup_camera", "navigation", "leather_seats", "sunroof", "xenon", "led", "alloy_wheels"]}', false, true, false);

-- Связываем атрибуты с автомобильными категориями
INSERT INTO category_attribute_mapping (category_id, attribute_id, is_required, sort_order)
SELECT 
    c.id as category_id,
    a.id as attribute_id,
    CASE 
        WHEN a.name IN ('car_make', 'car_model', 'car_year', 'car_mileage', 'car_condition') THEN true
        ELSE false
    END as is_required,
    CASE
        WHEN a.name = 'car_make' THEN 1
        WHEN a.name = 'car_model' THEN 2
        WHEN a.name = 'car_year' THEN 3
        WHEN a.name = 'car_mileage' THEN 4
        WHEN a.name = 'car_condition' THEN 5
        ELSE 10 + ROW_NUMBER() OVER (ORDER BY a.id)
    END as sort_order
FROM marketplace_categories c
CROSS JOIN category_attributes a
WHERE c.slug IN ('cars', 'automobiles', 'licni-automobili')
AND a.name LIKE 'car_%';
```

### 3.2 Интеграция с OpenSearch

```go
// backend/internal/storage/opensearch/car_mapping.go
package opensearch

func GetCarMapping() map[string]interface{} {
    return map[string]interface{}{
        "properties": map[string]interface{}{
            // Основные поля
            "make": map[string]interface{}{
                "type": "keyword",
                "fields": map[string]interface{}{
                    "text": map[string]interface{}{
                        "type": "text",
                        "analyzer": "standard",
                    },
                },
            },
            "model": map[string]interface{}{
                "type": "keyword",
                "fields": map[string]interface{}{
                    "text": map[string]interface{}{
                        "type": "text",
                        "analyzer": "standard",
                    },
                },
            },
            "year": map[string]interface{}{
                "type": "integer",
            },
            "mileage": map[string]interface{}{
                "type": "integer",
            },
            "price": map[string]interface{}{
                "type": "scaled_float",
                "scaling_factor": 100,
            },
            
            // Технические характеристики
            "engine": map[string]interface{}{
                "properties": map[string]interface{}{
                    "type": map[string]interface{}{
                        "type": "keyword",
                    },
                    "size": map[string]interface{}{
                        "type": "float",
                    },
                    "power_hp": map[string]interface{}{
                        "type": "integer",
                    },
                    "power_kw": map[string]interface{}{
                        "type": "integer",
                    },
                },
            },
            
            // Геолокация для поиска
            "location": map[string]interface{}{
                "type": "geo_point",
            },
            
            // Вложенные атрибуты
            "attributes": map[string]interface{}{
                "type": "nested",
                "properties": map[string]interface{}{
                    "name": map[string]interface{}{
                        "type": "keyword",
                    },
                    "value": map[string]interface{}{
                        "type": "keyword",
                        "fields": map[string]interface{}{
                            "text": map[string]interface{}{
                                "type": "text",
                            },
                            "numeric": map[string]interface{}{
                                "type": "double",
                                "ignore_malformed": true,
                            },
                        },
                    },
                },
            },
            
            // Для фасетного поиска
            "facets": map[string]interface{}{
                "properties": map[string]interface{}{
                    "body_type": map[string]interface{}{
                        "type": "keyword",
                    },
                    "fuel_type": map[string]interface{}{
                        "type": "keyword",
                    },
                    "transmission": map[string]interface{}{
                        "type": "keyword",
                    },
                    "drive_type": map[string]interface{}{
                        "type": "keyword",
                    },
                    "color": map[string]interface{}{
                        "type": "keyword",
                    },
                    "features": map[string]interface{}{
                        "type": "keyword",
                    },
                },
            },
        },
    }
}
```

### 3.3 Поисковые запросы

```go
// backend/internal/storage/opensearch/car_search.go
package opensearch

func BuildCarSearchQuery(params CarSearchParams) map[string]interface{} {
    must := []map[string]interface{}{}
    filter := []map[string]interface{}{}
    
    // Текстовый поиск
    if params.Query != "" {
        must = append(must, map[string]interface{}{
            "multi_match": map[string]interface{}{
                "query": params.Query,
                "fields": []string{
                    "title^3",
                    "make.text^2",
                    "model.text^2",
                    "description",
                },
                "type": "best_fields",
                "fuzziness": "AUTO",
            },
        })
    }
    
    // Фильтры
    if len(params.Makes) > 0 {
        filter = append(filter, map[string]interface{}{
            "terms": map[string]interface{}{
                "make": params.Makes,
            },
        })
    }
    
    if len(params.Models) > 0 {
        filter = append(filter, map[string]interface{}{
            "terms": map[string]interface{}{
                "model": params.Models,
            },
        })
    }
    
    // Диапазон года
    if params.YearFrom > 0 || params.YearTo > 0 {
        yearRange := map[string]interface{}{}
        if params.YearFrom > 0 {
            yearRange["gte"] = params.YearFrom
        }
        if params.YearTo > 0 {
            yearRange["lte"] = params.YearTo
        }
        filter = append(filter, map[string]interface{}{
            "range": map[string]interface{}{
                "year": yearRange,
            },
        })
    }
    
    // Диапазон цены
    if params.PriceFrom > 0 || params.PriceTo > 0 {
        priceRange := map[string]interface{}{}
        if params.PriceFrom > 0 {
            priceRange["gte"] = params.PriceFrom
        }
        if params.PriceTo > 0 {
            priceRange["lte"] = params.PriceTo
        }
        filter = append(filter, map[string]interface{}{
            "range": map[string]interface{}{
                "price": priceRange,
            },
        })
    }
    
    // Максимальный пробег
    if params.MileageMax > 0 {
        filter = append(filter, map[string]interface{}{
            "range": map[string]interface{}{
                "mileage": map[string]interface{}{
                    "lte": params.MileageMax,
                },
            },
        })
    }
    
    // Геолокация
    if params.Location != nil && params.RadiusKm > 0 {
        filter = append(filter, map[string]interface{}{
            "geo_distance": map[string]interface{}{
                "distance": fmt.Sprintf("%dkm", params.RadiusKm),
                "location": map[string]interface{}{
                    "lat": params.Location.Lat,
                    "lon": params.Location.Lon,
                },
            },
        })
    }
    
    // Nested атрибуты
    for _, attr := range params.Attributes {
        filter = append(filter, map[string]interface{}{
            "nested": map[string]interface{}{
                "path": "attributes",
                "query": map[string]interface{}{
                    "bool": map[string]interface{}{
                        "must": []map[string]interface{}{
                            {
                                "term": map[string]interface{}{
                                    "attributes.name": attr.Name,
                                },
                            },
                            {
                                "terms": map[string]interface{}{
                                    "attributes.value": attr.Values,
                                },
                            },
                        },
                    },
                },
            },
        })
    }
    
    // Построение финального запроса
    query := map[string]interface{}{
        "bool": map[string]interface{}{
            "must": must,
            "filter": filter,
        },
    }
    
    // Агрегации для фасетов
    aggs := map[string]interface{}{
        "makes": map[string]interface{}{
            "terms": map[string]interface{}{
                "field": "make",
                "size": 50,
            },
        },
        "models": map[string]interface{}{
            "terms": map[string]interface{}{
                "field": "model",
                "size": 100,
            },
        },
        "price_ranges": map[string]interface{}{
            "range": map[string]interface{}{
                "field": "price",
                "ranges": []map[string]interface{}{
                    {"to": 5000},
                    {"from": 5000, "to": 10000},
                    {"from": 10000, "to": 20000},
                    {"from": 20000, "to": 50000},
                    {"from": 50000},
                },
            },
        },
        "year_histogram": map[string]interface{}{
            "histogram": map[string]interface{}{
                "field": "year",
                "interval": 5,
                "min_doc_count": 1,
            },
        },
    }
    
    return map[string]interface{}{
        "query": query,
        "aggs": aggs,
        "sort": buildSortCriteria(params.Sort),
        "from": params.Offset,
        "size": params.Limit,
    }
}
```

## Фаза 4: AI и автоматизация (3 недели)

### 4.1 AI для анализа фотографий

```go
// backend/internal/integrations/vision/car_analyzer.go
package vision

import (
    "context"
    "cloud.google.com/go/vision/apiv1"
)

type CarPhotoAnalyzer struct {
    client *vision.ImageAnnotatorClient
    storage storage.Interface
}

func (a *CarPhotoAnalyzer) AnalyzeCarPhoto(ctx context.Context, imageURL string) (*CarPhotoAnalysis, error) {
    image := vision.NewImageFromURI(imageURL)
    
    // Детекция объектов
    objects, err := a.client.LocalizeObjects(ctx, image, nil)
    if err != nil {
        return nil, fmt.Errorf("localize objects: %w", err)
    }
    
    // Анализ текста (номера, VIN)
    texts, err := a.client.DetectTexts(ctx, image, nil, 10)
    if err != nil {
        return nil, fmt.Errorf("detect texts: %w", err)
    }
    
    // Детекция логотипов (марка авто)
    logos, err := a.client.DetectLogos(ctx, image, nil, 10)
    if err != nil {
        return nil, fmt.Errorf("detect logos: %w", err)
    }
    
    // Анализ качества изображения
    properties, err := a.client.DetectImageProperties(ctx, image, nil)
    if err != nil {
        return nil, fmt.Errorf("detect properties: %w", err)
    }
    
    analysis := &CarPhotoAnalysis{
        ImageURL: imageURL,
        Objects: a.extractCarObjects(objects),
        DetectedText: a.extractRelevantText(texts),
        DetectedLogos: a.extractCarLogos(logos),
        Quality: a.assessImageQuality(properties),
        Suggestions: []string{},
    }
    
    // Генерация рекомендаций
    analysis.Suggestions = a.generateSuggestions(analysis)
    
    return analysis, nil
}

func (a *CarPhotoAnalyzer) extractCarObjects(annotations *pb.LocalizeObjectsResponse) []DetectedObject {
    objects := []DetectedObject{}
    
    carRelatedNames := map[string]bool{
        "Car": true, "Vehicle": true, "Wheel": true,
        "Tire": true, "Window": true, "Door": true,
        "Bumper": true, "Mirror": true, "Headlight": true,
    }
    
    for _, obj := range annotations.LocalizedObjectAnnotations {
        if carRelatedNames[obj.Name] {
            objects = append(objects, DetectedObject{
                Name: obj.Name,
                Score: obj.Score,
                BoundingBox: obj.BoundingPoly,
            })
        }
    }
    
    return objects
}

func (a *CarPhotoAnalyzer) generateSuggestions(analysis *CarPhotoAnalysis) []string {
    suggestions := []string{}
    
    // Проверка качества
    if analysis.Quality.Score < 0.7 {
        suggestions = append(suggestions, "Улучшите освещение для более четкого изображения")
    }
    
    // Проверка углов съемки
    if !a.hasRequiredAngles(analysis.Objects) {
        suggestions = append(suggestions, "Добавьте фото с разных углов (спереди, сзади, сбоку)")
    }
    
    // Проверка интерьера
    if !a.hasInteriorPhotos(analysis.Objects) {
        suggestions = append(suggestions, "Добавьте фотографии салона автомобиля")
    }
    
    // Проверка деталей
    if len(analysis.Objects) < 5 {
        suggestions = append(suggestions, "Сделайте крупные планы важных деталей (колеса, фары, решетка)")
    }
    
    return suggestions
}
```

### 4.2 Умное ценообразование

```go
// backend/internal/proj/marketplace/service/car_pricing.go
package service

type CarPricingService struct {
    storage storage.Interface
    ml      MLService
}

func (s *CarPricingService) SuggestPrice(ctx context.Context, car CarDetails) (*PriceSuggestion, error) {
    // Получаем похожие объявления
    similar, err := s.storage.FindSimilarCars(ctx, SimilarCarQuery{
        Make:      car.Make,
        Model:     car.Model,
        YearRange: [2]int{car.Year - 2, car.Year + 2},
        MileageRange: [2]int{
            int(float64(car.Mileage) * 0.8),
            int(float64(car.Mileage) * 1.2),
        },
        Limit: 100,
    })
    if err != nil {
        return nil, fmt.Errorf("find similar cars: %w", err)
    }
    
    // Базовая статистика
    prices := extractPrices(similar)
    stats := calculatePriceStats(prices)
    
    // ML предсказание
    features := s.extractFeatures(car)
    mlPrice, confidence := s.ml.PredictPrice(ctx, features)
    
    // Факторы влияющие на цену
    factors := s.calculatePriceFactors(car, similar)
    
    suggestion := &PriceSuggestion{
        RecommendedPrice: s.calculateRecommendedPrice(stats, mlPrice, confidence),
        PriceRange: PriceRange{
            Min: stats.Percentile25,
            Max: stats.Percentile75,
        },
        MarketAverage: stats.Mean,
        MLPrediction: mlPrice,
        Confidence: confidence,
        Factors: factors,
        SimilarListings: len(similar),
        LastUpdated: time.Now(),
    }
    
    // Рекомендации по позиционированию
    suggestion.PositioningAdvice = s.generatePositioningAdvice(suggestion, car)
    
    return suggestion, nil
}

func (s *CarPricingService) calculatePriceFactors(car CarDetails, similar []CarListing) []PriceFactor {
    factors := []PriceFactor{}
    
    // Пробег
    avgMileage := calculateAvgMileage(similar)
    if car.Mileage < avgMileage {
        factors = append(factors, PriceFactor{
            Name: "Низкий пробег",
            Impact: "+5-10%",
            Description: fmt.Sprintf("Пробег ниже среднего на %.0f%%", 
                (1 - float64(car.Mileage)/float64(avgMileage)) * 100),
        })
    }
    
    // Комплектация
    if car.Trim != "" && isHighTrim(car.Trim) {
        factors = append(factors, PriceFactor{
            Name: "Топовая комплектация",
            Impact: "+10-15%",
            Description: "Максимальная комплектация с полным набором опций",
        })
    }
    
    // Сервисная история
    if car.ServiceBook {
        factors = append(factors, PriceFactor{
            Name: "Сервисная книжка",
            Impact: "+3-5%",
            Description: "Полная сервисная история",
        })
    }
    
    // Сезонность
    season := getCurrentSeason()
    if season == "spring" && car.BodyType == "convertible" {
        factors = append(factors, PriceFactor{
            Name: "Сезонный спрос",
            Impact: "+5-7%",
            Description: "Повышенный спрос на кабриолеты весной",
        })
    }
    
    return factors
}
```

### 4.3 Чат-бот помощник

```typescript
// frontend/svetu/src/components/cars/CarAssistantChat.tsx
import React, { useState, useRef, useEffect } from 'react';
import { useTranslations } from 'next-intl';
import { sendMessage } from '@/services/chat';

interface Message {
  id: string;
  role: 'user' | 'assistant';
  content: string;
  timestamp: Date;
  actions?: ChatAction[];
}

interface ChatAction {
  type: 'search' | 'filter' | 'compare' | 'schedule';
  label: string;
  data: any;
}

export function CarAssistantChat() {
  const t = useTranslations('cars.assistant');
  const [messages, setMessages] = useState<Message[]>([
    {
      id: '1',
      role: 'assistant',
      content: t('welcome'),
      timestamp: new Date(),
      actions: [
        { type: 'search', label: t('find_car'), data: {} },
        { type: 'compare', label: t('compare_cars'), data: {} },
      ]
    }
  ]);
  const [input, setInput] = useState('');
  const [isTyping, setIsTyping] = useState(false);
  const messagesEndRef = useRef<HTMLDivElement>(null);

  const scrollToBottom = () => {
    messagesEndRef.current?.scrollIntoView({ behavior: 'smooth' });
  };

  useEffect(() => {
    scrollToBottom();
  }, [messages]);

  const handleSend = async () => {
    if (!input.trim()) return;

    const userMessage: Message = {
      id: Date.now().toString(),
      role: 'user',
      content: input,
      timestamp: new Date(),
    };

    setMessages(prev => [...prev, userMessage]);
    setInput('');
    setIsTyping(true);

    try {
      const response = await sendMessage({
        message: input,
        context: 'car_search',
        history: messages,
      });

      const assistantMessage: Message = {
        id: Date.now().toString(),
        role: 'assistant',
        content: response.content,
        timestamp: new Date(),
        actions: response.actions,
      };

      setMessages(prev => [...prev, assistantMessage]);
    } catch (error) {
      console.error('Chat error:', error);
    } finally {
      setIsTyping(false);
    }
  };

  const quickQuestions = [
    t('quick.budget_car'),
    t('quick.family_suv'),
    t('quick.fuel_efficient'),
    t('quick.first_car'),
  ];

  return (
    <div className="flex flex-col h-[600px] bg-base-100 rounded-lg shadow-lg">
      {/* Header */}
      <div className="p-4 bg-primary text-primary-content">
        <h3 className="text-lg font-bold">{t('title')}</h3>
        <p className="text-sm opacity-90">{t('subtitle')}</p>
      </div>

      {/* Messages */}
      <div className="flex-1 overflow-y-auto p-4 space-y-4">
        {messages.map((message) => (
          <div
            key={message.id}
            className={cn(
              "flex",
              message.role === 'user' ? "justify-end" : "justify-start"
            )}
          >
            <div
              className={cn(
                "max-w-[80%] rounded-lg p-3",
                message.role === 'user'
                  ? "bg-primary text-primary-content"
                  : "bg-base-200"
              )}
            >
              <p className="whitespace-pre-wrap">{message.content}</p>
              
              {message.actions && (
                <div className="mt-2 space-y-1">
                  {message.actions.map((action, idx) => (
                    <button
                      key={idx}
                      onClick={() => handleAction(action)}
                      className="btn btn-sm btn-ghost w-full justify-start"
                    >
                      {action.label}
                    </button>
                  ))}
                </div>
              )}
              
              <p className="text-xs opacity-70 mt-1">
                {message.timestamp.toLocaleTimeString()}
              </p>
            </div>
          </div>
        ))}
        
        {isTyping && (
          <div className="flex justify-start">
            <div className="bg-base-200 rounded-lg p-3">
              <span className="loading loading-dots loading-sm"></span>
            </div>
          </div>
        )}
        
        <div ref={messagesEndRef} />
      </div>

      {/* Quick questions */}
      {messages.length === 1 && (
        <div className="p-4 border-t">
          <p className="text-sm text-base-content/70 mb-2">
            {t('try_asking')}
          </p>
          <div className="flex flex-wrap gap-2">
            {quickQuestions.map((question) => (
              <button
                key={question}
                onClick={() => setInput(question)}
                className="btn btn-sm btn-outline"
              >
                {question}
              </button>
            ))}
          </div>
        </div>
      )}

      {/* Input */}
      <div className="p-4 border-t">
        <form
          onSubmit={(e) => {
            e.preventDefault();
            handleSend();
          }}
          className="flex gap-2"
        >
          <input
            type="text"
            value={input}
            onChange={(e) => setInput(e.target.value)}
            placeholder={t('type_message')}
            className="input input-bordered flex-1"
            disabled={isTyping}
          />
          <button
            type="submit"
            className="btn btn-primary"
            disabled={!input.trim() || isTyping}
          >
            {t('send')}
          </button>
        </form>
      </div>
    </div>
  );
}
```

## Фаза 5: Мобильная оптимизация (2 недели)

### 5.1 PWA функциональность

```typescript
// frontend/svetu/src/app/manifest.ts
import { MetadataRoute } from 'next';

export default function manifest(): MetadataRoute.Manifest {
  return {
    name: 'SveTu Auto',
    short_name: 'SveTu Auto',
    description: 'Лучший автомобильный маркетплейс в регионе',
    start_url: '/cars',
    display: 'standalone',
    background_color: '#ffffff',
    theme_color: '#1d4ed8',
    icons: [
      {
        src: '/icons/icon-192x192.png',
        sizes: '192x192',
        type: 'image/png',
      },
      {
        src: '/icons/icon-512x512.png',
        sizes: '512x512',
        type: 'image/png',
      },
    ],
    shortcuts: [
      {
        name: 'Поиск авто',
        short_name: 'Поиск',
        description: 'Найти автомобиль',
        url: '/cars/search',
        icons: [{ src: '/icons/search.png', sizes: '192x192' }],
      },
      {
        name: 'Создать объявление',
        short_name: 'Продать',
        description: 'Продать автомобиль',
        url: '/cars/sell',
        icons: [{ src: '/icons/sell.png', sizes: '192x192' }],
      },
    ],
  };
}
```

### 5.2 Оффлайн функциональность

```typescript
// frontend/svetu/public/service-worker.js
const CACHE_NAME = 'svetu-auto-v1';
const urlsToCache = [
  '/',
  '/cars',
  '/cars/search',
  '/offline.html',
  '/icons/icon-192x192.png',
  '/icons/icon-512x512.png',
];

// Кеширование при установке
self.addEventListener('install', (event) => {
  event.waitUntil(
    caches.open(CACHE_NAME).then((cache) => {
      return cache.addAll(urlsToCache);
    })
  );
});

// Стратегия Network First с fallback на кеш
self.addEventListener('fetch', (event) => {
  if (event.request.url.includes('/api/')) {
    // API запросы - network first
    event.respondWith(
      fetch(event.request)
        .then((response) => {
          // Кешируем успешные GET запросы
          if (response.ok && event.request.method === 'GET') {
            const responseToCache = response.clone();
            caches.open(CACHE_NAME).then((cache) => {
              cache.put(event.request, responseToCache);
            });
          }
          return response;
        })
        .catch(() => {
          // Fallback на кеш для API
          return caches.match(event.request);
        })
    );
  } else {
    // Статика - cache first
    event.respondWith(
      caches.match(event.request).then((response) => {
        return response || fetch(event.request);
      })
    );
  }
});

// Фоновая синхронизация
self.addEventListener('sync', (event) => {
  if (event.tag === 'sync-favorites') {
    event.waitUntil(syncFavorites());
  }
});

async function syncFavorites() {
  const db = await openDB();
  const favorites = await db.getAll('pending-favorites');
  
  for (const favorite of favorites) {
    try {
      await fetch('/api/v1/favorites', {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify(favorite),
      });
      await db.delete('pending-favorites', favorite.id);
    } catch (error) {
      console.error('Sync failed:', error);
    }
  }
}
```

## Фаза 6: Аналитика и мониторинг (1 неделя)

### 6.1 Метрики и дашборды

```go
// backend/internal/metrics/car_metrics.go
package metrics

import (
    "github.com/prometheus/client_golang/prometheus"
    "github.com/prometheus/client_golang/prometheus/promauto"
)

var (
    // Поиск
    carSearchTotal = promauto.NewCounterVec(
        prometheus.CounterOpts{
            Name: "car_search_total",
            Help: "Total number of car searches",
        },
        []string{"make", "model", "year_range"},
    )
    
    carSearchDuration = promauto.NewHistogramVec(
        prometheus.HistogramOpts{
            Name: "car_search_duration_seconds",
            Help: "Duration of car search requests",
            Buckets: prometheus.DefBuckets,
        },
        []string{"status"},
    )
    
    // Создание объявлений
    carListingCreated = promauto.NewCounterVec(
        prometheus.CounterOpts{
            Name: "car_listing_created_total",
            Help: "Total number of car listings created",
        },
        []string{"make", "model", "source"},
    )
    
    carListingCompletionTime = promauto.NewHistogram(
        prometheus.HistogramOpts{
            Name: "car_listing_completion_seconds",
            Help: "Time to complete car listing creation",
            Buckets: []float64{30, 60, 120, 180, 300, 600},
        },
    )
    
    // VIN декодирование
    vinDecodeTotal = promauto.NewCounterVec(
        prometheus.CounterOpts{
            Name: "vin_decode_total",
            Help: "Total VIN decode requests",
        },
        []string{"status", "source"},
    )
    
    // API использование
    carAPIUsage = promauto.NewCounterVec(
        prometheus.CounterOpts{
            Name: "car_api_usage_total",
            Help: "External car API usage",
        },
        []string{"endpoint", "status"},
    )
    
    carAPICacheHit = promauto.NewCounterVec(
        prometheus.CounterOpts{
            Name: "car_api_cache_hit_total",
            Help: "Car API cache hits",
        },
        []string{"endpoint"},
    )
)

func RecordCarSearch(make, model, yearRange string, duration float64, err error) {
    status := "success"
    if err != nil {
        status = "error"
    }
    
    carSearchTotal.WithLabelValues(make, model, yearRange).Inc()
    carSearchDuration.WithLabelValues(status).Observe(duration)
}

func RecordCarListingCreated(make, model, source string, duration float64) {
    carListingCreated.WithLabelValues(make, model, source).Inc()
    carListingCompletionTime.Observe(duration)
}
```

### 6.2 Grafana дашборд

```json
{
  "dashboard": {
    "title": "Автомобильный маркетплейс",
    "panels": [
      {
        "title": "Популярные марки (поиск)",
        "targets": [{
          "expr": "sum by (make) (rate(car_search_total[5m]))"
        }],
        "type": "piechart"
      },
      {
        "title": "Время создания объявления",
        "targets": [{
          "expr": "histogram_quantile(0.95, car_listing_completion_seconds)"
        }],
        "type": "gauge"
      },
      {
        "title": "Использование API",
        "targets": [{
          "expr": "sum(rate(car_api_usage_total[5m])) by (endpoint)"
        }],
        "type": "timeseries"
      },
      {
        "title": "Эффективность кеша",
        "targets": [{
          "expr": "sum(rate(car_api_cache_hit_total[5m])) / sum(rate(car_api_usage_total[5m])) * 100"
        }],
        "type": "stat"
      }
    ]
  }
}
```

## План миграции и развертывания

### Этап 1: Подготовка (1 неделя)
1. Настройка окружения для CarAPI.app
2. Создание и тестирование миграций БД
3. Настройка Redis для кеширования
4. Подготовка CI/CD пайплайнов

### Этап 2: Backend развертывание (1 неделя)
1. Развертывание сервиса синхронизации
2. Первичная загрузка данных из CarAPI
3. Настройка cron задач
4. Мониторинг и алерты

### Этап 3: Frontend развертывание (2 недели)
1. Развертывание новых компонентов
2. A/B тестирование
3. Постепенный rollout
4. Сбор обратной связи

### Этап 4: Оптимизация (ongoing)
1. Анализ метрик
2. Улучшение производительности
3. Расширение функциональности
4. Масштабирование

## Риски и митигация

### Технические риски:
1. **Превышение лимитов API**
   - Митигация: Агрессивное кеширование, мониторинг использования
   
2. **Производительность поиска**
   - Митигация: Оптимизация индексов OpenSearch, кеширование популярных запросов
   
3. **Качество данных**
   - Митигация: Валидация на всех уровнях, ручная модерация

### Бизнес риски:
1. **Низкая конверсия**
   - Митигация: A/B тестирование, итеративные улучшения UX
   
2. **Конкуренция**
   - Митигация: Уникальные функции (VIN сканер, AI помощник)

## Заключение

Данный план представляет комплексный подход к созданию современного автомобильного маркетплейса. Ключевые преимущества:

1. **Автоматизация** - VIN сканер, AI анализ фото, умное ценообразование
2. **Удобство** - Пошаговый wizard, чат-помощник, мобильная оптимизация  
3. **Точность** - Интеграция с CarAPI.app, актуальные данные
4. **Производительность** - Кеширование, оптимизированный поиск
5. **Масштабируемость** - Модульная архитектура, микросервисы

Реализация плана позволит создать лучший автомобильный маркетплейс в регионе с высокой конверсией и удовлетворенностью пользователей.