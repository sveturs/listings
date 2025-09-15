'use client';

import { useState, useEffect } from 'react';
import { useTranslations } from 'next-intl';
import { contactsService } from '@/services/contacts';
import { toast } from '@/utils/toast';
import { formatDistanceToNow } from 'date-fns';
import { ru, enUS, sr } from 'date-fns/locale';
import { useParams } from 'next/navigation';

interface ContactRequest {
  id: number;
  user_id: number;
  contact_user_id: number;
  status: string;
  notes?: string;
  added_from_chat_id?: number;
  created_at: string;
  user?: {
    id: number;
    name: string;
    email: string;
    picture_url?: string;
  };
}

interface IncomingContactRequestProps {
  otherUserId: number;
  chatId?: number;
  onRequestHandled?: () => void;
}

export default function IncomingContactRequest({
  otherUserId,
  chatId,
  onRequestHandled,
}: IncomingContactRequestProps) {
  const t = useTranslations('chat');
  const params = useParams();
  const locale = params?.locale as string;
  const [request, setRequest] = useState<ContactRequest | null>(null);
  const [isLoading, setIsLoading] = useState(true);
  const [isHandling, setIsHandling] = useState(false);

  // Загружаем входящие запросы при монтировании
  useEffect(() => {
    const loadIncomingRequests = async () => {
      try {
        setIsLoading(true);
        const response = await contactsService.getIncomingRequests({
          page: 1,
          limit: 100,
        });

        // Ищем запрос от этого пользователя
        const foundRequest = response.contacts.find(
          (req) => req.user_id === otherUserId
        );

        if (foundRequest) {
          setRequest(foundRequest);
        }
      } catch (error) {
        console.error('Error loading incoming requests:', error);
      } finally {
        setIsLoading(false);
      }
    };

    loadIncomingRequests();
  }, [otherUserId]);

  const handleAccept = async () => {
    if (!request || isHandling) return;

    setIsHandling(true);
    try {
      await contactsService.updateContactStatus(request.user_id, {
        status: 'accepted',
        notes: `Accepted from chat ${chatId ? `#${chatId}` : ''}`,
      });

      toast.success(t('contactRequestAccepted'));
      setRequest(null);
      onRequestHandled?.();
    } catch (error) {
      console.error('Error accepting contact request:', error);
      toast.error(t('failedToAcceptRequest'));
    } finally {
      setIsHandling(false);
    }
  };

  const handleReject = async () => {
    if (!request || isHandling) return;

    setIsHandling(true);
    try {
      await contactsService.updateContactStatus(request.user_id, {
        status: 'blocked',
        notes: `Rejected from chat ${chatId ? `#${chatId}` : ''}`,
      });

      toast.info(t('contactRequestRejected'));
      setRequest(null);
      onRequestHandled?.();
    } catch (error) {
      console.error('Error rejecting contact request:', error);
      toast.error(t('failedToRejectRequest'));
    } finally {
      setIsHandling(false);
    }
  };

  // Если нет запроса или загружается - не показываем
  if (isLoading || !request) {
    return null;
  }

  const getLocale = () => {
    switch (locale) {
      case 'ru':
        return ru;
      case 'sr':
        return sr;
      default:
        return enUS;
    }
  };

  const timeAgo = formatDistanceToNow(new Date(request.created_at), {
    addSuffix: true,
    locale: getLocale(),
  });

  return (
    <div className="alert alert-info shadow-lg mb-3 mx-3 sm:mx-4 lg:mx-8">
      <div className="flex flex-col sm:flex-row items-start sm:items-center justify-between w-full gap-3">
        <div className="flex items-start gap-3">
          <svg
            className="w-5 h-5 sm:w-6 sm:h-6 flex-shrink-0 mt-1"
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
          <div className="flex-1">
            <div className="font-semibold text-sm sm:text-base">
              {t('incomingContactRequest')}
            </div>
            <div className="text-xs sm:text-sm opacity-90">
              {request.user?.name || t('unknownUser')} • {timeAgo}
            </div>
            {request.notes && (
              <div className="text-xs opacity-80 mt-1 italic">
                "{request.notes}"
              </div>
            )}
          </div>
        </div>
        <div className="flex gap-2 self-end sm:self-center">
          <button
            className="btn btn-success btn-xs sm:btn-sm"
            onClick={handleAccept}
            disabled={isHandling}
          >
            {isHandling ? (
              <span className="loading loading-spinner loading-xs"></span>
            ) : (
              t('accept')
            )}
          </button>
          <button
            className="btn btn-ghost btn-xs sm:btn-sm"
            onClick={handleReject}
            disabled={isHandling}
          >
            {t('reject')}
          </button>
        </div>
      </div>
    </div>
  );
}