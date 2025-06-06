'use client';

import { useState } from 'react';
import { Category } from '@/services/admin';
import { useTranslations } from 'next-intl';

interface CategoryTreeProps {
  categories: Category[];
  onEdit: (category: Category) => void;
  onDelete: (category: Category) => void;
  onManageAttributes?: (category: Category) => void;
  onReorder: (orderedIds: number[]) => void;
  onMove: (categoryId: number, newParentId: number) => void;
}

interface CategoryNodeProps {
  category: Category;
  level: number;
  onEdit: (category: Category) => void;
  onDelete: (category: Category) => void;
  onManageAttributes?: (category: Category) => void;
  categories: Category[];
}

const CategoryNode: React.FC<CategoryNodeProps> = ({
  category,
  level,
  onEdit,
  onDelete,
  onManageAttributes,
  categories,
}) => {
  const [expanded, setExpanded] = useState(true);
  const t = useTranslations('admin');

  const childCategories = categories.filter((c) => c.parent_id === category.id);
  const hasChildren = childCategories.length > 0;

  return (
    <div className="category-node">
      <div
        className={`flex items-center gap-2 p-2 hover:bg-base-200 rounded-lg cursor-pointer`}
        style={{ paddingLeft: `${level * 20}px` }}
      >
        {hasChildren && (
          <button
            onClick={() => setExpanded(!expanded)}
            className="btn btn-ghost btn-xs"
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

        {!hasChildren && <div className="w-8" />}

        {category.icon && <span className="text-lg">{category.icon}</span>}

        <span className="flex-1 font-medium">{category.name}</span>

        {category.items_count !== undefined && (
          <span className="badge badge-sm">{category.items_count}</span>
        )}

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
        <div className="ml-2">
          {childCategories.map((child) => (
            <CategoryNode
              key={child.id}
              category={child}
              level={level + 1}
              onEdit={onEdit}
              onDelete={onDelete}
              onManageAttributes={onManageAttributes}
              categories={categories}
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
  onReorder: _onReorder,
  onMove: _onMove,
}: CategoryTreeProps) {
  const t = useTranslations('admin');

  // Build tree structure
  const rootCategories = categories.filter((c) => !c.parent_id);

  if (categories.length === 0) {
    return (
      <div className="text-center py-8">
        <p className="text-base-content/60">{t('common.noData')}</p>
      </div>
    );
  }

  return (
    <div className="category-tree">
      {rootCategories.map((category) => (
        <CategoryNode
          key={category.id}
          category={category}
          level={0}
          onEdit={onEdit}
          onDelete={onDelete}
          onManageAttributes={onManageAttributes}
          categories={categories}
        />
      ))}
    </div>
  );
}
