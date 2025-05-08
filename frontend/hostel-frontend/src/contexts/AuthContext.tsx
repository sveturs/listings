// frontend/hostel-frontend/src/contexts/AuthContext.tsx

import React, { createContext, useState, useContext, useEffect, ReactNode } from 'react';
import axios from '../api/axios';

// Define the User interface
export interface User {
  id: number;
  name: string;
  email: string;
  avatar?: string;
  is_admin?: boolean;
  city?: string;
  country?: string;
  phone?: string;
  [key: string]: any; // For any additional properties
}

// Define the Authentication Context interface
export interface AuthContextType {
  user: User | null;
  loading: boolean;
  login: (params?: string) => void;
  logout: () => Promise<void>;
  checkAuth: () => Promise<void>;
}

// Define AuthProvider Props interface
interface AuthProviderProps {
  children: ReactNode;
}

// Create the Auth Context with initial null value
const AuthContext = createContext<AuthContextType | null>(null);

export const AuthProvider: React.FC<AuthProviderProps> = ({ children }) => {
  const [user, setUser] = useState<User | null>(null);
  const [isLoading, setIsLoading] = useState<boolean>(true);

  // Save session to localStorage
  const saveSession = (userData: User): void => {
    localStorage.setItem('user_session', JSON.stringify(userData));
  };

  // Load session from localStorage
  const loadSession = (): User | null => {
    try {
      const session = localStorage.getItem('user_session');
      if (session) {
        return JSON.parse(session);
      }
    } catch (error) {
      console.error('Error loading session:', error);
    }
    return null;
  };

  // Check authentication status
  const checkAuth = async (): Promise<void> => {
    try {
      // First check URL for session token
      const urlParams = new URLSearchParams(window.location.search);
      const sessionTokenFromUrl = urlParams.get('session_token');
      
      // If token found in URL, save it
      if (sessionTokenFromUrl) {
        localStorage.setItem('user_session_token', sessionTokenFromUrl);
      }
  
      // Then check local session
      const savedSession = loadSession();
      if (savedSession) {
        setUser(savedSession);
      }
  
      // Then verify with server
      const response = await axios.get('/auth/session');
      if (response.data.authenticated) {
        setUser(response.data.user);
        saveSession(response.data.user);
      } else {
        // If server says user is not authenticated, clear local session
        localStorage.removeItem('user_session');
        setUser(null);
      }
    } catch (error) {
      console.error('Error checking auth status:', error);
      // In case of error, don't remove local session - it might be just a network issue
    } finally {
      setIsLoading(false);
    }
  };

  useEffect(() => {
    // On load, immediately set user from localStorage
    const savedSession = loadSession();
    if (savedSession) {
      setUser(savedSession);
      setIsLoading(false);
    }
    
    // Then check authentication status
    checkAuth();
  }, []);

  const login = (params: string = ''): void => {
    // Use type assertion to access ENV property
    const backendUrl = (window as any).ENV?.REACT_APP_BACKEND_URL || '';
    const authUrl = (window as any).ENV?.REACT_APP_AUTH_URL || '/auth';
    window.location.href = `${backendUrl}${authUrl}/google${params}`;
  };

  const logout = async (): Promise<void> => {
    try {
      await axios.get('/auth/logout', { withCredentials: true });
      localStorage.removeItem('user_session');
      setUser(null);
    } catch (error) {
      console.error('Logout failed:', error);
      // Still clear local session
      localStorage.removeItem('user_session');
      setUser(null);
    }
  };

  const value: AuthContextType = {
    user,
    loading: isLoading,
    login,
    logout,
    checkAuth // Export method to call after successful payment
  };

  return (
    <AuthContext.Provider value={value}>
      {children}
    </AuthContext.Provider>
  );
};

export const useAuth = (): AuthContextType => {
  const context = useContext(AuthContext);
  if (!context) {
    throw new Error('useAuth must be used within an AuthProvider');
  }
  return context;
};