'use client';

import { useEffect, useState, useRef, useCallback } from 'react';
import { useChat } from '@/hooks/useChat';
import { useAuth } from '@/contexts/AuthContext';
import { useTranslations } from 'next-intl';
import { MarketplaceChat } from '@/types/chat';
import { UserContact } from '@/types/contacts';
import { contactsService } from '@/services/contacts';
import Image from 'next/image';
import { useLocale } from 'next-intl';
import configManager from '@/config';

interface ChatListProps {
  onChatSelect: (chat: MarketplaceChat) => void;
}

export default function ChatList({ onChatSelect }: ChatListProps) {
  const t = useTranslations('chat');
  const locale = useLocale();
  const { user, isAuthenticated } = useAuth();
  const [searchQuery, setSearchQuery] = useState('');
  const [activeTab, setActiveTab] = useState<'chats' | 'contacts'>('chats');
  const [contacts, setContacts] = useState<UserContact[]>([]);
  const [contactsLoading, setContactsLoading] = useState(false);
  const loadMoreRef = useRef<HTMLDivElement>(null);

  const {
    chats,
    currentChat,
    isLoading,
    hasMoreChats,
    loadChats,
    onlineUsers,
  } = useChat();

  // Функция для загрузки контактов
  const loadContacts = useCallback(async () => {
    if (!user) return;

    setContactsLoading(true);
    try {
      const response = await contactsService.getContacts('accepted');
      setContacts(response.contacts);
    } catch (error) {
      // Если ошибка авторизации, тихо устанавливаем пустой массив
      if (error instanceof Error && error.message === 'Unauthorized') {
        setContacts([]);
      } else {
        // Логируем только другие ошибки
        console.error('Error loading contacts:', error);
      }
    } finally {
      setContactsLoading(false);
    }
  }, [user]);

  // Загрузка чатов при монтировании
  useEffect(() => {
    // Загружаем чаты только если пользователь авторизован
    if (isAuthenticated && user && chats.length === 0) {
      loadChats(1);
    }
  }, [isAuthenticated, user, chats.length, loadChats]);

  // Загрузка контактов при переключении на вкладку контактов
  useEffect(() => {
    if (
      activeTab === 'contacts' &&
      isAuthenticated &&
      user &&
      (!contacts || contacts.length === 0)
    ) {
      loadContacts();
    }
  }, [activeTab, isAuthenticated, user, contacts, loadContacts]);

  // Бесконечная прокрутка
  useEffect(() => {
    // Не настраиваем observer если пользователь не авторизован
    if (!user) return;

    const observer = new IntersectionObserver(
      (entries) => {
        if (entries[0].isIntersecting && hasMoreChats && !isLoading && user) {
          loadChats(chats.length / 20 + 1);
        }
      },
      { threshold: 0.1 }
    );

    if (loadMoreRef.current) {
      observer.observe(loadMoreRef.current);
    }

    return () => observer.disconnect();
  }, [hasMoreChats, isLoading, chats.length, loadChats, user]);

  // Фильтрация чатов по поиску
  const filteredChats = chats.filter((chat) => {
    if (!searchQuery) return true;

    const query = searchQuery.toLowerCase();
    const otherUserName = chat.other_user?.name?.toLowerCase() || '';
    const listingTitle = chat.listing?.title?.toLowerCase() || '';
    const lastMessage = chat.last_message?.content?.toLowerCase() || '';

    return (
      otherUserName.includes(query) ||
      listingTitle.includes(query) ||
      lastMessage.includes(query)
    );
  });

  // Сортировка чатов по дате последнего сообщения (новые сверху)
  const sortedChats = [...filteredChats].sort((a, b) => {
    const aTime = new Date(a.last_message_at || a.created_at).getTime();
    const bTime = new Date(b.last_message_at || b.created_at).getTime();
    return bTime - aTime; // Новые сообщения сверху
  });

  // Фильтрация контактов по поиску
  const filteredContacts = (contacts || []).filter((contact) => {
    if (!searchQuery) return true;

    const query = searchQuery.toLowerCase();
    const contactName = contact.contact_user?.name?.toLowerCase() || '';
    const contactEmail = contact.contact_user?.email?.toLowerCase() || '';

    return contactName.includes(query) || contactEmail.includes(query);
  });

  const getChatTitle = (chat: MarketplaceChat) => {
    if (chat.listing) {
      // Проверяем специальные маркеры от бэкенда
      if (chat.listing.title === '__DIRECT_MESSAGE__') {
        return t('directMessage');
      }
      if (chat.listing.title === '__DELETED_LISTING__') {
        return t('deletedListing');
      }
      // Если это чат товара витрины, показываем название товара
      if (chat.storefront_product_id && chat.storefront_product_id > 0) {
        return chat.listing.title; // Backend уже отправляет название товара витрины в listing.title
      }
      return chat.listing.title;
    }
    return t('deletedListing');
  };

  const getChatSubtitle = (chat: MarketplaceChat) => {
    return chat.other_user?.name || t('unknownUser');
  };

  const getChatAvatar = (chat: MarketplaceChat) => {
    if (chat.listing?.images?.[0]?.public_url) {
      const imageUrl = configManager.buildImageUrl(
        chat.listing.images[0].public_url
      );
      return imageUrl;
    }
    return '/placeholder-listing.jpg';
  };

  // Функция для создания/открытия чата с контактом
  const handleContactSelect = async (contact: UserContact) => {
    if (!contact.contact_user?.id || !user?.id) return;

    // Сначала ищем существующий прямой чат с этим контактом
    const existingChat = chats.find(
      (chat) =>
        !chat.listing_id && // Прямой чат (без привязки к объявлению)
        contact.contact_user &&
        ((chat.buyer_id === user.id &&
          chat.seller_id === contact.contact_user.id) ||
          (chat.seller_id === user.id &&
            chat.buyer_id === contact.contact_user.id))
    );

    if (existingChat) {
      // Если чат уже существует, открываем его
      onChatSelect(existingChat);
    } else {
      // Создаем виртуальный чат для прямого общения с контактом
      const directChat: MarketplaceChat = {
        id: 0, // Временный ID для нового чата
        listing_id: 0,
        buyer_id: user.id,
        seller_id: contact.contact_user.id,
        last_message: undefined,
        last_message_at: new Date().toISOString(),
        unread_count: 0,
        is_archived: false,
        created_at: new Date().toISOString(),
        updated_at: new Date().toISOString(),
        other_user: {
          id: contact.contact_user.id,
          name: contact.contact_user.name || '',
          email: contact.contact_user.email || '',
          picture_url: '',
          provider: '',
        },
      };

      onChatSelect(directChat);
    }
  };

  // Если пользователь не авторизован, показываем сообщение
  if (!isAuthenticated || !user) {
    return (
      <div className="flex flex-col h-full bg-base-100 items-center justify-center p-8">
        <div className="text-center">
          <svg
            className="w-24 h-24 mx-auto mb-4 text-base-content/20"
            fill="none"
            viewBox="0 0 24 24"
            stroke="currentColor"
          >
            <path
              strokeLinecap="round"
              strokeLinejoin="round"
              strokeWidth={1}
              d="M8 7h12m0 0l-4-4m4 4l-4 4m0 6H4m0 0l4 4m-4-4l4-4"
            />
          </svg>
          <h3 className="text-lg font-semibold text-base-content/70 mb-2">
            {t('loginRequired')}
          </h3>
          <p className="text-base-content/50">{t('loginToViewChats')}</p>
        </div>
      </div>
    );
  }

  return (
    <div className="flex flex-col h-full bg-base-200">
      {/* Заголовок и поиск */}
      <div className="p-4 bg-base-100 border-b border-base-300">
        {/* Поиск с иконкой */}
        <div className="form-control mb-4">
          <div className="relative">
            <input
              type="text"
              placeholder={t('searchPlaceholder')}
              className="input input-bordered w-full pl-10"
              value={searchQuery}
              onChange={(e) => setSearchQuery(e.target.value)}
            />
            <svg
              className="absolute left-3 top-1/2 transform -translate-y-1/2 w-4 h-4 text-base-content/50"
              fill="none"
              stroke="currentColor"
              viewBox="0 0 24 24"
            >
              <path
                strokeLinecap="round"
                strokeLinejoin="round"
                strokeWidth={2}
                d="M21 21l-6-6m2-5a7 7 0 11-14 0 7 7 0 0114 0z"
              />
            </svg>
          </div>
        </div>

        {/* Табы DaisyUI */}
        <div className="tabs tabs-boxed">
          <button
            className={`tab tab-sm flex-1 gap-1 ${
              activeTab === 'chats' ? 'tab-active' : ''
            }`}
            onClick={() => setActiveTab('chats')}
          >
            <svg
              className="w-4 h-4 mr-1"
              fill="none"
              stroke="currentColor"
              viewBox="0 0 24 24"
              strokeWidth={2}
            >
              <path
                strokeLinecap="round"
                strokeLinejoin="round"
                d="M8.625 12a.375.375 0 11-.75 0 .375.375 0 01.75 0zm0 0H8.25m4.125 0a.375.375 0 11-.75 0 .375.375 0 01.75 0zm0 0H12m4.125 0a.375.375 0 11-.75 0 .375.375 0 01.75 0zm0 0h-.375M21 12c0 4.556-4.03 8.25-9 8.25a9.764 9.764 0 01-2.555-.337A5.972 5.972 0 015.41 20.97a5.969 5.969 0 01-.474-.065 4.48 4.48 0 00.978-2.025c.09-.457-.133-.901-.467-1.226C3.93 16.178 3 14.189 3 12c0-4.556 4.03-8.25 9-8.25s9 3.694 9 8.25z"
              />
            </svg>
            {t('chatsTab')}
          </button>
          <button
            className={`tab tab-sm flex-1 gap-1 ${
              activeTab === 'contacts' ? 'tab-active' : ''
            }`}
            onClick={() => setActiveTab('contacts')}
          >
            <svg
              className="w-4 h-4 mr-1"
              fill="none"
              stroke="currentColor"
              viewBox="0 0 24 24"
              strokeWidth={2}
            >
              <path
                strokeLinecap="round"
                strokeLinejoin="round"
                d="M15 19.128a9.38 9.38 0 002.625.372 9.337 9.337 0 004.121-.952 4.125 4.125 0 00-7.533-2.493M15 19.128v-.003c0-1.113-.285-2.16-.786-3.07M15 19.128v.106A12.318 12.318 0 018.624 21c-2.331 0-4.512-.645-6.374-1.766l-.001-.109a6.375 6.375 0 0111.964-3.07M12 6.375a3.375 3.375 0 11-6.75 0 3.375 3.375 0 016.75 0zm8.25 2.25a2.625 2.625 0 11-5.25 0 2.625 2.625 0 015.25 0z"
              />
            </svg>
            {t('contactsTab')}
          </button>
        </div>
      </div>

      {/* Список чатов или контактов */}
      <div className="flex-1 overflow-y-auto p-3 space-y-2">
        {activeTab === 'chats' ? (
          // Список чатов с card компонентами
          <>
            {isLoading && chats.length === 0 ? (
              <div className="flex justify-center p-8">
                <span className="loading loading-spinner loading-lg text-primary"></span>
              </div>
            ) : sortedChats.length === 0 ? (
              <div className="card bg-base-100 shadow-sm">
                <div className="card-body text-center py-12">
                  <div className="text-4xl mb-4">💬</div>
                  <h3 className="card-title justify-center text-base-content/70">
                    {searchQuery ? t('noSearchResults') : t('noChats')}
                  </h3>
                  <p className="text-base-content/50">
                    {searchQuery
                      ? t('tryDifferentSearch')
                      : t('startNewConversation')}
                  </p>
                </div>
              </div>
            ) : (
              <>
                {sortedChats.map((chat) => (
                  <div
                    key={chat.id}
                    onClick={() => onChatSelect(chat)}
                    className={`card card-compact cursor-pointer transition-all ${
                      currentChat?.id === chat.id
                        ? 'bg-primary/10 border-primary shadow-md'
                        : 'bg-base-100 hover:bg-base-200 shadow-sm'
                    }`}
                  >
                    <div className="card-body">
                      <div className="flex items-center gap-4">
                        {/* Аватар с индикатором */}
                        <div className="avatar indicator">
                          {chat.unread_count > 0 && (
                            <span className="indicator-item badge badge-primary badge-sm">
                              {chat.unread_count}
                            </span>
                          )}
                          <div className="w-12 rounded-lg">
                            <Image
                              src={getChatAvatar(chat)}
                              alt={getChatTitle(chat)}
                              width={48}
                              height={48}
                              className="object-cover"
                            />
                          </div>
                        </div>

                        {/* Информация о чате */}
                        <div className="flex-1 min-w-0">
                          <div className="flex justify-between items-start">
                            <h3 className="font-semibold text-sm truncate">
                              {getChatTitle(chat)}
                            </h3>
                            {chat.other_user &&
                              onlineUsers.includes(chat.other_user.id) && (
                                <div className="badge badge-success badge-xs">
                                  online
                                </div>
                              )}
                          </div>

                          <div className="flex justify-between items-center">
                            <p className="text-xs opacity-70 truncate">
                              {getChatSubtitle(chat)}
                            </p>
                            {chat.listing_id > 0 &&
                              chat.listing?.price !== undefined && (
                                <span className="badge badge-ghost badge-sm">
                                  {new Intl.NumberFormat(
                                    locale === 'ru' ? 'ru-RU' : 'en-US',
                                    {
                                      style: 'currency',
                                      currency: 'RSD',
                                      minimumFractionDigits: 0,
                                      maximumFractionDigits: 0,
                                    }
                                  ).format(chat.listing.price)}
                                </span>
                              )}
                          </div>

                          {chat.last_message && (
                            <div className="flex items-center gap-1 mt-1">
                              {chat.last_message.sender_id === user?.id && (
                                <span className="text-xs text-primary">✓</span>
                              )}
                              <p className="text-xs opacity-50 truncate">
                                {chat.last_message.content}
                              </p>
                            </div>
                          )}
                        </div>
                      </div>
                    </div>
                  </div>
                ))}

                {/* Индикатор загрузки следующей страницы */}
                {hasMoreChats && (
                  <div ref={loadMoreRef} className="flex justify-center p-4">
                    {isLoading && (
                      <span className="loading loading-spinner loading-md text-primary"></span>
                    )}
                  </div>
                )}
              </>
            )}
          </>
        ) : (
          // Список контактов с card компонентами
          <>
            {contactsLoading ? (
              <div className="flex justify-center p-8">
                <span className="loading loading-spinner loading-lg text-secondary"></span>
              </div>
            ) : filteredContacts.length === 0 ? (
              <div className="card bg-base-100 shadow-sm">
                <div className="card-body text-center py-12">
                  <div className="text-4xl mb-4">👥</div>
                  <h3 className="card-title justify-center text-base-content/70">
                    {searchQuery ? t('noSearchResults') : t('noContacts')}
                  </h3>
                  <p className="text-base-content/50">
                    {searchQuery
                      ? t('tryDifferentSearch')
                      : t('addContactsMessage')}
                  </p>
                </div>
              </div>
            ) : (
              filteredContacts.map((contact) => (
                <div
                  key={contact.id}
                  onClick={() => handleContactSelect(contact)}
                  className="card card-compact bg-base-100 hover:bg-base-200 cursor-pointer shadow-sm transition-all"
                >
                  <div className="card-body">
                    <div className="flex items-center gap-4">
                      {/* Аватар контакта */}
                      <div className="avatar placeholder">
                        <div className="bg-neutral text-neutral-content rounded-full w-12">
                          <span className="text-lg">
                            {contact.contact_user?.name
                              ?.charAt(0)
                              .toUpperCase() || '?'}
                          </span>
                        </div>
                      </div>

                      {/* Информация о контакте */}
                      <div className="flex-1 min-w-0">
                        <h3 className="font-semibold text-sm truncate">
                          {contact.contact_user?.name ||
                            `User #${contact.contact_user_id}`}
                        </h3>
                        <p className="text-xs opacity-70 truncate">
                          {contact.contact_user?.email}
                        </p>
                        {contact.notes && (
                          <p className="text-xs opacity-50 truncate mt-1">
                            {contact.notes}
                          </p>
                        )}
                      </div>

                      {/* Статус онлайн */}
                      {contact.contact_user &&
                        onlineUsers.includes(contact.contact_user.id) && (
                          <div className="badge badge-success badge-xs">
                            online
                          </div>
                        )}
                    </div>
                  </div>
                </div>
              ))
            )}
          </>
        )}
      </div>
    </div>
  );
}
