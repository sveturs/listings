// IMPORTANT: This root layout is intentionally minimal.
// The actual <html> and <body> tags are in app/[locale]/layout.tsx
// because all routes go through the locale segment.
// Having <html>/<body> in both layouts causes hydration errors.

export default function RootLayout({
  children,
}: {
  children: React.ReactNode;
}) {
  return children;
}
