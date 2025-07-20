'use client';

import { useState, useEffect } from 'react';
import { useTranslations } from 'next-intl';

interface ValidationRule {
  type: string;
  value?: string | number | boolean;
  message?: string;
}

interface ValidationRulesEditorProps {
  value?: Record<string, any>;
  onChange: (rules: Record<string, any>) => void;
  attributeType?: string;
}

const PRESET_RULES = {
  email: {
    type: 'pattern',
    value: '^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\\.[a-zA-Z]{2,}$',
    message: 'validationRules.messages.invalidEmail',
  },
  phone: {
    type: 'pattern',
    value: '^\\+?[1-9]\\d{1,14}$',
    message: 'validationRules.messages.invalidPhone',
  },
  url: {
    type: 'pattern',
    value:
      '^(https?:\\/\\/)?(www\\.)?[-a-zA-Z0-9@:%._\\+~#=]{1,256}\\.[a-zA-Z0-9()]{1,6}\\b([-a-zA-Z0-9()@:%_\\+.~#?&//=]*)$',
    message: 'validationRules.messages.invalidUrl',
  },
  alphanumeric: {
    type: 'pattern',
    value: '^[a-zA-Z0-9]+$',
    message: 'validationRules.messages.alphanumericOnly',
  },
  letters: {
    type: 'pattern',
    value: '^[a-zA-Z]+$',
    message: 'validationRules.messages.lettersOnly',
  },
  numbers: {
    type: 'pattern',
    value: '^[0-9]+$',
    message: 'validationRules.messages.numbersOnly',
  },
};

export default function ValidationRulesEditor({
  value = {},
  onChange,
  attributeType = 'text',
}: ValidationRulesEditorProps) {
  const t = useTranslations('admin.attributes');
  const [rules, setRules] = useState<Record<string, ValidationRule>>({});
  const [customPattern, setCustomPattern] = useState('');
  const [customMessage, setCustomMessage] = useState('');
  const [selectedPreset, setSelectedPreset] = useState('');

  useEffect(() => {
    if (value && typeof value === 'object') {
      const parsedRules: Record<string, ValidationRule> = {};

      Object.entries(value).forEach(([key, val]) => {
        if (key === 'pattern' && typeof val === 'string') {
          parsedRules.pattern = {
            type: 'pattern',
            value: val,
            message: value.patternMessage || '',
          };
        } else if (key === 'min' || key === 'max') {
          parsedRules[key] = {
            type: key,
            value: val as number,
            message: value[`${key}Message`] || '',
          };
        } else if (key === 'required' && val === true) {
          parsedRules.required = {
            type: 'required',
            value: true,
            message: value.requiredMessage || '',
          };
        } else if (key === 'minLength' || key === 'maxLength') {
          parsedRules[key] = {
            type: key,
            value: val as number,
            message: value[`${key}Message`] || '',
          };
        }
      });

      setRules(parsedRules);

      if (parsedRules.pattern?.value) {
        setCustomPattern(parsedRules.pattern.value as string);
        setCustomMessage(parsedRules.pattern.message || '');
      }
    }
  }, [value]);

  const updateRules = (newRules: Record<string, ValidationRule>) => {
    setRules(newRules);

    const output: Record<string, any> = {};

    Object.entries(newRules).forEach(([key, rule]) => {
      if (rule.value !== undefined && rule.value !== '') {
        output[key] = rule.value;
        if (rule.message) {
          output[`${key}Message`] = rule.message;
        }
      }
    });

    onChange(output);
  };

  const addRule = (type: string) => {
    const newRules = { ...rules };

    switch (type) {
      case 'required':
        newRules.required = {
          type: 'required',
          value: true,
          message: '',
        };
        break;
      case 'minLength':
      case 'maxLength':
        newRules[type] = {
          type,
          value: type === 'minLength' ? 1 : 100,
          message: '',
        };
        break;
      case 'min':
      case 'max':
        newRules[type] = {
          type,
          value: type === 'min' ? 0 : 100,
          message: '',
        };
        break;
      case 'pattern':
        newRules.pattern = {
          type: 'pattern',
          value: customPattern,
          message: customMessage,
        };
        break;
    }

    updateRules(newRules);
  };

  const removeRule = (type: string) => {
    const newRules = { ...rules };
    delete newRules[type];
    updateRules(newRules);
  };

  const updateRuleValue = (
    type: string,
    field: 'value' | 'message',
    value: any
  ) => {
    const newRules = { ...rules };
    if (newRules[type]) {
      newRules[type] = { ...newRules[type], [field]: value };
      updateRules(newRules);
    }
  };

  const applyPreset = (presetKey: string) => {
    const preset = PRESET_RULES[presetKey as keyof typeof PRESET_RULES];
    if (preset) {
      setCustomPattern(preset.value);
      setCustomMessage(preset.message);
      setSelectedPreset(presetKey);

      const newRules = { ...rules };
      newRules.pattern = {
        type: 'pattern',
        value: preset.value,
        message: preset.message,
      };
      updateRules(newRules);
    }
  };

  const availableRules = () => {
    const rules: Array<{ type: string; label: string }> = [];

    if (!rules.find((r) => r.type === 'required')) {
      rules.push({ type: 'required', label: t('validationRules.required') });
    }

    if (attributeType === 'text') {
      if (!rules.find((r) => r.type === 'minLength')) {
        rules.push({
          type: 'minLength',
          label: t('validationRules.minLength'),
        });
      }
      if (!rules.find((r) => r.type === 'maxLength')) {
        rules.push({
          type: 'maxLength',
          label: t('validationRules.maxLength'),
        });
      }
      if (!rules.find((r) => r.type === 'pattern')) {
        rules.push({ type: 'pattern', label: t('validationRules.pattern') });
      }
    }

    if (attributeType === 'number' || attributeType === 'range') {
      if (!rules.find((r) => r.type === 'min')) {
        rules.push({ type: 'min', label: t('validationRules.min') });
      }
      if (!rules.find((r) => r.type === 'max')) {
        rules.push({ type: 'max', label: t('validationRules.max') });
      }
    }

    return rules;
  };

  return (
    <div className="space-y-4">
      <div className="flex items-center justify-between">
        <h3 className="text-lg font-medium">{t('validationRules.title')}</h3>
        <div className="dropdown dropdown-end">
          <label tabIndex={0} className="btn btn-sm btn-outline">
            {t('validationRules.addRule')}
          </label>
          <ul
            tabIndex={0}
            className="dropdown-content menu p-2 shadow bg-base-100 rounded-box w-52 z-50"
          >
            {availableRules().map((rule) => (
              <li key={rule.type}>
                <a onClick={() => addRule(rule.type)}>{rule.label}</a>
              </li>
            ))}
          </ul>
        </div>
      </div>

      <div className="space-y-3">
        {/* Required Rule */}
        {rules.required && (
          <div className="card bg-base-200">
            <div className="card-body p-4">
              <div className="flex items-center justify-between">
                <h4 className="font-medium">{t('validationRules.required')}</h4>
                <button
                  type="button"
                  onClick={() => removeRule('required')}
                  className="btn btn-ghost btn-xs"
                >
                  ✕
                </button>
              </div>
              <div className="form-control mt-2">
                <label className="label">
                  <span className="label-text text-sm">
                    {t('validationRules.errorMessage')}
                  </span>
                </label>
                <input
                  type="text"
                  value={rules.required.message || ''}
                  onChange={(e) =>
                    updateRuleValue('required', 'message', e.target.value)
                  }
                  className="input input-sm input-bordered"
                  placeholder={t('validationRules.defaultMessages.required')}
                />
              </div>
            </div>
          </div>
        )}

        {/* Min/Max Length Rules */}
        {(rules.minLength || rules.maxLength) && (
          <div className="card bg-base-200">
            <div className="card-body p-4">
              <div className="flex items-center justify-between">
                <h4 className="font-medium">
                  {t('validationRules.lengthConstraints')}
                </h4>
                <button
                  type="button"
                  onClick={() => {
                    removeRule('minLength');
                    removeRule('maxLength');
                  }}
                  className="btn btn-ghost btn-xs"
                >
                  ✕
                </button>
              </div>
              <div className="grid grid-cols-2 gap-2 mt-2">
                {rules.minLength && (
                  <div>
                    <label className="label">
                      <span className="label-text text-sm">
                        {t('validationRules.minLength')}
                      </span>
                    </label>
                    <input
                      type="number"
                      value={
                        typeof rules.minLength.value === 'number'
                          ? rules.minLength.value
                          : ''
                      }
                      onChange={(e) =>
                        updateRuleValue(
                          'minLength',
                          'value',
                          parseInt(e.target.value)
                        )
                      }
                      className="input input-sm input-bordered w-full"
                      min="0"
                    />
                  </div>
                )}
                {rules.maxLength && (
                  <div>
                    <label className="label">
                      <span className="label-text text-sm">
                        {t('validationRules.maxLength')}
                      </span>
                    </label>
                    <input
                      type="number"
                      value={
                        typeof rules.maxLength.value === 'number'
                          ? rules.maxLength.value
                          : ''
                      }
                      onChange={(e) =>
                        updateRuleValue(
                          'maxLength',
                          'value',
                          parseInt(e.target.value)
                        )
                      }
                      className="input input-sm input-bordered w-full"
                      min="0"
                    />
                  </div>
                )}
              </div>
              {rules.minLength && (
                <div className="form-control mt-2">
                  <label className="label">
                    <span className="label-text text-sm">
                      {t('validationRules.minLengthMessage')}
                    </span>
                  </label>
                  <input
                    type="text"
                    value={rules.minLength.message || ''}
                    onChange={(e) =>
                      updateRuleValue('minLength', 'message', e.target.value)
                    }
                    className="input input-sm input-bordered"
                    placeholder={t('validationRules.defaultMessages.minLength')}
                  />
                </div>
              )}
              {rules.maxLength && (
                <div className="form-control mt-2">
                  <label className="label">
                    <span className="label-text text-sm">
                      {t('validationRules.maxLengthMessage')}
                    </span>
                  </label>
                  <input
                    type="text"
                    value={rules.maxLength.message || ''}
                    onChange={(e) =>
                      updateRuleValue('maxLength', 'message', e.target.value)
                    }
                    className="input input-sm input-bordered"
                    placeholder={t('validationRules.defaultMessages.maxLength')}
                  />
                </div>
              )}
            </div>
          </div>
        )}

        {/* Min/Max Value Rules */}
        {(rules.min || rules.max) && (
          <div className="card bg-base-200">
            <div className="card-body p-4">
              <div className="flex items-center justify-between">
                <h4 className="font-medium">
                  {t('validationRules.valueConstraints')}
                </h4>
                <button
                  type="button"
                  onClick={() => {
                    removeRule('min');
                    removeRule('max');
                  }}
                  className="btn btn-ghost btn-xs"
                >
                  ✕
                </button>
              </div>
              <div className="grid grid-cols-2 gap-2 mt-2">
                {rules.min && (
                  <div>
                    <label className="label">
                      <span className="label-text text-sm">
                        {t('validationRules.min')}
                      </span>
                    </label>
                    <input
                      type="number"
                      value={
                        typeof rules.min.value === 'number'
                          ? rules.min.value
                          : ''
                      }
                      onChange={(e) =>
                        updateRuleValue(
                          'min',
                          'value',
                          parseFloat(e.target.value)
                        )
                      }
                      className="input input-sm input-bordered w-full"
                    />
                  </div>
                )}
                {rules.max && (
                  <div>
                    <label className="label">
                      <span className="label-text text-sm">
                        {t('validationRules.max')}
                      </span>
                    </label>
                    <input
                      type="number"
                      value={
                        typeof rules.max.value === 'number'
                          ? rules.max.value
                          : ''
                      }
                      onChange={(e) =>
                        updateRuleValue(
                          'max',
                          'value',
                          parseFloat(e.target.value)
                        )
                      }
                      className="input input-sm input-bordered w-full"
                    />
                  </div>
                )}
              </div>
              {rules.min && (
                <div className="form-control mt-2">
                  <label className="label">
                    <span className="label-text text-sm">
                      {t('validationRules.minMessage')}
                    </span>
                  </label>
                  <input
                    type="text"
                    value={rules.min.message || ''}
                    onChange={(e) =>
                      updateRuleValue('min', 'message', e.target.value)
                    }
                    className="input input-sm input-bordered"
                    placeholder={t('validationRules.defaultMessages.min')}
                  />
                </div>
              )}
              {rules.max && (
                <div className="form-control mt-2">
                  <label className="label">
                    <span className="label-text text-sm">
                      {t('validationRules.maxMessage')}
                    </span>
                  </label>
                  <input
                    type="text"
                    value={rules.max.message || ''}
                    onChange={(e) =>
                      updateRuleValue('max', 'message', e.target.value)
                    }
                    className="input input-sm input-bordered"
                    placeholder={t('validationRules.defaultMessages.max')}
                  />
                </div>
              )}
            </div>
          </div>
        )}

        {/* Pattern Rule */}
        {rules.pattern && (
          <div className="card bg-base-200">
            <div className="card-body p-4">
              <div className="flex items-center justify-between">
                <h4 className="font-medium">{t('validationRules.pattern')}</h4>
                <button
                  type="button"
                  onClick={() => removeRule('pattern')}
                  className="btn btn-ghost btn-xs"
                >
                  ✕
                </button>
              </div>

              <div className="form-control mt-2">
                <label className="label">
                  <span className="label-text text-sm">
                    {t('validationRules.presets')}
                  </span>
                </label>
                <select
                  value={selectedPreset}
                  onChange={(e) => applyPreset(e.target.value)}
                  className="select select-sm select-bordered"
                >
                  <option value="">{t('validationRules.selectPreset')}</option>
                  <option value="email">
                    {t('validationRules.presetTypes.email')}
                  </option>
                  <option value="phone">
                    {t('validationRules.presetTypes.phone')}
                  </option>
                  <option value="url">
                    {t('validationRules.presetTypes.url')}
                  </option>
                  <option value="alphanumeric">
                    {t('validationRules.presetTypes.alphanumeric')}
                  </option>
                  <option value="letters">
                    {t('validationRules.presetTypes.letters')}
                  </option>
                  <option value="numbers">
                    {t('validationRules.presetTypes.numbers')}
                  </option>
                </select>
              </div>

              <div className="form-control mt-2">
                <label className="label">
                  <span className="label-text text-sm">
                    {t('validationRules.regularExpression')}
                  </span>
                </label>
                <input
                  type="text"
                  value={customPattern}
                  onChange={(e) => {
                    setCustomPattern(e.target.value);
                    updateRuleValue('pattern', 'value', e.target.value);
                    setSelectedPreset('');
                  }}
                  className="input input-sm input-bordered font-mono"
                  placeholder="^[A-Z][0-9]+$"
                />
              </div>

              <div className="form-control mt-2">
                <label className="label">
                  <span className="label-text text-sm">
                    {t('validationRules.errorMessage')}
                  </span>
                </label>
                <input
                  type="text"
                  value={customMessage}
                  onChange={(e) => {
                    setCustomMessage(e.target.value);
                    updateRuleValue('pattern', 'message', e.target.value);
                  }}
                  className="input input-sm input-bordered"
                  placeholder={t('validationRules.defaultMessages.pattern')}
                />
              </div>
            </div>
          </div>
        )}
      </div>

      {Object.keys(rules).length === 0 && (
        <div className="alert">
          <svg
            xmlns="http://www.w3.org/2000/svg"
            fill="none"
            viewBox="0 0 24 24"
            className="stroke-info shrink-0 w-6 h-6"
          >
            <path
              strokeLinecap="round"
              strokeLinejoin="round"
              strokeWidth="2"
              d="M13 16h-1v-4h-1m1-4h.01M21 12a9 9 0 11-18 0 9 9 0 0118 0z"
            ></path>
          </svg>
          <span>{t('validationRules.noRules')}</span>
        </div>
      )}
    </div>
  );
}
