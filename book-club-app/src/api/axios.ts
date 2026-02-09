import axios from 'axios';

const api = axios.create({
  //baseURL: "https://bookclub-backend.redwater-26f8bbd2.centralus.azurecontainerapps.io/",
  baseURL: 'http://localhost:5000/',
  headers: {
    'Content-Type': 'application/json',
  },
});

export default api;
