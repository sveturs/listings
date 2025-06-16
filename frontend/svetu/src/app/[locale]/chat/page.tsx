'use client';

import { useEffect, useState } from 'react';
import { useChat } from '@/hooks/useChat';
import { useAuth } from '@/contexts/AuthContext';
import { useLocale } from 'next-intl';
import { useRouter, useSearchParams } from 'next/navigation';
import ChatLayout from '@/components/Chat/ChatLayout';

export default function ChatPage() {
  const locale = useLocale();
  const router = useRouter();
  const searchParams = useSearchParams();
  const { user, isLoading: authLoading } = useAuth();
  const [mounted, setMounted] = useState(false);
  const { loadChats, chats, setCurrentChat, selectLatestChat, pendingChatId } =
    useChat();
  const [chatInitialized, setChatInitialized] = useState(false);

  // Параметры для создания нового чата
  const listingId = searchParams.get('listing_id');
  const sellerId = searchParams.get('seller_id');

  useEffect(() => {
    setMounted(true);
  }, []);

  // Автоматически выбираем чат после его создания
  useEffect(() => {
    if (pendingChatId && chats.length > 0) {
      const newChat = chats.find((c) => c.id === pendingChatId);
      if (newChat) {
        setCurrentChat(newChat);
        setChatInitialized(true);
      }
    }
  }, [pendingChatId, chats, setCurrentChat]);

  useEffect(() => {
    if (!authLoading && !user) {
      router.push('/');
      return;
    }

    if (user) {
      console.log('User authenticated, loading chats for user:', user);
      // Инициализация чата
      loadChats();
      // WebSocket теперь инициализируется глобально в WebSocketManager
    }

    return () => {
      // Не закрываем WebSocket при размонтировании, он управляется глобально
    };
  }, [user, authLoading, loadChats, router]);

  // Обработка параметров URL для открытия/создания чата
  useEffect(() => {
    if (!user || chatInitialized) return;

    // Если есть параметры URL для конкретного чата
    if (listingId && sellerId) {
      // Проверяем, что это не собственное объявление
      if (user.id.toString() === sellerId) {
        router.replace(`/${locale}/chat`);
        return;
      }

      // Ищем существующий чат с этим продавцом по этому объявлению
      const existingChat = chats.find(
        (chat) =>
          chat.listing_id === parseInt(listingId) &&
          ((chat.buyer_id === user.id &&
            chat.seller_id === parseInt(sellerId)) ||
            (chat.seller_id === user.id &&
              chat.buyer_id === parseInt(sellerId)))
      );

      if (existingChat) {
        // Если чат существует, выбираем его
        setCurrentChat(existingChat);
        setChatInitialized(true);
        router.replace(`/${locale}/chat`);
      } else {
        // Если чата нет, оставляем параметры для создания нового чата
        setChatInitialized(true);
        // НЕ очищаем URL, чтобы параметры остались для ChatLayout
      }
    } else if (chats.length > 0 && !chatInitialized) {
      // Если нет параметров URL и есть чаты, открываем самый свежий
      console.log('Selecting latest chat from:', chats);
      selectLatestChat();
      setChatInitialized(true);
    }
  }, [
    user,
    listingId,
    sellerId,
    chats,
    chatInitialized,
    locale,
    router,
    setCurrentChat,
    selectLatestChat,
  ]);

  // Закомментирован автоматический выбор последнего чата при обновлении списка
  // Теперь новые сообщения просто поднимают чат вверх, но не переключают фокус
  /*
  useEffect(() => {
    if (chats.length > 0 && !listingId && !sellerId && chatInitialized) {
      // Автоматически выбираем самый свежий чат при обновлении списка
      selectLatestChat();
    }
  }, [chats, chatInitialized, listingId, sellerId, selectLatestChat]);
  */

  if (!mounted || authLoading || !user) {
    return (
      <div className="flex items-center justify-center h-screen">
        <span className="loading loading-spinner loading-lg"></span>
      </div>
    );
  }

  return (
    <div className="absolute inset-0 top-16 flex flex-col">
      {/* Chat content */}
      <div className="flex-1 overflow-hidden px-2 sm:px-4 pb-2">
        <ChatLayout
          initialListingId={listingId ? parseInt(listingId) : undefined}
          initialSellerId={sellerId ? parseInt(sellerId) : undefined}
        />
      </div>
    </div>
  );
}
