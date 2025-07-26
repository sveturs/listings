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
  AlertCircle,
  Loader2,
  ArrowRight,
  ImageIcon,
} from 'lucide-react';
import { useRouter } from 'next/navigation';
import { useAuthContext } from '@/contexts/AuthContext';
import { toast } from '@/utils/toast';
import { useTranslations } from 'next-intl';
import { claudeAI } from '@/services/ai/claude.service';

export default function AIPoweredListingCreationPage() {
  const router = useRouter();
  const t = useTranslations();
  const { user } = useAuthContext();
  const [currentView, setCurrentView] = useState<
    'upload' | 'process' | 'enhance' | 'publish'
  >('upload');
  const [images, setImages] = useState<string[]>([]);
  const [isProcessing, setIsProcessing] = useState(false);
  const [voiceRecording, setVoiceRecording] = useState(false);
  const [error, setError] = useState<string | null>(null);

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
  const [selectedVariant, setSelectedVariant] = useState(0);

  const fileInputRef = useRef<HTMLInputElement>(null);

  useEffect(() => {
    if (!user) {
      toast.error(t('create_listing.auth_required'));
      router.push('/');
    }
  }, [user, router, t]);

  // Convert image to base64
  const convertToBase64 = async (imageUrl: string): Promise<string> => {
    return new Promise((resolve, reject) => {
      const img = new window.Image();
      img.crossOrigin = 'anonymous';
      img.onload = () => {
        const canvas = document.createElement('canvas');
        canvas.width = img.width;
        canvas.height = img.height;
        const ctx = canvas.getContext('2d');
        ctx?.drawImage(img, 0, 0);
        const base64 = canvas.toDataURL('image/jpeg', 0.8);
        resolve(base64.split(',')[1]); // Remove data:image/jpeg;base64, prefix
      };
      img.onerror = reject;
      img.src = imageUrl;
    });
  };

  // Simulate AI processing
  const processImages = async () => {
    setIsProcessing(true);
    setCurrentView('process');
    setError(null);

    try {
      // Convert first image to base64
      const base64Image = await convertToBase64(images[0]);

      // Call Claude AI service
      const analysis = await claudeAI.analyzeProduct(base64Image);

      // Update state with AI analysis
      setAiData({
        ...analysis,
        selectedTitleIndex: 0,
        publishTime: '19:00',
      });

      // Perform A/B testing on title variants
      if (analysis.titleVariants.length > 1) {
        const abTestResult = await claudeAI.performABTesting(
          analysis.titleVariants
        );
        const bestIndex = analysis.titleVariants.findIndex(
          (t) => t === abTestResult.bestVariant
        );
        setSelectedVariant(bestIndex >= 0 ? bestIndex : 0);
      }

      setIsProcessing(false);
      setCurrentView('enhance');
    } catch (err) {
      console.error('AI processing error:', err);
      setError('Произошла ошибка при анализе. Проверьте подключение к интернету и попробуйте еще раз.');
      setIsProcessing(false);
      // Не переходим к следующему шагу при ошибке
    }
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

  const regenerateDescription = async () => {
    if (!aiData.title || !aiData.category) return;

    setIsProcessing(true);
    try {
      const newDescription = await claudeAI.generateOptimizedDescription(
        aiData.title,
        aiData.category,
        aiData.attributes
      );
      setAiData({ ...aiData, description: newDescription });
      toast.success('Описание обновлено!');
    } catch {
      toast.error('Ошибка генерации описания');
    } finally {
      setIsProcessing(false);
    }
  };

  const handleVoiceInput = () => {
    toast.info('Голосовой ввод будет доступен в ближайшее время');
    setVoiceRecording(!voiceRecording);
  };

  const handleSocialImport = (platform: string) => {
    toast.info(`Импорт из ${platform} будет доступен в ближайшее время`);
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
              <button
                onClick={() => handleSocialImport('Instagram')}
                className="btn btn-outline gap-2"
              >
                <Instagram className="w-4 h-4" />
                Импорт из Instagram
              </button>
              <button
                onClick={() => handleSocialImport('Facebook')}
                className="btn btn-outline gap-2"
              >
                <Facebook className="w-4 h-4" />
                Импорт из Facebook
              </button>
              <button
                onClick={handleVoiceInput}
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

            {error && (
              <div className="alert alert-error mb-4">
                <AlertCircle className="w-4 h-4" />
                <span>{error}</span>
              </div>
            )}

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
            <Loader2 className="w-5 h-5 animate-spin text-primary" />
            <span>Распознавание товара...</span>
          </motion.div>
          <motion.div
            initial={{ x: -20, opacity: 0 }}
            animate={{ x: 0, opacity: 1 }}
            transition={{ delay: 0.4 }}
            className="flex items-center gap-3"
          >
            <Loader2 className="w-5 h-5 animate-spin text-primary" />
            <span>Анализ рынка и цен...</span>
          </motion.div>
          <motion.div
            initial={{ x: -20, opacity: 0 }}
            animate={{ x: 0, opacity: 1 }}
            transition={{ delay: 0.6 }}
            className="flex items-center gap-3"
          >
            <Loader2 className="w-5 h-5 animate-spin text-primary" />
            <span>Генерация описания...</span>
          </motion.div>
          <motion.div
            initial={{ x: -20, opacity: 0 }}
            animate={{ x: 0, opacity: 1 }}
            transition={{ delay: 0.8 }}
            className="flex items-center gap-3"
          >
            <Loader2 className="w-5 h-5 animate-spin text-primary" />
            <span>SEO оптимизация...</span>
          </motion.div>
          <motion.div
            initial={{ x: -20, opacity: 0 }}
            animate={{ x: 0, opacity: 1 }}
            transition={{ delay: 1.0 }}
            className="flex items-center gap-3"
          >
            <Loader2 className="w-5 h-5 animate-spin text-primary" />
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
                  disabled={isProcessing}
                >
                  <RefreshCw className="w-4 h-4" />
                  Изменить
                </button>
              </div>

              <div className="space-y-2">
                {aiData.titleVariants.map((variant, index) => (
                  <div
                    key={index}
                    onClick={() =>
                      setAiData({ ...aiData, selectedTitleIndex: index })
                    }
                    className={`p-3 rounded-lg cursor-pointer transition-all ${
                      aiData.selectedTitleIndex === index
                        ? 'bg-primary/20 ring-2 ring-primary'
                        : 'bg-base-100 hover:bg-base-300'
                    }`}
                  >
                    <div className="flex items-center justify-between">
                      <span className="flex-1">{variant}</span>
                      {index === selectedVariant && (
                        <div className="badge badge-secondary badge-sm">
                          AI выбор
                        </div>
                      )}
                    </div>
                  </div>
                ))}
              </div>

              <div className="alert alert-info mt-4">
                <TestTube2 className="w-4 h-4" />
                <span className="text-sm">
                  AI будет тестировать варианты и покажет самый эффективный
                </span>
              </div>
            </div>
          </div>

          {/* Description */}
          <div className="card bg-base-200 mb-6">
            <div className="card-body">
              <div className="flex items-center justify-between mb-2">
                <h3 className="card-title">
                  <Brain className="w-5 h-5" />
                  Описание (AI оптимизированное)
                </h3>
                <button
                  onClick={regenerateDescription}
                  className="btn btn-ghost btn-sm gap-1"
                  disabled={isProcessing}
                >
                  <RefreshCw className="w-4 h-4" />
                  Обновить
                </button>
              </div>
              <textarea
                className="textarea textarea-bordered h-40"
                value={aiData.description}
                onChange={(e) =>
                  setAiData({ ...aiData, description: e.target.value })
                }
              />
              <div className="flex gap-2 mt-2">
                <div className="badge badge-success gap-1">
                  <Check className="w-3 h-3" />
                  SEO оптимизировано
                </div>
                <div className="badge badge-info gap-1">
                  <Sparkles className="w-3 h-3" />
                  Ключевые слова
                </div>
              </div>
            </div>
          </div>

          {/* Price with AI analysis */}
          <div className="card bg-base-200 mb-6">
            <div className="card-body">
              <h3 className="card-title">
                <BarChart3 className="w-5 h-5" />
                Цена (AI анализ рынка)
              </h3>
              <div className="form-control">
                <label className="input-group">
                  <input
                    type="number"
                    className="input input-bordered flex-1"
                    value={aiData.price}
                    onChange={(e) =>
                      setAiData({ ...aiData, price: e.target.value })
                    }
                  />
                  <span>РСД</span>
                </label>
              </div>

              {/* Price range visualization */}
              <div className="mt-4">
                <div className="flex justify-between text-sm mb-2">
                  <span>
                    Минимум: {aiData.priceRange.min.toLocaleString()} РСД
                  </span>
                  <span>
                    Максимум: {aiData.priceRange.max.toLocaleString()} РСД
                  </span>
                </div>
                <input
                  type="range"
                  min={aiData.priceRange.min}
                  max={aiData.priceRange.max}
                  value={aiData.price}
                  onChange={(e) =>
                    setAiData({ ...aiData, price: e.target.value })
                  }
                  className="range range-primary"
                />
              </div>

              <div className="alert alert-success mt-4">
                <TrendingUp className="w-4 h-4" />
                <span className="text-sm">
                  AI проанализировал 50+ похожих объявлений для оптимальной цены
                </span>
              </div>
            </div>
          </div>

          {/* Attributes */}
          {Object.keys(aiData.attributes).length > 0 && (
            <div className="card bg-base-200 mb-6">
              <div className="card-body">
                <h3 className="card-title">
                  <Package className="w-5 h-5" />
                  Характеристики (AI распознал)
                </h3>
                <div className="grid grid-cols-2 gap-4">
                  {Object.entries(aiData.attributes).map(([key, value]) => (
                    <div key={key} className="form-control">
                      <label className="label">
                        <span className="label-text capitalize">{key}</span>
                      </label>
                      <input
                        type="text"
                        className="input input-bordered input-sm"
                        value={value}
                        onChange={(e) =>
                          setAiData({
                            ...aiData,
                            attributes: {
                              ...aiData.attributes,
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
          )}

          {/* Translations */}
          <div className="card bg-base-200 mb-6">
            <div className="card-body">
              <h3 className="card-title">
                <Languages className="w-5 h-5" />
                Мультиязычность (автоперевод)
              </h3>
              <div className="space-y-4">
                {Object.entries(aiData.translations).map(([lang, trans]) => (
                  <div key={lang} className="border rounded-lg p-4">
                    <div className="flex items-center gap-2 mb-2">
                      <Globe className="w-4 h-4" />
                      <span className="font-semibold">
                        {lang === 'en'
                          ? 'English'
                          : lang === 'sr'
                            ? 'Српски'
                            : lang}
                      </span>
                    </div>
                    <p className="font-medium">{trans.title}</p>
                    <p className="text-sm text-base-content/70">
                      {trans.description}
                    </p>
                  </div>
                ))}
              </div>
            </div>
          </div>

          {/* Social posts */}
          <div className="card bg-base-200 mb-6">
            <div className="card-body">
              <h3 className="card-title">
                <Share2 className="w-5 h-5" />
                Посты для соцсетей (готовы к публикации)
              </h3>
              <div className="grid grid-cols-1 lg:grid-cols-3 gap-4">
                {Object.entries(aiData.socialPosts).map(([platform, post]) => (
                  <div key={platform} className="border rounded-lg p-4">
                    <div className="flex items-center gap-2 mb-2">
                      {platform === 'whatsapp' && (
                        <MessageCircle className="w-4 h-4 text-green-500" />
                      )}
                      {platform === 'telegram' && (
                        <Send className="w-4 h-4 text-blue-500" />
                      )}
                      {platform === 'instagram' && (
                        <Instagram className="w-4 h-4 text-pink-500" />
                      )}
                      <span className="font-semibold capitalize">
                        {platform}
                      </span>
                    </div>
                    <p className="text-sm whitespace-pre-wrap">{post}</p>
                  </div>
                ))}
              </div>
            </div>
          </div>

          {/* Effectiveness prediction */}
          <div className="card bg-gradient-to-r from-success/10 to-success/5 border-2 border-success/20 mb-6">
            <div className="card-body">
              <h3 className="card-title">
                <TrendingUp className="w-5 h-5" />
                Прогноз эффективности
              </h3>
              <div className="grid grid-cols-2 lg:grid-cols-4 gap-4">
                <div className="text-center">
                  <div className="text-2xl font-bold text-success">250+</div>
                  <div className="text-sm text-base-content/60">
                    просмотров в день
                  </div>
                </div>
                <div className="text-center">
                  <div className="text-2xl font-bold text-success">15-20</div>
                  <div className="text-sm text-base-content/60">сообщений</div>
                </div>
                <div className="text-center">
                  <div className="text-2xl font-bold text-success">3-5</div>
                  <div className="text-sm text-base-content/60">
                    дней до продажи
                  </div>
                </div>
                <div className="text-center">
                  <div className="text-2xl font-bold text-success">95%</div>
                  <div className="text-sm text-base-content/60">
                    вероятность продажи
                  </div>
                </div>
              </div>
            </div>
          </div>

          {/* Actions */}
          <div className="flex gap-3">
            <button
              onClick={() => setCurrentView('publish')}
              className="btn btn-primary flex-1"
            >
              Продолжить к публикации
              <ArrowRight className="w-4 h-4 ml-1" />
            </button>
            <button className="btn btn-ghost">Сохранить черновик</button>
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
      className="min-h-screen bg-base-200"
    >
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
              AI создал идеальное объявление!
            </h1>
            <p className="text-base-content/70">
              Готово к публикации с максимальной эффективностью
            </p>
          </motion.div>

          {/* Preview Card */}
          <motion.div
            initial={{ y: 20, opacity: 0 }}
            animate={{ y: 0, opacity: 1 }}
            transition={{ delay: 0.2 }}
            className="card bg-base-100 shadow-xl mb-6"
          >
            {images.length > 0 && (
              <figure className="relative">
                <div className="relative w-full h-96">
                  <Image
                    src={images[0]}
                    alt={aiData.title}
                    fill
                    className="object-cover"
                  />
                </div>
                {images.length > 1 && (
                  <div className="absolute bottom-4 right-4 badge badge-neutral gap-1">
                    <ImageIcon className="w-3 h-3" />+{images.length - 1}
                  </div>
                )}
              </figure>
            )}

            <div className="card-body">
              <h2 className="card-title text-2xl">
                {aiData.titleVariants[aiData.selectedTitleIndex] ||
                  aiData.title}
              </h2>

              <div className="text-3xl font-bold text-primary mb-4">
                {aiData.price} РСД
              </div>

              <p className="text-base-content/80 mb-4 whitespace-pre-wrap">
                {aiData.description}
              </p>

              <div className="flex items-center gap-4 text-sm text-base-content/60 mb-4">
                <span className="flex items-center gap-1">
                  <Eye className="w-4 h-4" />
                  250+ просмотров/день
                </span>
                <span className="flex items-center gap-1">
                  <Heart className="w-4 h-4" />
                  Высокий интерес
                </span>
              </div>

              {/* Tags */}
              <div className="flex flex-wrap gap-2 mb-4">
                {aiData.tags.map((tag, index) => (
                  <div key={index} className="badge badge-secondary">
                    {tag}
                  </div>
                ))}
              </div>
            </div>
          </motion.div>

          {/* Publishing options */}
          <div className="grid grid-cols-1 lg:grid-cols-2 gap-4 mb-8">
            <motion.div
              initial={{ x: -20, opacity: 0 }}
              animate={{ x: 0, opacity: 1 }}
              transition={{ delay: 0.3 }}
              className="card bg-base-100"
            >
              <div className="card-body">
                <h3 className="font-bold flex items-center gap-2">
                  <Clock className="w-5 h-5" />
                  Оптимальное время
                </h3>
                <p className="text-sm text-base-content/70">
                  AI рекомендует опубликовать в {aiData.publishTime} для
                  максимального охвата
                </p>
                <button className="btn btn-primary btn-sm mt-2">
                  <Calendar className="w-4 h-4 mr-1" />
                  Запланировать
                </button>
              </div>
            </motion.div>

            <motion.div
              initial={{ x: 20, opacity: 0 }}
              animate={{ x: 0, opacity: 1 }}
              transition={{ delay: 0.3 }}
              className="card bg-base-100"
            >
              <div className="card-body">
                <h3 className="font-bold flex items-center gap-2">
                  <Share2 className="w-5 h-5" />
                  Автопубликация в соцсети
                </h3>
                <p className="text-sm text-base-content/70">
                  Опубликовать в WhatsApp, Telegram, Instagram
                </p>
                <button className="btn btn-secondary btn-sm mt-2">
                  <Sparkles className="w-4 h-4 mr-1" />
                  Включить
                </button>
              </div>
            </motion.div>
          </div>

          {/* AI insights */}
          <motion.div
            initial={{ y: 20, opacity: 0 }}
            animate={{ y: 0, opacity: 1 }}
            transition={{ delay: 0.4 }}
            className="card bg-gradient-to-r from-primary/10 to-secondary/10 mb-8"
          >
            <div className="card-body">
              <h3 className="font-bold mb-4 flex items-center gap-2">
                <Brain className="w-5 h-5" />
                AI инсайты для вашего объявления
              </h3>
              <div className="space-y-3">
                <div className="flex items-start gap-3">
                  <TrendingUp className="w-5 h-5 text-success mt-0.5" />
                  <div>
                    <p className="font-semibold">Высокий спрос</p>
                    <p className="text-sm text-base-content/70">
                      iPhone 13 Pro входит в топ-5 самых популярных телефонов
                    </p>
                  </div>
                </div>
                <div className="flex items-start gap-3">
                  <Users className="w-5 h-5 text-info mt-0.5" />
                  <div>
                    <p className="font-semibold">Целевая аудитория</p>
                    <p className="text-sm text-base-content/70">
                      25-45 лет, средний и высокий доход, ценят качество
                    </p>
                  </div>
                </div>
                <div className="flex items-start gap-3">
                  <ThumbsUp className="w-5 h-5 text-warning mt-0.5" />
                  <div>
                    <p className="font-semibold">Рекомендации</p>
                    <p className="text-sm text-base-content/70">
                      Отвечайте быстро в первые 24 часа для лучшей конверсии
                    </p>
                  </div>
                </div>
              </div>
            </div>
          </motion.div>

          {/* Publish Actions */}
          <motion.div
            initial={{ y: 20, opacity: 0 }}
            animate={{ y: 0, opacity: 1 }}
            transition={{ delay: 0.5 }}
            className="flex gap-3"
          >
            <button
              onClick={() => {
                toast.success('Объявление опубликовано с AI оптимизацией!');
                router.push('/profile/listings');
              }}
              className="btn btn-primary btn-lg flex-1"
            >
              Опубликовать сейчас
              <Brain className="w-5 h-5 ml-1" />
            </button>
            <button
              onClick={() => setCurrentView('enhance')}
              className="btn btn-outline btn-lg"
            >
              Вернуться к редактированию
            </button>
          </motion.div>

          {/* Social proof */}
          <motion.div
            initial={{ opacity: 0 }}
            animate={{ opacity: 1 }}
            transition={{ delay: 0.6 }}
            className="text-center mt-8"
          >
            <div className="flex items-center justify-center gap-2 text-sm text-base-content/60">
              <Users className="w-4 h-4" />
              <span>
                <span className="font-semibold">834</span> продавца использовали
                AI сегодня
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
          <Link href="/sr/create-listing-choice" className="btn btn-ghost">
            <ChevronLeft className="w-5 h-5" />
            Назад к выбору
          </Link>
        </div>
        <div className="flex-none">
          <div className="badge badge-secondary badge-lg gap-1">
            <Brain className="w-3 h-3" />
            AI-Powered
          </div>
        </div>
      </div>

      {/* Main Content with Padding for Fixed Navbar */}
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
