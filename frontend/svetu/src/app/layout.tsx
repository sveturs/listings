import './globals.css';

// Корневой layout просто передает children
// Все HTML теги и компоненты определены в /app/[locale]/layout.tsx
export default function RootLayout({
  children,
}: {
  children: React.ReactNode;
}) {
  return children;
}
