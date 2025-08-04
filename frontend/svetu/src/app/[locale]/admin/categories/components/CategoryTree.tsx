'use client';

import { useState } from 'react';
import { Category } from '@/services/admin';
import { useTranslations } from 'next-intl';
import { TranslationStatus } from '@/components/attributes/TranslationStatus';
import { InlineTranslationEditor } from '@/components/attributes/InlineTranslationEditor';
import { adminApi } from '@/services/admin';

interface CategoryTreeProps {
  categories: Category[];
  onEdit: (category: Category) => void;
  onDelete: (category: Category) => void;
  onManageAttributes?: (category: Category) => void;
  onManageKeywords?: (category: Category) => void;
  onManageVariantAttributes?: (category: Category) => void;
  onReorder: (orderedIds: number[]) => void;
  onMove: (categoryId: number, newParentId: number) => void;
}

interface CategoryNodeProps {
  category: Category;
  level: number;
  onEdit: (category: Category) => void;
  onDelete: (category: Category) => void;
  onManageAttributes?: (category: Category) => void;
  onManageKeywords?: (category: Category) => void;
  onManageVariantAttributes?: (category: Category) => void;
  categories: Category[];
  isLast: boolean;
  parentLines: boolean[];
}

const CategoryNode: React.FC<CategoryNodeProps> = ({
  category,
  level,
  onEdit,
  onDelete,
  onManageAttributes,
  onManageKeywords,
  onManageVariantAttributes,
  categories,
  isLast,
  parentLines,
}) => {
  const [expanded, setExpanded] = useState(true);
  const [translations, setTranslations] = useState<Record<string, string>>({
    en: category.translations?.en || category.name,
    ru: category.translations?.ru || category.name,
    sr: category.translations?.sr || category.name,
  });
  const t = useTranslations('admin');

  const childCategories = categories.filter((c) => c.parent_id === category.id);
  const hasChildren = childCategories.length > 0;

  return (
    <div className="category-node">
      <div
        className={`flex items-center gap-2 p-2 hover:bg-base-200 rounded-lg transition-colors ${
          !category.is_active ? 'opacity-50' : ''
        }`}
      >
        {/* Hierarchy lines */}
        <div className="flex items-center">
          {level > 0 && (
            <>
              {parentLines.map((showLine, index) => (
                <div
                  key={index}
                  className={`w-5 h-full ${
                    showLine && index < parentLines.length - 1
                      ? 'border-l-2 border-base-300'
                      : ''
                  }`}
                />
              ))}
              <div className="relative w-5 h-full">
                <div
                  className={`absolute top-0 left-0 w-full h-1/2 ${
                    !isLast ? 'border-l-2 border-base-300' : ''
                  }`}
                />
                <div className="absolute top-1/2 left-0 w-full h-px border-t-2 border-base-300" />
                {!isLast && (
                  <div className="absolute top-1/2 left-0 w-full h-1/2 border-l-2 border-base-300" />
                )}
              </div>
            </>
          )}
        </div>

        {/* Expand/Collapse button */}
        {hasChildren && (
          <button
            onClick={() => setExpanded(!expanded)}
            className="btn btn-ghost btn-xs p-0 min-h-0 h-6 w-6"
            aria-label={expanded ? t('common.collapse') : t('common.expand')}
          >
            <svg
              xmlns="http://www.w3.org/2000/svg"
              className={`h-4 w-4 transition-transform ${expanded ? 'rotate-90' : ''}`}
              fill="none"
              viewBox="0 0 24 24"
              stroke="currentColor"
            >
              <path
                strokeLinecap="round"
                strokeLinejoin="round"
                strokeWidth={2}
                d="M9 5l7 7-7 7"
              />
            </svg>
          </button>
        )}

        {!hasChildren && <div className="w-6" />}

        {/* Category icon */}
        {category.icon && <span className="text-lg">{category.icon}</span>}

        {/* Category name with active indicator */}
        <div className="flex-1 font-medium flex items-center gap-2">
          <InlineTranslationEditor
            entityType="category"
            entityId={category.id}
            fieldName="name"
            translations={translations}
            onSave={async (newTranslations) => {
              await adminApi.updateFieldTranslation(
                'category',
                category.id,
                'name',
                newTranslations
              );
              setTranslations(newTranslations);
              // TODO: обновить название в UI без перезагрузки
            }}
            compact={true}
          />
          {!category.is_active && (
            <span className="badge badge-ghost badge-sm">
              {t('common.inactive')}
            </span>
          )}
        </div>

        {/* Translation status */}
        <TranslationStatus
          entityType="category"
          entityId={category.id}
          compact={true}
        />

        {/* Items count */}
        {category.items_count !== undefined && (
          <span className="badge badge-sm badge-primary">
            {category.items_count}
          </span>
        )}

        {/* Active status indicator */}
        <div
          className={`w-2 h-2 rounded-full ${category.is_active ? 'bg-success' : 'bg-base-300'}`}
          title={category.is_active ? t('common.active') : t('common.inactive')}
        />

        {/* Actions dropdown */}
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
              <a onClick={() => onEdit(category)}>
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
            {onManageAttributes && (
              <li>
                <a onClick={() => onManageAttributes(category)}>
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
                      d="M7 7h.01M7 3h5c.512 0 1.024.195 1.414.586l7 7a2 2 0 010 2.828l-7 7a2 2 0 01-2.828 0l-7-7A1.994 1.994 0 013 12V7a4 4 0 014-4z"
                    />
                  </svg>
                  {t('sections.attributes')}
                </a>
              </li>
            )}
            {onManageKeywords && (
              <li>
                <a onClick={() => onManageKeywords(category)}>
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
                      d="M7 20l4-16m2 16l4-16M6 9h14M4 15h14"
                    />
                  </svg>
                  {t('categories.keywords.title')}
                </a>
              </li>
            )}
            {onManageVariantAttributes && (
              <li>
                <a onClick={() => onManageVariantAttributes(category)}>
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
                      d="M4 7v10c0 2.21 3.582 4 8 4s8-1.79 8-4V7M4 7c0 2.21 3.582 4 8 4s8-1.79 8-4M4 7c0-2.21 3.582-4 8-4s8 1.79 8 4m0 5c0 2.21-3.582 4-8 4s-8-1.79-8-4"
                    />
                  </svg>
                  {t('categories.variantAttributes')}
                </a>
              </li>
            )}
            <li className="divider my-0"></li>
            <li>
              <a
                onClick={() =>
                  onEdit({ ...category, is_active: !category.is_active })
                }
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
                    d={
                      category.is_active
                        ? 'M18.364 18.364A9 9 0 005.636 5.636m12.728 12.728A9 9 0 015.636 5.636m12.728 12.728L5.636 5.636'
                        : 'M9 12l2 2 4-4m6 2a9 9 0 11-18 0 9 9 0 0118 0z'
                    }
                  />
                </svg>
                {category.is_active
                  ? t('common.deactivate')
                  : t('common.activate')}
              </a>
            </li>
            <li>
              <a onClick={() => onDelete(category)} className="text-error">
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
      </div>

      {hasChildren && expanded && (
        <div className="category-children">
          {childCategories.map((child, index) => (
            <CategoryNode
              key={child.id}
              category={child}
              level={level + 1}
              onEdit={onEdit}
              onDelete={onDelete}
              onManageAttributes={onManageAttributes}
              onManageKeywords={onManageKeywords}
              onManageVariantAttributes={onManageVariantAttributes}
              categories={categories}
              isLast={index === childCategories.length - 1}
              parentLines={[...parentLines, !isLast]}
            />
          ))}
        </div>
      )}
    </div>
  );
};

export default function CategoryTree({
  categories,
  onEdit,
  onDelete,
  onManageAttributes,
  onManageKeywords,
  onManageVariantAttributes,
  onReorder: _onReorder,
  onMove: _onMove,
}: CategoryTreeProps) {
  const t = useTranslations('admin');
  const [showInactive, setShowInactive] = useState(true);

  // Filter categories based on active status
  const filteredCategories = showInactive
    ? categories
    : categories.filter((c) => c.is_active);

  // Build tree structure
  const rootCategories = filteredCategories.filter((c) => !c.parent_id);

  if (categories.length === 0) {
    return (
      <div className="text-center py-8">
        <p className="text-base-content/60">{t('common.noData')}</p>
      </div>
    );
  }

  return (
    <div className="category-tree space-y-1">
      {/* Filter controls */}
      <div className="flex justify-between items-center mb-4 p-2 bg-base-200 rounded-lg">
        <div className="flex items-center gap-2">
          <span className="text-sm font-medium">{t('common.filters')}:</span>
          <label className="label cursor-pointer gap-2">
            <input
              type="checkbox"
              className="checkbox checkbox-sm"
              checked={showInactive}
              onChange={(e) => setShowInactive(e.target.checked)}
            />
            <span className="label-text">{t('common.showInactive')}</span>
          </label>
        </div>
        <div className="flex items-center gap-2 text-sm">
          <div className="flex items-center gap-1">
            <div className="w-2 h-2 rounded-full bg-success" />
            <span>{t('common.active')}</span>
          </div>
          <div className="flex items-center gap-1">
            <div className="w-2 h-2 rounded-full bg-base-300" />
            <span>{t('common.inactive')}</span>
          </div>
        </div>
      </div>

      {/* Category tree */}
      {rootCategories.map((category, index) => (
        <CategoryNode
          key={category.id}
          category={category}
          level={0}
          onEdit={onEdit}
          onDelete={onDelete}
          onManageAttributes={onManageAttributes}
          onManageKeywords={onManageKeywords}
          onManageVariantAttributes={onManageVariantAttributes}
          categories={filteredCategories}
          isLast={index === rootCategories.length - 1}
          parentLines={[]}
        />
      ))}
    </div>
  );
}
