import { createSignal, createEffect, onCleanup } from 'solid-js';
import type { Component } from 'solid-js';

interface ImageDropProps {
  onFilesDrop: (files: string[]) => void;
}

const ImageDrop: Component<ImageDropProps> = (props) => {
  const [isDragging, setIsDragging] = createSignal(false);
  let fileInputRef: HTMLInputElement | undefined;

  const handleDragOver = (e: DragEvent) => {
    e.preventDefault();
    e.stopPropagation();
  };

  const handleDragEnter = (e: DragEvent) => {
    e.preventDefault();
    e.stopPropagation();
    setIsDragging(true);
  };

  const handleDragLeave = (e: DragEvent) => {
    e.preventDefault();
    e.stopPropagation();
    setIsDragging(false);
  };

  const convertToBase64 = (file: File): Promise<string> => {
    return new Promise((resolve, reject) => {
      const reader = new FileReader();
      reader.readAsDataURL(file);
      reader.onload = () => resolve(reader.result as string);
      reader.onerror = (error) => reject(error);
    });
  };

  const processFiles = (files: FileList | null) => {
    if (files && files.length) {
      Promise.all(Array.from(files).map(convertToBase64))
        .then(base64Files => {
          props.onFilesDrop(base64Files);
        })
        .catch(error => {
          console.error('Error converting files to base64:', error);
        });
    } else {
      console.warn('No files selected');
    }
  };

  const handleDrop = (e: DragEvent) => {
    e.preventDefault();
    e.stopPropagation();
    setIsDragging(false);
    
    let files: File[] = [];

    if (e.dataTransfer?.items) {
      files = Array.from(e.dataTransfer.items)
        .filter(item => item.kind === 'file')
        .map(item => item.getAsFile())
        .filter((file): file is File => file !== null);
    } else if (e.dataTransfer?.files) {
      files = Array.from(e.dataTransfer.files);
    }

    processFiles(files.length > 0 ? (files as unknown as FileList) : null);
  };

  const handleClick = () => {
    fileInputRef?.click();
  };

  const handleFileInputChange = (e: Event) => {
    const target = e.target as HTMLInputElement;
    processFiles(target.files);
  };

  createEffect(() => {
    const div = document.getElementById('image-drop-zone');
    if (div) {
      div.addEventListener('dragover', handleDragOver);
      div.addEventListener('dragenter', handleDragEnter);
      div.addEventListener('dragleave', handleDragLeave);
      div.addEventListener('drop', handleDrop);

      onCleanup(() => {
        div.removeEventListener('dragover', handleDragOver);
        div.removeEventListener('dragenter', handleDragEnter);
        div.removeEventListener('dragleave', handleDragLeave);
        div.removeEventListener('drop', handleDrop);
      });
    }
  });

  return (
    <div
      id="image-drop-zone"
      onClick={handleClick}
      style={{
        border: `2px dashed ${isDragging() ? 'blue' : 'gray'}`,
        'border-radius': '4px',
        padding: '20px',
        'text-align': 'center',
        cursor: 'pointer',
      }}
    >
      {isDragging() ? 'Drop files here' : 'Drag and drop files here or click to select'}
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

export default ImageDrop;