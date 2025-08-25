'use client';

/**
 * Utility функции для определения типа устройства и браузера
 */

import { DeviceType } from '../types/behavior';

/**
 * Информация о браузере
 */
export interface BrowserInfo {
  name: string;
  version: string;
  full: string;
}

/**
 * Информация об устройстве
 */
export interface DeviceInfo {
  type: DeviceType;
  os: string;
  osVersion: string;
  isMobile: boolean;
  isTablet: boolean;
  isDesktop: boolean;
}

/**
 * Полная информация о пользовательском агенте
 */
export interface UserAgentInfo {
  browser: BrowserInfo;
  device: DeviceInfo;
  userAgent: string;
  platform: string;
}

/**
 * Определяет тип устройства на основе user agent
 */
export function getDeviceType(
  userAgent: string = navigator.userAgent
): DeviceType {
  const ua = userAgent.toLowerCase();

  // Проверка на мобильные устройства
  const mobileRegex =
    /android|webos|iphone|ipad|ipod|blackberry|iemobile|opera mini/i;
  const isMobile = mobileRegex.test(ua);

  // Проверка на планшеты
  const tabletRegex =
    /ipad|android(?!.*mobile)|tablet|kindle|silk|playbook|bb10/i;
  const isTablet = tabletRegex.test(ua);

  // Специфические проверки
  const isIOS = /iphone|ipad|ipod/i.test(ua);
  const isAndroid = /android/i.test(ua);
  const isAndroidTablet = isAndroid && !/mobile/i.test(ua);
  const isIPad = /ipad/i.test(ua);

  // Определение типа устройства
  if (isIPad || isAndroidTablet || isTablet) {
    return 'tablet';
  }

  if (isIOS || (isAndroid && /mobile/i.test(ua)) || isMobile) {
    return 'mobile';
  }

  return 'desktop';
}

/**
 * Определяет информацию о браузере
 */
export function getBrowserInfo(
  userAgent: string = navigator.userAgent
): BrowserInfo {
  const ua = userAgent.toLowerCase();

  // Список браузеров для определения (порядок важен!)
  const browsers = [
    { name: 'Edge', regex: /edg\/([\d.]+)/ },
    { name: 'Chrome', regex: /chrome\/([\d.]+)/ },
    { name: 'Firefox', regex: /firefox\/([\d.]+)/ },
    { name: 'Safari', regex: /version\/([\d.]+).*safari/ },
    { name: 'Opera', regex: /opera\/([\d.]+)/ },
    { name: 'Opera', regex: /opr\/([\d.]+)/ },
    { name: 'Internet Explorer', regex: /msie\s([\d.]+)/ },
    { name: 'Internet Explorer', regex: /trident.*rv:([\d.]+)/ },
  ];

  for (const browser of browsers) {
    const match = ua.match(browser.regex);
    if (match) {
      const version = (match[1] || 'Unknown').replace(
        /[\x00-\x1F\x7F-\x9F]/g,
        ''
      );
      return {
        name: browser.name,
        version: version,
        full: `${browser.name} ${version}`,
      };
    }
  }

  return {
    name: 'Unknown',
    version: '0.0.0',
    full: 'Unknown Browser',
  };
}

/**
 * Определяет информацию об операционной системе
 */
export function getOSInfo(userAgent: string = navigator.userAgent): {
  os: string;
  osVersion: string;
} {
  const ua = userAgent.toLowerCase();

  // Список операционных систем
  const systems = [
    {
      name: 'Windows 11',
      regex: /windows nt 10\.0.*wow64|windows nt 10\.0.*win64/,
    },
    { name: 'Windows 10', regex: /windows nt 10\.0/ },
    { name: 'Windows 8.1', regex: /windows nt 6\.3/ },
    { name: 'Windows 8', regex: /windows nt 6\.2/ },
    { name: 'Windows 7', regex: /windows nt 6\.1/ },
    { name: 'Windows Vista', regex: /windows nt 6\.0/ },
    { name: 'Windows XP', regex: /windows nt 5\.1/ },
    { name: 'macOS', regex: /mac os x ([\d_.]+)/ },
    { name: 'iOS', regex: /iphone os ([\d_.]+)/ },
    { name: 'iOS', regex: /ipad.*os ([\d_.]+)/ },
    { name: 'Android', regex: /android ([\d.]+)/ },
    { name: 'Linux', regex: /linux/ },
    { name: 'Chrome OS', regex: /cros/ },
  ];

  for (const system of systems) {
    const match = ua.match(system.regex);
    if (match) {
      const version = match[1] ? match[1].replace(/_/g, '.') : 'Unknown';
      return {
        os: system.name,
        osVersion: version,
      };
    }
  }

  return {
    os: 'Unknown',
    osVersion: 'Unknown',
  };
}

/**
 * Получает полную информацию о пользовательском агенте
 */
export function getUserAgentInfo(
  userAgent: string = typeof navigator !== 'undefined'
    ? navigator.userAgent
    : ''
): UserAgentInfo {
  const browser = getBrowserInfo(userAgent);
  const { os, osVersion } = getOSInfo(userAgent);
  const deviceType = getDeviceType(userAgent);

  // Санитизация строк для предотвращения проблем с кодировкой
  const sanitizeString = (str: string): string => {
    // Удаляем непечатаемые символы и заменяем проблемные символы
    return str.replace(/[\x00-\x1F\x7F-\x9F]/g, '').trim();
  };

  return {
    browser,
    device: {
      type: deviceType,
      os: sanitizeString(os),
      osVersion: sanitizeString(osVersion),
      isMobile: deviceType === 'mobile',
      isTablet: deviceType === 'tablet',
      isDesktop: deviceType === 'desktop',
    },
    userAgent: sanitizeString(userAgent),
    platform:
      typeof navigator !== 'undefined'
        ? sanitizeString(navigator.platform)
        : 'Unknown',
  };
}

/**
 * Проверяет, является ли устройство мобильным
 */
export function isMobileDevice(
  userAgent: string = typeof navigator !== 'undefined'
    ? navigator.userAgent
    : ''
): boolean {
  return getDeviceType(userAgent) === 'mobile';
}

/**
 * Проверяет, является ли устройство планшетом
 */
export function isTabletDevice(
  userAgent: string = navigator.userAgent
): boolean {
  return getDeviceType(userAgent) === 'tablet';
}

/**
 * Проверяет, является ли устройство настольным компьютером
 */
export function isDesktopDevice(
  userAgent: string = navigator.userAgent
): boolean {
  return getDeviceType(userAgent) === 'desktop';
}

/**
 * Получает информацию о размере экрана
 */
export function getScreenInfo() {
  if (typeof window === 'undefined') {
    return {
      width: 0,
      height: 0,
      pixelRatio: 1,
    };
  }

  return {
    width: window.screen.width,
    height: window.screen.height,
    pixelRatio: window.devicePixelRatio || 1,
  };
}

/**
 * Получает информацию о размере viewport
 */
export function getViewportInfo() {
  if (typeof window === 'undefined') {
    return {
      width: 0,
      height: 0,
    };
  }

  return {
    width: window.innerWidth,
    height: window.innerHeight,
  };
}

/**
 * Определяет, поддерживает ли браузер touch events
 */
export function isTouchDevice(): boolean {
  if (typeof window === 'undefined') {
    return false;
  }

  return (
    'ontouchstart' in window ||
    navigator.maxTouchPoints > 0 ||
    // @ts-expect-error Legacy property for old IE
    navigator.msMaxTouchPoints > 0
  );
}

/**
 * Получает информацию о поддерживаемых технологиях
 */
export function getSupportedFeatures() {
  if (typeof window === 'undefined') {
    return {
      localStorage: false,
      sessionStorage: false,
      cookies: false,
      webWorkers: false,
      serviceWorkers: false,
      webGL: false,
      canvas: false,
      geolocation: false,
      notifications: false,
    };
  }

  return {
    localStorage: !!window.localStorage,
    sessionStorage: !!window.sessionStorage,
    cookies: navigator.cookieEnabled,
    webWorkers: !!window.Worker,
    serviceWorkers: 'serviceWorker' in navigator,
    webGL: !!window.WebGLRenderingContext,
    canvas: !!window.HTMLCanvasElement,
    geolocation: 'geolocation' in navigator,
    notifications: 'Notification' in window,
  };
}

/**
 * Получает краткую строку для идентификации браузера/устройства
 */
export function getBrowserFingerprint(): string {
  const info = getUserAgentInfo();
  const screen = getScreenInfo();
  const viewport = getViewportInfo();

  return `${info.browser.name}_${info.browser.version}_${info.device.type}_${info.device.os}_${screen.width}x${screen.height}_${viewport.width}x${viewport.height}`;
}

/**
 * Экспорт для тестирования
 */
export const __testing__ = {
  browsers: [
    { name: 'Edge', regex: /edg\/([\d.]+)/ },
    { name: 'Chrome', regex: /chrome\/([\d.]+)/ },
    { name: 'Firefox', regex: /firefox\/([\d.]+)/ },
    { name: 'Safari', regex: /version\/([\d.]+).*safari/ },
    { name: 'Opera', regex: /opera\/([\d.]+)/ },
    { name: 'Opera', regex: /opr\/([\d.]+)/ },
    { name: 'Internet Explorer', regex: /msie\s([\d.]+)/ },
    { name: 'Internet Explorer', regex: /trident.*rv:([\d.]+)/ },
  ],
};
