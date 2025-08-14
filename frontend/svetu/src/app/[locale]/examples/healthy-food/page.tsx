'use client';

import { useState } from 'react';
import { AnimatedSection } from '@/components/ui/AnimatedSection';

interface Product {
  id: string;
  name: string;
  category: string;
  description: string;
  benefits: string[];
  price: number;
  unit: string;
  image: string;
  organic: boolean;
  local: boolean;
  seasonal: boolean;
  nutrients: {
    calories?: number;
    protein?: number;
    carbs?: number;
    fiber?: number;
    vitamins?: string[];
    sugar?: number;
  };
  harmfulIngredients?: string[];
  warnings?: string[];
}

const healthyProducts: Product[] = [
  {
    id: '1',
    name: 'Авокадо',
    category: 'fruits',
    description: 'Богат полезными жирами и клетчаткой',
    benefits: [
      'Снижает холестерин',
      'Улучшает пищеварение',
      'Источник витамина E',
    ],
    price: 132,
    unit: 'шт',
    image: '🥑',
    organic: true,
    local: false,
    seasonal: false,
    nutrients: {
      calories: 160,
      protein: 2,
      carbs: 9,
      fiber: 7,
      vitamins: ['K', 'C', 'E', 'B5', 'B6'],
    },
  },
  {
    id: '2',
    name: 'Черника',
    category: 'berries',
    description: 'Антиоксиданты для здоровья мозга',
    benefits: ['Улучшает память', 'Защищает от старения', 'Снижает давление'],
    price: 275,
    unit: 'кг',
    image: '🫐',
    organic: true,
    local: true,
    seasonal: true,
    nutrients: {
      calories: 57,
      protein: 0.7,
      carbs: 14,
      fiber: 2.4,
      vitamins: ['C', 'K', 'B6'],
    },
  },
  {
    id: '3',
    name: 'Брокколи',
    category: 'vegetables',
    description: 'Суперфуд с витаминами и минералами',
    benefits: [
      'Укрепляет иммунитет',
      'Детоксикация организма',
      'Профилактика рака',
    ],
    price: 94,
    unit: 'кг',
    image: '🥦',
    organic: true,
    local: true,
    seasonal: false,
    nutrients: {
      calories: 34,
      protein: 2.8,
      carbs: 7,
      fiber: 2.6,
      vitamins: ['C', 'K', 'A', 'B9'],
    },
  },
  {
    id: '4',
    name: 'Киноа',
    category: 'grains',
    description: 'Полноценный белок без глютена',
    benefits: ['Все аминокислоты', 'Без глютена', 'Высокое содержание магния'],
    price: 352,
    unit: 'кг',
    image: '🌾',
    organic: true,
    local: false,
    seasonal: false,
    nutrients: {
      calories: 120,
      protein: 4.4,
      carbs: 21,
      fiber: 2.8,
      vitamins: ['B1', 'B2', 'B6', 'B9'],
    },
  },
  {
    id: '5',
    name: 'Грецкие орехи',
    category: 'nuts',
    description: 'Омега-3 для сердца и мозга',
    benefits: [
      'Здоровье сердца',
      'Улучшает работу мозга',
      'Снижает воспаление',
    ],
    price: 495,
    unit: 'кг',
    image: '🌰',
    organic: true,
    local: true,
    seasonal: true,
    nutrients: {
      calories: 654,
      protein: 15,
      carbs: 14,
      fiber: 7,
      vitamins: ['E', 'B6', 'B9'],
    },
  },
  {
    id: '6',
    name: 'Шпинат',
    category: 'vegetables',
    description: 'Железо и фолиевая кислота',
    benefits: ['Укрепляет кости', 'Улучшает зрение', 'Богат железом'],
    price: 105,
    unit: 'кг',
    image: '🥬',
    organic: true,
    local: true,
    seasonal: true,
    nutrients: {
      calories: 23,
      protein: 2.9,
      carbs: 3.6,
      fiber: 2.2,
      vitamins: ['K', 'A', 'C', 'B9'],
    },
  },
  {
    id: '7',
    name: 'Лосось дикий',
    category: 'fish',
    description: 'Омега-3 и витамин D',
    benefits: [
      'Здоровье сердца',
      'Противовоспалительный',
      'Улучшает настроение',
    ],
    price: 979,
    unit: 'кг',
    image: '🐟',
    organic: false,
    local: false,
    seasonal: false,
    nutrients: {
      calories: 208,
      protein: 20,
      carbs: 0,
      fiber: 0,
      vitamins: ['D', 'B12', 'B6'],
    },
  },
  {
    id: '8',
    name: 'Чечевица',
    category: 'legumes',
    description: 'Растительный белок и клетчатка',
    benefits: ['Стабилизирует сахар', 'Высокий белок', 'Улучшает пищеварение'],
    price: 132,
    unit: 'кг',
    image: '🥘',
    organic: true,
    local: false,
    seasonal: false,
    nutrients: {
      calories: 116,
      protein: 9,
      carbs: 20,
      fiber: 8,
      vitamins: ['B9', 'B1', 'B6'],
    },
  },
  {
    id: '9',
    name: 'Батат',
    category: 'vegetables',
    description: 'Бета-каротин и клетчатка',
    benefits: ['Улучшает зрение', 'Регулирует сахар', 'Антиоксиданты'],
    price: 83,
    unit: 'кг',
    image: '🍠',
    organic: true,
    local: true,
    seasonal: true,
    nutrients: {
      calories: 86,
      protein: 1.6,
      carbs: 20,
      fiber: 3,
      vitamins: ['A', 'C', 'B6'],
    },
  },
  {
    id: '10',
    name: 'Семена чиа',
    category: 'seeds',
    description: 'Омега-3 и антиоксиданты',
    benefits: ['Снижает вес', 'Укрепляет кости', 'Энергия на весь день'],
    price: 638,
    unit: 'кг',
    image: '🌱',
    organic: true,
    local: false,
    seasonal: false,
    nutrients: {
      calories: 486,
      protein: 17,
      carbs: 42,
      fiber: 34,
      vitamins: ['B1', 'B3'],
    },
  },
  {
    id: '11',
    name: 'Гранат',
    category: 'fruits',
    description: 'Антиоксиданты для сердца',
    benefits: ['Защищает сердце', 'Снижает давление', 'Противовоспалительный'],
    price: 198,
    unit: 'кг',
    image: '🍎',
    organic: true,
    local: false,
    seasonal: true,
    nutrients: {
      calories: 83,
      protein: 1.7,
      carbs: 19,
      fiber: 4,
      vitamins: ['C', 'K', 'B9'],
    },
  },
  {
    id: '12',
    name: 'Куркума',
    category: 'spices',
    description: 'Противовоспалительная специя',
    benefits: ['Снижает воспаление', 'Улучшает память', 'Антиоксидант'],
    price: 462,
    unit: 'кг',
    image: '🌶️',
    organic: true,
    local: false,
    seasonal: false,
    nutrients: {
      calories: 312,
      protein: 9.7,
      carbs: 67,
      fiber: 22,
      vitamins: ['C', 'B6'],
    },
  },
];

// Примеры "якобы здоровых" продуктов с вредными добавками
const fakeHealthyProducts: Product[] = [
  {
    id: 'fake1',
    name: 'Фитнес батончик "Энергия"',
    category: 'snacks',
    description: 'Злаковый батончик с витаминами',
    benefits: ['Добавлены витамины B', 'Содержит злаки'],
    price: 98,
    unit: 'шт',
    image: '🍫',
    organic: false,
    local: false,
    seasonal: false,
    nutrients: {
      calories: 450,
      protein: 5,
      carbs: 65,
      sugar: 42,
      fiber: 2,
    },
    harmfulIngredients: [
      'Кукурузный сироп',
      'Пальмовое масло',
      'Мальтодекстрин',
      'E621',
    ],
    warnings: ['42г сахара - это 10 ложек!', 'Больше калорий чем в Сникерсе'],
  },
  {
    id: 'fake2',
    name: 'Йогурт 0% жирности "Imlek"',
    category: 'dairy',
    description: 'Обезжиренный йогурт с фруктами',
    benefits: ['Без жира', 'С кусочками фруктов'],
    price: 72,
    unit: 'шт',
    image: '🥛',
    organic: false,
    local: false,
    seasonal: false,
    nutrients: {
      calories: 120,
      protein: 4,
      carbs: 24,
      sugar: 20,
      fiber: 0,
    },
    harmfulIngredients: [
      'Аспартам',
      'Ацесульфам К',
      'Модифицированный крахмал',
      'Ароматизаторы',
    ],
    warnings: [
      'Сахарозаменители могут вызывать диабет',
      'Без жира = больше сахара',
    ],
  },
  {
    id: 'fake3',
    name: 'Диетическая кола',
    category: 'drinks',
    description: 'Без калорий, без сахара',
    benefits: ['0 калорий', 'Без сахара'],
    price: 98,
    unit: 'л',
    image: '🥤',
    organic: false,
    local: false,
    seasonal: false,
    nutrients: {
      calories: 0,
      protein: 0,
      carbs: 0,
      sugar: 0,
      fiber: 0,
    },
    harmfulIngredients: [
      'Аспартам',
      'Ортофосфорная кислота',
      'Бензоат натрия',
      'Кофеин',
    ],
    warnings: [
      'Аспартам токсичен при нагревании',
      'Разрушает зубную эмаль',
      'Вымывает кальций из костей',
    ],
  },
  {
    id: 'fake4',
    name: 'Мюсли "Zlato Polje Fit"',
    category: 'grains',
    description: 'С сухофруктами и орехами',
    benefits: ['Цельные злаки', 'Витамины и минералы'],
    price: 352,
    unit: 'кг',
    image: '🥣',
    organic: false,
    local: false,
    seasonal: false,
    nutrients: {
      calories: 380,
      protein: 8,
      carbs: 68,
      sugar: 35,
      fiber: 5,
    },
    harmfulIngredients: [
      'Сахар',
      'Пальмовое масло',
      'Глюкозно-фруктозный сироп',
      'Консерванты E220',
    ],
    warnings: [
      '35г сахара на 100г - почти как в конфетах!',
      'Пальмовое масло закупоривает сосуды',
    ],
  },
  {
    id: 'fake5',
    name: 'Сок "100% натуральный"',
    category: 'drinks',
    description: 'Апельсиновый сок без добавок',
    benefits: ['100% сок', 'Витамин C'],
    price: 165,
    unit: 'л',
    image: '🧃',
    organic: false,
    local: false,
    seasonal: false,
    nutrients: {
      calories: 110,
      protein: 1.7,
      carbs: 26,
      sugar: 24,
      fiber: 0,
    },
    harmfulIngredients: [
      'Концентрированный сок',
      'Добавленный сахар',
      'Лимонная кислота',
    ],
    warnings: [
      'Стакан сока = 6 ложек сахара',
      'Без клетчатки - чистый сахар в кровь',
    ],
  },
  {
    id: 'fake6',
    name: 'Детски сир "Moja Kravica"',
    category: 'dairy',
    description: 'Для здорового роста',
    benefits: ['Кальций', 'Для детей'],
    price: 50,
    unit: 'шт',
    image: '🍮',
    organic: false,
    local: false,
    seasonal: false,
    nutrients: {
      calories: 140,
      protein: 5,
      carbs: 18,
      sugar: 15,
      fiber: 0,
    },
    harmfulIngredients: [
      'Сахар',
      'Модифицированный крахмал',
      'Ароматизаторы',
      'Красители',
    ],
    warnings: [
      '15г сахара в маленькой баночке',
      'Формирует сахарную зависимость у детей',
    ],
  },
  {
    id: 'fake7',
    name: 'Кракерси "ZdravoŽиво"',
    category: 'snacks',
    description: 'Диетические цельнозерновые',
    benefits: ['Мало калорий', 'Цельное зерно'],
    price: 98,
    unit: 'упак',
    image: '🍘',
    organic: false,
    local: false,
    seasonal: false,
    nutrients: {
      calories: 360,
      protein: 10,
      carbs: 70,
      sugar: 8,
      fiber: 12,
    },
    harmfulIngredients: [
      'Пальмовое масло',
      'Усилители вкуса',
      'Дрожжевой экстракт (скрытый глутамат)',
    ],
    warnings: ['Калорийнее хлеба!', 'Дрожжевой экстракт = глутамат натрия'],
  },
  {
    id: 'fake8',
    name: 'Брза овсена каша "Nestle"',
    category: 'grains',
    description: 'Овсянка с фруктами за 1 минуту',
    benefits: ['Быстро готовится', 'С фруктами'],
    price: 28,
    unit: 'порция',
    image: '🥘',
    organic: false,
    local: false,
    seasonal: false,
    nutrients: {
      calories: 350,
      protein: 8,
      carbs: 65,
      sugar: 28,
      fiber: 4,
    },
    harmfulIngredients: [
      'Сахар',
      'Мальтодекстрин',
      'Ароматизаторы',
      'Пальмовое масло',
    ],
    warnings: [
      '28г сахара в одной порции',
      'Мальтодекстрин повышает сахар в крови быстрее сахара',
    ],
  },
  {
    id: 'fake9',
    name: 'Соевое "мясо"',
    category: 'proteins',
    description: 'Растительный белок',
    benefits: ['Без холестерина', 'Высокий белок'],
    price: 308,
    unit: 'кг',
    image: '🍖',
    organic: false,
    local: false,
    seasonal: false,
    nutrients: {
      calories: 290,
      protein: 52,
      carbs: 30,
      sugar: 5,
      fiber: 4,
    },
    harmfulIngredients: [
      'Глутамат натрия',
      'Трансжиры',
      'Красители',
      'Консерванты',
    ],
    warnings: ['ГМО соя', 'Глутамат вызывает переедание'],
  },
  {
    id: 'fake10',
    name: 'Смузи "Детокс"',
    category: 'drinks',
    description: 'Очищающий напиток',
    benefits: ['Детокс эффект', 'Витамины'],
    price: 275,
    unit: 'бутылка',
    image: '🥤',
    organic: false,
    local: false,
    seasonal: false,
    nutrients: {
      calories: 180,
      protein: 2,
      carbs: 42,
      sugar: 38,
      fiber: 1,
    },
    harmfulIngredients: ['Концентрат сока', 'Сахарный сироп', 'Консерванты'],
    warnings: ['38г сахара = 9 ложек!', 'Никакого детокса, только сахар'],
  },
  {
    id: 'fake11',
    name: 'Протеински штапић "Protein Plus"',
    category: 'snacks',
    description: 'Для спортсменов',
    benefits: ['25г белка', 'Для мышц'],
    price: 165,
    unit: 'шт',
    image: '🍫',
    organic: false,
    local: false,
    seasonal: false,
    nutrients: {
      calories: 380,
      protein: 25,
      carbs: 35,
      sugar: 22,
      fiber: 2,
    },
    harmfulIngredients: [
      'Пальмовое масло',
      'Сукралоза',
      'Мальтитол',
      'Ароматизаторы',
    ],
    warnings: ['Сукралоза убивает кишечную флору', 'Мальтитол вызывает диарею'],
  },
  {
    id: 'fake12',
    name: 'Хрскави кракерси "Carnex"',
    category: 'snacks',
    description: 'Запеченные, не жареные',
    benefits: ['Без масла', 'Запеченные'],
    price: 50,
    unit: 'пачка',
    image: '🥖',
    organic: false,
    local: false,
    seasonal: false,
    nutrients: {
      calories: 410,
      protein: 11,
      carbs: 72,
      sugar: 6,
      fiber: 3,
    },
    harmfulIngredients: [
      'Глутамат натрия',
      'Усилители вкуса',
      'Трансжиры',
      'Акриламид',
    ],
    warnings: ['Акриламид - канцероген', 'Глутамат вызывает зависимость'],
  },
];

// Перемешиваем продукты в случайном порядке
const shuffleArray = (array: Product[]) => {
  const shuffled = [...array];
  for (let i = shuffled.length - 1; i > 0; i--) {
    const j = Math.floor(Math.random() * (i + 1));
    [shuffled[i], shuffled[j]] = [shuffled[j], shuffled[i]];
  }
  return shuffled;
};

const allProducts = shuffleArray([...healthyProducts, ...fakeHealthyProducts]);

const categories = [
  { id: 'all', name: 'Все продукты', icon: '🌿' },
  { id: 'vegetables', name: 'Овощи', icon: '🥬' },
  { id: 'fruits', name: 'Фрукты', icon: '🍎' },
  { id: 'berries', name: 'Ягоды', icon: '🫐' },
  { id: 'nuts', name: 'Орехи', icon: '🌰' },
  { id: 'seeds', name: 'Семена', icon: '🌱' },
  { id: 'grains', name: 'Злаки', icon: '🌾' },
  { id: 'legumes', name: 'Бобовые', icon: '🥘' },
  { id: 'fish', name: 'Рыба', icon: '🐟' },
  { id: 'spices', name: 'Специи', icon: '🌶️' },
  { id: 'snacks', name: 'Снеки', icon: '🍫' },
  { id: 'dairy', name: 'Молочные', icon: '🥛' },
  { id: 'drinks', name: 'Напитки', icon: '🥤' },
  { id: 'proteins', name: 'Белковые', icon: '🍖' },
];

interface UserProfile {
  age: number;
  weight: number;
  height: number;
  gender: 'male' | 'female';
  activity: 'sedentary' | 'light' | 'moderate' | 'active' | 'very_active';
  goal: 'lose' | 'maintain' | 'gain';
  restrictions?: string[];
}

interface Supermarket {
  name: string;
  logo: string;
  color: string;
  priceMultiplier: number;
}

const supermarkets: Supermarket[] = [
  { name: 'Maxi', logo: '🟢', color: 'bg-green-100', priceMultiplier: 1.0 },
  { name: 'Idea', logo: '🟡', color: 'bg-yellow-100', priceMultiplier: 0.85 },
  { name: 'Tempo', logo: '🔵', color: 'bg-blue-100', priceMultiplier: 0.9 },
  { name: 'Roda', logo: '🌿', color: 'bg-emerald-100', priceMultiplier: 1.5 },
  { name: 'Aman', logo: '🛒', color: 'bg-orange-100', priceMultiplier: 0.8 },
  {
    name: 'Univerexport',
    logo: 'Ⓜ️',
    color: 'bg-indigo-100',
    priceMultiplier: 0.95,
  },
];

export default function HealthyFoodPage() {
  const [selectedCategory, setSelectedCategory] = useState('all');
  const [showOrganic, setShowOrganic] = useState(false);
  const [showLocal, setShowLocal] = useState(false);
  const [showSeasonal, setShowSeasonal] = useState(false);
  const [showOnlyHealthy, setShowOnlyHealthy] = useState(false);
  const [showNoSugar, setShowNoSugar] = useState(false);
  const [showNoPalmOil, setShowNoPalmOil] = useState(false);
  const [showNoAspartame, setShowNoAspartame] = useState(false);
  const [showNoTransFats, setShowNoTransFats] = useState(false);
  const [showNoMSG, setShowNoMSG] = useState(false);
  const [userProfile, setUserProfile] = useState<UserProfile>({
    age: 30,
    weight: 70,
    height: 170,
    gender: 'male',
    activity: 'moderate',
    goal: 'maintain',
    restrictions: [],
  });
  const [showProfileForm, setShowProfileForm] = useState(false);
  const [generatedBasket, setGeneratedBasket] = useState<any>(null);
  const [selectedPeriod, setSelectedPeriod] = useState<'day' | 'week'>('day');

  // Подсчет активных фильтров вредных добавок
  const activeHarmfulFilters = [
    showNoSugar,
    showNoPalmOil,
    showNoAspartame,
    showNoTransFats,
    showNoMSG,
  ].filter(Boolean).length;

  const filteredProducts = allProducts.filter((product) => {
    if (selectedCategory !== 'all' && product.category !== selectedCategory) {
      return false;
    }
    if (showOrganic && !product.organic) return false;
    if (showLocal && !product.local) return false;
    if (showSeasonal && !product.seasonal) return false;

    // Фильтр только здоровых продуктов
    if (
      showOnlyHealthy &&
      product.harmfulIngredients &&
      product.harmfulIngredients.length > 0
    ) {
      return false;
    }

    // Фильтры вредных добавок
    if (showNoSugar && product.nutrients.sugar && product.nutrients.sugar > 5)
      return false;
    if (
      showNoPalmOil &&
      product.harmfulIngredients?.includes('Пальмовое масло')
    )
      return false;
    if (
      showNoAspartame &&
      product.harmfulIngredients?.some(
        (ing) => ing.includes('Аспартам') || ing.includes('Ацесульфам')
      )
    )
      return false;
    if (
      showNoTransFats &&
      product.harmfulIngredients?.some(
        (ing) =>
          ing.includes('Трансжиры') || ing.includes('Гидрогенизированное')
      )
    )
      return false;
    if (
      showNoMSG &&
      product.harmfulIngredients?.some(
        (ing) => ing.includes('E621') || ing.includes('Глутамат')
      )
    )
      return false;

    return true;
  });

  // Расчет базового метаболизма (BMR) по формуле Миффлина-Сан Жеора
  const calculateBMR = () => {
    const { weight, height, age, gender } = userProfile;
    let bmr = 10 * weight + 6.25 * height - 5 * age;
    bmr = gender === 'male' ? bmr + 5 : bmr - 161;
    return bmr;
  };

  // Расчет дневной нормы калорий с учетом активности
  const calculateDailyCalories = () => {
    const bmr = calculateBMR();
    const activityMultipliers = {
      sedentary: 1.2,
      light: 1.375,
      moderate: 1.55,
      active: 1.725,
      very_active: 1.9,
    };

    let calories = bmr * activityMultipliers[userProfile.activity];

    // Корректировка для цели
    if (userProfile.goal === 'lose') calories *= 0.85; // Дефицит 15%
    if (userProfile.goal === 'gain') calories *= 1.15; // Профицит 15%

    return Math.round(calories);
  };

  // Расчет макронутриентов
  const calculateMacros = () => {
    const calories = calculateDailyCalories();
    return {
      calories,
      protein: Math.round((calories * 0.25) / 4), // 25% от калорий, 4 ккал/г
      carbs: Math.round((calories * 0.45) / 4), // 45% от калорий, 4 ккал/г
      fats: Math.round((calories * 0.3) / 9), // 30% от калорий, 9 ккал/г
      fiber: userProfile.gender === 'male' ? 38 : 25, // г/день
    };
  };

  // Генерация оптимальной корзины
  const generateOptimalBasket = () => {
    const macros = calculateMacros();
    const days = selectedPeriod === 'week' ? 7 : 1;

    // Логика подбора продуктов для достижения целевых макросов
    const basket: any = {
      period: selectedPeriod,
      targetMacros: {
        ...macros,
        calories: macros.calories * days,
        protein: macros.protein * days,
        carbs: macros.carbs * days,
        fats: macros.fats * days,
        fiber: macros.fiber * days,
      },
      items: [] as any[],
      totals: {
        calories: 0,
        protein: 0,
        carbs: 0,
        fats: 0,
        fiber: 0,
        price: 0,
      },
      supermarketPrices: [] as any[],
    };

    // Упрощенный алгоритм подбора продуктов
    const categories = [
      'vegetables',
      'fruits',
      'grains',
      'legumes',
      'nuts',
      'fish',
    ];
    const selectedProducts: any[] = [];

    categories.forEach((category) => {
      const categoryProducts = healthyProducts.filter(
        (p) => p.category === category
      );
      if (categoryProducts.length > 0) {
        // Выбираем 1-2 продукта из каждой категории
        const selected = categoryProducts.slice(0, 2);
        selected.forEach((product) => {
          const quantity =
            category === 'vegetables' || category === 'fruits'
              ? 300 * days
              : 150 * days;
          selectedProducts.push({
            ...product,
            quantity: quantity,
            totalPrice: (product.price * quantity) / 1000,
          });
        });
      }
    });

    basket.items = selectedProducts;

    // Подсчет итогов
    selectedProducts.forEach((item) => {
      const multiplier = item.quantity / 100; // Нутриенты даны на 100г
      basket.totals.calories += (item.nutrients.calories || 0) * multiplier;
      basket.totals.protein += (item.nutrients.protein || 0) * multiplier;
      basket.totals.carbs += (item.nutrients.carbs || 0) * multiplier;
      basket.totals.fiber += (item.nutrients.fiber || 0) * multiplier;
      basket.totals.price += item.totalPrice;
    });

    // Расчет цен по супермаркетам
    basket.supermarketPrices = supermarkets.map((market) => ({
      ...market,
      totalPrice: Math.round(basket.totals.price * market.priceMultiplier),
      savings: Math.round(basket.totals.price * (1 - market.priceMultiplier)),
    }));

    setGeneratedBasket(basket);
  };

  // Экспорт списка для сборщика
  const exportShoppingList = () => {
    if (!generatedBasket) return;

    let shoppingList = `📱 СПИСОК ПОКУПОК НА ${selectedPeriod === 'day' ? 'ДЕНЬ' : 'НЕДЕЛЮ'}\n`;
    shoppingList += `👤 Для: ${userProfile.gender === 'male' ? 'Мужчина' : 'Женщина'}, ${userProfile.age} лет, ${userProfile.weight}кг\n`;
    shoppingList += `🎯 Цель: ${userProfile.goal === 'lose' ? 'Похудение' : userProfile.goal === 'maintain' ? 'Поддержание' : 'Набор массы'}\n`;
    shoppingList += `━━━━━━━━━━━━━━━━━━━━━━━━\n\n`;

    generatedBasket.items.forEach((item: any, index: number) => {
      shoppingList += `${index + 1}. ${item.name}\n`;
      shoppingList += `   Количество: ${item.quantity}г\n`;
      shoppingList += `   Примерная цена: ${Math.round(item.totalPrice)} RSD\n\n`;
    });

    shoppingList += `━━━━━━━━━━━━━━━━━━━━━━━━\n`;
    shoppingList += `💰 ИТОГО: ~${Math.round(generatedBasket.totals.price)} RSD\n\n`;
    shoppingList += `📊 ПИЩЕВАЯ ЦЕННОСТЬ:\n`;
    shoppingList += `• Калории: ${Math.round(generatedBasket.totals.calories)} ккал\n`;
    shoppingList += `• Белки: ${Math.round(generatedBasket.totals.protein)}г\n`;
    shoppingList += `• Углеводы: ${Math.round(generatedBasket.totals.carbs)}г\n`;
    shoppingList += `• Клетчатка: ${Math.round(generatedBasket.totals.fiber)}г\n\n`;
    shoppingList += `⚠️ ВАЖНО: Избегайте продуктов с:\n`;
    shoppingList += `• Добавленным сахаром\n`;
    shoppingList += `• Пальмовым маслом\n`;
    shoppingList += `• Аспартамом\n`;
    shoppingList += `• Трансжирами\n`;
    shoppingList += `• Глутаматом натрия (E621)\n`;

    // Создаем и скачиваем файл
    const blob = new Blob([shoppingList], { type: 'text/plain;charset=utf-8' });
    const url = window.URL.createObjectURL(blob);
    const a = document.createElement('a');
    a.href = url;
    a.download = `shopping-list-${new Date().toISOString().split('T')[0]}.txt`;
    document.body.appendChild(a);
    a.click();
    document.body.removeChild(a);
    window.URL.revokeObjectURL(url);

    // Копируем в буфер обмена
    navigator.clipboard.writeText(shoppingList).then(() => {
      alert('Список скопирован в буфер обмена и сохранен в файл!');
    });
  };

  return (
    <div className="min-h-screen bg-gradient-to-br from-green-50 via-white to-emerald-50">
      {/* Hero Section */}
      <AnimatedSection animation="fadeIn">
        <div className="hero min-h-[40vh] bg-gradient-to-r from-green-600 to-emerald-600">
          <div className="hero-content text-center text-white">
            <div className="max-w-4xl">
              <h1 className="text-5xl font-bold mb-6">
                🌱 Настоящее Здоровое Питание
              </h1>
              <p className="text-xl mb-4">
                Только натуральные продукты без маркетинговых уловок
              </p>
              <p className="text-lg opacity-90">
                Никаких &ldquo;диетических&rdquo; колы, &ldquo;полезных&rdquo;
                чипсов или &ldquo;витаминизированных&rdquo; конфет. Только то,
                что действительно питает ваш организм.
              </p>
            </div>
          </div>
        </div>
      </AnimatedSection>

      <div className="container mx-auto px-4 py-8 max-w-7xl">
        {/* Personal Profile Section */}
        <AnimatedSection animation="slideUp" delay={0.1}>
          <div className="card bg-gradient-to-r from-primary/10 to-secondary/10 shadow-xl mb-8">
            <div className="card-body">
              <div className="flex justify-between items-start">
                <div>
                  <h2 className="card-title text-2xl mb-4">
                    🎯 Персональный план питания
                  </h2>
                  {!showProfileForm && (
                    <div className="grid grid-cols-2 md:grid-cols-3 gap-4 text-sm">
                      <div>
                        <span className="font-semibold">Возраст:</span>{' '}
                        {userProfile.age} лет
                      </div>
                      <div>
                        <span className="font-semibold">Вес:</span>{' '}
                        {userProfile.weight} кг
                      </div>
                      <div>
                        <span className="font-semibold">Рост:</span>{' '}
                        {userProfile.height} см
                      </div>
                      <div>
                        <span className="font-semibold">Пол:</span>{' '}
                        {userProfile.gender === 'male' ? 'Мужской' : 'Женский'}
                      </div>
                      <div>
                        <span className="font-semibold">Активность:</span>{' '}
                        {
                          {
                            sedentary: 'Сидячий',
                            light: 'Легкая',
                            moderate: 'Умеренная',
                            active: 'Активная',
                            very_active: 'Очень активная',
                          }[userProfile.activity]
                        }
                      </div>
                      <div>
                        <span className="font-semibold">Цель:</span>{' '}
                        {
                          {
                            lose: 'Похудение',
                            maintain: 'Поддержание',
                            gain: 'Набор массы',
                          }[userProfile.goal]
                        }
                      </div>
                    </div>
                  )}
                </div>
                <button
                  onClick={() => setShowProfileForm(!showProfileForm)}
                  className="btn btn-sm btn-primary"
                >
                  {showProfileForm ? 'Закрыть' : 'Изменить'}
                </button>
              </div>

              {showProfileForm && (
                <div className="mt-4 grid grid-cols-1 md:grid-cols-3 gap-4">
                  <div className="form-control">
                    <label className="label">
                      <span className="label-text">Возраст</span>
                    </label>
                    <input
                      type="number"
                      className="input input-bordered"
                      value={userProfile.age}
                      onChange={(e) =>
                        setUserProfile({
                          ...userProfile,
                          age: parseInt(e.target.value),
                        })
                      }
                      min="10"
                      max="100"
                    />
                  </div>
                  <div className="form-control">
                    <label className="label">
                      <span className="label-text">Вес (кг)</span>
                    </label>
                    <input
                      type="number"
                      className="input input-bordered"
                      value={userProfile.weight}
                      onChange={(e) =>
                        setUserProfile({
                          ...userProfile,
                          weight: parseInt(e.target.value),
                        })
                      }
                      min="30"
                      max="200"
                    />
                  </div>
                  <div className="form-control">
                    <label className="label">
                      <span className="label-text">Рост (см)</span>
                    </label>
                    <input
                      type="number"
                      className="input input-bordered"
                      value={userProfile.height}
                      onChange={(e) =>
                        setUserProfile({
                          ...userProfile,
                          height: parseInt(e.target.value),
                        })
                      }
                      min="100"
                      max="250"
                    />
                  </div>
                  <div className="form-control">
                    <label className="label">
                      <span className="label-text">Пол</span>
                    </label>
                    <select
                      className="select select-bordered"
                      value={userProfile.gender}
                      onChange={(e) =>
                        setUserProfile({
                          ...userProfile,
                          gender: e.target.value as 'male' | 'female',
                        })
                      }
                    >
                      <option value="male">Мужской</option>
                      <option value="female">Женский</option>
                    </select>
                  </div>
                  <div className="form-control">
                    <label className="label">
                      <span className="label-text">Уровень активности</span>
                    </label>
                    <select
                      className="select select-bordered"
                      value={userProfile.activity}
                      onChange={(e) =>
                        setUserProfile({
                          ...userProfile,
                          activity: e.target.value as any,
                        })
                      }
                    >
                      <option value="sedentary">Сидячий образ жизни</option>
                      <option value="light">
                        Легкая активность (1-3 дня/нед)
                      </option>
                      <option value="moderate">Умеренная (3-5 дней/нед)</option>
                      <option value="active">Активная (6-7 дней/нед)</option>
                      <option value="very_active">Очень активная</option>
                    </select>
                  </div>
                  <div className="form-control">
                    <label className="label">
                      <span className="label-text">Цель</span>
                    </label>
                    <select
                      className="select select-bordered"
                      value={userProfile.goal}
                      onChange={(e) =>
                        setUserProfile({
                          ...userProfile,
                          goal: e.target.value as any,
                        })
                      }
                    >
                      <option value="lose">Похудение</option>
                      <option value="maintain">Поддержание веса</option>
                      <option value="gain">Набор массы</option>
                    </select>
                  </div>
                </div>
              )}

              <div className="divider"></div>

              {/* Calculated Macros Display */}
              <div className="grid grid-cols-2 md:grid-cols-5 gap-4 text-center">
                <div className="stat bg-base-100 rounded-lg p-3">
                  <div className="stat-title text-xs">Калории/день</div>
                  <div className="stat-value text-2xl">
                    {calculateDailyCalories()}
                  </div>
                </div>
                <div className="stat bg-base-100 rounded-lg p-3">
                  <div className="stat-title text-xs">Белки (г)</div>
                  <div className="stat-value text-2xl">
                    {calculateMacros().protein}
                  </div>
                </div>
                <div className="stat bg-base-100 rounded-lg p-3">
                  <div className="stat-title text-xs">Углеводы (г)</div>
                  <div className="stat-value text-2xl">
                    {calculateMacros().carbs}
                  </div>
                </div>
                <div className="stat bg-base-100 rounded-lg p-3">
                  <div className="stat-title text-xs">Жиры (г)</div>
                  <div className="stat-value text-2xl">
                    {calculateMacros().fats}
                  </div>
                </div>
                <div className="stat bg-base-100 rounded-lg p-3">
                  <div className="stat-title text-xs">Клетчатка (г)</div>
                  <div className="stat-value text-2xl">
                    {calculateMacros().fiber}
                  </div>
                </div>
              </div>

              <div className="card-actions justify-center mt-6">
                <div className="btn-group">
                  <button
                    className={`btn ${selectedPeriod === 'day' ? 'btn-active' : ''}`}
                    onClick={() => setSelectedPeriod('day')}
                  >
                    На день
                  </button>
                  <button
                    className={`btn ${selectedPeriod === 'week' ? 'btn-active' : ''}`}
                    onClick={() => setSelectedPeriod('week')}
                  >
                    На неделю
                  </button>
                </div>
                <button
                  onClick={generateOptimalBasket}
                  className="btn btn-primary btn-lg"
                >
                  🛒 Сформировать оптимальную корзину
                </button>
              </div>
            </div>
          </div>
        </AnimatedSection>

        {/* Generated Basket Section */}
        {generatedBasket && (
          <AnimatedSection animation="slideUp" delay={0.2}>
            <div className="card bg-base-100 shadow-xl mb-8">
              <div className="card-body">
                <h2 className="card-title text-2xl mb-4">
                  📊 Оптимальная корзина на{' '}
                  {selectedPeriod === 'day' ? 'день' : 'неделю'}
                </h2>

                {/* Basket Items */}
                <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-4 mb-6">
                  {generatedBasket.items.map((item: any, index: number) => (
                    <div key={index} className="card bg-base-200 compact">
                      <div className="card-body">
                        <div className="flex justify-between items-center">
                          <div className="flex items-center gap-2">
                            <span className="text-2xl">{item.image}</span>
                            <div>
                              <h4 className="font-semibold">{item.name}</h4>
                              <p className="text-xs text-base-content/60">
                                {item.quantity}г
                              </p>
                            </div>
                          </div>
                          <div className="text-right">
                            <div className="font-semibold">
                              {Math.round(item.totalPrice)} RSD
                            </div>
                            <div className="text-xs text-base-content/60">
                              {(item.nutrients.calories * item.quantity) / 100}{' '}
                              ккал
                            </div>
                          </div>
                        </div>
                      </div>
                    </div>
                  ))}
                </div>

                {/* Nutritional Summary */}
                <div className="bg-base-200 rounded-lg p-4 mb-6">
                  <h3 className="font-semibold mb-3">
                    Пищевая ценность корзины:
                  </h3>
                  <div className="grid grid-cols-2 md:grid-cols-3 gap-4">
                    <div>
                      <div className="text-sm text-base-content/60">
                        Калории
                      </div>
                      <div className="font-semibold">
                        {Math.round(generatedBasket.totals.calories)} /{' '}
                        {generatedBasket.targetMacros.calories} ккал
                      </div>
                      <progress
                        className="progress progress-primary"
                        value={generatedBasket.totals.calories}
                        max={generatedBasket.targetMacros.calories}
                      ></progress>
                    </div>
                    <div>
                      <div className="text-sm text-base-content/60">Белки</div>
                      <div className="font-semibold">
                        {Math.round(generatedBasket.totals.protein)} /{' '}
                        {generatedBasket.targetMacros.protein} г
                      </div>
                      <progress
                        className="progress progress-success"
                        value={generatedBasket.totals.protein}
                        max={generatedBasket.targetMacros.protein}
                      ></progress>
                    </div>
                    <div>
                      <div className="text-sm text-base-content/60">
                        Углеводы
                      </div>
                      <div className="font-semibold">
                        {Math.round(generatedBasket.totals.carbs)} /{' '}
                        {generatedBasket.targetMacros.carbs} г
                      </div>
                      <progress
                        className="progress progress-warning"
                        value={generatedBasket.totals.carbs}
                        max={generatedBasket.targetMacros.carbs}
                      ></progress>
                    </div>
                  </div>
                </div>

                {/* Export Button */}
                <div className="card-actions justify-center mb-6">
                  <button
                    onClick={exportShoppingList}
                    className="btn btn-primary btn-lg"
                  >
                    📱 Экспорт списка для сборщика
                  </button>
                  <div className="text-xs text-base-content/60 text-center mt-2">
                    Список будет скопирован в буфер обмена и сохранен в файл
                  </div>
                </div>

                {/* Supermarket Comparison */}
                <div>
                  <h3 className="font-semibold mb-3">
                    💰 Сравнение цен по супермаркетам:
                  </h3>
                  <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-3">
                    {generatedBasket.supermarketPrices
                      ?.sort((a: any, b: any) => a.totalPrice - b.totalPrice)
                      .map((market: any, index: number) => (
                        <div
                          key={index}
                          className={`card ${market.color} ${index === 0 ? 'ring-2 ring-success' : ''}`}
                        >
                          <div className="card-body p-4">
                            <div className="flex justify-between items-center">
                              <div className="flex items-center gap-2">
                                <span className="text-2xl">{market.logo}</span>
                                <div>
                                  <h4 className="font-semibold">
                                    {market.name}
                                  </h4>
                                  {index === 0 && (
                                    <div className="badge badge-success badge-sm">
                                      Лучшая цена
                                    </div>
                                  )}
                                </div>
                              </div>
                              <div className="text-right">
                                <div className="text-2xl font-bold">
                                  {market.totalPrice} RSD
                                </div>
                                {market.savings > 0 && (
                                  <div className="text-sm text-success">
                                    -{market.savings} RSD
                                  </div>
                                )}
                              </div>
                            </div>
                          </div>
                        </div>
                      ))}
                  </div>
                </div>
              </div>
            </div>
          </AnimatedSection>
        )}

        {/* Info Cards */}
        <AnimatedSection animation="slideUp" delay={0.2}>
          <div className="grid grid-cols-1 md:grid-cols-3 gap-6 mb-8">
            <div className="card bg-base-100 shadow-xl">
              <div className="card-body">
                <div className="text-4xl mb-2">🚫</div>
                <h3 className="card-title">Без обмана</h3>
                <p>
                  Никаких продуктов с добавленным сахаром, красителями и
                  консервантами под видом &ldquo;здоровых&rdquo;
                </p>
              </div>
            </div>
            <div className="card bg-base-100 shadow-xl">
              <div className="card-body">
                <div className="text-4xl mb-2">✅</div>
                <h3 className="card-title">Только польза</h3>
                <p>
                  Каждый продукт проверен диетологами и содержит реальные
                  питательные вещества
                </p>
              </div>
            </div>
            <div className="card bg-base-100 shadow-xl">
              <div className="card-body">
                <div className="text-4xl mb-2">🌍</div>
                <h3 className="card-title">Локальное и сезонное</h3>
                <p>
                  Приоритет местным фермерам и сезонным продуктам для
                  максимальной пользы
                </p>
              </div>
            </div>
          </div>
        </AnimatedSection>

        {/* Categories */}
        <AnimatedSection animation="slideUp" delay={0.3}>
          <div className="mb-8">
            <h2 className="text-2xl font-bold mb-4">Категории</h2>
            <div className="flex flex-wrap gap-2">
              {categories.map((category) => (
                <button
                  key={category.id}
                  onClick={() => setSelectedCategory(category.id)}
                  className={`btn ${
                    selectedCategory === category.id
                      ? 'btn-primary'
                      : 'btn-outline'
                  }`}
                >
                  <span className="text-xl mr-1">{category.icon}</span>
                  {category.name}
                </button>
              ))}
            </div>
          </div>
        </AnimatedSection>

        {/* Filters */}
        <AnimatedSection animation="slideUp" delay={0.4}>
          <div className="card bg-base-100 shadow-xl mb-8">
            <div className="card-body">
              <div className="flex justify-between items-center mb-4">
                <h3 className="card-title text-2xl">🛡️ Фильтры безопасности</h3>
                <div className="form-control">
                  <label className="label cursor-pointer">
                    <span className="label-text font-bold text-lg mr-2">
                      Только здоровые
                    </span>
                    <input
                      type="checkbox"
                      className="toggle toggle-success toggle-lg"
                      checked={showOnlyHealthy}
                      onChange={(e) => setShowOnlyHealthy(e.target.checked)}
                    />
                  </label>
                </div>
              </div>

              {/* Main harmful filters with bigger UI */}
              <div className="bg-error/10 rounded-xl p-6 mb-6">
                <div className="flex justify-between items-center mb-4">
                  <h4 className="text-lg font-bold text-error">
                    ☠️ ИСКЛЮЧИТЬ ВРЕДНЫЕ ДОБАВКИ
                  </h4>
                  {activeHarmfulFilters > 0 && (
                    <div className="badge badge-error badge-lg">
                      Активно фильтров: {activeHarmfulFilters}
                    </div>
                  )}
                </div>
                <div className="grid grid-cols-1 md:grid-cols-2 gap-4">
                  <label className="label cursor-pointer bg-base-100 rounded-lg p-4 hover:bg-error/20 transition-colors">
                    <div className="flex items-center gap-3">
                      <input
                        type="checkbox"
                        className="checkbox checkbox-error checkbox-lg"
                        checked={showNoSugar}
                        onChange={(e) => setShowNoSugar(e.target.checked)}
                      />
                      <div>
                        <span className="text-lg font-bold">🍬 БЕЗ САХАРА</span>
                        <div className="text-sm text-base-content/60">
                          Скрывает продукты с сахаром &gt;5г
                        </div>
                      </div>
                    </div>
                  </label>
                  <label className="label cursor-pointer bg-base-100 rounded-lg p-4 hover:bg-error/20 transition-colors">
                    <div className="flex items-center gap-3">
                      <input
                        type="checkbox"
                        className="checkbox checkbox-error checkbox-lg"
                        checked={showNoPalmOil}
                        onChange={(e) => setShowNoPalmOil(e.target.checked)}
                      />
                      <div>
                        <span className="text-lg font-bold">
                          🌴 БЕЗ ПАЛЬМОВОГО МАСЛА
                        </span>
                        <div className="text-sm text-base-content/60">
                          Закупоривает сосуды
                        </div>
                      </div>
                    </div>
                  </label>
                  <label className="label cursor-pointer bg-base-100 rounded-lg p-4 hover:bg-error/20 transition-colors">
                    <div className="flex items-center gap-3">
                      <input
                        type="checkbox"
                        className="checkbox checkbox-error checkbox-lg"
                        checked={showNoAspartame}
                        onChange={(e) => setShowNoAspartame(e.target.checked)}
                      />
                      <div>
                        <span className="text-lg font-bold">
                          💊 БЕЗ АСПАРТАМА
                        </span>
                        <div className="text-sm text-base-content/60">
                          Токсичный сахарозаменитель
                        </div>
                      </div>
                    </div>
                  </label>
                  <label className="label cursor-pointer bg-base-100 rounded-lg p-4 hover:bg-error/20 transition-colors">
                    <div className="flex items-center gap-3">
                      <input
                        type="checkbox"
                        className="checkbox checkbox-error checkbox-lg"
                        checked={showNoTransFats}
                        onChange={(e) => setShowNoTransFats(e.target.checked)}
                      />
                      <div>
                        <span className="text-lg font-bold">
                          🚫 БЕЗ ТРАНСЖИРОВ
                        </span>
                        <div className="text-sm text-base-content/60">
                          Вызывают инфаркты
                        </div>
                      </div>
                    </div>
                  </label>
                  <label className="label cursor-pointer bg-base-100 rounded-lg p-4 hover:bg-error/20 transition-colors">
                    <div className="flex items-center gap-3">
                      <input
                        type="checkbox"
                        className="checkbox checkbox-error checkbox-lg"
                        checked={showNoMSG}
                        onChange={(e) => setShowNoMSG(e.target.checked)}
                      />
                      <div>
                        <span className="text-lg font-bold">
                          🧪 БЕЗ ГЛУТАМАТА (E621)
                        </span>
                        <div className="text-sm text-base-content/60">
                          Вызывает зависимость
                        </div>
                      </div>
                    </div>
                  </label>
                </div>

                {/* Quick filter buttons */}
                <div className="flex flex-wrap gap-2 mt-4">
                  <button
                    onClick={() => {
                      setShowNoSugar(false);
                      setShowNoPalmOil(false);
                      setShowNoAspartame(false);
                      setShowNoTransFats(false);
                      setShowNoMSG(false);
                      setShowOnlyHealthy(false);
                    }}
                    className="btn btn-warning btn-sm"
                  >
                    🚨 Показать ВСЕ продукты (включая вредные)
                  </button>
                  <button
                    onClick={() => {
                      setShowNoSugar(true);
                      setShowNoPalmOil(true);
                      setShowNoAspartame(true);
                      setShowNoTransFats(true);
                      setShowNoMSG(true);
                    }}
                    className="btn btn-success btn-sm"
                  >
                    ✅ Включить ВСЕ защитные фильтры
                  </button>
                </div>
              </div>

              <div className="divider">Качество продуктов</div>
              <div className="flex flex-wrap gap-4 mb-4">
                <label className="label cursor-pointer">
                  <input
                    type="checkbox"
                    className="checkbox checkbox-success mr-2"
                    checked={showOrganic}
                    onChange={(e) => setShowOrganic(e.target.checked)}
                  />
                  <span className="label-text">🌿 Органические продукты</span>
                </label>
                <label className="label cursor-pointer">
                  <input
                    type="checkbox"
                    className="checkbox checkbox-success mr-2"
                    checked={showLocal}
                    onChange={(e) => setShowLocal(e.target.checked)}
                  />
                  <span className="label-text">📍 Местное производство</span>
                </label>
                <label className="label cursor-pointer">
                  <input
                    type="checkbox"
                    className="checkbox checkbox-success mr-2"
                    checked={showSeasonal}
                    onChange={(e) => setShowSeasonal(e.target.checked)}
                  />
                  <span className="label-text">🗓️ Сезонные продукты</span>
                </label>
              </div>

              {!showOnlyHealthy && (
                <div className="alert alert-warning mt-4">
                  <svg
                    xmlns="http://www.w3.org/2000/svg"
                    className="stroke-current shrink-0 h-6 w-6"
                    fill="none"
                    viewBox="0 0 24 24"
                  >
                    <path
                      strokeLinecap="round"
                      strokeLinejoin="round"
                      strokeWidth="2"
                      d="M12 9v2m0 4h.01m-6.938 4h13.856c1.54 0 2.502-1.667 1.732-3L13.732 4c-.77-1.333-2.694-1.333-3.464 0L3.34 16c-.77 1.333.192 3 1.732 3z"
                    />
                  </svg>
                  <span>
                    Внимание! Отображаются продукты с вредными добавками.
                    Смотрите предупреждения на карточках!
                  </span>
                </div>
              )}
            </div>
          </div>
        </AnimatedSection>

        {/* Products Grid */}
        <div>
          {/* Products counter */}
          <div className="flex justify-between items-center mb-6">
            <h2 className="text-2xl font-bold">
              {showOnlyHealthy ? '✅ Здоровые продукты' : '🛒 Все продукты'}
            </h2>
            <div className="flex gap-2">
              <div className="badge badge-lg badge-success">
                Показано: {filteredProducts.length}
              </div>
              {(activeHarmfulFilters > 0 || showOnlyHealthy) && (
                <div className="badge badge-lg badge-error">
                  Скрыто вредных: {allProducts.length - filteredProducts.length}
                </div>
              )}
            </div>
          </div>
          <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-6">
            {filteredProducts.map((product, _index) => (
              <div key={product.id}>
                <div
                  className={`card shadow-xl hover:shadow-2xl transition-all h-full flex flex-col ${
                    product.harmfulIngredients &&
                    product.harmfulIngredients.length > 0
                      ? 'bg-error/10 border-2 border-error'
                      : 'bg-base-100'
                  }`}
                >
                  <div className="card-body">
                    <div className="flex items-start justify-between">
                      <div className="text-5xl mb-4">{product.image}</div>
                      <div className="flex gap-1 flex-wrap">
                        {product.harmfulIngredients &&
                          product.harmfulIngredients.length > 0 && (
                            <div className="badge badge-error">⚠️ Опасно</div>
                          )}
                        {product.organic && (
                          <div className="badge badge-success">Органик</div>
                        )}
                        {product.local && (
                          <div className="badge badge-info">Местное</div>
                        )}
                        {product.seasonal && (
                          <div className="badge badge-warning">Сезон</div>
                        )}
                      </div>
                    </div>

                    <h3 className="card-title">{product.name}</h3>
                    <p className="text-base-content/70 mb-2">
                      {product.description}
                    </p>

                    {/* Warnings for harmful products */}
                    {product.warnings && product.warnings.length > 0 && (
                      <div className="alert alert-error mb-2">
                        <svg
                          xmlns="http://www.w3.org/2000/svg"
                          className="stroke-current shrink-0 h-4 w-4"
                          fill="none"
                          viewBox="0 0 24 24"
                        >
                          <path
                            strokeLinecap="round"
                            strokeLinejoin="round"
                            strokeWidth="2"
                            d="M10 14l2-2m0 0l2-2m-2 2l-2-2m2 2l2 2m7-2a9 9 0 11-18 0 9 9 0 0118 0z"
                          />
                        </svg>
                        <div className="text-xs">
                          {product.warnings.map((warning, i) => (
                            <div key={i}>• {warning}</div>
                          ))}
                        </div>
                      </div>
                    )}

                    {/* Harmful ingredients */}
                    {product.harmfulIngredients &&
                      product.harmfulIngredients.length > 0 && (
                        <div className="bg-error/20 rounded-lg p-2 mb-2">
                          <h5 className="text-xs font-bold text-error mb-1">
                            ☠️ Вредные добавки:
                          </h5>
                          <div className="text-xs text-error">
                            {product.harmfulIngredients.join(', ')}
                          </div>
                        </div>
                      )}

                    <div className="divider my-2"></div>

                    <div className="space-y-2">
                      <h4 className="font-semibold">Польза:</h4>
                      <ul className="text-sm space-y-1">
                        {product.benefits.map((benefit, i) => (
                          <li key={i} className="flex items-center gap-2">
                            <span className="text-success">✓</span>
                            {benefit}
                          </li>
                        ))}
                      </ul>
                    </div>

                    {product.nutrients && (
                      <>
                        <div className="divider my-2"></div>
                        <div className="text-xs space-y-1">
                          {product.nutrients.calories && (
                            <div>Калории: {product.nutrients.calories}</div>
                          )}
                          {product.nutrients.protein && (
                            <div>Белки: {product.nutrients.protein}г</div>
                          )}
                          {product.nutrients.carbs && (
                            <div>Углеводы: {product.nutrients.carbs}г</div>
                          )}
                          {product.nutrients.fiber && (
                            <div>Клетчатка: {product.nutrients.fiber}г</div>
                          )}
                          {product.nutrients.vitamins && (
                            <div>
                              Витамины: {product.nutrients.vitamins.join(', ')}
                            </div>
                          )}
                        </div>
                      </>
                    )}

                    <div className="card-actions justify-between items-center mt-4">
                      <div className="text-2xl font-bold text-primary">
                        {product.price} RSD/{product.unit}
                      </div>
                      <button className="btn btn-primary btn-sm">
                        В корзину
                      </button>
                    </div>
                  </div>
                </div>
              </div>
            ))}
          </div>
        </div>

        {/* Educational Section about Harmful Additives */}
        <AnimatedSection animation="slideUp" delay={0.55}>
          <div className="mt-8 card bg-error/5 border-2 border-error">
            <div className="card-body">
              <h2 className="card-title text-2xl mb-4">
                ☠️ Опасные добавки - враги здоровья
              </h2>
              <div className="grid md:grid-cols-2 lg:grid-cols-3 gap-4">
                <div className="card bg-base-100">
                  <div className="card-body compact">
                    <h3 className="font-bold text-error">🍬 Скрытый сахар</h3>
                    <p className="text-sm">
                      <strong>Другие названия:</strong> Кукурузный сироп,
                      фруктоза, декстроза, мальтодекстрин, сироп агавы
                    </p>
                    <p className="text-sm text-error">
                      <strong>Вред:</strong> Диабет, ожирение, кариес,
                      зависимость
                    </p>
                  </div>
                </div>
                <div className="card bg-base-100">
                  <div className="card-body compact">
                    <h3 className="font-bold text-error">🌴 Пальмовое масло</h3>
                    <p className="text-sm">
                      <strong>Другие названия:</strong> Растительный жир,
                      стеарин, пальмитиновая кислота
                    </p>
                    <p className="text-sm text-error">
                      <strong>Вред:</strong> Закупорка сосудов, рак, повышение
                      холестерина
                    </p>
                  </div>
                </div>
                <div className="card bg-base-100">
                  <div className="card-body compact">
                    <h3 className="font-bold text-error">💊 Аспартам (E951)</h3>
                    <p className="text-sm">
                      <strong>Другие названия:</strong> NutraSweet, Equal,
                      Canderel
                    </p>
                    <p className="text-sm text-error">
                      <strong>Вред:</strong> Головные боли, депрессия, риск рака
                    </p>
                  </div>
                </div>
                <div className="card bg-base-100">
                  <div className="card-body compact">
                    <h3 className="font-bold text-error">🧈 Трансжиры</h3>
                    <p className="text-sm">
                      <strong>Другие названия:</strong> Гидрогенизированное
                      масло, маргарин
                    </p>
                    <p className="text-sm text-error">
                      <strong>Вред:</strong> Инфаркт, инсульт, воспаления
                    </p>
                  </div>
                </div>
                <div className="card bg-base-100">
                  <div className="card-body compact">
                    <h3 className="font-bold text-error">
                      🧪 Глутамат натрия (E621)
                    </h3>
                    <p className="text-sm">
                      <strong>Другие названия:</strong> MSG, усилитель вкуса
                    </p>
                    <p className="text-sm text-error">
                      <strong>Вред:</strong> Переедание, головные боли,
                      повреждение нервов
                    </p>
                  </div>
                </div>
                <div className="card bg-base-100">
                  <div className="card-body compact">
                    <h3 className="font-bold text-error">🎨 Красители</h3>
                    <p className="text-sm">
                      <strong>Примеры:</strong> E102 (тартразин), E110, E124,
                      E129
                    </p>
                    <p className="text-sm text-error">
                      <strong>Вред:</strong> Гиперактивность у детей, аллергии,
                      астма
                    </p>
                  </div>
                </div>
              </div>

              <div className="divider"></div>

              <div className="alert alert-info">
                <svg
                  xmlns="http://www.w3.org/2000/svg"
                  fill="none"
                  viewBox="0 0 24 24"
                  className="stroke-current shrink-0 w-6 h-6"
                >
                  <path
                    strokeLinecap="round"
                    strokeLinejoin="round"
                    strokeWidth="2"
                    d="M13 16h-1v-4h-1m1-4h.01M21 12a9 9 0 11-18 0 9 9 0 0118 0z"
                  ></path>
                </svg>
                <div>
                  <h3 className="font-bold">💡 Как защитить себя:</h3>
                  <ul className="text-sm mt-2 space-y-1">
                    <li>
                      • Читайте состав ВСЕГДА - даже у &ldquo;здоровых&rdquo;
                      продуктов
                    </li>
                    <li>• Если больше 5 ингредиентов - это не еда, а химия</li>
                    <li>• Не можете произнести название - не покупайте</li>
                    <li>• Выбирайте цельные продукты без упаковки</li>
                    <li>• Готовьте дома из простых ингредиентов</li>
                  </ul>
                </div>
              </div>
            </div>
          </div>
        </AnimatedSection>

        {/* Educational Section */}
        <AnimatedSection animation="fadeIn" delay={0.6}>
          <div className="mt-12 card bg-warning/10 border-2 border-warning">
            <div className="card-body">
              <h2 className="card-title text-2xl mb-4">
                ⚠️ Остерегайтесь маркетинговых уловок!
              </h2>
              <div className="grid md:grid-cols-2 gap-6">
                <div>
                  <h3 className="font-bold mb-2 text-error">
                    ❌ НЕ здоровая еда:
                  </h3>
                  <ul className="space-y-2">
                    <li>
                      • &ldquo;Диетическая&rdquo; кола - химия без сахара, но с
                      аспартамом
                    </li>
                    <li>
                      • &ldquo;Фитнес&rdquo; батончики - сахар и пальмовое масло
                    </li>
                    <li>
                      • &ldquo;Обезжиренный&rdquo; йогурт - сахар вместо жира
                    </li>
                    <li>
                      • &ldquo;Витаминизированные&rdquo; конфеты - сахар с
                      витаминами
                    </li>
                    <li>
                      • &ldquo;Полезные&rdquo; чипсы из овощей - жареные в масле
                    </li>
                    <li>
                      • &ldquo;Натуральные&rdquo; соки - концентрат сахара без
                      клетчатки
                    </li>
                  </ul>
                </div>
                <div>
                  <h3 className="font-bold mb-2 text-success">
                    ✅ Настоящая здоровая еда:
                  </h3>
                  <ul className="space-y-2">
                    <li>• Свежие овощи и фрукты без обработки</li>
                    <li>• Цельнозерновые крупы и бобовые</li>
                    <li>• Орехи и семена без соли и сахара</li>
                    <li>• Дикая рыба, богатая омега-3</li>
                    <li>• Ферментированные продукты для кишечника</li>
                    <li>• Чистая вода вместо любых напитков</li>
                  </ul>
                </div>
              </div>
            </div>
          </div>
        </AnimatedSection>

        {/* Tips Section */}
        <AnimatedSection animation="slideUp" delay={0.7}>
          <div className="mt-8 grid md:grid-cols-3 gap-6">
            <div className="card bg-primary/10">
              <div className="card-body">
                <h3 className="card-title">💡 Совет дня</h3>
                <p>
                  Читайте состав! Если больше 5 ингредиентов или есть непонятные
                  названия - это не еда, а продукт пищевой промышленности.
                </p>
              </div>
            </div>
            <div className="card bg-success/10">
              <div className="card-body">
                <h3 className="card-title">🌈 Правило радуги</h3>
                <p>
                  Ешьте овощи и фрукты всех цветов радуги каждую неделю - каждый
                  цвет дает уникальные питательные вещества.
                </p>
              </div>
            </div>
            <div className="card bg-info/10">
              <div className="card-body">
                <h3 className="card-title">⏰ Время приема</h3>
                <p>
                  Фрукты - утром, овощи - днем, белки - вечером. Орехи и семена
                  - идеальный перекус между приемами пищи.
                </p>
              </div>
            </div>
          </div>
        </AnimatedSection>
      </div>
    </div>
  );
}
