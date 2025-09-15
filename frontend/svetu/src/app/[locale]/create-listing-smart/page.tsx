'use client';

import React, { useState, useRef, useEffect, useCallback } from 'react';
import Link from 'next/link';
import Image from 'next/image';
import { motion, AnimatePresence } from 'framer-motion';
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
  Clock as ClockIcon,
  FileText,
  Users,
  Loader2,
  Car,
} from 'lucide-react';
import { useRouter, useParams } from 'next/navigation';
import { useAuthContext } from '@/contexts/AuthContext';
import { toast } from '@/utils/toast';
import { useTranslations } from 'next-intl';
import { ListingsService } from '@/services/listings';
import type { CreateListingState } from '@/contexts/CreateListingContext';
import type { components } from '@/types/generated/api';
import CategorySelector from '@/components/listing/CategorySelector';
import CategoryAttributes from '@/components/listing/CategoryAttributes';
import { CarSelectorCompact } from '@/components/cars';
import configManager from '@/config';

type MarketplaceCategory =
  components['schemas']['backend_internal_domain_models.MarketplaceCategory'];

interface AttributeValue {
  attribute_id: number;
  attribute_name: string;
  display_name: string;
  attribute_type: string;
  text_value?: string;
  numeric_value?: number;
  boolean_value?: boolean;
  unit?: string;
}

export default function CreateListingSmartPage() {
  const router = useRouter();
  const params = useParams();
  const locale = params.locale as string;
  const t = useTranslations('create_listing');
  const { user } = useAuthContext();
  const [currentView, setCurrentView] = useState<
    'start' | 'create' | 'preview'
  >('start');
  const [quickMode, setQuickMode] = useState(false);
  const [isSubmitting, setIsSubmitting] = useState(false);
  const [categories, setCategories] = useState<MarketplaceCategory[]>([]);
  const [selectedCategory, setSelectedCategory] =
    useState<MarketplaceCategory | null>(null);
  const [categorySelectedManually, setCategorySelectedManually] =
    useState(false);
  const [suggestedCategories, setSuggestedCategories] = useState<
    Array<{
      categoryId: number;
      categoryName: string;
      categorySlug: string;
      confidenceScore: number;
    }>
  >([]);
  const [imageFiles, setImageFiles] = useState<File[]>([]);
  const [formData, setFormData] = useState({
    images: [] as string[],
    category: '',
    categoryId: 0,
    title: '',
    price: '',
    description: '',
    location: '',
    city: 'Белград',
    country: 'Србија',
    deliveryMethods: ['pickup'],
    attributes: {} as Record<number, AttributeValue>,
  });
  const _suggestions = {
    title: '',
    category: '',
    price: '',
    description: '',
  };
  // Состояние для сравнения с похожими
  const [showPriceComparison, setShowPriceComparison] = useState(false);
  const [similarListings, setSimilarListings] = useState<any[]>([]);
  const [isLoadingSimilar, setIsLoadingSimilar] = useState(false);

  // Drag & Drop состояние
  const [isDragging, setIsDragging] = useState(false);
  const [draggedIndex, setDraggedIndex] = useState<number | null>(null);

  // Состояние для выбора автомобиля
  const [carSelection, setCarSelection] = useState<{
    make?: any;
    model?: any;
    generation?: any;
  }>({});

  const fileInputRef = useRef<HTMLInputElement>(null);

  useEffect(() => {
    if (!user) {
      toast.error(t('auth_required'));
      router.push('/');
    }
  }, [user, router, t]);

  const fetchCategories = useCallback(async () => {
    try {
      const apiUrl = configManager.getApiUrl();
      const response = await fetch(
        `${apiUrl}/api/v1/marketplace/categories?lang=${locale}`
      );
      if (response.ok) {
        const data = await response.json();
        if (data.data && Array.isArray(data.data)) {
          setCategories(data.data);
        }
      }
    } catch (error) {
      console.error('Error fetching categories:', error);
    }
  }, [locale]);

  // Загружаем категории при монтировании
  useEffect(() => {
    fetchCategories();
  }, [fetchCategories]);

  // Category-specific attributes (legacy, now using API)
  const _categoryAttributes: Record<
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

    fashion: `👟 Размер: [точный размер]
📏 Стелька: [длина в см]
🧵 Материал: [кожа/текстиль/синтетика]
✨ Состояние: [новое/б/у]
📸 Носил(а): [сколько раз/период]`,

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

  // Поиск реальных похожих объявлений
  const fetchSimilarListings = useCallback(
    async (categoryId: number, title: string) => {
      console.log('fetchSimilarListings called with:', { categoryId, title });
      if (!categoryId || !title || title.length < 3) {
        console.log('Skipping fetch: invalid params');
        return;
      }

      setIsLoadingSimilar(true);
      try {
        // Извлекаем ключевые слова из заголовка
        const keywords = title
          .toLowerCase()
          .split(' ')
          .filter((word) => word.length > 2)
          .slice(0, 3)
          .join(' ');

        // Поиск через API
        const searchParams = new URLSearchParams({
          query: keywords,
          category_id: categoryId.toString(),
          page: '1',
          limit: '5',
          sort_by: 'date',
          sort_order: 'desc',
          language: locale,
        });

        const apiUrl = configManager.getApiUrl();
        const response = await fetch(`${apiUrl}/api/v1/search?${searchParams}`);
        if (response.ok) {
          const data = await response.json();
          console.log('Search API response:', data); // Для отладки
          if (data.items && Array.isArray(data.items)) {
            // Преобразуем данные для отображения
            const listings = data.items.slice(0, 3).map((item: any) => {
              const createdDate = new Date(item.created_at);
              const now = new Date();
              const daysAgo = Math.floor(
                (now.getTime() - createdDate.getTime()) / (1000 * 60 * 60 * 24)
              );

              return {
                id: item.product_id,
                title: item.name,
                price: item.price,
                views: Math.floor(Math.random() * 200) + 50, // Пока нет views в ответе search
                daysAgo: daysAgo || 0,
                sold: false, // В search API нет статуса
                image: item.images?.[0]?.url,
              };
            });
            setSimilarListings(listings);
          }
        }
      } catch (error) {
        console.error('Error fetching similar listings:', error);
      } finally {
        setIsLoadingSimilar(false);
      }
    },
    [locale]
  );

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

  // Получаем популярные категории из загруженных (legacy, replaced by CategorySelector)
  const _getPopularCategories = () => {
    const popularSlugs = [
      'electronics',
      'fashion',
      'home-garden',
      'automotive',
    ];
    const iconMap: Record<string, string> = {
      electronics: '📱',
      fashion: '👗',
      'home-garden': '🏠',
      automotive: '🚗',
    };
    const gradientMap: Record<string, string> = {
      electronics: 'from-blue-500 to-purple-500',
      fashion: 'from-pink-500 to-rose-500',
      'home-garden': 'from-green-500 to-emerald-500',
      automotive: 'from-orange-500 to-red-500',
    };

    return categories
      .filter((cat) => popularSlugs.includes(cat.slug || ''))
      .map((cat) => ({
        id: cat.id,
        slug: cat.slug || '',
        name: cat.translations?.name || cat.name || '',
        icon: iconMap[cat.slug || ''] || '📦',
        gradient: gradientMap[cat.slug || ''] || 'from-gray-500 to-gray-600',
      }));
  };

  // Умное определение категории через API
  const detectCategoryByTitle = useCallback(
    async (title: string) => {
      if (!title || title.length < 3) return;

      try {
        // Извлекаем ключевые слова из названия
        const keywords = title
          .toLowerCase()
          .split(' ')
          .filter((word) => word.length > 2);

        const apiUrl = configManager.getApiUrl();
        const response = await fetch(
          `${apiUrl}/api/v1/marketplace/categories/detect`,
          {
            method: 'POST',
            headers: {
              'Content-Type': 'application/json',
            },
            body: JSON.stringify({
              title: title,
              keywords: keywords,
              language: locale,
            }),
          }
        );

        if (response.ok) {
          const data = await response.json();
          if (data.data && data.data.category_id) {
            const detectedCategory = categories.find(
              (c) => c.id === data.data.category_id
            );

            if (detectedCategory) {
              // Проверяем, изменилась ли категория
              const categoryChanged =
                selectedCategory?.id !== detectedCategory.id;

              // Автоматически обновляем категорию только если:
              // 1. Категория не была выбрана вручную
              // 2. И категория действительно изменилась
              if (!categorySelectedManually && categoryChanged) {
                setSelectedCategory(detectedCategory);
                setFormData((prev) => ({
                  ...prev,
                  category: detectedCategory.slug || '',
                  categoryId: detectedCategory.id || 0,
                }));

                // Показываем уведомление только при изменении категории
                const confidence = Math.round(data.data.confidence_score * 100);
                toast.info(
                  `Категория определена: ${detectedCategory.translations?.name || detectedCategory.name} (${confidence}% уверенность)`
                );
              }

              // Всегда обновляем список предложенных категорий
              if (
                data.data.alternative_categories &&
                data.data.alternative_categories.length > 0
              ) {
                const alternatives = data.data.alternative_categories.map(
                  (alt: any) => ({
                    categoryId: alt.category_id,
                    categoryName: alt.category_name,
                    categorySlug: alt.category_slug,
                    confidenceScore: alt.confidence_score,
                  })
                );
                setSuggestedCategories([
                  {
                    categoryId: data.data.category_id,
                    categoryName: data.data.category_name,
                    categorySlug: data.data.category_slug,
                    confidenceScore: data.data.confidence_score,
                  },
                  ...alternatives,
                ]);
              } else {
                // Только основная категория
                setSuggestedCategories([
                  {
                    categoryId: data.data.category_id,
                    categoryName: data.data.category_name,
                    categorySlug: data.data.category_slug,
                    confidenceScore: data.data.confidence_score,
                  },
                ]);
              }
            }
          } else {
            // Если категория не определена - очищаем предложения
            setSuggestedCategories([]);
          }
        }
      } catch (error) {
        console.error('Ошибка определения категории:', error);
      }
    },
    [categories, locale, selectedCategory, categorySelectedManually]
  );

  // Эффект для умного определения категории при изменении названия
  useEffect(() => {
    // Определяем категорию при изменении названия
    if (formData.title && formData.title.length > 3) {
      const timeoutId = setTimeout(() => {
        detectCategoryByTitle(formData.title);
      }, 800); // Задержка 800мс после последнего ввода

      return () => clearTimeout(timeoutId);
    } else if (formData.title.length === 0) {
      // Если название очищено - очищаем предложенные категории и сбрасываем флаг ручного выбора
      setSuggestedCategories([]);
      setCategorySelectedManually(false);
      // Можно также очистить выбранную категорию если нужно
      // setSelectedCategory(null);
      // setFormData(prev => ({ ...prev, category: '', categoryId: 0 }));
    }
  }, [formData.title, detectCategoryByTitle]);

  useEffect(() => {
    // Обновляем похожие объявления при изменении категории или заголовка
    console.log('useEffect triggered, formData:', {
      title: formData.title,
      categoryId: formData.categoryId,
    });
    if (formData.title && formData.categoryId) {
      // Используем debounce для оптимизации запросов
      const timeoutId = setTimeout(() => {
        fetchSimilarListings(formData.categoryId, formData.title);
      }, 500);

      // Показываем сравнение цен если есть заголовок
      if (formData.title.length > 3) {
        setShowPriceComparison(true);
      }

      return () => clearTimeout(timeoutId);
    } else {
      setSimilarListings([]);
      setShowPriceComparison(false);
    }
  }, [formData.categoryId, formData.title, fetchSimilarListings]);

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
      const newFiles = Array.from(files);
      const newImages = newFiles.map((file) => URL.createObjectURL(file));

      // Проверка качества изображений
      newImages.forEach((imgUrl, index) => {
        const img = new window.Image();
        img.src = imgUrl;
        img.onload = () => {
          if (img.width < 800 || img.height < 600) {
            toast.warning(
              `Изображение ${index + 1} имеет низкое качество. Рекомендуем загрузить фото минимум 800x600 пикселей`
            );
          }
        };
      });

      // Сохраняем и файлы, и превью
      setImageFiles([...imageFiles, ...newFiles].slice(0, 8));
      setFormData({
        ...formData,
        images: [...formData.images, ...newImages].slice(0, 8),
      });
      if (newImages.length > 0 && currentView === 'start') {
        setCurrentView('create');
      }
    }
  };

  const removeImage = (index: number) => {
    setImageFiles(imageFiles.filter((_, i) => i !== index));
    setFormData({
      ...formData,
      images: formData.images.filter((_, i) => i !== index),
    });
  };

  // Drag & Drop handlers
  const handleDragStart = (e: React.DragEvent, index: number) => {
    setDraggedIndex(index);
    e.dataTransfer.effectAllowed = 'move';
  };

  const handleDragOver = (e: React.DragEvent) => {
    e.preventDefault();
    e.dataTransfer.dropEffect = 'move';
  };

  const handleDrop = (e: React.DragEvent, dropIndex: number) => {
    e.preventDefault();
    if (draggedIndex === null) return;

    const draggedImage = formData.images[draggedIndex];
    const newImages = [...formData.images];

    // Remove the dragged image
    newImages.splice(draggedIndex, 1);

    // Insert it at the new position
    newImages.splice(dropIndex, 0, draggedImage);

    setFormData({ ...formData, images: newImages });
    setDraggedIndex(null);
  };

  const applyDescriptionTemplate = () => {
    if (formData.category && descriptionTemplates[formData.category]) {
      setFormData({
        ...formData,
        description: descriptionTemplates[formData.category],
      });
      toast.success('Шаблон применен! Отредактируйте детали');
    }
  };

  const handleAttributeChange = (
    attributeId: number,
    value: AttributeValue
  ) => {
    setFormData((prev) => ({
      ...prev,
      attributes: {
        ...prev.attributes,
        [attributeId]: value,
      },
    }));
  };

  const isCarCategory = (category: MarketplaceCategory | null): boolean => {
    if (!category) return false;
    // Проверяем slug категории и родительской категории
    const carSlugs = ['automotive', 'cars', 'licni-automobili'];
    return (
      carSlugs.includes(category.slug || '') ||
      (!!category.parent_id &&
        categories.some(
          (c) => c.id === category.parent_id && carSlugs.includes(c.slug || '')
        ))
    );
  };

  const handleCategorySelect = (category: MarketplaceCategory) => {
    setSelectedCategory(category);
    setCategorySelectedManually(true); // Помечаем, что категория выбрана вручную
    setFormData((prev) => ({
      ...prev,
      category: category.slug || '',
      categoryId: category.id || 0,
      // Очищаем атрибуты при смене категории
      attributes: {},
    }));
    // Сбрасываем выбор автомобиля при смене категории
    setCarSelection({});

    toast.success(
      `Выбрана категория: ${category.translations?.name || category.name}`
    );
  };

  const handlePublish = async () => {
    if (isSubmitting) return;

    // Валидация
    if (!formData.title || !formData.price || !formData.categoryId) {
      toast.error('Заполните все обязательные поля');
      return;
    }

    if (imageFiles.length === 0) {
      toast.error('Добавьте хотя бы одно фото');
      return;
    }

    setIsSubmitting(true);

    try {
      // Подготовка данных для создания объявления
      const listingData: CreateListingState = {
        title: formData.title,
        description: formData.description || '',
        price: parseFloat(formData.price) || 0,
        currency: 'RSD',
        category: selectedCategory
          ? {
              id: selectedCategory.id || 0,
              name: selectedCategory.name || '',
              slug: selectedCategory.slug || '',
            }
          : undefined,
        condition: 'used',
        location: {
          address: formData.location || '',
          city: formData.city,
          region: '',
          country: formData.country,
          latitude: 0,
          longitude: 0,
        },
        attributes: formData.attributes,
        images: [],
        mainImageIndex: 0,
        payment: {
          methods: formData.deliveryMethods,
          codEnabled: false,
          codPrice: 0,
          personalMeeting: false,
          negotiablePrice: false,
          bundleDeals: false,
          deliveryOptions: ['pickup'],
        },
        trust: {
          phoneVerified: false,
          preferredMeetingType: 'personal',
          meetingLocations: [],
          availableHours: '',
          localReputation: 0,
        },
        localization: {
          script: 'cyrillic',
          language: 'sr',
          traditionalUnits: false,
          regionalPhrases: [],
        },
        pijaca: {
          vendorStallStyle: '',
          regularCustomers: false,
          traditionalStyle: false,
        },
        isPublished: true,
        isDraft: false,
        originalLanguage: locale,
        translations: {},
      };

      // Создаем объявление
      const createResponse = await ListingsService.createListing(listingData);

      if (!createResponse.data?.id) {
        throw new Error('Не получен ID созданного объявления');
      }

      const listingId = createResponse.data.id;
      console.log('Объявление создано с ID:', listingId);

      // Загружаем изображения
      if (imageFiles.length > 0) {
        try {
          await ListingsService.uploadImages(listingId, imageFiles, 0);
          console.log('Изображения загружены успешно');
        } catch (uploadError) {
          console.error('Ошибка загрузки изображений:', uploadError);
          toast.warning(
            'Объявление создано, но не удалось загрузить изображения'
          );
        }
      }

      toast.success('Объявление опубликовано успешно!');
      router.push(`/${locale}/profile/listings`);
    } catch (error) {
      console.error('Ошибка при создании объявления:', error);
      toast.error(
        error instanceof Error
          ? error.message
          : 'Произошла ошибка при создании объявления'
      );
    } finally {
      setIsSubmitting(false);
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
            Шаблоны, сравнение цен, умные подсказки — всё для успешной продажи
          </p>

          {/* Stats */}
          <div className="flex justify-center gap-8 mb-8">
            <motion.div
              initial={{ scale: 0 }}
              animate={{ scale: 1 }}
              transition={{ delay: 0.2 }}
              className="text-center"
            >
              <div className="text-3xl font-bold text-primary">2-4 мин</div>
              <div className="text-sm text-base-content/60">создание</div>
            </motion.div>
            <motion.div
              initial={{ scale: 0 }}
              animate={{ scale: 1 }}
              transition={{ delay: 0.3 }}
              className="text-center"
            >
              <div className="text-3xl font-bold text-success">85%</div>
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
              onDragOver={(e) => {
                e.preventDefault();
                setIsDragging(true);
              }}
              onDragLeave={() => setIsDragging(false)}
              onDrop={(e) => {
                e.preventDefault();
                setIsDragging(false);
                const files = e.dataTransfer.files;
                if (files && files.length > 0) {
                  const newFiles = Array.from(files);
                  const newImages = newFiles.map((file) =>
                    URL.createObjectURL(file)
                  );
                  setImageFiles([...imageFiles, ...newFiles].slice(0, 8));
                  setFormData({
                    ...formData,
                    images: [...formData.images, ...newImages].slice(0, 8),
                  });
                  if (newImages.length > 0) {
                    setCurrentView('create');
                  }
                }
              }}
            >
              <div
                className={`card-body text-center py-12 ${isDragging ? 'opacity-70' : ''}`}
              >
                <Camera className="w-16 h-16 mx-auto mb-4" />
                <h2 className="text-2xl font-bold mb-2">
                  {isDragging ? 'Отпустите фото здесь' : 'Начните с фото'}
                </h2>
                <p className="opacity-90 mb-4">
                  Загрузите или перетащите фото товара
                </p>
                <div className="flex gap-2 justify-center">
                  <div className="badge badge-lg badge-warning gap-2">
                    <Zap className="w-4 h-4" />
                    Быстрый старт
                  </div>
                  <div className="badge badge-lg badge-info gap-2">
                    <FileText className="w-4 h-4" />
                    Умные шаблоны
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
          {/* Photo Upload Section with Drag & Drop */}
          <div className="card bg-base-200">
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
                    className="relative aspect-square group cursor-move"
                    draggable
                    onDragStart={(e: any) => handleDragStart(e, index)}
                    onDragOver={(e: any) => handleDragOver(e)}
                    onDrop={(e: any) => handleDrop(e, index)}
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
                  <label
                    className="aspect-square border-2 border-dashed border-base-300 rounded-lg flex flex-col items-center justify-center cursor-pointer hover:border-primary transition-colors"
                    onDragOver={(e) => {
                      e.preventDefault();
                      e.currentTarget.classList.add('border-primary');
                    }}
                    onDragLeave={(e) => {
                      e.currentTarget.classList.remove('border-primary');
                    }}
                    onDrop={(e) => {
                      e.preventDefault();
                      e.currentTarget.classList.remove('border-primary');
                      const files = e.dataTransfer.files;
                      if (files && files.length > 0) {
                        const newFiles = Array.from(files);
                        const newImages = newFiles.map((file) =>
                          URL.createObjectURL(file)
                        );
                        setImageFiles([...imageFiles, ...newFiles].slice(0, 8));
                        setFormData({
                          ...formData,
                          images: [...formData.images, ...newImages].slice(
                            0,
                            8
                          ),
                        });
                      }
                    }}
                  >
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

              {/* Совет по порядку фото */}
              {formData.images.length > 1 && (
                <div className="alert alert-info mt-4">
                  <Info className="w-4 h-4" />
                  <span className="text-sm">
                    Перетащите фото, чтобы изменить их порядок. Первое фото -
                    главное
                  </span>
                </div>
              )}
            </div>
          </div>

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
                  onChange={(e) => {
                    const newTitle = e.target.value;
                    setFormData({ ...formData, title: newTitle });
                  }}
                  maxLength={80}
                />

                {/* Предложенные категории */}
                {suggestedCategories.length > 0 && (
                  <div className="mt-3 p-3 bg-base-100 rounded-lg border border-base-300">
                    <p className="text-sm font-semibold mb-2 flex items-center gap-2">
                      <Sparkles className="w-4 h-4 text-primary" />
                      {categorySelectedManually
                        ? 'Альтернативные категории (на основе названия):'
                        : 'Предложенные категории (на основе названия):'}
                    </p>
                    <div className="flex flex-wrap gap-2">
                      {suggestedCategories.map((cat, index) => {
                        const category = categories.find(
                          (c) => c.id === cat.categoryId
                        );
                        const isSelected =
                          selectedCategory?.id === cat.categoryId;
                        const confidence = Math.round(
                          cat.confidenceScore * 100
                        );

                        return (
                          <button
                            key={cat.categoryId}
                            onClick={() => {
                              if (category) {
                                setSelectedCategory(category);
                                setCategorySelectedManually(true); // Помечаем, что категория выбрана вручную
                                setFormData((prev) => ({
                                  ...prev,
                                  category: category.slug || '',
                                  categoryId: category.id || 0,
                                }));
                                toast.success(
                                  `Выбрана категория: ${category.translations?.name || category.name}`
                                );
                              }
                            }}
                            className={`btn btn-sm ${
                              isSelected
                                ? 'btn-primary'
                                : index === 0
                                  ? 'btn-outline btn-primary'
                                  : 'btn-outline'
                            }`}
                          >
                            {category?.translations?.name || cat.categoryName}
                            <span className="badge badge-xs ml-1">
                              {confidence}%
                            </span>
                          </button>
                        );
                      })}
                    </div>
                    <p className="text-xs text-base-content/60 mt-2">
                      Нажмите на категорию, чтобы выбрать её
                    </p>
                  </div>
                )}
              </div>

              {/* Category Selector */}
              {!quickMode && (
                <CategorySelector
                  categories={categories}
                  selectedCategory={selectedCategory}
                  onCategorySelect={handleCategorySelect}
                  locale={locale}
                  compact={true}
                />
              )}

              {/* Price with comparison */}
              <div className="form-control">
                <label className="label">
                  <span className="label-text font-semibold">Цена</span>
                  <button
                    onClick={() => {
                      console.log(
                        'Price comparison button clicked, current state:',
                        showPriceComparison
                      );
                      setShowPriceComparison(!showPriceComparison);
                      // Если включаем сравнение - сразу загружаем похожие
                      if (
                        !showPriceComparison &&
                        formData.title &&
                        formData.categoryId
                      ) {
                        fetchSimilarListings(
                          formData.categoryId,
                          formData.title
                        );
                      }
                    }}
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
                {showPriceComparison && (
                  <div className="mt-4 space-y-2">
                    <h4 className="text-sm font-semibold flex items-center gap-2">
                      Похожие объявления:
                      {isLoadingSimilar && (
                        <Loader2 className="w-3 h-3 animate-spin" />
                      )}
                    </h4>
                    {!isLoadingSimilar && similarListings.length > 0 ? (
                      similarListings.map((listing) => (
                        <div
                          key={listing.id}
                          className="flex items-center justify-between text-sm p-2 bg-base-100 rounded hover:bg-base-100/70 transition-colors cursor-pointer"
                          onClick={() =>
                            window.open(
                              `/${locale}/marketplace/${listing.id}`,
                              '_blank'
                            )
                          }
                        >
                          <div className="flex-1">
                            <p className="font-medium">{listing.title}</p>
                            <p className="text-xs text-base-content/60">
                              <Eye className="w-3 h-3 inline mr-1" />
                              {listing.views} просмотров •
                              {listing.daysAgo === 0
                                ? 'сегодня'
                                : listing.daysAgo === 1
                                  ? 'вчера'
                                  : `${listing.daysAgo} дн. назад`}
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
                      ))
                    ) : !isLoadingSimilar && formData.title.length > 3 ? (
                      <p className="text-xs text-base-content/60">
                        Похожих объявлений не найдено. Ваша цена может быть
                        уникальной!
                      </p>
                    ) : (
                      <p className="text-xs text-base-content/60">
                        Введите название товара для поиска похожих
                      </p>
                    )}
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
                      Использовать шаблон для{' '}
                      {formData.category === 'fashion'
                        ? 'одежды/обуви'
                        : formData.category === 'electronics'
                          ? 'электроники'
                          : formData.category === 'home'
                            ? 'дома'
                            : formData.category === 'auto'
                              ? 'авто'
                              : 'товара'}
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

          {/* Car Selector for automotive categories */}
          {selectedCategory && isCarCategory(selectedCategory) && (
            <div className="card bg-base-200 mb-6">
              <div className="card-body">
                <h3 className="card-title text-base mb-4">
                  <Car className="w-5 h-5" />
                  Выберите автомобиль
                </h3>
                <CarSelectorCompact
                  value={carSelection}
                  onChange={(selection) => {
                    setCarSelection(selection);
                    // Добавляем марку и модель только как атрибуты, НЕ в название
                    if (selection.make && selection.model) {
                      setFormData((prev) => ({
                        ...prev,
                        // Добавляем марку и модель как атрибуты
                        attributes: {
                          ...prev.attributes,
                          // Марка
                          make: {
                            attribute_id: 0,
                            attribute_name: 'make',
                            display_name: 'Марка',
                            attribute_type: 'text',
                            text_value: selection.make?.name || '',
                          },
                          // Модель
                          model: {
                            attribute_id: 0,
                            attribute_name: 'model',
                            display_name: 'Модель',
                            attribute_type: 'text',
                            text_value: selection.model?.name || '',
                          },
                        },
                      }));
                    }
                  }}
                />
              </div>
            </div>
          )}

          {/* Category Attributes */}
          {selectedCategory && (
            <CategoryAttributes
              selectedCategory={selectedCategory}
              attributes={formData.attributes}
              onAttributeChange={handleAttributeChange}
              locale={locale}
            />
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
              {Object.keys(formData.attributes).length > 0 && (
                <div className="grid grid-cols-2 gap-4 mb-4">
                  {Object.values(formData.attributes).map((attr) => {
                    const value =
                      attr.text_value ||
                      attr.numeric_value ||
                      (attr.boolean_value ? 'Да' : 'Нет');
                    if (!value) return null;

                    return (
                      <div
                        key={attr.attribute_id}
                        className="flex justify-between py-2 border-b border-base-200"
                      >
                        <span className="text-sm text-base-content/60">
                          {attr.display_name}
                        </span>
                        <span className="text-sm font-medium">
                          {value}
                          {attr.unit && ` ${attr.unit}`}
                        </span>
                      </div>
                    );
                  })}
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
                  Умная оптимизация увеличит охват
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
            <button
              onClick={handlePublish}
              disabled={isSubmitting}
              className="btn btn-primary btn-lg flex-1"
            >
              {isSubmitting ? (
                <>
                  <Loader2 className="w-5 h-5 animate-spin" />
                  Публикация...
                </>
              ) : (
                <>
                  Опубликовать сейчас
                  <Sparkles className="w-5 h-5 ml-1" />
                </>
              )}
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
                <span className="font-semibold">1,234</span> продавцов уже
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
            href={`/${locale}/create-listing-choice`}
            className="btn btn-ghost"
          >
            <ChevronLeft className="w-5 h-5" />
            Назад к выбору
          </Link>
        </div>
        <div className="flex-none">
          <div className="badge badge-success badge-lg">Бесплатная версия</div>
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
