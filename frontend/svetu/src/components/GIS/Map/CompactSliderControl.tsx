import React, { useEffect, useRef, useState } from 'react';
import { IControl, Map as MapboxMap } from 'mapbox-gl';

interface CompactSliderControlProps {
  mode: 'radius' | 'walking';
  onModeChange: (mode: 'radius' | 'walking') => void;
  walkingTime: number; // в минутах (5-60)
  onWalkingTimeChange: (time: number) => void;
  searchRadius: number; // в метрах (500-50000)
  onRadiusChange: (radius: number) => void;
  position?: 'top-left' | 'top-right' | 'bottom-left' | 'bottom-right';
  isFullscreen?: boolean;
  isMobile?: boolean;
  translations?: {
    walkingAccessibility: string;
    searchRadius: string;
    minutes: string;
    km: string;
    m: string;
    changeModeHint: string;
  };
}

class CompactSliderControlClass implements IControl {
  private map: MapboxMap | undefined;
  private container: HTMLDivElement;
  private props: CompactSliderControlProps;
  private onPropsChange: (props: CompactSliderControlProps) => void;
  private isExpanded: boolean = false;
  private longPressTimer: NodeJS.Timeout | null = null;
  private lastTapTime: number = 0;
  private tempWalkingTime: number | null = null;
  private tempSearchRadius: number | null = null;

  constructor(
    props: CompactSliderControlProps,
    onPropsChange: (props: CompactSliderControlProps) => void
  ) {
    this.props = props;
    this.onPropsChange = onPropsChange;
    this.container = document.createElement('div');
    // Инициализируем временные значения
    this.tempWalkingTime = props.walkingTime;
    this.tempSearchRadius = props.searchRadius;
  }

  onAdd(map: MapboxMap): HTMLElement {
    console.log('[CompactSliderControlClass] onAdd called');
    this.map = map;
    // Добавляем классы как у нативных контролов
    this.container.className = 'mapboxgl-ctrl mapboxgl-ctrl-group';

    // Основные стили контейнера
    this.container.style.transition = 'all 0.3s ease-in-out';
    this.container.style.overflow = 'visible';
    this.container.style.userSelect = 'none';
    this.container.style.touchAction = 'manipulation';
    this.container.style.pointerEvents = 'auto'; // Явно разрешаем события
    this.container.style.zIndex = '999'; // Высокий z-index

    this.updateContainerStyle();
    this.render();

    console.log('[CompactSliderControlClass] Container created and styled');
    return this.container;
  }

  onRemove(): void {
    if (this.container.parentNode) {
      this.container.parentNode.removeChild(this.container);
    }
    if (this.longPressTimer) {
      clearTimeout(this.longPressTimer);
    }
    this.removeOutsideClickHandler();
    this.map = undefined;
  }

  updateProps(newProps: CompactSliderControlProps): void {
    this.props = newProps;
    // Обновляем временные значения при изменении props
    if (this.props.mode === 'walking') {
      this.tempWalkingTime = newProps.walkingTime;
    } else {
      this.tempSearchRadius = newProps.searchRadius;
    }
    this.updateContainerStyle();
    this.render();
  }

  private updateContainerStyle(): void {
    const { isMobile } = this.props;

    // Базовые стили
    this.container.style.position = 'relative';
    this.container.style.transition = 'all 0.3s ease-in-out';

    // Добавляем отступ между контролами (как у других нативных контролов)
    this.container.style.marginTop = '10px';
    this.container.style.marginBottom = '0';

    if (this.isExpanded) {
      // Развернутое состояние - переопределяем стили для расширенной панели
      this.container.style.background = 'white';
      this.container.style.borderRadius = '8px';
      this.container.style.boxShadow = '0 2px 10px rgba(0,0,0,0.15)';
      this.container.style.padding = '12px';
      this.container.style.width = isMobile ? '260px' : '300px';
      this.container.style.height = 'auto';
      this.container.style.cursor = 'default';
      this.container.style.zIndex = '1000'; // Высокий z-index для развернутого состояния
    } else {
      // Свернутое состояние - используем стандартные стили mapbox контролов
      this.container.style.background = '';
      this.container.style.borderRadius = '';
      this.container.style.boxShadow = '';
      this.container.style.padding = '0';
      this.container.style.width = '29px';
      this.container.style.height = '29px';
      this.container.style.cursor = 'pointer';
    }
  }

  private render(): void {
    if (!this.isExpanded) {
      // Свернутое состояние - только иконка
      this.renderCompact();
    } else {
      // Развернутое состояние - иконка + слайдер
      this.renderExpanded();
    }

    this.addEventListeners();
  }

  private renderCompact(): void {
    const { mode } = this.props;
    const icon = mode === 'walking' ? '🚶' : '📏';
    const bgColor = mode === 'walking' ? '#10B981' : '#3B82F6';

    this.container.innerHTML = `
      <button id="compact-icon" style="
        width: 29px;
        height: 29px;
        display: flex;
        align-items: center;
        justify-content: center;
        font-size: 18px;
        position: relative;
        transition: all 0.2s ease;
        border: none;
        background: transparent;
        padding: 0;
        cursor: pointer;
        outline: none;
        -webkit-tap-highlight-color: transparent;
        touch-action: manipulation;
        pointer-events: auto;
        z-index: 1000;
        user-select: none;
        -webkit-user-select: none;
      ">
        <span style="
          display: flex;
          align-items: center;
          justify-content: center;
          width: 100%;
          height: 100%;
        ">${icon}</span>

        <!-- Индикатор значения -->
        <div style="
          position: absolute;
          bottom: -2px;
          right: -2px;
          background: ${bgColor};
          color: white;
          border-radius: 8px;
          padding: 1px 4px;
          font-size: 9px;
          font-weight: 600;
          line-height: 1;
          box-shadow: 0 1px 2px rgba(0,0,0,0.3);
          white-space: nowrap;
        ">
          ${this.getCompactValue()}
        </div>
      </button>
    `;
  }

  private renderExpanded(): void {
    const { mode, walkingTime, searchRadius } = this.props;

    // Используем временные значения если они есть, иначе берем из props
    const currentWalkingTime = this.tempWalkingTime ?? walkingTime;
    const currentSearchRadius = this.tempSearchRadius ?? searchRadius;

    // Конвертируем значения в проценты
    const radiusPercent = ((currentSearchRadius - 500) / (50000 - 500)) * 100;
    const walkingPercent = ((currentWalkingTime - 5) / (60 - 5)) * 100;
    const currentPercent = mode === 'walking' ? walkingPercent : radiusPercent;

    const icon = mode === 'walking' ? '🚶' : '📏';
    const color = mode === 'walking' ? '#10B981' : '#3B82F6';
    const displayValue = this.getDisplayValue();

    this.container.innerHTML = `
      <!-- Заголовок с иконкой и кнопкой закрытия -->
      <div style="display: flex; justify-content: space-between; align-items: center; margin-bottom: 8px;">
        <div style="display: flex; align-items: center; gap: 8px;">
          <span style="font-size: 18px;">${icon}</span>
          <span style="font-size: 14px; font-weight: 500; color: #374151;">
            ${mode === 'walking' ? this.props.translations?.walkingAccessibility || 'Walking accessibility' : this.props.translations?.searchRadius || 'Search radius'}
          </span>
        </div>
        <button id="close-btn" style="
          background: none;
          border: none;
          cursor: pointer;
          padding: 4px;
          color: #9CA3AF;
          font-size: 18px;
          line-height: 1;
          border-radius: 4px;
          transition: all 0.2s;
        " onmouseover="this.style.background='#F3F4F6'; this.style.color='#6B7280';"
           onmouseout="this.style.background='none'; this.style.color='#9CA3AF';">
          ×
        </button>
      </div>

      <!-- Слайдер -->
      <div style="position: relative;">
        <input
          type="range"
          id="distance-slider"
          min="0"
          max="100"
          value="${currentPercent}"
          style="
            width: 100%;
            height: 6px;
            border-radius: 3px;
            background: linear-gradient(to right,
              ${color} 0%,
              ${color} ${currentPercent}%,
              #e5e7eb ${currentPercent}%,
              #e5e7eb 100%);
            outline: none;
            -webkit-appearance: none;
            appearance: none;
            cursor: pointer;
            touch-action: manipulation;
          "
        />
        <style>
          #distance-slider::-webkit-slider-thumb {
            -webkit-appearance: none;
            appearance: none;
            width: 16px;
            height: 16px;
            border-radius: 50%;
            background: ${color};
            cursor: pointer;
            border: 2px solid white;
            box-shadow: 0 2px 4px rgba(0,0,0,0.2);
          }
          #distance-slider::-moz-range-thumb {
            width: 16px;
            height: 16px;
            border-radius: 50%;
            background: ${color};
            cursor: pointer;
            border: 2px solid white;
            box-shadow: 0 2px 4px rgba(0,0,0,0.2);
          }
        </style>
      </div>

      <!-- Значение и подсказка -->
      <div style="display: flex; justify-content: space-between; align-items: center; margin-top: 6px;">
        <span style="font-size: 11px; color: #6B7280;">
          ${this.props.translations?.changeModeHint || 'Click to change mode'}
        </span>
        <span style="font-size: 13px; font-weight: 600; color: ${color};">
          ${displayValue}
        </span>
      </div>
    `;
  }

  private getCompactValue(): string {
    const { mode, walkingTime, searchRadius, translations } = this.props;
    const t = translations || {
      minutes: 'min',
      km: 'km',
      m: 'm',
    };

    if (mode === 'walking') {
      return `${walkingTime}'`;
    } else {
      if (searchRadius >= 1000) {
        const km = (searchRadius / 1000).toFixed(0);
        return `${km}${t.km}`;
      }
      return `${searchRadius}${t.m}`;
    }
  }

  private getDisplayValue(): string {
    const { mode, walkingTime, searchRadius, translations } = this.props;
    const t = translations || {
      minutes: 'min',
      km: 'km',
      m: 'm',
    };

    // Используем временные значения если они есть
    const currentWalkingTime = this.tempWalkingTime ?? walkingTime;
    const currentSearchRadius = this.tempSearchRadius ?? searchRadius;

    if (mode === 'walking') {
      return `${currentWalkingTime} ${t.minutes}`;
    } else {
      if (currentSearchRadius >= 1000) {
        return `${(currentSearchRadius / 1000).toFixed(1)} ${t.km}`;
      }
      return `${currentSearchRadius} ${t.m}`;
    }
  }

  private addEventListeners(): void {
    // Добавляем слушатель на контейнер для отладки
    this.container.addEventListener(
      'touchstart',
      (e) => {
        console.log('[CompactSliderControl] Container touchstart detected!');
        console.log('Event target:', e.target);
        console.log('Current target:', e.currentTarget);
      },
      { capture: true }
    );

    if (!this.isExpanded) {
      this.addCompactListeners();
    } else {
      this.addExpandedListeners();
    }
  }

  private addCompactListeners(): void {
    const icon = this.container.querySelector('#compact-icon') as HTMLElement;
    if (!icon) return;

    // Обработка одиночного клика/тапа - переключение режима
    const handleSingleTap = () => {
      const newMode = this.props.mode === 'walking' ? 'radius' : 'walking';
      console.log(
        '[CompactSliderControl] Single tap - switching mode from',
        this.props.mode,
        'to',
        newMode
      );
      this.props.onModeChange(newMode);

      // Визуальная обратная связь
      icon.style.transform = 'scale(0.95)';
      setTimeout(() => {
        icon.style.transform = 'scale(1)';
      }, 100);
    };

    // Обработка двойного тапа
    const handleDoubleTap = () => {
      console.log(
        '[CompactSliderControl] Double tap - toggling expanded state'
      );
      this.toggleExpanded();
    };

    // Обработчик мыши для десктопа - отдельная логика для двойного клика
    let mouseClickCount = 0;
    let mouseClickTimer: NodeJS.Timeout | null = null;

    icon.addEventListener('click', (e) => {
      e.preventDefault();
      e.stopPropagation();

      mouseClickCount++;
      console.log(
        '[CompactSliderControl] Mouse click detected, count:',
        mouseClickCount
      );

      if (mouseClickTimer) {
        clearTimeout(mouseClickTimer);
      }

      if (mouseClickCount === 1) {
        mouseClickTimer = setTimeout(() => {
          if (mouseClickCount === 1) {
            // Одиночный клик - переключение режима
            console.log(
              '[CompactSliderControl] Single mouse click - switching mode'
            );
            handleSingleTap();
          }
          mouseClickCount = 0;
        }, 250);
      } else if (mouseClickCount === 2) {
        // Двойной клик - открытие/закрытие слайдера
        console.log(
          '[CompactSliderControl] Double mouse click - toggling expanded'
        );
        handleDoubleTap();
        mouseClickCount = 0;
        if (mouseClickTimer) {
          clearTimeout(mouseClickTimer);
        }
      }
    });

    // Long press и double tap для мобильных
    let touchStartTime = 0;
    let touchTimer: NodeJS.Timeout | null = null;
    let lastTouchTime = 0;

    // Используем pointerdown вместо touchstart для лучшей совместимости
    icon.addEventListener('pointerdown', (e) => {
      console.log('[CompactSliderControl] Pointer down detected');
      e.preventDefault();
      e.stopPropagation();
      touchStartTime = Date.now();

      // Визуальная обратная связь
      icon.style.opacity = '0.7';
      icon.style.transform = 'scale(0.95)';

      touchTimer = setTimeout(() => {
        // Long press detected
        console.log(
          '[CompactSliderControl] Long press detected - toggling expanded'
        );
        this.toggleExpanded();

        // Вибрация на мобильных (если поддерживается)
        if ('vibrate' in navigator) {
          navigator.vibrate(50);
        }

        // Возвращаем визуальное состояние
        icon.style.opacity = '1';
        icon.style.transform = 'scale(1)';
      }, 500); // 500ms для long press
    });

    icon.addEventListener('pointerup', (e) => {
      console.log('[CompactSliderControl] Pointer up detected');
      e.preventDefault();
      e.stopPropagation();

      // Возвращаем визуальное состояние
      icon.style.opacity = '1';
      icon.style.transform = 'scale(1)';

      if (touchTimer) {
        clearTimeout(touchTimer);
      }

      const touchDuration = Date.now() - touchStartTime;
      const currentTime = Date.now();

      if (touchDuration < 500) {
        // Проверка на двойной тап
        if (currentTime - lastTouchTime < 300) {
          console.log(
            '[CompactSliderControl] Double tap detected - toggling expanded'
          );
          this.toggleExpanded();
          lastTouchTime = 0; // Сбрасываем для избежания тройного тапа
        } else {
          // Короткое касание - переключаем режим
          console.log('[CompactSliderControl] Short tap - switching mode');
          handleSingleTap();
          lastTouchTime = currentTime;
        }
      }
    });

    icon.addEventListener('pointercancel', () => {
      if (touchTimer) {
        clearTimeout(touchTimer);
      }
      // Возвращаем визуальное состояние при отмене
      icon.style.opacity = '1';
      icon.style.transform = 'scale(1)';
    });

    // Hover эффект для десктопа
    icon.addEventListener('mouseenter', () => {
      icon.style.transform = 'scale(1.05)';
    });

    icon.addEventListener('mouseleave', () => {
      icon.style.transform = 'scale(1)';
    });
  }

  private addExpandedListeners(): void {
    // Кнопка закрытия
    const closeBtn = this.container.querySelector('#close-btn') as HTMLElement;
    if (closeBtn) {
      closeBtn.addEventListener('click', (e) => {
        e.preventDefault();
        e.stopPropagation();
        this.toggleExpanded();
      });
    }

    // Слайдер
    const slider = this.container.querySelector(
      '#distance-slider'
    ) as HTMLInputElement;
    if (slider) {
      // Обработка движения слайдера - обновляем только временные значения
      const handleSliderInput = (e: Event) => {
        const target = e.target as HTMLInputElement;
        const percent = parseFloat(target.value);

        if (this.props.mode === 'walking') {
          const minutes = Math.round(5 + (percent / 100) * (60 - 5));
          this.tempWalkingTime = minutes;
          console.log(
            '[CompactSliderControl] Slider input - temp walking time:',
            minutes,
            'min'
          );
        } else {
          const meters = Math.round(500 + (percent / 100) * (50000 - 500));
          this.tempSearchRadius = meters;
          console.log(
            '[CompactSliderControl] Slider input - temp radius:',
            meters,
            'm'
          );
        }

        this.updateSliderBackground(percent);
        // Обновляем отображение значения
        this.updateDisplayValue();
      };

      // Обработка окончания движения слайдера - применяем изменения
      const handleSliderEnd = () => {
        if (this.props.mode === 'walking' && this.tempWalkingTime !== null) {
          console.log(
            '[CompactSliderControl] Applying walking time:',
            this.tempWalkingTime,
            'min'
          );
          this.props.onWalkingTimeChange(this.tempWalkingTime);
          this.tempWalkingTime = null;
        } else if (
          this.props.mode === 'radius' &&
          this.tempSearchRadius !== null
        ) {
          console.log(
            '[CompactSliderControl] Applying radius:',
            this.tempSearchRadius,
            'm'
          );
          this.props.onRadiusChange(this.tempSearchRadius);
          this.tempSearchRadius = null;
        }
      };

      // События для отслеживания движения слайдера
      slider.addEventListener('input', handleSliderInput);

      // События для применения изменений только при отпускании
      slider.addEventListener('change', handleSliderEnd);

      // Предотвращаем сброс значения при клике
      slider.addEventListener('mousedown', (e) => {
        e.stopPropagation();
      });

      slider.addEventListener('touchstart', (e) => {
        e.stopPropagation();
      });
    }

    // Клик по заголовку для переключения режима
    const header = this.container.querySelector('div') as HTMLElement;
    if (header) {
      header.style.cursor = 'pointer';
      header.addEventListener('click', (e) => {
        const target = e.target as HTMLElement;
        // Не переключаем если кликнули на кнопку закрытия
        if (!target.closest('#close-btn')) {
          const newMode = this.props.mode === 'walking' ? 'radius' : 'walking';
          this.props.onModeChange(newMode);
          // Сбрасываем временные значения при смене режима
          this.tempWalkingTime = null;
          this.tempSearchRadius = null;
        }
      });
    }
  }

  private updateSliderBackground(percent: number): void {
    const slider = this.container.querySelector(
      '#distance-slider'
    ) as HTMLInputElement;
    if (slider) {
      const color = this.props.mode === 'walking' ? '#10B981' : '#3B82F6';
      slider.style.background = `linear-gradient(to right,
        ${color} 0%,
        ${color} ${percent}%,
        #e5e7eb ${percent}%,
        #e5e7eb 100%)`;
    }
  }

  private updateDisplayValue(): void {
    // Находим элемент отображения значения
    const valueElement = this.container.querySelector(
      'span:last-child'
    ) as HTMLElement;
    if (
      valueElement &&
      valueElement.parentElement?.style.justifyContent === 'space-between'
    ) {
      valueElement.textContent = this.getDisplayValue();
    }
  }

  private toggleExpanded(): void {
    this.isExpanded = !this.isExpanded;

    // Сбрасываем временные значения при закрытии/открытии
    this.tempWalkingTime = null;
    this.tempSearchRadius = null;

    // Управляем обработчиком кликов вне контрола
    if (this.isExpanded) {
      this.setupOutsideClickHandler();
    } else {
      this.removeOutsideClickHandler();
    }

    // Принудительное обновление стилей и содержимого
    setTimeout(() => {
      this.updateContainerStyle();
      this.render();
    }, 10);

    console.log('[CompactSliderControl] Toggled expanded:', this.isExpanded);
  }

  private outsideClickHandler = (e: MouseEvent) => {
    const target = e.target as HTMLElement;
    if (!this.container.contains(target) && this.isExpanded) {
      this.toggleExpanded();
    }
  };

  private setupOutsideClickHandler(): void {
    // Добавляем слушатель с небольшой задержкой, чтобы не сработало сразу при открытии
    setTimeout(() => {
      document.addEventListener('click', this.outsideClickHandler, true);
    }, 100);
  }

  private removeOutsideClickHandler(): void {
    document.removeEventListener('click', this.outsideClickHandler, true);
  }
}

interface CompactSliderControlComponentProps extends CompactSliderControlProps {
  map: MapboxMap | null;
}

const CompactSliderControl: React.FC<CompactSliderControlComponentProps> = ({
  map,
  position: _position = 'top-right',
  isFullscreen = false,
  isMobile = false,
  ...props
}) => {
  const controlRef = useRef<CompactSliderControlClass | null>(null);
  const [isAdded, setIsAdded] = useState(false);

  // Позиционирование справа под кнопкой fullscreen
  const adaptivePosition = 'top-right';

  useEffect(() => {
    if (!map) {
      console.log('[CompactSliderControl] Map not ready yet');
      return;
    }

    console.log('[CompactSliderControl] Creating control with props:', {
      mode: props.mode,
      walkingTime: props.walkingTime,
      searchRadius: props.searchRadius,
      position: adaptivePosition,
      isFullscreen,
      isMobile,
    });

    const control = new CompactSliderControlClass(
      { ...props, isMobile, isFullscreen },
      (newProps) => {
        console.log(
          '[CompactSliderControl] Props update from control:',
          newProps
        );
      }
    );

    map.addControl(control, adaptivePosition);
    controlRef.current = control;
    setIsAdded(true);

    console.log('[CompactSliderControl] Control added to map successfully');

    return () => {
      if (controlRef.current && map) {
        try {
          map.removeControl(controlRef.current);
          console.log('[CompactSliderControl] Control removed from map');
        } catch (error) {
          console.warn('[CompactSliderControl] Error removing control:', error);
        }
      }
      controlRef.current = null;
      setIsAdded(false);
    };
  }, [map, adaptivePosition, isFullscreen, isMobile, props]);

  useEffect(() => {
    if (controlRef.current && isAdded) {
      controlRef.current.updateProps({ ...props, isMobile, isFullscreen });
    }
  }, [props, isAdded, isMobile, isFullscreen]);

  return null;
};

export default CompactSliderControl;
