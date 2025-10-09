import {
  createContext,
  useCallback,
  useContext,
  useMemo,
  useRef,
  useState,
  type PropsWithChildren,
} from 'react';
import { ConfirmModal } from '@/ui/ConfirmModal';
import { ToastCenter } from '@/ui/ToastCenter';

type ToastVariant = 'primary' | 'success' | 'warning' | 'danger' | 'info';

export interface ToastMessage {
  id: string;
  title?: string;
  message: string;
  variant: ToastVariant;
}

export interface ConfirmDialogOptions {
  title: string;
  message: string;
  confirmLabel?: string;
  cancelLabel?: string;
  confirmVariant?: ToastVariant;
}

interface UiContextValue {
  showToast: (toast: Omit<ToastMessage, 'id'>) => void;
  dismissToast: (id: string) => void;
  confirm: (options: ConfirmDialogOptions) => Promise<boolean>;
}

const UiContext = createContext<UiContextValue | undefined>(undefined);

export function UiProvider({ children }: PropsWithChildren) {
  const [toasts, setToasts] = useState<ToastMessage[]>([]);
  const [confirmOptions, setConfirmOptions] = useState<ConfirmDialogOptions | null>(null);
  const resolverRef = useRef<(value: boolean) => void>();

  const showToast = useCallback((toast: Omit<ToastMessage, 'id'>) => {
    const id = `${Date.now()}-${Math.random().toString(16).slice(2)}`;
    setToasts((prev) => [...prev, { ...toast, id }]);
  }, []);

  const dismissToast = useCallback((id: string) => {
    setToasts((prev) => prev.filter((item) => item.id !== id));
  }, []);

  const confirm = useCallback((options: ConfirmDialogOptions) => {
    setConfirmOptions(options);
    return new Promise<boolean>((resolve) => {
      resolverRef.current = resolve;
    });
  }, []);

  const handleConfirm = useCallback(
    (result: boolean) => {
      setConfirmOptions(null);
      resolverRef.current?.(result);
      resolverRef.current = undefined;
    },
    []
  );

  const value = useMemo<UiContextValue>(
    () => ({
      showToast,
      dismissToast,
      confirm,
    }),
    [showToast, dismissToast, confirm]
  );

  return (
    <UiContext.Provider value={value}>
      {children}
      <ToastCenter toasts={toasts} onClose={dismissToast} />
      <ConfirmModal
        show={Boolean(confirmOptions)}
        title={confirmOptions?.title ?? ''}
        message={confirmOptions?.message ?? ''}
        confirmLabel={confirmOptions?.confirmLabel}
        cancelLabel={confirmOptions?.cancelLabel}
        confirmVariant={confirmOptions?.confirmVariant}
        onConfirm={() => handleConfirm(true)}
        onCancel={() => handleConfirm(false)}
      />
    </UiContext.Provider>
  );
}

export function useUi(): UiContextValue {
  const context = useContext(UiContext);
  if (!context) {
    throw new Error('useUi must be used within UiProvider');
  }
  return context;
}
