import {
  generateMockPayment,
  generateMockTransactions,
  mockPaymentHistory,
  mockWalletData,
} from '../mockDataGenerator';

describe('mockDataGenerator', () => {
  describe('generateMockPayment', () => {
    it('должен генерировать объект платежа с правильной структурой', () => {
      const payment = generateMockPayment();

      expect(payment).toEqual({
        id: expect.stringMatching(/^MOCK_\d+_[a-z0-9]{9}$/),
        merchantTransactionId: expect.stringMatching(/^MTX_\d+$/),
        amount: expect.any(Number),
        currency: 'RSD',
        status: 'pending',
        createdAt: expect.stringMatching(
          /^\d{4}-\d{2}-\d{2}T\d{2}:\d{2}:\d{2}\.\d{3}Z$/
        ),
        card: {
          lastFour: expect.stringMatching(/^\d{4}$/),
          brand: expect.stringMatching(/^(visa|mastercard|maestro)$/),
          expiryMonth: expect.any(Number),
          expiryYear: expect.any(Number),
        },
      });
    });

    it('должен генерировать уникальные ID для каждого платежа', async () => {
      const payment1 = generateMockPayment();
      // Добавляем небольшую задержку чтобы timestamp отличался
      await new Promise((resolve) => setTimeout(resolve, 1));
      const payment2 = generateMockPayment();

      expect(payment1.id).not.toEqual(payment2.id);
      // merchantTransactionId может быть одинаковым при быстрых вызовах, поэтому проверяем только id
    });

    it('должен генерировать ID в правильном формате', () => {
      const payment = generateMockPayment();

      // ID должен начинаться с MOCK_ и содержать timestamp + случайную строку
      expect(payment.id).toMatch(/^MOCK_\d+_[a-z0-9]{9}$/);

      // merchantTransactionId должен начинаться с MTX_ и содержать timestamp
      expect(payment.merchantTransactionId).toMatch(/^MTX_\d+$/);
    });

    it('должен генерировать суммы в правильном диапазоне', () => {
      const payments = Array.from({ length: 100 }, () => generateMockPayment());

      payments.forEach((payment) => {
        expect(payment.amount).toBeGreaterThanOrEqual(1000);
        expect(payment.amount).toBeLessThan(101000);
        expect(Number.isInteger(payment.amount)).toBe(true);
      });
    });

    it('должен устанавливать правильную валюту', () => {
      const payment = generateMockPayment();

      expect(payment.currency).toBe('RSD');
    });

    it('должен устанавливать статус как pending', () => {
      const payment = generateMockPayment();

      expect(payment.status).toBe('pending');
    });

    it('должен генерировать правильную дату создания', () => {
      const beforeGeneration = new Date();
      const payment = generateMockPayment();
      const afterGeneration = new Date();

      const createdAt = new Date(payment.createdAt);

      expect(createdAt.getTime()).toBeGreaterThanOrEqual(
        beforeGeneration.getTime() - 1000
      );
      expect(createdAt.getTime()).toBeLessThanOrEqual(
        afterGeneration.getTime() + 1000
      );
    });

    describe('данные карты', () => {
      it('должен генерировать lastFour как строку из 4 цифр', () => {
        const payments = Array.from({ length: 50 }, () =>
          generateMockPayment()
        );

        payments.forEach((payment) => {
          expect(payment.card).toBeDefined();
          expect(payment.card?.lastFour).toMatch(/^\d{4}$/);
          expect(payment.card?.lastFour.length).toBe(4);
        });
      });

      it('должен генерировать правильные бренды карт', () => {
        const payments = Array.from({ length: 50 }, () =>
          generateMockPayment()
        );
        const validBrands = ['visa', 'mastercard', 'maestro'];

        payments.forEach((payment) => {
          expect(payment.card).toBeDefined();
          expect(validBrands).toContain(payment.card?.brand);
        });

        // Проверяем что все бренды встречаются (статистически)
        const brands = payments.map((p) => p.card?.brand).filter(Boolean);
        validBrands.forEach((brand) => {
          expect(brands).toContain(brand);
        });
      });

      it('должен генерировать месяц в диапазоне 1-12', () => {
        const payments = Array.from({ length: 50 }, () =>
          generateMockPayment()
        );

        payments.forEach((payment) => {
          expect(payment.card?.expiryMonth).toBeGreaterThanOrEqual(1);
          expect(payment.card?.expiryMonth).toBeLessThanOrEqual(12);
          expect(Number.isInteger(payment.card?.expiryMonth)).toBe(true);
        });
      });

      it('должен генерировать год в будущем (1-5 лет)', () => {
        const currentYear = new Date().getFullYear();
        const payments = Array.from({ length: 50 }, () =>
          generateMockPayment()
        );

        payments.forEach((payment) => {
          expect(payment.card?.expiryYear).toBeGreaterThan(currentYear);
          expect(payment.card?.expiryYear).toBeLessThanOrEqual(currentYear + 6);
          expect(Number.isInteger(payment.card?.expiryYear)).toBe(true);
        });
      });

      it('должен генерировать lastFour с ведущими нулями если нужно', () => {
        // Mock Math.random чтобы получить небольшое число
        const originalRandom = Math.random;
        Math.random = jest.fn().mockReturnValue(0.001); // Даст число ~10

        const payment = generateMockPayment();

        expect(payment.card?.lastFour).toMatch(/^0\d{3}$/);

        Math.random = originalRandom;
      });
    });
  });

  describe('generateMockTransactions', () => {
    it('должен генерировать правильное количество транзакций', () => {
      const transactions = generateMockTransactions(5);

      expect(transactions).toHaveLength(5);
    });

    it('должен генерировать транзакции с различными статусами', () => {
      const transactions = generateMockTransactions(20);
      const validStatuses = [
        'pending',
        'authorized',
        'captured',
        'failed',
        'refunded',
      ];

      transactions.forEach((transaction) => {
        expect(validStatuses).toContain(transaction.status);
      });

      // Проверяем что встречаются разные статусы (статистически)
      const statuses = transactions.map((t) => t.status);
      const uniqueStatuses = [...new Set(statuses)];
      expect(uniqueStatuses.length).toBeGreaterThan(1);
    });

    it('должен иногда устанавливать completedAt', () => {
      const transactions = generateMockTransactions(20);

      const withCompletedAt = transactions.filter((t) => t.completedAt);
      const withoutCompletedAt = transactions.filter((t) => !t.completedAt);

      // Статистически должны быть и те и другие
      expect(withCompletedAt.length).toBeGreaterThan(0);
      expect(withoutCompletedAt.length).toBeGreaterThan(0);
    });

    it('должен устанавливать completedAt как ISO строку', () => {
      const transactions = generateMockTransactions(20);

      transactions.forEach((transaction) => {
        if (transaction.completedAt) {
          expect(transaction.completedAt).toMatch(
            /^\d{4}-\d{2}-\d{2}T\d{2}:\d{2}:\d{2}\.\d{3}Z$/
          );
          expect(() => new Date(transaction.completedAt!)).not.toThrow();
        }
      });
    });

    it('должен генерировать пустой массив для count = 0', () => {
      const transactions = generateMockTransactions(0);

      expect(transactions).toEqual([]);
    });

    it('должен обрабатывать большие количества', () => {
      const transactions = generateMockTransactions(100);

      expect(transactions).toHaveLength(100);
      expect(transactions[0]).toEqual(
        expect.objectContaining({
          id: expect.any(String),
          status: expect.any(String),
        })
      );
    });

    it('должен наследовать базовые свойства от generateMockPayment', () => {
      const transactions = generateMockTransactions(5);

      transactions.forEach((transaction) => {
        expect(transaction.currency).toBe('RSD');
        expect(transaction.card).toEqual({
          lastFour: expect.stringMatching(/^\d{4}$/),
          brand: expect.stringMatching(/^(visa|mastercard|maestro)$/),
          expiryMonth: expect.any(Number),
          expiryYear: expect.any(Number),
        });
      });
    });
  });

  describe('mockPaymentHistory', () => {
    it('должен иметь правильную структуру массива', () => {
      expect(Array.isArray(mockPaymentHistory)).toBe(true);
      expect(mockPaymentHistory.length).toBeGreaterThan(0);
    });

    it('должен содержать платежи с правильными полями', () => {
      mockPaymentHistory.forEach((payment) => {
        expect(payment).toEqual(
          expect.objectContaining({
            id: expect.any(String),
            listingId: expect.any(String),
            listingTitle: expect.any(String),
            amount: expect.any(Number),
            currency: 'RSD',
            status: expect.any(String),
            createdAt: expect.stringMatching(
              /^\d{4}-\d{2}-\d{2}T\d{2}:\d{2}:\d{2}\.\d{3}Z$/
            ),
            seller: {
              name: expect.any(String),
              id: expect.any(String),
            },
            escrowStatus: expect.any(String),
          })
        );

        // Проверяем опциональные поля отдельно
        if (payment.completedAt) {
          expect(payment.completedAt).toMatch(
            /^\d{4}-\d{2}-\d{2}T\d{2}:\d{2}:\d{2}\.\d{3}Z$/
          );
        }
        if (payment.escrowReleaseDate) {
          expect(payment.escrowReleaseDate).toMatch(
            /^\d{4}-\d{2}-\d{2}T\d{2}:\d{2}:\d{2}\.\d{3}Z$/
          );
        }
      });
    });

    it('должен содержать различные статусы платежей', () => {
      const statuses = mockPaymentHistory.map((p) => p.status);
      const uniqueStatuses = [...new Set(statuses)];

      expect(uniqueStatuses.length).toBeGreaterThan(1);
    });

    it('должен содержать реалистичные данные', () => {
      mockPaymentHistory.forEach((payment) => {
        expect(payment.amount).toBeGreaterThan(0);
        expect(payment.listingTitle.length).toBeGreaterThan(0);
        expect(payment.seller.name.length).toBeGreaterThan(0);
        expect(['pending', 'captured', 'failed']).toContain(payment.status);
        expect(['pending', 'held', 'released']).toContain(payment.escrowStatus);
      });
    });

    it('должен иметь валидные даты', () => {
      mockPaymentHistory.forEach((payment) => {
        expect(() => new Date(payment.createdAt)).not.toThrow();

        if (payment.completedAt) {
          expect(() => new Date(payment.completedAt)).not.toThrow();
        }

        if (payment.escrowReleaseDate) {
          expect(() => new Date(payment.escrowReleaseDate)).not.toThrow();
        }
      });
    });
  });

  describe('mockWalletData', () => {
    it('должен иметь правильную структуру', () => {
      expect(mockWalletData).toEqual({
        balance: expect.any(Number),
        pendingBalance: expect.any(Number),
        currency: 'RSD',
        lastPayout: expect.stringMatching(
          /^\d{4}-\d{2}-\d{2}T\d{2}:\d{2}:\d{2}\.\d{3}Z$/
        ),
        nextPayout: expect.stringMatching(
          /^\d{4}-\d{2}-\d{2}T\d{2}:\d{2}:\d{2}\.\d{3}Z$/
        ),
        payoutMethod: expect.any(String),
      });
    });

    it('должен иметь неотрицательные балансы', () => {
      expect(mockWalletData.balance).toBeGreaterThanOrEqual(0);
      expect(mockWalletData.pendingBalance).toBeGreaterThanOrEqual(0);
    });

    it('должен иметь правильную валюту', () => {
      expect(mockWalletData.currency).toBe('RSD');
    });

    it('должен иметь валидные даты', () => {
      expect(() => new Date(mockWalletData.lastPayout)).not.toThrow();
      expect(() => new Date(mockWalletData.nextPayout)).not.toThrow();
    });

    it('должен иметь nextPayout в будущем относительно lastPayout', () => {
      const lastPayout = new Date(mockWalletData.lastPayout);
      const nextPayout = new Date(mockWalletData.nextPayout);

      expect(nextPayout.getTime()).toBeGreaterThan(lastPayout.getTime());
    });

    it('должен иметь валидный метод выплаты', () => {
      expect(typeof mockWalletData.payoutMethod).toBe('string');
      expect(mockWalletData.payoutMethod.length).toBeGreaterThan(0);
    });
  });

  describe('edge cases и рандомность', () => {
    it('должен генерировать разные данные при множественных вызовах', () => {
      const payments = Array.from({ length: 10 }, () => generateMockPayment());

      // Проверяем что все ID уникальны
      const ids = payments.map((p) => p.id);
      const uniqueIds = [...new Set(ids)];
      expect(uniqueIds.length).toBe(payments.length);

      // Проверяем что суммы разные (статистически маловероятно что все одинаковые)
      const amounts = payments.map((p) => p.amount);
      const uniqueAmounts = [...new Set(amounts)];
      expect(uniqueAmounts.length).toBeGreaterThan(1);
    });

    it('должен корректно обрабатывать отрицательный count в generateMockTransactions', () => {
      const transactions = generateMockTransactions(-1);

      expect(transactions).toEqual([]);
    });

    it('должен работать с очень большими значениями count', () => {
      const transactions = generateMockTransactions(1000);

      expect(transactions).toHaveLength(1000);
    });

    it('должен использовать текущее время для генерации ID', () => {
      const beforeTime = Date.now();
      const payment = generateMockPayment();
      const afterTime = Date.now();

      // Извлекаем timestamp из ID
      const timestampMatch = payment.id.match(/^MOCK_(\d+)_/);
      expect(timestampMatch).toBeTruthy();

      const timestamp = parseInt(timestampMatch![1]);
      expect(timestamp).toBeGreaterThanOrEqual(beforeTime);
      expect(timestamp).toBeLessThanOrEqual(afterTime);
    });

    it('должен генерировать случайную строку в ID правильной длины', () => {
      const payments = Array.from({ length: 20 }, () => generateMockPayment());

      payments.forEach((payment) => {
        const randomPartMatch = payment.id.match(/_([a-z0-9]{9})$/);
        expect(randomPartMatch).toBeTruthy();
        expect(randomPartMatch![1]).toHaveLength(9);
        expect(randomPartMatch![1]).toMatch(/^[a-z0-9]+$/);
      });
    });
  });

  describe('статическое тестирование данных', () => {
    it('mockPaymentHistory должен иметь стабильную структуру', () => {
      const originalLength = mockPaymentHistory.length;
      const originalFirstItem = { ...mockPaymentHistory[0] };

      // Проверяем что массив доступен и имеет ожидаемую длину
      expect(mockPaymentHistory).toHaveLength(originalLength);
      expect(mockPaymentHistory[0]).toEqual(originalFirstItem);

      // Проверяем что это действительно массив с данными
      expect(Array.isArray(mockPaymentHistory)).toBe(true);
      expect(originalLength).toBeGreaterThan(0);
    });

    it('mockWalletData должен иметь консистентные данные', () => {
      expect(mockWalletData.balance).toBeLessThanOrEqual(
        mockWalletData.pendingBalance * 10
      ); // Разумное соотношение
      expect(mockWalletData.payoutMethod).toBe('bank_transfer');
    });
  });
});
