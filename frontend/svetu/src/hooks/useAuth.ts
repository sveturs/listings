import { useSelector, useDispatch } from 'react-redux';
import { useRouter } from 'next/navigation';
import { RootState, AppDispatch } from '@/store/store';
import { 
  loginUser, 
  registerUser, 
  logoutUser, 
  updateUser,
  googleLogin,
  clearError 
} from '@/store/authSlice';
import { LoginRequest, RegisterRequest, User } from '@/types/auth';

export const useAuth = () => {
  const router = useRouter();
  const dispatch = useDispatch<AppDispatch>();
  
  const { user, isLoading, isAuthenticated, error } = useSelector(
    (state: RootState) => state.auth
  );

  const login = async (credentials: LoginRequest) => {
    try {
      const result = await dispatch(loginUser(credentials)).unwrap();
      router.push('/');
      return result;
    } catch (error) {
      console.error('Login failed:', error);
      throw error;
    }
  };

  const loginWithGoogle = (params?: string) => {
    dispatch(googleLogin(params));
  };

  const register = async (data: RegisterRequest) => {
    try {
      const result = await dispatch(registerUser(data)).unwrap();
      router.push('/');
      return result;
    } catch (error) {
      console.error('Registration failed:', error);
      throw error;
    }
  };

  const logout = async () => {
    await dispatch(logoutUser()).unwrap();
    router.push('/');
  };

  const updateProfile = (userData: Partial<User>) => {
    if (user) {
      dispatch(updateUser({ ...user, ...userData } as User));
    }
  };

  const clearAuthError = () => {
    dispatch(clearError());
  };

  return {
    user,
    isLoading,
    isAuthenticated,
    error,
    login,
    loginWithGoogle,
    register,
    logout,
    updateProfile,
    clearAuthError,
  };
};