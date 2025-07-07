'use client';

import { useState, useEffect } from 'react';
import { useTranslations } from 'next-intl';
import { toast } from 'react-hot-toast';

interface TransliterationRule {
  id: string;
  from: string;
  to: string;
  language: 'sr' | 'en';
  active: boolean;
}

export default function TransliterationConfig() {
  const t = useTranslations();
  const [rules, setRules] = useState<TransliterationRule[]>([]);
  const [loading, setLoading] = useState(true);
  const [saving, setSaving] = useState(false);
  const [newRule, setNewRule] = useState<Omit<TransliterationRule, 'id'>>({
    from: '',
    to: '',
    language: 'sr',
    active: true,
  });

  useEffect(() => {
    fetchRules();
  }, []);

  const fetchRules = async () => {
    try {
      setLoading(true);
      // Mock data for now since API may not be implemented
      const mockRules: TransliterationRule[] = [
        { id: '1', from: 'č', to: 'c', language: 'sr', active: true },
        { id: '2', from: 'ć', to: 'c', language: 'sr', active: true },
        { id: '3', from: 'š', to: 's', language: 'sr', active: true },
        { id: '4', from: 'ž', to: 'z', language: 'sr', active: true },
        { id: '5', from: 'đ', to: 'd', language: 'sr', active: true },
      ];
      setRules(mockRules);
    } catch (error) {
      console.error('Error fetching transliteration rules:', error);
      toast.error(t('admin.search.transliteration.fetchError'));
    } finally {
      setLoading(false);
    }
  };

  const handleAddRule = async () => {
    if (!newRule.from || !newRule.to) {
      toast.error(t('admin.search.transliteration.fillFields'));
      return;
    }

    try {
      setSaving(true);
      // Mock save - in real implementation, this would call an API
      const rule: TransliterationRule = {
        ...newRule,
        id: Date.now().toString(),
      };
      setRules([...rules, rule]);
      setNewRule({ from: '', to: '', language: 'sr', active: true });
      toast.success(t('admin.search.transliteration.ruleAdded'));
    } catch (error) {
      console.error('Error adding rule:', error);
      toast.error(t('admin.search.transliteration.saveError'));
    } finally {
      setSaving(false);
    }
  };

  const handleDeleteRule = async (id: string) => {
    try {
      setRules(rules.filter((rule) => rule.id !== id));
      toast.success(t('admin.search.transliteration.ruleDeleted'));
    } catch (error) {
      console.error('Error deleting rule:', error);
      toast.error(t('admin.search.transliteration.deleteError'));
    }
  };

  const handleToggleRule = async (id: string) => {
    try {
      setRules(
        rules.map((rule) =>
          rule.id === id ? { ...rule, active: !rule.active } : rule
        )
      );
      toast.success(t('admin.search.transliteration.ruleUpdated'));
    } catch (error) {
      console.error('Error updating rule:', error);
      toast.error(t('admin.search.transliteration.updateError'));
    }
  };

  if (loading) {
    return (
      <div className="flex justify-center items-center h-64">
        <div className="loading loading-spinner loading-lg"></div>
      </div>
    );
  }

  return (
    <div className="space-y-6">
      <div className="card bg-base-100 shadow-xl">
        <div className="card-body">
          <h2 className="card-title">
            {t('admin.search.transliteration.addRule')}
          </h2>
          <div className="grid grid-cols-1 md:grid-cols-4 gap-4">
            <div className="form-control">
              <label className="label">
                <span className="label-text">
                  {t('admin.search.transliteration.from')}
                </span>
              </label>
              <input
                type="text"
                className="input input-bordered"
                value={newRule.from}
                onChange={(e) =>
                  setNewRule({ ...newRule, from: e.target.value })
                }
                placeholder="č"
              />
            </div>
            <div className="form-control">
              <label className="label">
                <span className="label-text">
                  {t('admin.search.transliteration.to')}
                </span>
              </label>
              <input
                type="text"
                className="input input-bordered"
                value={newRule.to}
                onChange={(e) => setNewRule({ ...newRule, to: e.target.value })}
                placeholder="c"
              />
            </div>
            <div className="form-control">
              <label className="label">
                <span className="label-text">
                  {t('admin.search.transliteration.language')}
                </span>
              </label>
              <select
                className="select select-bordered"
                value={newRule.language}
                onChange={(e) =>
                  setNewRule({
                    ...newRule,
                    language: e.target.value as 'sr' | 'en',
                  })
                }
              >
                <option value="sr">Српски</option>
                <option value="en">English</option>
              </select>
            </div>
            <div className="form-control">
              <label className="label">
                <span className="label-text">&nbsp;</span>
              </label>
              <button
                className={`btn btn-primary ${saving ? 'loading' : ''}`}
                onClick={handleAddRule}
                disabled={saving}
              >
                {saving ? '' : t('admin.search.transliteration.addRule')}
              </button>
            </div>
          </div>
        </div>
      </div>

      <div className="card bg-base-100 shadow-xl">
        <div className="card-body">
          <h2 className="card-title">
            {t('admin.search.transliteration.existingRules')}
          </h2>
          <div className="overflow-x-auto">
            <table className="table table-zebra">
              <thead>
                <tr>
                  <th>{t('admin.search.transliteration.from')}</th>
                  <th>{t('admin.search.transliteration.to')}</th>
                  <th>{t('admin.search.transliteration.language')}</th>
                  <th>{t('admin.search.transliteration.status')}</th>
                  <th>{t('admin.search.transliteration.actions')}</th>
                </tr>
              </thead>
              <tbody>
                {rules.map((rule) => (
                  <tr key={rule.id}>
                    <td className="font-mono text-lg">{rule.from}</td>
                    <td className="font-mono text-lg">{rule.to}</td>
                    <td>
                      <div className="badge badge-outline">
                        {rule.language === 'sr' ? 'Српски' : 'English'}
                      </div>
                    </td>
                    <td>
                      <input
                        type="checkbox"
                        className="toggle toggle-primary"
                        checked={rule.active}
                        onChange={() => handleToggleRule(rule.id)}
                      />
                    </td>
                    <td>
                      <button
                        className="btn btn-error btn-sm"
                        onClick={() => handleDeleteRule(rule.id)}
                      >
                        {t('admin.search.transliteration.delete')}
                      </button>
                    </td>
                  </tr>
                ))}
              </tbody>
            </table>
          </div>
        </div>
      </div>

      <div className="alert alert-info">
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
        <span>{t('admin.search.transliteration.info')}</span>
      </div>
    </div>
  );
}
