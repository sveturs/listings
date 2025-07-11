'use client';

import { MapClusterProps } from '../types/gis';

export function MapCluster({ count, onClick }: MapClusterProps) {
  const size = 48 + Math.min(count / 10, 20);

  return (
    <div
      className="relative flex items-center justify-center cursor-pointer transition-all duration-200 hover:scale-110"
      onClick={onClick}
      style={{
        width: `${size}px`,
        height: `${size}px`,
      }}
    >
      <div className="absolute inset-0 bg-primary rounded-full opacity-20 animate-pulse" />
      <div className="relative flex items-center justify-center w-full h-full bg-primary rounded-full shadow-lg border-2 border-white">
        <span className="text-white font-bold text-sm md:text-base">
          {count > 99 ? '99+' : count}
        </span>
      </div>
    </div>
  );
}
