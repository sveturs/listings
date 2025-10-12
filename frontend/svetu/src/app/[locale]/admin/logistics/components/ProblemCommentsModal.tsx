'use client';

import { useState, useEffect, useCallback } from 'react';
import { useTranslations } from 'next-intl';
import { apiClient } from '@/services/api-client';
import {
  FiX,
  FiMessageSquare,
  FiUser,
  FiClock,
  FiCheckCircle,
} from 'react-icons/fi';

interface Comment {
  id: number;
  admin_id: number;
  comment: string;
  comment_type: string;
  metadata: any;
  created_at: string;
  admin_name?: string;
  admin_email?: string;
}

interface ProblemCommentsModalProps {
  isOpen: boolean;
  onClose: () => void;
  problemId: number;
  problemDescription: string;
}

export default function ProblemCommentsModal({
  isOpen,
  onClose,
  problemId,
  problemDescription,
}: ProblemCommentsModalProps) {
  const t = useTranslations('admin');
  const [comments, setComments] = useState<Comment[]>([]);
  const [newComment, setNewComment] = useState('');
  const [loading, setLoading] = useState(false);
  const [submitting, setSubmitting] = useState(false);

  // Загрузка комментариев при открытии модального окна
  const loadComments = useCallback(async () => {
    setLoading(true);
    try {
      const response = await apiClient.get(
        `/admin/logistics/problems/${problemId}/comments`
      );
      setComments(response.data || []);
    } catch (error) {
      console.error('Error loading comments:', error);
    }
    setLoading(false);
  }, [problemId]);

  useEffect(() => {
    if (isOpen) {
      loadComments();
    }
  }, [isOpen, loadComments]);

  const handleAddComment = async () => {
    if (!newComment.trim()) return;

    setSubmitting(true);
    try {
      await apiClient.post(`/admin/logistics/problems/${problemId}/comments`, {
        comment: newComment,
        comment_type: 'comment',
      });

      setNewComment('');
      loadComments(); // Перезагружаем комментарии
    } catch (error) {
      console.error('Error adding comment:', error);
    }
    setSubmitting(false);
  };

  const formatDate = (dateString: string) => {
    return new Date(dateString).toLocaleString();
  };

  const getCommentTypeIcon = (type: string) => {
    switch (type) {
      case 'status_change':
        return <FiClock className="w-4 h-4 text-info" />;
      case 'assignment':
        return <FiUser className="w-4 h-4 text-warning" />;
      case 'resolution':
        return <FiCheckCircle className="w-4 h-4 text-success" />;
      default:
        return <FiMessageSquare className="w-4 h-4 text-base-content" />;
    }
  };

  const getCommentTypeBadge = (type: string) => {
    switch (type) {
      case 'status_change':
        return <span className="badge badge-info badge-xs">Status Change</span>;
      case 'assignment':
        return <span className="badge badge-warning badge-xs">Assignment</span>;
      case 'resolution':
        return <span className="badge badge-success badge-xs">Resolution</span>;
      default:
        return <span className="badge badge-ghost badge-xs">Comment</span>;
    }
  };

  if (!isOpen) return null;

  return (
    <div className="modal modal-open">
      <div className="modal-box max-w-4xl">
        <div className="flex justify-between items-center mb-4">
          <h3 className="font-bold text-lg flex items-center gap-2">
            <FiMessageSquare />
            {t('logistics.problems.problemComments.title')}
          </h3>
          <button className="btn btn-sm btn-circle btn-ghost" onClick={onClose}>
            <FiX />
          </button>
        </div>

        {/* Описание проблемы */}
        <div className="bg-base-200 p-4 rounded-lg mb-4">
          <h4 className="font-semibold text-sm text-base-content/70 mb-2">
            {t('logistics.problems.problemComments.problemDescription')}
          </h4>
          <p className="text-sm">{problemDescription}</p>
        </div>

        {/* Список комментариев */}
        <div className="max-h-96 overflow-y-auto mb-4">
          {loading ? (
            <div className="flex justify-center py-8">
              <span className="loading loading-spinner loading-md"></span>
            </div>
          ) : comments.length === 0 ? (
            <div className="text-center py-8 text-base-content/50">
              <FiMessageSquare className="w-12 h-12 mx-auto mb-2 opacity-50" />
              <p>{t('logistics.problems.problemComments.noComments')}</p>
            </div>
          ) : (
            <div className="space-y-4">
              {comments.map((comment) => (
                <div
                  key={comment.id}
                  className="bg-base-100 p-4 rounded-lg border"
                >
                  <div className="flex justify-between items-start mb-2">
                    <div className="flex items-center gap-2">
                      {getCommentTypeIcon(comment.comment_type)}
                      <span className="font-medium text-sm">
                        {comment.admin_name || `Admin ${comment.admin_id}`}
                      </span>
                      {getCommentTypeBadge(comment.comment_type)}
                    </div>
                    <span className="text-xs text-base-content/50">
                      {formatDate(comment.created_at)}
                    </span>
                  </div>
                  <p className="text-sm whitespace-pre-wrap">
                    {comment.comment}
                  </p>

                  {/* Метаданные для системных комментариев */}
                  {comment.metadata &&
                    Object.keys(comment.metadata).length > 0 && (
                      <div className="mt-2 text-xs bg-base-200 p-2 rounded">
                        <details>
                          <summary className="cursor-pointer">Metadata</summary>
                          <pre className="mt-1 text-xs">
                            {JSON.stringify(comment.metadata, null, 2)}
                          </pre>
                        </details>
                      </div>
                    )}
                </div>
              ))}
            </div>
          )}
        </div>

        {/* Форма добавления комментария */}
        <div className="border-t pt-4">
          <h4 className="font-semibold mb-2">
            {t('logistics.problems.problemComments.addComment')}
          </h4>
          <textarea
            className="textarea textarea-bordered w-full"
            rows={3}
            value={newComment}
            onChange={(e) => setNewComment(e.target.value)}
            placeholder={t(
              'logistics.problems.problemComments.commentPlaceholder'
            )}
            disabled={submitting}
          />
          <div className="flex justify-end gap-2 mt-2">
            <button className="btn btn-ghost btn-sm" onClick={onClose}>
              {t('common.cancel')}
            </button>
            <button
              className="btn btn-primary btn-sm"
              onClick={handleAddComment}
              disabled={!newComment.trim() || submitting}
            >
              {submitting && (
                <span className="loading loading-spinner loading-xs"></span>
              )}
              {t('logistics.problems.problemComments.addComment')}
            </button>
          </div>
        </div>
      </div>
      <div className="modal-backdrop" onClick={onClose}></div>
    </div>
  );
}
