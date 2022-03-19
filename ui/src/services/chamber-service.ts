import { Chamber } from "../types/Chamber";
import { HTTP } from "./http-common";

class ChamberService {
  getAll() {
    return HTTP.get<Array<Chamber>>("/chambers");
  }
  get(id: string) {
    return HTTP.get<Chamber>(`/chambers/${id}`);
  }
  save(chamber: Chamber) {
    return HTTP.post<Chamber>("/chambers", chamber);
  }
  delete(id: string) {
    return HTTP.delete<string>(`/chambers/${id}`);
  }
  startFermentation(id: string, step: string) {
    return HTTP.post<string>(`/chambers/${id}/start?step=${step}`);
  }
  stopFermentation(id: string) {
    return HTTP.post<string>(`/chambers/${id}/stop`);
  }
}

export default new ChamberService();
