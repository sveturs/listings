'use client';

import React, { useState } from 'react';
import { SveTuLogoStatic } from '@/components/logos/SveTuLogoStatic';
import { AnimatedSection } from '@/components/ui/AnimatedSection';

const StorefrontDashboard = () => {
  const [activeTab, setActiveTab] = useState<'overview' | 'products' | 'analytics' | 'settings'>('overview');
  const [timeRange, setTimeRange] = useState<'day' | 'week' | 'month'>('week');

  const stats = {
    revenue: { value: '€12,549', change: '+23%', trend: 'up' },
    orders: { value: '234', change: '+15%', trend: 'up' },
    visitors: { value: '3,451', change: '+8%', trend: 'up' },
    conversion: { value: '6.8%', change: '-2%', trend: 'down' },
  };

  const products = [
    {
      id: 1,
      name: 'iPhone 14 Pro Max 256GB',
      price: 899,
      stock: 5,
      sold: 23,
      image: '/api/minio/download?fileName=listings/0a47e66f-d8da-459f-a2ba-8e2b85ae0163/38ad29e6-7b07-4bfc-9db2-d965cb6b966f.jpg',
      status: 'active',
    },
    {
      id: 2,
      name: 'MacBook Pro M2 13"',
      price: 1299,
      stock: 3,
      sold: 12,
      image: '/api/minio/download?fileName=listings/0c91d2f7-53f7-4bff-87fe-d7e82dc3e2f0/3b26f07f-c5d6-4ff7-ba56-06ec69bb7f4d.jpg',
      status: 'active',
    },
    {
      id: 3,
      name: 'AirPods Pro 2',
      price: 249,
      stock: 0,
      sold: 45,
      image: '/api/minio/download?fileName=listings/0e17c3be-e76e-433a-a6d4-86bb8b7a0e29/23bb3da7-38ef-44f7-8c1d-1c14eaaafeb5.jpg',
      status: 'out_of_stock',
    },
  ];

  const chartData = [
    { day: 'Пн', sales: 1200 },
    { day: 'Вт', sales: 1900 },
    { day: 'Ср', sales: 1600 },
    { day: 'Чт', sales: 2200 },
    { day: 'Пт', sales: 2800 },
    { day: 'Сб', sales: 3200 },
    { day: 'Вс', sales: 2400 },
  ];

  return (
    <div className="min-h-screen bg-gradient-to-br from-base-100 to-base-200">
      {/* Header */}
      <div className="navbar bg-base-100 shadow-lg">
        <div className="navbar-start">
          <SveTuLogoStatic variant="gradient" width={120} height={40} />
        </div>
        <div className="navbar-center">
          <h1 className="text-xl font-bold">🏪 TechStore - Панель управления</h1>
        </div>
        <div className="navbar-end">
          <div className="dropdown dropdown-end">
            <label tabIndex={0} className="btn btn-ghost btn-circle avatar">
              <div className="w-10 rounded-full">
                <img src="https://ui-avatars.com/api/?name=Tech+Store&background=6366f1&color=fff" alt="Store" />
              </div>
            </label>
          </div>
        </div>
      </div>

      <div className="container mx-auto px-4 py-6 max-w-7xl">
        {/* Tabs */}
        <AnimatedSection animation="fadeIn">
          <div className="tabs tabs-boxed mb-6">
            <a 
              className={`tab ${activeTab === 'overview' ? 'tab-active' : ''}`}
              onClick={() => setActiveTab('overview')}
            >
              📊 Обзор
            </a>
            <a 
              className={`tab ${activeTab === 'products' ? 'tab-active' : ''}`}
              onClick={() => setActiveTab('products')}
            >
              📦 Товары
            </a>
            <a 
              className={`tab ${activeTab === 'analytics' ? 'tab-active' : ''}`}
              onClick={() => setActiveTab('analytics')}
            >
              📈 Аналитика
            </a>
            <a 
              className={`tab ${activeTab === 'settings' ? 'tab-active' : ''}`}
              onClick={() => setActiveTab('settings')}
            >
              ⚙️ Настройки
            </a>
          </div>
        </AnimatedSection>

        {/* Overview Tab */}
        {activeTab === 'overview' && (
          <>
            {/* Stats Grid */}
            <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-4 gap-6 mb-8">
              <AnimatedSection animation="slideUp" delay={0}>
                <div className="card bg-base-100 shadow-xl">
                  <div className="card-body">
                    <div className="flex justify-between items-start">
                      <div>
                        <p className="text-sm text-base-content/60">Выручка</p>
                        <p className="text-3xl font-bold">{stats.revenue.value}</p>
                        <p className={`text-sm ${stats.revenue.trend === 'up' ? 'text-success' : 'text-error'}`}>
                          {stats.revenue.change} к прошлой неделе
                        </p>
                      </div>
                      <div className="text-3xl">💰</div>
                    </div>
                  </div>
                </div>
              </AnimatedSection>

              <AnimatedSection animation="slideUp" delay={0.1}>
                <div className="card bg-base-100 shadow-xl">
                  <div className="card-body">
                    <div className="flex justify-between items-start">
                      <div>
                        <p className="text-sm text-base-content/60">Заказы</p>
                        <p className="text-3xl font-bold">{stats.orders.value}</p>
                        <p className={`text-sm ${stats.orders.trend === 'up' ? 'text-success' : 'text-error'}`}>
                          {stats.orders.change} к прошлой неделе
                        </p>
                      </div>
                      <div className="text-3xl">📦</div>
                    </div>
                  </div>
                </div>
              </AnimatedSection>

              <AnimatedSection animation="slideUp" delay={0.2}>
                <div className="card bg-base-100 shadow-xl">
                  <div className="card-body">
                    <div className="flex justify-between items-start">
                      <div>
                        <p className="text-sm text-base-content/60">Посетители</p>
                        <p className="text-3xl font-bold">{stats.visitors.value}</p>
                        <p className={`text-sm ${stats.visitors.trend === 'up' ? 'text-success' : 'text-error'}`}>
                          {stats.visitors.change} к прошлой неделе
                        </p>
                      </div>
                      <div className="text-3xl">👥</div>
                    </div>
                  </div>
                </div>
              </AnimatedSection>

              <AnimatedSection animation="slideUp" delay={0.3}>
                <div className="card bg-base-100 shadow-xl">
                  <div className="card-body">
                    <div className="flex justify-between items-start">
                      <div>
                        <p className="text-sm text-base-content/60">Конверсия</p>
                        <p className="text-3xl font-bold">{stats.conversion.value}</p>
                        <p className={`text-sm ${stats.conversion.trend === 'up' ? 'text-success' : 'text-error'}`}>
                          {stats.conversion.change} к прошлой неделе
                        </p>
                      </div>
                      <div className="text-3xl">🎯</div>
                    </div>
                  </div>
                </div>
              </AnimatedSection>
            </div>

            {/* Chart */}
            <AnimatedSection animation="fadeIn" delay={0.4}>
              <div className="card bg-base-100 shadow-xl mb-8">
                <div className="card-body">
                  <div className="flex justify-between items-center mb-4">
                    <h3 className="text-xl font-bold">График продаж</h3>
                    <div className="btn-group">
                      <button 
                        className={`btn btn-sm ${timeRange === 'day' ? 'btn-active' : ''}`}
                        onClick={() => setTimeRange('day')}
                      >
                        День
                      </button>
                      <button 
                        className={`btn btn-sm ${timeRange === 'week' ? 'btn-active' : ''}`}
                        onClick={() => setTimeRange('week')}
                      >
                        Неделя
                      </button>
                      <button 
                        className={`btn btn-sm ${timeRange === 'month' ? 'btn-active' : ''}`}
                        onClick={() => setTimeRange('month')}
                      >
                        Месяц
                      </button>
                    </div>
                  </div>
                  <div className="h-64 flex items-end justify-between gap-2">
                    {chartData.map((data, idx) => (
                      <div key={idx} className="flex-1 flex flex-col items-center">
                        <div 
                          className="w-full bg-primary rounded-t transition-all hover:bg-primary-focus"
                          style={{ height: `${(data.sales / 3200) * 100}%` }}
                        ></div>
                        <span className="text-xs mt-2">{data.day}</span>
                      </div>
                    ))}
                  </div>
                </div>
              </div>
            </AnimatedSection>

            {/* Recent Orders */}
            <AnimatedSection animation="slideUp" delay={0.5}>
              <div className="card bg-base-100 shadow-xl">
                <div className="card-body">
                  <h3 className="text-xl font-bold mb-4">Последние заказы</h3>
                  <div className="overflow-x-auto">
                    <table className="table">
                      <thead>
                        <tr>
                          <th>№ Заказа</th>
                          <th>Клиент</th>
                          <th>Товар</th>
                          <th>Сумма</th>
                          <th>Статус</th>
                          <th>Действия</th>
                        </tr>
                      </thead>
                      <tbody>
                        <tr>
                          <td>#2341</td>
                          <td>Иван Петров</td>
                          <td>iPhone 14 Pro Max</td>
                          <td>€899</td>
                          <td><span className="badge badge-warning">В обработке</span></td>
                          <td><button className="btn btn-ghost btn-xs">Детали</button></td>
                        </tr>
                        <tr>
                          <td>#2340</td>
                          <td>Анна Смирнова</td>
                          <td>MacBook Pro M2</td>
                          <td>€1299</td>
                          <td><span className="badge badge-success">Доставлен</span></td>
                          <td><button className="btn btn-ghost btn-xs">Детали</button></td>
                        </tr>
                        <tr>
                          <td>#2339</td>
                          <td>Петр Сидоров</td>
                          <td>AirPods Pro 2</td>
                          <td>€249</td>
                          <td><span className="badge badge-info">Отправлен</span></td>
                          <td><button className="btn btn-ghost btn-xs">Детали</button></td>
                        </tr>
                      </tbody>
                    </table>
                  </div>
                </div>
              </div>
            </AnimatedSection>
          </>
        )}

        {/* Products Tab */}
        {activeTab === 'products' && (
          <AnimatedSection animation="fadeIn">
            <div className="card bg-base-100 shadow-xl">
              <div className="card-body">
                <div className="flex justify-between items-center mb-4">
                  <h3 className="text-xl font-bold">Управление товарами</h3>
                  <button className="btn btn-primary">
                    <svg className="w-5 h-5 mr-2" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                      <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M12 6v6m0 0v6m0-6h6m-6 0H6" />
                    </svg>
                    Добавить товар
                  </button>
                </div>
                <div className="overflow-x-auto">
                  <table className="table">
                    <thead>
                      <tr>
                        <th>Товар</th>
                        <th>Цена</th>
                        <th>Остаток</th>
                        <th>Продано</th>
                        <th>Статус</th>
                        <th>Действия</th>
                      </tr>
                    </thead>
                    <tbody>
                      {products.map((product) => (
                        <tr key={product.id}>
                          <td>
                            <div className="flex items-center gap-3">
                              <div className="avatar">
                                <div className="mask mask-squircle w-12 h-12">
                                  <img src={product.image} alt={product.name} />
                                </div>
                              </div>
                              <div>
                                <div className="font-bold">{product.name}</div>
                                <div className="text-sm opacity-50">ID: {product.id}</div>
                              </div>
                            </div>
                          </td>
                          <td>€{product.price}</td>
                          <td className={product.stock === 0 ? 'text-error' : ''}>{product.stock}</td>
                          <td>{product.sold}</td>
                          <td>
                            <span className={`badge ${product.status === 'active' ? 'badge-success' : 'badge-error'}`}>
                              {product.status === 'active' ? 'Активен' : 'Нет в наличии'}
                            </span>
                          </td>
                          <td>
                            <div className="flex gap-2">
                              <button className="btn btn-ghost btn-xs">Изменить</button>
                              <button className="btn btn-ghost btn-xs text-error">Удалить</button>
                            </div>
                          </td>
                        </tr>
                      ))}
                    </tbody>
                  </table>
                </div>
              </div>
            </div>
          </AnimatedSection>
        )}

        {/* Analytics Tab */}
        {activeTab === 'analytics' && (
          <div className="grid grid-cols-1 lg:grid-cols-2 gap-6">
            <AnimatedSection animation="slideLeft">
              <div className="card bg-base-100 shadow-xl">
                <div className="card-body">
                  <h3 className="text-xl font-bold mb-4">Топ товары</h3>
                  <div className="space-y-4">
                    {products.map((product, idx) => (
                      <div key={product.id} className="flex items-center gap-4">
                        <div className="text-2xl font-bold text-base-content/30">#{idx + 1}</div>
                        <div className="avatar">
                          <div className="w-12 rounded">
                            <img src={product.image} alt={product.name} />
                          </div>
                        </div>
                        <div className="flex-1">
                          <div className="font-semibold">{product.name}</div>
                          <div className="text-sm text-base-content/60">{product.sold} продаж</div>
                        </div>
                        <div className="text-lg font-bold">€{product.price * product.sold}</div>
                      </div>
                    ))}
                  </div>
                </div>
              </div>
            </AnimatedSection>

            <AnimatedSection animation="slideRight">
              <div className="card bg-base-100 shadow-xl">
                <div className="card-body">
                  <h3 className="text-xl font-bold mb-4">Источники трафика</h3>
                  <div className="space-y-4">
                    <div>
                      <div className="flex justify-between mb-1">
                        <span>Поиск Google</span>
                        <span className="font-semibold">45%</span>
                      </div>
                      <progress className="progress progress-primary" value="45" max="100"></progress>
                    </div>
                    <div>
                      <div className="flex justify-between mb-1">
                        <span>Прямые заходы</span>
                        <span className="font-semibold">30%</span>
                      </div>
                      <progress className="progress progress-secondary" value="30" max="100"></progress>
                    </div>
                    <div>
                      <div className="flex justify-between mb-1">
                        <span>Социальные сети</span>
                        <span className="font-semibold">20%</span>
                      </div>
                      <progress className="progress progress-accent" value="20" max="100"></progress>
                    </div>
                    <div>
                      <div className="flex justify-between mb-1">
                        <span>Email рассылка</span>
                        <span className="font-semibold">5%</span>
                      </div>
                      <progress className="progress progress-info" value="5" max="100"></progress>
                    </div>
                  </div>
                </div>
              </div>
            </AnimatedSection>

            <AnimatedSection animation="slideUp" delay={0.2} className="lg:col-span-2">
              <div className="card bg-base-100 shadow-xl">
                <div className="card-body">
                  <h3 className="text-xl font-bold mb-4">Демография покупателей</h3>
                  <div className="grid grid-cols-1 md:grid-cols-3 gap-6">
                    <div className="text-center">
                      <div className="radial-progress text-primary" style={{"--value": 68, "--size": "8rem"} as any}>68%</div>
                      <p className="mt-2 font-semibold">Мужчины</p>
                    </div>
                    <div className="text-center">
                      <div className="radial-progress text-secondary" style={{"--value": 32, "--size": "8rem"} as any}>32%</div>
                      <p className="mt-2 font-semibold">Женщины</p>
                    </div>
                    <div className="space-y-2">
                      <h4 className="font-semibold">Возрастные группы:</h4>
                      <div className="text-sm space-y-1">
                        <div>18-24: 15%</div>
                        <div>25-34: 35%</div>
                        <div>35-44: 30%</div>
                        <div>45+: 20%</div>
                      </div>
                    </div>
                  </div>
                </div>
              </div>
            </AnimatedSection>
          </div>
        )}

        {/* Settings Tab */}
        {activeTab === 'settings' && (
          <AnimatedSection animation="fadeIn">
            <div className="grid grid-cols-1 lg:grid-cols-2 gap-6">
              <div className="card bg-base-100 shadow-xl">
                <div className="card-body">
                  <h3 className="text-xl font-bold mb-4">Информация о магазине</h3>
                  <div className="space-y-4">
                    <div>
                      <label className="label">
                        <span className="label-text">Название магазина</span>
                      </label>
                      <input type="text" className="input input-bordered w-full" defaultValue="TechStore" />
                    </div>
                    <div>
                      <label className="label">
                        <span className="label-text">Описание</span>
                      </label>
                      <textarea className="textarea textarea-bordered w-full" rows={3} defaultValue="Лучшая электроника по доступным ценам" />
                    </div>
                    <div>
                      <label className="label">
                        <span className="label-text">Контактный email</span>
                      </label>
                      <input type="email" className="input input-bordered w-full" defaultValue="info@techstore.rs" />
                    </div>
                    <div>
                      <label className="label">
                        <span className="label-text">Телефон</span>
                      </label>
                      <input type="tel" className="input input-bordered w-full" defaultValue="+381 11 123 4567" />
                    </div>
                    <button className="btn btn-primary">Сохранить изменения</button>
                  </div>
                </div>
              </div>

              <div className="space-y-6">
                <div className="card bg-base-100 shadow-xl">
                  <div className="card-body">
                    <h3 className="text-xl font-bold mb-4">Настройки витрины</h3>
                    <div className="space-y-4">
                      <label className="flex items-center justify-between cursor-pointer">
                        <span>Показывать рейтинги товаров</span>
                        <input type="checkbox" className="toggle toggle-primary" defaultChecked />
                      </label>
                      <label className="flex items-center justify-between cursor-pointer">
                        <span>Включить отзывы</span>
                        <input type="checkbox" className="toggle toggle-primary" defaultChecked />
                      </label>
                      <label className="flex items-center justify-between cursor-pointer">
                        <span>Автоматическая модерация</span>
                        <input type="checkbox" className="toggle toggle-primary" />
                      </label>
                      <label className="flex items-center justify-between cursor-pointer">
                        <span>Уведомления о заказах</span>
                        <input type="checkbox" className="toggle toggle-primary" defaultChecked />
                      </label>
                    </div>
                  </div>
                </div>

                <div className="card bg-warning/20 border border-warning">
                  <div className="card-body">
                    <h3 className="text-xl font-bold mb-2">⚡ Premium функции</h3>
                    <p className="text-sm mb-4">Разблокируйте дополнительные возможности для вашего бизнеса</p>
                    <ul className="space-y-2 text-sm">
                      <li className="flex items-center gap-2">
                        <svg className="w-5 h-5 text-success" fill="currentColor" viewBox="0 0 20 20">
                          <path fillRule="evenodd" d="M10 18a8 8 0 100-16 8 8 0 000 16zm3.707-9.293a1 1 0 00-1.414-1.414L9 10.586 7.707 9.293a1 1 0 00-1.414 1.414l2 2a1 1 0 001.414 0l4-4z" clipRule="evenodd" />
                        </svg>
                        Расширенная аналитика
                      </li>
                      <li className="flex items-center gap-2">
                        <svg className="w-5 h-5 text-success" fill="currentColor" viewBox="0 0 20 20">
                          <path fillRule="evenodd" d="M10 18a8 8 0 100-16 8 8 0 000 16zm3.707-9.293a1 1 0 00-1.414-1.414L9 10.586 7.707 9.293a1 1 0 00-1.414 1.414l2 2a1 1 0 001.414 0l4-4z" clipRule="evenodd" />
                        </svg>
                        Email маркетинг
                      </li>
                      <li className="flex items-center gap-2">
                        <svg className="w-5 h-5 text-success" fill="currentColor" viewBox="0 0 20 20">
                          <path fillRule="evenodd" d="M10 18a8 8 0 100-16 8 8 0 000 16zm3.707-9.293a1 1 0 00-1.414-1.414L9 10.586 7.707 9.293a1 1 0 00-1.414 1.414l2 2a1 1 0 001.414 0l4-4z" clipRule="evenodd" />
                        </svg>
                        Приоритетная поддержка
                      </li>
                    </ul>
                    <button className="btn btn-warning mt-4">Обновить до Premium</button>
                  </div>
                </div>
              </div>
            </div>
          </AnimatedSection>
        )}
      </div>
    </div>
  );
};

export default StorefrontDashboard;