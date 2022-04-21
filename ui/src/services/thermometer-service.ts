import axios from "axios";
import { authHeader, getUrl } from "./common";
class ThermometerService {
  getAll() {
    return axios.get<Array<string>>(getUrl("thermometers"), {
      headers: authHeader(),
    });
  }
}

export default new ThermometerService();
