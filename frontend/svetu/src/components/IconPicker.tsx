'use client';

import { useState } from 'react';

interface IconPickerProps {
  value: string;
  onChange: (icon: string) => void;
  placeholder?: string;
}

// Расширенный набор иконок для категорий и атрибутов
const iconCategories = {
  transport: {
    name: 'Транспорт',
    icons: [
      '🚗',
      '🚙',
      '🚌',
      '🚐',
      '🏎️',
      '🚓',
      '🚑',
      '🚒',
      '🚜',
      '🛺',
      '🚲',
      '🛵',
      '🏍️',
      '✈️',
      '🚁',
      '🛸',
      '🚀',
      '🛥️',
      '⛵',
      '🚢',
    ],
  },
  electronics: {
    name: 'Электроника',
    icons: [
      '📱',
      '💻',
      '🖥️',
      '⌨️',
      '🖱️',
      '🖨️',
      '📷',
      '📹',
      '📺',
      '📻',
      '🎮',
      '🕹️',
      '💿',
      '💾',
      '💽',
      '📀',
      '🔋',
      '🔌',
      '💡',
      '🔦',
    ],
  },
  home: {
    name: 'Дом и быт',
    icons: [
      '🏠',
      '🏡',
      '🏢',
      '🏬',
      '🏭',
      '🛏️',
      '🛋️',
      '🪑',
      '🚪',
      '🪟',
      '🛁',
      '🚿',
      '🚽',
      '🧹',
      '🧽',
      '🧴',
      '🧷',
      '📌',
      '✂️',
      '🔧',
    ],
  },
  clothing: {
    name: 'Одежда',
    icons: [
      '👕',
      '👔',
      '👗',
      '👘',
      '🥻',
      '👖',
      '👚',
      '🧥',
      '🧦',
      '🩱',
      '👙',
      '👟',
      '👞',
      '🥾',
      '👑',
      '👒',
      '🧢',
      '🎩',
      '🧣',
      '🧤',
    ],
  },
  food: {
    name: 'Еда и напитки',
    icons: [
      '🍎',
      '🍌',
      '🍇',
      '🍊',
      '🍋',
      '🥭',
      '🍅',
      '🥑',
      '🥦',
      '🥕',
      '🌽',
      '🍞',
      '🥖',
      '🧀',
      '🥩',
      '🍗',
      '☕',
      '🍺',
      '🍷',
      '🥤',
    ],
  },
  sports: {
    name: 'Спорт и отдых',
    icons: [
      '⚽',
      '🏀',
      '🏈',
      '⚾',
      '🎾',
      '🏐',
      '🏉',
      '🎱',
      '🏓',
      '🏸',
      '🥅',
      '⛳',
      '🏹',
      '🎣',
      '🥊',
      '🥋',
      '🎿',
      '⛷️',
      '🏂',
      '🏋️',
    ],
  },
  beauty: {
    name: 'Красота и здоровье',
    icons: [
      '💄',
      '💅',
      '💋',
      '👄',
      '👀',
      '👂',
      '👃',
      '🧴',
      '🧼',
      '🧽',
      '🪒',
      '💊',
      '🩹',
      '🩺',
      '💉',
      '🌡️',
      '🧬',
      '🔬',
      '⚗️',
      '💎',
    ],
  },
  books: {
    name: 'Книги и обучение',
    icons: [
      '📚',
      '📖',
      '📝',
      '📄',
      '📃',
      '📑',
      '📊',
      '📈',
      '📉',
      '🗂️',
      '📁',
      '📂',
      '🗃️',
      '🗄️',
      '📋',
      '📌',
      '📍',
      '📎',
      '🖇️',
      '📐',
    ],
  },
  nature: {
    name: 'Природа и животные',
    icons: [
      '🌱',
      '🌿',
      '🍀',
      '🌸',
      '🌺',
      '🌻',
      '🌷',
      '🌹',
      '🏵️',
      '💐',
      '🌳',
      '🌲',
      '🌴',
      '🐶',
      '🐱',
      '🐭',
      '🐹',
      '🐰',
      '🦊',
      '🐻',
    ],
  },
  tools: {
    name: 'Инструменты',
    icons: [
      '🔨',
      '🪓',
      '⛏️',
      '🔧',
      '🔩',
      '🪚',
      '🔗',
      '⛓️',
      '📎',
      '📏',
      '📐',
      '✂️',
      '📌',
      '📍',
      '🔍',
      '🔎',
      '💡',
      '🔦',
      '🕯️',
      '💰',
    ],
  },
  numbers: {
    name: 'Числа и символы',
    icons: [
      '🔢',
      '📊',
      '📈',
      '📉',
      '💹',
      '💰',
      '💵',
      '💴',
      '💶',
      '💷',
      '🪙',
      '💳',
      '🧮',
      '⚖️',
      '📏',
      '📐',
      '🔺',
      '🔻',
      '💯',
      '🎯',
    ],
  },
  attributes: {
    name: 'Атрибуты',
    icons: [
      '📝',
      '🔤',
      '🔢',
      '✅',
      '❌',
      '📅',
      '📍',
      '📁',
      '🖼️',
      '🎨',
      '🏷️',
      '⭐',
      '❤️',
      '🔥',
      '💎',
      '🎁',
      '🎈',
      '🎀',
      '🎊',
      '🎉',
    ],
  },
};

export default function IconPicker({
  value,
  onChange,
  placeholder = 'Выберите иконку',
}: IconPickerProps) {
  const [isOpen, setIsOpen] = useState(false);
  const [activeCategory, setActiveCategory] = useState('transport');

  const handleIconSelect = (icon: string) => {
    onChange(icon);
    setIsOpen(false);
  };

  const handleInputChange = (e: React.ChangeEvent<HTMLInputElement>) => {
    onChange(e.target.value);
  };

  return (
    <div className="form-control relative">
      <div className="flex gap-2">
        <input
          type="text"
          value={value}
          onChange={handleInputChange}
          className="input input-bordered flex-1"
          placeholder={placeholder}
        />
        <button
          type="button"
          onClick={() => setIsOpen(!isOpen)}
          className="btn btn-outline btn-square"
        >
          {value || '🎨'}
        </button>
      </div>

      {isOpen && (
        <div className="absolute z-50 mt-1 bg-base-100 border border-base-300 rounded-lg shadow-lg p-4 w-80 right-0">
          {/* Category tabs */}
          <div className="tabs tabs-boxed mb-4">
            <div className="flex flex-wrap gap-1">
              {Object.entries(iconCategories).map(([key, category]) => (
                <button
                  key={key}
                  type="button"
                  onClick={() => setActiveCategory(key)}
                  className={`tab tab-sm ${activeCategory === key ? 'tab-active' : ''}`}
                >
                  {category.name}
                </button>
              ))}
            </div>
          </div>

          {/* Icon grid */}
          <div className="grid grid-cols-8 gap-2 max-h-48 overflow-y-auto">
            {iconCategories[
              activeCategory as keyof typeof iconCategories
            ]?.icons.map((icon) => (
              <button
                key={icon}
                type="button"
                onClick={() => handleIconSelect(icon)}
                className={`btn btn-sm btn-ghost hover:btn-primary text-lg ${
                  value === icon ? 'btn-primary' : ''
                }`}
              >
                {icon}
              </button>
            ))}
          </div>

          {/* Close button */}
          <div className="flex justify-end mt-4">
            <button
              type="button"
              onClick={() => setIsOpen(false)}
              className="btn btn-sm btn-ghost"
            >
              Закрыть
            </button>
          </div>
        </div>
      )}
    </div>
  );
}
