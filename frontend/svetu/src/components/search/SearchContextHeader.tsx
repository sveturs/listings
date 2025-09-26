'use client';

import React from 'react';
import { useTranslations } from 'next-intl';
import { SearchContextConfig } from '@/types/searchContext';
import { motion } from 'framer-motion';

interface SearchContextHeaderProps {
  context: SearchContextConfig;
  totalResults?: number;
}

export const SearchContextHeader: React.FC<SearchContextHeaderProps> = ({
  context,
  totalResults = 0,
}) => {
  const t = useTranslations();

  // Если нет специального контекста, не показываем заголовок
  if (context.id === 'default') {
    return null;
  }

  return (
    <motion.div
      initial={{ opacity: 0, y: -20 }}
      animate={{ opacity: 1, y: 0 }}
      transition={{ duration: 0.5 }}
      className={`hero min-h-[300px] bg-gradient-to-r ${context.bgGradient} text-white mb-8 rounded-xl overflow-hidden`}
    >
      <div className="hero-content text-center">
        <div className="max-w-md">
          {/* Иконка контекста */}
          <motion.div
            initial={{ scale: 0 }}
            animate={{ scale: 1 }}
            transition={{
              delay: 0.2,
              type: 'spring',
              stiffness: 260,
              damping: 20,
            }}
            className="text-6xl mb-4"
          >
            {context.heroIcon}
          </motion.div>

          {/* Заголовок */}
          <h1 className="text-5xl font-bold mb-4">{t(context.heroTitle)}</h1>

          {/* Описание */}
          <p className="text-xl mb-6 opacity-90">
            {t(context.heroDescription)}
          </p>

          {/* Статистика */}
          {totalResults > 0 && (
            <motion.div
              initial={{ opacity: 0 }}
              animate={{ opacity: 1 }}
              transition={{ delay: 0.4 }}
              className="badge badge-lg badge-ghost text-white border-white/30"
            >
              {totalResults.toLocaleString()} {t('search.resultsFound')}
            </motion.div>
          )}

          {/* Кастомный баннер если есть */}
          {context.customBanner && (
            <motion.div
              initial={{ opacity: 0, x: -20 }}
              animate={{ opacity: 1, x: 0 }}
              transition={{ delay: 0.6 }}
              className="mt-6 p-4 bg-white/10 backdrop-blur-sm rounded-lg"
            >
              <div className="flex items-center gap-3">
                {context.customBanner.icon && (
                  <span className="text-2xl">{context.customBanner.icon}</span>
                )}
                <div className="text-left">
                  <h3 className="font-semibold">
                    {t(context.customBanner.title)}
                  </h3>
                  <p className="text-sm opacity-90">
                    {t(context.customBanner.subtitle)}
                  </p>
                </div>
              </div>
              {context.customBanner.cta && (
                <button className="btn btn-sm btn-ghost text-white mt-2">
                  {t(context.customBanner.cta)}
                </button>
              )}
            </motion.div>
          )}
        </div>
      </div>

      {/* Декоративные элементы для автомобильного контекста */}
      {context.id === 'automotive' && (
        <>
          <div className="absolute top-0 left-0 w-32 h-32 bg-white/5 rounded-full blur-3xl"></div>
          <div className="absolute bottom-0 right-0 w-48 h-48 bg-white/5 rounded-full blur-3xl"></div>
          <div className="absolute top-1/2 left-1/4 w-24 h-24 bg-white/5 rounded-full blur-2xl"></div>
        </>
      )}

      {/* Декоративные элементы для недвижимости */}
      {context.id === 'real-estate' && (
        <>
          <div className="absolute top-0 right-0 w-40 h-40 bg-white/5 rounded-full blur-3xl"></div>
          <div className="absolute bottom-0 left-0 w-32 h-32 bg-white/5 rounded-full blur-3xl"></div>
        </>
      )}

      {/* Декоративные элементы для электроники */}
      {context.id === 'electronics' && (
        <>
          <div className="absolute top-1/3 right-1/4 w-20 h-20 bg-white/10 rounded blur-xl"></div>
          <div className="absolute bottom-1/3 left-1/4 w-16 h-16 bg-white/10 rounded blur-xl"></div>
          <div className="absolute top-2/3 right-1/3 w-24 h-24 bg-white/5 rounded-full blur-2xl"></div>
        </>
      )}
    </motion.div>
  );
};
