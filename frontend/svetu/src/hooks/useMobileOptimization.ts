import { useState, useEffect, useCallback, useRef } from 'react';

interface MobileOptimizationOptions {
  // Дебаунс для карты (больше на мобильных для экономии ресурсов)
  mapDebounceTime?: number;
  // Лимит маркеров для отображения
  maxMarkersCount?: number;
  // Уменьшенное качество изображений на мобильных
  imageQuality?: 'low' | 'medium' | 'high';
  // Включить ленивую загрузку списков
  enableLazyLoading?: boolean;
  // Размер страницы для пагинации
  pageSize?: number;
  // Интервал для throttling scroll событий
  scrollThrottleTime?: number;
}

interface MobileOptimizationReturn {
  // Детекция мобильного устройства
  isMobile: boolean;
  isTablet: boolean;

  // Информация о производительности устройства
  deviceInfo: {
    isLowEndDevice: boolean;
    memory: number | null;
    cores: number;
    connectionType: string | null;
  };

  // Оптимизированные настройки
  settings: {
    mapDebounceTime: number;
    maxMarkersCount: number;
    imageQuality: 'low' | 'medium' | 'high';
    enableLazyLoading: boolean;
    pageSize: number;
    scrollThrottleTime: number;
  };

  // Утилиты для оптимизации
  optimizeImageUrl: (url: string, width?: number, height?: number) => string;
  shouldRenderItem: (
    index: number,
    viewportStart: number,
    viewportEnd: number
  ) => boolean;
  throttledScrollHandler: (callback: () => void) => () => void;
}

const useMobileOptimization = (
  options: MobileOptimizationOptions = {}
): MobileOptimizationReturn => {
  const [isMobile, setIsMobile] = useState(false);
  const [isTablet, setIsTablet] = useState(false);
  const [deviceInfo, setDeviceInfo] = useState({
    isLowEndDevice: false,
    memory: null as number | null,
    cores: 0,
    connectionType: null as string | null,
  });

  const scrollTimeoutRef = useRef<NodeJS.Timeout | null>(null);

  // Детекция типа устройства и его характеристик
  useEffect(() => {
    const detectDevice = () => {
      const width = window.innerWidth;
      const _height = window.innerHeight;

      setIsMobile(width < 768);
      setIsTablet(width >= 768 && width < 1024);

      // Определяем производительность устройства
      const navigator = window.navigator as any;

      // Память устройства (если доступно)
      const memory = navigator.deviceMemory || null;

      // Количество ядер процессора
      const cores = navigator.hardwareConcurrency || 4;

      // Тип подключения
      const connection =
        navigator.connection ||
        navigator.mozConnection ||
        navigator.webkitConnection;
      const connectionType = connection?.effectiveType || null;

      // Определяем слабое устройство
      const isLowEndDevice =
        (memory && memory <= 2) || // Менее 2GB RAM
        cores <= 2 || // Менее 2 ядер
        connectionType === 'slow-2g' ||
        connectionType === '2g' || // Медленное соединение
        width < 400; // Очень маленький экран

      setDeviceInfo({
        isLowEndDevice,
        memory,
        cores,
        connectionType,
      });
    };

    detectDevice();
    window.addEventListener('resize', detectDevice);

    return () => window.removeEventListener('resize', detectDevice);
  }, []);

  // Определяем оптимизированные настройки на основе устройства
  const settings = {
    mapDebounceTime:
      options.mapDebounceTime ||
      (deviceInfo.isLowEndDevice ? 800 : isMobile ? 500 : 300),

    maxMarkersCount:
      options.maxMarkersCount ||
      (deviceInfo.isLowEndDevice ? 50 : isMobile ? 100 : 200),

    imageQuality:
      options.imageQuality ||
      ((deviceInfo.isLowEndDevice ? 'low' : isMobile ? 'medium' : 'high') as
        | 'low'
        | 'medium'
        | 'high'),

    enableLazyLoading:
      options.enableLazyLoading !== undefined
        ? options.enableLazyLoading
        : isMobile,

    pageSize:
      options.pageSize || (deviceInfo.isLowEndDevice ? 10 : isMobile ? 20 : 30),

    scrollThrottleTime:
      options.scrollThrottleTime || (deviceInfo.isLowEndDevice ? 200 : 100),
  };

  // Оптимизация URL изображений
  const optimizeImageUrl = useCallback(
    (url: string, width?: number, height?: number): string => {
      if (!url) return url;

      // Определяем размеры на основе качества и устройства
      let targetWidth = width;
      const targetHeight = height;

      if (!targetWidth) {
        switch (settings.imageQuality) {
          case 'low':
            targetWidth = isMobile ? 150 : 200;
            break;
          case 'medium':
            targetWidth = isMobile ? 300 : 400;
            break;
          case 'high':
            targetWidth = isMobile ? 600 : 800;
            break;
        }
      }

      // Если URL уже содержит параметры размера, заменяем их
      if (url.includes('?w=') || url.includes('&w=')) {
        return url.replace(/([?&])w=\d+/, `$1w=${targetWidth}`);
      }

      // Добавляем параметры размера
      const separator = url.includes('?') ? '&' : '?';
      return `${url}${separator}w=${targetWidth}${targetHeight ? `&h=${targetHeight}` : ''}&q=${settings.imageQuality === 'low' ? 60 : settings.imageQuality === 'medium' ? 80 : 90}`;
    },
    [settings.imageQuality, isMobile]
  );

  // Определение видимости элементов для виртуализации
  const shouldRenderItem = useCallback(
    (index: number, viewportStart: number, viewportEnd: number): boolean => {
      if (!settings.enableLazyLoading) return true;

      // Добавляем буфер для предзагрузки
      const buffer = isMobile ? 5 : 10;
      return index >= viewportStart - buffer && index <= viewportEnd + buffer;
    },
    [settings.enableLazyLoading, isMobile]
  );

  // Throttled обработчик скролла
  const throttledScrollHandler = useCallback(
    (callback: () => void) => {
      return () => {
        if (scrollTimeoutRef.current) {
          clearTimeout(scrollTimeoutRef.current);
        }

        scrollTimeoutRef.current = setTimeout(() => {
          callback();
        }, settings.scrollThrottleTime);
      };
    },
    [settings.scrollThrottleTime]
  );

  // Очистка таймеров при размонтировании
  useEffect(() => {
    return () => {
      if (scrollTimeoutRef.current) {
        clearTimeout(scrollTimeoutRef.current);
      }
    };
  }, []);

  return {
    isMobile,
    isTablet,
    deviceInfo,
    settings,
    optimizeImageUrl,
    shouldRenderItem,
    throttledScrollHandler,
  };
};

export default useMobileOptimization;
