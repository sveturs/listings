'use client';

import { useEffect, useState } from 'react';
import { useTranslations } from 'next-intl';
import { tokenManager } from '@/utils/tokenManager';
import configManager from '@/config';

interface Provider {
  id: number;
  code: string;
  name: string;
  logo_url?: string;
  is_active: boolean;
  supports_cod: boolean;
  supports_insurance: boolean;
  supports_tracking: boolean;
  api_configured: boolean;
}

export default function DeliveryProviders() {
  const t = useTranslations('admin.delivery.providers');
  const [providers, setProviders] = useState<Provider[]>([]);
  const [loading, setLoading] = useState(true);
  const [configModalOpen, setConfigModalOpen] = useState(false);
  const [selectedProvider, setSelectedProvider] = useState<Provider | null>(
    null
  );

  useEffect(() => {
    fetchProviders();
  }, []);

  const fetchProviders = async () => {
    try {
      const token = tokenManager.getAccessToken();
      const headers: HeadersInit = {
        'Content-Type': 'application/json',
      };

      if (token) {
        headers['Authorization'] = `Bearer ${token}`;
      }

      const response = await fetch(
        `${configManager.getApiUrl()}/api/v1/admin/delivery/providers`,
        {
          credentials: 'include',
          headers,
        }
      );

      if (response.ok) {
        const data = await response.json();
        setProviders(data.data || []);
      } else {
        // Использовать mock данные если API не доступен
        setProviders([
          {
            id: 1,
            code: 'postexpress',
            name: 'Post Express',
            logo_url: '/images/providers/postexpress.png',
            is_active: true,
            supports_cod: true,
            supports_insurance: true,
            supports_tracking: true,
            api_configured: true,
          },
          {
            id: 2,
            code: 'bex',
            name: 'BEX Express',
            logo_url: '/images/providers/bex.png',
            is_active: true,
            supports_cod: true,
            supports_insurance: false,
            supports_tracking: true,
            api_configured: false,
          },
          {
            id: 3,
            code: 'aks',
            name: 'AKS',
            logo_url: '/images/providers/aks.png',
            is_active: false,
            supports_cod: true,
            supports_insurance: true,
            supports_tracking: true,
            api_configured: false,
          },
          {
            id: 4,
            code: 'dexpress',
            name: 'D Express',
            logo_url: '/images/providers/dexpress.png',
            is_active: false,
            supports_cod: false,
            supports_insurance: true,
            supports_tracking: true,
            api_configured: false,
          },
          {
            id: 5,
            code: 'cityexpress',
            name: 'City Express',
            logo_url: '/images/providers/cityexpress.png',
            is_active: false,
            supports_cod: true,
            supports_insurance: false,
            supports_tracking: true,
            api_configured: false,
          },
          {
            id: 6,
            code: 'dhl',
            name: 'DHL',
            logo_url: '/images/providers/dhl.png',
            is_active: false,
            supports_cod: false,
            supports_insurance: true,
            supports_tracking: true,
            api_configured: false,
          },
        ]);
      }
    } catch (error) {
      console.error('Failed to fetch providers:', error);
    } finally {
      setLoading(false);
    }
  };

  const toggleProviderStatus = async (provider: Provider) => {
    try {
      const token = tokenManager.getAccessToken();
      const headers: HeadersInit = {
        'Content-Type': 'application/json',
      };

      if (token) {
        headers['Authorization'] = `Bearer ${token}`;
      }

      const response = await fetch(
        `/api/v1/admin/delivery/providers/${provider.id}/toggle`,
        {
          method: 'POST',
          credentials: 'include',
          headers,
        }
      );

      if (response.ok) {
        setProviders((prev) =>
          prev.map((p) =>
            p.id === provider.id ? { ...p, is_active: !p.is_active } : p
          )
        );
      }
    } catch (error) {
      console.error('Failed to toggle provider status:', error);
      // В случае ошибки, все равно обновляем локально для демонстрации
      setProviders((prev) =>
        prev.map((p) =>
          p.id === provider.id ? { ...p, is_active: !p.is_active } : p
        )
      );
    }
  };

  const openConfigModal = (provider: Provider) => {
    setSelectedProvider(provider);
    setConfigModalOpen(true);
  };

  if (loading) {
    return (
      <div className="flex items-center justify-center h-64">
        <span className="loading loading-spinner loading-lg"></span>
      </div>
    );
  }

  return (
    <div>
      <div className="mb-6">
        <h2 className="text-xl font-semibold mb-2">{t('title')}</h2>
        <p className="text-base-content/70">{t('description')}</p>
      </div>

      <div className="overflow-x-auto">
        <table className="table table-zebra w-full">
          <thead>
            <tr>
              <th>{t('table.logo')}</th>
              <th>{t('table.code')}</th>
              <th>{t('table.name')}</th>
              <th>{t('table.status')}</th>
              <th>{t('table.capabilities')}</th>
              <th>{t('table.actions')}</th>
            </tr>
          </thead>
          <tbody>
            {providers.map((provider) => (
              <tr key={provider.id}>
                <td>
                  <div className="avatar">
                    <div className="w-12 h-12 rounded bg-base-200 flex items-center justify-center">
                      {provider.logo_url ? (
                        <img
                          src={provider.logo_url}
                          alt={provider.name}
                          className="w-10 h-10 object-contain"
                        />
                      ) : (
                        <span className="text-2xl">📦</span>
                      )}
                    </div>
                  </div>
                </td>
                <td>
                  <code className="text-sm">{provider.code}</code>
                </td>
                <td className="font-medium">{provider.name}</td>
                <td>
                  <div className="flex flex-col gap-1">
                    <span
                      className={`badge ${provider.is_active ? 'badge-success' : 'badge-ghost'}`}
                    >
                      {provider.is_active
                        ? t('status.active')
                        : t('status.inactive')}
                    </span>
                    {provider.api_configured && (
                      <span className="badge badge-info badge-sm">
                        API настроен
                      </span>
                    )}
                  </div>
                </td>
                <td>
                  <div className="flex flex-wrap gap-1">
                    {provider.supports_cod && (
                      <span
                        className="badge badge-sm"
                        title={t('capabilities.cod')}
                      >
                        COD
                      </span>
                    )}
                    {provider.supports_insurance && (
                      <span
                        className="badge badge-sm"
                        title={t('capabilities.insurance')}
                      >
                        📋
                      </span>
                    )}
                    {provider.supports_tracking && (
                      <span
                        className="badge badge-sm"
                        title={t('capabilities.tracking')}
                      >
                        📍
                      </span>
                    )}
                  </div>
                </td>
                <td>
                  <div className="flex gap-2">
                    <button
                      className="btn btn-ghost btn-xs"
                      onClick={() => openConfigModal(provider)}
                    >
                      {t('actions.configure')}
                    </button>
                    <button
                      className={`btn btn-xs ${provider.is_active ? 'btn-error' : 'btn-success'}`}
                      onClick={() => toggleProviderStatus(provider)}
                    >
                      {provider.is_active
                        ? t('actions.deactivate')
                        : t('actions.activate')}
                    </button>
                    <button className="btn btn-ghost btn-xs">
                      {t('actions.test')}
                    </button>
                  </div>
                </td>
              </tr>
            ))}
          </tbody>
        </table>
      </div>

      {/* Configuration Modal */}
      {configModalOpen && selectedProvider && (
        <dialog className="modal modal-open">
          <div className="modal-box">
            <h3 className="font-bold text-lg mb-4">
              {t('configModal.title')}: {selectedProvider.name}
            </h3>

            <div className="space-y-4">
              <div className="form-control">
                <label className="label">
                  <span className="label-text">{t('configModal.apiUrl')}</span>
                </label>
                <input
                  type="text"
                  className="input input-bordered"
                  placeholder="https://api.provider.com"
                />
              </div>

              <div className="form-control">
                <label className="label">
                  <span className="label-text">{t('configModal.apiKey')}</span>
                </label>
                <input
                  type="text"
                  className="input input-bordered"
                  placeholder="your-api-key"
                />
              </div>

              <div className="form-control">
                <label className="label">
                  <span className="label-text">
                    {t('configModal.apiSecret')}
                  </span>
                </label>
                <input
                  type="password"
                  className="input input-bordered"
                  placeholder="your-api-secret"
                />
              </div>

              <div className="form-control">
                <label className="label cursor-pointer">
                  <span className="label-text">
                    {t('configModal.testMode')}
                  </span>
                  <input type="checkbox" className="toggle toggle-primary" />
                </label>
              </div>

              <button className="btn btn-sm btn-outline btn-block">
                {t('configModal.testConnection')}
              </button>
            </div>

            <div className="modal-action">
              <button className="btn btn-primary">
                {t('configModal.save')}
              </button>
              <button className="btn" onClick={() => setConfigModalOpen(false)}>
                {t('configModal.cancel')}
              </button>
            </div>
          </div>
          <form
            method="dialog"
            className="modal-backdrop"
            onClick={() => setConfigModalOpen(false)}
          >
            <button>close</button>
          </form>
        </dialog>
      )}
    </div>
  );
}
