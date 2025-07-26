"use client";

import React, { useState, useRef, useEffect } from "react";
import { useTranslations } from "next-intl";
import Link from "next/link";
import { motion, AnimatePresence, useAnimation } from "framer-motion";
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
  CreditCard,
  Smartphone,
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
} from "lucide-react";

export default function NoBackendListingCreationPage() {
  const t = useTranslations();
  const [currentView, setCurrentView] = useState<"start" | "create" | "preview">("start");
  const [quickMode, setQuickMode] = useState(false);
  const [formData, setFormData] = useState({
    images: [] as string[],
    category: "",
    title: "",
    price: "",
    description: "",
    location: "",
    deliveryMethods: ["pickup"],
  });
  const [suggestions, setSuggestions] = useState({
    title: "",
    category: "",
    price: "",
  });
  const fileInputRef = useRef<HTMLInputElement>(null);
  const controls = useAnimation();

  // Simulated quick templates
  const quickTemplates = [
    {
      id: "phone",
      icon: "📱",
      title: "Продаю телефон",
      fields: ["Модель", "Память", "Состояние"],
    },
    {
      id: "clothes",
      icon: "👕",
      title: "Одежда/Обувь",
      fields: ["Размер", "Бренд", "Состояние"],
    },
    {
      id: "electronics",
      icon: "💻",
      title: "Электроника",
      fields: ["Бренд", "Модель", "Год"],
    },
    {
      id: "furniture",
      icon: "🛋️",
      title: "Мебель",
      fields: ["Тип", "Размеры", "Материал"],
    },
  ];

  const popularCategories = [
    { id: "electronics", name: "Электроника", icon: "📱", gradient: "from-blue-500 to-purple-500" },
    { id: "fashion", name: "Мода", icon: "👗", gradient: "from-pink-500 to-rose-500" },
    { id: "home", name: "Дом", icon: "🏠", gradient: "from-green-500 to-emerald-500" },
    { id: "auto", name: "Авто", icon: "🚗", gradient: "from-orange-500 to-red-500" },
  ];

  useEffect(() => {
    // Simulate AI suggestions when image is uploaded
    if (formData.images.length > 0 && !suggestions.title) {
      setTimeout(() => {
        setSuggestions({
          title: "iPhone 13 Pro, 256GB, Pacific Blue",
          category: "electronics",
          price: "65000",
        });
      }, 1000);
    }
  }, [formData.images]);

  const handleImageUpload = (e: React.ChangeEvent<HTMLInputElement>) => {
    const files = e.target.files;
    if (files) {
      const newImages = Array.from(files).map((file) => URL.createObjectURL(file));
      setFormData({ ...formData, images: [...formData.images, ...newImages].slice(0, 8) });
      if (newImages.length > 0) {
        setCurrentView("create");
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
    });
    controls.start({
      scale: [1, 1.05, 1],
      transition: { duration: 0.3 },
    });
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
            Продайте быстрее, чем заварится кофе ☕
          </h1>
          <p className="text-xl text-base-content/70 mb-8">
            Новый опыт создания объявлений — проще, быстрее, умнее
          </p>
          
          {/* Stats */}
          <div className="flex justify-center gap-8 mb-8">
            <motion.div
              initial={{ scale: 0 }}
              animate={{ scale: 1 }}
              transition={{ delay: 0.2 }}
              className="text-center"
            >
              <div className="text-3xl font-bold text-primary">3 мин</div>
              <div className="text-sm text-base-content/60">в среднем</div>
            </motion.div>
            <motion.div
              initial={{ scale: 0 }}
              animate={{ scale: 1 }}
              transition={{ delay: 0.3 }}
              className="text-center"
            >
              <div className="text-3xl font-bold text-success">95%</div>
              <div className="text-sm text-base-content/60">завершают</div>
            </motion.div>
            <motion.div
              initial={{ scale: 0 }}
              animate={{ scale: 1 }}
              transition={{ delay: 0.4 }}
              className="text-center"
            >
              <div className="text-3xl font-bold text-secondary">5x</div>
              <div className="text-sm text-base-content/60">больше просмотров</div>
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
                <div className="badge badge-lg badge-warning gap-2">
                  <Zap className="w-4 h-4" />
                  Быстрый старт
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

          {/* Alternative Options */}
          <div className="grid grid-cols-1 lg:grid-cols-2 gap-4 mb-8">
            <motion.div
              initial={{ x: -20, opacity: 0 }}
              animate={{ x: 0, opacity: 1 }}
              transition={{ delay: 0.3 }}
            >
              <button
                onClick={() => {
                  setCurrentView("create");
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
              transition={{ delay: 0.4 }}
            >
              <button
                onClick={() => {
                  setQuickMode(true);
                  setCurrentView("create");
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
            transition={{ delay: 0.5 }}
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
                    setCurrentView("create");
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
              onClick={() => setCurrentView("start")}
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
                  {suggestions.title} • Рекомендуемая цена: {suggestions.price} РСД
                </p>
              </div>
              <button onClick={applySuggestions} className="btn btn-sm btn-primary">
                Применить
              </button>
            </motion.div>
          )}

          {/* Photo Upload Section */}
          <motion.div
            animate={controls}
            className="card bg-base-200"
          >
            <div className="card-body">
              <h2 className="card-title">
                <Camera className="w-5 h-5" />
                Фотографии
                {formData.images.length > 0 && (
                  <span className="badge badge-primary">{formData.images.length}/8</span>
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
                    <img
                      src={img}
                      alt={`Photo ${index + 1}`}
                      className="w-full h-full object-cover rounded-lg"
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
                    <span className="text-xs text-base-content/50 mt-1">Добавить</span>
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
            </div>
          </motion.div>

          {/* Quick Info Section */}
          <div className="card bg-base-200">
            <div className="card-body space-y-4">
              {/* Title */}
              <div className="form-control">
                <label className="label">
                  <span className="label-text font-semibold">Название</span>
                  <span className="label-text-alt">{formData.title.length}/80</span>
                </label>
                <input
                  type="text"
                  placeholder="Что вы продаете?"
                  className="input input-bordered"
                  value={formData.title}
                  onChange={(e) => setFormData({ ...formData, title: e.target.value })}
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
                        onClick={() => setFormData({ ...formData, category: cat.id })}
                        className={`btn btn-sm ${
                          formData.category === cat.id ? "btn-primary" : "btn-outline"
                        } gap-1`}
                      >
                        <span>{cat.icon}</span>
                        {cat.name}
                      </button>
                    ))}
                  </div>
                </div>
              )}

              {/* Price */}
              <div className="form-control">
                <label className="label">
                  <span className="label-text font-semibold">Цена</span>
                  <span className="label-text-alt text-success">
                    📊 Средняя: 45.000 РСД
                  </span>
                </label>
                <label className="input-group">
                  <input
                    type="number"
                    placeholder="0"
                    className="input input-bordered flex-1"
                    value={formData.price}
                    onChange={(e) => setFormData({ ...formData, price: e.target.value })}
                  />
                  <span>РСД</span>
                </label>
              </div>

              {/* Quick Description */}
              {!quickMode && (
                <div className="form-control">
                  <label className="label">
                    <span className="label-text font-semibold">Описание</span>
                    <span className="label-text-alt">Опционально</span>
                  </label>
                  <textarea
                    className="textarea textarea-bordered h-20"
                    placeholder="Добавьте детали..."
                    value={formData.description}
                    onChange={(e) => setFormData({ ...formData, description: e.target.value })}
                  />
                </div>
              )}
            </div>
          </div>

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
                onChange={(e) => setFormData({ ...formData, location: e.target.value })}
              />
              <div className="flex items-center gap-2 mt-2">
                <Shield className="w-4 h-4 text-success" />
                <span className="text-sm text-base-content/70">
                  Точный адрес виден только после договоренности
                </span>
              </div>
            </div>
          </div>

          {/* Quick Actions */}
          <div className="flex gap-3">
            <button
              onClick={() => setCurrentView("preview")}
              className="btn btn-primary flex-1"
              disabled={!formData.title || !formData.price || formData.images.length === 0}
            >
              Предпросмотр
              <ArrowRight className="w-4 h-4 ml-1" />
            </button>
            <button className="btn btn-ghost">
              Сохранить черновик
            </button>
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
                Объявления с 3+ фото продаются в 2 раза быстрее
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
            onClick={() => setCurrentView("create")}
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
            transition={{ type: "spring", stiffness: 200 }}
            className="text-center mb-8"
          >
            <div className="inline-flex items-center justify-center w-20 h-20 bg-success/20 rounded-full mb-4">
              <Check className="w-10 h-10 text-success" />
            </div>
            <h1 className="text-2xl font-bold mb-2">Отлично! Ваше объявление готово</h1>
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
                <img
                  src={formData.images[0]}
                  alt={formData.title}
                  className="w-full h-96 object-cover"
                />
                {formData.images.length > 1 && (
                  <div className="absolute bottom-4 right-4 badge badge-neutral gap-1">
                    <ImageIcon className="w-3 h-3" />
                    +{formData.images.length - 1}
                  </div>
                )}
              </figure>
            )}

            <div className="card-body">
              <h2 className="card-title text-2xl">{formData.title || "Название товара"}</h2>
              
              <div className="text-3xl font-bold text-primary mb-4">
                {formData.price ? `${formData.price} РСД` : "Цена не указана"}
              </div>

              {formData.description && (
                <p className="text-base-content/80 mb-4">{formData.description}</p>
              )}

              <div className="flex items-center gap-4 text-sm text-base-content/60 mb-4">
                <span className="flex items-center gap-1">
                  <MapPin className="w-4 h-4" />
                  {formData.location || "Местоположение"}
                </span>
                <span className="flex items-center gap-1">
                  <Eye className="w-4 h-4" />
                  0 просмотров
                </span>
                <span className="flex items-center gap-1">
                  <Heart className="w-4 h-4" />
                  0 в избранном
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

          {/* Benefits Cards */}
          <div className="grid grid-cols-1 lg:grid-cols-3 gap-4 mb-8">
            <motion.div
              initial={{ y: 20, opacity: 0 }}
              animate={{ y: 0, opacity: 1 }}
              transition={{ delay: 0.3 }}
              className="card bg-primary/10 border-2 border-primary/20"
            >
              <div className="card-body text-center py-6">
                <TrendingUp className="w-8 h-8 text-primary mx-auto mb-2" />
                <h3 className="font-bold">Больше просмотров</h3>
                <p className="text-sm text-base-content/70">
                  Ваше объявление увидят тысячи покупателей
                </p>
              </div>
            </motion.div>

            <motion.div
              initial={{ y: 20, opacity: 0 }}
              animate={{ y: 0, opacity: 1 }}
              transition={{ delay: 0.4 }}
              className="card bg-success/10 border-2 border-success/20"
            >
              <div className="card-body text-center py-6">
                <Shield className="w-8 h-8 text-success mx-auto mb-2" />
                <h3 className="font-bold">Безопасная сделка</h3>
                <p className="text-sm text-base-content/70">
                  Мы защищаем ваши данные и помогаем с оплатой
                </p>
              </div>
            </motion.div>

            <motion.div
              initial={{ y: 20, opacity: 0 }}
              animate={{ y: 0, opacity: 1 }}
              transition={{ delay: 0.5 }}
              className="card bg-secondary/10 border-2 border-secondary/20"
            >
              <div className="card-body text-center py-6">
                <Award className="w-8 h-8 text-secondary mx-auto mb-2" />
                <h3 className="font-bold">Премиум размещение</h3>
                <p className="text-sm text-base-content/70">
                  Поднимите объявление в топ за 99 РСД
                </p>
              </div>
            </motion.div>
          </div>

          {/* Publish Actions */}
          <motion.div
            initial={{ y: 20, opacity: 0 }}
            animate={{ y: 0, opacity: 1 }}
            transition={{ delay: 0.6 }}
            className="flex gap-3"
          >
            <button className="btn btn-primary btn-lg flex-1">
              Опубликовать бесплатно
              <Sparkles className="w-5 h-5 ml-1" />
            </button>
            <button className="btn btn-outline btn-lg">
              Сохранить черновик
            </button>
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
          <Link href="/ru/examples/listing-creation-ux" className="btn btn-ghost">
            <ChevronLeft className="w-5 h-5" />
            Назад к примерам
          </Link>
        </div>
        <div className="flex-none">
          <div className="badge badge-warning badge-lg">Без изменений Backend</div>
        </div>
      </div>

      {/* Main Content with Padding for Fixed Navbar */}
      <div className="pt-16">
        <AnimatePresence mode="wait">
          {currentView === "start" && renderStartView()}
          {currentView === "create" && renderCreateView()}
          {currentView === "preview" && renderPreviewView()}
        </AnimatePresence>
      </div>
    </>
  );
}