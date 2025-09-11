'use client';

import { useState } from 'react';
import { useTranslations } from 'next-intl';
import InlineTranslationEditor from './InlineTranslationEditor';
import TranslationSearch from './TranslationSearch';
import type { Translation } from '@/types/translation';
import type { Translation as TranslationApiType } from '@/services/translationAdminApi';

export default function TranslationEditorDemo() {
  const t = useTranslations('admin');
  const [selectedTranslation, setSelectedTranslation] =
    useState<Translation | null>(null);

  // Demo translation for testing
  const demoTranslation: Translation = {
    id: 1,
    entity_type: 'category',
    entity_id: 1602,
    field_name: 'name',
    language: 'all',
    key: 'category.1602.name',
    value_en: 'Seeds and fertilizers',
    value_ru: 'Семена и удобрения',
    value_sr: 'Семе и ђубрива',
    context: 'Agricultural category name',
    is_verified: true,
    created_at: new Date().toISOString(),
    updated_at: new Date().toISOString(),
  };

  const handleSearchSelect = (translation: TranslationApiType) => {
    // Convert API translation type to local type
    const convertedTranslation: Translation = {
      ...translation,
      key: `${translation.entity_type}.${translation.entity_id}.${translation.field_name}`,
      value_en:
        translation.language === 'en' ? translation.translated_text : undefined,
      value_ru:
        translation.language === 'ru' ? translation.translated_text : undefined,
      value_sr:
        translation.language === 'sr' ? translation.translated_text : undefined,
      created_at: translation.created_at || new Date().toISOString(),
      updated_at: translation.updated_at || new Date().toISOString(),
    };
    setSelectedTranslation(convertedTranslation);
  };

  const handleSave = (translation: Translation) => {
    console.log('Saved translation:', translation);
    // In real app, this would update the translation in the backend
  };

  return (
    <div className="container mx-auto px-4 py-8">
      <div className="mb-8">
        <h2 className="text-2xl font-bold mb-2">{t('editorDemo.title')}</h2>
        <p className="text-base-content/60">{t('editorDemo.description')}</p>
      </div>

      {/* Search Component */}
      <div className="mb-8">
        <h3 className="text-lg font-semibold mb-4">
          {t('editorDemo.searchTitle')}
        </h3>
        <div className="max-w-2xl">
          <TranslationSearch onResultSelect={handleSearchSelect} />
        </div>
      </div>

      {/* Selected Translation Editor */}
      {selectedTranslation && (
        <div className="mb-8">
          <h3 className="text-lg font-semibold mb-4">
            {t('editorDemo.selectedTranslation')}
          </h3>
          <div className="card bg-base-100 shadow-sm">
            <div className="card-body">
              <div className="mb-4">
                <span className="font-mono text-sm text-primary">
                  {selectedTranslation.key}
                </span>
                {selectedTranslation.context && (
                  <p className="text-sm text-base-content/60 mt-1">
                    {selectedTranslation.context}
                  </p>
                )}
              </div>

              <div className="grid md:grid-cols-3 gap-4">
                <InlineTranslationEditor
                  translation={selectedTranslation}
                  language="en"
                  onSave={handleSave}
                />
                <InlineTranslationEditor
                  translation={selectedTranslation}
                  language="ru"
                  onSave={handleSave}
                />
                <InlineTranslationEditor
                  translation={selectedTranslation}
                  language="sr"
                  onSave={handleSave}
                />
              </div>
            </div>
          </div>
        </div>
      )}

      {/* Demo Translation Editor */}
      <div>
        <h3 className="text-lg font-semibold mb-4">
          {t('editorDemo.demoTitle')}
        </h3>
        <div className="card bg-base-100 shadow-sm">
          <div className="card-body">
            <div className="mb-4">
              <span className="font-mono text-sm text-primary">
                {demoTranslation.key}
              </span>
              {demoTranslation.context && (
                <p className="text-sm text-base-content/60 mt-1">
                  {demoTranslation.context}
                </p>
              )}
            </div>

            <div className="grid md:grid-cols-3 gap-4">
              <InlineTranslationEditor
                translation={demoTranslation}
                language="en"
                onSave={handleSave}
              />
              <InlineTranslationEditor
                translation={demoTranslation}
                language="ru"
                onSave={handleSave}
              />
              <InlineTranslationEditor
                translation={demoTranslation}
                language="sr"
                onSave={handleSave}
              />
            </div>
          </div>
        </div>
      </div>
    </div>
  );
}
