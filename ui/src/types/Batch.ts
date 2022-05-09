export interface BatchSummary {
  id: string;
  number: number;
  recipeName: string;
}

export interface BatchDetail {
  id: string;
  number: number;
  recipe: Recipe;
}

export interface Recipe {
  name: string;
  fermentation: Fermentation;
  originalGravity: number;
  finalGravity: number;
}

export interface Fermentation {
  name: string;
  steps: FermentationStep[];
}

export interface FermentationStep {
  name: string;
  temperature: number;
  duration: number;
}
