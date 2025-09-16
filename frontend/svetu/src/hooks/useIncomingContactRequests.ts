import { useState, useEffect, useCallback } from 'react';
import { contactsService } from '@/services/contacts';
import { useAuth } from '@/contexts/AuthContext';

interface IncomingRequest {
  user_id: number;
  contact_user_id: number;
  status: string;
  created_at: string;
  user?: {
    id: number;
    name: string;
    email: string;
  };
}

export function useIncomingContactRequests() {
  const { user } = useAuth();
  const [requests, setRequests] = useState<IncomingRequest[]>([]);
  const [isLoading, setIsLoading] = useState(true);
  const [requestsByUserId, setRequestsByUserId] = useState<
    Record<number, IncomingRequest>
  >({});

  const loadRequests = useCallback(async () => {
    if (!user) {
      setRequests([]);
      setRequestsByUserId({});
      setIsLoading(false);
      return;
    }

    try {
      setIsLoading(true);
      const response = await contactsService.getIncomingRequests({
        page: 1,
        limit: 100,
      });

      setRequests(response.contacts);

      // Создаем мапу по user_id для быстрого доступа
      const byUserId: Record<number, IncomingRequest> = {};
      response.contacts.forEach((req) => {
        byUserId[req.user_id] = req;
      });
      setRequestsByUserId(byUserId);
    } catch (error) {
      console.error('Error loading incoming contact requests:', error);
      setRequests([]);
      setRequestsByUserId({});
    } finally {
      setIsLoading(false);
    }
  }, [user]);

  useEffect(() => {
    loadRequests();

    // Обновляем каждую минуту
    const interval = setInterval(loadRequests, 60000);

    return () => clearInterval(interval);
  }, [user?.id, loadRequests]);

  const hasRequestFromUser = (userId: number): boolean => {
    return !!requestsByUserId[userId];
  };

  const getRequestFromUser = (userId: number): IncomingRequest | null => {
    return requestsByUserId[userId] || null;
  };

  const removeRequest = (userId: number) => {
    setRequests((prev) => prev.filter((r) => r.user_id !== userId));
    setRequestsByUserId((prev) => {
      const newMap = { ...prev };
      delete newMap[userId];
      return newMap;
    });
  };

  return {
    requests,
    isLoading,
    hasRequestFromUser,
    getRequestFromUser,
    removeRequest,
    totalCount: requests.length,
    reload: loadRequests,
  };
}
