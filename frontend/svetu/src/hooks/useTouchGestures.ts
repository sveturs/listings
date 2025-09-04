// Touch Gestures Hook для мобильной оптимизации
// День 25: Мобильная оптимизация системы атрибутов

import { useEffect, useRef, useState, useCallback } from 'react';

interface TouchPoint {
  x: number;
  y: number;
  timestamp: number;
}

interface SwipeResult {
  direction: 'left' | 'right' | 'up' | 'down' | null;
  velocity: number;
  distance: number;
}

interface PinchResult {
  scale: number;
  center: { x: number; y: number };
}

interface TouchGesturesOptions {
  onSwipe?: (result: SwipeResult) => void;
  onPinch?: (result: PinchResult) => void;
  onTap?: (point: TouchPoint) => void;
  onDoubleTap?: (point: TouchPoint) => void;
  onLongPress?: (point: TouchPoint) => void;
  swipeThreshold?: number;
  swipeVelocityThreshold?: number;
  tapTimeout?: number;
  longPressTimeout?: number;
  preventScroll?: boolean;
}

export const useTouchGestures = (
  elementRef: React.RefObject<HTMLElement>,
  options: TouchGesturesOptions = {}
) => {
  const {
    onSwipe,
    onPinch,
    onTap,
    onDoubleTap,
    onLongPress,
    swipeThreshold = 50,
    swipeVelocityThreshold = 0.3,
    tapTimeout = 300,
    longPressTimeout = 500,
    preventScroll = false,
  } = options;

  const [isGesturing, setIsGesturing] = useState(false);
  const [currentGesture, setCurrentGesture] = useState<string | null>(null);

  const touchStartRef = useRef<TouchPoint | null>(null);
  const touchEndRef = useRef<TouchPoint | null>(null);
  const lastTapRef = useRef<number>(0);
  const longPressTimerRef = useRef<NodeJS.Timeout | null>(null);
  const pinchStartDistanceRef = useRef<number>(0);
  const rafRef = useRef<number | null>(null);

  // Вычисление расстояния между двумя точками
  const getDistance = useCallback((p1: TouchPoint, p2: TouchPoint): number => {
    const dx = p2.x - p1.x;
    const dy = p2.y - p1.y;
    return Math.sqrt(dx * dx + dy * dy);
  }, []);

  // Вычисление расстояния между двумя пальцами (pinch)
  const getPinchDistance = useCallback((touches: TouchList): number => {
    if (touches.length < 2) return 0;
    const dx = touches[0].clientX - touches[1].clientX;
    const dy = touches[0].clientY - touches[1].clientY;
    return Math.sqrt(dx * dx + dy * dy);
  }, []);

  // Определение направления свайпа
  const getSwipeDirection = useCallback(
    (start: TouchPoint, end: TouchPoint): SwipeResult['direction'] => {
      const dx = end.x - start.x;
      const dy = end.y - start.y;
      const absDx = Math.abs(dx);
      const absDy = Math.abs(dy);

      if (absDx < swipeThreshold && absDy < swipeThreshold) {
        return null;
      }

      if (absDx > absDy) {
        return dx > 0 ? 'right' : 'left';
      } else {
        return dy > 0 ? 'down' : 'up';
      }
    },
    [swipeThreshold]
  );

  // Обработчик начала касания
  const handleTouchStart = useCallback(
    (e: TouchEvent) => {
      if (preventScroll) {
        e.preventDefault();
      }

      const touch = e.touches[0];
      const point: TouchPoint = {
        x: touch.clientX,
        y: touch.clientY,
        timestamp: Date.now(),
      };

      touchStartRef.current = point;
      setIsGesturing(true);

      // Проверка на double tap
      const timeSinceLastTap = Date.now() - lastTapRef.current;
      if (timeSinceLastTap < tapTimeout) {
        onDoubleTap?.(point);
        setCurrentGesture('double-tap');
        lastTapRef.current = 0;
      } else {
        lastTapRef.current = Date.now();
      }

      // Long press detection
      if (onLongPress) {
        longPressTimerRef.current = setTimeout(() => {
          onLongPress(point);
          setCurrentGesture('long-press');
        }, longPressTimeout);
      }

      // Pinch gesture начало
      if (e.touches.length === 2) {
        pinchStartDistanceRef.current = getPinchDistance(e.touches);
        setCurrentGesture('pinch');
      }
    },
    [
      preventScroll,
      tapTimeout,
      longPressTimeout,
      onDoubleTap,
      onLongPress,
      getPinchDistance,
    ]
  );

  // Обработчик движения пальца
  const handleTouchMove = useCallback(
    (e: TouchEvent) => {
      if (!touchStartRef.current) return;

      // Отмена long press при движении
      if (longPressTimerRef.current) {
        clearTimeout(longPressTimerRef.current);
        longPressTimerRef.current = null;
      }

      // Pinch gesture
      if (e.touches.length === 2 && pinchStartDistanceRef.current > 0) {
        if (rafRef.current) cancelAnimationFrame(rafRef.current);

        rafRef.current = requestAnimationFrame(() => {
          const currentDistance = getPinchDistance(e.touches);
          const scale = currentDistance / pinchStartDistanceRef.current;
          const center = {
            x: (e.touches[0].clientX + e.touches[1].clientX) / 2,
            y: (e.touches[0].clientY + e.touches[1].clientY) / 2,
          };

          onPinch?.({ scale, center });
        });
      }

      // Swipe tracking
      if (e.touches.length === 1) {
        const touch = e.touches[0];
        touchEndRef.current = {
          x: touch.clientX,
          y: touch.clientY,
          timestamp: Date.now(),
        };
        setCurrentGesture('swipe');
      }
    },
    [onPinch, getPinchDistance]
  );

  // Обработчик окончания касания
  const handleTouchEnd = useCallback(
    (e: TouchEvent) => {
      setIsGesturing(false);
      setCurrentGesture(null);

      // Отмена long press
      if (longPressTimerRef.current) {
        clearTimeout(longPressTimerRef.current);
        longPressTimerRef.current = null;
      }

      // Отмена RAF
      if (rafRef.current) {
        cancelAnimationFrame(rafRef.current);
        rafRef.current = null;
      }

      if (!touchStartRef.current) return;

      const endPoint = touchEndRef.current || {
        x: e.changedTouches[0].clientX,
        y: e.changedTouches[0].clientY,
        timestamp: Date.now(),
      };

      const distance = getDistance(touchStartRef.current, endPoint);
      const timeDiff = endPoint.timestamp - touchStartRef.current.timestamp;
      const velocity = timeDiff > 0 ? distance / timeDiff : 0;

      // Определение типа жеста
      if (distance < swipeThreshold) {
        // Tap gesture
        if (timeDiff < tapTimeout && onTap) {
          onTap(endPoint);
        }
      } else if (
        velocity > swipeVelocityThreshold ||
        distance > swipeThreshold
      ) {
        // Swipe gesture
        const direction = getSwipeDirection(touchStartRef.current, endPoint);
        if (direction && onSwipe) {
          onSwipe({ direction, velocity, distance });
        }
      }

      // Reset
      touchStartRef.current = null;
      touchEndRef.current = null;
      pinchStartDistanceRef.current = 0;
    },
    [
      getDistance,
      swipeThreshold,
      swipeVelocityThreshold,
      tapTimeout,
      getSwipeDirection,
      onTap,
      onSwipe,
    ]
  );

  // Установка обработчиков событий
  useEffect(() => {
    const element = elementRef.current;
    if (!element) return;

    // Passive listeners для лучшей производительности
    const options = { passive: !preventScroll };

    element.addEventListener('touchstart', handleTouchStart, options);
    element.addEventListener('touchmove', handleTouchMove, options);
    element.addEventListener('touchend', handleTouchEnd, options);
    element.addEventListener('touchcancel', handleTouchEnd, options);

    return () => {
      element.removeEventListener('touchstart', handleTouchStart);
      element.removeEventListener('touchmove', handleTouchMove);
      element.removeEventListener('touchend', handleTouchEnd);
      element.removeEventListener('touchcancel', handleTouchEnd);

      if (longPressTimerRef.current) {
        clearTimeout(longPressTimerRef.current);
      }
      if (rafRef.current) {
        cancelAnimationFrame(rafRef.current);
      }
    };
  }, [
    elementRef,
    handleTouchStart,
    handleTouchMove,
    handleTouchEnd,
    preventScroll,
  ]);

  return {
    isGesturing,
    currentGesture,
  };
};

// Хук для свайпа между страницами/секциями
export const useSwipeNavigation = (
  elementRef: React.RefObject<HTMLElement>,
  onNavigate: (direction: 'prev' | 'next') => void,
  options: { threshold?: number; enabled?: boolean } = {}
) => {
  const { threshold = 50, enabled = true } = options;

  useTouchGestures(elementRef, {
    onSwipe: ({ direction, distance }) => {
      if (!enabled) return;

      if (distance > threshold) {
        if (direction === 'left') {
          onNavigate('next');
        } else if (direction === 'right') {
          onNavigate('prev');
        }
      }
    },
    swipeThreshold: threshold,
  });
};

// Хук для pull-to-refresh
export const usePullToRefresh = (
  elementRef: React.RefObject<HTMLElement>,
  onRefresh: () => Promise<void>,
  options: { threshold?: number; enabled?: boolean } = {}
) => {
  const { threshold = 100, enabled = true } = options;
  const [isPulling, setIsPulling] = useState(false);
  const [pullDistance, setPullDistance] = useState(0);
  const startYRef = useRef<number>(0);

  useEffect(() => {
    const element = elementRef.current;
    if (!element || !enabled) return;

    let isRefreshing = false;

    const handleTouchStart = (e: TouchEvent) => {
      if (element.scrollTop === 0) {
        startYRef.current = e.touches[0].clientY;
        setIsPulling(true);
      }
    };

    const handleTouchMove = (e: TouchEvent) => {
      if (!isPulling || isRefreshing) return;

      const currentY = e.touches[0].clientY;
      const distance = Math.max(0, currentY - startYRef.current);

      setPullDistance(Math.min(distance, threshold * 1.5));

      if (distance > 0) {
        e.preventDefault();
      }
    };

    const handleTouchEnd = async () => {
      if (!isPulling || isRefreshing) return;

      if (pullDistance > threshold) {
        isRefreshing = true;
        try {
          await onRefresh();
        } finally {
          isRefreshing = false;
        }
      }

      setIsPulling(false);
      setPullDistance(0);
    };

    element.addEventListener('touchstart', handleTouchStart, { passive: true });
    element.addEventListener('touchmove', handleTouchMove, { passive: false });
    element.addEventListener('touchend', handleTouchEnd);
    element.addEventListener('touchcancel', handleTouchEnd);

    return () => {
      element.removeEventListener('touchstart', handleTouchStart);
      element.removeEventListener('touchmove', handleTouchMove);
      element.removeEventListener('touchend', handleTouchEnd);
      element.removeEventListener('touchcancel', handleTouchEnd);
    };
  }, [elementRef, enabled, isPulling, pullDistance, threshold, onRefresh]);

  return {
    isPulling,
    pullDistance,
    pullProgress: Math.min(pullDistance / threshold, 1),
  };
};
