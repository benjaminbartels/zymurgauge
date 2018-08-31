'use strict'
const merge = require('webpack-merge')
const prodEnv = require('./prod.env')

module.exports = merge(prodEnv, {
  NODE_ENV: '"development"',
  API_URL: '"http://localhost:3000/api/v1/"',
  AUTH_CLIENT_ID: '"h91yYYiVM3x0oLpBL7K4dE30G3Jo2VeP"',
  AUTH_CALLBACK_URL: '"http://localhost:8080/auth"'
})
