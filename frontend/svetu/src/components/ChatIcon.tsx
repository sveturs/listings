'use client';

import { Link } from '@/i18n/routing';
import { useParams } from 'next/navigation';
import { FiMessageCircle } from 'react-icons/fi';
import { useAppSelector } from '@/store/hooks';
import { useTranslations } from 'next-intl';

export default function ChatIcon() {
  const params = useParams();
  const _locale = params?.locale || 'en';
  const unreadCount = useAppSelector((state) => state.chat.unreadCount);
  const t = useTranslations('common');

  return (
    <Link
      href="/chat"
      className="btn btn-ghost btn-circle relative hidden sm:inline-flex"
      aria-label={t('chat.openChat')}
    >
      <FiMessageCircle className="w-5 h-5" aria-hidden="true" />
      {unreadCount > 0 && (
        <span className="badge badge-sm badge-error absolute -top-1 -right-1">
          {unreadCount > 99 ? '99+' : unreadCount}
        </span>
      )}
    </Link>
  );
}
