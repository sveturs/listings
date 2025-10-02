import { apiClient } from '@/services/api-client';
import type { components } from '@/types/generated/api';

// Type aliases for better readability
type SubscriptionPlan =
  components['schemas']['backend_internal_domain_models.SubscriptionPlanDetails'];
type UserSubscription =
  components['schemas']['backend_internal_domain_models.UserSubscription'];
type UserSubscriptionInfo =
  components['schemas']['backend_internal_domain_models.UserSubscriptionInfo'];
type CheckLimitRequest =
  components['schemas']['backend_internal_domain_models.CheckLimitRequest'];
type CheckLimitResponse =
  components['schemas']['backend_internal_domain_models.CheckLimitResponse'];

export interface CreateSubscriptionRequest {
  plan_code: string;
  billing_cycle: 'monthly' | 'yearly';
  payment_method: string;
  start_trial?: boolean;
}

export interface UpgradeSubscriptionRequest {
  plan_code: string;
  billing_cycle?: 'monthly' | 'yearly';
}

export interface InitiatePaymentRequest {
  plan_code: string;
  billing_cycle?: 'monthly' | 'yearly';
  return_url?: string;
}

export interface PaymentInitiationResponse {
  payment_required: boolean;
  payment_intent_id?: string;
  redirect_url?: string;
  amount?: number;
  currency?: string;
  message?: string;
}

class SubscriptionService {
  /**
   * Get all available subscription plans
   */
  async getPlans(): Promise<SubscriptionPlan[]> {
    try {
      const response = await apiClient.get('/subscriptions/plans');
      return response.data.data;
    } catch (error) {
      console.error('Failed to get subscription plans:', error);
      throw error;
    }
  }

  /**
   * Get current user's subscription
   */
  async getCurrentSubscription(): Promise<UserSubscriptionInfo | null> {
    try {
      const response = await apiClient.get('/subscriptions/current');
      return response.data.data;
    } catch (error: any) {
      if (error.response?.status === 401) {
        // User not authenticated
        return null;
      }
      console.error('Failed to get current subscription:', error);
      throw error;
    }
  }

  /**
   * Create new subscription
   */
  async createSubscription(
    request: CreateSubscriptionRequest
  ): Promise<UserSubscription> {
    try {
      const response = await apiClient.post('/subscriptions', request);
      return response.data.data;
    } catch (error) {
      console.error('Failed to create subscription:', error);
      throw error;
    }
  }

  /**
   * Upgrade existing subscription
   */
  async upgradeSubscription(
    request: UpgradeSubscriptionRequest
  ): Promise<UserSubscription> {
    try {
      const response = await apiClient.post('/subscriptions/upgrade', request);
      return response.data.data;
    } catch (error) {
      console.error('Failed to upgrade subscription:', error);
      throw error;
    }
  }

  /**
   * Cancel subscription
   */
  async cancelSubscription(reason?: string): Promise<void> {
    try {
      await apiClient.post('/subscriptions/cancel', { reason });
    } catch (error) {
      console.error('Failed to cancel subscription:', error);
      throw error;
    }
  }

  /**
   * Check subscription limits for a resource
   */
  async checkLimits(request: CheckLimitRequest): Promise<CheckLimitResponse> {
    try {
      const response = await apiClient.post(
        '/subscriptions/check-limits',
        request
      );
      return response.data.data;
    } catch (error) {
      console.error('Failed to check limits:', error);
      throw error;
    }
  }

  /**
   * Initiate payment for subscription
   */
  async initiatePayment(
    request: InitiatePaymentRequest
  ): Promise<PaymentInitiationResponse> {
    try {
      const response = await apiClient.post(
        '/subscriptions/initiate-payment',
        request
      );
      return response.data.data;
    } catch (error) {
      console.error('Failed to initiate payment:', error);
      throw error;
    }
  }

  /**
   * Complete payment after returning from payment gateway
   */
  async completePayment(paymentIntentId: string): Promise<UserSubscription> {
    try {
      const response = await apiClient.post(
        `/subscriptions/complete-payment?payment_intent=${paymentIntentId}`,
        {}
      );
      return response.data.data;
    } catch (error) {
      console.error('Failed to complete payment:', error);
      throw error;
    }
  }

  /**
   * Format price for display
   */
  formatPrice(amount: number, currency: string = 'EUR'): string {
    return new Intl.NumberFormat('en-US', {
      style: 'currency',
      currency: currency,
      minimumFractionDigits: 0,
      maximumFractionDigits: 2,
    }).format(amount);
  }

  /**
   * Get plan by code from cached plans
   */
  getPlanByCode(
    plans: SubscriptionPlan[],
    code: string
  ): SubscriptionPlan | undefined {
    return plans.find((plan) => plan.code === code);
  }

  /**
   * Check if user can create storefront
   */
  async canCreateStorefront(): Promise<CheckLimitResponse> {
    return this.checkLimits({
      resource_type: 'storefront',
      count: 1,
    });
  }

  /**
   * Check if user can add product
   */
  async canAddProduct(): Promise<CheckLimitResponse> {
    return this.checkLimits({
      resource_type: 'product',
      count: 1,
    });
  }

  /**
   * Check if user can add staff member
   */
  async canAddStaff(): Promise<CheckLimitResponse> {
    return this.checkLimits({
      resource_type: 'staff',
      count: 1,
    });
  }

  /**
   * Check if user can upload images
   */
  async canUploadImages(count: number = 1): Promise<CheckLimitResponse> {
    return this.checkLimits({
      resource_type: 'image',
      count,
    });
  }

  /**
   * Get recommended plan for user based on needs
   */
  getRecommendedPlan(
    plans: SubscriptionPlan[],
    needs: {
      storefronts?: number;
      products?: number;
      staff?: number;
      ai?: boolean;
      customDomain?: boolean;
    }
  ): SubscriptionPlan | undefined {
    // Filter plans that meet requirements
    const suitablePlans = plans.filter((plan) => {
      if (
        needs.storefronts &&
        (plan.max_storefronts ?? 0) !== -1 &&
        (plan.max_storefronts ?? 0) < needs.storefronts
      ) {
        return false;
      }
      if (
        needs.products &&
        (plan.max_products_per_storefront ?? 0) !== -1 &&
        (plan.max_products_per_storefront ?? 0) < needs.products
      ) {
        return false;
      }
      if (
        needs.staff &&
        (plan.max_staff_per_storefront ?? 0) !== -1 &&
        (plan.max_staff_per_storefront ?? 0) < needs.staff
      ) {
        return false;
      }
      if (needs.ai && !plan.has_ai_assistant) {
        return false;
      }
      if (needs.customDomain && !plan.has_custom_domain) {
        return false;
      }
      return true;
    });

    // Return cheapest suitable plan
    return suitablePlans.sort((a, b) => {
      const priceA = a.price_monthly || 0;
      const priceB = b.price_monthly || 0;
      return priceA - priceB;
    })[0];
  }
}

// Export singleton instance
export const subscriptionService = new SubscriptionService();
