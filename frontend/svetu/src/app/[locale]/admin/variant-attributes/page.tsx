import { Metadata } from 'next';
import { getTranslations } from 'next-intl/server';
import VariantAttributesClient from './VariantAttributesClient';

export const metadata: Metadata = {
  title: 'Вариативные атрибуты - Админ панель',
  description: 'Управление вариативными атрибутами товаров',
};

export default async function VariantAttributesPage() {
  const _t = await getTranslations('admin');

  return (
    <div className="container mx-auto px-4 py-8">
      <div className="mb-8">
        <h1 className="text-3xl font-bold mb-2">🔄 Вариативные атрибуты</h1>
        <p className="text-base-content/70">
          Управление атрибутами для создания вариантов товаров
        </p>
      </div>

      <VariantAttributesClient />
    </div>
  );
}
