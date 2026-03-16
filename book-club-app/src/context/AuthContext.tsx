import React, { createContext, useContext, useEffect, useState } from 'react';

import { getCurrentUser } from '../api/authApi';
import { AuthState } from '../types/auth';
import { User } from '../types/user';

type AuthContextType = {
  auth: AuthState;
  login: (token: string, expiry: string, user: User) => void;
  logout: () => void;
  initializing: boolean;
};

const AuthContext = createContext<AuthContextType | undefined>(undefined);

export const AuthProvider: React.FC<{ children: React.ReactNode }> = ({
  children,
}) => {
  const [auth, setAuth] = useState<AuthState>({
    token: '',
    user: null,
    isAuthenticated: false,
  });

  const [initializing, setInitializing] = useState(true);

  // 🔄 Restore session on app load
  useEffect(() => {
    const restoreSession = async () => {
      const storedToken = localStorage.getItem('authToken');
      const storedExpiry = localStorage.getItem('authExpiry');

      if (!storedToken || !storedExpiry) {
        setInitializing(false);
        return;
      }

      const expiryDate = new Date(storedExpiry);

      if (expiryDate <= new Date()) {
        localStorage.removeItem('authToken');
        localStorage.removeItem('authExpiry');
        localStorage.removeItem('user');
        setInitializing(false);
        return;
      }

      try {
        const user = await getCurrentUser();

        setAuth({
          token: storedToken,
          user,
          isAuthenticated: true,
        });
      } catch (error) {
        localStorage.removeItem('authToken');
      } finally {
        setInitializing(false);
      }
    };

    restoreSession();
  }, []);

  const login = (token: string, expiry: string, user: User) => {
    localStorage.setItem('authToken', token);
    localStorage.setItem('authExpiry', expiry);
    localStorage.setItem('user', JSON.stringify(user));

    setAuth({
      token,
      user,
      isAuthenticated: true,
    });
  };

  const logout = () => {
    localStorage.removeItem('authToken');

    setAuth({
      token: '',
      user: null,
      isAuthenticated: false,
    });
  };

  return (
    <AuthContext.Provider value={{ auth, login, logout, initializing }}>
      {children}
    </AuthContext.Provider>
  );
};

// Custom hook
export const useAuth = () => {
  const context = useContext(AuthContext);
  if (!context) {
    throw new Error('useAuth must be used within AuthProvider');
  }
  return context;
};
