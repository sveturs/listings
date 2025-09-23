import { Metadata } from "next";
import { getTranslations } from "next-intl/server";
import CarsPageClient from "./CarsPageClient";

export async function generateMetadata({
  params: { locale },
}: {
  params: { locale: string };
}): Promise<Metadata> {
  const t = await getTranslations({ locale, namespace: "cars" });

  return {
    title: t("pageTitle"),
    description: t("pageDescription"),
  };
}

export default async function CarsPage({
  params: { locale },
}: {
  params: { locale: string };
}) {
  const t = await getTranslations({ locale, namespace: "cars" });

  return <CarsPageClient locale={locale} />;
}