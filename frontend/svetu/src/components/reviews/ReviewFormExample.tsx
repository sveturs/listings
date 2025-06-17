'use client';

import React, { useState } from 'react';
import { ReviewForm } from './ReviewForm';

/**
 * Пример использования обновленного компонента ReviewForm
 * с двухэтапным процессом создания отзывов
 */
export const ReviewFormExample: React.FC = () => {
  const [showForm, setShowForm] = useState(false);

  const handleSuccess = () => {
    console.log('Отзыв успешно создан и опубликован!');
    setShowForm(false);
    // Здесь можно добавить уведомление об успехе
  };

  const handleCancel = () => {
    setShowForm(false);
  };

  return (
    <div className="max-w-2xl mx-auto p-6">
      <h2 className="text-2xl font-bold mb-6">
        Пример использования ReviewForm
      </h2>

      {!showForm ? (
        <div className="space-y-4">
          <p className="text-base-content/70">
            Новый компонент ReviewForm поддерживает двухэтапный процесс создания
            отзывов:
          </p>
          <ul className="list-disc list-inside space-y-2 text-sm">
            <li>Этап 1: Создание черновика отзыва с текстовыми данными</li>
            <li>Этап 2a: Загрузка фотографий (если есть)</li>
            <li>Этап 2b: Публикация отзыва</li>
          </ul>

          <button onClick={() => setShowForm(true)} className="btn btn-primary">
            Открыть форму отзыва
          </button>

          <div className="divider">Новый API</div>

          <div className="mockup-code">
            <pre data-prefix="1">
              <code>{`<ReviewForm`}</code>
            </pre>
            <pre data-prefix="2">
              <code>{`  entityType="listing"`}</code>
            </pre>
            <pre data-prefix="3">
              <code>{`  entityId={123}`}</code>
            </pre>
            <pre data-prefix="4">
              <code>{`  storefrontId={456}`}</code>
            </pre>
            <pre data-prefix="5">
              <code>{`  onSuccess={() => console.log('Success!')}`}</code>
            </pre>
            <pre data-prefix="6">
              <code>{`  onCancel={() => setShowForm(false)}`}</code>
            </pre>
            <pre data-prefix="7">
              <code>{`/>`}</code>
            </pre>
          </div>

          <div className="divider">Обратная совместимость</div>

          <div className="mockup-code">
            <pre data-prefix="1">
              <code>{`<ReviewForm`}</code>
            </pre>
            <pre data-prefix="2">
              <code>{`  entityType="listing"`}</code>
            </pre>
            <pre data-prefix="3">
              <code>{`  entityId={123}`}</code>
            </pre>
            <pre data-prefix="4">
              <code>{`  legacyOnSubmit={handleOldSubmit}`}</code>
            </pre>
            <pre data-prefix="5">
              <code>{`  onCancel={() => setShowForm(false)}`}</code>
            </pre>
            <pre data-prefix="6">
              <code>{`/>`}</code>
            </pre>
          </div>
        </div>
      ) : (
        <div className="space-y-4">
          <h3 className="text-lg font-semibold">Создать отзыв</h3>
          <ReviewForm
            entityType="listing"
            entityId={123}
            storefrontId={456}
            onSuccess={handleSuccess}
            onCancel={handleCancel}
          />
        </div>
      )}
    </div>
  );
};
