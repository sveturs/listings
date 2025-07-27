'use client';

import { MarketplaceMessage } from '@/types/chat';
import { format } from 'date-fns';
import { ru, enUS } from 'date-fns/locale';
import { useLocale } from 'next-intl';
import dynamic from 'next/dynamic';
import { ChatAttachments } from '@/components/Chat/ChatAttachments';
import { useChat } from '@/hooks/useChat';
import DOMPurify from 'isomorphic-dompurify';
import Image from 'next/image';

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
  const isEmojiOnly = isOnlyEmoji(message.content);
  const { deleteAttachment } = useChat();

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
    <div
      className={`chat ${isOwn ? 'chat-end' : 'chat-start'} mb-1 w-full max-w-full ${
        isOwn ? 'lg:pr-[5cm]' : 'lg:pl-[5cm]'
      }`}
    >
      {/* Аватар отправителя (только для чужих сообщений) */}
      {!isOwn && (
        <div className="chat-image avatar">
          <div className="w-10 rounded-full bg-base-300 flex items-center justify-center">
            {message.sender?.picture_url ? (
              <Image
                src={message.sender.picture_url}
                alt={message.sender?.name || 'User'}
                fill
                className="rounded-full object-cover"
              />
            ) : (
              <span className="text-sm font-semibold text-base-content">
                {getInitials(message.sender?.name || 'U')}
              </span>
            )}
          </div>
        </div>
      )}

      {/* Контейнер сообщения */}
      <div className="max-w-[95%] sm:max-w-[85%] md:max-w-[70%]">
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
              className={`relative mt-[-8px] rounded-b-lg px-3 py-1.5 text-sm sm:text-base transition-all duration-200 ease-in-out hover:shadow-md ${
                isOwn
                  ? 'bg-blue-500/80 text-white photo-caption-own'
                  : 'bg-green-100/70 text-base-content border border-green-300/70 photo-caption-with-border'
              }`}
              style={{ width: 'calc(100% - 4px)', marginLeft: '2px' }}
            >
              <div className="message-content-wrapper">
                <span
                  className={`message-time-status text-xs whitespace-nowrap ${isOwn ? 'text-white/70' : 'text-base-content/50'}`}
                >
                  {formatTime(message.created_at)}
                  {isOwn &&
                    (message.is_read ? (
                      <span className="ml-1 text-blue-200">✓✓</span>
                    ) : (
                      <span className="ml-1 text-white/60">✓</span>
                    ))}
                </span>
                <p
                  className="whitespace-pre-wrap break-words inline"
                  dangerouslySetInnerHTML={{
                    __html: DOMPurify.sanitize(message.content, {
                      ALLOWED_TAGS: [], // Не разрешаем никакие HTML теги
                      KEEP_CONTENT: true, // Сохраняем текстовое содержимое
                    }),
                  }}
                />
              </div>
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
                  className={`chat-bubble transition-all duration-200 ease-in-out ${
                    isEmojiOnly
                      ? 'bg-transparent p-0 shadow-none hover:shadow-none'
                      : isOwn
                        ? 'bg-blue-500/80 text-white hover:shadow-md'
                        : 'bg-green-100/70 text-base-content with-border hover:shadow-md'
                  }`}
                  style={{ minWidth: 'min-content' }}
                >
                  {isEmojiOnly ? (
                    <div className="leading-none">
                      {animatedEmojis.includes(message.content.trim()) ? (
                        <AnimatedEmoji
                          emoji={message.content.trim()}
                          size={64}
                        />
                      ) : (
                        <span className="text-4xl">
                          {DOMPurify.sanitize(message.content, {
                            ALLOWED_TAGS: [],
                            KEEP_CONTENT: true,
                          })}
                        </span>
                      )}
                    </div>
                  ) : (
                    <div className="message-content-wrapper">
                      <span
                        className={`message-time-status text-xs whitespace-nowrap ${isOwn ? 'text-white/70' : 'text-base-content/50'}`}
                      >
                        {formatTime(message.created_at)}
                        {isOwn &&
                          (message.is_read ? (
                            <span className="ml-1 text-blue-200">✓✓</span>
                          ) : (
                            <span className="ml-1 text-white/60">✓</span>
                          ))}
                      </span>
                      <p className="whitespace-pre-wrap break-normal text-sm sm:text-base inline">
                        <span
                          dangerouslySetInnerHTML={{
                            __html: DOMPurify.sanitize(message.content, {
                              ALLOWED_TAGS: [],
                              KEEP_CONTENT: true,
                            }),
                          }}
                        />
                      </p>
                    </div>
                  )}
                </div>
              )}
          </>
        )}
      </div>
    </div>
  );
}
