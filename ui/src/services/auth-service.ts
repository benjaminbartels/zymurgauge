import axios from "axios";
import { AuthRequest, AuthResponse } from "../types/Auth";
import { getUrl } from "./common";

class AuthService {
  async login(username: string, password: string) {
    const auth: AuthRequest = {
      username: username,
      password: password,
    };

    try {
      const response = await axios.post<AuthResponse>(getUrl("login"), auth);
      console.debug("Login success: ", response.data.token);
      localStorage.setItem("token", response.data.token);
    } catch (e: any) {
      console.error("Login Error:", e);
      throw e;
    }
  }

  logout() {
    localStorage.removeItem("token");
    localStorage.removeItem("influxDbUrl");
  }

  getToken() {
    return localStorage.getItem("token");
  }
}

export default new AuthService();
