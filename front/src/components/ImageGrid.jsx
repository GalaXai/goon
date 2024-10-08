import React from 'react';
import './ImageGrid.css';
import { loadImage } from '../services/imageApi';

const ImageGrid = ({ mainImage }) => {
  // Create an array of 8 elements, first are the provided images, rest are placeholders
  const imageArray = [mainImage, ...Array(8).fill(null)].slice(0, 8);

  loadImage(mainImage)
  .then(result => {
    console.log('Image matrix:', result.matrix);
    console.log('Image shape:', result.shape);
    // Process the image matrix and shape as needed
  })
  .catch(error => {
    console.error('Failed to load image:', error);
    // Handle the error appropriately
  });

  return (
    <div className="image-grid">
      {imageArray.map((image, index) => (
        <div key={index} className="image-container">
          <h3>Image {index + 1}</h3>
          {image ? (
            <img src={image} alt={`Converted ${index + 1}`} className="grid-image" />
          ) : (
            <div className="placeholder">Placeholder</div>
          )}
        </div>
      ))}
    </div>
  );
};

export default ImageGrid;