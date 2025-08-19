import { useEffect } from 'react';
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

  useEffect(() => {
    // Синхронизируем корзины только для авторизованных пользователей
    if (!isAuthenticated || !user) return;

    const syncCarts = async () => {
      try {
        // ВАЖНО: Сначала загружаем корзины с сервера
        // Это нужно чтобы не потерять существующие товары
        await dispatch(fetchUserCarts(user.id)).unwrap();

        // Теперь мигрируем локальную корзину на сервер
        // Товары из локальной корзины добавятся к существующим на сервере
        if (localCart && localCart.items && localCart.items.length > 0) {
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
      }
    };

    syncCarts();
  }, [user, isAuthenticated, dispatch, localCart]); // Добавляем зависимости обратно
}
