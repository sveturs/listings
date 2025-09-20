'use client';

import { useEffect, useState } from 'react';
import { useTranslations } from 'next-intl';
import { tokenManager } from '@/utils/tokenManager';

interface PricingRule {
  id: number;
  provider_id: number;
  provider_name: string;
  rule_type: string;
  priority: number;
  is_active: boolean;
  min_price: number;
  max_price: number;
  weight_ranges?: Array<{from: number; to: number; price_per_kg: number}>;
  volume_ranges?: Array<{from: number; to: number; price_per_m3: number}>;
  fragile_surcharge?: number;
  oversized_surcharge?: number;
  special_handling_surcharge?: number;
}

export default function PricingRules() {
  const t = useTranslations('admin.delivery.pricingRules');
  const [rules, setRules] = useState<PricingRule[]>([]);
  const [loading, setLoading] = useState(true);
  const [editModalOpen, setEditModalOpen] = useState(false);
  const [selectedRule, setSelectedRule] = useState<PricingRule | null>(null);

  useEffect(() => {
    fetchPricingRules();
  }, []);

  const fetchPricingRules = async () => {
    try {
      // Mock data для демонстрации
      setTimeout(() => {
        setRules([
          {
            id: 1,
            provider_id: 1,
            provider_name: 'Post Express',
            rule_type: 'weight_based',
            priority: 1,
            is_active: true,
            min_price: 5,
            max_price: 500,
            weight_ranges: [
              { from: 0, to: 1, price_per_kg: 5 },
              { from: 1, to: 5, price_per_kg: 4 },
              { from: 5, to: 10, price_per_kg: 3.5 },
              { from: 10, to: 50, price_per_kg: 3 },
              { from: 50, to: 100, price_per_kg: 2.5 },
            ],
            fragile_surcharge: 2,
            oversized_surcharge: 5
          },
          {
            id: 2,
            provider_id: 2,
            provider_name: 'BEX Express',
            rule_type: 'zone_based',
            priority: 1,
            is_active: true,
            min_price: 4,
            max_price: 400
          },
          {
            id: 3,
            provider_id: 3,
            provider_name: 'AKS',
            rule_type: 'combined',
            priority: 2,
            is_active: true,
            min_price: 6,
            max_price: 450
          },
          {
            id: 4,
            provider_id: 4,
            provider_name: 'D Express',
            rule_type: 'volume_based',
            priority: 1,
            is_active: false,
            min_price: 8,
            max_price: 350,
            volume_ranges: [
              { from: 0, to: 0.01, price_per_m3: 100 },
              { from: 0.01, to: 0.05, price_per_m3: 90 },
              { from: 0.05, to: 0.1, price_per_m3: 80 }
            ]
          }
        ]);
        setLoading(false);
      }, 500);
    } catch (error) {
      console.error('Failed to fetch pricing rules:', error);
      setLoading(false);
    }
  };

  const openEditModal = (rule: PricingRule) => {
    setSelectedRule(rule);
    setEditModalOpen(true);
  };

  const toggleRuleStatus = (ruleId: number) => {
    setRules(prev => prev.map(rule =>
      rule.id === ruleId ? { ...rule, is_active: !rule.is_active } : rule
    ));
  };

  const getRuleTypeLabel = (type: string) => {
    const labels: Record<string, string> = {
      'weight_based': t('ruleTypes.weight_based'),
      'volume_based': t('ruleTypes.volume_based'),
      'zone_based': t('ruleTypes.zone_based'),
      'combined': t('ruleTypes.combined')
    };
    return labels[type] || type;
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
      <div className="mb-6 flex justify-between items-center">
        <div>
          <h2 className="text-xl font-semibold mb-2">{t('title')}</h2>
          <p className="text-base-content/70">{t('description')}</p>
        </div>
        <button className="btn btn-primary">
          {t('addRule')}
        </button>
      </div>

      <div className="overflow-x-auto">
        <table className="table table-zebra w-full">
          <thead>
            <tr>
              <th>{t('table.provider')}</th>
              <th>{t('table.ruleType')}</th>
              <th>{t('table.priority')}</th>
              <th>{t('table.minPrice')}</th>
              <th>{t('table.maxPrice')}</th>
              <th>{t('table.status')}</th>
              <th>{t('table.actions')}</th>
            </tr>
          </thead>
          <tbody>
            {rules.map((rule) => (
              <tr key={rule.id}>
                <td className="font-medium">{rule.provider_name}</td>
                <td>
                  <span className="badge badge-outline">
                    {getRuleTypeLabel(rule.rule_type)}
                  </span>
                </td>
                <td>
                  <span className="text-lg font-semibold">{rule.priority}</span>
                </td>
                <td>€{rule.min_price}</td>
                <td>€{rule.max_price}</td>
                <td>
                  <input
                    type="checkbox"
                    className="toggle toggle-success"
                    checked={rule.is_active}
                    onChange={() => toggleRuleStatus(rule.id)}
                  />
                </td>
                <td>
                  <div className="flex gap-2">
                    <button
                      className="btn btn-ghost btn-xs"
                      onClick={() => openEditModal(rule)}
                    >
                      Редактировать
                    </button>
                    <button className="btn btn-ghost btn-xs text-error">
                      Удалить
                    </button>
                  </div>
                </td>
              </tr>
            ))}
          </tbody>
        </table>
      </div>

      {/* Edit Modal */}
      {editModalOpen && selectedRule && (
        <dialog className="modal modal-open">
          <div className="modal-box max-w-3xl">
            <h3 className="font-bold text-lg mb-4">
              {t('editModal.title')}: {selectedRule.provider_name}
            </h3>

            <div className="space-y-6">
              {/* Weight Ranges */}
              {selectedRule.weight_ranges && (
                <div>
                  <h4 className="font-semibold mb-3">{t('editModal.weightRanges')}</h4>
                  <div className="space-y-2">
                    {selectedRule.weight_ranges.map((range, index) => (
                      <div key={index} className="flex gap-3 items-center">
                        <input
                          type="number"
                          className="input input-bordered input-sm w-24"
                          defaultValue={range.from}
                          placeholder={t('editModal.from')}
                        />
                        <span>-</span>
                        <input
                          type="number"
                          className="input input-bordered input-sm w-24"
                          defaultValue={range.to}
                          placeholder={t('editModal.to')}
                        />
                        <span className="text-sm">кг</span>
                        <input
                          type="number"
                          className="input input-bordered input-sm w-24"
                          defaultValue={range.price_per_kg}
                          placeholder={t('editModal.pricePerKg')}
                        />
                        <span className="text-sm">€/кг</span>
                        <button className="btn btn-ghost btn-xs text-error">✕</button>
                      </div>
                    ))}
                    <button className="btn btn-sm btn-outline">
                      {t('editModal.addRange')}
                    </button>
                  </div>
                </div>
              )}

              {/* Volume Ranges */}
              {selectedRule.volume_ranges && (
                <div>
                  <h4 className="font-semibold mb-3">{t('editModal.volumeRanges')}</h4>
                  <div className="space-y-2">
                    {selectedRule.volume_ranges.map((range, index) => (
                      <div key={index} className="flex gap-3 items-center">
                        <input
                          type="number"
                          className="input input-bordered input-sm w-24"
                          defaultValue={range.from}
                          placeholder={t('editModal.from')}
                        />
                        <span>-</span>
                        <input
                          type="number"
                          className="input input-bordered input-sm w-24"
                          defaultValue={range.to}
                          placeholder={t('editModal.to')}
                        />
                        <span className="text-sm">м³</span>
                        <input
                          type="number"
                          className="input input-bordered input-sm w-24"
                          defaultValue={range.price_per_m3}
                          placeholder={t('editModal.pricePerM3')}
                        />
                        <span className="text-sm">€/м³</span>
                        <button className="btn btn-ghost btn-xs text-error">✕</button>
                      </div>
                    ))}
                  </div>
                </div>
              )}

              {/* Surcharges */}
              <div>
                <h4 className="font-semibold mb-3">{t('editModal.surcharges')}</h4>
                <div className="grid grid-cols-1 md:grid-cols-3 gap-4">
                  <div className="form-control">
                    <label className="label">
                      <span className="label-text text-sm">{t('editModal.fragile')}</span>
                    </label>
                    <input
                      type="number"
                      className="input input-bordered input-sm"
                      defaultValue={selectedRule.fragile_surcharge || 0}
                      placeholder="0"
                    />
                  </div>
                  <div className="form-control">
                    <label className="label">
                      <span className="label-text text-sm">{t('editModal.oversized')}</span>
                    </label>
                    <input
                      type="number"
                      className="input input-bordered input-sm"
                      defaultValue={selectedRule.oversized_surcharge || 0}
                      placeholder="0"
                    />
                  </div>
                  <div className="form-control">
                    <label className="label">
                      <span className="label-text text-sm">{t('editModal.specialHandling')}</span>
                    </label>
                    <input
                      type="number"
                      className="input input-bordered input-sm"
                      defaultValue={selectedRule.special_handling_surcharge || 0}
                      placeholder="0"
                    />
                  </div>
                </div>
              </div>

              {/* Custom Formula */}
              <div className="form-control">
                <label className="label">
                  <span className="label-text">{t('editModal.formula')}</span>
                </label>
                <textarea
                  className="textarea textarea-bordered h-24"
                  placeholder="(weight * 2.5) + (volume * 100) + base_price"
                />
              </div>
            </div>

            <div className="modal-action">
              <button className="btn btn-primary">Сохранить</button>
              <button className="btn" onClick={() => setEditModalOpen(false)}>
                Отмена
              </button>
            </div>
          </div>
          <form method="dialog" className="modal-backdrop" onClick={() => setEditModalOpen(false)}>
            <button>close</button>
          </form>
        </dialog>
      )}
    </div>
  );
}