import '../styles/toast.css';
import { createContext, useContext, useState, ReactNode } from 'react';

type Toast = {
  id: number;
  message: string;
  color: 'success' | 'danger' | 'primary';
};

type ToastContextType = {
  show: (message: string, color?: Toast['color']) => void;
};

const ToastContext = createContext<ToastContextType | undefined>(undefined);

export function ToastProvider({ children }: { children: ReactNode }) {
  const [toasts, setToasts] = useState<Toast[]>([]);

  const show = (msg: string, color: Toast['color'] = 'success') => {
    const id = Date.now();

    setToasts((prev) => [...prev, { id, message: msg, color }]);

    setTimeout(() => {
      setToasts((prev) => prev.filter((t) => t.id !== id));
    }, 3000);
  };

  return (
    <ToastContext.Provider value={{ show }}>
      {children}

      <div className="toast-container">
        {toasts.map((toast) => (
          <div key={toast.id} className="toast">
            <span className={`toast-dot ${toast.color}`} />
            {toast.message}
          </div>
        ))}
      </div>
    </ToastContext.Provider>
  );
}

export const useToast = () => {
  const ctx = useContext(ToastContext);
  if (!ctx) throw new Error('useToast must be used within ToastProvider');
  return ctx;
};
