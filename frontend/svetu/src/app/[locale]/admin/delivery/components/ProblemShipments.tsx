'use client';

import { useEffect, useState } from 'react';
import { useTranslations } from 'next-intl';
import { tokenManager } from '@/utils/tokenManager';

interface ProblemShipment {
  id: number;
  tracking_number: string;
  provider_name: string;
  problem_type: string;
  priority: string;
  status: string;
  assigned_to?: string;
  created_at: string;
  customer_name: string;
  customer_phone: string;
  description: string;
}

export default function ProblemShipments() {
  const t = useTranslations('admin.delivery.problemShipments');
  const [problems, setProblems] = useState<ProblemShipment[]>([]);
  const [loading, setLoading] = useState(true);
  const [filterType, setFilterType] = useState('all');
  const [filterStatus, setFilterStatus] = useState('all');
  const [selectedProblem, setSelectedProblem] = useState<ProblemShipment | null>(null);
  const [resolveModalOpen, setResolveModalOpen] = useState(false);

  useEffect(() => {
    fetchProblemShipments();
  }, [filterType, filterStatus]);

  const fetchProblemShipments = async () => {
    try {
      // Mock data
      setTimeout(() => {
        const mockData: ProblemShipment[] = [
          {
            id: 1,
            tracking_number: 'PE2024112301',
            provider_name: 'Post Express',
            problem_type: 'delayed',
            priority: 'high',
            status: 'investigating',
            assigned_to: 'Admin #2',
            created_at: '2024-11-23 14:30',
            customer_name: 'Марко Петровић',
            customer_phone: '+381 65 123 4567',
            description: 'Задержка доставки более 48 часов'
          },
          {
            id: 2,
            tracking_number: 'BEX2024112302',
            provider_name: 'BEX Express',
            problem_type: 'lost',
            priority: 'critical',
            status: 'open',
            created_at: '2024-11-23 10:15',
            customer_name: 'Ана Јовановић',
            customer_phone: '+381 64 987 6543',
            description: 'Посылка не найдена в системе отслеживания'
          },
          {
            id: 3,
            tracking_number: 'AKS2024112303',
            provider_name: 'AKS',
            problem_type: 'damaged',
            priority: 'medium',
            status: 'resolved',
            assigned_to: 'Admin #1',
            created_at: '2024-11-22 16:45',
            customer_name: 'Милан Стојановић',
            customer_phone: '+381 66 555 4444',
            description: 'Товар поврежден при транспортировке'
          },
          {
            id: 4,
            tracking_number: 'DE2024112304',
            provider_name: 'D Express',
            problem_type: 'wrongAddress',
            priority: 'low',
            status: 'investigating',
            assigned_to: 'Admin #3',
            created_at: '2024-11-23 09:00',
            customer_name: 'Јелена Николић',
            customer_phone: '+381 63 222 3333',
            description: 'Неверный адрес доставки'
          },
          {
            id: 5,
            tracking_number: 'CE2024112305',
            provider_name: 'City Express',
            problem_type: 'refused',
            priority: 'medium',
            status: 'closed',
            assigned_to: 'Admin #2',
            created_at: '2024-11-22 12:30',
            customer_name: 'Драган Ђорђевић',
            customer_phone: '+381 65 777 8888',
            description: 'Получатель отказался от посылки'
          }
        ];

        let filtered = mockData;
        if (filterType !== 'all') {
          filtered = filtered.filter(p => p.problem_type === filterType);
        }
        if (filterStatus !== 'all') {
          filtered = filtered.filter(p => p.status === filterStatus);
        }

        setProblems(filtered);
        setLoading(false);
      }, 500);
    } catch (error) {
      console.error('Failed to fetch problem shipments:', error);
      setLoading(false);
    }
  };

  const getPriorityBadge = (priority: string) => {
    const classes: Record<string, string> = {
      low: 'badge-ghost',
      medium: 'badge-warning',
      high: 'badge-error',
      critical: 'badge-error animate-pulse'
    };
    return `badge ${classes[priority] || 'badge-ghost'}`;
  };

  const getStatusBadge = (status: string) => {
    const classes: Record<string, string> = {
      open: 'badge-info',
      investigating: 'badge-warning',
      resolved: 'badge-success',
      closed: 'badge-ghost'
    };
    return `badge ${classes[status] || 'badge-ghost'}`;
  };

  const openResolveModal = (problem: ProblemShipment) => {
    setSelectedProblem(problem);
    setResolveModalOpen(true);
  };

  if (loading) {
    return (
      <div className="flex items-center justify-center h-64">
        <span className="loading loading-spinner loading-lg"></span>
      </div>
    );
  }

  return (
    <div>
      <div className="mb-6">
        <h2 className="text-xl font-semibold mb-2">{t('title')}</h2>

        {/* Stats Cards */}
        <div className="grid grid-cols-1 md:grid-cols-4 gap-4 mb-6">
          <div className="stat bg-base-100 rounded-lg shadow">
            <div className="stat-title text-sm">Открытые</div>
            <div className="stat-value text-2xl text-info">2</div>
          </div>
          <div className="stat bg-base-100 rounded-lg shadow">
            <div className="stat-title text-sm">В работе</div>
            <div className="stat-value text-2xl text-warning">2</div>
          </div>
          <div className="stat bg-base-100 rounded-lg shadow">
            <div className="stat-title text-sm">Решенные</div>
            <div className="stat-value text-2xl text-success">1</div>
          </div>
          <div className="stat bg-base-100 rounded-lg shadow">
            <div className="stat-title text-sm">Критические</div>
            <div className="stat-value text-2xl text-error">1</div>
          </div>
        </div>

        {/* Filters */}
        <div className="flex gap-4 mb-4">
          <select
            className="select select-bordered select-sm"
            value={filterType}
            onChange={(e) => setFilterType(e.target.value)}
          >
            <option value="all">{t('filters.allTypes')}</option>
            <option value="delayed">{t('filters.delayed')}</option>
            <option value="lost">{t('filters.lost')}</option>
            <option value="damaged">{t('filters.damaged')}</option>
            <option value="wrongAddress">{t('filters.wrongAddress')}</option>
            <option value="refused">{t('filters.refused')}</option>
          </select>

          <select
            className="select select-bordered select-sm"
            value={filterStatus}
            onChange={(e) => setFilterStatus(e.target.value)}
          >
            <option value="all">Все статусы</option>
            <option value="open">Открыто</option>
            <option value="investigating">Расследуется</option>
            <option value="resolved">Решено</option>
            <option value="closed">Закрыто</option>
          </select>
        </div>
      </div>

      <div className="overflow-x-auto">
        <table className="table table-zebra w-full">
          <thead>
            <tr>
              <th>{t('table.tracking')}</th>
              <th>{t('table.provider')}</th>
              <th>{t('table.problemType')}</th>
              <th>{t('table.priority')}</th>
              <th>{t('table.status')}</th>
              <th>{t('table.assignedTo')}</th>
              <th>{t('table.createdAt')}</th>
              <th>{t('table.actions')}</th>
            </tr>
          </thead>
          <tbody>
            {problems.map((problem) => (
              <tr key={problem.id}>
                <td>
                  <div>
                    <div className="font-mono text-sm">{problem.tracking_number}</div>
                    <div className="text-xs text-base-content/60">{problem.customer_name}</div>
                  </div>
                </td>
                <td>{problem.provider_name}</td>
                <td>{problem.problem_type}</td>
                <td>
                  <span className={getPriorityBadge(problem.priority)}>
                    {t(`priority.${problem.priority}`)}
                  </span>
                </td>
                <td>
                  <span className={getStatusBadge(problem.status)}>
                    {t(`status.${problem.status}`)}
                  </span>
                </td>
                <td>{problem.assigned_to || '-'}</td>
                <td className="text-sm">{problem.created_at}</td>
                <td>
                  <div className="dropdown dropdown-end">
                    <label tabIndex={0} className="btn btn-ghost btn-xs">
                      ⋮
                    </label>
                    <ul tabIndex={0} className="dropdown-content z-[1] menu p-2 shadow bg-base-100 rounded-box w-52">
                      <li><a>{t('actions.view')}</a></li>
                      <li><a>{t('actions.assign')}</a></li>
                      {problem.status !== 'resolved' && problem.status !== 'closed' && (
                        <li><a onClick={() => openResolveModal(problem)}>{t('actions.resolve')}</a></li>
                      )}
                      <li><a>{t('actions.escalate')}</a></li>
                      <li><a>{t('actions.contact')}</a></li>
                    </ul>
                  </div>
                </td>
              </tr>
            ))}
          </tbody>
        </table>
      </div>

      {/* Resolve Modal */}
      {resolveModalOpen && selectedProblem && (
        <dialog className="modal modal-open">
          <div className="modal-box">
            <h3 className="font-bold text-lg mb-4">
              Решить проблему: {selectedProblem.tracking_number}
            </h3>

            <div className="space-y-4">
              <div>
                <p className="text-sm text-base-content/70 mb-2">Описание проблемы:</p>
                <p className="p-3 bg-base-200 rounded">{selectedProblem.description}</p>
              </div>

              <div className="form-control">
                <label className="label">
                  <span className="label-text">Решение</span>
                </label>
                <textarea
                  className="textarea textarea-bordered h-32"
                  placeholder="Опишите, как была решена проблема..."
                />
              </div>

              <div className="form-control">
                <label className="label cursor-pointer">
                  <span className="label-text">Отправить уведомление клиенту</span>
                  <input type="checkbox" className="checkbox" defaultChecked />
                </label>
              </div>
            </div>

            <div className="modal-action">
              <button className="btn btn-success">Решить</button>
              <button className="btn" onClick={() => setResolveModalOpen(false)}>
                Отмена
              </button>
            </div>
          </div>
          <form method="dialog" className="modal-backdrop" onClick={() => setResolveModalOpen(false)}>
            <button>close</button>
          </form>
        </dialog>
      )}
    </div>
  );
}