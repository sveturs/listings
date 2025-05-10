// frontend/hostel-frontend/src/utils/adminUtils.ts
import axios from '../api/axios';

// Интерфейс для ответа API проверки админа
interface AdminCheckResponse {
    is_admin: boolean;
    [key: string]: any;
}

// Кеш для хранения результатов проверки админских прав
const adminStatusCache: Map<string, boolean> = new Map();

// Проверка админа через API с кешированием
export const checkAdminStatus = async (email: string | null | undefined): Promise<boolean> => {
    if (!email) return false;

    // Проверяем кеш сначала
    if (adminStatusCache.has(email)) {
        return adminStatusCache.get(email) || false;
    }

    try {
        const response = await axios.get<AdminCheckResponse>(`/api/v1/admin/admins/check/${email}`);
        const isAdmin = response.data.is_admin;

        // Сохраняем результат в кеш
        adminStatusCache.set(email, isAdmin);

        return isAdmin;
    } catch (error) {
        console.error('Error checking admin status:', error);
        return false;
    }
};

// Синхронная функция для быстрой первичной проверки - используется только для UI
// ВАЖНО: Это не является авторизующей функцией, а только для быстрой проверки в UI
export const isAdmin = (email: string | null | undefined): boolean => {
    if (!email) return false;

    // Проверяем кеш сначала
    if (adminStatusCache.has(email)) {
        return adminStatusCache.get(email) || false;
    }

    // Для первого рендера, можем использовать env переменную, если она есть
    if (process.env.REACT_APP_ADMIN_EMAILS) {
        return process.env.REACT_APP_ADMIN_EMAILS.split(',').includes(email);
    }

    // По умолчанию возвращаем false, реальная проверка произойдет асинхронно
    return false;
};