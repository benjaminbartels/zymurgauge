package controller_test

import (
	"github.com/benjaminbartels/zymurgauge/internal"
	"github.com/benjaminbartels/zymurgauge/internal/storage"
)

type chamberResourceMock struct {
	GetFn         func(string) (*storage.Chamber, error)
	SaveFn        func(*storage.Chamber) error
	SubscribeFn   func(string, chan storage.Chamber) error
	UnsubscribeFn func(string)

	GetInvoked         bool
	SaveInvoked        bool
	SubscribeInvoked   bool
	UnsubscribeInvoked bool
}

func (m *chamberResourceMock) Get(mac string) (*storage.Chamber, error) {
	m.GetInvoked = true
	return m.GetFn(mac)
}

func (m *chamberResourceMock) Save(c *storage.Chamber) error {
	m.SaveInvoked = true
	return m.SaveFn(c)
}

func (m *chamberResourceMock) Subscribe(mac string, ch chan storage.Chamber) error {
	m.SubscribeInvoked = true
	return m.SubscribeFn(mac, ch)
}

func (m *chamberResourceMock) Unsubscribe(mac string) {
	m.UnsubscribeInvoked = true
	m.UnsubscribeFn(mac)
}

type fermentationResourceMock struct {
	GetFn                   func(uint64) (*storage.Fermentation, error)
	SaveFn                  func(*storage.Fermentation) error
	SaveTemperatureChangeFn func(*internal.TemperatureChange) error

	GetInvoked                   bool
	SaveInvoked                  bool
	SaveTemperatureChangeInvoked bool
}

func (m *fermentationResourceMock) Get(id uint64) (*storage.Fermentation, error) {
	m.GetInvoked = true
	return m.GetFn(id)
}

func (m *fermentationResourceMock) Save(c *storage.Fermentation) error {
	m.SaveInvoked = true
	return m.SaveFn(c)
}

func (m *fermentationResourceMock) SaveTemperatureChange(t *internal.TemperatureChange) error {
	m.SaveTemperatureChangeInvoked = true
	return m.SaveTemperatureChangeFn(t)
}
