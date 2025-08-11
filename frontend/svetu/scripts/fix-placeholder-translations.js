#!/usr/bin/env node

const fs = require('fs');
const path = require('path');

// Маппинг английских фраз на русские переводы
const translationMap = {
  // Storefronts
  'No Variants To Manage': 'Нет вариантов для управления',
  'Total Variants': 'Всего вариантов',
  'Total Stock': 'Общий запас',
  'Average Price': 'Средняя цена',
  'Manage Variants': 'Управление вариантами',
  Variant: 'Вариант',
  Price: 'Цена',
  Main: 'Основной',
  'Out Of Stock': 'Нет в наличии',
  'Quick Actions': 'Быстрые действия',
  'Set Stock Prompt': 'Введите количество на складе',
  'Set Stock For All': 'Установить запас для всех',
  'Set Price Prompt': 'Введите цену',
  'Set Price For All': 'Установить цену для всех',
  ',': ',',
  'Configure Stock And Prices': 'Настроить запасы и цены',
  'Back To Settings': 'Вернуться к настройкам',
  'Confirm Variants': 'Подтвердить варианты',
  'No Variant Attributes': 'Нет атрибутов вариантов',
  'Generate Variants': 'Сгенерировать варианты',
  'All Products': 'Все товары',
  'Active Only': 'Только активные',
  'Inactive Only': 'Только неактивные',
  Inventory: 'Инвентарь',
  'Price Range': 'Диапазон цен',
  Min: 'Мин',
  Max: 'Макс',
  'Product Name': 'Название товара',
  Actions: 'Действия',
  'No Products Found': 'Товары не найдены',
  'No Products': 'Нет товаров',
  Label: 'Метка',
  'No Variant Attributes For Category': 'Нет атрибутов вариантов для категории',
  'Category Does Not Support Variants': 'Категория не поддерживает варианты',
  'Select Values For Variants': 'Выберите значения для вариантов',
  'Affects Stock': 'Влияет на запас',
  'Attributes From Previous Step': 'Атрибуты из предыдущего шага',
  'Stock Settings': 'Настройки запасов',
  'Default Stock Quantity': 'Количество по умолчанию',
  'Use Individual Quantities': 'Использовать индивидуальные количества',
  'Error Loading Products': 'Ошибка загрузки товаров',
  'Total Revenue': 'Общая выручка',
  'Conversion Rate': 'Коэффициент конверсии',
  'New Customers': 'Новые покупатели',
  'Returning Customers': 'Постоянные покупатели',
  'Popular Products': 'Популярные товары',
  'Low Stock Products': 'Товары с низким запасом',
  'Recent Activity': 'Недавняя активность',
  'Analytics Overview': 'Обзор аналитики',
  'Sales Performance': 'Эффективность продаж',
  'Customer Demographics': 'Демография покупателей',
  'Inventory Status': 'Статус инвентаря',
  'Marketing Campaigns': 'Маркетинговые кампании',

  // Reviews
  'Review Product': 'Оценить товар',
  'Review Service': 'Оценить сервис',
  'Review Title': 'Заголовок отзыва',
  'Review Text': 'Текст отзыва',
  'Review Photos': 'Фотографии к отзыву',
  'Submit Review': 'Отправить отзыв',
  'Cancel Review': 'Отменить отзыв',
  'Edit Review': 'Редактировать отзыв',
  'Delete Review': 'Удалить отзыв',
  'Report Review': 'Пожаловаться на отзыв',
  'Helpful Review': 'Полезный отзыв',
  'Not Helpful Review': 'Не полезный отзыв',
  'Sort By Most Recent': 'Сортировать по новизне',
  'Sort By Most Helpful': 'Сортировать по полезности',
  'Sort By Highest Rating': 'Сортировать по высшему рейтингу',
  'Sort By Lowest Rating': 'Сортировать по низшему рейтингу',
  'Filter By Rating': 'Фильтр по рейтингу',
  'Filter By Verified Purchase': 'Фильтр по подтвержденным покупкам',
  'Verified Purchase': 'Подтвержденная покупка',
  'Seller Response': 'Ответ продавца',
  'Write A Response': 'Написать ответ',
  'Response Placeholder': 'Напишите ваш ответ здесь...',
  'Submit Response': 'Отправить ответ',
  'Edit Response': 'Редактировать ответ',
  'Delete Response': 'Удалить ответ',
  'Response Submitted': 'Ответ отправлен',
  'Response Updated': 'Ответ обновлен',
  'Response Deleted': 'Ответ удален',
  'No Reviews Yet': 'Пока нет отзывов',
  'Be The First To Review': 'Будьте первым, кто оставит отзыв',
  'Total Reviews': 'Всего отзывов',
  'Average Rating': 'Средний рейтинг',
  Star: 'звезда',
  Stars: 'звезд',
  'Recommend This Product': 'Рекомендую этот товар',
  'Would Buy Again': 'Куплю снова',
  'Value For Money': 'Соотношение цена/качество',
  'Product Quality': 'Качество товара',
  'Customer Service': 'Обслуживание клиентов',
  'Delivery Speed': 'Скорость доставки',
  'Packaging Quality': 'Качество упаковки',
  'Product As Described': 'Товар соответствует описанию',
  Communication: 'Коммуникация',
  'Load More Reviews': 'Загрузить больше отзывов',
  'Showing Reviews': 'Показано отзывов',
  'Of Total': 'из всех',
  'Review Guidelines': 'Правила отзывов',
  'Inappropriate Content': 'Неподходящий контент',
  'Spam Or Fake': 'Спам или подделка',
  'Offensive Language': 'Оскорбительный язык',
  'Personal Information': 'Личная информация',
  'Report Submitted': 'Жалоба отправлена',
  'Thank You For Report': 'Спасибо за вашу жалобу',

  // Cart
  'Empty Cart': 'Корзина пуста',
  'Add To Cart': 'Добавить в корзину',
  'Remove From Cart': 'Удалить из корзины',
  'Update Quantity': 'Обновить количество',
  'Cart Total': 'Итого в корзине',
  'Proceed To Checkout': 'Перейти к оформлению',
  'Continue Shopping': 'Продолжить покупки',
  'Apply Coupon': 'Применить купон',
  'Coupon Code': 'Код купона',
  'Invalid Coupon': 'Недействительный купон',
  'Coupon Applied': 'Купон применен',
  'Remove Coupon': 'Удалить купон',
  Subtotal: 'Промежуточный итог',
  Discount: 'Скидка',
  Shipping: 'Доставка',
  Tax: 'Налог',
  'Grand Total': 'Общий итог',
  'Save For Later': 'Сохранить на потом',
  'Move To Cart': 'Переместить в корзину',
  'Saved Items': 'Сохраненные товары',

  // Search
  'Search Results': 'Результаты поиска',
  'No Results Found': 'Ничего не найдено',
  'Try Different Keywords': 'Попробуйте другие ключевые слова',
  'Popular Searches': 'Популярные запросы',
  'Recent Searches': 'Недавние поиски',
  'Clear Search History': 'Очистить историю поиска',
  'Search Filters': 'Фильтры поиска',
  'Apply Filters': 'Применить фильтры',
  'Clear Filters': 'Очистить фильтры',
  'Sort By': 'Сортировать по',
  Relevance: 'Релевантности',
  'Price Low To High': 'Цене: по возрастанию',
  'Price High To Low': 'Цене: по убыванию',
  'Newest First': 'Новые сначала',
  'Best Sellers': 'Бестселлеры',
  'Customer Rating': 'Рейтингу покупателей',

  // Checkout
  Checkout: 'Оформление заказа',
  'Billing Address': 'Адрес для выставления счета',
  'Shipping Address': 'Адрес доставки',
  'Same As Billing': 'Совпадает с адресом для счета',
  'Payment Method': 'Способ оплаты',
  'Credit Card': 'Кредитная карта',
  'Debit Card': 'Дебетовая карта',
  PayPal: 'PayPal',
  'Bank Transfer': 'Банковский перевод',
  'Cash On Delivery': 'Наложенный платеж',
  'Place Order': 'Разместить заказ',
  'Order Summary': 'Сводка заказа',
  'Edit Cart': 'Редактировать корзину',
  'Promo Code': 'Промокод',
  'Gift Message': 'Подарочное сообщение',
  'Estimated Delivery': 'Ожидаемая доставка',
  'Express Shipping': 'Экспресс-доставка',
  'Standard Shipping': 'Стандартная доставка',
  'Free Shipping': 'Бесплатная доставка',

  // Common
  Loading: 'Загрузка',
  Error: 'Ошибка',
  Success: 'Успешно',
  Warning: 'Предупреждение',
  Info: 'Информация',
  Confirm: 'Подтвердить',
  Cancel: 'Отмена',
  Save: 'Сохранить',
  Delete: 'Удалить',
  Edit: 'Редактировать',
  Close: 'Закрыть',
  Submit: 'Отправить',
  Reset: 'Сбросить',
  Back: 'Назад',
  Next: 'Далее',
  Previous: 'Предыдущий',
  Finish: 'Завершить',
  Download: 'Скачать',
  Upload: 'Загрузить',
  Refresh: 'Обновить',
  Retry: 'Повторить',
  Share: 'Поделиться',
  Copy: 'Копировать',
  Paste: 'Вставить',
  Cut: 'Вырезать',
  'Select All': 'Выбрать все',
  'Select None': 'Снять выделение',
  Clear: 'Очистить',
  Search: 'Поиск',
  Filter: 'Фильтр',
  Sort: 'Сортировка',
  Export: 'Экспорт',
  Import: 'Импорт',
  Print: 'Печать',
  Help: 'Помощь',
  Settings: 'Настройки',
  Logout: 'Выход',
  Login: 'Вход',
  Register: 'Регистрация',
  Profile: 'Профиль',
  Account: 'Аккаунт',
  Dashboard: 'Панель управления',
  Home: 'Главная',

  // Orders
  'Order Number': 'Номер заказа',
  'Order Date': 'Дата заказа',
  'Order Status': 'Статус заказа',
  'Order Total': 'Сумма заказа',
  'Order Details': 'Детали заказа',
  'Track Order': 'Отследить заказ',
  'Cancel Order': 'Отменить заказ',
  'Return Order': 'Вернуть заказ',
  Reorder: 'Заказать снова',
  Invoice: 'Счет-фактура',
  'Download Invoice': 'Скачать счет',
  Pending: 'В ожидании',
  Processing: 'В обработке',
  Shipped: 'Отправлен',
  Delivered: 'Доставлен',
  Cancelled: 'Отменен',
  Refunded: 'Возвращен',
  Failed: 'Неудачно',
  'On Hold': 'Приостановлен',
  Completed: 'Завершен',
  'Payment Pending': 'Ожидание оплаты',
  'Payment Failed': 'Ошибка оплаты',
  'Payment Complete': 'Оплата завершена',
};

// Маппинг для сербских переводов
const serbianMap = {
  // Storefronts
  'No Variants To Manage': 'Nema varijanti za upravljanje',
  'Total Variants': 'Ukupno varijanti',
  'Total Stock': 'Ukupne zalihe',
  'Average Price': 'Prosečna cena',
  'Manage Variants': 'Upravljanje varijantama',
  Variant: 'Varijanta',
  Price: 'Cena',
  Main: 'Glavno',
  'Out Of Stock': 'Nema na stanju',
  'Quick Actions': 'Brze akcije',
  'Set Stock Prompt': 'Unesite količinu na stanju',
  'Set Stock For All': 'Postavi zalihe za sve',
  'Set Price Prompt': 'Unesite cenu',
  'Set Price For All': 'Postavi cenu za sve',
  ',': ',',
  'Configure Stock And Prices': 'Konfiguriši zalihe i cene',
  'Back To Settings': 'Nazad na podešavanja',
  'Confirm Variants': 'Potvrdi varijante',
  'No Variant Attributes': 'Nema atributa varijanti',
  'Generate Variants': 'Generiši varijante',
  'All Products': 'Svi proizvodi',
  'Active Only': 'Samo aktivni',
  'Inactive Only': 'Samo neaktivni',
  Inventory: 'Inventar',
  'Price Range': 'Opseg cena',
  Min: 'Min',
  Max: 'Maks',
  'Product Name': 'Naziv proizvoda',
  Actions: 'Akcije',
  'No Products Found': 'Proizvodi nisu pronađeni',
  'No Products': 'Nema proizvoda',
  Label: 'Oznaka',
  // Add more Serbian translations as needed...
};

function fixPlaceholders(filePath, isRussian = true) {
  try {
    const content = fs.readFileSync(filePath, 'utf8');
    const json = JSON.parse(content);

    const map = isRussian ? translationMap : serbianMap;
    const prefix = isRussian ? '[RU]' : '[SR]';

    function replaceInObject(obj) {
      for (const key in obj) {
        if (typeof obj[key] === 'string' && obj[key].startsWith(prefix)) {
          // Extract the English phrase from placeholder
          const englishPhrase = obj[key].replace(prefix + ' ', '').trim();

          // Find translation
          if (map[englishPhrase]) {
            obj[key] = map[englishPhrase];
            console.log(
              `  ✓ Replaced: "${englishPhrase}" -> "${map[englishPhrase]}"`
            );
          } else {
            console.log(`  ⚠ No translation for: "${englishPhrase}"`);
          }
        } else if (typeof obj[key] === 'object' && obj[key] !== null) {
          replaceInObject(obj[key]);
        }
      }
    }

    replaceInObject(json);

    fs.writeFileSync(filePath, JSON.stringify(json, null, 2) + '\n', 'utf8');
    console.log(`✅ Fixed: ${filePath}`);
  } catch (error) {
    console.error(`❌ Error processing ${filePath}:`, error.message);
  }
}

// Process Russian files
console.log('Processing Russian translations...');
const ruFiles = [
  'ru/storefronts.json',
  'ru/reviews.json',
  'ru/cart.json',
  'ru/search.json',
  'ru/checkout.json',
  'ru/common.json',
  'ru/orders.json',
];

ruFiles.forEach((file) => {
  const fullPath = path.join(__dirname, '..', 'src', 'messages', file);
  if (fs.existsSync(fullPath)) {
    fixPlaceholders(fullPath, true);
  }
});

// Process Serbian files
console.log('\nProcessing Serbian translations...');
const srFiles = [
  'sr/storefronts.json',
  'sr/reviews.json',
  'sr/cart.json',
  'sr/search.json',
  'sr/checkout.json',
  'sr/common.json',
  'sr/orders.json',
];

srFiles.forEach((file) => {
  const fullPath = path.join(__dirname, '..', 'src', 'messages', file);
  if (fs.existsSync(fullPath)) {
    fixPlaceholders(fullPath, false);
  }
});

console.log('\n✅ Translation fix complete!');
