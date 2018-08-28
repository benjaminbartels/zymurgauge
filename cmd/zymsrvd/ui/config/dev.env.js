'use strict'
const merge = require('webpack-merge')
const prodEnv = require('./prod.env')

module.exports = merge(prodEnv, {
  NODE_ENV: '"development"',
  API_URL: '"http://localhost:3000/api/v1/"',
  AUTH_CALLBACK_URL: '"http://localhost:8080/auth"'
})
