import React, { useRef, useEffect, useState } from 'react';
// Remove the import for colorDownsample

const AsciiArtTextBox = ({ asciiArt }) => {
  const textareaRef = useRef(null);
  const containerRef = useRef(null);
  const [fontSize, setFontSize] = useState(0.4);
  // Remove colorImage state

  useEffect(() => {
    if (asciiArt && textareaRef.current && containerRef.current) {
      const { data, rows } = asciiArt;

    // Convert numeric values to characters
    const convertToChar = (value) => {
      if (value === 0) return ' ';
      return  String.fromCharCode(value);
    };

      // Set textarea content
      const asciiString = data.map((layer, layerIndex) => {
        const layerString = layer.map((row, rowIndex) => {
          const rowString = row.map(convertToChar).join('');
          return rowString;
        }).join(' ');
        return layerString;
      }).join('\n\n');

      textareaRef.current.value = asciiString;
      const containerHeight = `${rows * (fontSize * 1.5)}em`;
      containerRef.current.style.height = containerHeight;
    }
  }, [asciiArt, fontSize]);

  const containerStyle = {
    width: '100%',
    display: 'flex',
    justifyContent: 'center',
    alignItems: 'center',
  };

  const textareaStyle = {
    fontFamily: 'monospace',
    whiteSpace: 'pre',
    overflowWrap: 'normal',
    overflow: 'hidden',
    resize: 'none',
    fontSize: `${fontSize}em`,
    lineHeight: '1',
    padding: '0',
    border: 'none',
    width: '100%',
    height: '100%',
  };

  return (
    <div ref={containerRef} style={containerStyle}>
      <textarea ref={textareaRef} style={textareaStyle} readOnly />
    </div>
  );
};

export default AsciiArtTextBox;