import React, { useState, useCallback, useRef } from 'react';

const DragAndDrop = ({ onFilesDrop }) => {
  const [isDragging, setIsDragging] = useState(false);
  const fileInputRef = useRef(null);

  const handleDragOver = useCallback((e) => {
    e.preventDefault();
    e.stopPropagation();
  }, []);

  const handleDragEnter = useCallback((e) => {
    e.preventDefault();
    e.stopPropagation();
    setIsDragging(true);
  }, []);

  const handleDragLeave = useCallback((e) => {
    e.preventDefault();
    e.stopPropagation();
    setIsDragging(false);
  }, []);

  const convertToBase64 = (file) => {
    return new Promise((resolve, reject) => {
      const reader = new FileReader();
      reader.readAsDataURL(file);
      reader.onload = () => resolve(reader.result);
      reader.onerror = (error) => reject(error);
    });
  };

  const processFiles = (files) => {
    if (files.length) {
      Promise.all([...files].map(convertToBase64))
        .then(base64Files => {
          onFilesDrop(base64Files);
        })
        .catch(error => {
          console.error('Error converting files to base64:', error);
        });
    } else {
      console.warn('No files selected');
    }
  };

  const handleDrop = useCallback((e) => {
    e.preventDefault();
    e.stopPropagation();
    setIsDragging(false);
    
    let files = [];

    if (e.dataTransfer.items) {
      files = [...e.dataTransfer.items]
        .filter(item => item.kind === 'file')
        .map(item => item.getAsFile());
    } else if (e.dataTransfer.files) {
      files = [...e.dataTransfer.files];
    }

    processFiles(files);
  }, [onFilesDrop]);

  const handleClick = () => {
    fileInputRef.current.click();
  };

  const handleFileInputChange = (e) => {
    processFiles(e.target.files);
  };

  return (
    <div
      onClick={handleClick}
      onDragOver={handleDragOver}
      onDragEnter={handleDragEnter}
      onDragLeave={handleDragLeave}
      onDrop={handleDrop}
      style={{
        border: `2px dashed ${isDragging ? 'blue' : 'gray'}`,
        borderRadius: '4px',
        padding: '20px',
        textAlign: 'center',
        cursor: 'pointer',
      }}
    >
      {isDragging ? 'Drop files here' : 'Drag and drop files here or click to select'}
      <input
        type="file"
        ref={fileInputRef}
        onChange={handleFileInputChange}
        style={{ display: 'none' }}
        multiple
      />
    </div>
  );
};

export default DragAndDrop;