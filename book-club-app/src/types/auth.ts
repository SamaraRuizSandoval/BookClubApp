import { User } from './user';

export type AuthState = {
  token: AuthTokenResponse | null;
  user: User | null;
  isAuthenticated: boolean;
};

export type AuthTokenResponse = {
  auth_token: {
    token: string;
    expiry: string;
  };
};
