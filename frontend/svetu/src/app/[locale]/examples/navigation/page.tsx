'use client';

import React, { useState } from 'react';
import { MobileBottomNav } from '@/components/navigation/MobileBottomNav';
import { EnhancedMobileBottomNav } from '@/components/navigation/EnhancedMobileBottomNav';
import { SmartMobileBottomNav } from '@/components/navigation/SmartMobileBottomNav';

export default function NavigationExamplesPage() {
  const [showSmart, setShowSmart] = useState(false);

  return (
    <div className="min-h-screen pb-20">
      <div className="container mx-auto p-6 max-w-4xl">
        <h1 className="text-3xl font-bold mb-8">Mobile Navigation Examples</h1>

        {/* Описание компонентов */}
        <section className="space-y-8 mb-12">
          <div className="card bg-base-200 p-6">
            <h2 className="text-2xl font-semibold mb-4">1. Базовая навигация</h2>
            <p className="mb-4 text-base-content/70">
              Стандартная мобильная навигация с активными индикаторами и бейджами.
            </p>
            <ul className="list-disc list-inside space-y-1 text-sm">
              <li>Активный индикатор сверху</li>
              <li>Поддержка бейджей для уведомлений</li>
              <li>Условное отображение для авторизованных пользователей</li>
              <li>Мультиязычность</li>
            </ul>
          </div>

          <div className="card bg-base-200 p-6">
            <h2 className="text-2xl font-semibold mb-4">2. Улучшенная навигация</h2>
            <p className="mb-4 text-base-content/70">
              Версия с дополнительными анимациями и визуальными эффектами.
            </p>
            <ul className="list-disc list-inside space-y-1 text-sm">
              <li>Анимированный индикатор перемещения</li>
              <li>Hover эффекты и масштабирование</li>
              <li>Pulse анимация для бейджей</li>
              <li>Ripple эффект для кнопки создания</li>
              <li>Плавные переходы состояний</li>
            </ul>
          </div>

          <div className="card bg-base-200 p-6">
            <h2 className="text-2xl font-semibold mb-4">3. Умная навигация</h2>
            <p className="mb-4 text-base-content/70">
              Навигация, которая скрывается при скролле вниз и появляется при скролле вверх.
            </p>
            <button
              className="btn btn-primary"
              onClick={() => setShowSmart(!showSmart)}
            >
              {showSmart ? 'Скрыть' : 'Показать'} умную навигацию
            </button>
          </div>
        </section>

        {/* Сравнение функций */}
        <section className="mb-12">
          <h2 className="text-2xl font-semibold mb-6">Сравнение функций</h2>
          <div className="overflow-x-auto">
            <table className="table table-zebra">
              <thead>
                <tr>
                  <th>Функция</th>
                  <th>Базовая</th>
                  <th>Улучшенная</th>
                  <th>Умная</th>
                </tr>
              </thead>
              <tbody>
                <tr>
                  <td>Активные индикаторы</td>
                  <td>✅</td>
                  <td>✅</td>
                  <td>✅</td>
                </tr>
                <tr>
                  <td>Бейджи уведомлений</td>
                  <td>✅</td>
                  <td>✅ + анимация</td>
                  <td>✅ + анимация</td>
                </tr>
                <tr>
                  <td>Hover эффекты</td>
                  <td>❌</td>
                  <td>✅</td>
                  <td>✅</td>
                </tr>
                <tr>
                  <td>Анимированный индикатор</td>
                  <td>❌</td>
                  <td>✅</td>
                  <td>✅</td>
                </tr>
                <tr>
                  <td>Скрытие при скролле</td>
                  <td>❌</td>
                  <td>❌</td>
                  <td>✅</td>
                </tr>
                <tr>
                  <td>Ripple эффекты</td>
                  <td>❌</td>
                  <td>✅</td>
                  <td>✅</td>
                </tr>
              </tbody>
            </table>
          </div>
        </section>

        {/* Примеры кода */}
        <section className="card bg-base-200 p-6 mb-12">
          <h2 className="text-2xl font-semibold mb-4">Примеры использования</h2>
          <div className="mockup-code">
            <pre data-prefix="1"><code>{`// Базовая навигация`}</code></pre>
            <pre data-prefix="2"><code>{`import { MobileBottomNav } from '@/components/navigation';`}</code></pre>
            <pre data-prefix="3"><code>{`<MobileBottomNav />`}</code></pre>
            <pre data-prefix="4"><code>{``}</code></pre>
            <pre data-prefix="5"><code>{`// Улучшенная навигация`}</code></pre>
            <pre data-prefix="6"><code>{`import { EnhancedMobileBottomNav } from '@/components/navigation';`}</code></pre>
            <pre data-prefix="7"><code>{`<EnhancedMobileBottomNav />`}</code></pre>
            <pre data-prefix="8"><code>{``}</code></pre>
            <pre data-prefix="9"><code>{`// Умная навигация с автоскрытием`}</code></pre>
            <pre data-prefix="10"><code>{`import { SmartMobileBottomNav } from '@/components/navigation';`}</code></pre>
            <pre data-prefix="11"><code>{`<SmartMobileBottomNav />`}</code></pre>
          </div>
        </section>

        {/* Заполнитель контента для демонстрации скролла */}
        {showSmart && (
          <section className="space-y-4">
            <div className="alert alert-info">
              <svg xmlns="http://www.w3.org/2000/svg" fill="none" viewBox="0 0 24 24" className="stroke-info shrink-0 w-6 h-6">
                <path strokeLinecap="round" strokeLinejoin="round" strokeWidth="2" d="M13 16h-1v-4h-1m1-4h.01M21 12a9 9 0 11-18 0 9 9 0 0118 0z"></path>
              </svg>
              <span>Прокрутите страницу вниз/вверх чтобы увидеть умную навигацию в действии</span>
            </div>
            
            {[...Array(20)].map((_, i) => (
              <div key={i} className="card bg-base-100 shadow-sm p-6">
                <h3 className="text-lg font-medium mb-2">Контент блок {i + 1}</h3>
                <p className="text-base-content/70">
                  Это демонстрационный контент для проверки скрытия навигации при скролле.
                  Навигация скрывается при прокрутке вниз и появляется при прокрутке вверх.
                </p>
              </div>
            ))}
          </section>
        )}
      </div>

      {/* Навигации */}
      <div className="fixed bottom-0 left-0 right-0 md:hidden">
        {!showSmart ? (
          <div className="relative">
            {/* Базовая навигация (скрыта) */}
            <div className="opacity-0 pointer-events-none">
              <MobileBottomNav />
            </div>
            
            {/* Улучшенная навигация (показана) */}
            <div className="absolute inset-0">
              <EnhancedMobileBottomNav />
            </div>
          </div>
        ) : (
          <SmartMobileBottomNav />
        )}
      </div>
    </div>
  );
}