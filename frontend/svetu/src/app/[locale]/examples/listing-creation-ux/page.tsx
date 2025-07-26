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
  Camera,
  Package,
} from 'lucide-react';

export default function ListingCreationUXPage() {
  const examples = [
    {
      id: 'basic',
      title: 'Базовые улучшения UX',
      subtitle: 'Без AI, только оптимизация интерфейса',
      description:
        'Упрощенный процесс из 5 шагов с улучшенной мобильной версией',
      features: [
        'Сокращение с 8 до 5 шагов',
        'Популярные категории на главной',
        'Улучшенная мобильная навигация',
        'Контекстные подсказки',
        'Автосохранение черновиков',
      ],
      stats: {
        steps: 5,
        time: '7-10 мин',
        conversion: '40-50%',
      },
      gradient: 'from-blue-500 to-blue-600',
      icon: Package,
      path: '/ru/examples/listing-creation-ux/basic',
    },
    {
      id: 'no-backend',
      title: 'Продвинутый UX без изменений Backend',
      subtitle: 'Максимум возможностей на текущей архитектуре',
      description:
        'Революционный интерфейс с быстрым стартом и умными шаблонами',
      features: [
        'Старт с фото или шаблона',
        'Быстрый режим (3 шага)',
        'Прогрессивное раскрытие полей',
        'Визуальный предпросмотр',
        'Анимации и микровзаимодействия',
      ],
      stats: {
        steps: '3-4',
        time: '3-5 мин',
        conversion: '60-70%',
      },
      gradient: 'from-purple-500 to-pink-500',
      icon: Zap,
      path: '/ru/examples/listing-creation-ux/no-backend',
    },
    {
      id: 'ai-powered',
      title: 'AI-Powered создание объявлений',
      subtitle: 'Полная интеграция искусственного интеллекта',
      description: 'Футуристичный опыт с AI-ассистентом на каждом шаге',
      features: [
        'Распознавание товара по фото',
        'Автогенерация описания',
        'Умный подбор цены',
        'SEO-оптимизация',
        'Голосовой ввод',
        'Предсказание эффективности',
      ],
      stats: {
        steps: '1-2',
        time: '30 сек - 2 мин',
        conversion: '85-95%',
      },
      gradient: 'from-green-500 to-teal-500',
      icon: Brain,
      path: '/ru/examples/listing-creation-ux/ai-powered',
    },
  ];

  const improvements = [
    {
      category: 'Скорость',
      current: '15 минут',
      improved: '30 секунд',
      improvement: '30x быстрее',
      icon: Clock,
    },
    {
      category: 'Конверсия',
      current: '20%',
      improved: '95%',
      improvement: '+375%',
      icon: TrendingUp,
    },
    {
      category: 'Мобильные',
      current: '10%',
      improved: '85%',
      improvement: '+750%',
      icon: Smartphone,
    },
    {
      category: 'Повторные',
      current: '5%',
      improved: '80%',
      improvement: '+1500%',
      icon: Users,
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
      </div>

      {/* Hero Section */}
      <div className="container mx-auto px-4 py-12">
        <motion.div
          initial={{ opacity: 0, y: 20 }}
          animate={{ opacity: 1, y: 0 }}
          className="text-center mb-12"
        >
          <h1 className="text-4xl lg:text-5xl font-bold mb-4">
            Эволюция создания объявлений
          </h1>
          <p className="text-xl text-base-content/70 max-w-3xl mx-auto">
            От классического подхода до AI-powered решения. Посмотрите, как
            можно трансформировать процесс создания объявлений для максимальной
            конверсии.
          </p>
        </motion.div>

        {/* Stats Comparison */}
        <motion.div
          initial={{ opacity: 0, y: 20 }}
          animate={{ opacity: 1, y: 0 }}
          transition={{ delay: 0.1 }}
          className="grid grid-cols-2 lg:grid-cols-4 gap-4 mb-12 max-w-4xl mx-auto"
        >
          {improvements.map((stat, index) => {
            const Icon = stat.icon;
            return (
              <motion.div
                key={stat.category}
                initial={{ scale: 0 }}
                animate={{ scale: 1 }}
                transition={{ delay: 0.1 + index * 0.05 }}
                className="card bg-base-100 shadow-lg"
              >
                <div className="card-body text-center p-4">
                  <Icon className="w-8 h-8 text-primary mx-auto mb-2" />
                  <h3 className="font-bold text-sm">{stat.category}</h3>
                  <div className="text-xs text-base-content/60">
                    <div className="line-through">{stat.current}</div>
                    <div className="text-lg font-bold text-primary">
                      {stat.improved}
                    </div>
                    <div className="badge badge-success badge-sm">
                      {stat.improvement}
                    </div>
                  </div>
                </div>
              </motion.div>
            );
          })}
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
                transition={{ delay: 0.2 + index * 0.1 }}
                className="card bg-base-100 shadow-xl hover:shadow-2xl transition-all"
              >
                <div className="card-body">
                  <div
                    className={`w-16 h-16 rounded-full bg-gradient-to-br ${example.gradient} flex items-center justify-center mb-4`}
                  >
                    <Icon className="w-8 h-8 text-white" />
                  </div>

                  <h2 className="card-title text-xl mb-1">{example.title}</h2>
                  <p className="text-sm text-base-content/70 mb-2">
                    {example.subtitle}
                  </p>
                  <p className="text-sm mb-4">{example.description}</p>

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
                    className="btn btn-primary btn-block"
                  >
                    Посмотреть демо
                    <ArrowRight className="w-4 h-4 ml-1" />
                  </Link>
                </div>
              </motion.div>
            );
          })}
        </div>

        {/* Comparison Table */}
        <motion.div
          initial={{ opacity: 0 }}
          animate={{ opacity: 1 }}
          transition={{ delay: 0.5 }}
          className="card bg-base-100 shadow-xl"
        >
          <div className="card-body">
            <h2 className="card-title text-2xl mb-6">
              Детальное сравнение подходов
            </h2>

            <div className="overflow-x-auto">
              <table className="table">
                <thead>
                  <tr>
                    <th>Характеристика</th>
                    <th>Текущая версия</th>
                    <th>Базовые улучшения</th>
                    <th>Без изменений Backend</th>
                    <th className="bg-primary/10">AI-Powered</th>
                  </tr>
                </thead>
                <tbody>
                  <tr>
                    <td className="font-medium">Количество шагов</td>
                    <td>8</td>
                    <td>5</td>
                    <td>3-4</td>
                    <td className="bg-primary/10 font-bold">1-2</td>
                  </tr>
                  <tr>
                    <td className="font-medium">Время создания</td>
                    <td>15 мин</td>
                    <td>7-10 мин</td>
                    <td>3-5 мин</td>
                    <td className="bg-primary/10 font-bold">30 сек - 2 мин</td>
                  </tr>
                  <tr>
                    <td className="font-medium">Конверсия</td>
                    <td>20%</td>
                    <td>40-50%</td>
                    <td>60-70%</td>
                    <td className="bg-primary/10 font-bold">85-95%</td>
                  </tr>
                  <tr>
                    <td className="font-medium">Мобильная конверсия</td>
                    <td>10%</td>
                    <td>30-40%</td>
                    <td>50-60%</td>
                    <td className="bg-primary/10 font-bold">80-90%</td>
                  </tr>
                  <tr>
                    <td className="font-medium">Старт с фото</td>
                    <td>❌</td>
                    <td>❌</td>
                    <td>✅</td>
                    <td className="bg-primary/10">✅</td>
                  </tr>
                  <tr>
                    <td className="font-medium">Автозаполнение</td>
                    <td>❌</td>
                    <td>Частично</td>
                    <td>Шаблоны</td>
                    <td className="bg-primary/10 font-bold">AI генерация</td>
                  </tr>
                  <tr>
                    <td className="font-medium">Умная цена</td>
                    <td>❌</td>
                    <td>Средняя цена</td>
                    <td>Средняя цена</td>
                    <td className="bg-primary/10 font-bold">AI анализ рынка</td>
                  </tr>
                  <tr>
                    <td className="font-medium">SEO оптимизация</td>
                    <td>❌</td>
                    <td>❌</td>
                    <td>Базовая</td>
                    <td className="bg-primary/10 font-bold">AI оптимизация</td>
                  </tr>
                </tbody>
              </table>
            </div>
          </div>
        </motion.div>

        {/* Features Showcase */}
        <motion.div
          initial={{ opacity: 0 }}
          animate={{ opacity: 1 }}
          transition={{ delay: 0.6 }}
          className="mt-12"
        >
          <h2 className="text-3xl font-bold text-center mb-8">
            Ключевые инновации
          </h2>

          <div className="grid grid-cols-1 lg:grid-cols-2 gap-6">
            <div className="card bg-gradient-to-br from-primary/10 to-primary/5 border-2 border-primary/20">
              <div className="card-body">
                <Camera className="w-12 h-12 text-primary mb-4" />
                <h3 className="card-title">Старт с фото</h3>
                <p className="text-base-content/70">
                  Пользователь просто загружает фото, а система сама определяет
                  категорию, генерирует название и описание. Это снижает барьер
                  входа и ускоряет процесс.
                </p>
              </div>
            </div>

            <div className="card bg-gradient-to-br from-secondary/10 to-secondary/5 border-2 border-secondary/20">
              <div className="card-body">
                <Brain className="w-12 h-12 text-secondary mb-4" />
                <h3 className="card-title">AI-ассистент</h3>
                <p className="text-base-content/70">
                  Искусственный интеллект анализирует фото, предлагает
                  оптимальную цену на основе рынка и генерирует
                  SEO-оптимизированное описание.
                </p>
              </div>
            </div>

            <div className="card bg-gradient-to-br from-success/10 to-success/5 border-2 border-success/20">
              <div className="card-body">
                <Smartphone className="w-12 h-12 text-success mb-4" />
                <h3 className="card-title">Mobile-first дизайн</h3>
                <p className="text-base-content/70">
                  Интерфейс оптимизирован для мобильных устройств с крупными
                  элементами, жестами и минимальным вводом текста.
                </p>
              </div>
            </div>

            <div className="card bg-gradient-to-br from-warning/10 to-warning/5 border-2 border-warning/20">
              <div className="card-body">
                <Zap className="w-12 h-12 text-warning mb-4" />
                <h3 className="card-title">Быстрые шаблоны</h3>
                <p className="text-base-content/70">
                  Готовые шаблоны для популярных категорий позволяют создать
                  объявление буквально в несколько кликов.
                </p>
              </div>
            </div>
          </div>
        </motion.div>

        {/* CTA Section */}
        <motion.div
          initial={{ opacity: 0, y: 20 }}
          animate={{ opacity: 1, y: 0 }}
          transition={{ delay: 0.7 }}
          className="text-center mt-12"
        >
          <h2 className="text-2xl font-bold mb-4">Готовы увидеть будущее?</h2>
          <p className="text-base-content/70 mb-6">
            Выберите любой пример выше и посмотрите, как может выглядеть
            создание объявлений
          </p>
          <div className="flex flex-wrap gap-4 justify-center">
            <Link
              href="/ru/examples/listing-creation-ux/basic"
              className="btn btn-outline"
            >
              Базовые улучшения
            </Link>
            <Link
              href="/ru/examples/listing-creation-ux/no-backend"
              className="btn btn-outline"
            >
              Без изменений Backend
            </Link>
            <Link
              href="/ru/examples/listing-creation-ux/ai-powered"
              className="btn btn-primary"
            >
              <Sparkles className="w-4 h-4 mr-1" />
              AI-Powered версия
            </Link>
          </div>
        </motion.div>
      </div>
    </div>
  );
}

function Check({ className }: { className?: string }) {
  return (
    <svg
      className={className}
      fill="none"
      viewBox="0 0 24 24"
      stroke="currentColor"
    >
      <path
        strokeLinecap="round"
        strokeLinejoin="round"
        strokeWidth={2}
        d="M5 13l4 4L19 7"
      />
    </svg>
  );
}
