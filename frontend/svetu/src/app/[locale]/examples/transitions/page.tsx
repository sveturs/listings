'use client';

import { useState } from 'react';
import { motion, AnimatePresence } from 'framer-motion';
import { AnimatedSection } from '@/components/ui/AnimatedSection';
import Link from 'next/link';

export default function TransitionsExamplePage() {
  const [currentTab, setCurrentTab] = useState(0);
  const [showContent, setShowContent] = useState(true);

  const tabs = ['Overview', 'Features', 'Examples', 'Code'];

  const tabContent = [
    {
      title: 'Overview',
      content:
        'Smooth page transitions enhance user experience by providing visual continuity between different views.',
    },
    {
      title: 'Features',
      content:
        'Multiple animation modes, customizable duration, easy integration with Next.js App Router.',
    },
    {
      title: 'Examples',
      content:
        'See various transition effects in action below. Each animation can be customized to match your design.',
    },
    {
      title: 'Code',
      content:
        'Simple implementation using HOC pattern or direct component usage.',
    },
  ];

  return (
    <div className="container mx-auto p-4 max-w-6xl">
      <AnimatedSection animation="fadeIn">
        <h1 className="text-4xl font-bold mb-8">
          Page Transitions & Animations
        </h1>
      </AnimatedSection>

      {/* Tab Navigation Example */}
      <AnimatedSection animation="slideUp" delay={0.1}>
        <section className="card bg-base-100 shadow-xl mb-8">
          <div className="card-body">
            <h2 className="card-title mb-4">Tab Transitions</h2>

            <div className="tabs tabs-boxed mb-4">
              {tabs.map((tab, index) => (
                <button
                  key={tab}
                  className={`tab ${currentTab === index ? 'tab-active' : ''}`}
                  onClick={() => setCurrentTab(index)}
                >
                  {tab}
                </button>
              ))}
            </div>

            <AnimatePresence mode="wait">
              <motion.div
                key={currentTab}
                initial={{ opacity: 0, y: 10 }}
                animate={{ opacity: 1, y: 0 }}
                exit={{ opacity: 0, y: -10 }}
                transition={{ duration: 0.2 }}
                className="p-4 bg-base-200 rounded-lg"
              >
                <h3 className="text-lg font-semibold mb-2">
                  {tabContent[currentTab].title}
                </h3>
                <p>{tabContent[currentTab].content}</p>
              </motion.div>
            </AnimatePresence>
          </div>
        </section>
      </AnimatedSection>

      {/* Animation Types */}
      <AnimatedSection animation="slideUp" delay={0.2}>
        <section className="card bg-base-100 shadow-xl mb-8">
          <div className="card-body">
            <h2 className="card-title mb-4">Animation Types</h2>

            <div className="grid grid-cols-1 md:grid-cols-2 gap-4">
              <AnimatedSection animation="fadeIn" delay={0.3}>
                <div className="p-4 bg-primary/10 rounded-lg">
                  <h3 className="font-semibold mb-2">Fade Animation</h3>
                  <p className="text-sm">Simple opacity transition</p>
                </div>
              </AnimatedSection>

              <AnimatedSection animation="slideIn" delay={0.4}>
                <div className="p-4 bg-secondary/10 rounded-lg">
                  <h3 className="font-semibold mb-2">Slide Animation</h3>
                  <p className="text-sm">Horizontal slide with fade</p>
                </div>
              </AnimatedSection>

              <AnimatedSection animation="slideUp" delay={0.5}>
                <div className="p-4 bg-accent/10 rounded-lg">
                  <h3 className="font-semibold mb-2">Slide Up Animation</h3>
                  <p className="text-sm">Vertical slide from bottom</p>
                </div>
              </AnimatedSection>

              <AnimatedSection animation="zoomIn" delay={0.6}>
                <div className="p-4 bg-info/10 rounded-lg">
                  <h3 className="font-semibold mb-2">Zoom Animation</h3>
                  <p className="text-sm">Scale with opacity effect</p>
                </div>
              </AnimatedSection>
            </div>
          </div>
        </section>
      </AnimatedSection>

      {/* Interactive Demo */}
      <AnimatedSection animation="slideUp" delay={0.3}>
        <section className="card bg-base-100 shadow-xl mb-8">
          <div className="card-body">
            <h2 className="card-title mb-4">Interactive Demo</h2>

            <button
              className="btn btn-primary mb-4"
              onClick={() => setShowContent(!showContent)}
            >
              Toggle Content
            </button>

            <AnimatePresence>
              {showContent && (
                <motion.div
                  initial={{ height: 0, opacity: 0 }}
                  animate={{ height: 'auto', opacity: 1 }}
                  exit={{ height: 0, opacity: 0 }}
                  transition={{ duration: 0.3 }}
                  className="overflow-hidden"
                >
                  <div className="p-4 bg-base-200 rounded-lg">
                    <h3 className="font-semibold mb-2">Animated Content</h3>
                    <p>
                      This content smoothly animates in and out when toggled.
                    </p>
                    <div className="mt-4 space-y-2">
                      <div className="skeleton h-4 w-full"></div>
                      <div className="skeleton h-4 w-3/4"></div>
                      <div className="skeleton h-4 w-1/2"></div>
                    </div>
                  </div>
                </motion.div>
              )}
            </AnimatePresence>
          </div>
        </section>
      </AnimatedSection>

      {/* Stagger Animation */}
      <AnimatedSection animation="slideUp" delay={0.4}>
        <section className="card bg-base-100 shadow-xl mb-8">
          <div className="card-body">
            <h2 className="card-title mb-4">Stagger Animation</h2>
            <p className="mb-4">Elements animate in sequence</p>

            <div className="space-y-2">
              {[1, 2, 3, 4, 5].map((item, index) => (
                <AnimatedSection
                  key={item}
                  animation="slideIn"
                  delay={0.1 * index}
                  className="p-3 bg-base-200 rounded"
                >
                  <div className="flex items-center gap-3">
                    <div className="w-12 h-12 bg-primary rounded-full flex items-center justify-center text-white font-bold">
                      {item}
                    </div>
                    <div>
                      <h4 className="font-semibold">Item {item}</h4>
                      <p className="text-sm text-base-content/70">
                        This item animates with a {0.1 * index}s delay
                      </p>
                    </div>
                  </div>
                </AnimatedSection>
              ))}
            </div>
          </div>
        </section>
      </AnimatedSection>

      {/* Code Examples */}
      <AnimatedSection animation="slideUp" delay={0.5}>
        <section className="card bg-base-100 shadow-xl">
          <div className="card-body">
            <h2 className="card-title mb-4">Implementation Examples</h2>

            <div className="space-y-4">
              <div>
                <h3 className="font-semibold mb-2">
                  Using PageTransition Component
                </h3>
                <div className="mockup-code">
                  <pre data-prefix="1">
                    <code>{`import { PageTransition } from '@/components/ui/PageTransition';`}</code>
                  </pre>
                  <pre data-prefix="2">
                    <code>{``}</code>
                  </pre>
                  <pre data-prefix="3">
                    <code>{`export default function MyPage() {`}</code>
                  </pre>
                  <pre data-prefix="4">
                    <code>{`  return (`}</code>
                  </pre>
                  <pre data-prefix="5">
                    <code>{`    <PageTransition mode="slide" duration={0.3}>`}</code>
                  </pre>
                  <pre data-prefix="6">
                    <code>{`      <div>Your page content</div>`}</code>
                  </pre>
                  <pre data-prefix="7">
                    <code>{`    </PageTransition>`}</code>
                  </pre>
                  <pre data-prefix="8">
                    <code>{`  );`}</code>
                  </pre>
                  <pre data-prefix="9">
                    <code>{`}`}</code>
                  </pre>
                </div>
              </div>

              <div>
                <h3 className="font-semibold mb-2">Using HOC Pattern</h3>
                <div className="mockup-code">
                  <pre data-prefix="1">
                    <code>{`import { withPageTransition } from '@/components/ui/withPageTransition';`}</code>
                  </pre>
                  <pre data-prefix="2">
                    <code>{``}</code>
                  </pre>
                  <pre data-prefix="3">
                    <code>{`function MyPage() {`}</code>
                  </pre>
                  <pre data-prefix="4">
                    <code>{`  return <div>Your page content</div>;`}</code>
                  </pre>
                  <pre data-prefix="5">
                    <code>{`}`}</code>
                  </pre>
                  <pre data-prefix="6">
                    <code>{``}</code>
                  </pre>
                  <pre data-prefix="7">
                    <code>{`export default withPageTransition(MyPage, {`}</code>
                  </pre>
                  <pre data-prefix="8">
                    <code>{`  mode: 'fade',`}</code>
                  </pre>
                  <pre data-prefix="9">
                    <code>{`  duration: 0.5`}</code>
                  </pre>
                  <pre data-prefix="10">
                    <code>{`});`}</code>
                  </pre>
                </div>
              </div>

              <div>
                <h3 className="font-semibold mb-2">Using AnimatedSection</h3>
                <div className="mockup-code">
                  <pre data-prefix="1">
                    <code>{`import { AnimatedSection } from '@/components/ui/AnimatedSection';`}</code>
                  </pre>
                  <pre data-prefix="2">
                    <code>{``}</code>
                  </pre>
                  <pre data-prefix="3">
                    <code>{`<AnimatedSection animation="slideUp" delay={0.2}>`}</code>
                  </pre>
                  <pre data-prefix="4">
                    <code>{`  <h2>This will slide up when scrolled into view</h2>`}</code>
                  </pre>
                  <pre data-prefix="5">
                    <code>{`</AnimatedSection>`}</code>
                  </pre>
                </div>
              </div>
            </div>
          </div>
        </section>
      </AnimatedSection>

      {/* Navigation Links */}
      <AnimatedSection animation="fadeIn" delay={0.6}>
        <div className="mt-8 flex flex-wrap gap-2">
          <Link href="/examples" className="btn btn-ghost">
            ← Back to Examples
          </Link>
          <Link href="/examples/toast" className="btn btn-ghost">
            Toast Examples →
          </Link>
        </div>
      </AnimatedSection>
    </div>
  );
}
