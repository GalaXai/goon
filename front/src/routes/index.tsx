import { Title } from "@solidjs/meta";
import { createSignal, Show, For } from "solid-js";

import Slider from "../components/Slider/Slider";
import ImageGrid from "~/components/ImageGrid/ImageGrid";
import ImageDrop from "~/components/ImageDrop/ImageDrop";

export default function Home() {
  const [base64Images, setBase64Images] = createSignal<string[]>([]);
  const [sliderValue, setSliderValue] = createSignal(50);
  const [refreshKey, setRefreshKey] = createSignal(0);

  const handleFilesDrop = (base64Files: string[]) => {
    setBase64Images(base64Files);
    setRefreshKey(
      prev => prev + 1);  // Increment refresh key

  };

  return (
    <main>
      <Title>Image Filters</Title>
      <h1>Image Filters!</h1>
      <Slider value={sliderValue} onChange={setSliderValue} />
      <ImageDrop onFilesDrop={handleFilesDrop} />
      <Show when={base64Images().length > 0}>
        <For each={base64Images()}>
          {(base64, index) => (
            <div>
              <h2>Image Set {index() + 1}</h2>
              <ImageGrid 
                mainImage={base64} 
                gradientThreshold={sliderValue()} 
                key={refreshKey()}
              />
            </div>
          )}
        </For>
      </Show>
    </main>
  );
}