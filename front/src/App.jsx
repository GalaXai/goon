import React, { useState } from 'react';
import DragAndDrop from './components/DragAndDrop';
import ImageGrid from './components/ImageGrid';
import './App.css';

const App = () => {
  const [base64Images, setBase64Images] = useState([]);
  const [sliderValue, setSliderValue] = useState(0);

  const handleFilesDrop = (base64Files) => {
    setBase64Images(base64Files);
  };

  const handleSliderChange = (event) => {
    setSliderValue(Number(event.target.value));
  };

  return (
    <div>
      <h1>Image to Base64 Converter</h1>
      <DragAndDrop onFilesDrop={handleFilesDrop} />
      <div>
        <label htmlFor="threshold-slider">Gradient Threshold: {sliderValue}</label>
        <input
          type="range"
          id="threshold-slider"
          min="0"
          max="255"
          step="5"
          value={sliderValue}
          onChange={handleSliderChange}
        />
      </div>
      {base64Images.map((base64, index) => (
        <div key={index}> 
          <h2>Image Set {index + 1}</h2>
          <ImageGrid mainImage={base64} gradientThreshold={sliderValue}/>
        </div>
      ))}
    </div>
  );
};

export default App;