'use client';

import { useState, useEffect, useRef } from 'react';
import { useTranslations } from 'next-intl';
import {
  CheckIcon,
  XMarkIcon,
  PencilIcon,
  LanguageIcon,
  SparklesIcon,
} from '@heroicons/react/24/outline';
import { translationAdminApi } from '@/services/translationAdminApi';
import type { Translation } from '@/types/translation';

interface InlineTranslationEditorProps {
  translation: Translation;
  language: 'en' | 'ru' | 'sr';
  onSave?: (translation: Translation) => void;
  onCancel?: () => void;
  className?: string;
  autoFocus?: boolean;
  showAIButton?: boolean;
}

export default function InlineTranslationEditor({
  translation,
  language,
  onSave,
  onCancel,
  className = '',
  autoFocus = false,
  showAIButton = true,
}: InlineTranslationEditorProps) {
  const t = useTranslations('admin.translations.editor');

  const [isEditing, setIsEditing] = useState(false);
  const [value, setValue] = useState('');
  const [originalValue, setOriginalValue] = useState('');
  const [isSaving, setIsSaving] = useState(false);
  const [isGeneratingAI, setIsGeneratingAI] = useState(false);
  const textareaRef = useRef<HTMLTextAreaElement>(null);

  // Initialize value based on language
  useEffect(() => {
    const currentValue =
      language === 'en'
        ? translation.value_en
        : language === 'ru'
          ? translation.value_ru
          : translation.value_sr;

    setValue(currentValue || '');
    setOriginalValue(currentValue || '');
  }, [translation, language]);

  // Auto-focus when editing starts
  useEffect(() => {
    if (isEditing && autoFocus && textareaRef.current) {
      textareaRef.current.focus();
      textareaRef.current.select();
    }
  }, [isEditing, autoFocus]);

  // Auto-resize textarea
  const adjustTextareaHeight = () => {
    if (textareaRef.current) {
      textareaRef.current.style.height = 'auto';
      textareaRef.current.style.height = `${textareaRef.current.scrollHeight}px`;
    }
  };

  useEffect(() => {
    adjustTextareaHeight();
  }, [value, isEditing]);

  const handleEdit = () => {
    setIsEditing(true);
  };

  const handleCancel = () => {
    setValue(originalValue);
    setIsEditing(false);
    onCancel?.();
  };

  const handleSave = async () => {
    if (value === originalValue) {
      setIsEditing(false);
      return;
    }

    setIsSaving(true);
    try {
      const updateData: any = {};

      if (language === 'en') {
        updateData.value_en = value;
      } else if (language === 'ru') {
        updateData.value_ru = value;
      } else {
        updateData.value_sr = value;
      }

      const response = await translationAdminApi.updateTranslation(
        translation.id,
        {
          translated_text: value,
          is_verified: true,
        }
      );

      if (response.success) {
        setOriginalValue(value);
        setIsEditing(false);

        // Update the translation object
        const updatedTranslation = {
          ...translation,
          ...updateData,
        };

        onSave?.(updatedTranslation);
      }
    } catch (error) {
      console.error('Failed to save translation:', error);
    } finally {
      setIsSaving(false);
    }
  };

  const handleAITranslate = async () => {
    setIsGeneratingAI(true);
    try {
      // Get source text for AI translation
      let sourceText = '';
      let _sourceLanguage = '';

      if (language !== 'en' && translation.value_en) {
        sourceText = translation.value_en;
        _sourceLanguage = 'en';
      } else if (language !== 'ru' && translation.value_ru) {
        sourceText = translation.value_ru;
        _sourceLanguage = 'ru';
      } else if (language !== 'sr' && translation.value_sr) {
        sourceText = translation.value_sr;
        _sourceLanguage = 'sr';
      }

      if (!sourceText) {
        return;
      }

      const response = await translationAdminApi.translateWithAI(
        sourceText,
        language
      );

      if (response.success && response.data) {
        setValue(response.data.translated_text);
      }
    } catch (error) {
      console.error('AI translation failed:', error);
    } finally {
      setIsGeneratingAI(false);
    }
  };

  const handleKeyDown = (e: React.KeyboardEvent) => {
    if (e.key === 'Enter' && e.ctrlKey) {
      e.preventDefault();
      handleSave();
    } else if (e.key === 'Escape') {
      e.preventDefault();
      handleCancel();
    }
  };

  const getLanguageLabel = () => {
    switch (language) {
      case 'en':
        return 'English';
      case 'ru':
        return 'Русский';
      case 'sr':
        return 'Српски';
      default:
        return language as string;
    }
  };

  if (!isEditing) {
    return (
      <div className={`group relative ${className}`} onClick={handleEdit}>
        <div className="flex items-start gap-2 p-3 rounded-lg border border-base-300 hover:border-primary/50 cursor-pointer transition-colors">
          <LanguageIcon className="h-5 w-5 text-base-content/50 mt-0.5" />

          <div className="flex-1 min-w-0">
            <div className="text-xs font-semibold text-base-content/70 mb-1">
              {getLanguageLabel()}
            </div>
            <div className="text-sm text-base-content break-words">
              {value || (
                <span className="text-base-content/40 italic">
                  {t('noTranslation')}
                </span>
              )}
            </div>
          </div>

          <button
            className="opacity-0 group-hover:opacity-100 transition-opacity"
            onClick={(e) => {
              e.stopPropagation();
              handleEdit();
            }}
          >
            <PencilIcon className="h-4 w-4 text-primary" />
          </button>
        </div>
      </div>
    );
  }

  return (
    <div className={`${className}`}>
      <div className="border border-primary rounded-lg p-3 bg-base-100">
        {/* Header */}
        <div className="flex items-center justify-between mb-2">
          <div className="flex items-center gap-2">
            <LanguageIcon className="h-5 w-5 text-primary" />
            <span className="text-sm font-semibold">{getLanguageLabel()}</span>
          </div>

          {showAIButton && (
            <button
              onClick={handleAITranslate}
              disabled={isGeneratingAI || isSaving}
              className="btn btn-ghost btn-xs gap-1"
            >
              {isGeneratingAI ? (
                <span className="loading loading-spinner loading-xs"></span>
              ) : (
                <SparklesIcon className="h-4 w-4" />
              )}
              {t('aiTranslate')}
            </button>
          )}
        </div>

        {/* Editor */}
        <textarea
          ref={textareaRef}
          value={value}
          onChange={(e) => setValue(e.target.value)}
          onKeyDown={handleKeyDown}
          className="textarea textarea-bordered w-full min-h-[80px] text-sm"
          placeholder={t('enterTranslation')}
          disabled={isSaving || isGeneratingAI}
        />

        {/* Help text */}
        <div className="text-xs text-base-content/60 mt-1">
          {t('shortcuts')}
        </div>

        {/* Actions */}
        <div className="flex items-center justify-end gap-2 mt-3">
          <button
            onClick={handleCancel}
            disabled={isSaving}
            className="btn btn-ghost btn-sm"
          >
            <XMarkIcon className="h-4 w-4" />
            {t('cancel')}
          </button>

          <button
            onClick={handleSave}
            disabled={isSaving || value === originalValue}
            className="btn btn-primary btn-sm"
          >
            {isSaving ? (
              <span className="loading loading-spinner loading-xs"></span>
            ) : (
              <CheckIcon className="h-4 w-4" />
            )}
            {t('save')}
          </button>
        </div>
      </div>
    </div>
  );
}
