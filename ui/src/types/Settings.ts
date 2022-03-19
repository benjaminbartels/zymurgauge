export enum TemperatureUnits {
  Celsius = 1,
  Fahrenheit,
}

export interface Settings {
  brewfatherApiUserId: string;
  brewfatherApiKey: string;
  brewfatherLogUrl: string;
  temperatureUnits: TemperatureUnits;
}
