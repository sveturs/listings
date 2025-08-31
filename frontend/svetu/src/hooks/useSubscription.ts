'use client';

import { useState, useEffect, useCallback } from 'react';
import { useRouter } from 'next/navigation';
import { useLocale } from 'next-intl';
import { useAuth } from '@/contexts/AuthContext';
import { subscriptionService } from '@/services/subscriptions';
import type { components } from '@/types/generated/api';

type SubscriptionPlan =
  components['schemas']['backend_internal_domain_models.SubscriptionPlanDetails'];
type UserSubscriptionInfo =
  components['schemas']['backend_internal_domain_models.UserSubscriptionInfo'];

export interface UseSubscriptionResult {
  // Data
  plans: SubscriptionPlan[];
  currentSubscription: UserSubscriptionInfo | null;
  currentPlan: SubscriptionPlan | null;

  // Loading states
  isLoading: boolean;
  isLoadingPlans: boolean;
  isLoadingSubscription: boolean;
  isProcessing: boolean;

  // Error state
  error: string | null;

  // Actions
  refreshSubscription: () => Promise<void>;
  selectPlan: (
    planCode: string,
    billingCycle: 'monthly' | 'yearly'
  ) => Promise<void>;
  upgradePlan: (
    planCode: string,
    billingCycle?: 'monthly' | 'yearly'
  ) => Promise<void>;
  cancelSubscription: (reason?: string) => Promise<void>;
  canUseFeature: (
    feature: 'ai' | 'live' | 'export' | 'customDomain'
  ) => boolean;
  getRemainingLimits: () => {
    storefronts: number;
    products: number;
    staff: number;
    images: number;
  };
}

export function useSubscription(): UseSubscriptionResult {
  const router = useRouter();
  const locale = useLocale();
  const auth = useAuth();
  const [plans, setPlans] = useState<SubscriptionPlan[]>([]);
  const [currentSubscription, setCurrentSubscription] =
    useState<UserSubscriptionInfo | null>(null);
  const [isLoadingPlans, setIsLoadingPlans] = useState(false);
  const [isLoadingSubscription, setIsLoadingSubscription] = useState(false);
  const [isProcessing, setIsProcessing] = useState(false);
  const [error, setError] = useState<string | null>(null);

  // Load plans on mount
  useEffect(() => {
    loadPlans();
    loadCurrentSubscription();
  }, []);

  const loadPlans = async () => {
    setIsLoadingPlans(true);
    setError(null);
    try {
      const data = await subscriptionService.getPlans();
      setPlans(data);
    } catch (err: any) {
      setError(err.response?.data?.message || 'Failed to load plans');
      console.error('Failed to load plans:', err);
    } finally {
      setIsLoadingPlans(false);
    }
  };

  const loadCurrentSubscription = async () => {
    setIsLoadingSubscription(true);
    try {
      const data = await subscriptionService.getCurrentSubscription();
      setCurrentSubscription(data);
    } catch (err: any) {
      // Ignore auth errors - user may not be logged in
      if (err.response?.status !== 401) {
        console.error('Failed to load subscription:', err);
      }
    } finally {
      setIsLoadingSubscription(false);
    }
  };

  const refreshSubscription = useCallback(async () => {
    await loadCurrentSubscription();
  }, []);

  const selectPlan = useCallback(
    async (planCode: string, billingCycle: 'monthly' | 'yearly') => {
      setIsProcessing(true);
      setError(null);

      try {
        // Check if user is authenticated
        if (!auth?.isAuthenticated) {
          // Redirect to login with return URL
          const returnUrl = encodeURIComponent(
            `/${locale}/pricing?plan=${planCode}&cycle=${billingCycle}`
          );
          router.push(`/${locale}/auth/login?returnUrl=${returnUrl}`);
          return;
        }

        // Check if this is free plan
        const plan = plans.find((p) => p.code === planCode);
        if (plan && plan.price_monthly === 0 && plan.price_yearly === 0) {
          // Create free subscription directly
          await subscriptionService.createSubscription({
            plan_code: planCode,
            billing_cycle: billingCycle,
            payment_method: 'free',
          });

          // Refresh subscription data
          await loadCurrentSubscription();

          // Redirect to success page
          router.push(`/${locale}/subscription/success`);
        } else {
          // Initiate payment for paid plan
          const response = await subscriptionService.initiatePayment({
            plan_code: planCode,
            billing_cycle: billingCycle,
            return_url: `${window.location.origin}/subscription/payment-complete`,
          });

          if (response.payment_required && response.redirect_url) {
            // Redirect to payment gateway
            window.location.href = response.redirect_url;
          } else {
            // No payment required (shouldn't happen for paid plans)
            await loadCurrentSubscription();
            router.push(`/${locale}/subscription/success`);
          }
        }
      } catch (err: any) {
        setError(err.response?.data?.message || 'Failed to select plan');
        console.error('Failed to select plan:', err);
      } finally {
        setIsProcessing(false);
      }
    },
    [auth, locale, plans, router]
  );

  const upgradePlan = useCallback(
    async (planCode: string, billingCycle?: 'monthly' | 'yearly') => {
      setIsProcessing(true);
      setError(null);

      try {
        // Check if user has subscription
        if (!currentSubscription?.subscription_id) {
          // No subscription - create new one
          await selectPlan(planCode, billingCycle || 'monthly');
          return;
        }

        // Upgrade existing subscription
        await subscriptionService.upgradeSubscription({
          plan_code: planCode,
          billing_cycle: billingCycle,
        });

        // Check if payment is needed
        const plan = plans.find((p) => p.code === planCode);
        if (
          plan &&
          ((plan.price_monthly ?? 0) > 0 || (plan.price_yearly ?? 0) > 0)
        ) {
          // Initiate payment for upgrade
          const response = await subscriptionService.initiatePayment({
            plan_code: planCode,
            billing_cycle: (billingCycle ||
              currentSubscription.billing_cycle ||
              'monthly') as 'monthly' | 'yearly',
            return_url: `${window.location.origin}/subscription/upgrade-complete`,
          });

          if (response.payment_required && response.redirect_url) {
            window.location.href = response.redirect_url;
          }
        } else {
          // Refresh subscription data
          await loadCurrentSubscription();
          router.push('/subscription/success');
        }
      } catch (err: any) {
        setError(err.response?.data?.message || 'Failed to upgrade plan');
        console.error('Failed to upgrade plan:', err);
      } finally {
        setIsProcessing(false);
      }
    },
    [currentSubscription, plans, router, selectPlan]
  );

  const cancelSubscription = useCallback(
    async (reason?: string) => {
      setIsProcessing(true);
      setError(null);

      try {
        await subscriptionService.cancelSubscription(reason);
        await loadCurrentSubscription();
        router.push(`/${locale}/subscription/cancelled`);
      } catch (err: any) {
        setError(
          err.response?.data?.message || 'Failed to cancel subscription'
        );
        console.error('Failed to cancel subscription:', err);
      } finally {
        setIsProcessing(false);
      }
    },
    [locale, router]
  );

  const canUseFeature = useCallback(
    (feature: 'ai' | 'live' | 'export' | 'customDomain'): boolean => {
      if (!currentSubscription) return false;

      switch (feature) {
        case 'ai':
          return currentSubscription.has_ai || false;
        case 'live':
          return currentSubscription.has_live || false;
        case 'export':
          return currentSubscription.has_export || false;
        case 'customDomain':
          return currentSubscription.has_custom_domain || false;
        default:
          return false;
      }
    },
    [currentSubscription]
  );

  const getRemainingLimits = useCallback(() => {
    const defaults = {
      storefronts: 0,
      products: 0,
      staff: 0,
      images: 0,
    };

    if (!currentSubscription) return defaults;

    const maxStorefronts = currentSubscription.max_storefronts || 0;
    const usedStorefronts = currentSubscription.used_storefronts || 0;

    return {
      storefronts:
        maxStorefronts === -1
          ? Infinity
          : Math.max(0, maxStorefronts - usedStorefronts),
      products:
        currentSubscription.max_products === -1
          ? Infinity
          : currentSubscription.max_products || 0,
      staff:
        currentSubscription.max_staff === -1
          ? Infinity
          : currentSubscription.max_staff || 0,
      images:
        currentSubscription.max_images === -1
          ? Infinity
          : currentSubscription.max_images || 0,
    };
  }, [currentSubscription]);

  // Get current plan object
  const currentPlan =
    plans.find((p) => p.code === currentSubscription?.plan_code) || null;

  return {
    plans,
    currentSubscription,
    currentPlan,
    isLoading: isLoadingPlans || isLoadingSubscription,
    isLoadingPlans,
    isLoadingSubscription,
    isProcessing,
    error,
    refreshSubscription,
    selectPlan,
    upgradePlan,
    cancelSubscription,
    canUseFeature,
    getRemainingLimits,
  };
}
