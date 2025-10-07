import type { components } from '@/types/generated/api';

export type CarMake =
  components['schemas']['models.CarMake'];
export type CarModel =
  components['schemas']['models.CarModel'];
export type CarGeneration =
  components['schemas']['models.CarGeneration'];

export interface CarMakeWithModels extends CarMake {
  models?: CarModel[];
}

export interface CarModelWithGenerations extends CarModel {
  generations?: CarGeneration[];
}

export interface CarSelection {
  make?: CarMake;
  model?: CarModel;
  generation?: CarGeneration;
}

export interface CarSelectorProps {
  value?: CarSelection;
  onChange: (selection: CarSelection) => void;
  required?: boolean;
  disabled?: boolean;
  className?: string;
  showGenerations?: boolean;
  placeholder?: {
    make?: string;
    model?: string;
    generation?: string;
  };
}
