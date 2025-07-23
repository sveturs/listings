'use client';

import React, { useState } from 'react';
import { BentoGrid } from '@/components/ui/BentoGrid';
import { DynamicBentoGrid, BentoGridItem } from '@/components/ui/DynamicBentoGrid';
import { 
  Home, 
  Search, 
  ShoppingBag, 
  Users, 
  TrendingUp,
  Zap,
  Globe,
  Shield,
  Heart,
  MessageSquare,
  Camera,
  Palette
} from 'lucide-react';

export default function BentoGridExamplesPage() {
  const [variant, setVariant] = useState<'default' | 'compact' | 'hero'>('default');

  // Примеры данных для BentoGrid
  const categories = [
    { id: '1', name: 'Электроника', count: 342 },
    { id: '2', name: 'Одежда', count: 567 },
    { id: '3', name: 'Дом и сад', count: 234 },
    { id: '4', name: 'Авто', count: 123 },
  ];

  const featuredListing = {
    id: '123',
    title: 'MacBook Pro 16" M2 Max',
    price: '€2,500',
    image: 'https://via.placeholder.com/400x300',
    category: 'Электроника',
  };

  const stats = {
    totalListings: 12543,
    activeUsers: 3421,
    successfulDeals: 9876,
  };

  // Динамические элементы для DynamicBentoGrid
  const dynamicItems: BentoGridItem[] = [
    {
      id: '1',
      title: 'Добро пожаловать!',
      description: 'Исследуйте мир возможностей с нашей платформой',
      icon: <Home className="w-8 h-8 text-primary" />,
      colSpan: 2,
      rowSpan: 2,
      bgColor: 'bg-gradient-to-br from-primary/20 to-secondary/20',
      href: '/',
    },
    {
      id: '2',
      title: 'Поиск товаров',
      description: 'Найдите именно то, что ищете',
      icon: <Search className="w-6 h-6 text-info" />,
      colSpan: 1,
      bgColor: 'bg-info/10',
      href: '/search',
    },
    {
      id: '3',
      title: 'Активные пользователи',
      content: (
        <div className="text-center">
          <p className="text-3xl font-bold text-success">1,234</p>
          <p className="text-sm text-base-content/60">онлайн сейчас</p>
        </div>
      ),
      icon: <Users className="w-6 h-6 text-success" />,
      bgColor: 'bg-success/10',
    },
    {
      id: '4',
      title: 'Тренды недели',
      icon: <TrendingUp className="w-6 h-6 text-warning" />,
      content: (
        <div className="space-y-2">
          <div className="flex items-center justify-between">
            <span className="text-sm">iPhone 15</span>
            <span className="text-xs text-warning">↑ 23%</span>
          </div>
          <div className="flex items-center justify-between">
            <span className="text-sm">PS5</span>
            <span className="text-xs text-success">↑ 15%</span>
          </div>
        </div>
      ),
      bgColor: 'bg-warning/10',
    },
    {
      id: '5',
      title: 'Быстрые действия',
      icon: <Zap className="w-6 h-6 text-accent" />,
      colSpan: 1,
      content: (
        <div className="flex flex-col gap-2">
          <button className="btn btn-sm btn-accent">Создать</button>
          <button className="btn btn-sm btn-ghost">Обзор</button>
        </div>
      ),
      bgColor: 'bg-accent/10',
    },
    {
      id: '6',
      title: 'Глобальный охват',
      description: 'Доступно в 50+ странах',
      icon: <Globe className="w-6 h-6 text-secondary" />,
      bgColor: 'bg-secondary/10',
    },
    {
      id: '7',
      title: 'Безопасность',
      icon: <Shield className="w-6 h-6 text-success" />,
      content: (
        <div>
          <div className="rating rating-sm">
            <input type="radio" name="rating-1" className="mask mask-star-2 bg-success" disabled checked />
            <input type="radio" name="rating-1" className="mask mask-star-2 bg-success" disabled checked />
            <input type="radio" name="rating-1" className="mask mask-star-2 bg-success" disabled checked />
            <input type="radio" name="rating-1" className="mask mask-star-2 bg-success" disabled checked />
            <input type="radio" name="rating-1" className="mask mask-star-2 bg-success" disabled checked />
          </div>
          <p className="text-xs mt-1">Защита 24/7</p>
        </div>
      ),
      bgColor: 'bg-base-200',
    },
    {
      id: '8',
      title: 'Сообщество',
      description: 'Присоединяйтесь к тысячам пользователей',
      icon: <Heart className="w-6 h-6 text-error" />,
      bgColor: 'bg-error/10',
    },
  ];

  // Компактные элементы
  const compactItems: BentoGridItem[] = [
    { id: '1', icon: <Camera className="w-8 h-8" />, bgColor: 'bg-primary/20' },
    { id: '2', icon: <ShoppingBag className="w-8 h-8" />, bgColor: 'bg-secondary/20' },
    { id: '3', icon: <MessageSquare className="w-8 h-8" />, bgColor: 'bg-accent/20' },
    { id: '4', icon: <Heart className="w-8 h-8" />, bgColor: 'bg-error/20' },
    { id: '5', icon: <Globe className="w-8 h-8" />, bgColor: 'bg-warning/20' },
    { id: '6', icon: <Palette className="w-8 h-8" />, bgColor: 'bg-info/20' },
  ];

  return (
    <div className="min-h-screen py-8">
      <div className="container mx-auto px-4">
        <h1 className="text-3xl font-bold mb-8">Bento Grid Examples</h1>

        {/* Статический BentoGrid */}
        <section className="mb-12">
          <h2 className="text-2xl font-semibold mb-6">1. Статический Bento Grid</h2>
          <p className="text-base-content/70 mb-6">
            Готовый компонент для главной страницы с предустановленными карточками
          </p>
          <BentoGrid
            categories={categories}
            featuredListing={featuredListing}
            stats={stats}
          />
        </section>

        {/* Динамический BentoGrid */}
        <section className="mb-12">
          <h2 className="text-2xl font-semibold mb-6">2. Динамический Bento Grid</h2>
          <p className="text-base-content/70 mb-6">
            Гибкий компонент для создания кастомных макетов
          </p>
          
          {/* Переключатель вариантов */}
          <div className="mb-6 flex justify-center">
            <div className="join">
              <button
                className={`join-item btn ${variant === 'default' ? 'btn-active' : ''}`}
                onClick={() => setVariant('default')}
              >
                Default
              </button>
              <button
                className={`join-item btn ${variant === 'compact' ? 'btn-active' : ''}`}
                onClick={() => setVariant('compact')}
              >
                Compact
              </button>
              <button
                className={`join-item btn ${variant === 'hero' ? 'btn-active' : ''}`}
                onClick={() => setVariant('hero')}
              >
                Hero
              </button>
            </div>
          </div>

          <DynamicBentoGrid
            items={variant === 'compact' ? compactItems : dynamicItems}
            variant={variant}
          />
        </section>

        {/* Особенности */}
        <section className="mb-12">
          <h2 className="text-2xl font-semibold mb-6">Особенности Bento Grid</h2>
          <div className="grid grid-cols-1 md:grid-cols-2 gap-6">
            <div className="card bg-base-200 p-6">
              <h3 className="text-lg font-semibold mb-3">🎨 Визуальная иерархия</h3>
              <ul className="list-disc list-inside space-y-2 text-sm">
                <li>Разные размеры карточек для приоритизации контента</li>
                <li>Цветовое кодирование для различных типов информации</li>
                <li>Градиенты и тени для глубины</li>
              </ul>
            </div>
            
            <div className="card bg-base-200 p-6">
              <h3 className="text-lg font-semibold mb-3">🚀 Интерактивность</h3>
              <ul className="list-disc list-inside space-y-2 text-sm">
                <li>Hover эффекты с масштабированием</li>
                <li>Каскадная анимация появления</li>
                <li>Кликабельные карточки с навигацией</li>
              </ul>
            </div>
            
            <div className="card bg-base-200 p-6">
              <h3 className="text-lg font-semibold mb-3">📱 Адаптивность</h3>
              <ul className="list-disc list-inside space-y-2 text-sm">
                <li>Респонсивная сетка для всех устройств</li>
                <li>Автоматическая перекомпоновка элементов</li>
                <li>Оптимизация для мобильных экранов</li>
              </ul>
            </div>
            
            <div className="card bg-base-200 p-6">
              <h3 className="text-lg font-semibold mb-3">⚡ Производительность</h3>
              <ul className="list-disc list-inside space-y-2 text-sm">
                <li>CSS анимации для плавности</li>
                <li>Ленивая загрузка контента</li>
                <li>Оптимизированный рендеринг</li>
              </ul>
            </div>
          </div>
        </section>

        {/* Примеры кода */}
        <section className="card bg-base-200 p-6">
          <h2 className="text-2xl font-semibold mb-4">Примеры использования</h2>
          <div className="mockup-code">
            <pre data-prefix="1"><code>{`// Статический BentoGrid`}</code></pre>
            <pre data-prefix="2"><code>{`import { BentoGrid } from '@/components/ui/BentoGrid';`}</code></pre>
            <pre data-prefix="3"><code>{``}</code></pre>
            <pre data-prefix="4"><code>{`<BentoGrid`}</code></pre>
            <pre data-prefix="5"><code>{`  categories={categories}`}</code></pre>
            <pre data-prefix="6"><code>{`  featuredListing={featuredListing}`}</code></pre>
            <pre data-prefix="7"><code>{`  stats={stats}`}</code></pre>
            <pre data-prefix="8"><code>{`/>`}</code></pre>
            <pre data-prefix="9"><code>{``}</code></pre>
            <pre data-prefix="10"><code>{`// Динамический BentoGrid`}</code></pre>
            <pre data-prefix="11"><code>{`import { DynamicBentoGrid } from '@/components/ui/DynamicBentoGrid';`}</code></pre>
            <pre data-prefix="12"><code>{``}</code></pre>
            <pre data-prefix="13"><code>{`const items = [`}</code></pre>
            <pre data-prefix="14"><code>{`  { id: '1', title: 'Card 1', colSpan: 2, rowSpan: 2 },`}</code></pre>
            <pre data-prefix="15"><code>{`  { id: '2', title: 'Card 2', bgColor: 'bg-primary/20' },`}</code></pre>
            <pre data-prefix="16"><code>{`];`}</code></pre>
            <pre data-prefix="17"><code>{``}</code></pre>
            <pre data-prefix="18"><code>{`<DynamicBentoGrid items={items} variant="default" />`}</code></pre>
          </div>
        </section>
      </div>
    </div>
  );
}