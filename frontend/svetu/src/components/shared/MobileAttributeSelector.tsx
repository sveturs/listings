// Mobile Attribute Selector Component
// День 25: Мобильная оптимизация системы атрибутов

import React, { useState, useRef, useEffect, useCallback } from 'react';
import { motion, AnimatePresence } from 'framer-motion';
import {
  ChevronRightIcon,
  CheckIcon,
  XMarkIcon,
  MagnifyingGlassIcon,
  FunnelIcon,
  ArrowLeftIcon,
  SparklesIcon,
  ClockIcon,
} from '@heroicons/react/24/outline';
import { useTouchGestures, useSwipeNavigation } from '@/hooks/useTouchGestures';
import type { components } from '@/types/generated/api';

type UnifiedAttribute =
  components['schemas']['backend_internal_domain_models.UnifiedAttribute'];
type AttributeValue =
  components['schemas']['backend_internal_domain_models.UnifiedAttributeValue'];

interface MobileAttributeSelectorProps {
  attributes: UnifiedAttribute[];
  selectedValues: Record<string, any>;
  onValueChange: (attributeId: string, value: any) => void;
  onClose: () => void;
  categoryName?: string;
}

export const MobileAttributeSelector: React.FC<
  MobileAttributeSelectorProps
> = ({
  attributes,
  selectedValues,
  onValueChange,
  onClose,
  categoryName = 'Attributes',
}) => {
  const [activeAttribute, setActiveAttribute] =
    useState<UnifiedAttribute | null>(null);
  const [searchQuery, setSearchQuery] = useState('');
  const [recentSelections, setRecentSelections] = useState<string[]>([]);
  const [popularAttributes, setPopularAttributes] = useState<string[]>([]);
  const containerRef = useRef<HTMLDivElement>(null);
  const [viewStack, setViewStack] = useState<'list' | 'detail'>('list');

  // Touch gestures для навигации
  useSwipeNavigation(containerRef, (direction) => {
    if (direction === 'next' && viewStack === 'list' && attributes.length > 0) {
      setActiveAttribute(attributes[0]);
      setViewStack('detail');
    } else if (direction === 'prev' && viewStack === 'detail') {
      setActiveAttribute(null);
      setViewStack('list');
    }
  });

  // Фильтрация атрибутов по поиску
  const filteredAttributes = attributes.filter(
    (attr) =>
      attr.label?.toLowerCase().includes(searchQuery.toLowerCase()) ||
      attr.key?.toLowerCase().includes(searchQuery.toLowerCase())
  );

  // Группировка атрибутов по секциям
  const groupedAttributes = filteredAttributes.reduce(
    (acc, attr) => {
      const section = attr.section || 'general';
      if (!acc[section]) acc[section] = [];
      acc[section].push(attr);
      return acc;
    },
    {} as Record<string, UnifiedAttribute[]>
  );

  // Загрузка популярных атрибутов
  useEffect(() => {
    // Симуляция загрузки популярных атрибутов
    const popular = attributes
      .filter((attr) => attr.is_required)
      .slice(0, 5)
      .map((attr) => attr.id.toString());
    setPopularAttributes(popular);
  }, [attributes]);

  // Сохранение недавних выборов
  const handleAttributeSelect = useCallback((attr: UnifiedAttribute) => {
    setActiveAttribute(attr);
    setViewStack('detail');

    // Обновление недавних выборов
    setRecentSelections((prev) => {
      const updated = [
        attr.id.toString(),
        ...prev.filter((id) => id !== attr.id.toString()),
      ];
      return updated.slice(0, 5);
    });
  }, []);

  // Быстрый выбор значения
  const handleQuickSelect = useCallback(
    (attributeId: string, value: any) => {
      onValueChange(attributeId, value);

      // Haptic feedback на мобильных устройствах
      if ('vibrate' in navigator) {
        navigator.vibrate(10);
      }
    },
    [onValueChange]
  );

  // Рендер значения атрибута
  const renderAttributeValue = (attr: UnifiedAttribute) => {
    const value = selectedValues[attr.id];

    switch (attr.type) {
      case 'select':
        return (
          <div className="space-y-2">
            {attr.values?.map((option) => (
              <motion.button
                key={option.id}
                whileTap={{ scale: 0.95 }}
                onClick={() =>
                  handleQuickSelect(attr.id.toString(), option.value)
                }
                className={`
                  w-full p-4 rounded-xl flex items-center justify-between
                  transition-all duration-200
                  ${
                    value === option.value
                      ? 'bg-primary text-white shadow-lg'
                      : 'bg-base-200 hover:bg-base-300'
                  }
                `}
              >
                <span className="font-medium">{option.display_value}</span>
                {value === option.value && <CheckIcon className="w-5 h-5" />}
              </motion.button>
            ))}
          </div>
        );

      case 'checkbox':
        return (
          <motion.button
            whileTap={{ scale: 0.95 }}
            onClick={() => handleQuickSelect(attr.id.toString(), !value)}
            className={`
              w-full p-6 rounded-xl flex items-center justify-center
              transition-all duration-200
              ${value ? 'bg-primary text-white' : 'bg-base-200'}
            `}
          >
            <div className="flex items-center gap-4">
              <div
                className={`
                w-8 h-8 rounded-lg border-2 flex items-center justify-center
                ${value ? 'border-white bg-white' : 'border-base-content/30'}
              `}
              >
                {value && <CheckIcon className="w-5 h-5 text-primary" />}
              </div>
              <span className="text-lg font-medium">
                {value ? 'Enabled' : 'Disabled'}
              </span>
            </div>
          </motion.button>
        );

      case 'range':
        const [min, max] = attr.validation?.range || [0, 100];
        return (
          <div className="space-y-4">
            <input
              type="range"
              min={min}
              max={max}
              value={value || min}
              onChange={(e) =>
                handleQuickSelect(attr.id.toString(), Number(e.target.value))
              }
              className="w-full range range-primary"
            />
            <div className="flex justify-between text-sm">
              <span>{min}</span>
              <span className="font-bold text-primary">{value || min}</span>
              <span>{max}</span>
            </div>
          </div>
        );

      default:
        return (
          <input
            type="text"
            value={value || ''}
            onChange={(e) =>
              handleQuickSelect(attr.id.toString(), e.target.value)
            }
            placeholder={`Enter ${attr.label}`}
            className="w-full p-4 rounded-xl bg-base-200 focus:bg-base-100 
                     focus:ring-2 focus:ring-primary transition-all"
          />
        );
    }
  };

  return (
    <div className="fixed inset-0 z-50 bg-base-100">
      <AnimatePresence mode="wait">
        {viewStack === 'list' ? (
          <motion.div
            key="list"
            initial={{ x: -20, opacity: 0 }}
            animate={{ x: 0, opacity: 1 }}
            exit={{ x: -20, opacity: 0 }}
            className="h-full flex flex-col"
            ref={containerRef}
          >
            {/* Header */}
            <div className="sticky top-0 z-10 bg-base-100 border-b border-base-300">
              <div className="flex items-center justify-between p-4">
                <div className="flex items-center gap-3">
                  <FunnelIcon className="w-6 h-6 text-primary" />
                  <h2 className="text-lg font-bold">{categoryName} Filters</h2>
                </div>
                <button onClick={onClose} className="btn btn-ghost btn-circle">
                  <XMarkIcon className="w-6 h-6" />
                </button>
              </div>

              {/* Search Bar */}
              <div className="px-4 pb-3">
                <div className="relative">
                  <MagnifyingGlassIcon className="absolute left-3 top-1/2 -translate-y-1/2 w-5 h-5 text-base-content/50" />
                  <input
                    type="text"
                    value={searchQuery}
                    onChange={(e) => setSearchQuery(e.target.value)}
                    placeholder="Search filters..."
                    className="w-full pl-10 pr-4 py-3 rounded-xl bg-base-200 
                             focus:bg-base-100 focus:ring-2 focus:ring-primary"
                  />
                </div>
              </div>
            </div>

            {/* Content */}
            <div className="flex-1 overflow-y-auto pb-20">
              {/* Quick Filters */}
              {searchQuery === '' && (
                <>
                  {/* Recent Selections */}
                  {recentSelections.length > 0 && (
                    <div className="px-4 py-3">
                      <div className="flex items-center gap-2 mb-3">
                        <ClockIcon className="w-4 h-4 text-base-content/50" />
                        <span className="text-sm font-medium text-base-content/70">
                          Recent
                        </span>
                      </div>
                      <div className="flex gap-2 overflow-x-auto pb-2">
                        {recentSelections.map((id) => {
                          const attr = attributes.find(
                            (a) => a.id.toString() === id
                          );
                          if (!attr) return null;
                          return (
                            <motion.button
                              key={id}
                              whileTap={{ scale: 0.95 }}
                              onClick={() => handleAttributeSelect(attr)}
                              className="px-4 py-2 rounded-full bg-base-200 
                                       whitespace-nowrap text-sm font-medium"
                            >
                              {attr.label}
                            </motion.button>
                          );
                        })}
                      </div>
                    </div>
                  )}

                  {/* Popular Filters */}
                  {popularAttributes.length > 0 && (
                    <div className="px-4 py-3">
                      <div className="flex items-center gap-2 mb-3">
                        <SparklesIcon className="w-4 h-4 text-primary" />
                        <span className="text-sm font-medium text-base-content/70">
                          Popular
                        </span>
                      </div>
                      <div className="grid grid-cols-2 gap-2">
                        {popularAttributes.slice(0, 4).map((id) => {
                          const attr = attributes.find(
                            (a) => a.id.toString() === id
                          );
                          if (!attr) return null;
                          return (
                            <motion.button
                              key={id}
                              whileTap={{ scale: 0.95 }}
                              onClick={() => handleAttributeSelect(attr)}
                              className="p-3 rounded-xl bg-primary/10 text-primary 
                                       font-medium text-sm"
                            >
                              {attr.label}
                              {selectedValues[attr.id] && (
                                <CheckIcon className="w-4 h-4 ml-1 inline" />
                              )}
                            </motion.button>
                          );
                        })}
                      </div>
                    </div>
                  )}
                </>
              )}

              {/* Attribute Groups */}
              <div className="px-4 space-y-4">
                {Object.entries(groupedAttributes).map(([section, attrs]) => (
                  <div key={section}>
                    <h3 className="text-sm font-bold text-base-content/70 uppercase mb-2">
                      {section}
                    </h3>
                    <div className="space-y-2">
                      {attrs.map((attr) => (
                        <motion.button
                          key={attr.id}
                          whileTap={{ scale: 0.98 }}
                          onClick={() => handleAttributeSelect(attr)}
                          className="w-full p-4 rounded-xl bg-base-200 hover:bg-base-300 
                                   flex items-center justify-between group"
                        >
                          <div className="flex items-start gap-3">
                            <div className="flex-1 text-left">
                              <div className="font-medium">
                                {attr.label}
                                {attr.is_required && (
                                  <span className="ml-1 text-error">*</span>
                                )}
                              </div>
                              {attr.description && (
                                <div className="text-xs text-base-content/60 mt-1">
                                  {attr.description}
                                </div>
                              )}
                              {selectedValues[attr.id] && (
                                <div className="text-sm text-primary mt-1 font-medium">
                                  {typeof selectedValues[attr.id] === 'boolean'
                                    ? 'Enabled'
                                    : selectedValues[attr.id]}
                                </div>
                              )}
                            </div>
                          </div>
                          <ChevronRightIcon
                            className="w-5 h-5 text-base-content/50 
                                                       group-hover:text-primary transition-colors"
                          />
                        </motion.button>
                      ))}
                    </div>
                  </div>
                ))}
              </div>
            </div>

            {/* Bottom Action Bar */}
            <div className="sticky bottom-0 bg-base-100 border-t border-base-300 p-4">
              <div className="flex gap-3">
                <button
                  onClick={() => {
                    // Clear all filters
                    Object.keys(selectedValues).forEach((key) => {
                      onValueChange(key, null);
                    });
                  }}
                  className="btn btn-outline flex-1"
                >
                  Clear All
                </button>
                <button onClick={onClose} className="btn btn-primary flex-1">
                  Apply Filters
                </button>
              </div>
            </div>
          </motion.div>
        ) : (
          <motion.div
            key="detail"
            initial={{ x: 20, opacity: 0 }}
            animate={{ x: 0, opacity: 1 }}
            exit={{ x: 20, opacity: 0 }}
            className="h-full flex flex-col"
            ref={containerRef}
          >
            {activeAttribute && (
              <>
                {/* Header */}
                <div className="sticky top-0 z-10 bg-base-100 border-b border-base-300 p-4">
                  <div className="flex items-center gap-3">
                    <button
                      onClick={() => {
                        setActiveAttribute(null);
                        setViewStack('list');
                      }}
                      className="btn btn-ghost btn-circle"
                    >
                      <ArrowLeftIcon className="w-5 h-5" />
                    </button>
                    <div className="flex-1">
                      <h2 className="text-lg font-bold">
                        {activeAttribute.label}
                      </h2>
                      {activeAttribute.description && (
                        <p className="text-sm text-base-content/60">
                          {activeAttribute.description}
                        </p>
                      )}
                    </div>
                  </div>
                </div>

                {/* Content */}
                <div className="flex-1 overflow-y-auto p-4">
                  {renderAttributeValue(activeAttribute)}
                </div>

                {/* Bottom Action */}
                <div className="sticky bottom-0 bg-base-100 border-t border-base-300 p-4">
                  <button
                    onClick={() => {
                      setActiveAttribute(null);
                      setViewStack('list');
                    }}
                    className="btn btn-primary btn-block"
                  >
                    Done
                  </button>
                </div>
              </>
            )}
          </motion.div>
        )}
      </AnimatePresence>
    </div>
  );
};
