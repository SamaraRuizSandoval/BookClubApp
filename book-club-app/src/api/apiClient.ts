import axios from 'axios';

export const api = axios.create({
  baseURL:
    'https://bookclub-backend.redwater-26f8bbd2.centralus.azurecontainerapps.io/',
  headers: {
    'Content-Type': 'application/json',
  },
});

api.interceptors.request.use((config) => {
  const token = localStorage.getItem('authToken');

  if (token) {
    config.headers.Authorization = `Bearer ${token}`;
  }

  return config;
});
