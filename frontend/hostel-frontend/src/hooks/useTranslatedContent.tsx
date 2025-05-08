// frontend/hostel-frontend/src/hooks/useTranslatedContent.tsx
import { useMemo } from 'react';
import { useLanguage } from '../contexts/LanguageContext';

// Определяем интерфейсы для типизации
interface TranslationMap {
  [language: string]: {
    [field: string]: string;
  };
}

interface TranslatableContent {
  original_language?: string;
  translations?: TranslationMap;
  [key: string]: any;
}

/**
 * Хук для автоматического перевода контента в зависимости от выбранного языка
 * @param content Контент с переводами
 * @returns Переведенный контент в соответствии с текущим языком
 */
export const useTranslatedContent = <T extends TranslatableContent>(content: T | null | undefined): T | null | undefined => {
  const { language } = useLanguage();

  return useMemo(() => {
    // Если контент отсутствует или нет переводов/исходного языка - возвращаем контент без изменений
    if (!content || !content.translations || !content.original_language) {
      return content;
    }

    // Если текущий язык совпадает с оригинальным - возвращаем контент без изменений
    if (language === content.original_language) {
      return content;
    }

    // Создаем копию контента
    const translatedContent = { ...content };

    // Если есть перевод для текущего языка, применяем его
    if (content.translations && content.translations[language]) {
      Object.entries(content.translations[language]).forEach(([field, text]) => {
        // Используем as any, чтобы избежать проблем с индексацией типа
        (translatedContent as any)[field] = text;
      });
    }

    return translatedContent;
  }, [content, language]);
};