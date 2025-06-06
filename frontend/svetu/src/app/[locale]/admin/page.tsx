'use client';

import { useTranslations } from 'next-intl';
import { Link } from '@/i18n/routing';

export default function AdminPage() {
  const t = useTranslations('admin');

  return (
    <div className="container mx-auto px-4 py-8">
      <h1 className="text-3xl font-bold mb-8">{t('title')}</h1>

      <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-6">
        <div className="card bg-base-100 shadow-xl">
          <div className="card-body">
            <h2 className="card-title">{t('sections.users')}</h2>
            <p>{t('sections.usersDescription')}</p>
            <div className="card-actions justify-end">
              <button className="btn btn-primary" disabled title="Coming soon">
                {t('manage')}
              </button>
            </div>
          </div>
        </div>

        <div className="card bg-base-100 shadow-xl">
          <div className="card-body">
            <h2 className="card-title">{t('sections.listings')}</h2>
            <p>{t('sections.listingsDescription')}</p>
            <div className="card-actions justify-end">
              <button className="btn btn-primary" disabled title="Coming soon">
                {t('manage')}
              </button>
            </div>
          </div>
        </div>

        <div className="card bg-base-100 shadow-xl">
          <div className="card-body">
            <h2 className="card-title">{t('sections.categories')}</h2>
            <p>{t('sections.categoriesDescription')}</p>
            <div className="card-actions justify-end">
              <Link href="/admin/categories" className="btn btn-primary">
                {t('manage')}
              </Link>
            </div>
          </div>
        </div>

        <div className="card bg-base-100 shadow-xl">
          <div className="card-body">
            <h2 className="card-title">{t('sections.attributes')}</h2>
            <p>{t('sections.attributesDescription')}</p>
            <div className="card-actions justify-end">
              <Link href="/admin/attributes" className="btn btn-primary">
                {t('manage')}
              </Link>
            </div>
          </div>
        </div>

        <div className="card bg-base-100 shadow-xl">
          <div className="card-body">
            <h2 className="card-title">{t('sections.attributeGroups')}</h2>
            <p>{t('sections.attributeGroupsDescription')}</p>
            <div className="card-actions justify-end">
              <Link href="/admin/attribute-groups" className="btn btn-primary">
                {t('manage')}
              </Link>
            </div>
          </div>
        </div>
      </div>
    </div>
  );
}
