import { useEffect, useRef } from 'react';
import { useAppDispatch, useAppSelector } from '@/store/hooks';
import { useAuth } from '@/contexts/AuthContext';
import { fetchUserCarts, addToCart } from '@/store/slices/cartSlice';
import { selectLocalCart, clearLocalCart } from '@/store/slices/localCartSlice';

/**
 * Хук для синхронизации корзин пользователя при логине
 * Загружает все корзины с сервера и мигрирует локальную корзину
 */
export function useCartSync() {
  const dispatch = useAppDispatch();
  const { user, isAuthenticated } = useAuth();
  const localCart = useAppSelector(selectLocalCart);
  const lastSyncedUserId = useRef<number | null>(null);
  const syncAttempts = useRef(0);
  const localCartMigrated = useRef(false);
  const MAX_SYNC_ATTEMPTS = 3;

  useEffect(() => {
    // Синхронизируем корзины только для авторизованных пользователей
    if (!isAuthenticated || !user) {
      // Сбрасываем счетчики при выходе
      lastSyncedUserId.current = null;
      syncAttempts.current = 0;
      localCartMigrated.current = false;
      return;
    }

    // Проверяем, не синхронизировали ли мы уже этого пользователя
    if (lastSyncedUserId.current === user.id) {
      return; // Уже синхронизировано для этого пользователя
    }

    // Проверяем лимит попыток
    if (syncAttempts.current >= MAX_SYNC_ATTEMPTS) {
      console.warn('[CartSync] Max sync attempts reached, stopping');
      return;
    }

    const syncCarts = async () => {
      try {
        syncAttempts.current++;

        // ВАЖНО: Сначала загружаем корзины с сервера
        // Это нужно чтобы не потерять существующие товары
        await dispatch(fetchUserCarts(user.id)).unwrap();

        // Успешная синхронизация - запоминаем пользователя
        lastSyncedUserId.current = user.id;

        // Теперь мигрируем локальную корзину на сервер (только один раз)
        // Товары из локальной корзины добавятся к существующим на сервере
        if (
          !localCartMigrated.current &&
          localCart &&
          localCart.items &&
          localCart.items.length > 0
        ) {
          localCartMigrated.current = true; // Помечаем что миграция началась
          // Группируем товары по витринам
          const itemsByStorefront = localCart.items.reduce(
            (acc, item) => {
              const storefrontId = item.storefrontId || 0;
              if (!acc[storefrontId]) {
                acc[storefrontId] = [];
              }
              acc[storefrontId].push(item);
              return acc;
            },
            {} as Record<number, typeof localCart.items>
          );

          // Добавляем товары на сервер для каждой витрины
          let migrationSuccess = true;
          for (const [storefrontId, items] of Object.entries(
            itemsByStorefront
          )) {
            for (const item of items) {
              try {
                await dispatch(
                  addToCart({
                    storefrontId: parseInt(storefrontId),
                    item: {
                      product_id: item.productId,
                      variant_id: item.variantId,
                      quantity: item.quantity,
                    },
                  })
                ).unwrap();
              } catch (error) {
                console.error(
                  `[CartSync] Failed to migrate item ${item.productId}:`,
                  error
                );
                // Продолжаем миграцию других товаров даже если один не удался
                migrationSuccess = false;
              }
            }
          }

          // Очищаем локальную корзину только после миграции
          if (migrationSuccess) {
            dispatch(clearLocalCart());
          } else {
            console.warn(
              `[CartSync] Some items failed to migrate, keeping local cart`
            );
          }
        }

        // Загружаем обновленные корзины с сервера еще раз
        // чтобы синхронизировать состояние после миграции
        await dispatch(fetchUserCarts(user.id)).unwrap();
      } catch (error) {
        console.error('[CartSync] Failed to sync carts:', error);

        // Если это ошибка авторизации, не пытаемся снова
        if (error && typeof error === 'object' && 'message' in error) {
          const message = String(error.message).toLowerCase();
          if (message.includes('401') || message.includes('unauthorized')) {
            console.warn(
              '[CartSync] Authorization error, stopping sync attempts'
            );
            lastSyncedUserId.current = user.id; // Помечаем как синхронизированного чтобы не повторять
          }
        }
      }
    };

    syncCarts();
  }, [user?.id, isAuthenticated, dispatch]); // Убираем только localCart из зависимостей чтобы избежать бесконечных циклов
}
