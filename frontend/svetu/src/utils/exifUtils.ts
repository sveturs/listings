import { parse } from 'exifr';

export interface ExifLocation {
  latitude: number;
  longitude: number;
  altitude?: number;
  accuracy?: number;
}

export interface ExifData {
  location?: ExifLocation;
  timestamp?: Date;
  camera?: {
    make?: string;
    model?: string;
  };
}

/**
 * Извлекает EXIF данные из файла изображения
 */
export async function extractExifData(file: File): Promise<ExifData | null> {
  try {
    // Используем более полную конфигурацию для извлечения GPS
    const exifData = await parse(file, {
      gps: true,
      // Оптимизированная конфигурация для получения всех GPS данных
      translateKeys: true,
      translateValues: true,
      reviveValues: true,
      sanitize: true,
      mergeOutput: true,
    });

    console.log('Raw EXIF data:', exifData);
    console.log('EXIF data keys:', Object.keys(exifData || {}));

    if (!exifData) {
      console.log('No EXIF data found in image');
      return null;
    }

    const result: ExifData = {};

    // Извлекаем координаты GPS - проверяем несколько возможных полей
    if (exifData.latitude && exifData.longitude) {
      result.location = {
        latitude: exifData.latitude,
        longitude: exifData.longitude,
      };

      if (exifData.altitude !== undefined) {
        result.location.altitude = exifData.altitude;
      }

      console.log('GPS coordinates found in EXIF:', result.location);
    } else if (exifData.GPSLatitude && exifData.GPSLongitude) {
      // Альтернативный формат GPS данных
      console.log('Alternative GPS format found:', {
        GPSLatitude: exifData.GPSLatitude,
        GPSLongitude: exifData.GPSLongitude,
        GPSLatitudeRef: exifData.GPSLatitudeRef,
        GPSLongitudeRef: exifData.GPSLongitudeRef,
      });

      // Преобразуем DMS в десятичные градусы если нужно
      const lat = Array.isArray(exifData.GPSLatitude)
        ? exifData.GPSLatitude[0] +
          exifData.GPSLatitude[1] / 60 +
          exifData.GPSLatitude[2] / 3600
        : exifData.GPSLatitude;
      const lon = Array.isArray(exifData.GPSLongitude)
        ? exifData.GPSLongitude[0] +
          exifData.GPSLongitude[1] / 60 +
          exifData.GPSLongitude[2] / 3600
        : exifData.GPSLongitude;

      result.location = {
        latitude: exifData.GPSLatitudeRef === 'S' ? -lat : lat,
        longitude: exifData.GPSLongitudeRef === 'W' ? -lon : lon,
      };

      if (exifData.GPSAltitude !== undefined) {
        result.location.altitude = exifData.GPSAltitude;
      }

      console.log('Converted GPS coordinates:', result.location);
    } else {
      console.log(
        'No GPS data found in EXIF. Available fields:',
        Object.keys(exifData)
      );

      // Логируем все поля, которые могут содержать GPS
      const gpsFields = Object.keys(exifData).filter(
        (key) =>
          key.toLowerCase().includes('gps') ||
          key.toLowerCase().includes('lat') ||
          key.toLowerCase().includes('lon') ||
          key.toLowerCase().includes('location')
      );

      if (gpsFields.length > 0) {
        console.log('Found GPS-related fields:', gpsFields);
        gpsFields.forEach((field) => {
          console.log(`${field}:`, exifData[field]);
        });
      }
    }

    // Извлекаем дату съемки
    if (exifData.DateTimeOriginal || exifData.DateTime) {
      result.timestamp = new Date(
        exifData.DateTimeOriginal || exifData.DateTime
      );
      console.log('Photo timestamp found:', result.timestamp);
    }

    // Извлекаем информацию о камере
    if (exifData.Make || exifData.Model) {
      result.camera = {
        make: exifData.Make,
        model: exifData.Model,
      };
      console.log('Camera info found:', result.camera);
    }

    return result;
  } catch (error) {
    console.error('Error extracting EXIF data:', error);
    return null;
  }
}

/**
 * Извлекает геолокацию из нескольких изображений и возвращает наиболее подходящую
 */
export async function extractLocationFromImages(
  files: File[]
): Promise<ExifLocation | null> {
  const locations: (ExifLocation & { timestamp?: Date })[] = [];

  for (const file of files) {
    const exifData = await extractExifData(file);
    if (exifData?.location) {
      locations.push({
        ...exifData.location,
        timestamp: exifData.timestamp,
      });
    }
  }

  if (locations.length === 0) {
    console.log('No GPS coordinates found in any images');
    return null;
  }

  // Возвращаем самую свежую локацию (или первую, если нет временных меток)
  const sortedLocations = locations.sort((a, b) => {
    if (!a.timestamp || !b.timestamp) return 0;
    return b.timestamp.getTime() - a.timestamp.getTime();
  });

  console.log(
    `Found ${locations.length} locations, using:`,
    sortedLocations[0]
  );
  return sortedLocations[0];
}

/**
 * Проверяет, поддерживает ли браузер EXIF
 */
export function isExifSupported(): boolean {
  return typeof window !== 'undefined' && 'File' in window;
}
