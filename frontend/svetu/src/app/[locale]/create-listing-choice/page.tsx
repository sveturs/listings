'use client';

import React, { useState } from 'react';
import { useRouter } from 'next/navigation';
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
  const t = useTranslations();
  const { user } = useAuthContext();
  const [selectedOption, setSelectedOption] = useState<'free' | 'ai' | null>(
    null
  );

  const handleContinue = () => {
    if (!user) {
      toast.error(t('create_listing.auth_required'));
      router.push('/');
      return;
    }

    if (!selectedOption) {
      toast.error('Пожалуйста, выберите способ создания объявления');
      return;
    }

    if (selectedOption === 'free') {
      router.push('/sr/create-listing-smart');
    } else {
      router.push('/sr/create-listing-ai');
    }
  };

  const freeFeatures = [
    { icon: FileText, text: 'Умные шаблоны описаний' },
    { icon: ImageIcon, text: 'Drag & Drop для фото' },
    { icon: TrendingUp, text: 'Сравнение цен с похожими' },
    { icon: Clock, text: 'Оптимальное время публикации' },
    { icon: MapPin, text: 'Импорт из соцсетей' },
    { icon: Shield, text: 'Проверка на ошибки' },
  ];

  const aiFeatures = [
    { icon: Brain, text: 'AI распознает товар по фото' },
    { icon: Sparkles, text: 'Автогенерация описания' },
    { icon: DollarSign, text: 'AI анализ оптимальной цены' },
    { icon: FileText, text: 'A/B тестирование заголовков' },
    { icon: Check, text: 'Мультиязычность (5 языков)' },
    { icon: TrendingUp, text: 'Прогноз эффективности' },
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
            Выберите способ создания объявления
          </h1>
          <p className="text-xl text-base-content/70">
            Два варианта — выберите подходящий для вас
          </p>
        </motion.div>

        {/* Options */}
        <div className="grid lg:grid-cols-2 gap-8 max-w-6xl mx-auto">
          {/* Free Option */}
          <motion.div
            initial={{ x: -20, opacity: 0 }}
            animate={{ x: 0, opacity: 1 }}
            transition={{ delay: 0.1 }}
            onClick={() => setSelectedOption('free')}
            className={`card cursor-pointer transition-all ${
              selectedOption === 'free'
                ? 'ring-4 ring-primary shadow-2xl'
                : 'hover:shadow-xl'
            }`}
          >
            <div className="card-body">
              {/* Badge */}
              <div className="flex justify-between items-start mb-4">
                <div className="badge badge-success badge-lg">Бесплатно</div>
                {selectedOption === 'free' && (
                  <div className="badge badge-primary badge-lg">Выбрано</div>
                )}
              </div>

              {/* Title */}
              <div className="flex items-center gap-3 mb-4">
                <div className="p-3 bg-success/10 rounded-full">
                  <Package className="w-8 h-8 text-success" />
                </div>
                <div>
                  <h2 className="text-2xl font-bold">Умный помощник</h2>
                  <p className="text-base-content/70">
                    Без искусственного интеллекта
                  </p>
                </div>
              </div>

              {/* Stats */}
              <div className="grid grid-cols-3 gap-4 mb-6">
                <div className="text-center">
                  <div className="text-2xl font-bold">2-4 мин</div>
                  <div className="text-sm text-base-content/60">время</div>
                </div>
                <div className="text-center">
                  <div className="text-2xl font-bold">85%</div>
                  <div className="text-sm text-base-content/60">конверсия</div>
                </div>
                <div className="text-center">
                  <div className="text-2xl font-bold">0 РСД</div>
                  <div className="text-sm text-base-content/60">стоимость</div>
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
                <span className="text-sm">
                  Подходит для опытных продавцов, которые хотят полный контроль
                </span>
              </div>
            </div>
          </motion.div>

          {/* AI Option */}
          <motion.div
            initial={{ x: 20, opacity: 0 }}
            animate={{ x: 0, opacity: 1 }}
            transition={{ delay: 0.1 }}
            onClick={() => setSelectedOption('ai')}
            className={`card cursor-pointer transition-all ${
              selectedOption === 'ai'
                ? 'ring-4 ring-secondary shadow-2xl'
                : 'hover:shadow-xl'
            }`}
          >
            <div className="card-body">
              {/* Badge */}
              <div className="flex justify-between items-start mb-4">
                <div className="badge badge-secondary badge-lg">
                  <Zap className="w-3 h-3 mr-1" />
                  Premium
                </div>
                {selectedOption === 'ai' && (
                  <div className="badge badge-primary badge-lg">Выбрано</div>
                )}
              </div>

              {/* Title */}
              <div className="flex items-center gap-3 mb-4">
                <div className="p-3 bg-secondary/10 rounded-full">
                  <Brain className="w-8 h-8 text-secondary" />
                </div>
                <div>
                  <h2 className="text-2xl font-bold">AI-Powered</h2>
                  <p className="text-base-content/70">
                    Искусственный интеллект
                  </p>
                </div>
              </div>

              {/* Stats */}
              <div className="grid grid-cols-3 gap-4 mb-6">
                <div className="text-center">
                  <div className="text-2xl font-bold">30 сек</div>
                  <div className="text-sm text-base-content/60">время</div>
                </div>
                <div className="text-center">
                  <div className="text-2xl font-bold">99%</div>
                  <div className="text-sm text-base-content/60">конверсия</div>
                </div>
                <div className="text-center">
                  <div className="text-2xl font-bold">199 РСД</div>
                  <div className="text-sm text-base-content/60">
                    за объявление
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
                <span className="text-sm">
                  Идеально для быстрых продаж и максимального охвата
                </span>
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
            Сравнение возможностей
          </h3>
          <div className="overflow-x-auto">
            <table className="table table-zebra">
              <thead>
                <tr>
                  <th>Возможность</th>
                  <th className="text-center">Умный помощник</th>
                  <th className="text-center">AI-Powered</th>
                </tr>
              </thead>
              <tbody>
                <tr>
                  <td>Время создания</td>
                  <td className="text-center">2-4 минуты</td>
                  <td className="text-center font-bold text-success">
                    30 секунд
                  </td>
                </tr>
                <tr>
                  <td>Распознавание товара по фото</td>
                  <td className="text-center">❌</td>
                  <td className="text-center">✅</td>
                </tr>
                <tr>
                  <td>Автогенерация описания</td>
                  <td className="text-center">Шаблоны</td>
                  <td className="text-center font-bold">AI генерация</td>
                </tr>
                <tr>
                  <td>Анализ цены</td>
                  <td className="text-center">Сравнение с похожими</td>
                  <td className="text-center font-bold">AI анализ рынка</td>
                </tr>
                <tr>
                  <td>A/B тестирование</td>
                  <td className="text-center">❌</td>
                  <td className="text-center">✅</td>
                </tr>
                <tr>
                  <td>Мультиязычность</td>
                  <td className="text-center">❌</td>
                  <td className="text-center">5 языков</td>
                </tr>
                <tr>
                  <td>Стоимость</td>
                  <td className="text-center font-bold text-success">
                    Бесплатно
                  </td>
                  <td className="text-center">199 РСД</td>
                </tr>
              </tbody>
            </table>
          </div>
        </motion.div>

        {/* Action Button */}
        <motion.div
          initial={{ y: 20, opacity: 0 }}
          animate={{ y: 0, opacity: 1 }}
          transition={{ delay: 0.5 }}
          className="text-center mt-12"
        >
          <button
            onClick={handleContinue}
            disabled={!selectedOption}
            className="btn btn-primary btn-lg gap-2"
          >
            Продолжить
            <ChevronRight className="w-5 h-5" />
          </button>
          {!selectedOption && (
            <p className="text-sm text-base-content/60 mt-2">
              Выберите один из вариантов для продолжения
            </p>
          )}
        </motion.div>
      </div>
    </div>
  );
}
