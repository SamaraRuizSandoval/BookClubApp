import React, { createContext, useContext, useEffect, useState } from 'react';

import { AuthState } from '../types/auth';

type AuthContextType = {
  auth: AuthState;
  login: (token: string, user: any) => void;
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
    const storedToken = localStorage.getItem('authToken');

    if (!storedToken) {
      setInitializing(false);
      return;
    }

    fetch(`${import.meta.env.VITE_API_BASE_URL}/me`, {
      headers: {
        Authorization: `Bearer ${storedToken}`,
      },
    })
      .then((res) => {
        if (!res.ok) throw new Error('Invalid token');
        return res.json();
      })
      .then((user) => {
        setAuth({
          token: storedToken,
          user,
          isAuthenticated: true,
        });
      })
      .catch(() => {
        localStorage.removeItem('authToken');
      })
      .finally(() => {
        setInitializing(false);
      });
  }, []);

  const login = (token: string, user: any) => {
    localStorage.setItem('authToken', token);

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
