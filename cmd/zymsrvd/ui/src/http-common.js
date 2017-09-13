import axios from 'axios'
require('promise.prototype.finally').shim()

export const HTTP = axios.create({
    baseURL: `http://localhost:3000/v1/`
})
