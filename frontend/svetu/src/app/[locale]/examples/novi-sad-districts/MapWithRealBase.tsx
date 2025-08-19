'use client';

import React, { useEffect, useRef, useState } from 'react';

// Leaflet types and imports
declare global {
  interface Window {
    L: any;
  }
}

interface District {
  id: string;
  name: string;
  nameEn: string;
  coordinates: number[][];
  color: string;
  population: string;
  area: string;
  description: string;
}

interface MapWithRealBaseProps {
  districts: District[];
  selectedDistrict: District | null;
  onDistrictClick: (district: District) => void;
}

const MapWithRealBase: React.FC<MapWithRealBaseProps> = ({
  districts,
  selectedDistrict,
  onDistrictClick,
}) => {
  const mapRef = useRef<HTMLDivElement>(null);
  const leafletMapRef = useRef<any>(null);
  const polygonsRef = useRef<any[]>([]);
  const [mapReady, setMapReady] = useState(false);
  const [coordsIssues, setCoordsIssues] = useState<string[]>([]);

  useEffect(() => {
    if (typeof window !== 'undefined') {
      // Load Leaflet CSS and JS
      const loadLeaflet = async () => {
        // Load CSS
        if (!document.querySelector('link[href*="leaflet"]')) {
          const link = document.createElement('link');
          link.rel = 'stylesheet';
          link.href = 'https://unpkg.com/leaflet@1.9.4/dist/leaflet.css';
          document.head.appendChild(link);
        }

        // Load JS
        if (!window.L) {
          const script = document.createElement('script');
          script.src = 'https://unpkg.com/leaflet@1.9.4/dist/leaflet.js';
          script.onload = () => {
            initializeMap();
          };
          document.head.appendChild(script);
        } else {
          initializeMap();
        }
      };

      const initializeMap = () => {
        if (!mapRef.current || leafletMapRef.current) return;

        const L = window.L;
        
        // Initialize map centered on Novi Sad
        const map = L.map(mapRef.current).setView([45.2671, 19.8335], 12);

        // Add OpenStreetMap tile layer
        L.tileLayer('https://{s}.tile.openstreetmap.org/{z}/{x}/{y}.png', {
          attribution: '© OpenStreetMap contributors'
        }).addTo(map);

        leafletMapRef.current = map;
        setMapReady(true);
      };

      loadLeaflet();
    }

    return () => {
      if (leafletMapRef.current) {
        leafletMapRef.current.remove();
        leafletMapRef.current = null;
      }
    };
  }, []);

  // Update districts on map when they change
  useEffect(() => {
    if (!mapReady || !leafletMapRef.current) return;

    const L = window.L;
    const map = leafletMapRef.current;
    const issues: string[] = [];

    // Remove old polygons
    polygonsRef.current.forEach(polygon => {
      map.removeLayer(polygon);
    });
    polygonsRef.current = [];

    // Add new polygons
    districts.forEach((district) => {
          try {
            // Convert coordinates to Leaflet format [lat, lng]
            const leafletCoords = district.coordinates.map(coord => [coord[1], coord[0]]);
            
            // Check if coordinates are valid
            const validCoords = leafletCoords.every(coord => 
              coord[0] >= 44 && coord[0] <= 46 && // Valid latitude for Novi Sad area
              coord[1] >= 19 && coord[1] <= 21   // Valid longitude for Novi Sad area
            );

            if (!validCoords) {
              issues.push(`${district.name}: координаты вне области Нови-Сада`);
              return;
            }

            const polygon = L.polygon(leafletCoords, {
              color: district.color,
              weight: 2,
              opacity: 0.8,
              fillColor: district.color,
              fillOpacity: 0.3,
            }).addTo(map);

            // Add popup with district info
            polygon.bindPopup(`
              <div class="p-2">
                <h3 class="font-bold">${district.name}</h3>
                <p class="text-sm">${district.nameEn}</p>
                <p class="text-xs mt-1">${district.description}</p>
                <div class="mt-2 text-xs">
                  <div>Население: ${district.population}</div>
                  <div>Площадь: ${district.area}</div>
                </div>
              </div>
            `);

            // Add click handler
            polygon.on('click', () => {
              onDistrictClick(district);
            });

            // Add hover effect
            polygon.on('mouseover', (e: any) => {
              const layer = e.target;
              layer.setStyle({
                weight: 3,
                fillOpacity: 0.5,
              });
            });

            polygon.on('mouseout', (e: any) => {
              const layer = e.target;
              layer.setStyle({
                weight: 2,
                fillOpacity: 0.3,
              });
            });

            // Store polygon reference
            polygonsRef.current.push(polygon);
            
            // Add marker for first point (to show different starting points)
            if (district.id === 'veternik' && district.coordinates.length > 0) {
              const firstPoint = district.coordinates[0];
              const marker = L.circleMarker([firstPoint[1], firstPoint[0]], {
                radius: 8,
                fillColor: '#FFFFFF',
                color: district.color,
                weight: 3,
                opacity: 1,
                fillOpacity: 0.9
              }).addTo(map);
              
              marker.bindPopup(`
                <div class="p-2">
                  <strong>Начальная точка</strong><br/>
                  ${district.description.includes('OVERPASS') ? 'Overpass API' : 'Nominatim API'}<br/>
                  Координаты: [${firstPoint[0].toFixed(4)}, ${firstPoint[1].toFixed(4)}]
                </div>
              `);
              
              polygonsRef.current.push(marker);
            }

          } catch (error) {
            issues.push(`${district.name}: ошибка отображения полигона`);
            console.error(`Error adding district ${district.name}:`, error);
          }
        });

    setCoordsIssues(issues);
  }, [districts, onDistrictClick, mapReady]);

  // Update selected district
  useEffect(() => {
    if (leafletMapRef.current && selectedDistrict && mapReady) {
      const L = window.L;
      const leafletCoords = selectedDistrict.coordinates.map(coord => [coord[1], coord[0]]);
      const bounds = L.latLngBounds(leafletCoords);
      leafletMapRef.current.fitBounds(bounds, { padding: [20, 20] });
    }
  }, [selectedDistrict, mapReady]);

  return (
    <div className="relative">
      <div 
        ref={mapRef} 
        className="w-full h-[600px] border border-gray-300 rounded-lg shadow-lg"
      />
      
      {!mapReady && (
        <div className="absolute inset-0 flex items-center justify-center bg-gray-100 rounded-lg">
          <div className="text-center">
            <div className="animate-spin w-8 h-8 border-4 border-primary border-t-transparent rounded-full mx-auto mb-2"></div>
            <div className="text-sm">Загрузка карты...</div>
          </div>
        </div>
      )}

      {coordsIssues.length > 0 && (
        <div className="absolute top-4 right-4 bg-red-100 border border-red-400 text-red-700 px-3 py-2 rounded max-w-xs">
          <div className="font-bold text-xs mb-1">⚠️ Проблемы с координатами:</div>
          {coordsIssues.map((issue, index) => (
            <div key={index} className="text-xs">{issue}</div>
          ))}
        </div>
      )}

      <div className="absolute bottom-4 left-4 bg-white bg-opacity-90 p-2 rounded text-xs">
        <div className="font-bold">Источник карты:</div>
        <div>OpenStreetMap (реальная подложка)</div>
        <div className="mt-1 text-gray-600">
          Если полигоны не совпадают с реальностью - координаты неточные!
        </div>
      </div>
    </div>
  );
};

export default MapWithRealBase;