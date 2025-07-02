# Паспорт компонента: RegisterForm

## 📋 Метаданные
- **Путь**: `/frontend/svetu/src/components/auth/RegisterForm.tsx`
- **Роль**: Форма регистрации нового пользователя
- **Тип**: Presentational Component (React.memo)
- **Размер**: 223 строки

## 🎯 Назначение
Отображение формы регистрации с полями имени, email, телефона (опционально), пароля и подтверждения пароля. Поддерживает регистрацию через Google OAuth и показ успешного сообщения после регистрации.

## 🔧 Структура Props
```typescript
interface RegisterFormProps {
  formData: FormData;           // Данные формы
  errors: FormErrors;           // Ошибки валидации
  isLoading: boolean;           // Состояние загрузки
  successMessage: string;       // Сообщение успеха
  onFieldChange: (field: keyof FormData, value: string) => void;
  onSubmit: (e: React.FormEvent) => void;
  onGoogleLogin: () => void;
  onSwitchToLogin: () => void;
  onClose: () => void;          // Закрытие модала
  canSubmit: boolean;           // Флаг возможности отправки
}
```

## 🔗 Зависимости

### Компоненты:
- `FormField` - обертка для полей формы
- `GoogleIcon` - иконка Google

### Интернационализация:
- `next-intl` - переводы

## 🎨 UI структура

### 1. Success View (после успешной регистрации):
```tsx
<div className="text-center space-y-6">
  <div className="alert alert-success">
    <svg>✓</svg>
    <span>{successMessage}</span>
  </div>
  <p>{successDescription}</p>
  <buttons: switchToLogin | close>
</div>
```

### 2. Form View:
- **Name field** - обязательное, min 2 символа
- **Email field** - обязательное, email тип
- **Phone field** - опциональное, tel тип
- **Password field** - обязательное, min 6 символов
- **Confirm Password** - обязательное, должно совпадать
- **Submit button** - с loading состоянием
- **Google OAuth** - альтернативная регистрация
- **Switch to login** - ссылка на форму входа

## 📋 Поля формы

### Обязательные:
1. **Name** - `minLength={2}`, `autoComplete="given-name"`
2. **Email** - `type="email"`, `autoComplete="email"`
3. **Password** - `minLength={6}`, `autoComplete="new-password"`
4. **Confirm Password** - `minLength={6}`, `autoComplete="new-password"`

### Опциональные:
1. **Phone** - `type="tel"`, `autoComplete="tel"`

## 🎨 Стилизация
- DaisyUI классы для всех элементов
- Условные классы ошибок (`input-error`)
- Alert компоненты для сообщений
- Loading анимация на кнопке
- Responsive layout с `w-full`

## ⚡ Особенности

### 1. Два режима отображения:
- Form view - основная форма регистрации
- Success view - после успешной регистрации

### 2. Memoization:
```typescript
export default React.memo(RegisterForm);
```

### 3. Accessibility:
- Правильные `autoComplete` атрибуты
- `required` для обязательных полей
- `minLength` валидация
- Семантичные типы input

### 4. Success flow:
- Показывает успешное сообщение
- Предлагает перейти к входу или закрыть

## 🌍 Интернационализация

Используемые ключи:
- `auth.registerForm.name/namePlaceholder`
- `auth.registerForm.email/emailPlaceholder`
- `auth.registerForm.phone/phonePlaceholder`
- `auth.registerForm.password/passwordPlaceholder`
- `auth.registerForm.confirmPassword/confirmPasswordPlaceholder`
- `auth.registerForm.submit/submitting`
- `auth.registerForm.or`
- `auth.registerForm.googleRegister`
- `auth.registerForm.loginText/switchToLogin`
- `auth.registerForm.successMessage/successDescription`
- `auth.registerForm.close`

## 📝 Примеры использования

```tsx
<RegisterForm
  formData={formData}
  errors={errors}
  isLoading={false}
  successMessage=""
  onFieldChange={handleFieldChange}
  onSubmit={handleSubmit}
  onGoogleLogin={handleGoogleLogin}
  onSwitchToLogin={switchToLogin}
  onClose={closeModal}
  canSubmit={canSubmit}
/>
```

## ⚠️ Известные особенности

1. **Success view**:
   - Полностью заменяет форму после успеха
   - Не позволяет вернуться к форме

2. **Phone опционален**:
   - Единственное необязательное поле
   - Нет `required` атрибута

3. **Password matching**:
   - Валидация confirmPassword в родителе
   - HTML валидация только на длину

4. **Presentational**:
   - Нет внутреннего состояния
   - Вся логика снаружи

## 🔄 Связанные компоненты
- `LoginModal` - использует эту форму
- `LoginForm` - альтернативная форма
- `FormField` - обертка полей
- `GoogleIcon` - иконка OAuth