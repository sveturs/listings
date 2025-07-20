'use client';

import React, { useState, useEffect } from 'react';
import { useTranslations } from 'next-intl';

interface InlineTranslationEditorProps {
  entityType: 'category' | 'attribute';
  entityId: number;
  fieldName: string;
  translations: Record<string, string>;
  onSave: (translations: Record<string, string>) => Promise<void>;
  onCancel?: () => void;
  compact?: boolean;
}

const LANGUAGES = ['en', 'ru', 'sr'] as const;
const LANGUAGE_LABELS: Record<string, string> = {
  en: 'English',
  ru: 'Русский',
  sr: 'Српски',
};

export function InlineTranslationEditor({
  entityType: _entityType,
  entityId: _entityId,
  fieldName: _fieldName,
  translations: initialTranslations,
  onSave,
  onCancel,
  compact = false,
}: InlineTranslationEditorProps) {
  const t = useTranslations('admin.translations');
  const [isEditing, setIsEditing] = useState(false);
  const [translations, setTranslations] = useState(initialTranslations);
  const [saving, setSaving] = useState(false);
  const [error, setError] = useState<string | null>(null);

  useEffect(() => {
    setTranslations(initialTranslations);
  }, [initialTranslations]);

  const handleEdit = () => {
    setIsEditing(true);
    setError(null);
  };

  const handleCancel = () => {
    setIsEditing(false);
    setTranslations(initialTranslations);
    setError(null);
    onCancel?.();
  };

  const handleSave = async () => {
    try {
      setSaving(true);
      setError(null);
      await onSave(translations);
      setIsEditing(false);
    } catch (err) {
      setError(err instanceof Error ? err.message : t('saveFailed'));
    } finally {
      setSaving(false);
    }
  };

  const handleTranslationChange = (lang: string, value: string) => {
    setTranslations((prev) => ({
      ...prev,
      [lang]: value,
    }));
  };

  if (!isEditing) {
    return (
      <div className="group relative">
        <div
          className="cursor-pointer hover:bg-base-200 rounded p-1 transition-colors"
          onClick={handleEdit}
        >
          {compact ? (
            <div className="flex items-center gap-2">
              <span className="text-sm">
                {translations[LANGUAGES[0]] || t('notTranslated')}
              </span>
              <button
                className="btn btn-ghost btn-xs opacity-0 group-hover:opacity-100 transition-opacity"
                onClick={(e) => {
                  e.stopPropagation();
                  handleEdit();
                }}
              >
                ✏️
              </button>
            </div>
          ) : (
            <div className="space-y-1">
              {LANGUAGES.map((lang) => (
                <div key={lang} className="flex items-center gap-2">
                  <span className="text-xs font-medium w-12">
                    {lang.toUpperCase()}:
                  </span>
                  <span className="text-sm flex-1">
                    {translations[lang] || (
                      <span className="text-base-content/50">
                        {t('notTranslated')}
                      </span>
                    )}
                  </span>
                </div>
              ))}
              <button
                className="btn btn-ghost btn-xs w-full opacity-0 group-hover:opacity-100 transition-opacity"
                onClick={(e) => {
                  e.stopPropagation();
                  handleEdit();
                }}
              >
                ✏️ {t('edit')}
              </button>
            </div>
          )}
        </div>
      </div>
    );
  }

  return (
    <div className="p-2 bg-base-200 rounded-lg">
      <div className="space-y-2">
        {LANGUAGES.map((lang) => (
          <div key={lang}>
            <label className="label">
              <span className="label-text text-xs">
                {LANGUAGE_LABELS[lang]}
              </span>
            </label>
            <input
              type="text"
              value={translations[lang] || ''}
              onChange={(e) => handleTranslationChange(lang, e.target.value)}
              className="input input-bordered input-sm w-full"
              placeholder={t('enterTranslation')}
              disabled={saving}
            />
          </div>
        ))}
      </div>

      {error && (
        <div className="alert alert-error mt-2">
          <span className="text-xs">{error}</span>
        </div>
      )}

      <div className="flex gap-2 mt-3">
        <button
          onClick={handleSave}
          disabled={saving}
          className="btn btn-primary btn-sm flex-1"
        >
          {saving ? (
            <>
              <span className="loading loading-spinner loading-xs"></span>
              {t('saving')}
            </>
          ) : (
            <>💾 {t('save')}</>
          )}
        </button>
        <button
          onClick={handleCancel}
          disabled={saving}
          className="btn btn-ghost btn-sm flex-1"
        >
          {t('cancel')}
        </button>
      </div>
    </div>
  );
}
