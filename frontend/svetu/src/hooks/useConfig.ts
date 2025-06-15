import { useEffect, useState } from 'react';
import configManager from '@/config';
import { Config } from '@/config/types';

/**
 * React hook для доступа к конфигурации
 * Обеспечивает правильную работу с SSR/CSR
 */
export function useConfig(): Config {
  const [config, setConfig] = useState<Config>(configManager.getConfig());

  useEffect(() => {
    // Обновляем конфигурацию после гидратации
    setConfig(configManager.getConfig());
  }, []);

  return config;
}

/**
 * Hook для проверки доступности функций
 */
export function useFeature(feature: keyof Config['features']): boolean {
  const config = useConfig();
  return config.features[feature];
}
