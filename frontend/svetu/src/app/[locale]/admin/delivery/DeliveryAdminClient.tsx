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
      </div>
    </div>
  );
}
