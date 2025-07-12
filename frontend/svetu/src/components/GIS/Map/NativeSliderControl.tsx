import React, { useEffect, useRef, useState } from 'react';
import { IControl, Map as MapboxMap } from 'mapbox-gl';

interface NativeSliderControlProps {
  mode: 'radius' | 'walking';
  onModeChange: (mode: 'radius' | 'walking') => void;
  walkingTime: number; // в минутах (5-60)
  onWalkingTimeChange: (time: number) => void;
  searchRadius: number; // в метрах (500-50000)
  onRadiusChange: (radius: number) => void;
  position?: 'top-left' | 'top-right' | 'bottom-left' | 'bottom-right';
  isFullscreen?: boolean; // детекция fullscreen режима
  isMobile?: boolean; // детекция мобильного устройства
}

class SliderControl implements IControl {
  private map: MapboxMap | undefined;
  private container: HTMLDivElement;
  private props: NativeSliderControlProps;
  private onPropsChange: (props: NativeSliderControlProps) => void;

  constructor(
    props: NativeSliderControlProps,
    onPropsChange: (props: NativeSliderControlProps) => void
  ) {
    this.props = props;
    this.onPropsChange = onPropsChange;
    this.container = document.createElement('div');
  }

  onAdd(map: MapboxMap): HTMLElement {
    this.map = map;
    this.container.className = 'mapboxgl-ctrl mapboxgl-ctrl-group';
    this.container.style.background = 'white';
    this.container.style.borderRadius = '8px';
    this.container.style.boxShadow = '0 2px 10px rgba(0,0,0,0.15)';
    this.container.style.padding = this.props.isMobile
      ? '8px 10px'
      : '10px 12px';
    this.container.style.minWidth = this.props.isMobile ? '200px' : '240px';
    this.container.style.maxWidth = this.props.isMobile ? '240px' : '280px';
    this.container.style.userSelect = 'none';
    this.container.style.margin = this.props.isMobile ? '4px' : '6px';
    this.container.style.zIndex = '1002'; // выше навигационных контролов MapBox чтобы быть видимым
    this.container.style.position = 'relative';

    // Дополнительные стили для компактности в верхнем позиционировании
    if (this.props.isMobile) {
      this.container.style.fontSize = '13px';
      this.container.style.marginTop = '4px'; // минимальный отступ сверху на мобильных
    }

    this.render();
    return this.container;
  }

  onRemove(): void {
    if (this.container.parentNode) {
      this.container.parentNode.removeChild(this.container);
    }
    this.map = undefined;
  }

  updateProps(newProps: NativeSliderControlProps): void {
    this.props = newProps;
    this.render();
  }

  private render(): void {
    const { mode, walkingTime, searchRadius } = this.props;

    // Конвертируем значения в проценты для ползунка
    const radiusPercent = ((searchRadius - 500) / (50000 - 500)) * 100;
    const walkingPercent = ((walkingTime - 5) / (60 - 5)) * 100;

    const currentPercent = mode === 'walking' ? walkingPercent : radiusPercent;

    // Форматируем отображаемое значение
    const displayValue =
      mode === 'walking'
        ? `${walkingTime} мин`
        : searchRadius >= 1000
          ? `${(searchRadius / 1000).toFixed(1)} км`
          : `${searchRadius} м`;

    this.container.innerHTML = `
      <div style="display: flex; align-items: center; gap: 8px; font-family: -apple-system, BlinkMacSystemFont, 'Segoe UI', Roboto, sans-serif;">
        <!-- Иконка режима -->
        <div style="display: flex; align-items: center; cursor: pointer; padding: 4px;" id="mode-toggle">
          <span style="font-size: 16px;">${mode === 'walking' ? '🚶' : '📍'}</span>
        </div>
        
        <!-- Ползунок -->
        <div style="flex: 1; position: relative;">
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
                ${mode === 'walking' ? '#10B981' : '#3B82F6'} 0%, 
                ${mode === 'walking' ? '#10B981' : '#3B82F6'} ${currentPercent}%, 
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
              width: 20px;
              height: 20px;
              border-radius: 50%;
              background: ${mode === 'walking' ? '#10B981' : '#3B82F6'};
              cursor: pointer;
              border: 2px solid white;
              box-shadow: 0 2px 4px rgba(0,0,0,0.2);
            }
            #distance-slider::-moz-range-thumb {
              width: 20px;
              height: 20px;
              border-radius: 50%;
              background: ${mode === 'walking' ? '#10B981' : '#3B82F6'};
              cursor: pointer;
              border: 2px solid white;
              box-shadow: 0 2px 4px rgba(0,0,0,0.2);
            }
          </style>
        </div>
        
        <!-- Отображение значения -->
        <div style="
          min-width: 60px; 
          text-align: center; 
          font-size: 12px; 
          font-weight: 500;
          color: ${mode === 'walking' ? '#10B981' : '#3B82F6'};
        ">
          ${displayValue}
        </div>
      </div>
      
      <!-- Подсказка -->
      <div style="
        font-size: 10px; 
        color: #6b7280; 
        text-align: center; 
        margin-top: 4px;
        line-height: 1.2;
      ">
        ${
          mode === 'walking'
            ? 'Пешая доступность • Нажмите 📍 для радиуса'
            : 'Радиус поиска • Нажмите 🚶 для времени ходьбы'
        }
      </div>
    `;

    // Добавляем обработчики событий
    this.addEventListeners();
  }

  private addEventListeners(): void {
    const slider = this.container.querySelector(
      '#distance-slider'
    ) as HTMLInputElement;
    const modeToggle = this.container.querySelector(
      '#mode-toggle'
    ) as HTMLElement;

    if (slider) {
      // Обработка изменения ползунка
      const handleSliderChange = (e: Event) => {
        const target = e.target as HTMLInputElement;
        const percent = parseFloat(target.value);

        console.log(
          '[NativeSliderControl] Slider change:',
          percent,
          'mode:',
          this.props.mode
        );

        if (this.props.mode === 'walking') {
          // Конвертируем проценты в минуты (5-60)
          const minutes = Math.round(5 + (percent / 100) * (60 - 5));
          this.props.onWalkingTimeChange(minutes);
        } else {
          // Конвертируем проценты в метры (500-50000)
          const meters = Math.round(500 + (percent / 100) * (50000 - 500));
          this.props.onRadiusChange(meters);
        }

        // Обновляем визуальное отображение
        this.updateSliderBackground(percent);
      };

      // Обработка в реальном времени
      slider.addEventListener('input', handleSliderChange);
      slider.addEventListener('change', handleSliderChange);

      // Touch события для мобильных
      slider.addEventListener('touchstart', (e) => {
        e.stopPropagation();
      });

      slider.addEventListener('touchmove', (e) => {
        e.stopPropagation();
        handleSliderChange(e);
      });

      slider.addEventListener('touchend', (e) => {
        e.stopPropagation();
        handleSliderChange(e);
      });
    }

    if (modeToggle) {
      // Переключение режима
      const handleModeToggle = () => {
        const newMode = this.props.mode === 'walking' ? 'radius' : 'walking';
        console.log(
          '[NativeSliderControl] Mode toggle:',
          this.props.mode,
          '->',
          newMode
        );
        this.props.onModeChange(newMode);
      };

      modeToggle.addEventListener('click', handleModeToggle);
      modeToggle.addEventListener('touchend', (e) => {
        e.preventDefault();
        handleModeToggle();
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
}

interface NativeSliderControlComponentProps extends NativeSliderControlProps {
  map: MapboxMap | null;
}

const NativeSliderControl: React.FC<NativeSliderControlComponentProps> = ({
  map,
  position = 'top-right',
  isFullscreen = false,
  isMobile = false,
  ...props
}) => {
  const controlRef = useRef<SliderControl | null>(null);
  const [isAdded, setIsAdded] = useState(false);

  // Определяем оптимальную позицию на основе режима экрана и устройства
  const adaptivePosition = (() => {
    if (isFullscreen) {
      // В полноэкранном режиме - слева сверху, чтобы не конфликтовать с выбором стиля
      return 'top-left';
    }
    if (isMobile) {
      // На мобильных - справа сверху, чтобы всегда быть видимым
      return 'top-right';
    }
    // На десктопе - справа сверху, чтобы не уходить за границы
    return 'top-right';
  })();

  useEffect(() => {
    if (!map) return;

    // Создаем контрол с всеми пропсами включая isMobile
    const control = new SliderControl(
      { ...props, isMobile, isFullscreen },
      (newProps) => {
        // Коллбек для обновления props из контрола
        console.log(
          '[NativeSliderControl] Props update from control:',
          newProps
        );
      }
    );

    // Добавляем контрол на карту с адаптивной позицией
    map.addControl(control, adaptivePosition);
    controlRef.current = control;
    setIsAdded(true);

    console.log(
      '[NativeSliderControl] Added to map at position:',
      adaptivePosition,
      'isMobile:',
      isMobile,
      'isFullscreen:',
      isFullscreen
    );

    return () => {
      if (controlRef.current && map) {
        try {
          map.removeControl(controlRef.current);
          console.log('[NativeSliderControl] Removed from map');
        } catch (error) {
          console.warn('[NativeSliderControl] Error removing control:', error);
        }
      }
      controlRef.current = null;
      setIsAdded(false);
    };
  }, [map, adaptivePosition, isFullscreen, isMobile]);

  // Обновляем props контрола при их изменении
  useEffect(() => {
    if (controlRef.current && isAdded) {
      controlRef.current.updateProps({ ...props, isMobile, isFullscreen });
    }
  }, [props, isAdded, isMobile, isFullscreen]);

  return null; // Компонент ничего не рендерит в React DOM
};

export default NativeSliderControl;
