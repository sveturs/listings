'use client';

import React, { useState, useEffect } from 'react';
import { useTranslations } from 'next-intl';
import { ChevronDown, Car } from 'lucide-react';
import { CarsService } from '@/services/cars';
import type { CarMake, CarModel, CarSelection } from '@/types/cars';

interface CarSelectorCompactProps {
  value?: CarSelection;
  onChange: (selection: CarSelection) => void;
  disabled?: boolean;
  className?: string;
}

export const CarSelectorCompact: React.FC<CarSelectorCompactProps> = ({
  value = {},
  onChange,
  disabled = false,
  className = '',
}) => {
  const t = useTranslations('cars');

  const [makes, setMakes] = useState<CarMake[]>([]);
  const [models, setModels] = useState<CarModel[]>([]);
  const [loadingMakes, setLoadingMakes] = useState(false);
  const [loadingModels, setLoadingModels] = useState(false);

  useEffect(() => {
    loadMakes();
  }, []);

  useEffect(() => {
    if (value.make?.slug) {
      loadModels(value.make.slug);
    } else {
      setModels([]);
    }
  }, [value.make]);

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

  const handleMakeSelect = (make: CarMake) => {
    onChange({ make });
  };

  const handleModelSelect = (model: CarModel) => {
    onChange({ make: value.make, model });
  };

  return (
    <div className={`flex gap-2 ${className}`}>
      {/* Марка */}
      <div className="dropdown flex-1">
        <div
          tabIndex={0}
          role="button"
          className={`btn btn-outline w-full justify-between ${
            disabled || loadingMakes ? 'btn-disabled' : ''
          }`}
        >
          <div className="flex items-center gap-2 truncate">
            {value.make ? (
              <>
                {value.make.logo_url && (
                  <img
                    src={value.make.logo_url}
                    alt={value.make.name}
                    className="w-4 h-4 object-contain flex-shrink-0"
                  />
                )}
                <span className="truncate">{value.make.name}</span>
              </>
            ) : (
              <>
                <Car className="w-4 h-4 text-base-content/50 flex-shrink-0" />
                <span className="text-base-content/50 truncate">
                  {loadingMakes ? t('loading') : t('selectMake')}
                </span>
              </>
            )}
          </div>
          <ChevronDown className="w-4 h-4 flex-shrink-0" />
        </div>

        {!loadingMakes && makes.length > 0 && (
          <ul
            tabIndex={0}
            className="dropdown-content z-[1] menu p-2 shadow bg-base-100 rounded-box w-full max-h-60 overflow-y-auto border border-base-300"
          >
            {makes.map((make) => (
              <li key={make.id}>
                <button
                  type="button"
                  onClick={() => handleMakeSelect(make)}
                  className="w-full text-left"
                >
                  <div className="flex items-center gap-2">
                    {make.logo_url && (
                      <img
                        src={make.logo_url}
                        alt={make.name}
                        className="w-5 h-5 object-contain"
                      />
                    )}
                    <div className="flex-1">
                      <div className="font-medium">{make.name}</div>
                      {make.country && (
                        <div className="text-xs text-base-content/60">
                          {make.country}
                        </div>
                      )}
                    </div>
                    {make.is_domestic && (
                      <span className="badge badge-primary badge-xs">
                        {t('domestic')}
                      </span>
                    )}
                  </div>
                </button>
              </li>
            ))}
          </ul>
        )}
      </div>

      {/* Модель */}
      <div className="dropdown flex-1">
        <div
          tabIndex={0}
          role="button"
          className={`btn btn-outline w-full justify-between ${
            disabled || loadingModels || !value.make || models.length === 0
              ? 'btn-disabled'
              : ''
          }`}
        >
          <div className="flex items-center gap-2 truncate">
            {value.model ? (
              <span className="truncate">{value.model.name}</span>
            ) : (
              <span className="text-base-content/50 truncate">
                {!value.make
                  ? t('selectMakeFirst')
                  : loadingModels
                    ? t('loading')
                    : models.length === 0
                      ? t('noModels')
                      : t('selectModel')}
              </span>
            )}
          </div>
          <ChevronDown className="w-4 h-4 flex-shrink-0" />
        </div>

        {!loadingModels && models.length > 0 && value.make && (
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
                    {(model as any).year_start && (model as any).year_end && (
                      <div className="text-xs text-base-content/60">
                        {(model as any).year_start} - {(model as any).year_end}
                      </div>
                    )}
                  </div>
                </button>
              </li>
            ))}
          </ul>
        )}
      </div>
    </div>
  );
};
