'use client';

import { useState, useEffect } from 'react';
import { useTranslations } from 'next-intl';
import { FiX, FiUser } from 'react-icons/fi';
import { apiClientAuth } from '@/lib/api-client-auth';

interface Admin {
  id: number;
  name: string;
  email: string;
}

interface AssignProblemModalProps {
  isOpen: boolean;
  onClose: () => void;
  onAssign: (adminId: number) => void;
  problemId: number;
  currentAssignedTo?: number;
}

export default function AssignProblemModal({
  isOpen,
  onClose,
  onAssign,
  problemId,
  currentAssignedTo,
}: AssignProblemModalProps) {
  const t = useTranslations('admin');
  const [admins, setAdmins] = useState<Admin[]>([]);
  const [selectedAdminId, setSelectedAdminId] = useState<number | null>(
    currentAssignedTo || null
  );
  const [loading, setLoading] = useState(false);
  const [error, setError] = useState<string | null>(null);

  useEffect(() => {
    if (isOpen) {
      fetchAdmins();
    }
  }, [isOpen]);

  const fetchAdmins = async () => {
    try {
      setLoading(true);
      const response = await apiClientAuth.get('/admin/users?role=admin');
      if (response.data && Array.isArray(response.data)) {
        setAdmins(response.data);
      } else {
        // Fallback на тестовые данные если API не работает
        setAdmins([
          { id: 1, name: 'System Administrator', email: 'admin@system.local' },
        ]);
      }
    } catch (err) {
      console.error('Error fetching admins:', err);
      // Fallback на тестовые данные в случае ошибки
      setAdmins([
        { id: 1, name: 'System Administrator', email: 'admin@system.local' },
      ]);
      setError('Failed to load admins from API, using fallback data');
    } finally {
      setLoading(false);
    }
  };

  const handleAssign = () => {
    if (selectedAdminId) {
      onAssign(selectedAdminId);
      onClose();
    }
  };

  if (!isOpen) return null;

  return (
    <div className="modal modal-open">
      <div className="modal-box">
        <div className="flex justify-between items-center mb-4">
          <h3 className="font-bold text-lg">
            {t('logistics.problems.assignModal.title')}
          </h3>
          <button className="btn btn-sm btn-circle btn-ghost" onClick={onClose}>
            <FiX className="w-4 h-4" />
          </button>
        </div>

        {error && (
          <div className="alert alert-error mb-4">
            <span>{error}</span>
          </div>
        )}

        <div className="form-control">
          <label className="label">
            <span className="label-text">
              {t('logistics.problems.assignModal.selectAdmin')}
            </span>
          </label>

          {loading ? (
            <div className="flex justify-center py-4">
              <span className="loading loading-spinner"></span>
            </div>
          ) : (
            <div className="space-y-2">
              {admins.map((admin) => (
                <label
                  key={admin.id}
                  className="flex items-center gap-3 p-3 border rounded-lg cursor-pointer hover:bg-base-200"
                >
                  <input
                    type="radio"
                    name="admin"
                    className="radio"
                    value={admin.id}
                    checked={selectedAdminId === admin.id}
                    onChange={() => setSelectedAdminId(admin.id)}
                  />
                  <div className="flex items-center gap-2">
                    <FiUser className="w-4 h-4 text-gray-500" />
                    <div>
                      <div className="font-medium">{admin.name}</div>
                      <div className="text-sm text-gray-500">{admin.email}</div>
                    </div>
                  </div>
                </label>
              ))}
            </div>
          )}
        </div>

        <div className="modal-action">
          <button className="btn btn-ghost" onClick={onClose}>
            {t('common.cancel')}
          </button>
          <button
            className="btn btn-primary"
            onClick={handleAssign}
            disabled={!selectedAdminId || loading}
          >
            {t('logistics.problems.assignModal.assign')}
          </button>
        </div>
      </div>
    </div>
  );
}
