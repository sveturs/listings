/**
 * Добавляет водяной знак к изображению
 * @param image - Исходное изображение (Blob или File)
 * @param watermarkText - Текст водяного знака
 * @returns Promise с Blob изображения с водяным знаком
 */
export const addWatermark = async (image: Blob | File, watermarkText: string = 'SveTu.rs'): Promise<Blob> => {
    return new Promise((resolve) => {
        const canvas = document.createElement('canvas');
        const ctx = canvas.getContext('2d');
        
        if (!ctx) {
            // Если контекст не удалось получить, возвращаем исходное изображение
            resolve(image);
            return;
        }

        const img = new Image();
        img.onload = () => {
            canvas.width = img.width;
            canvas.height = img.height;

            // Рисуем исходное изображение
            ctx.drawImage(img, 0, 0, img.width, img.height);

            // Настраиваем стиль водяного знака
            const minDimension = Math.min(img.width, img.height);
            const fontSize = minDimension * 0.03; // Уменьшаем размер до 3% от размера изображения
            ctx.font = `300 ${fontSize}px SF Pro Display, -apple-system, BlinkMacSystemFont, Segoe UI, Roboto, sans-serif`;
            
            // Задаем текст и позицию
            const padding = fontSize; // Отступ от края
            const text = watermarkText;
            const metrics = ctx.measureText(text);
            const textWidth = metrics.width;

            // Создаем полупрозрачный фон для текста
            ctx.fillStyle = 'rgba(0, 0, 0, 0.3)';
            const backgroundHeight = fontSize * 1.4;
            ctx.fillRect(
                padding - fontSize * 0.2,
                canvas.height - padding - backgroundHeight,
                textWidth + fontSize * 0.4,
                backgroundHeight
            );

            // Рисуем текст
            ctx.fillStyle = 'rgba(255, 255, 255, 0.85)';
            ctx.textBaseline = 'middle';
            ctx.fillText(
                text,
                padding,
                canvas.height - padding - backgroundHeight/2
            );

            // Конвертируем canvas в blob
            canvas.toBlob((blob) => {
                if (blob) {
                    resolve(blob);
                } else {
                    // В случае ошибки возвращаем исходное изображение
                    resolve(image);
                }
            }, 'image/jpeg', 0.95);
        };

        img.onerror = () => {
            // В случае ошибки загрузки возвращаем исходное изображение
            resolve(image);
        };

        img.src = URL.createObjectURL(image);
    });
};