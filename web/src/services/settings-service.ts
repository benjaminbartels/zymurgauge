import { Settings } from "../types/Settings";
import { HTTP } from "./http-common";

class SettingsService {
  get() {
    return HTTP.get<Settings>("/settings");
  }
  save(settings: Settings) {
    return HTTP.post<Settings>("/settings", settings);
  }
}

export default new SettingsService();
