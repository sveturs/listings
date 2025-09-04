'use client';

import React, { useState, useEffect } from 'react';
import { useTranslations } from 'next-intl';
import {
  Home,
  Maximize,
  Layers,
  DoorOpen,
  MapPin,
  Building,
} from 'lucide-react';

interface RealEstateFiltersProps {
  onFiltersChange: (filters: Record<string, any>) => void;
  className?: string;
}

export const RealEstateFilters: React.FC<RealEstateFiltersProps> = ({
  onFiltersChange,
  className = '',
}) => {
  const t = useTranslations('realEstate');

  const [propertyType, setPropertyType] = useState<string>('');
  const [rooms, setRooms] = useState<string>('');
  const [areaMin, setAreaMin] = useState<string>('');
  const [areaMax, setAreaMax] = useState<string>('');
  const [floorMin, setFloorMin] = useState<string>('');
  const [floorMax, setFloorMax] = useState<string>('');
  const [buildingType, setBuildingType] = useState<string>('');
  const [district, setDistrict] = useState<string>('');
  const [hasParking, setHasParking] = useState(false);
  const [hasElevator, setHasElevator] = useState(false);
  const [hasTerrace, setHasTerrace] = useState(false);
  const [isNewBuilding, setIsNewBuilding] = useState(false);

  useEffect(() => {
    const filters: Record<string, any> = {};

    if (propertyType) filters.propertyType = propertyType;
    if (rooms) filters.rooms = rooms;
    if (areaMin) filters.areaMin = parseInt(areaMin);
    if (areaMax) filters.areaMax = parseInt(areaMax);
    if (floorMin) filters.floorMin = parseInt(floorMin);
    if (floorMax) filters.floorMax = parseInt(floorMax);
    if (buildingType) filters.buildingType = buildingType;
    if (district) filters.district = district;
    if (hasParking) filters.hasParking = hasParking;
    if (hasElevator) filters.hasElevator = hasElevator;
    if (hasTerrace) filters.hasTerrace = hasTerrace;
    if (isNewBuilding) filters.isNewBuilding = isNewBuilding;

    onFiltersChange(filters);
  }, [
    propertyType,
    rooms,
    areaMin,
    areaMax,
    floorMin,
    floorMax,
    buildingType,
    district,
    hasParking,
    hasElevator,
    hasTerrace,
    isNewBuilding,
    onFiltersChange,
  ]);

  return (
    <div className={`space-y-4 ${className}`}>
      <div>
        <label className="text-xs font-medium text-base-content/70 mb-2 flex items-center gap-1">
          <Home className="w-3 h-3" />
          {t('propertyType')}
        </label>
        <select
          value={propertyType}
          onChange={(e) => setPropertyType(e.target.value)}
          className="select select-bordered select-sm w-full"
        >
          <option value="">{t('allTypes')}</option>
          <option value="apartment">{t('types.apartment')}</option>
          <option value="house">{t('types.house')}</option>
          <option value="land">{t('types.land')}</option>
          <option value="commercial">{t('types.commercial')}</option>
          <option value="garage">{t('types.garage')}</option>
        </select>
      </div>

      <div>
        <label className="text-xs font-medium text-base-content/70 mb-2 flex items-center gap-1">
          <DoorOpen className="w-3 h-3" />
          {t('rooms')}
        </label>
        <select
          value={rooms}
          onChange={(e) => setRooms(e.target.value)}
          className="select select-bordered select-sm w-full"
        >
          <option value="">{t('anyRooms')}</option>
          <option value="studio">{t('roomOptions.studio')}</option>
          <option value="1">{t('roomOptions.one')}</option>
          <option value="2">{t('roomOptions.two')}</option>
          <option value="3">{t('roomOptions.three')}</option>
          <option value="4">{t('roomOptions.four')}</option>
          <option value="5+">{t('roomOptions.fivePlus')}</option>
        </select>
      </div>

      <div>
        <label className="text-xs font-medium text-base-content/70 mb-2 flex items-center gap-1">
          <Maximize className="w-3 h-3" />
          {t('area')} (m²)
        </label>
        <div className="flex gap-2">
          <input
            type="number"
            placeholder={t('from')}
            value={areaMin}
            onChange={(e) => setAreaMin(e.target.value)}
            className="input input-bordered input-sm w-full"
          />
          <input
            type="number"
            placeholder={t('to')}
            value={areaMax}
            onChange={(e) => setAreaMax(e.target.value)}
            className="input input-bordered input-sm w-full"
          />
        </div>
      </div>

      <div>
        <label className="text-xs font-medium text-base-content/70 mb-2 flex items-center gap-1">
          <Layers className="w-3 h-3" />
          {t('floor')}
        </label>
        <div className="flex gap-2">
          <input
            type="number"
            placeholder={t('from')}
            value={floorMin}
            onChange={(e) => setFloorMin(e.target.value)}
            className="input input-bordered input-sm w-full"
          />
          <input
            type="number"
            placeholder={t('to')}
            value={floorMax}
            onChange={(e) => setFloorMax(e.target.value)}
            className="input input-bordered input-sm w-full"
          />
        </div>
      </div>

      <div>
        <label className="text-xs font-medium text-base-content/70 mb-2 flex items-center gap-1">
          <Building className="w-3 h-3" />
          {t('buildingType')}
        </label>
        <select
          value={buildingType}
          onChange={(e) => setBuildingType(e.target.value)}
          className="select select-bordered select-sm w-full"
        >
          <option value="">{t('anyBuilding')}</option>
          <option value="brick">{t('buildingTypes.brick')}</option>
          <option value="panel">{t('buildingTypes.panel')}</option>
          <option value="monolithic">{t('buildingTypes.monolithic')}</option>
          <option value="wood">{t('buildingTypes.wood')}</option>
        </select>
      </div>

      <div>
        <label className="text-xs font-medium text-base-content/70 mb-2 flex items-center gap-1">
          <MapPin className="w-3 h-3" />
          {t('district')}
        </label>
        <input
          type="text"
          placeholder={t('enterDistrict')}
          value={district}
          onChange={(e) => setDistrict(e.target.value)}
          className="input input-bordered input-sm w-full"
        />
      </div>

      <div className="space-y-2">
        <label className="text-xs font-medium text-base-content/70">
          {t('features')}
        </label>
        <div className="space-y-2">
          <label className="flex items-center gap-2 cursor-pointer">
            <input
              type="checkbox"
              checked={hasParking}
              onChange={(e) => setHasParking(e.target.checked)}
              className="checkbox checkbox-sm checkbox-primary"
            />
            <span className="text-sm">{t('featuresOptions.parking')}</span>
          </label>
          <label className="flex items-center gap-2 cursor-pointer">
            <input
              type="checkbox"
              checked={hasElevator}
              onChange={(e) => setHasElevator(e.target.checked)}
              className="checkbox checkbox-sm checkbox-primary"
            />
            <span className="text-sm">{t('featuresOptions.elevator')}</span>
          </label>
          <label className="flex items-center gap-2 cursor-pointer">
            <input
              type="checkbox"
              checked={hasTerrace}
              onChange={(e) => setHasTerrace(e.target.checked)}
              className="checkbox checkbox-sm checkbox-primary"
            />
            <span className="text-sm">{t('featuresOptions.terrace')}</span>
          </label>
          <label className="flex items-center gap-2 cursor-pointer">
            <input
              type="checkbox"
              checked={isNewBuilding}
              onChange={(e) => setIsNewBuilding(e.target.checked)}
              className="checkbox checkbox-sm checkbox-primary"
            />
            <span className="text-sm">{t('featuresOptions.newBuilding')}</span>
          </label>
        </div>
      </div>
    </div>
  );
};
