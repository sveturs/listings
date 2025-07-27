'use client';

import { useState } from 'react';
import Link from 'next/link';
import { AnimatedSection } from '@/components/ui/AnimatedSection';

type ExampleView = 'modern' | 'smart' | 'ai-powered';

export default function ListingEditExamplesPage() {
  const [activeExample, setActiveExample] = useState<ExampleView>('modern');

  const renderModernExample = () => (
    <AnimatedSection animation="fadeIn" className="space-y-6">
      <div className="bg-base-100 rounded-xl p-6 shadow-xl">
        <h3 className="text-2xl font-bold mb-6">✨ Современный редактор</h3>

        {/* Tabs для секций */}
        <div className="tabs tabs-boxed mb-6">
          <a className="tab tab-active">Основное</a>
          <a className="tab">Фото</a>
          <a className="tab">Атрибуты</a>
          <a className="tab">Локация</a>
          <a className="tab">SEO</a>
        </div>

        {/* Основная информация */}
        <div className="space-y-4">
          <div className="form-control">
            <label className="label">
              <span className="label-text font-medium">Заголовок</span>
              <span className="label-text-alt">0/100</span>
            </label>
            <input
              type="text"
              className="input input-bordered w-full"
              defaultValue="Volkswagen Touran"
            />
            <label className="label">
              <span className="label-text-alt text-success">SEO: Отлично</span>
            </label>
          </div>

          <div className="form-control">
            <label className="label">
              <span className="label-text font-medium">Описание</span>
              <div className="dropdown dropdown-end">
                <label tabIndex={0} className="btn btn-ghost btn-xs">
                  AI помощник
                </label>
                <ul
                  tabIndex={0}
                  className="dropdown-content menu p-2 shadow bg-base-100 rounded-box w-52"
                >
                  <li>
                    <a>Улучшить текст</a>
                  </li>
                  <li>
                    <a>Добавить детали</a>
                  </li>
                  <li>
                    <a>Оптимизировать для SEO</a>
                  </li>
                </ul>
              </div>
            </label>
            <textarea
              className="textarea textarea-bordered h-32"
              defaultValue="Volkswagen Touran 2.0 TDI • 2012 • Идеальное состояние"
            />
          </div>

          <div className="grid grid-cols-1 md:grid-cols-2 gap-4">
            <div className="form-control">
              <label className="label">
                <span className="label-text font-medium">Цена</span>
              </label>
              <div className="join">
                <input
                  type="number"
                  className="input input-bordered join-item flex-1"
                  defaultValue="600000"
                />
                <select className="select select-bordered join-item">
                  <option>RSD</option>
                  <option>EUR</option>
                  <option>USD</option>
                </select>
              </div>
              <label className="label">
                <span className="label-text-alt">
                  Средняя цена: 550,000 RSD
                </span>
              </label>
            </div>

            <div className="form-control">
              <label className="label">
                <span className="label-text font-medium">Состояние</span>
              </label>
              <select className="select select-bordered w-full">
                <option>Новое</option>
                <option selected>Б/у</option>
                <option>Восстановленное</option>
              </select>
            </div>
          </div>

          {/* Живая превью карточки */}
          <div className="divider">Предпросмотр</div>
          <div className="mockup-browser border bg-base-300">
            <div className="mockup-browser-toolbar">
              <div className="input">Ваше объявление на сайте</div>
            </div>
            <div className="px-4 pb-4 bg-base-200">
              <div className="card bg-base-100 shadow-xl">
                <figure className="h-48 bg-base-300"></figure>
                <div className="card-body">
                  <h2 className="card-title">Volkswagen Touran</h2>
                  <p>Volkswagen Touran 2.0 TDI • 2012 • Идеальное...</p>
                  <div className="flex justify-between items-center">
                    <span className="text-2xl font-bold text-primary">
                      600,000 RSD
                    </span>
                    <div className="badge badge-outline">Б/у</div>
                  </div>
                </div>
              </div>
            </div>
          </div>
        </div>
      </div>
    </AnimatedSection>
  );

  const renderSmartExample = () => (
    <AnimatedSection animation="fadeIn" className="space-y-6">
      <div className="bg-base-100 rounded-xl p-6 shadow-xl">
        <h3 className="text-2xl font-bold mb-6">🧠 Умный редактор</h3>

        {/* Прогресс заполнения */}
        <div className="mb-6">
          <div className="flex justify-between items-center mb-2">
            <span className="text-sm font-medium">Заполнено на 85%</span>
            <span className="text-sm text-base-content/70">
              Осталось 3 поля
            </span>
          </div>
          <progress
            className="progress progress-primary w-full"
            value="85"
            max="100"
          ></progress>
        </div>

        {/* Умные подсказки */}
        <div className="alert alert-info mb-6">
          <svg
            xmlns="http://www.w3.org/2000/svg"
            fill="none"
            viewBox="0 0 24 24"
            className="stroke-current shrink-0 w-6 h-6"
          >
            <path
              strokeLinecap="round"
              strokeLinejoin="round"
              strokeWidth="2"
              d="M13 16h-1v-4h-1m1-4h.01M21 12a9 9 0 11-18 0 9 9 0 0118 0z"
            ></path>
          </svg>
          <div>
            <h3 className="font-bold">Рекомендация</h3>
            <div className="text-xs">
              Добавьте фото интерьера для увеличения просмотров на 40%
            </div>
          </div>
        </div>

        <div className="grid grid-cols-1 lg:grid-cols-2 gap-6">
          {/* Левая колонка - форма */}
          <div className="space-y-4">
            <div className="form-control">
              <label className="label">
                <span className="label-text font-medium">
                  Быстрое заполнение
                </span>
              </label>
              <div className="join w-full">
                <input
                  type="text"
                  className="input input-bordered join-item flex-1"
                  placeholder="Вставьте ссылку на похожее объявление"
                />
                <button className="btn btn-primary join-item">Импорт</button>
              </div>
            </div>

            <div className="divider">или заполните вручную</div>

            {/* Умная категоризация */}
            <div className="form-control">
              <label className="label">
                <span className="label-text font-medium">Категория</span>
                <span className="badge badge-success badge-sm">
                  Определена автоматически
                </span>
              </label>
              <select className="select select-bordered w-full">
                <option selected>Автомобили / Минивэны</option>
                <option>Автомобили / Седаны</option>
                <option>Автомобили / Внедорожники</option>
              </select>
            </div>

            {/* Динамические атрибуты */}
            <div className="card bg-base-200">
              <div className="card-body p-4">
                <h4 className="font-medium mb-3">Атрибуты автомобиля</h4>
                <div className="grid grid-cols-2 gap-3">
                  <input
                    type="text"
                    placeholder="Марка"
                    className="input input-sm input-bordered"
                    defaultValue="Volkswagen"
                  />
                  <input
                    type="text"
                    placeholder="Модель"
                    className="input input-sm input-bordered"
                    defaultValue="Touran"
                  />
                  <input
                    type="number"
                    placeholder="Год"
                    className="input input-sm input-bordered"
                    defaultValue="2012"
                  />
                  <input
                    type="text"
                    placeholder="Пробег"
                    className="input input-sm input-bordered"
                    defaultValue="150,000 км"
                  />
                  <select className="select select-sm select-bordered">
                    <option>Бензин</option>
                    <option selected>Дизель</option>
                    <option>Электро</option>
                    <option>Гибрид</option>
                  </select>
                  <select className="select select-sm select-bordered">
                    <option>Механика</option>
                    <option selected>Автомат</option>
                  </select>
                </div>
              </div>
            </div>

            {/* Умное ценообразование */}
            <div className="form-control">
              <label className="label">
                <span className="label-text font-medium">Цена</span>
                <span className="label-text-alt">Анализ рынка</span>
              </label>
              <input
                type="range"
                min="400000"
                max="800000"
                defaultValue="600000"
                className="range range-primary"
              />
              <div className="w-full flex justify-between text-xs px-2">
                <span>400K</span>
                <span className="text-primary font-bold">600K</span>
                <span>800K</span>
              </div>
              <div className="stats stats-vertical shadow mt-2">
                <div className="stat py-2">
                  <div className="stat-title text-xs">Минимальная</div>
                  <div className="stat-value text-lg">420,000 RSD</div>
                </div>
                <div className="stat py-2">
                  <div className="stat-title text-xs">Средняя</div>
                  <div className="stat-value text-lg text-primary">
                    550,000 RSD
                  </div>
                </div>
                <div className="stat py-2">
                  <div className="stat-title text-xs">Максимальная</div>
                  <div className="stat-value text-lg">780,000 RSD</div>
                </div>
              </div>
            </div>
          </div>

          {/* Правая колонка - помощники */}
          <div className="space-y-4">
            {/* Фото с AI анализом */}
            <div className="card bg-base-200">
              <div className="card-body p-4">
                <h4 className="font-medium mb-3">Фотографии</h4>
                <div className="grid grid-cols-3 gap-2">
                  <div className="aspect-square bg-base-300 rounded-lg relative">
                    <div className="absolute inset-0 flex items-center justify-center">
                      <span className="text-4xl">📷</span>
                    </div>
                    <div className="absolute top-1 right-1">
                      <div className="badge badge-success badge-xs">
                        Главное
                      </div>
                    </div>
                  </div>
                  <div className="aspect-square bg-base-300 rounded-lg relative">
                    <div className="absolute inset-0 flex items-center justify-center">
                      <span className="text-4xl">🚗</span>
                    </div>
                  </div>
                  <div className="aspect-square border-2 border-dashed border-base-300 rounded-lg relative cursor-pointer hover:border-primary">
                    <div className="absolute inset-0 flex items-center justify-center">
                      <span className="text-2xl">+</span>
                    </div>
                  </div>
                </div>
                <div className="text-xs text-base-content/70 mt-2">
                  AI обнаружил: Volkswagen Touran, серый цвет, хорошее состояние
                </div>
              </div>
            </div>

            {/* График активности */}
            <div className="card bg-base-200">
              <div className="card-body p-4">
                <h4 className="font-medium mb-3">Лучшее время публикации</h4>
                <div className="flex justify-between items-end h-20">
                  <div
                    className="w-8 bg-base-300 rounded"
                    style={{ height: '40%' }}
                  ></div>
                  <div
                    className="w-8 bg-base-300 rounded"
                    style={{ height: '60%' }}
                  ></div>
                  <div
                    className="w-8 bg-primary rounded"
                    style={{ height: '100%' }}
                  ></div>
                  <div
                    className="w-8 bg-primary rounded"
                    style={{ height: '90%' }}
                  ></div>
                  <div
                    className="w-8 bg-base-300 rounded"
                    style={{ height: '70%' }}
                  ></div>
                  <div
                    className="w-8 bg-base-300 rounded"
                    style={{ height: '50%' }}
                  ></div>
                  <div
                    className="w-8 bg-base-300 rounded"
                    style={{ height: '30%' }}
                  ></div>
                </div>
                <div className="text-xs text-center mt-2">
                  <span className="text-primary font-bold">Среда-Четверг</span>{' '}
                  - максимальная активность
                </div>
              </div>
            </div>

            {/* Проверка качества */}
            <div className="card bg-base-200">
              <div className="card-body p-4">
                <h4 className="font-medium mb-3">Качество объявления</h4>
                <div className="space-y-2">
                  <div className="flex items-center gap-2">
                    <span className="text-success">✓</span>
                    <span className="text-sm">Заголовок оптимальной длины</span>
                  </div>
                  <div className="flex items-center gap-2">
                    <span className="text-success">✓</span>
                    <span className="text-sm">Есть основное фото</span>
                  </div>
                  <div className="flex items-center gap-2">
                    <span className="text-warning">!</span>
                    <span className="text-sm">Добавьте больше фото</span>
                  </div>
                  <div className="flex items-center gap-2">
                    <span className="text-warning">!</span>
                    <span className="text-sm">Укажите точный адрес</span>
                  </div>
                </div>
                <div className="divider my-2"></div>
                <div className="flex justify-between items-center">
                  <span className="text-sm font-medium">Оценка качества</span>
                  <div
                    className="radial-progress text-primary"
                    style={{ '--value': 75 } as any}
                  >
                    75%
                  </div>
                </div>
              </div>
            </div>
          </div>
        </div>

        {/* Действия */}
        <div className="flex gap-4 mt-6">
          <button className="btn btn-ghost flex-1">Сохранить черновик</button>
          <button className="btn btn-primary flex-1">Опубликовать</button>
        </div>
      </div>
    </AnimatedSection>
  );

  const renderAIPoweredExample = () => (
    <AnimatedSection animation="fadeIn" className="space-y-6">
      <div className="bg-base-100 rounded-xl p-6 shadow-xl">
        <h3 className="text-2xl font-bold mb-6">🤖 AI-Powered редактор</h3>

        {/* AI Chat Interface */}
        <div className="chat chat-start mb-6">
          <div className="chat-bubble chat-bubble-primary">
            Привет! Я помогу обновить ваше объявление. Что бы вы хотели
            изменить?
          </div>
        </div>

        <div className="join w-full mb-6">
          <input
            type="text"
            className="input input-bordered join-item flex-1"
            placeholder="Например: 'Снизь цену на 10%' или 'Сделай описание более привлекательным'"
          />
          <button className="btn btn-primary join-item">
            <svg
              className="w-5 h-5"
              fill="none"
              stroke="currentColor"
              viewBox="0 0 24 24"
            >
              <path
                strokeLinecap="round"
                strokeLinejoin="round"
                strokeWidth={2}
                d="M13 5l7 7-7 7M5 5l7 7-7 7"
              />
            </svg>
          </button>
        </div>

        {/* Quick AI Actions */}
        <div className="grid grid-cols-2 md:grid-cols-4 gap-2 mb-6">
          <button className="btn btn-sm btn-outline">
            🎯 Оптимизировать цену
          </button>
          <button className="btn btn-sm btn-outline">
            ✨ Улучшить описание
          </button>
          <button className="btn btn-sm btn-outline">📸 Анализ фото</button>
          <button className="btn btn-sm btn-outline">🌐 Перевести</button>
        </div>

        {/* AI-Generated Content */}
        <div className="space-y-4">
          <div className="card bg-gradient-to-r from-primary/10 to-secondary/10 border-2 border-primary/20">
            <div className="card-body">
              <div className="flex items-start gap-2">
                <span className="text-2xl">🤖</span>
                <div className="flex-1">
                  <h4 className="font-bold mb-2">
                    AI предлагает 3 варианта заголовка:
                  </h4>
                  <div className="space-y-2">
                    <label className="cursor-pointer flex items-center gap-2 p-2 rounded hover:bg-base-100">
                      <input
                        type="radio"
                        name="title"
                        className="radio radio-primary"
                        checked
                      />
                      <span>
                        Volkswagen Touran 2.0 TDI - Семейный минивэн в отличном
                        состоянии
                      </span>
                      <div className="badge badge-success badge-sm">
                        +45% CTR
                      </div>
                    </label>
                    <label className="cursor-pointer flex items-center gap-2 p-2 rounded hover:bg-base-100">
                      <input
                        type="radio"
                        name="title"
                        className="radio radio-primary"
                      />
                      <span>
                        🚗 VW Touran 2012 | Автомат | Дизель | Идеал для семьи
                      </span>
                      <div className="badge badge-info badge-sm">+30% CTR</div>
                    </label>
                    <label className="cursor-pointer flex items-center gap-2 p-2 rounded hover:bg-base-100">
                      <input
                        type="radio"
                        name="title"
                        className="radio radio-primary"
                      />
                      <span>Срочно! Volkswagen Touran по цене ниже рынка</span>
                      <div className="badge badge-warning badge-sm">
                        +25% CTR
                      </div>
                    </label>
                  </div>
                </div>
              </div>
            </div>
          </div>

          {/* AI Analysis Dashboard */}
          <div className="grid grid-cols-1 md:grid-cols-3 gap-4">
            <div className="card bg-base-200">
              <div className="card-body p-4">
                <h4 className="text-sm font-medium mb-2">Анализ конкурентов</h4>
                <div className="text-2xl font-bold text-primary">12</div>
                <p className="text-xs text-base-content/70">
                  похожих объявлений
                </p>
                <div className="text-xs mt-2">
                  Ваша цена на{' '}
                  <span className="text-success font-bold">8%</span> ниже
                  средней
                </div>
              </div>
            </div>

            <div className="card bg-base-200">
              <div className="card-body p-4">
                <h4 className="text-sm font-medium mb-2">Прогноз просмотров</h4>
                <div className="text-2xl font-bold text-secondary">450+</div>
                <p className="text-xs text-base-content/70">в первую неделю</p>
                <progress
                  className="progress progress-secondary w-full mt-2"
                  value="85"
                  max="100"
                ></progress>
              </div>
            </div>

            <div className="card bg-base-200">
              <div className="card-body p-4">
                <h4 className="text-sm font-medium mb-2">Конверсия</h4>
                <div className="text-2xl font-bold text-accent">12%</div>
                <p className="text-xs text-base-content/70">
                  вероятность продажи
                </p>
                <div className="text-xs mt-2">
                  <span className="text-success">↑ 3%</span> после оптимизации
                </div>
              </div>
            </div>
          </div>

          {/* AI Content Generator */}
          <div className="card bg-base-200">
            <div className="card-body">
              <h4 className="font-bold mb-3">AI-генератор контента</h4>

              <div className="tabs tabs-boxed mb-4">
                <a className="tab tab-active">Описание</a>
                <a className="tab">Теги</a>
                <a className="tab">SEO</a>
              </div>

              <div className="space-y-3">
                <div className="form-control">
                  <label className="label">
                    <span className="label-text">Тон описания</span>
                  </label>
                  <select className="select select-bordered w-full">
                    <option>Профессиональный</option>
                    <option>Дружелюбный</option>
                    <option>Эмоциональный</option>
                    <option>Технический</option>
                  </select>
                </div>

                <div className="form-control">
                  <label className="label">
                    <span className="label-text">Акцент на</span>
                  </label>
                  <div className="flex flex-wrap gap-2">
                    <label className="cursor-pointer">
                      <input
                        type="checkbox"
                        className="checkbox checkbox-primary checkbox-sm"
                        checked
                      />
                      <span className="ml-2 text-sm">Экономичность</span>
                    </label>
                    <label className="cursor-pointer">
                      <input
                        type="checkbox"
                        className="checkbox checkbox-primary checkbox-sm"
                      />
                      <span className="ml-2 text-sm">Комфорт</span>
                    </label>
                    <label className="cursor-pointer">
                      <input
                        type="checkbox"
                        className="checkbox checkbox-primary checkbox-sm"
                        checked
                      />
                      <span className="ml-2 text-sm">Надежность</span>
                    </label>
                    <label className="cursor-pointer">
                      <input
                        type="checkbox"
                        className="checkbox checkbox-primary checkbox-sm"
                      />
                      <span className="ml-2 text-sm">Вместительность</span>
                    </label>
                  </div>
                </div>

                <button className="btn btn-primary w-full">
                  Сгенерировать описание
                </button>
              </div>
            </div>
          </div>

          {/* Multivariate Testing */}
          <div className="card bg-gradient-to-r from-purple-500/10 to-pink-500/10 border-2 border-purple-500/20">
            <div className="card-body">
              <h4 className="font-bold mb-3">A/B/C тестирование активно</h4>
              <div className="grid grid-cols-3 gap-2">
                <div className="text-center p-2 bg-base-100 rounded">
                  <div className="text-sm font-medium">Вариант A</div>
                  <div className="text-2xl font-bold text-primary">32%</div>
                  <div className="text-xs">конверсия</div>
                </div>
                <div className="text-center p-2 bg-base-100 rounded">
                  <div className="text-sm font-medium">Вариант B</div>
                  <div className="text-2xl font-bold text-secondary">28%</div>
                  <div className="text-xs">конверсия</div>
                </div>
                <div className="text-center p-2 bg-base-100 rounded">
                  <div className="text-sm font-medium">Вариант C</div>
                  <div className="text-2xl font-bold text-accent">41%</div>
                  <div className="text-xs">конверсия</div>
                </div>
              </div>
              <p className="text-sm mt-3">
                AI автоматически показывает разные версии вашего объявления и
                выбирает самую эффективную
              </p>
            </div>
          </div>
        </div>

        {/* Smart Actions */}
        <div className="flex gap-4 mt-6">
          <button className="btn btn-ghost flex-1">
            <span className="loading loading-spinner loading-xs mr-2"></span>
            AI оптимизирует...
          </button>
          <button className="btn btn-primary flex-1">
            Применить все рекомендации
          </button>
        </div>
      </div>
    </AnimatedSection>
  );

  return (
    <div className="container mx-auto p-4 max-w-6xl">
      <AnimatedSection animation="fadeIn">
        <div className="mb-8">
          <Link href="/examples" className="btn btn-ghost btn-sm mb-4">
            ← Назад к примерам
          </Link>
          <h1 className="text-4xl font-bold mb-4">
            Революционные подходы к редактированию объявлений
          </h1>
          <p className="text-lg text-base-content/70">
            От простого редактора до AI-powered системы с автоматической
            оптимизацией
          </p>
        </div>
      </AnimatedSection>

      {/* Example Selector */}
      <AnimatedSection animation="slideUp" delay={0.2}>
        <div className="flex flex-wrap gap-4 mb-8">
          <button
            onClick={() => setActiveExample('modern')}
            className={`btn ${activeExample === 'modern' ? 'btn-primary' : 'btn-outline'}`}
          >
            ✨ Современный
          </button>
          <button
            onClick={() => setActiveExample('smart')}
            className={`btn ${activeExample === 'smart' ? 'btn-secondary' : 'btn-outline'}`}
          >
            🧠 Умный
          </button>
          <button
            onClick={() => setActiveExample('ai-powered')}
            className={`btn ${activeExample === 'ai-powered' ? 'btn-accent' : 'btn-outline'}`}
          >
            🤖 AI-Powered
          </button>
        </div>
      </AnimatedSection>

      {/* Active Example */}
      <div className="mb-12">
        {activeExample === 'modern' && renderModernExample()}
        {activeExample === 'smart' && renderSmartExample()}
        {activeExample === 'ai-powered' && renderAIPoweredExample()}
      </div>

      {/* Feature Comparison */}
      <AnimatedSection animation="fadeIn" delay={0.4}>
        <div className="bg-base-200 rounded-xl p-6">
          <h2 className="text-2xl font-bold mb-6">Сравнение возможностей</h2>
          <div className="overflow-x-auto">
            <table className="table">
              <thead>
                <tr>
                  <th>Функция</th>
                  <th>Современный</th>
                  <th>Умный</th>
                  <th>AI-Powered</th>
                </tr>
              </thead>
              <tbody>
                <tr>
                  <td>Базовое редактирование</td>
                  <td>✅</td>
                  <td>✅</td>
                  <td>✅</td>
                </tr>
                <tr>
                  <td>Предпросмотр в реальном времени</td>
                  <td>✅</td>
                  <td>✅</td>
                  <td>✅</td>
                </tr>
                <tr>
                  <td>SEO оптимизация</td>
                  <td>✅</td>
                  <td>✅</td>
                  <td>✅</td>
                </tr>
                <tr>
                  <td>Умные подсказки</td>
                  <td>❌</td>
                  <td>✅</td>
                  <td>✅</td>
                </tr>
                <tr>
                  <td>Анализ рынка</td>
                  <td>❌</td>
                  <td>✅</td>
                  <td>✅</td>
                </tr>
                <tr>
                  <td>Импорт данных</td>
                  <td>❌</td>
                  <td>✅</td>
                  <td>✅</td>
                </tr>
                <tr>
                  <td>AI генерация контента</td>
                  <td>❌</td>
                  <td>❌</td>
                  <td>✅</td>
                </tr>
                <tr>
                  <td>A/B тестирование</td>
                  <td>❌</td>
                  <td>❌</td>
                  <td>✅</td>
                </tr>
                <tr>
                  <td>Автоматическая оптимизация</td>
                  <td>❌</td>
                  <td>❌</td>
                  <td>✅</td>
                </tr>
              </tbody>
            </table>
          </div>
        </div>
      </AnimatedSection>
    </div>
  );
}
