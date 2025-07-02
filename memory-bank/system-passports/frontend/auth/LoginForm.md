# Паспорт компонента: LoginForm

## 📋 Метаданные
- **Путь**: `/frontend/svetu/src/components/auth/LoginForm.tsx`
- **Роль**: Форма входа в систему
- **Тип**: Presentational Component (React.memo)
- **Размер**: 122 строки

## 🎯 Назначение
Отображение формы входа с полями email/пароль и альтернативным входом через Google OAuth. Чистый presentational компонент без собственной логики.

## 🔧 Структура Props
```typescript
interface LoginFormProps {
  formData: FormData;         // Данные формы
  errors: FormErrors;         // Ошибки валидации
  isLoading: boolean;         // Состояние загрузки
  onFieldChange: (field: keyof FormData, value: string) => void;
  onSubmit: (e: React.FormEvent) => void;
  onGoogleLogin: () => void;
  onSwitchToRegister: () => void;
  canSubmit: boolean;         // Флаг возможности отправки
}
```

## 🔗 Зависимости

### Компоненты:
- `FormField` - обертка для полей формы
- `GoogleIcon` - иконка Google

### Интернационализация:
- `next-intl` - переводы

## 🎨 UI структура

### 1. Email поле:
```tsx
<FormField label={t('loginForm.email')} required error={errors.email}>
  <input type="email" 
         className="input input-bordered w-full"
         autoComplete="email" />
</FormField>
```

### 2. Password поле:
```tsx
<FormField label={t('loginForm.password')} required error={errors.password}>
  <input type="password" 
         className="input input-bordered w-full"
         autoComplete="current-password"
         minLength={6} />
</FormField>
```

### 3. Общие ошибки:
```tsx
{errors.general && (
  <div className="alert alert-error">
    <span>{t(errors.general)}</span>
  </div>
)}
```

### 4. Кнопки действий:
- Submit button с loading состоянием
- Google OAuth button
- Ссылка на регистрацию

## 🎨 Стилизация
- DaisyUI классы для всех элементов
- Условные классы для ошибок (`input-error`)
- Loading состояние кнопки
- Responsive layout с `w-full`

## ⚡ Особенности

### 1. Memoization:
```typescript
export default React.memo(LoginForm);
```
Предотвращает лишние рендеры при изменении родителя

### 2. Accessibility:
- `autoComplete` атрибуты для полей
- `required` атрибуты
- `minLength` для пароля
- `type="email"` для валидации

### 3. Disabled состояния:
- Все поля и кнопки блокируются при `isLoading`
- Submit блокируется также при `!canSubmit`

## 🌍 Интернационализация

Используемые ключи:
- `auth.loginForm.email`
- `auth.loginForm.emailPlaceholder`
- `auth.loginForm.password`
- `auth.loginForm.passwordPlaceholder`
- `auth.loginForm.submit`
- `auth.loginForm.loggingIn`
- `auth.loginForm.or`
- `auth.loginForm.googleLogin`
- `auth.loginForm.registerText`
- `auth.loginForm.register`

## 📝 Примеры использования

```tsx
<LoginForm
  formData={formData}
  errors={errors}
  isLoading={false}
  onFieldChange={handleFieldChange}
  onSubmit={handleSubmit}
  onGoogleLogin={handleGoogleLogin}
  onSwitchToRegister={switchToRegister}
  canSubmit={canSubmit}
/>
```

## ⚠️ Известные особенности

1. **Presentational only**: 
   - Нет внутреннего состояния
   - Вся логика в родительском компоненте

2. **Password минимум**:
   - HTML5 валидация `minLength={6}`
   - Должна совпадать с серверной

3. **Error display**:
   - Ошибки полей через FormField
   - Общие ошибки в отдельном alert

## 🔄 Связанные компоненты
- `LoginModal` - использует эту форму
- `RegisterForm` - альтернативная форма
- `FormField` - обертка полей
- `GoogleIcon` - иконка OAuth