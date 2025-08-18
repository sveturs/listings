'use client';

import React, { useState, useEffect, useRef } from 'react';
import { motion, AnimatePresence } from 'framer-motion';
import { FiChevronDown, FiChevronRight, FiCheck } from 'react-icons/fi';
import {
  BsHouseDoor,
  BsLaptop,
  BsBriefcase,
  BsPalette,
  BsTools,
  BsPhone,
  BsGem,
  BsHandbag,
} from 'react-icons/bs';
import { FaCar, FaTshirt } from 'react-icons/fa';

interface Category {
  id: string | number;
  name: string;
  slug?: string;
  count?: number;
  icon?: React.ComponentType<any>;
  color?: string;
  children?: Category[];
  parent_id?: string | number | null;
}

interface NestedCategorySelectorProps {
  categories: Category[];
  selectedCategory?: string | number;
  onChange: (categoryId: string | number) => void;
  placeholder?: string;
  showCounts?: boolean;
  className?: string;
}

// Маппинг иконок для категорий
const categoryIconMap: { [key: string]: React.ComponentType<any> } = {
  'real-estate': BsHouseDoor,
  realestate: BsHouseDoor,
  automotive: FaCar,
  auto: FaCar,
  electronics: BsLaptop,
  fashion: FaTshirt,
  jobs: BsBriefcase,
  job: BsBriefcase,
  services: BsTools,
  'hobbies-entertainment': BsPalette,
  hobby: BsPalette,
  'home-garden': BsHandbag,
  home: BsHandbag,
  industrial: BsTools,
  'food-beverages': BsPhone,
  'books-stationery': BsGem,
  'antiques-art': BsPalette,
};

// Маппинг цветов для категорий
const categoryColorMap: { [key: string]: string } = {
  'real-estate': 'text-blue-600',
  realestate: 'text-blue-600',
  automotive: 'text-red-600',
  auto: 'text-red-600',
  electronics: 'text-purple-600',
  fashion: 'text-pink-600',
  jobs: 'text-green-600',
  job: 'text-green-600',
  services: 'text-orange-600',
  'hobbies-entertainment': 'text-indigo-600',
  hobby: 'text-indigo-600',
  'home-garden': 'text-yellow-600',
  home: 'text-yellow-600',
  industrial: 'text-gray-600',
  'food-beverages': 'text-teal-600',
  'books-stationery': 'text-cyan-600',
  'antiques-art': 'text-rose-600',
};

export const NestedCategorySelector: React.FC<NestedCategorySelectorProps> = ({
  categories,
  selectedCategory,
  onChange,
  placeholder = 'Все категории',
  showCounts = true,
  className = '',
}) => {
  const [isOpen, setIsOpen] = useState(false);
  const [expandedCategories, setExpandedCategories] = useState<
    Set<string | number>
  >(new Set());
  const [hierarchicalCategories, setHierarchicalCategories] = useState<
    Category[]
  >([]);
  const [selectedCategoryData, setSelectedCategoryData] =
    useState<Category | null>(null);
  const dropdownRef = useRef<HTMLDivElement>(null);

  // Построение иерархической структуры категорий
  useEffect(() => {
    const buildHierarchy = (cats: Category[]): Category[] => {
      const categoryMap = new Map<string | number, Category>();
      const rootCategories: Category[] = [];

      // Создаем копии категорий с иконками и цветами
      cats.forEach((cat) => {
        const IconComponent =
          categoryIconMap[cat.slug || ''] ||
          categoryIconMap[cat.id.toString()] ||
          BsHandbag;

        const enrichedCat = {
          ...cat,
          icon: IconComponent,
          color:
            cat.color ||
            categoryColorMap[cat.slug || ''] ||
            categoryColorMap[cat.id.toString()] ||
            'text-gray-600',
          children: [],
        };
        categoryMap.set(cat.id, enrichedCat);
      });

      // Строим иерархию
      categoryMap.forEach((cat) => {
        if (!cat.parent_id || cat.parent_id === null) {
          rootCategories.push(cat);
        } else {
          const parent = categoryMap.get(cat.parent_id);
          if (parent) {
            if (!parent.children) parent.children = [];
            parent.children.push(cat);
          } else {
            // Если родитель не найден, добавляем как корневую
            rootCategories.push(cat);
          }
        }
      });

      // Сортируем категории по имени
      const sortCategories = (cats: Category[]): Category[] => {
        return cats
          .sort((a, b) => a.name.localeCompare(b.name))
          .map((cat) => ({
            ...cat,
            children: cat.children ? sortCategories(cat.children) : [],
          }));
      };

      return sortCategories(rootCategories);
    };

    setHierarchicalCategories(buildHierarchy(categories));
  }, [categories]);

  // Находим выбранную категорию
  useEffect(() => {
    if (selectedCategory && selectedCategory !== 'all') {
      const findCategory = (cats: Category[]): Category | null => {
        for (const cat of cats) {
          if (cat.id === selectedCategory) return cat;
          if (cat.children) {
            const found = findCategory(cat.children);
            if (found) return found;
          }
        }
        return null;
      };

      const found = findCategory(hierarchicalCategories);
      setSelectedCategoryData(found);
    } else {
      setSelectedCategoryData(null);
    }
  }, [selectedCategory, hierarchicalCategories]);

  // Закрытие при клике вне компонента
  useEffect(() => {
    const handleClickOutside = (event: MouseEvent) => {
      if (
        dropdownRef.current &&
        !dropdownRef.current.contains(event.target as Node)
      ) {
        setIsOpen(false);
      }
    };

    document.addEventListener('mousedown', handleClickOutside);
    return () => document.removeEventListener('mousedown', handleClickOutside);
  }, []);

  // Переключение развернутого состояния категории
  const toggleExpanded = (
    categoryId: string | number,
    event: React.MouseEvent
  ) => {
    event.stopPropagation();
    const newExpanded = new Set(expandedCategories);
    if (newExpanded.has(categoryId)) {
      newExpanded.delete(categoryId);
    } else {
      newExpanded.add(categoryId);
    }
    setExpandedCategories(newExpanded);
  };

  // Выбор категории
  const handleSelectCategory = (category: Category | null) => {
    onChange(category ? category.id : 'all');
    setIsOpen(false);
  };

  // Подсчет общего количества товаров в категории и подкатегориях
  const getTotalCount = (category: Category): number => {
    let total = category.count || 0;
    if (category.children) {
      category.children.forEach((child) => {
        total += getTotalCount(child);
      });
    }
    return total;
  };

  // Рендер категории
  const renderCategory = (
    category: Category,
    level: number = 0
  ): React.ReactElement => {
    const hasChildren = category.children && category.children.length > 0;
    const isExpanded = expandedCategories.has(category.id);
    const isSelected = selectedCategory === category.id;
    const Icon = category.icon || BsHandbag;
    const totalCount = getTotalCount(category);

    return (
      <div key={category.id} className="w-full">
        <div
          className={`w-full flex items-center gap-2 px-3 py-2 rounded-lg hover:bg-base-200 transition-colors cursor-pointer ${
            isSelected ? 'bg-primary/10 text-primary font-medium' : ''
          }`}
          style={{ paddingLeft: `${12 + level * 16}px` }}
          onClick={() => handleSelectCategory(category)}
        >
          {/* Кнопка развертывания */}
          {hasChildren && (
            <div
              onClick={(e) => {
                e.stopPropagation();
                toggleExpanded(category.id, e);
              }}
              className="btn btn-ghost btn-circle btn-xs"
            >
              {isExpanded ? (
                <FiChevronDown className="w-3 h-3" />
              ) : (
                <FiChevronRight className="w-3 h-3" />
              )}
            </div>
          )}

          {/* Пустое место вместо кнопки для выравнивания */}
          {!hasChildren && <div className="w-6" />}

          {/* Иконка категории */}
          <Icon className={`w-4 h-4 ${category.color}`} />

          {/* Название категории */}
          <span className="flex-1 text-left">{category.name}</span>

          {/* Количество товаров */}
          {showCounts && totalCount > 0 && (
            <span className="text-xs text-base-content/60">
              {totalCount > 999
                ? `${Math.floor(totalCount / 1000)}K+`
                : totalCount}
            </span>
          )}

          {/* Галочка для выбранной категории */}
          {isSelected && <FiCheck className="w-4 h-4 text-primary" />}
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
              {category.children!.map((child) =>
                renderCategory(child, level + 1)
              )}
            </motion.div>
          )}
        </AnimatePresence>
      </div>
    );
  };

  return (
    <div ref={dropdownRef} className={`relative ${className}`}>
      {/* Кнопка открытия */}
      <button
        onClick={() => setIsOpen(!isOpen)}
        className="select select-bordered w-full flex items-center justify-between"
      >
        <div className="flex items-center gap-2">
          {selectedCategoryData ? (
            <>
              {selectedCategoryData.icon && (
                <selectedCategoryData.icon
                  className={`w-4 h-4 ${selectedCategoryData.color}`}
                />
              )}
              <span>{selectedCategoryData.name}</span>
            </>
          ) : (
            <span>{placeholder}</span>
          )}
        </div>
        <FiChevronDown
          className={`w-4 h-4 transition-transform ${isOpen ? 'rotate-180' : ''}`}
        />
      </button>

      {/* Выпадающий список */}
      <AnimatePresence>
        {isOpen && (
          <motion.div
            initial={{ opacity: 0, y: -10 }}
            animate={{ opacity: 1, y: 0 }}
            exit={{ opacity: 0, y: -10 }}
            transition={{ duration: 0.2 }}
            className="absolute top-full left-0 right-0 z-50 mt-1 bg-base-100 rounded-lg shadow-xl border border-base-300 max-h-[400px] overflow-y-auto min-w-[250px]"
          >
            <div className="p-2">
              {/* Опция "Все категории" */}
              <div
                onClick={() => handleSelectCategory(null)}
                className={`w-full flex items-center gap-2 px-3 py-2 rounded-lg hover:bg-base-200 transition-colors cursor-pointer ${
                  !selectedCategory || selectedCategory === 'all'
                    ? 'bg-primary/10 text-primary font-medium'
                    : ''
                }`}
              >
                <div className="w-6" />
                <span className="flex-1 text-left">{placeholder}</span>
                {(!selectedCategory || selectedCategory === 'all') && (
                  <FiCheck className="w-4 h-4 text-primary" />
                )}
              </div>

              <div className="divider my-1"></div>

              {/* Список категорий */}
              {hierarchicalCategories.map((category) =>
                renderCategory(category)
              )}
            </div>
          </motion.div>
        )}
      </AnimatePresence>
    </div>
  );
};
