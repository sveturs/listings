import { render, screen } from '@testing-library/react';
import PaymentProcessing from '../PaymentProcessing';

describe('PaymentProcessing', () => {
  const defaultProps = {
    status: 'checking',
    attempts: 1,
    maxAttempts: 10,
  };

  describe('рендеринг компонента', () => {
    it('должен отображать базовую структуру', () => {
      render(<PaymentProcessing {...defaultProps} />);

      expect(screen.getByRole('progressbar')).toBeInTheDocument();
      expect(screen.getByText('Проверяем статус')).toBeInTheDocument();
      expect(screen.getByText('1/10')).toBeInTheDocument();
      expect(
        screen.getByText(
          'Пожалуйста, не закрывайте это окно во время обработки платежа'
        )
      ).toBeInTheDocument();
    });

    it('должен отображать этапы обработки', () => {
      render(<PaymentProcessing {...defaultProps} />);

      expect(screen.getByText('Этапы обработки:')).toBeInTheDocument();
      expect(screen.getByText('Проверка')).toBeInTheDocument();
      expect(screen.getByText('Обработка')).toBeInTheDocument();
      expect(screen.getByText('Авторизация')).toBeInTheDocument();
      expect(screen.getByText('Завершение')).toBeInTheDocument();
    });

    it('должен отображать анимированные точки загрузки', () => {
      const { container } = render(<PaymentProcessing {...defaultProps} />);

      const dots = container.querySelectorAll('.animate-bounce');
      expect(dots).toHaveLength(3); // Только 3 точки загрузки
    });
  });

  describe('статусы платежей', () => {
    it('должен отображать правильную информацию для статуса "checking"', () => {
      render(<PaymentProcessing {...defaultProps} status="checking" />);

      expect(screen.getByText('🔍')).toBeInTheDocument();
      expect(screen.getByText('Проверка платежа')).toBeInTheDocument();
      expect(
        screen.getByText('Получаем информацию о статусе транзакции')
      ).toBeInTheDocument();
    });

    it('должен отображать правильную информацию для статуса "pending"', () => {
      render(<PaymentProcessing {...defaultProps} status="pending" />);

      expect(screen.getByText('⏳')).toBeInTheDocument();
      expect(screen.getByText('Ожидание обработки')).toBeInTheDocument();
      expect(
        screen.getByText('Платеж находится в очереди на обработку')
      ).toBeInTheDocument();
    });

    it('должен отображать правильную информацию для статуса "processing"', () => {
      render(<PaymentProcessing {...defaultProps} status="processing" />);

      expect(screen.getByText('⚡')).toBeInTheDocument();
      expect(screen.getByText('Обработка платежа')).toBeInTheDocument();
      expect(
        screen.getByText('Банк обрабатывает ваш платеж')
      ).toBeInTheDocument();
    });

    it('должен отображать правильную информацию для статуса "authorized"', () => {
      render(<PaymentProcessing {...defaultProps} status="authorized" />);

      expect(screen.getByText('✅')).toBeInTheDocument();
      expect(screen.getByText('Платеж авторизован')).toBeInTheDocument();
      expect(
        screen.getByText('Средства заблокированы, завершаем транзакцию')
      ).toBeInTheDocument();
    });

    it('должен отображать дефолтную информацию для неизвестного статуса', () => {
      render(<PaymentProcessing {...defaultProps} status="unknown_status" />);

      expect(screen.getByText('❓')).toBeInTheDocument();
      expect(screen.getByText('Unknown_status')).toBeInTheDocument();
      expect(
        screen.getByText('Обрабатываем статус платежа')
      ).toBeInTheDocument();
    });

    it('должен правильно капитализировать первую букву в дефолтном заголовке', () => {
      render(<PaymentProcessing {...defaultProps} status="test_status" />);

      expect(screen.getByText('Test_status')).toBeInTheDocument();
    });
  });

  describe('прогресс бар', () => {
    it('должен рассчитывать прогресс правильно', () => {
      render(
        <PaymentProcessing {...defaultProps} attempts={3} maxAttempts={10} />
      );

      const progressBar = screen.getByRole('progressbar');
      expect(progressBar).toHaveAttribute('value', '30');
      expect(progressBar).toHaveAttribute('max', '100');
      expect(screen.getByText('3/10')).toBeInTheDocument();
    });

    it('должен ограничивать прогресс до 100%', () => {
      render(
        <PaymentProcessing {...defaultProps} attempts={15} maxAttempts={10} />
      );

      const progressBar = screen.getByRole('progressbar');
      expect(progressBar).toHaveAttribute('value', '100');
    });

    it('должен обрабатывать нулевые попытки', () => {
      render(
        <PaymentProcessing {...defaultProps} attempts={0} maxAttempts={10} />
      );

      const progressBar = screen.getByRole('progressbar');
      expect(progressBar).toHaveAttribute('value', '0');
      expect(screen.getByText('0/10')).toBeInTheDocument();
    });

    it('должен обрабатывать максимальные попытки', () => {
      render(
        <PaymentProcessing {...defaultProps} attempts={10} maxAttempts={10} />
      );

      const progressBar = screen.getByRole('progressbar');
      expect(progressBar).toHaveAttribute('value', '100');
      expect(screen.getByText('10/10')).toBeInTheDocument();
    });
  });

  describe('этапы обработки', () => {
    it('должен активировать первый этап для статуса "checking"', () => {
      const { container } = render(
        <PaymentProcessing {...defaultProps} status="checking" />
      );

      const steps = container.querySelectorAll('.step');
      expect(steps[0]).toHaveClass('step-primary');
      expect(steps[1]).not.toHaveClass('step-primary');
      expect(steps[2]).not.toHaveClass('step-primary');
      expect(steps[3]).not.toHaveClass('step-primary');
    });

    it('должен активировать первые два этапа для статуса "pending"', () => {
      const { container } = render(
        <PaymentProcessing {...defaultProps} status="pending" />
      );

      const steps = container.querySelectorAll('.step');
      expect(steps[0]).toHaveClass('step-primary');
      expect(steps[1]).toHaveClass('step-primary');
      expect(steps[2]).not.toHaveClass('step-primary');
      expect(steps[3]).not.toHaveClass('step-primary');
    });

    it('должен активировать первые три этапа для статуса "processing"', () => {
      const { container } = render(
        <PaymentProcessing {...defaultProps} status="processing" />
      );

      const steps = container.querySelectorAll('.step');
      expect(steps[0]).toHaveClass('step-primary');
      expect(steps[1]).toHaveClass('step-primary');
      expect(steps[2]).toHaveClass('step-primary');
      expect(steps[3]).not.toHaveClass('step-primary');
    });

    it('должен активировать все этапы для статуса "authorized"', () => {
      const { container } = render(
        <PaymentProcessing {...defaultProps} status="authorized" />
      );

      const steps = container.querySelectorAll('.step');
      expect(steps[0]).toHaveClass('step-primary');
      expect(steps[1]).toHaveClass('step-primary');
      expect(steps[2]).toHaveClass('step-primary');
      expect(steps[3]).toHaveClass('step-primary');
    });

    it('не должен активировать этапы для неизвестного статуса', () => {
      const { container } = render(
        <PaymentProcessing {...defaultProps} status="unknown" />
      );

      const steps = container.querySelectorAll('.step');
      steps.forEach((step) => {
        expect(step).not.toHaveClass('step-primary');
      });
    });
  });

  describe('CSS классы и стили', () => {
    it('должен применять правильные цвета для разных статусов', () => {
      const { rerender } = render(
        <PaymentProcessing {...defaultProps} status="checking" />
      );
      expect(screen.getByText('🔍')).toHaveClass('text-info');

      rerender(<PaymentProcessing {...defaultProps} status="pending" />);
      expect(screen.getByText('⏳')).toHaveClass('text-warning');

      rerender(<PaymentProcessing {...defaultProps} status="processing" />);
      expect(screen.getByText('⚡')).toHaveClass('text-primary');

      rerender(<PaymentProcessing {...defaultProps} status="authorized" />);
      expect(screen.getByText('✅')).toHaveClass('text-success');

      rerender(<PaymentProcessing {...defaultProps} status="unknown" />);
      expect(screen.getByText('❓')).toHaveClass('text-base-content');
    });

    it('должен применять animate-pulse к иконке', () => {
      render(<PaymentProcessing {...defaultProps} />);

      const icon = screen.getByText('🔍');
      expect(icon).toHaveClass('animate-pulse');
    });

    it('должен применять правильные задержки анимации к точкам', () => {
      const { container } = render(<PaymentProcessing {...defaultProps} />);

      const animationDots = container.querySelectorAll('.animate-bounce');
      const dots = Array.from(animationDots).filter(
        (dot) => dot.classList.contains('w-2') && dot.classList.contains('h-2')
      );

      expect(dots[0]).toHaveStyle('animation-delay: 0ms');
      expect(dots[1]).toHaveStyle('animation-delay: 150ms');
      expect(dots[2]).toHaveStyle('animation-delay: 300ms');
    });
  });

  describe('доступность', () => {
    it('должен иметь правильные ARIA атрибуты для прогресс бара', () => {
      render(
        <PaymentProcessing {...defaultProps} attempts={3} maxAttempts={10} />
      );

      const progressBar = screen.getByRole('progressbar');
      expect(progressBar).toHaveAttribute('value', '30');
      expect(progressBar).toHaveAttribute('max', '100');
    });

    it('должен иметь правильную структуру заголовков', () => {
      render(<PaymentProcessing {...defaultProps} />);

      expect(screen.getByRole('heading', { level: 2 })).toHaveTextContent(
        'Проверка платежа'
      );
      expect(screen.getByRole('heading', { level: 3 })).toHaveTextContent(
        'Этапы обработки:'
      );
    });

    it('должен иметь описательный текст для всех статусов', () => {
      const statuses = ['checking', 'pending', 'processing', 'authorized'];

      statuses.forEach((status) => {
        const { unmount } = render(
          <PaymentProcessing {...defaultProps} status={status} />
        );

        // Проверяем что есть описание
        const descriptions = [
          'Получаем информацию о статусе транзакции',
          'Платеж находится в очереди на обработку',
          'Банк обрабатывает ваш платеж',
          'Средства заблокированы, завершаем транзакцию',
        ];

        const statusIndex = statuses.indexOf(status);
        if (statusIndex !== -1) {
          expect(
            screen.getByText(descriptions[statusIndex])
          ).toBeInTheDocument();
        }

        unmount();
      });
    });
  });

  describe('edge cases', () => {
    it('должен обрабатывать отрицательные попытки', () => {
      render(
        <PaymentProcessing {...defaultProps} attempts={-1} maxAttempts={10} />
      );

      const progressBar = screen.getByRole('progressbar');
      expect(progressBar).toHaveAttribute('value', '-10'); // Компонент показывает реальные значения
      expect(screen.getByText('-1/10')).toBeInTheDocument();
    });

    it('должен обрабатывать нулевые maxAttempts', () => {
      render(
        <PaymentProcessing {...defaultProps} attempts={1} maxAttempts={0} />
      );

      const progressBar = screen.getByRole('progressbar');
      expect(progressBar).toHaveAttribute('value', '100'); // Infinity будет ограничено до 100
      expect(screen.getByText('1/0')).toBeInTheDocument();
    });

    it('должен обрабатывать пустую строку статуса', () => {
      render(<PaymentProcessing {...defaultProps} status="" />);

      expect(screen.getByText('❓')).toBeInTheDocument();
      // Пустой статус будет капитализирован как пустая строка
      expect(
        screen.getByText('Обрабатываем статус платежа')
      ).toBeInTheDocument();
    });

    it('должен обрабатывать очень длинные статусы', () => {
      const longStatus = 'very_long_status_name_that_might_break_layout';
      render(<PaymentProcessing {...defaultProps} status={longStatus} />);

      expect(
        screen.getByText('Very_long_status_name_that_might_break_layout')
      ).toBeInTheDocument();
    });
  });
});
