'use client';

import React, { useState, useRef, useEffect } from 'react';
import Link from 'next/link';
import Image from 'next/image';
import { motion, AnimatePresence } from 'framer-motion';
import {
  ChevronLeft,
  Camera,
  Sparkles,
  Check,
  Mic,
  X,
  TrendingUp,
  Clock,
  Eye,
  Heart,
  MessageCircle,
  Share2,
  Shield,
  Award,
  Brain,
  Zap,
  Plus,
  RefreshCw,
  Globe,
  BarChart3,
  Users,
  ThumbsUp,
  Instagram,
  Facebook,
  Send,
  Calendar,
  Languages,
  TestTube2,
  Lightbulb,
  Package,
} from 'lucide-react';

export default function AIPoweredEnhancedListingCreationPage() {
  const [currentView, setCurrentView] = useState<
    'upload' | 'process' | 'enhance' | 'publish'
  >('upload');
  const [images, setImages] = useState<string[]>([]);
  const [isProcessing, setIsProcessing] = useState(false);
  const [voiceRecording, setVoiceRecording] = useState(false);

  // AI generated data
  const [aiData, setAiData] = useState({
    title: '',
    titleVariants: [] as string[],
    selectedTitleIndex: 0,
    description: '',
    category: '',
    categoryProbabilities: [] as { name: string; probability: number }[],
    price: '',
    priceRange: { min: 0, max: 0 },
    attributes: {} as Record<string, string>,
    tags: [] as string[],
    suggestedPhotos: [] as string[],
    translations: {} as Record<string, { title: string; description: string }>,
    publishTime: '',
    socialPosts: {} as Record<string, string>,
  });

  // A/B testing state
  const [abTestEnabled, setAbTestEnabled] = useState(false);
  const [selectedVariant, setSelectedVariant] = useState(0);

  const fileInputRef = useRef<HTMLInputElement>(null);

  // Simulate AI processing
  const processImages = async () => {
    setIsProcessing(true);
    setCurrentView('process');

    // Simulate processing time
    await new Promise((resolve) => setTimeout(resolve, 2000));

    // Generate AI data
    setAiData({
      title: 'iPhone 13 Pro, 256GB, Pacific Blue',
      titleVariants: [
        'iPhone 13 Pro, 256GB, Pacific Blue - Идеальное состояние',
        'Продаю iPhone 13 Pro 256GB (Pacific Blue) - как новый!',
        '📱 iPhone 13 Pro | 256GB | Pacific Blue | Гарантия',
      ],
      selectedTitleIndex: 0,
      description: `Продаю iPhone 13 Pro в идеальном состоянии!

📱 Модель: iPhone 13 Pro
💾 Память: 256GB
🎨 Цвет: Pacific Blue
🔋 Батарея: 92% (отличное состояние)
📦 Комплект: полный (коробка, зарядка, документы)

✅ Без царапин и сколов
✅ Всегда в чехле и с защитным стеклом
✅ Никогда не падал и не ремонтировался
✅ Все функции работают идеально

Причина продажи: переход на новую модель.`,
      category: 'electronics',
      categoryProbabilities: [
        { name: 'Электроника', probability: 98 },
        { name: 'Телефоны', probability: 95 },
        { name: 'Apple', probability: 92 },
      ],
      price: '65000',
      priceRange: { min: 60000, max: 70000 },
      attributes: {
        brand: 'Apple',
        model: 'iPhone 13 Pro',
        storage: '256GB',
        color: 'Pacific Blue',
        condition: 'Как новый',
        warranty: 'Нет',
        battery: '92%',
      },
      tags: ['iPhone', 'Apple', '256GB', 'Pro', 'Синий', 'Смартфон'],
      suggestedPhotos: [
        'Фото экрана включенного',
        'Фото задней панели',
        'Фото с комплектом',
        'Фото в чехле',
      ],
      translations: {
        en: {
          title: 'iPhone 13 Pro, 256GB, Pacific Blue',
          description: 'Selling iPhone 13 Pro in perfect condition!',
        },
        sr: {
          title: 'iPhone 13 Pro, 256GB, Pacific Blue',
          description: 'Prodajem iPhone 13 Pro u savršenom stanju!',
        },
      },
      publishTime: '19:00',
      socialPosts: {
        whatsapp:
          '📱 Продаю iPhone 13 Pro, 256GB\n💙 Pacific Blue\n✨ Идеальное состояние\n💰 65.000 РСД',
        telegram:
          '📱 iPhone 13 Pro на продажу!\n\n• 256GB, Pacific Blue\n• Состояние: как новый\n• Батарея: 92%\n• Цена: 65.000 РСД\n\nПодробности в личку 📩',
        instagram:
          '#ПродамiPhone #iPhone13Pro #Belgrade #Сербия\n\n📱 iPhone 13 Pro, 256GB\n💙 Цвет: Pacific Blue\n⚡ Состояние: идеальное\n💰 Цена: 65.000 РСД\n\nDM для деталей! 📩',
      },
    });

    setIsProcessing(false);
    setCurrentView('enhance');
  };

  const handleImageUpload = (e: React.ChangeEvent<HTMLInputElement>) => {
    const files = e.target.files;
    if (files) {
      const newImages = Array.from(files).map((file) =>
        URL.createObjectURL(file)
      );
      setImages([...images, ...newImages].slice(0, 8));
    }
  };

  const removeImage = (index: number) => {
    setImages(images.filter((_, i) => i !== index));
  };

  const regenerateTitle = () => {
    const newIndex =
      (aiData.selectedTitleIndex + 1) % aiData.titleVariants.length;
    setAiData({ ...aiData, selectedTitleIndex: newIndex });
  };

  const renderUploadView = () => (
    <motion.div
      initial={{ opacity: 0 }}
      animate={{ opacity: 1 }}
      exit={{ opacity: 0 }}
      className="min-h-screen bg-gradient-to-br from-base-100 to-base-200"
    >
      <div className="container mx-auto px-4 py-16">
        <motion.div
          initial={{ y: 20, opacity: 0 }}
          animate={{ y: 0, opacity: 1 }}
          className="text-center mb-12"
        >
          <div className="inline-flex items-center justify-center w-20 h-20 bg-gradient-to-br from-primary to-secondary rounded-full mb-6">
            <Brain className="w-10 h-10 text-primary-content" />
          </div>
          <h1 className="text-4xl lg:text-5xl font-bold mb-4">
            AI создаст объявление за вас
          </h1>
          <p className="text-xl text-base-content/70 mb-8">
            Просто загрузите фото — остальное сделает искусственный интеллект
          </p>

          <div className="flex justify-center gap-6 mb-8">
            <div className="text-center">
              <div className="text-3xl font-bold text-primary">30 сек</div>
              <div className="text-sm text-base-content/60">создание</div>
            </div>
            <div className="text-center">
              <div className="text-3xl font-bold text-success">95%</div>
              <div className="text-sm text-base-content/60">точность AI</div>
            </div>
            <div className="text-center">
              <div className="text-3xl font-bold text-secondary">5 языков</div>
              <div className="text-sm text-base-content/60">перевод</div>
            </div>
          </div>
        </motion.div>

        {images.length === 0 ? (
          <motion.div
            initial={{ scale: 0.9, opacity: 0 }}
            animate={{ scale: 1, opacity: 1 }}
            transition={{ delay: 0.2 }}
            className="max-w-2xl mx-auto"
          >
            <label
              htmlFor="ai-upload"
              className="card bg-gradient-to-br from-primary/10 to-secondary/10 border-2 border-dashed border-primary cursor-pointer hover:shadow-2xl transition-all"
            >
              <div className="card-body text-center py-16">
                <Camera className="w-20 h-20 mx-auto mb-4 text-primary" />
                <h2 className="text-2xl font-bold mb-2">
                  Загрузите фото товара
                </h2>
                <p className="text-base-content/70 mb-6">
                  AI распознает товар и создаст идеальное объявление
                </p>
                <div className="flex gap-4 justify-center">
                  <div className="badge badge-lg badge-primary gap-2">
                    <Brain className="w-4 h-4" />
                    AI распознавание
                  </div>
                  <div className="badge badge-lg badge-secondary gap-2">
                    <Zap className="w-4 h-4" />
                    30 секунд
                  </div>
                </div>
              </div>
            </label>
            <input
              id="ai-upload"
              ref={fileInputRef}
              type="file"
              multiple
              accept="image/*"
              className="hidden"
              onChange={handleImageUpload}
            />

            {/* Alternative input methods */}
            <div className="grid grid-cols-1 lg:grid-cols-3 gap-4 mt-6">
              <button className="btn btn-outline gap-2">
                <Instagram className="w-4 h-4" />
                Импорт из Instagram
              </button>
              <button className="btn btn-outline gap-2">
                <Facebook className="w-4 h-4" />
                Импорт из Facebook
              </button>
              <button
                onClick={() => setVoiceRecording(!voiceRecording)}
                className={`btn ${voiceRecording ? 'btn-error' : 'btn-outline'} gap-2`}
              >
                <Mic className="w-4 h-4" />
                {voiceRecording ? 'Остановить запись' : 'Голосовое описание'}
              </button>
            </div>
          </motion.div>
        ) : (
          <div className="max-w-4xl mx-auto">
            <div className="grid grid-cols-2 lg:grid-cols-4 gap-4 mb-8">
              {images.map((img, index) => (
                <motion.div
                  key={index}
                  initial={{ scale: 0 }}
                  animate={{ scale: 1 }}
                  className="relative aspect-square"
                >
                  <Image
                    src={img}
                    alt={`Photo ${index + 1}`}
                    fill
                    className="object-cover rounded-lg"
                  />
                  <button
                    onClick={() => removeImage(index)}
                    className="absolute top-2 right-2 btn btn-circle btn-sm btn-error"
                  >
                    <X className="w-4 h-4" />
                  </button>
                </motion.div>
              ))}
              {images.length < 8 && (
                <label className="aspect-square border-2 border-dashed border-base-300 rounded-lg flex flex-col items-center justify-center cursor-pointer hover:border-primary transition-colors">
                  <Plus className="w-8 h-8 text-base-content/50" />
                  <span className="text-sm text-base-content/50 mt-2">
                    Добавить еще
                  </span>
                  <input
                    type="file"
                    multiple
                    accept="image/*"
                    className="hidden"
                    onChange={handleImageUpload}
                  />
                </label>
              )}
            </div>

            <motion.button
              initial={{ y: 20, opacity: 0 }}
              animate={{ y: 0, opacity: 1 }}
              onClick={processImages}
              className="btn btn-primary btn-lg btn-block"
            >
              <Brain className="w-5 h-5 mr-2" />
              Создать объявление с помощью AI
            </motion.button>
          </div>
        )}
      </div>
    </motion.div>
  );

  const renderProcessView = () => (
    <motion.div
      initial={{ opacity: 0 }}
      animate={{ opacity: 1 }}
      exit={{ opacity: 0 }}
      className="min-h-screen bg-gradient-to-br from-base-100 to-base-200 flex items-center justify-center"
    >
      <div className="text-center">
        <motion.div
          animate={{
            rotate: 360,
          }}
          transition={{
            duration: 2,
            repeat: Infinity,
            ease: 'linear',
          }}
          className="inline-flex items-center justify-center w-24 h-24 bg-gradient-to-br from-primary to-secondary rounded-full mb-8"
        >
          <Brain className="w-12 h-12 text-primary-content" />
        </motion.div>

        <h2 className="text-2xl font-bold mb-4">AI анализирует ваши фото</h2>

        <div className="space-y-4 text-left max-w-md mx-auto">
          <motion.div
            initial={{ x: -20, opacity: 0 }}
            animate={{ x: 0, opacity: 1 }}
            transition={{ delay: 0.2 }}
            className="flex items-center gap-3"
          >
            <div className="loading loading-spinner loading-sm text-primary"></div>
            <span>Распознавание товара...</span>
          </motion.div>
          <motion.div
            initial={{ x: -20, opacity: 0 }}
            animate={{ x: 0, opacity: 1 }}
            transition={{ delay: 0.4 }}
            className="flex items-center gap-3"
          >
            <div className="loading loading-spinner loading-sm text-primary"></div>
            <span>Анализ рынка и цен...</span>
          </motion.div>
          <motion.div
            initial={{ x: -20, opacity: 0 }}
            animate={{ x: 0, opacity: 1 }}
            transition={{ delay: 0.6 }}
            className="flex items-center gap-3"
          >
            <div className="loading loading-spinner loading-sm text-primary"></div>
            <span>Генерация описания...</span>
          </motion.div>
          <motion.div
            initial={{ x: -20, opacity: 0 }}
            animate={{ x: 0, opacity: 1 }}
            transition={{ delay: 0.8 }}
            className="flex items-center gap-3"
          >
            <div className="loading loading-spinner loading-sm text-primary"></div>
            <span>SEO оптимизация...</span>
          </motion.div>
          <motion.div
            initial={{ x: -20, opacity: 0 }}
            animate={{ x: 0, opacity: 1 }}
            transition={{ delay: 1.0 }}
            className="flex items-center gap-3"
          >
            <div className="loading loading-spinner loading-sm text-primary"></div>
            <span>Создание переводов...</span>
          </motion.div>
        </div>
      </div>
    </motion.div>
  );

  const renderEnhanceView = () => (
    <motion.div
      initial={{ opacity: 0 }}
      animate={{ opacity: 1 }}
      exit={{ opacity: 0 }}
      className="min-h-screen bg-base-100"
    >
      <div className="container mx-auto px-4 py-8">
        <div className="max-w-4xl mx-auto">
          {/* Success banner */}
          <motion.div
            initial={{ y: -20, opacity: 0 }}
            animate={{ y: 0, opacity: 1 }}
            className="alert alert-success shadow-lg mb-8"
          >
            <Check className="w-6 h-6" />
            <div>
              <h3 className="font-bold">AI успешно создал ваше объявление!</h3>
              <p>Проверьте и отредактируйте при необходимости</p>
            </div>
          </motion.div>

          {/* Photos section */}
          <div className="card bg-base-200 mb-6">
            <div className="card-body">
              <h3 className="card-title">
                <Camera className="w-5 h-5" />
                Фотографии
                <span className="badge badge-primary">{images.length}/8</span>
              </h3>
              <div className="grid grid-cols-4 gap-3">
                {images.map((img, index) => (
                  <div key={index} className="relative aspect-square">
                    <Image
                      src={img}
                      alt={`Photo ${index + 1}`}
                      fill
                      className="object-cover rounded-lg"
                    />
                    {index === 0 && (
                      <div className="absolute top-1 left-1 badge badge-primary badge-sm">
                        Главное
                      </div>
                    )}
                  </div>
                ))}
              </div>

              {/* Suggested missing photos */}
              {aiData.suggestedPhotos.length > 0 && (
                <div className="alert alert-info mt-4">
                  <Lightbulb className="w-4 h-4" />
                  <div>
                    <p className="font-semibold text-sm">
                      AI рекомендует добавить:
                    </p>
                    <ul className="text-xs mt-1">
                      {aiData.suggestedPhotos.map((photo, index) => (
                        <li key={index}>• {photo}</li>
                      ))}
                    </ul>
                  </div>
                </div>
              )}
            </div>
          </div>

          {/* Category with probabilities */}
          <div className="card bg-base-200 mb-6">
            <div className="card-body">
              <h3 className="card-title">
                <Package className="w-5 h-5" />
                Категория
              </h3>
              <div className="space-y-2">
                {aiData.categoryProbabilities.map((cat, index) => (
                  <div
                    key={index}
                    className="flex items-center justify-between"
                  >
                    <span className={index === 0 ? 'font-semibold' : ''}>
                      {cat.name}
                    </span>
                    <div className="flex items-center gap-2">
                      <progress
                        className="progress progress-primary w-32"
                        value={cat.probability}
                        max="100"
                      ></progress>
                      <span className="text-sm">{cat.probability}%</span>
                    </div>
                  </div>
                ))}
              </div>
            </div>
          </div>

          {/* Title with A/B testing */}
          <div className="card bg-base-200 mb-6">
            <div className="card-body">
              <div className="flex items-center justify-between mb-2">
                <h3 className="card-title">
                  <TestTube2 className="w-5 h-5" />
                  Заголовок (A/B тестирование)
                </h3>
                <button
                  onClick={regenerateTitle}
                  className="btn btn-ghost btn-sm gap-1"
                >
                  <RefreshCw className="w-4 h-4" />
                  Другой вариант
                </button>
              </div>

              <div className="space-y-3">
                {aiData.titleVariants.map((variant, index) => (
                  <label
                    key={index}
                    className={`card cursor-pointer ${
                      aiData.selectedTitleIndex === index
                        ? 'ring-2 ring-primary'
                        : ''
                    }`}
                  >
                    <div className="card-body p-3">
                      <div className="flex items-start gap-3">
                        <input
                          type="radio"
                          name="title"
                          className="radio radio-primary"
                          checked={aiData.selectedTitleIndex === index}
                          onChange={() =>
                            setAiData({ ...aiData, selectedTitleIndex: index })
                          }
                        />
                        <div className="flex-1">
                          <p className="font-medium">{variant}</p>
                          <p className="text-xs text-base-content/60 mt-1">
                            Прогноз CTR: {95 - index * 5}%
                          </p>
                        </div>
                      </div>
                    </div>
                  </label>
                ))}
              </div>

              <div className="form-control form-control-sm mt-3">
                <label className="label cursor-pointer">
                  <span className="label-text">Включить A/B тестирование</span>
                  <input
                    type="checkbox"
                    className="toggle toggle-primary"
                    checked={abTestEnabled}
                    onChange={(e) => setAbTestEnabled(e.target.checked)}
                  />
                </label>
              </div>
            </div>
          </div>

          {/* Price with market analysis */}
          <div className="card bg-base-200 mb-6">
            <div className="card-body">
              <h3 className="card-title">
                <BarChart3 className="w-5 h-5" />
                Цена (AI анализ рынка)
              </h3>
              <div className="text-3xl font-bold text-primary mb-2">
                {aiData.price} РСД
              </div>
              <p className="text-sm text-base-content/60 mb-4">
                Рекомендуемый диапазон: {aiData.priceRange.min.toLocaleString()}{' '}
                - {aiData.priceRange.max.toLocaleString()} РСД
              </p>

              <div className="bg-base-100 p-3 rounded-lg">
                <p className="text-sm font-semibold mb-2">
                  Анализ конкурентов:
                </p>
                <div className="space-y-1 text-xs">
                  <div className="flex justify-between">
                    <span>Минимальная цена:</span>
                    <span className="font-medium">58.000 РСД</span>
                  </div>
                  <div className="flex justify-between">
                    <span>Средняя цена:</span>
                    <span className="font-medium">65.000 РСД</span>
                  </div>
                  <div className="flex justify-between">
                    <span>Максимальная цена:</span>
                    <span className="font-medium">72.000 РСД</span>
                  </div>
                </div>
              </div>
            </div>
          </div>

          {/* AI Generated Description */}
          <div className="card bg-base-200 mb-6">
            <div className="card-body">
              <h3 className="card-title">
                <Sparkles className="w-5 h-5" />
                Описание (AI-оптимизировано)
              </h3>
              <textarea
                className="textarea textarea-bordered h-48"
                value={aiData.description}
                onChange={(e) =>
                  setAiData({ ...aiData, description: e.target.value })
                }
              />
              <div className="flex flex-wrap gap-2 mt-3">
                {aiData.tags.map((tag, index) => (
                  <span key={index} className="badge badge-sm">
                    #{tag}
                  </span>
                ))}
              </div>
            </div>
          </div>

          {/* Multi-language support */}
          <div className="card bg-base-200 mb-6">
            <div className="card-body">
              <h3 className="card-title">
                <Languages className="w-5 h-5" />
                Мультиязычность
              </h3>
              <div className="grid grid-cols-1 lg:grid-cols-2 gap-4">
                {Object.entries(aiData.translations).map(([lang, trans]) => (
                  <div key={lang} className="border rounded-lg p-3">
                    <div className="flex items-center gap-2 mb-2">
                      <Globe className="w-4 h-4" />
                      <span className="font-semibold text-sm">
                        {lang === 'en' ? 'English' : 'Српски'}
                      </span>
                    </div>
                    <p className="font-medium text-sm mb-1">{trans.title}</p>
                    <p className="text-xs text-base-content/70">
                      {trans.description}
                    </p>
                  </div>
                ))}
              </div>
            </div>
          </div>

          {/* Publishing optimization */}
          <div className="card bg-gradient-to-r from-warning/10 to-warning/5 border-2 border-warning/20 mb-6">
            <div className="card-body">
              <h3 className="card-title">
                <Calendar className="w-5 h-5" />
                Оптимальное время публикации
              </h3>
              <p className="text-sm mb-3">
                AI проанализировал активность покупателей в вашей категории
              </p>
              <div className="grid grid-cols-3 gap-3">
                <div className="text-center p-3 bg-base-100 rounded-lg">
                  <p className="text-2xl font-bold text-warning">19:00</p>
                  <p className="text-xs">Лучшее время</p>
                  <p className="text-xs text-success">+45% просмотров</p>
                </div>
                <div className="text-center p-3 bg-base-100 rounded-lg">
                  <p className="text-2xl font-bold">12:00</p>
                  <p className="text-xs">Хорошее время</p>
                  <p className="text-xs text-info">+25% просмотров</p>
                </div>
                <div className="text-center p-3 bg-base-100 rounded-lg">
                  <p className="text-2xl font-bold">09:00</p>
                  <p className="text-xs">Среднее время</p>
                  <p className="text-xs">+10% просмотров</p>
                </div>
              </div>
            </div>
          </div>

          {/* Social media preview */}
          <div className="card bg-base-200 mb-6">
            <div className="card-body">
              <h3 className="card-title">
                <Share2 className="w-5 h-5" />
                Готовые посты для соцсетей
              </h3>
              <div className="grid grid-cols-1 lg:grid-cols-3 gap-4">
                {Object.entries(aiData.socialPosts).map(([platform, post]) => (
                  <div key={platform} className="border rounded-lg p-3">
                    <div className="flex items-center justify-between mb-2">
                      <span className="font-semibold text-sm capitalize">
                        {platform}
                      </span>
                      <button className="btn btn-ghost btn-xs">
                        <Send className="w-3 h-3" />
                      </button>
                    </div>
                    <p className="text-xs whitespace-pre-wrap">{post}</p>
                  </div>
                ))}
              </div>
            </div>
          </div>

          {/* Action buttons */}
          <div className="flex gap-3">
            <button
              onClick={() => setCurrentView('publish')}
              className="btn btn-primary btn-lg flex-1"
            >
              Опубликовать с AI-оптимизацией
              <Brain className="w-5 h-5 ml-1" />
            </button>
            <button className="btn btn-outline btn-lg">
              <Calendar className="w-5 h-5 mr-1" />
              Запланировать на {aiData.publishTime}
            </button>
          </div>
        </div>
      </div>
    </motion.div>
  );

  const renderPublishView = () => (
    <motion.div
      initial={{ opacity: 0 }}
      animate={{ opacity: 1 }}
      exit={{ opacity: 0 }}
      className="min-h-screen bg-gradient-to-br from-base-100 to-base-200"
    >
      <div className="container mx-auto px-4 py-16">
        <motion.div
          initial={{ scale: 0 }}
          animate={{ scale: 1 }}
          transition={{ type: 'spring', stiffness: 200 }}
          className="text-center mb-12"
        >
          <div className="inline-flex items-center justify-center w-24 h-24 bg-gradient-to-br from-success to-success/80 rounded-full mb-6">
            <Check className="w-12 h-12 text-success-content" />
          </div>
          <h1 className="text-3xl font-bold mb-4">
            Ваше объявление опубликовано!
          </h1>
          <p className="text-xl text-base-content/70 mb-8">
            AI будет оптимизировать его для максимальных продаж
          </p>
        </motion.div>

        {/* AI Features */}
        <div className="max-w-4xl mx-auto grid grid-cols-1 lg:grid-cols-2 gap-6 mb-12">
          <motion.div
            initial={{ x: -20, opacity: 0 }}
            animate={{ x: 0, opacity: 1 }}
            transition={{ delay: 0.2 }}
            className="card bg-gradient-to-br from-primary/10 to-primary/5 border-2 border-primary/20"
          >
            <div className="card-body">
              <BarChart3 className="w-10 h-10 text-primary mb-4" />
              <h3 className="card-title">AI мониторинг</h3>
              <p className="text-sm">
                AI отслеживает эффективность и автоматически корректирует цену
                для быстрой продажи
              </p>
            </div>
          </motion.div>

          <motion.div
            initial={{ x: 20, opacity: 0 }}
            animate={{ x: 0, opacity: 1 }}
            transition={{ delay: 0.3 }}
            className="card bg-gradient-to-br from-secondary/10 to-secondary/5 border-2 border-secondary/20"
          >
            <div className="card-body">
              <TrendingUp className="w-10 h-10 text-secondary mb-4" />
              <h3 className="card-title">Умное продвижение</h3>
              <p className="text-sm">
                Автоматическое поднятие в топ в оптимальное время для вашей
                категории
              </p>
            </div>
          </motion.div>

          <motion.div
            initial={{ x: -20, opacity: 0 }}
            animate={{ x: 0, opacity: 1 }}
            transition={{ delay: 0.4 }}
            className="card bg-gradient-to-br from-success/10 to-success/5 border-2 border-success/20"
          >
            <div className="card-body">
              <TestTube2 className="w-10 h-10 text-success mb-4" />
              <h3 className="card-title">A/B тестирование</h3>
              <p className="text-sm">
                {abTestEnabled
                  ? 'Активно: тестируем разные заголовки для лучшей конверсии'
                  : 'Доступно для активации в любой момент'}
              </p>
            </div>
          </motion.div>

          <motion.div
            initial={{ x: 20, opacity: 0 }}
            animate={{ x: 0, opacity: 1 }}
            transition={{ delay: 0.5 }}
            className="card bg-gradient-to-br from-warning/10 to-warning/5 border-2 border-warning/20"
          >
            <div className="card-body">
              <Globe className="w-10 h-10 text-warning mb-4" />
              <h3 className="card-title">Мультиязычность</h3>
              <p className="text-sm">
                Ваше объявление автоматически показывается на 3 языках для
                максимального охвата
              </p>
            </div>
          </motion.div>
        </div>

        {/* Stats preview */}
        <motion.div
          initial={{ y: 20, opacity: 0 }}
          animate={{ y: 0, opacity: 1 }}
          transition={{ delay: 0.6 }}
          className="max-w-4xl mx-auto card bg-base-100 shadow-xl mb-8"
        >
          <div className="card-body">
            <h3 className="card-title mb-4">
              Прогноз эффективности (AI анализ)
            </h3>
            <div className="grid grid-cols-2 lg:grid-cols-4 gap-4 text-center">
              <div>
                <p className="text-3xl font-bold text-primary">2,450</p>
                <p className="text-sm text-base-content/60">
                  просмотров за неделю
                </p>
              </div>
              <div>
                <p className="text-3xl font-bold text-secondary">125</p>
                <p className="text-sm text-base-content/60">
                  добавлений в избранное
                </p>
              </div>
              <div>
                <p className="text-3xl font-bold text-success">45</p>
                <p className="text-sm text-base-content/60">сообщений</p>
              </div>
              <div>
                <p className="text-3xl font-bold text-warning">5-7</p>
                <p className="text-sm text-base-content/60">дней до продажи</p>
              </div>
            </div>
          </div>
        </motion.div>

        {/* Action buttons */}
        <motion.div
          initial={{ y: 20, opacity: 0 }}
          animate={{ y: 0, opacity: 1 }}
          transition={{ delay: 0.7 }}
          className="max-w-4xl mx-auto flex flex-col lg:flex-row gap-4"
        >
          <Link
            href="/ru/my-listings"
            className="btn btn-primary btn-lg flex-1"
          >
            <Eye className="w-5 h-5 mr-2" />
            Посмотреть объявление
          </Link>
          <button className="btn btn-outline btn-lg flex-1">
            <Plus className="w-5 h-5 mr-2" />
            Создать еще одно
          </button>
          <button className="btn btn-ghost btn-lg flex-1">
            <Share2 className="w-5 h-5 mr-2" />
            Поделиться
          </button>
        </motion.div>

        {/* Social proof */}
        <motion.div
          initial={{ opacity: 0 }}
          animate={{ opacity: 1 }}
          transition={{ delay: 0.8 }}
          className="text-center mt-12"
        >
          <div className="flex items-center justify-center gap-2 text-sm text-base-content/60">
            <Users className="w-4 h-4" />
            <span>
              <span className="font-semibold">15,234</span> объявлений создано с
              AI за последний месяц
            </span>
          </div>
          <div className="flex items-center justify-center gap-2 mt-2 text-sm text-base-content/60">
            <ThumbsUp className="w-4 h-4" />
            <span>
              <span className="font-semibold">92%</span> продаются быстрее
              обычных
            </span>
          </div>
        </motion.div>
      </div>
    </motion.div>
  );

  return (
    <>
      {/* Navigation Bar */}
      <div className="navbar bg-base-100 border-b border-base-200 fixed top-0 z-50">
        <div className="flex-1">
          <Link
            href="/ru/examples/listing-creation-ux-v2"
            className="btn btn-ghost"
          >
            <ChevronLeft className="w-5 h-5" />
            Назад к примерам
          </Link>
        </div>
        <div className="flex-none">
          <div className="badge badge-primary badge-lg gap-1">
            <Brain className="w-4 h-4" />
            AI-Powered
          </div>
        </div>
      </div>

      {/* Main Content */}
      <div className="pt-16">
        <AnimatePresence mode="wait">
          {currentView === 'upload' && renderUploadView()}
          {currentView === 'process' && renderProcessView()}
          {currentView === 'enhance' && renderEnhanceView()}
          {currentView === 'publish' && renderPublishView()}
        </AnimatePresence>
      </div>
    </>
  );
}
