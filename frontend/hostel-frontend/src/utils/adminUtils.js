// frontend/hostel-frontend/src/utils/adminUtils.js
import axios from '../api/axios';

// Кеш для хранения результатов проверки админских прав
const adminStatusCache = new Map();

// Проверка админа через API с кешированием
export const checkAdminStatus = async (email) => {
    if (!email) return false;

    // Проверяем кеш сначала
    if (adminStatusCache.has(email)) {
        return adminStatusCache.get(email);
    }

    try {
        const response = await axios.get(`/api/v1/admin/admins/check/${email}`);
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
export const isAdmin = (email) => {
    if (!email) return false;

    // Проверяем кеш сначала
    if (adminStatusCache.has(email)) {
        return adminStatusCache.get(email);
    }

    // Для первого рендера, можем использовать env переменную, если она есть
    if (process.env.REACT_APP_ADMIN_EMAILS) {
        return process.env.REACT_APP_ADMIN_EMAILS.split(',').includes(email);
    }

    // По умолчанию возвращаем false, реальная проверка произойдет асинхронно
    return false;
};