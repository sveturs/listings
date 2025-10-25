# 🎯 ДЕТАЛЬНЫЙ ПЛАН УСТРАНЕНИЯ ПРОБЛЕМ ДОСТУПНОСТИ

**Дата аудита:** 2025-10-20
**Статус:** В РАБОТЕ
**Приоритет:** КРИТИЧЕСКИЙ (блокирует WCAG 2.1 AA compliance)

---

## 📊 EXECUTIVE SUMMARY

**Всего найдено проблем:** 27+
**WCAG violations:** Level A и AA
**Затронутых компонентов:** 10+
**Оценка времени:** 4-6 часов

### Категории проблем:
1. ❌ **user-scalable=no** (CRITICAL) - 1 issue
2. ❌ **Кнопки без aria-label** (HIGH) - 18+ issues
3. ❌ **title вместо aria-label** (MEDIUM) - 8+ issues
4. ⚠️ **Отсутствие aria-expanded** (MEDIUM) - 5+ issues
5. ⚠️ **Skeleton без role** (LOW-MEDIUM) - 2+ issues

---

## 🔥 PHASE 1: КРИТИЧНЫЕ ИСПРАВЛЕНИЯ (2 часа)

### ✅ TASK 1.1: Исправить viewport (CRITICAL - 5 минут)

**Файл:** `src/app/[locale]/layout.tsx:38-43`

**Текущий код:**
```typescript
export const viewport = {
  width: 'device-width',
  initialScale: 1,
  maximumScale: 1,
  userScalable: false,
};
```

**Исправленный код:**
```typescript
export const viewport = {
  width: 'device-width',
  initialScale: 1,
  // WCAG 2.1 AA: Reflow (1.4.10) - Allow user scaling
  maximumScale: 5,  // Allow up to 5x zoom
  userScalable: true,
};
```

**WCAG:** Fixes 1.4.4 (Resize Text) Level AA
**Тест:** После изменения все accessibility тесты должны показать на 1 critical ошибку меньше

---

### ✅ TASK 1.2: ChatIcon - Добавить aria-label (HIGH - 3 минуты)

**Файл:** `src/components/icons/ChatIcon.tsx:14-23`

**Текущий код:**
```tsx
<Link
  href={`/${locale}/chat`}
  className="btn btn-ghost btn-circle"
>
  <div className="w-5 h-5">
    {/* SVG icon */}
  </div>
</Link>
```

**Исправленный код:**
```tsx
<Link
  href={`/${locale}/chat`}
  className="btn btn-ghost btn-circle"
  aria-label={t('chat.openChat')}  // Add translation key
>
  <div className="w-5 h-5" aria-hidden="true">
    {/* SVG icon */}
  </div>
</Link>
```

**Добавить переводы:**
- `en/common.json`: `"chat": {"openChat": "Open chat"}`
- `ru/common.json`: `"chat": {"openChat": "Открыть чат"}`
- `sr/common.json`: `"chat": {"openChat": "Otvor ćaskanje"}`

**WCAG:** Fixes 4.1.2 (Name, Role, Value) Level A

---

### ✅ TASK 1.3: Header Mobile Menu - aria-label + aria-expanded (HIGH - 5 минут)

**Файл:** `src/components/Header.tsx:285-294` (примерная локация)

**Найти код:**
```tsx
<button className="btn btn-square btn-ghost lg:hidden">
  <svg>...</svg>
</button>
```

**Исправить на:**
```tsx
<button
  className="btn btn-square btn-ghost lg:hidden"
  aria-label={t('navigation.toggleMenu')}
  aria-expanded={isMenuOpen}
  aria-controls="mobile-menu"
  onClick={() => setIsMenuOpen(!isMenuOpen)}
>
  <svg aria-hidden="true">...</svg>
</button>

{/* Add id to menu */}
<div id="mobile-menu" className={...}>
```

**Переводы:**
- `"navigation": {"toggleMenu": "Toggle navigation menu"}`

**WCAG:** Fixes 4.1.2 (Name, Role, Value) + 4.1.3 (Status Messages)

---

### ✅ TASK 1.4: Header Close Button - aria-label (HIGH - 2 минуты)

**Файл:** `src/components/Header.tsx:336-341`

**Текущий код:**
```tsx
<button className="btn btn-square btn-ghost">
  <svg>...</svg>  {/* X icon */}
</button>
```

**Исправленный:**
```tsx
<button
  className="btn btn-square btn-ghost"
  aria-label={t('navigation.closeMenu')}
  onClick={() => setIsMenuOpen(false)}
>
  <svg aria-hidden="true">...</svg>
</button>
```

---

### ✅ TASK 1.5: ThemeToggle Skeleton - aria-label (MEDIUM - 3 минуты)

**Файл:** `src/components/ThemeToggle.tsx:36-39`

**Текущий код:**
```tsx
<div className="skeleton h-10 w-24" aria-label="Loading authentication status"></div>
```

**Исправленный:**
```tsx
<div
  className="skeleton h-10 w-24"
  role="status"
  aria-live="polite"
  aria-label={t('theme.loadingToggle')}
>
  <span className="sr-only">{t('common.loading')}</span>
</div>
```

**WCAG:** Fixes aria-prohibited-attr (SERIOUS)

---

## 🎨 PHASE 2: CAROUSEL & IMAGE NAVIGATION (1.5 часа)

### ✅ TASK 2.1: QuickView - Previous/Next Image Buttons (HIGH - 10 минут)

**Файл:** `src/components/QuickView.tsx:189-199`

**Исправления:**

```tsx
{/* Previous Button */}
<button
  onClick={handlePrevImage}
  className="..."
  aria-label={t('product.previousImage')}
  disabled={currentImageIndex === 0}
>
  <ChevronLeftIcon className="w-6 h-6" aria-hidden="true" />
</button>

{/* Next Button */}
<button
  onClick={handleNextImage}
  className="..."
  aria-label={t('product.nextImage')}
  disabled={currentImageIndex === product.images.length - 1}
>
  <ChevronRightIcon className="w-6 h-6" aria-hidden="true" />
</button>
```

**Переводы:**
```json
"product": {
  "previousImage": "Previous image",
  "nextImage": "Next image",
  "imageXofY": "Image {current} of {total}"
}
```

---

### ✅ TASK 2.2: QuickView - Thumbnail Buttons (MEDIUM - 8 минут)

**Файл:** `src/components/QuickView.tsx:218-239`

**Исправить:**
```tsx
{product.images.map((image, index) => (
  <button
    key={index}
    onClick={() => setCurrentImageIndex(index)}
    className={...}
    aria-label={t('product.imageXofY', {
      current: index + 1,
      total: product.images.length
    })}
    aria-current={currentImageIndex === index ? 'true' : 'false'}
  >
    <Image
      src={image}
      alt=""  // Empty alt since aria-label on button
      className="..."
    />
  </button>
))}
```

---

### ✅ TASK 2.3: QuickView - Close Button (HIGH - 3 минуты)

**Файл:** `src/components/QuickView.tsx:157-162`

```tsx
<button
  onClick={onClose}
  className="..."
  aria-label={t('common.close')}
>
  <XMarkIcon className="w-6 h-6" aria-hidden="true" />
</button>
```

---

### ✅ TASK 2.4: QuickView - Action Buttons (MEDIUM - 5 минут)

**Файл:** `src/components/QuickView.tsx:356-361`

```tsx
{/* Favorite Button */}
<button
  onClick={handleToggleFavorite}
  aria-label={isFavorite ? t('product.removeFromFavorites') : t('product.addToFavorites')}
  aria-pressed={isFavorite}
>
  <HeartIcon aria-hidden="true" />
</button>

{/* Share Button */}
<button
  onClick={handleShare}
  aria-label={t('product.shareProduct')}
>
  <ShareIcon aria-hidden="true" />
</button>
```

---

## 🛒 PHASE 3: PRODUCT CARDS (1 час)

### ✅ TASK 3.1: UnifiedProductCard - Grid View Buttons (HIGH - 15 минут)

**Файл:** `src/components/UnifiedProductCard.tsx:769-784`

**Заменить `title` на `aria-label`:**

```tsx
{/* Quick View - Grid */}
<button
  onClick={handleQuickView}
  className="..."
  aria-label={t('product.quickView')}  // Remove title
>
  <EyeIcon className="w-5 h-5" aria-hidden="true" />
</button>

{/* Favorite - Grid */}
<button
  onClick={handleToggleFavorite}
  className="..."
  aria-label={isFavorite ? t('product.removeFromFavorites') : t('product.addToFavorites')}
  aria-pressed={isFavorite}
>
  <HeartIcon className="w-5 h-5" aria-hidden="true" />
</button>
```

---

### ✅ TASK 3.2: UnifiedProductCard - List View Buttons (HIGH - 10 минут)

**Файл:** `src/components/UnifiedProductCard.tsx:540-600`

Аналогичные изменения для list view.

---

## 🎛️ PHASE 4: VIEW TOGGLES & CONTROLS (45 минут)

### ✅ TASK 4.1: GridColumnsToggle - Replace title with aria-label (MEDIUM - 8 минут)

**Файл:** `src/components/GridColumnsToggle.tsx`

**Найти все:**
```tsx
title="3 columns"
```

**Заменить на:**
```tsx
aria-label={t('view.threeColumns')}
aria-current={columns === 3 ? 'true' : 'false'}
```

**Добавить переводы:**
```json
"view": {
  "oneColumn": "One column",
  "twoColumns": "Two columns",
  "threeColumns": "Three columns",
  "fourColumns": "Four columns"
}
```

---

### ✅ TASK 4.2: ViewToggle - Replace title with aria-label (MEDIUM - 8 минут)

**Файл:** `src/components/ViewToggle.tsx`

```tsx
aria-label={t('view.gridView')}
aria-current={view === 'grid' ? 'true' : 'false'}

aria-label={t('view.listView')}
aria-current={view === 'list' ? 'true' : 'false'}
```

**Переводы:**
```json
"view": {
  "gridView": "Grid view",
  "listView": "List view"
}
```

---

## 🗂️ PHASE 5: MODALS & EXPANDABLE (45 минут)

### ✅ TASK 5.1: CategoryTreeModal - Close Button (MEDIUM - 5 минут)

**Файл:** `src/components/CategoryTreeModal.tsx:204-206`

```tsx
<button
  onClick={onClose}
  aria-label={t('categories.closeSelection')}
>
  <XMarkIcon aria-hidden="true" />
</button>
```

---

### ✅ TASK 5.2: CategoryTreeModal - Clear Search (MEDIUM - 5 минут)

**Файл:** `src/components/CategoryTreeModal.tsx:221-227`

```tsx
<button
  onClick={() => setSearchQuery('')}
  aria-label={t('common.clearSearch')}
>
  <XMarkIcon aria-hidden="true" />
</button>
```

---

### ✅ TASK 5.3: CategoryTreeModal - Expand/Collapse (HIGH - 12 минут)

**Файл:** `src/components/CategoryTreeModal.tsx:163-175`

```tsx
<button
  onClick={() => toggleExpand(item.id)}
  aria-label={t('categories.toggleCategory', { name: item.name })}
  aria-expanded={isExpanded}
  aria-controls={`category-${item.id}-children`}
>
  {isExpanded ? (
    <ChevronDownIcon aria-hidden="true" />
  ) : (
    <ChevronRightIcon aria-hidden="true" />
  )}
</button>

{/* Add id to children container */}
{isExpanded && (
  <div id={`category-${item.id}-children`}>
    {/* Children */}
  </div>
)}
```

**Переводы:**
```json
"categories": {
  "toggleCategory": "Toggle {name} category",
  "closeSelection": "Close category selection"
}
```

---

## 📱 PHASE 6: FLOATING ACTION BUTTONS (30 минут)

### ✅ TASK 6.1: FloatingActionButtons - Replace title with aria-label (MEDIUM - 15 минут)

**Файл:** `src/components/GIS/Mobile/FloatingActionButtons.tsx`

**Все кнопки:**

```tsx
{/* Main FAB */}
<button
  aria-label={t('map.fabMenu')}
  aria-expanded={isExpanded}
  aria-controls="fab-menu"
>

{/* Filters */}
<button
  aria-label={t('map.openFilters')}  // Not title
>

{/* Geolocation */}
<button
  aria-label={t('map.findMyLocation')}  // Not title
>

{/* Show All */}
<button
  aria-label={t('map.showAllListings')}  // Not title
>
```

---

## 🎭 PHASE 7: LANGUAGE SWITCHER (15 минут)

### ✅ TASK 7.1: LanguageSwitcher - Complete ARIA (HIGH - 15 минут)

**Файл:** `src/components/LanguageSwitcher.tsx:25-40`

```tsx
<button
  onClick={() => setIsOpen(!isOpen)}
  aria-label={t('language.switchLanguage')}
  aria-expanded={isOpen}
  aria-haspopup="listbox"
  aria-controls="language-menu"
>
  {/* Current language */}
</button>

<div
  id="language-menu"
  role="listbox"
  aria-label={t('language.selectLanguage')}
>
  {locales.map((loc) => (
    <button
      role="option"
      aria-selected={locale === loc}
      onClick={() => handleChange(loc)}
    >
      {/* Language option */}
    </button>
  ))}
</div>
```

**Переводы:**
```json
"language": {
  "switchLanguage": "Switch language",
  "selectLanguage": "Select language",
  "currentLanguage": "Current language: {lang}"
}
```

---

## 🧪 PHASE 8: ТЕСТИРОВАНИЕ И ВАЛИДАЦИЯ (1 час)

### ✅ TASK 8.1: Запустить accessibility тесты

```bash
cd /data/hostel-booking-system/frontend/svetu
yarn playwright test e2e/axe/ --project=chromium
```

**Ожидаемый результат:**
- ✅ Все 12 accessibility тестов должны пройти
- ✅ 0 WCAG violations (было 7)
- ✅ 0 critical issues (было 3)

---

### ✅ TASK 8.2: Ручное тестирование с клавиатуры

**Checklist:**
- [ ] Tab через все интерактивные элементы на homepage
- [ ] Enter/Space активируют кнопки
- [ ] Escape закрывает модальные окна
- [ ] Arrow keys работают в carousel
- [ ] Focus indicators видны на всех элементах

---

### ✅ TASK 8.3: Тестирование со screen reader

**Tools:** NVDA (Windows) или VoiceOver (Mac)

**Checklist:**
- [ ] Все кнопки имеют понятные названия
- [ ] Expanded/collapsed состояния объявляются
- [ ] Модальные окна правильно фокусируются
- [ ] Навигация интуитивна

---

## 📋 IMPLEMENTATION CHECKLIST

### Pre-Implementation:
- [ ] Создать feature branch: `fix/accessibility-wcag-compliance`
- [ ] Убедиться что все E2E тесты проходят (baseline)
- [ ] Создать backup текущего состояния

### Implementation Order:
- [ ] Phase 1: Критичные исправления (2ч)
- [ ] Phase 2: Carousel & Navigation (1.5ч)
- [ ] Phase 3: Product Cards (1ч)
- [ ] Phase 4: View Toggles (45мин)
- [ ] Phase 5: Modals (45мин)
- [ ] Phase 6: FABs (30мин)
- [ ] Phase 7: Language Switcher (15мин)
- [ ] Phase 8: Тестирование (1ч)

### Post-Implementation:
- [ ] Запустить все E2E тесты
- [ ] Запустить accessibility тесты
- [ ] Ручное тестирование
- [ ] Code review
- [ ] Создать PR с подробным описанием
- [ ] Обновить документацию

---

## 🎯 SUCCESS CRITERIA

### Обязательные:
✅ Все 12 accessibility тестов проходят
✅ 0 WCAG 2.1 Level A violations
✅ 0 WCAG 2.1 Level AA violations
✅ Viewport позволяет zoom до 200%
✅ Все кнопки имеют aria-label или visible text

### Желательные:
✅ Полная keyboard navigation
✅ Screen reader friendly
✅ Focus management в модальных окнах
✅ Консистентное использование ARIA атрибутов

---

## 📝 TRANSLATION KEYS SUMMARY

### Новые ключи для добавления:

**common.json:**
```json
{
  "loading": "Loading...",
  "close": "Close",
  "clearSearch": "Clear search"
}
```

**navigation.json:**
```json
{
  "toggleMenu": "Toggle navigation menu",
  "closeMenu": "Close menu"
}
```

**theme.json:**
```json
{
  "loadingToggle": "Loading theme toggle"
}
```

**product.json:**
```json
{
  "previousImage": "Previous image",
  "nextImage": "Next image",
  "imageXofY": "Image {current} of {total}",
  "quickView": "Quick view",
  "addToFavorites": "Add to favorites",
  "removeFromFavorites": "Remove from favorites",
  "shareProduct": "Share product"
}
```

**view.json:**
```json
{
  "gridView": "Grid view",
  "listView": "List view",
  "oneColumn": "One column",
  "twoColumns": "Two columns",
  "threeColumns": "Three columns",
  "fourColumns": "Four columns"
}
```

**categories.json:**
```json
{
  "toggleCategory": "Toggle {name} category",
  "closeSelection": "Close category selection"
}
```

**map.json:**
```json
{
  "fabMenu": "Map actions menu",
  "openFilters": "Open filters",
  "findMyLocation": "Find my location",
  "showAllListings": "Show all listings"
}
```

**language.json:**
```json
{
  "switchLanguage": "Switch language",
  "selectLanguage": "Select language",
  "currentLanguage": "Current language: {lang}"
}
```

**chat.json:**
```json
{
  "openChat": "Open chat"
}
```

---

## ⚠️ РИСКИ И МИТИГАЦИЯ

### Риск 1: Изменение viewport может сломать layout на мобильных
**Митигация:** Тестировать на реальных устройствах, использовать `maximumScale: 5` вместо удаления ограничения

### Риск 2: Слишком много aria-label может запутать screen reader users
**Митигация:** Использовать краткие, понятные labels; тестировать со screen readers

### Риск 3: Translation keys могут отсутствовать
**Митигация:** Добавить fallback в компонентах, проверить все локали

### Риск 4: Регрессия в существующих тестах
**Митигация:** Запускать полный test suite после каждой phase

---

## 📚 REFERENCES

- [WCAG 2.1 Guidelines](https://www.w3.org/WAI/WCAG21/quickref/)
- [ARIA Authoring Practices](https://www.w3.org/WAI/ARIA/apg/)
- [axe-core Rules](https://github.com/dequelabs/axe-core/blob/develop/doc/rule-descriptions.md)
- [MDN ARIA](https://developer.mozilla.org/en-US/docs/Web/Accessibility/ARIA)

---

**Создано:** 2025-10-20
**Автор:** AI Assistant (Claude)
**Статус:** READY FOR IMPLEMENTATION
