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
    toast.className = `alert ${this.getTypeClasses(type)} shadow-lg mb-2`;
    toast.innerHTML = `
      <div>
        <span>${message}</span>
      </div>
    `;

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
    setTimeout(() => {
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
    }, duration);
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
