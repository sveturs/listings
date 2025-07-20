'use client';

import { useState } from 'react';

interface CustomComponent {
  id: string;
  name: string;
  description: string;
  icon: string;
  preview?: React.ReactNode;
  compatibleTypes: string[];
}

interface CustomComponentSelectorProps {
  value?: string;
  onChange: (value: string) => void;
  attributeType?: string;
}

const customComponents: CustomComponent[] = [
  {
    id: 'color-picker',
    name: 'Выбор цвета',
    description: 'Палитра цветов с предустановленными вариантами',
    icon: '🎨',
    compatibleTypes: ['text', 'select'],
    preview: (
      <div className="flex gap-2">
        <div className="w-8 h-8 rounded bg-red-500"></div>
        <div className="w-8 h-8 rounded bg-blue-500"></div>
        <div className="w-8 h-8 rounded bg-green-500"></div>
        <div className="w-8 h-8 rounded bg-yellow-500"></div>
      </div>
    ),
  },
  {
    id: 'date-range-picker',
    name: 'Выбор диапазона дат',
    description: 'Календарь для выбора начальной и конечной даты',
    icon: '📅',
    compatibleTypes: ['date', 'text'],
    preview: (
      <div className="text-sm bg-base-200 px-3 py-2 rounded">
        01.01.2025 - 31.01.2025
      </div>
    ),
  },
  {
    id: 'file-upload',
    name: 'Загрузка файлов',
    description: 'Загрузка документов с предпросмотром',
    icon: '📄',
    compatibleTypes: ['file'],
    preview: (
      <div className="flex items-center gap-2 text-sm">
        <span className="text-2xl">📎</span>
        <span>document.pdf (2.5 MB)</span>
      </div>
    ),
  },
  {
    id: 'location-picker',
    name: 'Выбор местоположения',
    description: 'Интерактивная карта для выбора точки',
    icon: '📍',
    compatibleTypes: ['location', 'text'],
    preview: (
      <div className="bg-base-200 rounded p-4 text-center text-sm">
        <div className="text-2xl mb-1">🗺️</div>
        <div>55.7558° N, 37.6173° E</div>
      </div>
    ),
  },
  {
    id: 'image-gallery',
    name: 'Галерея изображений',
    description: 'Загрузка и организация нескольких изображений',
    icon: '🖼️',
    compatibleTypes: ['gallery', 'file'],
    preview: (
      <div className="flex gap-2">
        <div className="w-16 h-16 bg-base-200 rounded flex items-center justify-center">
          <span className="text-2xl">🖼️</span>
        </div>
        <div className="w-16 h-16 bg-base-200 rounded flex items-center justify-center">
          <span className="text-2xl">🖼️</span>
        </div>
        <div className="w-16 h-16 bg-base-200 rounded flex items-center justify-center">
          <span className="text-xl">+3</span>
        </div>
      </div>
    ),
  },
  {
    id: 'rating-stars',
    name: 'Звездный рейтинг',
    description: 'Выбор рейтинга от 1 до 5 звезд',
    icon: '⭐',
    compatibleTypes: ['number', 'range'],
    preview: (
      <div className="flex gap-1 text-xl">
        <span>⭐</span>
        <span>⭐</span>
        <span>⭐</span>
        <span>⭐</span>
        <span className="opacity-30">⭐</span>
      </div>
    ),
  },
  {
    id: 'tags-input',
    name: 'Ввод тегов',
    description: 'Множественный выбор с автодополнением',
    icon: '🏷️',
    compatibleTypes: ['multiselect', 'text'],
    preview: (
      <div className="flex gap-2 flex-wrap">
        <span className="badge badge-primary">React</span>
        <span className="badge badge-primary">TypeScript</span>
        <span className="badge badge-primary">Next.js</span>
      </div>
    ),
  },
  {
    id: 'slider-range',
    name: 'Слайдер диапазона',
    description: 'Выбор значения или диапазона с помощью слайдера',
    icon: '🎚️',
    compatibleTypes: ['range', 'number'],
    preview: (
      <div className="w-full">
        <input
          type="range"
          min="0"
          max="100"
          value="40"
          className="range range-primary range-sm"
          readOnly
        />
      </div>
    ),
  },
  {
    id: 'time-picker',
    name: 'Выбор времени',
    description: 'Выбор времени с точностью до минут',
    icon: '⏰',
    compatibleTypes: ['text', 'date'],
    preview: <div className="text-sm bg-base-200 px-3 py-2 rounded">14:30</div>,
  },
  {
    id: 'rich-text-editor',
    name: 'Текстовый редактор',
    description: 'Форматированный текст с поддержкой стилей',
    icon: '📝',
    compatibleTypes: ['text'],
    preview: (
      <div className="text-sm space-y-1">
        <div className="font-bold">Заголовок</div>
        <div>
          Обычный текст с <span className="font-bold">жирным</span> и{' '}
          <span className="italic">курсивом</span>
        </div>
      </div>
    ),
  },
];

export default function CustomComponentSelector({
  value,
  onChange,
  attributeType = 'text',
}: CustomComponentSelectorProps) {
  const [isOpen, setIsOpen] = useState(false);
  const [searchQuery, setSearchQuery] = useState('');

  const selectedComponent = customComponents.find((c) => c.id === value);

  const filteredComponents = customComponents.filter((component) => {
    const matchesSearch =
      component.name.toLowerCase().includes(searchQuery.toLowerCase()) ||
      component.description.toLowerCase().includes(searchQuery.toLowerCase());

    const isCompatible =
      !attributeType || component.compatibleTypes.includes(attributeType);

    return matchesSearch && isCompatible;
  });

  const handleSelect = (componentId: string) => {
    onChange(componentId);
    setIsOpen(false);
    setSearchQuery('');
  };

  const handleClear = () => {
    onChange('');
    setIsOpen(false);
  };

  return (
    <div className="relative">
      <div className="form-control">
        <label className="label">
          <span className="label-text">Кастомный компонент</span>
        </label>

        <div
          className="input input-bordered flex items-center justify-between cursor-pointer"
          onClick={() => setIsOpen(!isOpen)}
        >
          {selectedComponent ? (
            <div className="flex items-center gap-2">
              <span className="text-xl">{selectedComponent.icon}</span>
              <span>{selectedComponent.name}</span>
            </div>
          ) : (
            <span className="text-base-content/60">Выберите компонент</span>
          )}

          <div className="flex items-center gap-2">
            {selectedComponent && (
              <button
                type="button"
                onClick={(e) => {
                  e.stopPropagation();
                  handleClear();
                }}
                className="btn btn-ghost btn-xs btn-circle"
              >
                ✕
              </button>
            )}
            <svg
              className={`w-4 h-4 transition-transform ${isOpen ? 'rotate-180' : ''}`}
              fill="none"
              viewBox="0 0 24 24"
              stroke="currentColor"
            >
              <path
                strokeLinecap="round"
                strokeLinejoin="round"
                strokeWidth={2}
                d="M19 9l-7 7-7-7"
              />
            </svg>
          </div>
        </div>

        {selectedComponent && (
          <div className="text-sm text-base-content/60 mt-1">
            {selectedComponent.description}
          </div>
        )}
      </div>

      {isOpen && (
        <div className="absolute z-50 mt-2 w-full bg-base-100 rounded-lg shadow-lg border border-base-300 max-h-96 overflow-hidden">
          <div className="p-3 border-b border-base-300">
            <input
              type="text"
              placeholder="Поиск компонентов..."
              className="input input-bordered input-sm w-full"
              value={searchQuery}
              onChange={(e) => setSearchQuery(e.target.value)}
              onClick={(e) => e.stopPropagation()}
            />
          </div>

          <div className="overflow-y-auto max-h-80">
            {filteredComponents.length === 0 ? (
              <div className="p-4 text-center text-base-content/60">
                Нет подходящих компонентов для типа &quot;{attributeType}&quot;
              </div>
            ) : (
              filteredComponents.map((component) => (
                <div
                  key={component.id}
                  className="p-4 hover:bg-base-200 cursor-pointer border-b border-base-300 last:border-b-0"
                  onClick={() => handleSelect(component.id)}
                >
                  <div className="flex items-start gap-3">
                    <span className="text-2xl">{component.icon}</span>
                    <div className="flex-1">
                      <div className="font-medium">{component.name}</div>
                      <div className="text-sm text-base-content/60 mt-1">
                        {component.description}
                      </div>
                      {component.preview && (
                        <div className="mt-3 p-3 bg-base-200 rounded">
                          <div className="text-xs text-base-content/60 mb-2">
                            Предпросмотр:
                          </div>
                          {component.preview}
                        </div>
                      )}
                      <div className="mt-2 flex flex-wrap gap-1">
                        {component.compatibleTypes.map((type) => (
                          <span
                            key={type}
                            className="badge badge-ghost badge-sm"
                          >
                            {type}
                          </span>
                        ))}
                      </div>
                    </div>
                  </div>
                </div>
              ))
            )}
          </div>
        </div>
      )}
    </div>
  );
}
