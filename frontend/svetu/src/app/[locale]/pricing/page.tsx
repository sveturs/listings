'use client';

import { useState } from 'react';
import { useTranslations } from 'next-intl';
import Link from 'next/link';
import { useRouter } from 'next/navigation';
import {
  FiStar,
  FiTrendingUp,
  FiBriefcase,
  FiZap,
  FiCheck,
  FiX,
  FiArrowLeft,
  FiLoader,
  FiAlertCircle,
} from 'react-icons/fi';
import { useSubscription } from '@/hooks/useSubscription';
import { subscriptionService } from '@/services/subscriptions';

export default function PricingPage() {
  const t = useTranslations('misc');
  const tCommon = useTranslations('common');
  const router = useRouter();

  const {
    plans,
    currentSubscription,
    currentPlan,
    isLoading,
    isProcessing,
    error,
    selectPlan,
    upgradePlan,
    canUseFeature: _canUseFeature,
  } = useSubscription();

  const [selectedCycle, setSelectedCycle] = useState<'monthly' | 'yearly'>(
    'monthly'
  );
  const [showUpgradeModal, setShowUpgradeModal] = useState(false);
  const [selectedPlanCode, setSelectedPlanCode] = useState<string | null>(null);

  // Map backend plans to UI display
  const getPlanIcon = (code: string) => {
    switch (code) {
      case 'starter':
        return FiStar;
      case 'professional':
        return FiTrendingUp;
      case 'business':
        return FiBriefcase;
      case 'enterprise':
        return FiZap;
      default:
        return FiStar;
    }
  };

  const getPlanColor = (code: string) => {
    switch (code) {
      case 'starter':
        return 'bg-base-300';
      case 'professional':
        return 'bg-primary';
      case 'business':
        return 'bg-secondary';
      case 'enterprise':
        return 'bg-accent';
      default:
        return 'bg-base-300';
    }
  };

  const handleSelectPlan = (planCode: string) => {
    if (currentPlan?.code === planCode) {
      // Already on this plan
      return;
    }

    setSelectedPlanCode(planCode);
    setShowUpgradeModal(true);
  };

  const confirmUpgrade = async () => {
    if (!selectedPlanCode) return;

    setShowUpgradeModal(false);

    if (currentSubscription) {
      // Upgrade existing subscription
      await upgradePlan(selectedPlanCode, selectedCycle);
    } else {
      // Create new subscription
      await selectPlan(selectedPlanCode, selectedCycle);
    }
  };

  const handleContactSales = () => {
    router.push('/contact?subject=enterprise');
  };

  const formatPrice = (plan: any) => {
    if (!plan) return 'Free';

    const price =
      selectedCycle === 'yearly' ? plan.price_yearly : plan.price_monthly;

    if (!price || price === 0) {
      return plan.code === 'enterprise' ? 'Custom' : 'Free';
    }

    const formattedPrice = subscriptionService.formatPrice(price, 'EUR');
    const period = selectedCycle === 'yearly' ? '/year' : '/month';

    return formattedPrice + period;
  };

  const getYearlySavings = (plan: any) => {
    if (!plan || !plan.price_monthly || !plan.price_yearly) return 0;

    const monthlyTotal = plan.price_monthly * 12;
    const yearlyPrice = plan.price_yearly;
    const savings = monthlyTotal - yearlyPrice;

    return savings > 0 ? Math.round((savings / monthlyTotal) * 100) : 0;
  };

  if (isLoading) {
    return (
      <div className="container mx-auto px-4 py-8">
        <div className="flex justify-center items-center min-h-[400px]">
          <FiLoader className="animate-spin w-8 h-8" />
        </div>
      </div>
    );
  }

  return (
    <div className="container mx-auto px-4 py-8">
      {/* Back button */}
      <Link href="/" className="btn btn-ghost btn-sm mb-6">
        <FiArrowLeft className="w-4 h-4 mr-2" />
        {tCommon('back')}
      </Link>

      {/* Header */}
      <div className="text-center mb-12">
        <h1 className="text-4xl font-bold mb-4">{t('title')}</h1>
        <p className="text-lg text-base-content/70 mb-6">{t('subtitle')}</p>

        {/* Billing cycle toggle */}
        <div className="flex justify-center items-center gap-4">
          <span className={selectedCycle === 'monthly' ? 'font-bold' : ''}>
            Monthly
          </span>
          <label className="swap swap-flip">
            <input
              type="checkbox"
              checked={selectedCycle === 'yearly'}
              onChange={(e) =>
                setSelectedCycle(e.target.checked ? 'yearly' : 'monthly')
              }
            />
            <div className="swap-on">
              <div className="badge badge-primary">Yearly</div>
            </div>
            <div className="swap-off">
              <div className="badge badge-outline">Monthly</div>
            </div>
          </label>
          <span className={selectedCycle === 'yearly' ? 'font-bold' : ''}>
            Yearly
            {selectedCycle === 'yearly' && (
              <span className="text-success ml-2">Save up to 20%</span>
            )}
          </span>
        </div>
      </div>

      {/* Error display */}
      {error && (
        <div className="alert alert-error mb-6">
          <FiAlertCircle className="w-5 h-5" />
          <span>{error}</span>
        </div>
      )}

      {/* Current subscription info */}
      {currentSubscription && currentPlan && (
        <div className="alert alert-info mb-6">
          <FiAlertCircle className="w-5 h-5" />
          <div>
            <p className="font-semibold">
              Your current plan: {currentPlan.name}
            </p>
            {currentSubscription.expires_at && (
              <p className="text-sm">
                Expires:{' '}
                {new Date(currentSubscription.expires_at).toLocaleDateString()}
              </p>
            )}
          </div>
        </div>
      )}

      {/* Plans Grid */}
      <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-4 gap-6 mb-12">
        {plans.map((plan) => {
          const Icon = getPlanIcon(plan.code || '');
          const isCurrentPlan = currentPlan?.code === plan.code;
          const canUpgrade =
            currentPlan &&
            (plan.sort_order ?? 0) > (currentPlan.sort_order ?? 0);
          const canDowngrade =
            currentPlan &&
            (plan.sort_order ?? 0) < (currentPlan.sort_order ?? 0);
          const savings = getYearlySavings(plan);

          return (
            <div
              key={plan.id}
              className={`card ${getPlanColor(plan.code || '')} shadow-xl relative ${
                plan.is_recommended
                  ? 'ring-2 ring-primary transform scale-105'
                  : ''
              } ${isCurrentPlan ? 'ring-2 ring-success' : ''}`}
            >
              {plan.is_recommended && (
                <div className="absolute -top-4 left-1/2 -translate-x-1/2">
                  <span className="badge badge-primary px-4 py-2">
                    {t('recommended')}
                  </span>
                </div>
              )}

              {isCurrentPlan && (
                <div className="absolute -top-4 right-4">
                  <span className="badge badge-success px-4 py-2">
                    Current Plan
                  </span>
                </div>
              )}

              {selectedCycle === 'yearly' && savings > 0 && (
                <div className="absolute top-2 right-2">
                  <span className="badge badge-warning">Save {savings}%</span>
                </div>
              )}

              <div className="card-body">
                {/* Plan Header */}
                <div className="text-center mb-6">
                  <Icon className="w-12 h-12 mx-auto mb-3" />
                  <h3 className="text-2xl font-bold">{plan.name}</h3>
                  <div className="text-3xl font-bold mt-3">
                    {formatPrice(plan)}
                  </div>
                  {(plan.free_trial_days ?? 0) > 0 && !currentSubscription && (
                    <p className="text-sm text-success mt-2">
                      {plan.free_trial_days} day free trial
                    </p>
                  )}
                </div>

                {/* Features */}
                <ul className="space-y-3 mb-6">
                  <li className="flex items-center gap-2">
                    <FiCheck className="text-success shrink-0 w-5 h-5" />
                    <span>
                      {(plan.max_b2c_stores ?? 0) === -1
                        ? t('features.storefrontsUnlimited')
                        : t('features.storefronts', {
                            count: plan.max_b2c_stores ?? 0,
                          })}
                    </span>
                  </li>
                  <li className="flex items-center gap-2">
                    <FiCheck className="text-success shrink-0 w-5 h-5" />
                    <span>
                      {(plan.max_products_per_storefront ?? 0) === -1
                        ? t('features.productsUnlimited')
                        : t('features.products', {
                            count: plan.max_products_per_storefront ?? 0,
                          })}
                    </span>
                  </li>
                  <li className="flex items-center gap-2">
                    <FiCheck className="text-success shrink-0 w-5 h-5" />
                    <span>
                      {(plan.max_staff_per_storefront ?? 0) === -1
                        ? t('features.staffUnlimited')
                        : t('features.staff', {
                            count: plan.max_staff_per_storefront ?? 0,
                          })}
                    </span>
                  </li>
                  <li className="flex items-center gap-2">
                    {plan.has_custom_domain ? (
                      <FiCheck className="text-success shrink-0 w-5 h-5" />
                    ) : (
                      <FiX className="text-base-content/30 shrink-0 w-5 h-5" />
                    )}
                    <span
                      className={!plan.has_custom_domain ? 'opacity-50' : ''}
                    >
                      {t('features.customDomain')}
                    </span>
                  </li>
                  <li className="flex items-center gap-2">
                    {plan.has_ai_assistant ? (
                      <FiCheck className="text-success shrink-0 w-5 h-5" />
                    ) : (
                      <FiX className="text-base-content/30 shrink-0 w-5 h-5" />
                    )}
                    <span
                      className={!plan.has_ai_assistant ? 'opacity-50' : ''}
                    >
                      {t('features.aiAssistant')}
                    </span>
                  </li>
                  <li className="flex items-center gap-2">
                    {plan.has_live_shopping ? (
                      <FiCheck className="text-success shrink-0 w-5 h-5" />
                    ) : (
                      <FiX className="text-base-content/30 shrink-0 w-5 h-5" />
                    )}
                    <span
                      className={!plan.has_live_shopping ? 'opacity-50' : ''}
                    >
                      {t('features.liveShopping')}
                    </span>
                  </li>
                  <li className="flex items-center gap-2">
                    {plan.has_export_data ? (
                      <FiCheck className="text-success shrink-0 w-5 h-5" />
                    ) : (
                      <FiX className="text-base-content/30 shrink-0 w-5 h-5" />
                    )}
                    <span className={!plan.has_export_data ? 'opacity-50' : ''}>
                      {t('features.exportData') || 'Export Data'}
                    </span>
                  </li>
                  <li className="text-sm text-base-content/70 text-center">
                    Marketplace fee: {plan.commission_rate}%
                  </li>
                </ul>

                {/* Action Button */}
                {isCurrentPlan ? (
                  <button className="btn btn-disabled btn-block" disabled>
                    {t('currentPlan')}
                  </button>
                ) : plan.code === 'enterprise' ? (
                  <button
                    onClick={handleContactSales}
                    className="btn btn-outline btn-block"
                    disabled={isProcessing}
                  >
                    {t('contactSales')}
                  </button>
                ) : canUpgrade ? (
                  <button
                    onClick={() => handleSelectPlan(plan.code || '')}
                    className={`btn btn-block ${plan.is_recommended ? 'btn-primary' : 'btn-outline'}`}
                    disabled={isProcessing}
                  >
                    {isProcessing ? (
                      <FiLoader className="animate-spin" />
                    ) : (
                      t('upgrade')
                    )}
                  </button>
                ) : canDowngrade ? (
                  <button
                    onClick={() => handleSelectPlan(plan.code || '')}
                    className="btn btn-outline btn-block"
                    disabled={isProcessing}
                  >
                    {isProcessing ? (
                      <FiLoader className="animate-spin" />
                    ) : (
                      'Downgrade'
                    )}
                  </button>
                ) : (
                  <button
                    onClick={() => handleSelectPlan(plan.code || '')}
                    className={`btn btn-block ${plan.is_recommended ? 'btn-primary' : 'btn-outline'}`}
                    disabled={isProcessing}
                  >
                    {isProcessing ? (
                      <FiLoader className="animate-spin" />
                    ) : (
                      'Get Started'
                    )}
                  </button>
                )}
              </div>
            </div>
          );
        })}
      </div>

      {/* FAQ or Additional Info */}
      <div className="text-center">
        <h2 className="text-2xl font-bold mb-4">{t('needHelp')}</h2>
        <p className="text-base-content/70 mb-4">{t('contactSupport')}</p>
        <button onClick={handleContactSales} className="btn btn-primary">
          {t('contactSales')}
        </button>
      </div>

      {/* Upgrade Confirmation Modal */}
      {showUpgradeModal && selectedPlanCode && (
        <div className="modal modal-open">
          <div className="modal-box">
            <h3 className="font-bold text-lg mb-4">Confirm Plan Change</h3>

            <div className="py-4">
              {currentPlan && (
                <div className="mb-4">
                  <p className="text-sm text-base-content/70">Current plan:</p>
                  <p className="font-semibold">{currentPlan.name}</p>
                </div>
              )}

              <div className="mb-4">
                <p className="text-sm text-base-content/70">New plan:</p>
                <p className="font-semibold">
                  {plans.find((p) => p.code === selectedPlanCode)?.name}
                </p>
              </div>

              <div className="mb-4">
                <p className="text-sm text-base-content/70">Billing cycle:</p>
                <p className="font-semibold capitalize">{selectedCycle}</p>
              </div>

              <div className="alert alert-warning">
                <FiAlertCircle className="w-5 h-5" />
                <span className="text-sm">
                  {currentPlan &&
                  (plans.find((p) => p.code === selectedPlanCode)?.sort_order ??
                    0) < (currentPlan.sort_order ?? 0)
                    ? 'Downgrading will take effect at the end of your current billing period.'
                    : 'You will be charged immediately and the new plan will be activated.'}
                </span>
              </div>
            </div>

            <div className="modal-action">
              <button
                onClick={() => setShowUpgradeModal(false)}
                className="btn btn-ghost"
                disabled={isProcessing}
              >
                Cancel
              </button>
              <button
                onClick={confirmUpgrade}
                className="btn btn-primary"
                disabled={isProcessing}
              >
                {isProcessing ? (
                  <FiLoader className="animate-spin" />
                ) : (
                  'Confirm'
                )}
              </button>
            </div>
          </div>
        </div>
      )}
    </div>
  );
}
