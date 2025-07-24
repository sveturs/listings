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
  const t = useTranslations();

  const columns: GridColumns[] = [1, 2, 3];

  return (
    <div className="flex gap-1 bg-base-200 p-1 rounded-lg">
      {columns.map((col) => (
        <button
          key={col}
          onClick={() => onColumnsChange(col)}
          className={`btn btn-sm ${
            currentColumns === col ? 'btn-primary' : 'btn-ghost'
          } min-w-[2.5rem]`}
          title={t('common.columns', { count: col })}
        >
          <div className="flex gap-0.5">
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
