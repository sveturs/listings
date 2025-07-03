# Паспорт API Endpoints: Users (Пользователи)

## 📋 Метаданные
- **Группа API**: Users
- **Базовый путь**: `/api/v1/users`
- **Handler**: `backend/internal/proj/users/handler/routes.go`
- **Количество endpoints**: 3
- **Интеграции**: PostgreSQL, MinIO (аватары)

## 🎯 Назначение
Управление профилями пользователей:
- Получение и обновление личного профиля
- Просмотр публичных профилей других пользователей
- Управление аватарами и контактной информацией
- Настройки приватности

## 📡 Endpoints

### 🔒 Защищенные (требуют авторизации)

#### GET `/api/v1/users/profile`
**Назначение**: Получение собственного профиля пользователя
- **Handler**: `h.User.GetProfile`
- **Security**: Только авторизованный пользователь
- **Response**: Полная информация профиля включая приватные данные
- **Использование**: AuthContext, настройки профиля

#### PUT `/api/v1/users/profile`
**Назначение**: Обновление собственного профиля
- **Handler**: `h.User.UpdateProfile`
- **Security**: Только авторизованный пользователь
- **Body**: UpdateProfileRequest (частичное обновление)
- **Validation**: Email уникальность, форматы полей
- **Effect**: Обновление в БД + переиндексация в OpenSearch

### 🌐 Публичные (с ограничениями)

#### GET `/api/v1/users/:id/profile`
**Назначение**: Получение публичного профиля пользователя по ID
- **Handler**: `h.User.GetProfileByID`
- **Security**: Публичный доступ с фильтрацией приватных данных
- **Response**: Ограниченная информация согласно настройкам приватности
- **Использование**: Карточки продавцов, отзывы, чаты

## 🎭 Структуры данных

### Модель пользователя
```typescript
interface User {
  id: string;
  email: string;
  name: string;
  avatar_url?: string;
  phone?: string;
  bio?: string;
  location?: {
    city: string;
    country: string;
  };
  verification: {
    email_verified: boolean;
    phone_verified: boolean;
    identity_verified: boolean;
  };
  stats: {
    listings_count: number;
    sold_count: number;
    reviews_count: number;
    average_rating: number;
  };
  created_at: string;
  last_active: string;
  role: "user" | "admin";
}
```

### Запросы и ответы
```typescript
interface UpdateProfileRequest {
  name?: string;              // 2-50 символов
  bio?: string;               // до 500 символов
  phone?: string;             // валидный телефон
  location?: {
    city?: string;
    country?: string;
  };
  avatar?: File;              // изображение до 5MB
  privacy_settings?: {
    show_email: boolean;
    show_phone: boolean;
    show_last_active: boolean;
  };
}

interface PrivateProfileResponse {
  user: User;                 // полная информация
  privacy_settings: PrivacySettings;
  contact_preferences: ContactPreferences;
  notification_settings: NotificationSettings;
}

interface PublicProfileResponse {
  user: PublicUser;           // отфильтрованная информация
  stats: PublicStats;
  reviews_summary: ReviewsSummary;
}
```

### Настройки приватности
```typescript
interface PrivacySettings {
  id: string;
  user_id: string;
  show_email: boolean;        // показывать email в профиле
  show_phone: boolean;        // показывать телефон в профиле
  show_last_active: boolean;  // показывать время последней активности
  allow_messages: "everyone" | "contacts" | "none";
  show_listings: boolean;     // показывать активные объявления
  indexable: boolean;         // индексировать профиль в поиске
}

interface ContactPreferences {
  preferred_contact: "email" | "phone" | "chat";
  marketing_emails: boolean;
  transaction_emails: boolean;
  review_reminders: boolean;
}
```

## 🔐 Система приватности

### Фильтрация публичных данных
```typescript
function filterPublicProfile(user: User, privacy: PrivacySettings): PublicUser {
  return {
    id: user.id,
    name: user.name,
    avatar_url: user.avatar_url,
    bio: user.bio,
    location: user.location,
    
    // Условно показываемые поля
    email: privacy.show_email ? user.email : undefined,
    phone: privacy.show_phone ? user.phone : undefined,
    last_active: privacy.show_last_active ? user.last_active : undefined,
    
    // Всегда публичные
    verification: user.verification,
    stats: user.stats,
    created_at: user.created_at,
  };
}
```

### Уровни доступа
- **Собственный профиль**: Полный доступ ко всем данным
- **Контакты**: Расширенная информация (если разрешено)
- **Публичный**: Базовая информация согласно настройкам
- **Анонимный**: Только имя, аватар, рейтинг

## 🔄 Интеграции

### Database Schema
```sql
users (
  id, email, name, avatar_url, phone, bio,
  location_city, location_country,
  email_verified, phone_verified, identity_verified,
  created_at, updated_at, last_active,
  role
);

user_privacy_settings (
  id, user_id, show_email, show_phone, show_last_active,
  allow_messages, show_listings, indexable,
  created_at, updated_at
);

user_contacts (
  user_id, contact_user_id, status, created_at
);
```

### MinIO Integration
- **Bucket**: `users`
- **Avatar Path**: `/avatars/{user_id}/{timestamp}.{ext}`
- **Thumbnails**: 150x150, 300x300
- **Validation**: JPEG/PNG only, max 5MB

### OpenSearch Sync
- Профили индексируются для поиска (при indexable=true)
- Обновление рейтингов из reviews
- Синхронизация статистики объявлений

## 🎛️ Бизнес-логика

### Верификация
```typescript
interface VerificationStatus {
  email_verified: boolean;    // подтверждение email
  phone_verified: boolean;    // SMS верификация
  identity_verified: boolean; // ручная проверка документов
}

// Влияет на доверие и возможности платформы
function getTrustLevel(verification: VerificationStatus): number {
  let trust = 0;
  if (verification.email_verified) trust += 30;
  if (verification.phone_verified) trust += 40;
  if (verification.identity_verified) trust += 30;
  return trust; // 0-100
}
```

### Статистика профиля
```typescript
interface UserStats {
  listings_count: number;      // всего объявлений
  active_listings: number;     // активных сейчас
  sold_count: number;          // продано успешно
  reviews_count: number;       // получено отзывов
  average_rating: number;      // средний рейтинг (1-5)
  response_time: number;       // среднее время ответа в чате (минуты)
  response_rate: number;       // процент ответов на сообщения
  member_since: string;        // дата регистрации
}
```

### Автообновление статистики
- Пересчет после каждой транзакции
- Кеширование на 15 минут
- Фоновое обновление через cron

## 🛡️ Валидация и ограничения

### Поля профиля
```typescript
const VALIDATION_RULES = {
  name: {
    min: 2,
    max: 50,
    pattern: /^[a-zA-Zа-яА-Я\s-']+$/
  },
  bio: {
    max: 500,
    no_html: true
  },
  phone: {
    pattern: /^\+[1-9]\d{1,14}$/,  // E.164 format
    unique: true
  },
  avatar: {
    max_size: 5 * 1024 * 1024,      // 5MB
    types: ['image/jpeg', 'image/png'],
    dimensions: {min: 100, max: 2000}
  }
};
```

### Rate Limiting
- Обновление профиля: 5 раз в час
- Загрузка аватара: 3 раза в час
- Просмотр профилей: 100 в час

## ⚠️ Известные особенности

### Security
- Никогда не возвращаем пароли в API
- Email изменения требуют повторной верификации
- Админы могут видеть расширенную информацию
- Логирование всех изменений профиля

### Performance
- Кеширование публичных профилей на 10 минут
- Lazy loading статистики и отзывов
- Оптимизация запросов к БД через joins

### UX Features
- Автосохранение черновиков профиля
- Предварительный просмотр изменений
- Уведомления о важных изменениях
- История изменений профиля

## 🧪 Примеры использования

### Получение собственного профиля
```bash
curl -X GET /api/v1/users/profile \
  -H "Authorization: Bearer <token>"
```

### Обновление профиля
```bash
curl -X PUT /api/v1/users/profile \
  -H "Authorization: Bearer <token>" \
  -H "Content-Type: application/json" \
  -d '{
    "name": "John Doe",
    "bio": "Loving tech and good deals",
    "location": {"city": "Belgrade", "country": "Serbia"},
    "privacy_settings": {
      "show_email": false,
      "show_phone": true,
      "show_last_active": true
    }
  }'
```

### Просмотр публичного профиля
```bash
curl -X GET /api/v1/users/123/profile
```

### Загрузка аватара
```bash
curl -X PUT /api/v1/users/profile \
  -H "Authorization: Bearer <token>" \
  -F "avatar=@avatar.jpg"
```