'use client';

import { useEffect } from 'react';
import { useDispatch } from 'react-redux';
import { AppDispatch } from '@/store/store';
import { checkAuth, setAuthFromSession } from '@/store/authSlice';

export function AuthProvider({ children }: { children: React.ReactNode }) {
  const dispatch = useDispatch<AppDispatch>();

  useEffect(() => {
    // Load session from localStorage immediately
    if (typeof window !== 'undefined') {
      const savedSession = localStorage.getItem('user_session');
      if (savedSession) {
        try {
          const userData = JSON.parse(savedSession);
          dispatch(setAuthFromSession({ user: userData }));
        } catch (error) {
          console.error('Error parsing saved session:', error);
        }
      }
    }

    // Then verify with server
    dispatch(checkAuth());
  }, [dispatch]);

  return <>{children}</>;
}