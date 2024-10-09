import axiosInstance from '../config/axiosConfig';

export const loadImage = (base64Image, gradientThreshold) => {
  console.log('loadImage function called');
  return axiosInstance.post('/load-image', {
    base64Image: base64Image,
    gradientThreshold: gradientThreshold
  })
    .then(response => {
      return response.data;
    })
    .catch(error => {
      console.error('Error loading image:', error);
      throw error;
    });
};