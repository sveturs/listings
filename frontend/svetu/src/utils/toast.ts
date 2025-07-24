export type ToastType = 'success' | 'error' | 'warning' | 'info';

export interface ToastOptions {
  type?: ToastType;
  duration?: number;
  position?:
    | 'top'
    | 'bottom'
    | 'top-center'
    | 'top-start'
    | 'top-end'
    | 'middle'
    | 'bottom-start'
    | 'bottom-end';
}

class ToastManager {
  private container: HTMLDivElement | null = null;

  private getPositionClasses(
    position: ToastOptions['position'] = 'top'
  ): string {
    const positions = {
      top: 'toast-top',
      bottom: 'toast-bottom',
      'top-center': 'toast-top toast-center',
      'top-start': 'toast-top toast-start',
      'top-end': 'toast-top toast-end',
      middle: 'toast-middle',
      'bottom-start': 'toast-bottom toast-start',
      'bottom-end': 'toast-bottom toast-end',
    };
    return positions[position] || positions['top'];
  }

  private getTypeClasses(type: ToastType = 'info'): string {
    const types = {
      success: 'alert-success',
      error: 'alert-error',
      warning: 'alert-warning',
      info: 'alert-info',
    };
    return types[type] || types.info;
  }

  private getIcon(type: ToastType): string {
    const icons = {
      success: `<svg xmlns="http://www.w3.org/2000/svg" class="h-6 w-6 shrink-0 stroke-current" fill="none" viewBox="0 0 24 24">
        <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M9 12l2 2 4-4m6 2a9 9 0 11-18 0 9 9 0 0118 0z" />
      </svg>`,
      error: `<svg xmlns="http://www.w3.org/2000/svg" class="h-6 w-6 shrink-0 stroke-current" fill="none" viewBox="0 0 24 24">
        <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M10 14l2-2m0 0l2-2m-2 2l-2-2m2 2l2 2m7-2a9 9 0 11-18 0 9 9 0 0118 0z" />
      </svg>`,
      warning: `<svg xmlns="http://www.w3.org/2000/svg" class="h-6 w-6 shrink-0 stroke-current" fill="none" viewBox="0 0 24 24">
        <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 9v2m0 4h.01m-6.938 4h13.856c1.54 0 2.502-1.667 1.732-3L13.732 4c-.77-1.333-2.694-1.333-3.464 0L3.34 16c-.77 1.333.192 3 1.732 3z" />
      </svg>`,
      info: `<svg xmlns="http://www.w3.org/2000/svg" class="h-6 w-6 shrink-0 stroke-current" fill="none" viewBox="0 0 24 24">
        <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M13 16h-1v-4h-1m1-4h.01M21 12a9 9 0 11-18 0 9 9 0 0118 0z" />
      </svg>`,
    };
    return icons[type] || icons.info;
  }

  private close(toast: HTMLElement) {
    if (!toast || toast.style.opacity === '0') return;

    // Отменяем таймер если он есть
    const timeoutId = (toast as any).timeoutId;
    if (timeoutId) {
      clearTimeout(timeoutId);
    }

    toast.style.opacity = '0';
    toast.style.transform = 'scale(0.9)';

    setTimeout(() => {
      toast.remove();

      // Удаляем контейнер если он пустой
      if (this.container && this.container.children.length === 0) {
        this.container.remove();
        this.container = null;
      }
    }, 300);
  }

  show(message: string, options: ToastOptions = {}) {
    const { type = 'info', duration = 3000, position = 'top' } = options;

    // Создаем контейнер если его нет
    if (!this.container) {
      this.container = document.createElement('div');
      this.container.className = 'toast toast-top toast-center z-50';
      document.body.appendChild(this.container);
    }

    // Устанавливаем позицию контейнера
    this.container.className = `toast ${this.getPositionClasses(position)} z-[9999]`;

    // Создаем элемент уведомления
    const toast = document.createElement('div');
    toast.className = `alert ${this.getTypeClasses(type)} shadow-lg mb-2 cursor-pointer`;
    toast.innerHTML = `
      <div class="flex items-center gap-2">
        ${this.getIcon(type)}
        <span>${message}</span>
      </div>
    `;

    // Закрытие по клику
    toast.addEventListener('click', () => {
      this.close(toast);
    });

    // Добавляем в контейнер
    this.container.appendChild(toast);

    // Анимация появления
    toast.style.opacity = '0';
    toast.style.transform = 'scale(0.9)';
    toast.style.transition = 'all 0.3s ease-out';

    setTimeout(() => {
      toast.style.opacity = '1';
      toast.style.transform = 'scale(1)';
    }, 10);

    // Удаляем через заданное время
    const timeoutId = setTimeout(() => {
      this.close(toast);
    }, duration);

    // Сохраняем ID таймера для возможности отмены
    (toast as any).timeoutId = timeoutId;
  }

  success(message: string, options?: Omit<ToastOptions, 'type'>) {
    this.show(message, { ...options, type: 'success' });
  }

  error(message: string, options?: Omit<ToastOptions, 'type'>) {
    this.show(message, { ...options, type: 'error' });
  }

  warning(message: string, options?: Omit<ToastOptions, 'type'>) {
    this.show(message, { ...options, type: 'warning' });
  }

  info(message: string, options?: Omit<ToastOptions, 'type'>) {
    this.show(message, { ...options, type: 'info' });
  }
}

export const toast = new ToastManager();
