'use client';

import React, { useState, useEffect, useCallback } from 'react';
import { EnhancedMobileBottomNav } from './EnhancedMobileBottomNav';

export const SmartMobileBottomNav: React.FC = () => {
  const [isVisible, setIsVisible] = useState(true);
  const [lastScrollY, setLastScrollY] = useState(0);
  const [scrollDirection, setScrollDirection] = useState<'up' | 'down'>('up');

  const handleScroll = useCallback(() => {
    const currentScrollY = window.scrollY;
    
    // Определяем направление скролла
    if (currentScrollY > lastScrollY && currentScrollY > 100) {
      setScrollDirection('down');
      setIsVisible(false);
    } else {
      setScrollDirection('up');
      setIsVisible(true);
    }
    
    setLastScrollY(currentScrollY);
  }, [lastScrollY]);

  useEffect(() => {
    window.addEventListener('scroll', handleScroll, { passive: true });
    
    return () => {
      window.removeEventListener('scroll', handleScroll);
    };
  }, [handleScroll]);

  return (
    <div
      className={`
        fixed bottom-0 left-0 right-0 z-50
        transition-transform duration-300 ease-in-out
        ${isVisible ? 'translate-y-0' : 'translate-y-full'}
      `}
    >
      <EnhancedMobileBottomNav />
    </div>
  );
};