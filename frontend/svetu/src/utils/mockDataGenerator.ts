import type { MockPayment } from '@/types/payment';

export const generateMockPayment = (): MockPayment => ({
  id: `MOCK_${Date.now()}_${Math.random().toString(36).substr(2, 9)}`,
  merchantTransactionId: `MTX_${Date.now()}`,
  amount: Math.floor(Math.random() * 100000) + 1000,
  currency: 'RSD',
  status: 'pending',
  createdAt: new Date().toISOString(),
  card: {
    lastFour: Math.floor(Math.random() * 10000)
      .toString()
      .padStart(4, '0'),
    brand: ['visa', 'mastercard', 'maestro'][Math.floor(Math.random() * 3)],
    expiryMonth: Math.floor(Math.random() * 12) + 1,
    expiryYear: new Date().getFullYear() + Math.floor(Math.random() * 5) + 1,
  },
});

export const generateMockTransactions = (count: number): MockPayment[] => {
  return Array.from({ length: count }, () => ({
    ...generateMockPayment(),
    status: ['pending', 'authorized', 'captured', 'failed', 'refunded'][
      Math.floor(Math.random() * 5)
    ],
    completedAt: Math.random() > 0.5 ? new Date().toISOString() : undefined,
  }));
};

export const mockPaymentHistory = [
  {
    id: 'MOCK_1234567890',
    listingId: '123',
    listingTitle: 'iPhone 13 Pro Max',
    amount: 125000,
    currency: 'RSD',
    status: 'captured',
    createdAt: new Date(Date.now() - 86400000).toISOString(),
    completedAt: new Date(Date.now() - 86000000).toISOString(),
    seller: {
      name: 'TechStore Belgrade',
      id: 'seller_123',
    },
    escrowStatus: 'held',
    escrowReleaseDate: new Date(Date.now() + 518400000).toISOString(),
  },
  {
    id: 'MOCK_0987654321',
    listingId: '456',
    listingTitle: 'MacBook Air M2',
    amount: 180000,
    currency: 'RSD',
    status: 'pending',
    createdAt: new Date(Date.now() - 3600000).toISOString(),
    seller: {
      name: 'Apple Store Novi Sad',
      id: 'seller_456',
    },
    escrowStatus: 'pending',
  },
];

export const mockWalletData = {
  balance: 45000,
  pendingBalance: 125000,
  currency: 'RSD',
  lastPayout: new Date(Date.now() - 604800000).toISOString(),
  nextPayout: new Date(Date.now() + 172800000).toISOString(),
  payoutMethod: 'bank_transfer',
};
