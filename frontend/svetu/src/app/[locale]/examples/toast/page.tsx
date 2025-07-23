'use client';

import { toast } from '@/utils/toast';

export default function ToastExamplesPage() {
  return (
    <div className="container mx-auto p-6 max-w-4xl">
      <h1 className="text-3xl font-bold mb-8">Toast Notifications</h1>
      
      <div className="space-y-8">
        {/* Типы уведомлений */}
        <section className="card bg-base-200 p-6">
          <h2 className="text-2xl font-semibold mb-4">Типы уведомлений</h2>
          <div className="grid grid-cols-1 md:grid-cols-2 gap-4">
            <button
              className="btn btn-success"
              onClick={() => toast.success('Операция выполнена успешно!')}
            >
              Success Toast
            </button>
            
            <button
              className="btn btn-error"
              onClick={() => toast.error('Произошла ошибка при сохранении')}
            >
              Error Toast
            </button>
            
            <button
              className="btn btn-warning"
              onClick={() => toast.warning('Внимание! Проверьте введенные данные')}
            >
              Warning Toast
            </button>
            
            <button
              className="btn btn-info"
              onClick={() => toast.info('Новое обновление доступно')}
            >
              Info Toast
            </button>
          </div>
        </section>

        {/* Позиционирование */}
        <section className="card bg-base-200 p-6">
          <h2 className="text-2xl font-semibold mb-4">Позиционирование</h2>
          <div className="grid grid-cols-2 md:grid-cols-3 gap-4">
            <button
              className="btn btn-primary btn-sm"
              onClick={() => 
                toast.success('Top position', { position: 'top' })
              }
            >
              Top
            </button>
            
            <button
              className="btn btn-primary btn-sm"
              onClick={() => 
                toast.success('Top Center', { position: 'top-center' })
              }
            >
              Top Center
            </button>
            
            <button
              className="btn btn-primary btn-sm"
              onClick={() => 
                toast.success('Top Start', { position: 'top-start' })
              }
            >
              Top Start
            </button>
            
            <button
              className="btn btn-primary btn-sm"
              onClick={() => 
                toast.success('Top End', { position: 'top-end' })
              }
            >
              Top End
            </button>
            
            <button
              className="btn btn-primary btn-sm"
              onClick={() => 
                toast.info('Middle position', { position: 'middle' })
              }
            >
              Middle
            </button>
            
            <button
              className="btn btn-primary btn-sm"
              onClick={() => 
                toast.info('Bottom position', { position: 'bottom' })
              }
            >
              Bottom
            </button>
            
            <button
              className="btn btn-primary btn-sm"
              onClick={() => 
                toast.info('Bottom Start', { position: 'bottom-start' })
              }
            >
              Bottom Start
            </button>
            
            <button
              className="btn btn-primary btn-sm"
              onClick={() => 
                toast.info('Bottom End', { position: 'bottom-end' })
              }
            >
              Bottom End
            </button>
          </div>
        </section>

        {/* Длительность */}
        <section className="card bg-base-200 p-6">
          <h2 className="text-2xl font-semibold mb-4">Длительность отображения</h2>
          <div className="grid grid-cols-1 md:grid-cols-3 gap-4">
            <button
              className="btn btn-secondary"
              onClick={() => 
                toast.info('Быстрое уведомление (2 сек)', { duration: 2000 })
              }
            >
              2 секунды
            </button>
            
            <button
              className="btn btn-secondary"
              onClick={() => 
                toast.info('Стандартное уведомление (3 сек)', { duration: 3000 })
              }
            >
              3 секунды (по умолчанию)
            </button>
            
            <button
              className="btn btn-secondary"
              onClick={() => 
                toast.warning('Долгое уведомление (5 сек)', { duration: 5000 })
              }
            >
              5 секунд
            </button>
          </div>
        </section>

        {/* Реальные примеры */}
        <section className="card bg-base-200 p-6">
          <h2 className="text-2xl font-semibold mb-4">Примеры использования</h2>
          <div className="space-y-4">
            <button
              className="btn btn-primary"
              onClick={() => {
                // Имитация сохранения
                toast.info('Сохранение...');
                setTimeout(() => {
                  toast.success('Изменения сохранены!');
                }, 1500);
              }}
            >
              Сохранить изменения
            </button>
            
            <button
              className="btn btn-error"
              onClick={() => {
                toast.error('Ошибка подключения к серверу', {
                  duration: 5000,
                  position: 'top-end'
                });
              }}
            >
              Ошибка сети
            </button>
            
            <button
              className="btn btn-warning"
              onClick={() => {
                toast.warning('Пожалуйста, заполните все обязательные поля', {
                  position: 'top-center'
                });
              }}
            >
              Валидация формы
            </button>
            
            <button
              className="btn"
              onClick={() => {
                toast.info('Файл загружается...', { duration: 2000 });
                setTimeout(() => {
                  toast.success('Файл успешно загружен!', {
                    position: 'bottom-end'
                  });
                }, 2000);
              }}
            >
              Загрузка файла
            </button>
          </div>
        </section>

        {/* Код примеры */}
        <section className="card bg-base-200 p-6">
          <h2 className="text-2xl font-semibold mb-4">Примеры кода</h2>
          <div className="mockup-code">
            <pre data-prefix="1"><code>{`import { toast } from '@/utils/toast';`}</code></pre>
            <pre data-prefix="2"><code>{``}</code></pre>
            <pre data-prefix="3"><code>{`// Базовое использование`}</code></pre>
            <pre data-prefix="4"><code>{`toast.success('Успешно!');`}</code></pre>
            <pre data-prefix="5"><code>{`toast.error('Ошибка!');`}</code></pre>
            <pre data-prefix="6"><code>{`toast.warning('Внимание!');`}</code></pre>
            <pre data-prefix="7"><code>{`toast.info('Информация');`}</code></pre>
            <pre data-prefix="8"><code>{``}</code></pre>
            <pre data-prefix="9"><code>{`// С опциями`}</code></pre>
            <pre data-prefix="10"><code>{`toast.show('Сообщение', {`}</code></pre>
            <pre data-prefix="11"><code>{`  type: 'success',`}</code></pre>
            <pre data-prefix="12"><code>{`  duration: 5000,`}</code></pre>
            <pre data-prefix="13"><code>{`  position: 'top-end'`}</code></pre>
            <pre data-prefix="14"><code>{`});`}</code></pre>
          </div>
        </section>

        {/* Особенности */}
        <section className="card bg-base-200 p-6">
          <h2 className="text-2xl font-semibold mb-4">Особенности</h2>
          <ul className="list-disc list-inside space-y-2">
            <li>Иконки для каждого типа уведомления</li>
            <li>Анимированное появление и исчезновение</li>
            <li>Закрытие по клику</li>
            <li>Автоматическое удаление контейнера когда нет уведомлений</li>
            <li>Поддержка множественных уведомлений</li>
            <li>Настраиваемая позиция и длительность</li>
            <li>Использует DaisyUI стили для консистентности</li>
          </ul>
        </section>
      </div>
    </div>
  );
}