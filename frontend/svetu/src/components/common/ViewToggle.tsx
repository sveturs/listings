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
  const t = useTranslations();

  return (
    <div className="flex gap-1 bg-base-200 p-1 rounded-lg">
      <button
        onClick={() => onViewChange('grid')}
        className={`btn btn-sm ${
          currentView === 'grid' ? 'btn-primary' : 'btn-ghost'
        }`}
        title={t('common.viewGrid')}
      >
        <Squares2X2Icon className="w-4 h-4" />
      </button>
      <button
        onClick={() => onViewChange('list')}
        className={`btn btn-sm ${
          currentView === 'list' ? 'btn-primary' : 'btn-ghost'
        }`}
        title={t('common.viewList')}
      >
        <ListBulletIcon className="w-4 h-4" />
      </button>
    </div>
  );
}
