'use client';

import { useEffect, useState } from 'react';
import { useTranslations } from 'next-intl';
import { publicEnv } from '@/utils/env';

interface NeighborhoodStats {
  total_listings: number;
  new_today: number;
  within_radius: number;
  radius_km: number;
  center_lat?: number;
  center_lon?: number;
}

export default function NearbyStats() {
  const t = useTranslations('misc');
  const [stats, setStats] = useState<NeighborhoodStats | null>(null);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState<string | null>(null);

  useEffect(() => {
    fetchStats();
  }, []);

  const fetchStats = async () => {
    try {
      console.log('[NearbyStats] Starting to fetch stats...');
      setLoading(true);

      // Try to get user's location
      let lat = 44.8176; // Default Belgrade
      let lon = 20.4633;

      if (navigator.geolocation) {
        try {
          const position = await new Promise<GeolocationPosition>(
            (resolve, reject) => {
              navigator.geolocation.getCurrentPosition(resolve, reject, {
                timeout: 5000,
                enableHighAccuracy: false,
              });
            }
          );
          lat = position.coords.latitude;
          lon = position.coords.longitude;
          console.log('[NearbyStats] Got user location:', { lat, lon });
        } catch (geoError) {
          // Use default location if geolocation fails (user denied or not available)
          // Это нормальное поведение - не логируем как ошибку
        }
      }

      // Используем правильный базовый URL для API через утилиту
      const apiUrl = publicEnv.API_URL;
      const url = `${apiUrl}/api/v1/marketplace/neighborhood-stats?lat=${lat}&lon=${lon}&radius=5`;
      console.log('[NearbyStats] API URL:', apiUrl);
      console.log('[NearbyStats] Fetching from URL:', url);

      const response = await fetch(url);
      console.log('[NearbyStats] Response status:', response.status);

      if (!response.ok) {
        throw new Error(`Failed to fetch stats: ${response.status}`);
      }

      const data = await response.json();
      console.log('[NearbyStats] Received data:', data);
      setStats(data.data);
    } catch (err) {
      console.error('[NearbyStats] Error fetching neighborhood stats:', err);
      setError('Failed to load statistics');
    } finally {
      setLoading(false);
    }
  };

  if (loading) {
    return (
      <div className="card bg-base-100">
        <div className="card-body">
          <h3 className="card-title text-lg">{t('home.nearbyTitle')}</h3>
          <div className="animate-pulse">
            <div className="space-y-4">
              <div className="h-12 bg-base-200 rounded"></div>
              <div className="h-12 bg-base-200 rounded"></div>
              <div className="h-12 bg-base-200 rounded"></div>
            </div>
          </div>
        </div>
      </div>
    );
  }

  if (error || !stats) {
    return (
      <div className="card bg-base-100">
        <div className="card-body">
          <h3 className="card-title text-lg">{t('home.nearbyTitle')}</h3>
          <div className="text-base-content/60">{t('home.statsLoadError')}</div>
        </div>
      </div>
    );
  }

  return (
    <div className="card bg-base-100">
      <div className="card-body">
        <h3 className="card-title text-lg">{t('home.nearbyTitle')}</h3>
        <div className="stats stats-vertical">
          <div className="stat px-0">
            <div className="stat-title">{t('home.totalListings')}</div>
            <div className="stat-value text-primary">
              {stats.total_listings.toLocaleString()}
            </div>
          </div>
          <div className="stat px-0">
            <div className="stat-title">{t('home.newToday')}</div>
            <div className="stat-value text-success">+{stats.new_today}</div>
          </div>
          <div className="stat px-0">
            <div className="stat-title">
              {t('home.withinRadius', { radius: stats.radius_km })}
            </div>
            <div className="stat-value">
              {stats.within_radius.toLocaleString()}
            </div>
          </div>
        </div>
      </div>
    </div>
  );
}
