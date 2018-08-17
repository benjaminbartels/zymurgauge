var merge = require('webpack-merge')
var prodEnv = require('./prod.env')

module.exports = merge(prodEnv, {
  NODE_ENV: '"development"',
  'API_URL': "`http://orangesword.duckdns.org:3000/api/v1/`"
})
