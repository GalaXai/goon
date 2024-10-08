import React from 'react';
import './ImageGrid.css';

const ImageGrid = ({ mainImage }) => {
  // Create an array of 8 elements, first are the provided images, rest are placeholders
      const imageArray = [mainImage, ...Array(8).fill(null)].slice(0, 8);

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