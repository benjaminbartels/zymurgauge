import axios from "axios";
// require('promise.prototype.finally').shim()

export const HTTP = axios.create({
  baseURL: process.env.REACT_APP_ZYMURGAUGE_API_URL,
});

HTTP.interceptors.request.use((request) => {
  // request.headers.common['Authorization'] = 'Bearer ' + localStorage.getItem('access_token') // TODO: is this needed?
  return request;
});
