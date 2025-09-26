'use client';

import React from 'react';
import Link from 'next/link';
import { usePathname } from 'next/navigation';
import { ChevronRight, Home, X } from 'lucide-react';
import { useTranslations } from 'next-intl';

interface BreadcrumbItem {
  label: string;
  href?: string;
  isActive?: boolean;
}

interface ActiveFilter {
  key: string;
  value: string;
  label: string;
}

interface BreadcrumbsProps {
  items?: BreadcrumbItem[];
  activeFilters?: ActiveFilter[];
  onRemoveFilter?: (key: string) => void;
  className?: string;
}

export default function Breadcrumbs({
  items = [],
  activeFilters = [],
  onRemoveFilter,
  className = '',
}: BreadcrumbsProps) {
  const pathname = usePathname();
  const t = useTranslations('cars');

  // Generate breadcrumb items from pathname if not provided
  const breadcrumbItems =
    items.length > 0 ? items : generateFromPath(pathname, t);

  return (
    <nav className={`flex flex-col gap-2 ${className}`}>
      {/* Breadcrumb navigation */}
      <div className="flex items-center gap-2 text-sm">
        {/* Home link */}
        <Link
          href="/"
          className="text-gray-500 hover:text-primary transition-colors"
        >
          <Home className="w-4 h-4" />
        </Link>

        {/* Breadcrumb items */}
        {breadcrumbItems.map((item, index) => (
          <React.Fragment key={index}>
            <ChevronRight className="w-4 h-4 text-gray-400" />
            {item.href && !item.isActive ? (
              <Link
                href={item.href}
                className="text-gray-500 hover:text-primary transition-colors"
              >
                {item.label}
              </Link>
            ) : (
              <span
                className={`${
                  item.isActive ? 'text-gray-900 font-medium' : 'text-gray-500'
                }`}
              >
                {item.label}
              </span>
            )}
          </React.Fragment>
        ))}
      </div>

      {/* Active filters */}
      {activeFilters.length > 0 && (
        <div className="flex flex-wrap gap-2">
          <span className="text-sm text-gray-500">
            {t('breadcrumbs.filters')}:
          </span>
          {activeFilters.map((filter) => (
            <div
              key={filter.key}
              className="inline-flex items-center gap-1 px-2 py-1 bg-primary/10 text-primary rounded-md text-sm"
            >
              <span className="font-medium">{filter.label}:</span>
              <span>{filter.value}</span>
              {onRemoveFilter && (
                <button
                  onClick={() => onRemoveFilter(filter.key)}
                  className="ml-1 p-0.5 hover:bg-primary/20 rounded transition-colors"
                  aria-label={t('breadcrumbs.removeFilter')}
                >
                  <X className="w-3 h-3" />
                </button>
              )}
            </div>
          ))}
        </div>
      )}
    </nav>
  );
}

// Helper function to generate breadcrumb items from path
function generateFromPath(pathname: string, t: any): BreadcrumbItem[] {
  const segments = pathname.split('/').filter(Boolean);
  const items: BreadcrumbItem[] = [];

  // Skip locale segment
  const startIndex = segments[0]?.length === 2 ? 1 : 0;

  for (let i = startIndex; i < segments.length; i++) {
    const segment = segments[i];
    const isLast = i === segments.length - 1;

    // Build href for this segment
    const href = '/' + segments.slice(0, i + 1).join('/');

    // Translate segment
    let label = segment;
    switch (segment) {
      case 'cars':
        label = t('breadcrumbs.cars');
        break;
      case 'search':
        label = t('breadcrumbs.search');
        break;
      case 'favorites':
        label = t('breadcrumbs.favorites');
        break;
      case 'compare':
        label = t('breadcrumbs.compare');
        break;
      default:
        // Capitalize first letter for other segments
        label =
          segment.charAt(0).toUpperCase() + segment.slice(1).replace(/-/g, ' ');
    }

    items.push({
      label,
      href: isLast ? undefined : href,
      isActive: isLast,
    });
  }

  return items;
}
