import { render, screen, waitFor } from '@testing-library/react';
import userEvent from '@testing-library/user-event';
import Mock3DSForm from '../Mock3DSForm';
import { useForm } from 'react-hook-form';

// Mock react-hook-form
const mockRegister = jest.fn();
const mockHandleSubmit = jest.fn();
const mockWatch = jest.fn();
const mockGetValues = jest.fn();
const mockGetFieldState = jest.fn();
const mockSetError = jest.fn();
const mockClearErrors = jest.fn();
const mockSetValue = jest.fn();
const mockSetFocus = jest.fn();
const mockReset = jest.fn();
const mockResetField = jest.fn();
const mockTrigger = jest.fn();
const mockUnregister = jest.fn();
const mockSubscribe = jest.fn();
const mockControl = {};

jest.mock('react-hook-form', () => ({
  useForm: jest.fn(),
}));

const mockUseForm = jest.mocked(useForm);

beforeEach(() => {
  mockUseForm.mockReturnValue({
    register: mockRegister.mockReturnValue({}),
    handleSubmit: mockHandleSubmit.mockImplementation((fn) => fn),
    watch: mockWatch,
    getValues: mockGetValues,
    getFieldState: mockGetFieldState,
    setError: mockSetError,
    clearErrors: mockClearErrors,
    setValue: mockSetValue,
    setFocus: mockSetFocus,
    reset: mockReset,
    resetField: mockResetField,
    trigger: mockTrigger,
    unregister: mockUnregister,
    subscribe: mockSubscribe,
    control: mockControl as any,
    formState: {
      errors: {},
      isDirty: false,
      isLoading: false,
      isSubmitted: false,
      isSubmitSuccessful: false,
      isSubmitting: false,
      isValid: true,
      isValidating: false,
      submitCount: 0,
      touchedFields: {},
      dirtyFields: {},
      validatingFields: {},
      defaultValues: {},
      disabled: false,
      isReady: true,
    },
  });
});

describe('Mock3DSForm', () => {
  const mockOnSubmit = jest.fn();

  const defaultProps = {
    onSubmit: mockOnSubmit,
  };

  beforeEach(() => {
    jest.clearAllMocks();
  });

  describe('рендеринг компонента', () => {
    it('должен отображать заголовок 3D Secure', () => {
      render(<Mock3DSForm {...defaultProps} />);

      expect(screen.getByText('3D Secure Authentication')).toBeInTheDocument();
      expect(
        screen.getByText(/Ваш банк запросил дополнительную аутентификацию/)
      ).toBeInTheDocument();
    });

    it('должен отображать имитацию браузера банка', () => {
      render(<Mock3DSForm {...defaultProps} />);

      expect(
        screen.getByText('https://secure.bank.rs/3ds-auth')
      ).toBeInTheDocument();
      expect(screen.getByText('Банк Србије')).toBeInTheDocument();
      expect(screen.getByText('3D Secure Verification')).toBeInTheDocument();
    });

    it('должен отображать поле для ввода кода', () => {
      render(<Mock3DSForm {...defaultProps} />);

      expect(
        screen.getByLabelText('Унесите код из SMS-а:')
      ).toBeInTheDocument();
      expect(screen.getByPlaceholderText('123456')).toBeInTheDocument();
    });

    it('должен отображать кнопки подтверждения и отмены', () => {
      render(<Mock3DSForm {...defaultProps} />);

      expect(
        screen.getByRole('button', { name: 'Потврди' })
      ).toBeInTheDocument();
      expect(
        screen.getByRole('button', { name: 'Откажи' })
      ).toBeInTheDocument();
    });

    it('должен отображать подсказку с тестовым кодом', () => {
      render(<Mock3DSForm {...defaultProps} />);

      expect(
        screen.getByText(/За тестирање користите код:/)
      ).toBeInTheDocument();
      expect(screen.getByText('123')).toBeInTheDocument();
    });

    it('должен отображать информацию о безопасности', () => {
      render(<Mock3DSForm {...defaultProps} />);

      expect(
        screen.getByText('Заштићено 256-битним SSL шифровањем')
      ).toBeInTheDocument();
    });

    it('должен отображать справочную информацию', () => {
      render(<Mock3DSForm {...defaultProps} />);

      expect(screen.getByText('Информације о 3D Secure')).toBeInTheDocument();
      expect(
        screen.getByText(/3D Secure је додатни слој безбедности/)
      ).toBeInTheDocument();
    });
  });

  describe('отображение суммы', () => {
    it('должен отображать сумму платежа если передана', () => {
      render(<Mock3DSForm {...defaultProps} amount={10000} />);

      expect(screen.getByText('Сума за плаћање:')).toBeInTheDocument();
      expect(screen.getByText('10.000,00 RSD')).toBeInTheDocument();
    });

    it('не должен отображать сумму если не передана', () => {
      render(<Mock3DSForm {...defaultProps} />);

      expect(screen.queryByText('Сума за плаћање:')).not.toBeInTheDocument();
    });

    it('должен форматировать сумму в соответствии с локалью sr-RS', () => {
      render(<Mock3DSForm {...defaultProps} amount={123456.78} />);

      expect(screen.getByText('123.456,78 RSD')).toBeInTheDocument();
    });

    it('должен обрабатывать нулевую сумму', () => {
      render(<Mock3DSForm {...defaultProps} amount={0} />);

      expect(screen.getByText('0,00 RSD')).toBeInTheDocument();
    });

    it('должен обрабатывать большие суммы', () => {
      render(<Mock3DSForm {...defaultProps} amount={1000000} />);

      expect(screen.getByText('1.000.000,00 RSD')).toBeInTheDocument();
    });
  });

  describe('валидация формы', () => {
    it('должен регистрировать поле кода с валидацией', () => {
      render(<Mock3DSForm {...defaultProps} />);

      expect(mockRegister).toHaveBeenCalledWith('code', {
        required: 'Код је обавезан',
        minLength: {
          value: 3,
          message: 'Код мора имати најмање 3 цифре',
        },
      });
    });

    it('должен отображать ошибки валидации', () => {
      mockUseForm.mockReturnValue({
        register: mockRegister.mockReturnValue({}),
        handleSubmit: mockHandleSubmit,
        watch: mockWatch,
        getValues: mockGetValues,
        getFieldState: mockGetFieldState,
        setError: mockSetError,
        clearErrors: mockClearErrors,
        setValue: mockSetValue,
        setFocus: mockSetFocus,
        reset: mockReset,
        resetField: mockResetField,
        trigger: mockTrigger,
        unregister: mockUnregister,
        subscribe: mockSubscribe,
        control: mockControl as any,
        formState: {
          errors: {
            code: { type: 'required', message: 'Код је обавезан' },
          },
          isDirty: false,
          isLoading: false,
          isSubmitted: false,
          isSubmitSuccessful: false,
          isSubmitting: false,
          isValid: false,
          isValidating: false,
          submitCount: 0,
          touchedFields: {},
          dirtyFields: {},
          validatingFields: {},
          defaultValues: {},
          disabled: false,
          isReady: true,
        },
      });

      render(<Mock3DSForm {...defaultProps} />);

      expect(screen.getByText('Код је обавезан')).toBeInTheDocument();
    });

    it('должен добавлять CSS класс ошибки к поль с ошибкой', () => {
      mockUseForm.mockReturnValue({
        register: mockRegister.mockReturnValue({}),
        handleSubmit: mockHandleSubmit,
        watch: mockWatch,
        getValues: mockGetValues,
        getFieldState: mockGetFieldState,
        setError: mockSetError,
        clearErrors: mockClearErrors,
        setValue: mockSetValue,
        setFocus: mockSetFocus,
        reset: mockReset,
        resetField: mockResetField,
        trigger: mockTrigger,
        unregister: mockUnregister,
        subscribe: mockSubscribe,
        control: mockControl as any,
        formState: {
          errors: {
            code: { type: 'required', message: 'Error' },
          },
          isDirty: false,
          isLoading: false,
          isSubmitted: false,
          isSubmitSuccessful: false,
          isSubmitting: false,
          isValid: false,
          isValidating: false,
          submitCount: 0,
          touchedFields: {},
          dirtyFields: {},
          validatingFields: {},
          defaultValues: {},
          disabled: false,
          isReady: true,
        },
      });

      render(<Mock3DSForm {...defaultProps} />);

      expect(screen.getByPlaceholderText('123456')).toHaveClass('input-error');
    });
  });

  describe('отправка формы', () => {
    it('должен вызывать onSubmit с кодом при отправке формы', async () => {
      const user = userEvent.setup();
      mockHandleSubmit.mockImplementation((fn) => (e: any) => {
        e.preventDefault();
        fn({ code: '123' });
      });

      render(<Mock3DSForm {...defaultProps} />);

      const submitButton = screen.getByRole('button', { name: 'Потврди' });
      await user.click(submitButton);

      await waitFor(() => {
        expect(mockOnSubmit).toHaveBeenCalledWith('123');
      });
    });

    it('должен показывать состояние загрузки во время отправки', async () => {
      const user = userEvent.setup();
      let resolveSubmit: () => void;
      const submitPromise = new Promise<void>((resolve) => {
        resolveSubmit = resolve;
      });

      mockOnSubmit.mockReturnValue(submitPromise);
      mockHandleSubmit.mockImplementation((fn) => (e: any) => {
        e.preventDefault();
        fn({ code: '123' });
      });

      render(<Mock3DSForm {...defaultProps} />);

      const submitButton = screen.getByRole('button', { name: 'Потврди' });
      await user.click(submitButton);

      // Проверяем состояние загрузки
      expect(
        screen.getByRole('button', { name: 'Потврђивање...' })
      ).toBeInTheDocument();
      expect(submitButton).toBeDisabled();
      expect(submitButton).toHaveClass('loading');

      // Проверяем что кнопка отмены тоже отключена
      expect(screen.getByRole('button', { name: 'Откажи' })).toBeDisabled();

      // Завершаем отправку
      resolveSubmit!();
      await waitFor(() => {
        expect(
          screen.getByRole('button', { name: 'Потврди' })
        ).toBeInTheDocument();
        expect(submitButton).not.toBeDisabled();
      });
    });

    it('должен восстанавливать состояние кнопок после ошибки', async () => {
      const user = userEvent.setup();
      mockOnSubmit.mockImplementation(async () => {
        throw new Error('Submission failed');
      });
      mockHandleSubmit.mockImplementation((fn) => async (e: any) => {
        e.preventDefault();
        try {
          await fn({ code: '123' });
        } catch (error) {
          // Ошибка обработана в компоненте
        }
      });

      render(<Mock3DSForm {...defaultProps} />);

      const submitButton = screen.getByRole('button', { name: 'Потврди' });
      await user.click(submitButton);

      await waitFor(() => {
        expect(
          screen.getByRole('button', { name: 'Потврди' })
        ).toBeInTheDocument();
        expect(submitButton).not.toBeDisabled();
        expect(
          screen.getByRole('button', { name: 'Откажи' })
        ).not.toBeDisabled();
      });
    });
  });

  describe('кнопка отмены', () => {
    it('должен вызывать onSubmit с "cancel" при клике на Откажи', async () => {
      const user = userEvent.setup();
      // Явно сбрасываем mock чтобы избежать наследования от предыдущих тестов
      mockOnSubmit.mockReset();
      mockOnSubmit.mockResolvedValue(undefined);

      render(<Mock3DSForm {...defaultProps} />);

      const cancelButton = screen.getByRole('button', { name: 'Откажи' });
      await user.click(cancelButton);

      expect(mockOnSubmit).toHaveBeenCalledWith('cancel');
    });

    it('должен отключать кнопку отмены во время отправки', async () => {
      const user = userEvent.setup();
      let resolveSubmit: () => void;
      const submitPromise = new Promise<void>((resolve) => {
        resolveSubmit = resolve;
      });

      mockOnSubmit.mockReturnValue(submitPromise);
      mockHandleSubmit.mockImplementation((fn) => (e: any) => {
        e.preventDefault();
        fn({ code: '123' });
      });

      render(<Mock3DSForm {...defaultProps} />);

      const submitButton = screen.getByRole('button', { name: 'Потврди' });
      await user.click(submitButton);

      const cancelButton = screen.getByRole('button', { name: 'Откажи' });
      expect(cancelButton).toBeDisabled();

      resolveSubmit!();
      await waitFor(() => {
        expect(cancelButton).not.toBeDisabled();
      });
    });
  });

  describe('ограничения ввода', () => {
    it('должен иметь правильный maxLength для поля кода', () => {
      render(<Mock3DSForm {...defaultProps} />);

      expect(screen.getByPlaceholderText('123456')).toHaveAttribute(
        'maxLength',
        '6'
      );
    });

    it('должен иметь monospace и центрированный текст для поля кода', () => {
      render(<Mock3DSForm {...defaultProps} />);

      const codeInput = screen.getByPlaceholderText('123456');
      expect(codeInput).toHaveClass('font-mono');
      expect(codeInput).toHaveClass('text-center');
      expect(codeInput).toHaveClass('text-lg');
    });
  });

  describe('доступность', () => {
    it('должен иметь правильный label для поля кода', () => {
      render(<Mock3DSForm {...defaultProps} />);

      expect(
        screen.getByLabelText('Унесите код из SMS-а:')
      ).toBeInTheDocument();
    });

    it('должен иметь правильную структуру формы', () => {
      const { container } = render(<Mock3DSForm {...defaultProps} />);

      const form = container.querySelector('form');
      expect(form).toBeInTheDocument();
    });

    it('должен иметь иконку безопасности с правильными атрибутами', () => {
      const { container } = render(<Mock3DSForm {...defaultProps} />);

      const securityIcon = container.querySelector('svg[viewBox="0 0 20 20"]');
      expect(securityIcon).toBeInTheDocument();
      expect(securityIcon).toHaveClass('w-4');
      expect(securityIcon).toHaveClass('h-4');
    });

    it('должен иметь информационную иконку в справке', () => {
      const { container } = render(<Mock3DSForm {...defaultProps} />);

      const infoIcon = container.querySelector('svg[viewBox="0 0 24 24"]');
      expect(infoIcon).toBeInTheDocument();
      expect(infoIcon).toHaveClass('w-6');
      expect(infoIcon).toHaveClass('h-6');
    });
  });

  describe('CSS классы и стили', () => {
    it('должен применять правильные классы к основным элементам', () => {
      const { container } = render(<Mock3DSForm {...defaultProps} />);

      expect(container.querySelector('.mockup-browser')).toBeInTheDocument();
      expect(container.querySelector('.alert-info')).toBeInTheDocument();
      expect(screen.getByText('3D Secure Authentication')).toHaveClass(
        'text-primary'
      );
    });

    it('должен иметь flexbox layout для кнопок', () => {
      const { container } = render(<Mock3DSForm {...defaultProps} />);

      const buttonContainer = container.querySelector('.flex.gap-2');
      expect(buttonContainer).toBeInTheDocument();
      expect(screen.getByRole('button', { name: 'Потврди' })).toHaveClass(
        'flex-1'
      );
    });

    it('должен применять правильные классы к alert элементам', () => {
      const { container } = render(
        <Mock3DSForm {...defaultProps} amount={1000} />
      );

      const alerts = container.querySelectorAll('.alert');
      expect(alerts).toHaveLength(2); // alert-info для суммы и общий alert
    });
  });

  describe('edge cases', () => {
    it('должен обрабатывать отсутствие onSubmit без ошибок', () => {
      expect(() =>
        render(<Mock3DSForm onSubmit={undefined as any} />)
      ).not.toThrow();
    });

    it('должен обрабатывать отрицательные суммы', () => {
      render(<Mock3DSForm {...defaultProps} amount={-1000} />);

      expect(screen.getByText('-1.000,00 RSD')).toBeInTheDocument();
    });

    it('должен обрабатывать дробные суммы', () => {
      render(<Mock3DSForm {...defaultProps} amount={1234.56} />);

      expect(screen.getByText('1.234,56 RSD')).toBeInTheDocument();
    });

    it('должен обрабатывать очень маленькие суммы', () => {
      render(<Mock3DSForm {...defaultProps} amount={0.01} />);

      expect(screen.getByText('0,01 RSD')).toBeInTheDocument();
    });

    it('должен корректно отображать компонент без amount', () => {
      render(<Mock3DSForm {...defaultProps} amount={undefined} />);

      expect(screen.queryByText('Сума за плаћање:')).not.toBeInTheDocument();
      expect(screen.getByText('3D Secure Authentication')).toBeInTheDocument();
    });

    it('должен обрабатывать пустые строки в ошибках валидации', () => {
      mockUseForm.mockReturnValue({
        register: mockRegister.mockReturnValue({}),
        handleSubmit: mockHandleSubmit,
        watch: mockWatch,
        getValues: mockGetValues,
        getFieldState: mockGetFieldState,
        setError: mockSetError,
        clearErrors: mockClearErrors,
        setValue: mockSetValue,
        setFocus: mockSetFocus,
        reset: mockReset,
        resetField: mockResetField,
        trigger: mockTrigger,
        unregister: mockUnregister,
        subscribe: mockSubscribe,
        control: mockControl as any,
        formState: {
          errors: {
            code: { type: 'required', message: '' },
          },
          isDirty: false,
          isLoading: false,
          isSubmitted: false,
          isSubmitSuccessful: false,
          isSubmitting: false,
          isValid: false,
          isValidating: false,
          submitCount: 0,
          touchedFields: {},
          dirtyFields: {},
          validatingFields: {},
          defaultValues: {},
          disabled: false,
          isReady: true,
        },
      });

      const { container } = render(<Mock3DSForm {...defaultProps} />);

      // Пустое сообщение ошибки не должно ломать компонент
      expect(container.querySelector('.text-error')).toBeInTheDocument();
    });

    it('должен обрабатывать долгий async onSubmit', async () => {
      const user = userEvent.setup();
      const longRunningSubmit = jest.fn(
        () => new Promise((resolve) => setTimeout(resolve, 100))
      );

      mockHandleSubmit.mockImplementation((fn) => (e: any) => {
        e.preventDefault();
        fn({ code: '123' });
      });

      render(<Mock3DSForm onSubmit={longRunningSubmit} />);

      const submitButton = screen.getByRole('button', { name: 'Потврди' });
      await user.click(submitButton);

      // Проверяем что состояние загрузки активно
      expect(
        screen.getByRole('button', { name: 'Потврђивање...' })
      ).toBeInTheDocument();

      // Ждем завершения
      await waitFor(
        () => {
          expect(
            screen.getByRole('button', { name: 'Потврди' })
          ).toBeInTheDocument();
        },
        { timeout: 200 }
      );
    });
  });
});
