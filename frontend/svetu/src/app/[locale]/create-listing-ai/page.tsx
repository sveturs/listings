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
  MapPin as MapPinIcon,
  Shield,
} from 'lucide-react';
import { useRouter } from 'next/navigation';
import { useAuthContext } from '@/contexts/AuthContext';
import { toast } from '@/utils/toast';
import { useTranslations, useLocale } from 'next-intl';
import { claudeAI } from '@/services/ai/claude.service';
import { categoryDetector } from '@/services/categoryDetector';
import configManager from '@/config';
import type { CreateListingState } from '@/contexts/CreateListingContext';
import { useAddressGeocoding } from '@/hooks/useAddressGeocoding';
import { extractLocationFromImages } from '@/utils/exifUtils';
import { formatAddressWithPrivacy } from '@/utils/addressPrivacy';
import LocationPicker from '@/components/GIS/LocationPicker';
import LocationPrivacySettingsWithAddress, {
  LocationPrivacyLevel,
} from '@/components/GIS/LocationPrivacySettingsWithAddress';
import { CarSelectorCompact } from '@/components/cars';
import { Car } from 'lucide-react';

export default function AIPoweredListingCreationPage() {
  const router = useRouter();
  const t = useTranslations('create_listing');
  const locale = useLocale();
  const { user } = useAuthContext();
  const [currentView, setCurrentView] = useState<
    'upload' | 'process' | 'enhance' | 'publish'
  >('upload');
  const [images, setImages] = useState<string[]>([]);
  const [imageFiles, setImageFiles] = useState<File[]>([]);
  const [isProcessing, setIsProcessing] = useState(false);
  const [voiceRecording, setVoiceRecording] = useState(false);
  const [error, setError] = useState<string | null>(null);
  const [categories, setCategories] = useState<
    Array<{ id: number; name: string; slug: string; translations?: any }>
  >([]);

  // Состояние для выбора автомобиля
  const [carSelection, setCarSelection] = useState<{
    make?: any;
    model?: any;
    generation?: any;
  }>({});

  // Category attributes
  const [categoryAttributes, setCategoryAttributes] = useState<any[]>([]);

  // Location states
  const [showLocationPicker, setShowLocationPicker] = useState(false);
  const [detectedLocation, setDetectedLocation] = useState<{
    latitude: number;
    longitude: number;
    source: 'exif' | 'profile' | 'manual';
  } | null>(null);

  // Состояние для приватности адреса
  const [locationPrivacyLevel, setLocationPrivacyLevel] =
    useState<LocationPrivacyLevel['id']>('street');
  const [fullAddress, setFullAddress] = useState<string>(''); // Полный адрес для сохранения в БД
  const [showPrivacySettings, setShowPrivacySettings] = useState(false);

  // AI generated data
  const [aiData, setAiData] = useState({
    title: '',
    titleVariants: [] as string[],
    selectedTitleIndex: 0,
    description: '',
    category: '',
    categoryProbabilities: [] as { name: string; probability: number }[],
    categoryDetectionStatsId: null as number | null,
    matchedKeywords: [] as string[],
    price: '',
    priceRange: { min: 0, max: 0 },
    attributes: {} as Record<string, string>,
    tags: [] as string[],
    suggestedPhotos: [] as string[],
    translations: {} as Record<string, { title: string; description: string }>,
    publishTime: '',
    socialPosts: {} as Record<string, string>,
    location: {
      city: '',
      region: '',
      suggestedLocation: '',
    },
    condition: 'used' as 'new' | 'used' | 'refurbished',
    insights: {} as Record<
      string,
      {
        demand: string;
        audience: string;
        recommendations: string;
      }
    >,
  });

  // A/B testing state
  const [selectedVariant, setSelectedVariant] = useState(0);

  // Preview language state
  const [previewLanguage, setPreviewLanguage] = useState('ru');

  const fileInputRef = useRef<HTMLInputElement>(null);

  // Геокодирование
  const geocoding = useAddressGeocoding({
    country: ['rs'], // Сербия
    language: 'sr',
  });

  // Флаг для предотвращения повторной загрузки профиля
  const profileLoadedRef = useRef(false);

  useEffect(() => {
    if (!user) {
      toast.error(t('auth_required'));
      router.push('/');
      return;
    }

    // Предотвращаем повторную загрузку
    if (profileLoadedRef.current) {
      return;
    }

    // Загружаем профиль пользователя для получения адреса по умолчанию
    const loadUserProfile = async () => {
      try {
        const { tokenManager } = await import('@/utils/tokenManager');
        const token = tokenManager.getAccessToken();

        if (!token) {
          console.log('No access token available, skipping profile load');
          return;
        }

        // Используем абсолютный API URL
        const apiUrl = configManager.getApiUrl();
        console.log(
          'Making profile request to:',
          `${apiUrl}/api/v1/users/profile`
        );

        const response = await fetch(`${apiUrl}/api/v1/users/profile`, {
          headers: {
            Authorization: `Bearer ${token}`,
          },
        });

        if (response.ok) {
          const profileData = await response.json();
          if (profileData.data?.city && profileData.data?.country) {
            console.log(
              'Using default address from user profile:',
              profileData.data
            );
            setAiData((prev) => ({
              ...prev,
              location: {
                city: profileData.data.city,
                region: profileData.data.country,
                suggestedLocation: `${profileData.data.city}, ${profileData.data.country}`,
              },
            }));

            // Не геокодируем сразу, чтобы избежать циклических зависимостей
            // Геокодирование будет выполнено при необходимости позже
          }
        }
      } catch (error) {
        console.log('Failed to load user profile:', error);
      } finally {
        // Отмечаем, что профиль был загружен
        profileLoadedRef.current = true;
      }
    };

    loadUserProfile();
  }, [user, router, t]); // Убираем validateAddress из зависимостей

  // Загружаем категории при монтировании компонента
  useEffect(() => {
    const loadCategories = async () => {
      try {
        const apiUrl = configManager.getApiUrl();
        const response = await fetch(`${apiUrl}/api/v1/marketplace/categories`);
        if (response.ok) {
          const data = await response.json();
          if (data.data) {
            setCategories(data.data);
          }
        }
      } catch (error) {
        console.error('Failed to load categories:', error);
      }
    };

    loadCategories();
  }, []);

  // Загружаем атрибуты при изменении категории
  const loadCategoryAttributes = async (categoryId: number) => {
    try {
      const apiUrl = configManager.getApiUrl();
      const response = await fetch(
        `${apiUrl}/api/v1/marketplace/categories/${categoryId}/attributes`
      );
      if (response.ok) {
        const data = await response.json();
        if (data.data) {
          setCategoryAttributes(data.data);
          console.log('Loaded attributes for category:', categoryId, data.data);
        }
      }
    } catch (error) {
      console.error('Failed to load category attributes:', error);
    }
  };

  // Get category data by name/slug
  const getCategoryData = (
    categoryName: string
  ): { id: number; name: string; slug: string } => {
    // Проверка на undefined или пустую строку
    if (!categoryName) {
      return { id: 1008, name: 'Hrana i piće', slug: 'food-beverages' };
    }

    const normalizedName = categoryName.toLowerCase().trim();

    // 1. Обработка подкатегорий из AI анализа
    if (normalizedName.includes('/')) {
      // Разбираем путь категории
      const parts = normalizedName.split('/');

      // Проверяем для автозапчастей и шин
      if (parts.includes('auto-parts') && parts.includes('tires-wheels')) {
        // Ищем категорию "Gume i točkovi" (ID: 1304)
        const tiresCategory = categories.find(
          (cat) =>
            cat.id === 1304 ||
            cat.slug === 'tires-and-wheels' ||
            cat.name.toLowerCase().includes('gume') ||
            cat.name.toLowerCase().includes('točkovi')
        );
        if (tiresCategory) {
          return {
            id: tiresCategory.id,
            name: tiresCategory.name,
            slug: tiresCategory.slug,
          };
        }
      }

      // Для других автозапчастей
      if (parts.includes('auto-parts') && !parts.includes('tires-wheels')) {
        // Возвращаем категорию "Auto delovi" (ID: 1303)
        const autoPartsCategory = categories.find(
          (cat) =>
            cat.id === 1303 ||
            cat.slug === 'auto-parts' ||
            cat.name.toLowerCase().includes('auto delovi')
        );
        if (autoPartsCategory) {
          return {
            id: autoPartsCategory.id,
            name: autoPartsCategory.name,
            slug: autoPartsCategory.slug,
          };
        }
      }

      // Для автомобилей
      if (parts.includes('cars')) {
        const carsCategory = categories.find(
          (cat) =>
            cat.id === 1301 ||
            cat.slug === 'cars' ||
            cat.name.toLowerCase().includes('lični automobili')
        );
        if (carsCategory) {
          return {
            id: carsCategory.id,
            name: carsCategory.name,
            slug: carsCategory.slug,
          };
        }
      }
    }

    // 2. Поиск по имени/slug категории
    const category = categories.find(
      (cat) =>
        cat.slug === normalizedName ||
        cat.name.toLowerCase() === normalizedName ||
        normalizedName.includes(cat.slug)
    );

    if (category) {
      return { id: category.id, name: category.name, slug: category.slug };
    }

    // 3. Маппинг основных категорий
    const categoryMap: Record<string, number> = {
      electronics: 1001,
      fashion: 1002,
      automotive: 1003,
      'real-estate': 1004,
      'home-garden': 1005,
      agriculture: 1006,
      industrial: 1007,
      'food-beverages': 1008,
      services: 1009,
      'sports-recreation': 1010,
    };

    const mappedId = categoryMap[normalizedName];
    if (mappedId) {
      const mappedCategory = categories.find((cat) => cat.id === mappedId);
      if (mappedCategory) {
        return {
          id: mappedCategory.id,
          name: mappedCategory.name,
          slug: mappedCategory.slug,
        };
      }
    }

    // Возвращаем автомобильную категорию по умолчанию для automotive
    if (normalizedName.includes('automotive')) {
      return { id: 1003, name: 'Automobili', slug: 'automotive' };
    }

    return { id: 1008, name: 'Hrana i piće', slug: 'food-beverages' };
  };

  // Проверка, является ли категория автомобильной
  const isCarCategory = (categorySlug: string): boolean => {
    if (!categorySlug) return false;
    const carCategories = [
      'automotive',
      'cars',
      'automobili',
      'licni-automobili',
    ];
    const normalizedSlug = categorySlug.toLowerCase();

    // Проверяем основные автомобильные категории
    if (carCategories.includes(normalizedSlug)) {
      return true;
    }

    // Проверяем ID категории
    const categoryData = getCategoryData(categorySlug);
    if (categoryData.id >= 1301 && categoryData.id <= 1302) {
      return true;
    }

    // Проверяем по ID в диапазоне 10100-10199 (старый диапазон)
    if (categoryData.id >= 10100 && categoryData.id < 10200) {
      return true;
    }

    return false;
  };

  // Convert image to base64
  const convertToBase64 = async (imageUrl: string): Promise<string> => {
    return new Promise((resolve, reject) => {
      const img = new window.Image();
      // Don't set crossOrigin for blob URLs
      if (!imageUrl.startsWith('blob:')) {
        img.crossOrigin = 'anonymous';
      }
      img.onload = () => {
        const canvas = document.createElement('canvas');
        canvas.width = img.width;
        canvas.height = img.height;
        const ctx = canvas.getContext('2d');
        ctx?.drawImage(img, 0, 0);
        const base64 = canvas.toDataURL('image/jpeg', 0.8);
        resolve(base64.split(',')[1]); // Remove data:image/jpeg;base64, prefix
      };
      img.onerror = (error) => {
        console.error('Error loading image:', error);
        reject(error);
      };
      img.src = imageUrl;
    });
  };

  // Simulate AI processing
  const processImages = async () => {
    setIsProcessing(true);
    setCurrentView('process');
    setError(null);

    try {
      // Пытаемся извлечь геолокацию из EXIF данных перед AI анализом
      console.log('Processing images, checking EXIF data...');
      let exifLocationData = null;
      try {
        const exifLocation = await extractLocationFromImages(imageFiles);
        console.log('EXIF extraction during processing:', exifLocation);
        if (exifLocation) {
          console.log('Found EXIF location during processing:', exifLocation);
          setDetectedLocation({
            latitude: exifLocation.latitude,
            longitude: exifLocation.longitude,
            source: 'exif',
          });

          // Геокодируем координаты в адрес
          const geocodedAddress = await geocoding.reverseGeocode(
            exifLocation.latitude,
            exifLocation.longitude
          );

          if (geocodedAddress) {
            console.log('Geocoded address during processing:', geocodedAddress);

            // Сохраняем полный адрес
            const fullAddressStr =
              geocodedAddress.address || geocodedAddress.place_name;
            setFullAddress(fullAddressStr);

            // Форматируем адрес с учетом приватности по умолчанию (улица)
            const privateAddress = formatAddressWithPrivacy(fullAddressStr, {
              showHouseNumber: false,
              showStreet: true,
              showCity: true,
              showRegion: true,
            });

            // Сохраняем адрес для использования ниже
            exifLocationData = {
              city:
                geocodedAddress.city ||
                geocodedAddress.address_components.city ||
                'Белград',
              region:
                geocodedAddress.region ||
                geocodedAddress.address_components.district ||
                'Сербия',
              suggestedLocation: privateAddress,
            };
          }
        }
      } catch (error) {
        console.log('Error extracting EXIF during processing:', error);
      }

      // Convert first image to base64
      const base64Image = await convertToBase64(images[0]);

      // Call Claude AI service with user's language
      const analysis = await claudeAI.analyzeProduct(base64Image, locale);

      // Update state with AI analysis with validation
      setAiData({
        ...aiData,
        ...analysis,
        // Валидация categoryProbabilities
        categoryProbabilities: Array.isArray(analysis.categoryProbabilities)
          ? analysis.categoryProbabilities
          : [],
        // Валидация других массивов
        titleVariants: Array.isArray(analysis.titleVariants)
          ? analysis.titleVariants
          : [analysis.title || ''],
        tags: Array.isArray(analysis.tags) ? analysis.tags : [],
        suggestedPhotos: Array.isArray(analysis.suggestedPhotos)
          ? analysis.suggestedPhotos
          : [],
        // Извлекаем только число из цены, убираем валюту
        price: analysis.price?.replace(/[^\d.,]/g, '').replace(',', '.') || '',
        selectedTitleIndex: 0,
        publishTime: '19:00',
        location: exifLocationData || {
          city: analysis.location?.city || 'Белград',
          region: analysis.location?.region || 'Сербия',
          suggestedLocation: analysis.location?.suggestedLocation || '',
        },
        condition: analysis.condition || 'used',
        insights: analysis.insights || {},
      });

      // Используем новую систему определения категории если есть categoryHints
      if (analysis.categoryHints) {
        try {
          const detectionResult = await categoryDetector.detectCategory(
            analysis.categoryHints,
            analysis.title,
            analysis.description,
            locale
          );

          if (detectionResult.category_id) {
            // Находим категорию по ID
            const detectedCategory = categories.find(
              (cat) => cat.id === detectionResult.category_id
            );
            if (detectedCategory) {
              // Обновляем категорию в aiData
              setAiData((prev) => ({
                ...prev,
                category: detectedCategory.slug,
                categoryDetectionStatsId: detectionResult.stats_id || null,
                categoryProbabilities:
                  detectionResult.alternative_categories?.map((alt) => ({
                    name: alt.category_slug,
                    probability: alt.confidence_score,
                  })) || [],
                // Сохраняем найденные ключевые слова для последующей передачи на backend
                matchedKeywords:
                  detectionResult.debug_info?.matched_keywords || [],
              }));

              // Сбрасываем выбор автомобиля при автоматическом определении категории
              setCarSelection({});

              // Загружаем атрибуты для определенной категории
              await loadCategoryAttributes(detectionResult.category_id);

              // Показываем пользователю уровень уверенности
              if (detectionResult.confidence_score < 0.7) {
                toast.info(t('ai.category_detection.low_confidence'));
              } else {
                toast.success(
                  t('ai.category_detection.success', {
                    category: detectedCategory.name,
                  })
                );
              }
            }
          }
        } catch (error) {
          console.error('Category detection error:', error);
          // Fallback к старой системе
          const categoryData = getCategoryData(analysis.category);
          if (categoryData && categoryData.id) {
            await loadCategoryAttributes(categoryData.id);
          }
        }
      } else {
        // Fallback к старой системе если нет categoryHints
        const categoryData = getCategoryData(analysis.category);
        if (categoryData && categoryData.id) {
          await loadCategoryAttributes(categoryData.id);
        }
      }

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
      setError(t('ai.errors.ai_processing'));
      setIsProcessing(false);
      // Не переходим к следующему шагу при ошибке
    }
  };

  const handleImageUpload = async (e: React.ChangeEvent<HTMLInputElement>) => {
    console.log('handleImageUpload called');
    const files = e.target.files;
    if (files) {
      const filesArray = Array.from(files);
      console.log('Processing', filesArray.length, 'files');
      const newImages = filesArray.map((file) => URL.createObjectURL(file));
      setImages([...images, ...newImages].slice(0, 8));
      setImageFiles([...imageFiles, ...filesArray].slice(0, 8));

      // Пытаемся извлечь геолокацию из EXIF данных
      console.log('Attempting to extract EXIF location...');
      try {
        const exifLocation = await extractLocationFromImages(filesArray);
        console.log('EXIF extraction result:', exifLocation);
        if (exifLocation) {
          console.log('Detected location from EXIF:', exifLocation);
          setDetectedLocation({
            latitude: exifLocation.latitude,
            longitude: exifLocation.longitude,
            source: 'exif',
          });

          // Геокодируем координаты в адрес
          const geocodedAddress = await geocoding.reverseGeocode(
            exifLocation.latitude,
            exifLocation.longitude
          );

          if (geocodedAddress) {
            console.log('Geocoded address:', geocodedAddress);

            // Сохраняем полный адрес
            const fullAddressStr =
              geocodedAddress.address || geocodedAddress.place_name;
            setFullAddress(fullAddressStr);

            // Форматируем адрес с учетом приватности по умолчанию (улица)
            const privateAddress = formatAddressWithPrivacy(fullAddressStr, {
              showHouseNumber: false,
              showStreet: true,
              showCity: true,
              showRegion: true,
            });

            // Обновляем данные локации в AI данных с реальным адресом
            setAiData((prev) => ({
              ...prev,
              location: {
                city:
                  geocodedAddress.city ||
                  geocodedAddress.address_components.city ||
                  'Белград',
                region:
                  geocodedAddress.region ||
                  geocodedAddress.address_components.district ||
                  'Сербия',
                suggestedLocation: privateAddress,
              },
            }));

            // Используем нейтральное сообщение без упоминания AI
            toast.success(t('location_detected_from_photo'));
          } else {
            // Если не удалось геокодировать, показываем координаты
            setAiData((prev) => ({
              ...prev,
              location: {
                ...prev.location,
                suggestedLocation: `Координаты: ${exifLocation.latitude.toFixed(4)}, ${exifLocation.longitude.toFixed(4)}`,
              },
            }));
          }
        }
      } catch (error) {
        console.log('Error processing location data:', error);
      }
    }
  };

  const removeImage = (index: number) => {
    setImages(images.filter((_, i) => i !== index));
    setImageFiles(imageFiles.filter((_, i) => i !== index));
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
      toast.success(t('ai.success_messages.description_updated'));
    } catch {
      toast.error(t('ai.errors.description_generation'));
    } finally {
      setIsProcessing(false);
    }
  };

  const handleVoiceInput = () => {
    toast.info(t('ai.features.voice_coming_soon'));
    setVoiceRecording(!voiceRecording);
  };

  const handleSocialImport = (platform: string) => {
    toast.info(t('ai.features.import_coming_soon', { platform }));
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
            {t('ai.title')}
          </h1>
          <p className="text-xl text-base-content/70 mb-8">
            {t('ai.subtitle')}
          </p>

          <div className="flex justify-center gap-6 mb-8">
            <div className="text-center">
              <div className="text-3xl font-bold text-primary">
                {t('ai.stats.creation_time')}
              </div>
              <div className="text-sm text-base-content/60">
                {t('ai.stats.creation_label')}
              </div>
            </div>
            <div className="text-center">
              <div className="text-3xl font-bold text-success">
                {t('ai.stats.accuracy')}
              </div>
              <div className="text-sm text-base-content/60">
                {t('ai.stats.accuracy_label')}
              </div>
            </div>
            <div className="text-center">
              <div className="text-3xl font-bold text-secondary">
                {t('ai.stats.languages')}
              </div>
              <div className="text-sm text-base-content/60">
                {t('ai.stats.languages_label')}
              </div>
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
                  {t('ai.upload.title')}
                </h2>
                <p className="text-base-content/70 mb-6">
                  {t('ai.upload.subtitle')}
                </p>
                <div className="flex gap-4 justify-center">
                  <div className="badge badge-lg badge-primary gap-2">
                    <Brain className="w-4 h-4" />
                    {t('ai.upload.ai_recognition')}
                  </div>
                  <div className="badge badge-lg badge-secondary gap-2">
                    <Zap className="w-4 h-4" />
                    {t('ai.upload.time')}
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
                {t('ai.upload.import_instagram')}
              </button>
              <button
                onClick={() => handleSocialImport('Facebook')}
                className="btn btn-outline gap-2"
              >
                <Facebook className="w-4 h-4" />
                {t('ai.upload.import_facebook')}
              </button>
              <button
                onClick={handleVoiceInput}
                className={`btn ${voiceRecording ? 'btn-error' : 'btn-outline'} gap-2`}
              >
                <Mic className="w-4 h-4" />
                {voiceRecording
                  ? t('ai.upload.stop_recording')
                  : t('ai.upload.voice_input')}
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
                    {t('ai.upload.add_more')}
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
              {t('ai.upload.create_with_ai')}
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

        <h2 className="text-2xl font-bold mb-4">{t('ai.processing.title')}</h2>

        <div className="space-y-4 text-left max-w-md mx-auto">
          <motion.div
            initial={{ x: -20, opacity: 0 }}
            animate={{ x: 0, opacity: 1 }}
            transition={{ delay: 0.2 }}
            className="flex items-center gap-3"
          >
            <Loader2 className="w-5 h-5 animate-spin text-primary" />
            <span>{t('ai.processing.recognizing')}</span>
          </motion.div>
          <motion.div
            initial={{ x: -20, opacity: 0 }}
            animate={{ x: 0, opacity: 1 }}
            transition={{ delay: 0.4 }}
            className="flex items-center gap-3"
          >
            <Loader2 className="w-5 h-5 animate-spin text-primary" />
            <span>{t('ai.processing.analyzing_market')}</span>
          </motion.div>
          <motion.div
            initial={{ x: -20, opacity: 0 }}
            animate={{ x: 0, opacity: 1 }}
            transition={{ delay: 0.6 }}
            className="flex items-center gap-3"
          >
            <Loader2 className="w-5 h-5 animate-spin text-primary" />
            <span>{t('ai.processing.generating_description')}</span>
          </motion.div>
          <motion.div
            initial={{ x: -20, opacity: 0 }}
            animate={{ x: 0, opacity: 1 }}
            transition={{ delay: 0.8 }}
            className="flex items-center gap-3"
          >
            <Loader2 className="w-5 h-5 animate-spin text-primary" />
            <span>{t('ai.processing.seo_optimization')}</span>
          </motion.div>
          <motion.div
            initial={{ x: -20, opacity: 0 }}
            animate={{ x: 0, opacity: 1 }}
            transition={{ delay: 1.0 }}
            className="flex items-center gap-3"
          >
            <Loader2 className="w-5 h-5 animate-spin text-primary" />
            <span>{t('ai.processing.creating_translations')}</span>
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
              <h3 className="font-bold">{t('ai.enhance.success_title')}</h3>
              <p>{t('ai.enhance.success_subtitle')}</p>
            </div>
          </motion.div>

          {/* Photos section */}
          <div className="card bg-base-200 mb-6">
            <div className="card-body">
              <h3 className="card-title">
                <Camera className="w-5 h-5" />
                {t('ai.enhance.photos_title')}
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
                        {t('ai.enhance.main_photo')}
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
                      {t('ai.enhance.ai_recommends')}
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
                {t('ai.enhance.category_title')}
              </h3>
              <div className="space-y-2">
                {(aiData.categoryProbabilities || []).map((cat, index) => {
                  if (!cat || !cat.name) {
                    return null;
                  }
                  const categoryData = getCategoryData(cat.name);
                  const categoryName = categoryData?.name || cat.name;
                  const isSelected = index === 0; // Первая категория имеет наибольшую вероятность

                  return (
                    <div
                      key={index}
                      className={`flex items-center justify-between p-2 rounded-lg ${
                        isSelected ? 'bg-primary/10 ring-2 ring-primary' : ''
                      }`}
                    >
                      <span className={isSelected ? 'font-semibold' : ''}>
                        {categoryName}
                      </span>
                      <div className="flex items-center gap-2">
                        <progress
                          className={`progress ${isSelected ? 'progress-primary' : 'progress-base-300'} w-32`}
                          value={cat.probability}
                          max="100"
                        ></progress>
                        <span className="text-sm font-medium">
                          {cat.probability}%
                        </span>
                      </div>
                    </div>
                  );
                })}
              </div>

              {aiData.categoryProbabilities &&
                aiData.categoryProbabilities.length > 0 && (
                  <div className="mt-4">
                    <div className="alert alert-info">
                      <Check className="w-4 h-4" />
                      <span className="text-sm">
                        {t('ai.enhance.category_auto_selected', {
                          category:
                            getCategoryData(aiData.category)?.name ||
                            aiData.category,
                        })}
                      </span>
                    </div>
                    <div className="mt-2">
                      <label className="label">
                        <span className="label-text">
                          {t('ai.enhance.manual_category_selection')}
                        </span>
                      </label>
                      <select
                        className="select select-bordered w-full"
                        value={getCategoryData(aiData.category)?.id || ''}
                        onChange={(e) => {
                          const selectedCat = categories.find(
                            (cat) => cat.id === Number(e.target.value)
                          );
                          if (selectedCat) {
                            setAiData({
                              ...aiData,
                              category: selectedCat.slug,
                            });
                            loadCategoryAttributes(selectedCat.id);
                            // Сбрасываем выбор автомобиля при смене категории
                            setCarSelection({});
                          }
                        }}
                      >
                        <option value="">
                          {t('ai.enhance.select_category')}
                        </option>
                        {categories.map((cat) => (
                          <option key={cat.id} value={cat.id}>
                            {cat.name}
                          </option>
                        ))}
                      </select>
                    </div>
                  </div>
                )}
            </div>
          </div>

          {/* Title with A/B testing */}
          <div className="card bg-base-200 mb-6">
            <div className="card-body">
              <div className="flex items-center justify-between mb-2">
                <h3 className="card-title">
                  <TestTube2 className="w-5 h-5" />
                  {t('ai.enhance.title_ab_testing')}
                </h3>
                <button
                  onClick={regenerateTitle}
                  className="btn btn-ghost btn-sm gap-1"
                  disabled={isProcessing}
                >
                  <RefreshCw className="w-4 h-4" />
                  {t('ai.enhance.change')}
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
                          {t('ai.enhance.ai_choice')}
                        </div>
                      )}
                    </div>
                  </div>
                ))}
              </div>

              <div className="alert alert-info mt-4">
                <TestTube2 className="w-4 h-4" />
                <span className="text-sm">{t('ai.enhance.ai_will_test')}</span>
              </div>
            </div>
          </div>

          {/* Description */}
          <div className="card bg-base-200 mb-6">
            <div className="card-body">
              <div className="flex items-center justify-between mb-2">
                <h3 className="card-title">
                  <Brain className="w-5 h-5" />
                  {t('ai.enhance.description_title')}
                </h3>
                <button
                  onClick={regenerateDescription}
                  className="btn btn-ghost btn-sm gap-1"
                  disabled={isProcessing}
                >
                  <RefreshCw className="w-4 h-4" />
                  {t('ai.description_section.update')}
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
                  {t('ai.enhance.seo_optimized')}
                </div>
                <div className="badge badge-info gap-1">
                  <Sparkles className="w-3 h-3" />
                  {t('ai.enhance.keywords')}
                </div>
              </div>
            </div>
          </div>

          {/* Location and condition */}
          <div className="card bg-base-200 mb-6">
            <div className="card-body">
              <h3 className="card-title">
                <Globe className="w-5 h-5" />
                {t('ai.enhance.location_condition_title')}
              </h3>
              <div className="grid grid-cols-2 gap-4">
                <div>
                  <label className="label">
                    <span className="label-text">{t('ai.enhance.city')}</span>
                  </label>
                  <input
                    type="text"
                    className="input input-bordered"
                    value={aiData.location.city}
                    onChange={(e) =>
                      setAiData({
                        ...aiData,
                        location: { ...aiData.location, city: e.target.value },
                      })
                    }
                  />
                </div>
                <div>
                  <label className="label">
                    <span className="label-text">
                      {t('ai.enhance.condition')}
                    </span>
                  </label>
                  <select
                    className="select select-bordered"
                    value={aiData.condition}
                    onChange={(e) =>
                      setAiData({
                        ...aiData,
                        condition: e.target.value as
                          | 'new'
                          | 'used'
                          | 'refurbished',
                      })
                    }
                  >
                    <option value="new">{t('ai.enhance.condition_new')}</option>
                    <option value="used">
                      {t('ai.enhance.condition_used')}
                    </option>
                    <option value="refurbished">
                      {t('ai.enhance.condition_refurbished')}
                    </option>
                  </select>
                </div>
              </div>

              {/* Location source info and manual picker */}
              <div className="mt-4 space-y-3">
                {detectedLocation && (
                  <div className="alert alert-info">
                    <Globe className="w-4 h-4" />
                    <div>
                      <p className="font-semibold text-sm">
                        {t('ai.location.detected_from')}{' '}
                        {detectedLocation.source === 'exif'
                          ? t('ai.location.from_photo')
                          : detectedLocation.source === 'profile'
                            ? t('ai.location.from_profile')
                            : t('ai.location.manually')}
                      </p>
                      <p className="text-xs">
                        {t('ai.location.coordinates')}{' '}
                        {detectedLocation.latitude.toFixed(4)},{' '}
                        {detectedLocation.longitude.toFixed(4)}
                      </p>
                    </div>
                  </div>
                )}

                {!detectedLocation && (
                  <div className="alert alert-warning">
                    <AlertCircle className="w-4 h-4" />
                    <div>
                      <p className="font-semibold text-sm">
                        {t('ai.location.not_detected')}
                      </p>
                      <p className="text-xs">
                        {t('ai.location.not_detected_hint')}
                      </p>
                    </div>
                  </div>
                )}

                <button
                  onClick={() => setShowLocationPicker(true)}
                  className="btn btn-outline btn-sm gap-2"
                >
                  <MapPinIcon className="w-4 h-4" />
                  {detectedLocation
                    ? t('ai.location.refine_on_map')
                    : t('ai.location.select_on_map')}
                </button>
              </div>
              {aiData.location.suggestedLocation && (
                <>
                  <div className="alert alert-info mt-4">
                    <Lightbulb className="w-4 h-4" />
                    <span className="text-sm">
                      {fullAddress
                        ? `${t('location_detected_from_photo')}: ${aiData.location.suggestedLocation}`
                        : `${t('ai.location.ai_suggests')} ${aiData.location.suggestedLocation}`}
                    </span>
                  </div>

                  {/* Настройки приватности */}
                  {(detectedLocation || fullAddress) && (
                    <div className="mt-4">
                      <button
                        type="button"
                        onClick={() =>
                          setShowPrivacySettings(!showPrivacySettings)
                        }
                        className="btn btn-outline btn-sm gap-2"
                      >
                        <Shield className="w-4 h-4" />
                        {t('privacy_settings')}
                        <svg
                          className={`w-4 h-4 transition-transform ${showPrivacySettings ? 'rotate-180' : ''}`}
                          fill="none"
                          stroke="currentColor"
                          viewBox="0 0 24 24"
                        >
                          <path
                            strokeLinecap="round"
                            strokeLinejoin="round"
                            strokeWidth={2}
                            d="M19 9l-7 7-7-7"
                          />
                        </svg>
                      </button>

                      {showPrivacySettings && (
                        <div className="mt-4 p-4 bg-base-100 rounded-lg">
                          <LocationPrivacySettingsWithAddress
                            selectedLevel={locationPrivacyLevel}
                            onLevelChange={(level) => {
                              setLocationPrivacyLevel(level);
                              // Обновляем отображаемый адрес
                              if (fullAddress) {
                                const options = {
                                  showHouseNumber: level === 'exact',
                                  showStreet:
                                    level === 'exact' || level === 'street',
                                  showCity: true,
                                  showRegion: level !== 'city',
                                  showCountry: level === 'city',
                                };
                                const updatedAddress = formatAddressWithPrivacy(
                                  fullAddress,
                                  options
                                );
                                setAiData((prev) => ({
                                  ...prev,
                                  location: {
                                    ...prev.location,
                                    suggestedLocation: updatedAddress,
                                  },
                                }));
                              }
                            }}
                            location={
                              detectedLocation
                                ? {
                                    lat: detectedLocation.latitude,
                                    lng: detectedLocation.longitude,
                                  }
                                : undefined
                            }
                            fullAddress={fullAddress}
                            showPreview={true}
                          />
                        </div>
                      )}
                    </div>
                  )}
                </>
              )}
            </div>
          </div>

          {/* Price with AI analysis */}
          <div className="card bg-base-200 mb-6">
            <div className="card-body">
              <h3 className="card-title">
                <BarChart3 className="w-5 h-5" />
                {t('ai.price.title')}
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
                  <span>{t('ai.price.currency')}</span>
                </label>
              </div>

              {/* Price range visualization */}
              <div className="mt-4">
                <div className="flex justify-between text-sm mb-2">
                  <span>
                    {t('ai.price.min')} {aiData.priceRange.min.toLocaleString()}{' '}
                    {t('ai.price.currency')}
                  </span>
                  <span>
                    {t('ai.price.max')} {aiData.priceRange.max.toLocaleString()}{' '}
                    {t('ai.price.currency')}
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
                <span className="text-sm">{t('ai.price.ai_analyzed')}</span>
              </div>
            </div>
          </div>

          {/* Car Selector for automotive categories */}
          {aiData.category && isCarCategory(aiData.category) && (
            <div className="card bg-base-200 mb-6">
              <div className="card-body">
                <h3 className="card-title text-base mb-4">
                  <Car className="w-5 h-5" />
                  {t('cars.selectCar')}
                </h3>
                <CarSelectorCompact
                  value={carSelection}
                  onChange={(selection) => {
                    setCarSelection(selection);
                    // Обновляем атрибуты при выборе машины
                    if (selection.make && selection.model) {
                      setAiData((prev) => ({
                        ...prev,
                        attributes: {
                          ...prev.attributes,
                          // Марка
                          make: selection.make?.name || '',
                          // Модель
                          model: selection.model?.name || '',
                        },
                      }));

                      // Если заголовок пуст или это стандартный заголовок, обновляем его
                      const currentTitle =
                        aiData.titleVariants[aiData.selectedTitleIndex] ||
                        aiData.title;
                      if (
                        !currentTitle ||
                        currentTitle.toLowerCase().includes('автомобиль')
                      ) {
                        const newTitle =
                          `${selection.make?.name || ''} ${selection.model?.name || ''}`.trim();
                        setAiData((prev) => ({
                          ...prev,
                          title: newTitle,
                          titleVariants: [
                            newTitle,
                            ...prev.titleVariants.slice(1),
                          ],
                        }));
                      }
                    }
                  }}
                />
              </div>
            </div>
          )}

          {/* Attributes */}
          {(Object.keys(aiData.attributes).length > 0 ||
            categoryAttributes.length > 0) && (
            <div className="card bg-base-200 mb-6">
              <div className="card-body">
                <h3 className="card-title">
                  <Package className="w-5 h-5" />
                  {t('ai.attributes.title')}{' '}
                  {Object.keys(aiData.attributes).length > 0
                    ? t('ai.attributes.ai_detected')
                    : t('ai.attributes.category_based')}
                </h3>
                <div className="grid grid-cols-2 gap-4">
                  {/* Отображаем атрибуты от AI */}
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

                  {/* Отображаем атрибуты категории, которых нет в AI данных */}
                  {categoryAttributes
                    .filter((attr) => !aiData.attributes[attr.name])
                    .map((attr) => (
                      <div key={attr.id} className="form-control">
                        <label className="label">
                          <span className="label-text">
                            {attr.display_name || attr.name}
                            {attr.is_required && (
                              <span className="text-error">*</span>
                            )}
                          </span>
                        </label>
                        {attr.attribute_type === 'select' &&
                        attr.options &&
                        Array.isArray(attr.options) ? (
                          <select
                            className="select select-bordered select-sm"
                            onChange={(e) =>
                              setAiData({
                                ...aiData,
                                attributes: {
                                  ...aiData.attributes,
                                  [attr.name]: e.target.value,
                                },
                              })
                            }
                          >
                            <option value="">
                              {t('ai.attributes.choose')}
                            </option>
                            {attr.options.map((opt: any) => (
                              <option
                                key={opt.id || opt.value}
                                value={opt.value}
                              >
                                {opt.display_value || opt.value}
                              </option>
                            ))}
                          </select>
                        ) : attr.attribute_type === 'boolean' ? (
                          <input
                            type="checkbox"
                            className="checkbox"
                            onChange={(e) =>
                              setAiData({
                                ...aiData,
                                attributes: {
                                  ...aiData.attributes,
                                  [attr.name]: e.target.checked.toString(),
                                },
                              })
                            }
                          />
                        ) : attr.attribute_type === 'number' ? (
                          <input
                            type="number"
                            className="input input-bordered input-sm"
                            placeholder={attr.placeholder}
                            onChange={(e) =>
                              setAiData({
                                ...aiData,
                                attributes: {
                                  ...aiData.attributes,
                                  [attr.name]: e.target.value,
                                },
                              })
                            }
                          />
                        ) : (
                          <input
                            type="text"
                            className="input input-bordered input-sm"
                            placeholder={attr.placeholder}
                            onChange={(e) =>
                              setAiData({
                                ...aiData,
                                attributes: {
                                  ...aiData.attributes,
                                  [attr.name]: e.target.value,
                                },
                              })
                            }
                          />
                        )}
                      </div>
                    ))}
                </div>
              </div>
            </div>
          )}

          {/* Translations */}
          <div className="card bg-base-200 mb-6">
            <div className="card-body">
              <div className="flex items-center justify-between mb-2">
                <h3 className="card-title">
                  <Languages className="w-5 h-5" />
                  {t('ai.translations.title')}
                </h3>
                <button
                  onClick={async () => {
                    setIsProcessing(true);
                    try {
                      const translations = await claudeAI.translateContent(
                        {
                          title:
                            aiData.titleVariants[aiData.selectedTitleIndex] ||
                            aiData.title,
                          description: aiData.description,
                        },
                        ['en', 'sr', 'ru']
                      );
                      setAiData({ ...aiData, translations });
                      toast.success(
                        t('ai.success_messages.translations_updated')
                      );
                    } catch (error) {
                      console.error('Translation error:', error);
                      toast.error(t('ai.errors.translation_update'));
                    } finally {
                      setIsProcessing(false);
                    }
                  }}
                  className="btn btn-ghost btn-sm gap-1"
                  disabled={isProcessing}
                >
                  <RefreshCw className="w-4 h-4" />
                  {t('ai.translations.update')}
                </button>
              </div>
              <div className="space-y-4">
                {Object.entries(aiData.translations).map(([lang, trans]) => (
                  <div key={lang} className="border rounded-lg p-4">
                    <div className="flex items-center gap-2 mb-3">
                      <Globe className="w-4 h-4" />
                      <span className="font-semibold">
                        {lang === 'en'
                          ? 'English'
                          : lang === 'sr'
                            ? 'Српски'
                            : lang === 'ru'
                              ? 'Русский'
                              : lang}
                      </span>
                    </div>

                    <div className="space-y-3">
                      <div>
                        <label className="label">
                          <span className="label-text text-sm">
                            {t('ai.translations.language_title')}
                          </span>
                        </label>
                        <input
                          type="text"
                          className="input input-bordered w-full"
                          value={trans.title}
                          onChange={(e) => {
                            setAiData({
                              ...aiData,
                              translations: {
                                ...aiData.translations,
                                [lang]: {
                                  ...trans,
                                  title: e.target.value,
                                },
                              },
                            });
                          }}
                        />
                      </div>

                      <div>
                        <label className="label">
                          <span className="label-text text-sm">
                            {t('ai.translations.language_description')}
                          </span>
                        </label>
                        <textarea
                          className="textarea textarea-bordered w-full h-24"
                          value={trans.description}
                          onChange={(e) => {
                            setAiData({
                              ...aiData,
                              translations: {
                                ...aiData.translations,
                                [lang]: {
                                  ...trans,
                                  description: e.target.value,
                                },
                              },
                            });
                          }}
                        />
                      </div>
                    </div>
                  </div>
                ))}
              </div>

              <div className="alert alert-info mt-4">
                <Brain className="w-4 h-4" />
                <span className="text-sm">{t('ai.translations.info')}</span>
              </div>
            </div>
          </div>

          {/* Social posts */}
          <div className="card bg-base-200 mb-6">
            <div className="card-body">
              <h3 className="card-title">
                <Share2 className="w-5 h-5" />
                {t('ai.social.title')}
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
                {t('ai.prediction.title')}
              </h3>
              <div className="grid grid-cols-2 lg:grid-cols-4 gap-4">
                <div className="text-center">
                  <div className="text-2xl font-bold text-success">250+</div>
                  <div className="text-sm text-base-content/60">
                    {t('ai.prediction.views')}
                  </div>
                </div>
                <div className="text-center">
                  <div className="text-2xl font-bold text-success">15-20</div>
                  <div className="text-sm text-base-content/60">
                    {t('ai.prediction.messages')}
                  </div>
                </div>
                <div className="text-center">
                  <div className="text-2xl font-bold text-success">3-5</div>
                  <div className="text-sm text-base-content/60">
                    {t('ai.prediction.days_to_sell')}
                  </div>
                </div>
                <div className="text-center">
                  <div className="text-2xl font-bold text-success">95%</div>
                  <div className="text-sm text-base-content/60">
                    {t('ai.prediction.sell_probability')}
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
              {t('ai.actions.continue_to_publish')}
              <ArrowRight className="w-4 h-4 ml-1" />
            </button>
            <button className="btn btn-ghost">
              {t('ai.actions.save_draft')}
            </button>
          </div>
        </div>
      </div>
    </motion.div>
  );

  const publishListing = async () => {
    let listingData: CreateListingState | undefined;

    try {
      setIsProcessing(true);

      // Получаем данные категории
      const categoryData = getCategoryData(aiData.category);

      // Валидация обязательных атрибутов для выбранной категории
      const requiredAttributes = categoryAttributes.filter(
        (attr) => attr.is_required === true
      );

      const missingRequiredAttributes = requiredAttributes.filter((reqAttr) => {
        const attributeValue = aiData.attributes[reqAttr.name];
        const stringValue = String(attributeValue || '').trim();
        return !stringValue;
      });

      if (missingRequiredAttributes.length > 0) {
        const missingNames = missingRequiredAttributes
          .map((attr) => attr.display_name || attr.name)
          .join(', ');

        toast.error(
          t('ai.errors.missing_attributes', { attributes: missingNames })
        );
        setCurrentView('enhance');
        setIsProcessing(false);
        return;
      }

      // Используем координаты из detectedLocation или геокодируем адрес
      let latitude = 0;
      let longitude = 0;

      if (detectedLocation) {
        // Используем уже определенные координаты
        latitude = detectedLocation.latitude;
        longitude = detectedLocation.longitude;
        console.log(`Using ${detectedLocation.source} location:`, {
          latitude,
          longitude,
        });
      } else {
        // Геокодируем адрес как fallback
        const fullAddress =
          `${aiData.location.suggestedLocation || ''} ${aiData.location.city || 'Белград'} ${aiData.location.region || 'Сербия'}`.trim();

        try {
          const geoResult = await geocoding.validateAddress(fullAddress);
          if (geoResult.success && geoResult.location) {
            latitude = geoResult.location.lat;
            longitude = geoResult.location.lng;
            console.log('Geocoding successful:', { latitude, longitude });
          }
        } catch (error) {
          console.error('Geocoding failed:', error);
        }
      }

      // Преобразуем AI данные в формат CreateListingState
      listingData = {
        category: categoryData,
        title: aiData.titleVariants[aiData.selectedTitleIndex] || aiData.title,
        description: aiData.description,
        price: parseFloat(aiData.price) || 0,
        currency: 'RSD',
        condition: aiData.condition,

        // Локация
        location: {
          latitude,
          longitude,
          // Отправляем адрес с учетом приватности
          address:
            aiData.location.suggestedLocation ||
            aiData.location.city ||
            'Белград',
          city: aiData.location.city || 'Белград',
          region: aiData.location.region || 'Сербия',
          country: 'Сербия',
        },

        // Региональная система доверия
        trust: {
          phoneVerified: false,
          preferredMeetingType: 'personal',
          meetingLocations: [],
          availableHours: '',
          localReputation: 0,
        },

        // Платежи
        payment: {
          methods: ['cash'],
          codEnabled: false,
          codPrice: 0,
          personalMeeting: true,
          deliveryOptions: [],
          negotiablePrice: false,
          bundleDeals: false,
        },

        // Локализация
        localization: {
          script: 'cyrillic',
          language: 'sr',
          traditionalUnits: false,
          regionalPhrases: [],
        },

        // Pijaca 2.0
        pijaca: {
          vendorStallStyle: '',
          regularCustomers: false,
          traditionalStyle: false,
        },

        // Изображения и атрибуты
        images: images,
        mainImageIndex: 0,
        attributes: Object.entries(aiData.attributes).reduce(
          (acc, [key, value]) => {
            // Пропускаем пустые значения
            const stringValue = String(value || '');
            if (!stringValue || stringValue.trim() === '') {
              return acc;
            }

            // Находим соответствующий атрибут из загруженных
            const categoryAttr = categoryAttributes.find(
              (attr) => attr.name === key
            );

            if (categoryAttr && categoryAttr.id > 0) {
              // Проверяем, что значение подходит для типа атрибута
              const attributeValue: any = {
                attribute_id: categoryAttr.id,
                attribute_name: categoryAttr.name,
                display_name: categoryAttr.display_name,
                attribute_type: categoryAttr.attribute_type,
              };

              // Добавляем значение в зависимости от типа
              if (categoryAttr.attribute_type === 'text') {
                attributeValue.text_value = stringValue;
              } else if (categoryAttr.attribute_type === 'number') {
                const numericValue = parseFloat(stringValue);
                if (!isNaN(numericValue)) {
                  attributeValue.numeric_value = numericValue;
                } else {
                  // Если число невалидно, пропускаем атрибут
                  return acc;
                }
              } else if (categoryAttr.attribute_type === 'boolean') {
                attributeValue.boolean_value =
                  stringValue === 'true' ||
                  stringValue === 'да' ||
                  stringValue === 'yes';
              } else if (categoryAttr.attribute_type === 'select') {
                // Для select проверяем, что значение есть в options
                // Backend возвращает options как {values: [...]} объект
                if (
                  categoryAttr.options &&
                  categoryAttr.options.values &&
                  Array.isArray(categoryAttr.options.values)
                ) {
                  const validValue = categoryAttr.options.values.find(
                    (optValue: string) => optValue === stringValue
                  );
                  if (validValue) {
                    attributeValue.text_value = validValue;
                  } else {
                    // Если значение не найдено в options, пропускаем
                    console.warn(
                      `Invalid option value "${stringValue}" for attribute "${categoryAttr.name}". Valid options:`,
                      categoryAttr.options.values
                    );
                    return acc;
                  }
                } else {
                  // Fallback для случаев без options
                  attributeValue.text_value = stringValue;
                }
              } else if (categoryAttr.attribute_type === 'multiselect') {
                attributeValue.json_value = Array.isArray(value)
                  ? value
                  : [stringValue];
              }

              if (categoryAttr.unit) {
                attributeValue.unit = categoryAttr.unit;
              }

              acc[key] = attributeValue;
            }
            // Убираем блок else - больше не создаем атрибуты с ID=0

            return acc;
          },
          {} as Record<string, any>
        ),

        // Переводы
        translations: aiData.translations,

        // Язык оригинала
        originalLanguage: locale,

        // ID статистики определения категории
        categoryDetectionStatsId: aiData.categoryDetectionStatsId || undefined,

        // Ключевые слова, использованные для определения категории
        detectedKeywords: aiData.matchedKeywords || undefined,

        // Метаданные
        isPublished: true,
        isDraft: false,
      };

      // Создаем объявление
      const { ListingsService } = await import('@/services/listings');
      const response = await ListingsService.createListing(
        listingData,
        locationPrivacyLevel
      );

      if (response.data?.id) {
        // Загружаем изображения
        if (imageFiles.length > 0) {
          await ListingsService.uploadImages(response.data.id, imageFiles, 0);
        }

        toast.success(t('ai.success_messages.listing_created'));
        router.push(`/marketplace/${response.data.id}`);
      }
    } catch (error) {
      console.error('Error publishing listing:', error);

      // Улучшенная обработка ошибок
      let errorMessage = t('ai.errors.publish_error');

      if (error instanceof Error) {
        // Проверяем на специфические ошибки атрибутов
        if (
          error.message.includes('attribute') ||
          error.message.includes('атрибут')
        ) {
          errorMessage = t('ai.errors.attribute_error');
        } else if (
          error.message.includes('validation') ||
          error.message.includes('валидация')
        ) {
          errorMessage = t('ai.errors.validation_error');
        } else if (
          error.message.includes('network') ||
          error.message.includes('fetch')
        ) {
          errorMessage = t('ai.errors.network_error');
        }

        console.error('Detailed error information:', {
          message: error.message,
          stack: error.stack,
          listingData: listingData
            ? JSON.stringify(listingData, null, 2)
            : 'undefined',
        });
      }

      toast.error(errorMessage);
    } finally {
      setIsProcessing(false);
    }
  };

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
            <h1 className="text-2xl font-bold mb-2">{t('ai.publish.title')}</h1>
            <p className="text-base-content/70">{t('ai.publish.subtitle')}</p>
          </motion.div>

          {/* Language Switcher */}
          <div className="flex justify-center gap-2 mb-6">
            <button
              onClick={() => setPreviewLanguage('ru')}
              className={`btn btn-sm ${previewLanguage === 'ru' ? 'btn-primary' : 'btn-outline'}`}
            >
              🇷🇺 Русский
            </button>
            <button
              onClick={() => setPreviewLanguage('en')}
              className={`btn btn-sm ${previewLanguage === 'en' ? 'btn-primary' : 'btn-outline'}`}
            >
              🇬🇧 English
            </button>
            <button
              onClick={() => setPreviewLanguage('sr')}
              className={`btn btn-sm ${previewLanguage === 'sr' ? 'btn-primary' : 'btn-outline'}`}
            >
              🇷🇸 Српски
            </button>
          </div>

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
                {previewLanguage === 'ru'
                  ? aiData.titleVariants[aiData.selectedTitleIndex] ||
                    aiData.title
                  : aiData.translations[previewLanguage]?.title ||
                    aiData.titleVariants[aiData.selectedTitleIndex] ||
                    aiData.title}
              </h2>

              <div className="text-3xl font-bold text-primary mb-4">
                {aiData.price} РСД
              </div>

              <p className="text-base-content/80 mb-4 whitespace-pre-wrap">
                {previewLanguage === 'ru'
                  ? aiData.description
                  : aiData.translations[previewLanguage]?.description ||
                    aiData.description}
              </p>

              {/* Атрибуты */}
              {Object.keys(aiData.attributes).length > 0 && (
                <div className="mb-4">
                  <h3 className="font-semibold mb-2">
                    {t('ai.publish.specifications')}
                  </h3>
                  <div className="grid grid-cols-2 gap-2">
                    {Object.entries(aiData.attributes).map(([key, value]) => {
                      const categoryAttr = categoryAttributes.find(
                        (attr) => attr.name === key
                      );

                      // Получаем переведенное имя атрибута
                      let displayName = categoryAttr?.display_name || key;
                      if (
                        categoryAttr?.translations &&
                        categoryAttr.translations[previewLanguage]
                      ) {
                        displayName =
                          categoryAttr.translations[previewLanguage];
                      }

                      // Получаем переведенное значение для select атрибутов
                      let displayValue = String(value);
                      if (
                        categoryAttr?.attribute_type === 'select' &&
                        categoryAttr?.option_translations?.[previewLanguage]?.[
                          value
                        ]
                      ) {
                        displayValue =
                          categoryAttr.option_translations[previewLanguage][
                            value
                          ];
                      }

                      return (
                        <div
                          key={key}
                          className="flex justify-between py-1 border-b border-base-200"
                        >
                          <span className="text-base-content/70">
                            {displayName}:
                          </span>
                          <span className="font-medium">{displayValue}</span>
                        </div>
                      );
                    })}
                  </div>
                </div>
              )}

              <div className="flex items-center gap-4 text-sm text-base-content/60 mb-4">
                <span className="flex items-center gap-1">
                  <Eye className="w-4 h-4" />
                  250+ {t('ai.publish.views_per_day')}
                </span>
                <span className="flex items-center gap-1">
                  <Heart className="w-4 h-4" />
                  {t('ai.publish.high_interest')}
                </span>
              </div>

              {/* Tags */}
              {aiData.tags.length > 0 && (
                <div className="flex flex-wrap gap-2 mb-4">
                  {aiData.tags.map((tag, index) => (
                    <div key={index} className="badge badge-secondary">
                      {tag}
                    </div>
                  ))}
                </div>
              )}
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
                  {t('ai.publish.optimal_time')}
                </h3>
                <p className="text-sm text-base-content/70">
                  {t('ai.publish.optimal_time_desc', {
                    time: aiData.publishTime,
                  })}
                </p>
                <button className="btn btn-primary btn-sm mt-2">
                  <Calendar className="w-4 h-4 mr-1" />
                  {t('ai.publish.schedule')}
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
                  {t('ai.publish.social_auto')}
                </h3>
                <p className="text-sm text-base-content/70">
                  {t('ai.publish.social_auto_desc')}
                </p>
                <button className="btn btn-secondary btn-sm mt-2">
                  <Sparkles className="w-4 h-4 mr-1" />
                  {t('ai.publish.enable')}
                </button>
              </div>
            </motion.div>
          </div>

          {/* AI insights */}
          {aiData.insights && Object.keys(aiData.insights).length > 0 && (
            <motion.div
              initial={{ y: 20, opacity: 0 }}
              animate={{ y: 0, opacity: 1 }}
              transition={{ delay: 0.4 }}
              className="card bg-gradient-to-r from-primary/10 to-secondary/10 mb-8"
            >
              <div className="card-body">
                <h3 className="font-bold mb-4 flex items-center gap-2">
                  <Brain className="w-5 h-5" />
                  {t('ai.publish.ai_insights_title')}
                </h3>
                <div className="space-y-3">
                  {aiData.insights[previewLanguage] && (
                    <>
                      <div className="flex items-start gap-3">
                        <TrendingUp className="w-5 h-5 text-success mt-0.5" />
                        <div>
                          <p className="font-semibold">
                            {t('ai.publish.demand')}
                          </p>
                          <p className="text-sm text-base-content/70">
                            {aiData.insights[previewLanguage].demand}
                          </p>
                        </div>
                      </div>
                      <div className="flex items-start gap-3">
                        <Users className="w-5 h-5 text-info mt-0.5" />
                        <div>
                          <p className="font-semibold">
                            {t('ai.publish.target_audience')}
                          </p>
                          <p className="text-sm text-base-content/70">
                            {aiData.insights[previewLanguage].audience}
                          </p>
                        </div>
                      </div>
                      <div className="flex items-start gap-3">
                        <ThumbsUp className="w-5 h-5 text-warning mt-0.5" />
                        <div>
                          <p className="font-semibold">
                            {t('ai.publish.recommendations')}
                          </p>
                          <p className="text-sm text-base-content/70">
                            {aiData.insights[previewLanguage].recommendations}
                          </p>
                        </div>
                      </div>
                    </>
                  )}
                </div>
              </div>
            </motion.div>
          )}

          {/* Publish Actions */}
          <motion.div
            initial={{ y: 20, opacity: 0 }}
            animate={{ y: 0, opacity: 1 }}
            transition={{ delay: 0.5 }}
            className="flex gap-3"
          >
            <button
              onClick={publishListing}
              className="btn btn-primary btn-lg flex-1"
              disabled={isProcessing}
            >
              {isProcessing ? (
                <>
                  <Loader2 className="w-5 h-5 mr-2 animate-spin" />
                  {t('ai.actions.publishing')}
                </>
              ) : (
                <>
                  {t('ai.actions.publish_now')}
                  <Brain className="w-5 h-5 ml-1" />
                </>
              )}
            </button>
            <button
              onClick={() => setCurrentView('enhance')}
              className="btn btn-outline btn-lg"
            >
              {t('ai.actions.back_to_edit')}
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
              <span>{t('ai.publish.social_proof', { count: 834 })}</span>
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
            href={`/${locale}/create-listing-choice`}
            className="btn btn-ghost"
          >
            <ChevronLeft className="w-5 h-5" />
            {t('ai.back_to_choice')}
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

      {/* Location Picker Modal */}
      {showLocationPicker && (
        <div className="fixed inset-0 bg-black/50 z-50 flex items-center justify-center p-4">
          <div className="bg-base-100 rounded-lg max-w-4xl w-full max-h-[90vh] overflow-hidden">
            <div className="p-6">
              <div className="flex items-center justify-between mb-4">
                <h3 className="text-lg font-bold">
                  {t('ai.publish.location_picker_title')}
                </h3>
                <button
                  onClick={() => setShowLocationPicker(false)}
                  className="btn btn-ghost btn-sm btn-circle"
                >
                  <X className="w-4 h-4" />
                </button>
              </div>

              <LocationPicker
                value={
                  detectedLocation
                    ? {
                        latitude: detectedLocation.latitude,
                        longitude: detectedLocation.longitude,
                        address: aiData.location.suggestedLocation || '',
                        city: aiData.location.city || '',
                        region: aiData.location.region || '',
                        country: 'Сербия',
                        confidence: 0.8,
                      }
                    : undefined
                }
                onChange={(location) => {
                  setDetectedLocation({
                    latitude: location.latitude,
                    longitude: location.longitude,
                    source: 'manual',
                  });

                  // Сохраняем полный адрес
                  setFullAddress(location.address);

                  // Форматируем адрес с учетом текущего уровня приватности
                  const options = {
                    showHouseNumber: locationPrivacyLevel === 'exact',
                    showStreet:
                      locationPrivacyLevel === 'exact' ||
                      locationPrivacyLevel === 'street',
                    showCity: true,
                    showRegion: locationPrivacyLevel !== 'city',
                    showCountry: locationPrivacyLevel === 'city',
                  };
                  const privateAddress = formatAddressWithPrivacy(
                    location.address,
                    options
                  );

                  setAiData((prev) => ({
                    ...prev,
                    location: {
                      city: location.city,
                      region: location.region,
                      suggestedLocation: privateAddress,
                    },
                  }));

                  setShowLocationPicker(false);
                  toast.success(t('ai.success_messages.location_updated'));
                }}
                height="400px"
                showCurrentLocation={true}
              />
            </div>
          </div>
        </div>
      )}
    </>
  );
}
