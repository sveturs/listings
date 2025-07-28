'use client';

import React, { useState } from 'react';
import { SveTuLogoStatic } from '@/components/logos/SveTuLogoStatic';
import { AnimatedSection } from '@/components/ui/AnimatedSection';

const AIListingCreator = () => {
  const [selectedImages, setSelectedImages] = useState<string[]>([]);
  const [analyzing, setAnalyzing] = useState(false);
  const [aiResult, setAiResult] = useState<any>(null);
  const [dragActive, setDragActive] = useState(false);

  const demoImages = [
    '/api/minio/download?fileName=listings/0a47e66f-d8da-459f-a2ba-8e2b85ae0163/38ad29e6-7b07-4bfc-9db2-d965cb6b966f.jpg',
    '/api/minio/download?fileName=listings/0c1fc30d-5d84-485f-a86a-5c5dc37f8b97/4b8b8e48-ddd8-4c04-ad8e-00c4b4d10d26.jpg',
    '/api/minio/download?fileName=listings/0c91d2f7-53f7-4bff-87fe-d7e82dc3e2f0/3b26f07f-c5d6-4ff7-ba56-06ec69bb7f4d.jpg',
    '/api/minio/download?fileName=listings/0e17c3be-e76e-433a-a6d4-86bb8b7a0e29/23bb3da7-38ef-44f7-8c1d-1c14eaaafeb5.jpg',
  ];

  const handleDrop = (e: React.DragEvent) => {
    e.preventDefault();
    setDragActive(false);
    handleAnalyze();
  };

  const handleDragOver = (e: React.DragEvent) => {
    e.preventDefault();
    setDragActive(true);
  };

  const handleDragLeave = () => {
    setDragActive(false);
  };

  const handleAnalyze = () => {
    setAnalyzing(true);
    setSelectedImages(demoImages.slice(0, 3));
    
    setTimeout(() => {
      setAiResult({
        title: 'Современная 2-комнатная квартира с дизайнерским ремонтом',
        description: `✨ Уютная квартира в самом сердце города!

🏠 Характеристики:
• Площадь: 65 м²
• 2 спальни + гостиная
• Современная кухня с техникой Bosch
• Панорамные окна с видом на парк
• Дизайнерский ремонт 2023 года

🌟 Преимущества:
• 5 минут до метро
• Развитая инфраструктура
• Охраняемая территория
• Подземный паркинг
• Детская площадка во дворе

💰 Все включено в стоимость!`,
        category: 'Недвижимость / Квартиры',
        suggestedPrice: '850 €/месяц',
        attributes: {
          rooms: '2',
          area: '65 м²',
          floor: '7/12',
          furnished: 'Да',
          parking: 'Подземный',
          pets: 'По согласованию'
        },
        tags: ['Центр', 'Метро рядом', 'Новый ремонт', 'С мебелью', 'Парковка'],
        quality_score: 95,
        suggestions: [
          'Добавьте фото ванной комнаты',
          'Укажите точный адрес для лучшей видимости',
          'Добавьте виртуальный тур'
        ]
      });
      setAnalyzing(false);
    }, 3000);
  };

  return (
    <div className="min-h-screen bg-gradient-to-br from-base-100 to-base-200">
      {/* Header */}
      <div className="navbar bg-base-100 shadow-lg">
        <div className="navbar-start">
          <SveTuLogoStatic variant="gradient" width={120} height={40} />
        </div>
        <div className="navbar-center">
          <h1 className="text-xl font-bold">AI Создание объявлений</h1>
        </div>
        <div className="navbar-end">
          <button className="btn btn-ghost btn-circle">
            <svg className="w-6 h-6" fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M4 6h16M4 12h16M4 18h16" />
            </svg>
          </button>
        </div>
      </div>

      <div className="container mx-auto px-4 py-8 max-w-6xl">
        <AnimatedSection animation="fadeIn">
          <div className="text-center mb-8">
            <h2 className="text-3xl font-bold mb-4">
              🤖 AI-анализ фото для создания объявления
            </h2>
            <p className="text-lg text-base-content/70">
              Просто загрузите фото - AI сделает всё остальное!
            </p>
          </div>
        </AnimatedSection>

        <div className="grid grid-cols-1 lg:grid-cols-2 gap-8">
          {/* Upload Section */}
          <AnimatedSection animation="slideLeft">
            <div className="card bg-base-100 shadow-xl">
              <div className="card-body">
                <h3 className="card-title mb-4">📸 Загрузка фото</h3>
                
                <div
                  className={`border-2 border-dashed rounded-xl p-8 text-center transition-all ${
                    dragActive ? 'border-primary bg-primary/10' : 'border-base-300'
                  }`}
                  onDrop={handleDrop}
                  onDragOver={handleDragOver}
                  onDragLeave={handleDragLeave}
                >
                  <svg className="w-16 h-16 mx-auto mb-4 text-base-content/40" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                    <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M7 16a4 4 0 01-.88-7.903A5 5 0 1115.9 6L16 6a5 5 0 011 9.9M15 13l-3-3m0 0l-3 3m3-3v12" />
                  </svg>
                  <p className="text-lg mb-2">Перетащите фото сюда</p>
                  <p className="text-sm text-base-content/60 mb-4">или</p>
                  <button className="btn btn-primary" onClick={handleAnalyze}>
                    Использовать демо-фото
                  </button>
                </div>

                {selectedImages.length > 0 && (
                  <div className="mt-6">
                    <h4 className="font-semibold mb-3">Загруженные фото:</h4>
                    <div className="grid grid-cols-3 gap-2">
                      {selectedImages.map((img, idx) => (
                        <div key={idx} className="relative aspect-square rounded-lg overflow-hidden">
                          <img 
                            src={img} 
                            alt={`Photo ${idx + 1}`}
                            className="w-full h-full object-cover"
                          />
                          {analyzing && (
                            <div className="absolute inset-0 bg-black/50 flex items-center justify-center">
                              <span className="loading loading-spinner loading-lg text-white"></span>
                            </div>
                          )}
                        </div>
                      ))}
                    </div>
                  </div>
                )}

                {analyzing && (
                  <div className="mt-6">
                    <div className="flex items-center gap-3">
                      <span className="loading loading-dots loading-md"></span>
                      <span>AI анализирует фотографии...</span>
                    </div>
                    <progress className="progress progress-primary w-full mt-2" value="70" max="100"></progress>
                  </div>
                )}
              </div>
            </div>
          </AnimatedSection>

          {/* AI Result Section */}
          <AnimatedSection animation="slideRight" delay={0.2}>
            <div className="card bg-base-100 shadow-xl">
              <div className="card-body">
                <h3 className="card-title mb-4">✨ Результат AI-анализа</h3>
                
                {!aiResult ? (
                  <div className="text-center py-12 text-base-content/50">
                    <svg className="w-24 h-24 mx-auto mb-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                      <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M9.663 17h4.673M12 3v1m6.364 1.636l-.707.707M21 12h-1M4 12H3m3.343-5.657l-.707-.707m2.828 9.9a5 5 0 117.072 0l-.548.547A3.374 3.374 0 0014 18.469V19a2 2 0 11-4 0v-.531c0-.895-.356-1.754-.988-2.386l-.548-.547z" />
                    </svg>
                    <p>Загрузите фото для анализа</p>
                  </div>
                ) : (
                  <div className="space-y-4">
                    <div>
                      <label className="label">
                        <span className="label-text">Заголовок</span>
                        <span className="label-text-alt text-success">AI: 95% точность</span>
                      </label>
                      <input 
                        type="text" 
                        className="input input-bordered w-full" 
                        value={aiResult.title}
                        readOnly
                      />
                    </div>

                    <div>
                      <label className="label">
                        <span className="label-text">Описание</span>
                      </label>
                      <textarea 
                        className="textarea textarea-bordered w-full h-32" 
                        value={aiResult.description}
                        readOnly
                      />
                    </div>

                    <div className="grid grid-cols-2 gap-4">
                      <div>
                        <label className="label">
                          <span className="label-text">Категория</span>
                        </label>
                        <input 
                          type="text" 
                          className="input input-bordered w-full" 
                          value={aiResult.category}
                          readOnly
                        />
                      </div>
                      <div>
                        <label className="label">
                          <span className="label-text">Рекомендуемая цена</span>
                        </label>
                        <input 
                          type="text" 
                          className="input input-bordered w-full" 
                          value={aiResult.suggestedPrice}
                          readOnly
                        />
                      </div>
                    </div>

                    <div>
                      <label className="label">
                        <span className="label-text">Теги</span>
                      </label>
                      <div className="flex flex-wrap gap-2">
                        {aiResult.tags.map((tag: string, idx: number) => (
                          <span key={idx} className="badge badge-primary">{tag}</span>
                        ))}
                      </div>
                    </div>

                    <div className="divider"></div>

                    <div>
                      <h4 className="font-semibold mb-2">💡 Рекомендации AI:</h4>
                      <ul className="space-y-1">
                        {aiResult.suggestions.map((suggestion: string, idx: number) => (
                          <li key={idx} className="flex items-start gap-2">
                            <span className="text-warning">•</span>
                            <span className="text-sm">{suggestion}</span>
                          </li>
                        ))}
                      </ul>
                    </div>

                    <div className="card-actions justify-end mt-6">
                      <button className="btn btn-ghost">Редактировать</button>
                      <button className="btn btn-primary">
                        Опубликовать объявление
                      </button>
                    </div>
                  </div>
                )}
              </div>
            </div>
          </AnimatedSection>
        </div>

        {/* Features Section */}
        <AnimatedSection animation="fadeIn" delay={0.4}>
          <div className="mt-12 grid grid-cols-1 md:grid-cols-3 gap-6">
            <div className="card bg-primary/10">
              <div className="card-body text-center">
                <div className="text-4xl mb-4">🎯</div>
                <h3 className="font-bold text-lg">Точное распознавание</h3>
                <p className="text-sm">AI определяет тип товара, состояние и ключевые характеристики</p>
              </div>
            </div>
            <div className="card bg-secondary/10">
              <div className="card-body text-center">
                <div className="text-4xl mb-4">💰</div>
                <h3 className="font-bold text-lg">Умное ценообразование</h3>
                <p className="text-sm">Анализ рынка и рекомендация оптимальной цены</p>
              </div>
            </div>
            <div className="card bg-accent/10">
              <div className="card-body text-center">
                <div className="text-4xl mb-4">📝</div>
                <h3 className="font-bold text-lg">SEO-оптимизация</h3>
                <p className="text-sm">Автоматическая генерация тегов и ключевых слов</p>
              </div>
            </div>
          </div>
        </AnimatedSection>
      </div>
    </div>
  );
};

export default AIListingCreator;