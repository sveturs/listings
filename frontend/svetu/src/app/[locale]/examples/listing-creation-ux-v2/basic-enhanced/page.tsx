'use client';

import React, { useState, useRef, useEffect } from 'react';
import Link from 'next/link';
import Image from 'next/image';
import { motion, AnimatePresence, Reorder } from 'framer-motion';
import {
  ChevronLeft,
  ChevronRight,
  Camera,
  MapPin,
  Package,
  CreditCard,
  Check,
  Upload,
  X,
  Info,
  TrendingUp,
  Clock,
  Shield,
  Sparkles,
  Save,
  Cloud,
  CloudOff,
  GripVertical,
  AlertCircle,
  Volume2,
  History,
  Eye,
} from 'lucide-react';

// Тип для элемента изображения с id для Reorder
interface ImageItem {
  id: string;
  url: string;
  file?: File;
}

export default function BasicEnhancedListingCreationPage() {
  const [currentStep, setCurrentStep] = useState(0);
  const [formData, setFormData] = useState({
    category: '',
    title: '',
    description: '',
    price: '',
    condition: 'used',
    images: [] as ImageItem[],
    location: '',
    privacyLevel: 'district',
    deliveryMethods: [] as string[],
    paymentMethods: [] as string[],
  });

  // Состояние автосохранения
  const [saveStatus, setSaveStatus] = useState<
    'saved' | 'saving' | 'unsaved' | 'error'
  >('saved');
  const [lastSaved, setLastSaved] = useState<Date | null>(null);

  // История изменений
  const [history, setHistory] = useState<typeof formData[]>([]);
  const [historyIndex, setHistoryIndex] = useState(-1);

  const fileInputRef = useRef<HTMLInputElement>(null);

  const steps = [
    { id: 'category', title: 'Категория', icon: Package },
    { id: 'info', title: 'Информация', icon: Info },
    { id: 'photos', title: 'Фотографии', icon: Camera },
    { id: 'location', title: 'Местоположение', icon: MapPin },
    { id: 'payment', title: 'Оплата и доставка', icon: CreditCard },
  ];

  const popularCategories = [
    { id: 'electronics', name: 'Электроника', icon: '📱', count: '12.5k' },
    { id: 'fashion', name: 'Одежда и обувь', icon: '👕', count: '8.3k' },
    { id: 'home', name: 'Дом и сад', icon: '🏠', count: '6.7k' },
    { id: 'vehicles', name: 'Транспорт', icon: '🚗', count: '4.2k' },
    { id: 'services', name: 'Услуги', icon: '🛠️', count: '3.9k' },
    { id: 'other', name: 'Другое', icon: '📦', count: '15.1k' },
  ];

  const conditions = [
    { id: 'new', label: 'Новый', description: 'Не использовался', icon: '✨' },
    {
      id: 'like-new',
      label: 'Как новый',
      description: 'Использовался бережно',
      icon: '⭐',
    },
    {
      id: 'used',
      label: 'Б/у',
      description: 'Есть следы использования',
      icon: '👍',
    },
    {
      id: 'for-parts',
      label: 'На запчасти',
      description: 'Требует ремонта',
      icon: '🔧',
    },
  ];

  const privacyLevels = [
    {
      id: 'exact',
      label: 'Точный адрес',
      description: 'Покупатели увидят точное местоположение',
    },
    { id: 'street', label: 'Только улица', description: 'Без номера дома' },
    {
      id: 'district',
      label: 'Только район',
      description: 'Безопасный вариант',
      recommended: true,
    },
    {
      id: 'city',
      label: 'Только город',
      description: 'Максимальная приватность',
    },
  ];

  // Автосохранение черновика
  useEffect(() => {
    const saveTimer = setTimeout(() => {
      if (saveStatus === 'unsaved') {
        setSaveStatus('saving');
        // Симуляция сохранения
        setTimeout(() => {
          setSaveStatus('saved');
          setLastSaved(new Date());
        }, 1000);
      }
    }, 2000);

    return () => clearTimeout(saveTimer);
  }, [formData, saveStatus]);

  // Отслеживание изменений
  useEffect(() => {
    if (saveStatus === 'saved') {
      setSaveStatus('unsaved');
    }
  }, [formData]);

  // Добавление в историю
  const addToHistory = () => {
    const newHistory = history.slice(0, historyIndex + 1);
    newHistory.push(formData);
    setHistory(newHistory);
    setHistoryIndex(newHistory.length - 1);
  };

  // Откат изменений
  const undo = () => {
    if (historyIndex > 0) {
      setHistoryIndex(historyIndex - 1);
      setFormData(history[historyIndex - 1]);
    }
  };

  const redo = () => {
    if (historyIndex < history.length - 1) {
      setHistoryIndex(historyIndex + 1);
      setFormData(history[historyIndex + 1]);
    }
  };

  const nextStep = () => {
    if (currentStep < steps.length - 1) {
      setCurrentStep(currentStep + 1);
      addToHistory();
    }
  };

  const prevStep = () => {
    if (currentStep > 0) {
      setCurrentStep(currentStep - 1);
    }
  };

  const handleImageUpload = (e: React.ChangeEvent<HTMLInputElement>) => {
    const files = e.target.files;
    if (files) {
      const newImages = Array.from(files).map((file) => ({
        id: `img-${Date.now()}-${Math.random()}`,
        url: URL.createObjectURL(file),
        file,
      }));

      // Проверка качества изображений
      newImages.forEach((img) => {
        const image = new window.Image();
        image.src = img.url;
        image.onload = () => {
          if (image.width < 800 || image.height < 600) {
            // Показать предупреждение о низком качестве
            console.log('Low quality image detected');
          }
        };
      });

      setFormData({
        ...formData,
        images: [...formData.images, ...newImages].slice(0, 8),
      });
    }
  };

  const removeImage = (id: string) => {
    setFormData({
      ...formData,
      images: formData.images.filter((img) => img.id !== id),
    });
  };

  // Progress bar с мотивацией
  const getMotivationalMessage = () => {
    const completedFields = [
      formData.category,
      formData.title,
      formData.description,
      formData.price,
      formData.images.length > 0,
      formData.location,
      formData.deliveryMethods.length > 0,
      formData.paymentMethods.length > 0,
    ].filter(Boolean).length;

    const messages = [
      'Отличное начало!',
      'Продолжайте в том же духе!',
      'Уже больше половины!',
      'Почти готово!',
      'Последний рывок!',
      'Превосходно! Готово к публикации!',
    ];

    return messages[Math.floor((completedFields / 8) * messages.length)];
  };

  const renderStep = () => {
    switch (currentStep) {
      case 0: // Category Selection
        return (
          <motion.div
            initial={{ opacity: 0, y: 20 }}
            animate={{ opacity: 1, y: 0 }}
            exit={{ opacity: 0, y: -20 }}
            className="space-y-6"
          >
            <div>
              <h2 className="text-2xl font-bold mb-2">Выберите категорию</h2>
              <p className="text-base-content/70">
                Это поможет покупателям быстрее найти ваше объявление
              </p>
            </div>

            <div className="form-control">
              <input
                type="text"
                placeholder="🔍 Поиск категории..."
                className="input input-bordered input-lg"
              />
            </div>

            <div>
              <h3 className="font-semibold mb-4 flex items-center gap-2">
                <TrendingUp className="w-5 h-5" />
                Популярные категории
              </h3>
              <div className="grid grid-cols-2 lg:grid-cols-3 gap-4">
                {popularCategories.map((cat) => (
                  <motion.button
                    key={cat.id}
                    whileHover={{ scale: 1.02 }}
                    whileTap={{ scale: 0.98 }}
                    onClick={() =>
                      setFormData({ ...formData, category: cat.id })
                    }
                    className={`card ${
                      formData.category === cat.id ? 'ring-2 ring-primary' : ''
                    } hover:shadow-lg transition-all cursor-pointer`}
                  >
                    <div className="card-body p-4">
                      <div className="text-3xl mb-2">{cat.icon}</div>
                      <h4 className="font-semibold">{cat.name}</h4>
                      <p className="text-sm text-base-content/60">
                        {cat.count} объявлений
                      </p>
                    </div>
                  </motion.button>
                ))}
              </div>
            </div>

            <div className="divider">или</div>

            <button className="btn btn-outline btn-block">
              Показать все категории
            </button>
          </motion.div>
        );

      case 1: // Basic Information
        return (
          <motion.div
            initial={{ opacity: 0, y: 20 }}
            animate={{ opacity: 1, y: 0 }}
            exit={{ opacity: 0, y: -20 }}
            className="space-y-6"
          >
            <div>
              <h2 className="text-2xl font-bold mb-2">Основная информация</h2>
              <p className="text-base-content/70">
                Заполните ключевые данные о вашем товаре
              </p>
            </div>

            <div className="form-control">
              <label className="label">
                <span className="label-text font-semibold">
                  Название объявления
                </span>
                <span className="label-text-alt">
                  {formData.title.length}/80
                </span>
              </label>
              <input
                type="text"
                placeholder="Например: iPhone 13 Pro, 256GB, синий"
                className="input input-bordered input-lg"
                value={formData.title}
                onChange={(e) =>
                  setFormData({ ...formData, title: e.target.value })
                }
                maxLength={80}
              />
              <label className="label">
                <span className="label-text-alt text-info">
                  💡 Укажите бренд, модель и ключевые характеристики
                </span>
              </label>
            </div>

            <div className="form-control">
              <label className="label">
                <span className="label-text font-semibold">Описание</span>
                <span className="label-text-alt">
                  {formData.description.length}/1000
                </span>
              </label>
              <div className="relative">
                <textarea
                  className="textarea textarea-bordered h-32 w-full"
                  placeholder="Опишите состояние товара, комплектацию, причину продажи..."
                  value={formData.description}
                  onChange={(e) =>
                    setFormData({ ...formData, description: e.target.value })
                  }
                  maxLength={1000}
                />
                <button className="absolute bottom-2 right-2 btn btn-xs btn-ghost gap-1">
                  <Volume2 className="w-3 h-3" />
                  Голосовой ввод
                </button>
              </div>

              {/* Шаблоны описаний */}
              {formData.category && (
                <div className="mt-2">
                  <button className="btn btn-outline btn-sm gap-1">
                    <Sparkles className="w-3 h-3" />
                    Использовать шаблон для {formData.category}
                  </button>
                </div>
              )}
            </div>

            <div className="grid grid-cols-1 lg:grid-cols-2 gap-6">
              <div className="form-control">
                <label className="label">
                  <span className="label-text font-semibold">Цена</span>
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
                <label className="label">
                  <span className="label-text-alt text-success">
                    📊 Средняя цена в категории: 45.000 РСД
                  </span>
                  <button className="label-text-alt link link-primary">
                    Сравнить с похожими
                  </button>
                </label>
              </div>

              <div className="form-control">
                <label className="label">
                  <span className="label-text font-semibold">Состояние</span>
                </label>
                <div className="grid grid-cols-2 gap-2">
                  {conditions.map((cond) => (
                    <button
                      key={cond.id}
                      onClick={() =>
                        setFormData({ ...formData, condition: cond.id })
                      }
                      className={`btn ${
                        formData.condition === cond.id
                          ? 'btn-primary'
                          : 'btn-outline'
                      } btn-sm justify-start`}
                    >
                      <span className="text-lg mr-2">{cond.icon}</span>
                      <span>{cond.label}</span>
                    </button>
                  ))}
                </div>
              </div>
            </div>
          </motion.div>
        );

      case 2: // Photos with drag & drop
        return (
          <motion.div
            initial={{ opacity: 0, y: 20 }}
            animate={{ opacity: 1, y: 0 }}
            exit={{ opacity: 0, y: -20 }}
            className="space-y-6"
          >
            <div>
              <h2 className="text-2xl font-bold mb-2">Фотографии</h2>
              <p className="text-base-content/70">
                Объявления с фото получают в 5 раз больше просмотров
              </p>
            </div>

            <Reorder.Group
              axis="y"
              values={formData.images}
              onReorder={(newImages) =>
                setFormData({ ...formData, images: newImages })
              }
              className="grid grid-cols-2 lg:grid-cols-4 gap-4"
            >
              {formData.images.map((img, index) => (
                <Reorder.Item
                  key={img.id}
                  value={img}
                  className="relative aspect-square cursor-move"
                >
                  <motion.div
                    initial={{ scale: 0 }}
                    animate={{ scale: 1 }}
                    className="relative h-full w-full"
                  >
                    <Image
                      src={img.url}
                      alt={`Photo ${index + 1}`}
                      fill
                      className="object-cover rounded-lg"
                    />
                    {index === 0 && (
                      <div className="absolute top-2 left-2 badge badge-primary">
                        Главное фото
                      </div>
                    )}
                    <div className="absolute top-2 right-2 flex gap-1">
                      <div className="btn btn-circle btn-sm btn-neutral">
                        <GripVertical className="w-3 h-3" />
                      </div>
                      <button
                        onClick={() => removeImage(img.id)}
                        className="btn btn-circle btn-sm btn-error"
                      >
                        <X className="w-4 h-4" />
                      </button>
                    </div>
                    {/* Предупреждение о качестве */}
                    {Math.random() > 0.7 && (
                      <div className="absolute bottom-2 left-2 right-2 alert alert-warning p-2">
                        <AlertCircle className="w-3 h-3" />
                        <span className="text-xs">Низкое качество фото</span>
                      </div>
                    )}
                  </motion.div>
                </Reorder.Item>
              ))}

              {formData.images.length < 8 && (
                <label className="border-2 border-dashed border-base-300 rounded-lg aspect-square flex flex-col items-center justify-center cursor-pointer hover:border-primary transition-colors">
                  <Upload className="w-8 h-8 mb-2 text-base-content/50" />
                  <span className="text-sm text-base-content/70">
                    Добавить фото
                  </span>
                  <span className="text-xs text-base-content/50">
                    {8 - formData.images.length} осталось
                  </span>
                  <input
                    ref={fileInputRef}
                    type="file"
                    multiple
                    accept="image/*"
                    className="hidden"
                    onChange={handleImageUpload}
                  />
                </label>
              )}
            </Reorder.Group>

            <div className="alert alert-info">
              <Info className="w-5 h-5" />
              <div>
                <h3 className="font-bold">Советы для хороших фото</h3>
                <ul className="text-sm mt-1 space-y-1">
                  <li>• Перетащите фото для изменения порядка</li>
                  <li>• Первое фото - самое важное</li>
                  <li>• Снимайте при дневном свете</li>
                  <li>• Покажите товар с разных сторон</li>
                  <li>• Включите все дефекты, если есть</li>
                </ul>
              </div>
            </div>

            {/* Рекомендации по недостающим фото */}
            {formData.images.length > 0 && formData.images.length < 4 && (
              <div className="alert">
                <Sparkles className="w-5 h-5" />
                <div>
                  <h3 className="font-bold">Рекомендуем добавить</h3>
                  <p className="text-sm">
                    Фото сзади, фото деталей, фото в использовании
                  </p>
                </div>
              </div>
            )}
          </motion.div>
        );

      case 3: // Location
        return (
          <motion.div
            initial={{ opacity: 0, y: 20 }}
            animate={{ opacity: 1, y: 0 }}
            exit={{ opacity: 0, y: -20 }}
            className="space-y-6"
          >
            <div>
              <h2 className="text-2xl font-bold mb-2">Местоположение</h2>
              <p className="text-base-content/70">
                Укажите, где находится товар
              </p>
            </div>

            <div className="form-control">
              <label className="label">
                <span className="label-text font-semibold">
                  Адрес или район
                </span>
              </label>
              <input
                type="text"
                placeholder="Начните вводить адрес..."
                className="input input-bordered input-lg"
                value={formData.location}
                onChange={(e) =>
                  setFormData({ ...formData, location: e.target.value })
                }
              />
            </div>

            <div className="bg-base-200 h-64 rounded-lg flex items-center justify-center">
              <MapPin className="w-12 h-12 text-base-content/30" />
            </div>

            <div>
              <h3 className="font-semibold mb-4 flex items-center gap-2">
                <Shield className="w-5 h-5" />
                Настройки приватности
              </h3>
              <div className="space-y-3">
                {privacyLevels.map((level) => (
                  <label
                    key={level.id}
                    className={`card cursor-pointer ${
                      formData.privacyLevel === level.id
                        ? 'ring-2 ring-primary'
                        : ''
                    }`}
                  >
                    <div className="card-body p-4 flex-row items-start">
                      <input
                        type="radio"
                        name="privacy"
                        className="radio radio-primary"
                        checked={formData.privacyLevel === level.id}
                        onChange={() =>
                          setFormData({ ...formData, privacyLevel: level.id })
                        }
                      />
                      <div className="flex-1 ml-4">
                        <div className="flex items-center gap-2">
                          <span className="font-semibold">{level.label}</span>
                          {level.recommended && (
                            <span className="badge badge-success badge-sm">
                              Рекомендуется
                            </span>
                          )}
                        </div>
                        <p className="text-sm text-base-content/70">
                          {level.description}
                        </p>
                      </div>
                    </div>
                  </label>
                ))}
              </div>
            </div>
          </motion.div>
        );

      case 4: // Payment & Delivery
        return (
          <motion.div
            initial={{ opacity: 0, y: 20 }}
            animate={{ opacity: 1, y: 0 }}
            exit={{ opacity: 0, y: -20 }}
            className="space-y-6"
          >
            <div>
              <h2 className="text-2xl font-bold mb-2">Оплата и доставка</h2>
              <p className="text-base-content/70">
                Как покупатель сможет получить и оплатить товар
              </p>
            </div>

            <div>
              <h3 className="font-semibold mb-4">Способы получения</h3>
              <div className="space-y-3">
                {[
                  {
                    id: 'pickup',
                    label: 'Личная встреча',
                    icon: '🤝',
                    popular: true,
                  },
                  { id: 'delivery', label: 'Доставка курьером', icon: '🚚' },
                  { id: 'post', label: 'Отправка почтой', icon: '📦' },
                ].map((method) => (
                  <label key={method.id} className="flex items-center gap-3">
                    <input
                      type="checkbox"
                      className="checkbox checkbox-primary"
                      checked={formData.deliveryMethods.includes(method.id)}
                      onChange={(e) => {
                        if (e.target.checked) {
                          setFormData({
                            ...formData,
                            deliveryMethods: [
                              ...formData.deliveryMethods,
                              method.id,
                            ],
                          });
                        } else {
                          setFormData({
                            ...formData,
                            deliveryMethods: formData.deliveryMethods.filter(
                              (m) => m !== method.id
                            ),
                          });
                        }
                      }}
                    />
                    <span className="text-2xl">{method.icon}</span>
                    <span className="flex-1">{method.label}</span>
                    {method.popular && (
                      <span className="badge badge-primary badge-sm">
                        Популярно
                      </span>
                    )}
                  </label>
                ))}
              </div>
            </div>

            <div>
              <h3 className="font-semibold mb-4">Способы оплаты</h3>
              <div className="space-y-3">
                {[
                  { id: 'cash', label: 'Наличные', icon: '💵', safe: true },
                  { id: 'card', label: 'Перевод на карту', icon: '💳' },
                  {
                    id: 'online',
                    label: 'Онлайн оплата',
                    icon: '📱',
                    new: true,
                  },
                ].map((method) => (
                  <label key={method.id} className="flex items-center gap-3">
                    <input
                      type="checkbox"
                      className="checkbox checkbox-primary"
                      checked={formData.paymentMethods.includes(method.id)}
                      onChange={(e) => {
                        if (e.target.checked) {
                          setFormData({
                            ...formData,
                            paymentMethods: [
                              ...formData.paymentMethods,
                              method.id,
                            ],
                          });
                        } else {
                          setFormData({
                            ...formData,
                            paymentMethods: formData.paymentMethods.filter(
                              (m) => m !== method.id
                            ),
                          });
                        }
                      }}
                    />
                    <span className="text-2xl">{method.icon}</span>
                    <span className="flex-1">{method.label}</span>
                    {method.safe && (
                      <span className="badge badge-success badge-sm">
                        Безопасно
                      </span>
                    )}
                    {method.new && (
                      <span className="badge badge-info badge-sm">Новое</span>
                    )}
                  </label>
                ))}
              </div>
            </div>

            <div className="alert alert-success">
              <Sparkles className="w-5 h-5" />
              <div>
                <h3 className="font-bold">Готово к публикации!</h3>
                <p className="text-sm">
                  Ваше объявление готово. Нажмите "Опубликовать" для размещения.
                </p>
              </div>
            </div>

            {/* Предпросмотр шаринга */}
            <div className="card bg-base-200">
              <div className="card-body">
                <h3 className="card-title text-base">
                  <Eye className="w-4 h-4" />
                  Как будет выглядеть в соцсетях
                </h3>
                <div className="flex gap-2">
                  <div className="btn btn-sm btn-ghost">WhatsApp</div>
                  <div className="btn btn-sm btn-ghost">Telegram</div>
                  <div className="btn btn-sm btn-ghost">Facebook</div>
                </div>
              </div>
            </div>
          </motion.div>
        );

      default:
        return null;
    }
  };

  return (
    <div className="min-h-screen bg-base-100">
      {/* Header with save status */}
      <div className="navbar bg-base-100 border-b border-base-200">
        <div className="flex-1">
          <Link
            href="/ru/examples/listing-creation-ux-v2"
            className="btn btn-ghost"
          >
            <ChevronLeft className="w-5 h-5" />
            Назад к примерам
          </Link>
        </div>
        <div className="flex-none gap-2">
          {/* История изменений */}
          <div className="flex gap-1">
            <button
              onClick={undo}
              disabled={historyIndex <= 0}
              className="btn btn-ghost btn-sm"
              title="Отменить"
            >
              <History className="w-4 h-4 rotate-180" />
            </button>
            <button
              onClick={redo}
              disabled={historyIndex >= history.length - 1}
              className="btn btn-ghost btn-sm"
              title="Повторить"
            >
              <History className="w-4 h-4" />
            </button>
          </div>

          {/* Статус сохранения */}
          <div className="flex items-center gap-2">
            {saveStatus === 'saved' && (
              <>
                <Cloud className="w-4 h-4 text-success" />
                <span className="text-sm text-success">Сохранено</span>
              </>
            )}
            {saveStatus === 'saving' && (
              <>
                <Cloud className="w-4 h-4 text-warning animate-pulse" />
                <span className="text-sm text-warning">Сохраняется...</span>
              </>
            )}
            {saveStatus === 'unsaved' && (
              <>
                <CloudOff className="w-4 h-4 text-base-content/50" />
                <span className="text-sm text-base-content/50">
                  Изменения не сохранены
                </span>
              </>
            )}
            {saveStatus === 'error' && (
              <>
                <AlertCircle className="w-4 h-4 text-error" />
                <span className="text-sm text-error">Ошибка сохранения</span>
              </>
            )}
          </div>

          {lastSaved && (
            <span className="text-xs text-base-content/50">
              Последнее сохранение:{' '}
              {new Date(lastSaved).toLocaleTimeString('ru-RU', {
                hour: '2-digit',
                minute: '2-digit',
              })}
            </span>
          )}

          <div className="badge badge-primary badge-lg">
            Улучшенная версия
          </div>
        </div>
      </div>

      {/* Progress Bar with motivation */}
      <div className="bg-base-200 py-4">
        <div className="container mx-auto px-4">
          <div className="flex items-center justify-between mb-4">
            <h1 className="text-lg font-semibold">Создание объявления</h1>
            <div className="text-sm text-base-content/70">
              Шаг {currentStep + 1} из {steps.length}
            </div>
          </div>

          {/* Motivational message */}
          <div className="text-center mb-2">
            <p className="text-sm font-medium text-primary">
              {getMotivationalMessage()}
            </p>
          </div>

          {/* Desktop Progress */}
          <div className="hidden lg:flex items-center gap-2">
            {steps.map((step, index) => {
              const Icon = step.icon;
              const isActive = index === currentStep;
              const isCompleted = index < currentStep;

              return (
                <React.Fragment key={step.id}>
                  <div
                    className={`flex items-center gap-2 px-4 py-2 rounded-lg transition-all ${
                      isActive
                        ? 'bg-primary text-primary-content'
                        : isCompleted
                          ? 'bg-success text-success-content'
                          : 'bg-base-300 text-base-content/50'
                    }`}
                  >
                    {isCompleted ? (
                      <Check className="w-5 h-5" />
                    ) : (
                      <Icon className="w-5 h-5" />
                    )}
                    <span className="font-medium">{step.title}</span>
                  </div>
                  {index < steps.length - 1 && (
                    <div
                      className={`flex-1 h-1 ${
                        index < currentStep ? 'bg-success' : 'bg-base-300'
                      }`}
                    />
                  )}
                </React.Fragment>
              );
            })}
          </div>

          {/* Mobile Progress */}
          <div className="lg:hidden">
            <div className="flex items-center justify-between mb-2">
              {steps.map((step, index) => {
                const Icon = step.icon;
                const isActive = index === currentStep;
                const isCompleted = index < currentStep;

                return (
                  <div
                    key={step.id}
                    className={`w-10 h-10 rounded-full flex items-center justify-center transition-all ${
                      isActive
                        ? 'bg-primary text-primary-content'
                        : isCompleted
                          ? 'bg-success text-success-content'
                          : 'bg-base-300 text-base-content/50'
                    }`}
                  >
                    {isCompleted ? (
                      <Check className="w-5 h-5" />
                    ) : (
                      <Icon className="w-5 h-5" />
                    )}
                  </div>
                );
              })}
            </div>
            <div className="text-center text-sm font-medium">
              {steps[currentStep].title}
            </div>
          </div>
        </div>
      </div>

      {/* Content */}
      <div className="container mx-auto px-4 py-8">
        <div className="max-w-2xl mx-auto">
          <AnimatePresence mode="wait">{renderStep()}</AnimatePresence>
        </div>
      </div>

      {/* Footer Navigation */}
      <div className="fixed bottom-0 left-0 right-0 bg-base-100 border-t border-base-200 p-4">
        <div className="container mx-auto px-4">
          <div className="flex items-center justify-between max-w-2xl mx-auto">
            <button
              onClick={prevStep}
              disabled={currentStep === 0}
              className="btn btn-outline gap-2"
            >
              <ChevronLeft className="w-5 h-5" />
              Назад
            </button>

            <div className="flex items-center gap-2">
              <button className="btn btn-ghost btn-sm">
                <Save className="w-4 h-4 mr-1" />
                Сохранить черновик
              </button>
            </div>

            {currentStep === steps.length - 1 ? (
              <button className="btn btn-primary gap-2">
                Опубликовать
                <Check className="w-5 h-5" />
              </button>
            ) : (
              <button onClick={nextStep} className="btn btn-primary gap-2">
                Далее
                <ChevronRight className="w-5 h-5" />
              </button>
            )}
          </div>
        </div>
      </div>
    </div>
  );
}