import axiosInstance from '../config/axiosConfig';

export const loadImage = (base64Image, gradientThreshold) => {
  console.log('loadImage function called');
  return axiosInstance.post('/edge-detect-ascii', {
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

export const colorDownsample = (base64Image) => {
  console.log('loadImageColor function called');
  return axiosInstance.post('/color-downsample', {
    base64Image: base64Image
  })
    .then(response => {
      return response.data;
    })
    .catch(error => {
      console.error('Error loading color image:', error);
      throw error;
    });
};