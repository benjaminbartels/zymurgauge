import { BatchDetail } from "./Batch";

export interface Chamber {
  id: string | undefined;
  name: string;
  deviceConfig: DeviceConfig;
  chillerKp: number;
  chillerKi: number;
  chillerKd: number;
  heaterKp: number;
  heaterKi: number;
  heaterKd: number;
  currentBatch: BatchDetail | undefined;
  currentFermentationStep: string;
  readings: Readings | null;
}

export interface DeviceConfig {
  chillerGpio: string;
  heaterGpio: string;
  beerThermometerType: string;
  beerThermometerId: string;
  auxiliaryThermometerType: string;
  auxiliaryThermometerId: string;
  externalThermometerType: string;
  externalThermometerId: string;
  hydrometerType: string;
  hydrometerId: string;
}

export interface Readings {
  beerTemperature: number;
  auxiliaryTemperature: number;
  externalTemperature: number;
  hydrometerGravity: number;
}
