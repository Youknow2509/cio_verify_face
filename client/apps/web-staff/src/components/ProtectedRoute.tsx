import { Navigate, Outlet } from 'react-router-dom';
import { useAuthStore } from '@/stores/authStore';

export const ProtectedRoute: React.FC = () => {
    const { isAuthenticated } = useAuthStore();

    // Auth bypass disabled for production
    const bypassAuth = false;

    if (!isAuthenticated && !bypassAuth) {
        return <Navigate to="/login" replace />;
    }

    return <Outlet />;
};
