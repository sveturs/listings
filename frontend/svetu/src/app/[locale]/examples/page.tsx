'use client';

import Link from 'next/link';
import { AnimatedSection } from '@/components/ui/AnimatedSection';

export default function ExamplesPage() {
  const examples = [
    {
      title: 'Главная v2.0 - Лучшие практики',
      description:
        'Новый дизайн с лучшими практиками Avito, Amazon и Wildberries',
      href: '/examples/ideal-homepage-v2',
      color: 'bg-gradient-to-r from-purple-600 to-pink-600',
      icon: '🎆',
      badge: 'NEW',
    },
    {
      title: 'Идеальная главная страница',
      description:
        'Современный дизайн главной страницы со всеми ключевыми элементами',
      href: '/examples/ideal-homepage',
      color: 'bg-gradient-to-r from-amber-500 to-orange-500',
      icon: '🏠',
      badge: 'HOT',
    },
    {
      title: 'Toast Notifications',
      description:
        'Interactive toast messages with different types and positions',
      href: '/examples/toast',
      color: 'bg-primary',
      icon: '🔔',
    },
    {
      title: 'Skeleton Loaders',
      description: 'Beautiful loading states for better UX',
      href: '/examples/skeletons',
      color: 'bg-secondary',
      icon: '⚡',
    },
    {
      title: 'Mobile Navigation',
      description: 'Responsive bottom navigation for mobile devices',
      href: '/examples/navigation',
      color: 'bg-accent',
      icon: '📱',
    },
    {
      title: 'Bento Grid Layout',
      description: 'Modern grid layout for showcasing content',
      href: '/examples/bento-grid',
      color: 'bg-info',
      icon: '🎨',
    },
    {
      title: 'Distance Visualization',
      description: 'Interactive distance display with visual indicators',
      href: '/examples/distance',
      color: 'bg-success',
      icon: '📍',
    },
    {
      title: 'Quick View Modal',
      description: 'Fast preview of listings without navigation',
      href: '/examples/quick-view',
      color: 'bg-warning',
      icon: '👁️',
    },
    {
      title: 'Page Transitions',
      description: 'Smooth animations between pages and sections',
      href: '/examples/transitions',
      color: 'bg-error',
      icon: '✨',
    },
    {
      title: 'Discount System',
      description:
        'Interactive discount badges and price history visualization',
      href: '/examples/discounts',
      color: 'bg-gradient-to-r from-red-500 to-orange-500',
      icon: '🏷️',
    },
    {
      title: 'Interactive Logos',
      description:
        '3D animated logos with particles, springs and morphing effects',
      href: '/examples/logos',
      color: 'bg-gradient-to-r from-purple-500 via-pink-500 to-cyan-500',
      icon: '🎭',
    },
    {
      title: 'Listing Creation UX',
      description:
        'Three revolutionary approaches to creating listings: from basic to AI-powered',
      href: '/examples/listing-creation-ux',
      color: 'bg-gradient-to-r from-green-500 to-teal-500',
      icon: '🚀',
    },
    {
      title: 'Listing Creation UX v2.0',
      description:
        'Enhanced examples with drag&drop, smart templates, A/B testing and more',
      href: '/examples/listing-creation-ux-v2',
      color: 'bg-gradient-to-r from-yellow-500 to-orange-500',
      icon: '✨',
      badge: 'NEW',
    },
    {
      title: 'Listing Edit UX',
      description:
        'Modern listing editing: from basic to AI-powered with real-time preview',
      href: '/examples/listing-edit-ux',
      color: 'bg-gradient-to-r from-blue-500 to-purple-500',
      icon: '✏️',
      badge: 'NEW',
    },
    {
      title: 'AI Создание объявлений',
      description:
        'Умный анализ фото с автоматической генерацией описания и цены',
      href: '/examples/ai-listing-creator',
      color: 'bg-gradient-to-r from-violet-500 to-purple-500',
      icon: '🤖',
    },
    {
      title: 'Умный поиск',
      description:
        'Продвинутая система поиска с фильтрами и OpenSearch интеграцией',
      href: '/examples/smart-search',
      color: 'bg-gradient-to-r from-blue-500 to-cyan-500',
      icon: '🔍',
    },
    {
      title: 'Карта с приватностью',
      description:
        'Настройки приватности местоположения для безопасности пользователей',
      href: '/examples/map-privacy',
      color: 'bg-gradient-to-r from-green-500 to-emerald-500',
      icon: '🗺️',
    },
    {
      title: 'Витрина B2C',
      description:
        'Dashboard для управления магазином с аналитикой и статистикой',
      href: '/examples/storefront-dashboard',
      color: 'bg-gradient-to-r from-orange-500 to-red-500',
      icon: '🏪',
    },
    {
      title: 'Чат с эмодзи',
      description: 'Анимированные эмодзи и реакции для живого общения',
      href: '/examples/animated-chat',
      color: 'bg-gradient-to-r from-pink-500 to-rose-500',
      icon: '💬',
    },
    {
      title: 'Эскроу платежи',
      description: 'Безопасные сделки с защитой средств до получения товара',
      href: '/examples/escrow-payment',
      color: 'bg-gradient-to-r from-indigo-500 to-purple-500',
      icon: '🔒',
    },
    {
      title: 'Адаптивный дизайн',
      description: 'Демонстрация responsive дизайна для всех устройств',
      href: '/examples/adaptive-design',
      color: 'bg-gradient-to-r from-teal-500 to-cyan-500',
      icon: '📱',
    },
    {
      title: 'Доставка BEX Express',
      description: 'Интеграция курьерской службы для C2C и B2C сценариев',
      href: '/examples/delivery',
      color: 'bg-gradient-to-r from-blue-600 to-indigo-600',
      icon: '🚚',
      badge: 'NEW',
    },
    {
      title: 'Идеальный маркетплейс',
      description: 'Главная страница с лучшими практиками Avito, Amazon и Wildberries',
      href: '/examples/ideal-marketplace',
      color: 'bg-gradient-to-r from-purple-600 via-pink-600 to-orange-600',
      icon: '🛒',
      badge: 'HOT',
    },
    {
      title: 'Детальная карточка товара',
      description: 'Страница товара с галереей, отзывами, характеристиками и Black Friday',
      href: '/examples/product-detail',
      color: 'bg-gradient-to-r from-cyan-500 via-blue-500 to-purple-500',
      icon: '📦',
      badge: 'NEW',
    },
  ];

  return (
    <div className="container mx-auto p-4 max-w-6xl">
      <AnimatedSection animation="fadeIn">
        <h1 className="text-4xl font-bold mb-4">UI/UX Examples</h1>
        <p className="text-lg text-base-content/70 mb-4">
          Explore all the UI/UX improvements implemented in the Sve Tu platform
        </p>
        <div className="stats shadow mb-8">
          <div className="stat">
            <div className="stat-title">Total Examples</div>
            <div className="stat-value text-primary">{examples.length}</div>
            <div className="stat-desc">UI/UX improvements</div>
          </div>
          <div className="stat">
            <div className="stat-title">Categories</div>
            <div className="stat-value text-secondary">8+</div>
            <div className="stat-desc">Different types</div>
          </div>
          <div className="stat">
            <div className="stat-title">New</div>
            <div className="stat-value text-accent">4</div>
            <div className="stat-desc">Latest additions</div>
          </div>
        </div>
      </AnimatedSection>

      <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-6">
        {examples.map((example, index) => (
          <AnimatedSection
            key={example.href}
            animation="slideUp"
            delay={index * 0.1}
          >
            <Link href={example.href} className="block h-full">
              <div className="card bg-base-100 shadow-xl hover:shadow-2xl transition-all duration-300 hover:-translate-y-1 h-full relative">
                {example.badge && (
                  <div className="absolute top-2 right-2 badge badge-warning badge-lg z-10">
                    {example.badge}
                  </div>
                )}
                <div className="card-body">
                  <div
                    className={`w-16 h-16 rounded-xl ${example.color} flex items-center justify-center text-3xl mb-4`}
                  >
                    {example.icon}
                  </div>
                  <h2 className="card-title">{example.title}</h2>
                  <p className="text-base-content/70">{example.description}</p>
                  <div className="card-actions justify-end mt-4">
                    <span className="btn btn-sm btn-ghost">View Example →</span>
                  </div>
                </div>
              </div>
            </Link>
          </AnimatedSection>
        ))}
      </div>

      <AnimatedSection animation="fadeIn" delay={0.8}>
        <div className="mt-12 p-6 bg-base-200 rounded-xl">
          <h2 className="text-2xl font-semibold mb-4">
            Implementation Progress
          </h2>
          <div className="space-y-2">
            <div className="flex justify-between items-center">
              <span>Overall Progress</span>
              <span className="font-bold">100%</span>
            </div>
            <progress
              className="progress progress-success w-full"
              value="100"
              max="100"
            ></progress>
            <p className="text-sm text-base-content/70 mt-2">
              All 25 UI/UX improvements have been successfully implemented!
            </p>
          </div>
        </div>
      </AnimatedSection>
    </div>
  );
}
