import { Accessor, Setter } from "solid-js";

interface SliderProps {
  value: Accessor<number>;
  onChange: Setter<number>;
}

export default function Slider({ value, onChange }: SliderProps) {
  const handleSliderChange = (event: Event) => {
    const target = event.target as HTMLInputElement;
    onChange(parseInt(target.value, 10));
  };

  return (
    <div>
      <input
        type="range"
        min="0"
        max="255"
        step="5"
        value={value()}
        onInput={handleSliderChange}
      />
      <p>Value: {value()}</p>
    </div>
  );
}