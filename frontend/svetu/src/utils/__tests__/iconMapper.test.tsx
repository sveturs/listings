import React from 'react';
import { render } from '@testing-library/react';
import '@testing-library/jest-dom';
import { getCategoryIcon, renderCategoryIcon } from '../iconMapper';
import {
  Car,
  Truck,
  Bike,
  Ship,
  Factory,
  Tractor,
  Wheat,
  Cog,
  Wrench,
  Settings,
  Globe,
  Flag,
  CreditCard,
  Star,
  Clock,
  Calendar,
  Battery,
  Zap,
  Leaf,
  Crown,
  Shield,
  Music,
  Map,
  Mountain,
  Building2,
  Anchor,
  Sailboat,
  Package,
} from 'lucide-react';

describe('iconMapper', () => {
  describe('getCategoryIcon', () => {
    describe('Базовая функциональность', () => {
      test('возвращает null для пустого имени', () => {
        expect(getCategoryIcon('')).toBeNull();
        expect(getCategoryIcon(undefined)).toBeNull();
      });

      test('возвращает Package для неизвестного имени', () => {
        const IconComponent = getCategoryIcon('unknown-icon-name');
        expect(IconComponent).toBe(Package);
      });

      test('не чувствителен к регистру', () => {
        expect(getCategoryIcon('CAR')).toBe(getCategoryIcon('car'));
        expect(getCategoryIcon('Truck')).toBe(getCategoryIcon('truck'));
        expect(getCategoryIcon('BIKE')).toBe(getCategoryIcon('bike'));
      });
    });

    describe('Транспортные иконки', () => {
      test('возвращает правильную иконку для car', () => {
        expect(getCategoryIcon('car')).toBe(Car);
      });

      test('возвращает правильную иконку для truck', () => {
        expect(getCategoryIcon('truck')).toBe(Truck);
      });

      test('возвращает правильную иконку для motorcycle', () => {
        expect(getCategoryIcon('motorcycle')).toBe(Bike);
      });

      test('использует Truck для bus и van', () => {
        expect(getCategoryIcon('bus')).toBe(Truck);
        expect(getCategoryIcon('van')).toBe(Truck);
        expect(getCategoryIcon('trailer')).toBe(Truck);
      });

      test('использует Bike для scooter и bicycle', () => {
        expect(getCategoryIcon('scooter')).toBe(Bike);
        expect(getCategoryIcon('bicycle')).toBe(Bike);
        expect(getCategoryIcon('quad-bike')).toBe(Bike);
      });

      test('обрабатывает все транспортные иконки', () => {
        const transportIcons = [
          'car',
          'truck',
          'motorcycle',
          'bus',
          'van',
          'trailer',
        ];

        transportIcons.forEach((iconName) => {
          const IconComponent = getCategoryIcon(iconName);
          expect(IconComponent).toBeDefined();
          expect(IconComponent).not.toBe(Package);
        });
      });
    });

    describe('Водный транспорт', () => {
      test('возвращает правильную иконку для ship', () => {
        expect(getCategoryIcon('ship')).toBe(Ship);
      });

      test('возвращает правильную иконку для sailboat', () => {
        expect(getCategoryIcon('sailboat')).toBe(Sailboat);
      });

      test('возвращает правильную иконку для anchor', () => {
        expect(getCategoryIcon('anchor')).toBe(Anchor);
      });

      test('использует Ship для water', () => {
        expect(getCategoryIcon('water')).toBe(Ship);
      });
    });

    describe('Индустриальные иконки', () => {
      test('возвращает правильную иконку для factory', () => {
        expect(getCategoryIcon('factory')).toBe(Factory);
      });

      test('возвращает правильную иконку для tractor', () => {
        expect(getCategoryIcon('tractor')).toBe(Tractor);
      });

      test('возвращает правильную иконку для wheat', () => {
        expect(getCategoryIcon('wheat')).toBe(Wheat);
      });
    });

    describe('Технические иконки', () => {
      test('возвращает правильную иконку для cog', () => {
        expect(getCategoryIcon('cog')).toBe(Cog);
      });

      test('возвращает правильную иконку для wrench', () => {
        expect(getCategoryIcon('wrench')).toBe(Wrench);
      });

      test('возвращает правильную иконку для gear', () => {
        expect(getCategoryIcon('gear')).toBe(Settings);
      });

      test('использует Wrench для tools', () => {
        expect(getCategoryIcon('tools')).toBe(Wrench);
      });

      test('использует Cog для engine', () => {
        expect(getCategoryIcon('engine')).toBe(Cog);
      });
    });

    describe('Общие иконки', () => {
      test('возвращает правильную иконку для globe', () => {
        expect(getCategoryIcon('globe')).toBe(Globe);
      });

      test('возвращает правильную иконку для flag', () => {
        expect(getCategoryIcon('flag')).toBe(Flag);
      });

      test('возвращает правильную иконку для star', () => {
        expect(getCategoryIcon('star')).toBe(Star);
      });

      test('возвращает правильную иконку для clock', () => {
        expect(getCategoryIcon('clock')).toBe(Clock);
      });

      test('возвращает правильную иконку для calendar', () => {
        expect(getCategoryIcon('calendar')).toBe(Calendar);
      });

      test('возвращает правильную иконку для battery', () => {
        expect(getCategoryIcon('battery')).toBe(Battery);
      });

      test('возвращает правильную иконку для bolt/zap', () => {
        expect(getCategoryIcon('bolt')).toBe(Zap);
      });

      test('возвращает правильную иконку для leaf', () => {
        expect(getCategoryIcon('leaf')).toBe(Leaf);
      });

      test('возвращает правильную иконку для crown', () => {
        expect(getCategoryIcon('crown')).toBe(Crown);
      });

      test('возвращает правильную иконку для shield', () => {
        expect(getCategoryIcon('shield')).toBe(Shield);
      });

      test('возвращает правильную иконку для music', () => {
        expect(getCategoryIcon('music')).toBe(Music);
      });

      test('возвращает правильную иконку для map', () => {
        expect(getCategoryIcon('map')).toBe(Map);
      });

      test('возвращает правильную иконку для mountain', () => {
        expect(getCategoryIcon('mountain')).toBe(Mountain);
      });

      test('возвращает правильную иконку для city', () => {
        expect(getCategoryIcon('city')).toBe(Building2);
      });

      test('использует CreditCard для id-card', () => {
        expect(getCategoryIcon('id-card')).toBe(CreditCard);
      });
    });

    describe('Специальные транспортные средства', () => {
      test('использует Car для racing', () => {
        expect(getCategoryIcon('racing')).toBe(Car);
      });

      test('использует Car для car-side', () => {
        expect(getCategoryIcon('car-side')).toBe(Car);
      });

      test('использует Truck для caravan', () => {
        expect(getCategoryIcon('caravan')).toBe(Truck);
      });

      test('использует Zap для speed', () => {
        expect(getCategoryIcon('speed')).toBe(Zap);
      });

      test('использует Clock для vintage', () => {
        expect(getCategoryIcon('vintage')).toBe(Clock);
      });

      test('использует Star для gem и snowflake', () => {
        expect(getCategoryIcon('gem')).toBe(Star);
        expect(getCategoryIcon('snowflake')).toBe(Star);
      });

      test('использует Flag для golf', () => {
        expect(getCategoryIcon('golf')).toBe(Flag);
      });

      test('использует Shield для triangle', () => {
        expect(getCategoryIcon('triangle')).toBe(Shield);
      });
    });
  });

  describe('renderCategoryIcon', () => {
    describe('Базовая функциональность', () => {
      test('возвращает null для пустого имени', () => {
        expect(renderCategoryIcon('')).toBeNull();
        expect(renderCategoryIcon(undefined)).toBeNull();
      });

      test('рендерит иконку компонент', () => {
        const { container } = render(
          <>{renderCategoryIcon('car', 'w-6 h-6')}</>
        );

        const svg = container.querySelector('svg');
        expect(svg).toBeInTheDocument();
        expect(svg).toHaveClass('w-6', 'h-6');
      });

      test('применяет custom className', () => {
        const { container } = render(
          <>{renderCategoryIcon('car', 'custom-class')}</>
        );

        const svg = container.querySelector('svg');
        expect(svg).toHaveClass('custom-class');
      });

      test('применяет несколько классов', () => {
        const { container } = render(
          <>{renderCategoryIcon('car', 'w-8 h-8 text-blue-500')}</>
        );

        const svg = container.querySelector('svg');
        expect(svg).toHaveClass('w-8', 'h-8', 'text-blue-500');
      });
    });

    describe('Рендеринг иконок', () => {
      test('рендерит транспортные иконки', () => {
        const { container: container1 } = render(
          <>{renderCategoryIcon('car')}</>
        );
        const { container: container2 } = render(
          <>{renderCategoryIcon('truck')}</>
        );
        const { container: container3 } = render(
          <>{renderCategoryIcon('motorcycle')}</>
        );

        expect(container1.querySelector('svg')).toBeInTheDocument();
        expect(container2.querySelector('svg')).toBeInTheDocument();
        expect(container3.querySelector('svg')).toBeInTheDocument();
      });

      test('рендерит индустриальные иконки', () => {
        const { container } = render(<>{renderCategoryIcon('factory')}</>);

        expect(container.querySelector('svg')).toBeInTheDocument();
      });

      test('рендерит технические иконки', () => {
        const { container } = render(<>{renderCategoryIcon('cog')}</>);

        expect(container.querySelector('svg')).toBeInTheDocument();
      });

      test('рендерит общие иконки', () => {
        const { container } = render(<>{renderCategoryIcon('star')}</>);

        expect(container.querySelector('svg')).toBeInTheDocument();
      });

      test('рендерит Package для неизвестных иконок', () => {
        const { container } = render(<>{renderCategoryIcon('unknown-icon')}</>);

        expect(container.querySelector('svg')).toBeInTheDocument();
      });
    });

    describe('Поддержка эмодзи', () => {
      test('рендерит эмодзи как текст', () => {
        const { container } = render(
          <>{renderCategoryIcon('🚗', 'text-2xl')}</>
        );

        const span = container.querySelector('span');
        expect(span).toBeInTheDocument();
        expect(span).toHaveTextContent('🚗');
        expect(span).toHaveClass('text-2xl');
      });

      test('обрабатывает различные эмодзи', () => {
        const emojis = ['🚗', '🏠', '📱', '⚽', '🎮', '🍕'];

        emojis.forEach((emoji) => {
          const { container } = render(<>{renderCategoryIcon(emoji)}</>);

          const span = container.querySelector('span');
          expect(span).toHaveTextContent(emoji);
        });
      });

      test('применяет className к эмодзи', () => {
        const { container } = render(
          <>{renderCategoryIcon('🚗', 'custom-emoji-class')}</>
        );

        const span = container.querySelector('span');
        expect(span).toHaveClass('custom-emoji-class');
      });

      test('обрабатывает многобайтные эмодзи', () => {
        const complexEmojis = ['👨‍👩‍👧‍👦', '🏳️‍🌈', '👍🏻'];

        complexEmojis.forEach((emoji) => {
          const { container } = render(<>{renderCategoryIcon(emoji)}</>);

          const span = container.querySelector('span');
          expect(span).toBeInTheDocument();
        });
      });
    });

    describe('Без className', () => {
      test('рендерит без className', () => {
        const { container } = render(<>{renderCategoryIcon('car')}</>);

        const svg = container.querySelector('svg');
        expect(svg).toBeInTheDocument();
      });

      test('эмодзи без className', () => {
        const { container } = render(<>{renderCategoryIcon('🚗')}</>);

        const span = container.querySelector('span');
        expect(span).toBeInTheDocument();
        expect(span).toHaveTextContent('🚗');
      });
    });

    describe('Edge cases', () => {
      test('обрабатывает пробелы в иконке', () => {
        const { container } = render(<>{renderCategoryIcon('  car  ')}</>);

        expect(container.querySelector('svg')).toBeInTheDocument();
      });

      test('обрабатывает иконки с дефисами', () => {
        const { container } = render(<>{renderCategoryIcon('car-side')}</>);

        expect(container.querySelector('svg')).toBeInTheDocument();
      });

      test('обрабатывает смешанный регистр', () => {
        const { container } = render(<>{renderCategoryIcon('CaR')}</>);

        expect(container.querySelector('svg')).toBeInTheDocument();
      });

      test('обрабатывает спецсимволы (не эмодзи)', () => {
        const { container } = render(<>{renderCategoryIcon('@#$')}</>);

        // Должен вернуть Package icon
        expect(container.querySelector('svg')).toBeInTheDocument();
      });

      test('обрабатывает числа', () => {
        const { container } = render(<>{renderCategoryIcon('123')}</>);

        // Должен вернуть Package icon
        expect(container.querySelector('svg')).toBeInTheDocument();
      });
    });

    describe('Регистронезависимость', () => {
      test('работает с UPPERCASE', () => {
        const { container } = render(<>{renderCategoryIcon('CAR')}</>);

        expect(container.querySelector('svg')).toBeInTheDocument();
      });

      test('работает с lowercase', () => {
        const { container } = render(<>{renderCategoryIcon('car')}</>);

        expect(container.querySelector('svg')).toBeInTheDocument();
      });

      test('работает с MixedCase', () => {
        const { container } = render(<>{renderCategoryIcon('CaR')}</>);

        expect(container.querySelector('svg')).toBeInTheDocument();
      });
    });
  });

  describe('Интеграция getCategoryIcon и renderCategoryIcon', () => {
    test('renderCategoryIcon использует getCategoryIcon', () => {
      const iconName = 'car';
      const IconComponent = getCategoryIcon(iconName);
      const { container } = render(<>{renderCategoryIcon(iconName)}</>);

      expect(IconComponent).toBeDefined();
      expect(container.querySelector('svg')).toBeInTheDocument();
    });

    test('оба возвращают null для пустых значений', () => {
      expect(getCategoryIcon('')).toBeNull();
      expect(renderCategoryIcon('')).toBeNull();
    });

    test('оба обрабатывают неизвестные иконки', () => {
      const unknownIcon = 'totally-unknown-icon';
      const IconComponent = getCategoryIcon(unknownIcon);
      const { container } = render(<>{renderCategoryIcon(unknownIcon)}</>);

      expect(IconComponent).toBe(Package);
      expect(container.querySelector('svg')).toBeInTheDocument();
    });
  });

  describe('Полный набор иконок', () => {
    test('проверяет наличие всех основных категорий', () => {
      const categories = [
        'car',
        'truck',
        'motorcycle',
        'ship',
        'factory',
        'tractor',
        'cog',
        'wrench',
        'star',
        'clock',
      ];

      categories.forEach((category) => {
        const IconComponent = getCategoryIcon(category);
        expect(IconComponent).toBeDefined();
        expect(IconComponent).not.toBeNull();
      });
    });

    test('все иконки рендерятся без ошибок', () => {
      const allIcons = [
        'car',
        'truck',
        'bus',
        'van',
        'motorcycle',
        'scooter',
        'bicycle',
        'ship',
        'anchor',
        'sailboat',
        'factory',
        'tractor',
        'wheat',
        'cog',
        'gear',
        'wrench',
        'globe',
        'flag',
        'star',
        'clock',
        'calendar',
        'battery',
        'bolt',
        'leaf',
        'crown',
        'shield',
        'music',
        'map',
        'mountain',
        'city',
      ];

      allIcons.forEach((iconName) => {
        expect(() => {
          render(<>{renderCategoryIcon(iconName)}</>);
        }).not.toThrow();
      });
    });
  });
});
