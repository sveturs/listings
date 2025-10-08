'use client';

import { useEffect, useRef, useState, useMemo } from 'react';
import { useChat } from '@/hooks/useChat';
import { useAuth } from '@/contexts/AuthContext';
import { useTranslations } from 'next-intl';
import { useParams } from 'next/navigation';
import { MarketplaceChat } from '@/types/chat';
import MessageItem from './MessageItem';
import MessageInput from './MessageInput';
import Image from 'next/image';
import Link from 'next/link';
import configManager from '@/config';
import { contactsService } from '@/services/contacts';
import { getLastSeenText } from '@/utils/timeUtils';
import { toast } from '@/utils/toast';
import StorefrontProductQuickView from './StorefrontProductQuickView';
import IncomingContactRequest from './IncomingContactRequest';
import ChatSettings from './ChatSettings';
import { useAppDispatch } from '@/store/hooks';
import { addToCart } from '@/store/slices/cartSlice';
import { addItem } from '@/store/slices/localCartSlice';

interface ChatWindowProps {
  chat?: MarketplaceChat;
  initialListingId?: number;
  initialStorefrontProductId?: number;
  initialSellerId?: number;
  initialContactId?: number;
  onBack?: () => void;
  showBackButton?: boolean;
  onShowChat?: () => void;
}

interface ListingInfo {
  id: number;
  title: string;
  images?: Array<{
    id: number;
    public_url: string;
  }>;
  user_id: number;
}

interface StorefrontProduct {
  id: number;
  name: string;
  price: number;
  storefront_id: number;
  stock_quantity?: number;
  images?: Array<{ url: string }>;
  storefront?: {
    id: number;
    slug: string;
    name: string;
  };
}

export default function ChatWindow({
  chat,
  initialListingId,
  initialStorefrontProductId,
  initialSellerId,
  initialContactId,
  onBack,
  showBackButton,
  onShowChat,
}: ChatWindowProps) {
  const t = useTranslations('chat');
  const params = useParams();
  const locale = params?.locale as string;
  const { user } = useAuth();
  const messagesEndRef = useRef<HTMLDivElement>(null);
  const messagesContainerRef = useRef<HTMLDivElement>(null);
  const [isAtBottom, setIsAtBottom] = useState(true);
  const [isNewChat] = useState(
    !chat &&
      (((initialListingId || initialStorefrontProductId) && initialSellerId) ||
        initialContactId)
  );
  const [isContactChat] = useState(
    (!chat && initialContactId && !initialListingId) ||
      (chat?.id === 0 && !chat?.listing_id)
  );
  const [listingInfo, setListingInfo] = useState<ListingInfo | null>(null);
  const [isLoadingOldMessages, setIsLoadingOldMessages] = useState(false);
  const [isAddingContact, setIsAddingContact] = useState(false);
  const [isInitialLoad, setIsInitialLoad] = useState(true);
  const [showProductQuickView, setShowProductQuickView] = useState(false);
  const [storefrontProduct, setStorefrontProduct] =
    useState<StorefrontProduct | null>(null);
  const [isAddingToCart, setIsAddingToCart] = useState(false);
  const [showChatSettings, setShowChatSettings] = useState(false);
  const dispatch = useAppDispatch();

  const {
    messages,
    isLoading,
    hasMoreMessages,
    messagesLoaded,
    loadMessages,
    markMessagesAsRead,
    typingUsers,
    onlineUsers,
    userLastSeen,
  } = useChat();

  const chatMessages = useMemo(
    () => (chat ? messages[chat.id] || [] : []),
    [chat, messages]
  );
  const hasMore = chat ? hasMoreMessages[chat.id] || false : false;
  const typingInThisChat = useMemo(
    () =>
      chat ? (typingUsers[chat.id] || []).filter((id) => id !== user?.id) : [],
    [chat, typingUsers, user?.id]
  );
  const isOtherUserOnline =
    chat?.other_user && onlineUsers.includes(chat.other_user.id);

  // Загрузка сообщений при смене чата с отменой предыдущих запросов
  useEffect(() => {
    const abortController = new AbortController();

    if (chat) {
      setIsInitialLoad(true);
      // Всегда загружаем сообщения при смене чата
      // Для проверки используем messagesLoaded из Redux store
      const messagesAlreadyLoaded = messagesLoaded[chat.id];

      if (!messagesAlreadyLoaded) {
        loadMessages({
          chat_id: chat.id,
          page: 1,
          limit: 20,
          signal: abortController.signal,
        });
      }
    }

    // Отменяем предыдущий запрос при размонтировании или смене чата
    return () => {
      abortController.abort();
    };
  }, [chat, messagesLoaded, loadMessages]);

  // Прокрутка к последнему сообщению только при открытии чата или смене чата
  useEffect(() => {
    if (chat?.id && isInitialLoad && chatMessages.length > 0 && !isLoading) {
      let observer: MutationObserver | null = null;

      // Функция прокрутки
      const performScroll = () => {
        if (messagesContainerRef.current) {
          messagesContainerRef.current.scrollTop =
            messagesContainerRef.current.scrollHeight;
        }
      };

      // Используем MutationObserver для отслеживания изменений DOM
      if (messagesContainerRef.current) {
        observer = new MutationObserver(() => {
          performScroll();
        });

        observer.observe(messagesContainerRef.current, {
          childList: true,
          subtree: true,
        });

        // Делаем несколько попыток прокрутки
        const timeouts = [50, 150, 300, 500].map((delay) =>
          setTimeout(() => {
            performScroll();
          }, delay)
        );

        // Останавливаем наблюдение через 600мс
        setTimeout(() => {
          observer?.disconnect();
          setIsInitialLoad(false);
          timeouts.forEach(clearTimeout);
        }, 600);
      }

      return () => {
        observer?.disconnect();
      };
    }
  }, [chat, chatMessages.length, isInitialLoad, isLoading]);

  // Загрузка информации о товаре витрины для существующих чатов
  useEffect(() => {
    if (chat?.storefront_product_id && !storefrontProduct) {
      const apiUrl = configManager.getApiUrl();
      fetch(
        `${apiUrl}/api/v1/storefronts/products/${chat.storefront_product_id}`,
        {
          credentials: 'include',
          headers: {
            'Content-Type': 'application/json',
          },
        }
      )
        .then((res) => res.json())
        .then((result) => {
          const data = result.data || result;
          if (data && data.id) {
            setStorefrontProduct(data);
          }
        })
        .catch((err) => {
          console.error('Error loading storefront product:', err);
          // Fallback data for testing
          setStorefrontProduct({
            id: chat.storefront_product_id || 0,
            name: chat.listing?.title || 'Product',
            price: 1000,
            storefront_id: chat.seller_id || 1,
            stock_quantity: 10,
            storefront: {
              id: chat.seller_id || 1,
              slug: 'store', // Default slug, should be loaded from API
              name: 'Store',
            },
          });
        });
    }
  }, [
    chat?.storefront_product_id,
    storefrontProduct,
    chat?.listing?.title,
    chat?.seller_id,
  ]);

  // Загрузка информации об объявлении для новых чатов
  useEffect(() => {
    // Для обычных объявлений
    if (isNewChat && initialListingId && !listingInfo && !isContactChat) {
      const apiUrl = configManager.getApiUrl();
      fetch(`${apiUrl}/api/v1/marketplace/listings/${initialListingId}`, {
        credentials: 'include',
        headers: {
          'Content-Type': 'application/json',
        },
      })
        .then((res) => res.json())
        .then((result) => {
          // Проверяем обертку ответа
          const data = result.data || result;
          if (data && data.id) {
            setListingInfo(data);
          }
        })
        .catch((err) => console.error('Error loading listing info:', err));
    }

    // Для товаров витрин - используем заглушку пока endpoint не реализован
    if (
      isNewChat &&
      initialStorefrontProductId &&
      !listingInfo &&
      !isContactChat
    ) {
      console.log(
        'Storefront product loading - using placeholder data (endpoint not implemented):',
        initialStorefrontProductId
      );
      // TODO: Implement proper storefront product loading when endpoint is available
      // For now, set placeholder data with product ID
      setListingInfo({
        id: initialStorefrontProductId,
        title: t('storefrontProduct', { id: initialStorefrontProductId }),
        images: [],
        user_id: initialSellerId || 0,
      });
    }
  }, [
    isNewChat,
    initialListingId,
    initialStorefrontProductId,
    isContactChat,
    listingInfo,
    initialSellerId,
    t,
  ]);

  // Прокрутка к низу только при новых сообщениях (не при загрузке старых)
  useEffect(() => {
    if (isAtBottom && !isLoadingOldMessages && !isInitialLoad) {
      scrollToBottom();
    }
  }, [chatMessages.length, isAtBottom, isLoadingOldMessages, isInitialLoad]);

  // Пометка сообщений как прочитанных
  useEffect(() => {
    if (!chat) return;

    const unreadMessages = chatMessages
      .filter((msg) => !msg.is_read && msg.receiver_id === user?.id)
      .map((msg) => msg.id);

    if (unreadMessages.length > 0) {
      markMessagesAsRead(chat.id, unreadMessages);
    }
  }, [chatMessages, chat, user?.id, markMessagesAsRead]);

  // Отслеживание позиции скролла
  const handleScroll = async () => {
    if (!messagesContainerRef.current) return;

    const { scrollTop, scrollHeight, clientHeight } =
      messagesContainerRef.current;
    const isNowAtBottom = scrollHeight - scrollTop - clientHeight < 100;
    setIsAtBottom(isNowAtBottom);

    // Загрузка предыдущих сообщений при скролле вверх
    const currentPage = Math.ceil(chatMessages.length / 20);
    if (
      scrollTop < 100 &&
      hasMore &&
      !isLoading &&
      !isLoadingOldMessages &&
      chat
    ) {
      setIsLoadingOldMessages(true);

      // Сохраняем текущую высоту скролла
      const oldScrollHeight = scrollHeight;

      try {
        await loadMessages({
          chat_id: chat.id,
          page: currentPage + 1,
          limit: 20,
        });

        // После загрузки восстанавливаем позицию скролла
        setTimeout(() => {
          if (messagesContainerRef.current) {
            const newScrollHeight = messagesContainerRef.current.scrollHeight;
            const scrollDiff = newScrollHeight - oldScrollHeight;
            messagesContainerRef.current.scrollTop = scrollDiff;
          }
          setIsLoadingOldMessages(false);
        }, 100);
      } catch (error) {
        console.error('Error loading old messages:', error);
        setIsLoadingOldMessages(false);
      }
    }
  };

  const scrollToBottom = () => {
    messagesEndRef.current?.scrollIntoView({ behavior: 'smooth' });
  };

  const getChatTitle = () => {
    if (chat?.listing) {
      // Проверяем специальные маркеры от бэкенда
      if (chat.listing.title === '__DIRECT_MESSAGE__') {
        return t('directMessage');
      }
      if (chat.listing.title === '__DELETED_LISTING__') {
        return t('deletedListing');
      }
      // Если это чат товара витрины, показываем название товара
      if (chat.storefront_product_id && chat.storefront_product_id > 0) {
        return chat.listing.title; // Backend уже отправляет название товара витрины
      }
      return chat.listing.title;
    }
    if (listingInfo) {
      return listingInfo.title;
    }
    // Если это новый чат с товаром витрины, но информация еще загружается
    if (isNewChat && initialStorefrontProductId && !isContactChat) {
      return t('loadingProduct');
    }
    // Если это новый чат с обычным объявлением, но информация еще загружается
    if (isNewChat && initialListingId && !isContactChat) {
      return t('loadingListing');
    }
    if (isContactChat) {
      return t('directMessage');
    }
    return t('newChat');
  };

  const handleTitleClick = () => {
    // Открываем быстрый просмотр только для товаров витрин
    if (chat?.storefront_product_id && chat.storefront_product_id > 0) {
      setShowProductQuickView(true);
    }
  };

  const handleAddToContacts = async () => {
    if (!user || !chat?.other_user?.id || isAddingContact) return;

    setIsAddingContact(true);
    try {
      await contactsService.addContact({
        contact_user_id: chat.other_user.id,
        notes: `Added from chat about ${getChatTitle()}`,
        added_from_chat_id: chat.id,
      });

      toast.success(t('contactAdded'));
    } catch (error) {
      console.error('Error adding contact:', error);
      // Проверяем тип ошибки и показываем соответствующее сообщение
      if (error instanceof Error) {
        if (error.message.includes('already exists')) {
          toast.warning(t('contactAlreadyExists'));
        } else if (error.message.includes('cannot add yourself')) {
          toast.error(t('cannotAddYourself'));
        } else if (error.message.includes('does not allow contact requests')) {
          toast.info(t('userDoesNotAllowContacts'));
        } else if (error.message.includes('Unauthorized')) {
          toast.warning(t('pleaseLoginToAddContacts'));
        } else {
          toast.error(t('failedToAddContact'));
        }
      } else {
        toast.error(t('failedToAddContact'));
      }
    } finally {
      setIsAddingContact(false);
    }
  };

  const getContactUserId = () => {
    if (chat?.other_user?.id) return chat.other_user.id;
    if (initialContactId) return initialContactId;
    if (initialSellerId && user?.id !== initialSellerId) return initialSellerId;
    return null;
  };

  return (
    <div className="flex flex-col h-full w-full max-w-full overflow-hidden relative">
      {/* Заголовок чата */}
      <div className="navbar bg-base-100 border-b border-base-300 min-h-0 p-2 sm:p-3 relative z-30">
        <div className="navbar-start flex-1 gap-2">
          {showBackButton && (
            <button
              onClick={onBack}
              className="btn btn-ghost btn-circle btn-sm"
            >
              <svg
                className="w-4 h-4 sm:w-5 sm:h-5"
                fill="none"
                viewBox="0 0 24 24"
                stroke="currentColor"
              >
                <path
                  strokeLinecap="round"
                  strokeLinejoin="round"
                  strokeWidth={2}
                  d="M15 19l-7-7 7-7"
                />
              </svg>
            </button>
          )}

          {/* Информация о товаре или контакте */}
          {(chat?.listing || chat?.listing_id || listingInfo) &&
            !isContactChat && (
              <div
                className="avatar flex-shrink-0 cursor-pointer"
                onClick={
                  chat?.storefront_product_id ? handleTitleClick : undefined
                }
              >
                <div className="w-8 h-8 sm:w-10 sm:h-10 rounded">
                  <Image
                    src={
                      (chat?.listing?.images?.[0]?.public_url &&
                        configManager.buildImageUrl(
                          chat.listing.images[0].public_url
                        )) ||
                      (listingInfo?.images?.[0]?.public_url &&
                        configManager.buildImageUrl(
                          listingInfo.images[0].public_url
                        )) ||
                      '/placeholder-listing.jpg'
                    }
                    alt={chat?.listing?.title || listingInfo?.title || ''}
                    width={40}
                    height={40}
                    className="object-cover"
                  />
                </div>
              </div>
            )}

          {/* Аватар для прямого сообщения */}
          {isContactChat && (
            <div className="avatar flex-shrink-0">
              <div className="w-8 h-8 sm:w-10 sm:h-10 rounded-full bg-base-300 flex items-center justify-center">
                <span className="text-sm sm:text-lg font-semibold">
                  <svg
                    className="w-4 h-4 sm:w-6 sm:h-6"
                    fill="none"
                    viewBox="0 0 24 24"
                    stroke="currentColor"
                  >
                    <path
                      strokeLinecap="round"
                      strokeLinejoin="round"
                      strokeWidth={2}
                      d="M16 7a4 4 0 11-8 0 4 4 0 018 0zM12 14a7 7 0 00-7 7h14a7 7 0 00-7-7z"
                    />
                  </svg>
                </span>
              </div>
            </div>
          )}

          <div className="flex-1 min-w-0">
            <div className="flex items-center gap-2">
              <h2
                className={`text-sm sm:text-base font-semibold truncate ${
                  chat?.storefront_product_id
                    ? 'cursor-pointer hover:text-primary transition-colors'
                    : ''
                }`}
                onClick={handleTitleClick}
              >
                {getChatTitle()}
              </h2>
              {/* Цена для товара витрины */}
              {!!chat?.storefront_product_id && chat?.listing?.price ? (
                <span className="text-sm font-bold text-primary">
                  {chat.listing.price} RSD
                </span>
              ) : null}
            </div>
            <div className="flex items-center gap-2 text-xs sm:text-sm text-base-content/70">
              {isNewChat ? (
                <span>{t('startNewConversation')}</span>
              ) : (
                <>
                  <span className="truncate">
                    {chat?.other_user?.name || t('unknownUser')}
                  </span>
                  {chat?.other_user && (
                    <span className="text-xs">
                      {isOtherUserOnline ? (
                        <span className="text-gray-500">{t('online')}</span>
                      ) : userLastSeen[chat.other_user.id] ? (
                        <span className="text-base-content/50">
                          {(() => {
                            const time = getLastSeenText(
                              userLastSeen[chat.other_user.id],
                              (
                                key: string,
                                values?: Record<string, string | number | Date>
                              ) => {
                                // Безопасное преобразование для next-intl
                                const safeValues = values as Parameters<
                                  typeof t
                                >[1];
                                return t(key, safeValues);
                              }
                            );
                            return t('lastSeen', { time });
                          })()}
                        </span>
                      ) : null}
                    </span>
                  )}
                </>
              )}
            </div>
          </div>
        </div>
        {/* Кнопки действий */}
        <div className="navbar-end gap-1">
          {/* Кнопка добавить в корзину для товаров витрин */}
          {!!chat?.storefront_product_id && (
            <button
              className="btn btn-primary btn-xs sm:btn-sm"
              title={t('addToCart')}
              disabled={isAddingToCart || !storefrontProduct}
              onClick={async () => {
                if (!storefrontProduct || isAddingToCart) return;

                setIsAddingToCart(true);
                try {
                  if (user) {
                    // For authenticated users, save to backend
                    await dispatch(
                      addToCart({
                        storefrontId: storefrontProduct.storefront_id,
                        item: {
                          product_id: storefrontProduct.id,
                          quantity: 1,
                        },
                      })
                    ).unwrap();
                    toast.success(t('productAddedToCart'));
                  } else {
                    // For non-authenticated users, save to local storage
                    dispatch(
                      addItem({
                        productId: storefrontProduct.id,
                        name: storefrontProduct.name,
                        price: storefrontProduct.price,
                        currency: 'RSD',
                        quantity: 1,
                        stockQuantity: storefrontProduct.stock_quantity,
                        image: storefrontProduct.images?.[0]?.url,
                        storefrontId: storefrontProduct.storefront_id,
                      })
                    );
                    toast.success(t('productAddedToCart'));
                  }
                } catch (error) {
                  console.error('Failed to add to cart:', error);
                  toast.error(t('failedToAddToCart'));
                } finally {
                  setIsAddingToCart(false);
                }
              }}
            >
              <svg
                className="w-4 h-4 sm:w-5 sm:h-5"
                fill="none"
                viewBox="0 0 24 24"
                stroke="currentColor"
              >
                <path
                  strokeLinecap="round"
                  strokeLinejoin="round"
                  strokeWidth={2}
                  d="M3 3h2l.4 2M7 13h10l4-8H5.4M7 13L5.4 5M7 13l-2.293 2.293c-.63.63-.184 1.707.707 1.707H17m0 0a2 2 0 100 4 2 2 0 000-4zm-8 2a2 2 0 11-4 0 2 2 0 014 0z"
                />
              </svg>
            </button>
          )}

          {/* Кнопка просмотра детальной страницы */}
          {!!chat?.storefront_product_id ? (
            <Link
              href={`/${locale}/storefronts/product/${chat.storefront_product_id}`}
              className="btn btn-ghost btn-xs sm:btn-sm"
              title={t('viewDetails')}
            >
              <svg
                className="w-4 h-4 sm:w-5 sm:h-5"
                fill="none"
                viewBox="0 0 24 24"
                stroke="currentColor"
              >
                <path
                  strokeLinecap="round"
                  strokeLinejoin="round"
                  strokeWidth={2}
                  d="M10 6H6a2 2 0 00-2 2v10a2 2 0 002 2h10a2 2 0 002-2v-4M14 4h6m0 0v6m0-6L10 14"
                />
              </svg>
            </Link>
          ) : (
            !!(chat?.listing || chat?.listing_id || listingInfo) &&
            !isContactChat && (
              <Link
                href={`/${locale}/marketplace/${chat?.listing_id || chat?.listing?.id || listingInfo?.id}`}
                className="btn btn-ghost btn-xs sm:btn-sm"
                title={t('viewListing')}
              >
                <svg
                  className="w-4 h-4 sm:w-5 sm:h-5"
                  fill="none"
                  viewBox="0 0 24 24"
                  stroke="currentColor"
                >
                  <path
                    strokeLinecap="round"
                    strokeLinejoin="round"
                    strokeWidth={2}
                    d="M13.828 10.172a4 4 0 00-5.656 0l-4 4a4 4 0 105.656 5.656l1.102-1.101m-.758-4.899a4 4 0 005.656 0l4-4a4 4 0 00-5.656-5.656l-1.1 1.1"
                  />
                </svg>
              </Link>
            )
          )}

          {getContactUserId() && (
            <button
              className="btn btn-ghost btn-xs sm:btn-sm"
              title={t('addToContacts')}
              onClick={handleAddToContacts}
              disabled={isAddingContact}
            >
              {isAddingContact ? (
                <span className="loading loading-spinner loading-xs sm:loading-sm"></span>
              ) : (
                <svg
                  className="w-4 h-4 sm:w-5 sm:h-5"
                  fill="none"
                  viewBox="0 0 24 24"
                  stroke="currentColor"
                >
                  <path
                    strokeLinecap="round"
                    strokeLinejoin="round"
                    strokeWidth={2}
                    d="M18 9v3m0 0v3m0-3h3m-3 0h-3m-2-5a4 4 0 11-8 0 4 4 0 018 0zM3 20a6 6 0 0112 0v1H3v-1z"
                  />
                </svg>
              )}
            </button>
          )}

          {/* Кнопка настроек переводов */}
          <button
            className="btn btn-ghost btn-xs sm:btn-sm gap-0.5"
            title={t('translation.translationSettings')}
            onClick={() => setShowChatSettings(true)}
          >
            <div className="flex items-center gap-0.5">
              {/* Стильная иконка перевода АБВ↔ABC */}
              <span className="text-[10px] sm:text-xs font-semibold opacity-90">
                АБВ
              </span>
              <svg
                className="w-2.5 h-2.5 sm:w-3 sm:h-3"
                fill="none"
                viewBox="0 0 24 24"
                stroke="currentColor"
              >
                <path
                  strokeLinecap="round"
                  strokeLinejoin="round"
                  strokeWidth={2.5}
                  d="M8 7h12m0 0l-4-4m4 4l-4 4m0 6H4m0 0l4 4m-4-4l4-4"
                />
              </svg>
              <span className="text-[10px] sm:text-xs font-semibold opacity-90">
                ABC
              </span>
              {/* Маленькая шестеренка */}
              <svg
                className="w-2 h-2 sm:w-2.5 sm:h-2.5 ml-0.5 opacity-60"
                fill="none"
                viewBox="0 0 24 24"
                stroke="currentColor"
              >
                <path
                  strokeLinecap="round"
                  strokeLinejoin="round"
                  strokeWidth={2}
                  d="M10.325 4.317c.426-1.756 2.924-1.756 3.35 0a1.724 1.724 0 002.573 1.066c1.543-.94 3.31.826 2.37 2.37a1.724 1.724 0 001.065 2.572c1.756.426 1.756 2.924 0 3.35a1.724 1.724 0 00-1.066 2.573c.94 1.543-.826 3.31-2.37 2.37a1.724 1.724 0 00-2.572 1.065c-.426 1.756-2.924 1.756-3.35 0a1.724 1.724 0 00-2.573-1.066c-1.543.94-3.31-.826-2.37-2.37a1.724 1.724 0 00-1.065-2.572c-1.756-.426-1.756-2.924 0-3.35a1.724 1.724 0 001.066-2.573c-.94-1.543.826-3.31 2.37-2.37.996.608 2.296.07 2.572-1.065z"
                />
                <path
                  strokeLinecap="round"
                  strokeLinejoin="round"
                  strokeWidth={2}
                  d="M15 12a3 3 0 11-6 0 3 3 0 016 0z"
                />
              </svg>
            </div>
          </button>
        </div>
      </div>

      {/* Уведомление о входящем запросе в контакты */}
      {chat?.other_user?.id && (
        <IncomingContactRequest
          otherUserId={chat.other_user.id}
          chatId={chat?.id}
          onRequestHandled={() => {
            // Можно обновить UI или показать уведомление после обработки запроса
          }}
        />
      )}

      {/* Контейнер с фоном на всю высоту */}
      <div
        className="flex-1 relative"
        style={{
          backgroundColor: 'oklch(var(--b2))',
          backgroundImage: `url("data:image/svg+xml,%3Csvg xmlns='http://www.w3.org/2000/svg' width='28' height='49' viewBox='0 0 28 49'%3E%3Cg fill-rule='evenodd'%3E%3Cg id='hexagons' fill='%239C92AC' fill-opacity='0.1' fill-rule='nonzero'%3E%3Cpath d='M13.99 9.25l13 7.5v15l-13 7.5L1 31.75v-15l12.99-7.5zM3 17.9v12.7l10.99 6.34 11-6.35V17.9l-11-6.34L3 17.9zM0 15l12.98-7.5V0h-2v6.35L0 12.69v2.3zm0 18.5L12.98 41v8h-2v-6.85L0 35.81v-2.3zM15 0v7.5L27.99 15H28v-2.31h-.01L17 6.35V0h-2zm0 49v-8l12.99-7.5H28v2.31h-.01L17 42.15V49h-2z'/%3E%3C/g%3E%3C/g%3E%3C/svg%3E")`,
        }}
      >
        {/* Область сообщений */}
        <div
          ref={messagesContainerRef}
          onScroll={handleScroll}
          className="absolute inset-x-0 top-0 overflow-y-auto p-3 sm:p-4 lg:px-8"
          style={{
            bottom:
              'calc(4rem + 3.06rem + max(0px, env(safe-area-inset-bottom, 0px)))',
          }}
        >
          <style jsx>{`
            @media (min-width: 768px) {
              div {
                bottom: calc(
                  4rem + max(0px, env(safe-area-inset-bottom, 0px))
                ) !important;
              }
            }
          `}</style>
          {/* Индикатор загрузки старых сообщений */}
          {(hasMore || isLoadingOldMessages) && (
            <div className="text-center py-2">
              {isLoadingOldMessages ? (
                <div className="flex items-center justify-center gap-2">
                  <span className="loading loading-spinner loading-sm"></span>
                  <span className="text-sm opacity-50">
                    {t('loadingOldMessages')}
                  </span>
                </div>
              ) : (
                <button
                  className="btn btn-ghost btn-sm text-base-content/50"
                  onClick={() => handleScroll()}
                >
                  {t('scrollUpToLoadMore')}
                </button>
              )}
            </div>
          )}

          {/* Сообщения */}
          {chatMessages.length > 0 ? (
            chatMessages.map((message, index) => (
              <MessageItem
                key={`${message.id}-${index}`}
                message={message}
                isOwn={message.sender_id === user?.id}
              />
            ))
          ) : isNewChat ? (
            <div className="flex items-center justify-center h-full text-base-content/50">
              <div className="text-center">
                <p className="text-lg">{t('sendFirstMessage')}</p>
              </div>
            </div>
          ) : null}

          {/* Индикатор печатания */}
          {typingInThisChat.length > 0 && (
            <div className="chat chat-start">
              <div className="chat-bubble chat-bubble-secondary">
                <span className="loading loading-dots loading-xs"></span>
              </div>
              <div className="chat-footer opacity-50 text-xs">
                {typingInThisChat.length === 1
                  ? chat?.other_user?.name || t('userTyping')
                  : t('usersTyping', { count: typingInThisChat.length })}
              </div>
            </div>
          )}

          <div ref={messagesEndRef} />
        </div>

        {/* Кнопка прокрутки вниз */}
        {!isAtBottom && (
          <button
            onClick={scrollToBottom}
            className="absolute right-4 btn btn-circle btn-sm btn-primary shadow-lg z-10"
            style={{
              bottom:
                'calc(4.06rem + 3.06rem + max(0px, env(safe-area-inset-bottom, 0px)))',
            }}
          >
            <style jsx>{`
              @media (min-width: 768px) {
                button {
                  bottom: calc(
                    5rem + max(0px, env(safe-area-inset-bottom, 0px))
                  ) !important;
                }
              }
            `}</style>
            <svg
              className="w-3 h-3 sm:w-4 sm:h-4"
              fill="none"
              viewBox="0 0 24 24"
              stroke="currentColor"
            >
              <path
                strokeLinecap="round"
                strokeLinejoin="round"
                strokeWidth={2}
                d="M19 14l-7 7m0 0l-7-7m7 7V3"
              />
            </svg>
          </button>
        )}

        {/* Поле ввода - позиционировано внизу контейнера с фоном */}
        <div
          className="absolute left-0 right-0 z-20"
          style={{
            bottom:
              'calc(3.06rem + max(0px, env(safe-area-inset-bottom, 0px)))',
          }}
        >
          <style jsx>{`
            @media (min-width: 768px) {
              div {
                bottom: max(0px, env(safe-area-inset-bottom, 0px)) !important;
              }
            }
          `}</style>
          <MessageInput
            chat={chat}
            initialListingId={initialListingId}
            initialStorefrontProductId={initialStorefrontProductId}
            initialSellerId={initialSellerId || initialContactId}
            onShowChat={onShowChat}
          />
        </div>
      </div>

      {/* Quick View Modal for Storefront Product */}
      {!!chat?.storefront_product_id && (
        <StorefrontProductQuickView
          productId={chat.storefront_product_id}
          isOpen={showProductQuickView}
          onClose={() => setShowProductQuickView(false)}
        />
      )}

      {/* Chat Settings Modal */}
      <ChatSettings
        isOpen={showChatSettings}
        onClose={() => setShowChatSettings(false)}
      />
    </div>
  );
}
