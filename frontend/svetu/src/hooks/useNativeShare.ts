import { useState, useCallback } from 'react';
import { toast } from 'react-hot-toast';

interface ShareData {
  title?: string;
  text?: string;
  url?: string;
  files?: File[];
}

interface ShareResult {
  success: boolean;
  method?: 'native' | 'clipboard' | 'fallback';
  error?: string;
}

interface UseNativeShareReturn {
  canShare: boolean;
  canShareFiles: boolean;
  share: (data: ShareData) => Promise<ShareResult>;
  shareProduct: (product: any) => Promise<ShareResult>;
  shareImage: (imageUrl: string, title?: string) => Promise<ShareResult>;
  copyToClipboard: (text: string) => Promise<boolean>;
  generateShareableLink: (productId: string) => string;
  shareViaWhatsApp: (text: string, url?: string) => void;
  shareViaTelegram: (text: string, url?: string) => void;
  shareViaEmail: (subject: string, body: string) => void;
  isSharing: boolean;
}

export function useNativeShare(): UseNativeShareReturn {
  const [isSharing, setIsSharing] = useState(false);

  // Check if Web Share API is available
  const canShare = typeof navigator !== 'undefined' && 'share' in navigator;
  const canShareFiles =
    typeof navigator !== 'undefined' && 'canShare' in navigator;

  /**
   * Main share function using Web Share API
   */
  const share = useCallback(
    async (data: ShareData): Promise<ShareResult> => {
      setIsSharing(true);

      try {
        // Check if can share files
        if (data.files && canShareFiles) {
          const canShareResult = navigator.canShare(data);
          if (!canShareResult) {
            // Remove files if can't share them
            delete data.files;
          }
        }

        if (canShare) {
          // Use native share
          await navigator.share(data);
          return { success: true, method: 'native' };
        } else {
          // Fallback to clipboard
          const shareText =
            `${data.title || ''}\n${data.text || ''}\n${data.url || ''}`.trim();
          const copied = await copyToClipboard(shareText);

          if (copied) {
            toast.success('Скопировано в буфер обмена');
            return { success: true, method: 'clipboard' };
          } else {
            throw new Error('Failed to copy to clipboard');
          }
        }
      } catch (error: any) {
        // User cancelled share
        if (error.name === 'AbortError') {
          return { success: false, error: 'Share cancelled' };
        }

        // Try fallback methods
        return await shareFallback(data);
      } finally {
        setIsSharing(false);
      }
    },
    // eslint-disable-next-line react-hooks/exhaustive-deps
    [canShare, canShareFiles]
  );

  /**
   * Share a product
   */
  const shareProduct = useCallback(
    async (product: any): Promise<ShareResult> => {
      const url = generateShareableLink(product.id);
      const title = product.title || product.name;
      const description = product.description || '';
      const price = product.price
        ? `Цена: ${product.price} ${product.currency || 'RSD'}`
        : '';

      const text = `${description}\n${price}`.trim();

      // If product has images, try to share them
      if (product.images && product.images.length > 0 && canShareFiles) {
        try {
          // Fetch first image as file
          const response = await fetch(product.images[0].url);
          const blob = await response.blob();
          const file = new File([blob], 'product.jpg', { type: 'image/jpeg' });

          return await share({
            title,
            text,
            url,
            files: [file],
          });
        } catch {
          // If image fetch fails, share without files
        }
      }

      return await share({ title, text, url });
    },
    // eslint-disable-next-line react-hooks/exhaustive-deps
    [share, canShareFiles]
  );

  /**
   * Share an image with optional title
   */
  const shareImage = useCallback(
    async (imageUrl: string, title?: string): Promise<ShareResult> => {
      try {
        if (canShareFiles) {
          // Fetch image and convert to file
          const response = await fetch(imageUrl);
          const blob = await response.blob();
          const fileName = imageUrl.split('/').pop() || 'image.jpg';
          const file = new File([blob], fileName, { type: blob.type });

          return await share({
            title: title || 'Изображение',
            files: [file],
          });
        } else {
          // Share URL instead
          return await share({
            title: title || 'Изображение',
            url: imageUrl,
          });
        }
      } catch {
        return {
          success: false,
          error: 'Failed to share image',
        };
      }
    },
    [share, canShareFiles]
  );

  /**
   * Copy text to clipboard
   */
  const copyToClipboard = useCallback(
    async (text: string): Promise<boolean> => {
      try {
        if ('clipboard' in navigator) {
          await navigator.clipboard.writeText(text);
          return true;
        } else {
          // Fallback for older browsers
          const textArea = document.createElement('textarea');
          textArea.value = text;
          textArea.style.position = 'fixed';
          textArea.style.left = '-999999px';
          document.body.appendChild(textArea);
          textArea.focus();
          textArea.select();

          try {
            document.execCommand('copy');
            return true;
          } finally {
            document.body.removeChild(textArea);
          }
        }
      } catch {
        return false;
      }
    },
    []
  );

  /**
   * Generate shareable link for a product
   */
  const generateShareableLink = useCallback((productId: string): string => {
    const baseUrl =
      typeof window !== 'undefined'
        ? window.location.origin
        : 'https://svetu.rs';

    // Add tracking parameters
    const params = new URLSearchParams({
      utm_source: 'share',
      utm_medium: 'native',
      utm_campaign: 'product_share',
    });

    return `${baseUrl}/product/${productId}?${params.toString()}`;
  }, []);

  /**
   * Fallback share methods
   */
  const shareFallback = async (data: ShareData): Promise<ShareResult> => {
    try {
      // Create share menu
      const shareOptions = [];

      // WhatsApp
      if (isMobileDevice()) {
        shareOptions.push({
          name: 'WhatsApp',
          icon: '💬',
          action: () =>
            shareViaWhatsApp(`${data.title}\n${data.text}`.trim(), data.url),
        });
      }

      // Telegram
      shareOptions.push({
        name: 'Telegram',
        icon: '✈️',
        action: () =>
          shareViaTelegram(`${data.title}\n${data.text}`.trim(), data.url),
      });

      // Email
      shareOptions.push({
        name: 'Email',
        icon: '✉️',
        action: () =>
          shareViaEmail(data.title || '', `${data.text}\n\n${data.url}`.trim()),
      });

      // Copy link
      shareOptions.push({
        name: 'Copy Link',
        icon: '🔗',
        action: async () => {
          if (data.url) {
            await copyToClipboard(data.url);
            toast.success('Ссылка скопирована');
          }
        },
      });

      // Show custom share dialog
      showShareDialog(shareOptions);

      return { success: true, method: 'fallback' };
    } catch {
      return {
        success: false,
        error: 'Share failed',
        method: 'fallback',
      };
    }
  };

  /**
   * Share via WhatsApp
   */
  const shareViaWhatsApp = useCallback((text: string, url?: string) => {
    const message = url ? `${text}\n\n${url}` : text;
    const whatsappUrl = `https://wa.me/?text=${encodeURIComponent(message)}`;

    if (isMobileDevice()) {
      // Try to open WhatsApp app
      window.location.href = `whatsapp://send?text=${encodeURIComponent(message)}`;

      // Fallback to web WhatsApp after delay
      setTimeout(() => {
        window.open(whatsappUrl, '_blank');
      }, 500);
    } else {
      // Open WhatsApp Web
      window.open(whatsappUrl, '_blank');
    }
  }, []);

  /**
   * Share via Telegram
   */
  const shareViaTelegram = useCallback((text: string, url?: string) => {
    const telegramUrl = url
      ? `https://t.me/share/url?url=${encodeURIComponent(url)}&text=${encodeURIComponent(text)}`
      : `https://t.me/share/url?text=${encodeURIComponent(text)}`;

    if (isMobileDevice()) {
      // Try to open Telegram app
      window.location.href = `tg://msg?text=${encodeURIComponent(`${text}\n\n${url || ''}`)}`;

      // Fallback to web version
      setTimeout(() => {
        window.open(telegramUrl, '_blank');
      }, 500);
    } else {
      window.open(telegramUrl, '_blank');
    }
  }, []);

  /**
   * Share via Email
   */
  const shareViaEmail = useCallback((subject: string, body: string) => {
    const mailtoUrl = `mailto:?subject=${encodeURIComponent(subject)}&body=${encodeURIComponent(body)}`;
    window.location.href = mailtoUrl;
  }, []);

  /**
   * Helper function to detect mobile device
   */
  const isMobileDevice = (): boolean => {
    return /Android|webOS|iPhone|iPad|iPod|BlackBerry|IEMobile|Opera Mini/i.test(
      navigator.userAgent
    );
  };

  /**
   * Show custom share dialog (for fallback)
   */
  const showShareDialog = (options: any[]) => {
    // Create and show a custom modal with share options
    const modal = document.createElement('div');
    modal.className = 'fixed inset-0 z-50 flex items-end justify-center';
    modal.innerHTML = `
      <div class="absolute inset-0 bg-black/50" onclick="this.parentElement.remove()"></div>
      <div class="relative w-full max-w-md rounded-t-xl bg-white p-4">
        <h3 class="mb-4 text-lg font-semibold">Поделиться</h3>
        <div class="grid grid-cols-4 gap-4">
          ${options
            .map(
              (opt) => `
            <button class="share-option flex flex-col items-center p-2 hover:bg-gray-100 rounded" data-action="${opt.name}">
              <span class="text-2xl">${opt.icon}</span>
              <span class="mt-1 text-xs">${opt.name}</span>
            </button>
          `
            )
            .join('')}
        </div>
        <button class="btn btn-ghost mt-4 w-full" onclick="this.parentElement.parentElement.remove()">
          Отмена
        </button>
      </div>
    `;

    document.body.appendChild(modal);

    // Add click handlers
    modal.querySelectorAll('.share-option').forEach((btn, index) => {
      btn.addEventListener('click', () => {
        options[index].action();
        modal.remove();
      });
    });
  };

  return {
    canShare,
    canShareFiles,
    share,
    shareProduct,
    shareImage,
    copyToClipboard,
    generateShareableLink,
    shareViaWhatsApp,
    shareViaTelegram,
    shareViaEmail,
    isSharing,
  };
}
