'use client';

import React from 'react';
import { motion, AnimatePresence } from 'framer-motion';
import { FiChevronRight, FiChevronDown } from 'react-icons/fi';
import { getCategoryIcon } from '@/utils/categoryIcons';

interface Category {
  id: number;
  name: string;
  slug?: string;
  parent_id?: number | null;
  iconName?: string;
  color?: string;
  count?: number;
  listing_count?: number;
  translations?: {
    en: string;
    ru: string;
    sr: string;
  };
}

interface MobileCategoryTreeProps {
  categories: Category[];
  expandedCategories: Set<number>;
  onToggleExpand: (categoryId: number) => void;
  onSelectCategory: (categoryId: number) => void;
  locale?: string;
}

export const MobileCategoryTree: React.FC<MobileCategoryTreeProps> = ({
  categories,
  expandedCategories,
  onToggleExpand,
  onSelectCategory,
  locale = 'en',
}) => {
  // Построение иерархической структуры
  const buildHierarchy = (cats: Category[]): Map<number | null, Category[]> => {
    const hierarchy = new Map<number | null, Category[]>();

    cats.forEach((cat) => {
      const parentId = cat.parent_id || null;
      if (!hierarchy.has(parentId)) {
        hierarchy.set(parentId, []);
      }
      hierarchy.get(parentId)!.push(cat);
    });

    // Сортируем категории по имени
    hierarchy.forEach((children) => {
      children.sort((a, b) => a.name.localeCompare(b.name));
    });

    return hierarchy;
  };

  const hierarchy = buildHierarchy(categories);

  // Рекурсивный рендер категории
  const renderCategory = (
    category: Category,
    level: number = 0
  ): React.ReactElement => {
    const children = hierarchy.get(category.id) || [];
    const hasChildren = children.length > 0;
    const isExpanded = expandedCategories.has(category.id);
    const Icon = getCategoryIcon(category.iconName);

    // Получаем правильное имя категории в зависимости от локали
    const categoryName =
      category.translations?.[locale as keyof typeof category.translations] ||
      category.name;

    // Используем count или listing_count для отображения количества
    const itemCount = category.count || category.listing_count || 0;

    return (
      <div key={category.id}>
        <div
          className="flex items-center py-2 px-3 hover:bg-base-200 rounded-lg cursor-pointer"
          style={{ paddingLeft: `${12 + level * 20}px` }}
        >
          {/* Кнопка развертывания */}
          {hasChildren && (
            <button
              onClick={(e) => {
                e.stopPropagation();
                onToggleExpand(category.id);
              }}
              className="btn btn-ghost btn-circle btn-xs mr-2"
            >
              {isExpanded ? (
                <FiChevronDown className="w-3 h-3" />
              ) : (
                <FiChevronRight className="w-3 h-3" />
              )}
            </button>
          )}

          {/* Отступ для выравнивания категорий без детей */}
          {!hasChildren && <div className="w-8 mr-2" />}

          {/* Контент категории */}
          <div
            className="flex items-center flex-1 gap-2"
            onClick={() => onSelectCategory(category.id)}
          >
            <Icon className={`w-4 h-4 ${category.color || 'text-gray-600'}`} />
            <span className="text-sm">{categoryName}</span>
            {itemCount > 0 && (
              <span className="ml-auto text-xs text-base-content/40">
                {itemCount}
              </span>
            )}
          </div>
        </div>

        {/* Дочерние категории */}
        <AnimatePresence>
          {hasChildren && isExpanded && (
            <motion.div
              initial={{ height: 0, opacity: 0 }}
              animate={{ height: 'auto', opacity: 1 }}
              exit={{ height: 0, opacity: 0 }}
              transition={{ duration: 0.2 }}
              className="overflow-hidden"
            >
              {children.map((child) => renderCategory(child, level + 1))}
            </motion.div>
          )}
        </AnimatePresence>
      </div>
    );
  };

  // Корневые категории
  const rootCategories = hierarchy.get(null) || [];

  return (
    <div className="max-h-96 overflow-y-auto px-2">
      {rootCategories.map((category) => renderCategory(category, 0))}
    </div>
  );
};
