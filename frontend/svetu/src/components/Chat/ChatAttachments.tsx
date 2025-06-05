'use client';

import { ChatAttachment } from '@/types/chat';
import { useState } from 'react';
import Image from 'next/image';
import { FaFile, FaVideo, FaDownload, FaTrash, FaExpand } from 'react-icons/fa';
import configManager from '@/config';
import { useTranslations } from 'next-intl';

interface ChatAttachmentsProps {
  attachments: ChatAttachment[];
  onDelete?: (attachmentId: number) => void;
  canDelete?: boolean;
  hasTextBelow?: boolean;
}

export function ChatAttachments({
  attachments,
  onDelete,
  canDelete = false,
  hasTextBelow = false,
}: ChatAttachmentsProps) {
  const [selectedImage, setSelectedImage] = useState<string | null>(null);
  const t = useTranslations('Chat');

  const getImageUrl = (url: string) => {
    if (url.startsWith('/chat-files/')) {
      const fullUrl = configManager.buildImageUrl(url);
      console.log('Chat attachment URL:', url, '->', fullUrl);
      return fullUrl;
    }
    return url;
  };

  const formatFileSize = (bytes: number) => {
    if (bytes === 0) return '0 Bytes';
    const k = 1024;
    const sizes = ['Bytes', 'KB', 'MB', 'GB'];
    const i = Math.floor(Math.log(bytes) / Math.log(k));
    return parseFloat((bytes / Math.pow(k, i)).toFixed(2)) + ' ' + sizes[i];
  };

  const getFileIcon = (fileType: string) => {
    switch (fileType) {
      case 'image':
        return null; // Показываем превью для изображений
      case 'video':
        return <FaVideo className="text-4xl text-purple-500" />;
      case 'document':
        return <FaFile className="text-4xl text-blue-500" />;
      default:
        return <FaFile className="text-4xl text-gray-500" />;
    }
  };

  const handleDelete = async (attachmentId: number) => {
    if (onDelete && window.confirm(t('deleteAttachment'))) {
      onDelete(attachmentId);
    }
  };

  return (
    <>
      <div
        className={
          attachments.length === 1
            ? 'inline-block'
            : 'flex flex-col sm:flex-row sm:flex-wrap gap-2'
        }
      >
        {attachments.map((attachment) => (
          <div key={attachment.id} className="relative group">
            {attachment.file_type === 'image' ? (
              <div
                className={`relative w-48 h-48 border-2 border-base-300/50 overflow-hidden transition-all duration-200 ease-in-out hover:shadow-md ${hasTextBelow ? 'rounded-t-lg' : 'rounded-lg'}`}
              >
                <Image
                  src={getImageUrl(attachment.public_url)}
                  alt={attachment.file_name}
                  fill
                  sizes="(max-width: 640px) 100vw, 192px"
                  className="object-cover cursor-pointer"
                  onClick={() =>
                    setSelectedImage(getImageUrl(attachment.public_url))
                  }
                  unoptimized={attachment.public_url.startsWith('/chat-files/')}
                />
                {/* Панель действий для десктопа */}
                <div className="absolute bottom-2 right-2 opacity-0 group-hover:opacity-100 transition-all duration-200 hidden sm:block">
                  <div className="bg-white/90 dark:bg-gray-800/90 backdrop-blur-sm rounded-full shadow-lg px-2 py-1 flex items-center gap-1">
                    <button
                      onClick={() =>
                        setSelectedImage(getImageUrl(attachment.public_url))
                      }
                      className="btn btn-ghost btn-xs btn-circle hover:bg-gray-200 dark:hover:bg-gray-700"
                      title={t('expand')}
                    >
                      <FaExpand className="text-gray-700 dark:text-gray-300" />
                    </button>
                    <a
                      href={getImageUrl(attachment.public_url)}
                      download={attachment.file_name}
                      className="btn btn-ghost btn-xs btn-circle hover:bg-gray-200 dark:hover:bg-gray-700"
                      title={t('download')}
                    >
                      <FaDownload className="text-gray-700 dark:text-gray-300" />
                    </a>
                    {canDelete && (
                      <button
                        onClick={() => handleDelete(attachment.id)}
                        className="btn btn-ghost btn-xs btn-circle hover:bg-red-100 dark:hover:bg-red-900/30"
                        title={t('delete')}
                      >
                        <FaTrash className="text-red-500" />
                      </button>
                    )}
                  </div>
                </div>

                {/* Панель действий для мобильных устройств */}
                <div className="absolute bottom-2 right-2 opacity-0 active:opacity-100 transition-all duration-200 sm:hidden">
                  <div className="bg-white/90 dark:bg-gray-800/90 backdrop-blur-sm rounded-full shadow-lg px-2 py-1 flex items-center gap-1">
                    <button
                      onClick={() =>
                        setSelectedImage(getImageUrl(attachment.public_url))
                      }
                      className="btn btn-ghost btn-xs btn-circle"
                    >
                      <FaExpand className="text-gray-700 dark:text-gray-300" />
                    </button>
                    <a
                      href={getImageUrl(attachment.public_url)}
                      download={attachment.file_name}
                      className="btn btn-ghost btn-xs btn-circle"
                    >
                      <FaDownload className="text-gray-700 dark:text-gray-300" />
                    </a>
                    {canDelete && (
                      <button
                        onClick={() => handleDelete(attachment.id)}
                        className="btn btn-ghost btn-xs btn-circle"
                      >
                        <FaTrash className="text-red-500" />
                      </button>
                    )}
                  </div>
                </div>
              </div>
            ) : (
              <div className="w-full sm:w-48 h-48 p-4 flex flex-col items-center justify-center border-2 border-base-300/50 rounded-lg transition-all duration-200 ease-in-out hover:shadow-md">
                {getFileIcon(attachment.file_type)}
                <p className="text-sm font-medium mt-2 text-center truncate w-full">
                  {attachment.file_name}
                </p>
                <p className="text-xs text-base-content/60">
                  {formatFileSize(attachment.file_size)}
                </p>
                <div className="mt-2 flex gap-2">
                  <a
                    href={getImageUrl(attachment.public_url)}
                    download={attachment.file_name}
                    className="btn btn-xs btn-primary"
                  >
                    <FaDownload className="mr-1" /> {t('download')}
                  </a>
                  {canDelete && (
                    <button
                      onClick={() => handleDelete(attachment.id)}
                      className="btn btn-xs btn-error"
                    >
                      <FaTrash />
                    </button>
                  )}
                </div>
              </div>
            )}
          </div>
        ))}
      </div>

      {/* Модальное окно для просмотра изображений */}
      {selectedImage && (
        <div
          className="fixed inset-0 z-50 bg-black bg-opacity-90 flex items-center justify-center p-4"
          onClick={() => setSelectedImage(null)}
        >
          <div className="relative max-w-full max-h-full">
            <Image
              src={selectedImage}
              alt="Full size"
              width={1200}
              height={800}
              className="object-contain"
              unoptimized={selectedImage.startsWith('/chat-files/')}
            />
            <button
              onClick={() => setSelectedImage(null)}
              className="absolute top-4 right-4 btn btn-circle btn-ghost text-white"
            >
              ✕
            </button>
          </div>
        </div>
      )}
    </>
  );
}
