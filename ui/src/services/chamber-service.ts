import axios from "axios";
import { Chamber } from "../types/Chamber";
import { authHeader, getUrl } from "./common";

class ChamberService {
  getAll() {
    return axios.get<Array<Chamber>>(getUrl("chambers"), {
      headers: authHeader(),
    });
  }
  get(id: string) {
    return axios.get<Chamber>(getUrl(`chambers/${id}`), {
      headers: authHeader(),
    });
  }
  save(chamber: Chamber) {
    return axios.post<Chamber>(getUrl("chambers"), chamber, {
      headers: authHeader(),
    });
  }
  delete(id: string) {
    return axios.delete<string>(getUrl(`chambers/${id}`), {
      headers: authHeader(),
    });
  }
  startFermentation(id: string, step: string) {
    return axios.post<string>(
      getUrl(`chambers/${id}/start?step=${step}`),
      null,
      {
        headers: authHeader(),
      }
    );
  }
  stopFermentation(id: string) {
    return axios.post<string>(getUrl(`chambers/${id}/stop`), null, {
      headers: authHeader(),
    });
  }
}

export default new ChamberService();
