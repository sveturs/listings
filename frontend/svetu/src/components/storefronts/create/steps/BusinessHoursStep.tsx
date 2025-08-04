'use client';

import { useTranslations } from 'next-intl';
import { useCreateStorefrontContext } from '@/contexts/CreateStorefrontContext';

interface BusinessHoursStepProps {
  onNext: () => void;
  onBack: () => void;
}

const dayNames = [
  'sunday',
  'monday',
  'tuesday',
  'wednesday',
  'thursday',
  'friday',
  'saturday',
];

export default function BusinessHoursStep({
  onNext,
  onBack,
}: BusinessHoursStepProps) {
  const tCommon = useTranslations('common');
  const t = useTranslations('create_storefront');
  const { formData, updateFormData } = useCreateStorefrontContext();

  const handleHoursChange = (dayIndex: number, field: string, value: any) => {
    const newHours = [...formData.businessHours];
    newHours[dayIndex] = { ...newHours[dayIndex], [field]: value };
    updateFormData({ businessHours: newHours });
  };

  const handleNext = () => {
    onNext();
  };

  return (
    <div className="max-w-2xl mx-auto">
      <div className="card bg-base-100 shadow-xl">
        <div className="card-body">
          <h2 className="card-title text-2xl mb-4">
            {t('business_hours.title')}
          </h2>
          <p className="text-base-content/70 mb-6">
            {t('business_hours.subtitle')}
          </p>

          <div className="space-y-4">
            {formData.businessHours &&
              Array.isArray(formData.businessHours) &&
              formData.businessHours.map((hours, index) => (
                <div key={index} className="flex items-center gap-4">
                  <div className="w-32">
                    <span className="font-medium">
                      {t(
                        `common.days.${dayNames[hours.dayOfWeek] || 'monday'}`
                      )}
                    </span>
                  </div>

                  <div className="form-control">
                    <label className="cursor-pointer label">
                      <input
                        type="checkbox"
                        className="toggle toggle-primary"
                        checked={!hours.isClosed}
                        onChange={(e) =>
                          handleHoursChange(
                            index,
                            'isClosed',
                            !e.target.checked
                          )
                        }
                      />
                      <span className="label-text ml-2">
                        {hours.isClosed ? t('closed') : t('open')}
                      </span>
                    </label>
                  </div>

                  {!hours.isClosed && (
                    <>
                      <input
                        type="time"
                        className="input input-bordered input-sm"
                        value={hours.openTime}
                        onChange={(e) =>
                          handleHoursChange(index, 'openTime', e.target.value)
                        }
                      />
                      <span>-</span>
                      <input
                        type="time"
                        className="input input-bordered input-sm"
                        value={hours.closeTime}
                        onChange={(e) =>
                          handleHoursChange(index, 'closeTime', e.target.value)
                        }
                      />
                    </>
                  )}
                </div>
              ))}
          </div>

          <div className="divider"></div>

          <div className="alert alert-info">
            <svg
              xmlns="http://www.w3.org/2000/svg"
              fill="none"
              viewBox="0 0 24 24"
              className="stroke-current shrink-0 w-6 h-6"
            >
              <path
                strokeLinecap="round"
                strokeLinejoin="round"
                strokeWidth="2"
                d="M13 16h-1v-4h-1m1-4h.01M21 12a9 9 0 11-18 0 9 9 0 0118 0z"
              ></path>
            </svg>
            <span>{t('business_hours.tip')}</span>
          </div>

          <div className="card-actions justify-between mt-6">
            <button className="btn btn-ghost" onClick={onBack}>
              {t('back')}
            </button>
            <button className="btn btn-primary" onClick={handleNext}>
              {t('next')}
            </button>
          </div>
        </div>
      </div>
    </div>
  );
}
