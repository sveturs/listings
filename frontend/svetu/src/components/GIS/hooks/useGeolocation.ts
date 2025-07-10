import { useState, useEffect, useCallback } from 'react';
import { GeoLocation } from '../types/gis';

interface UseGeolocationOptions {
  enableHighAccuracy?: boolean;
  timeout?: number;
  maximumAge?: number;
  watch?: boolean;
}

interface UseGeolocationResult {
  location: GeoLocation | null;
  error: GeolocationPositionError | null;
  loading: boolean;
  supported: boolean;
  getCurrentPosition: () => Promise<GeoLocation>;
  watchPosition: () => void;
  clearWatch: () => void;
}

export const useGeolocation = (
  options: UseGeolocationOptions = {}
): UseGeolocationResult => {
  const [location, setLocation] = useState<GeoLocation | null>(null);
  const [error, setError] = useState<GeolocationPositionError | null>(null);
  const [loading, setLoading] = useState(false);
  const [watchId, setWatchId] = useState<number | null>(null);

  const {
    enableHighAccuracy = true,
    timeout = 15000,
    maximumAge = 600000, // 10 minutes
    watch = false,
  } = options;

  const supported =
    typeof navigator !== 'undefined' && 'geolocation' in navigator;

  const convertPosition = useCallback(
    (position: GeolocationPosition): GeoLocation => {
      return {
        latitude: position.coords.latitude,
        longitude: position.coords.longitude,
        accuracy: position.coords.accuracy,
        altitude: position.coords.altitude || undefined,
        altitudeAccuracy: position.coords.altitudeAccuracy || undefined,
        heading: position.coords.heading || undefined,
        speed: position.coords.speed || undefined,
      };
    },
    []
  );

  const handleSuccess = useCallback(
    (position: GeolocationPosition) => {
      const geoLocation = convertPosition(position);
      setLocation(geoLocation);
      setError(null);
      setLoading(false);
    },
    [convertPosition]
  );

  const handleError = useCallback((error: GeolocationPositionError) => {
    setError(error);
    setLoading(false);
  }, []);

  const getCurrentPosition = useCallback((): Promise<GeoLocation> => {
    return new Promise((resolve, reject) => {
      if (!supported) {
        reject(new Error('Geolocation is not supported'));
        return;
      }

      setLoading(true);
      setError(null);

      navigator.geolocation.getCurrentPosition(
        (position) => {
          const geoLocation = convertPosition(position);
          setLocation(geoLocation);
          setLoading(false);
          resolve(geoLocation);
        },
        (error) => {
          setError(error);
          setLoading(false);
          reject(error);
        },
        {
          enableHighAccuracy,
          timeout,
          maximumAge,
        }
      );
    });
  }, [supported, enableHighAccuracy, timeout, maximumAge, convertPosition]);

  const watchPosition = useCallback(() => {
    if (!supported) return;

    if (watchId !== null) {
      navigator.geolocation.clearWatch(watchId);
    }

    setLoading(true);
    setError(null);

    const id = navigator.geolocation.watchPosition(handleSuccess, handleError, {
      enableHighAccuracy,
      timeout,
      maximumAge,
    });

    setWatchId(id);
  }, [
    supported,
    watchId,
    handleSuccess,
    handleError,
    enableHighAccuracy,
    timeout,
    maximumAge,
  ]);

  const clearWatch = useCallback(() => {
    if (watchId !== null) {
      navigator.geolocation.clearWatch(watchId);
      setWatchId(null);
    }
  }, [watchId]);

  // Auto-start watching if enabled
  useEffect(() => {
    if (watch) {
      watchPosition();
    }

    return () => {
      clearWatch();
    };
  }, [watch, watchPosition, clearWatch]);

  return {
    location,
    error,
    loading,
    supported,
    getCurrentPosition,
    watchPosition,
    clearWatch,
  };
};

// Утилитные функции для работы с геолокацией
export const getLocationErrorMessage = (
  error: GeolocationPositionError
): string => {
  switch (error.code) {
    case error.PERMISSION_DENIED:
      return 'geolocation.permission_denied';
    case error.POSITION_UNAVAILABLE:
      return 'geolocation.position_unavailable';
    case error.TIMEOUT:
      return 'geolocation.timeout';
    default:
      return 'geolocation.unknown_error';
  }
};

export const calculateDistance = (
  lat1: number,
  lon1: number,
  lat2: number,
  lon2: number
): number => {
  const R = 6371; // Радиус Земли в км
  const dLat = ((lat2 - lat1) * Math.PI) / 180;
  const dLon = ((lon2 - lon1) * Math.PI) / 180;
  const a =
    Math.sin(dLat / 2) * Math.sin(dLat / 2) +
    Math.cos((lat1 * Math.PI) / 180) *
      Math.cos((lat2 * Math.PI) / 180) *
      Math.sin(dLon / 2) *
      Math.sin(dLon / 2);
  const c = 2 * Math.atan2(Math.sqrt(a), Math.sqrt(1 - a));
  const distance = R * c;
  return distance;
};

export const calculateBearing = (
  lat1: number,
  lon1: number,
  lat2: number,
  lon2: number
): number => {
  const dLon = ((lon2 - lon1) * Math.PI) / 180;
  const lat1Rad = (lat1 * Math.PI) / 180;
  const lat2Rad = (lat2 * Math.PI) / 180;

  const y = Math.sin(dLon) * Math.cos(lat2Rad);
  const x =
    Math.cos(lat1Rad) * Math.sin(lat2Rad) -
    Math.sin(lat1Rad) * Math.cos(lat2Rad) * Math.cos(dLon);

  const bearing = (Math.atan2(y, x) * 180) / Math.PI;
  return (bearing + 360) % 360;
};

export const formatDistance = (
  distance: number,
  unit: 'km' | 'm' = 'km'
): string => {
  if (unit === 'm') {
    return distance < 1000
      ? `${Math.round(distance)} м`
      : `${(distance / 1000).toFixed(1)} км`;
  }
  return distance < 1
    ? `${Math.round(distance * 1000)} м`
    : `${distance.toFixed(1)} км`;
};

export const isLocationInBounds = (
  location: GeoLocation,
  bounds: {
    north: number;
    south: number;
    east: number;
    west: number;
  }
): boolean => {
  return (
    location.latitude >= bounds.south &&
    location.latitude <= bounds.north &&
    location.longitude >= bounds.west &&
    location.longitude <= bounds.east
  );
};
