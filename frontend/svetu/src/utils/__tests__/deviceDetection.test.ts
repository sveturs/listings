/**
 * Unit тесты для utility функций определения устройства и браузера
 */

import {
  getDeviceType,
  getBrowserInfo,
  getOSInfo,
  getUserAgentInfo,
  isMobileDevice,
  isTabletDevice,
  isDesktopDevice,
} from '../deviceDetection';

describe('deviceDetection utilities', () => {
  describe('getDeviceType', () => {
    test('should detect mobile devices', () => {
      const mobileUserAgents = [
        'Mozilla/5.0 (iPhone; CPU iPhone OS 15_0 like Mac OS X) AppleWebKit/605.1.15',
        'Mozilla/5.0 (Linux; Android 11; SM-G991B) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/91.0.4472.124 Mobile Safari/537.36',
        'Mozilla/5.0 (Linux; Android 10; mobile) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/91.0.4472.124 Mobile Safari/537.36',
      ];

      mobileUserAgents.forEach((ua) => {
        expect(getDeviceType(ua)).toBe('mobile');
      });
    });

    test('should detect tablet devices', () => {
      const tabletUserAgents = [
        'Mozilla/5.0 (iPad; CPU OS 15_0 like Mac OS X) AppleWebKit/605.1.15',
        'Mozilla/5.0 (Linux; Android 11; SM-T860) AppleWebKit/537.36',
        'Mozilla/5.0 (Linux; Android 10; tablet) AppleWebKit/537.36',
      ];

      tabletUserAgents.forEach((ua) => {
        expect(getDeviceType(ua)).toBe('tablet');
      });
    });

    test('should detect desktop devices', () => {
      const desktopUserAgents = [
        'Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36',
        'Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36',
        'Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36',
      ];

      desktopUserAgents.forEach((ua) => {
        expect(getDeviceType(ua)).toBe('desktop');
      });
    });
  });

  describe('getBrowserInfo', () => {
    test('should detect Chrome browser', () => {
      const chromeUA =
        'Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/91.0.4472.124 Safari/537.36';
      const browserInfo = getBrowserInfo(chromeUA);

      expect(browserInfo.name).toBe('Chrome');
      expect(browserInfo.version).toBe('91.0.4472.124');
      expect(browserInfo.full).toBe('Chrome 91.0.4472.124');
    });

    test('should detect Firefox browser', () => {
      const firefoxUA =
        'Mozilla/5.0 (Windows NT 10.0; Win64; x64; rv:89.0) Gecko/20100101 Firefox/89.0';
      const browserInfo = getBrowserInfo(firefoxUA);

      expect(browserInfo.name).toBe('Firefox');
      expect(browserInfo.version).toBe('89.0');
      expect(browserInfo.full).toBe('Firefox 89.0');
    });

    test('should detect Safari browser', () => {
      const safariUA =
        'Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/605.1.15 (KHTML, like Gecko) Version/14.1.1 Safari/605.1.15';
      const browserInfo = getBrowserInfo(safariUA);

      expect(browserInfo.name).toBe('Safari');
      expect(browserInfo.version).toBe('14.1.1');
      expect(browserInfo.full).toBe('Safari 14.1.1');
    });

    test('should handle unknown browser', () => {
      const unknownUA = 'Some Unknown Browser 1.0';
      const browserInfo = getBrowserInfo(unknownUA);

      expect(browserInfo.name).toBe('Unknown');
      expect(browserInfo.version).toBe('0.0.0');
      expect(browserInfo.full).toBe('Unknown Browser');
    });
  });

  describe('getOSInfo', () => {
    test('should detect Windows', () => {
      const windowsUA =
        'Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36';
      const osInfo = getOSInfo(windowsUA);

      expect(osInfo.os).toBe('Windows 11'); // или 'Windows 10' в зависимости от логики определения
    });

    test('should detect macOS', () => {
      const macUA =
        'Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36';
      const osInfo = getOSInfo(macUA);

      expect(osInfo.os).toBe('macOS');
      expect(osInfo.osVersion).toBe('10.15.7');
    });

    test('should detect Android', () => {
      const androidUA =
        'Mozilla/5.0 (Linux; Android 11; SM-G991B) AppleWebKit/537.36';
      const osInfo = getOSInfo(androidUA);

      expect(osInfo.os).toBe('Android');
      expect(osInfo.osVersion).toBe('11');
    });

    test('should detect iOS', () => {
      const iOSUA =
        'Mozilla/5.0 (iPhone; CPU iPhone OS 15_0 like Mac OS X) AppleWebKit/605.1.15';
      const osInfo = getOSInfo(iOSUA);

      expect(osInfo.os).toBe('iOS');
      expect(osInfo.osVersion).toBe('15.0');
    });
  });

  describe('helper functions', () => {
    test('isMobileDevice should work correctly', () => {
      expect(
        isMobileDevice('Mozilla/5.0 (iPhone; CPU iPhone OS 15_0 like Mac OS X)')
      ).toBe(true);
      expect(isMobileDevice('Mozilla/5.0 (Windows NT 10.0; Win64; x64)')).toBe(
        false
      );
    });

    test('isTabletDevice should work correctly', () => {
      expect(
        isTabletDevice('Mozilla/5.0 (iPad; CPU OS 15_0 like Mac OS X)')
      ).toBe(true);
      expect(isTabletDevice('Mozilla/5.0 (Windows NT 10.0; Win64; x64)')).toBe(
        false
      );
    });

    test('isDesktopDevice should work correctly', () => {
      expect(isDesktopDevice('Mozilla/5.0 (Windows NT 10.0; Win64; x64)')).toBe(
        true
      );
      expect(
        isDesktopDevice(
          'Mozilla/5.0 (iPhone; CPU iPhone OS 15_0 like Mac OS X)'
        )
      ).toBe(false);
    });
  });

  describe('getUserAgentInfo', () => {
    test('should provide complete user agent information', () => {
      const chromeUA =
        'Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/91.0.4472.124 Safari/537.36';
      const info = getUserAgentInfo(chromeUA);

      expect(info.browser.name).toBe('Chrome');
      expect(info.device.type).toBe('desktop');
      expect(info.device.os).toBe('Windows 11');
      expect(info.device.isDesktop).toBe(true);
      expect(info.device.isMobile).toBe(false);
      expect(info.device.isTablet).toBe(false);
      expect(info.userAgent).toBe(chromeUA);
    });
  });
});
