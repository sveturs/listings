'use client';

import React from 'react';
import { useTranslations } from 'next-intl';
import { SveTuLogo3D } from '@/components/logos/SveTuLogo3D';
import { SveTuLogoSpring } from '@/components/logos/SveTuLogoSpring';
import { SveTuLogoMorphing } from '@/components/logos/SveTuLogoMorphing';
import { SveTuLogoParticles } from '@/components/logos/SveTuLogoParticles';
import { SveTuLogoRosePetals } from '@/components/logos/SveTuLogoRosePetals';
import { SveTuLogoStatic } from '@/components/logos/SveTuLogoStatic';

const LogosPage = () => {
  const _t = useTranslations();

  return (
    <div className="min-h-screen bg-gradient-to-br from-base-100 to-base-200 py-8">
      <div className="container mx-auto px-4">
        <div className="text-center mb-12">
          <h1 className="text-4xl font-bold text-base-content mb-4">
            Интерактивные логотипы SveTu
          </h1>
          <p className="text-lg text-base-content-secondary max-w-2xl mx-auto">
            Коллекция анимированных логотипов с 3D эффектами и
            мультипликационными движениями
          </p>
        </div>

        <div className="grid grid-cols-1 md:grid-cols-2 gap-8 max-w-6xl mx-auto">
          {/* 3D Floating Tiles */}
          <div className="bg-white rounded-2xl shadow-xl p-8 hover:shadow-2xl transition-shadow duration-300">
            <div className="text-center mb-6">
              <h2 className="text-2xl font-bold text-gray-800 mb-2">
                3D Плавающие плитки
              </h2>
              <p className="text-gray-600">
                Плитки парят в 3D пространстве с физикой движения
              </p>
            </div>
            <div className="flex justify-center items-center h-80 bg-gradient-to-br from-blue-50 to-purple-50 rounded-xl">
              <SveTuLogo3D width={200} height={200} />
            </div>
            <div className="mt-4 text-sm text-gray-500 text-center">
              Наведите курсор • Кликните для разлета
            </div>
          </div>

          {/* Spring Animation */}
          <div className="bg-white rounded-2xl shadow-xl p-8 hover:shadow-2xl transition-shadow duration-300">
            <div className="text-center mb-6">
              <h2 className="text-2xl font-bold text-gray-800 mb-2">
                Пружинная анимация
              </h2>
              <p className="text-gray-600">
                Плитки двигаются как настоящие пружины
              </p>
            </div>
            <div className="flex justify-center items-center h-80 bg-gradient-to-br from-green-50 to-emerald-50 rounded-xl">
              <SveTuLogoSpring width={200} height={200} />
            </div>
            <div className="mt-4 text-sm text-gray-500 text-center">
              Кликните для активации
            </div>
          </div>

          {/* Morphing Effects */}
          <div className="bg-white rounded-2xl shadow-xl p-8 hover:shadow-2xl transition-shadow duration-300">
            <div className="text-center mb-6">
              <h2 className="text-2xl font-bold text-gray-800 mb-2">
                Морфинг и волны
              </h2>
              <p className="text-gray-600">
                Плитки плавно трансформируются волнообразно
              </p>
            </div>
            <div className="flex justify-center items-center h-80 bg-gradient-to-br from-orange-50 to-red-50 rounded-xl">
              <SveTuLogoMorphing width={200} height={200} />
            </div>
            <div className="mt-4 text-sm text-gray-500 text-center">
              Автоматическая анимация
            </div>
          </div>

          {/* Particle Effects */}
          <div className="bg-white rounded-2xl shadow-xl p-8 hover:shadow-2xl transition-shadow duration-300">
            <div className="text-center mb-6">
              <h2 className="text-2xl font-bold text-gray-800 mb-2">
                Частицы и свечение
              </h2>
              <p className="text-gray-600">
                Магические частицы и светящиеся эффекты
              </p>
            </div>
            <div className="flex justify-center items-center h-80 bg-gradient-to-br from-purple-50 to-pink-50 rounded-xl">
              <SveTuLogoParticles width={200} height={200} />
            </div>
            <div className="mt-4 text-sm text-gray-500 text-center">
              Наведите для магии
            </div>
          </div>

          {/* Rose Petals Animation */}
          <div className="bg-white rounded-2xl shadow-xl p-8 hover:shadow-2xl transition-shadow duration-300">
            <div className="text-center mb-6">
              <h2 className="text-2xl font-bold text-gray-800 mb-2">
                Лепестки роз
              </h2>
              <p className="text-gray-600">
                Плитки летают как лепестки роз на ветру
              </p>
            </div>
            <div className="flex justify-center items-center h-80 bg-gradient-to-br from-rose-50 to-pink-50 rounded-xl">
              <SveTuLogoRosePetals width={200} height={200} />
            </div>
            <div className="mt-4 text-sm text-gray-500 text-center">
              Кликните для полета лепестков
            </div>
          </div>
        </div>

        {/* Static Logos Section */}
        <div className="mt-16 bg-white rounded-2xl shadow-xl p-8">
          <h2 className="text-3xl font-bold text-center text-gray-800 mb-8">
            Статичные версии
          </h2>

          {/* 100x100 size */}
          <div className="mb-8">
            <h3 className="text-xl font-semibold text-gray-700 mb-4 text-center">
              Размер 100×100
            </h3>
            <div className="grid grid-cols-2 md:grid-cols-3 lg:grid-cols-5 gap-6">
              <div className="text-center">
                <div className="flex justify-center mb-4 bg-gray-50 p-4 rounded-xl">
                  <SveTuLogoStatic
                    variant="gradient"
                    width={100}
                    height={100}
                  />
                </div>
                <h4 className="font-semibold text-gray-700">Градиент</h4>
              </div>
              <div className="text-center">
                <div className="flex justify-center mb-4 bg-gray-50 p-4 rounded-xl">
                  <SveTuLogoStatic variant="minimal" width={100} height={100} />
                </div>
                <h4 className="font-semibold text-gray-700">Минимализм</h4>
              </div>
              <div className="text-center">
                <div className="flex justify-center mb-4 bg-gray-50 p-4 rounded-xl">
                  <SveTuLogoStatic variant="retro" width={100} height={100} />
                </div>
                <h4 className="font-semibold text-gray-700">Ретро</h4>
              </div>
              <div className="text-center">
                <div className="flex justify-center mb-4 bg-gray-900 p-4 rounded-xl">
                  <SveTuLogoStatic variant="neon" width={100} height={100} />
                </div>
                <h4 className="font-semibold text-gray-700">Неон</h4>
              </div>
              <div className="text-center">
                <div className="flex justify-center mb-4 bg-gradient-to-br from-blue-100 to-purple-100 p-4 rounded-xl">
                  <SveTuLogoStatic
                    variant="glassmorphic"
                    width={100}
                    height={100}
                  />
                </div>
                <h4 className="font-semibold text-gray-700">Стекло</h4>
              </div>
            </div>
          </div>

          {/* 48x48 size */}
          <div className="mb-8">
            <h3 className="text-xl font-semibold text-gray-700 mb-4 text-center">
              Размер 48×48
            </h3>
            <div className="grid grid-cols-2 md:grid-cols-3 lg:grid-cols-5 gap-6">
              <div className="text-center">
                <div className="flex justify-center items-center mb-4 bg-gray-50 p-4 rounded-xl h-24">
                  <SveTuLogoStatic variant="gradient" width={48} height={48} />
                </div>
                <h4 className="font-semibold text-gray-700 text-sm">
                  Градиент
                </h4>
              </div>
              <div className="text-center">
                <div className="flex justify-center items-center mb-4 bg-gray-50 p-4 rounded-xl h-24">
                  <SveTuLogoStatic variant="minimal" width={48} height={48} />
                </div>
                <h4 className="font-semibold text-gray-700 text-sm">
                  Минимализм
                </h4>
              </div>
              <div className="text-center">
                <div className="flex justify-center items-center mb-4 bg-gray-50 p-4 rounded-xl h-24">
                  <SveTuLogoStatic variant="retro" width={48} height={48} />
                </div>
                <h4 className="font-semibold text-gray-700 text-sm">Ретро</h4>
              </div>
              <div className="text-center">
                <div className="flex justify-center items-center mb-4 bg-gray-900 p-4 rounded-xl h-24">
                  <SveTuLogoStatic variant="neon" width={48} height={48} />
                </div>
                <h4 className="font-semibold text-gray-700 text-sm">Неон</h4>
              </div>
              <div className="text-center">
                <div className="flex justify-center items-center mb-4 bg-gradient-to-br from-blue-100 to-purple-100 p-4 rounded-xl h-24">
                  <SveTuLogoStatic
                    variant="glassmorphic"
                    width={48}
                    height={48}
                  />
                </div>
                <h4 className="font-semibold text-gray-700 text-sm">Стекло</h4>
              </div>
            </div>
          </div>

          {/* 32x32 size */}
          <div className="mb-8">
            <h3 className="text-xl font-semibold text-gray-700 mb-4 text-center">
              Размер 32×32
            </h3>
            <div className="grid grid-cols-2 md:grid-cols-3 lg:grid-cols-5 gap-6">
              <div className="text-center">
                <div className="flex justify-center items-center mb-4 bg-gray-50 p-4 rounded-xl h-24">
                  <SveTuLogoStatic variant="gradient" width={32} height={32} />
                </div>
                <h4 className="font-semibold text-gray-700 text-sm">
                  Градиент
                </h4>
              </div>
              <div className="text-center">
                <div className="flex justify-center items-center mb-4 bg-gray-50 p-4 rounded-xl h-24">
                  <SveTuLogoStatic variant="minimal" width={32} height={32} />
                </div>
                <h4 className="font-semibold text-gray-700 text-sm">
                  Минимализм
                </h4>
              </div>
              <div className="text-center">
                <div className="flex justify-center items-center mb-4 bg-gray-50 p-4 rounded-xl h-24">
                  <SveTuLogoStatic variant="retro" width={32} height={32} />
                </div>
                <h4 className="font-semibold text-gray-700 text-sm">Ретро</h4>
              </div>
              <div className="text-center">
                <div className="flex justify-center items-center mb-4 bg-gray-900 p-4 rounded-xl h-24">
                  <SveTuLogoStatic variant="neon" width={32} height={32} />
                </div>
                <h4 className="font-semibold text-gray-700 text-sm">Неон</h4>
              </div>
              <div className="text-center">
                <div className="flex justify-center items-center mb-4 bg-gradient-to-br from-blue-100 to-purple-100 p-4 rounded-xl h-24">
                  <SveTuLogoStatic
                    variant="glassmorphic"
                    width={32}
                    height={32}
                  />
                </div>
                <h4 className="font-semibold text-gray-700 text-sm">Стекло</h4>
              </div>
            </div>
          </div>

          {/* 16x16 size */}
          <div>
            <h3 className="text-xl font-semibold text-gray-700 mb-4 text-center">
              Размер 16×16
            </h3>
            <div className="grid grid-cols-2 md:grid-cols-3 lg:grid-cols-5 gap-6">
              <div className="text-center">
                <div className="flex justify-center items-center mb-4 bg-gray-50 p-4 rounded-xl h-24">
                  <SveTuLogoStatic variant="gradient" width={16} height={16} />
                </div>
                <h4 className="font-semibold text-gray-700 text-sm">
                  Градиент
                </h4>
              </div>
              <div className="text-center">
                <div className="flex justify-center items-center mb-4 bg-gray-50 p-4 rounded-xl h-24">
                  <SveTuLogoStatic variant="minimal" width={16} height={16} />
                </div>
                <h4 className="font-semibold text-gray-700 text-sm">
                  Минимализм
                </h4>
              </div>
              <div className="text-center">
                <div className="flex justify-center items-center mb-4 bg-gray-50 p-4 rounded-xl h-24">
                  <SveTuLogoStatic variant="retro" width={16} height={16} />
                </div>
                <h4 className="font-semibold text-gray-700 text-sm">Ретро</h4>
              </div>
              <div className="text-center">
                <div className="flex justify-center items-center mb-4 bg-gray-900 p-4 rounded-xl h-24">
                  <SveTuLogoStatic variant="neon" width={16} height={16} />
                </div>
                <h4 className="font-semibold text-gray-700 text-sm">Неон</h4>
              </div>
              <div className="text-center">
                <div className="flex justify-center items-center mb-4 bg-gradient-to-br from-blue-100 to-purple-100 p-4 rounded-xl h-24">
                  <SveTuLogoStatic
                    variant="glassmorphic"
                    width={16}
                    height={16}
                  />
                </div>
                <h4 className="font-semibold text-gray-700 text-sm">Стекло</h4>
              </div>
            </div>
          </div>
        </div>

        {/* Comparison Section */}
        <div className="mt-16 bg-white rounded-2xl shadow-xl p-8">
          <h2 className="text-3xl font-bold text-center text-gray-800 mb-8">
            Сравнение динамических версий
          </h2>

          {/* 80x80 size */}
          <div className="mb-8">
            <h3 className="text-xl font-semibold text-gray-700 mb-4 text-center">
              Размер 80×80
            </h3>
            <div className="grid grid-cols-2 md:grid-cols-5 gap-6">
              <div className="text-center">
                <div className="flex justify-center mb-4">
                  <SveTuLogo3D width={80} height={80} />
                </div>
                <h3 className="font-semibold text-gray-700 text-sm">
                  3D Floating
                </h3>
              </div>
              <div className="text-center">
                <div className="flex justify-center mb-4">
                  <SveTuLogoSpring width={80} height={80} />
                </div>
                <h3 className="font-semibold text-gray-700 text-sm">Spring</h3>
              </div>
              <div className="text-center">
                <div className="flex justify-center mb-4">
                  <SveTuLogoMorphing width={80} height={80} />
                </div>
                <h3 className="font-semibold text-gray-700 text-sm">
                  Morphing
                </h3>
              </div>
              <div className="text-center">
                <div className="flex justify-center mb-4">
                  <SveTuLogoParticles width={80} height={80} />
                </div>
                <h3 className="font-semibold text-gray-700 text-sm">
                  Particles
                </h3>
              </div>
              <div className="text-center">
                <div className="flex justify-center mb-4">
                  <SveTuLogoRosePetals width={80} height={80} />
                </div>
                <h3 className="font-semibold text-gray-700 text-sm">
                  Rose Petals
                </h3>
              </div>
            </div>
          </div>

          {/* 48x48 size */}
          <div className="mt-8">
            <h3 className="text-xl font-semibold text-gray-700 mb-4 text-center">
              Размер 48×48
            </h3>
            <div className="grid grid-cols-2 md:grid-cols-5 gap-6">
              <div className="text-center">
                <div className="flex justify-center items-center mb-4 h-20">
                  <SveTuLogo3D width={48} height={48} />
                </div>
                <h4 className="font-semibold text-gray-700 text-sm">
                  3D Floating
                </h4>
              </div>
              <div className="text-center">
                <div className="flex justify-center items-center mb-4 h-20">
                  <SveTuLogoSpring width={48} height={48} />
                </div>
                <h4 className="font-semibold text-gray-700 text-sm">Spring</h4>
              </div>
              <div className="text-center">
                <div className="flex justify-center items-center mb-4 h-20">
                  <SveTuLogoMorphing width={48} height={48} />
                </div>
                <h4 className="font-semibold text-gray-700 text-sm">
                  Morphing
                </h4>
              </div>
              <div className="text-center">
                <div className="flex justify-center items-center mb-4 h-20">
                  <SveTuLogoParticles width={48} height={48} />
                </div>
                <h4 className="font-semibold text-gray-700 text-sm">
                  Particles
                </h4>
              </div>
              <div className="text-center">
                <div className="flex justify-center items-center mb-4 h-20">
                  <SveTuLogoRosePetals width={48} height={48} />
                </div>
                <h4 className="font-semibold text-gray-700 text-sm">
                  Rose Petals
                </h4>
              </div>
            </div>
          </div>
        </div>

        {/* Technical Info */}
        <div className="mt-12 bg-gradient-to-r from-blue-50 to-purple-50 rounded-2xl p-8">
          <h2 className="text-2xl font-bold text-center text-gray-800 mb-6">
            Технические особенности
          </h2>
          <div className="grid grid-cols-1 md:grid-cols-3 gap-6 text-center">
            <div>
              <div className="text-3xl mb-2">🎭</div>
              <h3 className="font-semibold text-gray-700 mb-2">CSS3 + SVG</h3>
              <p className="text-gray-600 text-sm">
                Использование современных веб-технологий для плавных анимаций
              </p>
            </div>
            <div>
              <div className="text-3xl mb-2">⚡</div>
              <h3 className="font-semibold text-gray-700 mb-2">60 FPS</h3>
              <p className="text-gray-600 text-sm">
                Оптимизированные анимации для максимальной производительности
              </p>
            </div>
            <div>
              <div className="text-3xl mb-2">📱</div>
              <h3 className="font-semibold text-gray-700 mb-2">Responsive</h3>
              <p className="text-gray-600 text-sm">
                Адаптивность под все устройства и разрешения экранов
              </p>
            </div>
          </div>
        </div>
      </div>
    </div>
  );
};

export default LogosPage;
