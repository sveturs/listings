'use client';

import Link from 'next/link';
import { AnimatedSection } from '@/components/ui/AnimatedSection';

export default function ExamplesPage() {
  const examples = [
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
  ];

  return (
    <div className="container mx-auto p-4 max-w-6xl">
      <AnimatedSection animation="fadeIn">
        <h1 className="text-4xl font-bold mb-4">UI/UX Examples</h1>
        <p className="text-lg text-base-content/70 mb-8">
          Explore all the UI/UX improvements implemented in the Sve Tu platform
        </p>
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
              All 15 UI/UX improvements have been successfully implemented!
            </p>
          </div>
        </div>
      </AnimatedSection>
    </div>
  );
}
