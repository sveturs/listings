'use client';

import { useState, useEffect } from 'react';
import { useTranslations, useLocale } from 'next-intl';
import { XMarkIcon, Cog6ToothIcon } from '@heroicons/react/24/outline';

interface ChatSettingsProps {
  isOpen: boolean;
  onClose: () => void;
}

export default function ChatSettings({ isOpen, onClose }: ChatSettingsProps) {
  const t = useTranslations('chat');
  const locale = useLocale();

  // Настройки из localStorage
  const [autoTranslate, setAutoTranslate] = useState(false);
  const [preferredLanguage, setPreferredLanguage] = useState(locale);

  // Загрузка настроек из localStorage при монтировании
  useEffect(() => {
    const savedAutoTranslate = localStorage.getItem('chat_auto_translate');
    const savedLanguage = localStorage.getItem('chat_preferred_language');

    if (savedAutoTranslate !== null) {
      setAutoTranslate(savedAutoTranslate === 'true');
    }

    if (savedLanguage) {
      setPreferredLanguage(savedLanguage);
    } else {
      setPreferredLanguage(locale);
    }
  }, [locale]);

  // Сохранение настроек
  const handleAutoTranslateChange = (checked: boolean) => {
    setAutoTranslate(checked);
    localStorage.setItem('chat_auto_translate', checked.toString());

    // Событие для уведомления других компонентов
    window.dispatchEvent(
      new CustomEvent('chat-settings-changed', {
        detail: { autoTranslate: checked, preferredLanguage },
      })
    );
  };

  const handleLanguageChange = (lang: string) => {
    setPreferredLanguage(lang);
    localStorage.setItem('chat_preferred_language', lang);

    // Событие для уведомления других компонентов
    window.dispatchEvent(
      new CustomEvent('chat-settings-changed', {
        detail: { autoTranslate, preferredLanguage: lang },
      })
    );
  };

  if (!isOpen) return null;

  return (
    <>
      {/* Backdrop */}
      <div
        className="fixed inset-0 bg-black bg-opacity-50 z-40"
        onClick={onClose}
      />

      {/* Settings Modal */}
      <div className="fixed inset-y-0 right-0 w-full sm:w-96 bg-base-100 shadow-xl z-50 transform transition-transform duration-300">
        {/* Header */}
        <div className="flex items-center justify-between p-4 border-b border-base-300">
          <div className="flex items-center gap-2">
            <Cog6ToothIcon className="w-6 h-6" />
            <h2 className="text-lg font-semibold">
              {t('translation.translationSettings')}
            </h2>
          </div>
          <button
            onClick={onClose}
            className="btn btn-ghost btn-sm btn-circle"
            aria-label="Close"
          >
            <XMarkIcon className="w-5 h-5" />
          </button>
        </div>

        {/* Content */}
        <div className="p-4 space-y-6">
          {/* Auto-translate toggle */}
          <div className="form-control">
            <label className="label cursor-pointer">
              <span className="label-text font-medium">
                {t('translation.autoTranslate')}
              </span>
              <input
                type="checkbox"
                className="toggle toggle-primary"
                checked={autoTranslate}
                onChange={(e) => handleAutoTranslateChange(e.target.checked)}
              />
            </label>
            <p className="text-sm text-base-content/70 mt-1">
              Automatically translate incoming messages to your preferred
              language
            </p>
          </div>

          {/* Language selection */}
          <div className="form-control">
            <label className="label">
              <span className="label-text font-medium">
                {t('translation.targetLanguage')}
              </span>
            </label>
            <select
              className="select select-bordered w-full"
              value={preferredLanguage}
              onChange={(e) => handleLanguageChange(e.target.value)}
            >
              <option value="en">{t('translation.languages.en')}</option>
              <option value="ru">{t('translation.languages.ru')}</option>
              <option value="sr">{t('translation.languages.sr')}</option>
            </select>
            <p className="text-sm text-base-content/70 mt-1">
              Messages will be translated to this language
            </p>
          </div>

          {/* Info */}
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
              />
            </svg>
            <div className="text-sm">
              <p>Translations are powered by Claude AI</p>
              <p className="mt-1">
                You can always view the original message by clicking "Show
                original"
              </p>
            </div>
          </div>
        </div>
      </div>
    </>
  );
}
