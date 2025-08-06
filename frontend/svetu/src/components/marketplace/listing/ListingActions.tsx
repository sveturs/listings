'use client';

import { useState } from 'react';
import { useTranslations } from 'next-intl';
import { toast } from '@/utils/toast';
import { useAuth } from '@/contexts/AuthContext';

interface ListingActionsProps {
  listing: {
    id: number;
    title: string;
    price: number;
    is_favorite?: boolean;
  };
}

export default function ListingActions({ listing }: ListingActionsProps) {
  const { user } = useAuth();
  const [isFavorite, setIsFavorite] = useState(listing.is_favorite || false);
  const [isInComparison, setIsInComparison] = useState(false);
  const [shareMenuOpen, setShareMenuOpen] = useState(false);
  const t = useTranslations('marketplace');

  const handleFavoriteToggle = async () => {
    if (!user) {
      toast.error(t('listingActions.signInToAddFavorites'));
      return;
    }

    try {
      setIsFavorite(!isFavorite);
      // TODO: API call to toggle favorite
      toast.success(
        isFavorite
          ? t('listingActions.removedFromFavorites')
          : t('listingActions.addedToFavorites')
      );
    } catch {
      setIsFavorite(!isFavorite); // Revert on error
      toast.error(t('listingActions.error'));
    }
  };

  const handleComparisonToggle = () => {
    setIsInComparison(!isInComparison);
    // TODO: Add to comparison store
    toast.success(
      isInComparison
        ? t('listingActions.removedFromComparison')
        : t('listingActions.addedToComparison')
    );
  };

  // Fallback функция для копирования в буфер обмена
  const fallbackCopyToClipboard = (text: string) => {
    const textArea = document.createElement('textarea');
    textArea.value = text;
    textArea.style.position = 'fixed';
    textArea.style.left = '-999999px';
    textArea.style.top = '-999999px';
    document.body.appendChild(textArea);
    textArea.focus();
    textArea.select();

    try {
      const successful = document.execCommand('copy');
      if (successful) {
        toast.success(t('listingActions.linkCopied'));
      } else {
        toast.error(t('listingActions.failedToCopyLink'));
      }
    } catch (err) {
      console.error('Fallback: Oops, unable to copy', err);
      toast.error(t('listingActions.copyNotSupported'));
    }

    document.body.removeChild(textArea);
    setShareMenuOpen(false);
  };

  const handleShare = (platform: string) => {
    const url = window.location.href;
    const text = `${listing.title} - ${listing.price}$`;

    const shareUrls: Record<string, string> = {
      whatsapp: `https://wa.me/?text=${encodeURIComponent(text + ' ' + url)}`,
      telegram: `https://t.me/share/url?url=${encodeURIComponent(url)}&text=${encodeURIComponent(text)}`,
      facebook: `https://www.facebook.com/sharer/sharer.php?u=${encodeURIComponent(url)}`,
      twitter: `https://twitter.com/intent/tweet?url=${encodeURIComponent(url)}&text=${encodeURIComponent(text)}`,
      viber: `viber://forward?text=${encodeURIComponent(text + ' ' + url)}`,
    };

    if (platform === 'copy') {
      // Проверяем поддержку clipboard API
      if (navigator.clipboard && window.isSecureContext) {
        navigator.clipboard
          .writeText(url)
          .then(() => {
            toast.success(t('listingActions.linkCopied'));
            setShareMenuOpen(false);
          })
          .catch(() => {
            // Fallback для старых браузеров
            fallbackCopyToClipboard(url);
          });
      } else {
        // Fallback для небезопасного контекста
        fallbackCopyToClipboard(url);
      }
    } else if (shareUrls[platform]) {
      window.open(shareUrls[platform], '_blank');
      setShareMenuOpen(false);
    }
  };

  return (
    <div className="flex items-center gap-2">
      {/* Favorite Button */}
      <button
        onClick={handleFavoriteToggle}
        className={`btn btn-circle ${isFavorite ? 'btn-error' : 'btn-ghost'}`}
        aria-label={t('listingActions.addToFavorites')}
      >
        <svg
          className="w-6 h-6"
          fill={isFavorite ? 'currentColor' : 'none'}
          viewBox="0 0 24 24"
          stroke="currentColor"
        >
          <path
            strokeLinecap="round"
            strokeLinejoin="round"
            strokeWidth={2}
            d="M4.318 6.318a4.5 4.5 0 000 6.364L12 20.364l7.682-7.682a4.5 4.5 0 00-6.364-6.364L12 7.636l-1.318-1.318a4.5 4.5 0 00-6.364 0z"
          />
        </svg>
      </button>

      {/* Comparison Button */}
      <button
        onClick={handleComparisonToggle}
        className={`btn btn-circle ${isInComparison ? 'btn-info' : 'btn-ghost'}`}
        aria-label={t('listingActions.addToComparison')}
      >
        <svg
          className="w-6 h-6"
          fill="none"
          viewBox="0 0 24 24"
          stroke="currentColor"
        >
          <path
            strokeLinecap="round"
            strokeLinejoin="round"
            strokeWidth={2}
            d="M9 19v-6a2 2 0 00-2-2H5a2 2 0 00-2 2v6a2 2 0 002 2h2a2 2 0 002-2zm0 0V9a2 2 0 012-2h2a2 2 0 012 2v10m-6 0a2 2 0 002 2h2a2 2 0 002-2m0 0V5a2 2 0 012-2h2a2 2 0 012 2v14a2 2 0 01-2 2h-2a2 2 0 01-2-2z"
          />
        </svg>
      </button>

      {/* Share Button */}
      <div className="dropdown dropdown-end">
        <button
          tabIndex={0}
          onClick={() => setShareMenuOpen(!shareMenuOpen)}
          className="btn btn-circle btn-ghost"
          aria-label={t('listingActions.share')}
        >
          <svg
            className="w-6 h-6"
            fill="none"
            viewBox="0 0 24 24"
            stroke="currentColor"
          >
            <path
              strokeLinecap="round"
              strokeLinejoin="round"
              strokeWidth={2}
              d="M8.684 13.342C8.886 12.938 9 12.482 9 12c0-.482-.114-.938-.316-1.342m0 2.684a3 3 0 110-2.684m9.632 4.316a3 3 0 00-5.367-2.684m5.367 2.684a3 3 0 01-5.367 2.684m0-5.368a3 3 0 015.367-2.684M21 12a9 9 0 11-18 0 9 9 0 0118 0z"
            />
          </svg>
        </button>

        {shareMenuOpen && (
          <ul
            tabIndex={0}
            className="dropdown-content menu p-2 shadow bg-base-200 rounded-box w-52 z-50"
          >
            <li>
              <a onClick={() => handleShare('whatsapp')}>
                <svg
                  className="w-5 h-5 text-success"
                  fill="currentColor"
                  viewBox="0 0 24 24"
                >
                  <path d="M17.472 14.382c-.297-.149-1.758-.867-2.03-.967-.273-.099-.471-.148-.67.15-.197.297-.767.966-.94 1.164-.173.199-.347.223-.644.075-.297-.15-1.255-.463-2.39-1.475-.883-.788-1.48-1.761-1.653-2.059-.173-.297-.018-.458.13-.606.134-.133.298-.347.446-.52.149-.174.198-.298.298-.497.099-.198.05-.371-.025-.52-.075-.149-.669-1.612-.916-2.207-.242-.579-.487-.5-.669-.51-.173-.008-.371-.01-.57-.01-.198 0-.52.074-.792.372-.272.297-1.04 1.016-1.04 2.479 0 1.462 1.065 2.875 1.213 3.074.149.198 2.096 3.2 5.077 4.487.709.306 1.262.489 1.694.625.712.227 1.36.195 1.871.118.571-.085 1.758-.719 2.006-1.413.248-.694.248-1.289.173-1.413-.074-.124-.272-.198-.57-.347m-5.421 7.403h-.004a9.87 9.87 0 01-5.031-1.378l-.361-.214-3.741.982.998-3.648-.235-.374a9.86 9.86 0 01-1.51-5.26c.001-5.45 4.436-9.884 9.888-9.884 2.64 0 5.122 1.03 6.988 2.898a9.825 9.825 0 012.893 6.994c-.003 5.45-4.437 9.884-9.885 9.884m8.413-18.297A11.815 11.815 0 0012.05 0C5.495 0 .16 5.335.157 11.892c0 2.096.547 4.142 1.588 5.945L.057 24l6.305-1.654a11.882 11.882 0 005.683 1.448h.005c6.554 0 11.89-5.335 11.893-11.893a11.821 11.821 0 00-3.48-8.413Z" />
                </svg>
                WhatsApp
              </a>
            </li>
            <li>
              <a onClick={() => handleShare('telegram')}>
                <svg
                  className="w-5 h-5 text-info"
                  fill="currentColor"
                  viewBox="0 0 24 24"
                >
                  <path d="M11.944 0A12 12 0 1 0 24 12a12 12 0 0 0-12.056-12zM16.906 7.224c.1-.002.321.023.465.14a.506.506 0 0 1 .171.325c.016.093.036.306.02.472-.18 1.898-.962 6.502-1.36 8.627-.168.9-.499 1.201-.82 1.23-.696.065-1.225-.46-1.9-.902-1.056-.693-1.653-1.124-2.678-1.8-1.185-.78-.417-1.21.258-1.91.177-.184 3.247-2.977 3.307-3.23.007-.032.014-.15-.056-.212s-.174-.041-.249-.024c-.106.024-1.793 1.14-5.061 3.345-.48.33-.913.49-1.302.48-.428-.008-1.252-.241-1.865-.44-.752-.245-1.349-.374-1.297-.789.027-.216.325-.437.893-.663 3.498-1.524 5.83-2.529 6.998-3.014 3.332-1.386 4.025-1.627 4.476-1.635z" />
                </svg>
                Telegram
              </a>
            </li>
            <li>
              <a onClick={() => handleShare('viber')}>
                <svg
                  className="w-5 h-5 text-secondary"
                  fill="currentColor"
                  viewBox="0 0 24 24"
                >
                  <path d="M11.398.002C9.473.028 5.331.344 3.014 2.467 1.294 4.177.693 6.698.623 9.82c-.06 3.11-.13 8.95 5.5 10.541v2.42s-.038.97.602 1.17c.79.25 1.24-.499 1.99-1.299l1.4-1.58c3.85.32 6.8-.419 7.14-.529.78-.25 5.181-.811 5.901-6.652.74-6.031-.36-9.831-2.34-11.551l-.01-.002c-.6-.55-3-2.3-8.37-2.32 0 0-.396-.025-1.038-.016zm.067 1.697c.545-.003.88.02.88.02 4.54.01 6.711 1.38 7.221 1.84 1.67 1.429 2.528 4.856 1.9 9.892-.6 4.88-4.17 5.19-4.83 5.4-.28.09-2.88.73-6.152.52 0 0-2.439 2.941-3.199 3.701-.12.13-.26.17-.35.15-.13-.03-.17-.19-.17-.41l.02-4.019c-4.771-1.32-4.491-6.302-4.441-8.902.06-2.6.55-4.732 2-6.172 1.957-1.77 5.475-2.01 7.11-2.02h.01zm.36 2.6a.299.299 0 0 0-.3.299.3.3 0 0 0 .3.3 5.631 5.631 0 0 1 4.03 1.59c1.09 1.06 1.621 2.48 1.641 4.34a.3.3 0 0 0 .3.3v-.009a.3.3 0 0 0 .3-.3 6.451 6.451 0 0 0-1.81-4.76 6.191 6.191 0 0 0-4.46-1.76zm-3.954.69a.955.955 0 0 0-.615.12l-.012.002c-.41.24-.788.54-1.148.94-.27.32-.421.639-.461.949-.008.062-.009.129-.009.199a1.77 1.77 0 0 0 .451 1.24c.332.407.741.863 1.218 1.362l.002.002c.64.657 1.339 1.292 2.062 1.869.39.31.79.599 1.2.87.708.465 1.312.75 1.811.859l.06.01c.155.029.304.043.444.043.22 0 .422-.038.604-.114.301-.127.572-.31.779-.542l.01-.01c.23-.26.41-.55.53-.87.108-.286.114-.57.015-.829l-.002-.009a1.006 1.006 0 0 0-.366-.399l-.01-.006-1.52-.95a1.06 1.06 0 0 0-.54-.141 1 1 0 0 0-.69.32l-.627.751a.734.734 0 0 1-.542.236.72.72 0 0 1-.264-.046c-.893-.278-1.669-.732-2.332-1.363-.602-.57-1.13-1.292-1.583-2.164a.76.76 0 0 1-.064-.315c0-.148.044-.279.136-.395l.76-.676.02-.018a.914.914 0 0 0 .287-.61.894.894 0 0 0-.168-.512l-.001-.001-.944-1.51a.997.997 0 0 0-.465-.392h-.002a.875.875 0 0 0-.363-.078zm4.475.748a.3.3 0 0 0 .002.6 3.78 3.78 0 0 1 2.65 1.065 3.5 3.5 0 0 1 .9 2.63.3.3 0 0 0 .3.299v.009a.3.3 0 0 0 .3-.3c.03-1.19-.4-2.182-1.1-2.992-.72-.83-1.73-1.311-2.948-1.31h-.103v-.001z" />
                </svg>
                Viber
              </a>
            </li>
            <li>
              <a onClick={() => handleShare('facebook')}>
                <svg
                  className="w-5 h-5 text-primary"
                  fill="currentColor"
                  viewBox="0 0 24 24"
                >
                  <path d="M24 12.073c0-6.627-5.373-12-12-12s-12 5.373-12 12c0 5.99 4.388 10.954 10.125 11.854v-8.385H7.078v-3.47h3.047V9.43c0-3.007 1.792-4.669 4.533-4.669 1.312 0 2.686.235 2.686.235v2.953H15.83c-1.491 0-1.956.925-1.956 1.874v2.25h3.328l-.532 3.47h-2.796v8.385C19.612 23.027 24 18.062 24 12.073z" />
                </svg>
                Facebook
              </a>
            </li>
            <li>
              <a onClick={() => handleShare('twitter')}>
                <svg
                  className="w-5 h-5"
                  fill="currentColor"
                  viewBox="0 0 24 24"
                >
                  <path d="M23.953 4.57a10 10 0 01-2.825.775 4.958 4.958 0 002.163-2.723c-.951.555-2.005.959-3.127 1.184a4.92 4.92 0 00-8.384 4.482C7.69 8.095 4.067 6.13 1.64 3.162a4.822 4.822 0 00-.666 2.475c0 1.71.87 3.213 2.188 4.096a4.904 4.904 0 01-2.228-.616v.06a4.923 4.923 0 003.946 4.827 4.996 4.996 0 01-2.212.085 4.936 4.936 0 004.604 3.417 9.867 9.867 0 01-6.102 2.105c-.39 0-.779-.023-1.17-.067a13.995 13.995 0 007.557 2.209c9.053 0 13.998-7.496 13.998-13.985 0-.21 0-.42-.015-.63A9.935 9.935 0 0024 4.59z" />
                </svg>
                Twitter / X
              </a>
            </li>
            <div className="divider my-1"></div>
            <li>
              <a onClick={() => handleShare('copy')}>
                <svg
                  className="w-5 h-5"
                  fill="none"
                  viewBox="0 0 24 24"
                  stroke="currentColor"
                >
                  <path
                    strokeLinecap="round"
                    strokeLinejoin="round"
                    strokeWidth={2}
                    d="M8 16H6a2 2 0 01-2-2V6a2 2 0 012-2h8a2 2 0 012 2v2m-6 12h8a2 2 0 002-2v-8a2 2 0 00-2-2h-8a2 2 0 00-2 2v8a2 2 0 002 2z"
                  />
                </svg>
                {t('listingActions.copyLink')}
              </a>
            </li>
          </ul>
        )}
      </div>

      {/* Report Button */}
      <button
        className="btn btn-circle btn-ghost text-error"
        aria-label={t('listingActions.report')}
      >
        <svg
          className="w-6 h-6"
          fill="none"
          viewBox="0 0 24 24"
          stroke="currentColor"
        >
          <path
            strokeLinecap="round"
            strokeLinejoin="round"
            strokeWidth={2}
            d="M12 9v2m0 4h.01m-6.938 4h13.856c1.54 0 2.502-1.667 1.732-3L13.732 4c-.77-1.333-2.694-1.333-3.464 0L3.34 16c-.77 1.333.192 3 1.732 3z"
          />
        </svg>
      </button>
    </div>
  );
}
