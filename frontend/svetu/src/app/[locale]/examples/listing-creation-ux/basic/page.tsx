'use client';

import React, { useState } from 'react';
import Link from 'next/link';
import Image from 'next/image';
import { motion, AnimatePresence } from 'framer-motion';
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
} from 'lucide-react';

export default function BasicListingCreationPage() {
  const [currentStep, setCurrentStep] = useState(0);
  const [formData, setFormData] = useState({
    category: '',
    title: '',
    description: '',
    price: '',
    condition: 'used',
    images: [] as string[],
    location: '',
    privacyLevel: 'district',
    deliveryMethods: [] as string[],
    paymentMethods: [] as string[],
  });

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

  const nextStep = () => {
    if (currentStep < steps.length - 1) {
      setCurrentStep(currentStep + 1);
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
      const newImages = Array.from(files).map((file) =>
        URL.createObjectURL(file)
      );
      setFormData({
        ...formData,
        images: [...formData.images, ...newImages].slice(0, 8),
      });
    }
  };

  const removeImage = (index: number) => {
    setFormData({
      ...formData,
      images: formData.images.filter((_, i) => i !== index),
    });
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
              <textarea
                className="textarea textarea-bordered h-32"
                placeholder="Опишите состояние товара, комплектацию, причину продажи..."
                value={formData.description}
                onChange={(e) =>
                  setFormData({ ...formData, description: e.target.value })
                }
                maxLength={1000}
              />
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

      case 2: // Photos
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

            <div className="grid grid-cols-2 lg:grid-cols-4 gap-4">
              {formData.images.map((img, index) => (
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
                  {index === 0 && (
                    <div className="absolute top-2 left-2 badge badge-primary">
                      Главное фото
                    </div>
                  )}
                  <button
                    onClick={() => removeImage(index)}
                    className="absolute top-2 right-2 btn btn-circle btn-sm btn-error"
                  >
                    <X className="w-4 h-4" />
                  </button>
                </motion.div>
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
                    type="file"
                    multiple
                    accept="image/*"
                    className="hidden"
                    onChange={handleImageUpload}
                  />
                </label>
              )}
            </div>

            <div className="alert alert-info">
              <Info className="w-5 h-5" />
              <div>
                <h3 className="font-bold">Советы для хороших фото</h3>
                <ul className="text-sm mt-1 space-y-1">
                  <li>• Снимайте при дневном свете</li>
                  <li>• Покажите товар с разных сторон</li>
                  <li>• Включите все дефекты, если есть</li>
                  <li>• Первое фото - самое важное</li>
                </ul>
              </div>
            </div>
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
                  Ваше объявление готово. Нажмите &quot;Опубликовать&quot; для
                  размещения.
                </p>
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
      {/* Header */}
      <div className="navbar bg-base-100 border-b border-base-200">
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
          <div className="badge badge-primary badge-lg">Без AI</div>
        </div>
      </div>

      {/* Progress Bar */}
      <div className="bg-base-200 py-4">
        <div className="container mx-auto px-4">
          <div className="flex items-center justify-between mb-4">
            <h1 className="text-lg font-semibold">Создание объявления</h1>
            <div className="text-sm text-base-content/70">
              Шаг {currentStep + 1} из {steps.length}
            </div>
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
                <Clock className="w-4 h-4 mr-1" />
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
