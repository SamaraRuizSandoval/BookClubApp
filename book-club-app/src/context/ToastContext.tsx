import { createContext, useContext, useState, ReactNode } from 'react';
import { useEffect, useRef } from 'react';

import '../styles/toast.css';

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
  const timeouts = useRef<number[]>([]);

  useEffect(() => {
    return () => {
      timeouts.current.forEach((t) => clearTimeout(t));
    };
  }, []);

  const show = (msg: string, color: Toast['color'] = 'success') => {
    const id = Date.now();

    setToasts((prev) => [...prev, { id, message: msg, color }]);

    // Track the timeout so it can be cleared on unmount
    const timeout = window.setTimeout(() => {
      setToasts((prev) => prev.filter((t) => t.id !== id));
      // Remove this timeout from the ref after it fires
      timeouts.current = timeouts.current.filter((t) => t !== timeout);
    }, 3000);

    timeouts.current.push(timeout);
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
