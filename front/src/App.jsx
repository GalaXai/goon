import React, { useState } from 'react';
import DragAndDrop from './components/DragAndDrop';
import ImageGrid from './components/ImageGrid';
import './App.css';

const App = () => {
  const [base64Images, setBase64Images] = useState([]);

  const handleFilesDrop = (base64Files) => {
    setBase64Images(base64Files);
  };

  return (
    <div>
      <h1>Image to Base64 Converter</h1>
      <DragAndDrop onFilesDrop={handleFilesDrop} />
      {base64Images.map((base64, index) => (

        <div key={index}>
          <h2>Image Set {index + 1}</h2>
          <ImageGrid mainImage={base64} />
          <textarea
            value={base64}
            readOnly
            style={{ width: '100%', height: '100px', marginTop: '10px' }}
          />
        </div>
      ))}
    </div>
  );
};

export default App;