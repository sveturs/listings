interface NetworkInfo {
  effectiveType: '4g' | '3g' | '2g' | 'slow-2g';
  downlink: number;
  rtt: number;
  saveData: boolean;
}

interface BatteryInfo {
  level: number;
  charging: boolean;
  chargingTime: number | null;
  dischargingTime: number | null;
}

interface PerformanceMetrics {
  fps: number;
  memory: {
    used: number;
    limit: number;
    percentage: number;
  };
  cpu: {
    cores: number;
    usage: number;
  };
  network: NetworkInfo;
  battery: BatteryInfo | null;
}

/**
 * Performance Optimization Service
 * Manages battery, network, and resource optimization for mobile devices
 */
class PerformanceOptimizationService {
  private static instance: PerformanceOptimizationService;
  private batteryManager: any = null;
  private networkConnection: any = null;
  private performanceObserver: PerformanceObserver | null = null;
  private metrics: PerformanceMetrics;
  private optimizationLevel: 'aggressive' | 'balanced' | 'performance' =
    'balanced';
  private frameRateMonitor: number | null = null;
  private lastFrameTime: number = 0;
  private frameCount: number = 0;

  private constructor() {
    this.metrics = this.getDefaultMetrics();
    this.initialize();
  }

  static getInstance(): PerformanceOptimizationService {
    if (!this.instance) {
      this.instance = new PerformanceOptimizationService();
    }
    return this.instance;
  }

  private getDefaultMetrics(): PerformanceMetrics {
    return {
      fps: 60,
      memory: { used: 0, limit: 0, percentage: 0 },
      cpu: { cores: navigator.hardwareConcurrency || 4, usage: 0 },
      network: {
        effectiveType: '4g',
        downlink: 10,
        rtt: 50,
        saveData: false,
      },
      battery: null,
    };
  }

  /**
   * Initialize monitoring systems
   */
  private async initialize() {
    await this.initBatteryMonitoring();
    this.initNetworkMonitoring();
    this.initPerformanceMonitoring();
    this.startFrameRateMonitoring();
    this.applyOptimizations();
  }

  /**
   * Battery monitoring
   */
  private async initBatteryMonitoring() {
    if ('getBattery' in navigator) {
      try {
        this.batteryManager = await (navigator as any).getBattery();
        this.updateBatteryInfo();

        // Listen for battery changes
        this.batteryManager.addEventListener('levelchange', () =>
          this.updateBatteryInfo()
        );
        this.batteryManager.addEventListener('chargingchange', () =>
          this.updateBatteryInfo()
        );
      } catch (error) {
        console.debug('Battery API not available');
      }
    }
  }

  private updateBatteryInfo() {
    if (!this.batteryManager) return;

    this.metrics.battery = {
      level: this.batteryManager.level * 100,
      charging: this.batteryManager.charging,
      chargingTime: this.batteryManager.chargingTime,
      dischargingTime: this.batteryManager.dischargingTime,
    };

    this.adjustOptimizationLevel();
  }

  /**
   * Network monitoring
   */
  private initNetworkMonitoring() {
    this.networkConnection =
      (navigator as any).connection ||
      (navigator as any).mozConnection ||
      (navigator as any).webkitConnection;

    if (this.networkConnection) {
      this.updateNetworkInfo();
      this.networkConnection.addEventListener('change', () => {
        this.updateNetworkInfo();
        this.adjustOptimizationLevel();
      });
    }
  }

  private updateNetworkInfo() {
    if (!this.networkConnection) return;

    this.metrics.network = {
      effectiveType: this.networkConnection.effectiveType || '4g',
      downlink: this.networkConnection.downlink || 10,
      rtt: this.networkConnection.rtt || 50,
      saveData: this.networkConnection.saveData || false,
    };
  }

  /**
   * Performance monitoring
   */
  private initPerformanceMonitoring() {
    if ('PerformanceObserver' in window) {
      this.performanceObserver = new PerformanceObserver((list) => {
        for (const entry of list.getEntries()) {
          this.processPerformanceEntry(entry);
        }
      });

      // Observe different performance metrics
      try {
        this.performanceObserver.observe({
          entryTypes: ['measure', 'navigation', 'resource'],
        });
      } catch (e) {
        // Some entry types might not be supported
      }
    }

    // Memory monitoring
    this.monitorMemory();
  }

  private processPerformanceEntry(entry: PerformanceEntry) {
    // Process different types of performance entries
    if (entry.entryType === 'navigation') {
      // Navigation timing
      const navEntry = entry as PerformanceNavigationTiming;
      this.trackNavigationMetrics(navEntry);
    } else if (entry.entryType === 'resource') {
      // Resource timing
      this.trackResourceMetrics(entry as PerformanceResourceTiming);
    }
  }

  private trackNavigationMetrics(entry: PerformanceNavigationTiming) {
    const metrics = {
      dns: entry.domainLookupEnd - entry.domainLookupStart,
      tcp: entry.connectEnd - entry.connectStart,
      ttfb: entry.responseStart - entry.requestStart,
      download: entry.responseEnd - entry.responseStart,
      domComplete: entry.domComplete - entry.domInteractive,
      loadComplete: entry.loadEventEnd - entry.loadEventStart,
    };

    // Adjust optimization based on metrics
    if (metrics.ttfb > 1000 || metrics.loadComplete > 3000) {
      this.optimizationLevel = 'aggressive';
      this.applyOptimizations();
    }
  }

  private trackResourceMetrics(entry: PerformanceResourceTiming) {
    // Track slow resources
    const duration = entry.responseEnd - entry.startTime;
    if (duration > 1000) {
      console.debug(`Slow resource: ${entry.name} took ${duration}ms`);
    }
  }

  /**
   * Memory monitoring
   */
  private monitorMemory() {
    if ('memory' in performance) {
      setInterval(() => {
        const memory = (performance as any).memory;
        this.metrics.memory = {
          used: memory.usedJSHeapSize,
          limit: memory.jsHeapSizeLimit,
          percentage: (memory.usedJSHeapSize / memory.jsHeapSizeLimit) * 100,
        };

        // Trigger GC hint if memory usage is high
        if (this.metrics.memory.percentage > 80) {
          this.requestGarbageCollection();
        }
      }, 5000);
    }
  }

  /**
   * Frame rate monitoring
   */
  private startFrameRateMonitoring() {
    const measureFPS = (timestamp: number) => {
      if (this.lastFrameTime) {
        const delta = timestamp - this.lastFrameTime;
        this.frameCount++;

        // Calculate FPS every second
        if (this.frameCount >= 60) {
          this.metrics.fps = Math.round(1000 / (delta / this.frameCount));
          this.frameCount = 0;

          // Adjust quality if FPS is low
          if (this.metrics.fps < 30) {
            this.reduceQuality();
          }
        }
      }

      this.lastFrameTime = timestamp;
      this.frameRateMonitor = requestAnimationFrame(measureFPS);
    };

    this.frameRateMonitor = requestAnimationFrame(measureFPS);
  }

  /**
   * Optimization level adjustment
   */
  private adjustOptimizationLevel() {
    const battery = this.metrics.battery;
    const network = this.metrics.network;

    // Aggressive optimization conditions
    if (
      (battery && battery.level < 20 && !battery.charging) ||
      network.saveData ||
      network.effectiveType === '2g' ||
      network.effectiveType === 'slow-2g' ||
      this.metrics.memory.percentage > 80
    ) {
      this.optimizationLevel = 'aggressive';
    }
    // Balanced optimization
    else if (
      (battery && battery.level < 50 && !battery.charging) ||
      network.effectiveType === '3g' ||
      this.metrics.memory.percentage > 60
    ) {
      this.optimizationLevel = 'balanced';
    }
    // Performance mode
    else {
      this.optimizationLevel = 'performance';
    }

    this.applyOptimizations();
  }

  /**
   * Apply optimizations based on current level
   */
  private applyOptimizations() {
    switch (this.optimizationLevel) {
      case 'aggressive':
        this.applyAggressiveOptimizations();
        break;
      case 'balanced':
        this.applyBalancedOptimizations();
        break;
      case 'performance':
        this.applyPerformanceOptimizations();
        break;
    }
  }

  private applyAggressiveOptimizations() {
    // Reduce animation frame rate
    this.setAnimationFrameRate(30);

    // Disable non-essential animations
    document.documentElement.style.setProperty('--animation-duration', '0s');

    // Reduce image quality
    this.setImageQuality(40);

    // Disable auto-play videos
    this.disableAutoplay();

    // Reduce network requests
    this.enableRequestBatching();

    // Disable background sync
    this.disableBackgroundSync();
  }

  private applyBalancedOptimizations() {
    // Normal animation frame rate
    this.setAnimationFrameRate(60);

    // Reduce some animations
    document.documentElement.style.setProperty('--animation-duration', '0.2s');

    // Moderate image quality
    this.setImageQuality(70);

    // Selective autoplay
    this.enableSelectiveAutoplay();

    // Normal network requests
    this.enableNormalRequests();

    // Limited background sync
    this.enableLimitedBackgroundSync();
  }

  private applyPerformanceOptimizations() {
    // Full animations
    this.setAnimationFrameRate(60);
    document.documentElement.style.setProperty('--animation-duration', '0.3s');

    // High image quality
    this.setImageQuality(90);

    // Enable all features
    this.enableAutoplay();
    this.enableNormalRequests();
    this.enableFullBackgroundSync();
  }

  /**
   * Specific optimization methods
   */
  private setAnimationFrameRate(fps: number) {
    // Implement frame rate limiting
    const frameTime = 1000 / fps;
    let lastFrame = 0;

    window.requestAnimationFrame = ((callback: FrameRequestCallback) => {
      const now = performance.now();
      const nextFrame = lastFrame + frameTime;

      if (now >= nextFrame) {
        lastFrame = now;
        return originalRAF(callback);
      } else {
        return setTimeout(() => {
          lastFrame = performance.now();
          callback(lastFrame);
        }, nextFrame - now);
      }
    }) as any;
  }

  private setImageQuality(quality: number) {
    // Set global image quality
    document.documentElement.style.setProperty(
      '--image-quality',
      quality.toString()
    );

    // Update existing images
    document.querySelectorAll('img').forEach((img) => {
      const src = img.src;
      if (src.includes('?')) {
        img.src = src.replace(/q=\d+/, `q=${quality}`);
      } else {
        img.src = `${src}?q=${quality}`;
      }
    });
  }

  private disableAutoplay() {
    document.querySelectorAll('video').forEach((video) => {
      video.autoplay = false;
      video.pause();
    });
  }

  private enableSelectiveAutoplay() {
    document.querySelectorAll('video').forEach((video) => {
      // Only autoplay if visible
      const rect = video.getBoundingClientRect();
      if (rect.top >= 0 && rect.bottom <= window.innerHeight) {
        video.autoplay = true;
      }
    });
  }

  private enableAutoplay() {
    document.querySelectorAll('video').forEach((video) => {
      video.autoplay = true;
    });
  }

  private enableRequestBatching() {
    // Implement request batching logic
    // This would batch multiple API requests together
  }

  private enableNormalRequests() {
    // Normal request behavior
  }

  private disableBackgroundSync() {
    if (
      'serviceWorker' in navigator &&
      'sync' in ServiceWorkerRegistration.prototype
    ) {
      // Disable background sync
    }
  }

  private enableLimitedBackgroundSync() {
    // Limited background sync (e.g., every 30 minutes)
  }

  private enableFullBackgroundSync() {
    // Full background sync
  }

  /**
   * Quality reduction methods
   */
  private reduceQuality() {
    // Reduce render quality for better performance
    const canvas = document.querySelector('canvas');
    if (canvas) {
      const ctx = canvas.getContext('2d');
      if (ctx) {
        ctx.imageSmoothingEnabled = false;
      }
    }

    // Reduce CSS effects
    document.body.classList.add('reduced-motion');
  }

  /**
   * Garbage collection hint
   */
  private requestGarbageCollection() {
    // Request idle callback for cleanup
    if ('requestIdleCallback' in window) {
      requestIdleCallback(() => {
        // Clear caches
        this.clearUnusedCaches();

        // Remove detached DOM nodes
        this.cleanupDetachedNodes();

        // Clear unused timers
        this.clearUnusedTimers();
      });
    }
  }

  private clearUnusedCaches() {
    // Clear application caches
    if ('caches' in window) {
      caches.keys().then((names) => {
        names.forEach((name) => {
          // Keep only essential caches
          if (!name.includes('essential')) {
            caches.delete(name);
          }
        });
      });
    }
  }

  private cleanupDetachedNodes() {
    // Find and remove detached DOM nodes
    const allElements = document.querySelectorAll('*');
    allElements.forEach((element) => {
      if (!document.body.contains(element)) {
        element.remove();
      }
    });
  }

  private clearUnusedTimers() {
    // Clear any leaked timers
    // This is a simplified version
  }

  /**
   * Public API
   */
  getMetrics(): PerformanceMetrics {
    return { ...this.metrics };
  }

  getOptimizationLevel(): string {
    return this.optimizationLevel;
  }

  forceOptimizationLevel(level: 'aggressive' | 'balanced' | 'performance') {
    this.optimizationLevel = level;
    this.applyOptimizations();
  }

  /**
   * Request optimization hints
   */
  requestIdleTask<T>(task: () => T, options?: IdleRequestOptions): Promise<T> {
    return new Promise((resolve) => {
      if ('requestIdleCallback' in window) {
        requestIdleCallback((deadline) => {
          if (deadline.timeRemaining() > 0) {
            resolve(task());
          } else {
            // Retry in next idle period
            requestIdleCallback(() => resolve(task()), options);
          }
        }, options);
      } else {
        // Fallback to setTimeout
        setTimeout(() => resolve(task()), 0);
      }
    });
  }

  /**
   * Network optimization
   */
  async prefetchResources(urls: string[]) {
    if (this.optimizationLevel === 'aggressive') return;

    const connection = this.metrics.network;
    if (connection.saveData || connection.effectiveType === '2g') return;

    urls.forEach((url) => {
      const link = document.createElement('link');
      link.rel = 'prefetch';
      link.href = url;
      document.head.appendChild(link);
    });
  }

  /**
   * Cleanup
   */
  destroy() {
    if (this.frameRateMonitor) {
      cancelAnimationFrame(this.frameRateMonitor);
    }

    if (this.performanceObserver) {
      this.performanceObserver.disconnect();
    }

    // Remove event listeners
    if (this.batteryManager) {
      this.batteryManager.removeEventListener(
        'levelchange',
        this.updateBatteryInfo
      );
      this.batteryManager.removeEventListener(
        'chargingchange',
        this.updateBatteryInfo
      );
    }

    if (this.networkConnection) {
      this.networkConnection.removeEventListener(
        'change',
        this.updateNetworkInfo
      );
    }
  }
}

// Export singleton instance
const originalRAF = window.requestAnimationFrame;
export const performanceOptimizer =
  PerformanceOptimizationService.getInstance();

// Export types
export type { PerformanceMetrics, NetworkInfo, BatteryInfo };
