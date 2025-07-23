'use client';

import React, { useState } from 'react';
import { QuickView } from '@/components/ui/QuickView';
import { Eye, Zap } from 'lucide-react';

export default function QuickViewExamplesPage() {
  const [selectedProduct, setSelectedProduct] = useState<any>(null);
  const [isQuickViewOpen, setIsQuickViewOpen] = useState(false);

  const sampleProducts = [
    {
      id: '1',
      title: 'iPhone 14 Pro Max 256GB',
      price: '€899',
      description: `Продаю iPhone 14 Pro Max в отличном состоянии.
      
Характеристики:
- Память: 256GB
- Цвет: Deep Purple
- Состояние: как новый
- Комплект: полный (коробка, зарядка, документы)
- Гарантия: до марта 2024

Телефон использовался аккуратно, всегда в чехле и с защитным стеклом.
Без царапин и сколов. Батарея держит отлично - 96% здоровья.

Причина продажи: переход на Android.`,
      images: [
        'https://via.placeholder.com/600x400/9333ea/ffffff?text=iPhone+1',
        'https://via.placeholder.com/600x400/7c3aed/ffffff?text=iPhone+2',
        'https://via.placeholder.com/600x400/6d28d9/ffffff?text=iPhone+3',
      ],
      category: 'Электроника',
      seller: {
        name: 'Александр П.',
        rating: 4.8,
        totalReviews: 127,
      },
      location: {
        address: 'Белград, Врачар',
        distance: 2.5,
      },
      stats: {
        views: 342,
        favorites: 28,
      },
      condition: 'used',
    },
    {
      id: '2',
      title: 'Кожаная куртка Zara, размер M',
      price: '€120',
      description: `Стильная кожаная куртка от Zara.
      
Детали:
- Размер: M (подойдет на 46-48)
- Материал: натуральная кожа
- Цвет: черный
- Состояние: отличное

Куртка практически новая, носилась несколько раз.`,
      images: [
        'https://via.placeholder.com/600x400/000000/ffffff?text=Jacket+1',
        'https://via.placeholder.com/600x400/171717/ffffff?text=Jacket+2',
      ],
      category: 'Одежда',
      seller: {
        name: 'Мария К.',
        rating: 4.9,
        totalReviews: 89,
      },
      location: {
        address: 'Нови Сад, Центр',
        distance: 0.8,
      },
      stats: {
        views: 156,
        favorites: 12,
      },
      condition: 'used',
    },
    {
      id: '3',
      title: 'MacBook Air M2 13" 512GB',
      price: '€1,299',
      description: `MacBook Air с процессором M2.
      
Конфигурация:
- Процессор: Apple M2
- Память: 16GB
- SSD: 512GB
- Цвет: Midnight
- Год: 2023

Идеальное состояние, покупался для учебы.`,
      images: [
        'https://via.placeholder.com/600x400/1e3a8a/ffffff?text=MacBook',
      ],
      category: 'Компьютеры',
      seller: {
        name: 'Tech Store',
        rating: 4.7,
        totalReviews: 234,
      },
      location: {
        address: 'Белград, Нови Београд',
        distance: 5.2,
      },
      stats: {
        views: 521,
        favorites: 45,
      },
      condition: 'refurbished',
    },
  ];

  const openQuickView = (product: any) => {
    setSelectedProduct(product);
    setIsQuickViewOpen(true);
  };

  return (
    <div className="container mx-auto p-6 max-w-6xl">
      <h1 className="text-3xl font-bold mb-8">Быстрый просмотр товара</h1>

      {/* Описание */}
      <section className="mb-12">
        <div className="card bg-base-200 p-6">
          <h2 className="text-2xl font-semibold mb-4">О компоненте QuickView</h2>
          <p className="text-base-content/80 mb-4">
            Компонент быстрого просмотра позволяет пользователям детально изучить товар без перехода на отдельную страницу.
          </p>
          <ul className="list-disc list-inside space-y-2 text-base-content/80">
            <li>Модальное окно с полной информацией о товаре</li>
            <li>Галерея изображений с навигацией</li>
            <li>Информация о продавце и местоположении</li>
            <li>Кнопки быстрых действий</li>
            <li>Адаптивный дизайн для всех устройств</li>
          </ul>
        </div>
      </section>

      {/* Примеры товаров */}
      <section className="mb-12">
        <h2 className="text-2xl font-semibold mb-6">Нажмите на кнопку для просмотра</h2>
        
        <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-6">
          {sampleProducts.map((product) => (
            <div key={product.id} className="card bg-base-100 shadow-sm hover:shadow-lg transition-shadow">
              <figure className="h-48 bg-base-200 relative">
                <img
                  src={product.images[0]}
                  alt={product.title}
                  className="w-full h-full object-cover"
                />
                {/* Quick View Button Overlay */}
                <div className="absolute inset-0 bg-black/40 opacity-0 hover:opacity-100 transition-opacity flex items-center justify-center">
                  <button
                    onClick={() => openQuickView(product)}
                    className="btn btn-primary btn-sm gap-2"
                  >
                    <Eye className="w-4 h-4" />
                    Быстрый просмотр
                  </button>
                </div>
              </figure>
              
              <div className="card-body p-4">
                <h3 className="font-semibold line-clamp-2">{product.title}</h3>
                <p className="text-xl font-bold text-primary">{product.price}</p>
                
                <div className="flex items-center justify-between mt-4">
                  <span className="text-sm text-base-content/60">
                    {product.stats.views} просмотров
                  </span>
                  <button
                    onClick={() => openQuickView(product)}
                    className="btn btn-ghost btn-sm gap-1"
                  >
                    <Zap className="w-4 h-4" />
                    Быстрый просмотр
                  </button>
                </div>
              </div>
            </div>
          ))}
        </div>
      </section>

      {/* Интеграция с карточками */}
      <section className="mb-12">
        <h2 className="text-2xl font-semibold mb-6">Интеграция с карточками товаров</h2>
        
        <div className="grid grid-cols-1 lg:grid-cols-2 gap-6">
          {/* Hover версия */}
          <div className="card bg-base-200 p-6">
            <h3 className="text-lg font-semibold mb-4">При наведении</h3>
            <div className="card bg-base-100 shadow-sm group">
              <figure className="h-40 bg-base-300 relative overflow-hidden">
                <div className="absolute inset-0 bg-black/60 opacity-0 group-hover:opacity-100 transition-all duration-300 flex items-center justify-center">
                  <button
                    onClick={() => openQuickView(sampleProducts[0])}
                    className="btn btn-primary btn-sm transform translate-y-4 group-hover:translate-y-0 transition-transform duration-300"
                  >
                    <Eye className="w-4 h-4 mr-2" />
                    Быстрый просмотр
                  </button>
                </div>
              </figure>
              <div className="card-body p-4">
                <h4 className="font-medium">Товар с hover эффектом</h4>
                <p className="text-sm text-base-content/60">Наведите для просмотра</p>
              </div>
            </div>
          </div>

          {/* Icon версия */}
          <div className="card bg-base-200 p-6">
            <h3 className="text-lg font-semibold mb-4">С иконкой</h3>
            <div className="card bg-base-100 shadow-sm">
              <figure className="h-40 bg-base-300 relative">
                <button
                  onClick={() => openQuickView(sampleProducts[1])}
                  className="absolute top-2 right-2 btn btn-circle btn-sm bg-base-100/80 hover:bg-base-100"
                >
                  <Eye className="w-4 h-4" />
                </button>
              </figure>
              <div className="card-body p-4">
                <h4 className="font-medium">Товар с иконкой просмотра</h4>
                <p className="text-sm text-base-content/60">Кнопка всегда видна</p>
              </div>
            </div>
          </div>
        </div>
      </section>

      {/* Особенности */}
      <section className="mb-12">
        <h2 className="text-2xl font-semibold mb-6">Особенности компонента</h2>
        
        <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-4">
          <div className="card bg-base-100 p-4">
            <h3 className="font-medium mb-2">🖼️ Галерея изображений</h3>
            <p className="text-sm text-base-content/70">
              Переключение между фото с превью и счетчиком
            </p>
          </div>
          
          <div className="card bg-base-100 p-4">
            <h3 className="font-medium mb-2">📱 Адаптивность</h3>
            <p className="text-sm text-base-content/70">
              Оптимизирован для мобильных и десктопов
            </p>
          </div>
          
          <div className="card bg-base-100 p-4">
            <h3 className="font-medium mb-2">⚡ Быстрая загрузка</h3>
            <p className="text-sm text-base-content/70">
              Ленивая загрузка изображений
            </p>
          </div>
          
          <div className="card bg-base-100 p-4">
            <h3 className="font-medium mb-2">🎨 Анимации</h3>
            <p className="text-sm text-base-content/70">
              Плавное появление и переходы
            </p>
          </div>
          
          <div className="card bg-base-100 p-4">
            <h3 className="font-medium mb-2">🔒 Блокировка скролла</h3>
            <p className="text-sm text-base-content/70">
              Предотвращение прокрутки фона
            </p>
          </div>
          
          <div className="card bg-base-100 p-4">
            <h3 className="font-medium mb-2">♿ Доступность</h3>
            <p className="text-sm text-base-content/70">
              Поддержка клавиатуры и скрин-ридеров
            </p>
          </div>
        </div>
      </section>

      {/* Примеры кода */}
      <section className="card bg-base-200 p-6">
        <h2 className="text-2xl font-semibold mb-4">Использование</h2>
        <div className="mockup-code">
          <pre data-prefix="1"><code>{`import { QuickView } from '@/components/ui/QuickView';`}</code></pre>
          <pre data-prefix="2"><code>{``}</code></pre>
          <pre data-prefix="3"><code>{`const [isOpen, setIsOpen] = useState(false);`}</code></pre>
          <pre data-prefix="4"><code>{`const [product, setProduct] = useState(null);`}</code></pre>
          <pre data-prefix="5"><code>{``}</code></pre>
          <pre data-prefix="6"><code>{`<QuickView`}</code></pre>
          <pre data-prefix="7"><code>{`  isOpen={isOpen}`}</code></pre>
          <pre data-prefix="8"><code>{`  onClose={() => setIsOpen(false)}`}</code></pre>
          <pre data-prefix="9"><code>{`  product={product}`}</code></pre>
          <pre data-prefix="10"><code>{`/>`}</code></pre>
        </div>
      </section>

      {/* Quick View Modal */}
      {selectedProduct && (
        <QuickView
          isOpen={isQuickViewOpen}
          onClose={() => setIsQuickViewOpen(false)}
          product={selectedProduct}
        />
      )}
    </div>
  );
}