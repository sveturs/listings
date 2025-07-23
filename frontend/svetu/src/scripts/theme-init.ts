// Этот скрипт должен выполняться до рендеринга страницы, чтобы избежать мерцания темы
export const themeInitScript = `
  (function() {
    try {
      const savedTheme = localStorage.getItem('theme');
      const systemTheme = window.matchMedia('(prefers-color-scheme: dark)').matches ? 'dark' : 'light';
      const theme = savedTheme || systemTheme;
      document.documentElement.setAttribute('data-theme', theme);
    } catch (e) {}
  })();
`;
