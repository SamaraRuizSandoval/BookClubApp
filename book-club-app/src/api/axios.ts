import axios from 'axios';

const api = axios.create({
  baseURL:
    'https://bookclub-backend.redwater-26f8bbd2.centralus.azurecontainerapps.io/',
  //baseURL: import.meta.env.VITE_API_BASE_URL,
  headers: {
    'Content-Type': 'application/json',
  },
});

export default api;
