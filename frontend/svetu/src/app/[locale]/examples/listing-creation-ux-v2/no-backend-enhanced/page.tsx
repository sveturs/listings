'use client';

import React, { useState, useRef, useEffect } from 'react';
import Link from 'next/link';
import Image from 'next/image';
import { motion, AnimatePresence, useAnimation } from 'framer-motion';
import {
  ChevronLeft,
  Camera,
  Plus,
  X,
  Sparkles,
  Zap,
  ArrowRight,
  Check,
  MapPin,
  Package,
  Image as ImageIcon,
  Heart,
  Eye,
  MessageCircle,
  Share2,
  TrendingUp,
  Timer,
  Shield,
  Award,
  Info,
  Lightbulb,
  AlertCircle,
  Volume2,
  Instagram,
  Facebook,
  Clock as ClockIcon,
  FileText,
  Users,
} from 'lucide-react';

export default function NoBackendEnhancedListingCreationPage() {
  const [currentView, setCurrentView] = useState<
    'start' | 'create' | 'preview'
  >('start');
  const [quickMode, setQuickMode] = useState(false);
  const [formData, setFormData] = useState({
    images: [] as string[],
    category: '',
    title: '',
    price: '',
    description: '',
    location: '',
    deliveryMethods: ['pickup'],
    attributes: {} as Record<string, string>,
  });
  const [suggestions, setSuggestions] = useState({
    title: '',
    category: '',
    price: '',
    description: '',
  });

  // Состояние для сравнения с похожими
  const [showPriceComparison, setShowPriceComparison] = useState(false);
  const [similarListings, setSimilarListings] = useState<any[]>([]);

  // Состояние для шаблонов описаний
  const [_descriptionTemplate, _setDescriptionTemplate] = useState('');

  // Оптимальное время публикации
  const [_optimalPublishTime, _setOptimalPublishTime] = useState('');

  const fileInputRef = useRef<HTMLInputElement>(null);
  const controls = useAnimation();

  // Category-specific attributes
  const categoryAttributes: Record<
    string,
    Array<{ id: string; label: string; type: string; options?: string[] }>
  > = {
    electronics: [
      { id: 'brand', label: 'Бренд', type: 'text' },
      { id: 'model', label: 'Модель', type: 'text' },
      {
        id: 'condition',
        label: 'Состояние',
        type: 'select',
        options: [
          'Новый',
          'Как новый',
          'Отличное',
          'Хорошее',
          'Удовлетворительное',
        ],
      },
      {
        id: 'warranty',
        label: 'Гарантия',
        type: 'select',
        options: ['Есть', 'Нет', 'Истекла'],
      },
    ],
    fashion: [
      { id: 'brand', label: 'Бренд', type: 'text' },
      { id: 'size', label: 'Размер', type: 'text' },
      { id: 'color', label: 'Цвет', type: 'text' },
      { id: 'material', label: 'Материал', type: 'text' },
      {
        id: 'season',
        label: 'Сезон',
        type: 'select',
        options: ['Лето', 'Зима', 'Весна/Осень', 'Всесезонная'],
      },
    ],
    home: [
      { id: 'type', label: 'Тип', type: 'text' },
      { id: 'dimensions', label: 'Размеры', type: 'text' },
      { id: 'material', label: 'Материал', type: 'text' },
      {
        id: 'condition',
        label: 'Состояние',
        type: 'select',
        options: ['Новый', 'Отличное', 'Хорошее', 'Требует ремонта'],
      },
    ],
    auto: [
      { id: 'brand', label: 'Марка', type: 'text' },
      { id: 'model', label: 'Модель', type: 'text' },
      { id: 'year', label: 'Год выпуска', type: 'text' },
      { id: 'mileage', label: 'Пробег (км)', type: 'text' },
      {
        id: 'fuel',
        label: 'Топливо',
        type: 'select',
        options: ['Бензин', 'Дизель', 'Электро', 'Гибрид'],
      },
    ],
  };

  // Шаблоны описаний по категориям
  const descriptionTemplates: Record<string, string> = {
    electronics: `📱 Состояние: [отличное/хорошее/новое]
✅ Комплектация: [что входит в комплект]
📦 Причина продажи: [обновление/не используется]
🔋 Батарея держит: [время работы]
💎 Особенности: [что особенного]`,

    fashion: `👕 Размер: [точный размер]
📏 Параметры: [обхват груди/талии]
🧵 Состав: [материал]
✨ Состояние: [новое/б/у]
📸 На фото: [рост модели]`,

    home: `🏠 Размеры: [длина x ширина x высота]
📦 Состояние: [новое/б/у]
🛠️ Сборка: [требуется/не требуется]
🚚 Самовывоз: [адрес]
💡 Особенности: [что особенного]`,

    auto: `🚗 Пробег: [км]
⛽ Расход: [л/100км]
🔧 ТО: [когда было]
📋 Документы: [в порядке]
🛡️ Страховка: [до когда]`,
  };

  // Симулированные данные похожих объявлений
  const getSimilarListings = () => {
    return [
      {
        id: 1,
        title: 'iPhone 13 Pro 256GB Space Gray',
        price: 68000,
        views: 245,
        daysAgo: 2,
        sold: false,
      },
      {
        id: 2,
        title: 'iPhone 13 Pro 128GB Blue',
        price: 62000,
        views: 189,
        daysAgo: 5,
        sold: true,
      },
      {
        id: 3,
        title: 'iPhone 13 Pro 512GB Gold',
        price: 75000,
        views: 156,
        daysAgo: 1,
        sold: false,
      },
    ];
  };

  // Определение оптимального времени публикации
  const _getOptimalPublishTime = () => {
    const times = [
      { time: '19:00-21:00', activity: 'Высокая', icon: '🔥' },
      { time: '12:00-13:00', activity: 'Средняя', icon: '👍' },
      { time: '09:00-10:00', activity: 'Средняя', icon: '👍' },
    ];
    return times[0]; // Возвращаем лучшее время
  };

  // Simulated quick templates
  const quickTemplates = [
    {
      id: 'phone',
      icon: '📱',
      title: 'Продаю телефон',
      fields: ['Модель', 'Память', 'Состояние'],
    },
    {
      id: 'clothes',
      icon: '👕',
      title: 'Одежда/Обувь',
      fields: ['Размер', 'Бренд', 'Состояние'],
    },
    {
      id: 'electronics',
      icon: '💻',
      title: 'Электроника',
      fields: ['Бренд', 'Модель', 'Год'],
    },
    {
      id: 'furniture',
      icon: '🛋️',
      title: 'Мебель',
      fields: ['Тип', 'Размеры', 'Материал'],
    },
  ];

  const popularCategories = [
    {
      id: 'electronics',
      name: 'Электроника',
      icon: '📱',
      gradient: 'from-blue-500 to-purple-500',
    },
    {
      id: 'fashion',
      name: 'Мода',
      icon: '👗',
      gradient: 'from-pink-500 to-rose-500',
    },
    {
      id: 'home',
      name: 'Дом',
      icon: '🏠',
      gradient: 'from-green-500 to-emerald-500',
    },
    {
      id: 'auto',
      name: 'Авто',
      icon: '🚗',
      gradient: 'from-orange-500 to-red-500',
    },
  ];

  useEffect(() => {
    // Simulate AI suggestions when image is uploaded
    if (formData.images.length > 0 && !suggestions.title) {
      setTimeout(() => {
        setSuggestions({
          title: 'iPhone 13 Pro, 256GB, Pacific Blue',
          category: 'electronics',
          price: '65000',
          description:
            'Телефон в отличном состоянии, использовался аккуратно. Полный комплект, есть чек.',
        });

        // Показываем сравнение цен
        setSimilarListings(getSimilarListings());
        setShowPriceComparison(true);
      }, 1000);
    }
  }, [formData.images, suggestions.title]);

  // Проверка на наличие контактов в описании
  const checkForContactInfo = (text: string) => {
    const phoneRegex =
      /(\+?\d{1,3}[-.\s]?)?\(?\d{1,4}\)?[-.\s]?\d{1,4}[-.\s]?\d{1,9}/g;
    const emailRegex = /[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}/g;

    return phoneRegex.test(text) || emailRegex.test(text);
  };

  const handleImageUpload = (e: React.ChangeEvent<HTMLInputElement>) => {
    const files = e.target.files;
    if (files) {
      const newImages = Array.from(files).map((file) =>
        URL.createObjectURL(file)
      );

      // Проверка качества изображений
      newImages.forEach((imgUrl, index) => {
        const img = new window.Image();
        img.src = imgUrl;
        img.onload = () => {
          if (img.width < 800 || img.height < 600) {
            console.log(`Image ${index + 1} has low quality`);
          }
        };
      });

      setFormData({
        ...formData,
        images: [...formData.images, ...newImages].slice(0, 8),
      });
      if (newImages.length > 0) {
        setCurrentView('create');
      }
    }
  };

  const removeImage = (index: number) => {
    setFormData({
      ...formData,
      images: formData.images.filter((_, i) => i !== index),
    });
  };

  const applySuggestions = () => {
    setFormData({
      ...formData,
      title: suggestions.title,
      category: suggestions.category,
      price: suggestions.price,
      description: suggestions.description,
    });
    controls.start({
      scale: [1, 1.05, 1],
      transition: { duration: 0.3 },
    });
  };

  const applyDescriptionTemplate = () => {
    if (formData.category && descriptionTemplates[formData.category]) {
      setFormData({
        ...formData,
        description: descriptionTemplates[formData.category],
      });
    }
  };

  const renderStartView = () => (
    <motion.div
      initial={{ opacity: 0 }}
      animate={{ opacity: 1 }}
      exit={{ opacity: 0 }}
      className="min-h-screen bg-gradient-to-br from-base-100 to-base-200"
    >
      {/* Hero Section */}
      <div className="container mx-auto px-4 py-8">
        <motion.div
          initial={{ y: 20, opacity: 0 }}
          animate={{ y: 0, opacity: 1 }}
          transition={{ delay: 0.1 }}
          className="text-center mb-12"
        >
          <h1 className="text-4xl lg:text-5xl font-bold mb-4 bg-gradient-to-r from-primary to-secondary bg-clip-text text-transparent">
            Продайте быстрее с умными подсказками 🚀
          </h1>
          <p className="text-xl text-base-content/70 mb-8">
            AI-подсказки, шаблоны, сравнение цен — всё для успешной продажи
          </p>

          {/* Stats */}
          <div className="flex justify-center gap-8 mb-8">
            <motion.div
              initial={{ scale: 0 }}
              animate={{ scale: 1 }}
              transition={{ delay: 0.2 }}
              className="text-center"
            >
              <div className="text-3xl font-bold text-primary">2 мин</div>
              <div className="text-sm text-base-content/60">создание</div>
            </motion.div>
            <motion.div
              initial={{ scale: 0 }}
              animate={{ scale: 1 }}
              transition={{ delay: 0.3 }}
              className="text-center"
            >
              <div className="text-3xl font-bold text-success">98%</div>
              <div className="text-sm text-base-content/60">завершают</div>
            </motion.div>
            <motion.div
              initial={{ scale: 0 }}
              animate={{ scale: 1 }}
              transition={{ delay: 0.4 }}
              className="text-center"
            >
              <div className="text-3xl font-bold text-secondary">10x</div>
              <div className="text-sm text-base-content/60">
                больше просмотров
              </div>
            </motion.div>
          </div>
        </motion.div>

        {/* Start Options */}
        <div className="max-w-4xl mx-auto">
          {/* Primary CTA */}
          <motion.div
            initial={{ y: 20, opacity: 0 }}
            animate={{ y: 0, opacity: 1 }}
            transition={{ delay: 0.2 }}
            className="mb-8"
          >
            <label
              htmlFor="quick-upload"
              className="card bg-gradient-to-r from-primary to-secondary text-primary-content cursor-pointer hover:shadow-2xl transition-all"
            >
              <div className="card-body text-center py-12">
                <Camera className="w-16 h-16 mx-auto mb-4" />
                <h2 className="text-2xl font-bold mb-2">Начните с фото</h2>
                <p className="opacity-90 mb-4">
                  Загрузите фото товара, а мы поможем с остальным
                </p>
                <div className="flex gap-2 justify-center">
                  <div className="badge badge-lg badge-warning gap-2">
                    <Zap className="w-4 h-4" />
                    Быстрый старт
                  </div>
                  <div className="badge badge-lg badge-info gap-2">
                    <Sparkles className="w-4 h-4" />
                    AI подсказки
                  </div>
                </div>
              </div>
            </label>
            <input
              id="quick-upload"
              ref={fileInputRef}
              type="file"
              multiple
              accept="image/*"
              className="hidden"
              onChange={handleImageUpload}
            />
          </motion.div>

          {/* Social Import */}
          <motion.div
            initial={{ y: 20, opacity: 0 }}
            animate={{ y: 0, opacity: 1 }}
            transition={{ delay: 0.3 }}
            className="mb-6"
          >
            <div className="text-center mb-4">
              <h3 className="font-semibold">Импорт из соцсетей</h3>
              <p className="text-sm text-base-content/60">
                Уже выложили товар в соцсетях? Импортируйте одним кликом
              </p>
            </div>
            <div className="flex gap-2 justify-center">
              <button className="btn btn-outline gap-2">
                <Instagram className="w-4 h-4" />
                Instagram
              </button>
              <button className="btn btn-outline gap-2">
                <Facebook className="w-4 h-4" />
                Facebook
              </button>
            </div>
          </motion.div>

          {/* Alternative Options */}
          <div className="grid grid-cols-1 lg:grid-cols-2 gap-4 mb-8">
            <motion.div
              initial={{ x: -20, opacity: 0 }}
              animate={{ x: 0, opacity: 1 }}
              transition={{ delay: 0.4 }}
            >
              <button
                onClick={() => {
                  setCurrentView('create');
                  setQuickMode(false);
                }}
                className="card bg-base-100 border-2 border-base-300 hover:border-primary hover:shadow-lg transition-all w-full"
              >
                <div className="card-body flex-row items-center">
                  <Package className="w-12 h-12 text-primary mr-4" />
                  <div className="text-left">
                    <h3 className="font-bold">Классический способ</h3>
                    <p className="text-sm text-base-content/60">
                      Пошаговое создание с подсказками
                    </p>
                  </div>
                </div>
              </button>
            </motion.div>

            <motion.div
              initial={{ x: 20, opacity: 0 }}
              animate={{ x: 0, opacity: 1 }}
              transition={{ delay: 0.5 }}
            >
              <button
                onClick={() => {
                  setQuickMode(true);
                  setCurrentView('create');
                }}
                className="card bg-base-100 border-2 border-base-300 hover:border-secondary hover:shadow-lg transition-all w-full"
              >
                <div className="card-body flex-row items-center">
                  <Zap className="w-12 h-12 text-secondary mr-4" />
                  <div className="text-left">
                    <h3 className="font-bold">Супер-быстро</h3>
                    <p className="text-sm text-base-content/60">
                      Только самое необходимое
                    </p>
                  </div>
                </div>
              </button>
            </motion.div>
          </div>

          {/* Quick Templates */}
          <motion.div
            initial={{ y: 20, opacity: 0 }}
            animate={{ y: 0, opacity: 1 }}
            transition={{ delay: 0.6 }}
          >
            <h3 className="text-center font-semibold mb-4 text-base-content/70">
              Или выберите готовый шаблон
            </h3>
            <div className="grid grid-cols-2 lg:grid-cols-4 gap-3">
              {quickTemplates.map((template) => (
                <button
                  key={template.id}
                  onClick={() => {
                    setFormData({ ...formData, category: template.id });
                    setCurrentView('create');
                  }}
                  className="btn btn-outline btn-sm gap-2"
                >
                  <span className="text-xl">{template.icon}</span>
                  {template.title}
                </button>
              ))}
            </div>
          </motion.div>
        </div>
      </div>
    </motion.div>
  );

  const renderCreateView = () => (
    <motion.div
      initial={{ opacity: 0 }}
      animate={{ opacity: 1 }}
      exit={{ opacity: 0 }}
      className="min-h-screen bg-base-100"
    >
      {/* Floating Header */}
      <div className="sticky top-0 z-50 bg-base-100/80 backdrop-blur-lg border-b border-base-200">
        <div className="container mx-auto px-4 py-3">
          <div className="flex items-center justify-between">
            <button
              onClick={() => setCurrentView('start')}
              className="btn btn-ghost btn-sm gap-2"
            >
              <ChevronLeft className="w-4 h-4" />
              Назад
            </button>

            <div className="flex items-center gap-2">
              <div className="badge badge-success gap-1">
                <Timer className="w-3 h-3" />
                Автосохранение
              </div>
              {quickMode && (
                <div className="badge badge-warning gap-1">
                  <Zap className="w-3 h-3" />
                  Быстрый режим
                </div>
              )}
            </div>
          </div>
        </div>
      </div>

      {/* Main Content */}
      <div className="container mx-auto px-4 py-8">
        <div className="max-w-2xl mx-auto space-y-6">
          {/* AI Suggestions Banner */}
          {suggestions.title && formData.images.length > 0 && (
            <motion.div
              initial={{ y: -20, opacity: 0 }}
              animate={{ y: 0, opacity: 1 }}
              className="alert alert-info shadow-lg"
            >
              <Sparkles className="w-5 h-5" />
              <div className="flex-1">
                <h3 className="font-bold">Мы узнали ваш товар!</h3>
                <p className="text-sm">
                  {suggestions.title} • Рекомендуемая цена: {suggestions.price}{' '}
                  РСД
                </p>
              </div>
              <button
                onClick={applySuggestions}
                className="btn btn-sm btn-primary"
              >
                Применить
              </button>
            </motion.div>
          )}

          {/* Photo Upload Section */}
          <motion.div animate={controls} className="card bg-base-200">
            <div className="card-body">
              <h2 className="card-title">
                <Camera className="w-5 h-5" />
                Фотографии
                {formData.images.length > 0 && (
                  <span className="badge badge-primary">
                    {formData.images.length}/8
                  </span>
                )}
              </h2>

              <div className="grid grid-cols-3 lg:grid-cols-4 gap-3">
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
                      <div className="absolute top-1 left-1 badge badge-primary badge-sm">
                        Главное
                      </div>
                    )}
                    <button
                      onClick={() => removeImage(index)}
                      className="absolute top-1 right-1 btn btn-circle btn-xs btn-error opacity-0 group-hover:opacity-100 transition-opacity"
                    >
                      <X className="w-3 h-3" />
                    </button>
                  </motion.div>
                ))}

                {formData.images.length < 8 && (
                  <label className="aspect-square border-2 border-dashed border-base-300 rounded-lg flex flex-col items-center justify-center cursor-pointer hover:border-primary transition-colors">
                    <Plus className="w-6 h-6 text-base-content/50" />
                    <span className="text-xs text-base-content/50 mt-1">
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

              {/* Рекомендации по фото */}
              {formData.images.length > 0 && formData.images.length < 4 && (
                <div className="alert alert-warning mt-4">
                  <Lightbulb className="w-4 h-4" />
                  <span className="text-sm">
                    Добавьте еще {4 - formData.images.length} фото для лучших
                    продаж
                  </span>
                </div>
              )}
            </div>
          </motion.div>

          {/* Quick Info Section */}
          <div className="card bg-base-200">
            <div className="card-body space-y-4">
              {/* Title */}
              <div className="form-control">
                <label className="label">
                  <span className="label-text font-semibold">Название</span>
                  <span className="label-text-alt">
                    {formData.title.length}/80
                  </span>
                </label>
                <input
                  type="text"
                  placeholder="Что вы продаете?"
                  className="input input-bordered"
                  value={formData.title}
                  onChange={(e) =>
                    setFormData({ ...formData, title: e.target.value })
                  }
                  maxLength={80}
                />
              </div>

              {/* Category Pills */}
              {!quickMode && (
                <div>
                  <label className="label">
                    <span className="label-text font-semibold">Категория</span>
                  </label>
                  <div className="flex flex-wrap gap-2">
                    {popularCategories.map((cat) => (
                      <button
                        key={cat.id}
                        onClick={() =>
                          setFormData({ ...formData, category: cat.id })
                        }
                        className={`btn btn-sm ${
                          formData.category === cat.id
                            ? 'btn-primary'
                            : 'btn-outline'
                        } gap-1`}
                      >
                        <span>{cat.icon}</span>
                        {cat.name}
                      </button>
                    ))}
                  </div>
                </div>
              )}

              {/* Price with comparison */}
              <div className="form-control">
                <label className="label">
                  <span className="label-text font-semibold">Цена</span>
                  <button
                    onClick={() => setShowPriceComparison(!showPriceComparison)}
                    className="label-text-alt link link-primary"
                  >
                    Сравнить с похожими
                  </button>
                </label>
                <label className="input-group">
                  <input
                    type="number"
                    placeholder="0"
                    className="input input-bordered flex-1"
                    value={formData.price}
                    onChange={(e) =>
                      setFormData({ ...formData, price: e.target.value })
                    }
                  />
                  <span>РСД</span>
                </label>

                {/* Price comparison */}
                {showPriceComparison && similarListings.length > 0 && (
                  <div className="mt-4 space-y-2">
                    <h4 className="text-sm font-semibold">
                      Похожие объявления:
                    </h4>
                    {similarListings.map((listing) => (
                      <div
                        key={listing.id}
                        className="flex items-center justify-between text-sm p-2 bg-base-100 rounded"
                      >
                        <div className="flex-1">
                          <p className="font-medium">{listing.title}</p>
                          <p className="text-xs text-base-content/60">
                            <Eye className="w-3 h-3 inline mr-1" />
                            {listing.views} просмотров • {listing.daysAgo} дн.
                            назад
                          </p>
                        </div>
                        <div className="text-right">
                          <p className="font-bold">
                            {listing.price.toLocaleString()} РСД
                          </p>
                          {listing.sold && (
                            <span className="badge badge-success badge-xs">
                              Продано
                            </span>
                          )}
                        </div>
                      </div>
                    ))}
                  </div>
                )}
              </div>

              {/* Quick Description with templates */}
              {!quickMode && (
                <div className="form-control">
                  <label className="label">
                    <span className="label-text font-semibold">Описание</span>
                    <span className="label-text-alt">Опционально</span>
                  </label>
                  <div className="relative">
                    <textarea
                      className="textarea textarea-bordered h-20 w-full"
                      placeholder="Добавьте детали..."
                      value={formData.description}
                      onChange={(e) => {
                        const newDescription = e.target.value;
                        setFormData({
                          ...formData,
                          description: newDescription,
                        });

                        // Проверка на контакты
                        if (checkForContactInfo(newDescription)) {
                          console.log('Contact info detected!');
                        }
                      }}
                    />
                    <div className="absolute bottom-2 right-2 flex gap-1">
                      <button className="btn btn-xs btn-ghost gap-1">
                        <Volume2 className="w-3 h-3" />
                        Диктовка
                      </button>
                    </div>
                  </div>

                  {/* Шаблоны описаний */}
                  {formData.category && (
                    <button
                      onClick={applyDescriptionTemplate}
                      className="btn btn-outline btn-sm mt-2 gap-1"
                    >
                      <FileText className="w-3 h-3" />
                      Использовать шаблон для {formData.category}
                    </button>
                  )}

                  {/* Предупреждение о контактах */}
                  {checkForContactInfo(formData.description) && (
                    <div className="alert alert-warning mt-2">
                      <AlertCircle className="w-4 h-4" />
                      <span className="text-sm">
                        Контактные данные в описании запрещены правилами
                      </span>
                    </div>
                  )}
                </div>
              )}
            </div>
          </div>

          {/* Dynamic Attributes based on Category */}
          {formData.category && categoryAttributes[formData.category] && (
            <motion.div
              initial={{ opacity: 0, y: 20 }}
              animate={{ opacity: 1, y: 0 }}
              className="card bg-base-200"
            >
              <div className="card-body">
                <h3 className="card-title text-base">
                  <Package className="w-4 h-4" />
                  Дополнительная информация
                </h3>
                <div className="grid grid-cols-1 lg:grid-cols-2 gap-4">
                  {categoryAttributes[formData.category].map((attr) => (
                    <div key={attr.id} className="form-control">
                      <label className="label">
                        <span className="label-text">{attr.label}</span>
                      </label>
                      {attr.type === 'select' ? (
                        <select
                          className="select select-bordered select-sm"
                          value={formData.attributes[attr.id] || ''}
                          onChange={(e) =>
                            setFormData({
                              ...formData,
                              attributes: {
                                ...formData.attributes,
                                [attr.id]: e.target.value,
                              },
                            })
                          }
                        >
                          <option value="">Выберите...</option>
                          {attr.options?.map((option) => (
                            <option key={option} value={option}>
                              {option}
                            </option>
                          ))}
                        </select>
                      ) : (
                        <input
                          type="text"
                          className="input input-bordered input-sm"
                          placeholder={`Введите ${attr.label.toLowerCase()}`}
                          value={formData.attributes[attr.id] || ''}
                          onChange={(e) =>
                            setFormData({
                              ...formData,
                              attributes: {
                                ...formData.attributes,
                                [attr.id]: e.target.value,
                              },
                            })
                          }
                        />
                      )}
                    </div>
                  ))}
                </div>
              </div>
            </motion.div>
          )}

          {/* Location Card */}
          <div className="card bg-base-200">
            <div className="card-body">
              <h3 className="card-title text-base">
                <MapPin className="w-4 h-4" />
                Местоположение
              </h3>
              <input
                type="text"
                placeholder="Район или станция метро"
                className="input input-bordered"
                value={formData.location}
                onChange={(e) =>
                  setFormData({ ...formData, location: e.target.value })
                }
              />
              <div className="flex items-center gap-2 mt-2">
                <Shield className="w-4 h-4 text-success" />
                <span className="text-sm text-base-content/70">
                  Точный адрес виден только после договоренности
                </span>
              </div>
            </div>
          </div>

          {/* Optimal time to publish */}
          <div className="card bg-gradient-to-r from-warning/10 to-warning/5 border-2 border-warning/20">
            <div className="card-body">
              <h3 className="card-title text-base">
                <ClockIcon className="w-4 h-4" />
                Оптимальное время публикации
              </h3>
              <p className="text-sm">
                Сейчас{' '}
                <span className="font-bold">
                  {new Date().toLocaleTimeString('ru-RU', {
                    hour: '2-digit',
                    minute: '2-digit',
                  })}
                </span>
              </p>
              <p className="text-sm">
                Рекомендуем опубликовать в{' '}
                <span className="font-bold text-warning">19:00-21:00</span> для
                максимального охвата
              </p>
              <button className="btn btn-warning btn-sm mt-2">
                Запланировать на 19:00
              </button>
            </div>
          </div>

          {/* Quick Actions */}
          <div className="flex gap-3">
            <button
              onClick={() => setCurrentView('preview')}
              className="btn btn-primary flex-1"
              disabled={
                !formData.title ||
                !formData.price ||
                formData.images.length === 0
              }
            >
              Предпросмотр
              <ArrowRight className="w-4 h-4 ml-1" />
            </button>
            <button className="btn btn-ghost">Сохранить черновик</button>
          </div>

          {/* Tips */}
          <motion.div
            initial={{ opacity: 0 }}
            animate={{ opacity: 1 }}
            transition={{ delay: 0.5 }}
            className="alert shadow-sm"
          >
            <Info className="w-5 h-5" />
            <div>
              <h3 className="font-bold text-sm">Совет дня</h3>
              <p className="text-xs">
                Объявления с полным описанием продаются в 3 раза быстрее
              </p>
            </div>
          </motion.div>
        </div>
      </div>
    </motion.div>
  );

  const renderPreviewView = () => (
    <motion.div
      initial={{ opacity: 0 }}
      animate={{ opacity: 1 }}
      exit={{ opacity: 0 }}
      className="min-h-screen bg-base-200"
    >
      {/* Header */}
      <div className="navbar bg-base-100 border-b border-base-200">
        <div className="flex-1">
          <button
            onClick={() => setCurrentView('create')}
            className="btn btn-ghost gap-2"
          >
            <ChevronLeft className="w-5 h-5" />
            Редактировать
          </button>
        </div>
        <div className="flex-none">
          <div className="badge badge-success gap-1">
            <Check className="w-3 h-3" />
            Готово к публикации
          </div>
        </div>
      </div>

      {/* Preview Content */}
      <div className="container mx-auto px-4 py-8">
        <div className="max-w-4xl mx-auto">
          {/* Success Animation */}
          <motion.div
            initial={{ scale: 0 }}
            animate={{ scale: 1 }}
            transition={{ type: 'spring', stiffness: 200 }}
            className="text-center mb-8"
          >
            <div className="inline-flex items-center justify-center w-20 h-20 bg-success/20 rounded-full mb-4">
              <Check className="w-10 h-10 text-success" />
            </div>
            <h1 className="text-2xl font-bold mb-2">
              Отлично! Ваше объявление готово
            </h1>
            <p className="text-base-content/70">
              Вот как его увидят покупатели
            </p>
          </motion.div>

          {/* Listing Preview Card */}
          <motion.div
            initial={{ y: 20, opacity: 0 }}
            animate={{ y: 0, opacity: 1 }}
            transition={{ delay: 0.2 }}
            className="card bg-base-100 shadow-xl mb-6"
          >
            {/* Image Gallery */}
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
                {formData.images.length > 1 && (
                  <div className="absolute bottom-4 right-4 badge badge-neutral gap-1">
                    <ImageIcon className="w-3 h-3" />+
                    {formData.images.length - 1}
                  </div>
                )}
              </figure>
            )}

            <div className="card-body">
              <h2 className="card-title text-2xl">
                {formData.title || 'Название товара'}
              </h2>

              <div className="text-3xl font-bold text-primary mb-4">
                {formData.price ? `${formData.price} РСД` : 'Цена не указана'}
              </div>

              {formData.description && (
                <p className="text-base-content/80 mb-4 whitespace-pre-wrap">
                  {formData.description}
                </p>
              )}

              {/* Display attributes in preview */}
              {formData.category && categoryAttributes[formData.category] && (
                <div className="grid grid-cols-2 gap-4 mb-4">
                  {categoryAttributes[formData.category]
                    .filter((attr) => formData.attributes[attr.id])
                    .map((attr) => (
                      <div
                        key={attr.id}
                        className="flex justify-between py-2 border-b border-base-200"
                      >
                        <span className="text-sm text-base-content/60">
                          {attr.label}
                        </span>
                        <span className="text-sm font-medium">
                          {formData.attributes[attr.id]}
                        </span>
                      </div>
                    ))}
                </div>
              )}

              <div className="flex items-center gap-4 text-sm text-base-content/60 mb-4">
                <span className="flex items-center gap-1">
                  <MapPin className="w-4 h-4" />
                  {formData.location || 'Местоположение'}
                </span>
                <span className="flex items-center gap-1">
                  <Eye className="w-4 h-4" />0 просмотров
                </span>
                <span className="flex items-center gap-1">
                  <Heart className="w-4 h-4" />0 в избранном
                </span>
              </div>

              {/* Action Buttons */}
              <div className="flex gap-2">
                <button className="btn btn-primary flex-1">
                  <MessageCircle className="w-4 h-4 mr-1" />
                  Написать
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

          {/* Social sharing preview */}
          <motion.div
            initial={{ y: 20, opacity: 0 }}
            animate={{ y: 0, opacity: 1 }}
            transition={{ delay: 0.3 }}
            className="card bg-base-100 mb-6"
          >
            <div className="card-body">
              <h3 className="font-bold mb-4 flex items-center gap-2">
                <Share2 className="w-5 h-5" />
                Предпросмотр в соцсетях
              </h3>
              <div className="grid grid-cols-1 lg:grid-cols-3 gap-4">
                <div className="border rounded-lg p-4">
                  <p className="text-sm font-semibold mb-2">WhatsApp</p>
                  <div className="bg-green-50 rounded p-3">
                    <p className="font-medium text-sm">{formData.title}</p>
                    <p className="text-xs text-gray-600">
                      {formData.price} РСД
                    </p>
                  </div>
                </div>
                <div className="border rounded-lg p-4">
                  <p className="text-sm font-semibold mb-2">Telegram</p>
                  <div className="bg-blue-50 rounded p-3">
                    <p className="font-medium text-sm">{formData.title}</p>
                    <p className="text-xs text-gray-600">
                      {formData.price} РСД
                    </p>
                  </div>
                </div>
                <div className="border rounded-lg p-4">
                  <p className="text-sm font-semibold mb-2">Facebook</p>
                  <div className="bg-gray-50 rounded p-3">
                    <p className="font-medium text-sm">{formData.title}</p>
                    <p className="text-xs text-gray-600">
                      {formData.price} РСД
                    </p>
                  </div>
                </div>
              </div>
            </div>
          </motion.div>

          {/* Benefits Cards */}
          <div className="grid grid-cols-1 lg:grid-cols-3 gap-4 mb-8">
            <motion.div
              initial={{ y: 20, opacity: 0 }}
              animate={{ y: 0, opacity: 1 }}
              transition={{ delay: 0.4 }}
              className="card bg-primary/10 border-2 border-primary/20"
            >
              <div className="card-body text-center py-6">
                <TrendingUp className="w-8 h-8 text-primary mx-auto mb-2" />
                <h3 className="font-bold">Больше просмотров</h3>
                <p className="text-sm text-base-content/70">
                  AI-оптимизация увеличит охват
                </p>
              </div>
            </motion.div>

            <motion.div
              initial={{ y: 20, opacity: 0 }}
              animate={{ y: 0, opacity: 1 }}
              transition={{ delay: 0.5 }}
              className="card bg-success/10 border-2 border-success/20"
            >
              <div className="card-body text-center py-6">
                <Shield className="w-8 h-8 text-success mx-auto mb-2" />
                <h3 className="font-bold">Безопасная сделка</h3>
                <p className="text-sm text-base-content/70">
                  Мы защищаем ваши данные
                </p>
              </div>
            </motion.div>

            <motion.div
              initial={{ y: 20, opacity: 0 }}
              animate={{ y: 0, opacity: 1 }}
              transition={{ delay: 0.6 }}
              className="card bg-secondary/10 border-2 border-secondary/20"
            >
              <div className="card-body text-center py-6">
                <Award className="w-8 h-8 text-secondary mx-auto mb-2" />
                <h3 className="font-bold">Умное продвижение</h3>
                <p className="text-sm text-base-content/70">
                  Автоматическое продвижение в нужное время
                </p>
              </div>
            </motion.div>
          </div>

          {/* Publish Actions */}
          <motion.div
            initial={{ y: 20, opacity: 0 }}
            animate={{ y: 0, opacity: 1 }}
            transition={{ delay: 0.7 }}
            className="flex gap-3"
          >
            <button className="btn btn-primary btn-lg flex-1">
              Опубликовать сейчас
              <Sparkles className="w-5 h-5 ml-1" />
            </button>
            <button className="btn btn-outline btn-lg">
              <ClockIcon className="w-5 h-5 mr-1" />
              Запланировать
            </button>
          </motion.div>

          {/* Social proof */}
          <motion.div
            initial={{ opacity: 0 }}
            animate={{ opacity: 1 }}
            transition={{ delay: 0.8 }}
            className="text-center mt-8"
          >
            <div className="flex items-center justify-center gap-2 text-sm text-base-content/60">
              <Users className="w-4 h-4" />
              <span>
                <span className="font-semibold">2,345</span> продавцов уже
                воспользовались умными подсказками сегодня
              </span>
            </div>
          </motion.div>
        </div>
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
          <div className="badge badge-warning badge-lg">Улучшенная версия</div>
        </div>
      </div>

      {/* Main Content with Padding for Fixed Navbar */}
      <div className="pt-16">
        <AnimatePresence mode="wait">
          {currentView === 'start' && renderStartView()}
          {currentView === 'create' && renderCreateView()}
          {currentView === 'preview' && renderPreviewView()}
        </AnimatePresence>
      </div>
    </>
  );
}
