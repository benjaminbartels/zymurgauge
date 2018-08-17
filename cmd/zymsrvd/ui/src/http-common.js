import axios from 'axios'
require('promise.prototype.finally').shim()

export const HTTP = axios.create({
  baseURL: process.env.API_URL
})
