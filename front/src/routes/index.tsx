import { Title } from "@solidjs/meta";
import { createSignal, Show, For } from "solid-js";

import Slider from "../components/Slider/Slider";
import ImageGrid from "~/components/ImageGrid/ImageGrid";
import ImageDrop from "~/components/ImageDrop/ImageDrop";

export default function Home() {
  const [base64Images, setBase64Images] = createSignal<string[]>([]);
  const [gradientThreshold, setGradientThreshold] = createSignal(50);
  const [tau, setTau] = createSignal(0.98);
  const [threshold, setThreshold] = createSignal(5.0);
  const [refreshKey, setRefreshKey] = createSignal(0);

  const handleFilesDrop = (base64Files: string[]) => {
    setBase64Images(base64Files);
    setRefreshKey(prev => prev + 1);
  };

  return (
    <main>
      <Title>Image Filters</Title>
      <h1>Image Filters!</h1>
      <div class="sliders-container">
        <Slider 
          value={gradientThreshold} 
          onChange={setGradientThreshold} 
          min={0} 
          max={255} 
          step={5} 
          label="Gradient Threshold"
        />
        <Slider 
          value={tau} 
          onChange={setTau} 
          min={0} 
          max={1} 
          step={0.01} 
          label="Tau Value"
        />
        <Slider
          value={threshold} 
          onChange={setThreshold} 
          min={0} 
          max={1000} 
          step={5} 
          label="Threshold Value"
        />
      </div>
      <ImageDrop onFilesDrop={handleFilesDrop} />
      <Show when={base64Images().length > 0}>
        <For each={base64Images()}>
          {(base64, index) => (
            <div>
              <h2>Image Set {index() + 1}</h2>
              <ImageGrid 
                mainImage={base64} 
                gradientThreshold={gradientThreshold()} 
                threshold={threshold()}
                tau={tau()}
                key={refreshKey()}
              />
            </div>
          )}
        </For>
      </Show>
    </main>
  );
}