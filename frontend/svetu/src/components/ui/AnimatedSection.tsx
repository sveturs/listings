'use client';

import { motion } from 'framer-motion';
import { ReactNode } from 'react';

interface AnimatedSectionProps {
  children: ReactNode;
  delay?: number;
  duration?: number;
  className?: string;
  animation?: 'fadeIn' | 'slideUp' | 'slideIn' | 'zoomIn';
}

const animations = {
  fadeIn: {
    initial: { opacity: 0 },
    whileInView: { opacity: 1 },
    viewport: { once: true, amount: 0.3 },
  },
  slideUp: {
    initial: { y: 50, opacity: 0 },
    whileInView: { y: 0, opacity: 1 },
    viewport: { once: true, amount: 0.3 },
  },
  slideIn: {
    initial: { x: -50, opacity: 0 },
    whileInView: { x: 0, opacity: 1 },
    viewport: { once: true, amount: 0.3 },
  },
  zoomIn: {
    initial: { scale: 0.8, opacity: 0 },
    whileInView: { scale: 1, opacity: 1 },
    viewport: { once: true, amount: 0.3 },
  },
};

export function AnimatedSection({
  children,
  delay = 0,
  duration = 0.5,
  className = '',
  animation = 'fadeIn',
}: AnimatedSectionProps) {
  const selectedAnimation = animations[animation];

  return (
    <motion.div
      className={className}
      initial={selectedAnimation.initial}
      whileInView={selectedAnimation.whileInView}
      viewport={selectedAnimation.viewport}
      transition={{
        duration: duration,
        delay: delay,
        ease: [0.22, 1, 0.36, 1],
      }}
    >
      {children}
    </motion.div>
  );
}
