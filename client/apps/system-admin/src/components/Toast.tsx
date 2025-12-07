import { Snackbar, Alert, AlertColor } from '@mui/material';
import { create } from 'zustand';

interface Toast {
    id: string;
    message: string;
    severity: AlertColor;
}

interface ToastState {
    toasts: Toast[];
    addToast: (message: string, severity?: AlertColor) => void;
    removeToast: (id: string) => void;
}

export const useToastStore = create<ToastState>((set) => ({
    toasts: [],
    addToast: (message, severity = 'info') => {
        const id = Date.now().toString();
        set((state) => ({
            toasts: [...state.toasts, { id, message, severity }],
        }));
        // Auto remove after 5 seconds
        setTimeout(() => {
            set((state) => ({
                toasts: state.toasts.filter((t) => t.id !== id),
            }));
        }, 5000);
    },
    removeToast: (id) =>
        set((state) => ({
            toasts: state.toasts.filter((t) => t.id !== id),
        })),
}));

export const ToastContainer: React.FC = () => {
    const { toasts, removeToast } = useToastStore();

    return (
        <>
            {toasts.map((toast, index) => (
                <Snackbar
                    key={toast.id}
                    open={true}
                    anchorOrigin={{ vertical: 'bottom', horizontal: 'right' }}
                    sx={{ bottom: { xs: 24 + index * 60 } }}
                >
                    <Alert
                        severity={toast.severity}
                        onClose={() => removeToast(toast.id)}
                        sx={{ width: '100%' }}
                    >
                        {toast.message}
                    </Alert>
                </Snackbar>
            ))}
        </>
    );
};
