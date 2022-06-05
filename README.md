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

Zymurgauge interfaces the Brewfather to allow the user to select the the batch to be fermeneted.  During fermentation,
data is collected and sent to your Brewfather account's streaming endpoint.  A premium membership to Brewfather is
required to use their API.  Data is collected by an Telegraf instance which send it to an InfluxDB instance via the
StatsD protocol.  A graph of the temperature and gravity readings can be view in the Zymurgauge UI.  The project
constains a Docker compose file which starts Zymurgauge along side InfluxDB and Telegraf all behind a nginx reverse
proxy.

In the future it will be extended to also control a HERMS (Heat Exchange Re-circulating Mash System).

## Getting Started

### Prerequisites

> **Note:**
> Before setting things up, you will need a Premium Brewfather account.  You will need a API UserID and Key to access
> their API.  You will also need to turn on the "Custom Stream" in the "Power-ups" section in your Brewfather settings.
> This will create a unique url used to log data to Brewfather.

Update and upgrade to the latest packages:

```sh
sudo apt-get update
sudo apt-get upgrade
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

> **Note:**
> If you plan on using a Tilt Hydrometer you will also need to sure bluetooh services are installed on you Raspberry Pi.

Install docker on you Raspberry Pi if needed and ensure the current user is in the docker group:

```sh
curl -sSL https://get.docker.com | sh
sudo usermod -aG docker ${USER}
```

Ensure that docker starts on boot:

```sh
sudo systemctl enable docker
```

### Installation

Run the setup script:

```sh
curl -sSL https://raw.githubusercontent.com/benjaminbartels/zymurgauge/master/scripts/setup.sh | sh
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
wget https://github.com/benjaminbartels/zymurgauge/blob/master/deployments/docker-compose.yml
```

Run this following to start the services:

```sh
GROUP_ID=$(stat -c '%g' /var/run/docker.sock) docker compose -d -f docker-compose.yml -p zymurgauge up
```

Once the services are up go to `https://<your-raspberry-pis-hostname>:8080` your web browser:

## Project Layout

## Roadmap

## Development Setup
