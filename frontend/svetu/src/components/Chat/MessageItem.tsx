'use client';

import { MarketplaceMessage } from '@/types/chat';
import { format } from 'date-fns';
import { ru, enUS } from 'date-fns/locale';
import { useLocale, useTranslations } from 'next-intl';
import dynamic from 'next/dynamic';
import { ChatAttachments } from '@/components/Chat/ChatAttachments';
import { useChat } from '@/hooks/useChat';
import { sanitizeText } from '@/utils/sanitize';
import Image from 'next/image';
import { useState } from 'react';

// Динамически импортируем AnimatedEmoji, чтобы избежать проблем с SSR
const AnimatedEmoji = dynamic(() => import('./AnimatedEmoji'), {
  ssr: false,
  loading: () => null,
});

interface MessageItemProps {
  message: MarketplaceMessage;
  isOwn: boolean;
}

// Функция для проверки, является ли текст одиночным эмодзи
const isOnlyEmoji = (text: string) => {
  const trimmed = text.trim();

  // Проверяем, что текст содержит только эмодзи (1-3 эмодзи подряд)
  const emojiOnlyRegex = /^(\p{Emoji_Presentation}|\p{Emoji}\uFE0F){1,3}$/u;

  return emojiOnlyRegex.test(trimmed);
};

// Список эмодзи, для которых у нас есть анимации
const animatedEmojis = [
  '😀',
  '😊',
  '❤️',
  '🔥',
  '👍',
  '😂',
  '🎉',
  '💕',
  '🥰',
  '😍',
  '🤗',
  '😘',
  '🙂',
  '😎',
  '😭',
  '😢',
  '😅',
  '🤔',
  '😱',
  '🤯',
  '😴',
  '🤩',
  '🥳',
  '🙏',
  '👌',
  '✌️',
  '🤞',
  '💪',
  '👏',
  '🙌',
];

// Функция для генерации инициалов из имени
const getInitials = (name: string) => {
  return name
    .split(' ')
    .map((word) => word.charAt(0))
    .join('')
    .toUpperCase()
    .slice(0, 2);
};

export default function MessageItem({ message, isOwn }: MessageItemProps) {
  const locale = useLocale();
  const t = useTranslations('chat');
  const isEmojiOnly = isOnlyEmoji(message.content);
  const { deleteAttachment } = useChat();

  // ✅ УПРОЩЕНО: Просто показываем готовый перевод из backend (БЕЗ API запросов!)
  const [showOriginal, setShowOriginal] = useState(false);

  // Проверяем готовый перевод из backend
  const hasTranslation = message.translations && message.translations[locale];

  // Отображаемый текст:
  // - Если это своё сообщение (isOwn), всегда показываем оригинал
  // - Если чужое сообщение, показываем перевод (если есть) или оригинал
  const displayText = isOwn
    ? message.content
    : showOriginal
      ? message.content
      : hasTranslation || message.content;

  // Кнопка переключения только если есть перевод
  const shouldShowToggleButton = !isOwn && !isEmojiOnly && hasTranslation;

  const formatTime = (date: string) => {
    return format(new Date(date), 'HH:mm', {
      locale: locale === 'ru' ? ru : enUS,
    });
  };

  const handleDeleteAttachment = async (attachmentId: number) => {
    try {
      await deleteAttachment(attachmentId);
    } catch (error) {
      console.error('Error deleting attachment:', error);
    }
  };

  // Современный DaisyUI chat bubble компонент
  return (
    <div className={`chat ${isOwn ? 'chat-end' : 'chat-start'} mb-2`}>
      {/* Аватар отправителя */}
      <div className="chat-image avatar">
        <div className="w-10 rounded-full">
          {!isOwn && message.sender?.picture_url ? (
            <Image
              src={message.sender.picture_url}
              alt={message.sender?.name || 'User'}
              width={40}
              height={40}
              className="rounded-full object-cover"
            />
          ) : !isOwn ? (
            <div className="w-10 h-10 rounded-full bg-neutral flex items-center justify-center">
              <span className="text-sm font-semibold text-neutral-content">
                {getInitials(message.sender?.name || 'U')}
              </span>
            </div>
          ) : null}
        </div>
      </div>

      {/* Контейнер сообщения */}
      <div className="chat-bubble-container">
        {/* Если есть вложения и текст - объединяем их */}
        {message.attachments &&
        message.attachments.length > 0 &&
        message.content &&
        !(
          message.content === message.attachments[0]?.file_name ||
          message.content.match(/^\d+ файла\(ов\)$/)
        ) ? (
          <div className="inline-block w-48">
            {/* Вложения */}
            <ChatAttachments
              attachments={message.attachments}
              onDelete={isOwn ? handleDeleteAttachment : undefined}
              canDelete={isOwn}
              hasTextBelow={true}
            />
            {/* Текст прилепленный снизу - точно под картинкой */}
            <div
              className={`chat-bubble ${isOwn ? 'chat-bubble-primary' : 'chat-bubble-accent'} mt-1`}
            >
              <p
                className="whitespace-pre-wrap"
                dangerouslySetInnerHTML={{
                  __html: sanitizeText(displayText),
                }}
              />
            </div>
          </div>
        ) : (
          <>
            {/* Только вложения */}
            {message.attachments && message.attachments.length > 0 && (
              <div className="w-full">
                <ChatAttachments
                  attachments={message.attachments}
                  onDelete={isOwn ? handleDeleteAttachment : undefined}
                  canDelete={isOwn}
                />
              </div>
            )}

            {/* Только текст */}
            {message.content &&
              !(
                message.attachments &&
                message.attachments.length > 0 &&
                (message.content === message.attachments[0]?.file_name ||
                  message.content.match(/^\d+ файла\(ов\)$/))
              ) &&
              !message.attachments?.length && (
                <div
                  className={`${
                    isEmojiOnly
                      ? 'text-6xl'
                      : `chat-bubble ${isOwn ? 'chat-bubble-primary' : 'chat-bubble-accent'}`
                  }`}
                >
                  {isEmojiOnly ? (
                    animatedEmojis.includes(message.content.trim()) ? (
                      <AnimatedEmoji emoji={message.content.trim()} size={64} />
                    ) : (
                      <span>{sanitizeText(message.content)}</span>
                    )
                  ) : (
                    <>
                      <p className="whitespace-pre-wrap">
                        <span
                          dangerouslySetInnerHTML={{
                            __html: sanitizeText(displayText),
                          }}
                        />
                      </p>
                    </>
                  )}
                </div>
              )}

            {/* ✅ УПРОЩЕНО: Кнопка только для переключения (БЕЗ API!) */}
            {shouldShowToggleButton && (
              <div className="mt-1 flex items-center gap-2">
                <button
                  onClick={() => setShowOriginal(!showOriginal)}
                  className="btn btn-xs btn-ghost text-base-content/50 hover:text-base-content/80 opacity-60 hover:opacity-100 transition-opacity"
                >
                  {showOriginal
                    ? t('translation.showTranslation')
                    : t('translation.showOriginal')}
                </button>
              </div>
            )}
          </>
        )}
      </div>

      {/* Время и статус прочтения */}
      <div className="chat-footer opacity-50 text-xs">
        <time>{formatTime(message.created_at)}</time>
        {isOwn && <span className="ml-1">{message.is_read ? '✓✓' : '✓'}</span>}
      </div>
    </div>
  );
}
