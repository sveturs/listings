'use client';

import React, { useState } from 'react';
import { useTranslations } from 'next-intl';
import { CarSelector, CarSelectorCompact } from '@/components/cars';
import type { CarSelection } from '@/types/cars';

export const CarSelectorDemo: React.FC = () => {
  const _t = useTranslations('cars');
  const [fullSelection, setFullSelection] = useState<CarSelection>({});
  const [compactSelection, setCompactSelection] = useState<CarSelection>({});

  return (
    <div className="p-6 space-y-8">
      <div className="max-w-2xl mx-auto">
        <h1 className="text-3xl font-bold mb-8 text-center">
          Демо каскадного селектора автомобилей
        </h1>

        {/* Полная версия */}
        <div className="card bg-base-100 shadow-xl">
          <div className="card-body">
            <h2 className="card-title">Полная версия с поколениями</h2>
            <CarSelector
              value={fullSelection}
              onChange={setFullSelection}
              required
              showGenerations
              placeholder={{
                make: 'Выберите марку автомобиля',
                model: 'Выберите модель',
                generation: 'Выберите поколение',
              }}
            />

            <div className="mt-4 p-4 bg-base-200 rounded-lg">
              <h3 className="font-semibold mb-2">Выбранные значения:</h3>
              <pre className="text-sm">
                {JSON.stringify(fullSelection, null, 2)}
              </pre>
            </div>
          </div>
        </div>

        {/* Компактная версия */}
        <div className="card bg-base-100 shadow-xl">
          <div className="card-body">
            <h2 className="card-title">
              Компактная версия (только марка и модель)
            </h2>
            <CarSelectorCompact
              value={compactSelection}
              onChange={setCompactSelection}
            />

            <div className="mt-4 p-4 bg-base-200 rounded-lg">
              <h3 className="font-semibold mb-2">Выбранные значения:</h3>
              <pre className="text-sm">
                {JSON.stringify(compactSelection, null, 2)}
              </pre>
            </div>
          </div>
        </div>

        {/* Пример использования в форме */}
        <div className="card bg-base-100 shadow-xl">
          <div className="card-body">
            <h2 className="card-title">Пример использования в форме</h2>
            <form className="space-y-4">
              <div className="form-control">
                <label className="label">
                  <span className="label-text">Название объявления</span>
                </label>
                <input
                  type="text"
                  placeholder="Введите название"
                  className="input input-bordered"
                />
              </div>

              <CarSelectorCompact
                value={compactSelection}
                onChange={setCompactSelection}
              />

              <div className="form-control">
                <label className="label">
                  <span className="label-text">Цена</span>
                </label>
                <input
                  type="number"
                  placeholder="Введите цену"
                  className="input input-bordered"
                />
              </div>

              <div className="form-control mt-6">
                <button type="submit" className="btn btn-primary">
                  Создать объявление
                </button>
              </div>
            </form>
          </div>
        </div>
      </div>
    </div>
  );
};
