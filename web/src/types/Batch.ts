export interface BatchSummary {
  id: string;
  name: string;
  number: number;
  recipeName: string;
}

export interface BatchDetail {
  id: string;
  name: string;
  number: number;
  recipeName: string;
  fermentation: Fermentation;
}

export interface Fermentation {
  name: string;
  steps: FermentationStep[];
}

export interface FermentationStep {
  type: string;
  actualTime: number;
  stepTemperature: number;
  stepTime: number;
}
