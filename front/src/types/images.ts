export interface Base64ImagesResponse {
  originalImage: string;
  desaturatedImage: string;
  downsampledImage: string;
  gaussiansDiffImage: string;
  horizontalSobel: string;
  verticalSobel: string;
}

export interface RuneMatrix3D {
  data: string[][][];
  cols: number;
  rows: number;
  depth: number;
}