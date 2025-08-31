'use client';

import { useRouter } from 'next/navigation';
import { useAuthContext } from '@/contexts/AuthContext';
import { useTranslations } from 'next-intl';
import { MessageCircle } from 'lucide-react';

interface ContactSellerButtonProps {
  sellerId: number;
  listingId?: number;
  storefrontProductId?: number;
  className?: string;
  size?: 'sm' | 'md' | 'lg';
}

export default function ContactSellerButton({
  sellerId,
  listingId,
  storefrontProductId,
  className = '',
  size = 'md',
}: ContactSellerButtonProps) {
  const router = useRouter();
  const { user } = useAuthContext();
  const t = useTranslations('common');

  const handleClick = () => {
    if (!user) {
      // Redirect to login with return URL
      const returnUrl = listingId
        ? `/chat?listing_id=${listingId}&seller_id=${sellerId}`
        : `/chat?storefront_product_id=${storefrontProductId}&seller_id=${sellerId}`;
      router.push(`/auth/login?redirect=${encodeURIComponent(returnUrl)}`);
      return;
    }

    // Prevent messaging yourself
    if (user.id === sellerId) {
      return;
    }

    // Navigate to chat with parameters to create/open chat
    const params = new URLSearchParams();
    params.append('seller_id', sellerId.toString());

    if (listingId) {
      params.append('listing_id', listingId.toString());
    } else if (storefrontProductId) {
      params.append('storefront_product_id', storefrontProductId.toString());
    }

    router.push(`/chat?${params.toString()}`);
  };

  // Don't show button if user is the seller
  if (user?.id === sellerId) {
    return null;
  }

  const sizeClasses = {
    sm: 'btn-sm',
    md: '',
    lg: 'btn-lg',
  };

  return (
    <button
      onClick={handleClick}
      className={`btn btn-primary ${sizeClasses[size]} ${className}`}
      aria-label={t('contactSeller')}
    >
      <MessageCircle className="w-4 h-4 mr-2" />
      {t('contactSeller')}
    </button>
  );
}
