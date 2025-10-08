'use client';

import React, {
  createContext,
  useContext,
  useEffect,
  useState,
  useCallback,
} from 'react';
import { usePathname } from 'next/navigation';

export interface ABTest {
  id: string;
  name: string;
  description?: string;
  variants: ABVariant[];
  targeting?: ABTargeting;
  allocation?: number; // процент трафика для теста (0-100)
  status: 'draft' | 'running' | 'paused' | 'completed';
  startDate?: Date;
  endDate?: Date;
  metrics: ABMetrics[];
}

export interface ABVariant {
  id: string;
  name: string;
  weight: number; // вес варианта (для распределения трафика)
  isControl?: boolean;
  component?: React.ComponentType<any>;
  props?: Record<string, any>;
  config?: Record<string, any>;
}

export interface ABTargeting {
  urls?: string[];
  segments?: string[];
  devices?: ('mobile' | 'tablet' | 'desktop')[];
  countries?: string[];
  languages?: string[];
}

export interface ABMetrics {
  name: string;
  type: 'conversion' | 'engagement' | 'revenue' | 'custom';
  goal?: number;
}

interface ABTestEvent {
  testId: string;
  variantId: string;
  event: string;
  value?: any;
  timestamp: Date;
}

interface ABTestContextValue {
  tests: ABTest[];
  activeTests: Map<string, string>; // testId -> variantId
  getVariant: (testId: string) => string | null;
  trackEvent: (testId: string, event: string, value?: any) => void;
  isInTest: (testId: string) => boolean;
  getTestComponent: (testId: string) => React.ComponentType<any> | null;
  forceVariant: (testId: string, variantId: string) => void;
  clearForceVariant: (testId: string) => void;
}

const ABTestContext = createContext<ABTestContextValue | null>(null);

const STORAGE_KEY = 'ab_test_assignments';
const EVENTS_KEY = 'ab_test_events';
const FORCE_KEY = 'ab_test_force';

export const ABTestProvider: React.FC<{
  children: React.ReactNode;
  tests: ABTest[];
  userId?: string;
  debug?: boolean;
}> = ({ children, tests, userId, debug = false }) => {
  const [activeTests, setActiveTests] = useState<Map<string, string>>(
    new Map()
  );
  const [forcedVariants, setForcedVariants] = useState<Map<string, string>>(
    new Map()
  );
  const pathname = usePathname();

  // Генерация уникального ID пользователя
  const generateUserId = (): string => {
    return Date.now().toString(36) + Math.random().toString(36).substr(2);
  };

  // Хэш-функция для распределения
  const hashCode = (str: string): number => {
    let hash = 0;
    for (let i = 0; i < str.length; i++) {
      const char = str.charCodeAt(i);
      hash = (hash << 5) - hash + char;
      hash = hash & hash;
    }
    return hash;
  };

  // Определение типа устройства
  const getDeviceType = (): 'mobile' | 'tablet' | 'desktop' => {
    if (typeof window === 'undefined') return 'desktop';

    const width = window.innerWidth;
    if (width < 768) return 'mobile';
    if (width < 1024) return 'tablet';
    return 'desktop';
  };

  // Загрузка сохраненных назначений
  const loadStoredAssignments = (): Record<string, string> => {
    if (typeof window === 'undefined') return {};

    try {
      const stored = localStorage.getItem(STORAGE_KEY);
      return stored ? JSON.parse(stored) : {};
    } catch (err) {
      console.error('Failed to load AB test assignments:', err);
      return {};
    }
  };

  // Сохранение назначений
  const saveAssignments = (assignments: Map<string, string>) => {
    if (typeof window === 'undefined') return;

    try {
      const obj = Object.fromEntries(assignments);
      localStorage.setItem(STORAGE_KEY, JSON.stringify(obj));
    } catch (err) {
      console.error('Failed to save AB test assignments:', err);
    }
  };

  // Загрузка принудительных вариантов из localStorage
  const loadForcedVariants = () => {
    if (typeof window === 'undefined') return;

    try {
      const stored = localStorage.getItem(FORCE_KEY);
      if (stored) {
        const parsed = JSON.parse(stored);
        setForcedVariants(new Map(Object.entries(parsed)));
      }
    } catch (err) {
      console.error('Failed to load forced variants:', err);
    }
  };

  // Получение или создание ID пользователя
  const getOrCreateUserId = useCallback((): string => {
    if (typeof window === 'undefined') return 'server';

    let id = localStorage.getItem('ab_user_id');
    if (!id) {
      id = generateUserId();
      localStorage.setItem('ab_user_id', id);
    }
    return id;
  }, []);

  // Проверка, должен ли тест запускаться
  const shouldRunTest = useCallback(
    (test: ABTest): boolean => {
      // Проверка статуса
      if (test.status !== 'running') return false;

      // Проверка дат
      const now = new Date();
      if (test.startDate && now < test.startDate) return false;
      if (test.endDate && now > test.endDate) return false;

      // Проверка allocation
      if (test.allocation !== undefined && test.allocation < 100) {
        const hash = hashCode(userId || getOrCreateUserId());
        const bucket = Math.abs(hash % 100);
        if (bucket >= test.allocation) return false;
      }

      // Проверка targeting
      if (test.targeting) {
        // Проверка URL
        if (
          test.targeting.urls &&
          !test.targeting.urls.some((url) => pathname.includes(url))
        ) {
          return false;
        }

        // Проверка устройства
        if (test.targeting.devices && typeof window !== 'undefined') {
          const device = getDeviceType();
          if (!test.targeting.devices.includes(device)) return false;
        }
      }

      return true;
    },
    [userId, pathname, getOrCreateUserId]
  );

  // Выбор варианта на основе весов
  const selectVariant = useCallback(
    (test: ABTest): ABVariant | null => {
      const totalWeight = test.variants.reduce((sum, v) => sum + v.weight, 0);
      if (totalWeight === 0) return null;

      const hash = hashCode(userId || getOrCreateUserId() + test.id);
      const bucket = Math.abs(hash % totalWeight);

      let cumWeight = 0;
      for (const variant of test.variants) {
        cumWeight += variant.weight;
        if (bucket < cumWeight) {
          return variant;
        }
      }

      return test.variants[test.variants.length - 1];
    },
    [userId, getOrCreateUserId]
  );

  // Инициализация активных тестов
  const initializeTests = useCallback(() => {
    const assignments = new Map<string, string>();
    const storedAssignments = loadStoredAssignments();

    tests.forEach((test) => {
      // Проверяем, подходит ли тест для текущих условий
      if (!shouldRunTest(test)) return;

      // Проверяем принудительный вариант
      const forcedVariant = forcedVariants.get(test.id);
      if (forcedVariant) {
        assignments.set(test.id, forcedVariant);
        if (debug) {
          console.log(
            `[ABTest] Using forced variant for ${test.id}: ${forcedVariant}`
          );
        }
        return;
      }

      // Проверяем сохраненное назначение
      const storedVariant = storedAssignments[test.id];
      if (storedVariant && test.variants.some((v) => v.id === storedVariant)) {
        assignments.set(test.id, storedVariant);
        return;
      }

      // Назначаем новый вариант
      const variant = selectVariant(test);
      if (variant) {
        assignments.set(test.id, variant.id);
      }
    });

    setActiveTests(assignments);
    saveAssignments(assignments);
  }, [tests, forcedVariants, debug, selectVariant, shouldRunTest]);

  // Получение варианта для теста
  const getVariant = (testId: string): string | null => {
    return activeTests.get(testId) || null;
  };

  // Проверка участия в тесте
  const isInTest = (testId: string): boolean => {
    return activeTests.has(testId);
  };

  // Получение компонента для варианта
  const getTestComponent = (
    testId: string
  ): React.ComponentType<any> | null => {
    const variantId = activeTests.get(testId);
    if (!variantId) return null;

    const test = tests.find((t) => t.id === testId);
    if (!test) return null;

    const variant = test.variants.find((v) => v.id === variantId);
    return variant?.component || null;
  };

  // Трекинг событий
  const trackEvent = (testId: string, event: string, value?: any) => {
    const variantId = activeTests.get(testId);
    if (!variantId) return;

    const eventData: ABTestEvent = {
      testId,
      variantId,
      event,
      value,
      timestamp: new Date(),
    };

    // Сохраняем событие локально
    saveEvent(eventData);

    // Отправляем на сервер (если есть endpoint)
    sendEventToServer(eventData);

    if (debug) {
      console.log('[ABTest] Event tracked:', eventData);
    }
  };

  // Сохранение события локально
  const saveEvent = (event: ABTestEvent) => {
    if (typeof window === 'undefined') return;

    try {
      const stored = localStorage.getItem(EVENTS_KEY);
      const events = stored ? JSON.parse(stored) : [];
      events.push(event);

      // Ограничиваем количество событий
      if (events.length > 1000) {
        events.splice(0, events.length - 1000);
      }

      localStorage.setItem(EVENTS_KEY, JSON.stringify(events));
    } catch (err) {
      console.error('Failed to save AB test event:', err);
    }
  };

  // Отправка события на сервер
  const sendEventToServer = async (_event: ABTestEvent) => {
    // TODO: Реализовать отправку на сервер
    // await fetch('/api/v1/abtest/events', {
    //   method: 'POST',
    //   headers: { 'Content-Type': 'application/json' },
    //   body: JSON.stringify(event)
    // });
  };

  // Принудительное назначение варианта (для QA)
  const forceVariant = (testId: string, variantId: string) => {
    const newForced = new Map(forcedVariants);
    newForced.set(testId, variantId);
    setForcedVariants(newForced);

    // Сохраняем в localStorage
    if (typeof window !== 'undefined') {
      const obj = Object.fromEntries(newForced);
      localStorage.setItem(FORCE_KEY, JSON.stringify(obj));
    }

    // Переинициализируем тесты
    initializeTests();
  };

  // Очистка принудительного варианта
  const clearForceVariant = (testId: string) => {
    const newForced = new Map(forcedVariants);
    newForced.delete(testId);
    setForcedVariants(newForced);

    // Обновляем localStorage
    if (typeof window !== 'undefined') {
      if (newForced.size > 0) {
        const obj = Object.fromEntries(newForced);
        localStorage.setItem(FORCE_KEY, JSON.stringify(obj));
      } else {
        localStorage.removeItem(FORCE_KEY);
      }
    }

    // Переинициализируем тесты
    initializeTests();
  };

  // Инициализация тестов при загрузке
  useEffect(() => {
    initializeTests();
    loadForcedVariants();
  }, [tests, pathname, initializeTests]);

  const contextValue: ABTestContextValue = {
    tests,
    activeTests,
    getVariant,
    trackEvent,
    isInTest,
    getTestComponent,
    forceVariant,
    clearForceVariant,
  };

  return (
    <ABTestContext.Provider value={contextValue}>
      {children}
      {debug && <ABTestDebugPanel />}
    </ABTestContext.Provider>
  );
};

// Хук для использования A/B тестов
export const useABTest = (testId: string) => {
  const context = useContext(ABTestContext);
  if (!context) {
    throw new Error('useABTest must be used within ABTestProvider');
  }

  return {
    variant: context.getVariant(testId),
    isInTest: context.isInTest(testId),
    trackEvent: (event: string, value?: any) =>
      context.trackEvent(testId, event, value),
    Component: context.getTestComponent(testId),
  };
};

// Компонент для A/B тестирования
export const ABTestComponent: React.FC<{
  testId: string;
  children: (variant: string | null) => React.ReactNode;
}> = ({ testId, children }) => {
  const { variant } = useABTest(testId);
  return <>{children(variant)}</>;
};

// Панель отладки для A/B тестов
const ABTestDebugPanel: React.FC = () => {
  const [isOpen, setIsOpen] = useState(false);
  const context = useContext(ABTestContext);
  if (!context) return null;

  return (
    <div className="fixed bottom-4 right-4 z-[9999]">
      <button
        onClick={() => setIsOpen(!isOpen)}
        className="btn btn-sm btn-primary"
      >
        A/B Tests ({context.activeTests.size})
      </button>

      {isOpen && (
        <div className="absolute bottom-12 right-0 w-96 bg-base-100 border border-base-300 rounded-lg shadow-xl p-4 max-h-96 overflow-y-auto">
          <h3 className="font-bold mb-4">Active A/B Tests</h3>

          {context.tests.map((test) => {
            const variantId = context.getVariant(test.id);
            const variant = test.variants.find((v) => v.id === variantId);

            return (
              <div key={test.id} className="mb-4 p-2 border rounded">
                <div className="font-semibold">{test.name}</div>
                <div className="text-sm text-base-content/70">
                  ID: {test.id}
                </div>
                {variant ? (
                  <>
                    <div className="text-sm">
                      Variant:{' '}
                      <span className="badge badge-sm">{variant?.name}</span>
                    </div>
                    <div className="mt-2 flex gap-2">
                      {test.variants.map((v) => (
                        <button
                          key={v.id}
                          onClick={() => context.forceVariant(test.id, v.id)}
                          className={`btn btn-xs ${v.id === variantId ? 'btn-primary' : 'btn-ghost'}`}
                        >
                          {v.name}
                        </button>
                      ))}
                      <button
                        onClick={() => context.clearForceVariant(test.id)}
                        className="btn btn-xs btn-error"
                      >
                        Clear
                      </button>
                    </div>
                  </>
                ) : (
                  <div className="text-sm text-base-content/50">
                    Not in test
                  </div>
                )}
              </div>
            );
          })}
        </div>
      )}
    </div>
  );
};

export default ABTestProvider;
