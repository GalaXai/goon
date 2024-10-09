import React, { useState, useEffect } from 'react';
import './ImageGrid.css';
import { loadImage } from '../services/imageApi';
import AsciiArtCanvas from './AsciiArtCanvas';

const ImageGrid = ({ mainImage, gradientThreshold }) => {
  const [imageArray, setImageArray] = useState(Array(6).fill(null));
  const [asciiArt, setAsciiArt] = useState(null);

  useEffect(() => {
    if (mainImage) {
      loadImage(mainImage, gradientThreshold)
        .then(results => {
          console.log('Images loaded successfully');
          const { imageResponse, asciiArt } = results;
          setAsciiArt(asciiArt);
          console.log(asciiArt)
          setImageArray([
            { src: `data:image/png;base64,${imageResponse.originalImage}`, label: 'Original' },
            { src: `data:image/png;base64,${imageResponse.desaturatedImage}`, label: 'Desaturated' },
            { src: `data:image/png;base64,${imageResponse.downsampledImage}`, label: 'Downsampled' },
            { src: `data:image/png;base64,${imageResponse.gaussiansDiffImage}`, label: 'Gaussians Diff' },
            { src: `data:image/png;base64,${imageResponse.horizontalSobel}`, label: 'Horizontal Sobel' },
            { src: `data:image/png;base64,${imageResponse.verticalSobel}`, label: 'Vertical Sobel' },
          ]);
        })
        .catch(error => {
          console.error('Failed to load image:', error);
          // Handle the error appropriately
        });
    }
  }, [mainImage, gradientThreshold]);

  return (
    <div className="image-grid-container">
      <div className="image-grid">
        {imageArray.map((item, index) => (
          <div key={index} className="image-container">
            <h3>{item ? item.label : `Image ${index + 1}`}</h3>
            {item ? (
              <img src={item.src} alt={item.label} className="grid-image" />
            ) : (
              <div className="placeholder">Placeholder</div>
            )}
          </div>
        ))}
      </div>
      {asciiArt && (
        <div className="ascii-art-container">
          <h3>ASCII Art</h3>
          <AsciiArtCanvas asciiArt={asciiArt} mainImage={mainImage}/>
        </div>
      )}
    </div>
  );
};

export default ImageGrid;