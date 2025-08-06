/**
 * Service Worker для кэширования модулей переводов
 * Обеспечивает offline доступ к переводам
 */

const CACHE_NAME = 'translations-v1';
const TRANSLATION_PATTERN = /\/_next\/static\/chunks\/src_messages_.*\.js/;
const JSON_PATTERN = /\/messages\/[a-z]{2}\/[a-zA-Z]+\.json/;

// Активация Service Worker
self.addEventListener('install', (event) => {
  console.log('[SW] Installing translations service worker');
  self.skipWaiting();
});

self.addEventListener('activate', (event) => {
  console.log('[SW] Activating translations service worker');
  event.waitUntil(
    caches.keys().then((cacheNames) => {
      return Promise.all(
        cacheNames
          .filter((cacheName) => cacheName.startsWith('translations-') && cacheName !== CACHE_NAME)
          .map((cacheName) => caches.delete(cacheName))
      );
    })
  );
  self.clients.claim();
});

// Перехват запросов
self.addEventListener('fetch', (event) => {
  const { request } = event;
  const url = new URL(request.url);
  
  // Кэшируем только модули переводов
  if (TRANSLATION_PATTERN.test(url.pathname) || JSON_PATTERN.test(url.pathname)) {
    event.respondWith(
      caches.match(request).then((cachedResponse) => {
        if (cachedResponse) {
          console.log('[SW] Serving from cache:', url.pathname);
          return cachedResponse;
        }
        
        console.log('[SW] Fetching:', url.pathname);
        return fetch(request).then((response) => {
          // Проверяем, что ответ успешный
          if (!response || response.status !== 200) {
            return response;
          }
          
          // Клонируем ответ для кэша
          const responseToCache = response.clone();
          
          caches.open(CACHE_NAME).then((cache) => {
            cache.put(request, responseToCache);
          });
          
          return response;
        });
      })
    );
  }
});

// Предзагрузка критичных модулей
self.addEventListener('message', (event) => {
  if (event.data && event.data.type === 'PRELOAD_MODULES') {
    const { locale, modules } = event.data;
    
    event.waitUntil(
      caches.open(CACHE_NAME).then(async (cache) => {
        const urls = modules.map(module => 
          `/_next/static/chunks/src_messages_${locale}_${module}.js`
        );
        
        console.log('[SW] Preloading modules:', urls);
        
        try {
          await cache.addAll(urls);
          self.clients.matchAll().then(clients => {
            clients.forEach(client => {
              client.postMessage({
                type: 'MODULES_PRELOADED',
                modules: modules
              });
            });
          });
        } catch (error) {
          console.error('[SW] Failed to preload modules:', error);
        }
      })
    );
  }
});