package mock

import "github.com/benjaminbartels/zymurgauge"

// ChamberService is the mock implementation of zymurgauge.ChamberService
type ChamberService struct {
	GetFn              func(mac string) (*zymurgauge.Chamber, error)
	SaveFn             func(controller *zymurgauge.Chamber) error
	SubscribeFn        func(mac string, ch chan zymurgauge.Chamber) error
	UnsubscribeFn      func(mac string)
	GetInvoked         bool
	SaveInvoked        bool
	SubscribeInvoked   bool
	UnsubscribeInvoked bool
}

// Get returns a controller by MAC address
func (s *ChamberService) Get(mac string) (*zymurgauge.Chamber, error) {
	s.GetInvoked = true
	return s.GetFn(mac)
}

// Save creates or updates the stored controller with the given Controller
func (s *ChamberService) Save(controller *zymurgauge.Chamber) error {
	s.SaveInvoked = true
	return s.SaveFn(controller)
}

// Subscribe registers the caller to receives updates to the given controller on the given channel
func (s *ChamberService) Subscribe(mac string, ch chan zymurgauge.Chamber) error {
	s.SubscribeInvoked = true
	return s.SubscribeFn(mac, ch)
}

// Unsubscribe unregisters the caller to receives updates to the given controller
func (s *ChamberService) Unsubscribe(mac string) {
	s.UnsubscribeInvoked = true
	s.UnsubscribeFn(mac)
}

// FermentationService is the mock implementation of zymurgauge.FermentationService
type FermentationService struct {
	GetFn           func(id uint64) (*zymurgauge.Fermentation, error)
	SaveFn          func(fermentation *zymurgauge.Fermentation) error
	LogEventFn      func(fermentationID uint64, event zymurgauge.FermentationEvent) error
	GetInvoked      bool
	SaveInvoked     bool
	LogEventInvoked bool
}

// Get returns a fermentation by ids
func (s *FermentationService) Get(id uint64) (*zymurgauge.Fermentation, error) {
	s.GetInvoked = true
	return s.GetFn(id)
}

// Save creates or updates the stored fermentation with the given Fermentation
func (s *FermentationService) Save(fermentation *zymurgauge.Fermentation) error {
	s.SaveInvoked = true
	return s.SaveFn(fermentation)
}

// LogEvent logs the given event for the given Fermentation by it FermentationID
func (s *FermentationService) LogEvent(fermentationID uint64, event zymurgauge.FermentationEvent) error {
	s.LogEventInvoked = true
	return s.LogEventFn(fermentationID, event)
}

// BeerService is the mock implementation of zymurgauge.BeerService
type BeerService struct {
	GetFn       func(id uint64) (*zymurgauge.Beer, error)
	SaveFn      func(beer *zymurgauge.Beer) error
	GetInvoked  bool
	SaveInvoked bool
}

// Get returns a beer by ids
func (s *BeerService) Get(id uint64) (*zymurgauge.Beer, error) {
	s.GetInvoked = true
	return s.GetFn(id)
}

// Save creates or updates the stored beer with the given Beer
func (s *BeerService) Save(beer *zymurgauge.Beer) error {
	s.SaveInvoked = true
	return s.SaveFn(beer)
}
