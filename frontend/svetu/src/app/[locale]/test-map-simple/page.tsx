'use client';

import React from 'react';
import dynamic from 'next/dynamic';

// Динамический импорт для избежания SSR проблем
const Map = dynamic(() => import('react-map-gl').then((mod) => mod.default), {
  ssr: false,
});

export default function TestMapSimplePage() {
  const mapboxToken = process.env.NEXT_PUBLIC_MAPBOX_ACCESS_TOKEN;

  return (
    <div className="container mx-auto p-4">
      <h1 className="text-2xl font-bold mb-4">Simple MapBox Test</h1>

      <div className="mb-4 p-4 bg-base-200 rounded">
        <p>
          <strong>MapBox Token:</strong> {mapboxToken ? 'Present' : 'Missing'}
        </p>
        <p>
          <strong>Token Length:</strong> {mapboxToken?.length || 0}
        </p>
      </div>

      <div className="h-96 bg-gray-200 rounded">
        {mapboxToken ? (
          <Map
            mapboxAccessToken={mapboxToken}
            initialViewState={{
              longitude: 20.4577,
              latitude: 44.8205,
              zoom: 12,
            }}
            style={{ width: '100%', height: '100%' }}
            mapStyle="mapbox://styles/mapbox/streets-v12"
          />
        ) : (
          <div className="flex items-center justify-center h-full">
            <p className="text-red-500">MapBox token not found</p>
          </div>
        )}
      </div>
    </div>
  );
}
