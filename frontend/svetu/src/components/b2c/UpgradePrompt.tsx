'use client';

import { useTranslations } from 'next-intl';
import { useRouter } from 'next/navigation';
import {
  FiCheck,
  FiX,
  FiStar,
  FiTrendingUp,
  FiBriefcase,
  FiZap,
} from 'react-icons/fi';

interface UpgradePromptProps {
  currentPlan?: 'starter' | 'professional' | 'business' | 'enterprise';
  onClose?: () => void;
}

const plans = [
  {
    id: 'starter',
    name: 'Starter',
    price: 'Free',
    icon: FiStar,
    color: 'bg-base-300',
    features: {
      storefronts: 1,
      products: 50,
      staff: 1,
      images: 100,
      ai: false,
      live: false,
      export: false,
      customDomain: false,
    },
  },
  {
    id: 'professional',
    name: 'Professional',
    price: '€29/month',
    icon: FiTrendingUp,
    color: 'bg-primary',
    recommended: true,
    features: {
      storefronts: 3,
      products: 500,
      staff: 5,
      images: 1000,
      ai: false,
      live: false,
      export: true,
      customDomain: true,
    },
  },
  {
    id: 'business',
    name: 'Business',
    price: '€99/month',
    icon: FiBriefcase,
    color: 'bg-secondary',
    features: {
      storefronts: 10,
      products: 5000,
      staff: 20,
      images: 10000,
      ai: true,
      live: true,
      export: true,
      customDomain: true,
    },
  },
  {
    id: 'enterprise',
    name: 'Enterprise',
    price: 'Custom',
    icon: FiZap,
    color: 'bg-accent',
    features: {
      storefronts: -1, // unlimited
      products: -1,
      staff: -1,
      images: -1,
      ai: true,
      live: true,
      export: true,
      customDomain: true,
    },
  },
];

export default function UpgradePrompt({
  currentPlan = 'starter',
  onClose,
}: UpgradePromptProps) {
  const t = useTranslations('misc');
  const tCommon = useTranslations('common');
  const router = useRouter();

  const handleUpgrade = (planId: string) => {
    // В будущем здесь будет переход на страницу оплаты
    console.log('Upgrading to plan:', planId);
    router.push('/pricing');
    if (onClose) onClose();
  };

  const handleContactSales = () => {
    // В будущем здесь будет страница контактов
    console.log('Contact sales');
    if (onClose) onClose();
  };

  return (
    <div className="fixed inset-0 bg-black bg-opacity-50 z-50 flex items-center justify-center p-4">
      <div className="bg-base-100 rounded-2xl max-w-6xl w-full max-h-[90vh] overflow-y-auto shadow-2xl">
        {/* Header */}
        <div className="sticky top-0 bg-base-100 border-b border-base-300 p-6 flex justify-between items-center">
          <div>
            <h2 className="text-2xl font-bold">{t('title')}</h2>
            <p className="text-base-content/70 mt-1">{t('subtitle')}</p>
          </div>
          {onClose && (
            <button
              onClick={onClose}
              className="btn btn-ghost btn-circle"
              aria-label={tCommon('close')}
            >
              <FiX className="w-5 h-5" />
            </button>
          )}
        </div>

        {/* Content */}
        <div className="p-6">
          {/* Alert */}
          <div className="alert alert-warning mb-6">
            <svg
              xmlns="http://www.w3.org/2000/svg"
              className="stroke-current shrink-0 h-6 w-6"
              fill="none"
              viewBox="0 0 24 24"
            >
              <path
                strokeLinecap="round"
                strokeLinejoin="round"
                strokeWidth="2"
                d="M12 9v2m0 4h.01m-6.938 4h13.856c1.54 0 2.502-1.667 1.732-3L13.732 4c-.77-1.333-2.694-1.333-3.464 0L3.34 16c-.77 1.333.192 3 1.732 3z"
              />
            </svg>
            <span>{t('limitReached')}</span>
          </div>

          {/* Plans Grid */}
          <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-4 gap-4 mb-8">
            {plans.map((plan) => {
              const isCurrentPlan = plan.id === currentPlan;
              const Icon = plan.icon;

              return (
                <div
                  key={plan.id}
                  className={`card ${plan.color} relative ${plan.recommended ? 'ring-2 ring-primary' : ''}`}
                >
                  {plan.recommended && (
                    <div className="absolute -top-3 left-1/2 -translate-x-1/2">
                      <span className="badge badge-primary badge-sm px-3">
                        {t('recommended')}
                      </span>
                    </div>
                  )}

                  <div className="card-body">
                    {/* Plan Header */}
                    <div className="text-center mb-4">
                      <Icon className="w-8 h-8 mx-auto mb-2" />
                      <h3 className="text-xl font-bold">{plan.name}</h3>
                      <div className="text-2xl font-bold mt-2">
                        {plan.price}
                      </div>
                      {isCurrentPlan && (
                        <span className="badge badge-ghost badge-sm mt-2">
                          {t('currentPlan')}
                        </span>
                      )}
                    </div>

                    {/* Features */}
                    <ul className="space-y-2 mb-4 text-sm">
                      <li className="flex items-center gap-2">
                        <FiCheck className="text-success shrink-0" />
                        <span>
                          {plan.features.storefronts === -1
                            ? t('features.storefrontsUnlimited')
                            : t('features.storefronts', {
                                count: plan.features.storefronts,
                              })}
                        </span>
                      </li>
                      <li className="flex items-center gap-2">
                        <FiCheck className="text-success shrink-0" />
                        <span>
                          {plan.features.products === -1
                            ? t('features.productsUnlimited')
                            : t('features.products', {
                                count: plan.features.products,
                              })}
                        </span>
                      </li>
                      <li className="flex items-center gap-2">
                        <FiCheck className="text-success shrink-0" />
                        <span>
                          {plan.features.staff === -1
                            ? t('features.staffUnlimited')
                            : t('features.staff', {
                                count: plan.features.staff,
                              })}
                        </span>
                      </li>
                      <li className="flex items-center gap-2">
                        {plan.features.customDomain ? (
                          <FiCheck className="text-success shrink-0" />
                        ) : (
                          <FiX className="text-base-content/30 shrink-0" />
                        )}
                        <span
                          className={
                            !plan.features.customDomain ? 'opacity-50' : ''
                          }
                        >
                          {t('features.customDomain')}
                        </span>
                      </li>
                      <li className="flex items-center gap-2">
                        {plan.features.ai ? (
                          <FiCheck className="text-success shrink-0" />
                        ) : (
                          <FiX className="text-base-content/30 shrink-0" />
                        )}
                        <span className={!plan.features.ai ? 'opacity-50' : ''}>
                          {t('features.aiAssistant')}
                        </span>
                      </li>
                      <li className="flex items-center gap-2">
                        {plan.features.live ? (
                          <FiCheck className="text-success shrink-0" />
                        ) : (
                          <FiX className="text-base-content/30 shrink-0" />
                        )}
                        <span
                          className={!plan.features.live ? 'opacity-50' : ''}
                        >
                          {t('features.liveShopping')}
                        </span>
                      </li>
                    </ul>

                    {/* Action Button */}
                    {isCurrentPlan ? (
                      <button className="btn btn-disabled btn-block" disabled>
                        {t('currentPlan')}
                      </button>
                    ) : plan.id === 'enterprise' ? (
                      <button
                        onClick={handleContactSales}
                        className="btn btn-outline btn-block"
                      >
                        {t('contactSales')}
                      </button>
                    ) : (
                      <button
                        onClick={() => handleUpgrade(plan.id)}
                        className={`btn btn-block ${plan.recommended ? 'btn-primary' : 'btn-outline'}`}
                      >
                        {t('upgrade')}
                      </button>
                    )}
                  </div>
                </div>
              );
            })}
          </div>

          {/* Features Comparison */}
          <div className="text-center text-sm text-base-content/70">
            <p>{t('needHelp')}</p>
            <button onClick={handleContactSales} className="link link-primary">
              {t('contactSupport')}
            </button>
          </div>
        </div>
      </div>
    </div>
  );
}
