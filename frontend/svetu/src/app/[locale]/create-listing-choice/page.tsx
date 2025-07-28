'use client';

import React, { useState } from 'react';
import { useRouter, useParams } from 'next/navigation';
import { useTranslations } from 'next-intl';
import { motion } from 'framer-motion';
import {
  Brain,
  Sparkles,
  Package,
  Clock,
  TrendingUp,
  Shield,
  Zap,
  Check,
  FileText,
  ImageIcon,
  MapPin,
  DollarSign,
  ChevronRight,
  Info,
} from 'lucide-react';
import { useAuthContext } from '@/contexts/AuthContext';
import { toast } from '@/utils/toast';

export default function CreateListingChoicePage() {
  const router = useRouter();
  const params = useParams();
  const locale = params.locale as string;
  const t = useTranslations('create_listing');
  const { user } = useAuthContext();
  const [selectedOption, setSelectedOption] = useState<'free' | 'ai' | null>(
    null
  );
  const [isLoading, setIsLoading] = useState(false);

  const handleOptionSelect = (option: 'free' | 'ai') => {
    if (!user) {
      toast.error(t('auth_required'));
      router.push('/');
      return;
    }

    setSelectedOption(option);
    setIsLoading(true);

    // Небольшая задержка для визуального эффекта выбора
    setTimeout(() => {
      if (option === 'free') {
        router.push(`/${locale}/create-listing-smart`);
      } else {
        router.push(`/${locale}/create-listing-ai`);
      }
    }, 300);
  };

  const freeFeatures = [
    { icon: FileText, text: t('choice.free_option.features.templates') },
    { icon: ImageIcon, text: t('choice.free_option.features.drag_drop') },
    { icon: TrendingUp, text: t('choice.free_option.features.price_compare') },
    { icon: Clock, text: t('choice.free_option.features.optimal_time') },
    { icon: MapPin, text: t('choice.free_option.features.social_import') },
    { icon: Shield, text: t('choice.free_option.features.error_check') },
  ];

  const aiFeatures = [
    { icon: Brain, text: t('choice.ai_option.features.photo_recognition') },
    { icon: Sparkles, text: t('choice.ai_option.features.auto_description') },
    { icon: DollarSign, text: t('choice.ai_option.features.price_analysis') },
    { icon: FileText, text: t('choice.ai_option.features.ab_testing') },
    { icon: Check, text: t('choice.ai_option.features.multilingual') },
    { icon: TrendingUp, text: t('choice.ai_option.features.forecast') },
  ];

  return (
    <div className="min-h-screen bg-gradient-to-br from-base-100 to-base-200">
      <div className="container mx-auto px-4 py-8">
        {/* Header */}
        <motion.div
          initial={{ y: -20, opacity: 0 }}
          animate={{ y: 0, opacity: 1 }}
          className="text-center mb-12"
        >
          <h1 className="text-4xl lg:text-5xl font-bold mb-4">
            {t('choice.title')}
          </h1>
          <p className="text-xl text-base-content/70">{t('choice.subtitle')}</p>
        </motion.div>

        {/* Options */}
        <div className="grid lg:grid-cols-2 gap-8 max-w-6xl mx-auto">
          {/* Free Option */}
          <motion.div
            initial={{ x: -20, opacity: 0 }}
            animate={{ x: 0, opacity: 1 }}
            transition={{ delay: 0.1 }}
            onClick={() => handleOptionSelect('free')}
            className={`card cursor-pointer transition-all ${
              selectedOption === 'free'
                ? 'ring-4 ring-primary shadow-2xl'
                : 'hover:shadow-xl'
            } ${isLoading ? 'opacity-50 pointer-events-none' : ''}`}
          >
            <div className="card-body">
              {/* Badge */}
              <div className="flex justify-between items-start mb-4">
                <div className="badge badge-success badge-lg">
                  {t('choice.free_option.badge')}
                </div>
                {selectedOption === 'free' && (
                  <div className="badge badge-primary badge-lg">
                    {isLoading ? (
                      <span className="loading loading-spinner loading-xs"></span>
                    ) : (
                      t('choice.selected')
                    )}
                  </div>
                )}
              </div>

              {/* Title */}
              <div className="flex items-center gap-3 mb-4">
                <div className="p-3 bg-success/10 rounded-full">
                  <Package className="w-8 h-8 text-success" />
                </div>
                <div>
                  <h2 className="text-2xl font-bold">
                    {t('choice.free_option.title')}
                  </h2>
                  <p className="text-base-content/70">
                    {t('choice.free_option.subtitle')}
                  </p>
                </div>
              </div>

              {/* Stats */}
              <div className="grid grid-cols-3 gap-4 mb-6">
                <div className="text-center">
                  <div className="text-2xl font-bold">
                    {t('choice.free_option.time_value')}
                  </div>
                  <div className="text-sm text-base-content/60">
                    {t('choice.free_option.time')}
                  </div>
                </div>
                <div className="text-center">
                  <div className="text-2xl font-bold">
                    {t('choice.free_option.conversion_value')}
                  </div>
                  <div className="text-sm text-base-content/60">
                    {t('choice.free_option.conversion')}
                  </div>
                </div>
                <div className="text-center">
                  <div className="text-2xl font-bold">
                    {t('choice.free_option.cost_value')}
                  </div>
                  <div className="text-sm text-base-content/60">
                    {t('choice.free_option.cost')}
                  </div>
                </div>
              </div>

              {/* Features */}
              <div className="space-y-3">
                {freeFeatures.map((feature, index) => (
                  <motion.div
                    key={index}
                    initial={{ x: -20, opacity: 0 }}
                    animate={{ x: 0, opacity: 1 }}
                    transition={{ delay: 0.2 + index * 0.05 }}
                    className="flex items-center gap-3"
                  >
                    <feature.icon className="w-5 h-5 text-success" />
                    <span>{feature.text}</span>
                  </motion.div>
                ))}
              </div>

              {/* Description */}
              <div className="alert alert-info mt-6">
                <Info className="w-4 h-4" />
                <span className="text-sm">{t('choice.free_option.info')}</span>
              </div>
            </div>
          </motion.div>

          {/* AI Option */}
          <motion.div
            initial={{ x: 20, opacity: 0 }}
            animate={{ x: 0, opacity: 1 }}
            transition={{ delay: 0.1 }}
            onClick={() => handleOptionSelect('ai')}
            className={`card cursor-pointer transition-all ${
              selectedOption === 'ai'
                ? 'ring-4 ring-secondary shadow-2xl'
                : 'hover:shadow-xl'
            } ${isLoading ? 'opacity-50 pointer-events-none' : ''}`}
          >
            <div className="card-body">
              {/* Badge */}
              <div className="flex justify-between items-start mb-4">
                <div className="badge badge-secondary badge-lg">
                  <Zap className="w-3 h-3 mr-1" />
                  {t('choice.ai_option.badge')}
                </div>
                {selectedOption === 'ai' && (
                  <div className="badge badge-primary badge-lg">
                    {isLoading ? (
                      <span className="loading loading-spinner loading-xs"></span>
                    ) : (
                      t('choice.selected')
                    )}
                  </div>
                )}
              </div>

              {/* Title */}
              <div className="flex items-center gap-3 mb-4">
                <div className="p-3 bg-secondary/10 rounded-full">
                  <Brain className="w-8 h-8 text-secondary" />
                </div>
                <div>
                  <h2 className="text-2xl font-bold">
                    {t('choice.ai_option.title')}
                  </h2>
                  <p className="text-base-content/70">
                    {t('choice.ai_option.subtitle')}
                  </p>
                </div>
              </div>

              {/* Stats */}
              <div className="grid grid-cols-3 gap-4 mb-6">
                <div className="text-center">
                  <div className="text-2xl font-bold">
                    {t('choice.ai_option.time_value')}
                  </div>
                  <div className="text-sm text-base-content/60">
                    {t('choice.ai_option.time')}
                  </div>
                </div>
                <div className="text-center">
                  <div className="text-2xl font-bold">
                    {t('choice.ai_option.conversion_value')}
                  </div>
                  <div className="text-sm text-base-content/60">
                    {t('choice.ai_option.conversion')}
                  </div>
                </div>
                <div className="text-center">
                  <div className="text-2xl font-bold">
                    {t('choice.ai_option.cost_value')}
                  </div>
                  <div className="text-sm text-base-content/60">
                    {t('choice.ai_option.cost')}
                  </div>
                </div>
              </div>

              {/* Features */}
              <div className="space-y-3">
                {aiFeatures.map((feature, index) => (
                  <motion.div
                    key={index}
                    initial={{ x: 20, opacity: 0 }}
                    animate={{ x: 0, opacity: 1 }}
                    transition={{ delay: 0.2 + index * 0.05 }}
                    className="flex items-center gap-3"
                  >
                    <feature.icon className="w-5 h-5 text-secondary" />
                    <span>{feature.text}</span>
                  </motion.div>
                ))}
              </div>

              {/* Description */}
              <div className="alert alert-warning mt-6">
                <Sparkles className="w-4 h-4" />
                <span className="text-sm">{t('choice.ai_option.info')}</span>
              </div>
            </div>
          </motion.div>
        </div>

        {/* Comparison Table */}
        <motion.div
          initial={{ y: 20, opacity: 0 }}
          animate={{ y: 0, opacity: 1 }}
          transition={{ delay: 0.4 }}
          className="max-w-4xl mx-auto mt-12"
        >
          <h3 className="text-2xl font-bold text-center mb-6">
            {t('choice.comparison.title')}
          </h3>
          <div className="overflow-x-auto">
            <table className="table table-zebra">
              <thead>
                <tr>
                  <th>{t('choice.comparison.feature')}</th>
                  <th className="text-center">
                    {t('choice.comparison.free_column')}
                  </th>
                  <th className="text-center">
                    {t('choice.comparison.ai_column')}
                  </th>
                </tr>
              </thead>
              <tbody>
                <tr>
                  <td>{t('choice.comparison.creation_time')}</td>
                  <td className="text-center">
                    {t('choice.comparison.creation_time_free')}
                  </td>
                  <td className="text-center font-bold text-success">
                    {t('choice.comparison.creation_time_ai')}
                  </td>
                </tr>
                <tr>
                  <td>{t('choice.comparison.photo_recognition')}</td>
                  <td className="text-center">❌</td>
                  <td className="text-center">✅</td>
                </tr>
                <tr>
                  <td>{t('choice.comparison.auto_description')}</td>
                  <td className="text-center">
                    {t('choice.comparison.auto_description_free')}
                  </td>
                  <td className="text-center font-bold">
                    {t('choice.comparison.auto_description_ai')}
                  </td>
                </tr>
                <tr>
                  <td>{t('choice.comparison.price_analysis')}</td>
                  <td className="text-center">
                    {t('choice.comparison.price_analysis_free')}
                  </td>
                  <td className="text-center font-bold">
                    {t('choice.comparison.price_analysis_ai')}
                  </td>
                </tr>
                <tr>
                  <td>{t('choice.comparison.ab_testing')}</td>
                  <td className="text-center">❌</td>
                  <td className="text-center">✅</td>
                </tr>
                <tr>
                  <td>{t('choice.comparison.multilingual')}</td>
                  <td className="text-center">❌</td>
                  <td className="text-center">
                    {t('choice.comparison.multilingual_ai')}
                  </td>
                </tr>
                <tr>
                  <td>{t('choice.comparison.cost')}</td>
                  <td className="text-center font-bold text-success">
                    {t('choice.comparison.cost_free')}
                  </td>
                  <td className="text-center">
                    {t('choice.comparison.cost_ai')}
                  </td>
                </tr>
              </tbody>
            </table>
          </div>
        </motion.div>
      </div>
    </div>
  );
}
