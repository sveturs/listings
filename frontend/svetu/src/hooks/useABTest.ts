import { useEffect, useState, useCallback, useMemo } from 'react';
import {
  abTesting,
  type Variant,
  type UserContext,
} from '../services/abTestingService';

interface UseABTestOptions {
  experimentId: string;
  userId?: string;
  trackImpression?: boolean;
  customProperties?: Record<string, any>;
}

interface UseABTestReturn {
  variant: Variant | null;
  isLoading: boolean;
  isControl: boolean;
  variantConfig: Record<string, any>;
  trackConversion: (eventName: string, value?: number) => void;
  trackEvent: (eventName: string, metadata?: Record<string, any>) => void;
}

/**
 * Hook for A/B testing
 */
export function useABTest({
  experimentId,
  userId,
  trackImpression = true,
  customProperties,
}: UseABTestOptions): UseABTestReturn {
  const [variant, setVariant] = useState<Variant | null>(null);
  const [isLoading, setIsLoading] = useState(true);

  // Get user context
  const userContext = useMemo((): UserContext => {
    const sessionId =
      sessionStorage.getItem('ab_session_id') ||
      Math.random().toString(36).substr(2, 9);

    if (!sessionStorage.getItem('ab_session_id')) {
      sessionStorage.setItem('ab_session_id', sessionId);
    }

    return {
      userId,
      sessionId,
      device: detectDevice(),
      browser: detectBrowser(),
      country: detectCountry(),
      language: navigator.language.split('-')[0],
      isNewUser: !localStorage.getItem('returning_user'),
      customProperties,
    };
  }, [userId, customProperties]);

  // Load variant assignment
  useEffect(() => {
    const loadVariant = async () => {
      setIsLoading(true);

      try {
        const assignedVariant = abTesting.getVariant(experimentId, userContext);
        setVariant(assignedVariant);

        // Track impression
        if (assignedVariant && trackImpression) {
          abTesting.trackConversion(experimentId, 'impression');
        }

        // Mark as returning user
        localStorage.setItem('returning_user', 'true');
      } catch (error) {
        console.error('Failed to load A/B test variant:', error);
      } finally {
        setIsLoading(false);
      }
    };

    loadVariant();
  }, [experimentId, userContext, trackImpression]);

  // Track conversion
  const trackConversion = useCallback(
    (eventName: string, value?: number) => {
      if (!variant) return;
      abTesting.trackConversion(experimentId, eventName, value);
    },
    [experimentId, variant]
  );

  // Track custom event
  const trackEvent = useCallback(
    (eventName: string, metadata?: Record<string, any>) => {
      if (!variant) return;
      abTesting.trackConversion(experimentId, eventName, undefined, metadata);
    },
    [experimentId, variant]
  );

  // Derived values
  const isControl = variant?.id === 'control';
  const variantConfig = variant?.config || {};

  return {
    variant,
    isLoading,
    isControl,
    variantConfig,
    trackConversion,
    trackEvent,
  };
}

/**
 * Hook for feature flags
 */
export function useFeatureFlag(
  flagName: string,
  defaultValue: boolean = false
): boolean {
  const [value, setValue] = useState<boolean>(() =>
    abTesting.isFeatureEnabled(flagName, defaultValue)
  );

  useEffect(() => {
    // Listen for feature flag updates
    const checkFlag = () => {
      setValue(abTesting.isFeatureEnabled(flagName, defaultValue));
    };

    // Check periodically for remote updates
    const interval = setInterval(checkFlag, 30000);

    return () => clearInterval(interval);
  }, [flagName, defaultValue]);

  return value;
}

/**
 * Hook for feature value (string/number)
 */
export function useFeatureValue<T = any>(flagName: string, defaultValue: T): T {
  const [value, setValue] = useState<T>(() =>
    abTesting.getFeatureValue(flagName, defaultValue)
  );

  useEffect(() => {
    const checkValue = () => {
      setValue(abTesting.getFeatureValue(flagName, defaultValue));
    };

    const interval = setInterval(checkValue, 30000);

    return () => clearInterval(interval);
  }, [flagName, defaultValue]);

  return value;
}

/**
 * Hook for multivariate testing
 */
export function useMultivariant(
  experiments: string[]
): Record<string, Variant | null> {
  const [variants, setVariants] = useState<Record<string, Variant | null>>({});

  useEffect(() => {
    const userContext = getUserContext();
    const loadedVariants: Record<string, Variant | null> = {};

    experiments.forEach((experimentId) => {
      loadedVariants[experimentId] = abTesting.getVariant(
        experimentId,
        userContext
      );
    });

    setVariants(loadedVariants);
  }, [experiments]);

  return variants;
}

/**
 * Hook for conversion tracking
 */
export function useConversionTracking(experimentId?: string) {
  const trackPageView = useCallback(
    (pageName: string) => {
      if (experimentId) {
        abTesting.trackConversion(experimentId, 'page_view', undefined, {
          page: pageName,
        });
      }
    },
    [experimentId]
  );

  const trackClick = useCallback(
    (elementName: string) => {
      if (experimentId) {
        abTesting.trackConversion(experimentId, 'click', undefined, {
          element: elementName,
        });
      }
    },
    [experimentId]
  );

  const trackFormSubmit = useCallback(
    (formName: string, success: boolean) => {
      if (experimentId) {
        abTesting.trackConversion(experimentId, 'form_submit', undefined, {
          form: formName,
          success,
        });
      }
    },
    [experimentId]
  );

  const trackPurchase = useCallback(
    (value: number, items: any[]) => {
      if (experimentId) {
        abTesting.trackConversion(experimentId, 'purchase', value, { items });
      }
    },
    [experimentId]
  );

  const trackCustomGoal = useCallback(
    (goalName: string, value?: number, metadata?: Record<string, any>) => {
      if (experimentId) {
        abTesting.trackConversion(experimentId, goalName, value, metadata);
      }
    },
    [experimentId]
  );

  return {
    trackPageView,
    trackClick,
    trackFormSubmit,
    trackPurchase,
    trackCustomGoal,
  };
}

/**
 * Utility functions
 */
function detectDevice(): string {
  const width = window.innerWidth;
  if (width < 768) return 'mobile';
  if (width < 1024) return 'tablet';
  return 'desktop';
}

function detectBrowser(): string {
  const userAgent = navigator.userAgent;
  if (userAgent.includes('Chrome')) return 'chrome';
  if (userAgent.includes('Safari')) return 'safari';
  if (userAgent.includes('Firefox')) return 'firefox';
  if (userAgent.includes('Edge')) return 'edge';
  return 'other';
}

function detectCountry(): string {
  // This would use a geolocation service
  // For now, use browser language as proxy
  const lang = navigator.language;
  if (lang.includes('US')) return 'US';
  if (lang.includes('GB')) return 'GB';
  if (lang.includes('DE')) return 'DE';
  if (lang.includes('FR')) return 'FR';
  if (lang.includes('RU')) return 'RU';
  if (lang.includes('RS')) return 'RS';
  return 'OTHER';
}

function getUserContext(): UserContext {
  const sessionId =
    sessionStorage.getItem('ab_session_id') ||
    Math.random().toString(36).substr(2, 9);

  return {
    sessionId,
    device: detectDevice(),
    browser: detectBrowser(),
    country: detectCountry(),
    language: navigator.language.split('-')[0],
    isNewUser: !localStorage.getItem('returning_user'),
  };
}

/**
 * Component for A/B test variants
 */
interface ABTestProps {
  experimentId: string;
  control: React.ReactNode;
  variant: React.ReactNode;
  onVariantShown?: (variantId: string) => void;
}

export function ABTest({
  experimentId,
  control,
  variant: variantComponent,
  onVariantShown,
}: ABTestProps) {
  const { variant, isLoading } = useABTest({ experimentId });

  useEffect(() => {
    if (variant && onVariantShown) {
      onVariantShown(variant.id);
    }
  }, [variant, onVariantShown]);

  if (isLoading) {
    return control; // Show control during loading
  }

  return variant?.id === 'control' ? control : variantComponent;
}

/**
 * Component for feature flags
 */
interface FeatureFlagProps {
  flag: string;
  children: React.ReactNode;
  fallback?: React.ReactNode;
}

export function FeatureFlag({
  flag,
  children,
  fallback = null,
}: FeatureFlagProps) {
  const isEnabled = useFeatureFlag(flag);
  return isEnabled ? children : fallback;
}
