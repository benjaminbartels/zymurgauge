# zymurgauge

Homebrewing automation system

[![Build Status](https://travis-ci.org/benjaminbartels/zymurgauge.svg?branch=master)](https://travis-ci.org/benjaminbartels/zymurgauge)

[![Code Climate](https://codeclimate.com/github/benjaminbartels/zymurgauge/badges/gpa.svg)](https://codeclimate.com/github/benjaminbartels/zymurgauge)

[![Test Coverage](https://codeclimate.com/github/benjaminbartels/zymurgauge/badges/coverage.svg)](https://codeclimate.com/github/benjaminbartels/zymurgauge/coverage)

[![Issue Count](https://codeclimate.com/github/benjaminbartels/zymurgauge/badges/issue_count.svg)](https://codeclimate.com/github/benjaminbartels/zymurgauge)

[![Go Report Card](https://goreportcard.com/badge/github.com/benjaminbartels/zymurgauge)](https://goreportcard.com/report/github.com/benjaminbartels/zymurgauge)

## Project Hierarchy

- zymurgauge - root
  - internal - Hides packages that should only imported from parent directory’s packages
    - handlers - HTTP handlers
    - lambdas - AWS lambdas
    - *model files* - Model structs used throughout project
    - *interfaces.go* - Interfaces for common hardware interaction functions
    - database - Contains database interfaces and implementations
      - dynamodb - AWS DynamoDB implementation
      - boltdb - BoltDB implementation
      - *interfaces* - interfaces for common database functions
    - platform - Platform specific packages that could be promoted to own repos
    - simulation - Simulators
    - client - HTTP client
  - cmd - Executable
    - fermmon - Fermentation Monitor App
    - zymsrv - App that Host the API and Web UI

## ToDo

- Add optional config Funcs to constructors

- Review error handling and separation `https://www.ardanlabs.com/blog/2017/02/package-oriented-design.html`
