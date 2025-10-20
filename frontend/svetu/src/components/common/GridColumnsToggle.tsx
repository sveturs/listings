'use client';

import { useTranslations } from 'next-intl';

export type GridColumns = 1 | 2 | 3;

interface GridColumnsToggleProps {
  currentColumns: GridColumns;
  onColumnsChange: (columns: GridColumns) => void;
}

export default function GridColumnsToggle({
  currentColumns,
  onColumnsChange,
}: GridColumnsToggleProps) {
  const t = useTranslations('common');

  const columns: GridColumns[] = [1, 2, 3];

  const getColumnLabel = (col: GridColumns) => {
    switch (col) {
      case 1:
        return t('view.oneColumn');
      case 2:
        return t('view.twoColumns');
      case 3:
        return t('view.threeColumns');
      default:
        return t('columns', { count: col });
    }
  };

  return (
    <div className="flex gap-1 bg-base-200 p-1 rounded-lg">
      {columns.map((col) => (
        <button
          key={col}
          onClick={() => onColumnsChange(col)}
          className={`btn btn-sm ${
            currentColumns === col ? 'btn-primary' : 'btn-ghost'
          } min-w-[2.5rem]`}
          aria-label={getColumnLabel(col)}
          aria-current={currentColumns === col ? 'true' : 'false'}
        >
          <div className="flex gap-0.5" aria-hidden="true">
            {Array.from({ length: col }).map((_, i) => (
              <div
                key={i}
                className={`w-1 h-3 ${
                  currentColumns === col
                    ? 'bg-primary-content'
                    : 'bg-base-content'
                } rounded-sm`}
              />
            ))}
          </div>
        </button>
      ))}
    </div>
  );
}
