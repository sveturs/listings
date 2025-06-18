'use client';

import { MapContainer, TileLayer, Marker, Popup } from 'react-leaflet';
import { LatLngExpression, Icon } from 'leaflet';
import 'leaflet/dist/leaflet.css';
import type { Storefront } from '@/types/storefront';

// Fix for default marker icons
if (typeof window !== 'undefined') {
  // eslint-disable-next-line @typescript-eslint/no-require-imports
  const L = require('leaflet');
  delete (L.Icon.Default.prototype as any)._getIconUrl;
  L.Icon.Default.mergeOptions({
    iconRetinaUrl: '/leaflet/marker-icon-2x.png',
    iconUrl: '/leaflet/marker-icon.png',
    shadowUrl: '/leaflet/marker-shadow.png',
  });
}

interface StorefrontLocationMapProps {
  location: {
    user_lat: number;
    user_lng: number;
    full_address: string;
  };
  storefront: Storefront;
}

export default function StorefrontLocationMap({ location, storefront }: StorefrontLocationMapProps) {
  const position: LatLngExpression = [location.user_lat, location.user_lng];

  // Custom icon for storefront
  const storefrontIcon = new Icon({
    iconUrl: '/leaflet/store-marker.png',
    iconRetinaUrl: '/leaflet/store-marker-2x.png',
    shadowUrl: '/leaflet/marker-shadow.png',
    iconSize: [32, 48],
    iconAnchor: [16, 48],
    popupAnchor: [0, -48],
  });

  const handleDirections = () => {
    const url = `https://www.google.com/maps/dir/?api=1&destination=${location.user_lat},${location.user_lng}`;
    window.open(url, '_blank');
  };

  return (
    <div className="relative h-full w-full">
      <MapContainer 
        center={position} 
        zoom={15} 
        className="h-full w-full rounded-lg"
        scrollWheelZoom={false}
      >
        <TileLayer
          attribution='&copy; <a href="https://www.openstreetmap.org/copyright">OpenStreetMap</a> contributors'
          url="https://{s}.tile.openstreetmap.org/{z}/{x}/{y}.png"
        />
        <Marker position={position} icon={storefrontIcon}>
          <Popup>
            <div className="p-2">
              <h3 className="font-bold text-lg">{storefront.name}</h3>
              <p className="text-sm mt-1">{location.full_address}</p>
              <button 
                className="btn btn-primary btn-sm mt-3 w-full"
                onClick={handleDirections}
              >
                <svg className="w-4 h-4" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                  <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M9 20l-5.447-2.724A1 1 0 013 16.382V5.618a1 1 0 011.447-.894L9 7m0 13l6-3m-6 3V7m6 10l4.553 2.276A1 1 0 0021 18.382V7.618a1 1 0 00-.553-.894L15 4m0 13V4m0 0L9 7" />
                </svg>
                Get Directions
              </button>
            </div>
          </Popup>
        </Marker>
      </MapContainer>

      {/* Overlay buttons */}
      <div className="absolute bottom-4 right-4 z-[1000] space-y-2">
        <button 
          className="btn btn-sm btn-circle bg-base-100 shadow-lg"
          onClick={handleDirections}
          title="Get directions"
        >
          <svg className="w-4 h-4" fill="none" viewBox="0 0 24 24" stroke="currentColor">
            <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M9 20l-5.447-2.724A1 1 0 013 16.382V5.618a1 1 0 011.447-.894L9 7m0 13l6-3m-6 3V7m6 10l4.553 2.276A1 1 0 0021 18.382V7.618a1 1 0 00-.553-.894L15 4m0 13V4m0 0L9 7" />
          </svg>
        </button>
      </div>

      {/* Delivery radius indicator (if applicable) */}
      {storefront.delivery_options && storefront.delivery_options.length > 0 && (
        <div className="absolute top-4 left-4 z-[1000] bg-base-100 rounded-lg shadow-lg p-3">
          <p className="text-sm font-semibold">Delivery Available</p>
          <p className="text-xs text-base-content/60">
            {storefront.delivery_options.find(opt => opt.is_enabled)?.name || 'Available'}
          </p>
        </div>
      )}
    </div>
  );
}