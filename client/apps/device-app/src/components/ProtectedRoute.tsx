import { Navigate } from 'react-router-dom';
import { useDeviceStore } from '@/stores/deviceStore';

interface ProtectedRouteProps {
    children: React.ReactNode;
}

export const ProtectedRoute: React.FC<ProtectedRouteProps> = ({ children }) => {
    const { isAuthenticated } = useDeviceStore();

    // Auth bypass disabled for production
    const bypassAuth = false;

    if (!isAuthenticated && !bypassAuth) {
        return <Navigate to="/token-auth" replace />;
    }

    return <>{children}</>;
};
