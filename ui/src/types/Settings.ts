export enum TemperatureUnits {
  Celsius = 1,
  Fahrenheit,
}

export interface AppSettings {
  temperatureUnits: TemperatureUnits;
  authSecret: string;
  brewfatherApiUserId: string;
  brewfatherApiKey: string;
  brewfatherLogUrl: string;
  influxDbUrl: string;
  influxDbReadToken: string;
  statsDAddress: string;
}
