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
  recipe: Recipe;
}

export interface Recipe {
  name: string;
  fermentation: Fermentation;
  og: number;
  fg: number;
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
