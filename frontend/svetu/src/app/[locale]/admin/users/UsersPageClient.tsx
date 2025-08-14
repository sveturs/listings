'use client';

import { useState, useEffect, useCallback } from 'react';
import { useTranslations } from 'next-intl';
import { useAuth } from '@/contexts/AuthContext';
import config from '@/config';
import { tokenManager } from '@/utils/tokenManager';
import { createAdminHeaders } from '@/utils/csrf';
import type { components } from '@/types/generated/api';

type UserProfile =
  components['schemas']['backend_internal_domain_models.UserProfile'];

export default function UsersPageClient() {
  const t = useTranslations('admin');
  const tCommon = useTranslations('common');
  const { user } = useAuth();

  const [users, setUsers] = useState<UserProfile[]>([]);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState<string | null>(null);
  const [currentPage, setCurrentPage] = useState(1);
  const [totalPages, setTotalPages] = useState(1);
  const [totalCount, setTotalCount] = useState(0);
  const [searchQuery, setSearchQuery] = useState('');
  const [statusFilter, setStatusFilter] = useState<string>('');
  const [selectedUser, setSelectedUser] = useState<UserProfile | null>(null);
  const [isModalOpen, setIsModalOpen] = useState(false);
  const [isDeleteModalOpen, setIsDeleteModalOpen] = useState(false);
  const [userToDelete, setUserToDelete] = useState<UserProfile | null>(null);
  const [updatingStatus, setUpdatingStatus] = useState<number | null>(null);
  const [deleting, setDeleting] = useState<number | null>(null);
  const [pageSize, setPageSize] = useState(25);
  const [sortField, setSortField] = useState<string>('id');
  const [sortOrder, setSortOrder] = useState<'asc' | 'desc'>('asc');
  const [roles, setRoles] = useState<any[]>([]);
  const [updatingRole, setUpdatingRole] = useState<number | null>(null);
  const [openRoleDropdown, setOpenRoleDropdown] = useState<number | null>(null);
  const [dropdownPosition, setDropdownPosition] = useState<{
    top?: string;
    bottom?: string;
  }>({});

  const pageSizeOptions = [25, 50, 100, 250, 1000];

  // Функция для вычисления позиции dropdown
  const calculateDropdownPosition = (userId: number) => {
    const buttonElement = document.querySelector(
      `[data-user-id="${userId}"] button`
    );
    if (!buttonElement) return { top: '100%' };

    const rect = buttonElement.getBoundingClientRect();
    const dropdownHeight = 320; // Примерная высота dropdown
    const spaceBelow = window.innerHeight - rect.bottom;
    const spaceAbove = rect.top;

    if (spaceBelow < dropdownHeight && spaceAbove > dropdownHeight) {
      return { bottom: '100%', top: 'auto' };
    } else {
      return { top: '100%', bottom: 'auto' };
    }
  };

  // Функция для получения локализованного описания роли
  const getRoleDescription = (roleName: string) => {
    const key = `users.roleDescriptions.${roleName}`;
    try {
      return t(key);
    } catch {
      // Fallback на оригинальное описание из API если перевод не найден
      const role = roles.find((r) => r.name === roleName);
      return role?.description || t('users.roleTooltip.noDescription');
    }
  };

  const fetchUsers = useCallback(async () => {
    try {
      setLoading(true);
      setError(null);

      const params = new URLSearchParams({
        page: currentPage.toString(),
        limit: pageSize.toString(),
        sort_by: sortField,
        sort_order: sortOrder,
      });

      if (statusFilter) {
        params.append('status', statusFilter);
      }

      const token = tokenManager.getAccessToken();
      if (!token) {
        throw new Error('No authentication token available');
      }

      const response = await fetch(
        `${config.getApiUrl()}/api/v1/admin/users?${params}`,
        {
          headers: {
            Authorization: `Bearer ${token}`,
          },
        }
      );

      if (!response.ok) {
        throw new Error(`Failed to fetch users: ${response.status}`);
      }

      const data = await response.json();

      // API возвращает структуру: { data: { data: [...], total: 9, ... }, success: true }
      const apiResponse = data.data;

      // Преобразуем ответ API в ожидаемый формат
      const users = apiResponse.data || [];
      const totalCount = apiResponse.total || 0;

      setUsers(users);
      setTotalCount(totalCount);
      setTotalPages(Math.ceil(totalCount / pageSize));
    } catch (err) {
      console.error('Error fetching users:', err);
      setError(t('users.error.fetchFailed'));
    } finally {
      setLoading(false);
    }
  }, [currentPage, statusFilter, pageSize, sortField, sortOrder, t]);

  const fetchRoles = useCallback(async () => {
    try {
      const token = tokenManager.getAccessToken();
      if (!token) {
        throw new Error('No authentication token available');
      }

      const response = await fetch(`${config.getApiUrl()}/api/v1/admin/roles`, {
        headers: {
          Authorization: `Bearer ${token}`,
        },
      });

      if (!response.ok) {
        throw new Error(`Failed to fetch roles: ${response.status}`);
      }

      const data = await response.json();
      setRoles(data.data || []);
    } catch (err) {
      console.error('Error fetching roles:', err);
    }
  }, []);

  useEffect(() => {
    fetchUsers();
    fetchRoles();
  }, [fetchUsers, fetchRoles]);

  // Закрытие dropdown при клике вне его и управление позиционированием
  useEffect(() => {
    const handleClickOutside = () => {
      setOpenRoleDropdown(null);
    };

    const handleScroll = () => {
      // Принудительно перерендерить dropdown для пересчета позиции
      if (openRoleDropdown !== null) {
        const currentDropdown = openRoleDropdown;
        setOpenRoleDropdown(null);
        setTimeout(() => setOpenRoleDropdown(currentDropdown), 0);
      }
    };

    if (openRoleDropdown !== null) {
      document.addEventListener('click', handleClickOutside);
      window.addEventListener('scroll', handleScroll, true);
      window.addEventListener('resize', handleScroll);
      return () => {
        document.removeEventListener('click', handleClickOutside);
        window.removeEventListener('scroll', handleScroll, true);
        window.removeEventListener('resize', handleScroll);
      };
    }
  }, [openRoleDropdown]);

  const handleStatusChange = async (userId: number, newStatus: string) => {
    if (userId === user?.id) {
      alert(t('users.error.cannotModifySelf'));
      return;
    }

    try {
      setUpdatingStatus(userId);

      const token = tokenManager.getAccessToken();
      if (!token) {
        throw new Error('No authentication token available');
      }

      const headers = await createAdminHeaders(token);
      const response = await fetch(
        `${config.getApiUrl()}/api/v1/admin/users/${userId}/status`,
        {
          method: 'PUT',
          headers,
          credentials: 'include',
          body: JSON.stringify({ status: newStatus }),
        }
      );

      if (!response.ok) {
        throw new Error(`Failed to update status: ${response.status}`);
      }

      // Update local state
      setUsers((prevUsers) =>
        prevUsers.map((u) =>
          u.id === userId ? { ...u, account_status: newStatus } : u
        )
      );
    } catch (err) {
      console.error('Error updating user status:', err);
      alert(t('users.error.updateFailed'));
    } finally {
      setUpdatingStatus(null);
    }
  };

  const handleRoleChange = async (userId: number, newRoleId: number) => {
    if (userId === user?.id) {
      alert(t('users.error.cannotModifySelf'));
      return;
    }

    try {
      setUpdatingRole(userId);

      const token = tokenManager.getAccessToken();
      if (!token) {
        throw new Error('No authentication token available');
      }

      const headers = await createAdminHeaders(token);
      const response = await fetch(
        `${config.getApiUrl()}/api/v1/admin/users/${userId}/role`,
        {
          method: 'PUT',
          headers,
          credentials: 'include', // Важно для отправки cookies
          body: JSON.stringify({ role_id: newRoleId }),
        }
      );

      if (!response.ok) {
        throw new Error(`Failed to update role: ${response.status}`);
      }

      // Update local state
      const updatedRole = roles.find((r) => r.id === newRoleId);
      setUsers((prevUsers) =>
        prevUsers.map((u) =>
          u.id === userId ? { ...u, role_id: newRoleId, role: updatedRole } : u
        )
      );
    } catch (err) {
      console.error('Error updating user role:', err);
      alert(t('users.error.updateFailed'));
    } finally {
      setUpdatingRole(null);
    }
  };

  const handleSort = (field: string) => {
    if (field === sortField) {
      setSortOrder(sortOrder === 'asc' ? 'desc' : 'asc');
    } else {
      setSortField(field);
      setSortOrder('asc');
    }
    setCurrentPage(1);
  };

  const handlePageSizeChange = (newSize: number) => {
    setPageSize(newSize);
    setCurrentPage(1);
  };

  const handleDeleteUser = async () => {
    if (!userToDelete || !userToDelete.id) return;

    if (userToDelete.id === user?.id) {
      alert(t('users.error.cannotDeleteSelf'));
      return;
    }

    try {
      setDeleting(userToDelete.id);

      const token = tokenManager.getAccessToken();
      if (!token) {
        throw new Error('No authentication token available');
      }

      const headers = await createAdminHeaders(token);
      const response = await fetch(
        `${config.getApiUrl()}/api/v1/admin/users/${userToDelete.id}`,
        {
          method: 'DELETE',
          headers,
          credentials: 'include',
        }
      );

      if (!response.ok) {
        throw new Error(`Failed to delete user: ${response.status}`);
      }

      // Remove from local state
      setUsers((prevUsers) =>
        prevUsers.filter((u) => u.id !== userToDelete.id)
      );
      setTotalCount((prev) => prev - 1);
      setIsDeleteModalOpen(false);
      setUserToDelete(null);
    } catch (err) {
      console.error('Error deleting user:', err);
      alert(t('users.error.deleteFailed'));
    } finally {
      setDeleting(null);
    }
  };

  const openUserModal = (user: UserProfile) => {
    setSelectedUser(user);
    setIsModalOpen(true);
  };

  const openDeleteModal = (user: UserProfile) => {
    setUserToDelete(user);
    setIsDeleteModalOpen(true);
  };

  const getStatusBadgeClass = (status: string) => {
    switch (status) {
      case 'active':
        return 'badge-success';
      case 'inactive':
        return 'badge-warning';
      case 'suspended':
        return 'badge-error';
      default:
        return 'badge-ghost';
    }
  };

  const formatDate = (date: string | undefined) => {
    if (!date) return '-';
    return new Date(date).toLocaleDateString();
  };

  if (loading && users.length === 0) {
    return (
      <div className="flex justify-center items-center h-64">
        <span className="loading loading-spinner loading-lg"></span>
      </div>
    );
  }

  return (
    <div>
      <div className="flex justify-between items-center mb-6">
        <h1 className="text-3xl font-bold">{t('users.title')}</h1>
        <div className="badge badge-neutral">
          {t('users.totalUsers', { count: totalCount })}
        </div>
      </div>

      {/* Filters */}
      <div className="card bg-base-100 shadow-xl mb-6">
        <div className="card-body">
          <div className="flex gap-4 flex-wrap items-end">
            <div className="form-control flex-1 min-w-[200px]">
              <input
                type="text"
                placeholder={t('users.searchPlaceholder')}
                className="input input-bordered"
                value={searchQuery}
                onChange={(e) => setSearchQuery(e.target.value)}
              />
            </div>
            <select
              className="select select-bordered min-w-[150px]"
              value={statusFilter}
              onChange={(e) => {
                setStatusFilter(e.target.value);
                setCurrentPage(1);
              }}
            >
              <option value="">{t('users.allStatuses')}</option>
              <option value="active">{t('users.status.active')}</option>
              <option value="inactive">{t('users.status.inactive')}</option>
              <option value="suspended">{t('users.status.suspended')}</option>
            </select>
            <div className="form-control">
              <label className="label">
                <span className="label-text">{t('users.itemsPerPage')}</span>
              </label>
              <select
                className="select select-bordered"
                value={pageSize}
                onChange={(e) => handlePageSizeChange(Number(e.target.value))}
              >
                {pageSizeOptions.map((size) => (
                  <option key={size} value={size}>
                    {size}
                  </option>
                ))}
              </select>
            </div>
          </div>
        </div>
      </div>

      {/* Users Table */}
      <div className="card bg-base-100 shadow-xl">
        <div className="card-body p-0">
          <div className="overflow-x-auto">
            <table className="table">
              <thead>
                <tr>
                  <th
                    className="cursor-pointer hover:bg-base-200 select-none"
                    onClick={() => handleSort('id')}
                  >
                    <div className="flex items-center gap-1">
                      {t('users.table.id')}
                      {sortField === 'id' && (
                        <span className="text-primary">
                          {sortOrder === 'asc' ? '↑' : '↓'}
                        </span>
                      )}
                    </div>
                  </th>
                  <th
                    className="cursor-pointer hover:bg-base-200 select-none"
                    onClick={() => handleSort('name')}
                  >
                    <div className="flex items-center gap-1">
                      {t('users.table.name')}
                      {sortField === 'name' && (
                        <span className="text-primary">
                          {sortOrder === 'asc' ? '↑' : '↓'}
                        </span>
                      )}
                    </div>
                  </th>
                  <th
                    className="cursor-pointer hover:bg-base-200 select-none"
                    onClick={() => handleSort('email')}
                  >
                    <div className="flex items-center gap-1">
                      {t('users.table.email')}
                      {sortField === 'email' && (
                        <span className="text-primary">
                          {sortOrder === 'asc' ? '↑' : '↓'}
                        </span>
                      )}
                    </div>
                  </th>
                  <th
                    className="cursor-pointer hover:bg-base-200 select-none"
                    onClick={() => handleSort('account_status')}
                  >
                    <div className="flex items-center gap-1">
                      {t('users.table.status')}
                      {sortField === 'account_status' && (
                        <span className="text-primary">
                          {sortOrder === 'asc' ? '↑' : '↓'}
                        </span>
                      )}
                    </div>
                  </th>
                  <th>{t('users.table.role')}</th>
                  <th>{t('users.table.provider')}</th>
                  <th
                    className="cursor-pointer hover:bg-base-200 select-none"
                    onClick={() => handleSort('created_at')}
                  >
                    <div className="flex items-center gap-1">
                      {t('users.table.createdAt')}
                      {sortField === 'created_at' && (
                        <span className="text-primary">
                          {sortOrder === 'asc' ? '↑' : '↓'}
                        </span>
                      )}
                    </div>
                  </th>
                  <th
                    className="cursor-pointer hover:bg-base-200 select-none"
                    onClick={() => handleSort('last_seen')}
                  >
                    <div className="flex items-center gap-1">
                      {t('users.table.lastSeen')}
                      {sortField === 'last_seen' && (
                        <span className="text-primary">
                          {sortOrder === 'asc' ? '↑' : '↓'}
                        </span>
                      )}
                    </div>
                  </th>
                  <th>{t('users.table.actions')}</th>
                </tr>
              </thead>
              <tbody>
                {users
                  .filter(
                    (u) =>
                      !searchQuery ||
                      u.name
                        ?.toLowerCase()
                        .includes(searchQuery.toLowerCase()) ||
                      u.email?.toLowerCase().includes(searchQuery.toLowerCase())
                  )
                  .map((userItem) => (
                    <tr
                      key={userItem.id}
                      className="hover"
                      data-user-id={userItem.id}
                    >
                      <td>{userItem.id}</td>
                      <td>
                        <div className="flex items-center space-x-3">
                          {userItem.picture_url && (
                            <div className="avatar">
                              <div className="mask mask-circle w-8 h-8">
                                <img
                                  src={userItem.picture_url}
                                  alt={userItem.name}
                                />
                              </div>
                            </div>
                          )}
                          <div>
                            <div className="font-bold">{userItem.name}</div>
                            {userItem.city && (
                              <div className="text-sm opacity-50">
                                {userItem.city}
                                {userItem.country && `, ${userItem.country}`}
                              </div>
                            )}
                          </div>
                        </div>
                      </td>
                      <td>{userItem.email}</td>
                      <td>
                        <select
                          className={`select select-bordered select-sm ${
                            updatingStatus === userItem.id ? 'loading' : ''
                          }`}
                          value={userItem.account_status}
                          onChange={(e) =>
                            userItem.id &&
                            handleStatusChange(userItem.id, e.target.value)
                          }
                          disabled={
                            updatingStatus === userItem.id ||
                            userItem.id === user?.id
                          }
                        >
                          <option value="active">
                            {t('users.status.active')}
                          </option>
                          <option value="inactive">
                            {t('users.status.inactive')}
                          </option>
                          <option value="suspended">
                            {t('users.status.suspended')}
                          </option>
                        </select>
                      </td>
                      <td>
                        <div className="relative">
                          <button
                            className={`btn btn-sm btn-outline ${
                              updatingRole === userItem.id ? 'loading' : ''
                            } ${
                              userItem.id === user?.id ? 'btn-disabled' : ''
                            }`}
                            onClick={(e) => {
                              e.stopPropagation();
                              if (
                                updatingRole === userItem.id ||
                                userItem.id === user?.id
                              )
                                return;

                              if (openRoleDropdown === userItem.id) {
                                setOpenRoleDropdown(null);
                              } else {
                                if (userItem.id) {
                                  setOpenRoleDropdown(userItem.id);
                                }
                                // Вычисляем позицию после небольшой задержки, чтобы dropdown появился
                                setTimeout(() => {
                                  if (userItem.id) {
                                    const position = calculateDropdownPosition(
                                      userItem.id
                                    );
                                    setDropdownPosition(position);
                                  }
                                }, 10);
                              }
                            }}
                            disabled={
                              updatingRole === userItem.id ||
                              userItem.id === user?.id
                            }
                          >
                            {userItem.role?.display_name ||
                              t('users.selectRole')}
                            <svg
                              className="w-4 h-4 ml-1"
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

                          {openRoleDropdown === userItem.id && (
                            <div
                              className="absolute right-0 bg-white rounded-lg shadow-xl border border-gray-200 z-[9999]"
                              style={{
                                width: '380px',
                                maxHeight: roles.length > 8 ? '320px' : 'auto',
                                overflowY:
                                  roles.length > 8 ? 'auto' : 'visible',
                                marginTop:
                                  dropdownPosition.top === '100%'
                                    ? '4px'
                                    : undefined,
                                marginBottom:
                                  dropdownPosition.bottom === '100%'
                                    ? '4px'
                                    : undefined,
                                ...dropdownPosition,
                              }}
                            >
                              {roles.map((role) => (
                                <div
                                  key={role.id}
                                  className="relative px-3 py-3 hover:bg-gray-50 cursor-pointer transition-colors border-b border-gray-100 last:border-b-0"
                                  onClick={(e) => {
                                    e.stopPropagation();
                                    if (userItem.id) {
                                      handleRoleChange(userItem.id, role.id);
                                    }
                                    setOpenRoleDropdown(null);
                                  }}
                                  title={getRoleDescription(role.name)}
                                  onMouseEnter={(e) => {
                                    const tooltip =
                                      e.currentTarget.querySelector(
                                        '.custom-tooltip'
                                      ) as HTMLElement;
                                    if (tooltip) {
                                      const rect =
                                        e.currentTarget.getBoundingClientRect();
                                      tooltip.style.position = 'fixed';
                                      tooltip.style.left = `${Math.max(10, rect.left - 320)}px`;
                                      tooltip.style.top = `${rect.top}px`;
                                      tooltip.style.zIndex = '9999';
                                      tooltip.style.opacity = '1';
                                      tooltip.style.visibility = 'visible';
                                      tooltip.style.transform = 'translateY(0)';
                                    }
                                  }}
                                  onMouseLeave={(e) => {
                                    const tooltip =
                                      e.currentTarget.querySelector(
                                        '.custom-tooltip'
                                      ) as HTMLElement;
                                    if (tooltip) {
                                      tooltip.style.opacity = '0';
                                      tooltip.style.visibility = 'hidden';
                                      tooltip.style.transform =
                                        'translateY(-10px)';
                                    }
                                  }}
                                >
                                  <div className="flex flex-col items-start w-full">
                                    <span className="font-medium text-sm">
                                      {role.display_name}
                                    </span>
                                    <span
                                      className="text-xs opacity-60 leading-tight mt-1"
                                      style={{
                                        display: '-webkit-box',
                                        WebkitLineClamp: 2,
                                        WebkitBoxOrient: 'vertical',
                                        overflow: 'hidden',
                                        wordBreak: 'break-word',
                                      }}
                                    >
                                      {getRoleDescription(role.name)}
                                    </span>
                                  </div>
                                  {/* Custom tooltip positioned outside dropdown with JavaScript */}
                                  <div
                                    className="custom-tooltip fixed bg-gray-800 text-white px-3 py-2 rounded-lg text-sm max-w-80 shadow-lg pointer-events-none whitespace-normal transition-all duration-200"
                                    style={{
                                      opacity: 0,
                                      visibility: 'hidden',
                                      transform: 'translateY(-10px)',
                                      zIndex: 9999,
                                    }}
                                  >
                                    <div className="font-medium mb-1">
                                      {role.display_name}
                                    </div>
                                    <div>{getRoleDescription(role.name)}</div>
                                    {/* Arrow pointing to the right */}
                                    <div className="absolute top-3 -right-1 w-2 h-2 bg-gray-800 transform rotate-45"></div>
                                  </div>
                                </div>
                              ))}
                            </div>
                          )}
                        </div>
                      </td>
                      <td>
                        <span
                          className={`badge ${userItem.google_id ? 'badge-primary' : 'badge-ghost'} gap-1`}
                        >
                          {userItem.google_id ? (
                            <>
                              <svg
                                className="w-3 h-3"
                                viewBox="0 0 24 24"
                                fill="currentColor"
                              >
                                <path d="M22.56 12.25c0-.78-.07-1.53-.2-2.25H12v4.26h5.92c-.26 1.37-1.04 2.53-2.21 3.31v2.77h3.57c2.08-1.92 3.28-4.74 3.28-8.09z" />
                                <path d="M12 23c2.97 0 5.46-.98 7.28-2.66l-3.57-2.77c-.98.66-2.23 1.06-3.71 1.06-2.86 0-5.29-1.93-6.16-4.53H2.18v2.84C3.99 20.53 7.7 23 12 23z" />
                                <path d="M5.84 14.09c-.22-.66-.35-1.36-.35-2.09s.13-1.43.35-2.09V7.07H2.18C1.43 8.55 1 10.22 1 12s.43 3.45 1.18 4.93l2.85-2.22.81-.62z" />
                                <path d="M12 5.38c1.62 0 3.06.56 4.21 1.64l3.15-3.15C17.45 2.09 14.97 1 12 1 7.7 1 3.99 3.47 2.18 7.07l3.66 2.84c.87-2.6 3.3-4.53 6.16-4.53z" />
                              </svg>
                              {t('users.providers.google')}
                            </>
                          ) : (
                            <>
                              <svg
                                className="w-3 h-3"
                                fill="none"
                                stroke="currentColor"
                                viewBox="0 0 24 24"
                              >
                                <path
                                  strokeLinecap="round"
                                  strokeLinejoin="round"
                                  strokeWidth={2}
                                  d="M3 8l7.89 5.26a2 2 0 002.22 0L21 8M5 19h14a2 2 0 002-2V7a2 2 0 00-2-2H5a2 2 0 00-2 2v10a2 2 0 002 2z"
                                />
                              </svg>
                              {t('users.providers.email')}
                            </>
                          )}
                        </span>
                      </td>
                      <td>{formatDate(userItem.created_at)}</td>
                      <td>{formatDate(userItem.last_seen)}</td>
                      <td>
                        <div className="flex gap-2">
                          <button
                            className="btn btn-ghost btn-xs"
                            onClick={() => openUserModal(userItem)}
                          >
                            {t('users.view')}
                          </button>
                          <button
                            className="btn btn-ghost btn-xs text-error"
                            onClick={() => openDeleteModal(userItem)}
                            disabled={userItem.id === user?.id}
                          >
                            {tCommon('delete')}
                          </button>
                        </div>
                      </td>
                    </tr>
                  ))}
              </tbody>
            </table>
          </div>

          {/* Pagination */}
          {totalPages > 1 && (
            <div className="flex justify-center p-4">
              <div className="join">
                <button
                  className="join-item btn"
                  onClick={() => setCurrentPage((p) => Math.max(1, p - 1))}
                  disabled={currentPage === 1}
                >
                  «
                </button>
                <button className="join-item btn">
                  {t('users.pageInfo', {
                    current: currentPage,
                    total: totalPages,
                  })}
                </button>
                <button
                  className="join-item btn"
                  onClick={() =>
                    setCurrentPage((p) => Math.min(totalPages, p + 1))
                  }
                  disabled={currentPage === totalPages}
                >
                  »
                </button>
              </div>
            </div>
          )}
        </div>
      </div>

      {/* User Details Modal */}
      {isModalOpen && selectedUser && (
        <dialog className="modal modal-open">
          <div className="modal-box max-w-2xl">
            <h3 className="font-bold text-lg mb-4">{t('users.userDetails')}</h3>

            <div className="space-y-4">
              <div className="flex items-center space-x-4">
                {selectedUser.picture_url && (
                  <div className="avatar">
                    <div className="w-24 rounded-full">
                      <img
                        src={selectedUser.picture_url}
                        alt={selectedUser.name}
                      />
                    </div>
                  </div>
                )}
                <div>
                  <h4 className="text-xl font-semibold">{selectedUser.name}</h4>
                  <p className="text-base-content/60">{selectedUser.email}</p>
                </div>
              </div>

              <div className="grid grid-cols-2 gap-4">
                <div>
                  <label className="label">
                    <span className="label-text font-semibold">
                      {t('users.details.id')}
                    </span>
                  </label>
                  <p>{selectedUser.id}</p>
                </div>
                <div>
                  <label className="label">
                    <span className="label-text font-semibold">
                      {t('users.details.status')}
                    </span>
                  </label>
                  <span
                    className={`badge ${getStatusBadgeClass(selectedUser.account_status || '')}`}
                  >
                    {selectedUser.account_status}
                  </span>
                </div>
                <div>
                  <label className="label">
                    <span className="label-text font-semibold">
                      {t('users.details.role')}
                    </span>
                  </label>
                  {selectedUser.role ? (
                    <div>
                      <span className="badge badge-secondary">
                        {selectedUser.role.display_name}
                      </span>
                      <p className="text-sm opacity-60 mt-1">
                        {getRoleDescription(selectedUser.role?.name || '')}
                      </p>
                      {selectedUser.role.priority && (
                        <p className="text-xs opacity-50 mt-1">
                          {t('users.roleTooltip.priority')}:{' '}
                          {selectedUser.role.priority}
                        </p>
                      )}
                    </div>
                  ) : (
                    <span className="badge badge-ghost">
                      {t('users.roleTooltip.noRole')}
                    </span>
                  )}
                </div>
                <div>
                  <label className="label">
                    <span className="label-text font-semibold">
                      {t('users.details.provider')}
                    </span>
                  </label>
                  <p>
                    {selectedUser.google_id
                      ? t('users.providers.google')
                      : t('users.providers.email')}
                  </p>
                </div>
                <div>
                  <label className="label">
                    <span className="label-text font-semibold">
                      {t('users.details.phone')}
                    </span>
                  </label>
                  <p>{selectedUser.phone || '-'}</p>
                </div>
                <div>
                  <label className="label">
                    <span className="label-text font-semibold">
                      {t('users.details.location')}
                    </span>
                  </label>
                  <p>
                    {selectedUser.city || '-'}
                    {selectedUser.country && `, ${selectedUser.country}`}
                  </p>
                </div>
                <div>
                  <label className="label">
                    <span className="label-text font-semibold">
                      {t('users.details.timezone')}
                    </span>
                  </label>
                  <p>{selectedUser.timezone || t('users.timezone.default')}</p>
                </div>
                <div>
                  <label className="label">
                    <span className="label-text font-semibold">
                      {t('users.details.created')}
                    </span>
                  </label>
                  <p>{formatDate(selectedUser.created_at)}</p>
                </div>
                <div>
                  <label className="label">
                    <span className="label-text font-semibold">
                      {t('users.details.lastSeen')}
                    </span>
                  </label>
                  <p>{formatDate(selectedUser.last_seen)}</p>
                </div>
              </div>

              {selectedUser.bio && (
                <div>
                  <label className="label">
                    <span className="label-text font-semibold">
                      {t('users.details.bio')}
                    </span>
                  </label>
                  <p className="text-base-content/80">{selectedUser.bio}</p>
                </div>
              )}

              <div>
                <label className="label">
                  <span className="label-text font-semibold">
                    {t('users.details.notifications')}
                  </span>
                </label>
                <div className="form-control">
                  <label className="label cursor-pointer justify-start gap-4">
                    <input
                      type="checkbox"
                      className="checkbox"
                      checked={selectedUser.notification_email}
                      disabled
                    />
                    <span className="label-text">
                      {t('users.details.emailNotifications')}
                    </span>
                  </label>
                </div>
              </div>

              {selectedUser.is_admin && (
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
                    />
                  </svg>
                  <span>{t('users.details.isAdmin')}</span>
                </div>
              )}
            </div>

            <div className="modal-action">
              <button className="btn" onClick={() => setIsModalOpen(false)}>
                {tCommon('close')}
              </button>
            </div>
          </div>
          <form method="dialog" className="modal-backdrop">
            <button onClick={() => setIsModalOpen(false)}>close</button>
          </form>
        </dialog>
      )}

      {/* Delete Confirmation Modal */}
      {isDeleteModalOpen && userToDelete && (
        <dialog className="modal modal-open">
          <div className="modal-box">
            <h3 className="font-bold text-lg">
              {t('users.deleteConfirm.title')}
            </h3>
            <p className="py-4">
              {t('users.deleteConfirm.message', {
                name: userToDelete.name || '',
              })}
            </p>
            <div className="alert alert-warning">
              <svg
                xmlns="http://www.w3.org/2000/svg"
                className="stroke-current shrink-0 h-6 w-6"
                fill="none"
                viewBox="0 0 24 24"
              >
                <path
                  strokeLinecap="round"
                  strokeLinejoin="round"
                  strokeWidth="2"
                  d="M12 9v2m0 4h.01m-6.938 4h13.856c1.54 0 2.502-1.667 1.732-3L13.732 4c-.77-1.333-2.694-1.333-3.464 0L3.34 16c-.77 1.333.192 3 1.732 3z"
                />
              </svg>
              <span>{t('users.deleteConfirm.warning')}</span>
            </div>
            <div className="modal-action">
              <button
                className="btn"
                onClick={() => {
                  setIsDeleteModalOpen(false);
                  setUserToDelete(null);
                }}
                disabled={deleting === userToDelete.id}
              >
                {tCommon('cancel')}
              </button>
              <button
                className={`btn btn-error ${
                  deleting === userToDelete.id ? 'loading' : ''
                }`}
                onClick={handleDeleteUser}
                disabled={deleting === userToDelete.id}
              >
                {tCommon('delete')}
              </button>
            </div>
          </div>
          <form method="dialog" className="modal-backdrop">
            <button
              onClick={() => {
                setIsDeleteModalOpen(false);
                setUserToDelete(null);
              }}
            >
              close
            </button>
          </form>
        </dialog>
      )}

      {error && (
        <div className="alert alert-error mt-4">
          <svg
            xmlns="http://www.w3.org/2000/svg"
            className="stroke-current shrink-0 h-6 w-6"
            fill="none"
            viewBox="0 0 24 24"
          >
            <path
              strokeLinecap="round"
              strokeLinejoin="round"
              strokeWidth="2"
              d="M10 14l2-2m0 0l2-2m-2 2l-2-2m2 2l2 2m7-2a9 9 0 11-18 0 9 9 0 0118 0z"
            />
          </svg>
          <span>{error}</span>
        </div>
      )}
    </div>
  );
}
