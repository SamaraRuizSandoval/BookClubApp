import { IonToast } from '@ionic/react';
import { createContext, useContext, useState, ReactNode } from 'react';
type ToastContextType = {
  show: (message: string, color?: 'success' | 'danger' | 'primary') => void;
};

const ToastContext = createContext<ToastContextType | undefined>(undefined);

export function ToastProvider({ children }: { children: ReactNode }) {
  const [isOpen, setIsOpen] = useState(false);
  const [message, setMessage] = useState('');
  const [color, setColor] = useState<'success' | 'danger' | 'primary'>(
    'success',
  );

  const show = (msg: string, c: typeof color = 'success') => {
    setMessage(msg);
    setColor(c);
    setIsOpen(true);
  };

  return (
    <ToastContext.Provider value={{ show }}>
      {children}
      <IonToast
        isOpen={isOpen}
        message={message}
        color={color}
        duration={3000}
        onDidDismiss={() => setIsOpen(false)}
        position="top"
        className="success-toast"
      />
    </ToastContext.Provider>
  );
}

export const useToast = () => {
  const ctx = useContext(ToastContext);
  if (!ctx) throw new Error('useToast must be used within ToastProvider');
  return ctx;
};
