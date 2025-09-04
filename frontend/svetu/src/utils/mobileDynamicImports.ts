import { lazy, ComponentType, LazyExoticComponent } from 'react';

interface DynamicImportOptions {
  prefetch?: boolean;
  preload?: boolean;
  chunkName?: string;
  fallback?: ComponentType<any>;
  retry?: number;
  timeout?: number;
}

interface DeviceDetection {
  isMobile: boolean;
  isTablet: boolean;
  isDesktop: boolean;
  isSlowConnection: boolean;
  isDataSaver: boolean;
  screenSize: 'xs' | 'sm' | 'md' | 'lg' | 'xl';
  deviceMemory?: number;
  cpuCores?: number;
}

/**
 * Detect device characteristics for optimal loading
 */
export function detectDevice(): DeviceDetection {
  const userAgent = navigator.userAgent;
  const screenWidth = window.innerWidth;

  // Device type detection
  const isMobile =
    /Android|webOS|iPhone|iPad|iPod|BlackBerry|IEMobile|Opera Mini/i.test(
      userAgent
    );
  const isTablet = /iPad|Android/i.test(userAgent) && screenWidth >= 768;
  const isDesktop = !isMobile && !isTablet;

  // Connection detection
  const connection =
    (navigator as any).connection ||
    (navigator as any).mozConnection ||
    (navigator as any).webkitConnection;

  const isSlowConnection =
    connection?.effectiveType === '2g' ||
    connection?.effectiveType === 'slow-2g' ||
    connection?.rtt > 400; // RTT > 400ms

  const isDataSaver = connection?.saveData || false;

  // Screen size
  let screenSize: DeviceDetection['screenSize'];
  if (screenWidth < 576) screenSize = 'xs';
  else if (screenWidth < 768) screenSize = 'sm';
  else if (screenWidth < 992) screenSize = 'md';
  else if (screenWidth < 1200) screenSize = 'lg';
  else screenSize = 'xl';

  // Hardware capabilities
  const deviceMemory = (navigator as any).deviceMemory; // GB of RAM
  const cpuCores = navigator.hardwareConcurrency;

  return {
    isMobile,
    isTablet,
    isDesktop,
    isSlowConnection,
    isDataSaver,
    screenSize,
    deviceMemory,
    cpuCores,
  };
}

/**
 * Dynamic import with retry logic and device-specific optimizations
 */
export function dynamicImport<T extends ComponentType<any>>(
  importFunc: () => Promise<{ default: T }>,
  options: DynamicImportOptions = {}
): LazyExoticComponent<T> {
  const {
    prefetch = false,
    preload = false,
    retry = 3,
    timeout = 10000,
  } = options;

  const device = detectDevice();

  // Retry logic wrapper
  const importWithRetry = async (): Promise<{ default: T }> => {
    let lastError: Error | null = null;

    for (let i = 0; i < retry; i++) {
      try {
        const timeoutPromise = new Promise<never>((_, reject) =>
          setTimeout(() => reject(new Error('Import timeout')), timeout)
        );

        const importPromise = importFunc();

        // Race between import and timeout
        const module = await Promise.race([importPromise, timeoutPromise]);
        return module;
      } catch (error) {
        lastError = error as Error;

        // Exponential backoff
        if (i < retry - 1) {
          await new Promise((resolve) =>
            setTimeout(resolve, Math.pow(2, i) * 1000)
          );
        }
      }
    }

    throw lastError || new Error('Failed to import module');
  };

  // Create lazy component
  const LazyComponent = lazy(importWithRetry);

  // Prefetch for desktop or fast connections
  if (prefetch && (device.isDesktop || !device.isSlowConnection)) {
    // Start loading in background
    requestIdleCallback(() => {
      importFunc().catch(() => {
        // Silent fail for prefetch
      });
    });
  }

  // Preload for critical components
  if (preload && !device.isSlowConnection && !device.isDataSaver) {
    // Start loading immediately but don't wait
    importFunc().catch(() => {
      // Silent fail for preload
    });
  }

  return LazyComponent;
}

/**
 * Mobile-optimized dynamic imports with automatic splitting
 */
export const MobileComponents = {
  // Heavy components - load on demand
  ARViewer: () =>
    dynamicImport(
      () =>
        import(
          /* webpackChunkName: "ar-viewer" */ '../components/shared/ARProductViewer'
        ),
      { prefetch: false, chunkName: 'ar-viewer' }
    ),

  QRScanner: () =>
    dynamicImport(
      () =>
        import(
          /* webpackChunkName: "qr-scanner" */ '../components/shared/QRBarcodeScanner'
        ),
      { prefetch: false, chunkName: 'qr-scanner' }
    ),

  // Medium components - prefetch on desktop
  MobileAttributeSelector: () =>
    dynamicImport(
      () =>
        import(
          /* webpackChunkName: "mobile-attrs" */ '../components/shared/MobileAttributeSelector'
        ),
      { prefetch: !detectDevice().isMobile, chunkName: 'mobile-attrs' }
    ),

  // Light components - preload
  TouchGestures: () =>
    dynamicImport(
      () =>
        import(/* webpackChunkName: "gestures" */ '../hooks/useTouchGestures'),
      { preload: true, chunkName: 'gestures' }
    ),

  VoiceSearch: () =>
    dynamicImport(
      () => import(/* webpackChunkName: "voice" */ '../hooks/useVoiceSearch'),
      { prefetch: true, chunkName: 'voice' }
    ),
};

/**
 * Route-based code splitting configuration
 */
export const RouteChunks = {
  // Critical routes - included in main bundle
  critical: ['/'], // '/login', '/register'

  // Primary routes - preload
  primary: ['/search', '/product/[id]', '/create'],

  // Secondary routes - prefetch
  secondary: ['/profile', '/settings', '/messages'],

  // Heavy routes - load on demand
  heavy: ['/ar-preview', '/scan', '/analytics'],

  // Mobile-specific routes
  mobile: ['/m/search', '/m/camera', '/m/quick-add'],
};

/**
 * Intelligent route preloading based on user behavior
 */
export class RoutePreloader {
  private static preloadedRoutes = new Set<string>();
  private static userPatterns: Map<string, number> = new Map();

  static analyzeUserBehavior(currentRoute: string) {
    // Track route visits
    const visits = this.userPatterns.get(currentRoute) || 0;
    this.userPatterns.set(currentRoute, visits + 1);

    // Predict next likely routes
    const predictions = this.predictNextRoutes(currentRoute);

    // Preload predicted routes
    predictions.forEach((route) => this.preloadRoute(route));
  }

  static predictNextRoutes(currentRoute: string): string[] {
    const predictions: string[] = [];
    const device = detectDevice();

    // Route-specific predictions
    if (currentRoute === '/') {
      predictions.push('/search', '/categories');
    } else if (currentRoute.startsWith('/product/')) {
      predictions.push('/create', '/messages');
    } else if (currentRoute === '/search') {
      predictions.push('/product/[id]', '/filters');
    }

    // Device-specific predictions
    if (device.isMobile) {
      predictions.push('/m/quick-add', '/m/camera');
    }

    // Popular routes based on history
    const popularRoutes = Array.from(this.userPatterns.entries())
      .sort(([, a], [, b]) => b - a)
      .slice(0, 3)
      .map(([route]) => route);

    predictions.push(...popularRoutes);

    return [...new Set(predictions)]; // Remove duplicates
  }

  static preloadRoute(route: string) {
    if (this.preloadedRoutes.has(route)) return;

    const device = detectDevice();

    // Skip preloading on slow connections
    if (device.isSlowConnection || device.isDataSaver) return;

    // Mark as preloaded
    this.preloadedRoutes.add(route);

    // Dynamic import based on route
    requestIdleCallback(() => {
      this.loadRouteChunk(route);
    });
  }

  private static loadRouteChunk(route: string) {
    // Route to chunk mapping
    const chunkMap: Record<string, () => Promise<any>> = {
      '/search': () =>
        import(/* webpackChunkName: "search" */ '../app/[locale]/search/page'),
      '/create': () =>
        import(
          /* webpackChunkName: "create" */ '../app/[locale]/create-listing/page'
        ),
      '/profile': () =>
        import(
          /* webpackChunkName: "profile" */ '../app/[locale]/profile/page'
        ),
      '/ar-preview': () =>
        import(
          /* webpackChunkName: "ar" */ '../components/shared/ARProductViewer'
        ),
      '/scan': () =>
        import(
          /* webpackChunkName: "scanner" */ '../components/shared/QRBarcodeScanner'
        ),
    };

    const loader = chunkMap[route];
    if (loader) {
      loader().catch(() => {
        // Silent fail for preload
      });
    }
  }
}

/**
 * Bundle size optimization utilities
 */
export const BundleOptimizer = {
  /**
   * Remove unused code for mobile builds
   */
  stripDesktopCode(code: string): string {
    const device = detectDevice();

    if (device.isMobile) {
      // Remove desktop-only code blocks
      code = code.replace(
        /\/\* desktop-only-start \*\/[\s\S]*?\/\* desktop-only-end \*\//g,
        ''
      );
    }

    return code;
  },

  /**
   * Load polyfills only when needed
   */
  async loadPolyfills() {
    const needed: string[] = [];

    // Check for missing features
    if (!window.IntersectionObserver) {
      needed.push('intersection-observer');
    }

    if (!window.ResizeObserver) {
      needed.push('resize-observer');
    }

    if (!Element.prototype.closest) {
      needed.push('element-closest');
    }

    // Load only needed polyfills
    if (needed.length > 0) {
      await import(
        /* webpackChunkName: "polyfills" */ `core-js/features/${needed.join(',')}`
      );
    }
  },

  /**
   * Tree-shake unused translations
   */
  optimizeTranslations(locale: string): Promise<any> {
    const device = detectDevice();

    // Load minimal translations for mobile
    if (device.isMobile && device.isSlowConnection) {
      return import(
        /* webpackChunkName: "i18n-mobile" */ `../messages/${locale}/mobile.json`
      );
    }

    // Load full translations
    return import(
      /* webpackChunkName: "i18n-full" */ `../messages/${locale}/common.json`
    );
  },
};

/**
 * Webpack magic comments for optimization
 */
export const WebpackHints = {
  // Prefetch - load when browser is idle
  prefetch: (chunkName: string) =>
    `/* webpackPrefetch: true, webpackChunkName: "${chunkName}" */`,

  // Preload - load in parallel with parent
  preload: (chunkName: string) =>
    `/* webpackPreload: true, webpackChunkName: "${chunkName}" */`,

  // Lazy mode
  lazy: (chunkName: string) =>
    `/* webpackMode: "lazy", webpackChunkName: "${chunkName}" */`,

  // Lazy once - single shared instance
  lazyOnce: (chunkName: string) =>
    `/* webpackMode: "lazy-once", webpackChunkName: "${chunkName}" */`,

  // Eager - no separate chunk
  eager: () => `/* webpackMode: "eager" */`,

  // Weak - don't follow import
  weak: () => `/* webpackMode: "weak" */`,
};
