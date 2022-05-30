import { AxiosRequestHeaders } from "axios";

export function getUrl(path: string) {
  // If the app is in development mode (npm start), proxy setting in package.json will be used.
  if (process.env.NODE_ENV === "development") {
    return "/api/v1/" + path;
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
