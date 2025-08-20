'use client';

import { useState, useEffect } from 'react';
import { useChat } from '@/hooks/useChat';
import { useTranslations } from 'next-intl';
import ChatList from './ChatList';
import ChatWindow from './ChatWindow';
import { MarketplaceChat } from '@/types/chat';

interface ChatLayoutProps {
  initialListingId?: number;
  initialStorefrontProductId?: number;
  initialSellerId?: number;
  initialContactId?: number;
}

export default function ChatLayout({
  initialListingId,
  initialStorefrontProductId,
  initialSellerId,
  initialContactId,
}: ChatLayoutProps) {
  const t = useTranslations('chat');
  const [isMobileListOpen, setIsMobileListOpen] = useState(true);
  const [isMobile, setIsMobile] = useState(false);
  const { currentChat, setCurrentChat } = useChat();

  useEffect(() => {
    const checkMobile = () => {
      setIsMobile(window.innerWidth < 768);
    };
    checkMobile();
    window.addEventListener('resize', checkMobile);
    return () => window.removeEventListener('resize', checkMobile);
  }, []);

  const handleChatSelect = (chat: MarketplaceChat) => {
    setCurrentChat(chat);
    // На мобильных устройствах скрываем список при выборе чата
    if (isMobile) {
      setIsMobileListOpen(false);
    }
  };

  const handleBackToList = () => {
    setIsMobileListOpen(true);
  };

  return (
    <div className="flex h-full bg-base-200 rounded-lg overflow-hidden max-w-full">
      {/* Боковая панель с чатами - на десктопе всегда видна, на мобильном - условно */}
      <div
        className={`
        w-full md:w-1/3 lg:w-1/4 min-w-0
        ${!isMobileListOpen && 'hidden md:block'}
        border-r border-base-300
      `}
      >
        <ChatList onChatSelect={handleChatSelect} />
      </div>

      {/* Окно чата - на десктопе всегда видно, на мобильном - когда чат выбран */}
      <div
        className={`
        flex-1 min-w-0
        ${isMobileListOpen && 'hidden md:flex'}
        flex flex-col
      `}
      >
        {currentChat ? (
          <ChatWindow
            chat={currentChat}
            initialContactId={
              currentChat.id === 0 && currentChat.seller_id
                ? currentChat.seller_id
                : undefined
            }
            onBack={handleBackToList}
            showBackButton={isMobile}
          />
        ) : (initialListingId || initialStorefrontProductId) &&
          initialSellerId ? (
          <ChatWindow
            initialListingId={initialListingId}
            initialStorefrontProductId={initialStorefrontProductId}
            initialSellerId={initialSellerId}
            onBack={handleBackToList}
            showBackButton={isMobile}
          />
        ) : initialContactId ? (
          <ChatWindow
            initialContactId={initialContactId}
            onBack={handleBackToList}
            showBackButton={isMobile}
          />
        ) : (
          <div className="flex items-center justify-center h-full text-base-content/50">
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
                  strokeWidth={1.5}
                  d="M8 12h.01M12 12h.01M16 12h.01M21 12c0 4.418-4.03 8-9 8a9.863 9.863 0 01-4.255-.949L3 20l1.395-3.72C3.512 15.042 3 13.574 3 12c0-4.418 4.03-8 9-8s9 3.582 9 8z"
                />
              </svg>
              <p className="text-lg">{t('selectChat')}</p>
            </div>
          </div>
        )}
      </div>
    </div>
  );
}
