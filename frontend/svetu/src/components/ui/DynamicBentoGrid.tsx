'use client';

import React, { useEffect, useState } from 'react';

export interface BentoGridItem {
  id: string;
  title?: string;
  description?: string;
  content?: React.ReactNode;
  className?: string;
  colSpan?: number;
  rowSpan?: number;
  onClick?: () => void;
  href?: string;
  bgColor?: string;
  icon?: React.ReactNode;
}

interface DynamicBentoGridProps {
  items: BentoGridItem[];
  className?: string;
  variant?: 'default' | 'compact' | 'hero';
}

export const DynamicBentoGrid: React.FC<DynamicBentoGridProps> = ({
  items,
  className,
  variant = 'default',
}) => {
  const [mounted, setMounted] = useState(false);

  useEffect(() => {
    setMounted(true);
  }, []);

  const gridVariants = {
    default: 'grid-cols-1 md:grid-cols-2 lg:grid-cols-4 gap-4',
    compact: 'grid-cols-2 md:grid-cols-3 lg:grid-cols-6 gap-2',
    hero: 'grid-cols-1 md:grid-cols-3 lg:grid-cols-4 gap-4',
  };

  const handleClick = (item: BentoGridItem) => {
    if (item.onClick) {
      item.onClick();
    } else if (item.href) {
      window.location.href = item.href;
    }
  };

  return (
    <div
      className={`grid ${gridVariants[variant]} ${className || ''}`}
    >
      {items.map((item, index) => (
        <div
          key={item.id}
          className={`
            relative group cursor-pointer
            animate-fadeInUp
            ${item.className || ''}
          `}
          style={{
            gridColumn: item.colSpan ? `span ${item.colSpan}` : undefined,
            gridRow: item.rowSpan ? `span ${item.rowSpan}` : undefined,
            animationDelay: mounted ? `${index * 100}ms` : '0ms',
          }}
          onClick={() => handleClick(item)}
        >
          <div
            className={`
              h-full w-full rounded-xl shadow-lg overflow-hidden
              transform transition-all duration-300
              hover:shadow-xl hover:scale-[1.02] active:scale-[0.98]
              ${item.bgColor || 'bg-base-100'}
            `}
          >
            {/* Gradient overlay for depth */}
            <div className="absolute inset-0 bg-gradient-to-br from-transparent to-black/5 pointer-events-none" />
            
            {/* Content */}
            <div className="relative h-full p-6 flex flex-col">
              {item.icon && (
                <div className="mb-4">
                  {item.icon}
                </div>
              )}
              
              {item.title && (
                <h3 className="font-semibold text-lg mb-2 line-clamp-2">
                  {item.title}
                </h3>
              )}
              
              {item.description && (
                <p className="text-sm text-base-content/70 mb-4 line-clamp-3">
                  {item.description}
                </p>
              )}
              
              {item.content && (
                <div className="flex-1">
                  {item.content}
                </div>
              )}
              
              {/* Hover indicator */}
              <div className="absolute bottom-4 right-4 opacity-0 group-hover:opacity-100 transition-opacity">
                <svg
                  className="w-5 h-5 text-base-content/50"
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
              </div>
            </div>
          </div>
        </div>
      ))}
    </div>
  );
};