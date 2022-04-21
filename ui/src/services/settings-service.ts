import axios from "axios";
import { AppSettings } from "../types/Settings";
import { authHeader, getUrl } from "./common";

class SettingsService {
  get() {
    return axios.get<AppSettings>(getUrl("settings"), {
      headers: authHeader(),
    });
  }
  save(settings: AppSettings) {
    return axios.post<AppSettings>(getUrl("settings"), settings, {
      headers: authHeader(),
    });
  }
}

export default new SettingsService();
