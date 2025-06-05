'use client';

import { useState, useRef, KeyboardEvent, useEffect } from 'react';
import { useChat } from '@/hooks/useChat';
import { useAuth } from '@/contexts/AuthContext';
import { useTranslations } from 'next-intl';
import { MarketplaceChat, MarketplaceMessage } from '@/types/chat';
import EmojiPicker from './EmojiPicker';
import { FileUploadProgress } from '@/components/Chat/FileUploadProgress';
import { toast } from '@/utils/toast';
import { validateFiles, formatFileSize } from '@/utils/fileValidation';

interface MessageInputProps {
  chat?: MarketplaceChat;
  initialListingId?: number;
  initialSellerId?: number;
}

export default function MessageInput({
  chat,
  initialListingId,
  initialSellerId,
}: MessageInputProps) {
  const t = useTranslations('Chat');
  const { user } = useAuth();
  const [message, setMessage] = useState('');
  const [showEmojiPicker, setShowEmojiPicker] = useState(false);
  const [isTyping, setIsTyping] = useState(false);
  const [selectedFiles, setSelectedFiles] = useState<File[]>([]);
  const inputRef = useRef<HTMLTextAreaElement>(null);
  const fileInputRef = useRef<HTMLInputElement>(null);
  const typingTimeoutRef = useRef<NodeJS.Timeout | null>(null);

  const {
    sendMessage,
    setUserTyping,
    uploadFiles,
    uploadingFiles,
    removeUploadingFile,
  } = useChat();

  const handleSend = async () => {
    if ((!message.trim() && selectedFiles.length === 0) || !user) return;

    // Если нет текста, но есть файлы, используем имена файлов как текст
    let messageContent = message.trim();
    if (!messageContent && selectedFiles.length > 0) {
      if (selectedFiles.length === 1) {
        messageContent = selectedFiles[0].name;
      } else {
        messageContent = `${selectedFiles.length} файла(ов)`;
      }
    }

    let payload;

    if (chat) {
      // Отправка в существующий чат
      payload = {
        chat_id: chat.id > 0 ? chat.id : undefined, // Если chat.id = 0, это виртуальный чат для контакта
        listing_id: chat.listing_id,
        receiver_id: chat.buyer_id === user.id ? chat.seller_id : chat.buyer_id,
        content: messageContent,
      };
    } else if (initialListingId && initialSellerId) {
      // Создание нового чата с объявлением
      payload = {
        listing_id: initialListingId,
        receiver_id: initialSellerId,
        content: messageContent,
      };
    } else if (initialSellerId) {
      // Прямое сообщение контакту без объявления
      payload = {
        receiver_id: initialSellerId,
        content: messageContent,
      };
    } else {
      return;
    }

    try {
      const result = await sendMessage(payload);
      const sentMessage = result.payload as MarketplaceMessage;
      setMessage('');

      // Загружаем файлы если они есть
      if (selectedFiles.length > 0 && sentMessage?.id) {
        try {
          await uploadFiles(sentMessage.id, selectedFiles);
          setSelectedFiles([]);
        } catch (uploadError: unknown) {
          const err = uploadError as Error & { status?: number };
          // Обработка ошибки загрузки файлов
          if (err.status === 429) {
            console.warn('File upload rate limit:', err.message);
            toast.warning(err.message || t('rateLimitFiles'));
          } else {
            console.error('File upload failed:', err);
            toast.error(t('uploadError'));
          }
        }
      }

      // Остановить индикатор печатания
      if (isTyping && chat) {
        setIsTyping(false);
        setUserTyping(chat.id, user.id, false);
      }

      // Если это был новый чат, получить созданный чат и выбрать его
      if (!chat && sentMessage) {
        console.log('New chat created, message sent:', sentMessage);
        // Chat selection will be handled by the parent component after loadChats updates
        // Clear URL parameters after creating chat
        window.history.replaceState({}, '', window.location.pathname);
      }
    } catch (error: unknown) {
      const err = error as Error & { status?: number };
      // Обрабатываем 429 ошибку как предупреждение
      if (err.status === 429) {
        console.warn('Rate limit exceeded:', err.message);
        toast.warning(err.message || t('rateLimitMessage'));
      } else {
        console.error('Failed to send message:', err);
        toast.error(t('sendError'));
      }
    }
  };

  const handleKeyDown = (e: KeyboardEvent<HTMLTextAreaElement>) => {
    if (e.key === 'Enter' && !e.shiftKey) {
      e.preventDefault();
      handleSend();
    }
  };

  const handleTyping = (value: string) => {
    setMessage(value);

    if (!user || !chat) return;

    // Начать показывать индикатор печатания
    if (!isTyping && value.trim()) {
      setIsTyping(true);
      setUserTyping(chat.id, user.id, true);
    }

    // Сбросить таймер
    if (typingTimeoutRef.current) {
      clearTimeout(typingTimeoutRef.current);
    }

    // Установить новый таймер для остановки индикатора
    typingTimeoutRef.current = setTimeout(() => {
      if (isTyping) {
        setIsTyping(false);
        setUserTyping(chat.id, user.id, false);
      }
    }, 3000);
  };

  const handleEmojiSelect = (emoji: string) => {
    const cursorPosition = inputRef.current?.selectionStart || message.length;
    const newMessage =
      message.slice(0, cursorPosition) + emoji + message.slice(cursorPosition);

    setMessage(newMessage);
    setShowEmojiPicker(false);
    inputRef.current?.focus();
  };

  const handleFileSelect = (e: React.ChangeEvent<HTMLInputElement>) => {
    const files = e.target.files;
    if (!files) return;

    // Валидация файлов
    const { valid, errors } = validateFiles(Array.from(files));

    // Показываем ошибки валидации
    errors.forEach((error, filename) => {
      if (filename === '_total') {
        toast.error(error);
      } else {
        toast.error(`${filename}: ${error}`);
      }
    });

    // Добавляем валидные файлы
    if (valid.length > 0) {
      // Проверяем общее количество файлов
      const totalFiles = selectedFiles.length + valid.length;
      if (totalFiles > 10) {
        toast.error('Можно прикрепить максимум 10 файлов к одному сообщению');
        return;
      }

      setSelectedFiles((prev) => [...prev, ...valid]);

      // Показываем сообщение об успешном добавлении
      if (valid.length === 1) {
        toast.success(
          `Файл "${valid[0].name}" добавлен (${formatFileSize(valid[0].size)})`
        );
      } else {
        toast.success(`Добавлено файлов: ${valid.length}`);
      }
    }

    // Очищаем input для возможности выбора тех же файлов
    if (fileInputRef.current) {
      fileInputRef.current.value = '';
    }
  };

  const removeSelectedFile = (index: number) => {
    setSelectedFiles((prev) => prev.filter((_, i) => i !== index));
  };

  // Автоматическая подстройка высоты textarea при изменении содержимого
  useEffect(() => {
    if (inputRef.current) {
      inputRef.current.style.height = 'auto';
      inputRef.current.style.height = inputRef.current.scrollHeight + 'px';
    }
  }, [message]);

  // Получаем массив загружаемых файлов
  const uploadingFilesList = Object.values(uploadingFiles);

  return (
    <div className="bg-gradient-to-r from-gray-50 to-white border-t border-gray-200 backdrop-blur-sm">
      {/* Отображение загружаемых файлов */}
      {uploadingFilesList.length > 0 && (
        <div className="p-4 pb-0">
          <div className="card bg-base-100/80 shadow-sm">
            <div className="card-body p-3">
              <FileUploadProgress
                uploadingFiles={uploadingFilesList}
                onRemove={removeUploadingFile}
              />
            </div>
          </div>
        </div>
      )}

      {/* Отображение выбранных файлов */}
      {selectedFiles.length > 0 && (
        <div className="p-4 pb-0">
          <div className="card bg-base-100/80 shadow-sm">
            <div className="card-body p-3">
              <div className="flex items-center gap-2 mb-2">
                <svg
                  className="w-4 h-4 text-blue-500"
                  fill="currentColor"
                  viewBox="0 0 20 20"
                >
                  <path
                    fillRule="evenodd"
                    d="M3 17a1 1 0 011-1h12a1 1 0 110 2H4a1 1 0 01-1-1zm3.293-7.707a1 1 0 011.414 0L9 10.586V3a1 1 0 112 0v7.586l1.293-1.293a1 1 0 111.414 1.414l-3 3a1 1 0 01-1.414 0l-3-3a1 1 0 010-1.414z"
                    clipRule="evenodd"
                  />
                </svg>
                <p className="text-sm font-medium text-blue-600">
                  {t('selectedFiles')}
                </p>
              </div>
              <div className="flex flex-wrap gap-2">
                {selectedFiles.map((file, index) => {
                  const isImage = file.type.startsWith('image/');
                  const fileUrl = URL.createObjectURL(file);

                  return (
                    <div
                      key={index}
                      className={
                        isImage
                          ? 'relative group'
                          : 'card card-sm bg-base-200/50 shadow-sm flex-1'
                      }
                    >
                      {isImage ? (
                        // Эскиз для изображений
                        <div className="relative w-20 h-20 sm:w-24 sm:h-24 rounded-lg overflow-hidden border-2 border-base-300/50">
                          <img
                            src={fileUrl}
                            alt={file.name}
                            className="w-full h-full object-cover"
                            onLoad={() => URL.revokeObjectURL(fileUrl)}
                          />
                          <button
                            onClick={() => removeSelectedFile(index)}
                            className="absolute top-1 right-1 btn btn-xs btn-circle bg-black/50 hover:bg-error/80 border-none text-white"
                          >
                            ✕
                          </button>
                          {/* Tooltip с названием файла */}
                          <div className="absolute bottom-0 left-0 right-0 bg-black/70 text-white text-xs p-1 truncate opacity-0 group-hover:opacity-100 transition-opacity">
                            {file.name}
                          </div>
                        </div>
                      ) : (
                        // Обычный вид для других файлов
                        <div className="card-body p-2 flex-row items-center justify-between">
                          <div className="flex items-center gap-2">
                            <div className="w-6 h-6 rounded bg-blue-100 flex items-center justify-center">
                              <svg
                                className="w-3 h-3 text-blue-500"
                                fill="currentColor"
                                viewBox="0 0 20 20"
                              >
                                <path
                                  fillRule="evenodd"
                                  d="M3 17a1 1 0 011-1h12a1 1 0 110 2H4a1 1 0 01-1-1zM6.293 6.707a1 1 0 010-1.414l3-3a1 1 0 011.414 0l3 3a1 1 0 01-1.414 1.414L11 5.414V13a1 1 0 11-2 0V5.414L7.707 6.707a1 1 0 01-1.414 0z"
                                  clipRule="evenodd"
                                />
                              </svg>
                            </div>
                            <span className="text-sm truncate">
                              {file.name}
                            </span>
                          </div>
                          <button
                            onClick={() => removeSelectedFile(index)}
                            className="btn btn-ghost btn-xs btn-circle hover:bg-error/20 hover:text-error"
                          >
                            ✕
                          </button>
                        </div>
                      )}
                    </div>
                  );
                })}
              </div>
            </div>
          </div>
        </div>
      )}

      <div className="p-2 sm:p-4">
        <div className="flex items-end gap-2">
          {/* Кнопка файлов */}
          <input
            ref={fileInputRef}
            type="file"
            multiple
            onChange={handleFileSelect}
            className="hidden"
            accept="image/*,video/*,.pdf,.doc,.docx,.txt"
          />
          <button
            onClick={() => fileInputRef.current?.click()}
            className="btn btn-circle btn-ghost hover:bg-violet-100 hover:text-violet-600 transition-all duration-200"
            title={t('attachFile')}
          >
            <svg
              className="w-5 h-5"
              fill="none"
              stroke="currentColor"
              viewBox="0 0 24 24"
              strokeWidth={2}
            >
              <path
                strokeLinecap="round"
                strokeLinejoin="round"
                d="M18.375 12.739l-7.693 7.693a4.5 4.5 0 01-6.364-6.364l10.94-10.94A3 3 0 1119.5 7.372L8.552 18.32m.009-.01l-.01.01m5.699-9.941l-7.81 7.81a1.5 1.5 0 002.112 2.13"
              />
            </svg>
          </button>

          {/* Поле ввода с эмодзи внутри */}
          <div className="flex-1 relative">
            <div className="relative bg-base-100 rounded-2xl border border-base-300 shadow-sm hover:shadow-md transition-shadow">
              <textarea
                ref={inputRef}
                value={message}
                onChange={(e) => {
                  handleTyping(e.target.value);
                  // Автоматическая высота
                  e.target.style.height = 'auto';
                  e.target.style.height = e.target.scrollHeight + 'px';
                }}
                onKeyDown={handleKeyDown}
                placeholder={t('messagePlaceholder')}
                className="textarea w-full resize-none border-0 focus:outline-none bg-transparent placeholder:text-base-content/50 pr-12"
                rows={1}
                style={{
                  minHeight: '44px',
                  maxHeight: '120px',
                  paddingTop: '0.75rem',
                  paddingBottom: '0.75rem',
                }}
              />

              {/* Кнопка эмодзи внутри поля */}
              <div className="absolute right-2 bottom-2">
                <button
                  onClick={() => setShowEmojiPicker(!showEmojiPicker)}
                  className="btn btn-ghost btn-sm btn-circle hover:bg-amber-100 hover:text-amber-600 transition-all duration-200"
                  title={t('addEmoji')}
                >
                  <span className="text-xl">😊</span>
                </button>
              </div>
            </div>

            {showEmojiPicker && (
              <div className="absolute bottom-full mb-2 right-0 z-[9999]">
                <EmojiPicker
                  onSelect={handleEmojiSelect}
                  onClose={() => setShowEmojiPicker(false)}
                />
              </div>
            )}
          </div>

          {/* Кнопка отправки */}
          <button
            onClick={handleSend}
            disabled={!message.trim() && selectedFiles.length === 0}
            className={`btn btn-square btn-ghost transition-all duration-200 ${
              !message.trim() && selectedFiles.length === 0
                ? 'text-base-content/30'
                : 'hover:bg-blue-100 hover:text-blue-600 text-blue-500'
            }`}
            title={t('sendMessage')}
          >
            <svg
              className="w-5 h-5"
              fill="none"
              stroke="currentColor"
              viewBox="0 0 24 24"
              strokeWidth={2}
            >
              <path
                strokeLinecap="round"
                strokeLinejoin="round"
                d="M6 12L3.269 3.126A59.768 59.768 0 0121.485 12 59.77 59.77 0 013.27 20.876L5.999 12zm0 0h7.5"
              />
            </svg>
          </button>
        </div>
      </div>
    </div>
  );
}
