'use client';

import { useState, useEffect, useCallback } from 'react';
import { useTranslations } from 'next-intl';
import { tokenManager } from '@/utils/tokenManager';
import config from '@/config';

interface Role {
  id: number;
  name: string;
  display_name: string;
  description: string;
  is_system: boolean;
}

interface UserRole {
  role_id: number;
  role_name: string;
  display_name: string;
  description?: string;
  granted_at: string;
  expires_at?: string;
  is_active: boolean;
  notes?: string;
}

interface RoleManagerProps {
  userId: number;
  userEmail: string;
  currentRoles?: UserRole[];
  onRoleUpdate?: () => void;
}

export default function RoleManager({
  userId,
  userEmail: _userEmail,
  currentRoles = [],
  onRoleUpdate,
}: RoleManagerProps) {
  const t = useTranslations('admin');
  const [availableRoles, setAvailableRoles] = useState<Role[]>([]);
  const [userRoles, setUserRoles] = useState<UserRole[]>(currentRoles);
  const [loading, setLoading] = useState(false);
  const [error, setError] = useState<string | null>(null);
  const [_selectedRole, setSelectedRole] = useState<string>('');
  const [isDropdownOpen, setIsDropdownOpen] = useState(false);

  const fetchAvailableRoles = async () => {
    try {
      const token = tokenManager.getAccessToken();
      if (!token) return;

      const response = await fetch(`${config.getApiUrl()}/api/v1/roles`, {
        headers: {
          Authorization: `Bearer ${token}`,
        },
      });

      if (response.ok) {
        const data = await response.json();
        setAvailableRoles(data.roles || []);
      }
    } catch (err) {
      console.error('Failed to fetch roles:', err);
    }
  };

  const fetchUserRoles = useCallback(async () => {
    try {
      const token = tokenManager.getAccessToken();
      if (!token) return;

      const response = await fetch(
        `${config.getApiUrl()}/api/v1/users/${userId}/roles`,
        {
          headers: {
            Authorization: `Bearer ${token}`,
          },
        }
      );

      if (response.ok) {
        const data = await response.json();
        setUserRoles(data.roles || []);
      }
    } catch (err) {
      console.error('Failed to fetch user roles:', err);
    }
  }, [userId]);

  // Fetch available roles from Auth Service
  useEffect(() => {
    fetchAvailableRoles();
  }, []);

  // Fetch user roles from Auth Service
  useEffect(() => {
    if (userId) {
      fetchUserRoles();
    }
  }, [userId, fetchUserRoles]);

  const assignRole = async (roleName: string) => {
    try {
      setLoading(true);
      setError(null);

      const token = tokenManager.getAccessToken();
      if (!token) {
        throw new Error('No authentication token');
      }

      const response = await fetch(
        `${config.getApiUrl()}/api/v1/roles/assign`,
        {
          method: 'POST',
          headers: {
            Authorization: `Bearer ${token}`,
            'Content-Type': 'application/json',
          },
          body: JSON.stringify({
            user_id: userId,
            role_name: roleName,
            notes: `Assigned via admin panel by admin`,
          }),
        }
      );

      if (!response.ok) {
        const errorData = await response.json();
        throw new Error(errorData.error || 'Failed to assign role');
      }

      // Refresh roles
      await fetchUserRoles();
      if (onRoleUpdate) {
        onRoleUpdate();
      }

      setIsDropdownOpen(false);
      setSelectedRole('');
    } catch (err: any) {
      setError(err.message || 'Failed to assign role');
    } finally {
      setLoading(false);
    }
  };

  const revokeRole = async (roleName: string) => {
    if (!confirm(t('users.confirmRevokeRole', { role: roleName }))) {
      return;
    }

    try {
      setLoading(true);
      setError(null);

      const token = tokenManager.getAccessToken();
      if (!token) {
        throw new Error('No authentication token');
      }

      const response = await fetch(
        `${config.getApiUrl()}/api/v1/roles/revoke`,
        {
          method: 'POST',
          headers: {
            Authorization: `Bearer ${token}`,
            'Content-Type': 'application/json',
          },
          body: JSON.stringify({
            user_id: userId,
            role_name: roleName,
            reason: 'Revoked via admin panel',
          }),
        }
      );

      if (!response.ok) {
        const errorData = await response.json();
        throw new Error(errorData.error || 'Failed to revoke role');
      }

      // Refresh roles
      await fetchUserRoles();
      if (onRoleUpdate) {
        onRoleUpdate();
      }
    } catch (err: any) {
      setError(err.message || 'Failed to revoke role');
    } finally {
      setLoading(false);
    }
  };

  const getRoleBadgeColor = (roleName: string) => {
    switch (roleName) {
      case 'admin':
        return 'badge-error';
      case 'moderator':
        return 'badge-warning';
      case 'support':
        return 'badge-info';
      case 'user':
        return 'badge-success';
      default:
        return 'badge-ghost';
    }
  };

  return (
    <div className="space-y-4">
      {/* Current Roles */}
      <div>
        <h4 className="text-sm font-semibold mb-2">
          {t('users.currentRoles')}
        </h4>
        <div className="flex flex-wrap gap-2">
          {userRoles.length === 0 ? (
            <span className="text-gray-500 text-sm">{t('users.noRoles')}</span>
          ) : (
            userRoles.map((role) => (
              <div key={role.role_name} className="flex items-center gap-1">
                <span className={`badge ${getRoleBadgeColor(role.role_name)}`}>
                  {role.display_name || role.role_name}
                </span>
                {role.role_name !== 'user' && (
                  <button
                    onClick={() => revokeRole(role.role_name)}
                    disabled={loading}
                    className="btn btn-ghost btn-xs text-error"
                    title={t('users.revokeRole')}
                  >
                    ×
                  </button>
                )}
              </div>
            ))
          )}
        </div>
      </div>

      {/* Add Role Dropdown */}
      <div className="relative">
        <button
          onClick={() => setIsDropdownOpen(!isDropdownOpen)}
          disabled={loading}
          className="btn btn-sm btn-primary"
        >
          {loading ? (
            <span className="loading loading-spinner loading-xs"></span>
          ) : (
            t('users.addRole')
          )}
        </button>

        {isDropdownOpen && (
          <div className="absolute z-10 mt-2 w-64 rounded-md shadow-lg bg-base-100 ring-1 ring-black ring-opacity-5">
            <div className="py-1">
              {availableRoles
                .filter(
                  (role) => !userRoles.some((ur) => ur.role_name === role.name)
                )
                .map((role) => (
                  <button
                    key={role.name}
                    onClick={() => assignRole(role.name)}
                    className="block w-full text-left px-4 py-2 text-sm hover:bg-base-200"
                  >
                    <div className="font-medium">{role.display_name}</div>
                    {role.description && (
                      <div className="text-xs text-gray-500">
                        {role.description}
                      </div>
                    )}
                  </button>
                ))}
            </div>
          </div>
        )}
      </div>

      {/* Error Message */}
      {error && (
        <div className="alert alert-error">
          <span>{error}</span>
        </div>
      )}
    </div>
  );
}
