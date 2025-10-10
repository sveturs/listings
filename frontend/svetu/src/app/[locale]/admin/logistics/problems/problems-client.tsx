'use client';

import { useState, useEffect, useCallback } from 'react';
import { useTranslations } from 'next-intl';
import { useRouter } from 'next/navigation';
import {
  FiAlertTriangle,
  FiClock,
  FiRotateCcw,
  FiX,
  FiEye,
  FiCheckCircle,
  FiUser,
  FiMessageSquare,
} from 'react-icons/fi';
import { apiClientAuth } from '@/lib/api-client-auth';
import AssignProblemModal from '../components/AssignProblemModal';
import ResolveProblemModal from '../components/ResolveProblemModal';
import ProblemCommentsModal from '../components/ProblemCommentsModal';

interface ProblemShipment {
  id: number;
  shipment_id: number;
  shipment_type: string;
  problem_type: 'delayed' | 'lost' | 'return' | 'complaint';
  description: string;
  status: 'open' | 'investigating' | 'resolved';
  assigned_to?: number;
  resolution?: string;
  created_at: string;
  resolved_at?: string;
  tracking_number?: string;
  recipient_name?: string;
  recipient_city?: string;
}

export default function ProblemShipmentsClient() {
  const t = useTranslations('admin');
  const router = useRouter();
  const [problems, setProblems] = useState<ProblemShipment[]>([]);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState<string | null>(null);
  const [filterType, setFilterType] = useState<string>('all');
  const [filterStatus, setFilterStatus] = useState<string>('all');
  const [assignModalOpen, setAssignModalOpen] = useState(false);
  const [resolveModalOpen, setResolveModalOpen] = useState(false);
  const [commentsModalOpen, setCommentsModalOpen] = useState(false);
  const [selectedProblem, setSelectedProblem] =
    useState<ProblemShipment | null>(null);

  const fetchProblems = useCallback(async () => {
    try {
      setLoading(true);
      setError(null);

      const params = new URLSearchParams();
      if (filterType !== 'all') params.append('problem_type', filterType);
      if (filterStatus !== 'all') params.append('status', filterStatus);

      const result = await apiClientAuth.get(
        `/admin/logistics/problems?${params}`
      );

      if (result.data) {
        setProblems(result.data.problems || []);
      } else {
        throw new Error(result.error?.message || 'Failed to fetch problems');
      }
    } catch (err) {
      console.error('Error fetching problems:', err);
      setError('Failed to load problems');
      // При ошибке используем пустой массив
      setProblems([]);
    } finally {
      setLoading(false);
    }
  }, [filterType, filterStatus]);

  useEffect(() => {
    fetchProblems();
  }, [fetchProblems]);

  const getProblemIcon = (type: string) => {
    switch (type) {
      case 'delayed':
        return <FiClock className="w-5 h-5 text-yellow-500" />;
      case 'lost':
        return <FiX className="w-5 h-5 text-red-500" />;
      case 'return':
        return <FiRotateCcw className="w-5 h-5 text-orange-500" />;
      case 'complaint':
        return <FiMessageSquare className="w-5 h-5 text-purple-500" />;
      default:
        return <FiAlertTriangle className="w-5 h-5 text-gray-500" />;
    }
  };

  const getStatusBadge = (status: string) => {
    const statusClasses = {
      open: 'badge badge-error',
      investigating: 'badge badge-warning',
      resolved: 'badge badge-success',
    };

    const statusLabels = {
      open: t('logistics.problems.status.open'),
      investigating: t('logistics.problems.status.investigating'),
      resolved: t('logistics.problems.status.resolved'),
    };

    return (
      <span
        className={
          statusClasses[status as keyof typeof statusClasses] || 'badge'
        }
      >
        {statusLabels[status as keyof typeof statusLabels] || status}
      </span>
    );
  };

  const getProblemTypeLabel = (type: string) => {
    const labels = {
      delayed: t('logistics.problems.types.delayed'),
      lost: t('logistics.problems.types.lost'),
      return: t('logistics.problems.types.return'),
      complaint: t('logistics.problems.types.complaint'),
    };
    return labels[type as keyof typeof labels] || type;
  };

  const handleViewShipment = (problem: ProblemShipment) => {
    router.push(
      `/admin/logistics/shipments/${problem.shipment_type}/${problem.shipment_id}`
    );
  };

  const handleResolveProblem = (problem: ProblemShipment) => {
    setSelectedProblem(problem);
    setResolveModalOpen(true);
  };

  const handleAssignProblem = (problem: ProblemShipment) => {
    setSelectedProblem(problem);
    setAssignModalOpen(true);
  };

  const onResolveProblem = async (resolution: string) => {
    if (!selectedProblem) return;

    try {
      const result = await apiClientAuth.post(
        `/admin/logistics/problems/${selectedProblem.id}/resolve`,
        {
          resolution: resolution,
        }
      );

      if (result.data) {
        fetchProblems();
      } else {
        setError(result.error?.message || 'Failed to resolve problem');
      }
    } catch (err) {
      console.error('Error resolving problem:', err);
      setError('Failed to resolve problem');
    }
  };

  const onShowComments = (problem: ProblemShipment) => {
    setSelectedProblem(problem);
    setCommentsModalOpen(true);
  };

  const onAssignProblem = async (adminId: number) => {
    if (!selectedProblem) return;

    try {
      const result = await apiClientAuth.post(
        `/admin/logistics/problems/${selectedProblem.id}/assign`,
        {
          assign_to: adminId,
        }
      );

      if (result.data) {
        fetchProblems();
      } else {
        setError(result.error?.message || 'Failed to assign problem');
      }
    } catch (err) {
      console.error('Error assigning problem:', err);
      setError('Failed to assign problem');
    }
  };

  if (loading) {
    return (
      <div className="flex justify-center items-center h-64">
        <span className="loading loading-spinner loading-lg"></span>
      </div>
    );
  }

  if (error) {
    return (
      <div className="alert alert-error">
        <FiAlertTriangle className="w-5 h-5" />
        <span>{error}</span>
      </div>
    );
  }

  return (
    <div className="space-y-6">
      {/* Фильтры */}
      <div className="card bg-base-100 shadow-sm">
        <div className="card-body">
          <h2 className="card-title">{t('logistics.filters.title')}</h2>
          <div className="grid grid-cols-1 md:grid-cols-2 gap-4">
            <div className="form-control">
              <label className="label">
                <span className="label-text">
                  {t('logistics.problems.problemType')}
                </span>
              </label>
              <select
                className="select select-bordered"
                value={filterType}
                onChange={(e) => setFilterType(e.target.value)}
              >
                <option value="all">{t('logistics.problems.allTypes')}</option>
                <option value="delayed">
                  {t('logistics.problems.types.delayed')}
                </option>
                <option value="lost">
                  {t('logistics.problems.types.lost')}
                </option>
                <option value="return">
                  {t('logistics.problems.types.return')}
                </option>
                <option value="complaint">
                  {t('logistics.problems.types.complaint')}
                </option>
              </select>
            </div>

            <div className="form-control">
              <label className="label">
                <span className="label-text">
                  {t('logistics.filters.status')}
                </span>
              </label>
              <select
                className="select select-bordered"
                value={filterStatus}
                onChange={(e) => setFilterStatus(e.target.value)}
              >
                <option value="all">
                  {t('logistics.filters.all_statuses')}
                </option>
                <option value="open">
                  {t('logistics.problems.status.open')}
                </option>
                <option value="investigating">
                  {t('logistics.problems.status.investigating')}
                </option>
                <option value="resolved">
                  {t('logistics.problems.status.resolved')}
                </option>
              </select>
            </div>
          </div>
        </div>
      </div>

      {/* Статистика */}
      <div className="stats shadow w-full">
        <div className="stat">
          <div className="stat-figure text-error">
            <FiAlertTriangle className="w-8 h-8" />
          </div>
          <div className="stat-title">{t('logistics.problems.stats.open')}</div>
          <div className="stat-value text-error">
            {problems.filter((p) => p.status === 'open').length}
          </div>
        </div>

        <div className="stat">
          <div className="stat-figure text-warning">
            <FiClock className="w-8 h-8" />
          </div>
          <div className="stat-title">
            {t('logistics.problems.stats.investigating')}
          </div>
          <div className="stat-value text-warning">
            {problems.filter((p) => p.status === 'investigating').length}
          </div>
        </div>

        <div className="stat">
          <div className="stat-figure text-success">
            <FiCheckCircle className="w-8 h-8" />
          </div>
          <div className="stat-title">
            {t('logistics.problems.stats.resolved')}
          </div>
          <div className="stat-value text-success">
            {problems.filter((p) => p.status === 'resolved').length}
          </div>
        </div>

        <div className="stat">
          <div className="stat-figure text-primary">
            <FiRotateCcw className="w-8 h-8" />
          </div>
          <div className="stat-title">
            {t('logistics.problems.stats.returns')}
          </div>
          <div className="stat-value text-primary">
            {problems.filter((p) => p.problem_type === 'return').length}
          </div>
        </div>
      </div>

      {/* Список проблем */}
      <div className="card bg-base-100 shadow-sm">
        <div className="card-body">
          <h2 className="card-title">
            {t('logistics.problems.list')} ({problems.length})
          </h2>

          {problems.length === 0 ? (
            <div className="text-center py-8 text-gray-500">
              {t('logistics.problems.noProblems')}
            </div>
          ) : (
            <div className="overflow-x-auto">
              <table className="table">
                <thead>
                  <tr>
                    <th>{t('logistics.problems.type')}</th>
                    <th>{t('logistics.table.tracking_number')}</th>
                    <th>{t('logistics.table.recipient')}</th>
                    <th>{t('logistics.problems.description')}</th>
                    <th>{t('logistics.filters.status')}</th>
                    <th>{t('logistics.problems.assignedTo')}</th>
                    <th>{t('logistics.table.created')}</th>
                    <th>{t('logistics.table.actions')}</th>
                  </tr>
                </thead>
                <tbody>
                  {problems.map((problem) => (
                    <tr key={problem.id} className="hover">
                      <td>
                        <div className="flex items-center gap-2">
                          {getProblemIcon(problem.problem_type)}
                          <span>
                            {getProblemTypeLabel(problem.problem_type)}
                          </span>
                        </div>
                      </td>
                      <td>
                        <div className="font-mono text-sm">
                          {problem.tracking_number || 'N/A'}
                        </div>
                        <div className="text-xs text-gray-500">
                          {problem.shipment_type}
                        </div>
                      </td>
                      <td>
                        <div>{problem.recipient_name || 'N/A'}</div>
                        <div className="text-xs text-gray-500">
                          {problem.recipient_city || 'N/A'}
                        </div>
                      </td>
                      <td>
                        <div
                          className="max-w-xs truncate"
                          title={problem.description}
                        >
                          {problem.description}
                        </div>
                      </td>
                      <td>{getStatusBadge(problem.status)}</td>
                      <td>
                        {problem.assigned_to ? (
                          <div className="flex items-center gap-1">
                            <FiUser className="w-4 h-4" />
                            <span>ID: {problem.assigned_to}</span>
                          </div>
                        ) : (
                          <span className="text-gray-400">-</span>
                        )}
                      </td>
                      <td>
                        {new Date(problem.created_at).toLocaleDateString()}
                      </td>
                      <td>
                        <div className="flex gap-2">
                          <button
                            className="btn btn-sm btn-ghost"
                            onClick={() => handleViewShipment(problem)}
                            title={t('logistics.actions.view')}
                          >
                            <FiEye />
                          </button>
                          <button
                            className="btn btn-sm btn-ghost"
                            onClick={() => onShowComments(problem)}
                            title={t('logistics.actions.comments')}
                          >
                            <FiMessageSquare />
                          </button>
                          {problem.status !== 'resolved' && (
                            <>
                              <button
                                className="btn btn-sm btn-ghost"
                                onClick={() => handleAssignProblem(problem)}
                                title={t('logistics.actions.assign')}
                              >
                                <FiUser />
                              </button>
                              <button
                                className="btn btn-sm btn-ghost text-success"
                                onClick={() => handleResolveProblem(problem)}
                                title={t('logistics.actions.resolve')}
                              >
                                <FiCheckCircle />
                              </button>
                            </>
                          )}
                        </div>
                      </td>
                    </tr>
                  ))}
                </tbody>
              </table>
            </div>
          )}
        </div>
      </div>

      {/* Модальные окна */}
      {selectedProblem && (
        <>
          <AssignProblemModal
            isOpen={assignModalOpen}
            onClose={() => setAssignModalOpen(false)}
            onAssign={onAssignProblem}
            problemId={selectedProblem.id}
            currentAssignedTo={selectedProblem.assigned_to}
          />
          <ResolveProblemModal
            isOpen={resolveModalOpen}
            onClose={() => setResolveModalOpen(false)}
            onResolve={onResolveProblem}
            problemId={selectedProblem.id}
            problemDescription={selectedProblem.description}
          />
          <ProblemCommentsModal
            isOpen={commentsModalOpen}
            onClose={() => setCommentsModalOpen(false)}
            problemId={selectedProblem.id}
            problemDescription={selectedProblem.description}
          />
        </>
      )}
    </div>
  );
}
