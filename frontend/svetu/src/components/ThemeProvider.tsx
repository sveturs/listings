'use client';

import { useEffect } from 'react';

export function ThemeProvider({ children }: { children: React.ReactNode }) {
  useEffect(() => {
    // Применяем тему только на клиенте после монтирования
    const applyTheme = () => {
      try {
        const savedTheme = localStorage.getItem('theme');
        const systemTheme =
          window.matchMedia &&
          window.matchMedia('(prefers-color-scheme: dark)').matches
            ? 'dark'
            : 'light';
        const theme = savedTheme || systemTheme || 'light';

        // Устанавливаем атрибут для DaisyUI
        document.documentElement.setAttribute('data-theme', theme);
      } catch {
        // Fallback to light theme
        document.documentElement.setAttribute('data-theme', 'light');
      }
    };

    applyTheme();

    // Слушаем изменения системной темы
    const mediaQuery = window.matchMedia('(prefers-color-scheme: dark)');
    const handleChange = () => {
      if (!localStorage.getItem('theme')) {
        applyTheme();
      }
    };

    mediaQuery.addEventListener('change', handleChange);
    return () => mediaQuery.removeEventListener('change', handleChange);
  }, []);

  return <>{children}</>;
}
