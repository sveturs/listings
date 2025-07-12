'use client';

import React, { useState, useCallback } from 'react';
import { InteractiveMap } from '@/components/GIS';
import {
  useMapboxGeoSearch,
  useMapboxGeocoding,
} from '@/components/GIS/hooks/useMapboxGeoSearch';
// Тип события клика по карте - `any` как в InteractiveMap

export default function TestGeocodingPage() {
  const [searchQuery, setSearchQuery] = useState('');
  const [clickedCoordinates, setClickedCoordinates] = useState<{
    lat: number;
    lng: number;
  } | null>(null);
  const [mapCenter, setMapCenter] = useState<[number, number]>([
    20.4577, 44.8205,
  ]); // Belgrade center
  const [mapZoom, setMapZoom] = useState(12);

  // Хук для поиска адресов
  const {
    results: searchResults,
    loading: isSearching,
    error: searchError,
    search,
  } = useMapboxGeoSearch();

  // Хук для обратного геокодирования
  const {
    loading: isReverseGeocoding,
    error: reverseGeocodingError,
    reverseGeocode,
  } = useMapboxGeocoding();

  // Состояние для результата обратного геокодирования
  const [reverseGeocodingResult, setReverseGeocodingResult] =
    useState<any>(null);

  // Обработчик поиска
  const handleSearch = useCallback(
    async (e: React.FormEvent) => {
      e.preventDefault();
      if (searchQuery.trim()) {
        await search({
          query: searchQuery,
          limit: 5,
          language: 'ru',
        });
      }
    },
    [searchQuery, search]
  );

  // Обработчик клика по карте
  const handleMapClick = useCallback(
    async (event: any) => {
      const { lng, lat } = event.lngLat;
      setClickedCoordinates({ lat, lng });
      const result = await reverseGeocode(lat, lng);
      setReverseGeocodingResult(result);
    },
    [reverseGeocode]
  );

  // Обработчик выбора результата поиска
  const handleSelectSearchResult = useCallback((result: any) => {
    const lat = parseFloat(result.lat);
    const lng = parseFloat(result.lon);
    setMapCenter([lng, lat]);
    setMapZoom(16);
    setClickedCoordinates({ lat, lng });
  }, []);

  return (
    <div className="container mx-auto p-4">
      <h1 className="text-2xl font-bold mb-6">MapBox Geocoding Test</h1>

      <div className="grid grid-cols-1 lg:grid-cols-2 gap-6">
        {/* Левая колонка - форма и результаты */}
        <div className="space-y-6">
          {/* Форма поиска */}
          <div className="card bg-base-100 shadow-xl">
            <div className="card-body">
              <h2 className="card-title">Forward Geocoding (Address Search)</h2>
              <form onSubmit={handleSearch} className="space-y-4">
                <div className="form-control">
                  <label className="label">
                    <span className="label-text">Enter address</span>
                  </label>
                  <input
                    type="text"
                    value={searchQuery}
                    onChange={(e) => setSearchQuery(e.target.value)}
                    placeholder="e.g., Knez Mihailova, Belgrade"
                    className="input input-bordered w-full"
                  />
                </div>
                <button
                  type="submit"
                  className={`btn btn-primary w-full ${isSearching ? 'loading' : ''}`}
                  disabled={isSearching || !searchQuery.trim()}
                >
                  {isSearching ? 'Searching...' : 'Search'}
                </button>
              </form>

              {/* Ошибка поиска */}
              {searchError && (
                <div className="alert alert-error mt-4">
                  <span>Error: {searchError}</span>
                </div>
              )}

              {/* Результаты поиска */}
              {searchResults.length > 0 && (
                <div className="mt-4">
                  <h3 className="font-semibold mb-2">Search Results:</h3>
                  <div className="space-y-2 max-h-64 overflow-y-auto">
                    {searchResults.map((result, index) => (
                      <div
                        key={index}
                        className="p-3 bg-base-200 rounded-lg cursor-pointer hover:bg-base-300 transition-colors"
                        onClick={() => handleSelectSearchResult(result)}
                      >
                        <div className="font-medium">{result.display_name}</div>
                        <div className="text-sm text-base-content/70">
                          Coordinates: {parseFloat(result.lat).toFixed(6)},{' '}
                          {parseFloat(result.lon).toFixed(6)}
                        </div>
                        {result.type && (
                          <div className="text-sm text-base-content/70">
                            Type: {result.type}
                          </div>
                        )}
                      </div>
                    ))}
                  </div>
                </div>
              )}
            </div>
          </div>

          {/* Обратное геокодирование */}
          <div className="card bg-base-100 shadow-xl">
            <div className="card-body">
              <h2 className="card-title">Reverse Geocoding</h2>
              <p className="text-base-content/70">
                Click on the map to get address
              </p>

              {/* Координаты клика */}
              {clickedCoordinates && (
                <div className="mt-4 p-3 bg-base-200 rounded-lg">
                  <div className="font-medium">Clicked Coordinates:</div>
                  <div className="text-sm">
                    Lat: {clickedCoordinates.lat.toFixed(6)}, Lng:{' '}
                    {clickedCoordinates.lng.toFixed(6)}
                  </div>
                </div>
              )}

              {/* Загрузка обратного геокодирования */}
              {isReverseGeocoding && (
                <div className="flex items-center justify-center mt-4">
                  <span className="loading loading-spinner loading-md"></span>
                  <span className="ml-2">Getting address...</span>
                </div>
              )}

              {/* Ошибка обратного геокодирования */}
              {reverseGeocodingError && (
                <div className="alert alert-error mt-4">
                  <span>Error: {reverseGeocodingError}</span>
                </div>
              )}

              {/* Результат обратного геокодирования */}
              {reverseGeocodingResult && !isReverseGeocoding && (
                <div className="mt-4 space-y-3">
                  <h3 className="font-semibold">Address Found:</h3>
                  <div className="p-3 bg-success/10 border border-success rounded-lg">
                    <div className="font-medium text-success">
                      {reverseGeocodingResult.display_name}
                    </div>
                  </div>

                  {/* Детали адреса */}
                  <div className="space-y-1">
                    <h4 className="font-medium text-sm">Details:</h4>
                    <div className="text-sm text-base-content/70">
                      Type: {reverseGeocodingResult.type || 'Unknown'}
                    </div>
                    <div className="text-sm text-base-content/70">
                      Class: {reverseGeocodingResult.class || 'Unknown'}
                    </div>
                    <div className="text-sm text-base-content/70">
                      Importance: {reverseGeocodingResult.importance || 0}
                    </div>
                    {reverseGeocodingResult.boundingbox && (
                      <div className="text-sm text-base-content/70">
                        Bounding Box: [
                        {reverseGeocodingResult.boundingbox.join(', ')}]
                      </div>
                    )}
                  </div>
                </div>
              )}
            </div>
          </div>

          {/* Инструкции */}
          <div className="card bg-info/10 border border-info">
            <div className="card-body">
              <h3 className="card-title text-info">Instructions</h3>
              <ul className="list-disc list-inside space-y-1 text-sm">
                <li>
                  Use the search form to find addresses (forward geocoding)
                </li>
                <li>
                  Click on the map to get address at that location (reverse
                  geocoding)
                </li>
                <li>
                  Search results will center the map on the selected location
                </li>
                <li>All geocoding requests are made to MapBox API</li>
              </ul>
            </div>
          </div>
        </div>

        {/* Правая колонка - карта */}
        <div className="card bg-base-100 shadow-xl">
          <div className="card-body p-0">
            <div className="h-[600px] lg:h-[800px]">
              <InteractiveMap
                initialViewState={{
                  longitude: mapCenter[0],
                  latitude: mapCenter[1],
                  zoom: mapZoom,
                }}
                onMapClick={handleMapClick}
                markers={
                  clickedCoordinates
                    ? [
                        {
                          id: 'clicked-point',
                          position: [
                            clickedCoordinates.lng,
                            clickedCoordinates.lat,
                          ],
                          longitude: clickedCoordinates.lng,
                          latitude: clickedCoordinates.lat,
                          title:
                            reverseGeocodingResult?.display_name ||
                            'Clicked Location',
                          type: 'poi' as const,
                        },
                      ]
                    : []
                }
              />
            </div>
          </div>
        </div>
      </div>

      {/* Debug информация */}
      <div className="mt-6 card bg-base-100 shadow-xl">
        <div className="card-body">
          <h3 className="card-title">Debug Information</h3>
          <div className="grid grid-cols-1 md:grid-cols-2 gap-4 text-sm font-mono">
            <div>
              <strong>Map Center:</strong> [{mapCenter[0].toFixed(6)},{' '}
              {mapCenter[1].toFixed(6)}]
            </div>
            <div>
              <strong>Map Zoom:</strong> {mapZoom}
            </div>
            <div>
              <strong>Search Results Count:</strong> {searchResults.length}
            </div>
            <div>
              <strong>API Status:</strong>{' '}
              {isSearching || isReverseGeocoding ? 'Loading' : 'Ready'}
            </div>
          </div>
        </div>
      </div>
    </div>
  );
}
