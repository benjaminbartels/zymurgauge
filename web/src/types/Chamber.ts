import { BatchDetail } from "./Batch";

export interface Chamber {
  id: string;
  name: string;
  deviceConfigs: DeviceConfig[];
  chillerKp: number;
  chillerKi: number;
  chillerKd: number;
  heaterKp: number;
  heaterKi: number;
  heaterKd: number;
  currentBatch: BatchDetail;
  currentFermentationStep: string;
  readings: Readings;
}

export interface DeviceConfig {
  id: string;
  roles: string[];
  type: string;
}

export interface Readings {
  beerTemperature: number;
  auxiliaryTemperature: number;
  externalTemperature: number;
  hydrometerGravity: number;
}
