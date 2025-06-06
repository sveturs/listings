'use client';

import { useTranslations } from 'next-intl';
import { Link } from '@/i18n/routing';
import { usePathname } from 'next/navigation';
import AdminGuard from '@/components/AdminGuard';

export default function AdminLayout({
  children,
}: {
  children: React.ReactNode;
}) {
  const t = useTranslations('admin');
  const pathname = usePathname();

  const isActive = (path: string) => pathname.includes(path);

  return (
    <AdminGuard>
      <div className="min-h-screen bg-base-200">
        <div className="drawer drawer-mobile lg:drawer-open">
          <input id="admin-drawer" type="checkbox" className="drawer-toggle" />

          <div className="drawer-content">
            {/* Navbar for mobile */}
            <div className="navbar lg:hidden bg-base-100 shadow-md">
              <div className="flex-none">
                <label
                  htmlFor="admin-drawer"
                  className="btn btn-square btn-ghost"
                >
                  <svg
                    xmlns="http://www.w3.org/2000/svg"
                    fill="none"
                    viewBox="0 0 24 24"
                    className="inline-block w-6 h-6 stroke-current"
                  >
                    <path
                      strokeLinecap="round"
                      strokeLinejoin="round"
                      strokeWidth="2"
                      d="M4 6h16M4 12h16M4 18h16"
                    ></path>
                  </svg>
                </label>
              </div>
              <div className="flex-1">
                <span className="text-xl font-bold">{t('title')}</span>
              </div>
            </div>

            {/* Main content */}
            <main className="p-4 lg:p-8">{children}</main>
          </div>

          {/* Sidebar */}
          <div className="drawer-side">
            <label htmlFor="admin-drawer" className="drawer-overlay"></label>
            <aside className="w-64 min-h-full bg-base-100">
              <div className="p-4 border-b">
                <h2 className="text-xl font-bold">{t('title')}</h2>
              </div>

              <ul className="menu p-4 w-full">
                <li>
                  <Link
                    href="/admin"
                    className={
                      isActive('/admin') && !isActive('/admin/') ? 'active' : ''
                    }
                  >
                    <svg
                      xmlns="http://www.w3.org/2000/svg"
                      className="h-5 w-5"
                      fill="none"
                      viewBox="0 0 24 24"
                      stroke="currentColor"
                    >
                      <path
                        strokeLinecap="round"
                        strokeLinejoin="round"
                        strokeWidth={2}
                        d="M3 12l2-2m0 0l7-7 7 7M5 10v10a1 1 0 001 1h3m10-11l2 2m-2-2v10a1 1 0 01-1 1h-3m-6 0a1 1 0 001-1v-4a1 1 0 011-1h2a1 1 0 011 1v4a1 1 0 001 1m-6 0h6"
                      />
                    </svg>
                    {t('sections.dashboard')}
                  </Link>
                </li>

                <li className="menu-title mt-4">
                  <span>{t('sections.catalog')}</span>
                </li>

                <li>
                  <Link
                    href="/admin/categories"
                    className={isActive('/admin/categories') ? 'active' : ''}
                  >
                    <svg
                      xmlns="http://www.w3.org/2000/svg"
                      className="h-5 w-5"
                      fill="none"
                      viewBox="0 0 24 24"
                      stroke="currentColor"
                    >
                      <path
                        strokeLinecap="round"
                        strokeLinejoin="round"
                        strokeWidth={2}
                        d="M19 11H5m14 0a2 2 0 012 2v6a2 2 0 01-2 2H5a2 2 0 01-2-2v-6a2 2 0 012-2m14 0V9a2 2 0 00-2-2M5 11V9a2 2 0 012-2m0 0V5a2 2 0 012-2h6a2 2 0 012 2v2M7 7h10"
                      />
                    </svg>
                    {t('sections.categories')}
                  </Link>
                </li>

                <li>
                  <Link
                    href="/admin/attributes"
                    className={isActive('/admin/attributes') ? 'active' : ''}
                  >
                    <svg
                      xmlns="http://www.w3.org/2000/svg"
                      className="h-5 w-5"
                      fill="none"
                      viewBox="0 0 24 24"
                      stroke="currentColor"
                    >
                      <path
                        strokeLinecap="round"
                        strokeLinejoin="round"
                        strokeWidth={2}
                        d="M7 7h.01M7 3h5c.512 0 1.024.195 1.414.586l7 7a2 2 0 010 2.828l-7 7a2 2 0 01-2.828 0l-7-7A1.994 1.994 0 013 12V7a4 4 0 014-4z"
                      />
                    </svg>
                    {t('sections.attributes')}
                  </Link>
                </li>

                <li>
                  <Link
                    href="/admin/attribute-groups"
                    className={
                      isActive('/admin/attribute-groups') ? 'active' : ''
                    }
                  >
                    <svg
                      xmlns="http://www.w3.org/2000/svg"
                      className="h-5 w-5"
                      fill="none"
                      viewBox="0 0 24 24"
                      stroke="currentColor"
                    >
                      <path
                        strokeLinecap="round"
                        strokeLinejoin="round"
                        strokeWidth={2}
                        d="M4 6a2 2 0 012-2h2a2 2 0 012 2v2a2 2 0 01-2 2H6a2 2 0 01-2-2V6zM14 6a2 2 0 012-2h2a2 2 0 012 2v2a2 2 0 01-2 2h-2a2 2 0 01-2-2V6zM4 16a2 2 0 012-2h2a2 2 0 012 2v2a2 2 0 01-2 2H6a2 2 0 01-2-2v-2zM14 16a2 2 0 012-2h2a2 2 0 012 2v2a2 2 0 01-2 2h-2a2 2 0 01-2-2v-2z"
                      />
                    </svg>
                    {t('sections.attributeGroups')}
                  </Link>
                </li>

                <li className="menu-title mt-4">
                  <span>{t('sections.content')}</span>
                </li>

                <li>
                  <Link
                    href="/admin/listings"
                    className={isActive('/admin/listings') ? 'active' : ''}
                  >
                    <svg
                      xmlns="http://www.w3.org/2000/svg"
                      className="h-5 w-5"
                      fill="none"
                      viewBox="0 0 24 24"
                      stroke="currentColor"
                    >
                      <path
                        strokeLinecap="round"
                        strokeLinejoin="round"
                        strokeWidth={2}
                        d="M9 5H7a2 2 0 00-2 2v12a2 2 0 002 2h10a2 2 0 002-2V7a2 2 0 00-2-2h-2M9 5a2 2 0 002 2h2a2 2 0 002-2M9 5a2 2 0 012-2h2a2 2 0 012 2"
                      />
                    </svg>
                    {t('sections.listings')}
                  </Link>
                </li>

                <li>
                  <Link
                    href="/admin/users"
                    className={isActive('/admin/users') ? 'active' : ''}
                  >
                    <svg
                      xmlns="http://www.w3.org/2000/svg"
                      className="h-5 w-5"
                      fill="none"
                      viewBox="0 0 24 24"
                      stroke="currentColor"
                    >
                      <path
                        strokeLinecap="round"
                        strokeLinejoin="round"
                        strokeWidth={2}
                        d="M12 4.354a4 4 0 110 5.292M15 21H3v-1a6 6 0 0112 0v1zm0 0h6v-1a6 6 0 00-9-5.197M13 7a4 4 0 11-8 0 4 4 0 018 0z"
                      />
                    </svg>
                    {t('sections.users')}
                  </Link>
                </li>
              </ul>
            </aside>
          </div>
        </div>
      </div>
    </AdminGuard>
  );
}
