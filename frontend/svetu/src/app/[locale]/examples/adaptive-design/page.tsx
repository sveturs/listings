'use client';

import React, { useState } from 'react';
import { SveTuLogoStatic } from '@/components/logos/SveTuLogoStatic';
import { AnimatedSection } from '@/components/ui/AnimatedSection';

const AdaptiveDesign = () => {
  const [deviceView, setDeviceView] = useState<'mobile' | 'tablet' | 'desktop'>('desktop');
  const [showGrid, setShowGrid] = useState(false);

  const deviceSizes = {
    mobile: { width: '375px', height: '667px', scale: 0.8 },
    tablet: { width: '768px', height: '1024px', scale: 0.6 },
    desktop: { width: '1440px', height: '900px', scale: 0.4 },
  };

  const demoProducts = [
    {
      id: 1,
      title: 'iPhone 14 Pro Max',
      price: 899,
      location: 'Белград',
      image: '/api/minio/download?fileName=listings/0a47e66f-d8da-459f-a2ba-8e2b85ae0163/38ad29e6-7b07-4bfc-9db2-d965cb6b966f.jpg',
    },
    {
      id: 2,
      title: 'MacBook Pro M2',
      price: 1299,
      location: 'Нови Сад',
      image: '/api/minio/download?fileName=listings/0c91d2f7-53f7-4bff-87fe-d7e82dc3e2f0/3b26f07f-c5d6-4ff7-ba56-06ec69bb7f4d.jpg',
    },
    {
      id: 3,
      title: 'AirPods Pro 2',
      price: 249,
      location: 'Ниш',
      image: '/api/minio/download?fileName=listings/0e17c3be-e76e-433a-a6d4-86bb8b7a0e29/23bb3da7-38ef-44f7-8c1d-1c14eaaafeb5.jpg',
    },
  ];

  return (
    <div className="min-h-screen bg-gradient-to-br from-base-100 to-base-200">
      {/* Header */}
      <div className="navbar bg-base-100 shadow-lg">
        <div className="navbar-start">
          <SveTuLogoStatic variant="gradient" width={120} height={40} />
        </div>
        <div className="navbar-center">
          <h1 className="text-xl font-bold">📱 Адаптивный дизайн</h1>
        </div>
        <div className="navbar-end">
          <label className="swap swap-rotate">
            <input 
              type="checkbox" 
              checked={showGrid}
              onChange={(e) => setShowGrid(e.target.checked)}
            />
            <svg className="swap-off fill-current w-6 h-6" xmlns="http://www.w3.org/2000/svg" viewBox="0 0 24 24">
              <path d="M3 4v16h18V4H3zm16 14H5V6h14v12z"/>
            </svg>
            <svg className="swap-on fill-current w-6 h-6" xmlns="http://www.w3.org/2000/svg" viewBox="0 0 24 24">
              <path d="M3 4v16h18V4H3zm8 14H5V6h6v12zm8 0h-6V6h6v12z"/>
            </svg>
          </label>
        </div>
      </div>

      <div className="container mx-auto px-4 py-8">
        {/* Device Selector */}
        <AnimatedSection animation="fadeIn">
          <div className="flex justify-center mb-8">
            <div className="join">
              <button 
                className={`join-item btn ${deviceView === 'mobile' ? 'btn-active' : ''}`}
                onClick={() => setDeviceView('mobile')}
              >
                <svg className="w-5 h-5 mr-2" fill="currentColor" viewBox="0 0 20 20">
                  <path d="M7 2a2 2 0 00-2 2v12a2 2 0 002 2h6a2 2 0 002-2V4a2 2 0 00-2-2H7zM9 14a1 1 0 100 2 1 1 0 000-2z" />
                </svg>
                Mobile
              </button>
              <button 
                className={`join-item btn ${deviceView === 'tablet' ? 'btn-active' : ''}`}
                onClick={() => setDeviceView('tablet')}
              >
                <svg className="w-5 h-5 mr-2" fill="currentColor" viewBox="0 0 20 20">
                  <path d="M5 2a2 2 0 00-2 2v12a2 2 0 002 2h10a2 2 0 002-2V4a2 2 0 00-2-2H5zm5 14a1 1 0 100 2 1 1 0 000-2z" />
                </svg>
                Tablet
              </button>
              <button 
                className={`join-item btn ${deviceView === 'desktop' ? 'btn-active' : ''}`}
                onClick={() => setDeviceView('desktop')}
              >
                <svg className="w-5 h-5 mr-2" fill="currentColor" viewBox="0 0 20 20">
                  <path fillRule="evenodd" d="M3 5a2 2 0 012-2h10a2 2 0 012 2v8a2 2 0 01-2 2h-2.22l.123.489.804.804A1 1 0 0113 18H7a1 1 0 01-.707-1.707l.804-.804L7.22 15H5a2 2 0 01-2-2V5zm5.771 7H5V5h10v7H8.771z" clipRule="evenodd" />
                </svg>
                Desktop
              </button>
            </div>
          </div>
        </AnimatedSection>

        {/* Device Preview */}
        <AnimatedSection animation="zoomIn">
          <div className="flex justify-center mb-8">
            <div 
              className={`relative bg-gray-900 rounded-3xl p-2 shadow-2xl transition-all duration-500 ${
                deviceView === 'mobile' ? 'w-[400px]' : 
                deviceView === 'tablet' ? 'w-[820px]' : 
                'w-full max-w-6xl'
              }`}
            >
              {/* Device Frame */}
              {deviceView === 'mobile' && (
                <>
                  <div className="absolute top-6 left-1/2 transform -translate-x-1/2 w-32 h-6 bg-black rounded-full"></div>
                  <div className="absolute bottom-2 left-1/2 transform -translate-x-1/2 w-32 h-1 bg-gray-700 rounded-full"></div>
                </>
              )}
              
              {/* Screen */}
              <div 
                className={`bg-white rounded-2xl overflow-hidden ${
                  deviceView === 'mobile' ? 'h-[667px]' : 
                  deviceView === 'tablet' ? 'h-[600px]' : 
                  'h-[500px]'
                }`}
                style={{
                  width: deviceView === 'desktop' ? '100%' : deviceSizes[deviceView].width,
                }}
              >
                {/* Demo App */}
                <div className="h-full overflow-y-auto">
                  {/* App Header */}
                  <div className="navbar bg-primary text-primary-content sticky top-0 z-10">
                    <div className="navbar-start">
                      {deviceView === 'mobile' && (
                        <button className="btn btn-ghost btn-circle">
                          <svg className="w-6 h-6" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                            <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M4 6h16M4 12h16M4 18h16" />
                          </svg>
                        </button>
                      )}
                      <SveTuLogoStatic variant="minimal" width={80} height={30} />
                    </div>
                    <div className="navbar-center">
                      {deviceView !== 'mobile' && (
                        <div className="form-control">
                          <input 
                            type="text" 
                            placeholder="Поиск..." 
                            className="input input-bordered input-sm w-64"
                          />
                        </div>
                      )}
                    </div>
                    <div className="navbar-end">
                      <button className="btn btn-ghost btn-circle">
                        <svg className="w-6 h-6" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                          <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M3 3h2l.4 2M7 13h10l4-8H5.4M7 13L5.4 5M7 13l-2.293 2.293c-.63.63-.184 1.707.707 1.707H17m0 0a2 2 0 100 4 2 2 0 000-4zm-8 2a2 2 0 11-4 0 2 2 0 014 0z" />
                        </svg>
                      </button>
                    </div>
                  </div>

                  {/* Search Bar for Mobile */}
                  {deviceView === 'mobile' && (
                    <div className="p-4 bg-base-100">
                      <input 
                        type="text" 
                        placeholder="Поиск товаров..." 
                        className="input input-bordered w-full"
                      />
                    </div>
                  )}

                  {/* Categories */}
                  <div className="p-4 bg-base-100">
                    <h3 className="font-bold mb-3">Категории</h3>
                    <div className={`grid gap-3 ${
                      deviceView === 'mobile' ? 'grid-cols-3' : 
                      deviceView === 'tablet' ? 'grid-cols-4' : 
                      'grid-cols-6'
                    }`}>
                      {['🏠', '🚗', '💻', '👕', '🎮', '📱'].map((emoji, idx) => (
                        <button key={idx} className="btn btn-outline btn-sm">
                          <span className="text-2xl">{emoji}</span>
                        </button>
                      ))}
                    </div>
                  </div>

                  {/* Products Grid */}
                  <div className="p-4">
                    <h3 className="font-bold mb-3">Популярные товары</h3>
                    <div className={`grid gap-4 ${
                      deviceView === 'mobile' ? 'grid-cols-1' : 
                      deviceView === 'tablet' ? 'grid-cols-2' : 
                      'grid-cols-3'
                    }`}>
                      {demoProducts.map((product) => (
                        <div key={product.id} className="card bg-base-100 shadow-xl">
                          <figure className={`${deviceView === 'mobile' ? 'h-48' : 'h-32'}`}>
                            <img 
                              src={product.image} 
                              alt={product.title}
                              className="w-full h-full object-cover"
                            />
                          </figure>
                          <div className="card-body p-4">
                            <h4 className={`font-bold ${deviceView === 'mobile' ? 'text-lg' : 'text-sm'}`}>
                              {product.title}
                            </h4>
                            <p className="text-xs text-base-content/60">{product.location}</p>
                            <div className="flex justify-between items-center mt-2">
                              <span className={`font-bold text-primary ${deviceView === 'mobile' ? 'text-xl' : 'text-lg'}`}>
                                €{product.price}
                              </span>
                              <button className="btn btn-primary btn-sm">Купить</button>
                            </div>
                          </div>
                        </div>
                      ))}
                    </div>
                  </div>

                  {/* Bottom Navigation for Mobile */}
                  {deviceView === 'mobile' && (
                    <div className="btm-nav">
                      <button className="active">
                        <svg className="w-5 h-5" fill="currentColor" viewBox="0 0 20 20">
                          <path d="M10.707 2.293a1 1 0 00-1.414 0l-7 7a1 1 0 001.414 1.414L4 10.414V17a1 1 0 001 1h2a1 1 0 001-1v-2a1 1 0 011-1h2a1 1 0 011 1v2a1 1 0 001 1h2a1 1 0 001-1v-6.586l.293.293a1 1 0 001.414-1.414l-7-7z" />
                        </svg>
                      </button>
                      <button>
                        <svg className="w-5 h-5" fill="currentColor" viewBox="0 0 20 20">
                          <path fillRule="evenodd" d="M8 4a4 4 0 100 8 4 4 0 000-8zM2 8a6 6 0 1110.89 3.476l4.817 4.817a1 1 0 01-1.414 1.414l-4.816-4.816A6 6 0 012 8z" clipRule="evenodd" />
                        </svg>
                      </button>
                      <button>
                        <svg className="w-5 h-5" fill="currentColor" viewBox="0 0 20 20">
                          <path d="M10 2a6 6 0 00-6 6v3.586l-.707.707A1 1 0 004 14h12a1 1 0 00.707-1.707L16 11.586V8a6 6 0 00-6-6zM10 18a3 3 0 01-3-3h6a3 3 0 01-3 3z" />
                        </svg>
                      </button>
                      <button>
                        <svg className="w-5 h-5" fill="currentColor" viewBox="0 0 20 20">
                          <path fillRule="evenodd" d="M10 9a3 3 0 100-6 3 3 0 000 6zm-7 9a7 7 0 1114 0H3z" clipRule="evenodd" />
                        </svg>
                      </button>
                    </div>
                  )}
                </div>
              </div>
            </div>
          </div>
        </AnimatedSection>

        {/* Features Grid */}
        <AnimatedSection animation="slideUp" delay={0.2}>
          <div className="grid grid-cols-1 md:grid-cols-3 gap-6 mt-12">
            <div className="card bg-base-100 shadow-xl">
              <div className="card-body text-center">
                <div className="text-4xl mb-4">📱</div>
                <h3 className="card-title justify-center">Mobile First</h3>
                <p className="text-sm">Оптимизация для мобильных устройств с приоритетом UX</p>
                <div className="mt-4 space-y-2 text-left">
                  <div className="flex items-center gap-2">
                    <svg className="w-5 h-5 text-success" fill="currentColor" viewBox="0 0 20 20">
                      <path fillRule="evenodd" d="M10 18a8 8 0 100-16 8 8 0 000 16zm3.707-9.293a1 1 0 00-1.414-1.414L9 10.586 7.707 9.293a1 1 0 00-1.414 1.414l2 2a1 1 0 001.414 0l4-4z" clipRule="evenodd" />
                    </svg>
                    <span className="text-sm">Touch-friendly интерфейс</span>
                  </div>
                  <div className="flex items-center gap-2">
                    <svg className="w-5 h-5 text-success" fill="currentColor" viewBox="0 0 20 20">
                      <path fillRule="evenodd" d="M10 18a8 8 0 100-16 8 8 0 000 16zm3.707-9.293a1 1 0 00-1.414-1.414L9 10.586 7.707 9.293a1 1 0 00-1.414 1.414l2 2a1 1 0 001.414 0l4-4z" clipRule="evenodd" />
                    </svg>
                    <span className="text-sm">Свайпы и жесты</span>
                  </div>
                  <div className="flex items-center gap-2">
                    <svg className="w-5 h-5 text-success" fill="currentColor" viewBox="0 0 20 20">
                      <path fillRule="evenodd" d="M10 18a8 8 0 100-16 8 8 0 000 16zm3.707-9.293a1 1 0 00-1.414-1.414L9 10.586 7.707 9.293a1 1 0 00-1.414 1.414l2 2a1 1 0 001.414 0l4-4z" clipRule="evenodd" />
                    </svg>
                    <span className="text-sm">Нижняя навигация</span>
                  </div>
                </div>
              </div>
            </div>

            <div className="card bg-base-100 shadow-xl">
              <div className="card-body text-center">
                <div className="text-4xl mb-4">💻</div>
                <h3 className="card-title justify-center">Responsive Grid</h3>
                <p className="text-sm">Гибкая сетка для всех размеров экранов</p>
                <div className="mt-4 space-y-2 text-left">
                  <div className="flex items-center gap-2">
                    <svg className="w-5 h-5 text-success" fill="currentColor" viewBox="0 0 20 20">
                      <path fillRule="evenodd" d="M10 18a8 8 0 100-16 8 8 0 000 16zm3.707-9.293a1 1 0 00-1.414-1.414L9 10.586 7.707 9.293a1 1 0 00-1.414 1.414l2 2a1 1 0 001.414 0l4-4z" clipRule="evenodd" />
                    </svg>
                    <span className="text-sm">12-колоночная система</span>
                  </div>
                  <div className="flex items-center gap-2">
                    <svg className="w-5 h-5 text-success" fill="currentColor" viewBox="0 0 20 20">
                      <path fillRule="evenodd" d="M10 18a8 8 0 100-16 8 8 0 000 16zm3.707-9.293a1 1 0 00-1.414-1.414L9 10.586 7.707 9.293a1 1 0 00-1.414 1.414l2 2a1 1 0 001.414 0l4-4z" clipRule="evenodd" />
                    </svg>
                    <span className="text-sm">Breakpoints: 640/768/1024/1280</span>
                  </div>
                  <div className="flex items-center gap-2">
                    <svg className="w-5 h-5 text-success" fill="currentColor" viewBox="0 0 20 20">
                      <path fillRule="evenodd" d="M10 18a8 8 0 100-16 8 8 0 000 16zm3.707-9.293a1 1 0 00-1.414-1.414L9 10.586 7.707 9.293a1 1 0 00-1.414 1.414l2 2a1 1 0 001.414 0l4-4z" clipRule="evenodd" />
                    </svg>
                    <span className="text-sm">Fluid typography</span>
                  </div>
                </div>
              </div>
            </div>

            <div className="card bg-base-100 shadow-xl">
              <div className="card-body text-center">
                <div className="text-4xl mb-4">⚡</div>
                <h3 className="card-title justify-center">Performance</h3>
                <p className="text-sm">Оптимизация для быстрой загрузки</p>
                <div className="mt-4 space-y-2 text-left">
                  <div className="flex items-center gap-2">
                    <svg className="w-5 h-5 text-success" fill="currentColor" viewBox="0 0 20 20">
                      <path fillRule="evenodd" d="M10 18a8 8 0 100-16 8 8 0 000 16zm3.707-9.293a1 1 0 00-1.414-1.414L9 10.586 7.707 9.293a1 1 0 00-1.414 1.414l2 2a1 1 0 001.414 0l4-4z" clipRule="evenodd" />
                    </svg>
                    <span className="text-sm">Lazy loading изображений</span>
                  </div>
                  <div className="flex items-center gap-2">
                    <svg className="w-5 h-5 text-success" fill="currentColor" viewBox="0 0 20 20">
                      <path fillRule="evenodd" d="M10 18a8 8 0 100-16 8 8 0 000 16zm3.707-9.293a1 1 0 00-1.414-1.414L9 10.586 7.707 9.293a1 1 0 00-1.414 1.414l2 2a1 1 0 001.414 0l4-4z" clipRule="evenodd" />
                    </svg>
                    <span className="text-sm">Минификация CSS/JS</span>
                  </div>
                  <div className="flex items-center gap-2">
                    <svg className="w-5 h-5 text-success" fill="currentColor" viewBox="0 0 20 20">
                      <path fillRule="evenodd" d="M10 18a8 8 0 100-16 8 8 0 000 16zm3.707-9.293a1 1 0 00-1.414-1.414L9 10.586 7.707 9.293a1 1 0 00-1.414 1.414l2 2a1 1 0 001.414 0l4-4z" clipRule="evenodd" />
                    </svg>
                    <span className="text-sm">PWA поддержка</span>
                  </div>
                </div>
              </div>
            </div>
          </div>
        </AnimatedSection>

        {/* Breakpoints Info */}
        <AnimatedSection animation="fadeIn" delay={0.4}>
          <div className="mt-12 card bg-base-100 shadow-xl">
            <div className="card-body">
              <h3 className="card-title mb-4">Breakpoints и медиа-запросы</h3>
              <div className="overflow-x-auto">
                <table className="table">
                  <thead>
                    <tr>
                      <th>Устройство</th>
                      <th>Breakpoint</th>
                      <th>Колонки</th>
                      <th>Отступы</th>
                      <th>Примеры</th>
                    </tr>
                  </thead>
                  <tbody>
                    <tr>
                      <td className="font-semibold">Mobile</td>
                      <td>&lt; 640px</td>
                      <td>1-2</td>
                      <td>16px</td>
                      <td>iPhone, Android</td>
                    </tr>
                    <tr>
                      <td className="font-semibold">Tablet</td>
                      <td>640px - 1024px</td>
                      <td>2-3</td>
                      <td>24px</td>
                      <td>iPad, Surface</td>
                    </tr>
                    <tr>
                      <td className="font-semibold">Desktop</td>
                      <td>1024px - 1280px</td>
                      <td>3-4</td>
                      <td>32px</td>
                      <td>Laptop, PC</td>
                    </tr>
                    <tr>
                      <td className="font-semibold">Wide</td>
                      <td>&gt; 1280px</td>
                      <td>4-6</td>
                      <td>48px</td>
                      <td>4K, Ultra-wide</td>
                    </tr>
                  </tbody>
                </table>
              </div>
            </div>
          </div>
        </AnimatedSection>
      </div>
    </div>
  );
};

export default AdaptiveDesign;