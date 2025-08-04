'use client';

import { useTranslations } from 'next-intl';
import { useState } from 'react';
import { loadMessages } from '@/lib/i18n/loadMessages';
import { useLocale } from 'next-intl';

/**
 * Пример компонента с динамической загрузкой модулей переводов
 */
export function ModularTranslationExample() {
  const locale = useLocale();
  const t = useTranslations();
  const [adminTranslations, setAdminTranslations] = useState<any>(null);
  const [loading, setLoading] = useState(false);

  // Пример динамической загрузки модуля
  const loadAdminModule = async () => {
    setLoading(true);
    try {
      const messages = await loadMessages(locale as any, ['admin']);
      setAdminTranslations(messages.admin);
    } catch (error) {
      console.error('Failed to load admin module:', error);
    } finally {
      setLoading(false);
    }
  };

  return (
    <div className="p-6 space-y-4">
      <h2 className="text-2xl font-bold">Пример модульных переводов</h2>

      {/* Базовые переводы (всегда доступны) */}
      <div className="card bg-base-100 shadow-xl">
        <div className="card-body">
          <h3 className="card-title">Базовые переводы (common)</h3>
          <p>{t('common.loading')}</p>
          <p>{t('common.save')}</p>
          <p>{t('common.cancel')}</p>
        </div>
      </div>

      {/* Динамически загружаемые переводы */}
      <div className="card bg-base-100 shadow-xl">
        <div className="card-body">
          <h3 className="card-title">Динамические переводы (admin)</h3>
          {!adminTranslations ? (
            <button
              className="btn btn-primary"
              onClick={loadAdminModule}
              disabled={loading}
            >
              {loading ? 'Загрузка...' : 'Загрузить админ переводы'}
            </button>
          ) : (
            <div>
              <p>Заголовок админки: {adminTranslations.title}</p>
              <p>Управление: {adminTranslations.manage}</p>
            </div>
          )}
        </div>
      </div>

      {/* Использование с namespace */}
      <div className="card bg-base-100 shadow-xl">
        <div className="card-body">
          <h3 className="card-title">Переводы с namespace</h3>
          <code className="text-sm bg-base-200 p-2 rounded">
            {`const t = useTranslations('marketplace');`}
          </code>
          <p className="mt-2">
            Это позволяет использовать короткие ключи: t(&apos;title&apos;)
            вместо t(&apos;marketplace.title&apos;)
          </p>
        </div>
      </div>
    </div>
  );
}
