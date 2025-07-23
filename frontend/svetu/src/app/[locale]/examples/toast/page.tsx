'use client';

import { toast } from '@/utils/toast';

export default function ToastExamplesPage() {
  const showAllPositions = () => {
    const positions = [
      'top',
      'bottom',
      'top-center',
      'top-start',
      'top-end',
      'middle',
      'bottom-start',
      'bottom-end',
    ] as const;

    positions.forEach((position, index) => {
      setTimeout(() => {
        toast.info(`Position: ${position}`, { position, duration: 5000 });
      }, index * 200);
    });
  };

  const showAllTypes = () => {
    toast.success('Operation completed successfully!', { position: 'top-end' });
    setTimeout(() => {
      toast.error('Something went wrong. Please try again.', {
        position: 'top-end',
      });
    }, 500);
    setTimeout(() => {
      toast.warning('Low disk space warning', { position: 'top-end' });
    }, 1000);
    setTimeout(() => {
      toast.info('New update available', { position: 'top-end' });
    }, 1500);
  };

  const showCustomDurations = () => {
    toast.info('Short toast (2s)', { duration: 2000 });
    toast.warning('Medium toast (5s)', { duration: 5000 });
    toast.error('Long toast (10s)', { duration: 10000 });
  };

  const showRealWorldExamples = () => {
    // Пример 1: Успешное сохранение
    toast.success('Profile updated successfully!', {
      position: 'top-center',
      duration: 3000,
    });

    // Пример 2: Ошибка валидации
    setTimeout(() => {
      toast.error('Please fill in all required fields', {
        position: 'top-center',
        duration: 4000,
      });
    }, 1000);

    // Пример 3: Информация о процессе
    setTimeout(() => {
      toast.info('Uploading files... This may take a moment', {
        position: 'bottom-end',
        duration: 5000,
      });
    }, 2000);

    // Пример 4: Предупреждение
    setTimeout(() => {
      toast.warning('Your session will expire in 5 minutes', {
        position: 'top-end',
        duration: 6000,
      });
    }, 3000);
  };

  return (
    <div className="container mx-auto p-4 max-w-4xl">
      <h1 className="text-3xl font-bold mb-8">Toast Notifications Examples</h1>

      <div className="space-y-8">
        {/* Типы уведомлений */}
        <section className="card bg-base-100 shadow-xl">
          <div className="card-body">
            <h2 className="card-title mb-4">Toast Types</h2>
            <p className="text-base-content/70 mb-4">
              Different types of notifications for various scenarios
            </p>
            <div className="flex flex-wrap gap-2">
              <button
                className="btn btn-success"
                onClick={() =>
                  toast.success('Success! Your changes have been saved.')
                }
              >
                Success Toast
              </button>
              <button
                className="btn btn-error"
                onClick={() => toast.error('Error! Something went wrong.')}
              >
                Error Toast
              </button>
              <button
                className="btn btn-warning"
                onClick={() =>
                  toast.warning('Warning! Please check your input.')
                }
              >
                Warning Toast
              </button>
              <button
                className="btn btn-info"
                onClick={() => toast.info('Info: New features available!')}
              >
                Info Toast
              </button>
            </div>
          </div>
        </section>

        {/* Позиционирование */}
        <section className="card bg-base-100 shadow-xl">
          <div className="card-body">
            <h2 className="card-title mb-4">Positioning</h2>
            <p className="text-base-content/70 mb-4">
              Show toasts at different positions on the screen
            </p>
            <div className="grid grid-cols-3 gap-2">
              <button
                className="btn btn-sm"
                onClick={() =>
                  toast.info('Top Start', { position: 'top-start' })
                }
              >
                Top Start
              </button>
              <button
                className="btn btn-sm"
                onClick={() =>
                  toast.info('Top Center', { position: 'top-center' })
                }
              >
                Top Center
              </button>
              <button
                className="btn btn-sm"
                onClick={() => toast.info('Top End', { position: 'top-end' })}
              >
                Top End
              </button>
              <button
                className="btn btn-sm"
                onClick={() => toast.info('Middle', { position: 'middle' })}
              >
                Middle
              </button>
              <button
                className="btn btn-sm"
                onClick={() => toast.info('Bottom', { position: 'bottom' })}
              >
                Bottom
              </button>
              <button
                className="btn btn-sm"
                onClick={() =>
                  toast.info('Bottom End', { position: 'bottom-end' })
                }
              >
                Bottom End
              </button>
            </div>
            <button className="btn btn-primary mt-4" onClick={showAllPositions}>
              Show All Positions
            </button>
          </div>
        </section>

        {/* Продолжительность */}
        <section className="card bg-base-100 shadow-xl">
          <div className="card-body">
            <h2 className="card-title mb-4">Duration</h2>
            <p className="text-base-content/70 mb-4">
              Control how long the toast stays visible
            </p>
            <div className="flex flex-wrap gap-2">
              <button
                className="btn"
                onClick={() => toast.info('Quick (1s)', { duration: 1000 })}
              >
                1 Second
              </button>
              <button
                className="btn"
                onClick={() => toast.info('Default (3s)', { duration: 3000 })}
              >
                3 Seconds
              </button>
              <button
                className="btn"
                onClick={() => toast.info('Long (7s)', { duration: 7000 })}
              >
                7 Seconds
              </button>
              <button className="btn btn-primary" onClick={showCustomDurations}>
                Show Different Durations
              </button>
            </div>
          </div>
        </section>

        {/* Интерактивность */}
        <section className="card bg-base-100 shadow-xl">
          <div className="card-body">
            <h2 className="card-title mb-4">Interactivity</h2>
            <p className="text-base-content/70 mb-4">
              Click on any toast to dismiss it immediately
            </p>
            <button
              className="btn btn-primary"
              onClick={() => {
                toast.info('Click me to dismiss!', { duration: 10000 });
                setTimeout(() => {
                  toast.warning(
                    "I'll stay for 10 seconds unless you click me",
                    { duration: 10000 }
                  );
                }, 500);
              }}
            >
              Show Clickable Toasts
            </button>
          </div>
        </section>

        {/* Реальные примеры */}
        <section className="card bg-base-100 shadow-xl">
          <div className="card-body">
            <h2 className="card-title mb-4">Real World Examples</h2>
            <p className="text-base-content/70 mb-4">
              Common use cases in the application
            </p>
            <div className="space-y-2">
              <button
                className="btn btn-primary"
                onClick={showRealWorldExamples}
              >
                Show Real Examples
              </button>
              <button className="btn" onClick={showAllTypes}>
                Show All Types at Once
              </button>
            </div>
          </div>
        </section>

        {/* Примеры кода */}
        <section className="card bg-base-100 shadow-xl">
          <div className="card-body">
            <h2 className="card-title mb-4">Code Examples</h2>
            <div className="mockup-code">
              <pre data-prefix="1">
                <code>{`import { toast } from '@/utils/toast';`}</code>
              </pre>
              <pre data-prefix="2">
                <code>{``}</code>
              </pre>
              <pre data-prefix="3">
                <code>{`// Basic usage`}</code>
              </pre>
              <pre data-prefix="4">
                <code>{`toast.success('Success message');`}</code>
              </pre>
              <pre data-prefix="5">
                <code>{`toast.error('Error message');`}</code>
              </pre>
              <pre data-prefix="6">
                <code>{`toast.warning('Warning message');`}</code>
              </pre>
              <pre data-prefix="7">
                <code>{`toast.info('Info message');`}</code>
              </pre>
              <pre data-prefix="8">
                <code>{``}</code>
              </pre>
              <pre data-prefix="9">
                <code>{`// With options`}</code>
              </pre>
              <pre data-prefix="10">
                <code>{`toast.success('Saved!', {`}</code>
              </pre>
              <pre data-prefix="11">
                <code>{`  position: 'top-end',`}</code>
              </pre>
              <pre data-prefix="12">
                <code>{`  duration: 5000`}</code>
              </pre>
              <pre data-prefix="13">
                <code>{`});`}</code>
              </pre>
            </div>
          </div>
        </section>

        {/* Интеграция с формами */}
        <section className="card bg-base-100 shadow-xl">
          <div className="card-body">
            <h2 className="card-title mb-4">Form Integration Example</h2>
            <form
              onSubmit={(e) => {
                e.preventDefault();
                toast.info('Submitting form...', { duration: 1000 });
                setTimeout(() => {
                  const success = Math.random() > 0.5;
                  if (success) {
                    toast.success('Form submitted successfully!');
                  } else {
                    toast.error('Failed to submit form. Please try again.');
                  }
                }, 1500);
              }}
              className="space-y-4"
            >
              <input
                type="text"
                placeholder="Enter your name"
                className="input input-bordered w-full"
                required
              />
              <input
                type="email"
                placeholder="Enter your email"
                className="input input-bordered w-full"
                required
              />
              <button type="submit" className="btn btn-primary">
                Submit Form (50% success rate)
              </button>
            </form>
          </div>
        </section>
      </div>
    </div>
  );
}
