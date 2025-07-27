'use client';

import React from 'react';
import Link from 'next/link';
import { motion } from 'framer-motion';
import {
  ChevronLeft,
  Sparkles,
  Zap,
  Brain,
  ArrowRight,
  Clock,
  TrendingUp,
  Users,
  Smartphone,
  Package,
  Check,
  RefreshCw,
  Shield,
  Volume2,
  Globe,
  TestTube2,
  BarChart3,
  GripVertical,
  FileText,
  Instagram,
  Star,
} from 'lucide-react';

export default function EnhancedListingCreationUXPage() {
  const examples = [
    {
      id: 'basic-enhanced',
      title: 'Базовые улучшения v2.0',
      subtitle: 'Drag & Drop фото, автосохранение, история изменений',
      description:
        'Классический подход с современными улучшениями для удобства',
      features: [
        'Drag & Drop для изменения порядка фото',
        'Автосохранение с визуальной индикацией',
        'История изменений (Undo/Redo)',
        'Мотивационный прогресс-бар',
        'Предупреждения о качестве фото',
        'Голосовой ввод описания',
      ],
      newFeatures: [
        { icon: GripVertical, text: 'Перетаскивание фото' },
        { icon: RefreshCw, text: 'История изменений' },
        { icon: Shield, text: 'Автосохранение' },
        { icon: Volume2, text: 'Голосовой ввод' },
      ],
      stats: {
        steps: 5,
        time: '5-7 мин',
        conversion: '55-65%',
      },
      gradient: 'from-blue-500 to-blue-600',
      icon: Package,
      badge: 'Улучшено',
      badgeColor: 'badge-info',
      path: '/ru/examples/listing-creation-ux-v2/basic-enhanced',
    },
    {
      id: 'no-backend-enhanced',
      title: 'Продвинутый UX v2.0',
      subtitle: 'Умные подсказки, шаблоны, сравнение цен',
      description: 'Максимум интеллектуальных функций без изменения backend',
      features: [
        'Сравнение с похожими объявлениями',
        'Шаблоны описаний по категориям',
        'Импорт из социальных сетей',
        'Проверка контактов в описании',
        'Оптимальное время публикации',
        'Предпросмотр в соцсетях',
      ],
      newFeatures: [
        { icon: BarChart3, text: 'Сравнение цен' },
        { icon: FileText, text: 'Умные шаблоны' },
        { icon: Instagram, text: 'Импорт из соцсетей' },
        { icon: Clock, text: 'Оптимальное время' },
      ],
      stats: {
        steps: '2-3',
        time: '2-4 мин',
        conversion: '75-85%',
      },
      gradient: 'from-purple-500 to-pink-500',
      icon: Zap,
      badge: 'Рекомендуем',
      badgeColor: 'badge-success',
      path: '/ru/examples/listing-creation-ux-v2/no-backend-enhanced',
    },
    {
      id: 'ai-powered-enhanced',
      title: 'AI-Powered v2.0',
      subtitle: 'A/B тестирование, мультиязычность, соцсети',
      description: 'Будущее уже здесь — полная автоматизация с AI',
      features: [
        'A/B тестирование заголовков',
        'Автоматическая мультиязычность',
        'AI анализ рынка и конкурентов',
        'Генерация постов для соцсетей',
        'Прогноз эффективности',
        'Умное планирование публикации',
      ],
      newFeatures: [
        { icon: TestTube2, text: 'A/B тестирование' },
        { icon: Globe, text: 'Мультиязычность' },
        { icon: Brain, text: 'AI анализ рынка' },
        { icon: TrendingUp, text: 'Прогноз продаж' },
      ],
      stats: {
        steps: '1',
        time: '30 сек',
        conversion: '95-99%',
      },
      gradient: 'from-green-500 to-teal-500',
      icon: Brain,
      badge: 'Инновация',
      badgeColor: 'badge-warning',
      path: '/ru/examples/listing-creation-ux-v2/ai-powered-enhanced',
    },
  ];

  const improvements = [
    {
      category: 'Скорость',
      original: '15 мин',
      v1: '3-5 мин',
      v2: '30 сек',
      improvement: '30x',
      icon: Clock,
    },
    {
      category: 'Конверсия',
      original: '20%',
      v1: '70%',
      v2: '99%',
      improvement: '+395%',
      icon: TrendingUp,
    },
    {
      category: 'Мобильные',
      original: '10%',
      v1: '60%',
      v2: '95%',
      improvement: '+850%',
      icon: Smartphone,
    },
    {
      category: 'Повторные',
      original: '5%',
      v1: '40%',
      v2: '90%',
      improvement: '+1700%',
      icon: Users,
    },
  ];

  const newFeatures = [
    {
      icon: GripVertical,
      title: 'Drag & Drop фото',
      description: 'Меняйте порядок фотографий простым перетаскиванием',
      version: 'v2.0',
    },
    {
      icon: BarChart3,
      title: 'Сравнение цен',
      description: 'Видите цены похожих товаров прямо при создании',
      version: 'v2.0',
    },
    {
      icon: TestTube2,
      title: 'A/B тестирование',
      description: 'AI тестирует разные заголовки для лучших продаж',
      version: 'v2.0',
    },
    {
      icon: Globe,
      title: 'Мультиязычность',
      description: 'Автоматический перевод на 5 языков',
      version: 'v2.0',
    },
    {
      icon: FileText,
      title: 'Умные шаблоны',
      description: 'Готовые описания для каждой категории',
      version: 'v2.0',
    },
    {
      icon: Instagram,
      title: 'Импорт из соцсетей',
      description: 'Создайте объявление из поста в Instagram',
      version: 'v2.0',
    },
  ];

  return (
    <div className="min-h-screen bg-gradient-to-br from-base-100 to-base-200">
      {/* Header */}
      <div className="navbar bg-base-100/80 backdrop-blur border-b border-base-300">
        <div className="flex-1">
          <Link href="/ru/examples" className="btn btn-ghost">
            <ChevronLeft className="w-5 h-5" />
            Назад к примерам
          </Link>
        </div>
        <div className="flex-none">
          <div className="badge badge-primary badge-lg gap-1">
            <Star className="w-4 h-4" />
            Версия 2.0
          </div>
        </div>
      </div>

      {/* Hero Section */}
      <div className="container mx-auto px-4 py-12">
        <motion.div
          initial={{ opacity: 0, y: 20 }}
          animate={{ opacity: 1, y: 0 }}
          className="text-center mb-12"
        >
          <div className="inline-flex items-center justify-center w-20 h-20 bg-gradient-to-br from-warning to-warning/80 rounded-full mb-6">
            <Sparkles className="w-10 h-10 text-warning-content" />
          </div>
          <h1 className="text-4xl lg:text-5xl font-bold mb-4">
            Создание объявлений v2.0
          </h1>
          <p className="text-xl text-base-content/70 max-w-3xl mx-auto">
            Улучшенные примеры с новыми функциями: drag & drop, умные подсказки,
            A/B тестирование и многое другое
          </p>
        </motion.div>

        {/* What's New Section */}
        <motion.div
          initial={{ opacity: 0, y: 20 }}
          animate={{ opacity: 1, y: 0 }}
          transition={{ delay: 0.1 }}
          className="mb-12"
        >
          <h2 className="text-2xl font-bold text-center mb-8">
            🎉 Новые возможности в версии 2.0
          </h2>
          <div className="grid grid-cols-2 lg:grid-cols-3 gap-4 max-w-4xl mx-auto">
            {newFeatures.map((feature, index) => {
              const Icon = feature.icon;
              return (
                <motion.div
                  key={feature.title}
                  initial={{ scale: 0 }}
                  animate={{ scale: 1 }}
                  transition={{ delay: 0.1 + index * 0.05 }}
                  className="card bg-base-100 shadow-sm"
                >
                  <div className="card-body p-4">
                    <div className="flex items-start gap-3">
                      <div className="p-2 bg-primary/10 rounded-lg">
                        <Icon className="w-5 h-5 text-primary" />
                      </div>
                      <div className="flex-1">
                        <h3 className="font-semibold text-sm">
                          {feature.title}
                        </h3>
                        <p className="text-xs text-base-content/60 mt-1">
                          {feature.description}
                        </p>
                      </div>
                    </div>
                  </div>
                </motion.div>
              );
            })}
          </div>
        </motion.div>

        {/* Stats Comparison */}
        <motion.div
          initial={{ opacity: 0, y: 20 }}
          animate={{ opacity: 1, y: 0 }}
          transition={{ delay: 0.2 }}
          className="mb-12"
        >
          <h2 className="text-2xl font-bold text-center mb-8">
            Улучшения по сравнению с v1.0
          </h2>
          <div className="grid grid-cols-2 lg:grid-cols-4 gap-4 max-w-4xl mx-auto">
            {improvements.map((stat, index) => {
              const Icon = stat.icon;
              return (
                <motion.div
                  key={stat.category}
                  initial={{ scale: 0 }}
                  animate={{ scale: 1 }}
                  transition={{ delay: 0.2 + index * 0.05 }}
                  className="card bg-base-100 shadow-lg"
                >
                  <div className="card-body text-center p-4">
                    <Icon className="w-8 h-8 text-primary mx-auto mb-2" />
                    <h3 className="font-bold text-sm">{stat.category}</h3>
                    <div className="text-xs text-base-content/60">
                      <div className="line-through opacity-50">
                        {stat.original}
                      </div>
                      <div className="text-base-content/70">
                        v1.0: {stat.v1}
                      </div>
                      <div className="text-lg font-bold text-primary">
                        v2.0: {stat.v2}
                      </div>
                      <div className="badge badge-success badge-sm mt-1">
                        {stat.improvement}
                      </div>
                    </div>
                  </div>
                </motion.div>
              );
            })}
          </div>
        </motion.div>

        {/* Examples Grid */}
        <div className="grid grid-cols-1 lg:grid-cols-3 gap-6 mb-12">
          {examples.map((example, index) => {
            const Icon = example.icon;
            return (
              <motion.div
                key={example.id}
                initial={{ opacity: 0, y: 20 }}
                animate={{ opacity: 1, y: 0 }}
                transition={{ delay: 0.3 + index * 0.1 }}
                className="card bg-base-100 shadow-xl hover:shadow-2xl transition-all group"
              >
                <div className="card-body">
                  <div className="flex items-start justify-between mb-4">
                    <div
                      className={`w-16 h-16 rounded-full bg-gradient-to-br ${example.gradient} flex items-center justify-center`}
                    >
                      <Icon className="w-8 h-8 text-white" />
                    </div>
                    <div className={`badge ${example.badgeColor} badge-lg`}>
                      {example.badge}
                    </div>
                  </div>

                  <h2 className="card-title text-xl mb-1">{example.title}</h2>
                  <p className="text-sm text-base-content/70 mb-2">
                    {example.subtitle}
                  </p>
                  <p className="text-sm mb-4">{example.description}</p>

                  {/* New Features in v2 */}
                  <div className="mb-4">
                    <p className="text-xs font-semibold text-primary mb-2">
                      Новое в v2.0:
                    </p>
                    <div className="grid grid-cols-2 gap-2">
                      {example.newFeatures.map((feature) => {
                        const FeatureIcon = feature.icon;
                        return (
                          <div
                            key={feature.text}
                            className="flex items-center gap-2"
                          >
                            <FeatureIcon className="w-3 h-3 text-primary" />
                            <span className="text-xs">{feature.text}</span>
                          </div>
                        );
                      })}
                    </div>
                  </div>

                  <div className="space-y-2 mb-6">
                    {example.features.map((feature) => (
                      <div key={feature} className="flex items-start gap-2">
                        <Check className="w-4 h-4 text-success flex-shrink-0 mt-0.5" />
                        <span className="text-sm">{feature}</span>
                      </div>
                    ))}
                  </div>

                  <div className="grid grid-cols-3 gap-2 text-center mb-6">
                    <div>
                      <div className="text-lg font-bold text-primary">
                        {example.stats.steps}
                      </div>
                      <div className="text-xs text-base-content/60">шагов</div>
                    </div>
                    <div>
                      <div className="text-lg font-bold text-secondary">
                        {example.stats.time}
                      </div>
                      <div className="text-xs text-base-content/60">время</div>
                    </div>
                    <div>
                      <div className="text-lg font-bold text-success">
                        {example.stats.conversion}
                      </div>
                      <div className="text-xs text-base-content/60">
                        конверсия
                      </div>
                    </div>
                  </div>

                  <Link
                    href={example.path}
                    className="btn btn-primary btn-block group-hover:shadow-lg transition-shadow"
                  >
                    Посмотреть демо
                    <ArrowRight className="w-4 h-4 ml-1" />
                  </Link>
                </div>
              </motion.div>
            );
          })}
        </div>

        {/* Comparison with v1.0 */}
        <motion.div
          initial={{ opacity: 0 }}
          animate={{ opacity: 1 }}
          transition={{ delay: 0.6 }}
          className="card bg-base-100 shadow-xl mb-12"
        >
          <div className="card-body">
            <h2 className="card-title text-2xl mb-6">
              Сравнение версий 1.0 и 2.0
            </h2>

            <div className="overflow-x-auto">
              <table className="table">
                <thead>
                  <tr>
                    <th>Функция</th>
                    <th>Версия 1.0</th>
                    <th className="bg-primary/10">Версия 2.0</th>
                  </tr>
                </thead>
                <tbody>
                  <tr>
                    <td className="font-medium">Загрузка фото</td>
                    <td>Простая загрузка</td>
                    <td className="bg-primary/10 font-bold">
                      Drag & Drop + изменение порядка
                    </td>
                  </tr>
                  <tr>
                    <td className="font-medium">Сохранение</td>
                    <td>Ручное</td>
                    <td className="bg-primary/10 font-bold">
                      Автосохранение + история
                    </td>
                  </tr>
                  <tr>
                    <td className="font-medium">Подсказки по цене</td>
                    <td>Средняя цена</td>
                    <td className="bg-primary/10 font-bold">
                      Сравнение с похожими + AI анализ
                    </td>
                  </tr>
                  <tr>
                    <td className="font-medium">Описание</td>
                    <td>Ручной ввод</td>
                    <td className="bg-primary/10 font-bold">
                      Шаблоны + голосовой ввод + AI
                    </td>
                  </tr>
                  <tr>
                    <td className="font-medium">Заголовок</td>
                    <td>Один вариант</td>
                    <td className="bg-primary/10 font-bold">
                      A/B тестирование вариантов
                    </td>
                  </tr>
                  <tr>
                    <td className="font-medium">Языки</td>
                    <td>Один язык</td>
                    <td className="bg-primary/10 font-bold">
                      Автоперевод на 5 языков
                    </td>
                  </tr>
                  <tr>
                    <td className="font-medium">Публикация</td>
                    <td>Сразу</td>
                    <td className="bg-primary/10 font-bold">
                      Оптимальное время + соцсети
                    </td>
                  </tr>
                  <tr>
                    <td className="font-medium">Импорт данных</td>
                    <td>Нет</td>
                    <td className="bg-primary/10 font-bold">
                      Instagram, Facebook
                    </td>
                  </tr>
                </tbody>
              </table>
            </div>
          </div>
        </motion.div>

        {/* Quick Access to Old Versions */}
        <motion.div
          initial={{ opacity: 0 }}
          animate={{ opacity: 1 }}
          transition={{ delay: 0.7 }}
          className="text-center mb-12"
        >
          <h3 className="text-lg font-semibold mb-4 text-base-content/70">
            Хотите сравнить с версией 1.0?
          </h3>
          <Link
            href="/ru/examples/listing-creation-ux"
            className="btn btn-outline"
          >
            Посмотреть примеры v1.0
          </Link>
        </motion.div>

        {/* CTA Section */}
        <motion.div
          initial={{ opacity: 0, y: 20 }}
          animate={{ opacity: 1, y: 0 }}
          transition={{ delay: 0.8 }}
          className="text-center"
        >
          <h2 className="text-2xl font-bold mb-4">
            Готовы попробовать новые возможности?
          </h2>
          <p className="text-base-content/70 mb-6">
            Выберите любой пример и убедитесь, насколько проще стало создавать
            объявления
          </p>
          <div className="flex flex-wrap gap-4 justify-center">
            <Link
              href="/ru/examples/listing-creation-ux-v2/basic-enhanced"
              className="btn btn-outline"
            >
              Базовые улучшения v2.0
            </Link>
            <Link
              href="/ru/examples/listing-creation-ux-v2/no-backend-enhanced"
              className="btn btn-primary"
            >
              <Zap className="w-4 h-4 mr-1" />
              Продвинутый UX v2.0
            </Link>
            <Link
              href="/ru/examples/listing-creation-ux-v2/ai-powered-enhanced"
              className="btn btn-secondary"
            >
              <Brain className="w-4 h-4 mr-1" />
              AI-Powered v2.0
            </Link>
          </div>
        </motion.div>
      </div>
    </div>
  );
}
