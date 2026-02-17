export type Role = 'admin' | 'user';

export type User = {
  id: number;
  username: string;
  email: string;
  role: Role;
};
