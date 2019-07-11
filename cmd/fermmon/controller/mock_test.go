package controller_test

import (
	"github.com/benjaminbartels/zymurgauge/internal"
)

type chamberResourceMock struct {
	GetFn         func(string) (*internal.Chamber, error)
	SaveFn        func(*internal.Chamber) error
	SubscribeFn   func(string, chan internal.Chamber) error
	UnsubscribeFn func(string)

	GetInvoked         bool
	SaveInvoked        bool
	SubscribeInvoked   bool
	UnsubscribeInvoked bool
}

func (m *chamberResourceMock) Get(mac string) (*internal.Chamber, error) {
	m.GetInvoked = true
	return m.GetFn(mac)
}

func (m *chamberResourceMock) Save(c *internal.Chamber) error {
	m.SaveInvoked = true
	return m.SaveFn(c)
}

func (m *chamberResourceMock) Subscribe(mac string, ch chan internal.Chamber) error {
	m.SubscribeInvoked = true
	return m.SubscribeFn(mac, ch)
}

func (m *chamberResourceMock) Unsubscribe(mac string) {
	m.UnsubscribeInvoked = true
	m.UnsubscribeFn(mac)
}

type fermentationResourceMock struct {
	GetFn                   func(string) (*internal.Fermentation, error)
	SaveFn                  func(*internal.Fermentation) error
	SaveTemperatureChangeFn func(*internal.TemperatureChange) error

	GetInvoked                   bool
	SaveInvoked                  bool
	SaveTemperatureChangeInvoked bool
}

func (m *fermentationResourceMock) Get(id string) (*internal.Fermentation, error) {
	m.GetInvoked = true
	return m.GetFn(id)
}

func (m *fermentationResourceMock) Save(c *internal.Fermentation) error {
	m.SaveInvoked = true
	return m.SaveFn(c)
}

func (m *fermentationResourceMock) SaveTemperatureChange(t *internal.TemperatureChange) error {
	m.SaveTemperatureChangeInvoked = true
	return m.SaveTemperatureChangeFn(t)
}
