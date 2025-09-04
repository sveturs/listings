// Service Worker для Progressive Web App
// День 25: PWA Features Implementation

const CACHE_NAME = 'svetu-v1.0.0';
const RUNTIME_CACHE = 'svetu-runtime';
const IMAGE_CACHE = 'svetu-images';
const API_CACHE = 'svetu-api';

// Статические ресурсы для кеширования
const STATIC_ASSETS = [
  '/',
  '/offline.html',
  '/manifest.json',
  '/icons/icon-192x192.png',
  '/icons/icon-512x512.png',
];

// API endpoints для кеширования
const API_ROUTES = [
  '/api/v1/marketplace/categories',
  '/api/v2/attributes',
];

// Стратегии кеширования
const CACHE_STRATEGIES = {
  staleWhileRevalidate: async (request) => {
    const cache = await caches.open(RUNTIME_CACHE);
    const cachedResponse = await cache.match(request);
    
    const fetchPromise = fetch(request).then(response => {
      if (response.ok) {
        cache.put(request, response.clone());
      }
      return response;
    });
    
    return cachedResponse || fetchPromise;
  },
  
  cacheFirst: async (request) => {
    const cache = await caches.open(RUNTIME_CACHE);
    const cachedResponse = await cache.match(request);
    
    if (cachedResponse) {
      return cachedResponse;
    }
    
    const response = await fetch(request);
    if (response.ok) {
      cache.put(request, response.clone());
    }
    return response;
  },
  
  networkFirst: async (request) => {
    try {
      const response = await fetch(request);
      if (response.ok) {
        const cache = await caches.open(RUNTIME_CACHE);
        cache.put(request, response.clone());
      }
      return response;
    } catch (error) {
      const cache = await caches.open(RUNTIME_CACHE);
      const cachedResponse = await cache.match(request);
      if (cachedResponse) {
        return cachedResponse;
      }
      throw error;
    }
  },
  
  networkOnly: async (request) => {
    return fetch(request);
  },
};

// Install event - кеширование статических ресурсов
self.addEventListener('install', (event) => {
  console.log('[SW] Installing service worker...');
  
  event.waitUntil(
    caches.open(CACHE_NAME).then(cache => {
      console.log('[SW] Caching static assets');
      return cache.addAll(STATIC_ASSETS);
    }).then(() => {
      console.log('[SW] Skip waiting');
      return self.skipWaiting();
    })
  );
});

// Activate event - очистка старых кешей
self.addEventListener('activate', (event) => {
  console.log('[SW] Activating service worker...');
  
  event.waitUntil(
    caches.keys().then(cacheNames => {
      return Promise.all(
        cacheNames
          .filter(cacheName => {
            return cacheName !== CACHE_NAME && 
                   cacheName !== RUNTIME_CACHE &&
                   cacheName !== IMAGE_CACHE &&
                   cacheName !== API_CACHE;
          })
          .map(cacheName => {
            console.log('[SW] Deleting old cache:', cacheName);
            return caches.delete(cacheName);
          })
      );
    }).then(() => {
      console.log('[SW] Claiming clients');
      return self.clients.claim();
    })
  );
});

// Fetch event - обработка запросов
self.addEventListener('fetch', (event) => {
  const { request } = event;
  const url = new URL(request.url);
  
  // Пропускаем non-GET запросы
  if (request.method !== 'GET') {
    return;
  }
  
  // Обработка изображений
  if (request.destination === 'image' || /\.(jpg|jpeg|png|gif|webp|svg)$/i.test(url.pathname)) {
    event.respondWith(
      caches.open(IMAGE_CACHE).then(cache => {
        return cache.match(request).then(cachedResponse => {
          if (cachedResponse) {
            return cachedResponse;
          }
          
          return fetch(request).then(response => {
            if (response.ok) {
              cache.put(request, response.clone());
            }
            return response;
          }).catch(() => {
            // Возвращаем placeholder изображение при ошибке
            return cache.match('/icons/placeholder.png');
          });
        });
      })
    );
    return;
  }
  
  // Обработка API запросов
  if (url.pathname.startsWith('/api/')) {
    event.respondWith(
      handleAPIRequest(request)
    );
    return;
  }
  
  // Обработка навигационных запросов (HTML страниц)
  if (request.mode === 'navigate') {
    event.respondWith(
      CACHE_STRATEGIES.networkFirst(request).catch(() => {
        return caches.match('/offline.html');
      })
    );
    return;
  }
  
  // Обработка остальных запросов (CSS, JS и т.д.)
  event.respondWith(
    CACHE_STRATEGIES.staleWhileRevalidate(request)
  );
});

// Обработка API запросов с умным кешированием
async function handleAPIRequest(request) {
  const cache = await caches.open(API_CACHE);
  const url = new URL(request.url);
  
  // Для критических API endpoints используем stale-while-revalidate
  if (API_ROUTES.some(route => url.pathname.includes(route))) {
    const cachedResponse = await cache.match(request);
    
    // Возвращаем кешированный ответ и обновляем в фоне
    if (cachedResponse) {
      // Обновляем кеш в фоне
      fetch(request).then(response => {
        if (response.ok) {
          cache.put(request, response.clone());
        }
      });
      
      return cachedResponse;
    }
  }
  
  // Для остальных API запросов - network first
  try {
    const response = await fetch(request);
    
    // Кешируем успешные ответы
    if (response.ok) {
      // Добавляем метку времени в заголовки
      const responseWithTimestamp = new Response(response.body, {
        status: response.status,
        statusText: response.statusText,
        headers: new Headers({
          ...response.headers,
          'sw-cache-timestamp': new Date().toISOString(),
        }),
      });
      
      cache.put(request, responseWithTimestamp.clone());
      return response;
    }
    
    return response;
  } catch (error) {
    // При ошибке сети пытаемся вернуть кешированный ответ
    const cachedResponse = await cache.match(request);
    if (cachedResponse) {
      // Добавляем заголовок, что данные из кеша
      return new Response(cachedResponse.body, {
        status: cachedResponse.status,
        statusText: cachedResponse.statusText,
        headers: new Headers({
          ...cachedResponse.headers,
          'X-From-Cache': 'true',
        }),
      });
    }
    
    // Если нет кеша, возвращаем ошибку
    return new Response(JSON.stringify({
      error: 'Network error and no cached data available',
    }), {
      status: 503,
      headers: { 'Content-Type': 'application/json' },
    });
  }
}

// Background Sync - синхронизация данных в фоне
self.addEventListener('sync', (event) => {
  console.log('[SW] Background sync event:', event.tag);
  
  if (event.tag === 'sync-pending-changes') {
    event.waitUntil(syncPendingChanges());
  }
});

// Синхронизация ожидающих изменений
async function syncPendingChanges() {
  try {
    // Получаем все ожидающие изменения из IndexedDB
    const pendingChanges = await getPendingChanges();
    
    for (const change of pendingChanges) {
      try {
        const response = await fetch(change.endpoint, {
          method: change.method,
          headers: change.headers,
          body: JSON.stringify(change.data),
        });
        
        if (response.ok) {
          // Удаляем успешно синхронизированное изменение
          await removePendingChange(change.id);
        }
      } catch (error) {
        console.error('[SW] Failed to sync change:', change, error);
      }
    }
  } catch (error) {
    console.error('[SW] Background sync failed:', error);
  }
}

// Push уведомления
self.addEventListener('push', (event) => {
  const options = {
    body: event.data ? event.data.text() : 'New update available',
    icon: '/icons/icon-192x192.png',
    badge: '/icons/badge-72x72.png',
    vibrate: [100, 50, 100],
    data: {
      dateOfArrival: Date.now(),
      primaryKey: 1,
    },
    actions: [
      {
        action: 'explore',
        title: 'Open',
        icon: '/icons/checkmark.png',
      },
      {
        action: 'close',
        title: 'Close',
        icon: '/icons/xmark.png',
      },
    ],
  };
  
  event.waitUntil(
    self.registration.showNotification('Sve Tu', options)
  );
});

// Обработка клика по уведомлению
self.addEventListener('notificationclick', (event) => {
  event.notification.close();
  
  if (event.action === 'explore') {
    event.waitUntil(
      clients.openWindow('/')
    );
  }
});

// Периодическая фоновая синхронизация (если поддерживается)
self.addEventListener('periodicsync', (event) => {
  if (event.tag === 'update-cache') {
    event.waitUntil(updateCache());
  }
});

// Обновление критических данных в кеше
async function updateCache() {
  const cache = await caches.open(API_CACHE);
  
  for (const route of API_ROUTES) {
    try {
      const response = await fetch(route);
      if (response.ok) {
        await cache.put(route, response);
      }
    } catch (error) {
      console.error('[SW] Failed to update cache for:', route, error);
    }
  }
}

// Helper функции для работы с IndexedDB
async function getPendingChanges() {
  // Здесь должна быть реализация получения данных из IndexedDB
  return [];
}

async function removePendingChange(id) {
  // Здесь должна быть реализация удаления данных из IndexedDB
}

// Обработка сообщений от клиента
self.addEventListener('message', (event) => {
  if (event.data && event.data.type === 'SKIP_WAITING') {
    self.skipWaiting();
  }
  
  if (event.data && event.data.type === 'CLEAR_CACHE') {
    event.waitUntil(
      caches.keys().then(cacheNames => {
        return Promise.all(
          cacheNames.map(cacheName => caches.delete(cacheName))
        );
      })
    );
  }
});