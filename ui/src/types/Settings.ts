export enum TemperatureUnits {
  Celsius = 1,
  Fahrenheit,
}

export interface Settings {
  temperatureUnits: TemperatureUnits;
  authSecret: string;
  brewfatherApiUserId: string;
  brewfatherApiKey: string;
  brewfatherLogUrl: string;
  influxDbUrl: string;
  statsDAddress: string;
}
