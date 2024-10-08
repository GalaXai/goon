import axiosInstance from '../config/axiosConfig';

export const fetchImage = () => {
  return axiosInstance.get('/image')
    .then(response => response.data)
    .catch(error => {
      console.error('Error fetching image:', error);
      throw error;
    });
};

export const createImage = (data) => {
  return axiosInstance.post('/image', data)
    .then(response => response.data)
    .catch(error => {
      console.error('Error creating image:', error);
      throw error;
    });
};