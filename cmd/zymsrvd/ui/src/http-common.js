import axios from 'axios'
require('promise.prototype.finally').shim()

export const HTTP = axios.create({
  baseURL: process.env.API_URL
})

HTTP.interceptors.request.use(request => {
  request.headers.common['Authorization'] = 'Bearer ' + localStorage.getItem('access_token')
  return request
})
