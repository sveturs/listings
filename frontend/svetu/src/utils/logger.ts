// Frontend logging utilities
const isDev = process.env.NODE_ENV === 'development';

// Configuration for different log levels and modules
const LOG_CONFIG = {
  // Global settings
  enabled: isDev,
  // Module-specific settings (enable/disable per module)
  modules: {
    auth: false,       // Very verbose, disabled by default
    api: false,        // Very verbose, disabled by default  
    search: false,     // Very verbose, disabled by default
    cart: false,       // Disabled by default
    chat: false,       // Disabled by default
    general: false,    // General debug logs, disabled by default
  }
};

export const logger = {
  // General logging
  debug: (...args: unknown[]) => {
    if (LOG_CONFIG.enabled && LOG_CONFIG.modules.general) console.log('[DEBUG]', ...args);
  },
  info: (...args: unknown[]) => {
    if (LOG_CONFIG.enabled) console.info('[INFO]', ...args);
  },
  warn: (...args: unknown[]) => {
    if (LOG_CONFIG.enabled) console.warn('[WARN]', ...args);
  },
  error: (...args: unknown[]) => {
    console.error('[ERROR]', ...args);
  },
  
  // Module-specific loggers
  auth: {
    debug: (...args: unknown[]) => {
      if (LOG_CONFIG.enabled && LOG_CONFIG.modules.auth) console.log('[AUTH]', ...args);
    },
    info: (...args: unknown[]) => {
      if (LOG_CONFIG.enabled && LOG_CONFIG.modules.auth) console.info('[AUTH]', ...args);
    },
    warn: (...args: unknown[]) => {
      if (LOG_CONFIG.enabled) console.warn('[AUTH]', ...args);
    },
    error: (...args: unknown[]) => {
      console.error('[AUTH]', ...args);
    },
  },
  
  api: {
    debug: (...args: unknown[]) => {
      if (LOG_CONFIG.enabled && LOG_CONFIG.modules.api) console.log('[API]', ...args);
    },
    info: (...args: unknown[]) => {
      if (LOG_CONFIG.enabled && LOG_CONFIG.modules.api) console.info('[API]', ...args);
    },
    warn: (...args: unknown[]) => {
      if (LOG_CONFIG.enabled) console.warn('[API]', ...args);
    },
    error: (...args: unknown[]) => {
      console.error('[API]', ...args);
    },
  },

  search: {
    debug: (...args: unknown[]) => {
      if (LOG_CONFIG.enabled && LOG_CONFIG.modules.search) console.log('[SEARCH]', ...args);
    },
    info: (...args: unknown[]) => {
      if (LOG_CONFIG.enabled && LOG_CONFIG.modules.search) console.info('[SEARCH]', ...args);
    },
    warn: (...args: unknown[]) => {
      if (LOG_CONFIG.enabled) console.warn('[SEARCH]', ...args);
    },
    error: (...args: unknown[]) => {
      console.error('[SEARCH]', ...args);
    },
  },

  cart: {
    debug: (...args: unknown[]) => {
      if (LOG_CONFIG.enabled && LOG_CONFIG.modules.cart) console.log('[CART]', ...args);
    },
    info: (...args: unknown[]) => {
      if (LOG_CONFIG.enabled && LOG_CONFIG.modules.cart) console.info('[CART]', ...args);
    },
    warn: (...args: unknown[]) => {
      if (LOG_CONFIG.enabled) console.warn('[CART]', ...args);
    },
    error: (...args: unknown[]) => {
      console.error('[CART]', ...args);
    },
  },

  chat: {
    debug: (...args: unknown[]) => {
      if (LOG_CONFIG.enabled && LOG_CONFIG.modules.chat) console.log('[CHAT]', ...args);
    },
    info: (...args: unknown[]) => {
      if (LOG_CONFIG.enabled && LOG_CONFIG.modules.chat) console.info('[CHAT]', ...args);
    },
    warn: (...args: unknown[]) => {
      if (LOG_CONFIG.enabled) console.warn('[CHAT]', ...args);
    },
    error: (...args: unknown[]) => {
      console.error('[CHAT]', ...args);
    },
  },
};