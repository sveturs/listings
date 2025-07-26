'use client';

import React, { useState, useRef } from 'react';
import Link from 'next/link';
import Image from 'next/image';
import { motion, AnimatePresence, useAnimation } from 'framer-motion';
import {
  ChevronLeft,
  Camera,
  Sparkles,
  ArrowRight,
  Check,
  MapPin,
  Mic,
  Brain,
  Image as ImageIcon,
  Wand2,
  RefreshCw,
  TrendingUp,
  Target,
  Globe,
  Calendar,
  MessageSquare,
  BarChart3,
  Rocket,
  Bot,
  Eye,
  Heart,
  Share2,
  Clock,
  Plus,
  X,
  Loader2,
} from 'lucide-react';

// Simulated AI responses
const AI_SUGGESTIONS = {
  iphone: {
    title: 'iPhone 13 Pro, 256GB, Pacific Blue',
    category: 'smartphones',
    subcategory: 'apple',
    price: { min: 60000, max: 70000, recommended: 65000 },
    description:
      'Продаю iPhone 13 Pro в отличном состоянии. Телефон использовался аккуратно, всегда в чехле и с защитным стеклом. Батарея держит отлично - 89% здоровья. В комплекте оригинальная коробка, документы и зарядное устройство. Причина продажи - переход на новую модель.',
    attributes: {
      brand: 'Apple',
      model: 'iPhone 13 Pro',
      memory: '256GB',
      color: 'Pacific Blue',
      condition: 'excellent',
      battery_health: '89%',
    },
    tags: ['premium', 'flagship', 'camera', '5G'],
    similar_listings: [
      { title: 'iPhone 13 Pro 128GB', price: 58000, sold_in_days: 3 },
      { title: 'iPhone 13 Pro 512GB', price: 72000, sold_in_days: 5 },
    ],
  },
};

export default function AIPoweredListingCreationPage() {
  const [stage, setStage] = useState<
    'welcome' | 'capture' | 'enhance' | 'publish'
  >('welcome');
  const [isProcessing, setIsProcessing] = useState(false);
  const [voiceActive, setVoiceActive] = useState(false);
  const [formData, setFormData] = useState({
    images: [] as string[],
    title: '',
    description: '',
    price: '',
    category: '',
    location: '',
    attributes: {} as any,
    tags: [] as string[],
    aiScore: 0,
  });
  const [aiSuggestions, setAiSuggestions] = useState<any>(null);
  const fileInputRef = useRef<HTMLInputElement>(null);
  const controls = useAnimation();

  // Simulate AI processing
  const processWithAI = async () => {
    setIsProcessing(true);

    // Simulate AI analysis delay
    await new Promise((resolve) => setTimeout(resolve, 2000));

    // Set AI suggestions
    setAiSuggestions(AI_SUGGESTIONS.iphone);
    setFormData({
      ...formData,
      title: AI_SUGGESTIONS.iphone.title,
      description: AI_SUGGESTIONS.iphone.description,
      price: AI_SUGGESTIONS.iphone.price.recommended.toString(),
      category: AI_SUGGESTIONS.iphone.category,
      attributes: AI_SUGGESTIONS.iphone.attributes,
      tags: AI_SUGGESTIONS.iphone.tags,
      aiScore: 95,
    });

    setIsProcessing(false);
    setStage('enhance');
  };

  const handleImageUpload = async (e: React.ChangeEvent<HTMLInputElement>) => {
    const files = e.target.files;
    if (files) {
      const newImages = Array.from(files).map((file) =>
        URL.createObjectURL(file)
      );
      setFormData({
        ...formData,
        images: [...formData.images, ...newImages].slice(0, 8),
      });

      if (newImages.length > 0) {
        setStage('capture');
        // Automatically start AI processing after image upload
        setTimeout(() => processWithAI(), 500);
      }
    }
  };

  const regenerateDescription = () => {
    controls.start({ rotate: 360 });
    setFormData({
      ...formData,
      description:
        'Предлагаю вашему вниманию iPhone 13 Pro в идеальном состоянии. Флагманский смартфон с профессиональной системой камер, мощным процессором A15 Bionic и дисплеем ProMotion 120Hz. Телефон прошел полную проверку, все функции работают безупречно. Идеальный выбор для тех, кто ценит качество и производительность.',
    });
  };

  const renderWelcomeStage = () => (
    <motion.div
      initial={{ opacity: 0 }}
      animate={{ opacity: 1 }}
      exit={{ opacity: 0 }}
      className="min-h-screen bg-gradient-to-br from-primary/5 via-base-100 to-secondary/5"
    >
      <div className="container mx-auto px-4 py-12">
        {/* Hero Section */}
        <motion.div
          initial={{ y: 20, opacity: 0 }}
          animate={{ y: 0, opacity: 1 }}
          className="text-center mb-12"
        >
          <motion.div
            animate={{
              scale: [1, 1.1, 1],
              rotate: [0, 5, -5, 0],
            }}
            transition={{
              duration: 3,
              repeat: Infinity,
              repeatType: 'reverse',
            }}
            className="inline-block mb-6"
          >
            <div className="w-24 h-24 bg-gradient-to-br from-primary to-secondary rounded-full flex items-center justify-center">
              <Brain className="w-12 h-12 text-primary-content" />
            </div>
          </motion.div>

          <h1 className="text-5xl lg:text-6xl font-bold mb-4">
            <span className="bg-gradient-to-r from-primary to-secondary bg-clip-text text-transparent">
              AI-Powered
            </span>{' '}
            Листинг
          </h1>
          <p className="text-xl text-base-content/70 mb-8 max-w-2xl mx-auto">
            Создайте идеальное объявление за 30 секунд с помощью искусственного
            интеллекта
          </p>

          {/* AI Features Grid */}
          <div className="grid grid-cols-2 lg:grid-cols-4 gap-4 max-w-4xl mx-auto mb-12">
            {[
              {
                icon: Camera,
                title: 'Умное распознавание',
                desc: 'AI определит товар по фото',
              },
              {
                icon: Wand2,
                title: 'Автогенерация',
                desc: 'Название и описание за секунду',
              },
              {
                icon: Target,
                title: 'Точная цена',
                desc: 'Анализ рынка в реальном времени',
              },
              {
                icon: Rocket,
                title: 'SEO оптимизация',
                desc: 'Максимум просмотров',
              },
            ].map((feature, index) => (
              <motion.div
                key={index}
                initial={{ y: 20, opacity: 0 }}
                animate={{ y: 0, opacity: 1 }}
                transition={{ delay: index * 0.1 }}
                className="card bg-base-200 hover:shadow-lg transition-shadow"
              >
                <div className="card-body items-center text-center p-4">
                  <feature.icon className="w-8 h-8 text-primary mb-2" />
                  <h3 className="font-bold text-sm">{feature.title}</h3>
                  <p className="text-xs text-base-content/60">{feature.desc}</p>
                </div>
              </motion.div>
            ))}
          </div>

          {/* CTA Buttons */}
          <div className="flex flex-col lg:flex-row gap-4 justify-center items-center">
            <label
              htmlFor="ai-upload"
              className="btn btn-primary btn-lg gap-2 cursor-pointer"
            >
              <Camera className="w-5 h-5" />
              Начать с фото
              <div className="badge badge-secondary">AI</div>
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

            <button
              onClick={() => setVoiceActive(true)}
              className="btn btn-outline btn-lg gap-2"
            >
              <Mic className="w-5 h-5" />
              Голосовой ввод
            </button>
          </div>
        </motion.div>

        {/* Live Stats */}
        <motion.div
          initial={{ opacity: 0 }}
          animate={{ opacity: 1 }}
          transition={{ delay: 0.5 }}
          className="grid grid-cols-3 gap-4 max-w-2xl mx-auto"
        >
          {[
            { label: 'Обработано сегодня', value: '12,847', icon: BarChart3 },
            { label: 'Точность AI', value: '98.5%', icon: Target },
            { label: 'Экономия времени', value: '92%', icon: Clock },
          ].map((stat, index) => (
            <div key={index} className="text-center">
              <stat.icon className="w-6 h-6 text-primary mx-auto mb-2" />
              <div className="text-2xl font-bold">{stat.value}</div>
              <div className="text-xs text-base-content/60">{stat.label}</div>
            </div>
          ))}
        </motion.div>
      </div>
    </motion.div>
  );

  const renderCaptureStage = () => (
    <motion.div
      initial={{ opacity: 0 }}
      animate={{ opacity: 1 }}
      exit={{ opacity: 0 }}
      className="min-h-screen bg-base-100"
    >
      <div className="container mx-auto px-4 py-8">
        <div className="max-w-4xl mx-auto">
          {/* AI Processing Animation */}
          {isProcessing && (
            <motion.div
              initial={{ opacity: 0 }}
              animate={{ opacity: 1 }}
              className="fixed inset-0 bg-base-100/90 backdrop-blur z-50 flex items-center justify-center"
            >
              <div className="text-center">
                <motion.div
                  animate={{
                    scale: [1, 1.2, 1],
                    rotate: [0, 180, 360],
                  }}
                  transition={{
                    duration: 2,
                    repeat: Infinity,
                  }}
                  className="w-20 h-20 bg-gradient-to-br from-primary to-secondary rounded-full flex items-center justify-center mb-4 mx-auto"
                >
                  <Brain className="w-10 h-10 text-primary-content" />
                </motion.div>
                <h2 className="text-xl font-bold mb-2">
                  AI анализирует ваш товар...
                </h2>
                <p className="text-base-content/70">
                  Это займет всего несколько секунд
                </p>

                {/* Progress Steps */}
                <div className="mt-8 space-y-2 text-left max-w-xs mx-auto">
                  {[
                    'Распознавание объекта',
                    'Анализ состояния',
                    'Поиск похожих товаров',
                    'Генерация описания',
                  ].map((step, index) => (
                    <motion.div
                      key={index}
                      initial={{ opacity: 0, x: -20 }}
                      animate={{ opacity: 1, x: 0 }}
                      transition={{ delay: index * 0.5 }}
                      className="flex items-center gap-2"
                    >
                      <Loader2 className="w-4 h-4 animate-spin text-primary" />
                      <span className="text-sm">{step}</span>
                    </motion.div>
                  ))}
                </div>
              </div>
            </motion.div>
          )}

          {/* Image Upload Area */}
          <div className="card bg-base-200 mb-6">
            <div className="card-body">
              <h2 className="card-title">
                <Camera className="w-5 h-5" />
                Фотографии товара
                {formData.images.length > 0 && (
                  <div className="badge badge-primary">
                    {formData.images.length}/8
                  </div>
                )}
              </h2>

              <div className="grid grid-cols-2 lg:grid-cols-4 gap-4">
                {formData.images.map((img, index) => (
                  <motion.div
                    key={index}
                    initial={{ scale: 0 }}
                    animate={{ scale: 1 }}
                    className="relative aspect-square group"
                  >
                    <Image
                      src={img}
                      alt={`Photo ${index + 1}`}
                      fill
                      className="object-cover rounded-lg"
                    />
                    {index === 0 && (
                      <div className="absolute top-2 left-2 badge badge-primary badge-sm">
                        Главное
                      </div>
                    )}
                    <button
                      onClick={() => {
                        setFormData({
                          ...formData,
                          images: formData.images.filter((_, i) => i !== index),
                        });
                      }}
                      className="absolute top-2 right-2 btn btn-circle btn-xs btn-error opacity-0 group-hover:opacity-100 transition-opacity"
                    >
                      <X className="w-3 h-3" />
                    </button>
                  </motion.div>
                ))}

                {formData.images.length < 8 && (
                  <label className="aspect-square border-2 border-dashed border-base-300 rounded-lg flex flex-col items-center justify-center cursor-pointer hover:border-primary transition-colors">
                    <Plus className="w-8 h-8 text-base-content/30" />
                    <span className="text-sm text-base-content/50 mt-2">
                      Добавить
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

              {/* AI Tips */}
              <div className="alert alert-info mt-4">
                <Bot className="w-5 h-5" />
                <div>
                  <h3 className="font-bold">Совет от AI</h3>
                  <p className="text-sm">
                    Сделайте фото на светлом фоне и покажите все важные детали -
                    это поможет AI точнее определить товар
                  </p>
                </div>
              </div>
            </div>
          </div>
        </div>
      </div>
    </motion.div>
  );

  const renderEnhanceStage = () => (
    <motion.div
      initial={{ opacity: 0 }}
      animate={{ opacity: 1 }}
      exit={{ opacity: 0 }}
      className="min-h-screen bg-base-100"
    >
      <div className="container mx-auto px-4 py-8">
        <div className="max-w-6xl mx-auto">
          {/* AI Success Banner */}
          <motion.div
            initial={{ y: -20, opacity: 0 }}
            animate={{ y: 0, opacity: 1 }}
            className="alert alert-success mb-6"
          >
            <Check className="w-5 h-5" />
            <div className="flex-1">
              <h3 className="font-bold">
                AI успешно проанализировал ваш товар!
              </h3>
              <p className="text-sm">
                Мы создали оптимизированное объявление на основе анализа{' '}
                {aiSuggestions?.similar_listings?.length || 0} похожих товаров
              </p>
            </div>
            <div className="flex items-center gap-2">
              <div
                className="radial-progress text-success"
                style={{ '--value': formData.aiScore } as any}
              >
                {formData.aiScore}%
              </div>
              <div className="text-sm">
                <div className="font-bold">AI Score</div>
                <div className="text-xs opacity-70">Отлично!</div>
              </div>
            </div>
          </motion.div>

          <div className="grid grid-cols-1 lg:grid-cols-3 gap-6">
            {/* Main Form */}
            <div className="lg:col-span-2 space-y-6">
              {/* Title with AI Enhancement */}
              <div className="card bg-base-200">
                <div className="card-body">
                  <div className="flex items-center justify-between mb-2">
                    <label className="label">
                      <span className="label-text font-bold text-lg">
                        Название
                      </span>
                    </label>
                    <button className="btn btn-ghost btn-xs gap-1">
                      <RefreshCw className="w-3 h-3" />
                      Другой вариант
                    </button>
                  </div>
                  <input
                    type="text"
                    className="input input-bordered input-lg"
                    value={formData.title}
                    onChange={(e) =>
                      setFormData({ ...formData, title: e.target.value })
                    }
                  />
                  <div className="flex gap-2 mt-2">
                    {aiSuggestions?.tags?.map((tag: string) => (
                      <span key={tag} className="badge badge-primary badge-sm">
                        {tag}
                      </span>
                    ))}
                  </div>
                </div>
              </div>

              {/* AI-Generated Description */}
              <div className="card bg-base-200">
                <div className="card-body">
                  <div className="flex items-center justify-between mb-2">
                    <label className="label">
                      <span className="label-text font-bold text-lg">
                        Описание
                      </span>
                    </label>
                    <motion.button
                      animate={controls}
                      onClick={regenerateDescription}
                      className="btn btn-ghost btn-xs gap-1"
                    >
                      <Wand2 className="w-3 h-3" />
                      Переписать
                    </motion.button>
                  </div>
                  <textarea
                    className="textarea textarea-bordered h-32"
                    value={formData.description}
                    onChange={(e) =>
                      setFormData({ ...formData, description: e.target.value })
                    }
                  />
                  <div className="text-xs text-base-content/60 mt-2">
                    ✨ AI оптимизировал текст для поисковых систем
                  </div>
                </div>
              </div>

              {/* Dynamic Attributes */}
              <div className="card bg-base-200">
                <div className="card-body">
                  <h3 className="font-bold text-lg mb-4">Характеристики</h3>
                  <div className="grid grid-cols-2 gap-4">
                    {Object.entries(formData.attributes).map(([key, value]) => (
                      <div key={key} className="form-control">
                        <label className="label">
                          <span className="label-text capitalize">
                            {key.replace(/_/g, ' ')}
                          </span>
                        </label>
                        <input
                          type="text"
                          className="input input-bordered input-sm"
                          value={value as string}
                          onChange={(e) =>
                            setFormData({
                              ...formData,
                              attributes: {
                                ...formData.attributes,
                                [key]: e.target.value,
                              },
                            })
                          }
                        />
                      </div>
                    ))}
                  </div>
                </div>
              </div>

              {/* Smart Location */}
              <div className="card bg-base-200">
                <div className="card-body">
                  <h3 className="font-bold text-lg mb-4 flex items-center gap-2">
                    <MapPin className="w-5 h-5" />
                    Местоположение
                  </h3>
                  <input
                    type="text"
                    placeholder="Начните вводить адрес..."
                    className="input input-bordered"
                    value={formData.location}
                    onChange={(e) =>
                      setFormData({ ...formData, location: e.target.value })
                    }
                  />
                  <div className="form-control mt-4">
                    <label className="label cursor-pointer">
                      <span className="label-text">
                        Скрыть точный адрес до сделки
                      </span>
                      <input
                        type="checkbox"
                        className="toggle toggle-primary"
                        defaultChecked
                      />
                    </label>
                  </div>
                </div>
              </div>
            </div>

            {/* AI Insights Sidebar */}
            <div className="space-y-6">
              {/* Price Intelligence */}
              <div className="card bg-gradient-to-br from-primary/10 to-primary/5 border-2 border-primary/20">
                <div className="card-body">
                  <h3 className="font-bold text-lg mb-4 flex items-center gap-2">
                    <TrendingUp className="w-5 h-5 text-primary" />
                    AI Ценовой анализ
                  </h3>

                  <div className="text-center mb-4">
                    <div className="text-3xl font-bold text-primary">
                      {formData.price} РСД
                    </div>
                    <div className="text-sm text-base-content/60">
                      Рекомендуемая цена
                    </div>
                  </div>

                  <div className="space-y-2 mb-4">
                    <div className="flex justify-between text-sm">
                      <span>Минимум</span>
                      <span className="font-semibold">
                        {aiSuggestions?.price?.min} РСД
                      </span>
                    </div>
                    <div className="flex justify-between text-sm">
                      <span>Максимум</span>
                      <span className="font-semibold">
                        {aiSuggestions?.price?.max} РСД
                      </span>
                    </div>
                  </div>

                  <input
                    type="range"
                    min={aiSuggestions?.price?.min}
                    max={aiSuggestions?.price?.max}
                    value={formData.price}
                    onChange={(e) =>
                      setFormData({ ...formData, price: e.target.value })
                    }
                    className="range range-primary"
                  />
                </div>
              </div>

              {/* Similar Listings */}
              <div className="card bg-base-200">
                <div className="card-body">
                  <h3 className="font-bold mb-4">Похожие объявления</h3>
                  <div className="space-y-3">
                    {aiSuggestions?.similar_listings?.map(
                      (listing: any, index: number) => (
                        <div
                          key={index}
                          className="flex justify-between items-center"
                        >
                          <div>
                            <div className="text-sm font-medium">
                              {listing.title}
                            </div>
                            <div className="text-xs text-base-content/60">
                              Продано за {listing.sold_in_days} дней
                            </div>
                          </div>
                          <div className="text-sm font-bold">
                            {listing.price} РСД
                          </div>
                        </div>
                      )
                    )}
                  </div>
                </div>
              </div>

              {/* Performance Prediction */}
              <div className="card bg-base-200">
                <div className="card-body">
                  <h3 className="font-bold mb-4 flex items-center gap-2">
                    <BarChart3 className="w-5 h-5" />
                    Прогноз эффективности
                  </h3>
                  <div className="space-y-3">
                    <div>
                      <div className="flex justify-between mb-1">
                        <span className="text-sm">Просмотры в день</span>
                        <span className="text-sm font-bold">120-150</span>
                      </div>
                      <progress
                        className="progress progress-primary"
                        value="80"
                        max="100"
                      ></progress>
                    </div>
                    <div>
                      <div className="flex justify-between mb-1">
                        <span className="text-sm">Вероятность продажи</span>
                        <span className="text-sm font-bold">92%</span>
                      </div>
                      <progress
                        className="progress progress-success"
                        value="92"
                        max="100"
                      ></progress>
                    </div>
                    <div>
                      <div className="flex justify-between mb-1">
                        <span className="text-sm">Время до продажи</span>
                        <span className="text-sm font-bold">3-5 дней</span>
                      </div>
                      <progress
                        className="progress progress-warning"
                        value="70"
                        max="100"
                      ></progress>
                    </div>
                  </div>
                </div>
              </div>

              {/* AI Tips */}
              <div className="card bg-gradient-to-br from-secondary/10 to-secondary/5 border-2 border-secondary/20">
                <div className="card-body">
                  <h3 className="font-bold mb-3 flex items-center gap-2">
                    <Bot className="w-5 h-5 text-secondary" />
                    AI Рекомендации
                  </h3>
                  <ul className="space-y-2 text-sm">
                    <li className="flex gap-2">
                      <Check className="w-4 h-4 text-success flex-shrink-0 mt-0.5" />
                      <span>Добавьте фото упаковки для +15% доверия</span>
                    </li>
                    <li className="flex gap-2">
                      <Check className="w-4 h-4 text-success flex-shrink-0 mt-0.5" />
                      <span>Укажите причину продажи в описании</span>
                    </li>
                    <li className="flex gap-2">
                      <Check className="w-4 h-4 text-success flex-shrink-0 mt-0.5" />
                      <span>Оптимальное время публикации: 19:00-21:00</span>
                    </li>
                  </ul>
                </div>
              </div>
            </div>
          </div>

          {/* Action Buttons */}
          <div className="flex gap-3 mt-8">
            <button
              onClick={() => setStage('publish')}
              className="btn btn-primary btn-lg flex-1"
            >
              Предпросмотр
              <ArrowRight className="w-5 h-5 ml-1" />
            </button>
            <button className="btn btn-ghost btn-lg">Сохранить черновик</button>
          </div>
        </div>
      </div>
    </motion.div>
  );

  const renderPublishStage = () => (
    <motion.div
      initial={{ opacity: 0 }}
      animate={{ opacity: 1 }}
      exit={{ opacity: 0 }}
      className="min-h-screen bg-gradient-to-br from-base-200 to-base-100"
    >
      <div className="container mx-auto px-4 py-8">
        <div className="max-w-6xl mx-auto">
          {/* Success Header */}
          <motion.div
            initial={{ scale: 0 }}
            animate={{ scale: 1 }}
            transition={{ type: 'spring', stiffness: 200 }}
            className="text-center mb-8"
          >
            <div className="inline-flex items-center justify-center w-24 h-24 bg-gradient-to-br from-success/20 to-success/10 rounded-full mb-4">
              <Check className="w-12 h-12 text-success" />
            </div>
            <h1 className="text-3xl font-bold mb-2">
              Великолепно! Объявление готово
            </h1>
            <p className="text-base-content/70">
              AI оптимизировал ваше объявление для максимальной эффективности
            </p>
          </motion.div>

          <div className="grid grid-cols-1 lg:grid-cols-3 gap-6">
            {/* Preview Card */}
            <div className="lg:col-span-2">
              <motion.div
                initial={{ y: 20, opacity: 0 }}
                animate={{ y: 0, opacity: 1 }}
                transition={{ delay: 0.1 }}
                className="card bg-base-100 shadow-xl"
              >
                {formData.images.length > 0 && (
                  <figure className="relative">
                    <div className="relative w-full h-96">
                      <Image
                        src={formData.images[0]}
                        alt={formData.title}
                        fill
                        className="object-cover"
                      />
                    </div>
                    <div className="absolute top-4 left-4 badge badge-primary gap-1">
                      <Sparkles className="w-3 h-3" />
                      AI Enhanced
                    </div>
                    {formData.images.length > 1 && (
                      <div className="absolute bottom-4 right-4 badge badge-neutral gap-1">
                        <ImageIcon className="w-3 h-3" />+
                        {formData.images.length - 1}
                      </div>
                    )}
                  </figure>
                )}

                <div className="card-body">
                  <h2 className="card-title text-2xl">{formData.title}</h2>

                  <div className="flex flex-wrap gap-2 mb-3">
                    {formData.tags.map((tag) => (
                      <span
                        key={tag}
                        className="badge badge-secondary badge-sm"
                      >
                        {tag}
                      </span>
                    ))}
                  </div>

                  <div className="text-3xl font-bold text-primary mb-4">
                    {formData.price} РСД
                  </div>

                  <p className="text-base-content/80 mb-4">
                    {formData.description}
                  </p>

                  <div className="grid grid-cols-2 gap-4 mb-4">
                    {Object.entries(formData.attributes)
                      .slice(0, 4)
                      .map(([key, value]) => (
                        <div
                          key={key}
                          className="flex justify-between py-2 border-b border-base-200"
                        >
                          <span className="text-sm text-base-content/60 capitalize">
                            {key.replace(/_/g, ' ')}
                          </span>
                          <span className="text-sm font-medium">
                            {value as string}
                          </span>
                        </div>
                      ))}
                  </div>

                  <div className="flex items-center gap-4 text-sm text-base-content/60">
                    <span className="flex items-center gap-1">
                      <MapPin className="w-4 h-4" />
                      {formData.location || 'Белград'}
                    </span>
                    <span className="flex items-center gap-1">
                      <Eye className="w-4 h-4" />0 просмотров
                    </span>
                    <span className="flex items-center gap-1">
                      <Heart className="w-4 h-4" />0 в избранном
                    </span>
                  </div>

                  <div className="divider"></div>

                  <div className="flex gap-2">
                    <button className="btn btn-primary flex-1">
                      <MessageSquare className="w-4 h-4 mr-1" />
                      Написать продавцу
                    </button>
                    <button className="btn btn-ghost">
                      <Heart className="w-4 h-4" />
                    </button>
                    <button className="btn btn-ghost">
                      <Share2 className="w-4 h-4" />
                    </button>
                  </div>
                </div>
              </motion.div>
            </div>

            {/* Publishing Options */}
            <div className="space-y-6">
              {/* AI Score Card */}
              <motion.div
                initial={{ x: 20, opacity: 0 }}
                animate={{ x: 0, opacity: 1 }}
                transition={{ delay: 0.2 }}
                className="card bg-gradient-to-br from-primary/10 to-secondary/10"
              >
                <div className="card-body text-center">
                  <h3 className="font-bold text-lg mb-4">AI Оценка качества</h3>
                  <div
                    className="radial-progress text-primary mx-auto mb-4"
                    style={
                      { '--value': formData.aiScore, '--size': '8rem' } as any
                    }
                  >
                    <span className="text-2xl font-bold">
                      {formData.aiScore}%
                    </span>
                  </div>
                  <p className="text-sm text-base-content/70">
                    Ваше объявление лучше 95% других в этой категории
                  </p>
                </div>
              </motion.div>

              {/* Publishing Plans */}
              <motion.div
                initial={{ x: 20, opacity: 0 }}
                animate={{ x: 0, opacity: 1 }}
                transition={{ delay: 0.3 }}
                className="space-y-3"
              >
                {/* Free Plan */}
                <div className="card bg-base-200">
                  <div className="card-body">
                    <h4 className="font-bold">Бесплатная публикация</h4>
                    <ul className="text-sm space-y-1 mb-4">
                      <li className="flex items-center gap-2">
                        <Check className="w-4 h-4 text-success" />
                        Базовое размещение
                      </li>
                      <li className="flex items-center gap-2">
                        <Check className="w-4 h-4 text-success" />
                        До 8 фотографий
                      </li>
                    </ul>
                    <button className="btn btn-outline btn-block">
                      Опубликовать бесплатно
                    </button>
                  </div>
                </div>

                {/* Premium Plan */}
                <div className="card bg-gradient-to-br from-warning/20 to-warning/10 border-2 border-warning">
                  <div className="card-body">
                    <div className="flex items-center justify-between mb-2">
                      <h4 className="font-bold">Premium размещение</h4>
                      <div className="badge badge-warning">Популярно</div>
                    </div>
                    <ul className="text-sm space-y-1 mb-4">
                      <li className="flex items-center gap-2">
                        <Check className="w-4 h-4 text-success" />
                        Топ выдачи на 7 дней
                      </li>
                      <li className="flex items-center gap-2">
                        <Check className="w-4 h-4 text-success" />
                        5x больше просмотров
                      </li>
                      <li className="flex items-center gap-2">
                        <Check className="w-4 h-4 text-success" />
                        Выделение цветом
                      </li>
                      <li className="flex items-center gap-2">
                        <Check className="w-4 h-4 text-success" />
                        AI продвижение
                      </li>
                    </ul>
                    <div className="text-center mb-4">
                      <div className="text-2xl font-bold">199 РСД</div>
                      <div className="text-xs text-base-content/60">
                        на 7 дней
                      </div>
                    </div>
                    <button className="btn btn-warning btn-block">
                      <Rocket className="w-4 h-4 mr-1" />
                      Разместить в топе
                    </button>
                  </div>
                </div>
              </motion.div>

              {/* AI Recommendations */}
              <motion.div
                initial={{ x: 20, opacity: 0 }}
                animate={{ x: 0, opacity: 1 }}
                transition={{ delay: 0.4 }}
                className="card bg-base-200"
              >
                <div className="card-body">
                  <h4 className="font-bold mb-3 flex items-center gap-2">
                    <Calendar className="w-5 h-5" />
                    Лучшее время публикации
                  </h4>
                  <div className="space-y-2">
                    <div className="flex justify-between items-center p-2 bg-success/10 rounded">
                      <span className="text-sm">Сейчас</span>
                      <span className="text-sm font-bold text-success">
                        Отлично
                      </span>
                    </div>
                    <div className="flex justify-between items-center p-2">
                      <span className="text-sm">Вечером (19:00)</span>
                      <span className="text-sm text-base-content/60">
                        Хорошо
                      </span>
                    </div>
                  </div>
                </div>
              </motion.div>

              {/* Share Options */}
              <motion.div
                initial={{ x: 20, opacity: 0 }}
                animate={{ x: 0, opacity: 1 }}
                transition={{ delay: 0.5 }}
                className="card bg-base-200"
              >
                <div className="card-body">
                  <h4 className="font-bold mb-3">
                    Поделиться после публикации
                  </h4>
                  <div className="flex gap-2">
                    <button className="btn btn-sm btn-ghost">
                      <Globe className="w-4 h-4" />
                    </button>
                    <button className="btn btn-sm btn-ghost">
                      <MessageSquare className="w-4 h-4" />
                    </button>
                    <button className="btn btn-sm btn-ghost">
                      <Share2 className="w-4 h-4" />
                    </button>
                  </div>
                </div>
              </motion.div>
            </div>
          </div>
        </div>
      </div>
    </motion.div>
  );

  return (
    <>
      {/* Navigation */}
      <div className="navbar bg-base-100 border-b border-base-200 sticky top-0 z-40">
        <div className="flex-1">
          <Link
            href="/ru/examples/listing-creation-ux"
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

      {/* Voice Input Modal */}
      <AnimatePresence>
        {voiceActive && (
          <motion.div
            initial={{ opacity: 0 }}
            animate={{ opacity: 1 }}
            exit={{ opacity: 0 }}
            className="fixed inset-0 bg-base-100/90 backdrop-blur z-50 flex items-center justify-center"
            onClick={() => setVoiceActive(false)}
          >
            <motion.div
              initial={{ scale: 0 }}
              animate={{ scale: 1 }}
              exit={{ scale: 0 }}
              className="card bg-base-200 w-96"
              onClick={(e) => e.stopPropagation()}
            >
              <div className="card-body text-center">
                <motion.div
                  animate={{
                    scale: [1, 1.2, 1],
                  }}
                  transition={{
                    duration: 1.5,
                    repeat: Infinity,
                  }}
                  className="w-24 h-24 bg-primary rounded-full flex items-center justify-center mx-auto mb-4"
                >
                  <Mic className="w-12 h-12 text-primary-content" />
                </motion.div>
                <h3 className="text-xl font-bold mb-2">Говорите...</h3>
                <p className="text-base-content/70 mb-4">
                  &quot;Продаю iPhone 13 Pro, 256 гигабайт, синий цвет,
                  состояние отличное&quot;
                </p>
                <button
                  onClick={() => setVoiceActive(false)}
                  className="btn btn-ghost"
                >
                  Отмена
                </button>
              </div>
            </motion.div>
          </motion.div>
        )}
      </AnimatePresence>

      {/* Main Content */}
      <AnimatePresence mode="wait">
        {stage === 'welcome' && renderWelcomeStage()}
        {stage === 'capture' && renderCaptureStage()}
        {stage === 'enhance' && renderEnhanceStage()}
        {stage === 'publish' && renderPublishStage()}
      </AnimatePresence>
    </>
  );
}
