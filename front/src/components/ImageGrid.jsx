import React, { useState, useEffect } from 'react';
import './ImageGrid.css';
import { loadImage } from '../services/imageApi';

const ImageGrid = ({ mainImage, gradientThreshold}) => {
  const [imageArray, setImageArray] = useState(Array(8).fill(null));

  useEffect(() => {
    if (mainImage) {
      loadImage(mainImage, gradientThreshold)
        .then(result => {
          console.log('Images loaded successfully');
          setImageArray([
            { src: `data:image/png;base64,${result.originalImage}`, label: 'Original' },
            { src: `data:image/png;base64,${result.desaturatedImage}`, label: 'Desaturated' },
            { src: `data:image/png;base64,${result.downsampledImage}`, label: 'Downsampled' },
            { src: `data:image/png;base64,${result.gaussiansDiffImage}`, label: 'Gaussians Diff' },
            { src: `data:image/png;base64,${result.sobelImage}`, label: 'Sobel' },
            { src: `data:image/png;base64,${result.gradientImage}`, label: 'Gradient' },
            // Fill remaining slots with null
            ...Array(2).fill(null)
          ]);
          console.log(imageArray)
        })
        .catch(error => {
          console.error('Failed to load image:', error);
          // Handle the error appropriately
        });
    }
  }, [mainImage, gradientThreshold]);

return (
  <div className="image-grid">
    {imageArray.map((image, index) => (
      <div key={index} className="image-container">
        <h3>{image ? image.label : `Image ${index + 1}`}</h3>
        {image ? (
          <img src={image.src} alt={image.label} className="grid-image" />
        ) : (
          <div className="placeholder">Placeholder</div>
        )}
      </div>
    ))}
  </div>
  );
};

export default ImageGrid;