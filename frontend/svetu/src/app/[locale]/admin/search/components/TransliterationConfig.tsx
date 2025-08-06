'use client';

import { useState, useEffect } from 'react';
import { useTranslations } from 'next-intl';

// Временные типы (до создания API)
interface _TransliterationRule {
  id: number;
  source_char: string;
  target_char: string;
  language: 'ru' | 'sr';
  enabled: boolean;
  created_at: string;
  updated_at: string;
}

// Встроенные правила транслитерации (из slug.go)
const BUILTIN_RULES = {
  ru: [
    { source_char: 'а', target_char: 'a' },
    { source_char: 'б', target_char: 'b' },
    { source_char: 'в', target_char: 'v' },
    { source_char: 'г', target_char: 'g' },
    { source_char: 'д', target_char: 'd' },
    { source_char: 'е', target_char: 'e' },
    { source_char: 'ё', target_char: 'yo' },
    { source_char: 'ж', target_char: 'zh' },
    { source_char: 'з', target_char: 'z' },
    { source_char: 'и', target_char: 'i' },
    { source_char: 'й', target_char: 'y' },
    { source_char: 'к', target_char: 'k' },
    { source_char: 'л', target_char: 'l' },
    { source_char: 'м', target_char: 'm' },
    { source_char: 'н', target_char: 'n' },
    { source_char: 'о', target_char: 'o' },
    { source_char: 'п', target_char: 'p' },
    { source_char: 'р', target_char: 'r' },
    { source_char: 'с', target_char: 's' },
    { source_char: 'т', target_char: 't' },
    { source_char: 'у', target_char: 'u' },
    { source_char: 'ф', target_char: 'f' },
    { source_char: 'х', target_char: 'h' },
    { source_char: 'ц', target_char: 'ts' },
    { source_char: 'ч', target_char: 'ch' },
    { source_char: 'ш', target_char: 'sh' },
    { source_char: 'щ', target_char: 'sch' },
    { source_char: 'ъ', target_char: '' },
    { source_char: 'ы', target_char: 'y' },
    { source_char: 'ь', target_char: '' },
    { source_char: 'э', target_char: 'e' },
    { source_char: 'ю', target_char: 'yu' },
    { source_char: 'я', target_char: 'ya' },
  ],
  sr: [
    { source_char: 'ђ', target_char: 'đ' },
    { source_char: 'ј', target_char: 'j' },
    { source_char: 'љ', target_char: 'lj' },
    { source_char: 'њ', target_char: 'nj' },
    { source_char: 'ћ', target_char: 'ć' },
    { source_char: 'џ', target_char: 'dž' },
    { source_char: 'ш', target_char: 'š' },
    { source_char: 'ж', target_char: 'ž' },
    { source_char: 'ч', target_char: 'č' },
    { source_char: 'Ђ', target_char: 'Đ' },
    { source_char: 'Ј', target_char: 'J' },
    { source_char: 'Љ', target_char: 'LJ' },
    { source_char: 'Њ', target_char: 'NJ' },
    { source_char: 'Ћ', target_char: 'Ć' },
    { source_char: 'Џ', target_char: 'DŽ' },
    { source_char: 'Ш', target_char: 'Š' },
    { source_char: 'Ж', target_char: 'Ž' },
    { source_char: 'Ч', target_char: 'Č' },
  ],
};

export default function TransliterationConfig() {
  const t = useTranslations('admin');
  const [selectedLanguage, setSelectedLanguage] = useState<'ru' | 'sr'>('ru');
  const [testText, setTestText] = useState('');
  const [transliteratedText, setTransliteratedText] = useState('');
  const [newRule, setNewRule] = useState({ source_char: '', target_char: '' });
  const [loading, setLoading] = useState(false);

  // Простая функция транслитерации (пока без API)
  const testTransliteration = (text: string, lang: 'ru' | 'sr') => {
    let result = text.toLowerCase();
    const rules = BUILTIN_RULES[lang];

    rules.forEach((rule) => {
      if (rule.target_char) {
        result = result.replaceAll(rule.source_char, rule.target_char);
      } else {
        result = result.replaceAll(rule.source_char, '');
      }
    });

    return result;
  };

  // Обновляем результат при изменении текста или языка
  useEffect(() => {
    if (testText) {
      setTransliteratedText(testTransliteration(testText, selectedLanguage));
    } else {
      setTransliteratedText('');
    }
  }, [testText, selectedLanguage]);

  const handleAddRule = async () => {
    if (!newRule.source_char || !newRule.target_char) return;

    setLoading(true);
    try {
      // TODO: API вызов для добавления правила
      console.log('Adding rule:', { ...newRule, language: selectedLanguage });
      setNewRule({ source_char: '', target_char: '' });
    } catch (error) {
      console.error('Failed to add rule:', error);
    } finally {
      setLoading(false);
    }
  };

  return (
    <div className="space-y-6">
      {/* Заголовок */}
      <div className="flex items-center justify-between">
        <h2 className="text-2xl font-bold">
          {t('title', { defaultValue: 'Настройка транслитерации' })}
        </h2>
        <div className="badge badge-info">
          {t('status', { defaultValue: 'В разработке' })}
        </div>
      </div>

      {/* Выбор языка */}
      <div className="card bg-base-100 shadow-md">
        <div className="card-body">
          <h3 className="card-title">
            {t('languageSelect', { defaultValue: 'Выбор языка' })}
          </h3>
          <div className="tabs tabs-boxed">
            <button
              className={`tab ${selectedLanguage === 'ru' ? 'tab-active' : ''}`}
              onClick={() => setSelectedLanguage('ru')}
            >
              Русский
            </button>
            <button
              className={`tab ${selectedLanguage === 'sr' ? 'tab-active' : ''}`}
              onClick={() => setSelectedLanguage('sr')}
            >
              Српски
            </button>
          </div>
        </div>
      </div>

      {/* Тестирование транслитерации */}
      <div className="card bg-base-100 shadow-md">
        <div className="card-body">
          <h3 className="card-title">
            {t('testing', { defaultValue: 'Тестирование транслитерации' })}
          </h3>
          <div className="space-y-4">
            <div className="form-control">
              <label className="label">
                <span className="label-text">
                  {t('inputText', {
                    defaultValue: 'Введите текст для тестирования',
                  })}
                </span>
              </label>
              <textarea
                className="textarea textarea-bordered h-20"
                placeholder={
                  selectedLanguage === 'ru'
                    ? 'Введите русский текст...'
                    : 'Унесите српски текст...'
                }
                value={testText}
                onChange={(e) => setTestText(e.target.value)}
              />
            </div>
            <div className="form-control">
              <label className="label">
                <span className="label-text">
                  {t('result', { defaultValue: 'Результат транслитерации' })}
                </span>
              </label>
              <textarea
                className="textarea textarea-bordered h-20 bg-base-200"
                value={transliteratedText}
                readOnly
                placeholder={t('resultPlaceholder', {
                  defaultValue: 'Результат появится здесь...',
                })}
              />
            </div>
          </div>
        </div>
      </div>

      {/* Таблица правил */}
      <div className="card bg-base-100 shadow-md">
        <div className="card-body">
          <h3 className="card-title">
            {t('rules', { defaultValue: 'Правила транслитерации' })} (
            {selectedLanguage.toUpperCase()})
          </h3>

          {/* Добавление нового правила */}
          <div className="flex gap-2 mb-4">
            <input
              type="text"
              className="input input-bordered input-sm flex-1"
              placeholder={t('sourceChar', { defaultValue: 'Исходный символ' })}
              value={newRule.source_char}
              onChange={(e) =>
                setNewRule({ ...newRule, source_char: e.target.value })
              }
              maxLength={2}
            />
            <input
              type="text"
              className="input input-bordered input-sm flex-1"
              placeholder={t('targetChar', {
                defaultValue: 'Целевой символ(ы)',
              })}
              value={newRule.target_char}
              onChange={(e) =>
                setNewRule({ ...newRule, target_char: e.target.value })
              }
              maxLength={5}
            />
            <button
              className={`btn btn-primary btn-sm ${loading ? 'loading' : ''}`}
              onClick={handleAddRule}
              disabled={!newRule.source_char || loading}
            >
              {t('add', { defaultValue: 'Добавить' })}
            </button>
          </div>

          {/* Таблица встроенных правил */}
          <div className="overflow-x-auto">
            <table className="table table-compact w-full">
              <thead>
                <tr>
                  <th>
                    {t('sourceChar', { defaultValue: 'Исходный символ' })}
                  </th>
                  <th>
                    {t('targetChar', { defaultValue: 'Целевой символ(ы)' })}
                  </th>
                  <th>{t('type', { defaultValue: 'Тип' })}</th>
                  <th>{t('actions', { defaultValue: 'Действия' })}</th>
                </tr>
              </thead>
              <tbody>
                {BUILTIN_RULES[selectedLanguage].map((rule, index) => (
                  <tr key={index}>
                    <td>
                      <code className="bg-base-200 px-2 py-1 rounded text-lg">
                        {rule.source_char}
                      </code>
                    </td>
                    <td>
                      <code className="bg-base-200 px-2 py-1 rounded">
                        {rule.target_char || '(удаление)'}
                      </code>
                    </td>
                    <td>
                      <div className="badge badge-ghost">
                        {t('builtin', { defaultValue: 'Встроенное' })}
                      </div>
                    </td>
                    <td>
                      <button
                        className="btn btn-ghost btn-xs"
                        disabled
                        title={t('builtinNotEditable', {
                          defaultValue:
                            'Встроенные правила нельзя редактировать',
                        })}
                      >
                        {t('edit', { defaultValue: 'Редактировать' })}
                      </button>
                    </td>
                  </tr>
                ))}
              </tbody>
            </table>
          </div>
        </div>
      </div>

      {/* Информация */}
      <div className="alert alert-info">
        <svg
          className="stroke-current shrink-0 h-6 w-6"
          fill="none"
          viewBox="0 0 24 24"
        >
          <path
            strokeLinecap="round"
            strokeLinejoin="round"
            strokeWidth="2"
            d="M13 16h-1v-4h-1m1-4h.01M21 12a9 9 0 11-18 0 9 9 0 0118 0z"
          ></path>
        </svg>
        <span>{t('infoLabel')}</span>
      </div>
    </div>
  );
}
