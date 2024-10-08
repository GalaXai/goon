import axiosInstance from '../config/axiosConfig';

export const loadImage = (base64Image) => {
  console.log('loadImage function called'); // Add this line
  return axiosInstance.post('/load-image', { image_base64: base64Image })
    .then(response => response.data)
    .catch(error => {
      console.error('Error loading image:', error);
      throw error;
    });
};