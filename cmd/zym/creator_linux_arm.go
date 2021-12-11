package main

func createThermometerRepo() mocks.ThermometerRepo {
	return raspberrypi.NewDs18b20Repo()
}
