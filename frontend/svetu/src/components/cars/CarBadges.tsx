'use client';

import React from 'react';
import { useTranslations } from 'next-intl';
import { TrendingDown, Award, Clock, Zap } from 'lucide-react';

interface CarBadgesProps {
  isNew?: boolean;
  priceReduction?: number;
  isTop?: boolean;
  isFeatured?: boolean;
  createdAt?: string;
}

export const CarBadges: React.FC<CarBadgesProps> = ({
  isNew,
  priceReduction,
  isTop,
  isFeatured,
  createdAt,
}) => {
  const t = useTranslations('cars');

  const isRecent = () => {
    if (!createdAt) return false;
    const created = new Date(createdAt);
    const now = new Date();
    const diffHours = (now.getTime() - created.getTime()) / (1000 * 60 * 60);
    return diffHours < 24;
  };

  const badges = [];

  if (isNew || isRecent()) {
    badges.push(
      <span key="new" className="badge badge-success badge-sm gap-1">
        <Clock className="w-3 h-3" />
        {t('badges.new')}
      </span>
    );
  }

  if (priceReduction && priceReduction > 0) {
    badges.push(
      <span key="discount" className="badge badge-error badge-sm gap-1">
        <TrendingDown className="w-3 h-3" />-{priceReduction}%
      </span>
    );
  }

  if (isTop) {
    badges.push(
      <span key="top" className="badge badge-warning badge-sm gap-1">
        <Award className="w-3 h-3" />
        {t('badges.top')}
      </span>
    );
  }

  if (isFeatured) {
    badges.push(
      <span key="featured" className="badge badge-primary badge-sm gap-1">
        <Zap className="w-3 h-3" />
        {t('badges.featured')}
      </span>
    );
  }

  if (badges.length === 0) return null;

  return (
    <div className="absolute top-2 left-2 z-10 flex flex-wrap gap-1">
      {badges}
    </div>
  );
};

export default CarBadges;
