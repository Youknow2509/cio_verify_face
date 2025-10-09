import Toast from 'react-bootstrap/Toast';
import ToastContainer from 'react-bootstrap/ToastContainer';
import type { ToastMessage } from '@/app/providers/UiProvider';

interface ToastCenterProps {
  toasts: ToastMessage[];
  onClose: (id: string) => void;
}

export function ToastCenter({ toasts, onClose }: ToastCenterProps) {
  return (
    <ToastContainer position="top-end" className="p-3">
      {toasts.map((toast) => (
        <Toast
          key={toast.id}
          bg={toast.variant === 'primary' ? 'primary' : undefined}
          onClose={() => onClose(toast.id)}
          delay={3500}
          autohide
        >
          {toast.title && (
            <Toast.Header closeButton={false}>
              <strong className="me-auto">{toast.title}</strong>
              <small>Vừa xong</small>
            </Toast.Header>
          )}
          <Toast.Body className={toast.variant === 'primary' ? 'text-white' : undefined}>
            {toast.message}
          </Toast.Body>
        </Toast>
      ))}
    </ToastContainer>
  );
}
