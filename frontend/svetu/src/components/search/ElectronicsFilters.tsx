'use client';

import React, { useState, useEffect } from 'react';
import { useTranslations } from 'next-intl';
import {
  Smartphone,
  Monitor,
  Cpu,
  HardDrive,
  Battery,
  Zap,
} from 'lucide-react';

interface ElectronicsFiltersProps {
  onFiltersChange: (filters: Record<string, any>) => void;
  className?: string;
}

const popularBrands = [
  'Apple',
  'Samsung',
  'Sony',
  'LG',
  'Huawei',
  'Xiaomi',
  'Dell',
  'HP',
  'Lenovo',
  'Asus',
  'MSI',
  'Acer',
];

export const ElectronicsFilters: React.FC<ElectronicsFiltersProps> = ({
  onFiltersChange,
  className = '',
}) => {
  const t = useTranslations('search');

  const [deviceType, setDeviceType] = useState<string>('');
  const [brand, setBrand] = useState<string>('');
  const [screenSizeMin, setScreenSizeMin] = useState<string>('');
  const [screenSizeMax, setScreenSizeMax] = useState<string>('');
  const [ramSize, setRamSize] = useState<string>('');
  const [storageSize, setStorageSize] = useState<string>('');
  const [batteryCapacity, setBatteryCapacity] = useState<string>('');
  const [processor, setProcessor] = useState<string>('');
  const [operatingSystem, setOperatingSystem] = useState<string>('');
  const [hasWarranty, setHasWarranty] = useState(false);
  const [isRefurbished, setIsRefurbished] = useState(false);
  const [has5G, setHas5G] = useState(false);

  useEffect(() => {
    const filters: Record<string, any> = {};

    if (deviceType) filters.deviceType = deviceType;
    if (brand) filters.brand = brand;
    if (screenSizeMin) filters.screenSizeMin = parseFloat(screenSizeMin);
    if (screenSizeMax) filters.screenSizeMax = parseFloat(screenSizeMax);
    if (ramSize) filters.ramSize = ramSize;
    if (storageSize) filters.storageSize = storageSize;
    if (batteryCapacity) filters.batteryCapacity = parseInt(batteryCapacity);
    if (processor) filters.processor = processor;
    if (operatingSystem) filters.operatingSystem = operatingSystem;
    if (hasWarranty) filters.hasWarranty = hasWarranty;
    if (isRefurbished) filters.isRefurbished = isRefurbished;
    if (has5G) filters.has5G = has5G;

    onFiltersChange(filters);
  }, [
    deviceType,
    brand,
    screenSizeMin,
    screenSizeMax,
    ramSize,
    storageSize,
    batteryCapacity,
    processor,
    operatingSystem,
    hasWarranty,
    isRefurbished,
    has5G,
    onFiltersChange,
  ]);

  return (
    <div className={`space-y-4 ${className}`}>
      <div>
        <label className="text-xs font-medium text-base-content/70 mb-2 flex items-center gap-1">
          <Monitor className="w-3 h-3" />
          {t('deviceType')}
        </label>
        <select
          value={deviceType}
          onChange={(e) => setDeviceType(e.target.value)}
          className="select select-bordered select-sm w-full"
        >
          <option value="">{t('allDevices')}</option>
          <option value="smartphone">{t('devices.smartphone')}</option>
          <option value="laptop">{t('devices.laptop')}</option>
          <option value="tablet">{t('devices.tablet')}</option>
          <option value="desktop">{t('devices.desktop')}</option>
          <option value="tv">{t('devices.tv')}</option>
          <option value="smartwatch">{t('devices.smartwatch')}</option>
          <option value="headphones">{t('devices.headphones')}</option>
          <option value="camera">{t('devices.camera')}</option>
          <option value="gaming">{t('devices.gaming')}</option>
        </select>
      </div>

      <div>
        <label className="text-xs font-medium text-base-content/70 mb-2 flex items-center gap-1">
          <Smartphone className="w-3 h-3" />
          {t('brand')}
        </label>
        <select
          value={brand}
          onChange={(e) => setBrand(e.target.value)}
          className="select select-bordered select-sm w-full"
        >
          <option value="">{t('allBrands')}</option>
          {popularBrands.map((b) => (
            <option key={b} value={b.toLowerCase()}>
              {b}
            </option>
          ))}
          <option value="other">{t('other')}</option>
        </select>
      </div>

      {(deviceType === 'smartphone' ||
        deviceType === 'tablet' ||
        deviceType === 'laptop' ||
        deviceType === 'tv' ||
        !deviceType) && (
        <div>
          <label className="text-xs font-medium text-base-content/70 mb-2 flex items-center gap-1">
            <Monitor className="w-3 h-3" />
            {t('screenSize')} ({deviceType === 'tv' ? 'inch' : 'inch'})
          </label>
          <div className="flex gap-2">
            <input
              type="number"
              step="0.1"
              placeholder={t('from')}
              value={screenSizeMin}
              onChange={(e) => setScreenSizeMin(e.target.value)}
              className="input input-bordered input-sm w-full"
            />
            <input
              type="number"
              step="0.1"
              placeholder={t('to')}
              value={screenSizeMax}
              onChange={(e) => setScreenSizeMax(e.target.value)}
              className="input input-bordered input-sm w-full"
            />
          </div>
        </div>
      )}

      {(deviceType === 'smartphone' ||
        deviceType === 'laptop' ||
        deviceType === 'tablet' ||
        deviceType === 'desktop' ||
        !deviceType) && (
        <>
          <div>
            <label className="text-xs font-medium text-base-content/70 mb-2 flex items-center gap-1">
              <Cpu className="w-3 h-3" />
              {t('ram')}
            </label>
            <select
              value={ramSize}
              onChange={(e) => setRamSize(e.target.value)}
              className="select select-bordered select-sm w-full"
            >
              <option value="">{t('anyRam')}</option>
              <option value="2">{t('ramOptions.2gb')}</option>
              <option value="4">{t('ramOptions.4gb')}</option>
              <option value="6">{t('ramOptions.6gb')}</option>
              <option value="8">{t('ramOptions.8gb')}</option>
              <option value="12">{t('ramOptions.12gb')}</option>
              <option value="16">{t('ramOptions.16gb')}</option>
              <option value="32">{t('ramOptions.32gb')}</option>
              <option value="64+">{t('ramOptions.64gbPlus')}</option>
            </select>
          </div>

          <div>
            <label className="text-xs font-medium text-base-content/70 mb-2 flex items-center gap-1">
              <HardDrive className="w-3 h-3" />
              {t('storage')}
            </label>
            <select
              value={storageSize}
              onChange={(e) => setStorageSize(e.target.value)}
              className="select select-bordered select-sm w-full"
            >
              <option value="">{t('anyStorage')}</option>
              <option value="32">{t('storageOptions.32gb')}</option>
              <option value="64">{t('storageOptions.64gb')}</option>
              <option value="128">{t('storageOptions.128gb')}</option>
              <option value="256">{t('storageOptions.256gb')}</option>
              <option value="512">{t('storageOptions.512gb')}</option>
              <option value="1024">{t('storageOptions.1tb')}</option>
              <option value="2048+">{t('storageOptions.2tbPlus')}</option>
            </select>
          </div>
        </>
      )}

      {(deviceType === 'smartphone' ||
        deviceType === 'tablet' ||
        deviceType === 'smartwatch' ||
        !deviceType) && (
        <div>
          <label className="text-xs font-medium text-base-content/70 mb-2 flex items-center gap-1">
            <Battery className="w-3 h-3" />
            {t('batteryCapacity')} (mAh)
          </label>
          <input
            type="number"
            placeholder={t('minBattery')}
            value={batteryCapacity}
            onChange={(e) => setBatteryCapacity(e.target.value)}
            className="input input-bordered input-sm w-full"
          />
        </div>
      )}

      <div>
        <label className="text-xs font-medium text-base-content/70 mb-2 flex items-center gap-1">
          <Zap className="w-3 h-3" />
          {t('processor')}
        </label>
        <input
          type="text"
          placeholder={t('enterProcessor')}
          value={processor}
          onChange={(e) => setProcessor(e.target.value)}
          className="input input-bordered input-sm w-full"
        />
      </div>

      {(deviceType === 'smartphone' ||
        deviceType === 'laptop' ||
        deviceType === 'tablet' ||
        deviceType === 'desktop' ||
        !deviceType) && (
        <div>
          <label className="text-xs font-medium text-base-content/70 mb-2">
            {t('operatingSystem')}
          </label>
          <select
            value={operatingSystem}
            onChange={(e) => setOperatingSystem(e.target.value)}
            className="select select-bordered select-sm w-full"
          >
            <option value="">{t('anyOS')}</option>
            <option value="ios">{t('osOptions.ios')}</option>
            <option value="android">{t('osOptions.android')}</option>
            <option value="windows">{t('osOptions.windows')}</option>
            <option value="macos">{t('osOptions.macos')}</option>
            <option value="linux">{t('osOptions.linux')}</option>
          </select>
        </div>
      )}

      <div className="space-y-2">
        <label className="text-xs font-medium text-base-content/70">
          {t('features')}
        </label>
        <div className="space-y-2">
          <label className="flex items-center gap-2 cursor-pointer">
            <input
              type="checkbox"
              checked={hasWarranty}
              onChange={(e) => setHasWarranty(e.target.checked)}
              className="checkbox checkbox-sm checkbox-primary"
            />
            <span className="text-sm">{t('featuresOptions.warranty')}</span>
          </label>
          <label className="flex items-center gap-2 cursor-pointer">
            <input
              type="checkbox"
              checked={isRefurbished}
              onChange={(e) => setIsRefurbished(e.target.checked)}
              className="checkbox checkbox-sm checkbox-primary"
            />
            <span className="text-sm">{t('featuresOptions.refurbished')}</span>
          </label>
          {(deviceType === 'smartphone' ||
            deviceType === 'tablet' ||
            !deviceType) && (
            <label className="flex items-center gap-2 cursor-pointer">
              <input
                type="checkbox"
                checked={has5G}
                onChange={(e) => setHas5G(e.target.checked)}
                className="checkbox checkbox-sm checkbox-primary"
              />
              <span className="text-sm">{t('featuresOptions.5g')}</span>
            </label>
          )}
        </div>
      </div>
    </div>
  );
};
