'use client';

import { useState } from 'react';
import { PageTransition } from '@/components/ui/PageTransition';
import { AnimatedSection } from '@/components/ui/AnimatedSection';
import { withPageTransition } from '@/components/ui/withPageTransition';
import { motion } from 'framer-motion';
import Link from 'next/link';

// Пример компонента с HOC
const ExampleCard = ({
  title,
  content,
}: {
  title: string;
  content: string;
}) => (
  <div className="card bg-base-200">
    <div className="card-body">
      <h3 className="card-title">{title}</h3>
      <p>{content}</p>
    </div>
  </div>
);

const _AnimatedCard = withPageTransition(ExampleCard, {
  mode: 'scale',
  duration: 0.5,
});

export default function TransitionsPage() {
  const [currentTransition, setCurrentTransition] = useState<
    'fade' | 'slide' | 'scale' | 'slideUp'
  >('fade');
  const [showContent, setShowContent] = useState(true);

  const toggleContent = () => {
    setShowContent(false);
    setTimeout(() => setShowContent(true), 100);
  };

  return (
    <div className="container mx-auto px-4 py-8">
      <PageTransition mode="slideUp">
        <h1 className="text-4xl font-bold mb-8">Анимации переходов</h1>
      </PageTransition>

      {/* Примеры PageTransition */}
      <section className="mb-12">
        <AnimatedSection animation="fadeIn">
          <h2 className="text-2xl font-semibold mb-6">Page Transitions</h2>
        </AnimatedSection>

        <AnimatedSection animation="slideUp" delay={0.1}>
          <p className="text-base-content/70 mb-6">
            Компонент PageTransition автоматически анимирует переходы между
            страницами
          </p>
        </AnimatedSection>

        <AnimatedSection animation="slideUp" delay={0.2}>
          <div className="flex gap-4 mb-6 flex-wrap">
            <button
              onClick={() => {
                setCurrentTransition('fade');
                toggleContent();
              }}
              className={`btn ${currentTransition === 'fade' ? 'btn-primary' : 'btn-outline'}`}
            >
              Fade
            </button>
            <button
              onClick={() => {
                setCurrentTransition('slide');
                toggleContent();
              }}
              className={`btn ${currentTransition === 'slide' ? 'btn-primary' : 'btn-outline'}`}
            >
              Slide
            </button>
            <button
              onClick={() => {
                setCurrentTransition('scale');
                toggleContent();
              }}
              className={`btn ${currentTransition === 'scale' ? 'btn-primary' : 'btn-outline'}`}
            >
              Scale
            </button>
            <button
              onClick={() => {
                setCurrentTransition('slideUp');
                toggleContent();
              }}
              className={`btn ${currentTransition === 'slideUp' ? 'btn-primary' : 'btn-outline'}`}
            >
              Slide Up
            </button>
          </div>
        </AnimatedSection>

        {showContent && (
          <PageTransition mode={currentTransition} duration={0.4}>
            <div className="mockup-window bg-base-300 mb-8">
              <div className="px-6 py-8">
                <h3 className="text-xl font-semibold mb-4">
                  Контент с анимацией: {currentTransition}
                </h3>
                <p className="text-base-content/70">
                  Этот контент анимируется при изменении типа перехода.
                  PageTransition автоматически применяет выбранную анимацию.
                </p>
              </div>
            </div>
          </PageTransition>
        )}
      </section>

      {/* Примеры AnimatedSection */}
      <section className="mb-12">
        <AnimatedSection animation="fadeIn">
          <h2 className="text-2xl font-semibold mb-6">Animated Sections</h2>
        </AnimatedSection>

        <AnimatedSection animation="slideUp" delay={0.1}>
          <p className="text-base-content/70 mb-6">
            Секции с анимацией при прокрутке страницы (viewport animations)
          </p>
        </AnimatedSection>

        <div className="grid md:grid-cols-2 gap-6">
          <AnimatedSection animation="slideIn" delay={0.2}>
            <div className="card bg-primary text-primary-content">
              <div className="card-body">
                <h3 className="card-title">Slide In Animation</h3>
                <p>Эта карточка появляется слева при прокрутке</p>
              </div>
            </div>
          </AnimatedSection>

          <AnimatedSection animation="zoomIn" delay={0.3}>
            <div className="card bg-secondary text-secondary-content">
              <div className="card-body">
                <h3 className="card-title">Zoom In Animation</h3>
                <p>Эта карточка увеличивается при появлении</p>
              </div>
            </div>
          </AnimatedSection>

          <AnimatedSection animation="fadeIn" delay={0.4}>
            <div className="card bg-accent text-accent-content">
              <div className="card-body">
                <h3 className="card-title">Fade In Animation</h3>
                <p>Эта карточка плавно появляется</p>
              </div>
            </div>
          </AnimatedSection>

          <AnimatedSection animation="slideUp" delay={0.5}>
            <div className="card bg-info text-info-content">
              <div className="card-body">
                <h3 className="card-title">Slide Up Animation</h3>
                <p>Эта карточка поднимается снизу</p>
              </div>
            </div>
          </AnimatedSection>
        </div>
      </section>

      {/* Примеры микро-анимаций */}
      <section className="mb-12">
        <AnimatedSection animation="fadeIn">
          <h2 className="text-2xl font-semibold mb-6">Микро-анимации</h2>
        </AnimatedSection>

        <AnimatedSection animation="slideUp" delay={0.1}>
          <p className="text-base-content/70 mb-6">
            Мелкие анимации для улучшения UX
          </p>
        </AnimatedSection>

        <div className="flex flex-wrap gap-4">
          <motion.button
            className="btn btn-primary"
            whileHover={{ scale: 1.05 }}
            whileTap={{ scale: 0.95 }}
          >
            Hover & Tap
          </motion.button>

          <motion.button
            className="btn btn-secondary"
            whileHover={{ rotate: 5 }}
            whileTap={{ rotate: -5 }}
          >
            Rotate on Hover
          </motion.button>

          <motion.div
            className="btn btn-accent"
            animate={{
              scale: [1, 1.1, 1],
              transition: { repeat: Infinity, duration: 2 },
            }}
          >
            Pulse Animation
          </motion.div>

          <motion.button
            className="btn btn-info"
            whileHover={{
              backgroundColor: '#3b82f6',
              color: '#ffffff',
              transition: { duration: 0.3 },
            }}
          >
            Color Change
          </motion.button>
        </div>
      </section>

      {/* Примеры кода */}
      <section className="mb-12">
        <AnimatedSection animation="fadeIn">
          <h2 className="text-2xl font-semibold mb-6">Примеры использования</h2>
        </AnimatedSection>

        <div className="space-y-6">
          <AnimatedSection animation="slideUp" delay={0.1}>
            <div className="mockup-code">
              <pre data-prefix="1">
                <code>{`// Page Transition`}</code>
              </pre>
              <pre data-prefix="2">
                <code>{`import { PageTransition } from '@/components/ui/PageTransition';`}</code>
              </pre>
              <pre data-prefix="3">
                <code>{``}</code>
              </pre>
              <pre data-prefix="4">
                <code>{`<PageTransition mode="slideUp">`}</code>
              </pre>
              <pre data-prefix="5">
                <code>{`  <YourPageContent />`}</code>
              </pre>
              <pre data-prefix="6">
                <code>{`</PageTransition>`}</code>
              </pre>
            </div>
          </AnimatedSection>

          <AnimatedSection animation="slideUp" delay={0.2}>
            <div className="mockup-code">
              <pre data-prefix="1">
                <code>{`// Animated Section`}</code>
              </pre>
              <pre data-prefix="2">
                <code>{`import { AnimatedSection } from '@/components/ui/AnimatedSection';`}</code>
              </pre>
              <pre data-prefix="3">
                <code>{``}</code>
              </pre>
              <pre data-prefix="4">
                <code>{`<AnimatedSection animation="slideUp" delay={0.2}>`}</code>
              </pre>
              <pre data-prefix="5">
                <code>{`  <Card />`}</code>
              </pre>
              <pre data-prefix="6">
                <code>{`</AnimatedSection>`}</code>
              </pre>
            </div>
          </AnimatedSection>

          <AnimatedSection animation="slideUp" delay={0.3}>
            <div className="mockup-code">
              <pre data-prefix="1">
                <code>{`// With HOC`}</code>
              </pre>
              <pre data-prefix="2">
                <code>{`import { withPageTransition } from '@/components/ui/withPageTransition';`}</code>
              </pre>
              <pre data-prefix="3">
                <code>{``}</code>
              </pre>
              <pre data-prefix="4">
                <code>{`const AnimatedComponent = withPageTransition(YourComponent, {`}</code>
              </pre>
              <pre data-prefix="5">
                <code>{`  mode: 'scale',`}</code>
              </pre>
              <pre data-prefix="6">
                <code>{`  duration: 0.5`}</code>
              </pre>
              <pre data-prefix="7">
                <code>{`});`}</code>
              </pre>
            </div>
          </AnimatedSection>
        </div>
      </section>

      {/* Навигация */}
      <AnimatedSection animation="fadeIn" delay={0.4}>
        <div className="flex gap-4 mt-8">
          <Link href="/examples" className="btn btn-outline">
            ← Назад к примерам
          </Link>
          <Link href="/examples/toast" className="btn btn-primary">
            Toast уведомления →
          </Link>
        </div>
      </AnimatedSection>
    </div>
  );
}
