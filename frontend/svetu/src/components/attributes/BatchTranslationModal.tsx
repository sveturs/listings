'use client';

import React, { useState } from 'react';
import { useTranslations } from 'next-intl';
import { adminApi } from '@/services/admin';

interface BatchTranslationModalProps {
  isOpen: boolean;
  onClose: () => void;
  entityType: 'category' | 'attribute';
  selectedIds: number[];
  selectedNames?: string[];
  onComplete?: () => void;
}

interface TranslationResult {
  id: number;
  name: string;
  status: 'pending' | 'translating' | 'success' | 'error';
  error?: string;
}

export function BatchTranslationModal({
  isOpen,
  onClose,
  entityType,
  selectedIds,
  selectedNames = [],
  onComplete,
}: BatchTranslationModalProps) {
  const t = useTranslations('admin');
  const [translating, setTranslating] = useState(false);
  const [results, setResults] = useState<TranslationResult[]>([]);
  const [overallProgress, setOverallProgress] = useState(0);

  const handleTranslate = async () => {
    if (selectedIds.length === 0) return;

    setTranslating(true);
    setResults(
      selectedIds.map((id, index) => ({
        id,
        name: selectedNames[index] || `${entityType} #${id}`,
        status: 'pending',
      }))
    );

    try {
      // Инициализируем результаты
      const initialResults = selectedIds.map((id, index) => ({
        id,
        name: selectedNames[index] || `${entityType} #${id}`,
        status: 'translating' as const,
      }));
      setResults(initialResults);

      // Выполняем массовый перевод
      const response = await (entityType === 'category'
        ? adminApi.batchTranslateCategories(selectedIds)
        : adminApi.batchTranslateAttributes(selectedIds));

      // Обрабатываем результаты
      // API возвращает success: true если запрос успешный
      if (response) {
        const finalResults = initialResults.map((result) => ({
          ...result,
          status: 'success' as const,
        }));
        setResults(finalResults);
      } else {
        // В случае ошибки помечаем все как ошибочные
        const finalResults = initialResults.map((result) => ({
          ...result,
          status: 'error' as const,
          error: t('translations.translationFailed'),
        }));
        setResults(finalResults);
      }

      setOverallProgress(100);

      // Обновляем данные после успешного перевода
      if (onComplete) {
        onComplete();
      }
    } catch (error) {
      // В случае общей ошибки помечаем все как ошибочные
      setResults((prev) =>
        prev.map((result) => ({
          ...result,
          status: 'error',
          error:
            error instanceof Error
              ? error.message
              : t('translations.unexpectedError'),
        }))
      );
      setOverallProgress(100);
    } finally {
      setTranslating(false);
    }
  };

  const getStatusIcon = (status: TranslationResult['status']) => {
    switch (status) {
      case 'pending':
        return '⏳';
      case 'translating':
        return '🔄';
      case 'success':
        return '✅';
      case 'error':
        return '❌';
    }
  };

  const getStatusColor = (status: TranslationResult['status']) => {
    switch (status) {
      case 'pending':
        return 'text-base-content/50';
      case 'translating':
        return 'text-info';
      case 'success':
        return 'text-success';
      case 'error':
        return 'text-error';
    }
  };

  const successCount = results.filter((r) => r.status === 'success').length;
  const errorCount = results.filter((r) => r.status === 'error').length;

  const handleRetryFailed = async () => {
    const failedIds = results
      .filter((r) => r.status === 'error')
      .map((r) => r.id);

    if (failedIds.length > 0) {
      // Сбрасываем статус ошибочных
      setResults((prev) =>
        prev.map((result) =>
          failedIds.includes(result.id)
            ? { ...result, status: 'pending', error: undefined }
            : result
        )
      );

      // Повторяем перевод только для неудачных
      await handleTranslate();
    }
  };

  return (
    <dialog className={`modal ${isOpen ? 'modal-open' : ''}`}>
      <div className="modal-box max-w-2xl">
        <h3 className="font-bold text-lg mb-4">
          {t('translations.batchTranslation')} -{' '}
          {t(entityType === 'category' ? 'categories' : 'attributes')}
        </h3>

        <div className="space-y-4">
          {/* Информация о выбранных элементах */}
          <div className="alert alert-info">
            <span>
              {t('translations.selected')}: {selectedIds.length}{' '}
              {t(entityType === 'category' ? 'categories' : 'attributes')}
            </span>
          </div>

          {/* Прогресс перевода */}
          {translating && (
            <div className="space-y-2">
              <div className="flex justify-between text-sm">
                <span>{t('translations.translating')}...</span>
                <span>{Math.round(overallProgress)}%</span>
              </div>
              <progress
                className="progress progress-primary w-full"
                value={overallProgress}
                max="100"
              ></progress>
            </div>
          )}

          {/* Результаты перевода */}
          {results.length > 0 && (
            <div className="max-h-64 overflow-y-auto">
              <table className="table table-sm">
                <thead>
                  <tr>
                    <th>{t('translations.name')}</th>
                    <th>{t('translations.status')}</th>
                  </tr>
                </thead>
                <tbody>
                  {results.map((result) => (
                    <tr key={result.id}>
                      <td>{result.name}</td>
                      <td>
                        <div className="flex items-center gap-2">
                          <span className={getStatusColor(result.status)}>
                            {getStatusIcon(result.status)}
                          </span>
                          {result.status === 'translating' && (
                            <span className="loading loading-spinner loading-xs"></span>
                          )}
                          {result.error && (
                            <span className="text-xs text-error">
                              {result.error}
                            </span>
                          )}
                        </div>
                      </td>
                    </tr>
                  ))}
                </tbody>
              </table>
            </div>
          )}

          {/* Статистика результатов */}
          {!translating && results.length > 0 && (
            <div className="stats shadow">
              <div className="stat">
                <div className="stat-title">{t('translations.successful')}</div>
                <div className="stat-value text-success">{successCount}</div>
              </div>
              <div className="stat">
                <div className="stat-title">{t('translations.failed')}</div>
                <div className="stat-value text-error">{errorCount}</div>
              </div>
            </div>
          )}
        </div>

        <div className="modal-action">
          {!translating && errorCount > 0 && (
            <button onClick={handleRetryFailed} className="btn btn-warning">
              🔄 {t('translations.retryFailed')}
            </button>
          )}
          {!translating && results.length === 0 && (
            <button onClick={handleTranslate} className="btn btn-primary">
              🌍 {t('translations.startTranslation')}
            </button>
          )}
          <button
            onClick={onClose}
            disabled={translating}
            className="btn btn-ghost"
          >
            {t('translations.close')}
          </button>
        </div>
      </div>
      <form method="dialog" className="modal-backdrop">
        <button onClick={onClose} disabled={translating}>
          close
        </button>
      </form>
    </dialog>
  );
}
