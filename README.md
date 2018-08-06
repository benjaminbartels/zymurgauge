# zymurgauge

Homebrewing automation system

[![Build Status](https://travis-ci.org/benjaminbartels/zymurgauge.svg?branch=master)](https://travis-ci.org/benjaminbartels/zymurgauge)

[![Code Climate](https://codeclimate.com/github/benjaminbartels/zymurgauge/badges/gpa.svg)](https://codeclimate.com/github/benjaminbartels/zymurgauge)

[![Test Coverage](https://codeclimate.com/github/benjaminbartels/zymurgauge/badges/coverage.svg)](https://codeclimate.com/github/benjaminbartels/zymurgauge/coverage)

[![Issue Count](https://codeclimate.com/github/benjaminbartels/zymurgauge/badges/issue_count.svg)](https://codeclimate.com/github/benjaminbartels/zymurgauge)

[![Go Report Card](https://goreportcard.com/badge/github.com/benjaminbartels/zymurgauge)](https://goreportcard.com/report/github.com/benjaminbartels/zymurgauge)

## Project Structure

* cmd
  * fermmond - *connects to http server*
    * http - client (counter part to handler in server)
  * zymsrvd - *invokes http servers*
    * handlers
* docker
* internal
  * *model hierarchy*
  * middleware
    * logger
    * metrics
    * auth
  * platform - kit
    * log - logger interface
    * bolt - boltdb interaction layer
    * web
      * app.go - *contains App struct w/ wrapped handle func*
      * middleware.go - *wrapMiddleware func*
      * response.go - *response func, errors, etc*