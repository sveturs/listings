# Исправление функционала лупы для вертикальных изображений

## Дата: 08.06.2025

## Проблема
При наведении курсора на вертикальные изображения лупа показывала неправильную область:
- Для горизонтальных фото всё работало корректно
- Для вертикальных фото область увеличения смещалась относительно курсора
- При наведении на пустое пространство слева/справа от вертикального фото лупа всё равно показывала увеличение

## Причина
Расчет позиции курсора производился относительно всего контейнера, а не относительно самого изображения. Для вертикальных изображений с `object-contain` изображение центрируется и не занимает всю ширину контейнера.

## Попытки решения

### Попытка 1: Поиск элемента img (НЕ СРАБОТАЛА)
Пытались найти реальный элемент img и использовать его getBoundingClientRect():
```typescript
const img = container.querySelector('img');
if (!img) return;
const imgRect = img.getBoundingClientRect();
```
Результат: Пользователь сообщил "не сработало, совсем".

### Попытка 2: Математический расчет размеров изображения (УСПЕШНАЯ)
Рассчитываем реальные размеры и позицию изображения на основе его естественных размеров и соотношения сторон:

```typescript
// Загружаем естественные размеры изображения
useEffect(() => {
  const currentImage = images[selectedIndex];
  if (!currentImage || currentImage.id === 0 || currentImage.is_video) {
    return;
  }

  const img = new Image();
  img.onload = () => {
    setImageNaturalDimensions({
      width: img.naturalWidth,
      height: img.naturalHeight
    });
  };
  img.src = config.buildImageUrl(currentImage.public_url);
}, [selectedIndex, images]);

// Новая функция handleMouseMove
const handleMouseMove = (e: React.MouseEvent<HTMLDivElement>) => {
  if (!imageRef.current || !imageNaturalDimensions.width || !imageNaturalDimensions.height) return;

  const container = imageRef.current;
  const rect = container.getBoundingClientRect();
  
  // Размеры контейнера
  const containerWidth = rect.width;
  const containerHeight = rect.height;
  
  // Соотношение сторон изображения
  const imageAspectRatio = imageNaturalDimensions.width / imageNaturalDimensions.height;
  // Соотношение сторон контейнера (aspect-[4/3] = 4/3 = 1.333...)
  const containerAspectRatio = 4 / 3;
  
  let renderWidth, renderHeight, offsetX = 0, offsetY = 0;
  
  if (imageAspectRatio > containerAspectRatio) {
    // Изображение шире контейнера - ограничено по ширине
    renderWidth = containerWidth;
    renderHeight = containerWidth / imageAspectRatio;
    offsetY = (containerHeight - renderHeight) / 2;
  } else {
    // Изображение выше контейнера - ограничено по высоте
    renderHeight = containerHeight;
    renderWidth = containerHeight * imageAspectRatio;
    offsetX = (containerWidth - renderWidth) / 2;
  }
  
  // Позиция курсора относительно контейнера
  const mouseX = e.clientX - rect.left;
  const mouseY = e.clientY - rect.top;
  
  // Проверяем, находится ли курсор над изображением
  const isOverImage = 
    mouseX >= offsetX && 
    mouseX <= offsetX + renderWidth && 
    mouseY >= offsetY && 
    mouseY <= offsetY + renderHeight;
  
  if (!isOverImage) {
    setIsZoomed(false);
    return;
  }
  
  // Рассчитываем позицию относительно изображения
  const x = ((mouseX - offsetX) / renderWidth) * 100;
  const y = ((mouseY - offsetY) / renderHeight) * 100;
  
  setIsZoomed(true);
  setZoomPosition({ x, y });
};
```

## Ключевые изменения:
1. Добавлено состояние `imageNaturalDimensions` для хранения естественных размеров изображения
2. Добавлен эффект для загрузки размеров при смене изображения
3. Рассчитываем реальные размеры отображаемого изображения на основе `object-contain`
4. Определяем отступы изображения от краев контейнера
5. Проверяем, находится ли курсор над изображением
6. Рассчитываем позицию относительно изображения, а не контейнера
7. Убрали `onMouseEnter`, так как теперь управляем состоянием в `handleMouseMove`

## Результат
- Лупа корректно работает как для горизонтальных, так и для вертикальных изображений
- При наведении на пустое пространство вокруг вертикальных фото лупа не показывается
- Увеличенная область точно соответствует позиции курсора

## Файлы изменены:
- `/frontend/svetu/src/components/marketplace/listing/ImageGallery.tsx`

## Тестирование
Создана тестовая HTML страница `/frontend/svetu/test-zoom.html` для проверки алгоритма расчета позиции на вертикальных и горизонтальных изображениях.