# zymurgauge

[![Build Status](https://github.com/benjaminbartels/zymurgauge/workflows/Build/badge.svg)](https://github.com/benjaminbartels/zymurgauge/actions?query=workflow%3ABuild)
[![codecov](https://codecov.io/gh/benjaminbartels/zymurgauge/branch/master/graph/badge.svg)](https://codecov.io/gh/benjaminbartels/zymurgauge)
[![Go Report Card](https://goreportcard.com/badge/github.com/benjaminbartels/zymurgauge)](https://goreportcard.com/report/github.com/benjaminbartels/zymurgauge)
[![TODOs](https://badgen.net/https/api.tickgit.com/badgen/github.com/benjaminbartels/zymurgauge)](https://www.tickgit.com/browse?repo=github.com/benjaminbartels/zymurgauge)

## About the project

Zymurgauge is a homebrewing automation system that interfaces with [Brewfather](https://brewfather.app/) for controlling
a fermentation chamber.  The service is designed to run in Docker container on a Raspberry Pi.  The backend is written
in Go while the frontend is in Typescript with the React framework.  

The system currently supports controlling chilling and heating devices, such as a chest freezer and heading lamp, via
the Raspberry Pi's GPIO interface.  It also uses DS18B20 sensors and/or [Tilt Hydrometers](https://tilthydrometer.com/)
to monitor temperatures.  Gravity readings from the Tilts can also be monitored.  

Zymurgauge interfaces the Brewfather to allow the user to select the batch to be fermeneted.  During fermentation,
data is collected and sent to your Brewfather account's streaming endpoint.  A premium membership to Brewfather is
required to use their API.  Data is collected by a Telegraf instance which sends data to an InfluxDB instance via the
StatsD protocol.  A graph of the temperature and gravity readings can be view in the Zymurgauge UI.  The project
contains a Docker compose file which starts Zymurgauge along side InfluxDB and Telegraf all behind an Nginx reverse
proxy.

In the future it will be extended to also control a HERMS (Heat Exchange Re-circulating Mash System).

## Getting Started

### Prerequisites

> **Note**
> Before setting things up, you will need a Premium Brewfather account.  You will need a API UserID and Key to access
> their API.  You will also need to turn on the "Custom Stream" in the "Power-ups" section in your Brewfather settings.
> This will create a unique url used to log data to Brewfather.

Update and upgrade to the latest packages:

```sh
sudo apt-get update && sudo apt-get upgrade
```

Install the GPIO package is it is not already installed:

```sh
sudo apt-get install rpi.gpio
```

Enable the One-Wire interface:

```sh
sudo nano /boot/config.txt
```

then add this to the bottom of the file and save:

```sh
dtoverlay=w1-gpio
```

> **Note**
> If you plan on using a Tilt Hydrometer you will also need to insure bluetooh services are installed on you Raspberry
> Pi.

Install docker on you Raspberry Pi if needed :

```sh
curl -sSL https://get.docker.com | sh
```

Ensure that the current user is in the docker group:

```sh
sudo usermod -aG docker ${USER}
```

Ensure that docker starts on boot:

```sh
sudo systemctl enable docker
```

### Installation

Download and run the setup script:

```sh
wget https://raw.githubusercontent.com/benjaminbartels/zymurgauge/master/scripts/setup.sh
chmod +x setup.sh
./setup.sh
```

The setup script will prompt you for the following:

- Zymurgauge Initial Admin Username
- Zymurgauge Initial Admin Password
- Brewfather API User ID
- Brewfather API Key
- Brewfather Log Stream URL
- InfluxDB Initial Admin Username
- InfluxDB Initial Admin Password

The script will initialize the local storage and conf files for zymurgauge, influxdb, nginx and telegraf in the
`~/.zymurgauge` directory

Download the Docker compose file to your home directory:

```sh
wget https://raw.githubusercontent.com/benjaminbartels/zymurgauge/master/deployments/docker-compose.yml
```

Run this following to start the services:

```sh
GROUP_ID=$(stat -c '%g' /var/run/docker.sock) docker compose -p zymurgauge up -d
```

Once the services are up go to `https://<your-raspberry-pis-hostname>:8080` your web browser:

## Project Layout

api - OpenAPI/Swagger specs
build - Packaging and Continuous Integration (Dockerfiles)

```sh
/
├─ api - OpenAPI/Swagger specs
├─ build - Packaging and Continuous Integration (Dockerfiles)
├─ cmd - Main GO applications for this project
│  ├─ zym - zymurgauge app
│  │  └─ handlers - HTTP request router and handlers
│  └─ zymsim - simlation apps
├─ config - config files
├─ deployments - container deployment configurations (docker-compose files)
├─ internal - private application and library code
│  ├─ auth - jwt authorization models and functions
│  ├─ batch - Batch models and functions
│  ├─ brewfather - Brewfather HTTPS client and models
│  ├─ chamber - fermentation chamber manager and models
│  ├─ database - bbolt database client and database models
│  ├─ device - device realted logic
│  │  ├─ gpio - GPIO Actuator device logic
│  │  ├─ onewire - One-Wire (ds18b20) device logic
│  │  └─ tilt - Tilt Hydrometer monitor and logic
│  ├─ middleware - HTTP request router middlewares
│  ├─ platform - foundational packages
│  │  ├─ bluetooth - bluetooth discoverer and ibeacon logic
│  │  ├─ clock - wrapper for Go time package
│  │  ├─ debug - pprof mux
│  │  ├─ metrics - metrics interface for statsd
│  │  └─ web - web server, api and 
│  ├─ settings - settings models
│  ├─ temperaturecontrol - temperature controller implementations
│  │  ├─ hysteresis - hysteresis temperature controller implementation
│  │  └─ pid - pid temperature controller implementation
│  └─ test - fakes, mocks and stubs for testing
├─ scripts - setup scripts
└─ ui - react ui source and ui file embed filesystem
   ├─ public - public static assets
   ├─ src - source code for web application
   │  ├─ components - web application views
   │  ├─ services - API clients
   │  └─ types - typescript type interfaces
```

## Roadmap

- [ ] Add HERMS management functionality
- [ ] Improve and polish React UI
- [ ] Add ISpendle support

## Development Setup

> **Note**
> Ensure that you have at least Go 1.20 and yarn 1.22 installed.

To initialize your local development environment run:

```sh
make init
```

To run the Go service locally with live reloading run:

```sh
make watch-go
```

To run the React UI locally with live reloading run:

```sh
make watch-react
```
