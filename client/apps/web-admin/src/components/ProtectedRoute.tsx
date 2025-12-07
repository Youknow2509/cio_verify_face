import { Navigate } from 'react-router-dom';
import { useAuthStore } from '@/stores/authStore';

interface ProtectedRouteProps {
  children: React.ReactNode;
}

export const ProtectedRoute: React.FC<ProtectedRouteProps> = ({ children }) => {
  const isAuthenticated = useAuthStore((state) => state.isAuthenticated);

  // Auth bypass disabled for production
  const bypassAuth = false;

  if (!isAuthenticated && !bypassAuth) {
    return <Navigate to="/login" replace />;
  }

  return <>{children}</>;
};
