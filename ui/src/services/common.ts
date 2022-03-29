import { AxiosRequestHeaders } from "axios";

export function getUrl(path: string) {
  // If the app was bundled with a development API URL (localhost:8080 or some other remote server) use that one.
  if (process.env.REACT_APP_DEVELOPMENT_API_URL) {
    return process.env.REACT_APP_DEVELOPMENT_API_URL + path;
  }

  return window.location.origin + "/api/v1/" + path;
}

export function authHeader(): AxiosRequestHeaders {
  const token = localStorage.getItem("token");

  if (token) {
    return { Authorization: "Bearer " + token };
  } else {
    return {};
  }
}
