// frontend/hostel-frontend/src/utils/dateUtils.ts

// Интерфейс для опций форматирования даты
interface DateFormatOptions extends Intl.DateTimeFormatOptions {
  [key: string]: any;
}

/**
 * Форматирует дату в локализованную строку
 * @param date - Дата для форматирования
 * @param options - Опции форматирования
 * @returns Отформатированная дата
 */
export const formatDate = (date: string | Date | null | undefined, options: DateFormatOptions = {}): string => {
  if (!date) return '-';
  
  const defaultOptions: DateFormatOptions = {
    year: 'numeric',
    month: 'short',
    day: 'numeric',
    hour: '2-digit',
    minute: '2-digit'
  };
  
  const mergedOptions: DateFormatOptions = { ...defaultOptions, ...options };
  
  try {
    const dateObj = typeof date === 'string' ? new Date(date) : date;
    return dateObj.toLocaleDateString(undefined, mergedOptions);
  } catch (error) {
    console.error('Error formatting date:', error);
    return String(date);
  }
};

/**
 * Форматирует дату в относительную строку (например, "2 дня назад")
 * @param date - Дата для форматирования
 * @returns Относительная дата
 */
export const formatRelativeDate = (date: string | Date | null | undefined): string => {
  if (!date) return '-';
  
  try {
    const dateObj = typeof date === 'string' ? new Date(date) : date;
    const now = new Date();
    const diffInSeconds = Math.floor((now.getTime() - dateObj.getTime()) / 1000);
    
    if (diffInSeconds < 60) {
      return 'только что';
    }
    
    const diffInMinutes = Math.floor(diffInSeconds / 60);
    if (diffInMinutes < 60) {
      return `${diffInMinutes} ${pluralize(diffInMinutes, 'минуту', 'минуты', 'минут')} назад`;
    }
    
    const diffInHours = Math.floor(diffInMinutes / 60);
    if (diffInHours < 24) {
      return `${diffInHours} ${pluralize(diffInHours, 'час', 'часа', 'часов')} назад`;
    }
    
    const diffInDays = Math.floor(diffInHours / 24);
    if (diffInDays < 30) {
      return `${diffInDays} ${pluralize(diffInDays, 'день', 'дня', 'дней')} назад`;
    }
    
    const diffInMonths = Math.floor(diffInDays / 30);
    if (diffInMonths < 12) {
      return `${diffInMonths} ${pluralize(diffInMonths, 'месяц', 'месяца', 'месяцев')} назад`;
    }
    
    const diffInYears = Math.floor(diffInMonths / 12);
    return `${diffInYears} ${pluralize(diffInYears, 'год', 'года', 'лет')} назад`;
  } catch (error) {
    console.error('Error formatting relative date:', error);
    return String(date);
  }
};

/**
 * Вспомогательная функция для правильного склонения слов
 * @param count - Количество
 * @param one - Форма для 1
 * @param few - Форма для 2-4
 * @param many - Форма для 5-20
 * @returns Правильная форма слова
 */
const pluralize = (count: number, one: string, few: string, many: string): string => {
  if (count % 10 === 1 && count % 100 !== 11) {
    return one;
  }
  if ([2, 3, 4].includes(count % 10) && ![12, 13, 14].includes(count % 100)) {
    return few;
  }
  return many;
};