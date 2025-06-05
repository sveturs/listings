'use client';

import { UploadingFile } from '@/types/chat';
import { FaFile, FaTimes } from 'react-icons/fa';

interface FileUploadProgressProps {
  uploadingFiles: UploadingFile[];
  onRemove: (fileId: string) => void;
}

export function FileUploadProgress({
  uploadingFiles,
  onRemove,
}: FileUploadProgressProps) {
  if (uploadingFiles.length === 0) return null;

  const formatFileSize = (bytes: number) => {
    if (bytes === 0) return '0 Bytes';
    const k = 1024;
    const sizes = ['Bytes', 'KB', 'MB', 'GB'];
    const i = Math.floor(Math.log(bytes) / Math.log(k));
    return parseFloat((bytes / Math.pow(k, i)).toFixed(2)) + ' ' + sizes[i];
  };

  const getStatusColor = (status: UploadingFile['status']) => {
    switch (status) {
      case 'pending':
        return 'text-base-content/60';
      case 'uploading':
        return 'text-primary';
      case 'success':
        return 'text-success';
      case 'error':
        return 'text-error';
    }
  };

  const getStatusText = (status: UploadingFile['status']) => {
    switch (status) {
      case 'pending':
        return 'Ожидание...';
      case 'uploading':
        return 'Загрузка...';
      case 'success':
        return 'Загружено';
      case 'error':
        return 'Ошибка';
    }
  };

  return (
    <div className="space-y-2 p-3 bg-base-200 rounded-lg">
      <p className="text-sm font-medium">Загрузка файлов:</p>
      {uploadingFiles.map((file) => (
        <div
          key={file.id}
          className="flex items-center gap-3 p-2 bg-base-100 rounded"
        >
          <FaFile className={`text-lg ${getStatusColor(file.status)}`} />
          <div className="flex-1 min-w-0">
            <div className="flex items-center justify-between">
              <p className="text-sm truncate pr-2">{file.file.name}</p>
              <p className="text-xs text-base-content/60">
                {formatFileSize(file.file.size)}
              </p>
            </div>
            <div className="mt-1">
              {file.status === 'uploading' ? (
                <div className="w-full bg-base-200 rounded-full h-2">
                  <div
                    className="bg-primary h-2 rounded-full transition-all duration-300"
                    style={{ width: `${file.progress}%` }}
                  />
                </div>
              ) : (
                <p className={`text-xs ${getStatusColor(file.status)}`}>
                  {file.error || getStatusText(file.status)}
                </p>
              )}
            </div>
          </div>
          {(file.status === 'error' || file.status === 'pending') && (
            <button
              onClick={() => onRemove(file.id)}
              className="btn btn-ghost btn-xs btn-circle"
            >
              <FaTimes />
            </button>
          )}
        </div>
      ))}
    </div>
  );
}
