'use client';

import React, { useState, useEffect } from 'react';
import { useTranslations } from 'next-intl';
import { Car, Calendar, Gauge, Fuel, Settings } from 'lucide-react';
import { CarsService } from '@/services/cars';
import type { CarMake, CarModel } from '@/types/cars';

interface CarFiltersProps {
  onFiltersChange: (filters: Record<string, any>) => void;
  className?: string;
}

export const CarFilters: React.FC<CarFiltersProps> = ({
  onFiltersChange,
  className = '',
}) => {
  const t = useTranslations('cars.filters');

  // Состояния для фильтров
  const [selectedMake, setSelectedMake] = useState<string>('');
  const [selectedModel, setSelectedModel] = useState<string>('');
  const [yearFrom, setYearFrom] = useState<string>('');
  const [yearTo, setYearTo] = useState<string>('');
  const [priceFrom, setPriceFrom] = useState<string>('');
  const [priceTo, setPriceTo] = useState<string>('');
  const [mileageMax, setMileageMax] = useState<string>('');
  const [fuelType, setFuelType] = useState<string>('');
  const [transmission, setTransmission] = useState<string>('');
  const [condition, setCondition] = useState<string>('');

  // Данные для селекторов
  const [makes, setMakes] = useState<CarMake[]>([]);
  const [models, setModels] = useState<CarModel[]>([]);
  const [loadingMakes, setLoadingMakes] = useState(false);
  const [loadingModels, setLoadingModels] = useState(false);

  // Загрузка марок при монтировании
  useEffect(() => {
    loadMakes();
  }, []);

  // Загрузка моделей при выборе марки
  useEffect(() => {
    if (selectedMake) {
      loadModels(selectedMake);
    } else {
      setModels([]);
      setSelectedModel('');
    }
  }, [selectedMake]);

  // Обновление фильтров при изменении любого параметра
  useEffect(() => {
    const filters: Record<string, any> = {};

    if (selectedMake) filters.make = selectedMake;
    if (selectedModel) filters.model = selectedModel;
    if (yearFrom) filters.yearFrom = parseInt(yearFrom);
    if (yearTo) filters.yearTo = parseInt(yearTo);
    if (priceFrom) filters.priceMin = parseInt(priceFrom);
    if (priceTo) filters.priceMax = parseInt(priceTo);
    if (mileageMax) filters.mileageMax = parseInt(mileageMax);
    if (fuelType) filters.fuelType = fuelType;
    if (transmission) filters.transmission = transmission;
    if (condition) filters.condition = condition;

    onFiltersChange(filters);
  }, [
    selectedMake,
    selectedModel,
    yearFrom,
    yearTo,
    priceFrom,
    priceTo,
    mileageMax,
    fuelType,
    transmission,
    condition,
    onFiltersChange,
  ]);

  const loadMakes = async () => {
    setLoadingMakes(true);
    try {
      const response = await CarsService.getMakes();
      if (response.success && response.data) {
        setMakes(response.data);
      }
    } catch (error) {
      console.error('Error loading makes:', error);
    } finally {
      setLoadingMakes(false);
    }
  };

  const loadModels = async (makeSlug: string) => {
    setLoadingModels(true);
    try {
      const response = await CarsService.getModelsByMake(makeSlug);
      if (response.success && response.data) {
        setModels(response.data);
      }
    } catch (error) {
      console.error('Error loading models:', error);
    } finally {
      setLoadingModels(false);
    }
  };

  const resetFilters = () => {
    setSelectedMake('');
    setSelectedModel('');
    setYearFrom('');
    setYearTo('');
    setPriceFrom('');
    setPriceTo('');
    setMileageMax('');
    setFuelType('');
    setTransmission('');
    setCondition('');
  };

  // Генерация списка годов
  const currentYear = new Date().getFullYear();
  const years = Array.from(
    { length: currentYear - 1990 + 1 },
    (_, i) => currentYear - i
  );

  return (
    <div className={`card bg-base-100 shadow-xl ${className}`}>
      <div className="card-body">
        <div className="flex items-center justify-between mb-4">
          <h3 className="card-title text-lg">
            <Car className="w-5 h-5" />
            {t('title')}
          </h3>
          <button
            onClick={resetFilters}
            className="btn btn-ghost btn-sm"
            aria-label={t('reset')}
          >
            {t('reset')}
          </button>
        </div>

        <div className="space-y-4">
          {/* Марка и модель */}
          <div className="space-y-2">
            <label className="label">
              <span className="label-text font-medium">{t('make')}</span>
            </label>
            <select
              value={selectedMake}
              onChange={(e) => setSelectedMake(e.target.value)}
              className="select select-bordered w-full"
              disabled={loadingMakes}
            >
              <option value="">{t('allMakes')}</option>
              {makes.map((make) => (
                <option key={make.id} value={make.slug}>
                  {make.name}
                </option>
              ))}
            </select>

            {selectedMake && (
              <>
                <label className="label">
                  <span className="label-text font-medium">{t('model')}</span>
                </label>
                <select
                  value={selectedModel}
                  onChange={(e) => setSelectedModel(e.target.value)}
                  className="select select-bordered w-full"
                  disabled={loadingModels || models.length === 0}
                >
                  <option value="">{t('allModels')}</option>
                  {models.map((model) => (
                    <option key={model.id} value={model.slug}>
                      {model.name}
                    </option>
                  ))}
                </select>
              </>
            )}
          </div>

          {/* Год выпуска */}
          <div>
            <label className="label">
              <span className="label-text font-medium">
                <Calendar className="w-4 h-4 inline mr-1" />
                {t('year')}
              </span>
            </label>
            <div className="grid grid-cols-2 gap-2">
              <select
                value={yearFrom}
                onChange={(e) => setYearFrom(e.target.value)}
                className="select select-bordered select-sm"
              >
                <option value="">{t('from')}</option>
                {years.map((year) => (
                  <option key={year} value={year}>
                    {year}
                  </option>
                ))}
              </select>
              <select
                value={yearTo}
                onChange={(e) => setYearTo(e.target.value)}
                className="select select-bordered select-sm"
              >
                <option value="">{t('to')}</option>
                {years.map((year) => (
                  <option key={year} value={year}>
                    {year}
                  </option>
                ))}
              </select>
            </div>
          </div>

          {/* Цена */}
          <div>
            <label className="label">
              <span className="label-text font-medium">{t('price')}</span>
            </label>
            <div className="grid grid-cols-2 gap-2">
              <input
                type="number"
                value={priceFrom}
                onChange={(e) => setPriceFrom(e.target.value)}
                placeholder={t('from')}
                className="input input-bordered input-sm"
              />
              <input
                type="number"
                value={priceTo}
                onChange={(e) => setPriceTo(e.target.value)}
                placeholder={t('to')}
                className="input input-bordered input-sm"
              />
            </div>
          </div>

          {/* Пробег */}
          <div>
            <label className="label">
              <span className="label-text font-medium">
                <Gauge className="w-4 h-4 inline mr-1" />
                {t('mileage')}
              </span>
            </label>
            <input
              type="number"
              value={mileageMax}
              onChange={(e) => setMileageMax(e.target.value)}
              placeholder={t('maxMileage')}
              className="input input-bordered input-sm w-full"
            />
          </div>

          {/* Тип топлива */}
          <div>
            <label className="label">
              <span className="label-text font-medium">
                <Fuel className="w-4 h-4 inline mr-1" />
                {t('fuelType')}
              </span>
            </label>
            <select
              value={fuelType}
              onChange={(e) => setFuelType(e.target.value)}
              className="select select-bordered select-sm w-full"
            >
              <option value="">{t('allFuelTypes')}</option>
              <option value="petrol">{t('petrol')}</option>
              <option value="diesel">{t('diesel')}</option>
              <option value="electric">{t('electric')}</option>
              <option value="hybrid">{t('hybrid')}</option>
              <option value="lpg">{t('lpg')}</option>
            </select>
          </div>

          {/* Коробка передач */}
          <div>
            <label className="label">
              <span className="label-text font-medium">
                <Settings className="w-4 h-4 inline mr-1" />
                {t('transmission')}
              </span>
            </label>
            <select
              value={transmission}
              onChange={(e) => setTransmission(e.target.value)}
              className="select select-bordered select-sm w-full"
            >
              <option value="">{t('allTransmissions')}</option>
              <option value="manual">{t('manual')}</option>
              <option value="automatic">{t('automatic')}</option>
              <option value="semi-automatic">{t('semiAutomatic')}</option>
            </select>
          </div>

          {/* Состояние */}
          <div>
            <label className="label">
              <span className="label-text font-medium">{t('condition')}</span>
            </label>
            <select
              value={condition}
              onChange={(e) => setCondition(e.target.value)}
              className="select select-bordered select-sm w-full"
            >
              <option value="">{t('allConditions')}</option>
              <option value="new">{t('new')}</option>
              <option value="used">{t('used')}</option>
              <option value="damaged">{t('damaged')}</option>
            </select>
          </div>
        </div>
      </div>
    </div>
  );
};
