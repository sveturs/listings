'use client';

import React from 'react';
import { MapContainer, TileLayer, Marker, Popup, Circle } from 'react-leaflet';
import 'leaflet/dist/leaflet.css';
import L from 'leaflet';

// Фикс для иконок Leaflet
if (typeof window !== 'undefined') {
  delete (L.Icon.Default.prototype as any)._getIconUrl;
  L.Icon.Default.mergeOptions({
    iconRetinaUrl: '/icons/marker-icon-2x.png',
    iconUrl: '/icons/marker-icon.png',
    shadowUrl: '/icons/marker-shadow.png',
  });
}

export default function MapSection() {
  return (
    <MapContainer 
      center={[44.8125, 20.4612]} 
      zoom={12} 
      style={{ height: '100%', width: '100%' }}
      className="rounded-lg"
    >
      <TileLayer
        attribution='&copy; <a href="https://www.openstreetmap.org/copyright">OpenStreetMap</a>'
        url="https://{s}.tile.openstreetmap.org/{z}/{x}/{y}.png"
      />
      {/* Круг, показывающий радиус поиска */}
      <Circle 
        center={[44.8125, 20.4612]} 
        radius={5000}
        pathOptions={{ 
          color: 'blue', 
          fillColor: '#3b82f6',
          fillOpacity: 0.1 
        }}
      />
      {/* Маркеры товаров */}
      <Marker position={[44.8125, 20.4612]}>
        <Popup>Вы здесь</Popup>
      </Marker>
      <Marker position={[44.8225, 20.4512]}>
        <Popup>
          <div className="text-sm">
            <strong>iPhone 15 Pro</strong><br/>
            €1099<br/>
            <span className="text-xs">2 км от вас</span>
          </div>
        </Popup>
      </Marker>
      <Marker position={[44.8025, 20.4712]}>
        <Popup>
          <div className="text-sm">
            <strong>MacBook Air M3</strong><br/>
            €1299<br/>
            <span className="text-xs">3 км от вас</span>
          </div>
        </Popup>
      </Marker>
      <Marker position={[44.8325, 20.4412]}>
        <Popup>
          <div className="text-sm">
            <strong>Квартира 2-комн</strong><br/>
            €85000<br/>
            <span className="text-xs">4 км от вас</span>
          </div>
        </Popup>
      </Marker>
      <Marker position={[44.7925, 20.4812]}>
        <Popup>
          <div className="text-sm">
            <strong>PlayStation 5</strong><br/>
            €549<br/>
            <span className="text-xs">2.5 км от вас</span>
          </div>
        </Popup>
      </Marker>
      <Marker position={[44.8425, 20.4312]}>
        <Popup>
          <div className="text-sm">
            <strong>Диван угловой</strong><br/>
            €899<br/>
            <span className="text-xs">5 км от вас</span>
          </div>
        </Popup>
      </Marker>
    </MapContainer>
  );
}