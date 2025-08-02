'use client';

import React, { useState, useEffect, useCallback } from 'react';
import { useTranslations } from 'next-intl';
import { ChevronDown, Car, Search, X } from 'lucide-react';
import { CarsService } from '@/services/cars';
import type {
  CarSelectorProps,
  CarMake,
  CarModel,
  CarGeneration,
  CarSelection,
} from '@/types/cars';

export const CarSelector: React.FC<CarSelectorProps> = ({
  value = {},
  onChange,
  required = false,
  disabled = false,
  className = '',
  showGenerations = false,
  placeholder = {},
}) => {
  const t = useTranslations('cars');

  const [makes, setMakes] = useState<CarMake[]>([]);
  const [models, setModels] = useState<CarModel[]>([]);
  const [generations, setGenerations] = useState<CarGeneration[]>([]);

  const [loadingMakes, setLoadingMakes] = useState(false);
  const [loadingModels, setLoadingModels] = useState(false);
  const [loadingGenerations, setLoadingGenerations] = useState(false);

  const [searchQuery, setSearchQuery] = useState('');
  const [showMakeSearch, setShowMakeSearch] = useState(false);

  // Загрузка марок при монтировании
  useEffect(() => {
    loadMakes();
  }, []);

  // Загрузка моделей при выборе марки
  useEffect(() => {
    if (value.make?.slug) {
      loadModels(value.make.slug);
    } else {
      setModels([]);
      setGenerations([]);
    }
  }, [value.make]);

  // Загрузка поколений при выборе модели
  useEffect(() => {
    if (value.model?.id && showGenerations) {
      loadGenerations(value.model.id);
    } else {
      setGenerations([]);
    }
  }, [value.model, showGenerations]);

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

  const loadGenerations = async (modelId: number) => {
    setLoadingGenerations(true);
    try {
      const response = await CarsService.getGenerationsByModel(modelId);
      if (response.success && response.data) {
        setGenerations(response.data);
      }
    } catch (error) {
      console.error('Error loading generations:', error);
    } finally {
      setLoadingGenerations(false);
    }
  };

  const handleMakeSelect = useCallback(
    (make: CarMake) => {
      const newSelection: CarSelection = { make };
      onChange(newSelection);
      setShowMakeSearch(false);
      setSearchQuery('');
    },
    [onChange]
  );

  const handleModelSelect = useCallback(
    (model: CarModel) => {
      const newSelection: CarSelection = {
        make: value.make,
        model,
      };
      onChange(newSelection);
    },
    [onChange, value.make]
  );

  const handleGenerationSelect = useCallback(
    (generation: CarGeneration) => {
      const newSelection: CarSelection = {
        make: value.make,
        model: value.model,
        generation,
      };
      onChange(newSelection);
    },
    [onChange, value.make, value.model]
  );

  const clearSelection = () => {
    onChange({});
  };

  const filteredMakes = searchQuery
    ? makes.filter((make) =>
        make.name?.toLowerCase().includes(searchQuery.toLowerCase())
      )
    : makes;

  return (
    <div className={`space-y-4 ${className}`}>
      {/* Марка автомобиля */}
      <div className="form-control">
        <label className="label">
          <span className="label-text font-medium">
            {t('make')} {required && <span className="text-error">*</span>}
          </span>
          {value.make && (
            <button
              type="button"
              onClick={clearSelection}
              className="btn btn-ghost btn-xs"
              title={t('clear')}
            >
              <X className="w-3 h-3" />
            </button>
          )}
        </label>

        {showMakeSearch ? (
          <div className="relative">
            <input
              type="text"
              placeholder={t('searchMake')}
              value={searchQuery}
              onChange={(e) => setSearchQuery(e.target.value)}
              className="input input-bordered w-full pr-10"
              autoFocus
            />
            <Search className="absolute right-3 top-1/2 transform -translate-y-1/2 w-4 h-4 text-base-content/50" />

            {filteredMakes.length > 0 && (
              <div className="absolute z-50 w-full mt-1 bg-base-100 border border-base-300 rounded-lg shadow-lg max-h-60 overflow-y-auto">
                {filteredMakes.map((make) => (
                  <button
                    key={make.id}
                    type="button"
                    onClick={() => handleMakeSelect(make)}
                    className="w-full px-4 py-2 text-left hover:bg-base-200 flex items-center gap-3"
                  >
                    {make.logo_url && (
                      <img
                        src={make.logo_url}
                        alt={make.name}
                        className="w-6 h-6 object-contain"
                      />
                    )}
                    <div>
                      <div className="font-medium">{make.name}</div>
                      {make.country && (
                        <div className="text-xs text-base-content/60">
                          {make.country}
                        </div>
                      )}
                    </div>
                    {make.is_domestic && (
                      <span className="badge badge-primary badge-xs ml-auto">
                        {t('domestic')}
                      </span>
                    )}
                  </button>
                ))}
              </div>
            )}
          </div>
        ) : (
          <div className="dropdown w-full">
            <div
              tabIndex={0}
              role="button"
              className={`btn btn-outline w-full justify-between ${disabled ? 'btn-disabled' : ''}`}
              onClick={() => !disabled && setShowMakeSearch(true)}
            >
              <div className="flex items-center gap-2">
                {value.make ? (
                  <>
                    {value.make.logo_url && (
                      <img
                        src={value.make.logo_url}
                        alt={value.make.name}
                        className="w-5 h-5 object-contain"
                      />
                    )}
                    <span>{value.make.name}</span>
                    {value.make.is_domestic && (
                      <span className="badge badge-primary badge-xs">
                        {t('domestic')}
                      </span>
                    )}
                  </>
                ) : (
                  <>
                    <Car className="w-4 h-4 text-base-content/50" />
                    <span className="text-base-content/50">
                      {placeholder.make || t('selectMake')}
                    </span>
                  </>
                )}
              </div>
              <ChevronDown className="w-4 h-4" />
            </div>
          </div>
        )}

        {loadingMakes && (
          <div className="mt-2">
            <span className="loading loading-spinner loading-sm"></span>
            <span className="ml-2 text-sm text-base-content/60">
              {t('loadingMakes')}
            </span>
          </div>
        )}
      </div>

      {/* Модель автомобиля */}
      {value.make && (
        <div className="form-control">
          <label className="label">
            <span className="label-text font-medium">
              {t('model')} {required && <span className="text-error">*</span>}
            </span>
          </label>

          <div className="dropdown w-full">
            <div
              tabIndex={0}
              role="button"
              className={`btn btn-outline w-full justify-between ${
                disabled || loadingModels || models.length === 0
                  ? 'btn-disabled'
                  : ''
              }`}
            >
              <div className="flex items-center gap-2">
                {value.model ? (
                  <span>{value.model.name}</span>
                ) : (
                  <span className="text-base-content/50">
                    {placeholder.model || t('selectModel')}
                  </span>
                )}
              </div>
              <ChevronDown className="w-4 h-4" />
            </div>

            {!loadingModels && models.length > 0 && (
              <ul
                tabIndex={0}
                className="dropdown-content z-[1] menu p-2 shadow bg-base-100 rounded-box w-full max-h-60 overflow-y-auto border border-base-300"
              >
                {models.map((model) => (
                  <li key={model.id}>
                    <button
                      type="button"
                      onClick={() => handleModelSelect(model)}
                      className="w-full text-left"
                    >
                      <div>
                        <div className="font-medium">{model.name}</div>
                        {(model as any).year_start &&
                          (model as any).year_end && (
                            <div className="text-xs text-base-content/60">
                              {(model as any).year_start} -{' '}
                              {(model as any).year_end}
                            </div>
                          )}
                      </div>
                    </button>
                  </li>
                ))}
              </ul>
            )}
          </div>

          {loadingModels && (
            <div className="mt-2">
              <span className="loading loading-spinner loading-sm"></span>
              <span className="ml-2 text-sm text-base-content/60">
                {t('loadingModels')}
              </span>
            </div>
          )}

          {!loadingModels && models.length === 0 && value.make && (
            <div className="mt-2 text-sm text-base-content/60">
              {t('noModelsFound')}
            </div>
          )}
        </div>
      )}

      {/* Поколение автомобиля */}
      {showGenerations && value.model && (
        <div className="form-control">
          <label className="label">
            <span className="label-text font-medium">{t('generation')}</span>
          </label>

          <div className="dropdown w-full">
            <div
              tabIndex={0}
              role="button"
              className={`btn btn-outline w-full justify-between ${
                disabled || loadingGenerations || generations.length === 0
                  ? 'btn-disabled'
                  : ''
              }`}
            >
              <div className="flex items-center gap-2">
                {value.generation ? (
                  <span>{value.generation.name}</span>
                ) : (
                  <span className="text-base-content/50">
                    {placeholder.generation || t('selectGeneration')}
                  </span>
                )}
              </div>
              <ChevronDown className="w-4 h-4" />
            </div>

            {!loadingGenerations && generations.length > 0 && (
              <ul
                tabIndex={0}
                className="dropdown-content z-[1] menu p-2 shadow bg-base-100 rounded-box w-full max-h-60 overflow-y-auto border border-base-300"
              >
                {generations.map((generation) => (
                  <li key={generation.id}>
                    <button
                      type="button"
                      onClick={() => handleGenerationSelect(generation)}
                      className="w-full text-left"
                    >
                      <div>
                        <div className="font-medium">{generation.name}</div>
                        {generation.year_start && generation.year_end && (
                          <div className="text-xs text-base-content/60">
                            {generation.year_start} - {generation.year_end}
                          </div>
                        )}
                      </div>
                    </button>
                  </li>
                ))}
              </ul>
            )}
          </div>

          {loadingGenerations && (
            <div className="mt-2">
              <span className="loading loading-spinner loading-sm"></span>
              <span className="ml-2 text-sm text-base-content/60">
                {t('loadingGenerations')}
              </span>
            </div>
          )}

          {!loadingGenerations && generations.length === 0 && value.model && (
            <div className="mt-2 text-sm text-base-content/60">
              {t('noGenerationsFound')}
            </div>
          )}
        </div>
      )}
    </div>
  );
};
