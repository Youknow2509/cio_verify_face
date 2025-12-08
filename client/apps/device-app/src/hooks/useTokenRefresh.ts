import { useEffect, useRef } from 'react';
import { useDeviceStore } from '@/stores/deviceStore';
import { deviceApi } from '@/services/api';

/**
 * Hook to automatically refresh device token every 50 minutes
 * Tokens typically expire after 1 hour, so we refresh at 50 minutes
 */
export const useTokenRefresh = () => {
    const { token, setDeviceToken } = useDeviceStore();
    const refreshIntervalRef = useRef<ReturnType<typeof setInterval> | null>(null);

    useEffect(() => {
        if (!token) return;

        const refreshToken = async () => {
            try {
                const response: any = await deviceApi.refreshToken();
                
                if (response.code === 200 && response.data?.device_token) {
                    console.log('Token refreshed successfully');
                    setDeviceToken(response.data.device_token);
                }
            } catch (error) {
                console.error('Failed to refresh token:', error);
                // Don't throw error, just log it. The app will handle expired tokens via API interceptor
            }
        };

        // Refresh token every 50 minutes (3000000 ms)
        refreshIntervalRef.current = setInterval(refreshToken, 50 * 60 * 1000);

        // Cleanup on unmount
        return () => {
            if (refreshIntervalRef.current) {
                clearInterval(refreshIntervalRef.current);
            }
        };
    }, [token, setDeviceToken]);
};
