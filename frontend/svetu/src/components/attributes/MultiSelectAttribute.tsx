'use client';

import { useState, useEffect, useRef } from 'react';
import { useTranslations } from 'next-intl';

interface MultiSelectAttributeProps {
  attribute: {
    id: number;
    name: string;
    display_name: string;
    options?: string | string[] | any[] | { values: string[] };
    translations?: Record<string, string>;
    option_translations?: Record<string, Record<string, string>>;
    is_required?: boolean;
  };
  value?: string[] | string;
  onChange: (value: string[]) => void;
  error?: string;
  locale?: string;
}

export default function MultiSelectAttribute({
  attribute,
  value = [],
  onChange,
  error,
  locale = 'ru',
}: MultiSelectAttributeProps) {
  const t = useTranslations('marketplace');
  const [selectedValues, setSelectedValues] = useState<string[]>([]);
  const [isOpen, setIsOpen] = useState(false);
  const dropdownRef = useRef<HTMLDivElement>(null);

  // Close dropdown when clicking outside
  useEffect(() => {
    const handleClickOutside = (event: MouseEvent) => {
      if (
        dropdownRef.current &&
        !dropdownRef.current.contains(event.target as Node)
      ) {
        setIsOpen(false);
      }
    };

    if (isOpen) {
      document.addEventListener('mousedown', handleClickOutside);
    }

    return () => {
      document.removeEventListener('mousedown', handleClickOutside);
    };
  }, [isOpen]);

  useEffect(() => {
    let newValues: string[] = [];

    if (value) {
      if (typeof value === 'string') {
        try {
          const parsed = JSON.parse(value);
          newValues = Array.isArray(parsed) ? parsed : [value];
        } catch {
          newValues = [value];
        }
      } else if (Array.isArray(value)) {
        newValues = value;
      }
    }

    // Only update if values actually changed (prevent infinite loop)
    const currentValuesStr = JSON.stringify(selectedValues.sort());
    const newValuesStr = JSON.stringify(newValues.sort());

    if (currentValuesStr !== newValuesStr) {
      setSelectedValues(newValues);
    }
  }, [value, selectedValues]);

  // Parse options
  let options: any[] = [];
  if (attribute.options) {
    if (typeof attribute.options === 'string') {
      try {
        const parsed = JSON.parse(attribute.options);
        options = Array.isArray(parsed) ? parsed : [];
      } catch {
        options = [];
      }
    } else if (Array.isArray(attribute.options)) {
      options = attribute.options;
    } else if (
      typeof attribute.options === 'object' &&
      attribute.options.values
    ) {
      // Handle options in format { values: [...] }
      options = Array.isArray(attribute.options.values)
        ? attribute.options.values
        : [];
    }
  }

  const toggleOption = (optionValue: string) => {
    const newValues = selectedValues.includes(optionValue)
      ? selectedValues.filter((v) => v !== optionValue)
      : [...selectedValues, optionValue];

    setSelectedValues(newValues);
    onChange(newValues);
  };

  const getOptionLabel = (option: any) => {
    const optionValue =
      typeof option === 'string' ? option : option.value || option;

    // Check for translations
    if (attribute.option_translations?.[optionValue]?.[locale]) {
      return attribute.option_translations[optionValue][locale];
    }

    // Return label or value
    return typeof option === 'object' ? option.label || option.value : option;
  };

  const getDisplayName = () => {
    if (attribute.translations?.[locale]) {
      return attribute.translations[locale];
    }
    return attribute.display_name;
  };

  const selectedCount = selectedValues.length;

  return (
    <div className="form-control">
      <label className="label">
        <span className="label-text">
          {getDisplayName()}
          {attribute.is_required && <span className="text-error"> *</span>}
        </span>
      </label>

      <div className="relative" ref={dropdownRef}>
        <button
          type="button"
          className={`btn btn-outline w-full justify-between ${error ? 'btn-error' : ''}`}
          onClick={() => setIsOpen(!isOpen)}
        >
          <span className="text-left truncate">
            {selectedCount > 0
              ? `${t('create.selected')}: ${selectedCount}`
              : t('create.selectOptions')}
          </span>
          <svg
            className={`w-4 h-4 ml-2 transition-transform ${isOpen ? 'rotate-180' : ''}`}
            fill="none"
            stroke="currentColor"
            viewBox="0 0 24 24"
          >
            <path
              strokeLinecap="round"
              strokeLinejoin="round"
              strokeWidth={2}
              d="M19 9l-7 7-7-7"
            />
          </svg>
        </button>

        {isOpen && (
          <div className="absolute z-10 w-full mt-1">
            <ul className="menu p-2 shadow-lg bg-base-100 rounded-box w-full max-h-60 overflow-auto border border-base-300">
              {options.length === 0 ? (
                <li className="text-base-content/50 p-2">
                  {t('create.noOptions')}
                </li>
              ) : (
                options.map((option, index) => {
                  const optionValue =
                    typeof option === 'string'
                      ? option
                      : option.value || option;
                  const isSelected = selectedValues.includes(optionValue);

                  return (
                    <li key={index}>
                      <label className="cursor-pointer flex items-center gap-2 p-2">
                        <input
                          type="checkbox"
                          className="checkbox checkbox-sm checkbox-primary"
                          checked={isSelected}
                          onChange={() => toggleOption(optionValue)}
                        />
                        <span className="flex-1">{getOptionLabel(option)}</span>
                      </label>
                    </li>
                  );
                })
              )}
            </ul>
          </div>
        )}
      </div>

      {selectedCount > 0 && (
        <div className="mt-2 flex flex-wrap gap-1">
          {selectedValues.map((val) => {
            const option = options.find((opt) =>
              typeof opt === 'string' ? opt === val : opt.value === val
            );
            return (
              <div key={val} className="badge badge-primary gap-1">
                <span>{option ? getOptionLabel(option) : val}</span>
                <button
                  type="button"
                  onClick={() => toggleOption(val)}
                  className="btn btn-ghost btn-xs p-0 h-4 w-4"
                >
                  ✕
                </button>
              </div>
            );
          })}
        </div>
      )}

      {error && (
        <label className="label">
          <span className="label-text-alt text-error">{error}</span>
        </label>
      )}
    </div>
  );
}
