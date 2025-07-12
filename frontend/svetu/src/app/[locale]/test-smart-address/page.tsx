'use client';

import React, { useState } from 'react';
import SmartAddressInput from '@/components/GIS/SmartAddressInput';
import AddressConfirmationMap from '@/components/GIS/AddressConfirmationMap';
import LocationPrivacySettings from '@/components/GIS/LocationPrivacySettings';
import { AddressGeocodingResult } from '@/hooks/useAddressGeocoding';

export default function TestSmartAddressPage() {
  const [address, setAddress] = useState('');
  const [location, setLocation] = useState<{ lat: number; lng: number } | undefined>();
  const [confidence, setConfidence] = useState(0);
  const [privacyLevel, setPrivacyLevel] = useState<'exact' | 'street' | 'district' | 'city'>('street');
  const [step, setStep] = useState<'input' | 'confirm' | 'privacy'>('input');

  const handleAddressChange = (value: string, result?: AddressGeocodingResult) => {
    setAddress(value);
    
    if (result) {
      setLocation({
        lat: result.location.lat,
        lng: result.location.lng,
      });
      setConfidence(result.confidence);
    }
  };

  const handleLocationSelect = (locationData: { lat: number; lng: number; address: string; confidence: number }) => {
    setLocation({ lat: locationData.lat, lng: locationData.lng });
    setAddress(locationData.address);
    setConfidence(locationData.confidence);
    setStep('confirm');
  };

  const handleLocationConfirm = (locationData: { lat: number; lng: number; address: string; confidence: number }) => {
    setLocation({ lat: locationData.lat, lng: locationData.lng });
    setAddress(locationData.address);
    setConfidence(locationData.confidence);
    setStep('privacy');
  };

  const handleLocationChange = (newLocation: { lat: number; lng: number }) => {
    setLocation(newLocation);
  };

  return (
    <div className="min-h-screen bg-base-200 py-8">
      <div className="container mx-auto px-4 max-w-4xl">
        {/* Заголовок */}
        <div className="mb-8 text-center">
          <h1 className="text-3xl font-bold mb-2">GIS Phase 2: Умный ввод адресов</h1>
          <p className="text-base-content/70">
            Тестирование компонентов SmartAddressInput, AddressConfirmationMap и LocationPrivacySettings
          </p>
        </div>

        {/* Навигация по шагам */}
        <div className="mb-8">
          <div className="flex justify-center">
            <div className="steps">
              <div 
                className={`step ${step === 'input' ? 'step-primary' : ''} ${location ? 'step-success' : ''}`}
                onClick={() => setStep('input')}
              >
                Ввод адреса
              </div>
              <div 
                className={`step ${step === 'confirm' ? 'step-primary' : ''} ${step === 'privacy' ? 'step-success' : ''}`}
                onClick={() => location && setStep('confirm')}
              >
                Подтверждение
              </div>
              <div 
                className={`step ${step === 'privacy' ? 'step-primary' : ''}`}
                onClick={() => location && setStep('privacy')}
              >
                Приватность
              </div>
            </div>
          </div>
        </div>

        {/* Контент в зависимости от шага */}
        <div className="space-y-8">
          {/* Шаг 1: Ввод адреса */}
          {step === 'input' && (
            <div className="card bg-base-100 shadow-xl">
              <div className="card-body">
                <h2 className="card-title mb-4">
                  <span className="text-2xl mr-2">📍</span>
                  Шаг 1: Введите адрес
                </h2>
                
                <div className="space-y-4">
                  <SmartAddressInput
                    value={address}
                    onChange={handleAddressChange}
                    onLocationSelect={handleLocationSelect}
                    placeholder="Начните вводить адрес (например: Београд, Кнез Михаилова)"
                    showCurrentLocation={true}
                    country={['rs', 'hr', 'ba', 'me']}
                    language="ru"
                  />
                  
                  {location && (
                    <div className="mt-4 p-4 bg-success/10 border border-success/20 rounded-lg">
                      <h3 className="font-medium text-success-content mb-2">✅ Адрес найден!</h3>
                      <div className="text-sm text-success-content/80 space-y-1">
                        <p><strong>Адрес:</strong> {address}</p>
                        <p><strong>Координаты:</strong> {location.lat.toFixed(6)}, {location.lng.toFixed(6)}</p>
                        <p><strong>Точность:</strong> {Math.round(confidence * 100)}%</p>
                      </div>
                      
                      <div className="mt-3">
                        <button 
                          className="btn btn-primary"
                          onClick={() => setStep('confirm')}
                        >
                          Перейти к подтверждению
                        </button>
                      </div>
                    </div>
                  )}
                </div>
              </div>
            </div>
          )}

          {/* Шаг 2: Подтверждение на карте */}
          {step === 'confirm' && location && (
            <div className="card bg-base-100 shadow-xl">
              <div className="card-body">
                <h2 className="card-title mb-4">
                  <span className="text-2xl mr-2">🗺️</span>
                  Шаг 2: Подтвердите местоположение на карте
                </h2>
                
                <AddressConfirmationMap
                  address={address}
                  initialLocation={location}
                  onLocationConfirm={handleLocationConfirm}
                  onLocationChange={handleLocationChange}
                  editable={true}
                  zoom={16}
                  height="500px"
                />
              </div>
            </div>
          )}

          {/* Шаг 3: Настройки приватности */}
          {step === 'privacy' && location && (
            <div className="card bg-base-100 shadow-xl">
              <div className="card-body">
                <h2 className="card-title mb-4">
                  <span className="text-2xl mr-2">🛡️</span>
                  Шаг 3: Настройки приватности
                </h2>
                
                <LocationPrivacySettings
                  selectedLevel={privacyLevel}
                  onLevelChange={setPrivacyLevel}
                  location={location}
                  showPreview={true}
                />
                
                <div className="mt-6 flex gap-2">
                  <button 
                    className="btn btn-outline"
                    onClick={() => setStep('confirm')}
                  >
                    ← Назад к карте
                  </button>
                  
                  <button 
                    className="btn btn-primary flex-1"
                    onClick={() => {
                      alert(`Настройки сохранены!\n\nАдрес: ${address}\nКоординаты: ${location.lat.toFixed(6)}, ${location.lng.toFixed(6)}\nПриватность: ${privacyLevel}\nТочность: ${Math.round(confidence * 100)}%`);
                    }}
                  >
                    Сохранить объявление
                  </button>
                </div>
              </div>
            </div>
          )}

          {/* Информационная панель */}
          <div className="card bg-base-100 shadow-xl">
            <div className="card-body">
              <h2 className="card-title mb-4">
                <span className="text-2xl mr-2">📊</span>
                Текущее состояние
              </h2>
              
              <div className="grid grid-cols-1 md:grid-cols-2 gap-4">
                <div>
                  <h3 className="font-medium mb-2">Адрес</h3>
                  <p className="text-sm text-base-content/70 bg-base-200 p-2 rounded">
                    {address || 'Не указан'}
                  </p>
                </div>
                
                <div>
                  <h3 className="font-medium mb-2">Координаты</h3>
                  <p className="text-sm text-base-content/70 bg-base-200 p-2 rounded font-mono">
                    {location ? `${location.lat.toFixed(6)}, ${location.lng.toFixed(6)}` : 'Не указаны'}
                  </p>
                </div>
                
                <div>
                  <h3 className="font-medium mb-2">Точность</h3>
                  <p className="text-sm text-base-content/70 bg-base-200 p-2 rounded">
                    {confidence ? `${Math.round(confidence * 100)}%` : 'Неизвестно'}
                  </p>
                </div>
                
                <div>
                  <h3 className="font-medium mb-2">Приватность</h3>
                  <p className="text-sm text-base-content/70 bg-base-200 p-2 rounded">
                    {privacyLevel}
                  </p>
                </div>
              </div>
            </div>
          </div>

          {/* Инструкции */}
          <div className="card bg-info/10 border border-info/20">
            <div className="card-body">
              <h2 className="card-title text-info-content mb-4">
                <span className="text-2xl mr-2">💡</span>
                Как пользоваться
              </h2>
              
              <div className="text-sm text-info-content/80 space-y-2">
                <p><strong>1. Ввод адреса:</strong> Начните вводить адрес и выберите из предложений или используйте кнопку геолокации</p>
                <p><strong>2. Подтверждение:</strong> Проверьте местоположение на карте, при необходимости скорректируйте перетаскиванием маркера</p>
                <p><strong>3. Приватность:</strong> Выберите уровень приватности для отображения местоположения другим пользователям</p>
              </div>
            </div>
          </div>
        </div>
      </div>
    </div>
  );
}