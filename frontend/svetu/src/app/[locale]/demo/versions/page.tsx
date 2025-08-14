'use client';

import { useState } from 'react';
import VersionHistoryViewer from '@/components/admin/translations/VersionHistoryViewer';
import { ClockIcon } from '@heroicons/react/24/outline';

export default function DemoVersionsPage() {
  const [showVersionHistory, setShowVersionHistory] = useState(false);
  const [entityType, setEntityType] = useState('attribute');
  const [entityId, setEntityId] = useState(3200);

  return (
    <div className="container mx-auto px-4 py-8">
      <h1 className="text-3xl font-bold mb-4">
        Демо: История версий переводов
      </h1>

      <div className="card bg-base-100 shadow-xl">
        <div className="card-body">
          <h2 className="card-title">Тестирование версионирования</h2>
          <p className="text-base-content/60 mb-4">
            Эта страница демонстрирует работу системы версионирования переводов
          </p>

          <div className="form-control">
            <label className="label">
              <span className="label-text">Тип сущности</span>
            </label>
            <select
              className="select select-bordered w-full"
              value={entityType}
              onChange={(e) => setEntityType(e.target.value)}
            >
              <option value="attribute">Атрибут</option>
              <option value="listing">Объявление</option>
              <option value="category">Категория</option>
            </select>
          </div>

          <div className="form-control">
            <label className="label">
              <span className="label-text">ID сущности</span>
            </label>
            <input
              type="number"
              className="input input-bordered w-full"
              value={entityId}
              onChange={(e) => setEntityId(parseInt(e.target.value) || 0)}
            />
          </div>

          <div className="card-actions justify-end mt-4">
            <button
              className="btn btn-primary"
              onClick={() => setShowVersionHistory(true)}
            >
              <ClockIcon className="h-4 w-4 mr-2" />
              Показать историю версий
            </button>
          </div>

          <div className="alert alert-info mt-4">
            <div>
              <h5 className="font-medium">Тестовые данные:</h5>
              <p className="text-sm mt-1">
                attribute ID: 3200 - имеет 2 версии изменений
              </p>
            </div>
          </div>
        </div>
      </div>

      {showVersionHistory && (
        <VersionHistoryViewer
          entityType={entityType}
          entityId={entityId}
          onClose={() => setShowVersionHistory(false)}
        />
      )}
    </div>
  );
}
