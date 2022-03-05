# zymurgauge

Homebrewing automation system

[![Build Status](https://github.com/benjaminbartels/zymurgauge/workflows/Build/badge.svg)](https://github.com/benjaminbartels/zymurgauge/actions?query=workflow%3ABuild)
[![codecov](https://codecov.io/gh/benjaminbartels/zymurgauge/branch/master/graph/badge.svg)](https://codecov.io/gh/benjaminbartels/zymurgauge)
[![Go Report Card](https://goreportcard.com/badge/github.com/benjaminbartels/zymurgauge)](https://goreportcard.com/report/github.com/benjaminbartels/zymurgauge)
[![TODOs](https://badgen.net/https/api.tickgit.com/badgen/github.com/benjaminbartels/zymurgauge)](https://www.tickgit.com/browse?repo=github.com/benjaminbartels/zymurgauge)

## Setup

Run the following

```sh
./scripts/setup.sh
```

## To run infrastucture

Replace `<YOUR_USERNAME>` and `<YOUR_PASSWORD>` with whatever admin credientials you want to use and run the following:

```sh
INFLUXDB_USER=<YOUR_USERNAME> INFLUXDB_PASSWORD=<YOUR_PASSWORD> GROUP_ID=$(stat -c '%g' /var/run/docker.sock) docker-compose -f deployments/docker-compose.yml -p zymurgauge up
```

## To run zymurgauge app

```sh
sudo zym
```
