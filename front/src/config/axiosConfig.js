import axios from 'axios';


const getEnvVariable = (key) => {
    if (import.meta.env) {
      return import.meta.env[key];
    }
    if (typeof process !== 'undefined' && process.env) {
      return process.env[key];
    }
    return undefined;
  };

  const baseURL = getEnvVariable('VITE_API_BASE_URL') || 'http://localhost:3000';

console.log('Environment Variables:');
console.log('VITE_API_BASE_URL:', getEnvVariable('VITE_API_BASE_URL'));

const axiosInstance = axios.create({
  baseURL: baseURL,
  // timeout: 5000,
  headers: {
    'Content-Type': 'application/json',
    'Accept': 'application/json',
  },    
});

// Request interceptor
axiosInstance.interceptors.request.use(
  (config) => {
    // You can add any request modifications here, such as adding auth tokens
    return config;
  },
  (error) => {m
    return Promise.reject(error);
  }
);

// Response interceptor
axiosInstance.interceptors.response.use(
  (response) => {
    return response;
  },
  (error) => {
    // Handle errors globally
    console.error('API Error:', error);
    return Promise.reject(error);
  }
);

export default axiosInstance;