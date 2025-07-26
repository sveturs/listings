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
    const exifData = await parse(file, ['gps', 'exif', 'ifd0']);

    if (!exifData) {
      console.log('No EXIF data found in image');
      return null;
    }

    const result: ExifData = {};

    // Извлекаем координаты GPS
    if (exifData.latitude && exifData.longitude) {
      result.location = {
        latitude: exifData.latitude,
        longitude: exifData.longitude,
      };

      if (exifData.altitude !== undefined) {
        result.location.altitude = exifData.altitude;
      }

      console.log('GPS coordinates found in EXIF:', result.location);
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
