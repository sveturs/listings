# Исправление функционала лупы для вертикальных изображений

## Проблема

При наведении курсора на вертикальные изображения лупа показывала неправильную область:

- Для горизонтальных фото всё работало корректно
- Для вертикальных фото область увеличения смещалась относительно курсора
- При наведении на пустое пространство слева/справа от вертикального фото лупа всё равно показывала увеличение

## Причина

Расчет позиции курсора производился относительно всего контейнера, а не относительно самого изображения. Для вертикальных изображений с `object-contain` изображение центрируется и не занимает всю ширину контейнера.

## Решение

### Изменения в `src/components/marketplace/listing/ImageGallery.tsx`

```typescript
const handleMouseMove = (e: React.MouseEvent<HTMLDivElement>) => {
  if (!imageRef.current) return;

  const container = imageRef.current;

  // Находим элемент img внутри контейнера
  const img = container.querySelector('img');
  if (!img) return;

  // Получаем реальные размеры изображения
  const imgRect = img.getBoundingClientRect();

  // Проверяем, находится ли курсор над изображением
  const isOverImage =
    e.clientX >= imgRect.left &&
    e.clientX <= imgRect.right &&
    e.clientY >= imgRect.top &&
    e.clientY <= imgRect.bottom;

  if (!isOverImage) {
    setIsZoomed(false);
    return;
  }

  // Рассчитываем позицию относительно изображения, а не контейнера
  const x = ((e.clientX - imgRect.left) / imgRect.width) * 100;
  const y = ((e.clientY - imgRect.top) / imgRect.height) * 100;

  const zoomX = Math.max(0, Math.min(100, x));
  const zoomY = Math.max(0, Math.min(100, y));

  setIsZoomed(true);
  setZoomPosition({ x: zoomX, y: zoomY });
};
```

## Ключевые изменения:

1. Находим реальный элемент `<img>` внутри контейнера
2. Получаем его точные размеры и позицию через `getBoundingClientRect()`
3. Проверяем, находится ли курсор над изображением
4. Рассчитываем позицию относительно изображения, а не контейнера
5. Убрали `onMouseEnter`, так как теперь управляем состоянием в `handleMouseMove`

## Результат

- Лупа корректно работает как для горизонтальных, так и для вертикальных изображений
- При наведении на пустое пространство вокруг вертикальных фото лупа не показывается
- Увеличенная область точно соответствует позиции курсора
