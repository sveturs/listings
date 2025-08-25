'use client';

import { useState } from 'react';
import { useTranslations } from 'next-intl';
import { FiX, FiCheckCircle } from 'react-icons/fi';

interface ResolveProblemModalProps {
  isOpen: boolean;
  onClose: () => void;
  onResolve: (resolution: string) => void;
  problemId: number;
  problemDescription: string;
}

export default function ResolveProblemModal({
  isOpen,
  onClose,
  onResolve,
  problemDescription,
}: ResolveProblemModalProps) {
  const t = useTranslations('admin');
  const [resolution, setResolution] = useState('');
  const [loading, setLoading] = useState(false);

  const handleResolve = async () => {
    if (!resolution.trim()) return;

    setLoading(true);
    try {
      await onResolve(resolution);
      onClose();
      setResolution('');
    } finally {
      setLoading(false);
    }
  };

  const handleClose = () => {
    onClose();
    setResolution('');
  };

  if (!isOpen) return null;

  return (
    <div className="modal modal-open">
      <div className="modal-box max-w-2xl">
        <div className="flex justify-between items-center mb-4">
          <h3 className="font-bold text-lg flex items-center gap-2">
            <FiCheckCircle className="w-5 h-5 text-success" />
            {t('logistics.problems.resolveModal.title')}
          </h3>
          <button
            className="btn btn-sm btn-circle btn-ghost"
            onClick={handleClose}
          >
            <FiX className="w-4 h-4" />
          </button>
        </div>

        <div className="space-y-4">
          <div>
            <label className="label">
              <span className="label-text font-medium">
                {t('logistics.problems.resolveModal.problemDescription')}
              </span>
            </label>
            <div className="p-3 bg-base-200 rounded-lg text-sm">
              {problemDescription}
            </div>
          </div>

          <div className="form-control">
            <label className="label">
              <span className="label-text font-medium">
                {t('logistics.problems.resolveModal.resolution')} *
              </span>
            </label>
            <textarea
              className="textarea textarea-bordered h-32"
              placeholder={t(
                'logistics.problems.resolveModal.resolutionPlaceholder'
              )}
              value={resolution}
              onChange={(e) => setResolution(e.target.value)}
              maxLength={500}
            />
            <label className="label">
              <span className="label-text-alt">{resolution.length}/500</span>
            </label>
          </div>
        </div>

        <div className="modal-action">
          <button
            className="btn btn-ghost"
            onClick={handleClose}
            disabled={loading}
          >
            {t('common.cancel')}
          </button>
          <button
            className="btn btn-success"
            onClick={handleResolve}
            disabled={!resolution.trim() || loading}
          >
            {loading && (
              <span className="loading loading-spinner loading-sm"></span>
            )}
            {t('logistics.problems.resolveModal.resolve')}
          </button>
        </div>
      </div>
    </div>
  );
}
