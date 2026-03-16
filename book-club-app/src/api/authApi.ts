import { api } from '../api/apiClient';
import { AuthTokenResponse } from '../types/auth';
import { User } from '../types/user';

export const getCurrentUser = async () => {
  const response = await api.get<User>('/me');
  return response.data;
};

export const loginUser = async (username: string, password: string) => {
  const authResponse = await api.post<AuthTokenResponse>(
    '/tokens/authentication',
    {
      username,
      password,
    },
  );

  const { token, expiry } = authResponse.data.auth_token;

  const userResponse = await api.get<User>('/me', {
    headers: {
      Authorization: `Bearer ${token}`,
    },
  });

  return {
    token,
    expiry,
    user: userResponse.data,
  };
};
