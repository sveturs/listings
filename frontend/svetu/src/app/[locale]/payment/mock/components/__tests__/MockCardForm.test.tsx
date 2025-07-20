import { render, screen, fireEvent, waitFor } from '@testing-library/react';
import userEvent from '@testing-library/user-event';
import MockCardForm from '../MockCardForm';
import { useForm } from 'react-hook-form';

// Mock react-hook-form
const mockRegister = jest.fn();
const mockHandleSubmit = jest.fn();
const mockSetValue = jest.fn();
const mockWatch = jest.fn();
const mockGetValues = jest.fn();
const mockGetFieldState = jest.fn();
const mockSetError = jest.fn();
const mockClearErrors = jest.fn();
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
    setValue: mockSetValue,
    watch: mockWatch,
    getValues: mockGetValues,
    getFieldState: mockGetFieldState,
    setError: mockSetError,
    clearErrors: mockClearErrors,
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

describe('MockCardForm', () => {
  const mockOnSubmit = jest.fn();
  const mockTestCards = [
    {
      number: '4111111111111111',
      type: 'success',
      description: 'Успешный платеж',
    },
    {
      number: '4000000000000002',
      type: 'decline',
      description: 'Отклоненный платеж',
    },
    {
      number: '4000000000003220',
      type: '3ds',
      description: '3D Secure',
    },
  ];

  const defaultProps = {
    onSubmit: mockOnSubmit,
    testCards: mockTestCards,
  };

  beforeEach(() => {
    jest.clearAllMocks();
  });

  describe('рендеринг компонента', () => {
    it('должен отображать все поля формы', () => {
      render(<MockCardForm {...defaultProps} />);

      expect(screen.getByLabelText('Номер карты')).toBeInTheDocument();
      expect(screen.getByLabelText('Имя на карте')).toBeInTheDocument();
      expect(screen.getByLabelText('Месяц')).toBeInTheDocument();
      expect(screen.getByLabelText('Год')).toBeInTheDocument();
      expect(screen.getByLabelText('CVV')).toBeInTheDocument();
    });

    it('должен отображать кнопку отправки', () => {
      render(<MockCardForm {...defaultProps} />);

      expect(
        screen.getByRole('button', { name: 'Оплатить' })
      ).toBeInTheDocument();
    });

    it('должен отображать поддерживаемые карты', () => {
      render(<MockCardForm {...defaultProps} />);

      expect(screen.getByText('VISA')).toBeInTheDocument();
      expect(screen.getByText('MasterCard')).toBeInTheDocument();
      expect(screen.getByText('Maestro')).toBeInTheDocument();
    });

    it("должен иметь правильные placeholder'ы", () => {
      render(<MockCardForm {...defaultProps} />);

      expect(
        screen.getByPlaceholderText('1234 5678 9012 3456')
      ).toBeInTheDocument();
      expect(screen.getByPlaceholderText('JOHN DOE')).toBeInTheDocument();
      expect(screen.getByPlaceholderText('MM')).toBeInTheDocument();
      expect(screen.getByPlaceholderText('YY')).toBeInTheDocument();
      expect(screen.getByPlaceholderText('123')).toBeInTheDocument();
    });
  });

  describe('тестовые карты', () => {
    it('должен отображать секцию тестовых карт по умолчанию', () => {
      render(<MockCardForm {...defaultProps} />);

      expect(screen.getByText('Тестовые карты')).toBeInTheDocument();
      expect(screen.getByText('4111111111111111')).toBeInTheDocument();
      expect(screen.getByText('Успешный платеж')).toBeInTheDocument();
      expect(screen.getByText('4000000000000002')).toBeInTheDocument();
      expect(screen.getByText('Отклоненный платеж')).toBeInTheDocument();
      expect(screen.getByText('4000000000003220')).toBeInTheDocument();
      expect(screen.getByText('3D Secure')).toBeInTheDocument();
    });

    it('должен скрывать секцию тестовых карт при клике на крестик', async () => {
      const user = userEvent.setup();
      render(<MockCardForm {...defaultProps} />);

      const closeButton = screen.getByRole('button', { name: '✕' });
      await user.click(closeButton);

      expect(screen.queryByText('Тестовые карты')).not.toBeInTheDocument();
    });

    it('должен заполнять форму при клике на тестовую карту', async () => {
      const user = userEvent.setup();
      render(<MockCardForm {...defaultProps} />);

      const firstTestCard = screen.getByRole('button', {
        name: /4111111111111111/,
      });
      await user.click(firstTestCard);

      expect(mockSetValue).toHaveBeenCalledWith(
        'cardNumber',
        '4111111111111111'
      );
      expect(mockSetValue).toHaveBeenCalledWith('cardHolder', 'TEST USER');
      expect(mockSetValue).toHaveBeenCalledWith('expiryMonth', '12');
      expect(mockSetValue).toHaveBeenCalledWith('expiryYear', '25');
      expect(mockSetValue).toHaveBeenCalledWith('cvv', '123');
    });

    it('должен заполнять форму для разных тестовых карт', async () => {
      const user = userEvent.setup();
      render(<MockCardForm {...defaultProps} />);

      // Клик на вторую тестовую карту
      const secondTestCard = screen.getByRole('button', {
        name: /4000000000000002/,
      });
      await user.click(secondTestCard);

      expect(mockSetValue).toHaveBeenCalledWith(
        'cardNumber',
        '4000000000000002'
      );
    });

    it('должен отображать все тестовые карты из props', () => {
      const customTestCards = [
        {
          number: '1234567890123456',
          type: 'custom',
          description: 'Кастомная карта',
        },
      ];

      render(<MockCardForm {...defaultProps} testCards={customTestCards} />);

      expect(screen.getByText('1234567890123456')).toBeInTheDocument();
      expect(screen.getByText('Кастомная карта')).toBeInTheDocument();
    });
  });

  describe('валидация формы', () => {
    it('должен использовать zod resolver для валидации', () => {
      render(<MockCardForm {...defaultProps} />);

      // Проверяем что useForm был вызван с zodResolver
      expect(mockUseForm).toHaveBeenCalledWith({
        resolver: expect.any(Function),
      });
    });

    it('должен регистрировать все поля формы', () => {
      render(<MockCardForm {...defaultProps} />);

      expect(mockRegister).toHaveBeenCalledWith('cardNumber');
      expect(mockRegister).toHaveBeenCalledWith('cardHolder');
      expect(mockRegister).toHaveBeenCalledWith('expiryMonth');
      expect(mockRegister).toHaveBeenCalledWith('expiryYear');
      expect(mockRegister).toHaveBeenCalledWith('cvv');
    });

    it('должен отображать ошибки валидации', () => {
      // Mock errors
      mockUseForm.mockReturnValue({
        register: mockRegister.mockReturnValue({}),
        handleSubmit: mockHandleSubmit,
        setValue: mockSetValue,
        watch: mockWatch,
        getValues: mockGetValues,
        getFieldState: mockGetFieldState,
        setError: mockSetError,
        clearErrors: mockClearErrors,
        setFocus: mockSetFocus,
        reset: mockReset,
        resetField: mockResetField,
        trigger: mockTrigger,
        unregister: mockUnregister,
        subscribe: mockSubscribe,
        control: mockControl as any,
        formState: {
          errors: {
            cardNumber: { type: 'required', message: 'Invalid card number' },
            cardHolder: {
              type: 'required',
              message: 'Cardholder name required',
            },
            expiryMonth: { type: 'required', message: 'Invalid month' },
            expiryYear: { type: 'required', message: 'Invalid year' },
            cvv: { type: 'required', message: 'Invalid CVV' },
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

      render(<MockCardForm {...defaultProps} />);

      expect(screen.getByText('Invalid card number')).toBeInTheDocument();
      expect(screen.getByText('Cardholder name required')).toBeInTheDocument();
      expect(screen.getByText('Неверный месяц')).toBeInTheDocument();
      expect(screen.getByText('Неверный год')).toBeInTheDocument();
      expect(screen.getByText('Неверный CVV')).toBeInTheDocument();
    });

    it('должен добавлять CSS класс ошибки к полям с ошибками', () => {
      mockUseForm.mockReturnValue({
        register: mockRegister.mockReturnValue({}),
        handleSubmit: mockHandleSubmit,
        setValue: mockSetValue,
        watch: mockWatch,
        getValues: mockGetValues,
        getFieldState: mockGetFieldState,
        setError: mockSetError,
        clearErrors: mockClearErrors,
        setFocus: mockSetFocus,
        reset: mockReset,
        resetField: mockResetField,
        trigger: mockTrigger,
        unregister: mockUnregister,
        subscribe: mockSubscribe,
        control: mockControl as any,
        formState: {
          errors: {
            cardNumber: { type: 'required', message: 'Error' },
            cardHolder: { type: 'required', message: 'Error' },
            expiryMonth: { type: 'required', message: 'Error' },
            expiryYear: { type: 'required', message: 'Error' },
            cvv: { type: 'required', message: 'Error' },
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

      render(<MockCardForm {...defaultProps} />);

      expect(screen.getByPlaceholderText('1234 5678 9012 3456')).toHaveClass(
        'input-error'
      );
      expect(screen.getByPlaceholderText('JOHN DOE')).toHaveClass(
        'input-error'
      );
      expect(screen.getByPlaceholderText('MM')).toHaveClass('input-error');
      expect(screen.getByPlaceholderText('YY')).toHaveClass('input-error');
      expect(screen.getByPlaceholderText('123')).toHaveClass('input-error');
    });
  });

  describe('отправка формы', () => {
    it('должен вызывать onSubmit при отправке формы', async () => {
      const user = userEvent.setup();
      mockHandleSubmit.mockImplementation((fn) => (e: any) => {
        e.preventDefault();
        fn({ cardNumber: '4111111111111111' });
      });

      render(<MockCardForm {...defaultProps} />);

      const submitButton = screen.getByRole('button', { name: 'Оплатить' });
      await user.click(submitButton);

      await waitFor(() => {
        expect(mockOnSubmit).toHaveBeenCalledWith({
          cardNumber: '4111111111111111',
        });
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
        fn({});
      });

      render(<MockCardForm {...defaultProps} />);

      const submitButton = screen.getByRole('button', { name: 'Оплатить' });
      await user.click(submitButton);

      // Проверяем состояние загрузки
      expect(
        screen.getByRole('button', { name: 'Обработка...' })
      ).toBeInTheDocument();
      expect(submitButton).toBeDisabled();
      expect(submitButton).toHaveClass('loading');

      // Завершаем отправку
      resolveSubmit!();
      await waitFor(() => {
        expect(
          screen.getByRole('button', { name: 'Оплатить' })
        ).toBeInTheDocument();
        expect(submitButton).not.toBeDisabled();
      });
    });

    it('должен восстанавливать состояние кнопки после ошибки', async () => {
      const user = userEvent.setup();
      mockOnSubmit.mockImplementation(() =>
        Promise.reject(new Error('Submission failed'))
      );
      mockHandleSubmit.mockImplementation((fn) => async (e: any) => {
        e.preventDefault();
        try {
          await fn({});
        } catch {
          // Ошибка обработана в компоненте
        }
      });

      render(<MockCardForm {...defaultProps} />);

      const submitButton = screen.getByRole('button', { name: 'Оплатить' });
      await user.click(submitButton);

      await waitFor(() => {
        expect(
          screen.getByRole('button', { name: 'Оплатить' })
        ).toBeInTheDocument();
        expect(submitButton).not.toBeDisabled();
      });
    });
  });

  describe('ограничения ввода', () => {
    it('должен иметь правильные maxLength для полей', () => {
      render(<MockCardForm {...defaultProps} />);

      expect(
        screen.getByPlaceholderText('1234 5678 9012 3456')
      ).toHaveAttribute('maxLength', '16');
      expect(screen.getByPlaceholderText('MM')).toHaveAttribute(
        'maxLength',
        '2'
      );
      expect(screen.getByPlaceholderText('YY')).toHaveAttribute(
        'maxLength',
        '2'
      );
      expect(screen.getByPlaceholderText('123')).toHaveAttribute(
        'maxLength',
        '4'
      );
    });

    it('должен иметь uppercase стиль для поля имени', () => {
      render(<MockCardForm {...defaultProps} />);

      const cardHolderInput = screen.getByPlaceholderText('JOHN DOE');
      expect(cardHolderInput).toHaveStyle('text-transform: uppercase');
    });

    it('должен иметь monospace шрифт для номера карты', () => {
      render(<MockCardForm {...defaultProps} />);

      expect(screen.getByPlaceholderText('1234 5678 9012 3456')).toHaveClass(
        'font-mono'
      );
    });
  });

  describe('доступность', () => {
    it('должен иметь правильные labels для всех полей', () => {
      render(<MockCardForm {...defaultProps} />);

      expect(screen.getByLabelText('Номер карты')).toBeInTheDocument();
      expect(screen.getByLabelText('Имя на карте')).toBeInTheDocument();
      expect(screen.getByLabelText('Месяц')).toBeInTheDocument();
      expect(screen.getByLabelText('Год')).toBeInTheDocument();
      expect(screen.getByLabelText('CVV')).toBeInTheDocument();
    });

    it('должен иметь правильный тип формы', () => {
      const { container } = render(<MockCardForm {...defaultProps} />);

      const form = container.querySelector('form');
      expect(form).toBeInTheDocument();
    });

    it('должен предотвращать отправку формы по умолчанию', () => {
      const { container } = render(<MockCardForm {...defaultProps} />);

      const form = container.querySelector('form');
      const submitEvent = new Event('submit');
      const preventDefaultSpy = jest.spyOn(submitEvent, 'preventDefault');

      fireEvent.submit(form!, submitEvent);

      expect(preventDefaultSpy).toHaveBeenCalled();
    });
  });

  describe('edge cases', () => {
    it('должен обрабатывать пустой массив тестовых карт', () => {
      render(<MockCardForm {...defaultProps} testCards={[]} />);

      expect(screen.getByText('Тестовые карты')).toBeInTheDocument();
      // Но не должно быть кнопок тестовых карт
      expect(
        screen.queryByRole('button', { name: /4111/ })
      ).not.toBeInTheDocument();
    });

    it('должен обрабатывать отсутствие onSubmit без ошибок', () => {
      const { container } = render(
        <MockCardForm onSubmit={undefined as any} testCards={mockTestCards} />
      );

      const form = container.querySelector('form');
      expect(() => fireEvent.submit(form!)).not.toThrow();
    });

    it('должен корректно отображать длинные описания карт', () => {
      const longDescriptionCards = [
        {
          number: '4111111111111111',
          type: 'test',
          description:
            'Очень длинное описание тестовой карты которое может не поместиться',
        },
      ];

      render(
        <MockCardForm {...defaultProps} testCards={longDescriptionCards} />
      );

      expect(
        screen.getByText(
          'Очень длинное описание тестовой карты которое может не поместиться'
        )
      ).toBeInTheDocument();
    });

    it('должен обрабатывать специальные символы в номерах карт', () => {
      const specialCards = [
        {
          number: '4111-1111-1111-1111',
          type: 'test',
          description: 'С дефисами',
        },
      ];

      render(<MockCardForm {...defaultProps} testCards={specialCards} />);

      expect(screen.getByText('4111-1111-1111-1111')).toBeInTheDocument();
    });
  });
});
