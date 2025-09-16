'use client';

import { useRouter } from 'next/navigation';
import { useLocale, useTranslations } from 'next-intl';
import { useState } from 'react';
import { useAuth } from '@/contexts/AuthContext';
import type { Storefront } from '@/types/storefront';

interface StorefrontActionsProps {
  storefront: Storefront;
  isOwner: boolean;
}

export default function StorefrontActions({ storefront, isOwner }: StorefrontActionsProps) {
  const t = useTranslations('storefronts');
  const tCommon = useTranslations('common');
  const locale = useLocale();
  const router = useRouter();
  const { user } = useAuth();
  const [isFollowing, setIsFollowing] = useState(false);
  const [shareOpen, setShareOpen] = useState(false);

  const handleContact = () => {
    if (!user) {
      router.push(`/${locale}/`);
      return;
    }
    
    router.push(`/${locale}/chat?storefront_id=${storefront.id}`);
  };

  const handleFollow = async () => {
    if (!user) {
      router.push(`/${locale}/`);
      return;
    }
    
    // TODO: Implement follow API
    setIsFollowing(!isFollowing);
  };

  const handleShare = () => {
    setShareOpen(true);
  };

  const handleCopyLink = () => {
    const url = window.location.href;
    navigator.clipboard.writeText(url);
    // TODO: Show toast notification
    setShareOpen(false);
  };

  return (
    <>
      <div className="card bg-base-200 shadow-xl">
        <div className="card-body">
          <h3 className="card-title">{tCommon('actionsTitle')}</h3>
          
          <div className="space-y-3">
            {isOwner ? (
              <>
                <button 
                  className="btn btn-primary btn-block"
                  onClick={() => router.push(`/${locale}/storefronts/${storefront.slug}/dashboard`)}
                >
                  <svg className="w-5 h-5" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                    <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M10.325 4.317c.426-1.756 2.924-1.756 3.35 0a1.724 1.724 0 002.573 1.066c1.543-.94 3.31.826 2.37 2.37a1.724 1.724 0 001.065 2.572c1.756.426 1.756 2.924 0 3.35a1.724 1.724 0 00-1.066 2.573c.94 1.543-.826 3.31-2.37 2.37a1.724 1.724 0 00-2.572 1.065c-.426 1.756-2.924 1.756-3.35 0a1.724 1.724 0 00-2.573-1.066c-1.543.94-3.31-.826-2.37-2.37a1.724 1.724 0 00-1.065-2.572c-1.756-.426-1.756-2.924 0-3.35a1.724 1.724 0 001.066-2.573c-.94-1.543.826-3.31 2.37-2.37.996.608 2.296.07 2.572-1.065z" />
                    <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M15 12a3 3 0 11-6 0 3 3 0 016 0z" />
                  </svg>
                  {t('manageDashboard')}
                </button>
                
                <button 
                  className="btn btn-outline btn-block"
                  onClick={() => router.push(`/${locale}/storefronts/${storefront.slug}/analytics`)}
                >
                  <svg className="w-5 h-5" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                    <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M9 19v-6a2 2 0 00-2-2H5a2 2 0 00-2 2v6a2 2 0 002 2h2a2 2 0 002-2zm0 0V9a2 2 0 012-2h2a2 2 0 012 2v10m-6 0a2 2 0 002 2h2a2 2 0 002-2m0 0V5a2 2 0 012-2h2a2 2 0 012 2v14a2 2 0 01-2 2h-2a2 2 0 01-2-2z" />
                  </svg>
                  {t('viewAnalytics')}
                </button>
              </>
            ) : (
              <>
                <button 
                  className="btn btn-primary btn-block"
                  onClick={handleContact}
                >
                  <svg className="w-5 h-5" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                    <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M8 12h.01M12 12h.01M16 12h.01M21 12c0 4.418-4.03 8-9 8a9.863 9.863 0 01-4.255-.949L3 20l1.395-3.72C3.512 15.042 3 13.574 3 12c0-4.418 4.03-8 9-8s9 3.582 9 8z" />
                  </svg>
                  {t('contactStore')}
                </button>
                
                <button 
                  className={`btn btn-block ${isFollowing ? 'btn-accent' : 'btn-outline'}`}
                  onClick={handleFollow}
                >
                  <svg className="w-5 h-5" fill={isFollowing ? 'currentColor' : 'none'} viewBox="0 0 24 24" stroke="currentColor">
                    <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M4.318 6.318a4.5 4.5 0 000 6.364L12 20.364l7.682-7.682a4.5 4.5 0 00-6.364-6.364L12 7.636l-1.318-1.318a4.5 4.5 0 00-6.364 0z" />
                  </svg>
                  {isFollowing ? t('following') : t('follow')}
                </button>
              </>
            )}
            
            <button 
              className="btn btn-ghost btn-block"
              onClick={handleShare}
            >
              <svg className="w-5 h-5" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M8.684 13.342C8.886 12.938 9 12.482 9 12c0-.482-.114-.938-.316-1.342m0 2.684a3 3 0 110-2.684m9.632 4.316C18.886 16.938 19 17.482 19 18c0 1.657-1.343 3-3 3s-3-1.343-3-3 1.343-3 3-3c.482 0 .938.114 1.342.316m0 0a3 3 0 00-4.684-4.684m4.684 4.684a3 3 0 014.684-4.684m0 0C20.886 8.938 21 8.482 21 8c0-1.657-1.343-3-3-3s-3 1.343-3 3 1.343 3 3 3c.482 0 .938-.114 1.342-.316" />
              </svg>
              {tCommon('share')}
            </button>
          </div>

          {/* Payment Methods */}
          {storefront.payment_methods && Array.isArray(storefront.payment_methods) && storefront.payment_methods.length > 0 && (
            <div className="mt-4 pt-4 border-t border-base-300">
              <p className="text-sm font-semibold mb-2">{t('acceptedPayments')}</p>
              <div className="flex flex-wrap gap-2">
                {storefront.payment_methods.map((method, index) => {
                  if (!method || typeof method.method_type !== 'string') return null;
                  return (
                    <div key={index} className="badge badge-outline">
                      {method.method_type === 'cash' && t('payment_methods.cash')}
                      {method.method_type === 'card' && t('payment_methods.card')}
                      {method.method_type === 'bank_transfer' && t('payment_methods.bank_transfer')}
                      {method.method_type === 'cod' && t('payment_methods.cod')}
                    </div>
                  );
                })}
              </div>
            </div>
          )}

          {/* Delivery Options */}
          {storefront.delivery_options && Array.isArray(storefront.delivery_options) && storefront.delivery_options.length > 0 && (
            <div className="mt-4 pt-4 border-t border-base-300">
              <p className="text-sm font-semibold mb-2">{t('deliveryOptions')}</p>
              <div className="space-y-1">
                {storefront.delivery_options.map((option, index) => {
                  if (!option || typeof option.name !== 'string') return null;
                  return (
                    <div key={index} className="text-sm">
                      <span className="font-medium">{option.name}</span>
                      {typeof option.base_price === 'number' && option.base_price > 0 && (
                        <span className="text-base-content/60"> - {option.base_price} RSD</span>
                      )}
                    </div>
                  );
                })}
              </div>
            </div>
          )}
        </div>
      </div>

      {/* Share Modal */}
      {shareOpen && (
        <div className="modal modal-open">
          <div className="modal-box">
            <h3 className="font-bold text-lg">{t('shareStore')}</h3>
            <div className="py-4">
              <div className="form-control">
                <label className="label">
                  <span className="label-text">{t('storeLink')}</span>
                </label>
                <div className="input-group">
                  <input 
                    type="text" 
                    value={window.location.href} 
                    readOnly 
                    className="input input-bordered flex-1"
                  />
                  <button className="btn btn-square" onClick={handleCopyLink}>
                    <svg className="w-5 h-5" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                      <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M8 16H6a2 2 0 01-2-2V6a2 2 0 012-2h8a2 2 0 012 2v2m-6 12h8a2 2 0 002-2v-8a2 2 0 00-2-2h-8a2 2 0 00-2 2v8a2 2 0 002 2z" />
                    </svg>
                  </button>
                </div>
              </div>
            </div>
            <div className="modal-action">
              <button className="btn" onClick={() => setShareOpen(false)}>
                {tCommon('close')}
              </button>
            </div>
          </div>
        </div>
      )}
    </>
  );
}