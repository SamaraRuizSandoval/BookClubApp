import { createContext, useContext, useEffect, useState } from 'react';

import { getUserBookStats } from '../api/userBooksApi';
import { UserBookStats } from '../types/userBooks';

import { useAuth } from './AuthContext';

type StatsContextType = {
  stats: UserBookStats | null;
  refreshStats: () => Promise<void>;
};

const UserStatsContext = createContext<StatsContextType | undefined>(undefined);

export function UserStatsProvider({ children }: { children: React.ReactNode }) {
  const { auth } = useAuth();
  const user = auth.user;

  const [stats, setStats] = useState<UserBookStats | null>(null);

  async function loadStats() {
    if (!user) return;

    try {
      const data = await getUserBookStats();
      setStats(data);
    } catch (err) {
      console.error(err);
    }
  }

  useEffect(() => {
    if (user) {
      loadStats();
    } else {
      setStats(null); // reset on logout
    }
  }, [user]);

  return (
    <UserStatsContext.Provider
      value={{
        stats,
        refreshStats: loadStats,
      }}
    >
      {children}
    </UserStatsContext.Provider>
  );
}

export function useUserStats() {
  const context = useContext(UserStatsContext);
  if (!context) {
    throw new Error('useUserStats must be used inside UserStatsProvider');
  }
  return context;
}
