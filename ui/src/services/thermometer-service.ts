import { HTTP } from "./http-common";

class ThermometerService {
  getAll() {
    return HTTP.get<Array<string>>("/thermometers");
  }
}

export default new ThermometerService();
