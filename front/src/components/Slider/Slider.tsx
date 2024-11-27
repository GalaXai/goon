import { Accessor, Setter } from "solid-js";

interface SliderProps {
  value: Accessor<number>;
  onChange: Setter<number>;
  min?: number;
  max?: number;
  step?: number;
  label?: string;
}

export default function Slider({ value, onChange, min = 0, max = 255, step = 1, label }: SliderProps) {
  const handleSliderChange = (event: Event) => {
    const target = event.target as HTMLInputElement;
    onChange(parseFloat(target.value));
  };

  return (
    <div class="slider-container">
      {label && <label>{label}</label>}
      <input
        type="range"
        min={min}
        max={max}
        step={step}
        value={value()}
        onInput={handleSliderChange}
      />
      <span class="slider-value">{value().toFixed(2)}</span>
    </div>
  );
}