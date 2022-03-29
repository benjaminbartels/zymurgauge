import axios from "axios";
import { Settings } from "../types/Settings";
import { authHeader, getUrl } from "./common";

class SettingsService {
  get() {
    return axios.get<Settings>(getUrl("settings"), {
      headers: authHeader(),
    });
  }
  save(settings: Settings) {
    return axios.post<Settings>(getUrl("settings"), settings, {
      headers: authHeader(),
    });
  }
}

export default new SettingsService();
