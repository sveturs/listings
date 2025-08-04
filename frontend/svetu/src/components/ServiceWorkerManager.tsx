'use client';

import { useEffect } from 'react';
import { useLocale } from 'next-intl';

export function ServiceWorkerManager() {
  const locale = useLocale();
  
  useEffect(() => {
    // Регистрируем Service Worker только в production
    if (
      process.env.NODE_ENV === 'production' && 
      'serviceWorker' in navigator &&
      process.env.USE_MODULAR_I18N === 'true'
    ) {
      navigator.serviceWorker
        .register('/sw-translations.js')
        .then((registration) => {
          console.log('SW registered:', registration);
          
          // Предзагружаем базовые модули
          if (registration.active) {
            registration.active.postMessage({
              type: 'PRELOAD_MODULES',
              locale,
              modules: ['common', 'marketplace', 'auth']
            });
          }
        })
        .catch((error) => {
          console.error('SW registration failed:', error);
        });
      
      // Слушаем сообщения от Service Worker
      navigator.serviceWorker.addEventListener('message', (event) => {
        if (event.data && event.data.type === 'MODULES_PRELOADED') {
          console.log('Modules preloaded:', event.data.modules);
        }
      });
    }
  }, [locale]);
  
  return null;
}