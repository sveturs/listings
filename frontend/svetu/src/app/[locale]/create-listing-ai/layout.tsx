export default async function CreateListingAILayout({
  children,
}: {
  children: React.ReactNode;
  params: Promise<{ locale: string }>;
}) {
  // Cars translations are already loaded in the main layout via ModularIntlProvider
  // No need for additional provider here as it causes conflicts
  return <>{children}</>;
}
