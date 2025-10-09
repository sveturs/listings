'use client';

import { useMemo } from 'react';
import DOMPurify from 'isomorphic-dompurify';

interface SafeHTMLProps {
  html: string;
  className?: string;
  /**
   * Дополнительные разрешенные теги (по умолчанию: p, br, strong, em, u, a, ul, ol, li, h1-h6)
   */
  allowedTags?: string[];
  /**
   * Дополнительные разрешенные атрибуты
   */
  allowedAttributes?: Record<string, string[]>;
}

/**
 * Безопасный компонент для рендеринга HTML контента
 * Использует DOMPurify для санитизации HTML и предотвращения XSS атак
 */
export default function SafeHTML({
  html,
  className = '',
  allowedTags,
  allowedAttributes,
}: SafeHTMLProps) {
  const sanitizedHTML = useMemo(() => {
    if (!html) return '';

    // Базовая конфигурация DOMPurify
    const config: any = {
      ALLOWED_TAGS: allowedTags || [
        'p',
        'br',
        'strong',
        'b',
        'em',
        'i',
        'u',
        'a',
        'ul',
        'ol',
        'li',
        'h1',
        'h2',
        'h3',
        'h4',
        'h5',
        'h6',
        'span',
        'div',
        'img',
      ],
      ALLOWED_ATTR: allowedAttributes || {
        a: ['href', 'title', 'target', 'rel'],
        img: ['src', 'alt', 'title', 'width', 'height'],
      },
      // Автоматически добавляет rel="noopener noreferrer" к внешним ссылкам
      ADD_ATTR: ['target'],
      ALLOW_DATA_ATTR: false,
      // Разрешаем только безопасные URI схемы
      ALLOWED_URI_REGEXP:
        /^(?:(?:(?:f|ht)tps?|mailto|tel|callto|sms|cid|xmpp):|[^a-z]|[a-z+.\-]+(?:[^a-z+.\-:]|$))/i,
    };

    return DOMPurify.sanitize(html, config);
  }, [html, allowedTags, allowedAttributes]);

  if (!sanitizedHTML) {
    return null;
  }

  return (
    <div
      className={className}
      dangerouslySetInnerHTML={{ __html: sanitizedHTML }}
    />
  );
}
