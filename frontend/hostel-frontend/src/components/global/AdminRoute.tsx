// frontend/hostel-frontend/src/components/global/AdminRoute.tsx
import React, { ReactNode } from 'react';
import { Navigate } from 'react-router-dom';
import { useAuth } from '../../contexts/AuthContext';
import { isAdmin } from '../../utils/adminUtils';

interface AdminRouteProps {
  children: ReactNode;
}

const AdminRoute: React.FC<AdminRouteProps> = ({ children }) => {
  const { user } = useAuth();

  // Check if user is authenticated
  if (!user) {
    return <Navigate to="/" replace />;
  }

  // Check if user's email is in the admin list
  if (!isAdmin(user.email)) {
    return <Navigate to="/" replace />;
  }

  return <>{children}</>;
};

export default AdminRoute;