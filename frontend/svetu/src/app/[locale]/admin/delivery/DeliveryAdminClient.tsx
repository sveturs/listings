'use client';

import { useState } from 'react';
import { useTranslations } from 'next-intl';
import { Link } from '@/i18n/routing';
import DeliveryProviders from './components/DeliveryProviders';
import PricingRules from './components/PricingRules';
import ProblemShipments from './components/ProblemShipments';
import DeliveryDashboard from './components/DeliveryDashboard';
import DeliveryAnalytics from './components/DeliveryAnalytics';
import DeliveryShipments from './components/DeliveryShipments';

export default function DeliveryAdminClient() {
  const t = useTranslations('admin.delivery');
  const [activeTab, setActiveTab] = useState('dashboard');

  const tabs = [
    { id: 'dashboard', label: t('tabs.dashboard'), icon: '📊' },
    { id: 'shipments', label: t('tabs.shipments'), icon: '📦' },
    { id: 'providers', label: t('tabs.providers'), icon: '🚚' },
    { id: 'pricingRules', label: t('tabs.pricingRules'), icon: '💰' },
    { id: 'problemShipments', label: t('tabs.problemShipments'), icon: '⚠️' },
    { id: 'analytics', label: t('tabs.analytics'), icon: '📈' },
    { id: 'testing', label: t('tabs.testing'), icon: '🧪' },
  ];

  return (
    <div className="container mx-auto px-4 py-8">
      <div className="flex items-center justify-between mb-8">
        <div>
          <h1 className="text-3xl font-bold">{t('title')}</h1>
          <p className="text-base-content/70 mt-2">{t('description')}</p>
        </div>
        <Link href="/admin" className="btn btn-ghost">
          ← Назад
        </Link>
      </div>

      {/* Tabs Navigation */}
      <div className="tabs tabs-boxed mb-8">
        {tabs.map((tab) => (
          <button
            key={tab.id}
            className={`tab tab-lg ${activeTab === tab.id ? 'tab-active' : ''}`}
            onClick={() => setActiveTab(tab.id)}
          >
            <span className="mr-2">{tab.icon}</span>
            {tab.label}
          </button>
        ))}
      </div>

      {/* Tab Content */}
      <div className="min-h-[600px]">
        {activeTab === 'dashboard' && <DeliveryDashboard />}
        {activeTab === 'shipments' && <DeliveryShipments />}
        {activeTab === 'providers' && <DeliveryProviders />}
        {activeTab === 'pricingRules' && <PricingRules />}
        {activeTab === 'problemShipments' && <ProblemShipments />}
        {activeTab === 'analytics' && <DeliveryAnalytics />}
        {activeTab === 'testing' && (
          <div className="card bg-base-100 shadow-xl">
            <div className="card-body">
              <h2 className="card-title text-2xl mb-4">
                <span className="mr-2">🧪</span>
                {t('testing.title')}
              </h2>
              <p className="text-base-content/70 mb-6">
                {t('testing.description')}
              </p>

              <div className="grid grid-cols-1 md:grid-cols-2 gap-4">
                {/* Post Express Testing */}
                <Link
                  href="/admin/postexpress/test"
                  className="card bg-base-200 hover:bg-base-300 transition-colors cursor-pointer"
                >
                  <div className="card-body">
                    <div className="flex items-center gap-3 mb-3">
                      <div className="w-12 h-12 rounded bg-primary/10 flex items-center justify-center">
                        <span className="text-2xl">📮</span>
                      </div>
                      <div>
                        <h3 className="font-semibold text-lg">Post Express</h3>
                        <div className="badge badge-success badge-sm">
                          {t('testing.configured')}
                        </div>
                      </div>
                    </div>
                    <p className="text-sm text-base-content/70">
                      {t('testing.postexpressDescription')}
                    </p>
                    <div className="card-actions justify-end mt-4">
                      <button className="btn btn-primary btn-sm">
                        {t('testing.openTest')} →
                      </button>
                    </div>
                  </div>
                </Link>

                {/* Placeholder for other providers */}
                <div className="card bg-base-200 opacity-50">
                  <div className="card-body">
                    <div className="flex items-center gap-3 mb-3">
                      <div className="w-12 h-12 rounded bg-base-300 flex items-center justify-center">
                        <span className="text-2xl">🚚</span>
                      </div>
                      <div>
                        <h3 className="font-semibold text-lg">
                          {t('testing.otherProviders')}
                        </h3>
                        <div className="badge badge-ghost badge-sm">
                          {t('testing.comingSoon')}
                        </div>
                      </div>
                    </div>
                    <p className="text-sm text-base-content/70">
                      {t('testing.otherProvidersDescription')}
                    </p>
                  </div>
                </div>
              </div>

              <div className="alert alert-info mt-6">
                <svg
                  xmlns="http://www.w3.org/2000/svg"
                  fill="none"
                  viewBox="0 0 24 24"
                  className="stroke-current shrink-0 w-6 h-6"
                >
                  <path
                    strokeLinecap="round"
                    strokeLinejoin="round"
                    strokeWidth="2"
                    d="M13 16h-1v-4h-1m1-4h.01M21 12a9 9 0 11-18 0 9 9 0 0118 0z"
                  ></path>
                </svg>
                <span>{t('testing.infoMessage')}</span>
              </div>
            </div>
          </div>
        )}
      </div>
    </div>
  );
}
