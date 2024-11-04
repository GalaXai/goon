import { createEffect, createSignal } from "solid-js";
import { RuneMatrix3D } from '../types/images';
import styles from './AsciiArtCanvas.module.css';

interface AsciiArtCanvasProps {
  asciiArt: RuneMatrix3D | null;
}

const AsciiArtCanvas = (props: AsciiArtCanvasProps) => {
  let textareaRef: HTMLTextAreaElement | undefined;
  let containerRef: HTMLDivElement | undefined;
  const [fontSize, setFontSize] = createSignal(0.4);

  const convertToChar = (value: number) => {
    if (value === 0) return '';
    return String.fromCharCode(value);
  };

  createEffect(() => {
    if (props.asciiArt && textareaRef && containerRef) {
      const { data, rows } = props.asciiArt;

      // Set textarea content
      const asciiString = data.map((layer) => {
        const layerString = layer.map((row) => {
          return row.map((value) => convertToChar(Number(value))).join('');
        }).join('');
        return layerString;
      }).join('\n\n');

      textareaRef.value = asciiString;
      const containerHeight = `${rows * (fontSize() * 1.5)}em`;
      containerRef.style.height = containerHeight;
    }
  });

  return (
    <div ref={containerRef} class={styles.container}>
      <textarea 
        ref={textareaRef} 
        class={styles.textarea} 
        style={{ 'font-size': `${fontSize()}em` }}
        readOnly 
      />
    </div>
  );
};

export default AsciiArtCanvas;