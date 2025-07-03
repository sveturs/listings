export interface PaymentRequest {
  listing_id: string;
  amount: number;
  currency: string;
  buyer_info: BuyerInfo;
  return_url: string;
  locale?: string;
}

export interface BuyerInfo {
  name: string;
  email: string;
  phone?: string;
  address?: string;
}

export interface PaymentResponse {
  id: string;
  redirectUrl: string;
  status: string;
}

export interface PaymentStatus {
  id: string;
  status:
    | 'pending'
    | 'processing'
    | 'authorized'
    | 'captured'
    | 'failed'
    | 'cancelled';
  amount: number;
  currency: string;
  createdAt: string;
  completedAt?: string;
  error_code?: string;
  card?: {
    lastFour: string;
    brand: string;
    expiryMonth: number;
    expiryYear: number;
  };
}

export interface MockPayment {
  id: string;
  merchantTransactionId: string;
  amount: number;
  currency: string;
  status: string;
  createdAt: string;
  completedAt?: string;
  card?: {
    lastFour: string;
    brand: string;
    expiryMonth: number;
    expiryYear: number;
  };
  listing_id?: string;
  buyer_info?: BuyerInfo;
}

export interface IPaymentService {
  createPayment(data: PaymentRequest): Promise<PaymentResponse>;
  getPaymentStatus(paymentId: string): Promise<PaymentStatus>;
  handle3DSecure?(paymentId: string, code: string): Promise<boolean>;
  simulateWebhook?(paymentId: string, status: string): Promise<void>;
}
