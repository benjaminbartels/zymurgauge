import axios from 'axios'
require('promise.prototype.finally').shim()

export const HTTP = axios.create({
  baseURL: `http://orangesword.duckdns.org:3000/api/v1/`
})
