'use client';

import React, { useState } from 'react';
import { DistanceBadge } from '@/components/ui/DistanceBadge';
import { DistanceVisualizer } from '@/components/ui/DistanceVisualizer';
import { DistanceIndicator } from '@/components/ui/DistanceIndicator';
import { Slider } from '@/components/ui/Slider';

export default function DistanceExamplesPage() {
  const [distance, setDistance] = useState(2.5);

  const exampleDistances = [0.3, 0.8, 1.5, 3.2, 7.5, 12, 25];
  const addresses = [
    'ул. Пушкина, д. 10',
    'пр. Ленина, д. 45, кв. 12',
    'Центральный район, ул. Мира',
    'Северный микрорайон, д. 7',
  ];

  return (
    <div className="container mx-auto p-6 max-w-6xl">
      <h1 className="text-3xl font-bold mb-8">Визуализация расстояния</h1>

      {/* Интерактивная демонстрация */}
      <section className="mb-12">
        <h2 className="text-2xl font-semibold mb-6">Интерактивная демонстрация</h2>
        
        <div className="card bg-base-200 p-6">
          <div className="mb-6">
            <label className="label">
              <span className="label-text">Расстояние: {distance.toFixed(1)} км</span>
            </label>
            <input
              type="range"
              min="0.1"
              max="30"
              step="0.1"
              value={distance}
              onChange={(e) => setDistance(parseFloat(e.target.value))}
              className="range range-primary"
            />
            <div className="w-full flex justify-between text-xs px-2 mt-1">
              <span>0.1 км</span>
              <span>5 км</span>
              <span>10 км</span>
              <span>20 км</span>
              <span>30 км</span>
            </div>
          </div>

          <div className="grid grid-cols-1 lg:grid-cols-2 gap-6">
            {/* Визуализатор */}
            <div>
              <h3 className="font-medium mb-3">DistanceVisualizer</h3>
              <DistanceVisualizer
                distance={distance}
                address="ул. Примерная, д. 123"
                showMap={true}
                showTravelTime={true}
              />
            </div>

            {/* Варианты бейджей */}
            <div className="space-y-4">
              <div>
                <h3 className="font-medium mb-3">DistanceBadge варианты</h3>
                <div className="space-y-3">
                  <div className="flex items-center gap-4">
                    <span className="text-sm text-base-content/60 w-20">Default:</span>
                    <DistanceBadge distance={distance} />
                  </div>
                  <div className="flex items-center gap-4">
                    <span className="text-sm text-base-content/60 w-20">Compact:</span>
                    <DistanceBadge distance={distance} variant="compact" />
                  </div>
                  <div className="flex items-center gap-4">
                    <span className="text-sm text-base-content/60 w-20">Detailed:</span>
                    <DistanceBadge distance={distance} variant="detailed" />
                  </div>
                </div>
              </div>

              <div>
                <h3 className="font-medium mb-3">DistanceIndicator размеры</h3>
                <div className="space-y-2">
                  <DistanceIndicator distance={distance} size="sm" />
                  <DistanceIndicator distance={distance} size="md" />
                  <DistanceIndicator distance={distance} size="lg" />
                </div>
              </div>
            </div>
          </div>
        </div>
      </section>

      {/* Примеры использования */}
      <section className="mb-12">
        <h2 className="text-2xl font-semibold mb-6">Примеры на карточках товаров</h2>
        
        <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-4">
          {exampleDistances.slice(0, 6).map((dist, index) => (
            <div key={index} className="card bg-base-100 shadow-sm">
              <figure className="h-40 bg-base-200"></figure>
              <div className="card-body p-4">
                <div className="flex justify-between items-start mb-2">
                  <h3 className="font-semibold">Товар {index + 1}</h3>
                  <DistanceBadge distance={dist} />
                </div>
                <p className="text-sm text-base-content/70 mb-3">
                  Описание товара для примера визуализации расстояния
                </p>
                <div className="flex items-center justify-between">
                  <span className="text-lg font-bold">€{(50 + index * 25).toFixed(0)}</span>
                  <DistanceIndicator distance={dist} size="sm" />
                </div>
              </div>
            </div>
          ))}
        </div>
      </section>

      {/* Цветовая схема */}
      <section className="mb-12">
        <h2 className="text-2xl font-semibold mb-6">Цветовая схема расстояний</h2>
        
        <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-4 gap-4">
          <div className="card bg-success/10 border border-success/20 p-4">
            <div className="flex items-center gap-2 mb-2">
              <div className="w-3 h-3 bg-success rounded-full"></div>
              <h3 className="font-medium text-success">Очень близко</h3>
            </div>
            <p className="text-sm">До 1 км</p>
            <div className="mt-3">
              <DistanceBadge distance={0.5} />
            </div>
          </div>

          <div className="card bg-info/10 border border-info/20 p-4">
            <div className="flex items-center gap-2 mb-2">
              <div className="w-3 h-3 bg-info rounded-full"></div>
              <h3 className="font-medium text-info">Рядом</h3>
            </div>
            <p className="text-sm">1-5 км</p>
            <div className="mt-3">
              <DistanceBadge distance={3} />
            </div>
          </div>

          <div className="card bg-warning/10 border border-warning/20 p-4">
            <div className="flex items-center gap-2 mb-2">
              <div className="w-3 h-3 bg-warning rounded-full"></div>
              <h3 className="font-medium text-warning">Недалеко</h3>
            </div>
            <p className="text-sm">5-15 км</p>
            <div className="mt-3">
              <DistanceBadge distance={10} />
            </div>
          </div>

          <div className="card bg-base-200 border border-base-300 p-4">
            <div className="flex items-center gap-2 mb-2">
              <div className="w-3 h-3 bg-base-300 rounded-full"></div>
              <h3 className="font-medium">Далеко</h3>
            </div>
            <p className="text-sm">Более 15 км</p>
            <div className="mt-3">
              <DistanceBadge distance={20} />
            </div>
          </div>
        </div>
      </section>

      {/* Детальный вид */}
      <section className="mb-12">
        <h2 className="text-2xl font-semibold mb-6">Детальный вид для страницы товара</h2>
        
        <div className="grid grid-cols-1 lg:grid-cols-2 gap-6">
          {addresses.map((address, index) => (
            <div key={index} className="card bg-base-100 shadow-md p-6">
              <h3 className="font-semibold mb-4">Местоположение товара</h3>
              <DistanceVisualizer
                distance={exampleDistances[index]}
                address={address}
                showMap={true}
                showTravelTime={true}
              />
            </div>
          ))}
        </div>
      </section>

      {/* Примеры кода */}
      <section className="card bg-base-200 p-6">
        <h2 className="text-2xl font-semibold mb-4">Примеры использования</h2>
        <div className="mockup-code">
          <pre data-prefix="1"><code>{`import { DistanceBadge, DistanceVisualizer, DistanceIndicator } from '@/components/ui';`}</code></pre>
          <pre data-prefix="2"><code>{``}</code></pre>
          <pre data-prefix="3"><code>{`// Простой бейдж`}</code></pre>
          <pre data-prefix="4"><code>{`<DistanceBadge distance={2.5} />`}</code></pre>
          <pre data-prefix="5"><code>{``}</code></pre>
          <pre data-prefix="6"><code>{`// Индикатор с тултипом`}</code></pre>
          <pre data-prefix="7"><code>{`<DistanceIndicator distance={1.2} size="sm" />`}</code></pre>
          <pre data-prefix="8"><code>{``}</code></pre>
          <pre data-prefix="9"><code>{`// Полная визуализация`}</code></pre>
          <pre data-prefix="10"><code>{`<DistanceVisualizer`}</code></pre>
          <pre data-prefix="11"><code>{`  distance={3.7}`}</code></pre>
          <pre data-prefix="12"><code>{`  address="ул. Примерная, д. 123"`}</code></pre>
          <pre data-prefix="13"><code>{`  showMap={true}`}</code></pre>
          <pre data-prefix="14"><code>{`  showTravelTime={true}`}</code></pre>
          <pre data-prefix="15"><code>{`/>`}</code></pre>
        </div>
      </section>
    </div>
  );
}