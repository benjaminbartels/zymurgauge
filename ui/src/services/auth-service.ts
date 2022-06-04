import axios from "axios";
import { Credentials, LoginResponse } from "../types/Auth";
import { authHeader, getUrl } from "./common";

class AuthService {
  async login(username: string, password: string) {
    const credentials: Credentials = {
      username: username,
      password: password,
    };

    try {
      const response = await axios.post<LoginResponse>(
        getUrl("auth/login"),
        credentials
      );
      console.debug("Login success: ", response.data.token);
      localStorage.setItem("token", response.data.token);
      localStorage.setItem("username", username);
    } catch (e: any) {
      console.error("Login Error:", e);
      throw e;
    }
  }

  logout() {
    localStorage.removeItem("token");
  }

  getToken() {
    return localStorage.getItem("token");
  }

  save(credentials: Credentials) {
    return axios.post<Credentials>(getUrl("auth/update"), credentials, {
      headers: authHeader(),
    });
  }
}

export default new AuthService();
