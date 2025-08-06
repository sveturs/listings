'use client';

import { useState } from 'react';
import { useTranslations } from 'next-intl';
import { useCreateStorefrontContext } from '@/contexts/CreateStorefrontContext';

interface StaffSetupStepProps {
  onNext: () => void;
  onBack: () => void;
}

interface StaffMember {
  email: string;
  role: string;
  canManageProducts: boolean;
  canManageOrders: boolean;
  canManageSettings: boolean;
}

export default function StaffSetupStep({
  onNext,
  onBack,
}: StaffSetupStepProps) {
  const t = useTranslations('create_storefront');
  const _tPermissions = useTranslations('permissions');
  const _tRoles = useTranslations('roles');
  const tCommon = useTranslations('common');
  const { formData, updateFormData } = useCreateStorefrontContext();
  const [newStaffMember, setNewStaffMember] = useState<StaffMember>({
    email: '',
    role: 'staff',
    canManageProducts: true,
    canManageOrders: true,
    canManageSettings: false,
  });
  const [errors, setErrors] = useState<Record<string, string>>({});

  const validateEmail = (email: string) => {
    return /^[^\s@]+@[^\s@]+\.[^\s@]+$/.test(email);
  };

  const addStaffMember = () => {
    const newErrors: Record<string, string> = {};

    if (!newStaffMember.email) {
      newErrors.email = t('errors.email_required');
    } else if (!validateEmail(newStaffMember.email)) {
      newErrors.email = t('errors.email_invalid');
    } else if (formData.staff?.some((s) => s.email === newStaffMember.email)) {
      newErrors.email = t('errors.email_duplicate');
    }

    setErrors(newErrors);

    if (Object.keys(newErrors).length === 0) {
      const currentStaff = formData.staff || [];
      updateFormData({
        staff: [...currentStaff, { ...newStaffMember }],
      });
      setNewStaffMember({
        email: '',
        role: 'staff',
        canManageProducts: true,
        canManageOrders: true,
        canManageSettings: false,
      });
    }
  };

  const removeStaffMember = (index: number) => {
    const currentStaff = formData.staff || [];
    updateFormData({
      staff: currentStaff.filter((_, i) => i !== index),
    });
  };

  const handleNext = () => {
    onNext();
  };

  return (
    <div className="max-w-3xl mx-auto">
      <div className="card bg-base-100 shadow-xl">
        <div className="card-body">
          <h2 className="card-title text-2xl mb-4">{t('title')}</h2>
          <p className="text-base-content/70 mb-6">{t('subtitle')}</p>

          {/* Existing staff members */}
          {formData.staff && formData.staff.length > 0 && (
            <div className="mb-6">
              <h3 className="text-lg font-semibold mb-3">
                {t('current_staff')}
              </h3>
              <div className="space-y-2">
                {formData.staff.map((member, index) => (
                  <div
                    key={index}
                    className="flex items-center justify-between p-4 bg-base-200 rounded-lg"
                  >
                    <div>
                      <div className="font-medium">{member.email}</div>
                      <div className="text-sm text-base-content/70">
                        {t(`roles.${member.role}`)}
                      </div>
                      <div className="flex gap-4 mt-1 text-xs">
                        {member.canManageProducts && (
                          <span className="badge badge-sm">
                            {t('products')}
                          </span>
                        )}
                        {member.canManageOrders && (
                          <span className="badge badge-sm">{t('orders')}</span>
                        )}
                        {member.canManageSettings && (
                          <span className="badge badge-sm">
                            {t('settings')}
                          </span>
                        )}
                      </div>
                    </div>
                    <button
                      className="btn btn-ghost btn-sm"
                      onClick={() => removeStaffMember(index)}
                    >
                      ✕
                    </button>
                  </div>
                ))}
              </div>
            </div>
          )}

          {/* Add new staff member */}
          <div className="card bg-base-200">
            <div className="card-body p-4">
              <h3 className="text-lg font-semibold mb-3">{t('add_member')}</h3>

              <div className="grid grid-cols-1 md:grid-cols-2 gap-4">
                <div className="form-control">
                  <label className="label">
                    <span className="label-text">{t('email')}</span>
                  </label>
                  <input
                    type="email"
                    className={`input input-bordered ${errors.email ? 'input-error' : ''}`}
                    placeholder={t(
                      'create_storefront.staff_setup.email_placeholder'
                    )}
                    value={newStaffMember.email}
                    onChange={(e) => {
                      setNewStaffMember({
                        ...newStaffMember,
                        email: e.target.value,
                      });
                      setErrors({});
                    }}
                  />
                  {errors.email && (
                    <label className="label">
                      <span className="label-text-alt text-error">
                        {errors.email}
                      </span>
                    </label>
                  )}
                </div>

                <div className="form-control">
                  <label className="label">
                    <span className="label-text">{t('role')}</span>
                  </label>
                  <select
                    className="select select-bordered"
                    value={newStaffMember.role}
                    onChange={(e) =>
                      setNewStaffMember({
                        ...newStaffMember,
                        role: e.target.value,
                      })
                    }
                  >
                    <option value="staff">{t('staff')}</option>
                    <option value="manager">{t('manager')}</option>
                    <option value="admin">{t('admin')}</option>
                  </select>
                </div>
              </div>

              <div className="mt-4">
                <h4 className="font-medium mb-2">{t('permissions')}</h4>
                <div className="space-y-2">
                  <label className="cursor-pointer label justify-start">
                    <input
                      type="checkbox"
                      className="checkbox checkbox-primary"
                      checked={newStaffMember.canManageProducts}
                      onChange={(e) =>
                        setNewStaffMember({
                          ...newStaffMember,
                          canManageProducts: e.target.checked,
                        })
                      }
                    />
                    <span className="label-text ml-2">{t('products')}</span>
                  </label>
                  <label className="cursor-pointer label justify-start">
                    <input
                      type="checkbox"
                      className="checkbox checkbox-primary"
                      checked={newStaffMember.canManageOrders}
                      onChange={(e) =>
                        setNewStaffMember({
                          ...newStaffMember,
                          canManageOrders: e.target.checked,
                        })
                      }
                    />
                    <span className="label-text ml-2">{t('orders')}</span>
                  </label>
                  <label className="cursor-pointer label justify-start">
                    <input
                      type="checkbox"
                      className="checkbox checkbox-primary"
                      checked={newStaffMember.canManageSettings}
                      onChange={(e) =>
                        setNewStaffMember({
                          ...newStaffMember,
                          canManageSettings: e.target.checked,
                        })
                      }
                    />
                    <span className="label-text ml-2">{t('settings')}</span>
                  </label>
                </div>
              </div>

              <button
                className="btn btn-primary mt-4"
                onClick={addStaffMember}
                disabled={!newStaffMember.email}
              >
                {tCommon('add')}
              </button>
            </div>
          </div>

          <div className="alert alert-info mt-4">
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
            <span>{t('tip')}</span>
          </div>

          <div className="card-actions justify-between mt-6">
            <button className="btn btn-ghost" onClick={onBack}>
              {tCommon('back')}
            </button>
            <button className="btn btn-primary" onClick={handleNext}>
              {tCommon('next')}
            </button>
          </div>
        </div>
      </div>
    </div>
  );
}
