import type { Metadata } from 'next';
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
    <html>
      <body>{children}</body>
    </html>
  );
}
