// src/utils/webPush.ts

/**
 * Конвертирует URL-safe base64 строку в Uint8Array
 * @param base64String - URL-safe base64 строка для конвертации
 * @returns Uint8Array, созданный из base64 строки
 */
export function urlBase64ToUint8Array(base64String: string): Uint8Array {
    const padding = '='.repeat((4 - base64String.length % 4) % 4);
    const base64 = (base64String + padding)
        .replace(/\-/g, '+')
        .replace(/_/g, '/');
    const rawData = window.atob(base64);
    // Используем Array.from вместо оператора spread для совместимости с ES5
    return Uint8Array.from(Array.from(rawData).map((char) => char.charCodeAt(0)));
}