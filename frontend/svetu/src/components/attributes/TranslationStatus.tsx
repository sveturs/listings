'use client';

import React, { useEffect, useState } from 'react';
import { adminApi } from '@/services/admin';
import type { components } from '@/types/generated/api';
import { useTranslations } from 'next-intl';

interface TranslationStatusProps {
  entityType: 'category' | 'attribute';
  entityId: number;
  onTranslateClick?: () => void;
  compact?: boolean;
}

type TranslationStatus = components['schemas']['handler.TranslationStatusItem'];
type FieldStatus = components['schemas']['handler.TranslationFieldStatus'];

const LANGUAGES = ['en', 'ru', 'sr'] as const;
const LANGUAGE_LABELS: Record<string, string> = {
  en: 'English',
  ru: 'Русский',
  sr: 'Српски',
};

export function TranslationStatus({
  entityType,
  entityId,
  onTranslateClick,
  compact = false,
}: TranslationStatusProps) {
  const t = useTranslations('admin');
  const tCommon = useTranslations('common');
  const [status, setStatus] = useState<TranslationStatus | null>(null);
  const [loading, setLoading] = useState(true);

  useEffect(() => {
    fetchTranslationStatus();
    // eslint-disable-next-line react-hooks/exhaustive-deps
  }, [entityType, entityId]);

  const fetchTranslationStatus = async () => {
    try {
      setLoading(true);
      const statuses = await adminApi.getTranslationStatus(entityType, [
        entityId,
      ]);
      if (statuses && statuses.length > 0) {
        setStatus(statuses[0]);
      }
    } catch (error) {
      // Игнорируем ошибки 401 - это нормально если пользователь не залогинен
      if ((error as any)?.status !== 401) {
        console.error('Failed to fetch translation status:', error);
      }
    } finally {
      setLoading(false);
    }
  };

  if (loading) {
    return (
      <div className="flex items-center gap-2">
        <span className="loading loading-spinner loading-xs"></span>
        <span className="text-xs text-gray-500">{tCommon('loading')}</span>
      </div>
    );
  }

  if (!status) {
    return null;
  }

  const getStatusIcon = (fieldStatus: FieldStatus | undefined) => {
    if (!fieldStatus) {
      return '❌'; // Not translated
    }
    if (fieldStatus.is_verified) {
      return '✅'; // Verified translation
    }
    if (fieldStatus.is_machine_translated) {
      return '🤖'; // Machine translation
    }
    return '✏️'; // Manual translation, not verified
  };

  const getStatusColor = (fieldStatus: FieldStatus | undefined) => {
    if (!fieldStatus) {
      return 'text-error';
    }
    if (fieldStatus.is_verified) {
      return 'text-success';
    }
    if (fieldStatus.is_machine_translated) {
      return 'text-warning';
    }
    return 'text-info';
  };

  const getStatusTooltip = (
    fieldStatus: FieldStatus | undefined,
    lang: string
  ) => {
    const langLabel = LANGUAGE_LABELS[lang] || lang;
    if (!fieldStatus) {
      return `${langLabel}: ${t('translations.notTranslated')}`;
    }
    if (fieldStatus.is_verified) {
      return `${langLabel}: ${t('translations.verified')}`;
    }
    if (fieldStatus.is_machine_translated) {
      return `${langLabel}: ${t('translations.machineTranslated')}`;
    }
    return `${langLabel}: ${t('translations.manualTranslated')}`;
  };

  const allTranslated = LANGUAGES.every(
    (lang) => status.languages?.[lang]?.is_translated
  );

  if (compact) {
    return (
      <div className="flex items-center gap-1">
        {LANGUAGES.map((lang) => {
          const fieldStatus = status.languages?.[lang];
          return (
            <div
              key={lang}
              className={`tooltip ${getStatusColor(fieldStatus)}`}
              data-tip={getStatusTooltip(fieldStatus, lang)}
            >
              <span className="text-xs">{getStatusIcon(fieldStatus)}</span>
            </div>
          );
        })}
        {onTranslateClick && !allTranslated && (
          <button
            onClick={onTranslateClick}
            className="btn btn-ghost btn-xs"
            title={t('translations.translate')}
          >
            🌍
          </button>
        )}
      </div>
    );
  }

  return (
    <div className="p-2 bg-base-200 rounded-lg">
      <div className="flex items-center justify-between mb-2">
        <h4 className="text-sm font-medium">{t('translations.status')}</h4>
        {onTranslateClick && !allTranslated && (
          <button
            onClick={onTranslateClick}
            className="btn btn-ghost btn-xs"
            title={t('translations.translate')}
          >
            🌍 {t('translations.translate')}
          </button>
        )}
      </div>
      <div className="grid grid-cols-3 gap-2">
        {LANGUAGES.map((lang) => {
          const fieldStatus = status.languages?.[lang];
          return (
            <div key={lang} className="flex items-center gap-1">
              <span className={`text-lg ${getStatusColor(fieldStatus)}`}>
                {getStatusIcon(fieldStatus)}
              </span>
              <span className="text-xs">{LANGUAGE_LABELS[lang]}</span>
            </div>
          );
        })}
      </div>
      {allTranslated && (
        <div className="mt-2 text-xs text-success">
          {t('translations.allTranslated')}
        </div>
      )}
    </div>
  );
}
