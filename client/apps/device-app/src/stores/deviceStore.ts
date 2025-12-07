import { create } from 'zustand';
import { persist } from 'zustand/middleware';

interface DeviceInfo {
    deviceId: string;
    deviceName: string;
    companyId: string;
    companyName: string;
    location?: string;
}

interface DeviceState {
    token: string | null;
    isAuthenticated: boolean;
    deviceInfo: DeviceInfo | null;
    setDeviceToken: (token: string) => void;
    setDeviceInfo: (info: DeviceInfo) => void;
    clearDevice: () => void;
}

export const useDeviceStore = create<DeviceState>()(
    persist(
        (set) => ({
            token: null,
            isAuthenticated: false,
            deviceInfo: null,
            setDeviceToken: (token: string) => {
                localStorage.setItem('device_token', token);
                set({ token, isAuthenticated: true });
            },
            setDeviceInfo: (info: DeviceInfo) => {
                set({ deviceInfo: info });
            },
            clearDevice: () => {
                localStorage.removeItem('device_token');
                set({ token: null, isAuthenticated: false, deviceInfo: null });
            },
        }),
        {
            name: 'device-storage',
        }
    )
);
