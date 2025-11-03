/**
 * Централизованное хранилище API эндпоинтов
 * Обеспечивает типобезопасность и автодополнение
 */

export const API_ENDPOINTS = {
  // Auth endpoints - обновлено для Auth Service микросервиса
  auth: {
    login: '/api/v1/auth/login',
    register: '/api/v1/auth/register',
    logout: '/api/v1/auth/logout',
    refresh: '/api/v1/auth/refresh',
    validate: '/api/v1/auth/validate',
    session: '/api/v1/auth/session',
    google: '/api/v1/auth/google',
    googleCallback: '/api/v1/auth/google/callback',
  },

  // User endpoints
  users: {
    profile: '/api/v1/users/profile',
    update: '/api/v1/users/update',
    delete: '/api/v1/users/delete',
    checkEmail: '/api/v1/users/check-email',
  },

  // Chat endpoints
  chat: {
    list: '/api/v1/chat/chats',
    messages: '/api/v1/chat/messages',
    send: '/api/v1/chat/send',
    read: '/api/v1/chat/read',
    unreadCount: '/api/v1/chat/unread-count',
    archive: (id: number) => `/api/v1/chat/${id}/archive`,
    attachments: {
      upload: (messageId: number) =>
        `/api/v1/chat/messages/${messageId}/attachments`,
      delete: (id: number) => `/api/v1/chat/attachments/${id}`,
    },
  },

  // Marketplace endpoints (unified c2c/b2c)
  marketplace: {
    listings: '/api/v1/marketplace/listings',
    listing: (id: number) => `/api/v1/marketplace/listings/${id}`,
    create: '/api/v1/marketplace/listings/create',
    update: (id: number) => `/api/v1/marketplace/listings/${id}`,
    delete: (id: number) => `/api/v1/marketplace/listings/${id}`,
    search: '/api/v1/marketplace/search',
    categories: '/api/v1/marketplace/categories',
    popularCategories: '/api/v1/marketplace/categories?popular=true&limit=8',
    favorites: {
      list: '/api/v1/marketplace/favorites',
      add: (id: number) => `/api/v1/marketplace/favorites/${id}`,
      remove: (id: number) => `/api/v1/marketplace/favorites/${id}`,
    },
  },

  // Storefronts endpoints (former b2c)
  storefronts: {
    list: '/api/v1/storefronts',
    detail: (slug: string) => `/api/v1/storefronts/${slug}`,
    products: (slug: string) => `/api/v1/storefronts/${slug}/products`,
  },

  // Contacts endpoints
  contacts: {
    list: '/api/v1/contacts',
    add: '/api/v1/contacts',
    update: (id: number) => `/api/v1/contacts/${id}`,
    delete: (id: number) => `/api/v1/contacts/${id}`,
    privacy: '/api/v1/contacts/privacy',
  },

  // WebSocket endpoints
  ws: {
    chat: '/ws/chat',
  },
} as const;

// Тип для извлечения всех эндпоинтов
export type ApiEndpoint =
  | (typeof API_ENDPOINTS.auth)[keyof typeof API_ENDPOINTS.auth]
  | (typeof API_ENDPOINTS.users)[keyof typeof API_ENDPOINTS.users]
  | typeof API_ENDPOINTS.chat.list
  | typeof API_ENDPOINTS.chat.messages
  | typeof API_ENDPOINTS.chat.send
  | typeof API_ENDPOINTS.chat.read
  | typeof API_ENDPOINTS.chat.unreadCount
  | typeof API_ENDPOINTS.marketplace.listings
  | typeof API_ENDPOINTS.marketplace.create
  | typeof API_ENDPOINTS.marketplace.search
  | typeof API_ENDPOINTS.marketplace.categories
  | typeof API_ENDPOINTS.contacts.list
  | typeof API_ENDPOINTS.contacts.add
  | typeof API_ENDPOINTS.contacts.privacy
  | typeof API_ENDPOINTS.ws.chat;

// Хелпер для построения полного URL
export function buildApiUrl(
  endpoint: string,
  params?: Record<string, string | number>
): string {
  let url = endpoint;

  // Заменяем параметры в URL (например, :id)
  if (params) {
    Object.entries(params).forEach(([key, value]) => {
      url = url.replace(`:${key}`, String(value));
    });
  }

  return url;
}
