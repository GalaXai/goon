import { createSignal, For, createResource, createMemo, createEffect, JSX } from "solid-js";
import { loadImage } from "~/services/imageApi";
import { RuneMatrix3D } from '../../types/images';
import AsciiArtCanvas from '../../AsciiArtCanvas/AsciiArtCanvas';
import './ImageGrid.css';

interface ImageGridProps {
    mainImage: string;
    gradientThreshold: number;
    threshold: number;
    tau: number;
    key?: number;
  }

export default function ImageGrid(props: ImageGridProps & JSX.HTMLAttributes<HTMLDivElement>) {
    const [imageArray, setImageArray] = createSignal(Array(6).fill(null));
    const [asciiArt, setAsciiArt] = createSignal<RuneMatrix3D | null>(null);
    
    const imageParams = createMemo(() => ({ 
        base64Image: props.mainImage, 
        gradientThreshold: props.gradientThreshold,
        threshold: props.threshold,
        tau: props.tau,
    }));

    const [imageData] = createResource(imageParams, async (params) => {
        if (!params.base64Image) return null;
        try {
            console.log('Gradient Threshold:', props.gradientThreshold);
            console.log('Tau Value:', props.tau);
            console.log('Threshold Value:', props.threshold);
            const results = await loadImage(params);
            const { imageResponse, asciiArt } = results;
            setAsciiArt(asciiArt);
            setImageArray([
                { src: `data:image/png;base64,${imageResponse.originalImage}`, label: 'Original' },
                { src: `data:image/png;base64,${imageResponse.desaturatedImage}`, label: 'Desaturated' },
                { src: `data:image/png;base64,${imageResponse.downsampledImage}`, label: 'Downsampled' },
                { src: `data:image/png;base64,${imageResponse.gaussiansDiffImage}`, label: 'Gaussians Diff' },
                { src: `data:image/png;base64,${imageResponse.horizontalSobel}`, label: 'Horizontal Sobel' },
                { src: `data:image/png;base64,${imageResponse.verticalSobel}`, label: 'Vertical Sobel' },
            ]);
            return results;
        } catch (error) {
            console.error('Failed to load image:', error)
        }
    })

    return (
        <div class="image-grid-container" {...props}>
            <div class="image-grid">
                <For each={imageArray()}>
                    {(item, index) => (
                        <div class="image-container">
                            <h3>{item ? item.label : `Item ${index() + 1}`}</h3>
                            {item ? (
                                <img src={item.src} about={item.label} class="grid-image"/>
                            ) : (
                                <div class="placeholder">Placeholder</div>
                            )}
                        </div>
                    )}
                </For>
            </div>
            {asciiArt() && (
                <div>
                    <h3>ASCII Art</h3>
                    <AsciiArtCanvas asciiArt={asciiArt()} />
                </div>
            )}
        </div>
    )
}