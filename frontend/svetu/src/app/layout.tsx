import type { Metadata } from 'next';
import { PublicEnvScript } from 'next-runtime-env';
import './globals.css';

export const metadata: Metadata = {
  title: 'SveTu - Marketplace',
  description: 'Your local marketplace',
};

export default function RootLayout({
  children,
}: {
  children: React.ReactNode;
}) {
  return (
    <html suppressHydrationWarning>
      <head>
        <PublicEnvScript />
      </head>
      <body suppressHydrationWarning>{children}</body>
    </html>
  );
}
