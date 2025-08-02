'use client';

import { useTranslations } from 'next-intl';
import { useRef } from 'react';
import { VariantAttribute } from '@/services/admin';

interface VariantAttributeListProps {
  attributes: VariantAttribute[];
  searchTerm: string;
  filterType: string;
  currentPage: number;
  totalPages: number;
  totalItems: number;
  pageSize: number;
  onSearchChange: (term: string) => void;
  onFilterChange: (type: string) => void;
  onPageChange: (page: number) => void;
  onEdit: (attribute: VariantAttribute) => void;
  onDelete: (attribute: VariantAttribute) => void;
}

export default function VariantAttributeList({
  attributes,
  searchTerm,
  filterType,
  currentPage,
  totalPages,
  totalItems,
  pageSize,
  onSearchChange,
  onFilterChange,
  onPageChange,
  onEdit,
  onDelete,
}: VariantAttributeListProps) {
  const t = useTranslations('admin');
  const searchInputRef = useRef<HTMLInputElement>(null);
  const filterSelectRef = useRef<HTMLSelectElement>(null);

  return (
    <div className="card bg-base-100 shadow-xl">
      <div className="card-body">
        {/* Filters */}
        <div className="flex gap-4 mb-4">
          <div className="form-control flex-1">
            <input
              ref={searchInputRef}
              type="text"
              placeholder={t('common.search')}
              className="input input-bordered"
              value={searchTerm}
              onChange={(e) => onSearchChange(e.target.value)}
            />
          </div>
          <div className="form-control">
            <select
              ref={filterSelectRef}
              className="select select-bordered"
              value={filterType}
              onChange={(e) => onFilterChange(e.target.value)}
            >
              <option value="">{t('variantAttributes.allTypes')}</option>
              <option value="text">{t('variantAttributes.types.text')}</option>
              <option value="number">
                {t('variantAttributes.types.number')}
              </option>
              <option value="select">
                {t('variantAttributes.types.select')}
              </option>
              <option value="multiselect">
                {t('variantAttributes.types.multiselect')}
              </option>
              <option value="boolean">
                {t('variantAttributes.types.boolean')}
              </option>
              <option value="date">{t('variantAttributes.types.date')}</option>
              <option value="range">
                {t('variantAttributes.types.range')}
              </option>
            </select>
          </div>
        </div>

        {/* Attributes Table */}
        <div className="overflow-x-auto">
          <table className="table table-zebra">
            <thead>
              <tr>
                <th>{t('variantAttributes.systemName')}</th>
                <th>{t('variantAttributes.displayName')}</th>
                <th>{t('variantAttributes.type')}</th>
                <th>{t('variantAttributes.settings')}</th>
                <th className="text-center">{t('common.actions')}</th>
              </tr>
            </thead>
            <tbody>
              {attributes.length === 0 ? (
                <tr>
                  <td colSpan={5} className="text-center">
                    {t('common.noData')}
                  </td>
                </tr>
              ) : (
                attributes.map((attr) => (
                  <tr key={attr.id}>
                    <td>
                      <code className="text-sm">{attr.name}</code>
                    </td>
                    <td>{attr.display_name}</td>
                    <td>
                      <span className="badge badge-outline">
                        {t(`variantAttributes.types.${attr.type}`)}
                      </span>
                    </td>
                    <td>
                      <div className="flex gap-1">
                        {attr.is_required && (
                          <span
                            className="badge badge-sm badge-error"
                            title={t('variantAttributes.isRequired')}
                          >
                            *
                          </span>
                        )}
                        {attr.affects_stock && (
                          <span
                            className="badge badge-sm badge-warning"
                            title={t('variantAttributes.affectsStock')}
                          >
                            📦
                          </span>
                        )}
                      </div>
                    </td>
                    <td className="text-center">
                      <div className="dropdown dropdown-end">
                        <label tabIndex={0} className="btn btn-ghost btn-xs">
                          <svg
                            xmlns="http://www.w3.org/2000/svg"
                            className="h-4 w-4"
                            fill="none"
                            viewBox="0 0 24 24"
                            stroke="currentColor"
                          >
                            <path
                              strokeLinecap="round"
                              strokeLinejoin="round"
                              strokeWidth={2}
                              d="M5 12h.01M12 12h.01M19 12h.01M6 12a1 1 0 11-2 0 1 1 0 012 0zm7 0a1 1 0 11-2 0 1 1 0 012 0zm7 0a1 1 0 11-2 0 1 1 0 012 0z"
                            />
                          </svg>
                        </label>
                        <ul
                          tabIndex={0}
                          className="dropdown-content z-[1] menu p-2 shadow bg-base-100 rounded-box w-52"
                        >
                          <li>
                            <a onClick={() => onEdit(attr)}>
                              <svg
                                xmlns="http://www.w3.org/2000/svg"
                                className="h-4 w-4"
                                fill="none"
                                viewBox="0 0 24 24"
                                stroke="currentColor"
                              >
                                <path
                                  strokeLinecap="round"
                                  strokeLinejoin="round"
                                  strokeWidth={2}
                                  d="M11 5H6a2 2 0 00-2 2v11a2 2 0 002 2h11a2 2 0 002-2v-5m-1.414-9.414a2 2 0 112.828 2.828L11.828 15H9v-2.828l8.586-8.586z"
                                />
                              </svg>
                              {t('common.edit')}
                            </a>
                          </li>
                          <li>
                            <a
                              onClick={() => onDelete(attr)}
                              className="text-error"
                            >
                              <svg
                                xmlns="http://www.w3.org/2000/svg"
                                className="h-4 w-4"
                                fill="none"
                                viewBox="0 0 24 24"
                                stroke="currentColor"
                              >
                                <path
                                  strokeLinecap="round"
                                  strokeLinejoin="round"
                                  strokeWidth={2}
                                  d="M19 7l-.867 12.142A2 2 0 0116.138 21H7.862a2 2 0 01-1.995-1.858L5 7m5 4v6m4-6v6m1-10V4a1 1 0 00-1-1h-4a1 1 0 00-1 1v3M4 7h16"
                                />
                              </svg>
                              {t('common.delete')}
                            </a>
                          </li>
                        </ul>
                      </div>
                    </td>
                  </tr>
                ))
              )}
            </tbody>
          </table>
        </div>

        {/* Pagination */}
        {totalPages > 1 && (
          <div className="flex justify-between items-center mt-4">
            <div className="text-sm text-base-content/70">
              {t('common.showing')} {(currentPage - 1) * pageSize + 1} -{' '}
              {Math.min(currentPage * pageSize, totalItems)} {t('common.of')}{' '}
              {totalItems} {t('common.items')}
            </div>
            <div className="join">
              <button
                className="join-item btn btn-sm"
                disabled={currentPage === 1}
                onClick={() => onPageChange(currentPage - 1)}
              >
                «
              </button>

              {/* Показываем страницы */}
              {Array.from({ length: Math.min(5, totalPages) }, (_, i) => {
                let pageNumber;
                if (totalPages <= 5) {
                  pageNumber = i + 1;
                } else if (currentPage <= 3) {
                  pageNumber = i + 1;
                } else if (currentPage >= totalPages - 2) {
                  pageNumber = totalPages - 4 + i;
                } else {
                  pageNumber = currentPage - 2 + i;
                }

                return (
                  <button
                    key={pageNumber}
                    className={`join-item btn btn-sm ${
                      pageNumber === currentPage ? 'btn-active' : ''
                    }`}
                    onClick={() => onPageChange(pageNumber)}
                  >
                    {pageNumber}
                  </button>
                );
              })}

              <button
                className="join-item btn btn-sm"
                disabled={currentPage === totalPages}
                onClick={() => onPageChange(currentPage + 1)}
              >
                »
              </button>
            </div>
          </div>
        )}
      </div>
    </div>
  );
}
