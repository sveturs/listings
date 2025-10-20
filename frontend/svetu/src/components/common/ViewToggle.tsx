'use client';

import { useTranslations } from 'next-intl';
import { Squares2X2Icon, ListBulletIcon } from '@heroicons/react/24/outline';

export type ViewMode = 'grid' | 'list';

interface ViewToggleProps {
  currentView: ViewMode;
  onViewChange: (view: ViewMode) => void;
}

export default function ViewToggle({
  currentView,
  onViewChange,
}: ViewToggleProps) {
  const t = useTranslations('common');

  return (
    <div className="flex gap-1 bg-base-200 p-1 rounded-lg">
      <button
        onClick={() => onViewChange('grid')}
        className={`btn btn-sm ${
          currentView === 'grid' ? 'btn-primary' : 'btn-ghost'
        }`}
        aria-label={t('view.gridView')}
        aria-current={currentView === 'grid' ? 'true' : 'false'}
      >
        <Squares2X2Icon className="w-4 h-4" aria-hidden="true" />
      </button>
      <button
        onClick={() => onViewChange('list')}
        className={`btn btn-sm ${
          currentView === 'list' ? 'btn-primary' : 'btn-ghost'
        }`}
        aria-label={t('view.listView')}
        aria-current={currentView === 'list' ? 'true' : 'false'}
      >
        <ListBulletIcon className="w-4 h-4" aria-hidden="true" />
      </button>
    </div>
  );
}
